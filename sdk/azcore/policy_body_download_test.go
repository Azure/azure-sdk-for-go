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

func TestDownloadBody(t *testing.T) {
	const message = "downloaded"
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte(message)))
	// download policy is automatically added during pipeline construction
	pl := NewPipeline(srv)
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodGet, srv.URL()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.payload()) == 0 {
		t.Fatal("missing payload")
	}
	if string(resp.payload()) != message {
		t.Fatalf("unexpected response: %s", string(resp.payload()))
	}
}

func TestSkipBodyDownload(t *testing.T) {
	const message = "not downloaded"
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte(message)))
	// download policy is automatically added during pipeline construction
	pl := NewPipeline(srv)
	req := NewRequest(http.MethodGet, srv.URL())
	req.SkipBodyDownload()
	resp, err := pl.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.payload()) > 0 {
		t.Fatalf("unexpected download: %s", string(resp.payload()))
	}
}
