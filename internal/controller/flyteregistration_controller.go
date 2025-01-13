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

// Package controller holds the logic for interacting with the K8s client & reconciling K8s resources
package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/adarga-ai/flyte-workflow-registration-operator/api/v1"
	"github.com/adarga-ai/flyte-workflow-registration-operator/internal"
	"github.com/adarga-ai/flyte-workflow-registration-operator/internal/command"
	"github.com/adarga-ai/flyte-workflow-registration-operator/internal/downloader"
	"github.com/adarga-ai/flyte-workflow-registration-operator/internal/flyte"
)

// K8sClient is the interface the K8s client used for mocking
//
//go:generate mockery --name=K8sClient
type K8sClient interface {
	Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error
}

// FlyteRegistrationReconciler reconciles a FlyteRegistration object
type FlyteRegistrationReconciler struct {
	K8sClient K8sClient
	Scheme    *runtime.Scheme
	// A configuration for the controller
	Config internal.Config
	// A downloader for downloading artifacts.  This is configured
	Downloader downloader.Client
	// A client for interacting with the flyte.backend cluster
	FlyteAdminClient flyte.Client
}

// NewFlyteRegistrationReconciler sets up the required dependencies for the controller
// and returns a new FlyteRegistrationReconciler instance
func NewFlyteRegistrationReconciler(config internal.Config, k8sClient client.Client, scheme *runtime.Scheme) (*FlyteRegistrationReconciler, error) {
	ctx := context.Background()

	// Set up the artifact Client
	d, err := downloader.NewClient(ctx, config)
	if err != nil {
		return nil, err
	}

	// Set up the Flyte Client
	fClient := flyte.NewClient(command.NewOSCommandExecutor())

	return &FlyteRegistrationReconciler{
		K8sClient:        k8sClient,
		Scheme:           scheme,
		Downloader:       d,
		FlyteAdminClient: fClient,
		Config:           config,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FlyteRegistrationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.FlyteRegistration{}).
		Complete(r)
}

//+kubebuilder:rbac:groups=flyte.backend,resources=flyteregistrations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flyte.backend,resources=flyteregistrations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=flyte.backend,resources=flyteregistrations/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the FlyteRegistration object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *FlyteRegistrationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Get the FlyteRegistration object
	var flyteWorkflow v1.FlyteRegistration
	if err := r.K8sClient.Get(ctx, req.NamespacedName, &flyteWorkflow); err != nil {
		log.Log.Info("unable to fetch any workflows")
		// There are no workflows to reconcile
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	workflowVersion := flyteWorkflow.Spec.WorkflowVersion
	workflowDomain := flyteWorkflow.Spec.WorkflowDomain
	workflowProject := flyteWorkflow.Spec.WorkflowProject
	workflowPackageURI := flyteWorkflow.Spec.WorkflowPackageURI

	meta := flyte.WorkflowMetadata{
		WorkflowVersion: workflowVersion,
		Domain:          workflowDomain,
		Project:         workflowProject,
	}

	flyteAuth := flyte.Auth{
		AdminEndpoint:      r.Config.FlyteAdminEndpoint,
		ClientID:           r.Config.FlyteClientID,
		ClientSecretEnvVar: "FLYTE_CLIENT_SECRET",
	}

	fullArtifactPath, err := r.Downloader.DownloadArtifact(ctx, workflowPackageURI, workflowVersion)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to download artifact: %w", err)
	}

	log.Log.Info("downloaded artifact", "path", fullArtifactPath)

	if err := r.FlyteAdminClient.RegisterWorkflow(ctx, fullArtifactPath, meta, flyteAuth); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to register workflow %w", err)
	}

	log.Log.Info("successfully registered workflow", "name", req.Name, "version", workflowVersion, "domain", workflowDomain, "project", workflowProject)
	return ctrl.Result{}, nil
}
