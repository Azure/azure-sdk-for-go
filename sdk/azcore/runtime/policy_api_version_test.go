//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestAPIVersionPolicy(t *testing.T) {
	name, version := "api-version", "42"
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()

	for _, header := range []bool{true, false} {
		s := "query param"
		if header {
			s = "header"
		}
		t.Run(s, func(t *testing.T) {
			var location APIVersionLocation = APIVersionLocationQueryParam
			if header {
				location = APIVersionLocationHeader
			}
			p := newAPIVersionPolicy(version, &APIVersionOptions{Location: location, Name: name})
			pl := newTestPipeline(&policy.ClientOptions{Transport: srv, PerCallPolicies: []policy.Policy{p}})

			// when the value isn't set, the policy should set it
			req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
			require.NoError(t, err)
			res, err := pl.Do(req)
			require.NoError(t, err)
			if header {
				require.Equal(t, version, res.Request.Header.Get(name))
			} else {
				require.Equal(t, version, res.Request.URL.Query().Get(name))
			}

			// the policy should override an existing value
			req, err = NewRequest(context.Background(), http.MethodGet, srv.URL())
			require.NoError(t, err)
			if header {
				req.Raw().Header.Set(s, "not-"+version)
			} else {
				q := req.Raw().URL.Query()
				q.Set(s, "not-"+version)
				req.Raw().URL.RawQuery = q.Encode()
			}
			res, err = pl.Do(req)
			require.NoError(t, err)
			if header {
				require.Equal(t, version, res.Request.Header.Get(name))
			} else {
				require.Equal(t, version, res.Request.URL.Query().Get(name))
			}
		})
	}

	for _, test := range []struct {
		err           bool
		location      APIVersionLocation
		name, version string
	}{
		// the policy should modify the request only when given both a version and parameter name
		{},
		{location: APIVersionLocationHeader, version: ""},
		{location: APIVersionLocationQueryParam, version: ""},

		// The policy must know which header/query param to set. This should come from the service client
		// ctor via NewPipeline(). The policy should return an error when the user specifies a version
		// the policy can't set because the service client didn't identify the header/query param.
		{version: version, err: true},
		{location: 2, version: version, err: true},
	} {
		t.Run("no-op", func(t *testing.T) {
			p := newAPIVersionPolicy(test.version, &APIVersionOptions{Location: test.location, Name: test.name})
			pl := newTestPipeline(&policy.ClientOptions{Transport: srv, PerCallPolicies: []policy.Policy{p}})
			req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
			require.NoError(t, err)
			res, err := pl.Do(req)
			if test.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			for _, p := range res.Request.URL.Query() {
				require.NotEqual(t, name, p)
				require.NotContains(t, p, version)
			}
			for _, h := range res.Request.Header {
				require.NotEqual(t, name, h)
				require.NotContains(t, h, version)
			}
		})
	}
}
