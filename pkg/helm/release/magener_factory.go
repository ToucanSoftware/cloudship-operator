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
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/kube"
	helmrelease "helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"

	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/helm/pkg/strvals"

	ctrl "sigs.k8s.io/controller-runtime"
	crmanager "sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/ToucanSoftware/cloudship-operator/internal/helm/client"
)

// ManagerFactory creates Managers that are specific to custom resources. It is
// used by the HelmOperatorReconciler during resource reconciliation, and it
// improves decoupling between reconciliation logic and the Helm backend
// components used to manage releases.
type ManagerFactory interface {
	NewManager(namespace string, overrideValues map[string]string) (Manager, error)
}

type managerFactory struct {
	mgr          crmanager.Manager
	chartName    string
	chartVersion string
	releaseName  string
	values       map[string]interface{}
	settings     *cli.EnvSettings
	action       ManagerAction
}

const (
	// This is the directory in the docker image where all the Helm Charts will be located
	defaultChartPathPrefix string = "/charts"
)

func (f managerFactory) NewManager(namespace string, overrideValues map[string]string) (Manager, error) {
	var log = ctrl.Log.WithName("helm").WithName("manager_factory")

	// Get both v2 and v3 storage backends
	clientv1, err := v1.NewForConfig(f.mgr.GetConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to get core/v1 client: %w", err)
	}
	storageBackend := storage.Init(driver.NewSecrets(clientv1.Secrets(namespace)))

	// Get the necessary clients and client getters. Use a client that injects the CR
	// as an owner reference into all resources templated by the chart.
	rcg, err := client.NewRESTClientGetter(f.mgr, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get REST client getter from manager: %w", err)
	}

	kubeClient := kube.New(rcg)
	restMapper := f.mgr.GetRESTMapper()
	ownerRefClient, err := client.NewOwnerRefInjectingClient(*kubeClient, restMapper, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to inject owner references: %w", err)
	}
	expOverrides, err := parseOverrides(overrideValues)
	if err != nil {
		return nil, fmt.Errorf("failed to parse override values: %w", err)
	}
	values := mergeMaps(f.values, expOverrides)
	if err != nil {
		return nil, fmt.Errorf("failed to parse override values: %w", err)
	}
	actionConfig := &action.Configuration{
		RESTClientGetter: rcg,
		Releases:         storageBackend,
		KubeClient:       ownerRefClient,
		Log:              func(_ string, _ ...interface{}) {},
	}

	var cp = f.buildChartPath()

	log.Info(fmt.Sprintf("Loading Chart: %s", cp))
	crChart, err := loader.Load(cp)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart dir: %w", err)
	}
	releaseName, err := getReleaseName(storageBackend, crChart.Name(), f.releaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to get helm release name: %w", err)
	}
	return &manager{
		actionConfig:   actionConfig,
		storageBackend: storageBackend,
		kubeClient:     ownerRefClient,

		releaseName: releaseName,
		namespace:   namespace,

		chart:  crChart,
		values: values,
		action: f.action,
		//status: types.StatusFor(cr),
	}, nil
}

// getReleaseName returns a release name for the CR.
//
// getReleaseName searches for a release using the CR name. If a release
// cannot be found, or if it is found and was created by the chart managed
// by this manager, the CR name is returned.
//
// If a release is found but it was created by another chart, that means we
// have a release name collision, so return an error. This case is possible
// because Kubernetes allows instances of different types to have the same name
// in the same namespace.
//
// TODO(jlanford): As noted above, using the CR name as the release name raises
//   the possibility of collision. We should move this logic to a validating
//   admission webhook so that the CR owner receives immediate feedback of the
//   collision. As is, the only indication of collision will be in the CR status
//   and operator logs.
func getReleaseName(storageBackend *storage.Storage, crChartName string,
	extectedReleaseName string) (string, error) {

	releaseName := extectedReleaseName
	history, exists, err := releaseHistory(storageBackend, releaseName)
	if err != nil {
		return "", err
	}
	if !exists {
		return releaseName, nil
	}

	// If a release name with the CR name exists, but the release's chart is
	// different than the chart managed by this operator, return an error
	// because something else created the existing release.
	if history[0].Chart == nil {
		return "", fmt.Errorf("could not find chart metadata in release with name %q", releaseName)
	}
	existingChartName := history[0].Chart.Name()
	if existingChartName != crChartName {
		return "", fmt.Errorf("duplicate release name: found existing release with name %q for chart %q",
			releaseName, existingChartName)
	}

	return releaseName, nil
}

func releaseHistory(storageBackend *storage.Storage, releaseName string) ([]*helmrelease.Release, bool, error) {
	releaseHistory, err := storageBackend.History(releaseName)
	if err != nil {
		if notFoundErr(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return releaseHistory, len(releaseHistory) > 0, nil
}

func parseOverrides(in map[string]string) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	for k, v := range in {
		val := fmt.Sprintf("%s=%s", k, v)
		if err := strvals.ParseIntoString(val, out); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

func (f managerFactory) buildChartPath() string {
	// using a chart path prefix we can use it for debuging
	var chartPathPrefix = os.Getenv("CHART_PATH_PREFIX")
	var chartFileName = fmt.Sprintf("%s-%s.tgz", f.chartName, f.chartVersion)
	if chartPathPrefix == "" {
		chartPathPrefix = defaultChartPathPrefix
	}
	return fmt.Sprintf("%s/%s", chartPathPrefix, chartFileName)
}
