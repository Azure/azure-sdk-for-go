//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestResponseUnmarshalXML(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// include UTF8 BOM
	srv.SetResponse(mock.WithBody([]byte("\xef\xbb\xbf<testXML><SomeInt>1</SomeInt><SomeString>s</SomeString></testXML>")))
	pl := NewPipeline(srv)
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasStatusCode(resp, http.StatusOK) {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	var tx testXML
	if err := UnmarshalAsXML(resp, &tx); err != nil {
		t.Fatalf("unexpected error unmarshalling: %v", err)
	}
	if tx.SomeInt != 1 || tx.SomeString != "s" {
		t.Fatal("unexpected value")
	}
}

func TestResponseFailureStatusCode(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusForbidden))
	pl := NewPipeline(srv)
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if HasStatusCode(resp, http.StatusOK) {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

func TestResponseUnmarshalJSON(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte(`{ "someInt": 1, "someString": "s" }`)))
	pl := NewPipeline(srv)
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasStatusCode(resp, http.StatusOK) {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	var tx testJSON
	if err := UnmarshalAsJSON(resp, &tx); err != nil {
		t.Fatalf("unexpected error unmarshalling: %v", err)
	}
	if tx.SomeInt != 1 || tx.SomeString != "s" {
		t.Fatal("unexpected value")
	}
}

func TestResponseUnmarshalJSONskipDownload(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte(`{ "someInt": 1, "someString": "s" }`)))
	pl := NewPipeline(srv)
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.SkipBodyDownload()
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasStatusCode(resp, http.StatusOK) {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	var tx testJSON
	if err := UnmarshalAsJSON(resp, &tx); err != nil {
		t.Fatalf("unexpected error unmarshalling: %v", err)
	}
	if tx.SomeInt != 1 || tx.SomeString != "s" {
		t.Fatal("unexpected value")
	}
}

func TestResponseUnmarshalJSONNoBody(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte{}))
	pl := NewPipeline(srv)
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasStatusCode(resp, http.StatusOK) {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if err := UnmarshalAsJSON(resp, nil); err != nil {
		t.Fatalf("unexpected error unmarshalling: %v", err)
	}
}

func TestResponseUnmarshalXMLNoBody(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte{}))
	pl := NewPipeline(srv)
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasStatusCode(resp, http.StatusOK) {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if err := UnmarshalAsXML(resp, nil); err != nil {
		t.Fatalf("unexpected error unmarshalling: %v", err)
	}
}

func TestRetryAfter(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{},
	}
	if d := RetryAfter(resp); d > 0 {
		t.Fatalf("unexpected retry-after value %d", d)
	}
	resp.Header.Set(headerRetryAfter, "300")
	d := RetryAfter(resp)
	if d <= 0 {
		t.Fatal("expected retry-after value from seconds")
	}
	if d != 300*time.Second {
		t.Fatalf("expected 300 seconds, got %d", d/time.Second)
	}
	atDate := time.Now().Add(600 * time.Second)
	resp.Header.Set(headerRetryAfter, atDate.Format(time.RFC1123))
	d = RetryAfter(resp)
	if d <= 0 {
		t.Fatal("expected retry-after value from date")
	}
	// d will not be exactly 600 seconds but it will be close
	if s := d / time.Second; s < 598 || s > 602 {
		t.Fatalf("expected ~600 seconds, got %d", s)
	}
}

func TestResponseUnmarshalAsByteArrayURLFormat(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte(`"YSBzdHJpbmcgdGhhdCBnZXRzIGVuY29kZWQgd2l0aCBiYXNlNjR1cmw"`)))
	pl := NewPipeline(srv)
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasStatusCode(resp, http.StatusOK) {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	var ba []byte
	if err := UnmarshalAsByteArray(resp, &ba, Base64URLFormat); err != nil {
		t.Fatalf("unexpected error unmarshalling: %v", err)
	}
	if string(ba) != "a string that gets encoded with base64url" {
		t.Fatalf("bad payload, got %s", string(ba))
	}
}

func TestResponseUnmarshalAsByteArrayStdFormat(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte(`"YSBzdHJpbmcgdGhhdCBnZXRzIGVuY29kZWQgd2l0aCBiYXNlNjR1cmw="`)))
	pl := NewPipeline(srv)
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasStatusCode(resp, http.StatusOK) {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	var ba []byte
	if err := UnmarshalAsByteArray(resp, &ba, Base64StdFormat); err != nil {
		t.Fatalf("unexpected error unmarshalling: %v", err)
	}
	if string(ba) != "a string that gets encoded with base64url" {
		t.Fatalf("bad payload, got %s", string(ba))
	}
}
