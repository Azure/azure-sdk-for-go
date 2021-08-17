//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func testRetryOptions() *RetryOptions {
	def := RetryOptions{}
	def.RetryDelay = 20 * time.Millisecond
	return &def
}

func TestRetryPolicySuccess(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))
	pl := NewPipeline(srv, NewRetryPolicy(nil))
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
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
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
	pl := NewPipeline(srv, policyFunc(func(r *Request) (*http.Response, error) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		r.Body = ioutil.NopCloser(bytes.NewReader(b))
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
	b, err := ioutil.ReadAll(resp.Body)
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
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
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
	pl := NewPipeline(nilInjector, NewRetryPolicy(testRetryOptions()))
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
	pl := NewPipeline(srv, NewRetryPolicy(&RetryOptions{MaxRetries: -1}))
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
	pl := NewPipeline(srv, NewRetryPolicy(opt))
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
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
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
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
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

type fatalError struct {
	s string
}

func (f fatalError) Error() string {
	return f.s
}

func (f fatalError) NonRetriable() {
	// marker method
}

var _ NonRetriableError = (*fatalError)(nil)

func TestRetryPolicyIsNotRetriable(t *testing.T) {
	theErr := fatalError{s: "it's dead Jim"}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusRequestTimeout))
	srv.AppendError(theErr)
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
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
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
	req, err := NewRequest(context.Background(), http.MethodPost, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.SkipBodyDownload()
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
	pl := NewPipeline(srv, NewRetryPolicy(nil))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.SkipBodyDownload()
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
	pl := NewPipeline(srv, NewRetryPolicy(nil))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.SkipBodyDownload()
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
	pl := NewPipeline(srv, NewRetryPolicy(testRetryOptions()))
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
	pl := NewPipeline(srv, NewRetryPolicy(nil))
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
	pl := NewPipeline(srv, NewRetryPolicy(opt))
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
	t Transporter
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
