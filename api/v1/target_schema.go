package v1

import (
	"encoding/json"
	"fmt"
	"sync"
)

var targetBuilder map[string]Target
var targetBuildLock sync.RWMutex

func init() {
	targetBuilder = make(map[string]Target)
}

// Register a target type. RegisterTarget panics if a
// target with the same target is already registered.
func RegisterTarget(t Target, targetSpec *RebalanceTarget) {
	targetName, err := getTargetName(targetSpec)
	if err != nil {
		panic(fmt.Sprintf("store err registring scheme: %s", err.Error()))
	}

	targetBuildLock.Lock()
	defer targetBuildLock.Unlock()
	_, exists := targetBuilder[targetName]
	if exists {
		panic(fmt.Sprintf("target %q already registerd", targetName))
	}

	targetBuilder[targetName] = t
}

func GetTargetByName(name string) (Target, bool) {
	targetBuildLock.RLock()
	f, ok := targetBuilder[name]
	targetBuildLock.RUnlock()
	return f, ok
}

// GetTarget returns the target from the rebalance
func GetTarget(r Rebalance) (Target, error) {
	spec := &r.Spec.Target
	targetName, err := getTargetName(spec)
	if err != nil {
		return nil, fmt.Errorf("target err for %s: %w", r.GetName(), err)
	}

	targetBuildLock.RLock()
	f, ok := targetBuilder[targetName]
	targetBuildLock.RUnlock()

	if !ok {
		return nil, fmt.Errorf("failed to find registerd target for type: %s, name: %s", targetName, r.GetName())
	}

	return f, nil
}

// getTargetName returns the name of the configured target
// or an error if the target is not configured
func getTargetName(targetSpec *RebalanceTarget) (string, error) {
	targetBytes, err := json.Marshal(targetSpec)
	if err != nil || targetBytes == nil {
		return "", fmt.Errorf("failed to marshal target spec: %w", err)
	}

	targetMap := make(map[string]interface{})
	err = json.Unmarshal(targetBytes, &targetMap)
	if err != nil {
		return "", fmt.Errorf("fialed to unmarshal target spec: %w", err)
	}

	if len(targetMap) != 1 {
		return "", fmt.Errorf("target must only have exactly one specified, found %d", len(targetMap))
	}

	for k := range targetMap {
		return k, nil
	}

	return "", fmt.Errorf("failed to find registerd target")
}
