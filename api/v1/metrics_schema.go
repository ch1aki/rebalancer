package v1

import (
	"encoding/json"
	"fmt"
	"sync"
)

var metricsBuilder map[string]Metrics
var metricsBuildLock sync.RWMutex

func init() {
	metricsBuilder = make(map[string]Metrics)
}

// Register a metrics type. Registermetrics panics if a
// metrics with the same metrics is already registered.
func RegisterMetrics(m Metrics, metricsSpec *RebalanceMetrics) {
	metricsName, err := getMetricsName(metricsSpec)
	if err != nil {
		panic(fmt.Sprintf("store err registring scheme: %s", err.Error()))
	}

	metricsBuildLock.Lock()
	defer metricsBuildLock.Unlock()
	_, exists := metricsBuilder[metricsName]
	if exists {
		panic(fmt.Sprintf("metrics %q already registerd", metricsName))
	}

	metricsBuilder[metricsName] = m
}

func GetMetricsByName(name string) (Metrics, bool) {
	metricsBuildLock.RLock()
	f, ok := metricsBuilder[name]
	metricsBuildLock.RUnlock()
	return f, ok
}

// Getmetrics returns the metrics from the rebalance
func GetMetrics(r Rebalance) (Metrics, error) {
	spec := &r.Spec.Metrics
	metricsName, err := getMetricsName(spec)
	if err != nil {
		return nil, fmt.Errorf("metrics err for %s: %w", r.GetName(), err)
	}

	metricsBuildLock.RLock()
	f, ok := metricsBuilder[metricsName]
	metricsBuildLock.RUnlock()

	if !ok {
		return nil, fmt.Errorf("failed to find registerd metrics for type: %s, name: %s", metricsName, r.GetName())
	}

	return f, nil
}

// getMetricsName returns the name of the configured metrics
// or an error if the metrics is not configured
func getMetricsName(metricsSpec *RebalanceMetrics) (string, error) {
	metricsBytes, err := json.Marshal(metricsSpec)
	if err != nil || metricsBytes == nil {
		return "", fmt.Errorf("failed to marshal metrics spec: %w", err)
	}

	metricsMap := make(map[string]interface{})
	err = json.Unmarshal(metricsBytes, &metricsMap)
	if err != nil {
		return "", fmt.Errorf("fialed to unmarshal metrics spec: %w", err)
	}

	if len(metricsMap) != 1 {
		return "", fmt.Errorf("metrics must only have exactly one specified, found %d", len(metricsMap))
	}

	for k := range metricsMap {
		return k, nil
	}

	return "", fmt.Errorf("failed to find registerd metrics")
}
