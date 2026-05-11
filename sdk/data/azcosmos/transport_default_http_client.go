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

// defaultCosmosHTTPClient is the default HTTP client used by the Cosmos
// SDK when the caller does not supply one via ClientOptions.Transport.
//
// It mirrors the default transport configured by sdk/azcore/runtime, but
// uses more aggressive HTTP/2 PING-based health checks tuned for Cosmos
// workloads. See https://github.com/golang/go/issues/59690 for the
// underlying Go issue this works around.
var defaultCosmosHTTPClient *http.Client

func init() {
	defaultCosmosHTTPClient, _ = newDefaultCosmosHTTPClient()
}

// newDefaultCosmosHTTPClient builds the Cosmos default *http.Client and
// returns the configured *http2.Transport so callers (in particular tests)
// can verify the HTTP/2 PING parameters.
//
// The returned *http2.Transport is nil when http2.ConfigureTransports
// failed; the *http.Client is always usable in that case (HTTP/1.1).
func newDefaultCosmosHTTPClient() (*http.Client, *http2.Transport) {
	defaultTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: defaultTransportDialContext(&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}),
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion:    tls.VersionTLS12,
			Renegotiation: tls.RenegotiateFreelyAsClient,
		},
	}
	// TODO: evaluate removing this once https://github.com/golang/go/issues/59690 has been fixed.
	http2Transport, err := http2.ConfigureTransports(defaultTransport)
	if err == nil {
		// ReadIdleTimeout triggers an HTTP/2 PING frame if no data is received
		// from the server within this duration. A value of 1s aggressively
		// detects broken connections before they are reused.
		http2Transport.ReadIdleTimeout = 1 * time.Second
		// PingTimeout is how long to wait for a PING response before the
		// connection is considered dead and closed.
		http2Transport.PingTimeout = 2 * time.Second
	} else {
		http2Transport = nil
	}
	return &http.Client{Transport: defaultTransport}, http2Transport
}
