// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// NewEnvironmentCredential creates an instance of the ClientSecretCredential type and reads credential details from environment variables.
// If the expected environment variables are not found at this time, then a CredentialUnavailableError will be returned.
// options: The options used to configure the management of the requests sent to Azure Active Directory.
func NewEnvironmentCredential(options *TokenCredentialOptions) (*ClientSecretCredential, error) {
	tenantID := os.Getenv("AZURE_TENANT_ID")
	if tenantID == "" {
		err := &CredentialUnavailableError{CredentialType: "Environment Credential", Message: "Missing environment variable AZURE_TENANT_ID"}
		logCredentialError(err.CredentialType, err)
		return nil, err
	}

	clientID := os.Getenv("AZURE_CLIENT_ID")
	if clientID == "" {
		err := &CredentialUnavailableError{CredentialType: "Environment Credential", Message: "Missing environment variable AZURE_CLIENT_ID"}
		logCredentialError(err.CredentialType, err)
		return nil, err
	}

	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	if clientSecret == "" {
		err := &CredentialUnavailableError{CredentialType: "Environment Credential", Message: "Missing environment variable AZURE_CLIENT_SECRET"}
		logCredentialError(err.CredentialType, err)
		return nil, err
	}
	azcore.Log().Write(LogCredential, "Azure Identity => NewEnvironmentCredential() invoking ClientSecretCredential")
	return NewClientSecretCredential(tenantID, clientID, clientSecret, options)
}
