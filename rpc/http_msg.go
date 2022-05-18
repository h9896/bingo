package rpc

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	timestampKey = "timestamp"
	signatureKey = "signature"
)

type reqMsg struct {
	method     string
	endpoint   string
	private    bool
	signature  RequestOption
	params     url.Values
	header     http.Header
	bodyString string
	fullURL    string
	signErr    error
}

type RequestOption func(req *reqMsg)

type HttpParameter struct {
	Key        string
	Val        string
	NotReplace bool
}

func SetHeader(headers ...*HttpParameter) RequestOption {
	return func(req *reqMsg) {
		if req.header == nil {
			req.header = http.Header{}
		}

		for _, h := range headers {
			if h.NotReplace {
				req.header.Add(h.Key, h.Val)
			} else {
				req.header.Set(h.Key, h.Val)
			}
		}

	}
}

func SetMethod(method string) RequestOption {
	return func(req *reqMsg) {
		switch strings.ToUpper(method) {
		case "PUT":
			req.method = http.MethodPut
		case "PATCH":
			req.method = http.MethodPatch
		case "DELETE":
			req.method = http.MethodDelete
		case "POST":
			req.method = http.MethodPost
		default:
			req.method = http.MethodGet
		}
	}
}

func SetParams(params ...*HttpParameter) RequestOption {
	return func(req *reqMsg) {
		if req.params == nil {
			req.params = url.Values{}
		}

		for _, p := range params {
			if p.NotReplace {
				req.params.Add(p.Key, p.Val)
			} else {
				req.params.Set(p.Key, p.Val)
			}
		}

	}
}

func SetEndpoint(endpoint string) RequestOption {
	return func(req *reqMsg) {
		req.endpoint = endpoint
	}
}

func SetPrivate() RequestOption {
	return func(req *reqMsg) {
		req.private = true
	}
}

func SetTimestamp() RequestOption {
	return func(req *reqMsg) {
		if req.params == nil {
			req.params = url.Values{}
		}
		req.params.Set(timestampKey, fmt.Sprintf("%d", time.Now().UnixMilli()))
	}
}

func SetSignature(secret string) RequestOption {
	return func(r *reqMsg) {
		r.signature = func(req *reqMsg) {
			if req.params != nil {
				bodyString := req.params.Encode()
				mac := hmac.New(sha256.New, []byte(secret))
				_, err := mac.Write([]byte(bodyString))
				if err != nil {
					req.signErr = err
				} else {
					signature := fmt.Sprintf("%x", mac.Sum(nil))
					req.bodyString = fmt.Sprintf("%s&%s=%s", bodyString, signatureKey, signature)
				}
			}
		}
	}
}
