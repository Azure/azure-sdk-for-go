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

func TestUniqueRequestIDPolicy(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	pl := NewPipeline(srv, NewUniqueRequestIDPolicy())
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodGet, srv.URL()))
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
	req := NewRequest(http.MethodGet, srv.URL())
	const customID = "my-custom-id"
	req.Header.Set(xMsClientRequestID, customID)
	resp, err := pl.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(xMsClientRequestID); v != customID {
		t.Fatalf("unexpected request ID value: %s", v)
	}
}
