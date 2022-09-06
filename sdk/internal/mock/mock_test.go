//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package mock

import (
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"golang.org/x/net/http2"
)

func TestStaticResponse(t *testing.T) {
	srv, close := NewServer()
	defer close()
	srv.SetResponse(WithStatusCode(http.StatusNoContent))
	if srv.static == nil {
		t.Fatal("missing static response")
	}
	if srv.isErrorResp() {
		t.Fatal("unexpected error response")
	}
	req, err := http.NewRequest(http.MethodGet, srv.URL(), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := srv.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	resp, err = srv.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}

func TestStaticError(t *testing.T) {
	staticError := errors.New("static error")
	srv, close := NewServer()
	defer close()
	srv.SetError(staticError)
	if srv.static == nil {
		t.Fatal("missing static response")
	}
	if !srv.isErrorResp() {
		t.Fatal("unexpected response")
	}
	req, err := http.NewRequest(http.MethodGet, srv.URL(), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := srv.Do(req)
	if !errors.Is(err, staticError) {
		t.Fatalf("unexpected error %v", err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	resp, err = srv.Do(req)
	if !errors.Is(err, staticError) {
		t.Fatalf("unexpected error %v", err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
}

func TestAppendedResponses(t *testing.T) {
	srv, close := NewServer()
	defer close()
	srv.AppendResponse(WithStatusCode(http.StatusAccepted))
	srv.AppendResponse(WithStatusCode(http.StatusOK))
	if l := len(srv.resp); l != 2 {
		t.Fatalf("unexpected response count %d", l)
	}
	req, err := http.NewRequest(http.MethodGet, srv.URL(), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := srv.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	resp, err = srv.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	if l := len(srv.resp); l != 0 {
		t.Fatalf("expected zero count, got %d", l)
	}
}

func TestAppendedErrors(t *testing.T) {
	staticError1 := errors.New("static error 1")
	staticError2 := errors.New("static error 2")
	srv, close := NewServer()
	defer close()
	srv.AppendError(staticError1)
	srv.AppendError(staticError2)
	if l := len(srv.resp); l != 2 {
		t.Fatalf("unexpected response count %d", l)
	}
	req, err := http.NewRequest(http.MethodGet, srv.URL(), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := srv.Do(req)
	if !errors.Is(err, staticError1) {
		t.Fatalf("unexpected error %v", err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	resp, err = srv.Do(req)
	if !errors.Is(err, staticError2) {
		t.Fatalf("unexpected error %v", err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
}

func TestRepeatedResponse(t *testing.T) {
	srv, close := NewServer()
	defer close()
	const repeatCount = 10
	srv.RepeatResponse(repeatCount, WithStatusCode(http.StatusInternalServerError))
	if l := len(srv.resp); l != repeatCount {
		t.Fatalf("unexpected response count %d", l)
	}
	req, err := http.NewRequest(http.MethodGet, srv.URL(), nil)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < repeatCount; i++ {
		resp, err := srv.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("unexpected status code %d", resp.StatusCode)
		}
	}
	if srv.Requests() != repeatCount {
		t.Fatalf("expected request count %d, got %d", repeatCount, srv.Requests())
	}
	if l := len(srv.resp); l != 0 {
		t.Fatalf("expected zero count, got %d", l)
	}
}

func TestRepeatedError(t *testing.T) {
	repeatError := errors.New("repeated error")
	srv, close := NewServer()
	defer close()
	const repeatCount = 10
	srv.RepeatError(repeatCount, repeatError)
	if l := len(srv.resp); l != repeatCount {
		t.Fatalf("unexpected response count %d", l)
	}
	req, err := http.NewRequest(http.MethodGet, srv.URL(), nil)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < repeatCount; i++ {
		resp, err := srv.Do(req)
		if !errors.Is(err, repeatError) {
			t.Fatalf("unexpected error %v", err)
		}
		if resp != nil {
			t.Fatal("expected nil response")
		}
	}
	if srv.Requests() != repeatCount {
		t.Fatalf("expected request count %d, got %d", repeatCount, srv.Requests())
	}
	if l := len(srv.resp); l != 0 {
		t.Fatalf("expected zero count, got %d", l)
	}
}

func TestComplexResponse(t *testing.T) {
	srv, close := NewServer()
	defer close()
	const body = "this is the response body"
	srv.AppendResponse(
		WithStatusCode(http.StatusOK),
		WithBody([]byte(body)),
		WithHeader("some", "value"),
		WithSlowResponse(2*time.Second),
	)
	req, err := http.NewRequest(http.MethodGet, srv.URL(), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := srv.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	if h := resp.Header.Get("some"); h != "value" {
		t.Fatalf("unexpected header value %s", h)
	}
	r, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(r) != body {
		t.Fatalf("unexpected response body %s", string(r))
	}
}

func TestComplexResponseTLS(t *testing.T) {
	srv, close := NewTLSServer()
	defer close()
	const body = "this is the response body"
	srv.AppendResponse(
		WithStatusCode(http.StatusOK),
		WithBody([]byte(body)),
		WithHeader("some", "value"),
		WithSlowResponse(2*time.Second),
	)
	req, err := http.NewRequest(http.MethodGet, srv.URL(), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := srv.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	if h := resp.Header.Get("some"); h != "value" {
		t.Fatalf("unexpected header value %s", h)
	}
	r, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(r) != body {
		t.Fatalf("unexpected response body %s", string(r))
	}
}

func TestTLSServerConfig(t *testing.T) {
	cfg := &http.Server{}
	if err := http2.ConfigureServer(cfg, new(http2.Server)); err != nil {
		t.Fatal(err)
	}
	srv, close := NewTLSServer(WithTLSConfig(cfg.TLSConfig))
	defer close()
	const body = "this is the response body"
	srv.AppendResponse(
		WithStatusCode(http.StatusOK),
		WithBody([]byte(body)),
		WithHeader("some", "value"),
		WithSlowResponse(2*time.Second),
	)
	req, err := http.NewRequest(http.MethodGet, srv.URL(), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := srv.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	if h := resp.Header.Get("some"); h != "value" {
		t.Fatalf("unexpected header value %s", h)
	}
	r, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(r) != body {
		t.Fatalf("unexpected response body %s", string(r))
	}
}

func TestBodyReadError(t *testing.T) {
	srv, close := NewServer()
	defer close()
	srv.AppendResponse(WithBodyReadError())
	req, err := http.NewRequest(http.MethodGet, srv.URL(), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := srv.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	_, err = io.ReadAll(resp.Body)
	if err == nil {
		t.Fatal("unexpected nil error reading response body")
	}
	resp.Body.Close()
}

func TestPredicateSuccess(t *testing.T) {
	srv, close := NewServer()
	defer close()
	const customHeader = "custom-header"
	const customValue = "custom-value"
	srv.AppendResponse(WithPredicate(func(r *http.Request) bool {
		return r.Header.Get(customHeader) == customValue
	}), WithStatusCode(http.StatusOK))
	// predicate failure response
	srv.AppendResponse(WithStatusCode(http.StatusBadRequest))
	srv.AppendResponse(WithStatusCode(http.StatusNoContent))
	req, err := http.NewRequest(http.MethodGet, srv.URL(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(customHeader, customValue)
	resp, err := srv.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	resp, err = srv.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}

func TestPredicateFail(t *testing.T) {
	srv, close := NewServer()
	defer close()
	const customHeader = "custom-header"
	const customValue = "custom-value"
	srv.AppendResponse(WithPredicate(func(r *http.Request) bool {
		return r.Header.Get(customHeader) == customValue
	}), WithStatusCode(http.StatusOK))
	// predicate failure response
	srv.AppendResponse(WithStatusCode(http.StatusBadRequest))
	srv.AppendResponse(WithStatusCode(http.StatusNoContent))
	req, err := http.NewRequest(http.MethodGet, srv.URL(), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := srv.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	resp, err = srv.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}

func TestPredicateMultiple(t *testing.T) {
	srv, close := NewServer()
	defer close()
	const customHeader = "custom-header"
	const customValue = "custom-value"
	srv.AppendResponse(WithPredicate(func(r *http.Request) bool {
		return r.Header.Get(customHeader) == customValue
	}), WithStatusCode(http.StatusOK))
	// predicate failure response
	srv.AppendResponse(WithStatusCode(http.StatusBadRequest))
	srv.AppendResponse(WithStatusCode(http.StatusNoContent))
	srv.AppendResponse(WithPredicate(func(r *http.Request) bool {
		return r.Header.Get(customHeader) == customValue
	}), WithStatusCode(http.StatusAccepted))
	// predicate failure response
	srv.AppendResponse(WithStatusCode(http.StatusForbidden))
	srv.AppendResponse(WithStatusCode(http.StatusOK))
	req, err := http.NewRequest(http.MethodGet, srv.URL(), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := srv.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	resp, err = srv.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	req.Header.Add(customHeader, customValue)
	resp, err = srv.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	resp, err = srv.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
}
