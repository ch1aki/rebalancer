package v1

import (
	route53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

type AWSAuth struct {
	SecretRef *AWSAuthSecretRef `json:"secretRef,omitempty"`
}

type AWSAuthSecretRef struct {
	// The AccessKeyID is used for authentication
	AccessKeyID SecretKeySelector `json:"accessKeyIDSecretRef,omitempty"`

	// The SecretAccessKey is used for authentication
	SecretAccessKey SecretKeySelector `json:"secretAccessKeySecretRef,omitempty"`
}

type Route53TargetRecord struct {
	Name string              `json:"name"`
	Type route53Types.RRType `json:"type"`

	// +optional
	Identifier string `json:"identifier,omitempty"`
}

type Route53Target struct {
	HostedZoneID string              `json:"hostedZoneID"`
	Resource     Route53TargetRecord `json:"resource"`

	// +optional
	Auth AWSAuth `json:"auth"`
}
