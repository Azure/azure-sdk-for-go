// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"os"
)

// NewEnvironmentCredential creates an instance of the EnvironmentCredential type and reads client secret details from environment variables.
// If the expected environment variables are not found at this time, the GetToken method will return the default AccessToken when invoked.
// options: The options used to configure the management of the requests sent to the Azure Active Directory service.
func NewEnvironmentCredential(options *TokenCredentialOptions) (*ClientSecretCredential, error) {
	tenantID := os.Getenv("AZURE_TENANT_ID")
	if tenantID == "" {
		return nil, &CredentialUnavailableError{CredentialType: "Environment Credential", Message: "Missing environment variable AZURE_TENANT_ID"}
	}

	clientID := os.Getenv("AZURE_CLIENT_ID")
	if clientID == "" {
		return nil, &CredentialUnavailableError{CredentialType: "Environment Credential", Message: "Missing environment variable AZURE_CLIENT_ID"}
	}

	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	if clientSecret == "" {
		return nil, &CredentialUnavailableError{CredentialType: "Environment Credential", Message: "Missing environment variable AZURE_CLIENT_SECRET"}
	}

	return NewClientSecretCredential(tenantID, clientID, clientSecret, options), nil
}
