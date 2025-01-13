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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FlyteRegistrationSpec defines the desired state of FlyteRegistration
type FlyteRegistrationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// WorkflowDomain is the domain of the workflow - we can have multiple domains on one flyte.backend cluster
	WorkflowDomain string `json:"workflowDomain"`

	// WorkflowProject is the project of the workflow - we can have multiple projects on one flyte.backend cluster
	WorkflowProject string `json:"workflowProject"`

	// WorkflowPackageURI is the URI of the workflow artifact packaged by pyflyte in CI and stored in s3 storage
	WorkflowPackageURI string `json:"workflowPackageUri"`

	// WorkflowVersion is the version of the workflow
	WorkflowVersion string `json:"workflowVersion"`
}

// FlyteRegistrationStatus defines the observed state of FlyteRegistration
type FlyteRegistrationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// WorkflowDomain is the domain of the workflow - we can have multiple domains on one flyte.backend cluster
	WorkflowDomain string `json:"workflowDomain"`

	// WorkflowProject is the project of the workflow - we can have multiple projects on one flyte.backend cluster
	WorkflowProject string `json:"workflowProject"`

	// WorkflowPackageURI is the URI of the workflow artifact packaged by pyflyte in CI and stored in s3 storage
	WorkflowPackageURI string `json:"workflowPackageUri"`

	// WorkflowVersion is the version of the workflow
	WorkflowVersion string `json:"workflowVersion"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FlyteRegistration is the Schema for the flyteregistrations API
type FlyteRegistration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FlyteRegistrationSpec   `json:"spec,omitempty"`
	Status FlyteRegistrationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FlyteRegistrationList contains a list of FlyteRegistration
type FlyteRegistrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FlyteRegistration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FlyteRegistration{}, &FlyteRegistrationList{})
}
