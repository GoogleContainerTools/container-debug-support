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

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestValidateDebugMode(t *testing.T) {
	tests := []struct {
		mode      string
		shouldErr bool
	}{
		{"debugpy", false},
		{"ptvsd", false},
		{"pydevd", false},
		{"pydevd-pycharm", false},
		{"", true},
		{"pydev", true},         // the 'd' is important
		{"pydev-pycharm", true}, // the 'd' is important
	}
	for _, test := range tests {
		t.Run(test.mode, func(t *testing.T) {
			result := validateDebugMode(test.mode)
			if test.shouldErr && result == nil {
				t.Error("should have errored")
			} else if !test.shouldErr && result != nil {
				t.Error("should not have errored")
			}
		})
	}
}

func TestIsEnabled(t *testing.T) {
	tests := []struct {
		env      env
		expected bool
	}{
		{
			env:      nil,
			expected: true,
		},
		{
			env:      env{"WRAPPER_ENABLED": "1"},
			expected: true,
		},
		{
			env:      env{"WRAPPER_ENABLED": "true"},
			expected: true,
		},
		{
			env:      env{"WRAPPER_ENABLED": "yes"},
			expected: true,
		},
		{
			env:      env{"WRAPPER_ENABLED": ""},
			expected: true,
		},
		{
			env:      env{"WRAPPER_ENABLED": "0"},
			expected: false,
		},
		{
			env:      env{"WRAPPER_ENABLED": "no"},
			expected: false,
		},
		{
			env:      env{"WRAPPER_ENABLED": "false"},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("env: %v", test.env), func(t *testing.T) {
			result := isEnabled(test.env)
			if test.expected != result {
				t.Errorf("expected %v but got %v", test.expected, result)
			}
		})
	}
}

func TestAlreadyConfigured(t *testing.T) {
	tests := []struct {
		description string
		pc          pythonContext
		expected    bool
	}{
		{"non-python", pythonContext{args: []string{"/app"}}, false},
		{"python with no debug", pythonContext{args: []string{"python", "app.py"}}, false},
		{"misconfigured python module", pythonContext{args: []string{"python", "-m"}}, false},
		{"python with app module", pythonContext{args: []string{"python", "-mapp"}}, false},
		{"python with app module 2", pythonContext{args: []string{"python", "-m", "app"}}, false},
		{"configured for pydevd", pythonContext{args: []string{"pydevd", "--server", "app"}}, true},
		{"configured for pydevd", pythonContext{args: []string{"/dbg/pydevd/bin/pydevd", "--server", "app"}}, true},
		{"configured for pydevd", pythonContext{args: []string{"python", "-mpydevd", "--server", "app"}}, true},
		{"configured for pydevd", pythonContext{args: []string{"python3.8", "-m", "pydevd", "--server", "app"}}, true},
		{"python with debugpy module", pythonContext{args: []string{"python", "-mdebugpy"}}, true},
		{"versioned python with debugpy module", pythonContext{args: []string{"/usr/bin/python3.9", "-m", "debugpy"}}, true},
		{"python with ptvsd module", pythonContext{args: []string{"python", "-mptvsd"}}, true},
		{"versioned python with ptvsd module", pythonContext{args: []string{"/usr/bin/python3.9", "-m", "ptvsd"}}, true},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := test.pc.alreadyConfigured()
			if test.expected != result {
				t.Errorf("expected %v but got %v", test.expected, result)
			}
		})
	}
}

func TestUnwrapLauncher(t *testing.T) {
	tests := []struct {
		description string
		filename    string
		contents    []byte
		shouldErr   bool
		expected    []string
	}{
		{
			description: "non-existent file",
			filename:    "d03$-n0t-3x1$t",
			shouldErr:   true,
		},
		{
			description: "empty file",
			contents:    nil,
			expected:    nil,
		},
		{
			description: "non-shebang",
			contents:    []byte{0, 1, 2, 3, 4, 5, 6, 7},
			expected:    nil,
		},
		{
			description: "python script",
			contents:    []byte("#!/bin/python\nprint \"Hello World\""),
			expected:    []string{"/bin/python"},
		},
		{
			description: "script with args",
			contents:    []byte("#!/bin/sh -x\necho \"Hello World\""),
			expected:    []string{"/bin/sh", "-x"},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			n := test.filename
			if n == "" {
				f, err := ioutil.TempFile(t.TempDir(), "script*")
				if err != nil {
					t.Fatal(err)
				}
				if _, err := f.Write(test.contents); err != nil {
					t.Fatal("error creating temp file", err)
				}
				f.Close()
				n = f.Name()
			}
			// for script files, the shebang should be extracted and parsed, and then
			// prepended to the current command-line.
			pc := pythonContext{args: []string{n, "arg1", "arg2"}}
			expected := []string{n, "arg1", "arg2"}
			if test.expected != nil {
				expected = append(test.expected, expected...)
			}
			err := pc.unwrapLauncher(nil)
			if test.shouldErr && err == nil {
				t.Error("expected an error")
			}
			if !test.shouldErr && err != nil {
				t.Error("should not error:", err)
			} else if diff := cmp.Diff(pc.args, expected); diff != "" {
				t.Errorf("%T differ (-got, +want): %s", test.expected, diff)
			}
		})
	}
}

func TestDeterminePythonMajorMinor(t *testing.T) {
	tests := []struct {
		description string
		env         env
		commands    commands
		shouldErr   bool
		major       int
		minor       int
	}{
		{description: "2.7", commands: RunCmdOut([]string{"python", "-V"}, "Python 2.7.8"), major: 2, minor: 7},
		{description: "2.7 and newline", commands: RunCmdOut([]string{"python", "-V"}, "Python 2.7.2\n"), major: 2, minor: 7},
		{description: "3.9 and newline", commands: RunCmdOut([]string{"python", "-V"}, "Python 3.9.14\n"), major: 3, minor: 9},
		{description: "4.13 from env", env: env{"WRAPPER_PYTHON_VERSION": "4.13.8888"}, major: 4, minor: 13},
		{description: "error", commands: RunCmdOutFail([]string{"python", "-V"}, "", 1), shouldErr: true, major: -1, minor: -1},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			test.commands.Setup(t)
			major, minor, err := determinePythonMajorMinor(context.TODO(), "python", test.env)
			if test.shouldErr && err == nil {
				t.Error("expected an error")
			} else if !test.shouldErr && err != nil {
				t.Error("unexpected error:", err)
			}
			if test.major != major {
				t.Errorf("expected major %d but got %d", test.major, major)
			}
			if test.minor != minor {
				t.Errorf("expected minor %d but got %d", test.minor, minor)
			}
		})
	}
}

func TestPrepare(t *testing.T) {
	dbgRoot = t.TempDir()

	tests := []struct {
		description string
		pc          pythonContext
		commands    commands
		shouldFail  bool
		expected    pythonContext
	}{
		{
			description: "debugpy",
			pc:          pythonContext{debugMode: "debugpy", port: 2345, wait: false, args: []string{"python", "app.py"}, env: nil},
			commands: RunCmdOut([]string{"python", "-V"}, "Python 3.7.4\n").
				AndRunCmd([]string{"python", "-m", "debugpy", "--listen", "2345", "app.py"}),
			expected: pythonContext{debugMode: "debugpy", port: 2345, wait: false, major: 3, minor: 7, args: []string{"python", "-m", "debugpy", "--listen", "2345", "app.py"}, env: env{"PYTHONPATH": dbgRoot + "/python/lib/python3.7/site-packages"}},
		},
		{
			description: "debugpy with wait",
			pc:          pythonContext{debugMode: "debugpy", port: 2345, wait: true, args: []string{"python", "app.py"}, env: nil},
			commands: RunCmdOut([]string{"python", "-V"}, "Python 3.7.4\n").
				AndRunCmd([]string{"python", "-m", "debugpy", "--listen", "2345", "--wait-for-client", "app.py"}),
			expected: pythonContext{debugMode: "debugpy", port: 2345, wait: true, major: 3, minor: 7, args: []string{"python", "-m", "debugpy", "--listen", "2345", "--wait-for-client", "app.py"}, env: env{"PYTHONPATH": dbgRoot + "/python/lib/python3.7/site-packages"}},
		},
		{
			description: "ptvsd",
			pc:          pythonContext{debugMode: "ptvsd", port: 2345, wait: false, args: []string{"python", "app.py"}, env: nil},
			commands: RunCmdOut([]string{"python", "-V"}, "Python 3.7.4\n").
				AndRunCmd([]string{"python", "-m", "ptvsd", "--host", "localhost", "--port", "2345", "app.py"}),
			expected: pythonContext{debugMode: "ptvsd", port: 2345, wait: false, major: 3, minor: 7, args: []string{"python", "-m", "ptvsd", "--host", "localhost", "--port", "2345", "app.py"}, env: env{"PYTHONPATH": dbgRoot + "/python/lib/python3.7/site-packages"}},
		},
		{
			description: "ptvsd with wait",
			pc:          pythonContext{debugMode: "ptvsd", port: 2345, wait: true, args: []string{"python", "app.py"}, env: nil},
			commands: RunCmdOut([]string{"python", "-V"}, "Python 3.7.4\n").
				AndRunCmd([]string{"python", "-m", "ptvsd", "--host", "localhost", "--port", "2345", "--wait", "app.py"}),
			expected: pythonContext{debugMode: "ptvsd", port: 2345, wait: true, major: 3, minor: 7, args: []string{"python", "-m", "ptvsd", "--host", "localhost", "--port", "2345", "--wait", "app.py"}, env: env{"PYTHONPATH": dbgRoot + "/python/lib/python3.7/site-packages"}},
		},
		{
			description: "pydevd",
			pc:          pythonContext{debugMode: "pydevd", port: 2345, wait: false, args: []string{"python", "app.py"}, env: nil},
			commands: RunCmdOut([]string{"python", "-V"}, "Python 3.7.4\n").
				AndRunCmd([]string{"python", "-m", "pydevd", "--server", "--port", "2345", "--continue", "--file", "app.py"}),
			expected: pythonContext{debugMode: "pydevd", port: 2345, wait: false, major: 3, minor: 7, args: []string{"python", "-m", "pydevd", "--server", "--port", "2345", "--continue", "--file", "app.py"}, env: env{"PYTHONPATH": dbgRoot + "/python/pydevd/python3.7/lib/python3.7/site-packages"}},
		},
		{
			description: "pydevd with wait",
			pc:          pythonContext{debugMode: "pydevd", port: 2345, wait: true, args: []string{"python", "app.py"}, env: nil},
			commands: RunCmdOut([]string{"python", "-V"}, "Python 3.7.4\n").
				AndRunCmd([]string{"python", "-m", "pydevd", "--server", "--port", "2345", "--file", "app.py"}),
			expected: pythonContext{debugMode: "pydevd", port: 2345, wait: true, major: 3, minor: 7, args: []string{"python", "-m", "pydevd", "--server", "--port", "2345", "--file", "app.py"}, env: env{"PYTHONPATH": dbgRoot + "/python/pydevd/python3.7/lib/python3.7/site-packages"}},
		},
		{
			description: "WRAPPER_ENABLED=false",
			pc:          pythonContext{debugMode: "pydevd", port: 2345, wait: true, args: []string{"python", "app.py"}, env: map[string]string{"WRAPPER_ENABLED": "false"}},
			shouldFail:  true,
			expected:    pythonContext{debugMode: "pydevd", port: 2345, wait: true, args: []string{"python", "app.py"}, env: map[string]string{"WRAPPER_ENABLED": "false"}},
		},
		{
			description: "already configured with debugpy",
			pc:          pythonContext{debugMode: "debugpy", port: 2345, wait: false, args: []string{"python", "-m", "debugpy", "--listen", "2345", "app.py"}},
			shouldFail:  true,
			expected:    pythonContext{debugMode: "debugpy", port: 2345, wait: false, args: []string{"python", "-m", "debugpy", "--listen", "2345", "app.py"}},
		},
		{
			description: "already configured with ptvsd",
			pc:          pythonContext{debugMode: "ptvsd", port: 2345, wait: true, args: []string{"python", "-m", "ptvsd", "--host", "localhost", "--port", "2345", "--wait", "app.py"}},
			shouldFail:  true,
			expected:    pythonContext{debugMode: "ptvsd", port: 2345, wait: true, args: []string{"python", "-m", "ptvsd", "--host", "localhost", "--port", "2345", "--wait", "app.py"}},
		},
		{
			description: "already configured with pydevd",
			pc:          pythonContext{debugMode: "pydevd", port: 2345, wait: true, args: []string{"python", "-m", "pydevd", "--server", "--port", "2345", "--file", "app.py"}},
			shouldFail:  true,
			expected:    pythonContext{debugMode: "pydevd", port: 2345, wait: true, args: []string{"python", "-m", "pydevd", "--server", "--port", "2345", "--file", "app.py"}},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			test.commands.Setup(t)
			pc := test.pc
			result := pc.prepare(context.TODO())

			if test.shouldFail && result == true {
				t.Error("prepare() should have failed")
			} else if !test.shouldFail && !result {
				t.Error("prepare() should have succeeded")
			} else if diff := cmp.Diff(test.expected, pc, cmp.AllowUnexported(test.expected)); diff != "" {
				_t.Errorf("%T differ (-got, +want): %s", pc, diff)
			}
		})
	}
}

func TestPathExists(t *testing.T) {
	if pathExists(filepath.Join("this", "should", "not", "exist")) {
		t.Error("pathExists should have failed on non-existent path")
	}
	if !pathExists(t.TempDir()) {
		t.Error("pathExists failed on real path")
	}
}

func TestHandlePydevModule(t *testing.T) {
	tmp := os.TempDir()

	tests := []struct {
		description string
		args        []string
		shouldErr   bool
		module      string
		file        string
		remaining   []string
	}{
		{
			description: "plain file",
			args:        []string{"app.py"},
			file:        "app.py",
		},
		{
			description: "-mmodule",
			args:        []string{"-mmodule"},
			file:        filepath.Join(tmp, "*", "skaffold_pydevd_launch.py"),
		},
		{
			description: "-m module",
			args:        []string{"-m", "module"},
			file:        filepath.Join(tmp, "*", "skaffold_pydevd_launch.py"),
		},
		{
			description: "- should error",
			args:        []string{"-", "module"},
			shouldErr:   true,
		},
		{
			description: "-x should error",
			args:        []string{"-x", "module"},
			shouldErr:   true,
		},
		{
			description: "lone -m should error",
			args:        []string{"-m"},
			shouldErr:   true,
		},
		{
			description: "no args should error",
			shouldErr:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			file, args, err := handlePydevModule(test.args)
			if test.shouldErr {
				if err == nil {
					t.Error("Expected an error")
				}
			} else {
				if !fileMatch(t, test.file, file) {
					t.Errorf("Wanted %q but got %q", test.file, file)
				}
				if diff := cmp.Diff(args, test.remaining, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("remaining args %T differ (-got, +want): %s", test.remaining, diff)
				}
			}
		})
	}
}

func fileMatch(t *testing.T, glob, file string) bool {
	if file == glob {
		return true
	}
	matches, err := filepath.Glob(glob)
	if err != nil {
		t.Errorf("Failed to expand globe %q: %v", glob, err)
		return false
	}
	for _, m := range matches {
		if file == m {
			return true
		}
	}
	return false
}
