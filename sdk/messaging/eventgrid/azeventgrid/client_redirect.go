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

// maxDefaultRedirects mirrors net/http's default limit and is only applied to
// the same-host redirects that this policy permits when the caller has not
// supplied their own CheckRedirect.
const maxDefaultRedirects = 10

// errCrossHostRedirect is returned to stop an HTTP redirect that would cross to
// a different host.
var errCrossHostRedirect = errors.New("azeventgrid: refusing to follow a redirect to a different host to avoid leaking the publishing credential")

// blockCrossHostRedirect returns an http.Client CheckRedirect function that
// refuses to follow a redirect whenever it would cross to a different host,
// then (for permitted same-host redirects) delegates to prior (the caller's
// existing CheckRedirect, if any).
//
// The Event Grid basic publisher sends its credential in one of the
// Authorization, aeg-sas-key or aeg-sas-token headers. azcore relies on
// net/http to follow redirects; net/http strips only a small, fixed set of
// headers (Authorization, WWW-Authenticate, Cookie, Cookie2) on a cross-host
// redirect and always forwards the custom aeg-sas-* headers. Rather than
// re-sending the request (and its body) to an unconfigured host at all, this
// blocks the cross-host redirect outright, so neither the credential nor the
// event payload is ever sent to a host the caller did not configure. Same-host
// redirects (for example HTTPS enforcement or path normalization performed by a
// gateway or custom domain in front of the topic) are still followed normally.
func blockCrossHostRedirect(prior func(req *http.Request, via []*http.Request) error) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if len(via) > 0 && !strings.EqualFold(req.URL.Hostname(), via[0].URL.Hostname()) {
			return errCrossHostRedirect
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
		transport := http.DefaultTransport
		if dt, ok := http.DefaultTransport.(*http.Transport); ok {
			transport = dt.Clone()
		}
		options.Transport = &http.Client{
			Transport:     transport,
			CheckRedirect: blockCrossHostRedirect(nil),
		}
	case *http.Client:
		clone := *t
		clone.CheckRedirect = blockCrossHostRedirect(t.CheckRedirect)
		options.Transport = &clone
	}
}
