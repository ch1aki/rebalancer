package targettracking

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	rebalancev1 "git.pepabo.com/akichan/rebalancer/api/v1"
	"github.com/thoas/go-funk"
)

type Policy struct {
	target              *rebalancev1.TargetClient
	metrics             *rebalancev1.MetricsClient
	trackingTargetValue int64
	baseValue           int64
	disableScaleIn      bool
	scheduled           []rebalancev1.Scheduled
}

func (p *Policy) New(rebalance *rebalancev1.Rebalance, target *rebalancev1.TargetClient,
	metrics *rebalancev1.MetricsClient) (rebalancev1.Estimator, error) {

	return &Policy{
		target:              target,
		metrics:             metrics,
		trackingTargetValue: rebalance.Spec.Policy.TargetTracking.TargetValue,
		baseValue:           rebalance.Spec.Policy.TargetTracking.BaseValue,
		disableScaleIn:      rebalance.Spec.Policy.TargetTracking.DisableScaleIn,
		scheduled:           rebalance.Spec.Policy.TargetTracking.Scheduled,
	}, nil
}

func (p *Policy) Estimate(ctx context.Context) (int64, error) {
	// get current metrics
	currentMetric, err := (*p.metrics).Fetch(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed get current metric: %w", err)
	}

	// process desired value
	val := processBestContrast(float64(p.baseValue), float64(p.trackingTargetValue), currentMetric)

	// check scheduled values
	if len(p.scheduled) > 0 {
		nowTime := time.Now()
		val, err = checkScheduledValue(p.scheduled, val, nowTime)
		if err != nil {
			return 0, fmt.Errorf("failed to check scheduled value")
		}
	}

	return val, nil
}

func processBestContrast(base float64, trackingTargetVal float64, current float64) int64 {
	rate := current/trackingTargetVal - 1
	if rate < 0 {
		rate = 0
	}
	return int64(math.Ceil(base * rate))
}

func checkScheduledValue(scheduled []rebalancev1.Scheduled, v int64, nowTime time.Time) (int64, error) {
	var values []int
	values = append(values, int(v))
	for _, s := range scheduled {
		startTime, err := parseTime(s.StartTime, nowTime)
		if err != nil {
			return 0, err
		}
		endTime, _ := parseTime(s.EndTime, nowTime)
		if err != err {
			return 0, err
		}
		if (nowTime.Equal(startTime) || nowTime.After(startTime)) && nowTime.Before(endTime) {
			values = append(values, int(s.Value))
		}
	}
	return int64(funk.MaxInt(values)), nil
}

func parseTime(t string, n time.Time) (time.Time, error) {
	ts := strings.Split(t, ":")
	hour, err := strconv.Atoi(ts[0])
	if err != nil {
		return time.Time{}, err
	}
	min, err := strconv.Atoi(ts[1])
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(n.Year(), n.Month(), n.Day(), hour, min, 0, 0, time.Local), nil
}

func init() {
	rebalancev1.RegisterPolicy(&Policy{}, &rebalancev1.RebalancePolicy{
		TargetTracking: &rebalancev1.TargetTrackingPolicy{},
	})
}
