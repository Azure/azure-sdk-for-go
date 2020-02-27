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
func NewAzureCLICredential(execManager execManager) *AzureCLICredential {
	var client = newAzureCLICredentialClient()

	if execManager == nil {
		execManager = &execManage{}
		return &AzureCLICredential{client: client, execManager: execManager}
	}

	return &AzureCLICredential{client: client, execManager: execManager}
}

// GetToken obtains a token from Azure CLI, using Azure CLI to generated an access token to authenticate.
// scopes: The list of scopes for which the token will have access.
// ctx: controlling the request lifetime.
// Returns an AccessToken which can be used to authenticate service Client calls.
func (c *AzureCLICredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	return c.client.authenticate(ctx, opts.Scopes, c.execManager)
}

// AuthenticationPolicy implements the azcore.Credential interface on AzureCLICredential.
func (c *AzureCLICredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}
