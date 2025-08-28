// Copyright 2011 The Go Authors. All rights reserved.
// Licensed under the BSD-style license that be found at https://golang.org/LICENSE .
//
// Copyright (c) Microsoft Corporation. All rights reserved.
// Portions of this file are modifications licensed under the MIT License.

package customtokenproxy

import (
	"net/http"
	"net/url"
	"strings"
)

// rewriteRequestURL rewrites the request URL to target the specified URL.
// Target is the token proxy URL in custom token endpoint mode.
//
// using a similar approach as httputil.ReverseProxy to ensure all edge cases are handled correctly
// ref: https://cs.opensource.google/go/go/+/refs/tags/go1.25.0:src/net/http/httputil/reverseproxy.go;l=282;drc=84e0061460d7c9a624a74e13f0212f443b079531
func rewriteRequestURL(req *http.Request, target *url.URL) {
	targetQuery := target.RawQuery
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
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
