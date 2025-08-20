/*
Copyright 2025 jwwangari6.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy at

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
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	infrav1alpha1 "github.com/Jwwangari6/schedule-scaler-operator/api/v1alpha1"
)

// ScheduleScalerReconciler reconciles a ScheduleScaler object
type ScheduleScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// -----------------------------------------------------------------------------
// RBAC permissions
// -----------------------------------------------------------------------------
// These markers tell Kubebuilder/controller-gen what RBAC rules to generate.
// They allow this controller to read/write ScheduleScaler CRDs and scale workloads.
//
// +kubebuilder:rbac:groups=infra.my.domain,resources=schedulescalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infra.my.domain,resources=schedulescalers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infra.my.domain,resources=schedulescalers/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments;statefulsets,verbs=get;list;watch;patch;update

// Reconcile is part of the main Kubernetes reconciliation loop.
// It compares the desired state (ScheduleScaler.Spec) with the actual cluster state,
// and performs operations to make the actual state match the desired schedule.
func (r *ScheduleScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	// Fetch the ScheduleScaler resource
	var scaler infrav1alpha1.ScheduleScaler
	if err := r.Get(ctx, req.NamespacedName, &scaler); err != nil {
		// Object not found (deleted) -> ignore
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Current time details
	now := time.Now()
	day := now.Weekday().String()[:3] // e.g. "Mon", "Tue"
	hourMin := now.Format("15:04")

	// Helper: check if the schedule time matches the current time
	shouldRun := func(t string) bool { return t == hourMin }

	// Iterate over all defined targets in the spec
	for _, target := range scaler.Spec.Targets {
		var replicas int32

		// Decide whether to scale up or down
		if shouldRun(scaler.Spec.Schedule.ScaleDown) && contains(scaler.Spec.Schedule.Days, day) {
			replicas = target.ScaleDown
		} else if shouldRun(scaler.Spec.Schedule.ScaleUp) && contains(scaler.Spec.Schedule.Days, day) {
			replicas = target.ScaleUp
		} else {
			// No action for this target this cycle
			continue
		}

		// Handle Deployments
		if strings.EqualFold(target.Kind, "Deployment") {
			var dep appsv1.Deployment
			if err := r.Get(ctx, types.NamespacedName{Name: target.Name, Namespace: scaler.Spec.Namespace}, &dep); err == nil {
				dep.Spec.Replicas = &replicas
				if err := r.Update(ctx, &dep); err != nil {
					logger.Error(err, "Failed to update Deployment", "name", target.Name)
				} else {
					logger.Info("Scaled Deployment", "name", target.Name, "replicas", replicas)
				}
			} else {
				logger.Error(err, "Failed to fetch Deployment", "name", target.Name)
			}
		}

		// TODO: Add support for StatefulSet scaling
	}

	// Requeue every minute to re-check the schedule
	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

// SetupWithManager registers the controller with the manager so it starts watching ScheduleScaler resources.
func (r *ScheduleScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1alpha1.ScheduleScaler{}).
		Named("schedulescaler").
		Complete(r)
}

// contains helper: checks if a string slice contains a value
func contains(list []string, v string) bool {
	for _, i := range list {
		if i == v {
			return true
		}
	}
	return false
}
