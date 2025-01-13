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

// Package internal holds the logic for getting environment variables loaded into internal config
package internal

import (
	"fmt"

	"github.com/alexflint/go-arg"
)

// DownloadStrategyOCI is a value for the Downloadstrategy env var that is set to ensure artifacts are downloaded from
// OCI
const DownloadStrategyOCI = "oci"

// DownloadStrategyJFrog is a value for the Downloadstrategy env var that is set to ensure artifacts are downloaded from
// JFrog
const DownloadStrategyJFrog = "jfrog"

// OCIAuthStrategyStatic is a value for the OCIAuthStrategy env var, when set to this OCI will be authenticated to
// statically
const OCIAuthStrategyStatic = "static"

// OCIAuthStrategyECR is a value for the OCIAuthStrategy env var, when set to this ECR will be used to authenticate with
// OCI
const OCIAuthStrategyECR = "ecr"

// Config is the configuration to run the service
// args are parsed from go-arg, https://github.com/alexflint/go-arg
// Add here service config arguments and add the specific arg tag
type Config struct {
	DownloadStrategy string `arg:"env:DOWNLOADER_STRATEGY" default:"oci"`
	LogLevel         string `arg:"env:LOG_LEVEL" default:"info"`

	// JFrog config
	JfrogURL      string `arg:"env:JFROG_ARTIFACTORY_URL"`
	JfrogUser     string `arg:"env:JFROG_USER"`
	JfrogPassword string `arg:"env:JFROG_PASSWORD"`

	// OCI config
	OCIRegistry     string `arg:"env:OCI_REGISTRY"`
	OCIAuthStrategy string `arg:"env:OCI_AUTH_STRATEGY" default:"static"`
	OCIUsername     string `arg:"env:OCI_USERNAME"`
	OCIPassword     string `arg:"env:OCI_PASSWORD"`

	// Flyte config
	FlyteAdminEndpoint string `arg:"env:FLYTE_ADMIN_ENDPOINT"`
	FlyteClientID      string `arg:"env:FLYTE_CLIENT_ID"`
	FlyteClientSecret  string `arg:"env:FLYTE_CLIENT_SECRET"`
}

// NewConfig return a new instance of Config
func NewConfig() (Config, error) {
	var config Config
	err := arg.Parse(&config)
	if err != nil {
		return Config{}, fmt.Errorf("could not parse config: %w", err)
	}

	// Validate the config
	if config.DownloadStrategy != DownloadStrategyOCI && config.DownloadStrategy != DownloadStrategyJFrog {
		return Config{}, fmt.Errorf("invalid download strategy: %s, only `oci` or `jfrog` allowed", config.DownloadStrategy)
	}

	if config.OCIAuthStrategy != OCIAuthStrategyStatic && config.OCIAuthStrategy != OCIAuthStrategyECR {
		return Config{}, fmt.Errorf("invalid OCI auth strategy: %s, only `static` or `ecr` allowed", config.OCIAuthStrategy)
	}

	return config, nil
}
