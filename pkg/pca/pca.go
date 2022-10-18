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
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/NaverCloudPlatform/ncloud-pca-issuer/pkg/api/v1alpha1"
	"github.com/NaverCloudPlatform/ncloud-pca-issuer/pkg/privateca"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
	"time"
)

// A Signer is an abstraction of a certificate authority
type Signer interface {
	// Sign signs a CSR and returns a cert and chain
	Sign(csr []byte, expiry time.Duration) (cert []byte, ca []byte, err error)
}

type pcaSigner struct {
	// caTag is id for ca
	caTag string
	// ncloudApiGw is url for ncloud api gw
	ncloudApiGw string
	// spec is a reference to the issuer Spec
	spec *v1alpha1.NcloudPCAIssuerSpec
	// namespace is the namespace to look for secrets in
	namespace string

	cfg    *ncloud.Configuration
	client client.Client
	ctx    context.Context
}

func (p *pcaSigner) Sign(csr []byte, expiry time.Duration) (cert []byte, ca []byte, err error) {
	pcaClient, err := p.creatPcaClient()
	if err != nil {
		return nil, nil, err
	}
	period := fmt.Sprintf("%d", int(expiry.Hours()/24))
	csrPem := strings.TrimSpace(string(csr))
	csrReq := &privateca.SignCsr{
		CsrPem:  &csrPem,
		Period:  &period,
		KeyType: &p.spec.KeyType,
		KeyBits: &p.spec.KeyBits,
	}

	csrResp, err := pcaClient.V1Api.CaCaTagCertSignPost(context.Background(), csrReq, &p.spec.CaTag, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	return extractCertAndCA(csrResp.Data)
}

func NewSigner(ctx context.Context, spec *v1alpha1.NcloudPCAIssuerSpec, namespace string, client client.Client) (Signer, error) {

	p, err := newSignerNoSelftest(ctx, spec, client, namespace)
	if err != nil {
		return p, err
	}
	pcaClient, err := p.creatPcaClient()
	if err != nil {
		return p, err
	}
	ca, err := pcaClient.V1Api.CaCaTagGet(ctx, &spec.CaTag)
	if err != nil {

	}

	pemBytes := []byte(ncloud.StringValue(ca.Data.CaCertInfo.CertPem))
	block, _ := pem.Decode(pemBytes)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	var keyType, keyBits string
	switch pub := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		keyType = "RSA"
		keyBits = strconv.Itoa(pub.N.BitLen())
	case *ecdsa.PublicKey:
		keyType = "ECDSA"
		keyBits = strconv.Itoa(pub.Params().N.BitLen())
	default:
		panic("unknown type of public key")
	}

	if p.spec.KeyType == "" {
		p.spec.KeyType = keyType
	}
	if p.spec.KeyBits == "" {
		p.spec.KeyBits = keyBits
	}
	if spec.KeyType != keyType || spec.KeyBits != keyBits {
		return nil, fmt.Errorf("KeyType, KeyBits are not match from CA Cert[%s,%s]", keyType, keyBits)
	}

	return p, nil
}

// newSignerNoSelftest creates a Signer without doing a self-check, useful for tests
func newSignerNoSelftest(ctx context.Context, spec *v1alpha1.NcloudPCAIssuerSpec, client client.Client, namespace string) (*pcaSigner, error) {
	if spec.CaTag == "" {
		return nil, fmt.Errorf("must specify a CaTag")
	}
	p := &pcaSigner{
		caTag:       spec.CaTag,
		ncloudApiGw: spec.NcloudApiGw,
		spec:        spec,
		namespace:   namespace,
		client:      client,
		ctx:         ctx,
	}
	return p, nil
}

func (p *pcaSigner) creatPcaClient() (pcaClient *privateca.APIClient, err error) {
	var apiKey *ncloud.APIKey
	if p.spec.SecretRef.Name != "" {
		if p.spec.SecretRef.Namespace == "" {
			p.spec.SecretRef.Namespace = p.namespace
		}
		secretNamespaceName := types.NamespacedName{
			Namespace: p.spec.SecretRef.Namespace,
			Name:      p.spec.SecretRef.Name,
		}

		secret := new(core.Secret)
		if err := p.client.Get(p.ctx, secretNamespaceName, secret); err != nil {
			return nil, fmt.Errorf("failed to retrieve secret: %v", err)
		}

		key := "NCLOUD_ACCESS_KEY"
		if p.spec.SecretRef.AccessKeyIDSelector.Key != "" {
			key = p.spec.SecretRef.AccessKeyIDSelector.Key
		}
		accessKey, ok := secret.Data[key]
		if !ok {
			return nil, errors.New("no NCLOUD Access Key Found")
		}

		key = "NCLOUD_SECRET_KEY"
		if p.spec.SecretRef.SecretAccessKeySelector.Key != "" {
			key = p.spec.SecretRef.SecretAccessKeySelector.Key
		}
		secretKey, ok := secret.Data[key]
		if !ok {
			return nil, errors.New("no NCLOUD Secret Key Found")
		}
		apiKey = &ncloud.APIKey{
			AccessKey: string(accessKey),
			SecretKey: string(secretKey),
		}

		os.Setenv("NCLOUD_API_GW", p.ncloudApiGw)
		pcaClient = privateca.NewAPIClient(privateca.NewConfiguration(apiKey))

	}
	return pcaClient, nil
}

func extractCertAndCA(data *privateca.SignCsrResponseData) (cert []byte, ca []byte, err error) {
	if data == nil {
		return nil, nil, errors.New("extractCertAndCA: certificate response is nil")
	}
	certBuf := &bytes.Buffer{}

	// Write the leaf to the buffer
	certBuf.WriteString(strings.TrimSpace(*data.Certificate))
	certBuf.WriteRune('\n')

	caBuf := &bytes.Buffer{}
	caBuf.WriteString(strings.TrimSpace(*data.Issuer))
	caBuf.WriteRune('\n')

	// Return the root-most certificate in the CA field.
	return certBuf.Bytes(), caBuf.Bytes(), nil
}
