/*
Copyright 2022.

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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	rebalancerv1 "git.pepabo.com/akichan/rebalancer/api/v1"
	_ "git.pepabo.com/akichan/rebalancer/controllers/metrics/register"
	_ "git.pepabo.com/akichan/rebalancer/controllers/policy/register"
	_ "git.pepabo.com/akichan/rebalancer/controllers/target/register"
)

// RebalanceReconciler reconciles a Rebalance object
type RebalanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=rebalancer.ch1aki.github.io,resources=rebalances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rebalancer.ch1aki.github.io,resources=rebalances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=rebalancer.ch1aki.github.io,resources=rebalances/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Rebalance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *RebalanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var rb rebalancerv1.Rebalance
	err := r.Get(ctx, req.NamespacedName, &rb)
	if errors.IsNotFound(err) {
		r.removeMetrics(rb)
		return ctrl.Result{}, nil
	}
	if err != nil {
		logger.Error(err, "unable to get Rebalance", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	if !rb.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	interval, err := time.ParseDuration(rb.Spec.Interval)
	if err != nil {
		logger.Error(err, "unable to parse interval string", "interval", rb.Spec.Interval)
		return ctrl.Result{}, err
	}

	desired, actual, err := r.rebalance(ctx, rb, r.Client)
	if err != nil {
		logger.Error(err, "rebalance operation failed", "interval", rb.Spec.Interval)
		return ctrl.Result{}, err
	}

	err = r.updateStatus(ctx, rb, desired, actual)
	if err != nil {
		logger.Error(err, "rebalance operation failed", "update status", rb.Spec)
	}

	return ctrl.Result{
		RequeueAfter: interval,
	}, nil
}

func (r *RebalanceReconciler) updateStatus(ctx context.Context, rb rebalancerv1.Rebalance, desired int64, actual int64) error {
	var status rebalancerv1.RebalanceStatus

	// rebalance status
	if desired == actual {
		status.Condition = rebalancerv1.RebalanceHealty
	} else if rb.Spec.DryRun {
		status.Condition = rebalancerv1.RebalanceUnhealthy
	} else {
		status.Condition = rebalancerv1.RebalanceError
	}

	// weight
	status.DesiredValue = desired
	status.ActualValue = actual

	// update
	if rb.Status != status {
		r.setMetrics(rb)
		status.LastUpdateAt = time.Now().Format(time.RFC3339)
		rb.Status = status
		err := r.Status().Update(ctx, &rb)
		if err != nil {
			return err
		}
	}

	if rb.Status.Condition != rebalancerv1.RebalanceUnhealthy &&
		rb.Status.Condition != rebalancerv1.RebalanceHealty {
		return nil
	}
	return nil
}

func (r *RebalanceReconciler) rebalance(ctx context.Context, rb rebalancerv1.Rebalance, c client.Client) (desired int64, actual int64, e error) {
	// get metrics client
	metrics, err := rebalancerv1.GetMetrics(rb)
	if err != nil {
		return 0, 0, fmt.Errorf("failes to get metrics: %w", err)
	}
	metricsClient, err := metrics.NewClient(ctx, rb)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to initialize metrics client")
	}

	// get target client
	target, err := rebalancerv1.GetTarget(rb)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get target: %w", err)
	}
	targetClient, err := target.NewClient(ctx, rb, c)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to initialize target client: %w", err)
	}

	// get policy
	p, err := rebalancerv1.GetPolicy(rb)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get policy: %w", err)
	}
	policy, err := p.New(&rb, &targetClient, &metricsClient)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to initialize policy")
	}

	// estimate target val
	desired, err = policy.Estimate(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to estimate targeet value: %w", err)
	}

	// get target actual value
	actual, err = targetClient.GetWeight(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed get current value: %w", err)
	}

	// set weight
	if actual != desired && !rb.Spec.DryRun {
		err = targetClient.SetWeight(ctx, desired)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to set target value: %w", err)
		}
	}

	// get target actual value
	actual, err = targetClient.GetWeight(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed get current value: %w", err)
	}

	return desired, actual, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RebalanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&rebalancerv1.Rebalance{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

func (r *RebalanceReconciler) setMetrics(rb rebalancerv1.Rebalance) {
	switch rb.Status.Condition {
	case rebalancerv1.RebalanceError:
		ErrorVec.WithLabelValues(rb.Name, rb.Namespace).Set(1)
		UnhealthyVec.WithLabelValues(rb.Name, rb.Namespace).Set(0)
		HealthyVec.WithLabelValues(rb.Name, rb.Namespace).Set(0)
	case rebalancerv1.RebalanceUnhealthy:
		ErrorVec.WithLabelValues(rb.Name, rb.Namespace).Set(0)
		UnhealthyVec.WithLabelValues(rb.Name, rb.Namespace).Set(1)
		HealthyVec.WithLabelValues(rb.Name, rb.Namespace).Set(0)
	case rebalancerv1.RebalanceHealty:
		ErrorVec.WithLabelValues(rb.Name, rb.Namespace).Set(0)
		UnhealthyVec.WithLabelValues(rb.Name, rb.Namespace).Set(0)
		HealthyVec.WithLabelValues(rb.Name, rb.Namespace).Set(1)
	}

	DesiredValVec.WithLabelValues(rb.Name, rb.Namespace).Set(float64(rb.Status.DesiredValue))
	ActualValVec.WithLabelValues(rb.Name, rb.Namespace).Set(float64(rb.Status.ActualValue))
}

func (r *RebalanceReconciler) removeMetrics(rb rebalancerv1.Rebalance) {
	ErrorVec.DeleteLabelValues(rb.Name, rb.Namespace)
	UnhealthyVec.DeleteLabelValues(rb.Name, rb.Namespace)
	HealthyVec.DeleteLabelValues(rb.Name, rb.Namespace)
}
