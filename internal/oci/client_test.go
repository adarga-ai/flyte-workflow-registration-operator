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

package oci

import (
	"context"
	"testing"

	"github.com/adarga-ai/flyte-workflow-registration-operator/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCredentials(t *testing.T) {
	t.Run("we can get static credentials successfully", func(t *testing.T) {
		cfg := internal.Config{
			OCIRegistry:     "localhost:1234",
			OCIUsername:     "test-username",
			OCIPassword:     "test-password",
			OCIAuthStrategy: "static",
		}

		credsFunc, err := getCredential(context.Background(), cfg)
		require.NoError(t, err)

		creds, err := credsFunc(context.Background(), "localhost:1234")
		require.NoError(t, err)

		assert.Equal(t, "test-username", creds.Username)
		assert.Equal(t, "test-password", creds.Password)
	})
}
