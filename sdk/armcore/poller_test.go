// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armcore

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/armcore/internal/pollers"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore/internal/pollers/async"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore/internal/pollers/body"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore/internal/pollers/loc"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
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

func getPipeline(srv *mock.Server) azcore.Pipeline {
	return azcore.NewPipeline(
		srv,
		azcore.NewLogPolicy(nil))
}

func handleError(resp *azcore.Response) error {
	var me mockError
	if err := resp.UnmarshalAsJSON(&me); err != nil {
		return err
	}
	return me
}

func initialResponse(method, u string, resp io.Reader) *azcore.Response {
	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		panic(err)
	}
	return &azcore.Response{
		Response: &http.Response{
			Body:          ioutil.NopCloser(resp),
			ContentLength: -1,
			Header:        http.Header{},
			Request:       req,
		},
	}
}

func TestNewLROPollerAsync(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(statusInProgress)))
	srv.AppendResponse(mock.WithBody([]byte(statusSucceeded)))
	srv.AppendResponse(mock.WithBody([]byte(successResp)))
	resp := initialResponse(http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(pollers.HeaderAzureAsync, srv.URL())
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewLROPoller("pollerID", "", resp, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := poller.lro.(*async.Poller); !ok {
		t.Fatalf("unexpected poller type %T", poller.lro)
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	poller, err = NewLROPollerFromResumeToken("pollerID", tk, pl, handleError)
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

func TestNewLROPollerBody(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(provStateUpdating)), mock.WithHeader("Retry-After", "1"))
	srv.AppendResponse(mock.WithBody([]byte(provStateSucceeded)))
	resp := initialResponse(http.MethodPatch, srv.URL(), strings.NewReader(provStateStarted))
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewLROPoller("pollerID", "", resp, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := poller.lro.(*body.Poller); !ok {
		t.Fatalf("unexpected poller type %T", poller.lro)
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	poller, err = NewLROPollerFromResumeToken("pollerID", tk, pl, handleError)
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

func TestNewLROPollerLoc(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithBody([]byte(successResp)))
	resp := initialResponse(http.MethodPatch, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(pollers.HeaderLocation, srv.URL())
	resp.StatusCode = http.StatusAccepted
	pl := getPipeline(srv)
	poller, err := NewLROPoller("pollerID", "", resp, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := poller.lro.(*loc.Poller); !ok {
		t.Fatalf("unexpected poller type %T", poller.lro)
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	poller, err = NewLROPollerFromResumeToken("pollerID", tk, pl, handleError)
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

func TestNewLROPollerNop(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	resp := initialResponse(http.MethodPost, srv.URL(), strings.NewReader(successResp))
	resp.StatusCode = http.StatusOK
	poller, err := NewLROPoller("pollerID", "", resp, getPipeline(srv), handleError)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := poller.lro.(*nopPoller); !ok {
		t.Fatalf("unexpected poller type %T", poller.lro)
	}
	tk, err := poller.ResumeToken()
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if tk != "" {
		t.Fatal("expected empty token")
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

func TestNewLROPollerInitialRetryAfter(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(statusInProgress)))
	srv.AppendResponse(mock.WithBody([]byte(statusSucceeded)))
	srv.AppendResponse(mock.WithBody([]byte(successResp)))
	resp := initialResponse(http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(pollers.HeaderAzureAsync, srv.URL())
	resp.Header.Set("Retry-After", "1")
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewLROPoller("pollerID", "", resp, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := poller.lro.(*async.Poller); !ok {
		t.Fatalf("unexpected poller type %T", poller.lro)
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

func TestNewLROPollerCanceled(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(statusInProgress)))
	srv.AppendResponse(mock.WithBody([]byte(statusCanceled)), mock.WithStatusCode(http.StatusOK))
	resp := initialResponse(http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(pollers.HeaderAzureAsync, srv.URL())
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewLROPoller("pollerID", "", resp, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := poller.lro.(*async.Poller); !ok {
		t.Fatalf("unexpected poller type %T", poller.lro)
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

func TestNewLROPollerFailedWithError(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(statusInProgress)))
	srv.AppendResponse(mock.WithBody([]byte(errorResp)), mock.WithStatusCode(http.StatusBadRequest))
	resp := initialResponse(http.MethodPut, srv.URL(), strings.NewReader(provStateStarted))
	resp.Header.Set(pollers.HeaderAzureAsync, srv.URL())
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewLROPoller("pollerID", "", resp, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := poller.lro.(*async.Poller); !ok {
		t.Fatalf("unexpected poller type %T", poller.lro)
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

func TestNewLROPollerSuccessNoContent(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(provStateUpdating)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusNoContent))
	resp := initialResponse(http.MethodPatch, srv.URL(), strings.NewReader(provStateStarted))
	resp.StatusCode = http.StatusCreated
	pl := getPipeline(srv)
	poller, err := NewLROPoller("pollerID", "", resp, pl, handleError)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := poller.lro.(*body.Poller); !ok {
		t.Fatalf("unexpected poller type %T", poller.lro)
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	poller, err = NewLROPollerFromResumeToken("pollerID", tk, pl, handleError)
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

func TestNewLROPollerFail202NoHeaders(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	resp := initialResponse(http.MethodDelete, srv.URL(), http.NoBody)
	resp.StatusCode = http.StatusAccepted
	pl := getPipeline(srv)
	poller, err := NewLROPoller("pollerID", "", resp, pl, handleError)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if poller != nil {
		t.Fatal("expected nil poller")
	}
}
