package pca

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

// testChain generates a test certificate chain: root -> intermediate(s) -> leaf.
// Returns PEMs in the order leaf, int1, int2, ..., root so they mirror what the
// NCloud API returns (leaf via `certificate`, rest via `caChain`).
type testCert struct {
	cert *x509.Certificate
	key  *rsa.PrivateKey
	pem  string
}

func mustGenCert(t *testing.T, subject string, parent *testCert, isCA bool) *testCert {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(time.Now().UnixNano()),
		Subject:               pkix.Name{CommonName: subject},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  isCA,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	var (
		parentCert *x509.Certificate
		signerKey  *rsa.PrivateKey
	)
	if parent == nil {
		parentCert = tmpl // self-signed
		signerKey = key
	} else {
		parentCert = parent.cert
		signerKey = parent.key
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, parentCert, &key.PublicKey, signerKey)
	if err != nil {
		t.Fatalf("create cert %q: %v", subject, err)
	}
	c, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("parse cert: %v", err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	return &testCert{cert: c, key: key, pem: string(pemBytes)}
}

func TestParsePEMChain(t *testing.T) {
	root := mustGenCert(t, "root", nil, true)
	intA := mustGenCert(t, "intA", root, true)
	intB := mustGenCert(t, "intB", intA, true)

	// Each entry a single PEM.
	chain, err := parsePEMChain([]string{intA.pem, intB.pem, root.pem})
	if err != nil {
		t.Fatalf("parsePEMChain: %v", err)
	}
	if len(chain) != 3 {
		t.Fatalf("want 3 certs, got %d", len(chain))
	}
	if chain[0].Subject.CommonName != "intA" || chain[2].Subject.CommonName != "root" {
		t.Fatalf("unexpected order: %v", []string{chain[0].Subject.CommonName, chain[1].Subject.CommonName, chain[2].Subject.CommonName})
	}

	// A single entry with multiple concatenated PEM blocks.
	concat := intA.pem + intB.pem + root.pem
	chain2, err := parsePEMChain([]string{concat})
	if err != nil {
		t.Fatalf("parsePEMChain concat: %v", err)
	}
	if len(chain2) != 3 {
		t.Fatalf("want 3 certs from concat, got %d", len(chain2))
	}

	// Trailing non-PEM whitespace/junk after a valid block is tolerated.
	withTrailing := root.pem + "\n# trailing comment\n"
	chain3, err := parsePEMChain([]string{withTrailing})
	if err != nil {
		t.Fatalf("parsePEMChain trailing: %v", err)
	}
	if len(chain3) != 1 {
		t.Fatalf("want 1 cert with trailing junk, got %d", len(chain3))
	}

	// Entirely unparsable entry errors out.
	if _, err := parsePEMChain([]string{"not a pem"}); err == nil {
		t.Fatal("expected error for entry without any PEM block")
	}
}

func TestSplitRoot(t *testing.T) {
	root := mustGenCert(t, "root", nil, true)
	intA := mustGenCert(t, "intA", root, true)
	intB := mustGenCert(t, "intB", intA, true)

	ints, got := splitRoot([]*x509.Certificate{intA.cert, intB.cert, root.cert})
	if got == nil || got.Subject.CommonName != "root" {
		t.Fatalf("expected root to be separated")
	}
	if len(ints) != 2 || ints[0].Subject.CommonName != "intA" || ints[1].Subject.CommonName != "intB" {
		t.Fatalf("intermediates preserved in order: got %v", ints)
	}

	// No root present.
	ints2, got2 := splitRoot([]*x509.Certificate{intA.cert, intB.cert})
	if got2 != nil {
		t.Fatalf("expected no root, got %v", got2.Subject)
	}
	if len(ints2) != 2 {
		t.Fatalf("expected 2 intermediates, got %d", len(ints2))
	}
}
