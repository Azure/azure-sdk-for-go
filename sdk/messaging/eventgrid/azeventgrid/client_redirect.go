//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package azeventgrid

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/http2"
)

// maxDefaultRedirects mirrors net/http's default limit and is only applied to
// the same-origin redirects that this policy permits when the caller has not
// supplied their own CheckRedirect.
const maxDefaultRedirects = 10

// errUnsafeRedirect is returned to stop an HTTP redirect that would leave the
// configured origin (a different host or port, or an https->http downgrade).
var errUnsafeRedirect = errors.New("azeventgrid: refusing to follow a redirect to a different origin (host, port, or an https-to-http downgrade) to avoid leaking the publishing credential")

// blockUnsafeRedirect returns an http.Client CheckRedirect function that refuses
// to follow a redirect whenever it would leave the configured origin, then (for
// permitted same-origin redirects) delegates to prior (the caller's existing
// CheckRedirect, if any).
//
// The Event Grid basic publisher sends its credential in one of the
// Authorization, aeg-sas-key or aeg-sas-token headers. azcore relies on
// net/http to follow redirects; net/http strips only a small, fixed set of
// headers (Authorization, WWW-Authenticate, Cookie, Cookie2) on a cross-host
// redirect and always forwards the custom aeg-sas-* headers. Because the
// credential (and the request body) have already been attached by the time
// net/http follows a redirect, this refuses to follow any redirect that changes
// the origin, so neither the credential nor the event payload is ever sent to a
// host/port the caller did not configure, and an https->http downgrade can never
// leak the credential in cleartext. Same-origin redirects (for example path
// normalization, or an http->https upgrade performed by a gateway or custom
// domain in front of the topic) are still followed normally.
func blockUnsafeRedirect(prior func(req *http.Request, via []*http.Request) error) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if len(via) > 0 && !isSafeRedirect(via[0].URL, req.URL) {
			return errUnsafeRedirect
		}

		if prior != nil {
			return prior(req, via)
		}

		if len(via) >= maxDefaultRedirects {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}
}

// isSafeRedirect reports whether a redirect from original to target stays within
// the caller-configured origin and therefore may be followed with the
// credential still attached. A redirect is considered safe only when it targets
// the same host and port, or performs a standard http->https upgrade of the same
// host. Any host change, port change, or https->http downgrade is unsafe.
func isSafeRedirect(original, target *url.URL) bool {
	if !strings.EqualFold(target.Hostname(), original.Hostname()) {
		return false // different host
	}

	origScheme := strings.ToLower(original.Scheme)
	targetScheme := strings.ToLower(target.Scheme)

	if origScheme == "https" && targetScheme == "http" {
		return false // scheme downgrade would leak the credential in cleartext
	}
	if origScheme == "http" && targetScheme == "https" {
		return true // standard upgrade to TLS on the same host
	}

	// Same scheme: require an identical effective port (full-authority match).
	return effectivePort(original) == effectivePort(target)
}

// effectivePort returns the port for u, defaulting to the well-known port for
// the URL's scheme when no explicit port is present.
func effectivePort(u *url.URL) string {
	if p := u.Port(); p != "" {
		return p
	}
	switch strings.ToLower(u.Scheme) {
	case "https":
		return "443"
	case "http":
		return "80"
	}
	return ""
}

// applyRedirectCredentialProtection ensures the HTTP client used by the pipeline
// refuses cross-host redirects, closing a credential-leak vector. It preserves a
// caller-supplied transport where possible:
//
//   - no transport supplied: a default client is installed with the protection.
//   - an *http.Client supplied: it is cloned and the protection is chained onto
//     any CheckRedirect the caller already set.
//   - any other custom transport: left untouched (such a transport is
//     responsible for its own redirect handling).
func applyRedirectCredentialProtection(options *ClientOptions) {
	switch t := options.Transport.(type) {
	case nil:
		options.Transport = &http.Client{
			Transport:     newDefaultTransport(),
			CheckRedirect: blockUnsafeRedirect(nil),
		}
	case *http.Client:
		clone := *t
		clone.CheckRedirect = blockUnsafeRedirect(t.CheckRedirect)
		options.Transport = &clone
	}
}

// newDefaultTransport returns an *http.Transport configured identically to
// azcore's default HTTP transport (see
// github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime/transport_default_http_client.go).
//
// When no transport is supplied we must install our own *http.Client to attach
// the redirect protection. This transport is duplicated (rather than falling
// back to the mutable process-global http.DefaultTransport) so that callers keep
// azcore's TLS floor, connection-pool sizing and HTTP/2 idle health checks. Keep
// this in sync with azcore's default if that ever changes.
func newDefaultTransport() *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
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
		// if the connection has been idle for 10 seconds, send a ping frame for a
		// health check; close the connection if there's no response within 5s.
		http2Transport.ReadIdleTimeout = 10 * time.Second
		http2Transport.PingTimeout = 5 * time.Second
	}
	return transport
}
