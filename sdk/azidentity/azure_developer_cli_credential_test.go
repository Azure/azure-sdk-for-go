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
)

var (
	mockAzdTokenProviderSuccess = func(context.Context, []string, string) ([]byte, error) {
		return []byte(`{
  "token": "mocktoken",
  "expiresOn": "2001-02-03T04:05:06Z"
}
`), nil
	}
	mockAzdTokenProviderFailure = func(context.Context, []string, string) ([]byte, error) {
		return nil, newAuthenticationFailedError(credNameAzureCLI, "mock provider error", nil, nil)
	}
)

func TestAzureDeveloperCLICredential_DefaultChainError(t *testing.T) {
	cred, err := NewAzureDeveloperCLICredential(&AzureDeveloperCLICredentialOptions{
		inDefaultChain: true,
		tokenProvider:  mockAzdTokenProviderFailure,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = cred.GetToken(context.Background(), testTRO)
	var cu credentialUnavailable
	if !errors.As(err, &cu) {
		t.Fatalf("expected %T, got %T: %q", cu, err, err)
	}
}

func TestAzureDeveloperCLICredential_Error(t *testing.T) {
	// GetToken shouldn't invoke the CLI a second time after a failure
	authNs := 0
	expected := newCredentialUnavailableError(credNameAzureDeveloperCLI, "it didn't work")
	o := AzureDeveloperCLICredentialOptions{
		tokenProvider: func(context.Context, []string, string) ([]byte, error) {
			authNs++
			return nil, expected
		},
	}
	cred, err := NewAzureDeveloperCLICredential(&o)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cred.GetToken(context.Background(), testTRO)
	if err == nil {
		t.Fatal("expected an error")
	}
	if err != expected {
		t.Fatalf("expected %v, got %v", expected, err)
	}
	if authNs != 1 {
		t.Fatalf("expected 1 authN, got %d", authNs)
	}
}

func TestAzureDeveloperCLICredential_GetTokenSuccess(t *testing.T) {
	cred, err := NewAzureDeveloperCLICredential(&AzureDeveloperCLICredentialOptions{
		tokenProvider: mockAzdTokenProviderSuccess,
	})
	if err != nil {
		t.Fatal(err)
	}
	at, err := cred.GetToken(context.Background(), testTRO)
	if err != nil {
		t.Fatal(err)
	}
	if at.Token != "mocktoken" {
		t.Fatalf("unexpected access token %q", at.Token)
	}
	expected := time.Date(2001, 2, 3, 4, 5, 6, 000, time.UTC)
	if actual := at.ExpiresOn; !actual.Equal(expected) || actual.Location() != time.UTC {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}

func TestAzureDeveloperCLICredential_GetTokenInvalidToken(t *testing.T) {
	options := AzureDeveloperCLICredentialOptions{}
	options.tokenProvider = mockAzdTokenProviderFailure
	cred, err := NewAzureDeveloperCLICredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), testTRO)
	if err == nil {
		t.Fatalf("Expected an error but did not receive one.")
	}
}

func TestAzureDeveloperCLICredential_TenantID(t *testing.T) {
	expected := "expected-tenant-id"
	called := false
	options := AzureDeveloperCLICredentialOptions{
		TenantID: expected,
		tokenProvider: func(ctx context.Context, scopes []string, tenant string) ([]byte, error) {
			called = true
			if tenant != expected {
				t.Fatalf("unexpected tenant %q", tenant)
			}
			return mockAzdTokenProviderSuccess(ctx, scopes, tenant)
		},
	}
	cred, err := NewAzureDeveloperCLICredential(&options)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = cred.GetToken(context.Background(), testTRO)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !called {
		t.Fatal("token provider wasn't called")
	}
}
