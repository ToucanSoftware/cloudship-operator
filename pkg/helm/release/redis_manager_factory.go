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
	redisRepositoryURL  string = "https://charts.bitnami.com/bitnami"
	redisRepositoryName string = "bitnami"
	redisChartName      string = "redis"
	redisChartVersion   string = "12.8.3"
)

var redisValues map[string]interface{} = map[string]interface{}{}

type EstrategiaRedis struct{}

func (e EstrategiaRedis) PreInstalacion() map[string]interface{} {
	fmt.Print("Soy la estrategia de redis")
	return redisValues
}

// NewRedisManagerFactory returns a new Helm manager factory capable of installing and uninstalling Redis releases.
func NewRedisManagerFactory(mgr crmanager.Manager) ManagerFactory {
	return &managerFactory{
		mgr:          mgr,
		chartName:    redisChartName,
		chartVersion: redisChartVersion,
		//values:       redisValues,
		releaseName: "cache-redis",
		settings:    cli.New(),
		estrategia:  EstrategiaRedis{},
	}
}
