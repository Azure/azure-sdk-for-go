// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armcore

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

const (
	mockSuccessResp = `{"field": "success"}`
)

type mockType struct {
	Field *string `json:"field,omitempty"`
}

func getPipeline(srv *mock.Server) azcore.Pipeline {
	return azcore.NewPipeline(
		srv,
		azcore.NewTelemetryPolicy(azcore.TelemetryOptions{}),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(nil),
		azcore.NewRequestLogPolicy(azcore.RequestLogOptions{}))
}

func handleError(resp *azcore.Response) error {
	return fmt.Errorf("error status: %d", resp.StatusCode)
}

func TestNewPollerTracker(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(mockSuccessResp)))
	p := getPipeline(srv)
	req := azcore.NewRequest(http.MethodPost, srv.URL())
	resp, err := p.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	poller, err := NewPoller("testPoller", "", resp, handleError)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	pt := poller.(*pollingTrackerPost)
	if pt.PollerType != "testPoller" {
		t.Fatal("wrong poller type assigned")
	}
	if pt.resp != resp {
		t.Fatal("wrong response assigned")
	}
	if pt.Method != "POST" {
		t.Fatal("wrong poller method")
	}
}

func TestResumeTokenFail(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(`{"field": "success"}`)))
	p := getPipeline(srv)
	req := azcore.NewRequest(http.MethodPost, srv.URL())
	resp, err := p.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	poller, err := NewPoller("testPoller", "", resp, handleError)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	tk, err := poller.ResumeToken()
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
	if tk != "" {
		t.Fatal("did not expect to receive resume token for a poller in a terminal state")
	}
}

func TestPollUntilDone(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithBody([]byte(mockSuccessResp)))
	p := getPipeline(srv)
	req := azcore.NewRequest(http.MethodPut, srv.URL())
	resp, err := p.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	poller, err := NewPoller("testPoller", "", resp, handleError)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	m := &mockType{}
	pollerResp, err := poller.PollUntilDone(context.Background(), 1*time.Millisecond, p, m)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if pollerResp == nil {
		t.Fatal("Unexpected nil response")
	}
	if pollerResp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected response status code")
	}
	if *m.Field != "success" {
		t.Fatalf("Unexpected value for MockType.Field: %s", *m.Field)
	}
}

func TestPutFinalResponseCheck(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithBody([]byte(`{"other": "other"}`)))
	srv.AppendResponse(mock.WithBody([]byte(mockSuccessResp)))
	p := getPipeline(srv)
	req := azcore.NewRequest(http.MethodPut, srv.URL())
	resp, err := p.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	poller, err := NewPoller("testPoller", "", resp, handleError)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	m := &mockType{}
	pollerResp, err := poller.PollUntilDone(context.Background(), 1*time.Millisecond, p, m)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if pollerResp == nil {
		t.Fatal("Unexpected nil response")
	}
	if pollerResp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected response status code")
	}
	if *m.Field != "success" {
		t.Fatalf("Unexpected value for MockType.Field: %s", *m.Field)
	}
}

func TestNewPollerFromResumeTokenTracker(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithBody([]byte(`{"field": "success"}`)))
	p := getPipeline(srv)
	req := azcore.NewRequest(http.MethodPut, srv.URL())
	resp, err := p.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	poller, err := NewPoller("testPoller", "", resp, handleError)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	tk, err := poller.ResumeToken()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	poller, err = NewPollerFromResumeToken("testPoller", tk, handleError)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	m := &mockType{}
	pollerResp, err := poller.PollUntilDone(context.Background(), 1*time.Millisecond, p, m)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if pollerResp == nil {
		t.Fatal("Unexpected nil response")
	}
	if pollerResp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected response status code")
	}
	if *m.Field != "success" {
		t.Fatalf("Unexpected value for MockType.Field: %s", *m.Field)
	}
}
