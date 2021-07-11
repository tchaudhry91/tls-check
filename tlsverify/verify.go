package tlsverify

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"time"
)

// CertDetails returns the gathered details of the certificate fetched from the remote
type CertDetails struct {
	URL       string    `json:"url,omitempty"`
	ValidFrom time.Time `json:"valid_from,omitempty"`
	ValidTill time.Time `json:"valid_till,omitempty"`
	Domains   []string  `json:"domains,omitempty"`
	Error     error     `json:"error,omitempty"`
}

func JSONPrint(cert CertDetails) {
	data, _ := json.Marshal(cert)
	fmt.Println(string(data))
}

// fetchFromURL converts a url into host/port combination
func fetchFromURL(rawURL string) (host, port string, err error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	return parsed.Hostname(), parsed.Port(), nil
}

// GetTLSCertDetails connects to the remote server, fetches the TLS Cert and returns selected values
func GetTLSCertDetails(rawURL string) (cert CertDetails) {
	cert = CertDetails{}
	cert.URL = rawURL
	host, port, err := fetchFromURL(rawURL)
	if err != nil {
		cert.Error = err
		return
	}
	if port == "" {
		// Default to SSL Check
		port = "443"
	}
	dialer := net.Dialer{Timeout: time.Second * 5}
	conn, err := tls.DialWithDialer(&dialer, "tcp", host+":"+port, nil)
	if err != nil {
		cert.Error = err
		return
	}
	defer conn.Close()
	remoteCert := conn.ConnectionState().PeerCertificates[0]
	cert.Domains = remoteCert.DNSNames
	cert.Error = remoteCert.VerifyHostname(host)
	cert.ValidTill = remoteCert.NotAfter
	cert.ValidFrom = remoteCert.NotBefore

	return
}
