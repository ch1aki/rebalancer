package v1

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:object:root=false
// +kubebuilder:object:generate:false
// +k8s:deepcopy-gen:interfaces=nil
// +k8s:deepcopy-gen=nil

type Target interface {
	NewClient(ctx context.Context, r Rebalance, c client.Client) (TargetClient, error)
}

// +kubebuilder:object:root=false
// +kubebuilder:object:generate:false
// +k8s:deepcopy-gen:interfaces=nil
// +k8s:deepcopy-gen=nil

type TargetClient interface {
	GetWeight(ctx context.Context) (int64, error)
	SetWeight(ctx context.Context, value int64) error
}
