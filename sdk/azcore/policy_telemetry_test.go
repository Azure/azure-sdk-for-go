// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestPolicyTelemetryDefault(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	pl := NewPipeline(srv, NewTelemetryPolicy(TelemetryOptions{}))
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodGet, srv.URL()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(HeaderUserAgent); v != platformInfo {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}

func TestPolicyTelemetryWithCustomInfo(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	const testValue = "azcore_test"
	pl := NewPipeline(srv, NewTelemetryPolicy(TelemetryOptions{Value: testValue}))
	resp, err := pl.Do(context.Background(), NewRequest(http.MethodGet, srv.URL()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(HeaderUserAgent); v != fmt.Sprintf("%s %s", testValue, platformInfo) {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}
