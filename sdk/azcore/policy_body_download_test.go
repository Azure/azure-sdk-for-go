// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestDownloadBody(t *testing.T) {
	const message = "downloaded"
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte(message)))
	// download policy is automatically added during pipeline construction
	pl := NewPipeline(srv)
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodGet, srv.URL()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.payload()) == 0 {
		t.Fatal("missing payload")
	}
	if string(resp.payload()) != message {
		t.Fatalf("unexpected response: %s", string(resp.payload()))
	}
}

func TestSkipBodyDownload(t *testing.T) {
	const message = "not downloaded"
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte(message)))
	// download policy is automatically added during pipeline construction
	pl := NewPipeline(srv)
	req := NewRequest(http.MethodGet, srv.URL())
	req.SkipBodyDownload()
	resp, err := pl.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.payload()) > 0 {
		t.Fatalf("unexpected download: %s", string(resp.payload()))
	}
}

func TestDownloadBodyFail(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBodyReadError())
	// download policy is automatically added during pipeline construction
	pl := NewPipeline(srv)
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodGet, srv.URL()))
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp.payload() != nil {
		t.Fatal("expected nil payload")
	}
}

func TestDownloadBodyWithRetryGet(t *testing.T) {
	const message = "downloaded"
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBodyReadError())
	srv.AppendResponse(mock.WithBodyReadError())
	srv.AppendResponse(mock.WithBody([]byte(message)))
	// download policy is automatically added during pipeline construction
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodGet, srv.URL()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.payload()) == 0 {
		t.Fatal("missing payload")
	}
	if string(resp.payload()) != message {
		t.Fatalf("unexpected response: %s", string(resp.payload()))
	}
	if r := srv.Requests(); r != 3 {
		t.Fatalf("expected %d requests, got %d", 3, r)
	}
}

func TestDownloadBodyWithRetryDelete(t *testing.T) {
	const message = "downloaded"
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBodyReadError())
	srv.AppendResponse(mock.WithBodyReadError())
	srv.AppendResponse(mock.WithBody([]byte(message)))
	// download policy is automatically added during pipeline construction
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodDelete, srv.URL()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.payload()) == 0 {
		t.Fatal("missing payload")
	}
	if string(resp.payload()) != message {
		t.Fatalf("unexpected response: %s", string(resp.payload()))
	}
	if r := srv.Requests(); r != 3 {
		t.Fatalf("expected %d requests, got %d", 3, r)
	}
}

func TestDownloadBodyWithRetryPut(t *testing.T) {
	const message = "downloaded"
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBodyReadError())
	srv.AppendResponse(mock.WithBodyReadError())
	srv.AppendResponse(mock.WithBody([]byte(message)))
	// download policy is automatically added during pipeline construction
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodPut, srv.URL()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.payload()) == 0 {
		t.Fatal("missing payload")
	}
	if string(resp.payload()) != message {
		t.Fatalf("unexpected response: %s", string(resp.payload()))
	}
	if r := srv.Requests(); r != 3 {
		t.Fatalf("expected %d requests, got %d", 3, r)
	}
}

func TestDownloadBodyWithRetryPatch(t *testing.T) {
	const message = "downloaded"
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBodyReadError())
	srv.AppendResponse(mock.WithBodyReadError())
	srv.AppendResponse(mock.WithBody([]byte(message)))
	// download policy is automatically added during pipeline construction
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodPatch, srv.URL()))
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if _, ok := err.(*bodyDownloadError); !ok {
		t.Fatal("expected *bodyDownloadError type")
	}
	if len(resp.payload()) != 0 {
		t.Fatal("unexpected payload")
	}
	// should be only one request, no retires
	if r := srv.Requests(); r != 1 {
		t.Fatalf("expected %d requests, got %d", 1, r)
	}
}

func TestDownloadBodyWithRetryPost(t *testing.T) {
	const message = "downloaded"
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBodyReadError())
	srv.AppendResponse(mock.WithBodyReadError())
	srv.AppendResponse(mock.WithBody([]byte(message)))
	// download policy is automatically added during pipeline construction
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodPost, srv.URL()))
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if _, ok := err.(*bodyDownloadError); !ok {
		t.Fatal("expected *bodyDownloadError type")
	}
	if len(resp.payload()) != 0 {
		t.Fatal("unexpected payload")
	}
	// should be only one request, no retires
	if r := srv.Requests(); r != 1 {
		t.Fatalf("expected %d requests, got %d", 1, r)
	}
}

func TestSkipBodyDownloadWith400(t *testing.T) {
	const message = "error should be downloaded"
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusBadRequest), mock.WithBody([]byte(message)))
	// download policy is automatically added during pipeline construction
	pl := NewPipeline(srv)
	req := NewRequest(http.MethodGet, srv.URL())
	req.SkipBodyDownload()
	resp, err := pl.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.payload()) == 0 {
		t.Fatal("missing payload")
	}
	if string(resp.payload()) != message {
		t.Fatalf("unexpected response: %s", string(resp.payload()))
	}
}
