// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/mock"
)

const retryDelay = 20 * time.Millisecond

func TestRetryPolicySuccess(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))
	pl := NewPipeline(srv, NewRetryPolicy(RetryOptions{}))
	req := pl.NewRequest(http.MethodGet, srv.URL())
	resp, err := req.Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

func TestRetryPolicyFailOnStatusCode(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusInternalServerError))
	pl := NewPipeline(srv, NewRetryPolicy(RetryOptions{RetryDelay: retryDelay}))
	req := pl.NewRequest(http.MethodGet, srv.URL())
	resp, err := req.Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if r := srv.Requests(); r != defaultMaxRetries {
		t.Fatalf("wrong retry count, got %d expected %d", r, defaultMaxRetries)
	}
}

func TestRetryPolicySuccessWithRetry(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	srv.AppendResponse()
	pl := NewPipeline(srv, NewRetryPolicy(RetryOptions{RetryDelay: retryDelay}))
	req := pl.NewRequest(http.MethodGet, srv.URL())
	resp, err := req.Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if r := srv.Requests(); r != 3 {
		t.Fatalf("wrong retry count, got %d expected %d", r, 3)
	}
}

func TestRetryPolicyFailOnError(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	fakeErr := errors.New("bogus error")
	srv.SetError(fakeErr)
	pl := NewPipeline(srv, NewRetryPolicy(RetryOptions{RetryDelay: retryDelay}))
	req := pl.NewRequest(http.MethodPost, srv.URL())
	resp, err := req.Do(context.Background())
	if !errors.Is(err, fakeErr) {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Fatal("unexpected response")
	}
	if r := srv.Requests(); r != defaultMaxRetries {
		t.Fatalf("wrong retry count, got %d expected %d", r, defaultMaxRetries)
	}
}

func TestRetryPolicySuccessWithRetryComplex(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendError(errors.New("bogus error"))
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	pl := NewPipeline(srv, NewRetryPolicy(RetryOptions{RetryDelay: retryDelay}))
	req := pl.NewRequest(http.MethodGet, srv.URL())
	resp, err := req.Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if r := srv.Requests(); r != defaultMaxRetries {
		t.Fatalf("wrong retry count, got %d expected %d", r, 3)
	}
}

func TestRetryPolicyRequestTimedOut(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetError(errors.New("bogus error"))
	pl := NewPipeline(srv, NewRetryPolicy(RetryOptions{}))
	req := pl.NewRequest(http.MethodPost, srv.URL())
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := req.Do(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Fatal("unexpected response")
	}
}

func TestRetryPolicyFixedDelaySuccessWithRetry(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	srv.AppendResponse()
	pl := NewPipeline(srv, NewRetryPolicy(RetryOptions{
		Policy:     RetryPolicyFixed,
		RetryDelay: retryDelay,
	}))
	req := pl.NewRequest(http.MethodGet, srv.URL())
	resp, err := req.Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if r := srv.Requests(); r != 3 {
		t.Fatalf("wrong retry count, got %d expected %d", r, 3)
	}
}
