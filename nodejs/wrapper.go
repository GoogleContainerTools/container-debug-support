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
	logrus.Debugf("Launched: %v", os.Args)

	env := envToMap(os.Environ())
	// suppress npm warnings when node on PATH isn't the node used for npm
	env["npm_config_scripts_prepend_node_path"] = "false"
	nc := nodeContext{program: os.Args[0], args: os.Args[1:], env: env}
	if err := run(&nc, os.Stdin, os.Stdout, os.Stderr); err != nil {
		logrus.Fatal(err)
	}
}

func logrusLevel() logrus.Level {
	switch os.Getenv("WRAPPER_VERBOSE") {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel

	default:
		return logrus.WarnLevel
	}
}

func run(nc *nodeContext, stdin io.Reader, stdout, stderr io.Writer) error {
	if !nc.unwrap() {
		return fmt.Errorf("unwrap could not find actual executable")
	}
	logrus.Debugf("unwrapped: %s\n", nc.program)

	script := findScript(nc.args)
	logrus.Debugf("script: %s\n", script)

	nodeDebugOption, hasNodeDebug := nc.env["NODE_DEBUG"]
	if hasNodeDebug {
		logrus.Debugf("NODE_DEBUG: %s", nodeDebugOption)
	}

	// if we're about to execute the application script, install the NODE_DEBUG
	// arguments if found and go
	if isApplicationScript(script) {
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
			logrus.Debugf("Setting NODE_DEBUG: %s", inspectArg)
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
func (nc *nodeContext) unwrap() bool {
	if nc == nil {
		return false
	}
	path := nc.env["PATH"]
	origInfo, err := os.Stat(nc.program)
	origFound := err == nil
	if err != nil && !os.IsNotExist(err) {
		logrus.Errorf("unable to stat %q: %v", nc.program, err)
		return false
	}
	base := filepath.Base(nc.program)
	for _, dir := range strings.Split(path, string(os.PathListSeparator)) {
		p := filepath.Join(dir, base)
		if pInfo, err := os.Stat(p); err == nil && (!origFound || !os.SameFile(origInfo, pInfo)) {
			nc.program = p
			return true
		}
	}
	logrus.Errorf("unable to unwrap %q: not in PATH", base)
	return false
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
	if len(nc.args) == 0 {
		nc.args = []string{nodeArg}
		return
	}
	for i, arg := range nc.args {
		if len(arg) > 0 && arg[0] != '-' {
			nc.args = append(nc.args, "")
			copy(nc.args[i+1:], nc.args[i:])
			nc.args[i] = nodeArg
			logrus.Debugf("added node arg: %q", nc.args)
			return
		}
	}
	nc.args = append(nc.args, nodeArg)
}

// exec runs the command, and returns an error should one occur.
func (nc *nodeContext) exec(in io.Reader, out, err io.Writer) error {
	logrus.Debugf("exec: %s %v (env: %v)", nc.program, nc.args, nc.env)
	cmd := exec.Command(nc.program, nc.args...)
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
	if len(entries) == 0 {
		return nil
	}
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
func stripInspectArg(args []string) (newArgs []string, inspectArg string) {
	// inspect directives are always a single argument: `node --inspect 9226` causes node to load 9226 as a file
	newArgs = nil
	inspectArg = "" // default case: no inspect arg found

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
	return
}
