// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"time"
)

const (
	defaultMaxConnectionsPerEndpoint = 32

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
// azcore default transport but uses Cosmos-specific timeouts and a bounded
// connection pool. Its budget is shared by clients using the default transport.
var defaultCosmosHTTPClient *http.Client

func init() {
	defaultCosmosHTTPClient = newDefaultCosmosHTTPClient()
}

func newDefaultCosmosHTTPClient() *http.Client {
	if transport := newCosmosHTTPTransport(); transport.DialContext == nil {
		// Browser/WASI networking is managed by the platform, not a TCP pool.
		return &http.Client{Transport: transport, Timeout: defaultHTTPRoundTripTimeout}
	}
	client, err := NewHTTPClient(nil)
	if err != nil {
		panic(err)
	}
	return client
}

// HTTPClientOptions configures the HTTP client returned by [NewHTTPClient].
type HTTPClientOptions struct {
	// MaxConnectionsPerEndpoint limits connections to each scheme/host/port.
	// The limit includes dialing, active, idle, and retiring connections, and
	// applies to both HTTP/1.1 and HTTP/2. HTTP/2 streams share connections.
	// Requests wait for capacity, subject to their context and the HTTP client's
	// timeout. The default (zero) is 32. Set -1 to use the unbounded standard
	// transport instead. Values less than -1 are invalid.
	MaxConnectionsPerEndpoint int
}

// NewHTTPClient creates an HTTP client with Cosmos-specific transport settings.
// Pass it as ClientOptions.Transport to customize the connection budget.
// A budget is shared by all Cosmos clients using the returned HTTP client, not
// by independently created HTTP clients. Explicitly supplied transports are
// never modified by Cosmos client constructors.
//
// The client has a 65-second timeout per HTTP attempt, including pool waits and
// response body reads. Close response bodies, and use CloseIdleConnections to
// release unused connections. Bounded pooling does not replay failed requests;
// retries remain the responsibility of the SDK pipeline.
//
// Bounded pooling is unavailable on WebAssembly, where the runtime owns
// networking. On that platform, specify -1 to use the platform transport.
func NewHTTPClient(options *HTTPClientOptions) (*http.Client, error) {
	limit := defaultMaxConnectionsPerEndpoint
	if options != nil && options.MaxConnectionsPerEndpoint != 0 {
		limit = options.MaxConnectionsPerEndpoint
	}
	if limit < -1 {
		return nil, errors.New("MaxConnectionsPerEndpoint must be -1, zero, or positive")
	}
	transport := newCosmosHTTPTransport()
	var roundTripper http.RoundTripper = transport
	if limit > 0 {
		if transport.DialContext == nil {
			return nil, errors.New("bounded Cosmos connection pooling requires a native TCP transport")
		}
		// Avoid parallel address attempts consuming more than one socket per slot.
		transport.DialContext = (&net.Dialer{
			Timeout:       defaultConnectTimeout,
			KeepAlive:     30 * time.Second,
			FallbackDelay: -1,
		}).DialContext
		roundTripper = newCosmosConnectionPool(transport, limit)
	}
	return &http.Client{
		Transport: roundTripper,
		Timeout:   defaultHTTPRoundTripTimeout,
	}, nil
}

func newCosmosHTTPTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: defaultCosmosTransportDialContext(&net.Dialer{
			Timeout:   defaultConnectTimeout,
			KeepAlive: 30 * time.Second,
		}),
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          1000,
		MaxIdleConnsPerHost:   100,
		MaxConnsPerHost:       1000,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion:    tls.VersionTLS12,
			Renegotiation: tls.RenegotiateFreelyAsClient,
		},
		HTTP2: &http.HTTP2Config{
			SendPingTimeout: 2 * time.Second,
			PingTimeout:     1 * time.Second,
		},
	}
}
