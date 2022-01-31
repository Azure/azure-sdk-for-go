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
		t.Fatal(err)
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
		t.Fatal(err)
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
		t.Fatal(err)
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

func TestChainedTokenCredential_MultipleCredentialsGetTokenUnavailable(t *testing.T) {
	credential1 := &TestCredential{responses: []testCredentialResponse{
		{err: newCredentialUnavailableError("unavailableCredential1", "Unavailable expected error")},
	}}
	credential2 := &TestCredential{responses: []testCredentialResponse{
		{err: newCredentialUnavailableError("unavailableCredential2", "Unavailable expected error")},
	}}
	credential3 := &TestCredential{responses: []testCredentialResponse{
		{err: newCredentialUnavailableError("unavailableCredential3", "Unavailable expected error")},
	}}
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{credential1, credential2, credential3}, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
	var authErr credentialUnavailableError
	if !errors.As(err, &authErr) {
		t.Fatalf("Expected CredentialUnavailableError, received %T", err)
	}
	expectedError := `ChainedTokenCredential: failed to acquire a token.
Attempted credentials:
	unavailableCredential1: Unavailable expected error
	unavailableCredential2: Unavailable expected error
	unavailableCredential3: Unavailable expected error`
	if err.Error() != expectedError {
		t.Fatalf("Did not create an appropriate error message.\n\nReceived:\n%s\n\nExpected:\n%s", err.Error(), expectedError)
	}
}

func TestChainedTokenCredential_MultipleCredentialsGetTokenAuthenticationFailed(t *testing.T) {
	credential1 := &TestCredential{responses: []testCredentialResponse{
		{err: newCredentialUnavailableError("unavailableCredential1", "Unavailable expected error")},
	}}
	credential2 := &TestCredential{responses: []testCredentialResponse{
		{err: newCredentialUnavailableError("unavailableCredential2", "Unavailable expected error")},
	}}
	credential3 := &TestCredential{responses: []testCredentialResponse{
		{err: newAuthenticationFailedError(newCredentialUnavailableError("authenticationFailedCredential3", "Authentication failed expected error"), nil)},
	}}
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{credential1, credential2, credential3}, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
	var authErr AuthenticationFailedError
	if !errors.As(err, &authErr) {
		t.Fatalf("Expected AuthenticationFailedError, received %T", err)
	}
	expectedError := `ChainedTokenCredential: failed to acquire a token.
Attempted credentials:
	unavailableCredential1: Unavailable expected error
	unavailableCredential2: Unavailable expected error
	authenticationFailedCredential3: Authentication failed expected error`
	if err.Error() != expectedError {
		t.Fatalf("Did not create an appropriate error message.\n\nReceived:\n%s\n\nExpected:\n%s", err.Error(), expectedError)
	}
}

func TestChainedTokenCredential_MultipleCredentialsGetTokenCustomName(t *testing.T) {
	credential1 := &TestCredential{responses: []testCredentialResponse{
		{err: newCredentialUnavailableError("unavailableCredential1", "Unavailable expected error")},
	}}
	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{credential1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	cred.name = "CustomNameCredential"
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{liveTestScope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
	var authErr credentialUnavailableError
	if !errors.As(err, &authErr) {
		t.Fatalf("Expected credentialUnavailableError, received %T", err)
	}
	expectedError := `CustomNameCredential: failed to acquire a token.
Attempted credentials:
	unavailableCredential1: Unavailable expected error`
	if err.Error() != expectedError {
		t.Fatalf("Did not create an appropriate error message.\n\nReceived:\n%s\n\nExpected:\n%s", err.Error(), expectedError)
	}
}

// TestCredential response
type testCredentialResponse struct {
	token *azcore.AccessToken
	err   error
}

// Credential used for testing
type TestCredential struct {
	getTokenCalls int
	responses     []testCredentialResponse
}

func (c *TestCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (token *azcore.AccessToken, err error) {
	index := c.getTokenCalls
	c.getTokenCalls += 1
	response := c.responses[index]
	return response.token, response.err
}

func testGoodGetTokenResponse(t *testing.T, token *azcore.AccessToken, err error) {
	if err != nil {
		t.Fatalf("Received an error when attempting to get a token but expected none")
	}
	if token.Token != tokenValue {
		t.Fatalf("Received an incorrect access token")
	}
	if token.ExpiresOn.IsZero() {
		t.Fatalf("Received an incorrect time in the response")
	}
}

func TestChainedTokenCredential_RepeatedGetTokenWithSuccessfulCredential(t *testing.T) {
	failedCredential := &TestCredential{responses: []testCredentialResponse{
		{err: newCredentialUnavailableError("MockCredential", "Mocking a credential unavailable error")},
		{err: newCredentialUnavailableError("MockCredential", "Mocking a credential unavailable error")},
	}}
	successfulCredential := &TestCredential{responses: []testCredentialResponse{
		{token: &azcore.AccessToken{Token: tokenValue, ExpiresOn: time.Now()}},
		{token: &azcore.AccessToken{Token: tokenValue, ExpiresOn: time.Now()}},
	}}

	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{failedCredential, successfulCredential}, nil)
	if err != nil {
		t.Fatal(err)
	}

	getTokenOptions := policy.TokenRequestOptions{Scopes: []string{liveTestScope}}

	tk, err := cred.GetToken(context.Background(), getTokenOptions)
	testGoodGetTokenResponse(t, tk, err)
	if failedCredential.getTokenCalls != 1 {
		t.Fatal("The failed credential getToken should have been called once")
	}
	if successfulCredential.getTokenCalls != 1 {
		t.Fatalf("The successful credential getToken should have been called once")
	}
	tk2, err2 := cred.GetToken(context.Background(), getTokenOptions)
	testGoodGetTokenResponse(t, tk2, err2)
	if failedCredential.getTokenCalls != 1 {
		t.Fatalf("The failed credential getToken should not have been called again")
	}
	if successfulCredential.getTokenCalls != 2 {
		t.Fatalf("The successful credential getToken should have been called twice")
	}
}

func TestChainedTokenCredential_RepeatedGetTokenWithSuccessfulCredentialWithRetrySources(t *testing.T) {
	failedCredential := &TestCredential{responses: []testCredentialResponse{
		{err: newCredentialUnavailableError("MockCredential", "Mocking a credential unavailable error")},
		{err: newCredentialUnavailableError("MockCredential", "Mocking a credential unavailable error")},
	}}
	successfulCredential := &TestCredential{responses: []testCredentialResponse{
		{token: &azcore.AccessToken{Token: tokenValue, ExpiresOn: time.Now()}},
		{token: &azcore.AccessToken{Token: tokenValue, ExpiresOn: time.Now()}},
	}}

	cred, err := NewChainedTokenCredential([]azcore.TokenCredential{failedCredential, successfulCredential}, &ChainedTokenCredentialOptions{RetrySources: true})
	if err != nil {
		t.Fatal(err)
	}

	getTokenOptions := policy.TokenRequestOptions{Scopes: []string{liveTestScope}}

	tk, err := cred.GetToken(context.Background(), getTokenOptions)
	testGoodGetTokenResponse(t, tk, err)
	if failedCredential.getTokenCalls != 1 {
		t.Fatalf("The failed credential getToken should have been called once")
	}
	if successfulCredential.getTokenCalls != 1 {
		t.Fatalf("The successful credential getToken should have been called once")
	}
	tk2, err2 := cred.GetToken(context.Background(), getTokenOptions)
	testGoodGetTokenResponse(t, tk2, err2)
	if failedCredential.getTokenCalls != 2 {
		t.Fatalf("The failed credential getToken should have been called twice")
	}
	if successfulCredential.getTokenCalls != 2 {
		t.Fatalf("The successful credential getToken should have been called twice")
	}
}
