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

package jfrog

import (
	"context"
	"testing"

	"github.com/adarga-ai/flyte-workflow-registration-operator/internal"
	"github.com/adarga-ai/flyte-workflow-registration-operator/internal/jfrog/mocks"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Used to generate tests for the ArtifactDownloader interface
//go:generate mockery --srcpkg=github.com/jfrog/jfrog-client-go/artifactory --name=ArtifactoryServicesManager

func TestDownloadArtifact(t *testing.T) {
	t.Run("we can download a file successfully", func(t *testing.T) {
		// Set up the mocks
		jFrogManagerMock := mocks.NewArtifactoryServicesManager(t)
		numDownloaded := 1
		numFailed := 0
		jFrogManagerMock.On("DownloadFiles", mock.Anything).Return(numDownloaded, numFailed, nil)
		j := Downloader{
			JFrogManager: jFrogManagerMock,
		}

		// Call the method we are testing
		result, err := j.DownloadArtifact(context.Background(), "packagePath", "1.2.3")
		assert.NoError(t, err)

		// Assert the result
		assert.Contains(t, result, "packagePath")
	})
}

func TestSetupDownloader(t *testing.T) {
	type fields struct {
		Config       internal.Config
		JFrogManager artifactory.ArtifactoryServicesManager
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "jfrog can be setup with correct values",
			fields: fields{
				Config: internal.Config{
					DownloadStrategy: "jfrog",
					JfrogURL:         "url",
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := Downloader{
				Config:       tt.fields.Config,
				JFrogManager: tt.fields.JFrogManager,
			}
			if err := j.SetupDownloader(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("JFrogDownloader.SetupDownloader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
