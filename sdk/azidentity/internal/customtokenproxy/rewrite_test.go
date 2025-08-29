// Copyright (c) Microsoft Corporation. All rights reserved.

package customtokenproxy

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRewriteProxyRequestURL(t *testing.T) {
	tests := []struct {
		name         string
		proxyURL     *url.URL
		reqURL       *url.URL
		wantScheme   string
		wantHost     string
		wantPath     string
		wantRawPath  string
		wantRawQuery string
	}{
		{
			name: "no RawPath on either; add slash between",
			proxyURL: &url.URL{
				Scheme:  "https",
				Host:    "proxy.example.com",
				Path:    "/base", // no trailing slash
				RawPath: "",      // explicitly empty
			},
			reqURL: &url.URL{
				Scheme:   "https",
				Host:     "orig.example.com",
				Path:     "login", // no leading slash
				RawPath:  "",
				RawQuery: "a=1&b=2",
			},
			wantScheme:   "https",
			wantHost:     "proxy.example.com",
			wantPath:     "/base/login",
			wantRawPath:  "",
			wantRawQuery: "a=1&b=2",
		},
		{
			name: "no RawPath; collapse double slash",
			proxyURL: &url.URL{
				Scheme:  "https",
				Host:    "proxy.example.com",
				Path:    "/v1/", // trailing slash
				RawPath: "",
			},
			reqURL: &url.URL{
				Scheme:   "https",
				Host:     "orig.example.com",
				Path:     "/oauth2/token", // leading slash
				RawPath:  "",
				RawQuery: "x=1",
			},
			wantScheme:   "https",
			wantHost:     "proxy.example.com",
			wantPath:     "/v1/oauth2/token",
			wantRawPath:  "",
			wantRawQuery: "x=1",
		},
		{
			name: "with RawPath; maintain escaped segments and collapse slash",
			proxyURL: &url.URL{
				Scheme:  "https",
				Host:    "proxy.example.com",
				Path:    "/base/",
				RawPath: "/base/",
			},
			reqURL: &url.URL{
				Scheme:   "https",
				Host:     "orig.example.com",
				Path:     "/a b",   // space in segment
				RawPath:  "/a%20b", // encoded form
				RawQuery: "q=1",
			},
			wantScheme:   "https",
			wantHost:     "proxy.example.com",
			wantPath:     "/base/a b",
			wantRawPath:  "/base/a%20b",
			wantRawQuery: "q=1",
		},
		{
			name: "with RawPath both sides no slashes; insert slash",
			proxyURL: &url.URL{
				Scheme:  "https",
				Host:    "proxy.example.com",
				Path:    "/api", // no trailing slash
				RawPath: "/api", // no trailing slash
			},
			reqURL: &url.URL{
				Scheme:  "https",
				Host:    "orig.example.com",
				Path:    "v1", // no leading slash
				RawPath: "v1", // no leading slash
			},
			wantScheme:   "https",
			wantHost:     "proxy.example.com",
			wantPath:     "/api/v1",
			wantRawPath:  "/api/v1",
			wantRawQuery: "",
		},
		{
			name: "with RawPath on proxy only; preserve encoded path",
			proxyURL: &url.URL{
				Scheme:  "https",
				Host:    "proxy.example.com",
				Path:    "/p a",
				RawPath: "/p%20a",
			},
			reqURL: &url.URL{
				Scheme:  "https",
				Host:    "orig.example.com",
				Path:    "/b",
				RawPath: "",
			},
			wantScheme:   "https",
			wantHost:     "proxy.example.com",
			wantPath:     "/p a/b",
			wantRawPath:  "/p%20a/b",
			wantRawQuery: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &http.Request{URL: tc.reqURL}
			rewriteProxyRequestURL(req, tc.proxyURL)

			require.Equal(t, tc.wantScheme, req.URL.Scheme, "scheme mismatch")
			require.Equal(t, tc.wantHost, req.URL.Host, "host mismatch")
			require.Equal(t, tc.wantPath, req.URL.Path, "path mismatch")
			require.Equal(t, tc.wantRawPath, req.URL.RawPath, "raw path mismatch")
			require.Equal(t, tc.wantRawQuery, req.URL.RawQuery, "query mismatch")
		})
	}
}
