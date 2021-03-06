/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

/*
This file (together with pkg/expapi/v1/types.go) contain the experimental
types in kubernetes. These API objects are experimental, meaning that the
APIs may be broken at any time by the kubernetes team.

DISCLAIMER: The implementation of the experimental API group itself is
a temporary one meant as a stopgap solution until kubernetes has proper
support for multiple API groups. The transition may require changes
beyond registration differences. In other words, experimental API group
support is experimental.
*/

package expapi

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/resource"
	"k8s.io/kubernetes/pkg/util"
)

// ScaleSpec describes the attributes a Scale subresource
type ScaleSpec struct {
	// Replicas is the number of desired replicas.
	Replicas int `json:"replicas,omitempty" description:"number of replicas desired;  http://releases.k8s.io/HEAD/docs/user-guide/replication-controller.md#what-is-a-replication-controller"`
}

// ScaleStatus represents the current status of a Scale subresource.
type ScaleStatus struct {
	// Replicas is the number of actual replicas.
	Replicas int `json:"replicas" description:"most recently oberved number of replicas; see http://releases.k8s.io/HEAD/docs/user-guide/replication-controller.md#what-is-a-replication-controller"`

	// Selector is a label query over pods that should match the replicas count.
	Selector map[string]string `json:"selector,omitempty" description:"label keys and values that must match in order to be controlled by this replication controller, if empty defaulted to labels on Pod template; see http://releases.k8s.io/HEAD/docs/user-guide/labels.md#label-selectors"`
}

// Scale subresource, applicable to ReplicationControllers and (in future) Deployment.
type Scale struct {
	api.TypeMeta   `json:",inline"`
	api.ObjectMeta `json:"metadata,omitempty" description:"standard object metadata; see http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata"`

	// Spec defines the behavior of the scale.
	Spec ScaleSpec `json:"spec,omitempty" description:"specification of the desired behavior of the scale; http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#spec-and-status"`

	// Status represents the current status of the scale.
	Status ScaleStatus `json:"status,omitempty" description:"most recently observed status of the service; populated by the system, read-only; http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#spec-and-status"`
}

// Dummy definition
type ReplicationControllerDummy struct {
	api.TypeMeta `json:",inline"`
}

// SubresourceReference contains enough information to let you inspect or modify the referred subresource.
type SubresourceReference struct {
	Kind        string `json:"kind,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
	Name        string `json:"name,omitempty"`
	APIVersion  string `json:"apiVersion,omitempty"`
	Subresource string `json:"subresource,omitempty"`
}

// ResourceConsumption is an object for specifying average resource consumption of a particular resource.
type ResourceConsumption struct {
	Resource api.ResourceName  `json:"resource,omitempty"`
	Quantity resource.Quantity `json:"quantity,omitempty"`
}

// HorizontalPodAutoscalerSpec is the specification of a horizontal pod autoscaler.
type HorizontalPodAutoscalerSpec struct {
	// ScaleRef is a reference to Scale subresource. HorizontalPodAutoscaler will learn the current resource consumption from its status,
	// and will set the desired number of pods by modyfying its spec.
	ScaleRef *SubresourceReference `json:"scaleRef"`
	// MinCount is the lower limit for the number of pods that can be set by the autoscaler.
	MinCount int `json:"minCount"`
	// MaxCount is the upper limit for the number of pods that can be set by the autoscaler. It cannot be smaller than MinCount.
	MaxCount int `json:"maxCount"`
	// Target is the target average consumption of the given resource that the autoscaler will try to maintain by adjusting the desired number of pods.
	// Currently two types of resources are supported: "cpu" and "memory".
	Target ResourceConsumption `json:"target"`
}

// HorizontalPodAutoscalerStatus contains the current status of a horizontal pod autoscaler
type HorizontalPodAutoscalerStatus struct {
	// CurrentReplicas is the number of replicas of pods managed by this autoscaler.
	CurrentReplicas int `json:"currentReplicas"`

	// DesiredReplicas is the desired number of replicas of pods managed by this autoscaler.
	// The number may be different because pod downscaling is someteimes delayed to keep the number
	// of pods stable.
	DesiredReplicas int `json:"desiredReplicas"`

	// CurrentConsumption is the current average consumption of the given resource that the autoscaler will
	// try to maintain by adjusting the desired number of pods.
	// Two types of resources are supported: "cpu" and "memory".
	CurrentConsumption ResourceConsumption `json:"currentConsumption"`

	// LastScaleTimestamp is the last time the HorizontalPodAutoscaler scaled the number of pods.
	// This is used by the autoscaler to controll how often the number of pods is changed.
	LastScaleTimestamp *util.Time `json:"lastScaleTimestamp,omitempty"`
}

// HorizontalPodAutoscaler represents the configuration of a horizontal pod autoscaler.
type HorizontalPodAutoscaler struct {
	api.TypeMeta   `json:",inline"`
	api.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behaviour of autoscaler.
	Spec HorizontalPodAutoscalerSpec `json:"spec,omitempty"`

	// Status represents the current information about the autoscaler.
	Status HorizontalPodAutoscalerStatus `json:"status,omitempty"`
}

// HorizontalPodAutoscaler is a collection of pod autoscalers.
type HorizontalPodAutoscalerList struct {
	api.TypeMeta `json:",inline"`
	api.ListMeta `json:"metadata,omitempty"`

	Items []HorizontalPodAutoscaler `json:"items"`
}

// A ThirdPartyResource is a generic representation of a resource, it is used by add-ons and plugins to add new resource
// types to the API.  It consists of one or more Versions of the api.
type ThirdPartyResource struct {
	api.TypeMeta   `json:",inline"`
	api.ObjectMeta `json:"metadata,omitempty" description:"standard object metadata"`

	Description string `json:"description,omitempty" description:"The description of this object"`

	Versions []APIVersion `json:"versions,omitempty" description:"The versions for this third party object"`
}

type ThirdPartyResourceList struct {
	api.TypeMeta `json:",inline"`
	api.ListMeta `json:"metadata,omitempty" description:"standard list metadata; see http://docs.k8s.io/api-conventions.md#metadata"`

	Items []ThirdPartyResource `json:"items" description:"items is a list of schema objects"`
}

// An APIVersion represents a single concrete version of an object model.
type APIVersion struct {
	Name     string `json:"name,omitempty" description:"name of this version (e.g. 'v1')"`
	APIGroup string `json:"apiGroup,omitempty" description:"The API group to add this object into, default 'experimental'"`
}

// An internal object, used for versioned storage in etcd.  Not exposed to the end user.
type ThirdPartyResourceData struct {
	api.TypeMeta   `json:",inline"`
	api.ObjectMeta `json:"metadata,omitempty" description:"standard object metadata"`

	Data []byte `json:"name,omitempty" description:"the raw JSON data for this data"`
}
