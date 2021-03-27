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

// DatabaseType are the types of database supported
// +kubebuilder:validation:Enum=MySQL;PostgreSQL
type DatabaseType string

const (
	// DatabaseTypeMySQL use MySQL for Database
	DatabaseTypeMySQL DatabaseType = "MySQL"
	// DatabaseTypePostgreSQL use PostgreSQL for Database
	DatabaseTypePostgreSQL DatabaseType = "PostgreSQL"
)

// DatabaseSpec is the definition for Database support for the service
type DatabaseSpec struct {
	// Type is the type of the database
	Type DatabaseType `json:"type,omitempty"`
}

// Service defines an Application Service
type Service struct {
	// Name of this service. Must be unique within its service.
	Name string `json:"name"`
	// Port is the number of the port
	Port int32 `json:"portNumber"`
}

// Container defines a OCI container
type Container struct {
	// Name of this container. Must be unique within its service.
	Name string `json:"name"`

	// Image this container should run. Must be a path-like or URI-like
	// representation of an OCI image. May be prefixed with a registry address
	// and should be suffixed with a tag.
	Image string `json:"image"`

	// Ports are the ports that this container exposes
	Ports []Service `json:"ports"`
}

// AppServiceSpec defines the desired state of AppService
type AppServiceSpec struct {
	// Containers of which this service consists.
	Containers []Container `json:"containers"`

	// DatabaseRef is the reference to database for the service
	// +optional
	DatabaseRef *DatabaseSpec `json:"databaseRef,omitempty"`
}

// AppServiceStatus defines the observed state of AppService
type AppServiceStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +genclient
// +kubebuilder:resource:path=services,scope=Namespaced,singular=service,shortName=css,categories=cloudship

// AppService is the Schema for the application services API
type AppService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppServiceSpec   `json:"spec,omitempty"`
	Status AppServiceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AppServiceList contains a list of AppService
type AppServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AppService{}, &AppServiceList{})
}
