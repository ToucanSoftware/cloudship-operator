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
	postgresqlChartName    string = "postgresql"
	postgresqlChartVersion string = "10.3.13"
)

var postgresqlValues map[string]interface{} = map[string]interface{}{
	"postgresqlPassword": "123456",
	"postgresqlDatabase": "cloudship",
	"postgresqlUsername": "cloudship",
}

type EstrategiaPostgresql struct{}

func (e EstrategiaPostgresql) PreInstalacion() map[string]interface{} {
	fmt.Print("Soy la estrategia de postgres")
	return postgresqlValues
}

// NewPostgreSQLManagerFactory returns a new Helm manager factory capable of installing and uninstalling PostgreSQL releases.
func NewPostgreSQLManagerFactory(mgr crmanager.Manager) ManagerFactory {
	return &managerFactory{
		mgr:          mgr,
		chartName:    postgresqlChartName,
		chartVersion: postgresqlChartVersion,
		//values:       postgresqlValues,
		releaseName: "db",
		settings:    cli.New(),
		estrategia:  EstrategiaPostgresql{},
	}
}
