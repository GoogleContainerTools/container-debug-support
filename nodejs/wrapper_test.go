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

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
)

func TestFindScript(t *testing.T) {
	tests := []struct {
		description string
		args        []string
		expected    string
	}{
		{
			description: "no args",
			args:        []string{},
			expected:    "",
		},
		{
			description: "single script",
			args:        []string{"index.js"},
			expected:    "index.js",
		},
		{
			description: "options but no script",
			args:        []string{"-i", "--help"},
			expected:    "",
		},
		{
			description: "options and script",
			args:        []string{"-i", "--help", "index.js"},
			expected:    "index.js",
		},
		{
			description: "options, script, and arguments",
			args:        []string{"-i", "--help", "index.js", "arg1", "arg2"},
			expected:    "index.js",
		},
		{
			description: "options, script path, and arguments",
			args:        []string{"-i", "--help", "node_modules/index.js", "arg1", "arg2"},
			expected:    "node_modules/index.js",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := findScript(test.args)
			if result != test.expected {
				t.Errorf("expected %s but got %s", test.expected, result)
			}
		})
	}
}

func TestIsApplicationScript(t *testing.T) {
	tests := []struct {
		script   string
		expected bool
	}{
		{"index.js", true},
		{"/usr/local/bin/npm", false},
		{"node_modules/nodemon/nodemon.js", false},
		{"lib/node_modules/nodemon/nodemon.js", false},
	}

	for _, test := range tests {
		t.Run(test.script, func(t *testing.T) {
			result := isApplicationScript(test.script)
			if result != test.expected {
				t.Errorf("expected %v but got %v", test.expected, result)
			}
		})
	}
}

func TestEnvFromMap(t *testing.T) {
	tests := []struct {
		description string
		env         map[string]string
		expected    []string
	}{
		{"nil", nil, nil},
		{"empty", map[string]string{}, nil},
		{"single", map[string]string{"a": "b"}, []string{"a=b"}},
		{"multiple", map[string]string{"a": "b", "c": "d"}, []string{"a=b", "c=d"}},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := envFromMap(test.env)
			sort.Strings(result)
			if len(result) != len(test.expected) {
				t.Errorf("expected %v but got %v", test.expected, result)
			} else {
				for i := 0; i < len(result); i++ {
					if result[i] != test.expected[i] {
						t.Errorf("expected %v but got %v", test.expected[i], result[i])
					}
				}
			}
		})
	}
}

func TestEnvToMap(t *testing.T) {
	tests := []struct {
		description string
		env         []string
		expected    map[string]string
	}{
		{"nil", nil, nil},
		{"empty", []string{}, nil},
		{"single", []string{"a=b"}, map[string]string{"a": "b"}},
		{"multiple", []string{"a=b", "c=d"}, map[string]string{"a": "b", "c": "d"}},
		{"collisions", []string{"a=b", "a=d"}, map[string]string{"a": "d"}},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := envToMap(test.env)
			if len(result) != len(test.expected) {
				t.Errorf("expected %v but got %v", test.expected, result)
			} else {
				for k, v := range result {
					if v != test.expected[k] {
						t.Errorf("for %v expected %v but got %v", k, test.expected[k], v)
					}
				}
			}
		})
	}
}

func TestStripInspectArg(t *testing.T) {
	tests := []struct {
		description string
		args        []string
		newArgs     []string
		inspectArg  string
	}{
		{"nil", nil, nil, ""},
		{"no args", []string{}, []string{}, ""},
		{"no inspect args", []string{"--foo", "bar"}, []string{"--foo", "bar"}, ""},
		{"lone <<<INSPECT>>> removed", []string{"<<<INSPECT>>>"}, []string{}, "<<<INSPECT>>>"},
		{"<<<INSPECT>>> at beginning removed", []string{"<<<INSPECT>>>", "--foo", "bar"}, []string{"--foo", "bar"}, "<<<INSPECT>>>"},
		{"<<<INSPECT>>> mid way removed", []string{"-c", "<<<INSPECT>>>", "--foo", "bar"}, []string{"-c", "--foo", "bar"}, "<<<INSPECT>>>"},
		{"<<<INSPECT>>> after script untouched", []string{"--foo", "bar", "<<<INSPECT>>>"}, []string{"--foo", "bar", "<<<INSPECT>>>"}, ""},
	}

	for _, test := range tests {
		// run the test for the difference inspect variants
		for _, inspect := range []string{"--inspect", "--inspect=9224", "--inspect-brk", "--inspect-brk=3452"} {
			test.description = strings.ReplaceAll(test.description, "<<<INSPECT>>>", inspect)
			for i := range test.args {
				if test.args[i] == "<<<INSPECT>>>" {
					test.args[i] = inspect
				}
			}
			for i := range test.newArgs {
				if test.newArgs[i] == "<<<INSPECT>>>" {
					test.newArgs[i] = inspect
				}
			}
			if test.inspectArg == "<<<INSPECT>>>" {
				test.inspectArg = inspect
			}
		}
		t.Run(test.description, func(t *testing.T) {
			newArgs, inspectArg := stripInspectArg(test.args)
			if len(newArgs) != len(test.newArgs) {
				t.Errorf("expected %v but got %v", test.newArgs, newArgs)
			} else {
				for i := 0; i < len(newArgs); i++ {
					if newArgs[i] != test.newArgs[i] {
						t.Errorf("expected %v but got %v", test.newArgs[i], newArgs[i])
					}
				}
			}
			if inspectArg != test.inspectArg {
				t.Errorf("expected %v but got %v", test.inspectArg, inspectArg)
			}
		})
	}
}

func TestNodeContext_unwrap(t *testing.T) {
	name := "foo" // ensure no code explicitly looks for "node"

	root, err := ioutil.TempDir("", "nc")
	if err != nil {
		t.Error(err)
	}
	originalNode := filepath.Join(root, name)
	if err := ioutil.WriteFile(originalNode, []byte{}, 0555); err != nil {
		t.Error(err)
	}
	binPath := filepath.Join(root, "bin")
	if err := os.Mkdir(binPath, 0777); err != nil {
		t.Error(err)
	}
	binNode := filepath.Join(binPath, name)
	if err := ioutil.WriteFile(binNode, []byte{}, 0555); err != nil {
		t.Error(err)
	}
	sbinPath := filepath.Join(root, "sbin")
	if err := os.Mkdir(sbinPath, 0777); err != nil {
		t.Error(err)
	}
	sbinNode := filepath.Join(sbinPath, name)
	if err := ioutil.WriteFile(sbinNode, []byte{}, 0555); err != nil {
		t.Error(err)
	}

	t.Cleanup(func() { os.RemoveAll(root) })

	tests := []struct {
		description string
		input       nodeContext
		result      bool
		expected    string
	}{
		{
			description: "no PATH leaves unchanged",
			input:       nodeContext{program: originalNode},
			result:      false,
			expected:    originalNode,
		},
		{
			description: "empty PATH leaves unchanged",
			input:       nodeContext{program: originalNode, env: map[string]string{"PATH": ""}},
			result:      false,
			expected:    originalNode,
		},
		{
			description: "no other node leaves unchanged",
			input:       nodeContext{program: originalNode, env: map[string]string{"PATH": root}},
			result:      false,
			expected:    originalNode,
		},
		{
			description: "first other node wins",
			input:       nodeContext{program: originalNode, env: map[string]string{"PATH": root + string(os.PathListSeparator) + binPath + string(os.PathListSeparator) + sbinPath}},
			result:      true,
			expected:    binNode,
		},
		{
			description: "first node wins when original not found",
			input:       nodeContext{program: name, env: map[string]string{"PATH": root + string(os.PathListSeparator) + binPath + string(os.PathListSeparator) + sbinPath}},
			result:      true,
			expected:    originalNode,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			copy := test.input
			result := copy.unwrap()
			if result != test.result {
				t.Errorf("expected unwrap() = %v but got %v", test.result, result)
			}
			if copy.program != test.expected {
				t.Errorf("expected %v but got %v", test.expected, copy.program)
			}
		})
	}
}

func TestNodeContext_StripInspectArgs(t *testing.T) {
	tests := []struct {
		description string
		input       nodeContext
		expected    nodeContext
		arg         string
	}{
		{
			description: "no inspect",
			input:       nodeContext{args: []string{"--no-warnings", "index.js"}, env: map[string]string{"NODE_OPTIONS": "--trace-sync-io"}},
			expected:    nodeContext{args: []string{"--no-warnings", "index.js"}, env: map[string]string{"NODE_OPTIONS": "--trace-sync-io"}},
			arg:         "",
		},
		{
			description: "inspect in args",
			input:       nodeContext{args: []string{"--inspect", "--no-warnings", "index.js"}, env: map[string]string{"NODE_OPTIONS": "--trace-sync-io"}},
			expected:    nodeContext{args: []string{"--no-warnings", "index.js"}, env: map[string]string{"NODE_OPTIONS": "--trace-sync-io"}},
			arg:         "--inspect",
		},
		{
			description: "inspect in NODE_OPTIONS",
			input:       nodeContext{args: []string{"--no-warnings", "index.js"}, env: map[string]string{"NODE_OPTIONS": "--trace-sync-io --inspect-brk"}},
			expected:    nodeContext{args: []string{"--no-warnings", "index.js"}, env: map[string]string{"NODE_OPTIONS": "--trace-sync-io"}},
			arg:         "--inspect-brk",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			copy := test.input
			arg := copy.stripInspectArgs()
			if arg != test.arg {
				t.Errorf("expected inspect args = %v but got %v", test.arg, arg)
			}
			if !reflect.DeepEqual(copy, test.expected) {
				t.Errorf("expected %v but got %v", test.expected, copy)
			}
		})
	}
}

func TestNodeContext_HandleNodemon(t *testing.T) {
	tests := []struct {
		description string
		input       nodeContext
		expected    nodeContext
	}{
		{
			description: "no nodemon",
			input:       nodeContext{args: []string{"--no-warnings", "index.js"}, env: map[string]string{"NODE_DEBUG": "--inspect=3333"}},
			expected:    nodeContext{args: []string{"--no-warnings", "index.js"}, env: map[string]string{"NODE_DEBUG": "--inspect=3333"}},
		},
		{
			description: "nodemon no args",
			input:       nodeContext{args: []string{"--no-warnings", "./node_modules/nodemon/bin/nodemon.js"}, env: map[string]string{"NODE_DEBUG": "--inspect=3333"}},
			expected:    nodeContext{args: []string{"--no-warnings", "./node_modules/nodemon/bin/nodemon.js", "--inspect=3333"}, env: map[string]string{}},
		},
		{
			description: "nodemon with args",
			input:       nodeContext{args: []string{"--no-warnings", "./node_modules/nodemon/bin/nodemon.js", "-v", "index.js"}, env: map[string]string{"NODE_DEBUG": "--inspect=3333"}},
			expected:    nodeContext{args: []string{"--no-warnings", "./node_modules/nodemon/bin/nodemon.js", "--inspect=3333", "-v", "index.js"}, env: map[string]string{}},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			copy := test.input
			copy.handleNodemon()
			if !reflect.DeepEqual(copy, test.expected) {
				t.Errorf("mismatch\nexpected: %v\n but got: %v", test.expected, copy)
			}
		})
	}
}

func TestNodeContext_AddNodeArg(t *testing.T) {
	tests := []struct {
		description string
		input       nodeContext
		arg         string
		expected    nodeContext
	}{
		{
			description: "nil args",
			input:       nodeContext{},
			arg:         "abc",
			expected:    nodeContext{args: []string{"abc"}},
		},
		{
			description: "no script",
			input:       nodeContext{args: []string{"--no-warnings"}},
			arg:         "abc",
			expected:    nodeContext{args: []string{"--no-warnings", "abc"}},
		},
		{
			description: "after options before script",
			input:       nodeContext{args: []string{"--no-warnings", "index.js"}},
			arg:         "abc",
			expected:    nodeContext{args: []string{"--no-warnings", "abc", "index.js"}},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			copy := test.input
			copy.addNodeArg(test.arg)
			if !reflect.DeepEqual(copy, test.expected) {
				t.Errorf("expected %v but got %v", test.expected, copy)
			}
		})
	}
}

func TestIntegration(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("we only support nix")
	}
	root, err := ioutil.TempDir("", "node")
	if err != nil {
		t.Error(err)
	}

	actualNode := filepath.Join(root, "nodeBin")
	script := `#!/bin/sh
if [ -n "$NODE_DEBUG" ]; then
  echo "NODE_DEBUG=$NODE_DEBUG"
fi
if [ -n "$NODE_OPTIONS" ]; then
  echo "NODE_OPTIONS=$NODE_OPTIONS"
fi
for arg in "$@"; do
  echo "$arg"
done
`
	if err := ioutil.WriteFile(actualNode, []byte(script), 0555); err != nil {
		t.Errorf("could not create node script: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(root) })

	tests := []struct {
		description string
		args        []string
		env         map[string]string
		expected    string
	}{
		// app scripts are terminal: commands should only affected if NODE_DEBUG is defined
		{
			description: "app script: passed through",
			args:        []string{"script.js"},
			expected:    "script.js\n",
		},
		{
			description: "app script: inspect arg passed through",
			args:        []string{"--inspect", "script.js"},
			expected:    "--inspect\nscript.js\n",
		},
		{
			description: "app script: inspect as app args left alone",
			args:        []string{"script.js", "--inspect"},
			expected:    "script.js\n--inspect\n",
		},
		{
			description: "app script with NODE_OPTIONS='--inspect': passed through",
			args:        []string{"script.js"},
			env:         map[string]string{"NODE_OPTIONS": "--inspect"},
			expected:    "NODE_OPTIONS=--inspect\nscript.js\n",
		},
		{
			description: "app script with NODE_OPTIONS='--foo --inspect --bar': passed through",
			args:        []string{"script.js"},
			env:         map[string]string{"NODE_OPTIONS": "--foo --inspect --bar"},
			expected:    "NODE_OPTIONS=--foo --inspect --bar\nscript.js\n",
		},
		{
			description: "app script with NODE_DEBUG='--inspect': installed",
			args:        []string{"script.js"},
			env:         map[string]string{"NODE_DEBUG": "--inspect"},
			expected:    "--inspect\nscript.js\n",
		},

		// node_module scripts should have --inspect stripped and propagated,
		// and NODE_DEBUG should never be overwritten
		{
			description: "node_modules script: passed through",
			args:        []string{"node_modules/script.js"},
			expected:    "node_modules/script.js\n",
		},
		{
			description: "node_modules script: inspect as app args left alone",
			args:        []string{"node_modules/script.js", "--inspect"},
			expected:    "node_modules/script.js\n--inspect\n",
		},
		{
			description: "node_modules script with inspect: seeds NODE_DEBUG",
			args:        []string{"--inspect", "node_modules/script.js"},
			expected:    "NODE_DEBUG=--inspect\nnode_modules/script.js\n",
		},
		{
			description: "node_modules script with NODE_OPTIONS='--inspect': seeds NODE_DEBUG",
			args:        []string{"node_modules/script.js"},
			env:         map[string]string{"NODE_OPTIONS": "--inspect"},
			expected:    "NODE_DEBUG=--inspect\nnode_modules/script.js\n",
		},
		{
			description: "node_modules script with NODE_OPTIONS='--foo --inspect --bar': seeds NODE_DEBUG",
			args:        []string{"node_modules/script.js"},
			env:         map[string]string{"NODE_OPTIONS": "--foo --inspect --bar"},
			expected:    "NODE_DEBUG=--inspect\nNODE_OPTIONS=--foo --bar\nnode_modules/script.js\n",
		},
		{
			description: "node_modules script with NODE_DEBUG='--inspect': passed through",
			args:        []string{"node_modules/script.js"},
			env:         map[string]string{"NODE_DEBUG": "--inspect"},
			expected:    "NODE_DEBUG=--inspect\nnode_modules/script.js\n",
		},
		{
			description: "node_modules script with NODE_DEBUG='--inspect' and inspect-brk arg: NODE_DEBUG wins",
			args:        []string{"--inspect-brk", "node_modules/script.js"},
			env:         map[string]string{"NODE_DEBUG": "--inspect"},
			expected:    "NODE_DEBUG=--inspect\nnode_modules/script.js\n",
		},
		{
			description: "node_modules script with NODE_DEBUG='--inspect' and inspect-brk NODE_OPTIONS: NODE_DEBUG wins",
			args:        []string{"node_modules/script.js"},
			env:         map[string]string{"NODE_DEBUG": "--inspect", "NODE_OPTIONS": "--inspect-brk"},
			expected:    "NODE_DEBUG=--inspect\nnode_modules/script.js\n",
		},
		{
			description: "nodemon script with NODE_DEBUG='--inspect': added to nodemon",
			args:        []string{"node_modules/nodemon/nodemon.js"},
			env:         map[string]string{"NODE_DEBUG": "--inspect"},
			expected:    "node_modules/nodemon/nodemon.js\n--inspect\n",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			env := map[string]string{"PATH": root + string(os.PathListSeparator) + os.Getenv("PATH")}
			for k, v := range test.env {
				env[k] = v
			}

			nc := nodeContext{program: "nodeBin", args: test.args, env: env}
			var in bytes.Buffer
			var out bytes.Buffer
			if err := run(&nc, &in, &out, &out); err != nil {
				t.Errorf("node exec failed: %v", err)
			}
			if nc.program != actualNode {
				t.Errorf("unwrap resolved to %q but wanted %q", nc.program, actualNode)
			}
			if out.String() != test.expected {
				t.Errorf("output mismatch\nexpected: %q\n but got: %q", test.expected, out.String())
			}
		})
	}

}
