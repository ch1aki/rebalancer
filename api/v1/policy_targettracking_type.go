package v1

type TargetTrackingPolicy struct {
	TargetValue int64 `json:"targetValue"`
	BaseValue   int64 `json:"baseValue"`

	// +optional
	DisableScaleIn bool `json:"disableScaleIn,omitempty"`

	// +optional
	Scheduled []Scheduled `json:"scheduled"`
}

type Scheduled struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Value     int64  `json:"value"`
}
