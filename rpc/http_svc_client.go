package rpc

import (
	"context"
	"fmt"
	"net/http"
)

type httpSvcClient struct {
	client   httpClient
	protocol string
	apiKey   string
}

func NewHttpClient(apiKey string, useSSL bool, c httpClient) HttpClient {
	client := &httpSvcClient{
		apiKey: apiKey,
	}

	if c == nil {
		client.client = http.DefaultClient
	} else {
		client.client = c
	}

	if useSSL {
		client.protocol = "https"
	} else {
		client.protocol = "http"
	}

	return client
}

// Execute a http request
func (c *httpSvcClient) ExecuteHttpOperation(ctx context.Context, request *reqMsg) (*http.Response, error) {

	if request.private {
		SetHeader(&HttpParameter{Key: "X-MBX-APIKEY", Val: c.apiKey})(request)
	}

	if request.params != nil {
		SetHeader(&HttpParameter{Key: "Content-Type", Val: "application/x-www-form-urlencoded"})(request)
	}

	if request.signature != nil {
		request.signature(request)
	} else {
		request.bodyString = request.params.Encode()
	}
	if request.signErr != nil {
		return nil, request.signErr
	}

	if request.bodyString != "" {
		request.fullURL = fmt.Sprintf("%s://%s?%s", c.protocol, request.endpoint, request.bodyString)
	} else {
		request.fullURL = fmt.Sprintf("%s://%s", c.protocol, request.endpoint)
	}

	req, err := http.NewRequestWithContext(ctx, request.method, request.fullURL, nil)

	if err != nil {
		return nil, err
	}

	if request.header != nil {
		req.Header = request.header
	}

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Create a http request
func (c *httpSvcClient) GetHttpRequest(opts ...RequestOption) *reqMsg {
	req := &reqMsg{
		private: false,
	}

	for _, opt := range opts {
		opt(req)
	}

	return req
}
