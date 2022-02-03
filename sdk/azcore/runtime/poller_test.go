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
	p, err := NewPoller("fake.poller", &http.Response{
		Body:       body,
		StatusCode: http.StatusBadRequest,
	}, newTestPipeline(nil))
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
			p, err := NewPollerFromResumeToken("fake.poller", test.token, newTestPipeline(nil))
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
	lro, err := NewPoller("fake.poller", firstResp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	var respFromCtx *http.Response
	ctxWithResp := WithCaptureResponse(context.Background(), &respFromCtx)
	resp, err := lro.PollUntilDone(ctxWithResp, time.Second, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	if respFromCtx != resp {
		t.Fatal("response from context doesn't match returned response")
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
	lro, err := NewPoller("fake.poller", firstResp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	var w widget
	resp, err := lro.PollUntilDone(context.Background(), time.Second, &w)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
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
	lro, err := NewPoller("fake.poller", firstResp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	var w widget
	resp, err := lro.PollUntilDone(context.Background(), time.Second, &w)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if _, ok := err.(*shared.ResponseError); !ok {
		t.Fatal("expected pollerError")
	}
	if resp != nil {
		t.Fatal("expected nil response")
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
	lro, err := NewPoller("fake.poller", firstResp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	var w widget
	resp, err := lro.PollUntilDone(context.Background(), time.Second, &w)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if e := err.Error(); e != "oops" {
		t.Fatalf("expected error %s", e)
	}
	if resp != nil {
		t.Fatal("expected nil response")
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
	lro, err := NewPoller("fake.poller", firstResp, pl)
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
	resp, err = lro.FinalResponse(context.Background(), nil)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	tk, err := lro.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	lro, err = NewPollerFromResumeToken("fake.poller", tk, pl)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = lro.PollUntilDone(context.Background(), time.Second, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
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
	lro, err := NewPoller("fake.poller", firstResp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	resp, err := lro.PollUntilDone(ctx, time.Second, nil)
	cancel()
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("expected nil response")
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
	lro, err := NewPoller("fake.poller", firstResp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	resp, err := lro.PollUntilDone(context.Background(), time.Second, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
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
	lro, err := NewPoller("fake.poller", firstResp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	var w widget
	resp, err := lro.PollUntilDone(context.Background(), time.Second, &w)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
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
	lro, err := NewPoller("fake.poller", firstResp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	var w widget
	resp, err := lro.PollUntilDone(context.Background(), time.Second, &w)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
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
	lro, err := NewPoller("fake.poller", firstResp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	var w widget
	resp, err := lro.PollUntilDone(context.Background(), time.Second, &w)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
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
	lro, err := NewPoller("fake.poller", firstResp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	var w widget
	resp, err := lro.PollUntilDone(context.Background(), time.Second, &w)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
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
	lro, err := NewPoller("fake.poller", firstResp, pl)
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
	resp, err = lro.FinalResponse(context.Background(), nil)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	tk, err := lro.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	lro, err = NewPollerFromResumeToken("fake.poller", tk, pl)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = lro.PollUntilDone(context.Background(), time.Second, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}

func TestNopPoller(t *testing.T) {
	firstResp := &http.Response{
		StatusCode: http.StatusOK,
	}
	body, closed := mock.NewTrackedCloser(http.NoBody)
	firstResp.Body = body
	pl := newTestPipeline(nil)
	lro, err := NewPoller("fake.poller", firstResp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := pollers.PollerType(lro); pt != reflect.TypeOf(&pollers.NopPoller{}) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	if !lro.Done() {
		t.Fatal("expected Done() for nopPoller")
	}
	resp, err := lro.FinalResponse(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp != firstResp {
		t.Fatal("unexpected response")
	}
	resp, err = lro.Poll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp != firstResp {
		t.Fatal("unexpected response")
	}
	resp, err = lro.PollUntilDone(context.Background(), time.Second, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp != firstResp {
		t.Fatal("unexpected response")
	}
	tk, err := lro.ResumeToken()
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if tk != "" {
		t.Fatal("expected empty token")
	}
}
