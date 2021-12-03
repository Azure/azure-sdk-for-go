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
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestAzurePublicCloudParse(t *testing.T) {
	_, err := url.Parse(string(AzurePublicCloud))
	if err != nil {
		t.Fatalf("Failed to parse default authority host: %v", err)
	}
}

func TestAzureChinaParse(t *testing.T) {
	_, err := url.Parse(string(AzureChina))
	if err != nil {
		t.Fatalf("Failed to parse AzureChina authority host: %v", err)
	}
}

func TestAzureGovernmentParse(t *testing.T) {
	_, err := url.Parse(string(AzureGovernment))
	if err != nil {
		t.Fatalf("Failed to parse AzureGovernment authority host: %v", err)
	}
}
func TestTelemetryDefaultUserAgent(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := azcore.ClientOptions{Transport: srv}
	client, err := newAADIdentityClient(srv.URL(), &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
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
	if ua := resp.Request.Header.Get(headerUserAgent); !strings.HasPrefix(ua, "azsdk-go-"+component+"/"+version) {
		t.Fatalf("unexpected User-Agent %s", ua)
	}
}

func TestTelemetryCustom(t *testing.T) {
	customTelemetry := "customvalue"
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	options := azcore.ClientOptions{Transport: srv}
	options.Telemetry.ApplicationID = customTelemetry
	client, err := newAADIdentityClient(srv.URL(), &options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
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
	if ua := resp.Request.Header.Get(headerUserAgent); !strings.HasPrefix(ua, customTelemetry+" "+"azsdk-go-"+component+"/"+version) {
		t.Fatalf("unexpected User-Agent %s", ua)
	}
}
