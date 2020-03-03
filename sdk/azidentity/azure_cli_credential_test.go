// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	scopeResource = "https://storage.azure.com/.default"
	mockScope     = "https://default.mock.auth.scope/.default"
)

func TestCLICredential_GetTokenSuccessMock(t *testing.T) {
	cred := &AzureCLICredential{}
	client := reflect.ValueOf(cred).Elem().FieldByName("Client")

	mockClient := newAzureCLICredentialClient(&azureCLIAccessTokenProviderGetTokenMock{})
	client.Set(reflect.ValueOf(mockClient))

	accessToken, err := cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{mockScope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}
	if accessToken.Token != "mocktoken" {
		t.Fatalf("Expected token equals 'mocktoken' but received: %v", accessToken.Token)
	}
}

func TestCLICredential_CredentialUnavailableMock(t *testing.T) {
	cred := &AzureCLICredential{}
	client := reflect.ValueOf(cred).Elem().FieldByName("Client")

	mockClient := newAzureCLICredentialClient(&azureCLIAccessTokenProviderCredentialUnavailableMock{})
	client.Set(reflect.ValueOf(mockClient))

	_, err := cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{mockScope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}

	var authFailed *CredentialUnavailableError
	if !errors.As(err, &authFailed) {
		t.Fatalf("Expected: CredentialUnavailableError, Received: %T", err)
	}
}

// azureCLIAccessTokenProviderGetTokenMock mock func getAzureCLIAuthResults return access token.
type azureCLIAccessTokenProviderGetTokenMock struct {
}

func (c *azureCLIAccessTokenProviderGetTokenMock) getAzureCLIAuthResults(command string) (*authResults, error) {
	out := []byte(" {\"accessToken\":\"mocktoken\" , " +
		"\"expiresOn\": \"2007-01-01 01:01:01.079627\"," +
		"\"subscription\": \"mocksub\"," +
		"\"tenant\": \"mocktenant\"," +
		"\"tokenType\": \"mocktype\"}")

	results := &authResults{}
	results.out = out
	results.errOut = ""

	return results, nil
}

// azureCLIAccessTokenProviderCredentialUnavailableMock mock func getAzureCLIAuthResults return error.
type azureCLIAccessTokenProviderCredentialUnavailableMock struct {
}

func (c *azureCLIAccessTokenProviderCredentialUnavailableMock) getAzureCLIAuthResults(command string) (*authResults, error) {
	results := &authResults{}
	results.out = nil
	results.errOut = "some error"

	return results, errors.New("Errors")
}
