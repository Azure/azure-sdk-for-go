// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/mock"
)

func TestUniqueRequestIDPolicy(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	pl := NewPipeline(srv, NewUniqueRequestIDPolicy())
	req := pl.NewRequest(http.MethodGet, srv.URL())
	resp, err := req.Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Request.Header.Get(xMsClientRequestID) == "" {
		t.Fatal("missing request ID header")
	}
}

func TestUniqueRequestIDPolicyUserDefined(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	pl := NewPipeline(srv, NewUniqueRequestIDPolicy())
	req := pl.NewRequest(http.MethodGet, srv.URL())
	const customID = "my-custom-id"
	req.Header.Set(xMsClientRequestID, customID)
	resp, err := req.Do(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(xMsClientRequestID); v != customID {
		t.Fatalf("unexpected request ID value: %s", v)
	}
}
