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
	"errors"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-pca-issuer/pkg/api/v1alpha1"
	"github.com/NaverCloudPlatform/ncloud-pca-issuer/pkg/pca"
	cmutil "github.com/cert-manager/cert-manager/pkg/api/util"
	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/client-go/tools/record"
	"time"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	reasonInvalidIssuer  = "InvalidIssuer"
	reasonSignerNotReady = "SignerNotReady"
	reasonCRInvalid      = "CRInvalid"
	reasonCertIssued     = "CertificateIssued"
	reasonCRNotApproved  = "CRNotApproved"
)

// CertificateRequestReconciler reconciles a NcloudPCAIssuer object
type CertificateRequestReconciler struct {
	Log logr.Logger
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder

	Clock                    clock.Clock
	ClusterResourceNamespace string
	CheckApprovedCondition   bool
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificaterequests,verbs=get;list;watch;update
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificaterequests/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;create;update

func (r *CertificateRequestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	log := r.Log.WithValues("certificaterequest", req.NamespacedName)

	// Get the CertificateRequest
	var certificateRequest cmapi.CertificateRequest
	if err := r.Get(ctx, req.NamespacedName, &certificateRequest); err != nil {
		if err := client.IgnoreNotFound(err); err != nil {
			return ctrl.Result{}, fmt.Errorf("unexpected get error: %v", err)
		}
		log.Info("Not found. Ignoring.")
		return ctrl.Result{}, nil
	}

	// Ignore CertificateRequest if issuerRef doesn't match our group
	if certificateRequest.Spec.IssuerRef.Group != v1alpha1.GroupVersion.Group {
		log.Info("CR is for a different Issuer. Ignoring.", "group", certificateRequest.Spec.IssuerRef.Group)
		return ctrl.Result{}, nil
	}

	// Ignore CertificateRequest if it is already Ready
	if cmutil.CertificateRequestHasCondition(&certificateRequest, cmapi.CertificateRequestCondition{
		Type:   cmapi.CertificateRequestConditionReady,
		Status: cmmeta.ConditionTrue,
	}) {
		log.Info("CertificateRequest is Ready. Ignoring.", "cr", req.NamespacedName)
		return ctrl.Result{}, nil
	}

	// Ignore CertificateRequest if it is already Failed
	if cmutil.CertificateRequestHasCondition(&certificateRequest, cmapi.CertificateRequestCondition{
		Type:   cmapi.CertificateRequestConditionReady,
		Status: cmmeta.ConditionFalse,
		Reason: cmapi.CertificateRequestReasonFailed,
	}) {
		log.Info("CertificateRequest is Failed. Ignoring.", "cr", req.NamespacedName)
		return ctrl.Result{}, nil
	}
	// Ignore CertificateRequest if it already has a Denied Ready Reason
	if cmutil.CertificateRequestHasCondition(&certificateRequest, cmapi.CertificateRequestCondition{
		Type:   cmapi.CertificateRequestConditionReady,
		Status: cmmeta.ConditionFalse,
		Reason: cmapi.CertificateRequestReasonDenied,
	}) {
		log.Info("CertificateRequest is Denied. Ignoring.", "cr", req.NamespacedName)
		return ctrl.Result{}, nil
	}

	// We now have a CertificateRequest that belongs to us so we are responsible
	// for updating its Ready condition.
	setReadyCondition := func(status cmmeta.ConditionStatus, reason, message string) {
		cmutil.SetCertificateRequestCondition(
			&certificateRequest,
			cmapi.CertificateRequestConditionReady,
			status,
			reason,
			message,
		)
	}

	// Always attempt to update the Ready condition
	defer func() {
		if err != nil {
			setReadyCondition(cmmeta.ConditionFalse, cmapi.CertificateRequestReasonPending, err.Error())
		}
		if updateErr := r.Status().Update(ctx, &certificateRequest); updateErr != nil {
			err = utilerrors.NewAggregate([]error{err, updateErr})
			result = ctrl.Result{}
		}
	}()

	// If CertificateRequest has been denied, mark the CertificateRequest as
	// Ready=Denied and set FailureTime if not already.
	if cmutil.CertificateRequestIsDenied(&certificateRequest) {
		log.Info("CertificateRequest has been denied yet. Marking as failed.")

		if certificateRequest.Status.FailureTime == nil {
			nowTime := metav1.NewTime(r.Clock.Now())
			certificateRequest.Status.FailureTime = &nowTime
		}

		message := "The CertificateRequest was denied by an approval controller"
		setReadyCondition(cmmeta.ConditionFalse, cmapi.CertificateRequestReasonDenied, message)
		return ctrl.Result{}, nil
	}

	if r.CheckApprovedCondition {
		// If CertificateRequest has not been approved, exit early.
		log.Info("Checking whether CR has been approved", "cr", req.NamespacedName)
		if !cmutil.CertificateRequestIsApproved(&certificateRequest) {
			msg := "certificate request is not approved yet"
			log.Info(msg, "cr", req.NamespacedName)
			r.Recorder.Event(&certificateRequest, eventTypeWarning, reasonCRNotApproved, msg)
			return ctrl.Result{}, nil
		}
	}

	// Add a Ready condition if one does not already exist
	if ready := cmutil.GetCertificateRequestCondition(&certificateRequest, cmapi.CertificateRequestConditionReady); ready == nil {
		log.Info("Initialising Ready condition")
		setReadyCondition(cmmeta.ConditionFalse, cmapi.CertificateRequestReasonPending, "Initialising")
		// re-reconcile
		return ctrl.Result{}, nil
	}

	// Ignore but log an error if the issuerRef.Kind is unrecognised
	issuerGVK := v1alpha1.GroupVersion.WithKind(certificateRequest.Spec.IssuerRef.Kind)
	issuerRO, err := r.Scheme.New(issuerGVK)
	var spec *v1alpha1.NcloudPCAIssuerSpec
	var ns string
	if err != nil {
		log.Error(err, "unknown issuer kind", "kind", certificateRequest.Spec.IssuerRef.Kind)
		msg := "The issuer kind " + certificateRequest.Spec.IssuerRef.Kind + " is invalid"
		setReadyCondition(cmmeta.ConditionFalse, cmapi.CertificateRequestReasonFailed, msg)
		r.Recorder.Event(&certificateRequest, eventTypeWarning, reasonInvalidIssuer, msg)
		return ctrl.Result{}, nil
	}

	//issuer := issuerRO.(client.Object)
	//// Create a Namespaced name for Issuer and a non-Namespaced name for ClusterIssuer
	//issuerName := types.NamespacedName{
	//	Name: certificateRequest.Spec.IssuerRef.Name,
	//}
	//var secretNamespace string

	switch t := issuerRO.(type) {
	case *v1alpha1.NcloudPCAIssuer:
		err := r.Client.Get(ctx, types.NamespacedName{Namespace: req.NamespacedName.Namespace,
			Name: certificateRequest.Spec.IssuerRef.Name}, t)
		if err != nil {
			if client.IgnoreNotFound(err) == nil {
				log.Info("issuer not found", "issuer", certificateRequest.Spec.IssuerRef.Name, "namespace", req.NamespacedName.Namespace)
				msg := fmt.Sprintf("The issuer %s was not found", certificateRequest.Spec.IssuerRef.Name)
				setReadyCondition(cmmeta.ConditionFalse, cmapi.CertificateRequestReasonFailed, msg)
				r.Recorder.Event(&certificateRequest, eventTypeWarning, reasonInvalidIssuer, msg)
				return ctrl.Result{}, nil
			}
			return ctrl.Result{}, err
		}
		spec = &t.Spec
		ns = req.NamespacedName.Namespace

	case *v1alpha1.NcloudPCAClusterIssuer:
		err := r.Client.Get(ctx, types.NamespacedName{
			Name: certificateRequest.Spec.IssuerRef.Name,
		}, t)
		if err != nil {
			if client.IgnoreNotFound(err) == nil {
				log.Info("ClusterIssuer not found", "CLusterIssuer", certificateRequest.Spec.IssuerRef.Name)
				msg := fmt.Sprintf("The ClusterIssuer %s was not found", certificateRequest.Spec.IssuerRef.Name)
				setReadyCondition(cmmeta.ConditionFalse, cmapi.CertificateRequestReasonFailed, msg)
				r.Recorder.Event(&certificateRequest, eventTypeWarning, reasonInvalidIssuer, msg)
				return ctrl.Result{}, nil
			}
			return ctrl.Result{}, err
		}
		spec = &t.Spec
		ns = r.ClusterResourceNamespace
	default:
		log.Error(err, "unknown issuer type", "object", t)
		msg := "Unknown issuer type"
		setReadyCondition(cmmeta.ConditionFalse, cmapi.CertificateRequestReasonFailed, msg)
		r.Recorder.Event(&certificateRequest, eventTypeWarning, reasonInvalidIssuer, msg)
		return ctrl.Result{}, nil
	}

	signer, err := pca.NewSigner(ctx, spec, ns, r.Client)
	if err != nil {
		log.Error(err, "could not create signer", "cr", req.NamespacedName)
		msg := fmt.Sprintf("Could not create signer, check ca %s is ready", spec.CaTag)
		r.Recorder.Event(&certificateRequest, eventTypeWarning, reasonSignerNotReady, msg)
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// Check for obvious errors, e.g. missing duration, malformed certificate request
	if err := sanitiseCertificateRequestSpec(&certificateRequest.Spec); err != nil {
		log.Error(err, "certificate request has issues", "cr", req.NamespacedName)
		msg := "certificate request has issues: " + err.Error()
		setReadyCondition(cmmeta.ConditionFalse, cmapi.CertificateRequestReasonFailed, msg)
		r.Recorder.Event(&certificateRequest, eventTypeWarning, reasonCRInvalid, msg)
		return ctrl.Result{}, nil
	}

	// Sign certificate
	cert, ca, err := signer.Sign(certificateRequest.Spec.Request, certificateRequest.Spec.Duration.Duration)
	if err != nil {
		return ctrl.Result{}, err
	}
	certificateRequest.Status.CA = ca
	certificateRequest.Status.Certificate = cert
	msg := "Certificate Issued"
	setReadyCondition(cmmeta.ConditionTrue, cmapi.CertificateRequestReasonIssued, msg)
	r.Recorder.Event(&certificateRequest, eventTypeNormal, reasonCertIssued, msg)
	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *CertificateRequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cmapi.CertificateRequest{}).
		Complete(r)
}

func sanitiseCertificateRequestSpec(spec *cmapi.CertificateRequestSpec) error {
	// Ensure there is a duration
	if spec.Duration == nil {
		spec.Duration = &metav1.Duration{
			Duration: cmapi.DefaultCertificateDuration,
		}
	}
	// Very short durations should be increased
	if spec.Duration.Duration < cmapi.MinimumCertificateDuration {
		spec.Duration = &metav1.Duration{
			Duration: cmapi.MinimumCertificateDuration,
		}
	}
	if len(spec.Request) == 0 {
		return errors.New("certificate request is empty")
	}
	return nil
}
