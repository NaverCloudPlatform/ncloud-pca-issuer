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

package controllers

import (
	"context"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-pca-issuer/pkg/api/v1alpha1"
	"github.com/NaverCloudPlatform/ncloud-pca-issuer/pkg/pca"
	"github.com/go-logr/logr"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	eventTypeWarning = "Warning"
	eventTypeNormal  = "Normal"

	reasonPCAClientOK         = "PrivateCaClientOK"
	reasonIssuerMisconfigured = "IssuerMisconfigured"
)

// NcloudPCAIssuerReconciler reconciles a NcloudPCAIssuer object
type NcloudPCAIssuerReconciler struct {
	// NcloudPCAIssuer or NcloudPCAClusterIssuer
	Kind string

	Log logr.Logger
	client.Client
	Scheme                   *runtime.Scheme
	Recorder                 record.EventRecorder
	ClusterResourceNamespace string
}

// +kubebuilder:rbac:groups=privateca-issuer.ncloud.com,resources=ncloudpcaissuers;ncloudpcaclusterissuers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=privateca-issuer.ncloud.com,resources=ncloudpcaissuers/status;ncloudpcaclusterissuers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;create;update

func (r *NcloudPCAIssuerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	log := r.Log.WithValues(r.Kind, req.NamespacedName)
	issuer, err := r.getIssuer()
	if err != nil {
		log.Error(err, "ignore invalid issuer kind")
		return ctrl.Result{}, nil
	}

	if err := r.Get(ctx, req.NamespacedName, issuer); err != nil {
		if err := client.IgnoreNotFound(err); err != nil {
			log.Error(err, "failed to get Issuer resource")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	spec, status, err := getIssuerSpecStatus(issuer)
	if err != nil {
		log.Error(err, "failed to get issuer spec and status")
		return ctrl.Result{}, err
	}

	defer func() {
		if err != nil {
			setReadyCondition(status, v1alpha1.ConditionFalse, "issuer failed to reconcile", err.Error())
		}
		// If the Issuer is deleted mid-reconcile, ignore it
		if updateErr := client.IgnoreNotFound(r.Status().Update(ctx, issuer)); updateErr != nil {
			log.Info("Couldn't update ready condition", "err", err)
			result = ctrl.Result{}
		}
	}()

	ns := req.NamespacedName.Namespace
	if len(ns) == 0 {
		ns = viper.GetString("cluster-resource-namespace")
	}

	_, err = pca.NewSigner(ctx, spec, ns, r.Client)
	if err != nil {
		log.Info("Issuer config failed", "info", err.Error())
		setReadyCondition(status, v1alpha1.ConditionFalse, "issuer config failed", err.Error())
		r.Recorder.Event(issuer, eventTypeWarning, reasonIssuerMisconfigured, err.Error())
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	log.Info("reconcile issuer success", "kind", issuer.GetObjectKind())
	msg := "Successfully init PCA Client"
	setReadyCondition(status, v1alpha1.ConditionTrue, "OK", msg)
	r.Recorder.Event(issuer, eventTypeNormal, reasonPCAClientOK, msg)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NcloudPCAIssuerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	issuer, err := r.getIssuer()
	if err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(issuer).
		Complete(r)
}

// convert a k8s.io/apimachinery/pkg/runtime.Object into a sigs.k8s.io/controller-runtime/pkg/client.Object
func (r *NcloudPCAIssuerReconciler) getIssuer() (client.Object, error) {
	issuer, err := r.Scheme.New(v1alpha1.GroupVersion.WithKind(r.Kind))
	if err != nil {
		return nil, err
	}
	switch t := issuer.(type) {
	case *v1alpha1.NcloudPCAIssuer:
		return t, nil
	case *v1alpha1.NcloudPCAClusterIssuer:
		return t, nil
	default:
		return nil, fmt.Errorf("unsupported kind %s", r.Kind)
	}
}

func getIssuerSpecStatus(object client.Object) (*v1alpha1.NcloudPCAIssuerSpec, *v1alpha1.NcloudPCAIssuerStatus, error) {
	switch t := object.(type) {
	case *v1alpha1.NcloudPCAIssuer:
		return &t.Spec, &t.Status, nil
	case *v1alpha1.NcloudPCAClusterIssuer:
		return &t.Spec, &t.Status, nil
	default:
		return nil, nil, fmt.Errorf("unexpected kind %T", t)
	}
}

func setReadyCondition(status *v1alpha1.NcloudPCAIssuerStatus, conditionStatus v1alpha1.ConditionStatus, reason, message string) {
	var ready *v1alpha1.NcloudPCAIssuerCondition
	for _, c := range status.Conditions {
		if c.Type == v1alpha1.IssuerConditionReady {
			ready = &c
			break
		}
	}
	if ready == nil {
		ready = &v1alpha1.NcloudPCAIssuerCondition{Type: v1alpha1.IssuerConditionReady}
	}
	if ready.Status != conditionStatus {
		ready.Status = conditionStatus
		now := metav1.Now()
		ready.LastTransitionTime = &now
	}
	ready.Reason = reason
	ready.Message = message

	for i, c := range status.Conditions {
		if c.Type == v1alpha1.IssuerConditionReady {
			status.Conditions[i] = *ready
			return
		}
	}

	status.Conditions = append(status.Conditions, *ready)
}
