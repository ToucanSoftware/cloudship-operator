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
	cloudshipv1alpha1 "github.com/ToucanSoftware/cloudship-operator/api/v1alpha1"
	"helm.sh/helm/v3/pkg/cli"
	corev1 "k8s.io/api/core/v1"

	crmanager "sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	mysqlChartName    string = "mysql"
	mysqlChartVersion string = "8.5.1"
)

var mysqlValues map[string]interface{} = map[string]interface{}{}

type mysqlAction struct{}

func (e mysqlAction) PreInstalacion() map[string]interface{} {
	return mysqlValues
}

func (e mysqlAction) EnvVars(as *cloudshipv1alpha1.Application) []corev1.EnvVar {
	return nil
}

func (e mysqlAction) Port() string {
	return ""
}

func (e mysqlAction) Hostname(as *cloudshipv1alpha1.Application) string {
	return ""
}

// NewMySQLManagerFactory returns a new Helm manager factory capable of installing and uninstalling MySQL releases.
func NewMySQLManagerFactory(mgr crmanager.Manager) ManagerFactory {
	return &managerFactory{
		mgr:          mgr,
		chartName:    mysqlChartName,
		chartVersion: mysqlChartVersion,
		values:       mysqlValues,
		releaseName:  "db-mysql",
		settings:     cli.New(),
		action:       mysqlAction{},
	}
}
