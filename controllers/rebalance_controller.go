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
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	rebalancerv1 "git.pepabo.com/akichan/rebalancer/api/v1"
	"git.pepabo.com/akichan/rebalancer/controllers/analysis"
	"git.pepabo.com/akichan/rebalancer/controllers/provider"
)

// RebalanceReconciler reconciles a Rebalance object
type RebalanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=rebalancer.ch1aki.github.io,resources=rebalances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rebalancer.ch1aki.github.io,resources=rebalances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=rebalancer.ch1aki.github.io,resources=rebalances/finalizers,verbs=update

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

	var rebalance rebalancerv1.Rebalance
	err := r.Get(ctx, req.NamespacedName, &rebalance)
	if errors.IsNotFound(err) {
		return ctrl.Result{}, nil
	}
	if err != nil {
		logger.Error(err, "unable to get Rebalance", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	if !rebalance.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	r.rebalance(ctx, rebalance)

	refreshInt, err := time.ParseDuration(rebalance.Spec.Rule.Interval)
	if err != nil {
		logger.Error(err, "unable to parse interval string", "interval", rebalance.Spec.Rule.Interval)
		return ctrl.Result{}, err
	}

	r.updateStatus(ctx, rebalance)

	return ctrl.Result{
		RequeueAfter: refreshInt,
	}, nil
}

func (r *RebalanceReconciler) updateStatus(ctx context.Context, rebalance rebalancerv1.Rebalance) (ctrl.Result, error) {
	current := time.Now()
	rebalance.Status.LastUpdateAt = current.Format(time.RFC3339)

	return ctrl.Result{}, nil
}

func (r *RebalanceReconciler) rebalance(ctx context.Context, rebalance rebalancerv1.Rebalance) error {
	max, err := strconv.ParseInt(rebalance.Spec.Rule.Flactation.Max, 10, 64)
	if err != nil {
		return err
	}
	min, err := strconv.ParseInt(rebalance.Spec.Rule.Flactation.Min, 10, 64)
	if err != nil {
		return err
	}
	valiation, err := strconv.ParseInt(rebalance.Spec.Rule.Flactation.Variation, 10, 64)
	if err != nil {
		return err
	}

	// fetch metrics
	a, err := analysis.NewClient(ctx, rebalance)
	if err != nil {
		return err
	}
	cond, err := a.Eval(ctx)
	if err != nil {
		return err
	}

	// get current weight
	p, err := provider.NewProvider(ctx, rebalance)
	if err != nil {
		return err
	}
	current, err := p.GetWeight(ctx)
	if err != nil {
		return err
	}

	// change weight
	var newWeight int64
	if cond {
		// increse weight
		if current < max-valiation {
			newWeight = current + valiation
		}
	} else {
		// decrese weight
		newWeight = current - valiation
		if newWeight < min {
			newWeight = min
		}
	}
	if current != newWeight {
		p.SetWeight(ctx, newWeight)
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RebalanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&rebalancerv1.Rebalance{}).
		Complete(r)
}
