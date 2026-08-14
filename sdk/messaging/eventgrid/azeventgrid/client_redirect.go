//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package azeventgrid

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

// noFollowRedirect is an http.Client CheckRedirect that refuses to follow any
// HTTP redirect.
//
// The Event Grid basic publisher sends its credential in one of the
// Authorization, aeg-sas-key or aeg-sas-token headers, and the credential (and
// the event payload) are already attached by the time net/http would follow a
// redirect. net/http strips only a small, fixed set of headers on a cross-host
// redirect and always forwards the custom aeg-sas-* headers, so following a
// redirect can leak the credential and the payload to the redirect target.
// Event Grid topic endpoints do not legitimately issue redirects, so the client
// simply does not follow them. This matches the default posture of
// azcore-based clients in other languages (for example .NET's RedirectPolicy
// does not auto-redirect).
//
// Returning http.ErrUseLastResponse leaves the 3xx response in place without
// contacting the redirect target; it surfaces as an *azcore.ResponseError from
// the pipeline.
func noFollowRedirect(_ *http.Request, _ []*http.Request) error {
	return http.ErrUseLastResponse
}

// withRedirectProtection returns a *ClientOptions that disables redirect
// following, without mutating the caller-supplied options. azcore.ClientOptions
// instances may be shared across client constructors, so we apply the redirect
// protection to a shallow copy and leave the caller's options untouched.
func withRedirectProtection(options *ClientOptions) *ClientOptions {
	local := ClientOptions{}
	if options != nil {
		local = *options
	}
	switch t := local.Transport.(type) {
	case nil:
		local.Transport = &http.Client{
			Transport:     newDefaultTransport(),
			CheckRedirect: noFollowRedirect,
		}
	case *http.Client:
		clone := *t
		clone.CheckRedirect = noFollowRedirect
		local.Transport = &clone
	}
	return &local
}

// newDefaultTransport returns an *http.Transport configured identically to
// azcore's default HTTP transport (see
// github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime/transport_default_http_client.go).
//
// When no transport is supplied we must install our own *http.Client to disable
// redirect following. This transport is duplicated (rather than falling back to
// the mutable process-global http.DefaultTransport) so that callers keep
// azcore's TLS floor, connection-pool sizing and HTTP/2 idle health checks. Keep
// this in sync with azcore's default if that ever changes.
func newDefaultTransport() *http.Transport {
	transport := &http.Transport{
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
	if http2Transport, err := http2.ConfigureTransports(transport); err == nil {
		// if the connection has been idle for 10 seconds, send a ping frame for a
		// health check; close the connection if there's no response within 5s.
		http2Transport.ReadIdleTimeout = 10 * time.Second
		http2Transport.PingTimeout = 5 * time.Second
	}
	return transport
}
