/*
Copyright 2026.

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
	"time"

	restartv1alpha1 "github.com/jgiornazi/k8s-podrestart-operator/api/v1alpha1"
	"github.com/robfig/cron/v3"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// ScheduledRestartReconciler reconciles a ScheduledRestart object
type ScheduledRestartReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=restart.jgiornazi.dev,resources=scheduledrestarts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=restart.jgiornazi.dev,resources=scheduledrestarts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=restart.jgiornazi.dev,resources=scheduledrestarts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ScheduledRestart object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.23.3/pkg/reconcile
func (r *ScheduledRestartReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	// TODO(user): your logic here
	sr := &restartv1alpha1.ScheduledRestart{}
	err := r.Get(ctx, req.NamespacedName, sr)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	if sr.Spec.Suspend {
		return ctrl.Result{}, nil
	}
	schedule, err := cron.ParseStandard(sr.Spec.Schedule)
	if err != nil {
		return ctrl.Result{}, err
	}
	lastRestart := time.Time{} // zero time = beginning of time
	if sr.Status.LastRestartTime != nil {
		lastRestart = sr.Status.LastRestartTime.Time
	}
	nextRun := schedule.Next(lastRestart)
	if time.Now().Before(nextRun) {
		return ctrl.Result{RequeueAfter: time.Until(nextRun)}, nil
	}
	podList := &corev1.PodList{}
	err = r.List(ctx, podList, client.InNamespace(sr.Namespace), client.MatchingLabels(sr.Spec.Selector.MatchLabels))
	if err != nil {
		return ctrl.Result{}, err
	}

	for _, pod := range podList.Items {
		if err := r.Delete(ctx, &pod); err != nil {
			return ctrl.Result{}, err
		}
	}
	now := metav1.Now()
	sr.Status.LastRestartTime = &now

	nextRunTime := metav1.NewTime(schedule.Next(time.Now()))
	sr.Status.NextRestartTime = &nextRunTime

	err = r.Status().Update(ctx, sr)

	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Until(nextRunTime.Time)}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *ScheduledRestartReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&restartv1alpha1.ScheduledRestart{}).
		Named("scheduledrestart").
		Complete(r)
}
