//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDelay(t *testing.T) {
	if err := Delay(context.Background(), 5*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := Delay(ctx, 5*time.Minute); err == nil {
		t.Fatal("unexpected nil error")
	}
}

func TestRetryAfter(t *testing.T) {
	if RetryAfter(nil) != 0 {
		t.Fatal("expected zero duration")
	}
	resp := &http.Response{
		Header: http.Header{},
	}
	if d := RetryAfter(resp); d > 0 {
		t.Fatalf("unexpected retry-after value %d", d)
	}
	resp.Header.Set(HeaderRetryAfter, "300")
	d := RetryAfter(resp)
	if d <= 0 {
		t.Fatal("expected retry-after value from seconds")
	}
	if d != 300*time.Second {
		t.Fatalf("expected 300 seconds, got %d", d/time.Second)
	}
	atDate := time.Now().Add(600 * time.Second)
	resp.Header.Set(HeaderRetryAfter, atDate.Format(time.RFC1123))
	d = RetryAfter(resp)
	if d <= 0 {
		t.Fatal("expected retry-after value from date")
	}
	// d will not be exactly 600 seconds but it will be close
	if s := d / time.Second; s < 598 || s > 602 {
		t.Fatalf("expected ~600 seconds, got %d", s)
	}
	resp.Header.Set(HeaderRetryAfter, "invalid")
	if d = RetryAfter(resp); d != 0 {
		t.Fatalf("expected zero for invalid value, got %d", d)
	}
	// verify that the ms-granularity headers are preferred
	resp.Header = http.Header{}
	resp.Header.Set(HeaderRetryAfterMS, "500")
	require.Equal(t, time.Duration(500)*time.Millisecond, RetryAfter(resp))
	resp.Header = http.Header{}
	resp.Header.Set(HeaderXMSRetryAfterMS, "400")
	require.Equal(t, time.Duration(400)*time.Millisecond, RetryAfter(resp))
	resp.Header = http.Header{}
	resp.Header.Set(HeaderRetryAfterMS, "500")
	resp.Header.Set(HeaderXMSRetryAfterMS, "400")
	resp.Header.Set(HeaderRetryAfter, "300")
	require.Equal(t, time.Duration(500)*time.Millisecond, RetryAfter(resp))
	resp.Header = http.Header{}
	resp.Header.Set(HeaderXMSRetryAfterMS, "400")
	resp.Header.Set(HeaderRetryAfter, "300")
	require.Equal(t, time.Duration(400)*time.Millisecond, RetryAfter(resp))
	resp.Header = http.Header{}
	resp.Header.Set(HeaderRetryAfterMS, "invalid")
	require.Zero(t, RetryAfter(resp))
	resp.Header = http.Header{}
	resp.Header.Set(HeaderXMSRetryAfterMS, "invalid")
	require.Zero(t, RetryAfter(resp))
}

func TestTypeOfT(t *testing.T) {
	if tt := TypeOfT[bool](); tt != reflect.TypeOf(true) {
		t.Fatalf("unexpected type %s", tt)
	}
	if tt := TypeOfT[int32](); tt == reflect.TypeOf(3.14) {
		t.Fatal("didn't expect types to match")
	}
}

func TestTransportFunc(t *testing.T) {
	resp, err := TransportFunc(func(req *http.Request) (*http.Response, error) {
		return nil, nil
	}).Do(nil)
	require.Nil(t, resp)
	require.NoError(t, err)
}

func TestValidateModVer(t *testing.T) {
	require.NoError(t, ValidateModVer("v1.2.3"))
	require.NoError(t, ValidateModVer("v1.2.3-beta.1"))
	require.Error(t, ValidateModVer("1.2.3"))
	require.Error(t, ValidateModVer("v1.2"))
}

func TestContextWithDeniedValues(t *testing.T) {
	type testKey struct{}
	const value = "value"

	ctx := context.WithValue(context.Background(), testKey{}, value)
	ctx = context.WithValue(ctx, CtxAPINameKey{}, value)
	ctx = context.WithValue(ctx, CtxWithCaptureResponse{}, value)
	ctx = context.WithValue(ctx, CtxWithHTTPHeaderKey{}, value)
	ctx = context.WithValue(ctx, CtxWithRetryOptionsKey{}, value)
	ctx = context.WithValue(ctx, CtxWithTracingTracer{}, value)
	ctx = &ContextWithDeniedValues{Context: ctx}

	require.Nil(t, ctx.Value(CtxAPINameKey{}))
	require.Nil(t, ctx.Value(CtxWithCaptureResponse{}))
	require.Nil(t, ctx.Value(CtxWithHTTPHeaderKey{}))
	require.Nil(t, ctx.Value(CtxWithRetryOptionsKey{}))
	require.Nil(t, ctx.Value(CtxWithTracingTracer{}))
	require.NotNil(t, ctx.Value(testKey{}))
}
