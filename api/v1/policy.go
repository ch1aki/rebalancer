package v1

import (
	"context"
)

// +kubebuilder:object:root=false
// +kubebuilder:object:generate:false
// +k8s:deepcopy-gen:interfaces=nil
// +k8s:deepcopy-gen=nil

// Policy is a common interface for scaling
type Policy interface {
	New(rebalance *Rebalance, target *TargetClient, metrics *MetricsClient) (Estimator, error)
}

// +kubebuilder:object:root=false
// +kubebuilder:object:generate:false
// +k8s:deepcopy-gen:interfaces=nil
// +k8s:deepcopy-gen=nil
type Estimator interface {
	Estimate(ctx context.Context) (int64, error)
}
