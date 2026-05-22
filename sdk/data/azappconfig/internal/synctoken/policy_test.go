// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package synctoken

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestPolicy(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()

	cache := NewCache()
	pl := runtime.NewPipeline("azappconfig", "v0.1.0", runtime.PipelineOptions{
		PerRetry: []policy.Policy{NewPolicy(cache)},
	}, &policy.ClientOptions{
		Transport: &transporter{
			real: srv,
			predicate: func(req *http.Request) {
				require.EqualValues(t, cache.Get(), req.Header.Get(syncTokenHeader))
			},
		},
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)

	srv.AppendResponse()
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	srv.AppendResponse(mock.WithHeader(syncTokenHeader, "id=val"))
	resp, err = pl.Do(req)
	require.Error(t, err) // malformed Sync-Token value
	require.Nil(t, resp)

	srv.AppendResponse(mock.WithHeader(syncTokenHeader, "id=val;sn=1"))
	resp, err = pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	srv.AppendResponse(mock.WithHeader(syncTokenHeader, "id=val;sn=1"))
	resp, err = pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestPolicyMultipleSyncTokensCommaSeparated(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()

	cache := NewCache()
	pl := runtime.NewPipeline("azappconfig", "v0.1.0", runtime.PipelineOptions{
		PerRetry: []policy.Policy{NewPolicy(cache)},
	}, &policy.ClientOptions{
		Transport: srv,
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)

	// server returns two sync tokens comma-separated in a single header value
	srv.AppendResponse(mock.WithHeader(syncTokenHeader, "id1=val1;sn=1,id2=val2;sn=2"))
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// verify both tokens are cached
	cached := cache.Get()
	require.Contains(t, cached, "id1=val1")
	require.Contains(t, cached, "id2=val2")

	// verify the next request carries both tokens in the header
	srv.AppendResponse()
	req2, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)
	pl2 := runtime.NewPipeline("azappconfig", "v0.1.0", runtime.PipelineOptions{
		PerRetry: []policy.Policy{NewPolicy(cache)},
	}, &policy.ClientOptions{
		Transport: &transporter{
			real: srv,
			predicate: func(req *http.Request) {
				hdr := req.Header.Get(syncTokenHeader)
				require.Contains(t, hdr, "id1=val1")
				require.Contains(t, hdr, "id2=val2")
			},
		},
	})
	resp, err = pl2.Do(req2)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestPolicyMultipleSyncTokenHeaders(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()

	cache := NewCache()
	pl := runtime.NewPipeline("azappconfig", "v0.1.0", runtime.PipelineOptions{
		PerRetry: []policy.Policy{NewPolicy(cache)},
	}, &policy.ClientOptions{
		Transport: srv,
	})

	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)

	// server returns two separate Sync-Token headers in the response
	srv.AppendResponse(
		mock.WithHeader(syncTokenHeader, "id1=val1;sn=1"),
		mock.WithHeader(syncTokenHeader, "id2=val2;sn=2"),
	)
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// verify both tokens are cached
	cached := cache.Get()
	require.Contains(t, cached, "id1=val1")
	require.Contains(t, cached, "id2=val2")
}

type transporter struct {
	predicate func(*http.Request)
	real      policy.Transporter
}

func (t *transporter) Do(req *http.Request) (*http.Response, error) {
	if t.predicate != nil {
		t.predicate(req)
	}
	return t.real.Do(req)
}
