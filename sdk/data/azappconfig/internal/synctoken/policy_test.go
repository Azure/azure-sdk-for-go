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
