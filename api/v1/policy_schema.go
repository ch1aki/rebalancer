package v1

import (
	"encoding/json"
	"fmt"
	"sync"
)

var policyBuilder map[string]Policy
var policyBuildLock sync.RWMutex

func init() {
	policyBuilder = make(map[string]Policy)
}

// Register a policy type. RegisterPolicy panics if a
// policy with the same policy is already registered.
func RegisterPolicy(p Policy, policySpec *RebalancePolicy) {
	policyName, err := getPolicyName(policySpec)
	if err != nil {
		panic(fmt.Sprintf("store err registring scheme: %s", err.Error()))
	}

	policyBuildLock.Lock()
	defer policyBuildLock.Unlock()
	_, exists := policyBuilder[policyName]
	if exists {
		panic(fmt.Sprintf("policy %q already registerd", policyName))
	}

	policyBuilder[policyName] = p
}

// GetPolicy returns the policy from the rebalance
func GetPolicy(r Rebalance) (Policy, error) {
	spec := &r.Spec.Policy
	policyName, err := getPolicyName(spec)
	if err != nil {
		return nil, fmt.Errorf("policy err for %s: %w", r.GetName(), err)
	}

	policyBuildLock.RLock()
	f, ok := policyBuilder[policyName]
	policyBuildLock.RUnlock()

	if !ok {
		return nil, fmt.Errorf("failed to find registerd policy for type: %s, name: %s", policyName, r.GetName())
	}

	return f, nil
}

// getPolicyName returns the name of the configured policy
// or an error if the policy is not configured
func getPolicyName(policySpec *RebalancePolicy) (string, error) {
	policyBytes, err := json.Marshal(policySpec)
	if err != nil || policyBytes == nil {
		return "", fmt.Errorf("failed to marshal policy spec: %w", err)
	}

	policyMap := make(map[string]interface{})
	err = json.Unmarshal(policyBytes, &policyMap)
	if err != nil {
		return "", fmt.Errorf("fialed to unmarshal policy spec: %w", err)
	}

	if len(policyMap) != 1 {
		return "", fmt.Errorf("policy must only have exactly one specified, found %d", len(policyMap))
	}

	for k := range policyMap {
		return k, nil
	}

	return "", fmt.Errorf("failed to find registerd policy")
}
