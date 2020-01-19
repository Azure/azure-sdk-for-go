// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// Enables authentication to Azure Active Directory using Azure cli to generated an access token.
type AzureCliCredential struct {
	client      *AzureCliCredentialClient
	shellClient ShellClient
}

// AzureCliCredentialOption contains parameter that Shell execute Azure Cli command to return a json output what include accessToken and expiresOn.
type AzureCliCredentialOption struct {
	shellClientOption ShellClient
}

// NewAzureCliCredential creates an instance of AzureCliCredential to authenticate against Azure Active Directory with Azure Cli Credential's token.
func NewAzureCliCredential(options *AzureCliCredentialOption) (*AzureCliCredential, error) {
	client := NewAzureCliCredentialClient()
	if options != nil {
		return &AzureCliCredential{client: client, shellClient: options.shellClientOption}, nil
	} else {
		var shellClient ShellClient
		shellClient = client

		return &AzureCliCredential{client: client, shellClient: shellClient}, nil
	}
}

// GetToken obtains a token from Azure Cli, using Azure CLI to generated an access token to authenticate.
// scopes: The list of scopes for which the token will have access.
// ctx: controlling the request lifetime.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *AzureCliCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	return c.client.authenticate(ctx, opts.Scopes, c.shellClient)
}

// AuthenticationPolicy implements the azcore.Credential interface on AzureCliCredential.
func (c *AzureCliCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}
