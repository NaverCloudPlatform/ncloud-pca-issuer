package pca

import (
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"

	"github.com/NaverCloudPlatform/ncloud-pca-issuer/pkg/privateca"
)

func strPtr(s string) *string     { return &s }
func slcPtr(s []string) *[]string { return &s }

func parseAllPEM(t *testing.T, in []byte) []*x509.Certificate {
	t.Helper()
	var out []*x509.Certificate
	rest := in
	for len(rest) > 0 {
		block, next := pem.Decode(rest)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			c, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				t.Fatalf("parse cert: %v", err)
			}
			out = append(out, c)
		}
		rest = next
	}
	return out
}

func subjects(certs []*x509.Certificate) []string {
	s := make([]string, 0, len(certs))
	for _, c := range certs {
		s = append(s, c.Subject.CommonName)
	}
	return s
}

func TestExtractCertAndCA_NilData(t *testing.T) {
	if _, _, err := extractCertAndCA(nil); err == nil {
		t.Fatal("expected error on nil data")
	}
	if _, _, err := extractCertAndCA(&privateca.SignCsrResponseData{}); err == nil {
		t.Fatal("expected error on missing certificate")
	}
}

func TestExtractCertAndCA_SingleIssuer(t *testing.T) {
	root := mustGenCert(t, "root", nil, true)
	leaf := mustGenCert(t, "leaf.example", root, false)

	data := &privateca.SignCsrResponseData{
		Certificate: strPtr(leaf.pem),
		Issuer:      strPtr(root.pem),
	}
	cert, ca, err := extractCertAndCA(data)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}

	certs := parseAllPEM(t, cert)
	if len(certs) != 1 || certs[0].Subject.CommonName != "leaf.example" {
		t.Fatalf("want leaf only in tls.crt, got %v", subjects(certs))
	}
	caCerts := parseAllPEM(t, ca)
	if len(caCerts) != 1 || caCerts[0].Subject.CommonName != "root" {
		t.Fatalf("want root in ca.crt, got %v", subjects(caCerts))
	}
}

func TestExtractCertAndCA_WithIntermediate(t *testing.T) {
	root := mustGenCert(t, "root", nil, true)
	inter := mustGenCert(t, "inter", root, true)
	leaf := mustGenCert(t, "leaf.example", inter, false)

	data := &privateca.SignCsrResponseData{
		Certificate: strPtr(leaf.pem),
		Issuer:      strPtr(inter.pem),
		CaChain:     slcPtr([]string{inter.pem, root.pem}),
	}
	cert, ca, err := extractCertAndCA(data)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}

	certs := parseAllPEM(t, cert)
	if got := subjects(certs); len(got) != 2 || got[0] != "leaf.example" || got[1] != "inter" {
		t.Fatalf("tls.crt want [leaf, inter], got %v", got)
	}
	caCerts := parseAllPEM(t, ca)
	if got := subjects(caCerts); len(got) != 1 || got[0] != "root" {
		t.Fatalf("ca.crt want [root], got %v", got)
	}
	// tls.crt must not contain the root.
	if strings.Contains(string(cert), strings.TrimSpace(root.pem)) {
		t.Fatal("tls.crt must not include root certificate")
	}
}

func TestExtractCertAndCA_TwoIntermediates(t *testing.T) {
	root := mustGenCert(t, "root", nil, true)
	int1 := mustGenCert(t, "int1", root, true)
	int2 := mustGenCert(t, "int2", int1, true)
	leaf := mustGenCert(t, "leaf.example", int2, false)

	data := &privateca.SignCsrResponseData{
		Certificate: strPtr(leaf.pem),
		Issuer:      strPtr(int2.pem),
		CaChain:     slcPtr([]string{int2.pem, int1.pem, root.pem}),
	}
	cert, ca, err := extractCertAndCA(data)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}

	certs := parseAllPEM(t, cert)
	if got := subjects(certs); len(got) != 3 || got[0] != "leaf.example" || got[1] != "int2" || got[2] != "int1" {
		t.Fatalf("tls.crt want [leaf, int2, int1], got %v", got)
	}
	caCerts := parseAllPEM(t, ca)
	if got := subjects(caCerts); len(got) != 1 || got[0] != "root" {
		t.Fatalf("ca.crt want [root], got %v", got)
	}
}

func TestExtractCertAndCA_SelfSignedInCaChain(t *testing.T) {
	// Shape observed from NCloud production for a single-tier CA:
	// Certificate = leaf, Issuer = root, CaChain = [root], root is self-signed.
	root := mustGenCert(t, "root", nil, true)
	leaf := mustGenCert(t, "leaf.example", root, false)

	data := &privateca.SignCsrResponseData{
		Certificate: strPtr(leaf.pem),
		Issuer:      strPtr(root.pem),
		CaChain:     slcPtr([]string{root.pem}),
	}
	cert, ca, err := extractCertAndCA(data)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	certs := parseAllPEM(t, cert)
	if got := subjects(certs); len(got) != 1 || got[0] != "leaf.example" {
		t.Fatalf("tls.crt want [leaf], got %v", got)
	}
	caCerts := parseAllPEM(t, ca)
	if got := subjects(caCerts); len(got) != 1 || got[0] != "root" {
		t.Fatalf("ca.crt want [root], got %v", got)
	}
	if strings.Contains(string(cert), strings.TrimSpace(root.pem)) {
		t.Fatal("tls.crt must not include self-signed root")
	}
}

func TestExtractCertAndCA_CaChainEmptyStringsFallBackToIssuer(t *testing.T) {
	root := mustGenCert(t, "root", nil, true)
	leaf := mustGenCert(t, "leaf.example", root, false)

	data := &privateca.SignCsrResponseData{
		Certificate: strPtr(leaf.pem),
		Issuer:      strPtr(root.pem),
		CaChain:     slcPtr([]string{"", "   "}),
	}
	cert, ca, err := extractCertAndCA(data)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	certs := parseAllPEM(t, cert)
	if got := subjects(certs); len(got) != 1 || got[0] != "leaf.example" {
		t.Fatalf("tls.crt want [leaf], got %v", got)
	}
	caCerts := parseAllPEM(t, ca)
	if got := subjects(caCerts); len(got) != 1 || got[0] != "root" {
		t.Fatalf("ca.crt want [root], got %v", got)
	}
}

func TestExtractCertAndCA_ChainMissingRoot(t *testing.T) {
	root := mustGenCert(t, "root", nil, true)
	inter := mustGenCert(t, "inter", root, true)
	leaf := mustGenCert(t, "leaf.example", inter, false)

	// Only intermediate in chain (root omitted — best-effort fallback).
	data := &privateca.SignCsrResponseData{
		Certificate: strPtr(leaf.pem),
		CaChain:     slcPtr([]string{inter.pem}),
	}
	cert, ca, err := extractCertAndCA(data)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}

	certs := parseAllPEM(t, cert)
	if got := subjects(certs); len(got) != 2 || got[0] != "leaf.example" || got[1] != "inter" {
		t.Fatalf("tls.crt want [leaf, inter], got %v", got)
	}
	caCerts := parseAllPEM(t, ca)
	if got := subjects(caCerts); len(got) != 1 || got[0] != "inter" {
		t.Fatalf("ca.crt fallback want [inter], got %v", got)
	}
}
