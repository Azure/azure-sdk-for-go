// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// AzureCLICredential enables authentication to Azure Active Directory using Azure CLI to generated an access token.
type AzureCLICredential struct {
	client      *azureCLICredentialClient
	execManager execManager
}

// NewAzureCLICredential creates an instance of AzureCLICredential to authenticate against Azure Active Directory with Azure CLI Credential's token.
// execManager: This allows for unit testing, by mocking execManger.
// Azure CLI command "az account get-access-token -o json --resource" will be executed when it set as nil.
func NewAzureCLICredential(execManager execManager) *AzureCLICredential {
	var client = newAzureCLICredentialClient()
	if execManager == nil {
		return &AzureCLICredential{client: client, execManager: &execManagerStruct{}}
	}

	return &AzureCLICredential{client: client, execManager: execManager}
}

// GetToken obtains a token from Azure CLI for development scenarios.
// ctx: controlling the request lifetime.
// scopes: The list of scopes for which the token will have access.
// execManager: The struct implements the interface execManager.
// Returns an AccessToken which can be used to authenticate service Client calls.
func (c *AzureCLICredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	return c.client.authenticate(ctx, opts.Scopes, c.execManager)
}

// AuthenticationPolicy implements the azcore.Credential interface on AzureCLICredential.
func (c *AzureCLICredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}
