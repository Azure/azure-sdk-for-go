//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestDownloadBody(t *testing.T) {
	const message = "downloaded"
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte(message)))
	// download policy is automatically added during pipeline construction
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	payload, err := Payload(resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payload) == 0 {
		t.Fatal("missing payload")
	}
	if string(payload) != message {
		t.Fatalf("unexpected response: %s", string(payload))
	}
}

func TestSkipBodyDownload(t *testing.T) {
	const message = "not downloaded"
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte(message)))
	// download policy is automatically added during pipeline construction
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	SkipBodyDownload(req)
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	payload, err := Payload(resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(payload) != message {
		t.Fatalf("unexpected body: %s", string(payload))
	}
}

func TestDownloadBodyFail(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBodyReadError())
	// download policy is automatically added during pipeline construction
	pl := newTestPipeline(&policy.ClientOptions{
		Transport: srv,
		Retry: policy.RetryOptions{
			RetryDelay: 10 * time.Millisecond,
		},
	})
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	payload, err := Payload(resp)
	if err == nil {
		t.Fatalf("expected an error")
	}
	if payload != nil {
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
	pl := newTestPipeline(&policy.ClientOptions{Retry: *testRetryOptions(), Transport: srv})
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	payload, err := Payload(resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payload) == 0 {
		t.Fatal("missing payload")
	}
	if string(payload) != message {
		t.Fatalf("unexpected response: %s", string(payload))
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
	pl := newTestPipeline(&policy.ClientOptions{Retry: *testRetryOptions(), Transport: srv})
	req, err := NewRequest(context.Background(), http.MethodDelete, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	payload, err := Payload(resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payload) == 0 {
		t.Fatal("missing payload")
	}
	if string(payload) != message {
		t.Fatalf("unexpected response: %s", string(payload))
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
	pl := newTestPipeline(&policy.ClientOptions{Retry: *testRetryOptions(), Transport: srv})
	req, err := NewRequest(context.Background(), http.MethodPut, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	payload, err := Payload(resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payload) == 0 {
		t.Fatal("missing payload")
	}
	if string(payload) != message {
		t.Fatalf("unexpected response: %s", string(payload))
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
	pl := newTestPipeline(&policy.ClientOptions{Retry: *testRetryOptions(), Transport: srv})
	req, err := NewRequest(context.Background(), http.MethodPatch, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if _, ok := err.(*bodyDownloadError); !ok {
		t.Fatal("expected *bodyDownloadError type")
	}
	payload, err := Payload(resp)
	if err == nil {
		t.Fatalf("expected an error")
	}
	if len(payload) != 0 {
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
	pl := newTestPipeline(&policy.ClientOptions{Retry: *testRetryOptions(), Transport: srv})
	req, err := NewRequest(context.Background(), http.MethodPost, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if s := err.Error(); s != "body download policy: mock read failure" {
		t.Fatalf("unexpected error message: %s", s)
	}
	payload, err := Payload(resp)
	if err == nil {
		t.Fatalf("expected an error")
	}
	if len(payload) != 0 {
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
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	SkipBodyDownload(req)
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	payload, err := Payload(resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payload) == 0 {
		t.Fatal("missing payload")
	}
	if string(payload) != message {
		t.Fatalf("unexpected response: %s", string(payload))
	}
}

func TestReadBodyAfterSeek(t *testing.T) {
	const message = "downloaded"
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(message)))
	srv.AppendResponse(mock.WithBody([]byte(message)))
	// download policy is automatically added during pipeline construction
	pl := newTestPipeline(&policy.ClientOptions{Retry: *testRetryOptions(), Transport: srv})
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	payload, err := Payload(resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(payload) != message {
		t.Fatal("incorrect payload")
	}
	nb, ok := resp.Body.(io.ReadSeekCloser)
	if !ok {
		t.Fatalf("unexpected body type: %t", resp.Body)
	}
	i, err := nb.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if i != 0 {
		t.Fatalf("did not seek correctly")
	}
	i, err = nb.Seek(5, io.SeekCurrent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if i != 5 {
		t.Fatalf("did not seek correctly")
	}
	i, err = nb.Seek(5, io.SeekCurrent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if i != 10 {
		t.Fatalf("did not seek correctly")
	}
}

func TestBodyDownloadError(t *testing.T) {
	bde := &bodyDownloadError{err: io.EOF}
	if !errors.Is(bde, io.EOF) {
		t.Fatal("unwrap should provide inner error")
	}
}
