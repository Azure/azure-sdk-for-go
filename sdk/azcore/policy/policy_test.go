//go:build go1.16
// +build go1.16

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

func TestWithHTTPHeader(t *testing.T) {
	const (
		key = "some"
		val = "thing"
	)
	input := http.Header{}
	input.Set(key, val)
	ctx := WithHTTPHeader(context.Background(), input)
	if ctx == nil {
		t.Fatal("nil context")
	}
	raw := ctx.Value(shared.CtxWithHTTPHeaderKey{})
	header, ok := raw.(http.Header)
	if !ok {
		t.Fatalf("unexpected type %T", raw)
	}
	if v := header.Get(key); v != val {
		t.Fatalf("unexpected value %s", v)
	}
}

func TestWithRetryOptions(t *testing.T) {
	ctx := WithRetryOptions(context.Background(), RetryOptions{
		MaxRetries: math.MaxInt32,
	})
	if ctx == nil {
		t.Fatal("nil context")
	}
	raw := ctx.Value(shared.CtxWithRetryOptionsKey{})
	opts, ok := raw.(RetryOptions)
	if !ok {
		t.Fatalf("unexpected type %T", raw)
	}
	if opts.MaxRetries != math.MaxInt32 {
		t.Fatalf("unexpected value %d", opts.MaxRetries)
	}
}

func TestIncludeResponse(t *testing.T) {
	ctx := IncludeResponse(context.Background())
	require.NotNil(t, ctx)
	raw := ctx.Value(shared.CtxIncludeResponseKey{})
	_, ok := raw.(**http.Response)
	require.Truef(t, ok, "unexpected type %T", raw)
}
