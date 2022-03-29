//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package pollers

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/pipeline"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestKindFromToken(t *testing.T) {
	const tk = `{ "type": "pollerID;kind" }`
	k, err := KindFromToken("pollerID", tk)
	if err != nil {
		t.Fatal(err)
	}
	if k != "kind" {
		t.Fatalf("unexpected kind %s", k)
	}
	k, err = KindFromToken("mismatched", tk)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if k != "" {
		t.Fatal("expected empty kind")
	}
}

func TestKindFromTokenInvalid(t *testing.T) {
	const tk1 = `{ "missing": "type" }`
	k, err := KindFromToken("mismatched", tk1)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if k != "" {
		t.Fatal("expected empty kind")
	}
	const tk2 = `{ "type": false }`
	k, err = KindFromToken("mismatched", tk2)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if k != "" {
		t.Fatal("expected empty kind")
	}
	const tk3 = `{ "type": "pollerID;kind;extra" }`
	k, err = KindFromToken("mismatched", tk3)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if k != "" {
		t.Fatal("expected empty kind")
	}
}

// simple status code-based poller
type fakePoller struct {
	Ep   string
	Fg   string
	Code int
}

func (f *fakePoller) Done() bool {
	return f.Code == http.StatusOK || f.Code == http.StatusNoContent
}

func (f *fakePoller) Update(resp *http.Response) error {
	f.Code = resp.StatusCode
	return nil
}

func (f *fakePoller) FinalGetURL() string {
	return f.Fg
}

func (f *fakePoller) URL() string {
	return f.Ep
}

func (f *fakePoller) Status() string {
	switch f.Code {
	case http.StatusAccepted:
		return StatusInProgress
	case http.StatusOK, http.StatusNoContent:
		return StatusSucceeded
	case http.StatusCreated:
		return StatusCanceled
	default:
		return StatusFailed
	}
}

func TestNewPoller(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithStatusCode(http.StatusNoContent)) // terminal
	defer close()
	pl := pipeline.NewPipeline(srv)
	firstResp := &http.Response{
		StatusCode: http.StatusAccepted,
		Header:     http.Header{},
	}
	firstResp.Header.Set(shared.HeaderRetryAfter, "1")
	p := NewPoller(&fakePoller{Ep: srv.URL()}, firstResp, pl)
	if p.Done() {
		t.Fatal("unexpected done")
	}
	resp, err := p.FinalResponse(context.Background(), nil)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	tk, err := p.ResumeToken()
	if err != nil {
		t.Fatal(err)
	}
	if tk == "" {
		t.Fatal("unexpected empty resume token")
	}
	resp, err = p.PollUntilDone(context.Background(), 1*time.Millisecond, nil)
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	resp, err = p.PollUntilDone(context.Background(), time.Second, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	tk, err = p.ResumeToken()
	if err == nil {
		t.Fatal("unexpected nil error")
	}
	if tk != "" {
		t.Fatal("expected empty resume token")
	}
}

func TestNewPollerWithFinalGET(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted), mock.WithHeader(shared.HeaderRetryAfter, "1"))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))                                                // terminal
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`{ "shape": "round" }`))) // final GET
	defer close()
	pl := pipeline.NewPipeline(srv)
	firstResp := &http.Response{
		StatusCode: http.StatusAccepted,
	}
	p := NewPoller(&fakePoller{Ep: srv.URL(), Fg: srv.URL()}, firstResp, pl)
	if p.Done() {
		t.Fatal("unexpected done")
	}
	type widget struct {
		Shape string `json:"shape"`
	}
	var w widget
	resp, err := p.PollUntilDone(context.Background(), time.Second, &w)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	if w.Shape != "round" {
		t.Fatalf("unexpected result %s", w.Shape)
	}
	resp, err = p.Poll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}

func TestNewPollerFail1(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithStatusCode(http.StatusConflict)) // terminal
	defer close()
	pl := pipeline.NewPipeline(srv)
	firstResp := &http.Response{
		StatusCode: http.StatusAccepted,
	}
	p := NewPoller(&fakePoller{Ep: srv.URL()}, firstResp, pl)
	resp, err := p.PollUntilDone(context.Background(), time.Second, nil)
	var respErr *shared.ResponseError
	if !errors.As(err, &respErr) {
		t.Fatalf("unexpected error type %T", err)
	}
	if respErr.StatusCode != http.StatusConflict {
		t.Fatalf("unexpected terminal status code %d", respErr.StatusCode)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
}

func TestNewPollerFail2(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(mock.WithStatusCode(http.StatusCreated)) // terminal
	defer close()
	pl := pipeline.NewPipeline(srv)
	firstResp := &http.Response{
		StatusCode: http.StatusAccepted,
	}
	p := NewPoller(&fakePoller{Ep: srv.URL()}, firstResp, pl)
	resp, err := p.PollUntilDone(context.Background(), time.Second, nil)
	var respErr *shared.ResponseError
	if !errors.As(err, &respErr) {
		t.Fatalf("unexpected error type %T", err)
	}
	if respErr.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected terminal status code %d", respErr.StatusCode)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
}

func TestNewPollerError(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithStatusCode(http.StatusAccepted))
	srv.AppendError(errors.New("fatal"))
	defer close()
	pl := pipeline.NewPipeline(srv)
	firstResp := &http.Response{
		StatusCode: http.StatusAccepted,
	}
	p := NewPoller(&fakePoller{Ep: srv.URL()}, firstResp, pl)
	resp, err := p.PollUntilDone(context.Background(), time.Second, nil)
	if err == nil {
		t.Fatal("unexpected nil error")
	} else if s := err.Error(); s != "fatal" {
		t.Fatalf("unexpected error %s", s)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
}
