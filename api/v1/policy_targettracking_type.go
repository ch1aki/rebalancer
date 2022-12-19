package v1

type TargetTrackingPolicy struct {
	TargetValue int64 `json:"targetValue"`
	BaseValue   int64 `json:"baseValue"`

	// +optional
	DisableScaleIn bool `json:"disableScaleIn,omitempty"`

	// +optional
	Scheduled []Scheduled `json:"scheduled,omitempty"`
}

type Scheduled struct {
	StartTime Time  `json:"startTime"`
	EndTime   Time  `json:"endTime"`
	Value     int64 `json:"value"`
}

type Time struct {
	Hour int64 `json:"hour"`

	// +optional
	Min int64 `json:"min,omitempty"`
}
