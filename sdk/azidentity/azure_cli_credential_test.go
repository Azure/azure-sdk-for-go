//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

var (
	mockCLITokenProviderSuccess = func(ctx context.Context, resource string, tenantID string) ([]byte, error) {
		return []byte(`{
  "accessToken": "mocktoken",
  "expiresOn": "2001-02-03 04:05:06.000007",
  "subscription": "mocksub",
  "tenant": "mocktenant",
  "tokenType": "Bearer"
}
`), nil
	}
	mockCLITokenProviderFailure = func(ctx context.Context, resource string, tenantID string) ([]byte, error) {
		return nil, errors.New("provider failure message")
	}
)

func TestAzureCLICredential_GetTokenSuccess(t *testing.T) {
	options := AzureCLICredentialOptions{}
	options.tokenProvider = mockCLITokenProviderSuccess
	cred, err := NewAzureCLICredential(&options)
	if err != nil {
		t.Fatal(err)
	}
	at, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatal(err)
	}
	if at.Token != "mocktoken" {
		t.Fatalf("unexpected access token %q", at.Token)
	}
	expected := time.Date(2001, 2, 3, 4, 5, 6, 7000, time.Local).UTC()
	if actual := at.ExpiresOn; !actual.Equal(expected) || actual.Location() != time.UTC {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}

func TestAzureCLICredential_GetTokenInvalidToken(t *testing.T) {
	options := AzureCLICredentialOptions{}
	options.tokenProvider = mockCLITokenProviderFailure
	cred, err := NewAzureCLICredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one.")
	}
}

func TestAzureCLICredential_TenantID(t *testing.T) {
	expected := "expected-tenant-id"
	called := false
	options := AzureCLICredentialOptions{
		TenantID: expected,
		tokenProvider: func(ctx context.Context, resource, tenantID string) ([]byte, error) {
			called = true
			if tenantID != expected {
				t.Fatal("Unexpected tenant ID: " + tenantID)
			}
			return mockCLITokenProviderSuccess(ctx, resource, tenantID)
		},
	}
	cred, err := NewAzureCLICredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !called {
		t.Fatal("token provider wasn't called")
	}
}
