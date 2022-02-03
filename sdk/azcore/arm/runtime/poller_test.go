//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/pollers/async"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/pollers/body"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/pollers/loc"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pipeline"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

const (
	provStateStarted   = `{ "properties": { "provisioningState": "Started" } }`
	provStateUpdating  = `{ "properties": { "provisioningState": "Updating" } }`
	provStateSucceeded = `{ "properties": { "provisioningState": "Succeeded" }, "field": "value" }`
	provStateFailed    = `{ "properties": { "provisioningState": "Failed" } }` //nolint
	statusInProgress   = `{ "status": "InProgress" }`
	statusSucceeded    = `{ "status": "Succeeded" }`
	statusCanceled     = `{ "status": "Canceled" }`
	successResp        = `{ "field": "value" }`
	errorResp          = `{ "error": "the operation failed" }`
)

type mockType struct {
	Field *string `json:"field,omitempty"`
}

func getPipeline(srv *mock.Server) pipeline.Pipeline {
	return runtime.NewPipeline(
		"test",
		"v0.1.0",
		runtime.PipelineOptions{PerRetry: []pipeline.Policy{runtime.NewLogPolicy(nil)}},
		&policy.ClientOptions{Transport: srv},
	)
}

func initialResponse(method, u string, resp io.Reader) (*http.Response, mock.TrackedClose) {
	req, err := http.NewRequest(method, u, nil)
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

func TestNewPollerAsync(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(statusInProgress)))
	srv.AppendResponse(mock.WithBody([]byte(statusSucceeded)))
	srv.AppendResponse(mock.WithBody([]byte(successResp)))
	resp, closed := initialResponse(http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderAzureAsync, srv.URL())
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller("pollerID", "", resp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := pollers.PollerType(poller); pt != reflect.TypeOf(&async.Poller{}) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	poller, err = NewPollerFromResumeToken("pollerID", tk, pl)
	if err != nil {
		t.Fatal(err)
	}
	var result mockType
	_, err = poller.PollUntilDone(context.Background(), time.Second, &result)
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
	resp, closed := initialResponse(http.MethodPatch, srv.URL(), strings.NewReader(provStateStarted))
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller("pollerID", "", resp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := pollers.PollerType(poller); pt != reflect.TypeOf(&body.Poller{}) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	poller, err = NewPollerFromResumeToken("pollerID", tk, pl)
	if err != nil {
		t.Fatal(err)
	}
	var result mockType
	_, err = poller.PollUntilDone(context.Background(), time.Second, &result)
	if err != nil {
		t.Fatal(err)
	}
	if v := *result.Field; v != "value" {
		t.Fatalf("unexpected value %s", v)
	}
}

func TestNewPollerLoc(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithBody([]byte(successResp)))
	resp, closed := initialResponse(http.MethodPatch, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderLocation, srv.URL())
	resp.StatusCode = http.StatusAccepted
	pl := getPipeline(srv)
	poller, err := NewPoller("pollerID", "", resp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := pollers.PollerType(poller); pt != reflect.TypeOf(&loc.Poller{}) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	poller, err = NewPollerFromResumeToken("pollerID", tk, pl)
	if err != nil {
		t.Fatal(err)
	}
	var result mockType
	_, err = poller.PollUntilDone(context.Background(), time.Second, &result)
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
	resp, closed := initialResponse(http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderAzureAsync, srv.URL())
	resp.Header.Set("Retry-After", "1")
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller("pollerID", "", resp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := pollers.PollerType(poller); pt != reflect.TypeOf(&async.Poller{}) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	var result mockType
	_, err = poller.PollUntilDone(context.Background(), time.Second, &result)
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
	resp, closed := initialResponse(http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderAzureAsync, srv.URL())
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller("pollerID", "", resp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := pollers.PollerType(poller); pt != reflect.TypeOf(&async.Poller{}) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	_, err = poller.Poll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	_, err = poller.Poll(context.Background())
	if err == nil {
		t.Fatal("unexpected nil error")
	}
}

func TestNewPollerFailedWithError(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(statusInProgress)))
	srv.AppendResponse(mock.WithBody([]byte(errorResp)), mock.WithStatusCode(http.StatusBadRequest))
	resp, closed := initialResponse(http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderAzureAsync, srv.URL())
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller("pollerID", "", resp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := pollers.PollerType(poller); pt != reflect.TypeOf(&async.Poller{}) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	var result mockType
	_, err = poller.PollUntilDone(context.Background(), time.Second, &result)
	if err == nil {
		t.Fatal(err)
	}
	if _, ok := err.(*shared.ResponseError); !ok {
		t.Fatalf("unexpected error type %T", err)
	}
}

func TestNewPollerSuccessNoContent(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(provStateUpdating)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusNoContent))
	resp, closed := initialResponse(http.MethodPatch, srv.URL(), strings.NewReader(provStateStarted))
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller("pollerID", "", resp, pl)
	if err != nil {
		t.Fatal(err)
	}
	if !closed() {
		t.Fatal("initial response body wasn't closed")
	}
	if pt := pollers.PollerType(poller); pt != reflect.TypeOf(&body.Poller{}) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	poller, err = NewPollerFromResumeToken("pollerID", tk, pl)
	if err != nil {
		t.Fatal(err)
	}
	var result mockType
	_, err = poller.PollUntilDone(context.Background(), time.Second, &result)
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
	resp, closed := initialResponse(http.MethodDelete, srv.URL(), http.NoBody)
	resp.StatusCode = http.StatusAccepted
	pl := getPipeline(srv)
	poller, err := NewPoller("pollerID", "", resp, pl)
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
