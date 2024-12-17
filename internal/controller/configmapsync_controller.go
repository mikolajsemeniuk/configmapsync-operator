/*
Copyright 2024.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	syncv1alpha1 "github.com/mikolajsemeniuk/configmapsync-operator/api/v1alpha1"
)

// ConfigMapSyncReconciler reconciles a ConfigMapSync object
type ConfigMapSyncReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=sync.example.com,resources=configmapsyncs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=sync.example.com,resources=configmapsyncs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=sync.example.com,resources=configmapsyncs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConfigMapSync object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *ConfigMapSyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	// 1. Fetch the ConfigMapSync CR
	var cms syncv1alpha1.ConfigMapSync
	err := r.Get(ctx, req.NamespacedName, &cms)
	if errors.IsNotFound(err) {
		// CR deleted, nothing to do.
		return ctrl.Result{}, nil
	}

	if err != nil {
		return ctrl.Result{}, err
	}

	// 2. Fetch the source ConfigMap
	var source corev1.ConfigMap
	key := types.NamespacedName{Namespace: cms.Spec.SourceNamespace, Name: cms.Spec.SourceName}
	if err := r.Get(ctx, key, &source); err != nil {
		// If source not found, we might requeue to try later.
		log.Error(err, "Source ConfigMap not found")
		return ctrl.Result{}, nil
	}

	// 3. Sync to target namespaces
	var synced []string
	for _, ns := range cms.Spec.TargetNamespaces {
		target := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cms.Spec.SourceName,
				Namespace: ns,
			},
			Data: source.Data,
		}

		var existing corev1.ConfigMap
		err := r.Get(ctx, types.NamespacedName{Name: target.Name, Namespace: target.Namespace}, &existing)
		if err != nil && !errors.IsNotFound(err) {
			// Not expected error.
			log.Error(err, "Error fetching ConfigMap in namespace", "namespace", ns)
			continue
		}

		// Create the ConfigMap in target namespace.
		if errors.IsNotFound(err) {
			if err := r.Create(ctx, target); err != nil {
				log.Error(err, "Failed to create ConfigMap in namespace", "namespace", ns)
				continue
			}

			synced = append(synced, ns)
			continue
		}

		// Update the ConfigMap in target namespace.
		existing.Data = source.Data
		if err := r.Update(ctx, &existing); err != nil {
			log.Error(err, "Failed to update ConfigMap in namespace", "namespace", ns)
			continue
		}
		synced = append(synced, ns)
	}

	// Update status
	cms.Status.SyncedNamespaces = synced
	if err := r.Status().Update(ctx, &cms); err != nil {
		log.Error(err, "Failed to update ConfigMapSync status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigMapSyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&syncv1alpha1.ConfigMapSync{}).
		Complete(r)
}
