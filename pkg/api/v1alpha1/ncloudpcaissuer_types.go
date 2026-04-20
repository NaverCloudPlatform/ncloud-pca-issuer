/*
Copyright 2022 Naver Cloud Platform.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NcloudPCAIssuerSpec defines the desired state of a Ncloud Private CA Issuer.
type NcloudPCAIssuerSpec struct {
	// CaTag is the Private CA tag ID shown on the NCloud console.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	CaTag string `json:"caTag"`

	// Region selects the NCloud realm. Determines the API gateway URL.
	// +kubebuilder:validation:Enum=public;gov;fin
	// +kubebuilder:default=public
	Region string `json:"region,omitempty"`

	// APIGatewayURL overrides the region-derived API gateway URL.
	// Advanced use only (e.g. private gateway, testing). Leave empty to use Region.
	// +optional
	APIGatewayURL string `json:"apiGatewayUrl,omitempty"`

	// CredentialsRef references the Secret containing NCloud API credentials.
	// Omit to authenticate using the NKS node's ServerRole (IAM role) — see
	// https://guide.ncloud-docs.com/docs/kubernetes-node-iam-role.
	// +optional
	CredentialsRef NcloudCredentialsRef `json:"credentialsRef,omitempty"`
}

// NcloudCredentialsRef points to the Secret that stores the NCloud API access and secret keys.
type NcloudCredentialsRef struct {
	// Name of the Secret. Required when CredentialsRef is set.
	// +optional
	Name string `json:"name,omitempty"`

	// Namespace of the Secret. Required for ClusterIssuer; ignored for namespaced Issuer.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// AccessKeyField is the Secret data key holding the NCloud access key.
	// +kubebuilder:default=NCLOUD_ACCESS_KEY
	AccessKeyField string `json:"accessKeyField,omitempty"`

	// SecretKeyField is the Secret data key holding the NCloud secret key.
	// +kubebuilder:default=NCLOUD_SECRET_KEY
	SecretKeyField string `json:"secretKeyField,omitempty"`
}

// NcloudPCAIssuerStatus defines the observed state of NcloudPCAIssuer.
type NcloudPCAIssuerStatus struct {
	Conditions []NcloudPCAIssuerCondition `json:"conditions,omitempty"`
}

// NcloudPCAIssuerConditionType represents an Issuer condition value.
type NcloudPCAIssuerConditionType string

const (
	// IssuerConditionReady represents the fact that a given Issuer condition
	// is in ready state and able to issue certificates.
	// If the `status` of this condition is `False`, CertificateRequest controllers
	// should prevent attempts to sign certificates.
	IssuerConditionReady NcloudPCAIssuerConditionType = "Ready"
)

// ConditionStatus represents a condition's status.
// +kubebuilder:validation:Enum=True;False;Unknown
type ConditionStatus string

const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
//+kubebuilder:printcolumn:name="Region",type=string,JSONPath=`.spec.region`
//+kubebuilder:printcolumn:name="CaTag",type=string,JSONPath=`.spec.caTag`
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// NcloudPCAIssuer is the Schema for the ncloudpcaissuers API
type NcloudPCAIssuer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NcloudPCAIssuerSpec   `json:"spec,omitempty"`
	Status NcloudPCAIssuerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NcloudPCAIssuerList contains a list of NcloudPCAIssuer
type NcloudPCAIssuerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NcloudPCAIssuer `json:"items"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
//+kubebuilder:printcolumn:name="Region",type=string,JSONPath=`.spec.region`
//+kubebuilder:printcolumn:name="CaTag",type=string,JSONPath=`.spec.caTag`
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// NcloudPCAClusterIssuer is the Schema for the ncloudpcaclusterissuers API
type NcloudPCAClusterIssuer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NcloudPCAIssuerSpec   `json:"spec,omitempty"`
	Status NcloudPCAIssuerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NcloudPCAClusterIssuerList contains a list of NcloudPCAClusterIssuer
type NcloudPCAClusterIssuerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NcloudPCAClusterIssuer `json:"items"`
}

// NcloudPCAIssuerCondition contains condition information for a PCA Issuer.
type NcloudPCAIssuerCondition struct {
	// Type of the condition, currently ('Ready').
	Type NcloudPCAIssuerConditionType `json:"type"`

	// Status of the condition, one of ('True', 'False', 'Unknown').
	Status ConditionStatus `json:"status"`

	// LastTransitionTime is the timestamp corresponding to the last status
	// change of this condition.
	// +optional
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`

	// Reason is a brief machine readable explanation for the condition's last
	// transition.
	// +optional
	Reason string `json:"reason,omitempty"`

	// Message is a human readable description of the details of the last
	// transition, complementing reason.
	// +optional
	Message string `json:"message,omitempty"`
}

func init() {
	SchemeBuilder.Register(&NcloudPCAIssuer{}, &NcloudPCAIssuerList{})
	SchemeBuilder.Register(&NcloudPCAClusterIssuer{}, &NcloudPCAClusterIssuerList{})
}
