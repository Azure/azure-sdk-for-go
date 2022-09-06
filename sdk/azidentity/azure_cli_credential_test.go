//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

var (
	mockCLITokenProviderSuccess = func(ctx context.Context, resource string, tenantID string) ([]byte, error) {
		return []byte(" {\"accessToken\":\"mocktoken\" , " +
			"\"expiresOn\": \"2007-01-01 01:01:01.079627\"," +
			"\"subscription\": \"mocksub\"," +
			"\"tenant\": \"mocktenant\"," +
			"\"tokenType\": \"mocktype\"}"), nil
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
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	at, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}
	if len(at.Token) == 0 {
		t.Fatalf(("Did not receive a token"))
	}
	if at.Token != "mocktoken" {
		t.Fatalf(("Did not receive the correct access token"))
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
