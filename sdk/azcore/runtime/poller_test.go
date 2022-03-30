//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

type widget struct {
	Size int `json:"size"`
}

func TestNewPollerFail(t *testing.T) {
	body, closed := mock.NewTrackedCloser(http.NoBody)
	p, err := NewPoller[widget]("fake.poller", &http.Response{
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
			p, err := NewPollerFromResumeToken[widget]("fake.poller", test.token, newTestPipeline(nil), nil)
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
	lro, err := NewPoller[struct{}]("fake.poller", firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	var respFromCtx *http.Response
	ctxWithResp := WithCaptureResponse(context.Background(), &respFromCtx)
	_, err = lro.PollUntilDone(ctxWithResp, time.Second)
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
	lro, err := NewPoller[widget]("fake.poller", firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	w, err := lro.PollUntilDone(context.Background(), time.Second)
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
	lro, err := NewPoller[widget]("fake.poller", firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	w, err := lro.PollUntilDone(context.Background(), time.Second)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if _, ok := err.(*shared.ResponseError); !ok {
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
	lro, err := NewPoller[widget]("fake.poller", firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	w, err := lro.PollUntilDone(context.Background(), time.Second)
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
	lro, err := NewPoller[struct{}]("fake.poller", firstResp, pl, nil)
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
	lro, err = NewPollerFromResumeToken[struct{}]("fake.poller", tk, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = lro.PollUntilDone(context.Background(), time.Second)
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
	lro, err := NewPoller[struct{}]("fake.poller", firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	_, err = lro.PollUntilDone(ctx, time.Second)
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
	lro, err := NewPoller[struct{}]("fake.poller", firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	_, err = lro.PollUntilDone(context.Background(), time.Second)
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
	lro, err := NewPoller[widget]("fake.poller", firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	w, err := lro.PollUntilDone(context.Background(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if w.Size != 2 {
		t.Fatalf("unexpected widget size %d", w.Size)
	}
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
	lro, err := NewPoller[widget]("fake.poller", firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	w, err := lro.PollUntilDone(context.Background(), time.Second)
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
	lro, err := NewPoller[widget]("fake.poller", firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	w, err := lro.PollUntilDone(context.Background(), time.Second)
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
	lro, err := NewPoller[widget]("fake.poller", firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	w, err := lro.PollUntilDone(context.Background(), time.Second)
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
	lro, err := NewPoller[struct{}]("fake.poller", firstResp, pl, nil)
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
	lro, err = NewPollerFromResumeToken[struct{}]("fake.poller", tk, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = lro.PollUntilDone(context.Background(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNopPoller(t *testing.T) {
	firstResp := &http.Response{
		StatusCode: http.StatusOK,
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp.Body = body
	pl := newTestPipeline(nil)
	lro, err := NewPoller[struct{}]("fake.poller", firstResp, pl, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := pollers.PollerType(lro.pt); pt != reflect.TypeOf(&pollers.NopPoller{}) {
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
	_, err = lro.PollUntilDone(context.Background(), time.Second)
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
