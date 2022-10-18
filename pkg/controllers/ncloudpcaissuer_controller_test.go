package controllers

import (
	privatecaissuerv1alpha1 "github.com/NaverCloudPlatform/ncloud-pca-issuer/pkg/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetReadyCondition(t *testing.T) {
	tests := []struct {
		name                 string
		inputStatus          *privatecaissuerv1alpha1.NcloudPCAIssuerStatus
		inputConditionStatus privatecaissuerv1alpha1.ConditionStatus
		inputReason          string
		inputMessage         string
		expectedStatus       *privatecaissuerv1alpha1.NcloudPCAIssuerStatus
	}{
		{
			name:                 "Status with nil condition should be set",
			inputStatus:          &privatecaissuerv1alpha1.NcloudPCAIssuerStatus{Conditions: nil},
			inputConditionStatus: privatecaissuerv1alpha1.ConditionTrue,
			inputReason:          "Test Ready Reason",
			inputMessage:         "Test Ready Message",
			expectedStatus: &privatecaissuerv1alpha1.NcloudPCAIssuerStatus{
				Conditions: []privatecaissuerv1alpha1.NcloudPCAIssuerCondition{{
					Type:               privatecaissuerv1alpha1.IssuerConditionReady,
					Status:             privatecaissuerv1alpha1.ConditionTrue,
					LastTransitionTime: nil,
					Reason:             "Test Ready Reason",
					Message:            "Test Ready Message",
				}},
			},
		},
		{
			name: "Status can transition from Ready to Not Ready",
			inputStatus: &privatecaissuerv1alpha1.NcloudPCAIssuerStatus{
				Conditions: []privatecaissuerv1alpha1.NcloudPCAIssuerCondition{
					{
						Type:               privatecaissuerv1alpha1.IssuerConditionReady,
						Status:             privatecaissuerv1alpha1.ConditionTrue,
						LastTransitionTime: nil,
						Reason:             "I was Ready before",
						Message:            "Test Ready Message",
					},
				},
			},
			inputConditionStatus: privatecaissuerv1alpha1.ConditionFalse,
			inputReason:          "I'm not ready now reason",
			inputMessage:         "I'm not ready now message",
			expectedStatus: &privatecaissuerv1alpha1.NcloudPCAIssuerStatus{
				Conditions: []privatecaissuerv1alpha1.NcloudPCAIssuerCondition{{
					Type:               privatecaissuerv1alpha1.IssuerConditionReady,
					Status:             privatecaissuerv1alpha1.ConditionFalse,
					LastTransitionTime: nil,
					Reason:             "I'm not ready now reason",
					Message:            "I'm not ready now message",
				}},
			},
		},
		{
			name: "Status can transition from Not Ready to Ready",
			inputStatus: &privatecaissuerv1alpha1.NcloudPCAIssuerStatus{
				Conditions: []privatecaissuerv1alpha1.NcloudPCAIssuerCondition{
					{
						Type:               privatecaissuerv1alpha1.IssuerConditionReady,
						Status:             privatecaissuerv1alpha1.ConditionFalse,
						LastTransitionTime: nil,
						Reason:             "I was not ready before",
						Message:            "Test Ready Message",
					},
				},
			},
			inputConditionStatus: privatecaissuerv1alpha1.ConditionTrue,
			inputReason:          "I'm ready now reason",
			inputMessage:         "I'm ready now message",
			expectedStatus: &privatecaissuerv1alpha1.NcloudPCAIssuerStatus{
				Conditions: []privatecaissuerv1alpha1.NcloudPCAIssuerCondition{{
					Type:               privatecaissuerv1alpha1.IssuerConditionReady,
					Status:             privatecaissuerv1alpha1.ConditionTrue,
					LastTransitionTime: nil,
					Reason:             "I'm ready now reason",
					Message:            "I'm ready now message",
				}},
			},
		},
	}
	for _, tt := range tests {
		status := tt.inputStatus.DeepCopy()
		setReadyCondition(status, tt.inputConditionStatus, tt.inputReason, tt.inputMessage)
		// ignore time.now
		for i := range status.Conditions {
			status.Conditions[i].LastTransitionTime = nil
		}
		assert.Equal(t, tt.expectedStatus, status, "%s failed", tt.name)
	}
}
