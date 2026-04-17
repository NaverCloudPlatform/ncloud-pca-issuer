package pca

import "testing"

func TestAPIGatewayURLFor(t *testing.T) {
	cases := []struct {
		name     string
		region   string
		override string
		want     string
		wantErr  bool
	}{
		{name: "empty defaults to public", region: "", want: "https://pca.apigw.ntruss.com"},
		{name: "public", region: "public", want: "https://pca.apigw.ntruss.com"},
		{name: "gov", region: "gov", want: "https://privateca.apigw.gov-ntruss.com"},
		{name: "fin", region: "fin", want: "https://pca.apigw.fin-ntruss.com"},
		{name: "override wins", region: "public", override: "https://example.test", want: "https://example.test"},
		{name: "override wins even on bad region", region: "bogus", override: "https://example.test", want: "https://example.test"},
		{name: "unknown region errors", region: "mars", wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := APIGatewayURLFor(tc.region, tc.override)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (result=%q)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("want %q, got %q", tc.want, got)
			}
		})
	}
}
