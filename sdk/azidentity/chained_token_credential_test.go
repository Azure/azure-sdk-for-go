// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func initChainedTokenCredentialTest() error {
	err := os.Setenv("AZURE_TENANT_ID", tenantID)
	if err != nil {
		return err
	}
	err = os.Setenv("AZURE_CLIENT_ID", clientID)
	if err != nil {
		return err
	}
	err = os.Setenv("AZURE_CLIENT_SECRET", secret)
	if err != nil {
		return err
	}
	return nil
}
func TestChainedTokenCredentialSuccess(t *testing.T) {
	err := initChainedTokenCredentialTest()
	if err != nil {
		t.Fatalf("Could not set environment variables for testing: %v", err)
	}
	secCred := NewClientSecretCredential("expected_tenant", "client", "secret", nil)
	envCred, err := NewEnvironmentCredential(nil)
	if err != nil {
		t.Fatalf("Could not find the appropriate environment credentials")
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
func TestNilCredentialInChain(t *testing.T) {
	var unavailableError *CredentialUnavailableError
	cred := NewClientSecretCredential("expected_tenant", "client", "secret", nil)

	_, err := NewChainedTokenCredential(cred, nil, cred)
	if err != nil {
		switch i := err.(type) {
		case *CredentialUnavailableError:
		default:
			t.Fatalf("Actual error: %v, Expected error: %v, wrong type %t", err, unavailableError, i)
		}
	}
}

func TestNilChain(t *testing.T) {
	var unavailableError *CredentialUnavailableError

	_, err := NewChainedTokenCredential()
	if err != nil {
		switch i := err.(type) {
		case *CredentialUnavailableError:
			fmt.Println("Received: ", err.Error())
		default:
			t.Errorf("Actual error: %v, Expected error: %v, wrong type %t", err, unavailableError, i)
		}
	}
}

func Test_ChainedGetToken_Success(t *testing.T) {
	err := initChainedTokenCredentialTest()
	if err != nil {
		t.Fatalf("Could not set environment variables for testing: %v", err)
	}
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(`{"access_token": "new_token", "expires_in": 3600}`)))
	srvURL := srv.URL()
	secCred := NewClientSecretCredential("expected_tenant", "client", "secret", &TokenCredentialOptions{PipelineOptions: azcore.PipelineOptions{HTTPClient: srv}, AuthorityHost: &srvURL})
	envCred, err := NewEnvironmentCredential(&TokenCredentialOptions{PipelineOptions: azcore.PipelineOptions{HTTPClient: srv}, AuthorityHost: &srvURL})
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
	if tk.Token != "new_token" {
		t.Fatalf("Received an incorrect access token")
	}
	if tk.ExpiresIn != "3600" {
		t.Fatalf("Received an incorrect time in the response")
	}
}
