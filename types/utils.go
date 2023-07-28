package types

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// GetTLSCertPool get certificates from target server
func GetTLSCertPool(gateWayURL string) ([]*x509.Certificate, error) {
	if !strings.Contains(strings.ToLower(gateWayURL), "https://") {
		return nil, errors.New("this function requires HTTPS protocol")
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(gateWayURL)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil && err == nil {
			err = resp.Body.Close()
		}
	}()

	return resp.TLS.PeerCertificates, err
}
