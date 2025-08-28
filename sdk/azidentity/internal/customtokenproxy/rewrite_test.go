// Copyright 2011 The Go Authors. All rights reserved.
// Licensed under the BSD-style license.

package customtokenproxy

import (
	"net/url"
	"testing"
)

func TestSingleJoinSlash(t *testing.T) {
	tests := []struct {
		slasha   string
		slashb   string
		expected string
	}{
		{"https://www.google.com/", "/favicon.ico", "https://www.google.com/favicon.ico"},
		{"https://www.google.com", "/favicon.ico", "https://www.google.com/favicon.ico"},
		{"https://www.google.com", "favicon.ico", "https://www.google.com/favicon.ico"},
		{"https://www.google.com", "", "https://www.google.com/"},
		{"", "favicon.ico", "/favicon.ico"},
	}
	for _, tt := range tests {
		if got := singleJoiningSlash(tt.slasha, tt.slashb); got != tt.expected {
			t.Errorf("singleJoiningSlash(%q,%q) want %q got %q",
				tt.slasha,
				tt.slashb,
				tt.expected,
				got)
		}
	}
}

func TestJoinURLPath(t *testing.T) {
	tests := []struct {
		a        *url.URL
		b        *url.URL
		wantPath string
		wantRaw  string
	}{
		{&url.URL{Path: "/a/b"}, &url.URL{Path: "/c"}, "/a/b/c", ""},
		{&url.URL{Path: "/a/b", RawPath: "badpath"}, &url.URL{Path: "c"}, "/a/b/c", "/a/b/c"},
		{&url.URL{Path: "/a/b", RawPath: "/a%2Fb"}, &url.URL{Path: "/c"}, "/a/b/c", "/a%2Fb/c"},
		{&url.URL{Path: "/a/b", RawPath: "/a%2Fb"}, &url.URL{Path: "/c"}, "/a/b/c", "/a%2Fb/c"},
		{&url.URL{Path: "/a/b/", RawPath: "/a%2Fb%2F"}, &url.URL{Path: "c"}, "/a/b//c", "/a%2Fb%2F/c"},
		{&url.URL{Path: "/a/b/", RawPath: "/a%2Fb/"}, &url.URL{Path: "/c/d", RawPath: "/c%2Fd"}, "/a/b/c/d", "/a%2Fb/c%2Fd"},
	}

	for _, tt := range tests {
		p, rp := joinURLPath(tt.a, tt.b)
		if p != tt.wantPath || rp != tt.wantRaw {
			t.Errorf("joinURLPath(URL(%q,%q),URL(%q,%q)) want (%q,%q) got (%q,%q)",
				tt.a.Path, tt.a.RawPath,
				tt.b.Path, tt.b.RawPath,
				tt.wantPath, tt.wantRaw,
				p, rp)
		}
	}
}
