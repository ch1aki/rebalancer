package v1

type TargetTrackingPolicy struct {
	TargetValue int64 `json:"targetValue"`
	BaseValue   int64 `json:"baseValue"`

	// +optional
	DisableScaleIn bool `json:"disableScaleIn,omitempty"`
}
