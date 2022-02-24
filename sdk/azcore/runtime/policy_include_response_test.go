//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestIncludeResponse(t *testing.T) {
	ctx, respFn := WithCaptureResponse(context.Background())
	require.NotNil(t, ctx)
	raw := ctx.Value(shared.CtxIncludeResponseKey{})
	_, ok := raw.(**http.Response)
	require.Truef(t, ok, "unexpected type %T", raw)
	require.Nil(t, respFn())
}

func TestIncludeResponsePolicy(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// add a generic HTTP 200 response
	srv.SetResponse()
	// include response policy is automatically added during pipeline construction
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	ctxWithResp, respFn := WithCaptureResponse(context.Background())
	req, err := NewRequest(ctxWithResp, http.MethodGet, srv.URL())
	require.NoError(t, err)
	resp, err := pl.Do(req)
	require.NoError(t, err)
	rawResp := respFn()
	require.NotNil(t, rawResp)
	require.Equal(t, rawResp, resp)
}
