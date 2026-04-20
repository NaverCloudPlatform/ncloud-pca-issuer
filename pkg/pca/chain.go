/*
Copyright 2022 Naver Cloud Platform.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package pca

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
)

// parsePEMChain decodes each PEM-encoded string in order and returns the parsed
// x509 certificates. A single input string may contain multiple concatenated
// PEM blocks; all of them are returned. Trailing whitespace or non-PEM content
// after the last valid block is tolerated. An entry yielding zero blocks is an
// error.
func parsePEMChain(pems []string) ([]*x509.Certificate, error) {
	var out []*x509.Certificate
	for i, p := range pems {
		rest := []byte(strings.TrimSpace(p))
		if len(rest) == 0 {
			continue
		}
		blocksInEntry := 0
		for len(rest) > 0 {
			var block *pem.Block
			block, rest = pem.Decode(rest)
			if block == nil {
				if blocksInEntry == 0 {
					return nil, fmt.Errorf("chain entry %d: no PEM block found", i)
				}
				break
			}
			blocksInEntry++
			if block.Type != "CERTIFICATE" {
				continue
			}
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("chain entry %d: %w", i, err)
			}
			out = append(out, cert)
		}
	}
	return out, nil
}

// splitRoot separates a certificate chain into intermediates and the root.
// A certificate is considered the root when it is self-signed. Intermediates
// retain their original order. If no self-signed certificate is present,
// root is nil.
func splitRoot(chain []*x509.Certificate) (intermediates []*x509.Certificate, root *x509.Certificate) {
	for _, c := range chain {
		if isSelfSigned(c) {
			if root == nil {
				root = c
			}
			continue
		}
		intermediates = append(intermediates, c)
	}
	return intermediates, root
}

func isSelfSigned(c *x509.Certificate) bool {
	if c == nil {
		return false
	}
	return c.CheckSignatureFrom(c) == nil
}
