// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// AuthorizationCodeCredentialOptions contains optional parameters for AuthorizationCodeCredential.
type AuthorizationCodeCredentialOptions struct {
	azcore.ClientOptions

	// ClientSecret is one of the application's client secrets.
	ClientSecret string
	// AuthorityHost is the base URL of an Azure Active Directory authority. Defaults
	// to the value of environment variable AZURE_AUTHORITY_HOST, if set, or AzurePublicCloud.
	AuthorityHost AuthorityHost
}

// AuthorizationCodeCredential authenticates by redeeming an authorization code previously
// obtained from Azure Active Directory. The authorization code flow is described in more detail
// in Azure Active Directory documentation: https://docs.microsoft.com/azure/active-directory/develop/v2-oauth2-auth-code-flow
type AuthorizationCodeCredential struct {
	client       *aadIdentityClient
	tenantID     string
	clientID     string
	authCode     string
	clientSecret string
	redirectURI  string
}

// NewAuthorizationCodeCredential constructs an AuthorizationCodeCredential.
// tenantID: The application's Azure Active Directory tenant or directory ID.
// clientID: The application's client ID.
// authCode: The authorization code received from the authorization code flow. Note that authorization codes are single-use.
// redirectURL: The application's redirect URL. Must match the redirect URL used to request the authorization code.
// options: Optional configuration.
func NewAuthorizationCodeCredential(tenantID string, clientID string, authCode string, redirectURL string, options *AuthorizationCodeCredentialOptions) (*AuthorizationCodeCredential, error) {
	if !validTenantID(tenantID) {
		return nil, errors.New(tenantIDValidationErr)
	}
	cp := AuthorizationCodeCredentialOptions{}
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
	return &AuthorizationCodeCredential{tenantID: tenantID, clientID: clientID, authCode: authCode, clientSecret: cp.ClientSecret, redirectURI: redirectURL, client: c}, nil
}

// GetToken obtains a token from Azure Active Directory by redeeming the authorization code. This method is called automatically by Azure SDK clients.
// ctx: Context controlling the request lifetime.
// opts: Options for the token request, in particular the desired scope of the access token.
func (c *AuthorizationCodeCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	tk, err := c.client.authenticateAuthCode(ctx, c.tenantID, c.clientID, c.authCode, c.clientSecret, "", c.redirectURI, opts.Scopes)
	if err != nil {
		addGetTokenFailureLogs("Authorization Code Credential", err, true)
		return nil, err
	}
	logGetTokenSuccess(c, opts)
	return tk, nil
}

var _ azcore.TokenCredential = (*AuthorizationCodeCredential)(nil)
