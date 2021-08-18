// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armcore

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

type mockTokenCred struct{}

func (mockTokenCred) NewAuthenticationPolicy(azcore.AuthenticationOptions) azcore.Policy {
	return azcore.PolicyFunc(func(req *azcore.Request) (*azcore.Response, error) {
		return req.Next()
	})
}

func (mockTokenCred) GetToken(context.Context, azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	return &azcore.AccessToken{
		Token:     "abc123",
		ExpiresOn: time.Now().Add(1 * time.Hour),
	}, nil
}

func TestNewDefaultConnection(t *testing.T) {
	opt := ConnectionOptions{}
	con := NewDefaultConnection(mockTokenCred{}, &opt)
	if ep := con.Endpoint(); ep != AzurePublicCloud {
		t.Fatalf("unexpected endpoint %s", ep)
	}
}

func TestNewConnection(t *testing.T) {
	const customEndpoint = "https://contoso.com/fake/endpoint"
	con := NewConnection(customEndpoint, mockTokenCred{}, nil)
	if ep := con.Endpoint(); ep != customEndpoint {
		t.Fatalf("unexpected endpoint %s", ep)
	}
}

func TestNewConnectionWithOptions(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse()
	opt := ConnectionOptions{}
	opt.HTTPClient = srv
	con := NewConnection(srv.URL(), mockTokenCred{}, &opt)
	if ep := con.Endpoint(); ep != srv.URL() {
		t.Fatalf("unexpected endpoint %s", ep)
	}
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	resp, err := con.Pipeline().Do(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if ua := resp.Request.Header.Get("User-Agent"); !strings.HasPrefix(ua, UserAgent) {
		t.Fatalf("unexpected User-Agent %s", ua)
	}
}

func TestNewConnectionWithCustomTelemetry(t *testing.T) {
	const myTelemetry = "something"
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse()
	opt := ConnectionOptions{}
	opt.HTTPClient = srv
	opt.Telemetry.Value = myTelemetry
	con := NewConnection(srv.URL(), mockTokenCred{}, &opt)
	if ep := con.Endpoint(); ep != srv.URL() {
		t.Fatalf("unexpected endpoint %s", ep)
	}
	if opt.Telemetry.Value != myTelemetry {
		t.Fatalf("telemetry was modified: %s", opt.Telemetry.Value)
	}
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	resp, err := con.Pipeline().Do(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if ua := resp.Request.Header.Get("User-Agent"); !strings.HasPrefix(ua, myTelemetry+" "+UserAgent) {
		t.Fatalf("unexpected User-Agent %s", ua)
	}
}

func TestScope(t *testing.T) {
	if s := endpointToScope(AzureGermany); s != "https://management.microsoftazure.de//.default" {
		t.Fatalf("unexpected scope %s", s)
	}
	if s := endpointToScope("https://management.usgovcloudapi.net"); s != "https://management.usgovcloudapi.net//.default" {
		t.Fatalf("unexpected scope %s", s)
	}
}

func TestDisableAutoRPRegistration(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// initial response that RP is unregistered
	srv.SetResponse(mock.WithStatusCode(http.StatusConflict), mock.WithBody([]byte(rpUnregisteredResp)))
	con := NewConnection(srv.URL(), mockTokenCred{}, &ConnectionOptions{DisableRPRegistration: true})
	if ep := con.Endpoint(); ep != srv.URL() {
		t.Fatalf("unexpected endpoint %s", ep)
	}
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// log only RP registration
	azcore.Log().SetClassifications(LogRPRegistration)
	defer func() {
		// reset logging
		azcore.Log().SetClassifications()
	}()
	logEntries := 0
	azcore.Log().SetListener(func(cls azcore.LogClassification, msg string) {
		logEntries++
	})
	resp, err := con.Pipeline().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("unexpected status code %d:", resp.StatusCode)
	}
	// shouldn't be any log entries
	if logEntries != 0 {
		t.Fatalf("expected 0 log entries, got %d", logEntries)
	}
}

// policy that tracks the number of times it was invoked
type countingPolicy struct {
	count int
}

func (p *countingPolicy) Do(req *azcore.Request) (*azcore.Response, error) {
	p.count++
	return req.Next()
}

func TestConnectionWithCustomPolicies(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	// initial response is a failure to trigger retry
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	perCallPolicy := countingPolicy{}
	perRetryPolicy := countingPolicy{}
	con := NewConnection(srv.URL(), mockTokenCred{}, &ConnectionOptions{
		DisableRPRegistration: true,
		PerCallPolicies:       []azcore.Policy{&perCallPolicy},
		PerRetryPolicies:      []azcore.Policy{&perRetryPolicy},
	})
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	resp, err := con.Pipeline().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code %d", resp.StatusCode)
	}
	if perCallPolicy.count != 1 {
		t.Fatalf("unexpected per call policy count %d", perCallPolicy.count)
	}
	if perRetryPolicy.count != 2 {
		t.Fatalf("unexpected per retry policy count %d", perRetryPolicy.count)
	}
}
