//go:build go1.18
// +build go1.18

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
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type userAgentValidatingPolicy struct {
	t     *testing.T
	appID string
}

func (p userAgentValidatingPolicy) Do(req *policy.Request) (*http.Response, error) {
	expected := "azsdk-go-" + component + "/" + version
	if p.appID != "" {
		expected = p.appID + " " + expected
	}
	if ua := req.Raw().Header.Get("User-Agent"); !strings.HasPrefix(ua, expected) {
		p.t.Fatalf("unexpected User-Agent %s", ua)
	}
	return req.Next()
}

func TestIMDSEndpointParse(t *testing.T) {
	_, err := url.Parse(imdsEndpoint)
	if err != nil {
		t.Fatalf("Failed to parse the IMDS endpoint: %v", err)
	}
}

func TestManagedIdentityClient_UserAgent(t *testing.T) {
	options := ManagedIdentityCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: &mockSTS{}, PerCallPolicies: []policy.Policy{userAgentValidatingPolicy{t: t}},
		},
	}
	client, err := newManagedIdentityClient(&options)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.authenticate(context.Background(), nil, []string{liveTestScope})
	if err != nil {
		t.Fatal(err)
	}
}

func TestManagedIdentityClient_ApplicationID(t *testing.T) {
	appID := "customvalue"
	options := ManagedIdentityCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: &mockSTS{}, PerCallPolicies: []policy.Policy{userAgentValidatingPolicy{t: t, appID: appID}},
		},
	}
	options.Telemetry.ApplicationID = appID
	client, err := newManagedIdentityClient(&options)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.authenticate(context.Background(), nil, []string{liveTestScope})
	if err != nil {
		t.Fatal(err)
	}
}
