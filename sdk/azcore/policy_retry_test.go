// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func testRetryOptions() *RetryOptions {
	def := DefaultRetryOptions()
	def.RetryDelay = 20 * time.Millisecond
	return &def
}

func TestRetryPolicySuccess(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))
	pl := NewPipeline(srv, NewRetryPolicy(nil))
	req := NewRequest(http.MethodGet, srv.URL())
	body := newRewindTrackingBody("stuff")
	req.SetBody(body)
	resp, err := pl.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if body.rcount > 0 {
		t.Fatalf("unexpected rewind count: %d", body.rcount)
	}
	if !body.closed {
		t.Fatal("request body wasn't closed")
	}
}

func TestRetryPolicyFailOnStatusCode(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusInternalServerError))
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
	req := NewRequest(http.MethodGet, srv.URL())
	body := newRewindTrackingBody("stuff")
	req.SetBody(body)
	resp, err := pl.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if r := srv.Requests(); r != defaultMaxRetries+1 {
		t.Fatalf("wrong request count, got %d expected %d", r, defaultMaxRetries+1)
	}
	if body.rcount != defaultMaxRetries {
		t.Fatalf("unexpected rewind count: %d", body.rcount)
	}
	if !body.closed {
		t.Fatal("request body wasn't closed")
	}
}

func TestRetryPolicySuccessWithRetry(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	srv.AppendResponse()
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
	req := NewRequest(http.MethodGet, srv.URL())
	body := newRewindTrackingBody("stuff")
	req.SetBody(body)
	resp, err := pl.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if r := srv.Requests(); r != 3 {
		t.Fatalf("wrong retry count, got %d expected %d", r, 3)
	}
	if body.rcount != 2 {
		t.Fatalf("unexpected rewind count: %d", body.rcount)
	}
	if !body.closed {
		t.Fatal("request body wasn't closed")
	}
}

func TestRetryPolicyFailOnError(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	fakeErr := errors.New("bogus error")
	srv.SetError(fakeErr)
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
	req := NewRequest(http.MethodPost, srv.URL())
	body := newRewindTrackingBody("stuff")
	req.SetBody(body)
	resp, err := pl.Do(context.Background(), req)
	if !errors.Is(err, fakeErr) {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Fatal("unexpected response")
	}
	if r := srv.Requests(); r != defaultMaxRetries+1 {
		t.Fatalf("wrong request count, got %d expected %d", r, defaultMaxRetries+1)
	}
	if body.rcount != defaultMaxRetries {
		t.Fatalf("unexpected rewind count: %d", body.rcount)
	}
	if !body.closed {
		t.Fatal("request body wasn't closed")
	}
}

func TestRetryPolicySuccessWithRetryComplex(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendError(errors.New("bogus error"))
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
	req := NewRequest(http.MethodGet, srv.URL())
	body := newRewindTrackingBody("stuff")
	req.SetBody(body)
	resp, err := pl.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if r := srv.Requests(); r != defaultMaxRetries+1 {
		t.Fatalf("wrong request count, got %d expected %d", r, defaultMaxRetries+1)
	}
	if body.rcount != defaultMaxRetries {
		t.Fatalf("unexpected rewind count: %d", body.rcount)
	}
	if !body.closed {
		t.Fatal("request body wasn't closed")
	}
}

func TestRetryPolicyRequestTimedOut(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetError(errors.New("bogus error"))
	pl := NewPipeline(srv, NewRetryPolicy(nil))
	req := NewRequest(http.MethodPost, srv.URL())
	body := newRewindTrackingBody("stuff")
	req.SetBody(body)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	resp, err := pl.Do(ctx, req)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Fatal("unexpected response")
	}
	if body.rcount > 0 {
		t.Fatalf("unexpected rewind count: %d", body.rcount)
	}
	if !body.closed {
		t.Fatal("request body wasn't closed")
	}
}

type fatalError struct {
	s string
}

func (f fatalError) Error() string {
	return f.s
}

func (f fatalError) IsNotRetriable() bool {
	return true
}

func TestRetryPolicyIsNotRetriable(t *testing.T) {
	theErr := fatalError{s: "it's dead Jim"}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendError(theErr)
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
	_, err := pl.Do(context.Background(), NewRequest(http.MethodGet, srv.URL()))
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if !errors.Is(err, theErr) {
		t.Fatalf("unexpected error type: got %v wanted %v", err, theErr)
	}
	if r := srv.Requests(); r != 2 {
		t.Fatalf("wrong retry count, got %d expected %d", r, 3)
	}
}

func TestWithRetryOptions(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.RepeatResponse(9, mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	defaultOptions := testRetryOptions()
	pl := NewPipeline(srv, NewRetryPolicy(defaultOptions))
	customOptions := *defaultOptions
	customOptions.MaxRetries = 10
	customOptions.MaxRetryDelay = 200 * time.Millisecond
	retryCtx := WithRetryOptions(context.Background(), customOptions)
	req := NewRequest(http.MethodGet, srv.URL())
	body := newRewindTrackingBody("stuff")
	req.SetBody(body)
	resp, err := pl.Do(retryCtx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if body.rcount != int(customOptions.MaxRetries-1) {
		t.Fatalf("unexpected rewind count: %d", body.rcount)
	}
	if !body.closed {
		t.Fatal("request body wasn't closed")
	}
}

// TODO: add test for retry failing to read response body

// TODO: add test for per-retry timeout failed but e2e succeeded

func newRewindTrackingBody(s string) *rewindTrackingBody {
	// there are two rewinds that happen before rewinding for a retry
	// 1. to get the body's size in SetBody()
	// 2. the first call to Do() in the retry policy
	// to offset this we init rcount with -2 so rcount is only > 0 on a rewind due to a retry
	return &rewindTrackingBody{
		body:   strings.NewReader(s),
		rcount: -2,
	}
}

// used to track the number of times a request body has been rewound
type rewindTrackingBody struct {
	body   *strings.Reader
	closed bool // indicates if the body was closed
	rcount int  // number of times a rewind happened
}

func (r *rewindTrackingBody) Close() error {
	r.closed = true
	return nil
}

func (r *rewindTrackingBody) Read(b []byte) (int, error) {
	return r.body.Read(b)
}

func (r *rewindTrackingBody) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		r.rcount++
	}
	return r.body.Seek(offset, whence)
}
