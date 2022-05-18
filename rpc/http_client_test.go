package rpc

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestOption(t *testing.T) {

	req := &reqMsg{}
	SetEndpoint("api/v1")(req)
	assert.EqualValues(t, "api/v1", req.endpoint)

	SetHeader(&HttpParameter{Key: "Content-Type", Val: "application/json"})(req)
	assert.EqualValues(t, http.Header{"Content-Type": []string{"application/json"}}, req.header)

	SetMethod("")(req)
	assert.EqualValues(t, http.MethodGet, req.method)
	SetMethod("post")(req)
	assert.EqualValues(t, http.MethodPost, req.method)
	SetMethod("patch")(req)
	assert.EqualValues(t, http.MethodPatch, req.method)
	SetMethod("delete")(req)
	assert.EqualValues(t, http.MethodDelete, req.method)
	SetMethod("put")(req)
	assert.EqualValues(t, http.MethodPut, req.method)

	params := []*HttpParameter{
		{Key: "abc", Val: "ww", NotReplace: true},
		{Key: "abc", Val: "q", NotReplace: true},
		{Key: "b", Val: "wwq"},
	}

	expectP := url.Values{}
	expectP.Add("abc", "ww")
	expectP.Add("abc", "q")
	expectP.Set("b", "wwq")
	SetParams(params...)(req)
	assert.EqualValues(t, expectP, req.params)

	SetPrivate()(req)
	assert.True(t, true, req.private)

	SetTimestamp()(req)
	assert.NotNil(t, req.params.Get(timestampKey))

}

func TestGetHttpRequest(t *testing.T) {
	client := NewHttpClient("apikey", true, nil)
	query := []*HttpParameter{
		{Key: "symbol", Val: "BTCUSD_PERP"},
		{Key: "limit", Val: "5"}}
	req := client.GetHttpRequest(
		SetEndpoint("dapi.binance.com/dapi/v1/depth"),
		SetMethod("get"),
		SetParams(query...))

	expectP := url.Values{}
	expectP.Set("symbol", "BTCUSD_PERP")
	expectP.Set("limit", "5")

	assert.EqualValues(t, "dapi.binance.com/dapi/v1/depth", req.endpoint)
	assert.EqualValues(t, http.MethodGet, req.method)
	assert.EqualValues(t, expectP, req.params)

}

type mockHttpClient struct{}

func mockClient() httpClient {
	return &mockHttpClient{}
}
func (c *mockHttpClient) Do(req *http.Request) (resp *http.Response, error error) {
	resp = &http.Response{Request: req}
	return
}

func TestExecuteHttpOperationWithSignature(t *testing.T) {
	client := NewHttpClient("apikey", false, mockClient())

	body := []*HttpParameter{
		{Key: "symbol", Val: "BTCUSD_PERP"},
		{Key: "side", Val: "BUY"},
		{Key: "recvWindow", Val: "600000"},
		{Key: "type", Val: "LIMIT"},
		{Key: "timeInForce", Val: "GTC"},
		{Key: "quantity", Val: "0.0035"},
		{Key: "price", Val: "28000.1"}}
	req := client.GetHttpRequest(SetEndpoint("dapi.binance.com/dapi/v1/order"),
		SetPrivate(),
		SetMethod("post"),
		SetParams(body...),
		SetTimestamp(),
		SetSignature("secret"))
	resp, _ := client.ExecuteHttpOperation(context.Background(), req)

	assert.EqualValues(t, req.header, resp.Request.Header)
	assert.EqualValues(t, req.method, resp.Request.Method)
	assert.EqualValues(t, req.fullURL, resp.Request.URL.String())
}

func TestExecuteHttpOperation(t *testing.T) {
	client := NewHttpClient("apikey", false, mockClient())
	req := client.GetHttpRequest(SetEndpoint("dapi.binance.com/dapi/v1"),
		SetPrivate(),
		SetMethod("get"))
	resp, _ := client.ExecuteHttpOperation(context.Background(), req)

	assert.EqualValues(t, req.header, resp.Request.Header)
	assert.EqualValues(t, req.method, resp.Request.Method)
	assert.EqualValues(t, req.fullURL, resp.Request.URL.String())
}
