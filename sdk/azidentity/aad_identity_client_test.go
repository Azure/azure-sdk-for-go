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

func TestAzurePublicCloudParse(t *testing.T) {
	_, err := url.Parse(AzurePublicCloud)
	if err != nil {
		t.Fatalf("Failed to parse default authority host: %v", err)
	}
}

func TestAzureChinaParse(t *testing.T) {
	_, err := url.Parse(AzureChina)
	if err != nil {
		t.Fatalf("Failed to parse AzureChina authority host: %v", err)
	}
}

func TestAzureGermanyParse(t *testing.T) {
	_, err := url.Parse(AzureGermany)
	if err != nil {
		t.Fatalf("Failed to parse AzureGermany authority host: %v", err)
	}
}

func TestAzureGovernmentParse(t *testing.T) {
	_, err := url.Parse(AzureGovernment)
	if err != nil {
		t.Fatalf("Failed to parse AzureGovernment authority host: %v", err)
	}
}
func TestTelemetryDefaultUserAgent(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := pipelineOptions{
		HTTPClient: srv,
	}
	client, err := newAADIdentityClient(srv.URL(), options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	resp, err := client.pipeline.Do(req)
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

func TestTelemetryCustom(t *testing.T) {
	customTelemetry := "customvalue"
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := pipelineOptions{
		HTTPClient: srv,
	}
	options.Telemetry.Value = customTelemetry
	client, err := newAADIdentityClient(srv.URL(), options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	resp, err := client.pipeline.Do(req)
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
