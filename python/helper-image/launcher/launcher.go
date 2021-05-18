/*
Copyright 2021 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// A `skaffold debug` launcher for Python.
//
// Python introduces some quirks.  There are now three
// methods for hooking up a debugging backend:
//
// - pydevd: the stock Python debugging backend
// - pydevd-pycharm: PyDev with modifications for IntelliJ/PyCharm
// - ptvsd: wraps pydevd with the debug-adapter protocol (obsolete)
// - debugpy: new and improved ptvsd
//
// pydevd has pyx libraries which are specific to particular versions of Python.
//
// Further complicating matters is that a number of Python packages
// use launcher scripts (e.g., gunicorn), and so we can't simply run
// `python -m ptvsd -- gunicorn` as ptvsd/debugpy/etc don't look for
// the script file in the PATH.
//
// Another wrinkle is that we cannot just provide a `python` wrapper
// executable that will hand off to the real `python` as `pip install
// hard-codes the python binary location in launcher scripts.  And it's
// not that unusual to have a `python`, `python3`, and `python2`
// scripts that invoke different python installations.
//
// This launcher is expected to be invoked as follows:
//
//    launcher --mode <pydevd|pydevd-pycharm|debugpy|ptvsd> \
//        --port p [--wait] -- original-command-line ...
//
// This launcher will determine the python executable based on the `original-command-line`.
// The launcher will configure the PYTHONPATH to point to the appropriate installation
// of pydevd/debugpy/ptvsd for the corresponding python binary.
//
// debugpy and ptvsd are pretty straightforward translations of the
// launcher command-line `python -m debugpy`.
//
// pydevd is more involved as pydevd does not support loading modules
// from the command-line (e.g., `python -m flask`).  This launcher
// instead creates a small module-loader script using runpy.
// So `launcher --mode pydevd --port 5678 -- python -m flask app.py`
// will create a temp file named `skaffold_pydevd_launch.py`:
// ```
// import sys
// import runpy
// runpy.run_module('flask', run_name="__main__",alter_sys=True)
// ```
// and will then invoke:
// ```
// python -m pydevd --server --port 5678 --DEBUG --continue \
//   --file /tmp/pydevd716531212/skaffold_pydevd_launch.py
// ```
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	shell "github.com/kballard/go-shellquote"
	"github.com/sirupsen/logrus"
)

var (
	// dbgRoot is the location where the skaffold-debug helpers should be installed.
	// The python helpers should be in dbgRoot + "/python"
	dbgRoot = "/dbg"
)

const (
	ModeDebugpy       string = "debugpy"
	ModePtvsd         string = "ptvsd"
	ModePydevd        string = "pydevd"
	ModePydevdPycharm string = "pydevd-pycharm"
)

// pythonContext represents the launch context.
type pythonContext struct {
	debugMode string
	port      uint
	wait      bool

	args []string
	env  env
}

func main() {
	ctx := context.Background()
	env := EnvFromPairs(os.Environ())
	logrus.SetLevel(logrusLevel(env))
	logrus.Debug("launcher args:", os.Args[1:])

	pc := pythonContext{env: env}
	flag.StringVar(&dbgRoot, "helpers", "/dbg", "base location for skaffold-debug helpers")
	flag.StringVar(&pc.debugMode, "mode", "", "debugger mode: debugpy, ptvsd, pydevd, pydevd-pycharm")
	flag.UintVar(&pc.port, "port", 9999, "port to listen for remote debug connections")
	flag.BoolVar(&pc.wait, "wait", false, "wait for debugger connection on start")

	flag.Parse()
	if err := validateDebugMode(pc.debugMode); err != nil {
		logrus.Fatal(err)
	}

	if len(flag.Args()) == 0 {
		logrus.Fatal("expected python command-line args")
	}
	pc.args = flag.Args()
	logrus.Tracef("command-line: %v", pc.args)

	if err := pc.launch(ctx); err != nil {
		logrus.Fatalf("error launching python debugging: %v", err)
	}
}

// validateDebugMode ensures the provided mode is a supported mode.
func validateDebugMode(mode string) error {
	switch mode {
	case ModeDebugpy, ModePtvsd, ModePydevd, ModePydevdPycharm:
		return nil
	default:
		return fmt.Errorf("unknown debugger mode %q; expecting one of %v", mode, []string{ModeDebugpy, ModePtvsd, ModePydevd, ModePydevdPycharm})
	}
}

func (pc *pythonContext) launch(ctx context.Context) error {
	if !isEnabled(pc.env) {
		logrus.Infof("wrapper disabled, launching %v", pc.args)
		cmd := newConsoleCommand(ctx, pc.args, pc.env)
		return cmd.Run()
	}
	if pc.alreadyConfigured() {
		logrus.Infof("already configured for debugging")
		cmd := newConsoleCommand(ctx, pc.args, pc.env)
		return cmd.Run()
	}

	// rewrite the command-line by expanding script shebangs to run python and launch the app
	if err := pc.unwrapLauncher(ctx); err != nil {
		return err
	}

	// set PYTHONPATH to point to the appropriate library for the given python version.
	if err := pc.updateEnv(ctx); err != nil {
		return err
	}
	// so pc.args[0] should be the python interpreter

	if err := pc.updateCommandLine(ctx); err != nil {
		return err
	}

	cmd := newConsoleCommand(ctx, pc.args, pc.env)
	return cmd.Run()
}

// alreadyConfigured tries to determine if the python command-line is already configured
// for debugging.  Only handles simple command-lines; users should set `WRAPPER_ENABLED=false`
// for more complicated situations.
func (pc *pythonContext) alreadyConfigured() bool {
	// TODO: consider handling `#!/usr/bin/env python` too, though `pip install` seems
	// to hard-code the python location instead.
	if filepath.Base(pc.args[0]) == "pydevd" {
		logrus.Debug("already configured to use pydevd")
		return true
	}
	if strings.HasPrefix(filepath.Base(pc.args[0]), "python") && len(pc.args) > 1 {
		if (pc.args[1] == "-m" && len(pc.args) > 2 && pc.args[2] == "debugpy") || pc.args[1] == "-mdebugpy" {
			logrus.Debug("already configured to use debugpy")
			return true
		}
		if (pc.args[1] == "-m" && len(pc.args) > 2 && pc.args[2] == "ptvsd") || pc.args[1] == "-mptvsd" {
			logrus.Debug("already configured to use ptvsd")
			return true
		}
		if (pc.args[1] == "-m" && len(pc.args) > 2 && pc.args[2] == "pydevd") || pc.args[1] == "-mpydevd" {
			logrus.Debug("already configured to use pydevd")
			return true
		}
	}
	return false
}

// unwrapLauncher attempts to expand the command-line in the given script,
// providing that it does not look like a `python` launcher.
// TODO: Windows .cmd and .bat files?
func (pc *pythonContext) unwrapLauncher(_ context.Context) error {
	p := pc.args[0]

	_, err := os.Stat(p)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("could not access launcher %q: %w", p, err)
		}
		// try looking through PATH
		l, err := exec.LookPath(p)
		if err != nil {
			return fmt.Errorf("could not find launcher %q: %w", p, err)
		}
		p = l
	}
	if strings.HasPrefix(filepath.Base(p), "python") {
		logrus.Debugf("no further unwrapping required: launcher appears to be python: %q", p)
		return nil
	}
	f, err := os.Open(p)
	if err != nil {
		return fmt.Errorf("could not open launcher %q: %w", p, err)
	}
	defer f.Close()

	shebang := make([]byte, 1024)
	if n, err := f.Read(shebang); err == io.EOF || n < 2 {
		logrus.Debugf("%q has no shebang", p)
		return nil
	} else if err != nil {
		return fmt.Errorf("error reading file header from %q: %w", p, err)
	} else if string(shebang[0:2]) != "#!" {
		logrus.Debugf("%q appears to be a binary", p)
		return nil
	}
	cl := strings.SplitN(string(shebang[2:]), "\n", 2)[0]
	logrus.Tracef("%q has shebang %q", p, cl)
	s, err := shell.Split(cl)
	if err != nil {
		logrus.Warnf("%q shebang %q seems odd: %v", p, cl, err)
		s = []string{cl}
	}
	pc.args[0] = p // ensure script is full path if resolved in PATH
	pc.args = append(s, pc.args...)
	logrus.Debugf("expanded command-line: %q -> %v", p, pc.args)
	return nil
}

func (pc *pythonContext) updateEnv(ctx context.Context) error {
	// Perhaps we should check PYTHONPATH or ~/.local to see if the user has already
	// installed one of our supported debug libraries
	if pc.env["WRAPPER_SKIP_ENV"] != "" {
		logrus.Debug("Skipping environment configuration by request")
		return nil
	}

	_, err := os.Stat(dbgRoot)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Warnf("skaffold-debug helpers not found at %q", dbgRoot)
			return nil
		}
		return fmt.Errorf("skaffold-debug helpers are inaccessible at %q: %w", dbgRoot, err)
	}

	major, minor, err := determinePythonMajorMinor(ctx, pc.args[0], pc.env)
	if err != nil {
		// We could skip setting PYTHONPATH in the hopes that the appropaite
		// debugger library was installed explicitly by the user?
		return fmt.Errorf("unable to determine python version from %q: %w", pc.args[0], err)
	}

	if pc.env == nil {
		pc.env = env{}
	}
	// The skaffold-debug-python helper image places pydevd and debugpy in /dbg/python/lib/pythonM.N,
	// but separates pydevd and pydevd-pycharm in separate directories to avoid possible leakage.
	var libraryPath string
	switch pc.debugMode {
	case ModePtvsd, ModeDebugpy:
		libraryPath = fmt.Sprintf(dbgRoot+"/python/lib/python%d.%d/site-packages", major, minor)

	case ModePydevd:
		libraryPath = fmt.Sprintf(dbgRoot+"/python/pydevd/python%d.%d/lib/python%d.%d/site-packages", major, minor, major, minor)

	case ModePydevdPycharm:
		libraryPath = fmt.Sprintf(dbgRoot+"/python/pydevd-pycharm/python%d.%d/lib/python%d.%d/site-packages", major, minor, major, minor)
	}
	if libraryPath != "" {
		if !pathExists(libraryPath) {
			// Warn as the user may have installed debugpy themselves
			logrus.Warnf("Debugging support for Python %d.%d not found: may require manually installing %q", major, minor, pc.debugMode)
		}
		// Append to ensure user-configured values are found first.
		pc.env.AppendFilepath("PYTHONPATH", libraryPath)
	}
	return nil
}

func (pc *pythonContext) updateCommandLine(ctx context.Context) error {
	var cmdline []string
	switch pc.debugMode {
	case ModePtvsd:
		cmdline = append(cmdline, pc.args[0])
		cmdline = append(cmdline, "-m", "ptvsd", "--host", "localhost", "--port", strconv.Itoa(int(pc.port)))
		if pc.wait {
			cmdline = append(cmdline, "--wait")
		}
		cmdline = append(cmdline, pc.args[1:]...)
		pc.args = cmdline

	case ModeDebugpy:
		cmdline = append(cmdline, pc.args[0])
		cmdline = append(cmdline, "-m", "debugpy", "--listen", strconv.Itoa(int(pc.port)))
		if pc.wait {
			cmdline = append(cmdline, "--wait-for-client")
		}
		cmdline = append(cmdline, pc.args[1:]...)
		pc.args = cmdline

	case ModePydevd, ModePydevdPycharm:
		// Appropriate location to resolve pydevd is set in updateEnv
		cmdline = append(cmdline, pc.args[0])
		cmdline = append(cmdline, "-m", "pydevd", "--server", "--port", strconv.Itoa(int(pc.port)))
		if pc.env["WRAPPER_VERBOSE"] != "" {
			cmdline = append(cmdline, "--DEBUG")
		}
		if !pc.wait {
			cmdline = append(cmdline, "--continue")
		}

		// --file is expected as last pydev argument, but it must be a file, and so launching with
		// a module requires some special handling.
		cmdline = append(cmdline, "--file")
		file, args, err := handlePydevModule(pc.args[1:])
		if err != nil {
			return err
		}
		cmdline = append(cmdline, file)
		cmdline = append(cmdline, args...)
		if pc.wait {
			logrus.Warn("pydevd does not support wait-for-client")
		}
		pc.args = cmdline
	}
	return nil
}

func determinePythonMajorMinor(ctx context.Context, launcherBin string, env env) (major, minor int, err error) {
	var versionString string
	if env["WRAPPER_PYTHON_VERSION"] != "" {
		versionString = env["WRAPPER_PYTHON_VERSION"]
	} else {
		cmd := newCommand(ctx, []string{launcherBin, "-V"}, env)
		out, err := cmd.CombinedOutput()
		if err != nil {
			logrus.Warnf("'%s -V' errored: %v", launcherBin, err)
			return -1, -1, err
		}
		versionString = string(out)
		logrus.Debugf("'%s -V' = %q", launcherBin, versionString)
		if !strings.HasPrefix(versionString, "Python ") {
			return -1, -1, fmt.Errorf("launcher is not a python interpreter: %q", launcherBin)
		}
		versionString = versionString[len("Python "):]
	}

	v := strings.Split(strings.TrimSpace(versionString), ".")
	major, err = strconv.Atoi(v[0])
	if err == nil {
		minor, err = strconv.Atoi(v[1])
	}
	return
}

// handlePydevModule applies special pydevd handling for a python module.  When a module is
// found, we write out a python script that uses runpy to invoke the module.
func handlePydevModule(args []string) (string, []string, error) {
	switch {
	case len(args) == 0:
		return "", nil, fmt.Errorf("no python command-line specified") // shouldn't happen
	case !strings.HasPrefix(args[0], "-"):
		// this is a file
		return args[0], args[1:], nil
	case !strings.HasPrefix(args[0], "-m"):
		// this is some other command-line flag
		return "", nil, fmt.Errorf("expected python module: %q", args)
	}
	module := args[0][2:]
	remaining := args[1:]
	if module == "" {
		if len(args) == 1 {
			return "", nil, fmt.Errorf("missing python module: %q", args)
		}
		module = args[1]
		remaining = args[2:]
	}

	snippet := strings.ReplaceAll(`import sys
import runpy
runpy.run_module('{module}', run_name="__main__",alter_sys=True)
`, `{module}`, module)

	// write out the temp location as other locations may not be writable
	d, err := ioutil.TempDir("", "pydevd*")
	if err != nil {
		return "", nil, err
	}
	// use a skaffold-specific file name to ensure no possibility of it matching a user import
	f := filepath.Join(d, "skaffold_pydevd_launch.py")
	if err := ioutil.WriteFile(f, []byte(snippet), 0755); err != nil {
		return "", nil, err
	}
	return f, remaining, nil
}

func isEnabled(env env) bool {
	v, found := env["WRAPPER_ENABLED"]
	return !found || (v != "0" && v != "false" && v != "no")
}

func logrusLevel(env env) logrus.Level {
	v := env["WRAPPER_VERBOSE"]
	if v != "" {
		if l, err := logrus.ParseLevel(v); err == nil {
			return l
		}
		logrus.Warnln("Unknown logging level: WRAPPER_VERBOSE=", v)
	}
	return logrus.WarnLevel
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil || !os.IsNotExist(err) {
		return true
	}
	return false
}
