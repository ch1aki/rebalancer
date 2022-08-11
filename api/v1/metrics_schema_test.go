package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const shouldBeRegisteredMetrics = "metrics should be registered"

type MT struct{}

// New constructs a SecretsManager Provider.
func (m *MT) NewClient(ctx context.Context, r Rebalance) (MetricsClient, error) {
	return m, nil
}

func (m *MT) Evaluate(ctx context.Context, expresion string) (bool, error) {
	return true, nil
}

func (m *MT) Fetch(ctx context.Context) (float64, error) {
	return 0, nil
}

// TestRegister tests if the Register function
// (1) panics if it tries to register something invalid
// (2) stores the correct provider.
func TestRegisterMetrics(t *testing.T) {
	tbl := []struct {
		test      string
		name      string
		expPanic  bool
		expExists bool
		metrics   *RebalanceMetrics
	}{
		{
			test:      "should panic when given an invalid metrics",
			name:      "prometheus",
			expPanic:  true,
			expExists: false,
			metrics:   &RebalanceMetrics{},
		},
		{
			test:      "should register an correct metrics",
			name:      "prometheus",
			expExists: false,
			metrics: &RebalanceMetrics{
				Prometheus: &PrometheusMetrics{},
			},
		},
		{
			test:      "should panic if already exists",
			name:      "prometheus",
			expPanic:  true,
			expExists: true,
			metrics: &RebalanceMetrics{
				Prometheus: &PrometheusMetrics{},
			},
		},
	}
	for i := range tbl {
		row := tbl[i]
		t.Run(row.test, func(t *testing.T) {
			runTestMetrics(t,
				row.name,
				row.metrics,
				row.expPanic,
			)
		})
	}
}

func runTestMetrics(t *testing.T, name string, metrics *RebalanceMetrics, expPanic bool) {
	testMetrics := &MT{}
	rebalance := &Rebalance{
		Spec: RebalanceSpec{
			Metrics: *metrics,
		},
	}
	if expPanic {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Register should panic")
			}
		}()
	}
	RegisterMetrics(testMetrics, &rebalance.Spec.Metrics)
	t1, ok := GetMetricsByName(name)
	assert.True(t, ok, shouldBeRegisteredMetrics)
	assert.Equal(t, testMetrics, t1)
	t2, err := GetMetrics(*rebalance)
	assert.Nil(t, err)
	assert.Equal(t, testMetrics, t2)
}
