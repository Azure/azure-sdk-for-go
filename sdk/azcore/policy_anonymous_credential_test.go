// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestAnonymousCredential(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))
	pl := NewPipeline(srv, AnonymousCredential().AuthenticationPolicy(AuthenticationPolicyOptions{}))
	req := NewRequest(http.MethodGet, srv.URL())
	resp, err := pl.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(req.Header, resp.Request.Header) {
		t.Fatal("unexpected modification to request headers")
	}
}
