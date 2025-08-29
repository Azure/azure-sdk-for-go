// Copyright (c) Microsoft Corporation. All rights reserved.

package customtokenproxy

import (
	"net/http"
	"net/url"
	"strings"
)

// rewrites the request URL to target the specified URL.
// Target is the token proxy URL in custom token endpoint mode.
func rewriteProxyRequestURL(req *http.Request, proxyURL *url.URL) {
	req.URL.Scheme = proxyURL.Scheme // always https
	req.URL.Host = proxyURL.Host     // host:port of the proxy

	// NOTE: proxyURL doesn't include query, req might include query
	// we just retain the RawQuery from req.URL

	req.URL.Path, req.URL.RawPath = joinURLPath(proxyURL, req.URL)
}

// using a similar approach as httputil.ReverseProxy to ensure all edge cases are handled correctly
// ref: https://cs.opensource.google/go/go/+/refs/tags/go1.25.0:src/net/http/httputil/reverseproxy.go;l=282;drc=84e0061460d7c9a624a74e13f0212f443b079531
func joinURLPath(a, b *url.URL) (string, string) {
	if a.RawPath != "" || b.RawPath != "" {
		return joinURLEscapedPath(a, b)
	}

	aslash := strings.HasSuffix(a.Path, "/")
	bslash := strings.HasPrefix(b.Path, "/")
	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], ""
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, ""
	}
	return a.Path + b.Path, ""
}

func joinURLEscapedPath(a, b *url.URL) (path, rawpath string) {
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}
