/*
Copyright 2025 jwwangari6.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied.
See the License for the specific language governing permissions
and limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Target defines what resource should be scaled
type Target struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	ScaleDown int32  `json:"scaleDown"`
	ScaleUp   int32  `json:"scaleUp"`
}

// Schedule defines the scaling schedule
type Schedule struct {
	ScaleDown string   `json:"scaleDown"`
	ScaleUp   string   `json:"scaleUp"`
	Days      []string `json:"days"`
}

// ScheduleScalerSpec defines the desired state of ScheduleScaler
type ScheduleScalerSpec struct {
	Namespace string   `json:"namespace"`
	Targets   []Target `json:"targets"`
	Schedule  Schedule `json:"schedule"`
}

// ScheduleScalerStatus defines the observed state of ScheduleScaler
type ScheduleScalerStatus struct {
	LastAction string `json:"lastAction,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ScheduleScaler is the Schema for the schedulescalers API
type ScheduleScaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScheduleScalerSpec   `json:"spec"`
	Status ScheduleScalerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ScheduleScalerList contains a list of ScheduleScaler
type ScheduleScalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScheduleScaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ScheduleScaler{}, &ScheduleScalerList{})
}
