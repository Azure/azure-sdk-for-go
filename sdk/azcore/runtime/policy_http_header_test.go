//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"net/http"
	"net/textproto"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func newTestPipeline(opts *policy.ClientOptions) Pipeline {
	return NewPipeline("testmodule", "v0.1.0", PipelineOptions{}, opts)
}

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

func TestAddCustomHTTPHeaderSuccess(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	const customHeader = "custom-header"
	const customValue = "custom-value"
	const preexistingHeader = "preexisting-header"
	const preexistingValue = "preexisting-value"
	srv.AppendResponse(mock.WithPredicate(func(r *http.Request) bool {
		// ensure preexisting header wasn't removed
		return r.Header.Get(customHeader) == customValue && r.Header.Get(preexistingHeader) == preexistingValue
	}), mock.WithStatusCode(http.StatusOK))
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
	// HTTP header policy is automatically added during pipeline construction
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	req, err := NewRequest(WithHTTPHeader(context.Background(), http.Header{
		customHeader: []string{customValue},
	}), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.Raw().Header.Set(preexistingHeader, preexistingValue)
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}

func TestAddCustomHTTPHeaderFail(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	const customHeader = "custom-header"
	const customValue = "custom-value"
	srv.AppendResponse(mock.WithPredicate(func(r *http.Request) bool {
		return r.Header.Get(customHeader) == customValue
	}), mock.WithStatusCode(http.StatusOK))
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
	// HTTP header policy is automatically added during pipeline construction
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}

func TestAddCustomHTTPHeaderOverwrite(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	const customHeader = "custom-header"
	const customValue = "custom-value"
	srv.AppendResponse(mock.WithPredicate(func(r *http.Request) bool {
		return r.Header.Get(customHeader) == customValue
	}), mock.WithStatusCode(http.StatusOK))
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
	// HTTP header policy is automatically added during pipeline construction
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	// overwrite the request ID with our own value
	req, err := NewRequest(WithHTTPHeader(context.Background(), http.Header{
		customHeader: []string{customValue},
	}), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}

func TestAddCustomHTTPHeaderMultipleValues(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	const customHeader = "custom-header"
	const customValue1 = "custom-value1"
	const customValue2 = "custom-value2"
	srv.AppendResponse(mock.WithPredicate(func(r *http.Request) bool {
		// Values() method is Go 1.14+
		//vals := r.Header.Values(customHeader)
		vals := r.Header[textproto.CanonicalMIMEHeaderKey(customHeader)]
		return vals[0] == customValue1 && vals[1] == customValue2
	}), mock.WithStatusCode(http.StatusOK))
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
	// HTTP header policy is automatically added during pipeline construction
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	// overwrite the request ID with our own value
	req, err := NewRequest(WithHTTPHeader(context.Background(), http.Header{
		customHeader: []string{customValue1, customValue2},
	}), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}
