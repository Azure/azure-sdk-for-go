// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// Enables authentication to Azure Active Directory using Azure CLI to generated an access token.
type CliCredential struct {
	client      *CliCredentialClient
	shellClient ShellClient
}

// CliCredentialOption contains parameter that Shell execute Azure Cli command to return a json output what include the accessToken and experisOn
type CliCredentialOption struct {
	shellClientOption ShellClient
}

// NewCliCredential creates an instance of CliCredential to authenticate against Azure Active Directory with Azure Cli Credential's token.
func NewCliCredential(options *CliCredentialOption) (*CliCredential, error) {
	client := NewCliCredentialClient()
	if options != nil {
		return &CliCredential{client: client, shellClient: options.shellClientOption}, nil
	} else {
		var shellClient ShellClient
		shellClient = client

		return &CliCredential{client: client, shellClient: shellClient}, nil
	}
}

// GetToken obtains a token from Azure Cli service, using Azure CLI to generated an access token to authenticate.
// scopes: The list of scopes for which the token will have access.
// ctx: controlling the request lifetime.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *CliCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	return c.client.authenticate(ctx, opts.Scopes, c.shellClient)
}

// AuthenticationPolicy implements the azcore.Credential interface on CliCredential.
func (c *CliCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}
