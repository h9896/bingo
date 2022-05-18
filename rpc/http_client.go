package rpc

import (
	"context"
	"net/http"
)

type HttpClient interface {
	// Execute a http request
	ExecuteHttpOperation(ctx context.Context, request *reqMsg) (*http.Response, error)

	// Create a http request
	GetHttpRequest(opts ...RequestOption) *reqMsg
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
