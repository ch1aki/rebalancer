package v1

type BasicAuth struct {
	SecretRef *BasicAuthSecretRef `json:"secretRef,omitempty"`
}

type BasicAuthSecretRef struct {
	// The User is used for authentication
	User SecretKeySelector `json:"userSecretRef,omitempty"`

	// The Password is used for authentication
	Password SecretKeySelector `json:"passwordSecretRef,omitempty"`
}

type PrometheusDataSource struct {
	Address string `json:"address"`
	Query   string `json:"query"`

	// +optional
	Timeout string `json:"timeout"`

	// +optional
	Auth BasicAuth `json:"auth"`
}
