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
package pca

import (
	"bytes"
	"context"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/NaverCloudPlatform/ncloud-pca-issuer/pkg/api/v1alpha1"
	"github.com/NaverCloudPlatform/ncloud-pca-issuer/pkg/privateca"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	defaultAccessKeyField = "NCLOUD_ACCESS_KEY"
	defaultSecretKeyField = "NCLOUD_SECRET_KEY"
)

// A Signer is an abstraction of a certificate authority.
type Signer interface {
	// Sign signs a CSR and returns the leaf+intermediate chain (cert) and the
	// root-most CA (ca) in PEM format.
	Sign(csr []byte, expiry time.Duration) (cert []byte, ca []byte, err error)
}

type pcaSigner struct {
	spec      *v1alpha1.NcloudPCAIssuerSpec
	namespace string

	gatewayURL string
	client     client.Client
	ctx        context.Context
}

// NewSigner returns a Signer after validating credentials and the CA tag by
// performing a read of the CA from the NCloud API.
func NewSigner(ctx context.Context, spec *v1alpha1.NcloudPCAIssuerSpec, namespace string, k8sClient client.Client) (Signer, error) {
	p, err := newSignerNoSelftest(ctx, spec, namespace, k8sClient)
	if err != nil {
		return nil, err
	}

	pcaClient, err := p.newPcaClient()
	if err != nil {
		return nil, err
	}
	if _, err := pcaClient.V1Api.CaCaTagGet(ctx, &spec.CaTag); err != nil {
		return nil, fmt.Errorf("get CA info (caTag=%s): %w", spec.CaTag, err)
	}
	return p, nil
}

func newSignerNoSelftest(ctx context.Context, spec *v1alpha1.NcloudPCAIssuerSpec, namespace string, k8sClient client.Client) (*pcaSigner, error) {
	if spec.CaTag == "" {
		return nil, errors.New("must specify a CaTag")
	}
	gw, err := APIGatewayURLFor(spec.Region, spec.APIGatewayURL)
	if err != nil {
		return nil, err
	}
	return &pcaSigner{
		spec:       spec,
		namespace:  namespace,
		gatewayURL: gw,
		client:     k8sClient,
		ctx:        ctx,
	}, nil
}

func (p *pcaSigner) Sign(csr []byte, expiry time.Duration) ([]byte, []byte, error) {
	pcaClient, err := p.newPcaClient()
	if err != nil {
		return nil, nil, err
	}
	period := fmt.Sprintf("%d", int(expiry.Hours()/24))
	csrPem := strings.TrimSpace(string(csr))
	csrReq := &privateca.SignCsr{
		CsrPem: &csrPem,
		Period: &period,
	}

	csrResp, err := pcaClient.V1Api.CaCaTagCertSignPost(p.ctx, csrReq, &p.spec.CaTag, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("sign CSR: %w", err)
	}
	return extractCertAndCA(csrResp.Data)
}

func (p *pcaSigner) newPcaClient() (*privateca.APIClient, error) {
	os.Setenv("NCLOUD_API_GW", p.gatewayURL)

	ref := p.spec.CredentialsRef
	if ref.Name == "" {
		return privateca.NewAPIClient(privateca.NewConfiguration()), nil
	}

	ns := ref.Namespace
	if ns == "" {
		ns = p.namespace
	}
	secret := &core.Secret{}
	if err := p.client.Get(p.ctx, types.NamespacedName{Namespace: ns, Name: ref.Name}, secret); err != nil {
		return nil, fmt.Errorf("failed to retrieve secret %s/%s: %w", ns, ref.Name, err)
	}

	accessField := ref.AccessKeyField
	if accessField == "" {
		accessField = defaultAccessKeyField
	}
	secretField := ref.SecretKeyField
	if secretField == "" {
		secretField = defaultSecretKeyField
	}

	accessKey, ok := secret.Data[accessField]
	if !ok {
		return nil, fmt.Errorf("secret %s/%s missing key %q", ns, ref.Name, accessField)
	}
	secretKey, ok := secret.Data[secretField]
	if !ok {
		return nil, fmt.Errorf("secret %s/%s missing key %q", ns, ref.Name, secretField)
	}

	apiKey := &ncloud.APIKey{
		AccessKey: string(accessKey),
		SecretKey: string(secretKey),
	}
	return privateca.NewAPIClient(privateca.NewConfiguration(apiKey)), nil
}

func extractCertAndCA(data *privateca.SignCsrResponseData) ([]byte, []byte, error) {
	if data == nil || data.Certificate == nil {
		return nil, nil, errors.New("extractCertAndCA: response missing certificate")
	}
	leafPEM := strings.TrimSpace(*data.Certificate)

	var chainPEMs []string
	if data.CaChain != nil {
		for _, entry := range *data.CaChain {
			if strings.TrimSpace(entry) != "" {
				chainPEMs = append(chainPEMs, entry)
			}
		}
	}
	if len(chainPEMs) == 0 && data.Issuer != nil && strings.TrimSpace(*data.Issuer) != "" {
		chainPEMs = []string{*data.Issuer}
	}

	parsed, err := parsePEMChain(chainPEMs)
	if err != nil {
		return nil, nil, fmt.Errorf("parse CA chain: %w", err)
	}
	intermediates, root := splitRoot(parsed)

	var certBuf bytes.Buffer
	certBuf.WriteString(leafPEM)
	certBuf.WriteByte('\n')
	for _, c := range intermediates {
		if err := pem.Encode(&certBuf, &pem.Block{Type: "CERTIFICATE", Bytes: c.Raw}); err != nil {
			return nil, nil, fmt.Errorf("encode intermediate: %w", err)
		}
	}

	var caBuf bytes.Buffer
	switch {
	case root != nil:
		if err := pem.Encode(&caBuf, &pem.Block{Type: "CERTIFICATE", Bytes: root.Raw}); err != nil {
			return nil, nil, fmt.Errorf("encode root: %w", err)
		}
	case len(intermediates) > 0:
		top := intermediates[len(intermediates)-1]
		if err := pem.Encode(&caBuf, &pem.Block{Type: "CERTIFICATE", Bytes: top.Raw}); err != nil {
			return nil, nil, fmt.Errorf("encode top intermediate: %w", err)
		}
	}
	return certBuf.Bytes(), caBuf.Bytes(), nil
}
