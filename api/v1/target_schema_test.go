package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const shouldBeRegistered = "target should be registered"

type TT struct{}

// New constructs a SecretsManager Provider.
func (t *TT) NewClient(ctx context.Context, r Rebalance, c client.Client) (TargetClient, error) {
	return t, nil
}

func (t *TT) GetWeight(ctx context.Context) (int64, error) {
	return 0, nil
}

func (t *TT) SetWeight(ctx context.Context, value int64) error {
	return nil
}

// TestRegister tests if the Register function
// (1) panics if it tries to register something invalid
// (2) stores the correct provider.
func TestRegisterTarget(t *testing.T) {
	tbl := []struct {
		test      string
		name      string
		expPanic  bool
		expExists bool
		target    *RebalanceTarget
	}{
		{
			test:      "should panic when given an invalid target",
			name:      "route53",
			expPanic:  true,
			expExists: false,
			target:    &RebalanceTarget{},
		},
		{
			test:      "should register an correct target",
			name:      "route53",
			expExists: false,
			target: &RebalanceTarget{
				Route53: &Route53Target{},
			},
		},
		{
			test:      "should panic if already exists",
			name:      "route53",
			expPanic:  true,
			expExists: true,
			target: &RebalanceTarget{
				Route53: &Route53Target{},
			},
		},
	}
	for i := range tbl {
		row := tbl[i]
		t.Run(row.test, func(t *testing.T) {
			runTest(t,
				row.name,
				row.target,
				row.expPanic,
			)
		})
	}
}

func runTest(t *testing.T, name string, target *RebalanceTarget, expPanic bool) {
	testTarget := &TT{}
	rebalance := &Rebalance{
		Spec: RebalanceSpec{
			Target: *target,
		},
	}
	if expPanic {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Register should panic")
			}
		}()
	}
	RegisterTarget(testTarget, &rebalance.Spec.Target)
	t1, ok := GetTargetByName(name)
	assert.True(t, ok, shouldBeRegistered)
	assert.Equal(t, testTarget, t1)
	t2, err := GetTarget(*rebalance)
	assert.Nil(t, err)
	assert.Equal(t, testTarget, t2)
}
