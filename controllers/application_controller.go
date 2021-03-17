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

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cloudshipv1alpha1 "github.com/ToucanSoftware/cloudship-operator/api/v1alpha1"
	"github.com/ToucanSoftware/cloudship-operator/pkg/helm/release"
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Log            logr.Logger
	Scheme         *runtime.Scheme
	ManagerFactory release.ManagerFactory
}

// +kubebuilder:rbac:groups=cloudship.toucansoft.io,resources=applications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cloudship.toucansoft.io,resources=applications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cloudship.toucansoft.io,resources=applications/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error
	log := r.Log.WithValues("application", req.NamespacedName)
	log.Info("Reconcilate Applicatio")

	var app cloudshipv1alpha1.Application

	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Application is deleted")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Reder Namespace base on the application name
	namespace := r.renderNamespace(&app)
	// Set application as the owner and controller
	ctrl.SetControllerReference(&app, namespace, r.Scheme)
	// server side apply, only the fields we set are touched
	applyOpts := []client.PatchOption{client.ForceOwnership, client.FieldOwner(app.GetUID())}
	if err := r.Patch(ctx, namespace, client.Apply, applyOpts...); err != nil {
		log.Error(err, "Failed to apply to a Namespace")
		//r.record.Event(eventObj, event.Warning(errApplyDeployment, err))
		return ReconcileWaitResult, client.IgnoreNotFound(err)
	}

	err = r.processCache(log, &app)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	err = r.processEventStream(log, &app)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudshipv1alpha1.Application{}).
		Complete(r)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) renderNamespace(app *cloudshipv1alpha1.Application) *corev1.Namespace {
	return &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       namespaceKind,
			APIVersion: namespaceAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: app.GetName(),
		},
	}
}

func (r *ApplicationReconciler) processCache(log logr.Logger, app *cloudshipv1alpha1.Application) error {
	if app.Spec.CacheRef == nil {
		log.Info(fmt.Sprintf("No cache for application %s", app.GetName()))
		return nil
	}
	log.Info(fmt.Sprintf("Processing cache for application %s", app.GetName()))
	return nil
}

func (r *ApplicationReconciler) processEventStream(log logr.Logger, app *cloudshipv1alpha1.Application) error {
	if app.Spec.EventStreamRefs == nil {
		log.Info(fmt.Sprintf("No event stream for application %s", app.GetName()))
		return nil
	}
	log.Info(fmt.Sprintf("Processing event stream for application %s", app.GetName()))
	return nil
}
