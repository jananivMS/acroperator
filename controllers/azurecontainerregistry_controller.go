/*

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

	azurev1 "github.com/jananiv/acroperator/api/v1"
	acrhelper "github.com/jananiv/acroperator/resourcehelper/acr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// AzureContainerRegistryReconciler reconciles a AzureContainerRegistry object
type AzureContainerRegistryReconciler struct {
	client.Client
	Log logr.Logger
}

const acrFinalizerName = "acr.finalizers.com"

func ignoreNotFound(err error) error {

	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}

// +kubebuilder:rbac:groups=azure.microsoft.com,resources=azurecontainerregistries,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=azure.microsoft.com,resources=azurecontainerregistries/status,verbs=get;update;patch

// Reconcile does this
func (r *AzureContainerRegistryReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("azurecontainerregistry", req.NamespacedName)

	// your logic here
	var instance azurev1.AzureContainerRegistry

	if err := r.Get(ctx, req.NamespacedName, &instance); err != nil {
		log.Error(err, "unable to fetch Eventhub")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, ignoreNotFound(err)
	}

	if instance.IsBeingDeleted() {
		err := r.handleFinalizer(&instance)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("error when handling finalizer: %v", err)
		}
		return ctrl.Result{}, nil
	}

	if !instance.HasFinalizer(acrFinalizerName) {
		err := r.addFinalizer(&instance)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("error when removing finalizer: %v", err)
		}
		return ctrl.Result{}, nil
	}

	if !instance.IsSubmitted() {
		err := r.createAzureContainerRegistry(&instance)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("error when creating resource in azure: %v", err)
		}
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager does this
func (r *AzureContainerRegistryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&azurev1.AzureContainerRegistry{}).
		Complete(r)
}

func (r *AzureContainerRegistryReconciler) createAzureContainerRegistry(instance *azurev1.AzureContainerRegistry) error {
	ctx := context.Background()

	var err error

	acrName := instance.ObjectMeta.Name
	resourcegroup := instance.Spec.ResourceGroup
	location := instance.Spec.Location
	adminuserenabled := instance.Spec.AdminUserEnabled
	sku := instance.Spec.Sku

	// write information back to instance
	instance.Status.Provisioning = true
	err = r.Update(ctx, instance)
	if err != nil {
		//log error and kill it
		fmt.Println("Unable to update instance")
	}
	_, err = acrhelper.CreateRegistry(ctx, resourcegroup, acrName, location, sku, adminuserenabled)
	if err != nil {
		fmt.Println("Couldn't create resource in azure")
		return err
	}
	// write information back to instance
	instance.Status.Provisioning = false
	instance.Status.Provisioned = true

	err = r.Update(ctx, instance)
	if err != nil {
		//log error and kill it
		fmt.Println("Unable to update instance")
	}
	return nil
}

func (r *AzureContainerRegistryReconciler) deleteAzureContainerRegistry(instance *azurev1.AzureContainerRegistry) error {

	ctx := context.Background()

	acrName := instance.ObjectMeta.Name
	resourcegroup := instance.Spec.ResourceGroup

	var err error
	_, err = acrhelper.DeleteRegistry(ctx, resourcegroup, acrName)
	if err != nil {
		fmt.Println("Couldn't delete resouce in azure")
		return err
	}
	return nil
}

func (r *AzureContainerRegistryReconciler) addFinalizer(instance *azurev1.AzureContainerRegistry) error {
	instance.AddFinalizer(acrFinalizerName)
	err := r.Update(context.Background(), instance)
	if err != nil {
		return fmt.Errorf("failed to update finalizer: %v", err)
	}
	return nil
}

func (r *AzureContainerRegistryReconciler) handleFinalizer(instance *azurev1.AzureContainerRegistry) error {
	if instance.HasFinalizer(acrFinalizerName) {
		// our finalizer is present, so lets handle our external dependency
		if err := r.deleteExternalDependency(instance); err != nil {
			return err
		}

		instance.RemoveFinalizer(acrFinalizerName)
		if err := r.Update(context.Background(), instance); err != nil {
			return err
		}
	}
	// Our finalizer has finished, so the reconciler can do nothing.
	return nil
}

func (r *AzureContainerRegistryReconciler) deleteExternalDependency(instance *azurev1.AzureContainerRegistry) error {

	return r.deleteAzureContainerRegistry(instance)
}
