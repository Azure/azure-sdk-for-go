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
	shellClient shellClient
}

// AzureCLICredentialOptions contains the parameters for the shellClient to execute the Azure CLI command that will return an Access Token.
type AzureCLICredentialOptions struct {
	shellClientOption shellClient
}

// NewAzureCLICredential creates an instance of AzureCLICredential to authenticate against Azure Active Directory with Azure CLI Credential's token.
func NewAzureCLICredential(options *AzureCLICredentialOptions) *AzureCLICredential {
	var client = newAzureCLICredentialClient()
	if options == nil {
		var shellClient shellClient
		shellClient = client

		return &AzureCLICredential{client: client, shellClient: shellClient}
	}

	return &AzureCLICredential{client: client, shellClient: options.shellClientOption}
}

// GetToken obtains a token from Azure CLI, using Azure CLI to generated an access token to authenticate.
// scopes: The list of scopes for which the token will have access.
// ctx: controlling the request lifetime.
// Returns an AccessToken which can be used to authenticate service Client calls.
func (c *AzureCLICredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	return c.client.authenticate(ctx, opts.Scopes, c.shellClient)
}

// AuthenticationPolicy implements the azcore.Credential interface on AzureCLICredential.
func (c *AzureCLICredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}
