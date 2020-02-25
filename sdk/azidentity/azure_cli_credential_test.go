// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	scopeResource = "https://storage.azure.com/.default"
	mockScope     = "https://default.mock.auth.scope/.default"
)

type getTokenMock struct {
}

// Mock func getAzureCLIAccessToken return resulting
func (c *getTokenMock) getAzureCLIAccessToken(command string) ([]byte, string, error) {
	return []byte(" {\"accessToken\":\"mocktoken\" , " +
		"\"expiresOn\": \"2007-01-01 01:01:01.079627\"," +
		"\"subscription\": \"mocksub\"," +
		"\"tenant\": \"mocktenant\"," +
		"\"tokenType\": \"mocktype\"}"), "", nil
}

func TestCLICredential_GetTokenSuccessMock(t *testing.T) {
	var shellClientMock shellClient
	shellClientMock = &getTokenMock{}

	var options *AzureCLICredentialOptions
	options = &AzureCLICredentialOptions{shellClientOption: shellClientMock}
	cred := NewAzureCLICredential(options)

	accessToken, err := cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{mockScope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}
	if accessToken.Token != "mocktoken" {
		t.Fatalf("Expected token equals 'mocktoken' but received: %v", accessToken.Token)
	}
}

type azNotLoginMock struct {
}

func (c *azNotLoginMock) getAzureCLIAccessToken(command string) ([]byte, string, error) {
	return nil, "ERROR: Please run 'az login'", errors.New("mockError")
}

func TestCLICredential_AzNotLogin(t *testing.T) {
	var err error
	var shellClientMock shellClient
	shellClientMock = &azNotLoginMock{}

	var options *AzureCLICredentialOptions
	options = &AzureCLICredentialOptions{shellClientOption: shellClientMock}

	cred := NewAzureCLICredential(options)

	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{mockScope}})
	if err == nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}

	var authFailed *AuthenticationFailedError
	if !errors.As(err, &authFailed) {
		t.Fatalf("Expected: AuthenticationFailedError, Received: %T", err)
	}
}

type winAzureCLINotInstalledMock struct {
}

func (c *winAzureCLINotInstalledMock) getAzureCLIAccessToken(command string) ([]byte, string, error) {
	return nil, "'az' is not recognized", errors.New("mockError")
}

func TestCLICredential_WinAzureCLINotInstalled(t *testing.T) {
	var err error
	var shellClientMock shellClient
	shellClientMock = &winAzureCLINotInstalledMock{}

	var options *AzureCLICredentialOptions
	options = &AzureCLICredentialOptions{shellClientOption: shellClientMock}

	cred := NewAzureCLICredential(options)

	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{mockScope}})
	if err == nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}

	var authFailed *AuthenticationFailedError
	if !errors.As(err, &authFailed) {
		t.Fatalf("Expected: AuthenticationFailedError, Received: %T", err)
	}
}

type linuxAzureCLINotInstalledMock struct {
}

func (c *linuxAzureCLINotInstalledMock) getAzureCLIAccessToken(command string) ([]byte, string, error) {
	return nil, "az: command not found", errors.New("mockError")
}

func TestCLICredential_LinuxAzureCLINotInstalled(t *testing.T) {
	var err error
	var shellClientMock shellClient
	shellClientMock = &linuxAzureCLINotInstalledMock{}

	var options *AzureCLICredentialOptions
	options = &AzureCLICredentialOptions{shellClientOption: shellClientMock}

	cred := NewAzureCLICredential(options)

	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{mockScope}})
	if err == nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}

	var authFailed *AuthenticationFailedError
	if !errors.As(err, &authFailed) {
		t.Fatalf("Expected: AuthenticationFailedError, Received: %T", err)
	}
}

type macAzureCLINotInstalledMock struct {
}

func (c *macAzureCLINotInstalledMock) getAzureCLIAccessToken(command string) ([]byte, string, error) {
	return nil, "az: not found", errors.New("mockError")
}

func TestCLICredential_MacAzureCLINotInstalled(t *testing.T) {
	var err error
	var shellClientMock shellClient
	shellClientMock = &macAzureCLINotInstalledMock{}

	var options *AzureCLICredentialOptions
	options = &AzureCLICredentialOptions{shellClientOption: shellClientMock}

	cred := NewAzureCLICredential(options)

	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{mockScope}})
	if err == nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}

	var authFailed *AuthenticationFailedError
	if !errors.As(err, &authFailed) {
		t.Fatalf("Expected: AuthenticationFailedError, Received: %T", err)
	}
}
