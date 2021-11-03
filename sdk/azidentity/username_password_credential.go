// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// UsernamePasswordCredentialOptions can be used to provide additional information to configure the UsernamePasswordCredential.
// Use these options to modify the default pipeline behavior through the TokenCredentialcp.
// All zero-value fields will be initialized with their default values.
type UsernamePasswordCredentialOptions struct {
	azcore.ClientOptions

	// The host of the Azure Active Directory authority. The default is AzurePublicCloud.
	// Leave empty to allow overriding the value from the AZURE_AUTHORITY_HOST environment variable.
	AuthorityHost AuthorityHost
}

// UsernamePasswordCredential enables authentication to Azure Active Directory using a user's  username and password. If the user has MFA enabled this
// credential will fail to get a token returning an AuthenticationFailureError. Also, this credential requires a high degree of trust and is not
// recommended outside of prototyping when more secure credentials can be used.
type UsernamePasswordCredential struct {
	client   *aadIdentityClient
	tenantID string // Gets the Azure Active Directory tenant (directory) ID of the service principal
	clientID string // Gets the client (application) ID of the service principal
	username string // Gets the user account's user name
	password string // Gets the user account's password
}

// NewUsernamePasswordCredential constructs a new UsernamePasswordCredential with the details needed to authenticate against Azure Active Directory with
// a simple username and password.
// tenantID: The Azure Active Directory tenant (directory) ID of the service principal.
// clientID: The client (application) ID of the service principal.
// username: A user's account username
// password: A user's account password
// options: UsernamePasswordCredentialOptions used to configure the pipeline for the requests sent to Azure Active Directory.
func NewUsernamePasswordCredential(tenantID string, clientID string, username string, password string, options *UsernamePasswordCredentialOptions) (*UsernamePasswordCredential, error) {
	if !validTenantID(tenantID) {
		return nil, errors.New(tenantIDValidationErr)
	}
	cp := UsernamePasswordCredentialOptions{}
	if options != nil {
		cp = *options
	}
	authorityHost, err := setAuthorityHost(cp.AuthorityHost)
	if err != nil {
		return nil, err
	}
	c, err := newAADIdentityClient(authorityHost, &cp.ClientOptions)
	if err != nil {
		return nil, err
	}
	return &UsernamePasswordCredential{tenantID: tenantID, clientID: clientID, username: username, password: password, client: c}, nil
}

// GetToken obtains a token from Azure Active Directory using the specified username and password.
// scopes: The list of scopes for which the token will have access.
// ctx: The context used to control the request lifetime.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *UsernamePasswordCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	tk, err := c.client.authenticateUsernamePassword(ctx, c.tenantID, c.clientID, c.username, c.password, opts.Scopes)
	if err != nil {
		addGetTokenFailureLogs("Username Password Credential", err, true)
		return nil, err
	}
	logGetTokenSuccess(c, opts)
	return tk, err
}

var _ azcore.TokenCredential = (*UsernamePasswordCredential)(nil)
