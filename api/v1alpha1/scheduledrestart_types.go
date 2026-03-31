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

// ScheduledRestartSpec defines the desired state of ScheduledRestart
type ScheduledRestartSpec struct {
	Selector      metav1.LabelSelector `json:"matchLabels,omitempty"`
	Schedule      string               `json:"schedule"`
	Suspend       bool                 `json:"suspend"`
	RestartPolicy string               `json:"restartPolicy"`
}

// ScheduledRestartStatus defines the observed state of ScheduledRestart.
type ScheduledRestartStatus struct {
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
