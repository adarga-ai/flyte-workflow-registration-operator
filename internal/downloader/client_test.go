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

package downloader

import (
	"context"
	"testing"

	"github.com/adarga-ai/flyte-workflow-registration-operator/internal"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		config := internal.Config{DownloadStrategy: "jfrog"}
		cxt := context.Background()
		_, err := NewClient(cxt, config)

		assert.NoError(t, err)
	})

	t.Run("failure case", func(t *testing.T) {
		config := internal.Config{DownloadStrategy: "not-correct"}
		cxt := context.Background()
		_, err := NewClient(cxt, config)

		assert.ErrorContains(t, err, "invalid downloader strategy")
	})
}
