//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package azeventgrid

import (
	"errors"
	"net/http"
	"strings"
)

// credentialRedirectHeaders are the authentication headers that must never be
// forwarded to a different host when an HTTP redirect is followed.
//
// The Event Grid basic publisher sends its credential in one of these headers.
// azcore relies on net/http to follow redirects, and net/http only strips a
// small, fixed set of headers (Authorization, WWW-Authenticate, Cookie, Cookie2)
// when a redirect crosses to a different host. The custom aeg-sas-key /
// aeg-sas-token headers are NOT in that set, so without this cleanup they would
// be leaked in full to a redirect target on a different host.
var credentialRedirectHeaders = []string{
	"Authorization",
	"aeg-sas-key",
	"aeg-sas-token",
}

// maxDefaultRedirects mirrors net/http's default limit and is only used when the
// caller has not supplied their own CheckRedirect.
const maxDefaultRedirects = 10

// stripCredentialsOnCrossHostRedirect returns an http.Client CheckRedirect
// function that removes credential headers whenever a redirect crosses to a
// different host, then delegates to prior (the caller's existing CheckRedirect,
// if any).
func stripCredentialsOnCrossHostRedirect(prior func(req *http.Request, via []*http.Request) error) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if len(via) > 0 {
			originalHost := via[0].URL.Hostname()
			if !strings.EqualFold(req.URL.Hostname(), originalHost) {
				for _, h := range credentialRedirectHeaders {
					req.Header.Del(h)
				}
			}
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

// applyRedirectCredentialProtection ensures the HTTP client used by the pipeline
// strips Event Grid credential headers on cross-host redirects, closing a
// credential-leak vector. It preserves a caller-supplied transport where
// possible:
//
//   - no transport supplied: a default client is installed with the protection.
//   - an *http.Client supplied: it is cloned and the protection is chained onto
//     any CheckRedirect the caller already set.
//   - any other custom transport: left untouched (such a transport is
//     responsible for its own redirect handling).
func applyRedirectCredentialProtection(options *ClientOptions) {
	switch t := options.Transport.(type) {
	case nil:
		transport := http.DefaultTransport
		if dt, ok := http.DefaultTransport.(*http.Transport); ok {
			transport = dt.Clone()
		}
		options.Transport = &http.Client{
			Transport:     transport,
			CheckRedirect: stripCredentialsOnCrossHostRedirect(nil),
		}
	case *http.Client:
		clone := *t
		clone.CheckRedirect = stripCredentialsOnCrossHostRedirect(t.CheckRedirect)
		options.Transport = &clone
	}
}
