// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"crypto/tls"
	"net/http"
)

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
	defaultHTTPClient = &http.Client{
		Transport: transport,
	}
}

// DefaultHTTPClientTransport returns the default Transport implementation.
// It uses http.DefaultTransport with a TLS minimum version of 1.2.
func DefaultHTTPClientTransport() Transport {
	return transportFunc(func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return defaultHTTPClient.Do(req.WithContext(ctx))
	})
}
