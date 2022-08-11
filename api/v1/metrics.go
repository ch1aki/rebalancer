package v1

import (
	"context"
)

// +kubebuilder:object:root=false
// +kubebuilder:object:generate:false
// +k8s:deepcopy-gen:interfaces=nil
// +k8s:deepcopy-gen=nil

type Metrics interface {
	NewClient(ctx context.Context, rebalance Rebalance) (MetricsClient, error)
}

// +kubebuilder:object:root=false
// +kubebuilder:object:generate:false
// +k8s:deepcopy-gen:interfaces=nil
// +k8s:deepcopy-gen=nil

type MetricsClient interface {
	Evaluate(ctx context.Context, expression string) (bool, error)
	Fetch(ctx context.Context) (float64, error)
}
