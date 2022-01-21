// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
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

func TestDefaultAzureCredential_defaultAzureCredentialLogOrThrowError(t *testing.T) {
	errorMessages := []string{
		"<credential-name>: <error-message>",
		"<credential-name>: <error-message>",
	}
	err := defaultAzureCredentialLogOrThrowError(0, errorMessages)
	if err == nil {
		t.Fatalf("Expected an error, but received none.")
	}
	expectedError := `Default Azure Credential:
	<credential-name>: <error-message>
	<credential-name>: <error-message>`
	if err.Error() != expectedError {
		t.Fatalf("Did not create an appropriate error message.\n\nReceived:\n%s\n\nExpected:\n%s", err.Error(), expectedError)
	}
}
