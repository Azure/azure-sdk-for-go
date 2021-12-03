// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// UsernamePasswordCredentialOptions contains optional parameters for UsernamePasswordCredential.
type UsernamePasswordCredentialOptions struct {
	azcore.ClientOptions

	// AuthorityHost is the base URL of an Azure Active Directory authority. Defaults
	// to the value of environment variable AZURE_AUTHORITY_HOST, if set, or AzurePublicCloud.
	AuthorityHost AuthorityHost
}

// UsernamePasswordCredential authenticates user with a password. Microsoft doesn't recommend this kind of authentication,
// because it's less secure than other authentication flows. This credential is not interactive, so it isn't compatible
// with any form of multi-factor authentication, and the application must already have user or admin consent.
// This credential can only authenticate work and school accounts; it can't authenticate Microsoft accounts.
type UsernamePasswordCredential struct {
	client   *aadIdentityClient
	tenantID string
	clientID string
	username string
	password string
}

// NewUsernamePasswordCredential creates a UsernamePasswordCredential.
// tenantID: The ID of the Azure Active Directory tenant the credential authenticates in.
// clientID: The ID of the application users will authenticate to.
// username: A username (typically an email address).
// password: That user's password.
// options: Optional configuration.
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

// GetToken obtains a token from Azure Active Directory. This method is called automatically by Azure SDK clients.
// ctx: Context used to control the request lifetime.
// opts: Options for the token request, in particular the desired scope of the access token.
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
