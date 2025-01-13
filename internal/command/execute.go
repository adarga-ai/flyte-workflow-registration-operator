/*
Copyright 2025 Adarga Limited.

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

// Package command contains the logic for executing os commands
package command

import (
	"context"
	"os/exec"
)

// Executor is an interface for executing os commands
//
//go:generate mockery --name=Executor
type Executor interface {
	ExecuteCommand(ctx context.Context, command string, args ...string) ([]byte, error)
}

// OSCommandExecutor is an instance of Executor
type OSCommandExecutor struct{}

// NewOSCommandExecutor creates a new instance of the command executor
func NewOSCommandExecutor() *OSCommandExecutor {
	return &OSCommandExecutor{}
}

// ExecuteCommand executes a given os command
func (e *OSCommandExecutor) ExecuteCommand(ctx context.Context, command string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	return cmd.CombinedOutput()
}
