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

package integration

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	inferencev1 "github.com/adarga-ai/flyte-workflow-registration-operator/api/v1"
	"github.com/adarga-ai/flyte-workflow-registration-operator/internal"
	"github.com/adarga-ai/flyte-workflow-registration-operator/internal/controller"
	dMocks "github.com/adarga-ai/flyte-workflow-registration-operator/internal/downloader/mocks"
	"github.com/adarga-ai/flyte-workflow-registration-operator/internal/flyte"
	fMocks "github.com/adarga-ai/flyte-workflow-registration-operator/internal/flyte/mocks"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	//+kubebuilder:scaffold:imports

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	k8sClient client.Client
	testEnv   *envtest.Environment
)

func TestIntegration(t *testing.T) {
	workflowVersion := "1.0.0"
	workflowDomain := "test-domain"
	workflowProject := "test-project"
	workflowPackageURI := "test-uri"

	artifactPath := "test-artifact-path"

	meta := flyte.WorkflowMetadata{
		WorkflowVersion: workflowVersion,
		Domain:          workflowDomain,
		Project:         workflowProject,
	}

	flyteClientID := "test-client-id"
	flyteAdminEndpoint := "test-endpoint"

	flyteAuth := flyte.Auth{
		AdminEndpoint:      flyteAdminEndpoint,
		ClientID:           flyteClientID,
		ClientSecretEnvVar: "FLYTE_CLIENT_SECRET",
	}

	mockDownloader := dMocks.Client{}
	mockFlyteClient := fMocks.Client{}

	mockDownloader.EXPECT().DownloadArtifact(mock.Anything, workflowPackageURI, workflowVersion).Return(artifactPath, nil).Once()

	mockFlyteClient.EXPECT().RegisterWorkflow(mock.Anything, artifactPath, meta, flyteAuth).Return(nil).Once()

	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,

		// The BinaryAssetsDirectory is only required if you want to run the tests directly
		// without call the makefile target test. If not informed it will look for the
		// default path defined in controller-runtime which is /usr/local/kubebuilder/.
		// Note that you must have the required binaries setup under the bin directory to perform
		// the tests directly. When we run make test it will be setup and used automatically.
		BinaryAssetsDirectory: filepath.Join("..", "bin", "k8s",
			fmt.Sprintf("1.28.0-%s-%s", runtime.GOOS, runtime.GOARCH)),
	}

	cfg, err := testEnv.Start()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	defer func() {
		err = testEnv.Stop()
		assert.NoError(t, err, "unexpected error stopping test environment")
	}()

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	require.NoError(t, err)

	err = inferencev1.AddToScheme(scheme.Scheme)
	require.NoError(t, err)

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	require.NoError(t, err)
	require.NotNil(t, k8sClient)

	r := &controller.FlyteRegistrationReconciler{
		K8sClient:        k8sClient,
		Scheme:           scheme.Scheme,
		Config:           internal.Config{FlyteClientID: flyteClientID, FlyteAdminEndpoint: flyteAdminEndpoint},
		Downloader:       &mockDownloader,
		FlyteAdminClient: &mockFlyteClient,
	}

	err = r.SetupWithManager(k8sManager)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		require.NoError(t, k8sManager.Start(ctx))
	}()

	testNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-namespace",
		},
	}
	err = k8sClient.Create(context.Background(), testNamespace)
	require.NoError(t, err)

	defer func() {
		err := k8sClient.Delete(context.Background(), testNamespace)
		require.NoError(t, err)
	}()

	flyteRegistration := &inferencev1.FlyteRegistration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "relationships-extraction",
		},
		Spec: inferencev1.FlyteRegistrationSpec{
			WorkflowDomain:     workflowDomain,
			WorkflowProject:    workflowProject,
			WorkflowPackageURI: workflowPackageURI,
			WorkflowVersion:    workflowVersion,
		},
	}

	flyteRegistration.Namespace = "test-namespace"

	err = k8sClient.Create(context.Background(), flyteRegistration)
	require.NoError(t, err)

	cancel()
}
