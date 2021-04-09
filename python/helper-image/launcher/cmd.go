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
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

// for testing
var newCommand = createCommand
var newConsoleCommand = createConsoleCommand

// commander is a subset of exec.Cmd
type commander interface {
	Run() error
	Output() ([]byte, error)
	CombinedOutput() ([]byte, error)
}

var _ commander = (*exec.Cmd)(nil)

// createCommand creates a normal exec.Cmd object
func createCommand(ctx context.Context, cmdline []string, env env) commander {
	logrus.Debugf("command: %v (env: %s)", cmdline, env)
	cmd := exec.CommandContext(ctx, cmdline[0], cmdline[1:]...)
	cmd.Env = env.AsPairs()
	return cmd
}

// createConsoleCommand creates an exec.Cmd object that connects to os.Stdin, os.Stdout, os.Stderr
func createConsoleCommand(ctx context.Context, cmdline []string, env env) commander {
	logrus.Debugf("command(stdin/out/err): %v (env: %s)", cmdline, env)
	cmd := exec.CommandContext(ctx, cmdline[0], cmdline[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env.AsPairs()
	return cmd
}
