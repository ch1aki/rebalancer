package controllers

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	metricsNamespace = "rebalancer"
)

var (
	ErrorVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "error",
		Help:      "The cluster status about error condition",
	}, []string{"name", "namespace"})

	UnhealthyVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "available",
		Help:      "The cluster status about available condition",
	}, []string{"name", "namespace"})

	HealthyVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "healthy",
		Help:      "The cluster status about healthy condition",
	}, []string{"name", "namespace"})

	DesiredValVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "desired",
		Help:      "desired value",
	}, []string{"name", "namespace"})

	ActualValVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "actual",
		Help:      "actual value",
	}, []string{"name", "namespace"})
)

func init() {
	metrics.Registry.MustRegister(ErrorVec, UnhealthyVec, HealthyVec, DesiredValVec, ActualValVec)
}
