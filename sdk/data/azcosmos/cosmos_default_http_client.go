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
	// connection to the Cosmos service. Cosmos accounts are expected to be
	// reachable from the configured/preferred region within a few hundred
	// milliseconds; failing fast at the dial layer lets the global endpoint
	// manager retry against the next preferred region instead of blocking on a
	// dead endpoint for the azcore default of 30 seconds.
	defaultConnectTimeout = 5 * time.Second

	// defaultHTTPRoundTripTimeout is the default wall-clock cap that the
	// underlying *http.Client applies to a single HTTP attempt, covering
	// connection setup, request write, response header read, and response body
	// read. It is intentionally set as http.Client.Timeout (not as a per-try
	// retry timeout) so it survives custom per-call policies and acts as a
	// hard backstop against runaway requests.
	//
	// Trade-offs to understand before changing this value:
	//
	//   * http.Client.Timeout is a wall-clock cap on the entire round-trip,
	//     including streaming the response body. A caller-supplied
	//     context.WithTimeout that is *longer* than this value will be
	//     truncated by the HTTP client; a *shorter* caller context still wins.
	//     This is the intended safety property: no Cosmos request should hang
	//     for an unbounded amount of time even if the caller forgot to set a
	//     deadline.
	//
	//   * The azcore retry policy is layered above the transport, so the
	//     65-second cap applies per HTTP attempt; the policy can still issue
	//     additional retries when one attempt exceeds the cap.
	//
	//   * 65 seconds was chosen to exceed the Cosmos gateway's own server-side
	//     request budget (~60s) by a small margin so the server gets a chance
	//     to return a structured error (which the retry policy can interpret)
	//     before the client gives up locally.
	//
	// Callers that legitimately need to drain very large query/change-feed
	// pages that take longer than this can override by supplying their own
	// Transport via ClientOptions.
	defaultHTTPRoundTripTimeout = 65 * time.Second
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
	if http2Transport, err := http2.ConfigureTransports(transport); err == nil {
		http2Transport.ReadIdleTimeout = 2 * time.Second
		http2Transport.PingTimeout = 1 * time.Second
	}
	return &http.Client{
		Transport: transport,
		Timeout:   defaultHTTPRoundTripTimeout,
	}
}
