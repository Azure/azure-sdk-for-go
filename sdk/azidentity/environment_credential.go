// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// EnvironmentCredential type enables authentication to Azure Active Directory using client secret
// details configured in the following environment variables:
// - AZURE_TENANT_ID: The Azure Active Directory tenant(directory) ID.
// - AZURE_CLIENT_ID: The client(application) ID of an App Registration in the tenant.
// - AZURE_CLIENT_SECRET: A client secret that was generated for the App Registration.
// This credential ultimately uses a ClientSecretCredential to
// perform the authentication using these details. Please consult the
// documentation of that class for more details.
type EnvironmentCredential struct {
	credential azcore.TokenCredential
}

// NewEnvironmentCredential creates an instance of the EnvironmentCredential type and reads client secret details from environment variables.
// If the expected environment variables are not found at this time, the GetToken method will return the default AccessToken when invoked.
// - options: The options used to configure the management of the requests sent to the Azure Active Directory service.
func NewEnvironmentCredential(options *IdentityClientOptions) (*EnvironmentCredential, error) {
	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")

	if tenantID == "" {
		return nil, &CredentialUnavailableError{Message: "Missing environment variable: AZURE_TENANT_ID"}
	}

	if clientID == "" {
		return nil, &CredentialUnavailableError{Message: "Missing environment variable: AZURE_CLIENT_ID"}
	}

	if clientSecret == "" {
		return nil, &CredentialUnavailableError{Message: "Missing environment variable: AZURE_CLIENT_SECRET"}
	}

	credential, err := NewClientSecretCredential(tenantID, clientID, clientSecret, options)
	if err != nil {
		return nil, err
	}
	return &EnvironmentCredential{credential: credential}, nil
}

// GetToken obtains a token from the Azure Active Directory service, using the specified client details specified in the environment variables
// AZURE_TENANT_ID, AZURE_CLIENT_ID, and AZURE_CLIENT_SECRET to authenticate.
// scopes: The list of scopes for which the token will have access.
// ctx: controlling the request lifetime.
// If the environment variables AZURE_TENANT_ID, AZURE_CLIENT_ID, and AZURE_CLIENT_SECRET are not specified, the default AccessToken is returned
func (c EnvironmentCredential) GetToken(ctx context.Context, scopes []string) (*azcore.AccessToken, error) {
	return c.credential.GetToken(ctx, scopes)
}

// TODO: credential methods by ref
