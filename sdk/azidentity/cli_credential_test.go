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

// Mock func getCliAccessToken return resulting
func (c *getTokenMock) getCliAccessToken(command string) ([]byte, string, error) {
	return []byte(" {\"accessToken\":\"mocktoken\" , " +
		"\"expiresOn\": \"2007-01-01 01:01:01.079627\"," +
		"\"subscription\": \"mocksub\"," +
		"\"tenant\": \"mocktenant\"," +
		"\"tokenType\": \"mocktype\"}"), "", nil
}

func TestCliCredential_GetTokenSuccessMock(t *testing.T) {
	var shellClientMock ShellClient
	shellClientMock = &getTokenMock{}

	var options *CliCredentialOption
	options = &CliCredentialOption{shellClientOption: shellClientMock}
	cred, err := NewCliCredential(options)
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}

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

func (c *azNotLoginMock) getCliAccessToken(command string) ([]byte, string, error) {
	return nil, "ERROR: Please run 'az login'", errors.New("mockError")
}

func TestCliCredential_AzNotLogin(t *testing.T) {
	var shellClientMock ShellClient
	shellClientMock = &azNotLoginMock{}

	var options *CliCredentialOption
	options = &CliCredentialOption{shellClientOption: shellClientMock}

	cred, err := NewCliCredential(options)
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}

	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{mockScope}})
	if err == nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}

	var authFailed *AuthenticationFailedError
	if !errors.As(err, &authFailed) {
		t.Fatalf("Expected: AuthenticationFailedError, Received: %T", err)
	}
}

type winAzureCliNotInstalledMock struct {
}

func (c *winAzureCliNotInstalledMock) getCliAccessToken(command string) ([]byte, string, error) {
	return nil, "'az' is not recognized", errors.New("mockError")
}

func TestCliCredential_WinAzureCLINotInstalled(t *testing.T) {
	var shellClientMock ShellClient
	shellClientMock = &winAzureCliNotInstalledMock{}

	var options *CliCredentialOption
	options = &CliCredentialOption{shellClientOption: shellClientMock}

	cred, err := NewCliCredential(options)
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}

	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{mockScope}})
	if err == nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}

	var authFailed *AuthenticationFailedError
	if !errors.As(err, &authFailed) {
		t.Fatalf("Expected: AuthenticationFailedError, Received: %T", err)
	}
}

type linuxAzureCliNotInstalledMock struct {
}

func (c *linuxAzureCliNotInstalledMock) getCliAccessToken(command string) ([]byte, string, error) {
	return nil, "az: command not found", errors.New("mockError")
}

func TestCliCredential_LinuxAzureCLINotInstalled(t *testing.T) {
	var shellClientMock ShellClient
	shellClientMock = &linuxAzureCliNotInstalledMock{}

	var options *CliCredentialOption
	options = &CliCredentialOption{shellClientOption: shellClientMock}

	cred, err := NewCliCredential(options)
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}

	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{mockScope}})
	if err == nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}

	var authFailed *AuthenticationFailedError
	if !errors.As(err, &authFailed) {
		t.Fatalf("Expected: AuthenticationFailedError, Received: %T", err)
	}
}

type macAzureCliNotInstalledMock struct {
}

func (c *macAzureCliNotInstalledMock) getCliAccessToken(command string) ([]byte, string, error) {
	return nil, "az: not found", errors.New("mockError")
}

func TestCliCredential_MacAzureCLINotInstalled(t *testing.T) {
	var shellClientMock ShellClient
	shellClientMock = &macAzureCliNotInstalledMock{}

	var options *CliCredentialOption
	options = &CliCredentialOption{shellClientOption: shellClientMock}

	cred, err := NewCliCredential(options)
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}

	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{mockScope}})
	if err == nil {
		t.Fatalf("Expected an empty error but received: %v", err)
	}

	var authFailed *AuthenticationFailedError
	if !errors.As(err, &authFailed) {
		t.Fatalf("Expected: AuthenticationFailedError, Received: %T", err)
	}
}
