// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"crypto/tls"
	"net/http"
)

type httpClientPolicy struct {
	// Intentionally left empty because all instance share the same defaultHTTPClient object
}

// DefaultHTTPClientPolicy ...
func DefaultHTTPClientPolicy() Policy { return &httpClientPolicy{} }

var defaultHTTPClient *http.Client

func init() {
	defaultTransport := http.DefaultTransport.(*http.Transport)
	transport := &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          defaultTransport.MaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	// TODO: in track 1 we created a cookiejar, do we need one here?  make it an option?  user-specified HTTP client policy?
	defaultHTTPClient = &http.Client{
		Transport: transport,
	}
}

// Do ...
func (p httpClientPolicy) Do(ctx context.Context, req *Request) (*Response, error) {
	response, err := defaultHTTPClient.Do(req.Request.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	return &Response{Response: response}, nil
}
