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
	postgresqlChartName    string = "postgresql"
	postgresqlChartVersion string = "10.3.13"
)

var postgresqlValues map[string]interface{} = map[string]interface{}{
	"postgresqlPassword": "123456",
	"postgresqlDatabase": "cloudship",
	"postgresqlUsername": "cloudship",
}

type postgresqlAction struct{}

func (e postgresqlAction) PreInstalacion() map[string]interface{} {
	fmt.Print("Soy la estrategia de postgres")
	return postgresqlValues
}

func (e postgresqlAction) EnvVars(as *cloudshipv1alpha1.Application) []corev1.EnvVar {
	var dbURL = fmt.Sprintf("db-postgresql-headless.%s.svc.cluster.local", as.GetName())
	return []corev1.EnvVar{
		{
			Name:  "DATABASE_NAME",
			Value: "cloudship",
		},
		{
			Name:  "DATABASE_HOST",
			Value: dbURL,
		},
		{
			Name:  "DATABASE_PORT",
			Value: "5432",
		},
		{
			Name:  "DATABASE_USERNAME",
			Value: "cloudship",
		},
	}
}

// NewPostgreSQLManagerFactory returns a new Helm manager factory capable of installing and uninstalling PostgreSQL releases.
func NewPostgreSQLManagerFactory(mgr crmanager.Manager) ManagerFactory {
	return &managerFactory{
		mgr:          mgr,
		chartName:    postgresqlChartName,
		chartVersion: postgresqlChartVersion,
		values:       postgresqlValues,
		releaseName:  "db",
		settings:     cli.New(),
		action:       postgresqlAction{},
	}
}
