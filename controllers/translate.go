/*
Copyright © 2020 ToucanSoftware

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

package controllers

import (
	"context"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cloudshipv1alpha1 "github.com/ToucanSoftware/cloudship-operator/api/v1alpha1"
	"github.com/ToucanSoftware/cloudship-operator/pkg/types"
)

var (
	deploymentKind       = reflect.TypeOf(appsv1.Deployment{}).Name()
	deploymentAPIVersion = appsv1.SchemeGroupVersion.String()
	serviceKind          = reflect.TypeOf(corev1.Service{}).Name()
	serviceAPIVersion    = corev1.SchemeGroupVersion.String()
)

// Reconcile error strings.
const (
	labelKey = "appservice.toucansoft.io"

	errNotContainerizedWorkload = "object is not a containerized workload"
)

//TranslateContainer transalate to kubernetes objects
func TranslateContainer(ctx context.Context, as *cloudshipv1alpha1.AppService) ([]types.Object, error) {

	d := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       deploymentKind,
			APIVersion: deploymentAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      as.GetName(),
			Namespace: as.GetNamespace(),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					labelKey: string(as.GetUID()),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						labelKey: string(as.GetUID()),
					},
				},
			},
		},
	}

	// pass through label and annotation from the workload to the deployment
	//util.PassLabelAndAnnotation(w, d)
	// pass through label from the workload to the pod template
	//util.PassLabel(w, &d.Spec.Template)

	return []types.Object{d}, nil
}
