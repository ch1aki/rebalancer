package prometheus

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/argoproj/argo-rollouts/utils/evaluate"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"

	rebalancerv1 "git.pepabo.com/akichan/rebalancer/api/v1"
)

type Metrics struct {
	api              v1.API
	queryString      string
	timeout          time.Duration
	successCondition string
	name             string
}

func (m *Metrics) NewClient(ctx context.Context, r rebalancerv1.Rebalance) (rebalancerv1.MetricsClient, error) {
	u, err := url.Parse(r.Spec.Metrics.Prometheus.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url in %s: %w", r.Name, err)
	} else if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("url must contain scheme and host: %w", err)
	}

	c, err := api.NewClient(api.Config{
		Address: u.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create prometheus client: %w", err)
	}

	return &Metrics{
		api:         v1.NewAPI(c),
		queryString: r.Spec.Metrics.Prometheus.Query,
		timeout:     time.Duration(r.Spec.Metrics.Prometheus.Timeout) * time.Second,
		name:        r.Name,
	}, nil
}

func (m *Metrics) Fetch(ctx context.Context) (float64, error) {
	responce, err := m.query(ctx)
	if err != nil {
		return 0, err
	}

	switch value := responce.(type) {
	case *model.Scalar:
		return float64(value.Value), nil
	default:
		return 0, fmt.Errorf("prometheus metric is expected to return scalar value: %s", m.name)
	}
}

func (m *Metrics) Evaluate(ctx context.Context, expression string) (bool, error) {
	responce, err := m.query(ctx)
	if err != nil {
		return false, err
	}

	switch value := responce.(type) {
	case *model.Scalar:
		result := float64(value.Value)
		return evaluate.EvalCondition(result, expression)
	case model.Vector:
		results := make([]float64, 0, len(value))
		for _, s := range value {
			if s != nil {
				results = append(results, float64(s.Value))
			}
		}
		return evaluate.EvalCondition(results, expression)
	default:
		return false, fmt.Errorf("prometheus metric type not supported")
	}
}

func (m *Metrics) query(ctx context.Context) (model.Value, error) {
	c, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	// query
	responce, warnings, err := m.api.Query(c, m.queryString, time.Now())
	if err != nil {
		return responce, err
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

	return responce, nil
}

func init() {
	rebalancerv1.RegisterMetrics(&Metrics{}, &rebalancerv1.RebalanceMetrics{
		Prometheus: &rebalancerv1.PrometheusMetrics{},
	})
}
