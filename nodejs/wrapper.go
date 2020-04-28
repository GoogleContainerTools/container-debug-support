/*
Copyright 2020 The Skaffold Authors

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

// A wrapper for node executables to support debugging of application scripts.
// Many NodeJS applications use NodeJS-based launch tools (e.g., npm,
// nodemon), and often use several in combination.  This makes it very
// difficult to start debugging the application as `--inspect`s are usually
// intercepted by one of the launch tools.  When executing a `node_modules`
// script, this wrapper strips out and propagates `--inspect`-like arguments
// via `NODE_DEBUG`.  When executing an app script, this wrapper then inlines
// the `NODE_DEBUG` when found.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	shell "github.com/kballard/go-shellquote"
	"github.com/sirupsen/logrus"
)

// nodeContext allows manipulating the launch context for node.
type nodeContext struct {
	program string
	args    []string
	env     map[string]string
}

func main() {
	logrus.SetLevel(logrusLevel())
	logrus.Debugln("Launched: ", os.Args)

	env := envToMap(os.Environ())
	// suppress npm warnings when node on PATH isn't the node used for npm
	env["npm_config_scripts_prepend_node_path"] = "false"
	nc := nodeContext{program: os.Args[0], args: os.Args[1:], env: env}
	if err := run(&nc, os.Stdin, os.Stdout, os.Stderr); err != nil {
		logrus.Fatal(err)
	}
}

func logrusLevel() logrus.Level {
	v := os.Getenv("WRAPPER_VERBOSE")
	if v != "" {
		if l, err := logrus.ParseLevel(v); err == nil {
			return l
		}
		logrus.Warnln("Unknown logging level: WRAPPER_VERBOSE=", v)
	}
	return logrus.WarnLevel
}

func run(nc *nodeContext, stdin io.Reader, stdout, stderr io.Writer) error {
	if err := nc.unwrap(); err != nil {
		return fmt.Errorf("could not unwrap: %w", err)
	}
	logrus.Debugln("unwrapped: ", nc.program)

	// Use an absolute path in case we're being run within a node_modules directory
	// If there's an error, then hand off immediately to the real node.
	script := findScript(nc.args)
	if abs, err := filepath.Abs(script); err == nil {
		script = abs
	} else {
		logrus.Warn("could not access script: ", err)
		return nc.exec(stdin, stdout, stderr)
	}		
	logrus.Debugln("script: ", script)

	nodeDebugOption, hasNodeDebug := nc.env["NODE_DEBUG"]
	if hasNodeDebug {
		logrus.Debugln("found NODE_DEBUG=", nodeDebugOption)
	}

	// if we're about to execute the application script, install the NODE_DEBUG
	// arguments if found and go
	if isApplicationScript(script) || script == "" {
		if hasNodeDebug {
			nc.stripInspectArgs() // top-level debug options win
			nc.addNodeArg(nodeDebugOption)
			delete(nc.env, "NODE_DEBUG")
		}
		return nc.exec(stdin, stdout, stderr)
	}

	// We're executing a node module: strip any --inspect args and propagate
	inspectArg := nc.stripInspectArgs()
	if inspectArg != "" {
		logrus.Debugf("Stripped %q as not an app script", inspectArg)
		if !hasNodeDebug {
			logrus.Debugln("Setting NODE_DEBUG=", inspectArg)
			nc.env["NODE_DEBUG"] = inspectArg
		}
	}

	// nodemon needs special handling as `nodemon --inspect` will use spawn to invoke a
	// child node, which picks up this wrapped node.  Otherwise nodemon uses fork to launch
	// the actual application script file directly, which circumvents the use of this node wrapper.
	nc.handleNodemon()

	return nc.exec(stdin, stdout, stderr)
}

// unwrap looks for the real node executable (not this wrapper).
func (nc *nodeContext) unwrap() error {
	if nc == nil {
		return fmt.Errorf("nil context")
	}

	// Here we try to find the original program.  When a program is
	// resolved from the PATH, most shells will set argv[0] to the
	// command and so it won't appear to exist and so the first file
	// resolved in the PATH should be this program.
	origInfo, err := os.Stat(nc.program)
	origFound := err == nil
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to stat %q: %v", nc.program, err)
	}

	path := nc.env["PATH"]
	base := filepath.Base(nc.program)
	for _, dir := range strings.Split(path, string(os.PathListSeparator)) {
		p := filepath.Join(dir, base)
		if pInfo, err := os.Stat(p); err == nil {
			if !origFound {
				// the original nc.program was not resolved, meaning this
				// it had been resolved in the PATH, so treat this first
				// instance as the original file and continue searching
				logrus.Debugln("unwrap: presumed wrapper at ", p)
				origInfo = pInfo
				origFound = true
			} else if !os.SameFile(origInfo, pInfo) {
				logrus.Debugf("unwrap: replacing %s -> %s", nc.program, p)
				nc.program = p
				return nil
			}
		}
	}
	return fmt.Errorf("could not find %q in PATH", base)
}

// stripInspectArgs removes all `--inspect*` args from both the command-line and from
// NODE_OPTIONS.  It returns the last inspect arg or "" if there were no inspect arguments.
func (nc *nodeContext) stripInspectArgs() string {
	foundOption := ""
	if options, found := nc.env["NODE_OPTIONS"]; found {
		if args, err := shell.Split(options); err != nil {
			logrus.Warnf("NODE_OPTIONS cannot be split: %v", err)
		} else {
			args, inspectArg := stripInspectArg(args)
			if inspectArg != "" {
				logrus.Debugf("Found %q in NODE_OPTIONS", inspectArg)
				nc.env["NODE_OPTIONS"] = shell.Join(args...)
				foundOption = inspectArg
			}
		}
	}
	strippedArgs, inspectArg := stripInspectArg(nc.args)
	if inspectArg != "" {
		logrus.Debugf("Found %q in command-line", inspectArg)
		nc.args = strippedArgs
		foundOption = inspectArg
	}
	return foundOption
}

func (nc *nodeContext) handleNodemon() {
	if nodeDebug, found := nc.env["NODE_DEBUG"]; found {
		// look for the nodemon script (if it appears) and insert the --inspect argument
		for i, arg := range nc.args {
			if len(arg) > 0 && arg[0] != '-' && strings.Contains(arg, "/nodemon") {
				nc.args = append(nc.args, "")
				copy(nc.args[i+2:], nc.args[i+1:])
				nc.args[i+1] = nodeDebug
				delete(nc.env, "NODE_DEBUG")
				logrus.Debugf("special handling for nodemon: %q", nc.args)
				return
			}
		}
	}
}

func (nc *nodeContext) addNodeArg(nodeArg string) {
	// find the script location and insert the provided argument
	for i, arg := range nc.args {
		if len(arg) > 0 && arg[0] != '-' {
			nc.args = append(nc.args, "")
			copy(nc.args[i+1:], nc.args[i:])
			nc.args[i] = nodeArg
			logrus.Debugf("added node arg: %q", nc.args)
			return
		}
	}
	// script not found so add at end
	nc.args = append(nc.args, nodeArg)
}

// exec runs the command, and returns an error should one occur.
func (nc *nodeContext) exec(in io.Reader, out, err io.Writer) error {
	logrus.Debugf("exec: %s %v (env: %v)", nc.program, nc.args, nc.env)
	cmd := exec.CommandContext(context.Background(), nc.program, nc.args...)
	cmd.Env = envFromMap(nc.env)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = err
	return cmd.Run()
}

// findScript returns the path to the node script that will be executed.
// Returns an empty string if no script was found.
func findScript(args []string) string {
	// a bit of a hack, but all node options are of the form `--arg=option`
	for _, arg := range args {
		if len(arg) > 0 && arg[0] != '-' {
			return arg
		}
	}
	return ""
}

// isApplicationScript return true if the script appears to be an application
// script, or false if a library (node_modules) script or `npm` (special case).
func isApplicationScript(path string) bool {
	// We could consider checking if the parent's base name is `bin`?
	return !strings.HasPrefix(path, "node_modules/") && !strings.Contains(path, "/node_modules/") &&
		!strings.HasSuffix(path, "/bin/npm")
}

// envToMap turns a set of VAR=VALUE strings to a map.
func envToMap(entries []string) map[string]string {
	m := make(map[string]string)
	for _, entry := range entries {
		kv := strings.SplitN(entry, "=", 2)
		m[kv[0]] = kv[1]
	}
	return m
}

// envToMap turns a map of variable:value pairs into a set of VAR=VALUE strings.
func envFromMap(env map[string]string) []string {
	var m []string
	for k, v := range env {
		m = append(m, k+"="+v)
	}
	return m
}

// stripInspectArg searches and removes all node `--inspect` style arguments, returning the
// altered arguments and the inspect argument.
func stripInspectArg(args []string) ([]string, string) {
	// inspect directives are always a single argument: `node --inspect 9226` causes node to load 9226 as a file
	var newArgs []string
	inspectArg := "" // default case: no inspect arg found

	for i, arg := range args {
		if strings.HasPrefix(arg, "--inspect") {
			// todo: we should coalesce --inspect-port=xxx
			inspectArg = arg
			continue
		}

		// if at end of node options, copy remaining arguments
		// "--" marks end of node options
		if arg == "--" || len(arg) == 0 || arg[0] != '-' {
			newArgs = append(newArgs, args[i:]...)
			break
		}
		newArgs = append(newArgs, arg)
	}
	return newArgs, inspectArg
}
