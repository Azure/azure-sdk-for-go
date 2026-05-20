// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

const (
	// defaultConnectTimeout is the default timeout for establishing a new TCP
	// connection to the Cosmos service.
	defaultConnectTimeout = 5 * time.Second

	// defaultRequestTimeout is the default end-to-end timeout applied by the
	// HTTP client to a single request, including connection setup, sending the
	// request, and reading the entire response body. Callers that need a
	// different bound can supply their own Transport via ClientOptions.
	defaultRequestTimeout = 65 * time.Second
)

// defaultCosmosHTTPClient is the http.Client used by Cosmos clients when the
// caller does not provide a custom Transport via ClientOptions. It mirrors the
// azcore default transport but uses Cosmos-specific connect and request
// timeouts.
var defaultCosmosHTTPClient *http.Client

func init() {
	defaultCosmosHTTPClient = newDefaultCosmosHTTPClient()
}

func newDefaultCosmosHTTPClient() *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   defaultConnectTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          1000,
		MaxIdleConnsPerHost:   1000,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion:    tls.VersionTLS12,
			Renegotiation: tls.RenegotiateFreelyAsClient,
		},
	}
	if http2Transport, err := http2.ConfigureTransports(transport); err == nil {
		http2Transport.ReadIdleTimeout = 2 * time.Second
		http2Transport.PingTimeout = 1 * time.Second
	}
	return &http.Client{
		Transport: transport,
		Timeout:   defaultRequestTimeout,
	}
}
