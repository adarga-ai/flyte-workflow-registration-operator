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

package controller

import (
	"context"
	"errors"
	"testing"

	v1 "github.com/adarga-ai/flyte-workflow-registration-operator/api/v1"
	"github.com/adarga-ai/flyte-workflow-registration-operator/internal"
	"github.com/adarga-ai/flyte-workflow-registration-operator/internal/controller/mocks"
	dMocks "github.com/adarga-ai/flyte-workflow-registration-operator/internal/downloader/mocks"
	"github.com/adarga-ai/flyte-workflow-registration-operator/internal/flyte"
	fMocks "github.com/adarga-ai/flyte-workflow-registration-operator/internal/flyte/mocks"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcile(t *testing.T) {
	// SHARED INPUTS
	testWorkflow := v1.FlyteRegistration{}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{Namespace: "test", Name: "test"},
	}

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

	// SHARED MOCKS
	mockK8sClient := &mocks.K8sClient{}
	mockDownloader := &dMocks.Client{}
	mockFlyteAdminClient := &fMocks.Client{}

	t.Run("success case", func(t *testing.T) {
		// MOCK BEHAVIOUR
		mockK8sClient.EXPECT().Get(mock.Anything, req.NamespacedName, &testWorkflow).Once().
			Run(func(args mock.Arguments) {
				arg := args.Get(2).(*v1.FlyteRegistration)
				arg.Spec.WorkflowVersion = workflowVersion
				arg.Spec.WorkflowDomain = workflowDomain
				arg.Spec.WorkflowProject = workflowProject
				arg.Spec.WorkflowPackageURI = workflowPackageURI
			}).Return(nil).Once()

		mockDownloader.EXPECT().DownloadArtifact(mock.Anything, workflowPackageURI, workflowVersion).Return(artifactPath, nil).Once()

		mockFlyteAdminClient.EXPECT().RegisterWorkflow(mock.Anything, artifactPath, meta, flyteAuth).Return(nil).Once()

		// EXECUTION
		reconciler := &FlyteRegistrationReconciler{
			K8sClient:        mockK8sClient,
			Downloader:       mockDownloader,
			FlyteAdminClient: mockFlyteAdminClient,
			Config:           internal.Config{FlyteClientID: flyteClientID, FlyteAdminEndpoint: flyteAdminEndpoint},
		}

		_, err := reconciler.Reconcile(context.Background(), req)

		// ASSERTIONS
		assert.NoError(t, err)
	})

	t.Run("failure case: k8s client get method", func(t *testing.T) {
		// MOCK BEHAVIOUR
		mockK8sClient.EXPECT().Get(mock.Anything, req.NamespacedName, &testWorkflow).Once().
			Run(func(args mock.Arguments) {
				arg := args.Get(2).(*v1.FlyteRegistration)
				arg.Spec.WorkflowVersion = workflowVersion
				arg.Spec.WorkflowDomain = workflowDomain
				arg.Spec.WorkflowProject = workflowProject
				arg.Spec.WorkflowPackageURI = workflowPackageURI
			}).Return(errors.New("test error")).Once()

		// EXECUTION
		reconciler := &FlyteRegistrationReconciler{
			K8sClient:        mockK8sClient,
			Downloader:       mockDownloader,
			FlyteAdminClient: mockFlyteAdminClient,
		}

		_, err := reconciler.Reconcile(context.Background(), req)

		// ASSERTIONS
		assert.ErrorContains(t, err, "test error")
	})

	t.Run("failure case: downloader downloadArtifact method", func(t *testing.T) {
		// MOCK BEHAVIOUR
		mockK8sClient.EXPECT().Get(mock.Anything, req.NamespacedName, &testWorkflow).Once().
			Run(func(args mock.Arguments) {
				arg := args.Get(2).(*v1.FlyteRegistration)
				arg.Spec.WorkflowVersion = workflowVersion
				arg.Spec.WorkflowDomain = workflowDomain
				arg.Spec.WorkflowProject = workflowProject
				arg.Spec.WorkflowPackageURI = workflowPackageURI
			}).Return(nil).Once()

		mockDownloader.EXPECT().DownloadArtifact(mock.Anything, workflowPackageURI, workflowVersion).Return("", errors.New("test error")).Once()

		// EXECUTION
		reconciler := &FlyteRegistrationReconciler{
			K8sClient:        mockK8sClient,
			Downloader:       mockDownloader,
			FlyteAdminClient: mockFlyteAdminClient,
		}

		_, err := reconciler.Reconcile(context.Background(), req)

		// ASSERTIONS
		assert.ErrorContains(t, err, "failed to download artifact")
	})

	t.Run("failure case: flyte client registerWorkflow", func(t *testing.T) {
		// MOCK BEHAVIOUR
		mockK8sClient.EXPECT().Get(mock.Anything, req.NamespacedName, &testWorkflow).Once().
			Run(func(args mock.Arguments) {
				arg := args.Get(2).(*v1.FlyteRegistration)
				arg.Spec.WorkflowVersion = workflowVersion
				arg.Spec.WorkflowDomain = workflowDomain
				arg.Spec.WorkflowProject = workflowProject
				arg.Spec.WorkflowPackageURI = workflowPackageURI
			}).Return(nil).Once()

		mockDownloader.EXPECT().DownloadArtifact(mock.Anything, workflowPackageURI, workflowVersion).Return(artifactPath, nil).Once()

		mockFlyteAdminClient.EXPECT().RegisterWorkflow(mock.Anything, artifactPath, meta, flyteAuth).Return(errors.New("test error")).Once()

		// EXECUTION
		reconciler := &FlyteRegistrationReconciler{
			K8sClient:        mockK8sClient,
			Downloader:       mockDownloader,
			FlyteAdminClient: mockFlyteAdminClient,
			Config:           internal.Config{FlyteClientID: flyteClientID, FlyteAdminEndpoint: flyteAdminEndpoint},
		}

		_, err := reconciler.Reconcile(context.Background(), req)

		// ASSERTIONS
		assert.ErrorContains(t, err, "failed to register workflow")
	})
}
