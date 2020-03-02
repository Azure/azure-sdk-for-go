// +build go1.13

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
	pl := NewPipeline(srv, NewTelemetryPolicy(TelemetryOptions{}))
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodGet, srv.URL()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.HasStatusCode(http.StatusOK) {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	var tx testXML
	if err := resp.UnmarshalAsXML(&tx); err != nil {
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
	pl := NewPipeline(srv, NewTelemetryPolicy(TelemetryOptions{}))
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodGet, srv.URL()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.HasStatusCode(http.StatusOK) {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
}

func TestResponseUnmarshalJSON(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte(`{ "someInt": 1, "someString": "s" }`)))
	pl := NewPipeline(srv, NewTelemetryPolicy(TelemetryOptions{}))
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodGet, srv.URL()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.HasStatusCode(http.StatusOK) {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	var tx testJSON
	if err := resp.UnmarshalAsJSON(&tx); err != nil {
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
	pl := NewPipeline(srv, NewTelemetryPolicy(TelemetryOptions{}))
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodGet, srv.URL()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.HasStatusCode(http.StatusOK) {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if err := resp.UnmarshalAsJSON(nil); err != nil {
		t.Fatalf("unexpected error unmarshalling: %v", err)
	}
}

func TestResponseUnmarshalXMLNoBody(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte{}))
	pl := NewPipeline(srv, NewTelemetryPolicy(TelemetryOptions{}))
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodGet, srv.URL()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.HasStatusCode(http.StatusOK) {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if err := resp.UnmarshalAsXML(nil); err != nil {
		t.Fatalf("unexpected error unmarshalling: %v", err)
	}
}

func TestRetryAfter(t *testing.T) {
	raw := &http.Response{
		Header: http.Header{},
	}
	resp := Response{raw}
	if d, ok := resp.RetryAfter(); ok {
		t.Fatalf("unexpected retry-after value %d", d)
	}
	raw.Header.Set(HeaderRetryAfter, "300")
	d, ok := resp.RetryAfter()
	if !ok {
		t.Fatal("expected retry-after value from seconds")
	}
	if d != 300*time.Second {
		t.Fatalf("expected 300 seconds, got %d", d/time.Second)
	}
	atDate := time.Now().Add(600 * time.Second)
	raw.Header.Set(HeaderRetryAfter, atDate.Format(time.RFC1123))
	d, ok = resp.RetryAfter()
	if !ok {
		t.Fatal("expected retry-after value from date")
	}
	// d will not be exactly 600 seconds but it will be close
	if d/time.Second != 599 {
		t.Fatalf("expected ~600 seconds, got %d", d/time.Second)
	}
}
