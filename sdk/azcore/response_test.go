// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestResponseUnmarshalXML(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte("<testXML><SomeInt>1</SomeInt><SomeString>s</SomeString></testXML>")))
	pl := NewPipeline(srv, NewTelemetryPolicy(TelemetryOptions{}))
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodGet, srv.URL()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := resp.CheckStatusCode(http.StatusOK); err != nil {
		t.Fatalf("unexpected status code error: %v", err)
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
	if err = resp.CheckStatusCode(http.StatusOK); err == nil {
		t.Fatal("unexpected nil status code error")
	}
	re, ok := err.(RequestError)
	if !ok {
		t.Fatal("expected RequestError type")
	}
	if re.Response().StatusCode != http.StatusForbidden {
		t.Fatal("unexpected response")
	}
}
