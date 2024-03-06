//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package policy

import (
	"context"
	"math"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/stretchr/testify/require"
)

func TestWithCaptureResponse(t *testing.T) {
	var httpResp *http.Response
	ctx := WithCaptureResponse(context.Background(), &httpResp)
	require.NotNil(t, ctx)
	raw := ctx.Value(shared.CtxWithCaptureResponse{})
	resp, ok := raw.(**http.Response)
	require.True(t, ok)
	require.Same(t, &httpResp, resp)
}

func TestWithHTTPHeader(t *testing.T) {
	const (
		key = "some"
		val = "thing"
	)
	input := http.Header{}
	input.Set(key, val)
	ctx := WithHTTPHeader(context.Background(), input)
	require.NotNil(t, ctx)
	raw := ctx.Value(shared.CtxWithHTTPHeaderKey{})
	header, ok := raw.(http.Header)
	require.True(t, ok)
	require.EqualValues(t, val, header.Get(key))
}

func TestWithRetryOptions(t *testing.T) {
	ctx := WithRetryOptions(context.Background(), RetryOptions{
		MaxRetries: math.MaxInt32,
	})
	require.NotNil(t, ctx)
	raw := ctx.Value(shared.CtxWithRetryOptionsKey{})
	opts, ok := raw.(RetryOptions)
	require.True(t, ok)
	require.EqualValues(t, math.MaxInt32, opts.MaxRetries)
}
