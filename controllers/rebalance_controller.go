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
	"net/url"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/argoproj/argo-rollouts/utils/evaluate"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"

	rebalancerv1 "github.com/ch1aki/rebalancer/api/v1"
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

	_, err = r.evalCondition(ctx, rebalance)
	if err != nil {
		logger.Error(err, "unable to eval condition", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	refreshInt, err := time.ParseDuration(rebalance.Spec.Rule.Interval)
	if err != nil {
		logger.Error(err, "unable to parse interval string", "interval", rebalance.Spec.Rule.Interval)
		return ctrl.Result{}, err
	}
	return ctrl.Result{
		RequeueAfter: refreshInt,
	}, nil
}

func (r *RebalanceReconciler) evalCondition(ctx context.Context, rebalance rebalancerv1.Rebalance) (bool, error) {
	logger := log.FromContext(ctx)

	u, err := url.Parse(rebalance.Spec.DataSource.Prometheus.Address)
	if err != nil {
		logger.Error(err, "Error in parsing url")
		return false, err
	} else if u.Scheme == "" || u.Host == "" {
		logger.Error(err, "scheme or host missing")
		return false, err
	}

	client, err := api.NewClient(api.Config{
		Address: rebalance.Spec.DataSource.Prometheus.Address,
	})
	if err != nil {
		logger.Error(err, "Error creating client")
		return false, err
	}
	api := v1.NewAPI(client)
	c, cancel := context.WithTimeout(context.Background(), time.Duration(rebalance.Spec.DataSource.Prometheus.Timeout)*time.Second)
	defer cancel()

	// query
	responce, warnings, err := api.Query(c, rebalance.Spec.DataSource.Prometheus.Query, time.Now())
	if err != nil {
		logger.Error(err, "Error get Targets")
		return false, err
	}

	// output warn message
	if len(warnings) > 0 {
		warningMetadata := ""
		for _, warning := range warnings {
			warningMetadata = fmt.Sprintf(`%s"%s", `, warningMetadata, warning)
		}
		warningMetadata = warningMetadata[:len(warningMetadata)-2]
		if warningMetadata != "" {
			fmt.Printf("Prometheus returned the following warnings: %s", warningMetadata)
		}
	}

	// eval
	switch value := responce.(type) {
	case *model.Scalar:
		result := float64(value.Value)
		return evaluate.EvalCondition(result, string(rebalance.Spec.Rule.Condition))
	case model.Vector:
		results := make([]float64, 0, len(value))
		for _, s := range value {
			if s != nil {
				results = append(results, float64(s.Value))
			}
		}
		return evaluate.EvalCondition(results, string(rebalance.Spec.Rule.Condition))
	default:
		return false, fmt.Errorf("Prometheus metric type not supported")
	}
}

func (r *RebalanceReconciler) updateStatus(ctx context.Context, rebalance rebalancerv1.Rebalance) (ctrl.Result, error) {
	// TODO

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RebalanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&rebalancerv1.Rebalance{}).
		Complete(r)
}
