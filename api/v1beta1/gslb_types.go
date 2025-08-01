package v1beta1

/*
Copyright 2021-2025 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Strategy defines Gslb behavior
// +k8s:openapi-gen=true
type Strategy struct {
	// Load balancing strategy type:(roundRobin|failover)
	Type string `json:"type" validate:"required,oneof=geoip roundRobin failover"`
	// Weight is defined by map region:weight
	Weight map[string]int `json:"weight,omitempty"`
	// Primary Geo Tag. Valid for failover strategy only
	PrimaryGeoTag string `json:"primaryGeoTag,omitempty"`
	// Defines DNS record TTL in seconds
	DNSTtlSeconds int `json:"dnsTtlSeconds,omitempty" validate:"gte=0"`
	// Split brain TXT record expiration in seconds. The field is deprecated and not used.
	SplitBrainThresholdSeconds int `json:"splitBrainThresholdSeconds,omitempty"`
}

// ResourceRef selects a resource defining the GSLB's load balancer and server
// +k8s:openapi-gen=true
type ResourceRef struct {
	corev1.ObjectReference `json:",inline"`
	// LabelSelector of the referenced resource
	metav1.LabelSelector `json:",inline"`
}

// GslbSpec defines the desired state of Gslb
// +k8s:openapi-gen=true
type GslbSpec struct {
	// Gslb-enabled Ingress Spec
	Ingress IngressSpec `json:"ingress,omitempty"`
	// Gslb Strategy spec
	Strategy Strategy `json:"strategy"`
	// ResourceRef spec
	ResourceRef ResourceRef `json:"resourceRef,omitempty"`
}

// LoadBalancer holds the GSLB's load balancer configuration
// +k8s:openapi-gen=true
type LoadBalancer struct {
	// ExposedIPs on the local Load Balancer
	ExposedIPs []string `json:"exposedIps,omitempty"`
}

// Servers holds the GSLB's servers' configuration
// +k8s:openapi-gen=true
type Server struct {
	// Hostname exposed by the GSLB
	Host string `json:"host,omitempty"`
	// Kubernetes Services backing the load balanced application
	Services []*NamespacedName `json:"services,omitempty"`
}

// NamespacedName holds a reference to a k8s resource
// +k8s:openapi-gen=true
type NamespacedName struct {
	// Namespace where the resource can be found
	Namespace string `json:"namespace"`
	// Name of the resource
	Name string `json:"name"`
}

// GslbStatus defines the observed state of Gslb
type GslbStatus struct {
	// Associated Service status
	ServiceHealth map[string]HealthStatus `json:"serviceHealth"`
	// Current Healthy DNS record structure
	HealthyRecords map[string][]string `json:"healthyRecords"`
	// Cluster Geo Tag
	GeoTag string `json:"geoTag"`
	// Comma-separated list of hosts
	Hosts string `json:"hosts,omitempty"`
	// LoadBalancer configuration
	LoadBalancer LoadBalancer `json:"loadBalancer,omitempty"`
	// Servers configuration
	Servers []*Server `json:"servers,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Gslb is the Schema for the gslbs API
// +kubebuilder:printcolumn:name="strategy",type=string,JSONPath=`.spec.strategy.type`
// +kubebuilder:printcolumn:name="geoTag",type=string,JSONPath=`.status.geoTag`
// +kubebuilder:printcolumn:name="hosts",type=string,JSONPath=`.status.hosts`,priority=1
type Gslb struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GslbSpec   `json:"spec,omitempty"`
	Status GslbStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GslbList contains a list of Gslb
type GslbList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gslb `json:"items"`
}

type HealthStatus string

const (
	Healthy   HealthStatus = "Healthy"
	Unhealthy HealthStatus = "Unhealthy"
	NotFound  HealthStatus = "NotFound"
)

func (h HealthStatus) String() string {
	return string(h)
}

func init() {
	SchemeBuilder.Register(&Gslb{}, &GslbList{})
}
