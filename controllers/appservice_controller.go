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

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	cloudshipv1alpha1 "github.com/ToucanSoftware/cloudship-operator/api/v1alpha1"
	"github.com/ToucanSoftware/cloudship-operator/pkg/types"
)

// AppServiceReconciler reconciles a AppService object
type AppServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var (
	// ReconcileWaitResult is the time to wait between reconciliation.
	ReconcileWaitResult = reconcile.Result{RequeueAfter: 30 * time.Second}
)

// Reconcile reconciles a AppService object
// +kubebuilder:rbac:groups=cloudship.toucansoft.io,resources=appservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cloudship.toucansoft.io,resources=appservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
func (r *AppServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("appservice", req.NamespacedName)
	log.Info("Reconcile container workload")

	var appService cloudshipv1alpha1.AppService
	if err := r.Get(ctx, req.NamespacedName, &appService); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Container workload is deleted")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deploy, err := r.renderDeployment(ctx, &appService)
	if err != nil {
		log.Error(err, "Failed to render a deployment")
		// r.record.Event(eventObj, event.Warning(errRenderWorkload, err))
		return ReconcileWaitResult, client.IgnoreNotFound(err)
	}
	// server side apply, only the fields we set are touched
	applyOpts := []client.PatchOption{client.ForceOwnership, client.FieldOwner(appService.GetUID())}
	if err := r.Patch(ctx, deploy, client.Apply, applyOpts...); err != nil {
		log.Error(err, "Failed to apply to a deployment")
		//r.record.Event(eventObj, event.Warning(errApplyDeployment, err))
		return ReconcileWaitResult, client.IgnoreNotFound(err)
	}

	service, err := r.renderService(ctx, &appService, deploy)
	if err != nil {
		log.Error(err, "Failed to render a service")
		//r.record.Event(eventObj, event.Warning(errRenderService, err))
		return ReconcileWaitResult, client.IgnoreNotFound(err)
	}
	// server side apply the service
	if err := r.Patch(ctx, service, client.Apply, applyOpts...); err != nil {
		log.Error(err, "Failed to apply a service")
		//r.record.Event(eventObj, event.Warning(errApplyDeployment, err))
		return ReconcileWaitResult, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudshipv1alpha1.AppService{}).
		Complete(r)
}

// create a corresponding deployment
func (r *AppServiceReconciler) renderDeployment(ctx context.Context,
	appService *cloudshipv1alpha1.AppService) (*appsv1.Deployment, error) {

	resources, err := TranslateContainer(ctx, appService)
	if err != nil {
		return nil, err
	}

	deploy, ok := resources[0].(*appsv1.Deployment)
	if !ok {
		return nil, fmt.Errorf("internal error, deployment is not rendered correctly")
	}
	// make sure we don't have opinion on the replica count
	deploy.Spec.Replicas = nil
	// k8s server-side patch complains if the protocol is not set
	for i := 0; i < len(deploy.Spec.Template.Spec.Containers); i++ {
		for j := 0; j < len(deploy.Spec.Template.Spec.Containers[i].Ports); j++ {
			if len(deploy.Spec.Template.Spec.Containers[i].Ports[j].Protocol) == 0 {
				deploy.Spec.Template.Spec.Containers[i].Ports[j].Protocol = corev1.ProtocolTCP
			}
		}
	}
	r.Log.Info("rendered a deployment", "deploy", deploy.Spec.Template.Spec)

	// set the controller reference so that we can watch this deployment and it will be deleted automatically
	if err := ctrl.SetControllerReference(appService, deploy, r.Scheme); err != nil {
		return nil, err
	}

	return deploy, nil
}

// create a corresponding service
func (r *AppServiceReconciler) renderService(ctx context.Context,
	appService *cloudshipv1alpha1.AppService, deploy *appsv1.Deployment) (*corev1.Service, error) {
	// create a service for the workload
	resources, err := ServiceInjector(ctx, appService, []types.Object{deploy})
	if err != nil {
		return nil, err
	}
	service, ok := resources[1].(*corev1.Service)
	if !ok {
		return nil, fmt.Errorf("internal error, service is not rendered correctly")
	}
	// the service injector lib doesn't set the namespace and serviceType
	service.Namespace = appService.Namespace
	service.Spec.Type = corev1.ServiceTypeClusterIP
	// k8s server-side patch complains if the protocol is not set
	for i := 0; i < len(service.Spec.Ports); i++ {
		service.Spec.Ports[i].Protocol = corev1.ProtocolTCP
	}
	// always set the controller reference so that we can watch this service and
	if err := ctrl.SetControllerReference(appService, service, r.Scheme); err != nil {
		return nil, err
	}
	return service, nil
}
