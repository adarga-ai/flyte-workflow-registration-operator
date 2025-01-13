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

package flyte

import (
	"context"
	"errors"
	"testing"

	"github.com/adarga-ai/flyte-workflow-registration-operator/internal/command/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewClient(t *testing.T) {
	mockCommandExecutor := mocks.Executor{}

	expectedClient := AdminClient{
		Executor: &mockCommandExecutor,
	}

	c := NewClient(&mockCommandExecutor)

	assert.Equal(t, &expectedClient, c)
}

func TestRegisterWorkflow(t *testing.T) {
	// SHARED INPUTS
	command := "flytectl"

	tgzPath := "test-tgz"

	meta := WorkflowMetadata{
		WorkflowVersion: "1.0.0",
		Domain:          "test-domain",
		Project:         "test-project",
	}

	flyteClientID := "test-client-id"
	flyteAdminEndpoint := "test-endpoint"
	flyteClientSecretEnvVar := "FLYTE_CLIENT_SECRET"

	flyteAuth := Auth{
		AdminEndpoint:      flyteAdminEndpoint,
		ClientID:           flyteClientID,
		ClientSecretEnvVar: flyteClientSecretEnvVar,
	}

	args := []interface{}{
		"register",
		"files",
		"--archive", tgzPath,
		"--project", meta.Project,
		"--domain", meta.Domain,
		"--version", meta.WorkflowVersion,
		"--admin.endpoint", flyteAdminEndpoint,
		"--admin.authType", "ClientSecret",
		"--admin.clientId", flyteClientID,
		"--admin.clientSecretEnvVar", flyteClientSecretEnvVar,
	}

	// SHARED MOCKS
	mockCommandExecutor := mocks.Executor{}

	t.Run("success case", func(t *testing.T) {
		// MOCK BEHAVIOUR
		mockCommandExecutor.EXPECT().ExecuteCommand(mock.Anything, command, args...).Return(nil, nil).Once()

		// EXECUTION
		c := NewClient(&mockCommandExecutor)

		err := c.RegisterWorkflow(context.Background(), tgzPath, meta, flyteAuth)

		// ASSERTIONS
		assert.NoError(t, err)
	})

	t.Run("failure case: execute command error", func(t *testing.T) {
		// MOCK BEHAVIOUR
		mockCommandExecutor.EXPECT().ExecuteCommand(mock.Anything, command, args...).Return(nil, errors.New("test error")).Once()

		// EXECUTION
		c := NewClient(&mockCommandExecutor)

		err := c.RegisterWorkflow(context.Background(), tgzPath, meta, flyteAuth)

		// ASSERTIONS
		assert.ErrorContains(t, err, "failed to execute flytectl")
	})
}
