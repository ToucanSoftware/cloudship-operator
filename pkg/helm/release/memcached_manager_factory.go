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

package release

import (
	"fmt"

	cloudshipv1alpha1 "github.com/ToucanSoftware/cloudship-operator/api/v1alpha1"
	"helm.sh/helm/v3/pkg/cli"
	corev1 "k8s.io/api/core/v1"
	crmanager "sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	memcachedChartName    string = "memcached"
	memcachedChartVersion string = "5.8.0"
)

var memcachedValues map[string]interface{} = map[string]interface{}{}

type memcachedActions struct{}

func (e memcachedActions) PreInstalacion() map[string]interface{} {
	fmt.Print("Soy la estrategia de memcached")
	return memcachedValues
}

func (e memcachedActions) EnvVars(as *cloudshipv1alpha1.Application) []corev1.EnvVar {
	return nil
}

// NewMemecachedManagerFactory returns a new Helm manager factory capable of installing and uninstalling Memcached releases.
func NewMemecachedManagerFactory(mgr crmanager.Manager) ManagerFactory {
	return &managerFactory{
		mgr:          mgr,
		chartName:    memcachedChartName,
		chartVersion: memcachedChartVersion,
		values:       memcachedValues,
		releaseName:  "cache-memcached",
		settings:     cli.New(),
		action:       memcachedActions{},
	}
}
