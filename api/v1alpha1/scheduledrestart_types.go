/*
Copyright 2026.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ScheduledRestartSpec defines the desired state of ScheduledRestart
type ScheduledRestartSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html
	Selector      metav1.LabelSelector `json:"matchLabels,omitempty"`
	Schedule      string               `json:"schedule"`
	Suspend       bool                 `json:"suspend"`
	RestartPolicy string               `json:"restartPolicy"`
}

// ScheduledRestartStatus defines the observed state of ScheduledRestart.
type ScheduledRestartStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	LastRestartTime    *metav1.Time `json:"lastRestartTime,omitempty"`
	NextRestartTime    *metav1.Time `json:"nextRestartTime,omitempty"`
	ObservedGeneration int64        `json:"observedGeneration"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ScheduledRestart is the Schema for the scheduledrestarts API
type ScheduledRestart struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of ScheduledRestart
	// +required
	Spec ScheduledRestartSpec `json:"spec"`

	// status defines the observed state of ScheduledRestart
	// +optional
	Status ScheduledRestartStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// ScheduledRestartList contains a list of ScheduledRestart
type ScheduledRestartList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []ScheduledRestart `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ScheduledRestart{}, &ScheduledRestartList{})
}
