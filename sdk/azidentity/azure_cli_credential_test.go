// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

var (
	mockCLITokenProviderSuccess = func(ctx context.Context, resource string) ([]byte, error) {
		return []byte(" {\"accessToken\":\"mocktoken\" , " +
			"\"expiresOn\": \"2007-01-01 01:01:01.079627\"," +
			"\"subscription\": \"mocksub\"," +
			"\"tenant\": \"mocktenant\"," +
			"\"tokenType\": \"mocktype\"}"), nil
	}
	mockCLITokenProviderFailure = func(ctx context.Context, resource string) ([]byte, error) {
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
	at, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
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
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one.")
	}
}

func TestBearerPolicy_AzureCLICredential(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	options := AzureCLICredentialOptions{}
	options.tokenProvider = mockCLITokenProviderSuccess
	cred, err := NewAzureCLICredential(&options)
	if err != nil {
		t.Fatalf("Did not expect an error but received: %v", err)
	}
	pipeline := defaultTestPipeline(srv, cred, scope)
	req, err := runtime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}
	_, err = pipeline.Do(req)
	if err != nil {
		t.Fatal("Expected nil error but received one")
	}
}
