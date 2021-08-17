//go:build go1.16
// +build go1.16

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

var defaultTelemetry = UserAgent + " " + platformInfo

func TestPolicyTelemetryDefault(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	pl := NewPipeline(srv, NewTelemetryPolicy(nil))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(headerUserAgent); v != defaultTelemetry {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}

func TestPolicyTelemetryWithCustomInfo(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	const testValue = "azcore_test"
	pl := NewPipeline(srv, NewTelemetryPolicy(&TelemetryOptions{Value: testValue}))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(headerUserAgent); v != fmt.Sprintf("%s %s", testValue, defaultTelemetry) {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}

func TestPolicyTelemetryPreserveExisting(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	pl := NewPipeline(srv, NewTelemetryPolicy(nil))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	const otherValue = "this should stay"
	req.Header.Set(headerUserAgent, otherValue)
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(headerUserAgent); v != fmt.Sprintf("%s %s", defaultTelemetry, otherValue) {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}

func TestPolicyTelemetryWithAppID(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	const appID = "my_application"
	pl := NewPipeline(srv, NewTelemetryPolicy(&TelemetryOptions{ApplicationID: appID}))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(headerUserAgent); v != fmt.Sprintf("%s %s", appID, defaultTelemetry) {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}

func TestPolicyTelemetryWithAppIDAndReqTelemetry(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	const appID = "my_application"
	pl := NewPipeline(srv, NewTelemetryPolicy(&TelemetryOptions{ApplicationID: appID}))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.Telemetry("TestPolicyTelemetryWithAppIDAndReqTelemetry")
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(headerUserAgent); v != fmt.Sprintf("%s %s %s", "TestPolicyTelemetryWithAppIDAndReqTelemetry", appID, defaultTelemetry) {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}

func TestPolicyTelemetryWithAppIDSanitized(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	const appID = "This will get the spaces removed and truncated."
	pl := NewPipeline(srv, NewTelemetryPolicy(&TelemetryOptions{ApplicationID: appID}))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	const newAppID = "This/will/get/the/spaces"
	if v := resp.Request.Header.Get(headerUserAgent); v != fmt.Sprintf("%s %s", newAppID, defaultTelemetry) {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}

func TestPolicyTelemetryPreserveExistingWithAppID(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	const appID = "my_application"
	pl := NewPipeline(srv, NewTelemetryPolicy(&TelemetryOptions{ApplicationID: appID}))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	const otherValue = "this should stay"
	req.Header.Set(headerUserAgent, otherValue)
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(headerUserAgent); v != fmt.Sprintf("%s %s %s", appID, defaultTelemetry, otherValue) {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}

func TestPolicyTelemetryDisabled(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse()
	const appID = "my_application"
	pl := NewPipeline(srv, NewTelemetryPolicy(&TelemetryOptions{ApplicationID: appID, Disabled: true}))
	req, err := NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req.Telemetry("this should be ignored")
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := resp.Request.Header.Get(headerUserAgent); v != "" {
		t.Fatalf("unexpected user agent value: %s", v)
	}
}
