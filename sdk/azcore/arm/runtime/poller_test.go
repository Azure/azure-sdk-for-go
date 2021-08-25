//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"io"
	"io/ioutil"
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

type mockError struct {
	Msg string `json:"error"`
}

func (m mockError) Error() string {
	return m.Msg
}

func getPipeline(srv *mock.Server) pipeline.Pipeline {
	return runtime.NewPipeline(
		srv,
		runtime.NewLogPolicy(nil))
}

func handleError(resp *http.Response) error {
	var me mockError
	if err := runtime.UnmarshalAsJSON(resp, &me); err != nil {
		return err
	}
	return me
}

func initialResponse(method, u string, resp io.Reader) *http.Response {
	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		panic(err)
	}
	return &http.Response{
		Body:          ioutil.NopCloser(resp),
		ContentLength: -1,
		Header:        http.Header{},
		Request:       req,
	}
}

func TestNewPollerAsync(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(statusInProgress)))
	srv.AppendResponse(mock.WithBody([]byte(statusSucceeded)))
	srv.AppendResponse(mock.WithBody([]byte(successResp)))
	resp := initialResponse(http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderAzureAsync, srv.URL())
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller("pollerID", "", resp, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	if pt := pollers.PollerType(poller); pt != reflect.TypeOf(&async.Poller{}) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	poller, err = NewPollerFromResumeToken("pollerID", tk, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	var result mockType
	_, err = poller.PollUntilDone(context.Background(), 10*time.Millisecond, &result)
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
	resp := initialResponse(http.MethodPatch, srv.URL(), strings.NewReader(provStateStarted))
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller("pollerID", "", resp, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	if pt := pollers.PollerType(poller); pt != reflect.TypeOf(&body.Poller{}) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	poller, err = NewPollerFromResumeToken("pollerID", tk, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	var result mockType
	_, err = poller.PollUntilDone(context.Background(), 10*time.Millisecond, &result)
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
	resp := initialResponse(http.MethodPatch, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderLocation, srv.URL())
	resp.StatusCode = http.StatusAccepted
	pl := getPipeline(srv)
	poller, err := NewPoller("pollerID", "", resp, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	if pt := pollers.PollerType(poller); pt != reflect.TypeOf(&loc.Poller{}) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	poller, err = NewPollerFromResumeToken("pollerID", tk, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	var result mockType
	_, err = poller.PollUntilDone(context.Background(), 10*time.Millisecond, &result)
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
	resp := initialResponse(http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderAzureAsync, srv.URL())
	resp.Header.Set("Retry-After", "1")
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller("pollerID", "", resp, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	if pt := pollers.PollerType(poller); pt != reflect.TypeOf(&async.Poller{}) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	var result mockType
	_, err = poller.PollUntilDone(context.Background(), 10*time.Millisecond, &result)
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
	resp := initialResponse(http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderAzureAsync, srv.URL())
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller("pollerID", "", resp, pl, handleError)
	if err != nil {
		t.Fatal(err)
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
	resp := initialResponse(http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(shared.HeaderAzureAsync, srv.URL())
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller("pollerID", "", resp, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	if pt := pollers.PollerType(poller); pt != reflect.TypeOf(&async.Poller{}) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	var result mockType
	_, err = poller.PollUntilDone(context.Background(), 10*time.Millisecond, &result)
	if err == nil {
		t.Fatal(err)
	}
	if _, ok := err.(mockError); !ok {
		t.Fatalf("unexpected error type %T", err)
	}
}

func TestNewPollerSuccessNoContent(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(provStateUpdating)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusNoContent))
	resp := initialResponse(http.MethodPatch, srv.URL(), strings.NewReader(provStateStarted))
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewPoller("pollerID", "", resp, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	if pt := pollers.PollerType(poller); pt != reflect.TypeOf(&body.Poller{}) {
		t.Fatalf("unexpected poller type %s", pt.String())
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	poller, err = NewPollerFromResumeToken("pollerID", tk, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	var result mockType
	_, err = poller.PollUntilDone(context.Background(), 10*time.Millisecond, &result)
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
	resp := initialResponse(http.MethodDelete, srv.URL(), http.NoBody)
	resp.StatusCode = http.StatusAccepted
	pl := getPipeline(srv)
	poller, err := NewPoller("pollerID", "", resp, pl, handleError)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if poller != nil {
		t.Fatal("expected nil poller")
	}
}
