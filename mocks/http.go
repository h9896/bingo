package mocks

import "net/http"

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
