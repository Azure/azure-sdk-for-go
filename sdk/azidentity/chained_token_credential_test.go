// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
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
	secCred := NewClientSecretCredential(tenantID, clientID, secret, nil)
	envCred, err := NewEnvironmentCredential(nil)
	if err != nil {
		t.Fatalf("Could not find appropriate environment credentials")
	}
	cred, err := NewChainedTokenCredential(secCred, envCred)
	if err != nil {
		t.Fatalf("Unable to instantiate ChainedTokenCredential")
	}
	if cred != nil {
		if len(cred.sources) != 2 {
			t.Fatalf("Expected 2 sources in the chained token credential, instead found %d", len(cred.sources))
		}
	}

}
func TestChainedTokenCredential_NilCredentialInChain(t *testing.T) {
	var unavailableError *CredentialUnavailableError
	cred := NewClientSecretCredential(tenantID, clientID, secret, nil)

	_, err := NewChainedTokenCredential(cred, nil)
	if err != nil {
		if !errors.As(err, &unavailableError) {
			t.Fatalf("Actual error: %v, Expected error: %v, wrong type %T", err, unavailableError, err)
		}
		if len(err.Error()) == 0 {
			t.Fatalf("Did not create an appropriate error message")
		}
	}
}

func TestChainedTokenCredential_NilChain(t *testing.T) {
	var unavailableError *CredentialUnavailableError
	_, err := NewChainedTokenCredential()
	if err != nil {
		if !errors.As(err, &unavailableError) {
			t.Fatalf("Actual error: %v, Expected error: %v, wrong type %T", err, unavailableError, err)
		}
	}
}

func TestChainedTokenCredential_GetTokenSuccess(t *testing.T) {
	err := initEnvironmentVarsForTest()
	if err != nil {
		t.Fatalf("Could not set environment variables for testing: %v", err)
	}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(`{"access_token": "new_token", "expires_in": 3600}`)))
	srvURL := srv.URL()
	secCred := NewClientSecretCredential(tenantID, clientID, secret, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	envCred, err := NewEnvironmentCredential(&TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Failed to create environment credential: %v", err)
	}
	cred, err := NewChainedTokenCredential(secCred, envCred)
	if err != nil {
		t.Fatalf("Failed to create ChainedTokenCredential: %v", err)
	}
	tk, err := cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Received an error when attempting to get a token but expected none")
	}
	if tk.Token != tokenValue {
		t.Fatalf("Received an incorrect access token")
	}
	if tk.ExpiresIn != tokenExpiresIn {
		t.Fatalf("Received an incorrect time in the response")
	}
}

func TestChainedTokenCredential_GetTokenFail(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized))
	testURL := srv.URL()
	secCred := NewClientSecretCredential(tenantID, clientID, wrongSecret, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &testURL})
	msiCred := NewManagedIdentityCredential("", nil)
	cred, err := NewChainedTokenCredential(msiCred, secCred)
	if err != nil {
		t.Fatalf("Failed to create ChainedTokenCredential: %v", err)
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
	var chainedError *ChainedCredentialError
	if !errors.As(err, &chainedError) {
		t.Fatalf("Expected Error Type: ChainedCredentialError, ReceivedErrorType: %T", err)
	}
	if len(err.Error()) == 0 {
		t.Fatalf("Did not create an appropriate error message")
	}
}
