//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestPolicyTelemetryDefault(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	pl := exported.NewPipeline(srv, NewTelemetryPolicy("test", "v1.2.3", nil))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(shared.HeaderUserAgent); v != "azsdk-go-test/v1.2.3 "+platformInfo {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}

func TestPolicyTelemetryDefaultFullQualified(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	pl := exported.NewPipeline(srv, NewTelemetryPolicy("github.com/foo/bar/test", "v1.2.3", nil))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(shared.HeaderUserAgent); v != "azsdk-go-test/v1.2.3 "+platformInfo {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}

func TestPolicyTelemetryPreserveExisting(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	pl := exported.NewPipeline(srv, NewTelemetryPolicy("test", "v1.2.3", nil))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	const otherValue = "this should stay"
	req.Raw().Header.Set(shared.HeaderUserAgent, otherValue)
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(shared.HeaderUserAgent); v != fmt.Sprintf("%s %s", "azsdk-go-test/v1.2.3 "+platformInfo, otherValue) {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}

func TestPolicyTelemetryWithAppID(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	const appID = "my_application"
	pl := exported.NewPipeline(srv, NewTelemetryPolicy("test", "v1.2.3", &policy.TelemetryOptions{ApplicationID: appID}))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(shared.HeaderUserAgent); v != fmt.Sprintf("%s %s", appID, "azsdk-go-test/v1.2.3 "+platformInfo) {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}

func TestPolicyTelemetryWithAppIDSanitized(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	const appID = "This will get the spaces removed and truncated."
	pl := exported.NewPipeline(srv, NewTelemetryPolicy("test", "v1.2.3", &policy.TelemetryOptions{ApplicationID: appID}))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	const newAppID = "This/will/get/the/spaces"
	if v := resp.Request.Header.Get(shared.HeaderUserAgent); v != fmt.Sprintf("%s %s", newAppID, "azsdk-go-test/v1.2.3 "+platformInfo) {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}

func TestPolicyTelemetryPreserveExistingWithAppID(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	const appID = "my_application"
	pl := exported.NewPipeline(srv, NewTelemetryPolicy("test", "v1.2.3", &policy.TelemetryOptions{ApplicationID: appID}))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	const otherValue = "this should stay"
	req.Raw().Header.Set(shared.HeaderUserAgent, otherValue)
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(shared.HeaderUserAgent); v != fmt.Sprintf("%s %s %s", appID, "azsdk-go-test/v1.2.3 "+platformInfo, otherValue) {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}

func TestPolicyTelemetryDisabled(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	const appID = "my_application"
	pl := exported.NewPipeline(srv, NewTelemetryPolicy("test", "v1.2.3", &policy.TelemetryOptions{ApplicationID: appID, Disabled: true}))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(shared.HeaderUserAgent); v != "" {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}
