// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// AuthorizationCodeCredentialOptions contain optional parameters that can be used to configure the AuthorizationCodeCredential.
// All zero-value fields will be initialized with their default values.
type AuthorizationCodeCredentialOptions struct {
	azcore.ClientOptions

	// Gets the client secret that was generated for the App Registration used to authenticate the client.
	ClientSecret string
	// The host of the Azure Active Directory authority. The default is AzurePublicCloud.
	// Leave empty to allow overriding the value from the AZURE_AUTHORITY_HOST environment variable.
	AuthorityHost AuthorityHost
}

// AuthorizationCodeCredential enables authentication to Azure Active Directory using an authorization code
// that was obtained through the authorization code flow, described in more detail in the Azure Active Directory
// documentation: https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow.
type AuthorizationCodeCredential struct {
	client       *aadIdentityClient
	tenantID     string // Gets the Azure Active Directory tenant (directory) ID of the service principal
	clientID     string // Gets the client (application) ID of the service principal
	authCode     string // The authorization code received from the authorization code flow. The authorization code must not have been used to obtain another token.
	clientSecret string // Gets the client secret that was generated for the App Registration used to authenticate the client.
	redirectURI  string // The redirect URI that was used to request the authorization code. Must be the same URI that is configured for the App Registration.
}

// NewAuthorizationCodeCredential constructs a new AuthorizationCodeCredential with the details needed to authenticate against Azure Active Directory with an authorization code.
// tenantID: The Azure Active Directory tenant (directory) ID of the service principal.
// clientID: The client (application) ID of the service principal.
// authCode: The authorization code received from the authorization code flow. The authorization code must not have been used to obtain another token.
// redirectURL: The redirect URL that was used to request the authorization code. Must be the same URL that is configured for the App Registration.
// options: Manage the configuration of the requests sent to Azure Active Directory, they can also include a client secret for web app authentication.
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

// GetToken obtains a token from Azure Active Directory, using the specified authorization code to authenticate.
// ctx: Context used to control the request lifetime.
// opts: TokenRequestOptions contains the list of scopes for which the token will have access.
// Returns an AccessToken which can be used to authenticate service client calls.
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
