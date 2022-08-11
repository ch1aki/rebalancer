package targettracking

import (
	"context"
	"fmt"
	"math"

	rebalancev1 "git.pepabo.com/akichan/rebalancer/api/v1"
)

type Policy struct {
	target              *rebalancev1.TargetClient
	metrics             *rebalancev1.MetricsClient
	trackingTargetValue int64
	baseValue           int64
	disableScaleIn      bool
}

func (p *Policy) New(rebalance *rebalancev1.Rebalance, target *rebalancev1.TargetClient,
	metrics *rebalancev1.MetricsClient) (rebalancev1.Estimator, error) {

	return &Policy{
		target:              target,
		metrics:             metrics,
		trackingTargetValue: rebalance.Spec.Policy.TargetTracking.TargetValue,
		baseValue:           rebalance.Spec.Policy.TargetTracking.BaseValue,
		disableScaleIn:      rebalance.Spec.Policy.TargetTracking.DisableScaleIn,
	}, nil
}

func (p *Policy) Estimate(ctx context.Context) (int64, error) {
	// get current metrics
	currentMetric, err := (*p.metrics).Fetch(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed get current metric: %w", err)
	}

	// process desired value
	v := processBestContrast(float64(p.baseValue), float64(p.trackingTargetValue), currentMetric)

	return v, nil
}

func processBestContrast(base float64, trackingTargetVal float64, current float64) int64 {
	rate := current/trackingTargetVal - 1
	if rate < 0 {
		rate = 0
	}
	return int64(math.Ceil(base * rate))
}

func init() {
	rebalancev1.RegisterPolicy(&Policy{}, &rebalancev1.RebalancePolicy{
		TargetTracking: &rebalancev1.TargetTrackingPolicy{},
	})
}
