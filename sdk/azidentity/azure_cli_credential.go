// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// AzureCLICredential enables authentication to Azure Active Directory using Azure CLI to generated an access token.
type AzureCLICredential struct {
	client *azureCLICredentialClient
}

// NewAzureCLICredential creates an instance of AzureCLICredential to authenticate against Azure Active Directory with Azure CLI Credential's token.
func NewAzureCLICredential() *AzureCLICredential {
	var client = newAzureCLICredentialClient(nil)
	return &AzureCLICredential{client: client}
}

// GetToken obtains a token from Azure CLI for development scenarios.
// ctx: controlling the request lifetime.
// scopes: The list of scopes for which the token will have access.
// Returns an AccessToken which can be used to authenticate service Client calls.
func (c *AzureCLICredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	return c.client.authenticate(ctx, opts.Scopes, c.client.azAccessTokenProvider)
}

// AuthenticationPolicy implements the azcore.Credential interface on AzureCLICredential.
func (c *AzureCLICredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}
