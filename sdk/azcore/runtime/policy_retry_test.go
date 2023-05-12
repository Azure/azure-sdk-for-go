//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"bytes"
	"context"
	"errors"
	"io"
	"math"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func testRetryOptions() *policy.RetryOptions {
	return &policy.RetryOptions{
		RetryDelay: time.Millisecond,
	}
}

func TestRetryPolicySuccess(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))
	pl := exported.NewPipeline(srv, NewRetryPolicy(nil))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := newRewindTrackingBody("stuff")
	if err := req.SetBody(body, "text/plain"); err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
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
	pl := exported.NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := newRewindTrackingBody("stuff")
	if err := req.SetBody(body, "text/plain"); err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
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

func TestRetryPolicyFailOnStatusCodeRespBodyPreserved(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	const respBody = "response body"
	srv.SetResponse(mock.WithStatusCode(http.StatusInternalServerError), mock.WithBody([]byte(respBody)))
	// add a per-request policy that reads and restores the request body.
	// this is to simulate how something like httputil.DumpRequest works.
	pl := exported.NewPipeline(srv, exported.PolicyFunc(func(r *policy.Request) (*http.Response, error) {
		b, err := io.ReadAll(r.Raw().Body)
		if err != nil {
			t.Fatal(err)
		}
		r.Raw().Body = io.NopCloser(bytes.NewReader(b))
		return r.Next()
	}), NewRetryPolicy(testRetryOptions()))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := newRewindTrackingBody("stuff")
	if err := req.SetBody(body, "text/plain"); err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
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
	// ensure response body hasn't been drained
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != respBody {
		t.Fatalf("unexpected response body: %s", string(b))
	}
}

func TestRetryPolicySuccessWithRetry(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	srv.AppendResponse()
	pl := exported.NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := newRewindTrackingBody("stuff")
	if err := req.SetBody(body, "text/plain"); err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
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

func TestRetryPolicySuccessRetryWithNilResponse(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	srv.AppendResponse()
	nilInjector := &nilRespInjector{
		t: srv,
		r: []int{2}, // send a nil on the second request
	}
	pl := exported.NewPipeline(nilInjector, NewRetryPolicy(testRetryOptions()))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := newRewindTrackingBody("stuff")
	if err := req.SetBody(body, "text/plain"); err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if r := srv.Requests(); r != 3 {
		t.Fatalf("wrong retry count, got %d expected %d", r, 3)
	}
	if body.rcount != 3 {
		t.Fatalf("unexpected rewind count: %d", body.rcount)
	}
	if !body.closed {
		t.Fatal("request body wasn't closed")
	}
}

func TestRetryPolicyNoRetries(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	srv.AppendResponse()
	pl := exported.NewPipeline(srv, NewRetryPolicy(&policy.RetryOptions{MaxRetries: -1}))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusRequestTimeout {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if r := srv.Requests(); r != 1 {
		t.Fatalf("wrong try count, got %d expected %d", r, 1)
	}
}

func TestRetryPolicyUnlimitedRetryDelay(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	srv.AppendResponse()
	opt := testRetryOptions()
	opt.MaxRetryDelay = -1
	pl := exported.NewPipeline(srv, NewRetryPolicy(opt))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if r := srv.Requests(); r != 3 {
		t.Fatalf("wrong try count, got %d expected %d", r, 3)
	}
}

func TestRetryPolicyFailOnError(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	fakeErr := errors.New("bogus error")
	srv.SetError(fakeErr)
	pl := exported.NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
	req, err := NewRequest(context.Background(), http.MethodPost, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := newRewindTrackingBody("stuff")
	if err := req.SetBody(body, "text/plain"); err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
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
	pl := exported.NewPipeline(srv, exported.PolicyFunc(includeResponsePolicy), NewRetryPolicy(testRetryOptions()))
	var respFromCtx *http.Response
	ctxWithResp := WithCaptureResponse(context.Background(), &respFromCtx)
	req, err := NewRequest(ctxWithResp, http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := newRewindTrackingBody("stuff")
	if err := req.SetBody(body, "text/plain"); err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if respFromCtx != resp {
		t.Fatal("response from context doesn't match returned response")
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
	pl := exported.NewPipeline(srv, NewRetryPolicy(nil))
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()
	req, err := NewRequest(ctx, http.MethodPost, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := newRewindTrackingBody("stuff")
	if err := req.SetBody(body, "text/plain"); err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
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

func (f fatalError) NonRetriable() {
	// marker method
}

var _ errorinfo.NonRetriable = (*fatalError)(nil)

func TestRetryPolicyIsNotRetriable(t *testing.T) {
	theErr := fatalError{s: "it's dead Jim"}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendError(theErr)
	pl := exported.NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = pl.Do(req)
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
	ctx := WithRetryOptions(context.Background(), policy.RetryOptions{
		MaxRetries: math.MaxInt32,
	})
	if ctx == nil {
		t.Fatal("nil context")
	}
	raw := ctx.Value(shared.CtxWithRetryOptionsKey{})
	opts, ok := raw.(policy.RetryOptions)
	if !ok {
		t.Fatalf("unexpected type %T", raw)
	}
	if opts.MaxRetries != math.MaxInt32 {
		t.Fatalf("unexpected value %d", opts.MaxRetries)
	}
}

func TestWithRetryOptionsE2E(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.RepeatResponse(9, mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	defaultOptions := testRetryOptions()
	pl := exported.NewPipeline(srv, NewRetryPolicy(defaultOptions))
	customOptions := *defaultOptions
	customOptions.MaxRetries = 10
	customOptions.MaxRetryDelay = 200 * time.Millisecond
	retryCtx := WithRetryOptions(context.Background(), customOptions)
	req, err := NewRequest(retryCtx, http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := newRewindTrackingBody("stuff")
	if err := req.SetBody(body, "text/plain"); err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
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

func TestRetryPolicyFailOnErrorNoDownload(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	fakeErr := errors.New("bogus error")
	srv.SetError(fakeErr)
	pl := exported.NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
	req, err := NewRequest(context.Background(), http.MethodPost, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	SkipBodyDownload(req)
	resp, err := pl.Do(req)
	if !errors.Is(err, fakeErr) {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Fatal("unexpected response")
	}
	if r := srv.Requests(); r != defaultMaxRetries+1 {
		t.Fatalf("wrong request count, got %d expected %d", r, defaultMaxRetries+1)
	}
}

func TestRetryPolicySuccessNoDownload(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte("response body")))
	pl := exported.NewPipeline(srv, NewRetryPolicy(nil))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	SkipBodyDownload(req)
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestRetryPolicySuccessNoDownloadNoBody(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))
	pl := exported.NewPipeline(srv, NewRetryPolicy(nil))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	SkipBodyDownload(req)
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestRetryPolicySuccessWithRetryReadingResponse(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBodyReadError())
	srv.AppendResponse(mock.WithBodyReadError())
	srv.AppendResponse()
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv, Retry: *testRetryOptions()})
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := newRewindTrackingBody("stuff")
	if err := req.SetBody(body, "text/plain"); err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
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

func TestRetryPolicyRequestTimedOutTooSlow(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithSlowResponse(5 * time.Second))
	pl := exported.NewPipeline(srv, NewRetryPolicy(nil))
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req, err := NewRequest(ctx, http.MethodPost, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := newRewindTrackingBody("stuff")
	if err := req.SetBody(body, "text/plain"); err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
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

func TestRetryPolicySuccessWithPerTryTimeout(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithSlowResponse(5 * time.Second))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	opt := testRetryOptions()
	opt.TryTimeout = 1 * time.Second
	pl := exported.NewPipeline(srv, NewRetryPolicy(opt))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body := newRewindTrackingBody("stuff")
	if err := req.SetBody(body, "text/plain"); err != nil {
		t.Fatal(err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if body.rcount != 1 {
		// should have been rewound once due to per-try timeout
		t.Fatalf("unexpected rewind count: %d", body.rcount)
	}
	if !body.closed {
		t.Fatal("request body wasn't closed")
	}
}

func TestRetryPolicySuccessWithPerTryTimeoutNoRetry(t *testing.T) {
	// ensure that the size of the payload is larger than the read buffer
	// on the underlying transport (defaults to 4KB).  this will ensure
	// that the writes will hit the network again so the bug will repro.
	const bodySize = 1024 * 8
	largeBody := make([]byte, bodySize)
	for i := 0; i < bodySize; i++ {
		largeBody[i] = byte(i % 256)
	}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody(largeBody))
	opt := testRetryOptions()
	opt.TryTimeout = 10 * time.Second
	pl := exported.NewPipeline(srv, NewRetryPolicy(opt))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	require.NoError(t, resp.Body.Close())
	require.Equal(t, largeBody, body)
}

func TestRetryPolicySuccessWithPerTryTimeoutNoRetryWithBodyDownload(t *testing.T) {
	// ensure that the size of the payload is larger than the read buffer
	// on the underlying transport (defaults to 4KB).  this will ensure
	// that the writes will hit the network again so the bug will repro.
	const bodySize = 1024 * 8
	largeBody := make([]byte, bodySize)
	for i := 0; i < bodySize; i++ {
		largeBody[i] = byte(i % 256)
	}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody(largeBody))
	opt := testRetryOptions()
	opt.TryTimeout = 10 * time.Second
	pl := exported.NewPipeline(srv, NewRetryPolicy(opt), exported.PolicyFunc(bodyDownloadPolicy))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	require.NoError(t, resp.Body.Close())
	require.Equal(t, largeBody, body)
}

func TestRetryPolicyWithShouldRetryNoRetry(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))

	pl := exported.NewPipeline(srv, NewRetryPolicy(&policy.RetryOptions{
		RetryDelay: time.Millisecond,
		ShouldRetry: func(r *http.Response, err error) bool {
			return r.StatusCode != http.StatusRequestTimeout
		},
	}))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusRequestTimeout, resp.StatusCode)
	require.EqualValues(t, 1, srv.Requests())
}

func TestRetryPolicyWithShouldRetryRetry(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendResponse()

	shouldRetryCalled := false
	pl := exported.NewPipeline(srv, NewRetryPolicy(&policy.RetryOptions{
		RetryDelay: time.Millisecond,
		ShouldRetry: func(r *http.Response, err error) bool {
			shouldRetryCalled = true
			return r.StatusCode == http.StatusRequestTimeout
		},
	}))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.True(t, shouldRetryCalled)
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
	require.EqualValues(t, 2, srv.Requests())
}

func TestPipelineNoRetryOn429(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// initial response is throttling with a long retry-after delay, it should not trigger a retry
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests), mock.WithHeader(shared.HeaderRetryAfter, "300"))
	perRetryPolicy := countingPolicy{}
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)
	pl := exported.NewPipeline(srv, NewRetryPolicy(nil), &perRetryPolicy)
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	require.Equal(t, 1, perRetryPolicy.count)
}

func TestPipelineRetryOn429(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests), mock.WithHeader(shared.HeaderRetryAfter, "1"))
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests), mock.WithHeader(shared.HeaderRetryAfter, "1"))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	perRetryPolicy := countingPolicy{}
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)
	opt := testRetryOptions()
	pl := exported.NewPipeline(srv, NewRetryPolicy(opt), &perRetryPolicy)
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 3, perRetryPolicy.count)
}

type readSeekerTracker struct {
	readCalled bool
	seekCalled bool
}

func (r *readSeekerTracker) Read([]byte) (int, error) {
	r.readCalled = true
	return 0, nil
}

func (r *readSeekerTracker) Seek(int64, int) (int64, error) {
	r.seekCalled = true
	return 0, nil
}

func TestRetryableRequestBodyNoCloser(t *testing.T) {
	tr := &readSeekerTracker{}
	rr := &retryableRequestBody{tr}
	_, err := rr.Read(nil)
	require.NoError(t, err)
	_, err = rr.Seek(0, 0)
	require.NoError(t, err)
	require.NoError(t, rr.Close())
	require.NoError(t, rr.realClose())
	require.True(t, tr.readCalled)
	require.True(t, tr.seekCalled)
}

type readSeekCloseerTracker struct {
	readCalled  bool
	seekCalled  bool
	closeCalled bool
}

func (r *readSeekCloseerTracker) Read([]byte) (int, error) {
	r.readCalled = true
	return 0, nil
}

func (r *readSeekCloseerTracker) Seek(int64, int) (int64, error) {
	r.seekCalled = true
	return 0, nil
}

func (r *readSeekCloseerTracker) Close() error {
	r.closeCalled = true
	return nil
}

func TestRetryableRequestBodyWithCloser(t *testing.T) {
	tr := &readSeekCloseerTracker{}
	rr := &retryableRequestBody{tr}
	_, err := rr.Read(nil)
	require.NoError(t, err)
	_, err = rr.Seek(0, 0)
	require.NoError(t, err)
	require.NoError(t, rr.Close())
	require.False(t, tr.closeCalled)
	require.True(t, tr.readCalled)
	require.True(t, tr.seekCalled)
	require.NoError(t, rr.realClose())
	require.True(t, tr.closeCalled)
}

func TestRetryPolicySuccessWithRetryPreserveHeaders(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendResponse()
	pl := exported.NewPipeline(srv, NewRetryPolicy(testRetryOptions()), exported.PolicyFunc(challengeLikePolicy))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	require.NoError(t, err)
	body := newRewindTrackingBody("stuff")
	require.NoError(t, req.SetBody(body, "text/plain"))
	resp, err := pl.Do(req)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
	require.EqualValues(t, 2, srv.Requests())
	require.EqualValues(t, 1, body.rcount)
	require.True(t, body.closed)
}

func challengeLikePolicy(req *policy.Request) (*http.Response, error) {
	if req.Body() == nil {
		return nil, errors.New("request body wasn't restored")
	}
	if req.Raw().Header.Get("content-type") != "text/plain" {
		return nil, errors.New("content-type header wasn't restored")
	}

	// remove the body and header. the retry policy should restore them
	if err := req.SetBody(nil, ""); err != nil {
		return nil, err
	}
	return req.Next()
}

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

// used to inject a nil response
type nilRespInjector struct {
	t policy.Transporter
	c int   // the current request number
	r []int // the list of request numbers to return a nil response (one-based)
}

func (n *nilRespInjector) Do(req *http.Request) (*http.Response, error) {
	n.c++
	// check if current request number n.c is in n.r
	for _, v := range n.r {
		if v == n.c {
			return nil, nil
		}
	}
	return n.t.Do(req)
}
