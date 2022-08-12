//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSAuth) DeepCopyInto(out *AWSAuth) {
	*out = *in
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(AWSAuthSecretRef)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSAuth.
func (in *AWSAuth) DeepCopy() *AWSAuth {
	if in == nil {
		return nil
	}
	out := new(AWSAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSAuthSecretRef) DeepCopyInto(out *AWSAuthSecretRef) {
	*out = *in
	in.AccessKeyID.DeepCopyInto(&out.AccessKeyID)
	in.SecretAccessKey.DeepCopyInto(&out.SecretAccessKey)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSAuthSecretRef.
func (in *AWSAuthSecretRef) DeepCopy() *AWSAuthSecretRef {
	if in == nil {
		return nil
	}
	out := new(AWSAuthSecretRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BasicAuth) DeepCopyInto(out *BasicAuth) {
	*out = *in
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(BasicAuthSecretRef)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BasicAuth.
func (in *BasicAuth) DeepCopy() *BasicAuth {
	if in == nil {
		return nil
	}
	out := new(BasicAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BasicAuthSecretRef) DeepCopyInto(out *BasicAuthSecretRef) {
	*out = *in
	in.User.DeepCopyInto(&out.User)
	in.Password.DeepCopyInto(&out.Password)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BasicAuthSecretRef.
func (in *BasicAuthSecretRef) DeepCopy() *BasicAuthSecretRef {
	if in == nil {
		return nil
	}
	out := new(BasicAuthSecretRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusDataSource) DeepCopyInto(out *PrometheusDataSource) {
	*out = *in
	in.Auth.DeepCopyInto(&out.Auth)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusDataSource.
func (in *PrometheusDataSource) DeepCopy() *PrometheusDataSource {
	if in == nil {
		return nil
	}
	out := new(PrometheusDataSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Rebalance) DeepCopyInto(out *Rebalance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rebalance.
func (in *Rebalance) DeepCopy() *Rebalance {
	if in == nil {
		return nil
	}
	out := new(Rebalance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Rebalance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RebalanceDataSource) DeepCopyInto(out *RebalanceDataSource) {
	*out = *in
	if in.Prometheus != nil {
		in, out := &in.Prometheus, &out.Prometheus
		*out = new(PrometheusDataSource)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RebalanceDataSource.
func (in *RebalanceDataSource) DeepCopy() *RebalanceDataSource {
	if in == nil {
		return nil
	}
	out := new(RebalanceDataSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RebalanceList) DeepCopyInto(out *RebalanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Rebalance, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RebalanceList.
func (in *RebalanceList) DeepCopy() *RebalanceList {
	if in == nil {
		return nil
	}
	out := new(RebalanceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RebalanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RebalanceRule) DeepCopyInto(out *RebalanceRule) {
	*out = *in
	out.Flactation = in.Flactation
	out.Threshold = in.Threshold
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RebalanceRule.
func (in *RebalanceRule) DeepCopy() *RebalanceRule {
	if in == nil {
		return nil
	}
	out := new(RebalanceRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RebalanceRuleFlactation) DeepCopyInto(out *RebalanceRuleFlactation) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RebalanceRuleFlactation.
func (in *RebalanceRuleFlactation) DeepCopy() *RebalanceRuleFlactation {
	if in == nil {
		return nil
	}
	out := new(RebalanceRuleFlactation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RebalanceRuleThreshold) DeepCopyInto(out *RebalanceRuleThreshold) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RebalanceRuleThreshold.
func (in *RebalanceRuleThreshold) DeepCopy() *RebalanceRuleThreshold {
	if in == nil {
		return nil
	}
	out := new(RebalanceRuleThreshold)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RebalanceSpec) DeepCopyInto(out *RebalanceSpec) {
	*out = *in
	in.Target.DeepCopyInto(&out.Target)
	in.DataSource.DeepCopyInto(&out.DataSource)
	out.Rule = in.Rule
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RebalanceSpec.
func (in *RebalanceSpec) DeepCopy() *RebalanceSpec {
	if in == nil {
		return nil
	}
	out := new(RebalanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RebalanceStatus) DeepCopyInto(out *RebalanceStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RebalanceStatus.
func (in *RebalanceStatus) DeepCopy() *RebalanceStatus {
	if in == nil {
		return nil
	}
	out := new(RebalanceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RebalanceTarget) DeepCopyInto(out *RebalanceTarget) {
	*out = *in
	if in.Route53 != nil {
		in, out := &in.Route53, &out.Route53
		*out = new(Route53Target)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RebalanceTarget.
func (in *RebalanceTarget) DeepCopy() *RebalanceTarget {
	if in == nil {
		return nil
	}
	out := new(RebalanceTarget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Route53Target) DeepCopyInto(out *Route53Target) {
	*out = *in
	out.Resource = in.Resource
	in.Auth.DeepCopyInto(&out.Auth)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Route53Target.
func (in *Route53Target) DeepCopy() *Route53Target {
	if in == nil {
		return nil
	}
	out := new(Route53Target)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Route53TargetRecord) DeepCopyInto(out *Route53TargetRecord) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Route53TargetRecord.
func (in *Route53TargetRecord) DeepCopy() *Route53TargetRecord {
	if in == nil {
		return nil
	}
	out := new(Route53TargetRecord)
	in.DeepCopyInto(out)
	return out
}
