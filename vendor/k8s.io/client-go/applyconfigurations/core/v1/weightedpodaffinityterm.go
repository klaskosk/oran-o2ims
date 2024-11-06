/*
Copyright The Kubernetes Authors.

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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

// WeightedPodAffinityTermApplyConfiguration represents a declarative configuration of the WeightedPodAffinityTerm type for use
// with apply.
type WeightedPodAffinityTermApplyConfiguration struct {
	Weight          *int32                             `json:"weight,omitempty"`
	PodAffinityTerm *PodAffinityTermApplyConfiguration `json:"podAffinityTerm,omitempty"`
}

// WeightedPodAffinityTermApplyConfiguration constructs a declarative configuration of the WeightedPodAffinityTerm type for use with
// apply.
func WeightedPodAffinityTerm() *WeightedPodAffinityTermApplyConfiguration {
	return &WeightedPodAffinityTermApplyConfiguration{}
}

// WithWeight sets the Weight field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Weight field is set to the value of the last call.
func (b *WeightedPodAffinityTermApplyConfiguration) WithWeight(value int32) *WeightedPodAffinityTermApplyConfiguration {
	b.Weight = &value
	return b
}

// WithPodAffinityTerm sets the PodAffinityTerm field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PodAffinityTerm field is set to the value of the last call.
func (b *WeightedPodAffinityTermApplyConfiguration) WithPodAffinityTerm(value *PodAffinityTermApplyConfiguration) *WeightedPodAffinityTermApplyConfiguration {
	b.PodAffinityTerm = value
	return b
}
