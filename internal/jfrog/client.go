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

// Package jfrog holds the logic for interacting with jfrog
package jfrog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adarga-ai/flyte-workflow-registration-operator/internal"
	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/config"
	"k8s.io/kube-openapi/pkg/validation/errors"
)

// Downloader holds the dependencies for downloading files from Jfrog
type Downloader struct {
	Config       internal.Config
	JFrogManager artifactory.ArtifactoryServicesManager
}

// DownloadArtifact downloads an artifact from JFrog
func (d *Downloader) DownloadArtifact(_ context.Context, uri string, version string) (string, error) {
	packagePath := fmt.Sprintf("%s_%s.tgz", uri, version)

	params := services.NewDownloadParams()

	// The full path to the artifact in Artifactory, including the repository name
	params.Pattern = packagePath

	// Dynamically determine the folder name from the packagePath
	folderPath := filepath.Dir(packagePath) // This extracts the folder path from the full package path
	folderName := filepath.Base(folderPath) // This gets the actual folder name

	// The local file system path where the downloaded files will be saved
	params.Target = filepath.Join(os.TempDir(), filepath.Base(packagePath))

	params.SplitCount = 2

	params.MinSplitSize = 7168

	totalDownloaded, _, err := d.JFrogManager.DownloadFiles(params)
	// Handler errors
	if err != nil {
		return "", fmt.Errorf("jfrog manager: downloading files: %w", err)
	}
	if totalDownloaded < 1 {
		return "", errors.New(1, fmt.Sprintf("no files to download: %s", packagePath))
	}

	return filepath.Join(os.TempDir(), folderName, filepath.Base(packagePath)), nil
}

// SetupDownloader sets up the JFrog downloader
func (d *Downloader) SetupDownloader(ctx context.Context) error {
	rtDetails := auth.NewArtifactoryDetails()
	rtDetails.SetUrl(d.Config.JfrogURL)
	rtDetails.SetUser(d.Config.JfrogUser)
	rtDetails.SetPassword(d.Config.JfrogPassword)

	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(rtDetails).
		SetDryRun(false).
		SetContext(ctx).
		Build()
	if err != nil {
		// Handle error
		return errors.New(1, "error creating service config")
	}

	rtManager, err := artifactory.New(serviceConfig)
	if err != nil {
		return err
	}

	d.JFrogManager = rtManager

	return nil
}
