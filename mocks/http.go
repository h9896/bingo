package mocks

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockClient is the mock client
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

var (
	// GetDoFunc fetches the mock client's `Do` func
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

// Do is the mock client's `Do` func
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

const (
	timestampKey = "timestamp"
	signatureKey = "signature"
	MockDomain   = "mock"
	MockApiKey   = "apikey"
	MockSecret   = "secret"
)

func CheckTimestampAndSignature(t *testing.T, params url.Values) {
	_, timeOk := params[timestampKey]
	_, signOk := params[signatureKey]
	assert.True(t, timeOk, "timestamp not found")
	assert.True(t, signOk, "signature not found")
}
func CheckHeader(t *testing.T, headers http.Header) {
	val, ok := headers["X-Mbx-Apikey"]
	assert.True(t, ok, "apikey header not found")
	assert.Contains(t, val, MockApiKey)
}
