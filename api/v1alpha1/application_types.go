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

// EventStreamType are the types of event stream supported
type EventStreamType string

const (
	// EventStreamTypeKafka use Kafka for Event Stream
	EventStreamTypeKafka EventStreamType = "Kafka"
	// EventStreamTypeRabbitMQ use RabbitMQ for Event Stream
	EventStreamTypeRabbitMQ EventStreamType = "RabbitMQ"
)

// EventStreamSpec is the definition for Event Stream support for and applicaction
type EventStreamSpec struct {
	// Type is the type of the cache
	Type EventStreamType `json:"type,omitempty"`
}

// CacheType are the types of cache supported
type CacheType string

const (
	// CacheTypeReddis use Reddis for Cache
	CacheTypeReddis CacheType = "Reddis"
	// CacheTypeMemcached use Memcached for Cache
	CacheTypeMemcached CacheType = "Memcached"
)

// CacheSpec is the definition for Cache support for and applicaction
type CacheSpec struct {
	// Type is the type of the cache
	Type CacheType `json:"type,omitempty"`
}

// ApplicationSpec defines the desired state of Application
type ApplicationSpec struct {
	// Description is the name of the application
	Description string `json:"description,omitempty"`

	// CacheRef is the reference to cache informacion for the applicacion
	// +optional
	CacheRef *CacheSpec `json:"cacheRef,omitempty"`

	// EventStreamRefs is the reference to event stream informacion for the applicacion
	// +optional
	EventStreamRefs *EventStreamSpec `json:"eventStreamRef,omitempty"`
}

// ApplicationStatus defines the observed state of Application
type ApplicationStatus struct {
	// Deployment is the status of the deployment of the application
	Deployment string `json:"description,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Application is the Schema for the applications API
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ApplicationList contains a list of Application
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
