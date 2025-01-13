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

// Package flyte contains all code for interacting with flyte
package flyte

import (
	"context"
	"fmt"

	"github.com/adarga-ai/flyte-workflow-registration-operator/internal/command"
)

// Client is an interface for the flyte admin client
//
//go:generate mockery --name=Client
type Client interface {
	RegisterWorkflow(ctx context.Context, tgzPath string, meta WorkflowMetadata, auth Auth) error
}

// AdminClient is a wrapper for interactions with FlyteAdmin
type AdminClient struct {
	Executor command.Executor
}

// NewClient creates a new instance of flyte wrapper
func NewClient(executor command.Executor) *AdminClient {
	return &AdminClient{
		Executor: executor,
	}
}

// WorkflowMetadata is a struct that holds the information
// to represent a single ticker
type WorkflowMetadata struct {
	WorkflowVersion string
	Domain          string
	Project         string
}

// Auth is a struct that contains the endpoint for the flyte admin server as well as the client credentials.
type Auth struct {
	AdminEndpoint      string
	ClientID           string
	ClientSecretEnvVar string
}

// RegisterWorkflow registers workflow using flytectl with a provided .tgz file.
func (a *AdminClient) RegisterWorkflow(ctx context.Context, tgzPath string, meta WorkflowMetadata, auth Auth) error {
	args := []string{
		"register",
		"files",
		"--archive", tgzPath,
		"--project", meta.Project,
		"--domain", meta.Domain,
		"--version", meta.WorkflowVersion,
		"--admin.endpoint", auth.AdminEndpoint,
		"--admin.authType", "ClientSecret",
		"--admin.clientId", auth.ClientID,
		"--admin.clientSecretEnvVar", auth.ClientSecretEnvVar,
	}
	output, err := a.Executor.ExecuteCommand(ctx, "flytectl", args...)
	if err != nil {
		return fmt.Errorf("failed to execute flytectl: %w, output: %s", err, output)
	}
	return nil
}
