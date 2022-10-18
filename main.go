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

package main

import (
	"errors"
	"flag"
	"fmt"
	privatecaissuerv1alpha1 "github.com/NaverCloudPlatform/ncloud-pca-issuer/pkg/api/v1alpha1"
	"github.com/NaverCloudPlatform/ncloud-pca-issuer/pkg/controllers"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/clock"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	//+kubebuilder:scaffold:imports

	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
)

const inClusterNamespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(privatecaissuerv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme

	utilruntime.Must(cmapi.AddToScheme(scheme))
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var clusterResourceNamespace string
	var printVersion bool
	var disableApprovedCheck bool

	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&clusterResourceNamespace, "cluster-resource-namespace", "cert-manager", "The namespace for secrets in which cluster-scoped resources are found.")
	flag.BoolVar(&printVersion, "version", false, "Print version to stdout and exit")
	flag.BoolVar(&disableApprovedCheck, "disable-approved-check", false,
		"Disables waiting for CertificateRequests to have an approved condition before signing.")

	opts := zap.Options{
		Development: false,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	if clusterResourceNamespace == "" {
		var err error
		clusterResourceNamespace, err = getInClusterNamespace()
		if err != nil {
			if errors.Is(err, errNotInCluster) {
				setupLog.Error(err, "please supply --cluster-resource-namespace")
			} else {
				setupLog.Error(err, "unexpected error while getting in-cluster Namespace")
			}
			os.Exit(1)
		}
	}

	setupLog.Info(
		"starting",
		//"version", version.Version,
		"enable-leader-election", enableLeaderElection,
		"metrics-addr", metricsAddr,
		"cluster-resource-namespace", clusterResourceNamespace,
	)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		//HealthProbeBindAddress: probeAddr,
		LeaderElection:   enableLeaderElection,
		LeaderElectionID: "c46c0d21.ncloud.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.NcloudPCAIssuerReconciler{
		Kind:                     "NcloudPCAIssuer",
		Client:                   mgr.GetClient(),
		Scheme:                   mgr.GetScheme(),
		Recorder:                 mgr.GetEventRecorderFor("ncloud-pca-issuer-issuer-controller"),
		Log:                      ctrl.Log.WithName("controllers").WithName("NcloudPCAIssuer"),
		ClusterResourceNamespace: clusterResourceNamespace,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "NcloudPCAIssuer")
		os.Exit(1)
	}

	if err = (&controllers.NcloudPCAIssuerReconciler{
		Kind:                     "NcloudPCAClusterIssuer",
		Client:                   mgr.GetClient(),
		Scheme:                   mgr.GetScheme(),
		Recorder:                 mgr.GetEventRecorderFor("ncloud-pca-issuer-clusterissuer-controller"),
		Log:                      ctrl.Log.WithName("controllers").WithName("NcloudPCAClusterIssuer"),
		ClusterResourceNamespace: clusterResourceNamespace,
		//HealthCheckerBuilder:     signer.ExampleHealthCheckerFromIssuerAndSecretData,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "NcloudPCAClusterIssuer")
		os.Exit(1)
	}

	if err = (&controllers.CertificateRequestReconciler{
		Client:                   mgr.GetClient(),
		Scheme:                   mgr.GetScheme(),
		Recorder:                 mgr.GetEventRecorderFor("ncloud-pca-issuer-certificaterequest-controller"),
		ClusterResourceNamespace: clusterResourceNamespace,
		Log:                      ctrl.Log.WithName("controllers").WithName("CertificateRequest"),
		CheckApprovedCondition:   !disableApprovedCheck,
		Clock:                    clock.RealClock{},
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CertificateRequest")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

var errNotInCluster = errors.New("not running in-cluster")

// Copied from controller-runtime/pkg/leaderelection
func getInClusterNamespace() (string, error) {
	// Check whether the namespace file exists.
	// If not, we are not running in cluster so can't guess the namespace.
	_, err := os.Stat(inClusterNamespacePath)
	if os.IsNotExist(err) {
		return "", errNotInCluster
	} else if err != nil {
		return "", fmt.Errorf("error checking namespace file: %w", err)
	}

	// Load the namespace file and return its content
	namespace, err := ioutil.ReadFile(inClusterNamespacePath)
	if err != nil {
		return "", fmt.Errorf("error reading namespace file: %w", err)
	}
	return string(namespace), nil
}
