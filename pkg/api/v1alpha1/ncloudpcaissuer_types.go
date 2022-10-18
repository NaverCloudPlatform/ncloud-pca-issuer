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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NcloudPCAIssuerSpec defines the desired state of NcloudPCAIssuer
type NcloudPCAIssuerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// CaTag is the id of the CA to issue certificates from
	CaTag string `json:"caTag,omitempty"`

	// KeyType is the Algorithm type of the CA Public key
	KeyType string `json:"keyType,omitempty"`

	// KeyBits is the bit length of the CA Public key
	KeyBits string `json:"keyBits,omitempty"`

	// NcloudApiGw is the URL for NCLOUD API Gateway
	NcloudApiGw string `json:"ncloudApiGw,omitempty"`

	// Needs to be specified if you want to authorize with AWS using an access and secret key
	// +optional
	SecretRef NcloudCredentialsSecretReference `json:"secretRef,omitempty"`
}

//AWSCredentialsSecretReference defines the secret used by the issuer
type NcloudCredentialsSecretReference struct {
	v1.SecretReference `json:""`
	// Specifies the secret key where the AWS Access Key ID exists
	// +optional
	AccessKeyIDSelector v1.SecretKeySelector `json:"accessKeyIDSelector,omitempty"`
	// Specifies the secret key where the AWS Secret Access Key exists
	// +optional
	SecretAccessKeySelector v1.SecretKeySelector `json:"secretAccessKeySelector,omitempty"`
}

// NcloudPCAIssuerStatus defines the observed state of NcloudPCAIssuer
type NcloudPCAIssuerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
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

// These are valid condition statuses. "ConditionTrue" means a resource is in
// the condition; "ConditionFalse" means a resource is not in the condition;
// "ConditionUnknown" means kubernetes can't decide if a resource is in the
// condition or not. In the future, we could add other intermediate
// conditions, e.g. ConditionDegraded.
const (
	// ConditionTrue represents the fact that a given condition is true
	ConditionTrue ConditionStatus = "True"

	// ConditionFalse represents the fact that a given condition is false
	ConditionFalse ConditionStatus = "False"

	// ConditionUnknown represents the fact that a given condition is unknown
	ConditionUnknown ConditionStatus = "Unknown"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

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

// IssuerCondition contains condition information for a PCA Issuer.
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
