/*
Copyright 2021 ToucanSoftware.

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

// AppResourceSpec defines the desired state of AppResource
type AppResourceSpec struct {
}

// AppResourceStatus defines the observed state of AppResource
type AppResourceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +genclient
// +kubebuilder:resource:path=resources,scope=Namespaced,singular=resource,shortName=csr,categories=cloudship

// AppResource is the Schema for the application resources API
type AppResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppResourceSpec   `json:"spec,omitempty"`
	Status AppResourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AppResourceList contains a list of AppResource
type AppResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AppResource{}, &AppResourceList{})
}
