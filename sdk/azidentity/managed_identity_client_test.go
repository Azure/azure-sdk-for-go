// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestIMDSEndpointParse(t *testing.T) {
	_, err := url.Parse(imdsEndpoint)
	if err != nil {
		t.Fatalf("Failed to parse the IMDS endpoint: %v", err)
	}
}

func TestMSITelemetryDefaultUserAgent(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ManagedIdentityCredentialOptions{
		HTTPClient: srv,
	}
	pipeline := newDefaultMSIPipeline(options)
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	resp, err := pipeline.Do(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if ua := resp.Request.Header.Get(headerUserAgent); !strings.HasPrefix(ua, UserAgent) {
		t.Fatalf("unexpected User-Agent %s", ua)
	}
}

func TestMSITelemetryCustom(t *testing.T) {
	customTelemetry := "customvalue"
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := ManagedIdentityCredentialOptions{
		HTTPClient: srv,
	}
	options.Telemetry.Value = customTelemetry
	pipeline := newDefaultMSIPipeline(options)
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	resp, err := pipeline.Do(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if ua := resp.Request.Header.Get(headerUserAgent); !strings.HasPrefix(ua, customTelemetry+" "+UserAgent) {
		t.Fatalf("unexpected User-Agent %s", ua)
	}
}
