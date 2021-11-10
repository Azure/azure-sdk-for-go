// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
)

func TestChainedTokenCredential_InstantiateSuccess(t *testing.T) {
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Could not set environment variables for testing: %v", err)
	}
	secCred, err := NewClientSecretCredential(fakeTenantID, fakeClientID, secret, nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	envCred, err := NewEnvironmentCredential(nil)
	if err != nil {
		t.Fatalf("Could not find appropriate environment credentials")
	}
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{secCred, envCred}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred != nil {
		if len(cred.sources) != 2 {
			t.Fatalf("Expected 2 sources in the chained token credential, instead found %d", len(cred.sources))
		}
	}
}

func TestChainedTokenCredential_InstantiateFailure(t *testing.T) {
	secCred, err := NewClientSecretCredential(fakeTenantID, fakeClientID, secret, nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	_, err = NewChainedTokenCredential([]azcore.TokenCredential{secCred, nil}, nil)
	if err == nil {
		t.Fatalf("Expected an error for sending a nil credential in the chain")
	}
	_, err = NewChainedTokenCredential([]azcore.TokenCredential{}, nil)
	if err == nil {
		t.Fatalf("Expected an error for not sending any credential sources")
	}
}

func TestChainedTokenCredential_GetTokenSuccess(t *testing.T) {
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Could not set environment variables for testing: %v", err)
	}
	secCred, err := NewClientSecretCredential(fakeTenantID, fakeClientID, secret, nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	secCred.client = fakeConfidentialClient{
		ar: confidential.AuthResult{
			AccessToken: tokenValue,
			ExpiresOn:   time.Now().Add(1 * time.Hour),
		},
	}
	envCred, err := NewEnvironmentCredential(nil)
	if err != nil {
		t.Fatalf("Failed to create environment credential: %v", err)
	}
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{secCred, envCred}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to get a token but expected none")
	}
	if tk.Token != tokenValue {
		t.Fatalf("Received an incorrect access token")
	}
	if tk.ExpiresOn.IsZero() {
		t.Fatalf("Received an incorrect time in the response")
	}
}

func TestChainedTokenCredential_GetTokenFail(t *testing.T) {
	secCred, err := NewClientSecretCredential(fakeTenantID, fakeClientID, secret, nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	secCred.client = fakeConfidentialClient{
		err: errors.New("invalid client secret"),
	}
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{secCred}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
	var authErr AuthenticationFailedError
	if !errors.As(err, &authErr) {
		t.Fatalf("Expected AuthenticationFailedError, received %T", err)
	}
	if len(err.Error()) == 0 {
		t.Fatalf("Did not create an appropriate error message")
	}
}

type unavailableCredential struct{}

func (*unavailableCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (token *azcore.AccessToken, err error) {
	return nil, newCredentialUnavailableError("unavailableCredential", "is unavailable")
}

func TestChainedTokenCredential_GetTokenWithUnavailableCredentialInChain(t *testing.T) {
	secCred, err := NewClientSecretCredential(fakeTenantID, fakeClientID, secret, nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	secCred.client = fakeConfidentialClient{ar: confidential.AuthResult{AccessToken: tokenValue, ExpiresOn: time.Now().Add(time.Hour)}}
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{&unavailableCredential{}, secCred}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to get a token but expected none")
	}
	if tk.Token != tokenValue {
		t.Fatalf("Received an incorrect access token")
	}
	if tk.ExpiresOn.IsZero() {
		t.Fatalf("Received an incorrect time in the response")
	}
}
