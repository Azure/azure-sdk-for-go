// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

func TestDefaultAzureCredential_GetTokenSuccess(t *testing.T) {
	env := map[string]string{"AZURE_TENANT_ID": fakeTenantID, "AZURE_CLIENT_ID": fakeClientID, "AZURE_CLIENT_SECRET": secret}
	setEnvironmentVariables(t, env)
	cred, err := NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("Unable to create credential. Received: %v", err)
	}
	c := cred.chain.sources[0].(*EnvironmentCredential)
	c.cred.(*ClientSecretCredential).client = fakeConfidentialClient{}
	_, err = cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{"scope"}})
	if err != nil {
		t.Fatalf("GetToken error: %v", err)
	}
}

func TestDefaultAzureCredential_defaultAzureCredentialConstructorErrorHandlerNoSuccessfulCredential(t *testing.T) {
	err := os.Setenv("AZURE_SDK_GO_LOGGING", "all")
	if err != nil {
		t.Fatal("Unexpected error", err.Error())
	}

	logMessages := []string{}
	log.SetListener(func(event log.Event, message string) {
		logMessages = append(logMessages, message)
	})

	errorMessages := []string{
		"<credential-name>: <error-message>",
		"<credential-name>: <error-message>",
	}
	err = defaultAzureCredentialConstructorErrorHandler(0, errorMessages)
	if err == nil {
		t.Fatalf("Expected an error, but received none.")
	}
	expectedError := `<credential-name>: <error-message>
	<credential-name>: <error-message>`
	if err.Error() != expectedError {
		t.Fatalf("Did not create an appropriate error message.\n\nReceived:\n%s\n\nExpected:\n%s", err.Error(), expectedError)
	}

	expectedLogs := `Azure Identity => Failed to initialize the Default Azure Credential:
	<credential-name>: <error-message>
	<credential-name>: <error-message>`
	if logMessages[0] != expectedLogs {
		t.Fatalf("Did not receive the expected logs.\n\nReceived:\n%s\n\nExpected:\n%s", logMessages[0], expectedLogs)
	}
}

func TestDefaultAzureCredential_defaultAzureCredentialConstructorErrorHandlerOneSuccessfulCredential(t *testing.T) {
	err := os.Setenv("AZURE_SDK_GO_LOGGING", "all")
	if err != nil {
		t.Fatal("Unexpected error", err.Error())
	}

	logMessages := []string{}
	log.SetListener(func(event log.Event, message string) {
		logMessages = append(logMessages, message)
	})

	errorMessages := []string{
		"<credential-name>: <error-message>",
		"<credential-name>: <error-message>",
	}
	err = defaultAzureCredentialConstructorErrorHandler(1, errorMessages)
	if err != nil {
		t.Fatal("Unexpected error", err.Error())
	}

	expectedLogs := `Azure Identity => Failed to initialize some credentials on the Default Azure Credential:
	<credential-name>: <error-message>
	<credential-name>: <error-message>`
	if logMessages[0] != expectedLogs {
		t.Fatalf("Did not receive the expected logs.\n\nReceived:\n%s\n\nExpected:\n%s", logMessages[0], expectedLogs)
	}
}
