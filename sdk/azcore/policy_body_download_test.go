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
	req := pl.NewRequest(http.MethodGet, srv.URL())
	resp, err := req.Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Payload) == 0 {
		t.Fatal("missing payload")
	}
	if string(resp.Payload) != message {
		t.Fatalf("unexpected response: %s", string(resp.Payload))
	}
}

func TestSkipBodyDownload(t *testing.T) {
	const message = "not downloaded"
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithBody([]byte(message)))
	// download policy is automatically added during pipeline construction
	pl := NewPipeline(srv)
	req := pl.NewRequest(http.MethodGet, srv.URL())
	req.SkipBodyDownload()
	resp, err := req.Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Payload) > 0 {
		t.Fatalf("unexpected download: %s", string(resp.Payload))
	}
}
