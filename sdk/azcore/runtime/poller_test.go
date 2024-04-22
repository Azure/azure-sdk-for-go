//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers/async"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers/body"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/poller"
	"github.com/stretchr/testify/require"
)

type none struct{}

type widget struct {
	Size int `json:"size"`
}

func TestNewPollerFail(t *testing.T) {
	body, closed := mock.NewTrackedCloser(http.NoBody)
	p, err := NewPoller[widget](&http.Response{
		Body:       body,
		StatusCode: http.StatusBadRequest,
	}, newTestPipeline(nil), nil)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if p != nil {
		t.Fatal("expected nil poller")
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
}

func TestNewPollerFromResumeTokenFail(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{"invalid", "invalid"},
		{"empty", "{}"},
		{"wrong type", `{"type": 1}`},
		{"missing type", `{"type": "fake.poller"}`},
		{"mismatched type", `{"type": "faker.poller;opPoller"}`},
		{"malformed type", `{"type": "fake.poller;dummy"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := NewPollerFromResumeToken[widget](test.token, newTestPipeline(nil), nil)
			if err == nil {
				t.Fatal("unexpected nil error")
			}
			if p != nil {
				t.Fatal("expected nil poller")
			}
		})
	}
}

func TestLocPollerSimple(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))

	firstResp := &http.Response{
		StatusCode: http.StatusAccepted,
		Header: http.Header{
			"Location":    []string{srv.URL()},
			"Retry-After": []string{"1"},
		},
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp.Body = body
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	lro, err := NewPoller[none](firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	var respFromCtx *http.Response
	ctxWithResp := WithCaptureResponse(context.Background(), &respFromCtx)
	_, err = lro.PollUntilDone(ctxWithResp, &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	if respFromCtx.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", respFromCtx.StatusCode)
	}
}

func TestLocPollerWithWidget(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"size": 3}`)))

	firstResp := &http.Response{
		StatusCode: http.StatusAccepted,
		Header: http.Header{
			"Location":    []string{srv.URL()},
			"Retry-After": []string{"1"},
		},
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp.Body = body
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	lro, err := NewPoller[widget](firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	w, err := lro.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	if w.Size != 3 {
		t.Fatalf("unexpected widget size %d", w.Size)
	}
}

func TestLocPollerCancelled(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithStatusCode(http.StatusConflict), mock.WithBody([]byte(`{"error": "cancelled"}`)))

	firstResp := &http.Response{
		StatusCode: http.StatusAccepted,
		Header: http.Header{
			"Location":    []string{srv.URL()},
			"Retry-After": []string{"1"},
		},
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp.Body = body
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	lro, err := NewPoller[widget](firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	w, err := lro.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if _, ok := err.(*exported.ResponseError); !ok {
		t.Fatal("expected pollerError")
	}
	if w.Size != 0 {
		t.Fatalf("unexpected widget size %d", w.Size)
	}
	w, err = lro.Result(context.Background())
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if _, ok := err.(*exported.ResponseError); !ok {
		t.Fatal("expected pollerError")
	}
	if w.Size != 0 {
		t.Fatalf("unexpected widget size %d", w.Size)
	}
}

func TestLocPollerWithError(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendError(errors.New("oops"))

	firstResp := &http.Response{
		StatusCode: http.StatusAccepted,
		Header: http.Header{
			"Location":    []string{srv.URL()},
			"Retry-After": []string{"1"},
		},
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp.Body = body
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv, Retry: policy.RetryOptions{MaxRetries: -1}})
	lro, err := NewPoller[widget](firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	w, err := lro.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if e := err.Error(); e != "oops" {
		t.Fatalf("expected error %s", e)
	}
	if w.Size != 0 {
		t.Fatalf("unexpected widget size %d", w.Size)
	}
}

func TestLocPollerWithResumeToken(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	defer close()

	firstResp := &http.Response{
		StatusCode: http.StatusAccepted,
		Header: http.Header{
			"Location":    []string{srv.URL()},
			"Retry-After": []string{"1"},
		},
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp.Body = body
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	lro, err := NewPoller[none](firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	resp, err := lro.Poll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	if lro.Done() {
		t.Fatal("poller shouldn't be done yet")
	}
	_, err = lro.Result(context.Background())
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	tk, err := lro.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	lro, err = NewPollerFromResumeToken[none](tk, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = lro.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
}

func TestLocPollerWithTimeout(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithSlowResponse(5 * time.Second))
	defer close()

	firstResp := &http.Response{
		StatusCode: http.StatusAccepted,
		Header: http.Header{
			"Location": []string{srv.URL()},
		},
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp.Body = body
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	lro, err := NewPoller[none](firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	_, err = lro.PollUntilDone(ctx, &PollUntilDoneOptions{Frequency: time.Millisecond})
	cancel()
	if err == nil {
		t.Fatal("unexpected nil error")
	}
}

func TestOpPollerSimple(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted), mock.WithBody([]byte(`{ "status": "InProgress"}`)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted), mock.WithBody([]byte(`{ "status": "InProgress"}`)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{ "status": "Succeeded"}`)))
	defer close()

	reqURL, err := url.Parse(srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp := &http.Response{
		Body:       body,
		StatusCode: http.StatusAccepted,
		Header: http.Header{
			"Operation-Location": []string{srv.URL()},
			"Retry-After":        []string{"1"},
		},
		Request: &http.Request{
			Method: http.MethodDelete,
			URL:    reqURL,
		},
	}
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	lro, err := NewPoller[none](firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	_, err = lro.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
}

func TestOpPollerWithWidgetPUT(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted), mock.WithBody([]byte(`{"status": "InProgress"}`)), mock.WithHeader("Retry-After", "1"))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted), mock.WithBody([]byte(`{"status": "InProgress"}`)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"status": "Succeeded"}`)))
	// PUT and PATCH state that a final GET will happen
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"size": 2}`)))
	defer close()

	reqURL, err := url.Parse(srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp := &http.Response{
		Body:       body,
		StatusCode: http.StatusAccepted,
		Header: http.Header{
			"Operation-Location": []string{srv.URL()},
			"Retry-After":        []string{"1"},
		},
		Request: &http.Request{
			Method: http.MethodPut,
			URL:    reqURL,
		},
	}
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	lro, err := NewPoller[widget](firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	w, err := lro.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	if w.Size != 2 {
		t.Fatalf("unexpected widget size %d", w.Size)
	}
}

func TestOpPollerWithWidgetFinalGetError(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted), mock.WithBody([]byte(`{"status": "InProgress"}`)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"status": "Succeeded"}`)))
	// PUT and PATCH state that a final GET will happen
	// the first attempt at a final GET returns an error
	srv.AppendError(errorinfo.NonRetriableError(errors.New("failed attempt")))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"size": 2}`)))
	defer close()

	reqURL, err := url.Parse(srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp := &http.Response{
		Body:       body,
		StatusCode: http.StatusAccepted,
		Header: http.Header{
			"Operation-Location": []string{srv.URL()},
		},
		Request: &http.Request{
			Method: http.MethodPut,
			URL:    reqURL,
		},
	}
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	lro, err := NewPoller[widget](firstResp, pl, nil)
	require.Nil(t, err)
	require.True(t, closed(), "initial response body wasn't closed")

	resp, err := lro.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, resp.StatusCode)
	require.False(t, lro.Done())

	resp, err = lro.Poll(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.True(t, lro.Done())

	w, err := lro.Result(context.Background())
	require.Error(t, err)
	require.Empty(t, w)

	w, err = lro.Result(context.Background())
	require.NoError(t, err)
	require.Equal(t, w.Size, 2)
}

func TestOpPollerWithWidgetPOSTLocation(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted), mock.WithBody([]byte(`{"status": "InProgress"}`)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted), mock.WithBody([]byte(`{"status": "InProgress"}`)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"status": "Succeeded"}`)))
	// POST state that a final GET will happen from the URL provided in the Location header if available
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"size": 2}`)))
	defer close()

	reqURL, err := url.Parse(srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp := &http.Response{
		Body:       body,
		StatusCode: http.StatusAccepted,
		Header: http.Header{
			"Operation-Location": []string{srv.URL()},
			"Location":           []string{srv.URL()},
			"Retry-After":        []string{"1"},
		},
		Request: &http.Request{
			Method: http.MethodPost,
			URL:    reqURL,
		},
	}
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	lro, err := NewPoller[widget](firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	w, err := lro.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	if w.Size != 2 {
		t.Fatalf("unexpected widget size %d", w.Size)
	}
}

func TestOpPollerWithWidgetPOST(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted), mock.WithBody([]byte(`{"status": "InProgress"}`)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted), mock.WithBody([]byte(`{"status": "InProgress"}`)))
	// POST with no location header means the success response returns the model
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"status": "Succeeded", "size": 2}`)))
	defer close()

	reqURL, err := url.Parse(srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp := &http.Response{
		Body:       body,
		StatusCode: http.StatusAccepted,
		Header: http.Header{
			"Operation-Location": []string{srv.URL()},
			"Retry-After":        []string{"1"},
		},
		Request: &http.Request{
			Method: http.MethodPost,
			URL:    reqURL,
		},
	}
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	lro, err := NewPoller[widget](firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	w, err := lro.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	if w.Size != 2 {
		t.Fatalf("unexpected widget size %d", w.Size)
	}
}

func TestOpPollerWithWidgetResourceLocation(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted), mock.WithBody([]byte(`{"status": "InProgress"}`)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted), mock.WithBody([]byte(`{"status": "InProgress"}`)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(
		fmt.Sprintf(`{"status": "Succeeded", "resourceLocation": "%s"}`, srv.URL()))))
	// final GET will happen from the URL provided in the resourceLocation
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"size": 2}`)))
	defer close()

	reqURL, err := url.Parse(srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp := &http.Response{
		Body:       body,
		StatusCode: http.StatusAccepted,
		Header: http.Header{
			"Operation-Location": []string{srv.URL()},
			"Location":           []string{srv.URL()},
			"Retry-After":        []string{"1"},
		},
		Request: &http.Request{
			Method: http.MethodPatch,
			URL:    reqURL,
		},
	}
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	lro, err := NewPoller[widget](firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	w, err := lro.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	if w.Size != 2 {
		t.Fatalf("unexpected widget size %d", w.Size)
	}
}

func TestOpPollerWithResumeToken(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted), mock.WithBody([]byte(`{ "status": "InProgress"}`)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted), mock.WithBody([]byte(`{ "status": "InProgress"}`)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{ "status": "Succeeded"}`)))
	defer close()

	reqURL, err := url.Parse(srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp := &http.Response{
		Body:       body,
		StatusCode: http.StatusAccepted,
		Header: http.Header{
			"Operation-Location": []string{srv.URL()},
			"Retry-After":        []string{"1"},
		},
		Request: &http.Request{
			Method: http.MethodDelete,
			URL:    reqURL,
		},
	}
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	lro, err := NewPoller[none](firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	resp, err := lro.Poll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	if lro.Done() {
		t.Fatal("poller shouldn't be done yet")
	}
	_, err = lro.Result(context.Background())
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	tk, err := lro.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	lro, err = NewPollerFromResumeToken[none](tk, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = lro.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNopPoller(t *testing.T) {
	reqURL, err := url.Parse("https://fake.endpoint/for/testing")
	if err != nil {
		t.Fatal(err)
	}
	firstResp := &http.Response{
		StatusCode: http.StatusOK,
		Request: &http.Request{
			Method: http.MethodDelete,
			URL:    reqURL,
		},
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp.Body = body
	pl := newTestPipeline(nil)
	lro, err := NewPoller[none](firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := reflect.TypeOf(lro.op); pt != reflect.TypeOf((*pollers.NopPoller[none])(nil)) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	if !lro.Done() {
		t.Fatal("expected Done() for nopPoller")
	}
	resp, err := lro.Poll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp != firstResp {
		t.Fatal("unexpected response")
	}
	_, err = lro.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	tk, err := lro.ResumeToken()
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if tk != "" {
		t.Fatal("expected empty token")
	}
}

type preconstructedWidget struct {
	Size           int `json:"size"`
	Preconstructed int
}

func TestOpPollerWithResponseType(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted), mock.WithBody([]byte(`{"status": "InProgress"}`)), mock.WithHeader("Retry-After", "1"))
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted), mock.WithBody([]byte(`{"status": "InProgress"}`)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"status": "Succeeded"}`)))
	// PUT and PATCH state that a final GET will happen
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{"size": 2}`)))
	defer close()

	reqURL, err := url.Parse(srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp := &http.Response{
		Body:       body,
		StatusCode: http.StatusAccepted,
		Header: http.Header{
			"Operation-Location": []string{srv.URL()},
			"Retry-After":        []string{"1"},
		},
		Request: &http.Request{
			Method: http.MethodPut,
			URL:    reqURL,
		},
	}
	pl := newTestPipeline(&policy.ClientOptions{Transport: srv})
	lro, err := NewPoller(firstResp, pl, &NewPollerOptions[preconstructedWidget]{
		Response: &preconstructedWidget{
			Preconstructed: 12345,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	w, err := lro.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	if w.Size != 2 {
		t.Fatalf("unexpected widget size %d", w.Size)
	}
	if w.Preconstructed != 12345 {
		t.Fatalf("unexpected widget size %d", w.Preconstructed)
	}
}

const (
	provStateStarted   = `{ "properties": { "provisioningState": "Started" } }`
	provStateUpdating  = `{ "properties": { "provisioningState": "Updating" } }`
	provStateSucceeded = `{ "properties": { "provisioningState": "Succeeded" }, "field": "value" }`
	provStateFailed    = `{ "properties": { "provisioningState": "Failed" } }`
	statusInProgress   = `{ "status": "InProgress" }`
	statusSucceeded    = `{ "status": "Succeeded" }`
	statusCanceled     = `{ "status": "Canceled", "error": { "code": "OperationCanceled", "message": "somebody canceled it" } }`
	successResp        = `{ "field": "value" }`
)

type mockType struct {
	Field *string `json:"field,omitempty"`
}

func getPipeline(srv *mock.Server) Pipeline {
	return NewPipeline(
		"test",
		"v0.1.0",
		PipelineOptions{PerRetry: []policy.Policy{NewLogPolicy(nil)}},
		&policy.ClientOptions{
			Retry: policy.RetryOptions{
				MaxRetryDelay: 1 * time.Second,
			},
			Transport: srv,
		},
	)
}

func initialResponse(ctx context.Context, method, u string, resp io.Reader) (*http.Response, mock.TrackedClose) {
	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		panic(err)
	}
	body, closed := mock.NewTrackedCloser(resp)
	return &http.Response{
		Body:          body,
		ContentLength: -1,
		Header:        http.Header{},
		Request:       req,
	}, closed
}

func typeOfOpField[T any](pl *Poller[T]) reflect.Type {
	return reflect.ValueOf(pl).Elem().FieldByName("op").Elem().Type()
}

func TestNewPollerAsync(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(statusInProgress)))
	srv.AppendResponse(mock.WithBody([]byte(statusSucceeded)))
	srv.AppendResponse(mock.WithBody([]byte(successResp)))
	resp, closed := initialResponse(context.Background(), http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderAzureAsync, srv.URL())
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller[mockType](resp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := typeOfOpField(poller); pt != reflect.TypeOf((*async.Poller[mockType])(nil)) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	poller, err = NewPollerFromResumeToken[mockType](tk, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	result, err := poller.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	if v := *result.Field; v != "value" {
		t.Fatalf("unexpected value %s", v)
	}
	result, err = poller.Result(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if v := *result.Field; v != "value" {
		t.Fatalf("unexpected value %s", v)
	}
}

func TestNewPollerBody(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(provStateUpdating)), mock.WithHeader("Retry-After", "1"))
	srv.AppendResponse(mock.WithBody([]byte(provStateSucceeded)))
	resp, closed := initialResponse(context.Background(), http.MethodPatch, srv.URL(), strings.NewReader(provStateStarted))
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller[mockType](resp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := typeOfOpField(poller); pt != reflect.TypeOf((*body.Poller[mockType])(nil)) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	poller, err = NewPollerFromResumeToken[mockType](tk, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	result, err := poller.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	if v := *result.Field; v != "value" {
		t.Fatalf("unexpected value %s", v)
	}
}

func TestNewPollerInitialRetryAfter(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(statusInProgress)))
	srv.AppendResponse(mock.WithBody([]byte(statusSucceeded)))
	srv.AppendResponse(mock.WithBody([]byte(successResp)))
	resp, closed := initialResponse(context.Background(), http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderAzureAsync, srv.URL())
	resp.Header.Set("Retry-After", "1")
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller[mockType](resp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := typeOfOpField(poller); pt != reflect.TypeOf((*async.Poller[mockType])(nil)) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	result, err := poller.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	if v := *result.Field; v != "value" {
		t.Fatalf("unexpected value %s", v)
	}
}

func TestNewPollerCanceled(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(statusInProgress)))
	srv.AppendResponse(mock.WithBody([]byte(statusCanceled)), mock.WithStatusCode(http.StatusOK))
	resp, closed := initialResponse(context.Background(), http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderAzureAsync, srv.URL())
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller[mockType](resp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := typeOfOpField(poller); pt != reflect.TypeOf((*async.Poller[mockType])(nil)) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	_, err = poller.Poll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if poller.Done() {
		t.Fatal("poller shouldn't be done yet")
	}
	_, err = poller.Poll(context.Background())
	if err != nil {
		t.Fatal("expected nil error")
	}
	if !poller.Done() {
		t.Fatal("poller should be done")
	}
	_, err = poller.Result(context.Background())
	if err == nil {
		t.Fatal("unexpected nil error")
	}
}

func TestNewPollerFailed(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(provStateFailed)))
	resp, closed := initialResponse(context.Background(), http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderAzureAsync, srv.URL())
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller[mockType](resp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := typeOfOpField(poller); pt != reflect.TypeOf((*async.Poller[mockType])(nil)) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	_, err = poller.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err == nil {
		t.Fatal(err)
	}
}

func TestNewPollerFailedWithError(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(statusInProgress)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
	resp, closed := initialResponse(context.Background(), http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderAzureAsync, srv.URL())
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller[mockType](resp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := typeOfOpField(poller); pt != reflect.TypeOf((*async.Poller[mockType])(nil)) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	_, err = poller.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err == nil {
		t.Fatal(err)
	}
}

func TestNewPollerSuccessNoContent(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(provStateUpdating)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusNoContent))
	resp, closed := initialResponse(context.Background(), http.MethodPatch, srv.URL(), strings.NewReader(provStateStarted))
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller[mockType](resp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := typeOfOpField(poller); pt != reflect.TypeOf((*body.Poller[mockType])(nil)) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	poller, err = NewPollerFromResumeToken[mockType](tk, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	result, err := poller.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	if result.Field != nil {
		t.Fatal("expected nil result")
	}
}

func TestNewPollerFail202NoHeaders(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	resp, closed := initialResponse(context.Background(), http.MethodDelete, srv.URL(), http.NoBody)
	resp.StatusCode = http.StatusAccepted
	pl := getPipeline(srv)
	poller, err := NewPoller[mockType](resp, pl, nil)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if poller != nil {
		t.Fatal("expected nil poller")
	}
}

type preconstructedMockType struct {
	Field          *string `json:"field,omitempty"`
	Preconstructed int
}

func TestNewPollerWithResponseType(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(provStateUpdating)), mock.WithHeader("Retry-After", "1"))
	srv.AppendResponse(mock.WithBody([]byte(provStateSucceeded)))
	resp, closed := initialResponse(context.Background(), http.MethodPatch, srv.URL(), strings.NewReader(provStateStarted))
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller[preconstructedMockType](resp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := typeOfOpField(poller); pt != reflect.TypeOf((*body.Poller[preconstructedMockType])(nil)) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	poller, err = NewPollerFromResumeToken(tk, pl, &NewPollerFromResumeTokenOptions[preconstructedMockType]{
		Response: &preconstructedMockType{
			Preconstructed: 12345,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err := poller.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	if v := *result.Field; v != "value" {
		t.Fatalf("unexpected value %s", v)
	}
	if result.Preconstructed != 12345 {
		t.Fatalf("unexpected value %d", result.Preconstructed)
	}
}

// purposefully looks like an async poller but isn't
type customHandler struct {
	PollURL string `json:"asyncURL"`
	State   string `json:"state"`
	p       Pipeline
}

func (c *customHandler) Done() bool {
	return c.State == "Succeeded"
}

func (c *customHandler) Poll(ctx context.Context) (*http.Response, error) {
	req, err := NewRequest(ctx, http.MethodGet, c.PollURL)
	if err != nil {
		return nil, err
	}
	resp, err := c.p.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := Payload(resp)
	if err != nil {
		return nil, err
	}
	type statusMon struct {
		Status string `json:"status"`
	}
	var sm statusMon
	if err = json.Unmarshal(body, &sm); err != nil {
		return nil, err
	}
	c.State = sm.Status
	return resp, nil
}

func (c *customHandler) Result(ctx context.Context, out *mockType) error {
	req, err := NewRequest(ctx, http.MethodGet, c.PollURL)
	if err != nil {
		return err
	}
	resp, err := c.p.Do(req)
	if err != nil {
		return err
	}
	body, err := Payload(resp)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, out); err != nil {
		return err
	}
	return nil
}

func TestNewPollerWithCustomHandler(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(statusInProgress)))
	srv.AppendResponse(mock.WithBody([]byte(statusSucceeded)))
	srv.AppendResponse(mock.WithBody([]byte(successResp)))
	resp, closed := initialResponse(context.Background(), http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderAzureAsync, srv.URL())
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller(resp, pl, &NewPollerOptions[mockType]{
		Handler: &customHandler{
			PollURL: srv.URL(),
			State:   "InProgress",
			p:       pl,
		},
	})
	require.NoError(t, err)
	require.False(t, closed())
	require.IsType(t, &customHandler{}, poller.op)
	tk, err := poller.ResumeToken()
	require.NoError(t, err)
	poller, err = NewPollerFromResumeToken(tk, pl, &NewPollerFromResumeTokenOptions[mockType]{
		Handler: &customHandler{
			p: pl,
		},
	})
	require.IsType(t, &customHandler{}, poller.op)
	require.NoError(t, err)
	result, err := poller.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	require.NoError(t, err)
	require.EqualValues(t, "value", *result.Field)
	result, err = poller.Result(context.Background())
	require.NoError(t, err)
	require.EqualValues(t, "value", *result.Field)
}

func TestShortenPollerTypeName(t *testing.T) {
	result := shortenTypeName("Poller[module/package.ClientOperationResponse].PollUntilDone")
	require.EqualValues(t, "Poller[ClientOperationResponse].PollUntilDone", result)

	result = shortenTypeName("Poller[package.ClientOperationResponse].PollUntilDone")
	require.EqualValues(t, "Poller[ClientOperationResponse].PollUntilDone", result)

	result = shortenTypeName("Poller[ClientOperationResponse].PollUntilDone")
	require.EqualValues(t, "Poller[ClientOperationResponse].PollUntilDone", result)

	result = shortenTypeName("Poller.PollUntilDone")
	require.EqualValues(t, "Poller.PollUntilDone", result)
}

func TestNewFakePoller(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithHeader(shared.HeaderFakePollerStatus, "FakePollerInProgress"))
	srv.AppendResponse(mock.WithHeader(shared.HeaderFakePollerStatus, poller.StatusSucceeded), mock.WithStatusCode(http.StatusNoContent))
	pollCtx := context.WithValue(context.Background(), shared.CtxAPINameKey{}, "FakeAPI")
	resp, closed := initialResponse(pollCtx, http.MethodPatch, srv.URL(), http.NoBody)
	resp.StatusCode = http.StatusCreated
	resp.Header.Set(shared.HeaderFakePollerStatus, "FakePollerInProgress")
	pl := getPipeline(srv)
	poller, err := NewPoller[mockType](resp, pl, nil)
	require.NoError(t, err)
	require.True(t, closed())
	if pt := typeOfOpField(poller); pt != reflect.TypeOf((*fake.Poller[mockType])(nil)) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	tk, err := poller.ResumeToken()
	require.NoError(t, err)
	poller, err = NewPollerFromResumeToken[mockType](tk, pl, nil)
	require.NoError(t, err)
	result, err := poller.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	require.NoError(t, err)
	require.Nil(t, result.Field)
}

func TestNewPollerWithThrottling(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(statusInProgress)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests))
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests))
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests))
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests))
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests))
	srv.AppendResponse(mock.WithBody([]byte(statusSucceeded)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests))
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests))
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests))
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests))
	srv.AppendResponse(mock.WithStatusCode(http.StatusTooManyRequests))
	srv.AppendResponse(mock.WithBody([]byte(successResp)))
	resp, closed := initialResponse(context.Background(), http.MethodPatch, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderAzureAsync, srv.URL())
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller[mockType](resp, pl, nil)
	require.NoError(t, err)
	require.True(t, closed())
	if pt := typeOfOpField(poller); pt != reflect.TypeOf((*async.Poller[mockType])(nil)) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	tk, err := poller.ResumeToken()
	require.NoError(t, err)
	poller, err = NewPollerFromResumeToken[mockType](tk, pl, nil)
	require.NoError(t, err)
	result, err := poller.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	require.Zero(t, result)
	var respErr *exported.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.EqualValues(t, http.StatusTooManyRequests, respErr.StatusCode)
	result, err = poller.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	require.Zero(t, result)
	require.ErrorAs(t, err, &respErr)
	require.EqualValues(t, http.StatusTooManyRequests, respErr.StatusCode)
	result, err = poller.PollUntilDone(context.Background(), &PollUntilDoneOptions{Frequency: time.Millisecond})
	require.NoError(t, err)
	require.NotNil(t, result.Field)
	require.EqualValues(t, "value", *result.Field)
}
