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

/*
Package downloader handles the downloading of a workflow artifacts from a download source

The idea is that we can easily define what downloader strategy to use and new download strategies
can be swapped. For now workflow packages will be stored in an OCI registry, but they could as easily be stored in jFrog.

We use a strategy pattern, e.g a `config.DownloadStrategy` is defined as an enum e.g `jfrog` or `oci`
which instantiates a given specific strategy and returns it from the NewDownloader function.
Each strategy should implement the DownloadArtifact function which downloaders the artifact from
its source.
*/
package downloader

import (
	"context"
	"errors"

	"github.com/adarga-ai/flyte-workflow-registration-operator/internal"
	"github.com/adarga-ai/flyte-workflow-registration-operator/internal/jfrog"
	"github.com/adarga-ai/flyte-workflow-registration-operator/internal/oci"
)

// Client is an interface for downloading artifacts
//
//go:generate mockery --name=Client
type Client interface {
	DownloadArtifact(ctx context.Context, uri string, version string) (string, error)
}

// NewClient returns the relevant downloader based on the given download strategy
func NewClient(ctx context.Context, cfg internal.Config) (Client, error) {
	switch cfg.DownloadStrategy {
	case internal.DownloadStrategyOCI:
		return oci.NewDownloader(cfg)
	case internal.DownloadStrategyJFrog:
		d := jfrog.Downloader{Config: cfg}
		if err := d.SetupDownloader(ctx); err != nil {
			return nil, errors.New("failed to setup JFrog downloader")
		}

		return &d, nil
	default:
		return nil, errors.New("invalid downloader strategy")
	}
}
