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
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCmd(t *testing.T) {
	RunCmdOut([]string{"hello"}, "abc").
		AndRunCmdFail([]string{"ls"}, 1).
		Setup(t)

	if out, err := newCommand(nil, []string{"hello"}, nil).Output(); err != nil {
		t.Error("command should not have failed")
	} else if string(out) != "abc" {
		t.Error("output should have been abc")
	}
	if newCommand(nil, []string{"ls"}, nil).Run() == nil {
		t.Error("command should have failed")
	}
}

type fakeCmd struct {
	mode     string
	cmdline  []string
	exitCode int
	output   string
}

var _ commander = (*fakeCmd)(nil)

func (f *fakeCmd) Run() error {
	_t.Helper()
	if f.mode != "Run" {
		_t.Errorf("Command%v: expected %s() not Run()", f.cmdline, f.mode)
	}
	if f.exitCode == 0 {
		return nil
	}
	// doesn't seem to be an easy way to set the exitcode
	return &exec.ExitError{}
}

func (f *fakeCmd) Output() ([]byte, error) {
	_t.Helper()
	if f.mode != "Output" {
		_t.Errorf("Command%v: expected %v() not Output()", f.cmdline, f.mode)
	}
	if f.exitCode == 0 {
		return []byte(f.output), nil
	}
	// doesn't seem to be an easy way to set the exitcode
	return []byte(f.output), &exec.ExitError{}
}

func (f *fakeCmd) CombinedOutput() ([]byte, error) {
	_t.Helper()
	if f.mode != "Output" {
		_t.Errorf("Command%v: expected %v() not CombinedOutput()", f.cmdline, f.mode)
	}
	if f.exitCode == 0 {
		return []byte(f.output), nil
	}
	// doesn't seem to be an easy way to set the exitcode
	return []byte(f.output), &exec.ExitError{}
}

type commands []*fakeCmd

var (
	_cmdStack commands
	_t        *testing.T
)

func fakeCommand(_ context.Context, cmdline []string, env env) commander {
	_t.Helper()
	if len(_cmdStack) == 0 {
		_t.Fatalf("test expected no further commands: %v", cmdline)
	}
	current := _cmdStack[0]
	_cmdStack = _cmdStack[1:]
	if diff := cmp.Diff(current.cmdline, cmdline); diff != "" {
		_t.Errorf("cmdlines differ (-got, +want): %s", diff)
	}
	return current
}

func (c commands) Setup(t *testing.T) {
	_t = t
	_cmdStack = c

	oldCommand := newCommand
	oldConsoleCommand := newConsoleCommand
	newCommand = fakeCommand
	newConsoleCommand = fakeCommand
	_t.Cleanup(func() {
		newCommand = oldCommand
		newConsoleCommand = oldConsoleCommand
	})
}

func RunCmd(cmdline []string) commands {
	return commands{}.AndRunCmd(cmdline)
}

func RunCmdFail(cmdline []string, exitCode int) commands {
	return commands{}.AndRunCmdFail(cmdline, exitCode)
}

func RunCmdOut(cmdline []string, output string) commands {
	return commands{}.AndRunCmdOut(cmdline, output)
}

func RunCmdOutFail(cmdline []string, output string, exitCode int) commands {
	return commands{}.AndRunCmdOutFail(cmdline, output, exitCode)
}

func (c commands) AndRunCmd(cmdline []string) commands {
	c = append(c, &fakeCmd{mode: "Run", cmdline: cmdline})
	return c
}

func (c commands) AndRunCmdFail(cmdline []string, exitCode int) commands {
	c = append(c, &fakeCmd{mode: "Run", cmdline: cmdline, exitCode: exitCode})
	return c
}

func (c commands) AndRunCmdOut(cmdline []string, output string) commands {
	c = append(c, &fakeCmd{mode: "Output", cmdline: cmdline, output: output})
	return c
}

func (c commands) AndRunCmdOutFail(cmdline []string, output string, exitCode int) commands {
	c = append(c, &fakeCmd{mode: "Output", cmdline: cmdline, output: output, exitCode: exitCode})
	return c
}
