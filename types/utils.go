package types

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

// SortedJSON takes any JSON and returns it sorted by keys. Also, all white-spaces
// are removed.
// This method can be used to canonicalize JSON to be returned by GetSignBytes,
// e.g. for the ledger integration.
// If the passed JSON isn't valid it will return an error.
func SortJSON(toSortJSON []byte) ([]byte, error) {
	var c interface{}
	if err := json.Unmarshal(toSortJSON, &c); err != nil {
		return nil, err
	}
	return json.Marshal(c)
}

// MustSortJSON is like SortJSON but panic if an error occurs, e.g., if
// the passed JSON isn't valid.
func MustSortJSON(toSortJSON []byte) []byte {
	js, err := SortJSON(toSortJSON)
	if err != nil {
		panic(err)
	}
	return js
}

// Uint64ToBigEndian - marshals uint64 to a bigendian byte slice so it can be sorted
func Uint64ToBigEndian(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

// BigEndianToUint64 returns an uint64 from big endian encoded bytes. If encoding
// is empty, zero is returned.
func BigEndianToUint64(bz []byte) uint64 {
	if len(bz) == 0 {
		return 0
	}

	return binary.BigEndian.Uint64(bz)
}

// Slight modification of the RFC3339Nano but it right pads all zeros and drops the time zone info
const SortableTimeFormat = "2006-01-02T15:04:05.000000000"

// Formats a time.Time into a []byte that can be sorted
func FormatTimeBytes(t time.Time) []byte {
	return []byte(t.UTC().Round(0).Format(SortableTimeFormat))
}

// Parses a []byte encoded using FormatTimeKey back into a time.Time
func ParseTimeBytes(bz []byte) (time.Time, error) {
	str := string(bz)
	t, err := time.Parse(SortableTimeFormat, str)
	if err != nil {
		return t, err
	}
	return t.UTC().Round(0), nil
}

// copy bytes
func CopyBytes(bz []byte) (ret []byte) {
	if bz == nil {
		return nil
	}
	ret = make([]byte, len(bz))
	copy(ret, bz)
	return ret
}

func GetTLSCertPool(gateWayURL string) ([]*x509.Certificate, error) {
	if !strings.Contains(strings.ToLower(gateWayURL), "https://") {
		return nil, errors.New("TLS is enabled, but the address is http")
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(gateWayURL)
	defer func() {
		closeErr := resp.Body.Close()
		if err == nil {
			err = closeErr
		}
	}()

	if err != nil {
		return nil, err
	}

	return resp.TLS.PeerCertificates, err
}
