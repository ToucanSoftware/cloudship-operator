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

package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	cloudshipv1alpha1 "github.com/ToucanSoftware/cloudship-operator/api/v1alpha1"
)

// An Object is a Kubernetes object.
type Object interface {
	metav1.Object
	runtime.Object
}

// StatusFor safely returns a typed status block from an application.
func StatusFor(app *cloudshipv1alpha1.Application) *cloudshipv1alpha1.ApplicationStatus {
	return &cloudshipv1alpha1.ApplicationStatus{}
}
