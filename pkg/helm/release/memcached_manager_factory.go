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
	"helm.sh/helm/v3/pkg/cli"

	crmanager "sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	memcachedChartName    string = "memcached"
	memcachedChartVersion string = "5.8.0"
)

var memcachedValues map[string]interface{} = map[string]interface{}{}

type EstrategiaMemcached struct{}

func (e EstrategiaMemcached) PreInstalacion() map[string]interface{} {
	fmt.Print("Soy la estrategia de memcached")
	return memcachedValues
}

// NewMemecachedManagerFactory returns a new Helm manager factory capable of installing and uninstalling Memcached releases.
func NewMemecachedManagerFactory(mgr crmanager.Manager) ManagerFactory {
	return &managerFactory{
		mgr:          mgr,
		chartName:    memcachedChartName,
		chartVersion: memcachedChartVersion,
		//values:       memcachedValues,
		releaseName: "cache-memcached",
		settings:    cli.New(),
		estrategia:  EstrategiaMemcached{},
	}
}
