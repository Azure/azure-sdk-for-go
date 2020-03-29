// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestChainedTokenCredential_InstantiateSuccess(t *testing.T) {
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Could not set environment variables for testing: %v", err)
	}
	secCred, err := NewClientSecretCredential(tenantID, clientID, secret, nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	envCred, err := NewEnvironmentCredential(nil)
	if err != nil {
		t.Fatalf("Could not find appropriate environment credentials")
	}
	cred, err := NewChainedTokenCredential(secCred, envCred)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred != nil {
		if len(cred.sources) != 2 {
			t.Fatalf("Expected 2 sources in the chained token credential, instead found %d", len(cred.sources))
		}
	}

}

func TestChainedTokenCredential_GetTokenSuccess(t *testing.T) {
	l := azcore.Log()
	l.SetListener(func(c azcore.LogClassification, s string) {
		fmt.Println(s)
	})
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Could not set environment variables for testing: %v", err)
	}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srvURL := srv.URL()
	secCred, err := NewClientSecretCredential(tenantID, clientID, secret, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	envCred, err := NewEnvironmentCredential(&TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Failed to create environment credential: %v", err)
	}
	cred, err := NewChainedTokenCredential(secCred, envCred)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
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
	l := azcore.Log()
	l.SetListener(func(c azcore.LogClassification, s string) {
		fmt.Println(s)
	})
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	testURL := srv.URL()
	secCred, err := NewClientSecretCredential(tenantID, clientID, wrongSecret, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &testURL})
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	cred, err := NewChainedTokenCredential(secCred)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
	var authErr *AuthenticationFailedError
	if !errors.As(err, &authErr) {
		t.Fatalf("Expected Error Type: AuthenticationFailedError, ReceivedErrorType: %T", err)
	}
	if len(err.Error()) == 0 {
		t.Fatalf("Did not create an appropriate error message")
	}
}

func TestChainedTokenCredential_GetTokenWithUnavailableCredentialInChain(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendError(&CredentialUnavailableError{CredentialType: "MockCredential", Message: "Mocking a credential unavailable error"})
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	testURL := srv.URL()
	secCred, err := NewClientSecretCredential(tenantID, clientID, wrongSecret, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &testURL})
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	// The chain has the same credential twice, since it doesn't matter what credential we add to the chain as long as it is not a nil credential.
	// Most credentials will not be instantiated if the conditions do not exist to allow them to be used, thus returning a
	// CredentialUnavailable error from the constructor. In order to test the CredentialUnavailable functionality for
	// ChainedTokenCredential we have to mock with two valid credentials, but the first will fail since the first response queued
	// in the test server is a CredentialUnavailable error.
	cred, err := NewChainedTokenCredential(secCred, secCred)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
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

func TestBearerPolicy_ChainedTokenCredential(t *testing.T) {
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Unable to initialize environment variables. Received: %v", err)
	}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	srvURL := srv.URL()
	cred, err := NewClientSecretCredential(tenantID, clientID, secret, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	chainedCred, err := NewChainedTokenCredential(cred)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pipeline := azcore.NewPipeline(
		srv,
		azcore.NewTelemetryPolicy(azcore.TelemetryOptions{}),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(azcore.RetryOptions{}),
		chainedCred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{scope}}}),
		azcore.NewRequestLogPolicy(azcore.RequestLogOptions{}))
	_, err = pipeline.Do(context.Background(), azcore.NewRequest(http.MethodGet, srv.URL()))
	if err != nil {
		t.Fatalf("Expected an empty error but receive: %v", err)
	}
}
