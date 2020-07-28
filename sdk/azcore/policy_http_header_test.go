// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"net/http"
	"net/textproto"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

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
	pl := NewPipeline(srv)
	req := NewRequest(http.MethodGet, srv.URL())
	req.Header.Set(preexistingHeader, preexistingValue)
	resp, err := pl.Do(WithHTTPHeader(context.Background(), http.Header{
		customHeader: []string{customValue},
	}), req)
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
	const customValue = "custom-value"
	srv.AppendResponse(mock.WithPredicate(func(r *http.Request) bool {
		return r.Header.Get(xMsClientRequestID) == customValue
	}), mock.WithStatusCode(http.StatusOK))
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
	// HTTP header policy is automatically added during pipeline construction
	pl := NewPipeline(srv)
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodGet, srv.URL()))
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
	const customValue = "custom-value"
	srv.AppendResponse(mock.WithPredicate(func(r *http.Request) bool {
		return r.Header.Get(xMsClientRequestID) == customValue
	}), mock.WithStatusCode(http.StatusOK))
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
	// HTTP header policy is automatically added during pipeline construction
	pl := NewPipeline(srv)
	// overwrite the request ID with our own value
	resp, err := pl.Do(WithHTTPHeader(context.Background(), http.Header{
		xMsClientRequestID: []string{customValue},
	}), NewRequest(http.MethodGet, srv.URL()))
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
	pl := NewPipeline(srv)
	// overwrite the request ID with our own value
	resp, err := pl.Do(WithHTTPHeader(context.Background(), http.Header{
		customHeader: []string{customValue1, customValue2},
	}), NewRequest(http.MethodGet, srv.URL()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}
