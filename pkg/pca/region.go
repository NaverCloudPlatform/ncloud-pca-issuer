/*
Copyright 2022 Naver Cloud Platform.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package pca

import "fmt"

const (
	// Verified against each realm's API gateway. Host prefix differs per realm:
	// public uses "pca.", gov uses "privateca.". fin follows the public pattern
	// (unverified end-to-end — override with apiGatewayUrl if incorrect).
	apiGatewayPublic = "https://pca.apigw.ntruss.com"
	apiGatewayGov    = "https://privateca.apigw.gov-ntruss.com"
	apiGatewayFin    = "https://pca.apigw.fin-ntruss.com"
)

// APIGatewayURLFor resolves the NCloud API gateway URL for the given region.
// An explicit override takes precedence; an empty region defaults to "public".
func APIGatewayURLFor(region, override string) (string, error) {
	if override != "" {
		return override, nil
	}
	switch region {
	case "", "public":
		return apiGatewayPublic, nil
	case "gov":
		return apiGatewayGov, nil
	case "fin":
		return apiGatewayFin, nil
	default:
		return "", fmt.Errorf("unknown region %q (expected one of: public, gov, fin)", region)
	}
}
