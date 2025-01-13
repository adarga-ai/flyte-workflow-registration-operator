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

// Package oci contains a struct used for downloading artifacts from OCI
package oci

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adarga-ai/flyte-workflow-registration-operator/internal"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Downloader is a struct that exposes a function for downloading an artifact from OCI.
type Downloader struct {
	cfg internal.Config
}

// NewDownloader returns a downloader from a given config.
func NewDownloader(cfg internal.Config) (*Downloader, error) {
	return &Downloader{
		cfg: cfg,
	}, nil
}

// DownloadArtifact downloads an artifact from the OCI registry.
// The uri is the path to the artifact in the OCI registry. For example adarga/ds-wf-relationships-extraction
// The version is the version of the artifact to download. For example 0.1.0
func (d *Downloader) DownloadArtifact(ctx context.Context, uri string, version string) (string, error) {
	logger := log.FromContext(ctx)
	logger.Info("downloading artifact...",
		"artifact", uri,
		"version", version,
	)

	// Create a new OCI repository. We need to create a new repository when we download a new artifact as each
	// workflow package could be stored in a separate repository
	repo, err := remote.NewRepository(fmt.Sprintf("%s/%s", d.cfg.OCIRegistry, uri))
	if err != nil {
		return "", fmt.Errorf("failed to create OCI repository: %w", err)
	}

	// Configure the authentication for the OCI repository. We need to do this each time we download an artifact
	// to prevent the credentials from expiring
	credential, err := getCredential(ctx, d.cfg)
	if err != nil {
		return "", fmt.Errorf("failed to get OCI credentials: %w", err)
	}

	repo.Client = &auth.Client{
		Client:     retry.DefaultClient,
		Cache:      auth.NewCache(),
		Credential: credential,
	}

	// Delete any existing files in the directory before downloading the new artifact
	base := filepath.Base(uri)
	basedir := filepath.Join(os.TempDir(), base)
	err = os.RemoveAll(basedir)
	if err != nil {
		return "", fmt.Errorf("removing existing directory: %w", err)
	}

	// Download the artifact to a local file
	fs, err := file.New(basedir)
	if err != nil {
		return "", fmt.Errorf("failed to create local file system: %w", err)
	}

	_, err = oras.Copy(ctx, repo, version, fs, version, oras.DefaultCopyOptions)
	if err != nil {
		return "", fmt.Errorf("failed to copy artifact from remote: %w", err)
	}

	err = fs.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close file store: %w", err)
	}

	// Obtain the name of the downloaded file. Because the downloaded artifact is opaque to us, we need to
	// read the directory and get the name of the file
	files, err := os.ReadDir(basedir)
	if err != nil {
		return "", fmt.Errorf("failed to read downloaded files: %w", err)
	}

	if len(files) != 1 {
		for i, file := range files {
			logger.Info("too many files downloaded",
				"file number", i+1,
				"file name", file.Name(),
			)
		}
		return "", fmt.Errorf("expected 1 file in the downloaded directory, got %d", len(files))
	}

	return filepath.Join(basedir, files[0].Name()), nil
}

func getCredential(ctx context.Context, cfg internal.Config) (auth.CredentialFunc, error) {
	var username string
	var password string

	switch cfg.OCIAuthStrategy {
	case internal.OCIAuthStrategyStatic:
		username = cfg.OCIUsername
		password = cfg.OCIPassword
	case internal.OCIAuthStrategyECR:
		username = "AWS"

		awsCfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS configuration: %w", err)
		}

		// Create an ECR client
		ecrClient := ecr.NewFromConfig(awsCfg)

		// Get an authorization token
		input := &ecr.GetAuthorizationTokenInput{}
		result, err := ecrClient.GetAuthorizationToken(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to get ECR authorization token: %w", err)
		}

		if len(result.AuthorizationData) == 0 {
			return nil, fmt.Errorf("no ECR authorization token found")
		}

		encodedToken := aws.ToString(result.AuthorizationData[0].AuthorizationToken)
		decodedToken, err := base64.StdEncoding.DecodeString(encodedToken)
		if err != nil {
			return nil, fmt.Errorf("failed to decode ECR authorization token: %w", err)
		}

		// Split the decoded token to get the password
		tokenParts := strings.SplitN(string(decodedToken), ":", 2)
		if len(tokenParts) != 2 {
			return nil, fmt.Errorf("invalid ECR authorization token: %s", decodedToken)
		}

		// Extract the password
		password = tokenParts[1]
	default:
		return nil, fmt.Errorf("unknown OCI auth strategy: %s, only `static` or `ecr` allowed", cfg.OCIAuthStrategy)
	}

	return auth.StaticCredential(cfg.OCIRegistry, auth.Credential{
		Username: username,
		Password: password,
	}), nil
}
