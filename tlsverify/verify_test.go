package tlsverify_test

import (
	"testing"

	"github.com/tchaudhry91/tls-check/tlsverify"
)

// TestGetCertDetails is a simple test to check the library. This is NOT guaranteed to always work and as such should not be relied on for correctness.
func TestGetCertDetails(t *testing.T) {
	type testCase struct {
		RawURL    string
		ValidConn bool
	}

	cases := []testCase{
		{"https://google.com", true},
		{"psql://hh-pgsql-public.ebi.ac.uk:5432", false},
		{"https://adobe.com", true},
		{"https://sdfasdfasdf.com", false},
	}

	for _, c := range cases {
		cert := tlsverify.GetTLSCertDetails(c.RawURL)
		if cert.Error != nil && c.ValidConn {
			t.Errorf("Failed on URL:%s, Error: %s", c.RawURL, cert.Error.Error())
		}
		t.Log(cert)
	}
}
