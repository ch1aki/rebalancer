package analysis

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

type Prometheus struct {
	client      api.Client
	queryString string
	timeout     time.Duration
	condition   string
}

func NewClient(ctx context.Context, r rebalancerv1.Rebalance) (Prometheus, error) {
	u, err := url.Parse(r.Spec.DataSource.Prometheus.Address)
	if err != nil {
		return Prometheus{}, err
	} else if u.Scheme == "" || u.Host == "" {
		return Prometheus{}, err
	}

	c, err := api.NewClient(api.Config{
		Address: r.Spec.DataSource.Prometheus.Address,
	})
	if err != nil {
		return Prometheus{}, err
	}

	p := Prometheus{
		client:      c,
		queryString: r.Spec.DataSource.Prometheus.Query,
		timeout:     time.Duration(r.Spec.DataSource.Prometheus.Timeout) * time.Second,
		condition:   r.Spec.Rule.Condition,
	}

	return p, nil
}

func (p *Prometheus) Eval(ctx context.Context) (bool, error) {
	responce, err := p.query(ctx)
	if err != nil {
		return false, err
	}

	switch value := responce.(type) {
	case *model.Scalar:
		result := float64(value.Value)
		return evaluate.EvalCondition(result, p.condition)
	case model.Vector:
		results := make([]float64, 0, len(value))
		for _, s := range value {
			if s != nil {
				results = append(results, float64(s.Value))
			}
		}
		return evaluate.EvalCondition(results, p.condition)
	default:
		return false, fmt.Errorf("prometheus metric type not supported")
	}
}

func (p *Prometheus) query(ctx context.Context) (model.Value, error) {
	api := v1.NewAPI(p.client)
	c, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	// query
	responce, warnings, err := api.Query(c, p.queryString, time.Now())
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
