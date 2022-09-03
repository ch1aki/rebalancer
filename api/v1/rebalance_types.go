/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RebalanceSpec defines the desired state of Rebalance
type RebalanceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Used to configure the target. Only one target may be set
	Target RebalanceTarget `json:"target"`

	// Used to configure the datasource. Only one data source may be set
	DataSource RebalanceDataSource `json:"dataSource"`

	// Used to configure the rule
	Rule RebalanceRule `json:"rule"`

	// DryRun is the flag of dry-run operation.
	// +kubebuilder:default=false
	// +optional
	DryRun bool `json:"dryRun,omitempty"`
}

type RebalanceTarget struct {
	Route53 *Route53Target `json:"route53,omitempty"`
}

type RebalanceDataSource struct {
	Prometheus *PrometheusDataSource `json:"prometheus,omitempty"`
}

type RebalanceRule struct {
	Flactation RebalanceRuleFlactation `json:"flactation"`
	Condition  string                  `json:"condition"`

	// +optional
	Interval string `json:"interval,omitempty"`
}

type RebalanceRuleFlactation struct {
	Variation string `json:"variation"`
	Max       string `json:"max"`
	Min       string `json:"min"`
}

type RebalanceCondition string

const (
	RebalanceNotReady  = RebalanceCondition("NotReady")
	RebalanceAvailable = RebalanceCondition("Available")
	RebalanceHealty    = RebalanceCondition("Health")
)

// RebalanceStatus defines the observed state of Rebalance
type RebalanceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Condition RebalanceCondition `json:"condition"`

	// +optional
	CurrentValue string `json:"currentValue"`

	// +optional
	LastUpdateAt string `json:"lastUpdateAt"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Rebalance is the Schema for the rebalances API
type Rebalance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RebalanceSpec   `json:"spec,omitempty"`
	Status RebalanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RebalanceList contains a list of Rebalance
type RebalanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Rebalance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Rebalance{}, &RebalanceList{})
}
