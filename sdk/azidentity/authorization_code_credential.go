// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// AuthorizationCodeCredentialOptions contain optional parameters that can be used to configure the AuthorizationCodeCredential.
type AuthorizationCodeCredentialOptions struct {
	// Gets the client secret that was generated for the App Registration used to authenticate the client.
	ClientSecret *string
	// Manage the configuration of the requests sent to Azure Active Directory.
	Options *TokenCredentialOptions
}

// AuthorizationCodeCredential enables authentication to Azure Active Directory using an authorization code
// that was obtained through the authorization code flow, described in more detail in the Azure Active Directory
// documentation: https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow.
type AuthorizationCodeCredential struct {
	client       *aadIdentityClient
	tenantID     string  // Gets the Azure Active Directory tenant (directory) ID of the service principal
	clientID     string  // Gets the client (application) ID of the service principal
	authCode     string  // The authorization code received from the authorization code flow. The authorization code must not have been used to obtain another token.
	clientSecret *string // Gets the client secret that was generated for the App Registration used to authenticate the client.
	redirectURI  string  // The redirect URI that was used to request the authorization code. Must be the same URI that is configured for the App Registration.
}

// DefaultAuthorizationCodeCredentialOptions returns an instance of AuthorizationCodeCredentialOptions initialized with default values.
func DefaultAuthorizationCodeCredentialOptions() AuthorizationCodeCredentialOptions {
	return AuthorizationCodeCredentialOptions{}
}

// NewAuthorizationCodeCredential constructs a new AuthorizationCodeCredential with the details needed to authenticate against Azure Active Directory with an authorization code.
// tenantID: The Azure Active Directory tenant (directory) ID of the service principal.
// clientID: The client (application) ID of the service principal.
// authCode: The authorization code received from the authorization code flow. The authorization code must not have been used to obtain another token.
// redirectURI: The redirect URI that was used to request the authorization code. Must be the same URI that is configured for the App Registration.
// options: Manage the configuration of the requests sent to Azure Active Directory, they can also include a client secret for web app authentication.
func NewAuthorizationCodeCredential(tenantID string, clientID string, authCode string, redirectURI string, options *AuthorizationCodeCredentialOptions) (*AuthorizationCodeCredential, error) {
	if !validTenantID(tenantID) {
		return nil, &CredentialUnavailableError{CredentialType: "Authorization Code Credential", Message: "invalid tenant ID passed to credential"}
	}
	if options == nil {
		temp := DefaultAuthorizationCodeCredentialOptions()
		options = &temp
	}
	c, err := newAADIdentityClient(options.Options)
	if err != nil {
		return nil, err
	}
	return &AuthorizationCodeCredential{tenantID: tenantID, clientID: clientID, authCode: authCode, clientSecret: options.ClientSecret, redirectURI: redirectURI, client: c}, nil
}

// GetToken obtains a token from Azure Active Directory, using the specified authorization code to authenticate.
// ctx: Context used to control the request lifetime.
// opts: TokenRequestOptions contains the list of scopes for which the token will have access.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *AuthorizationCodeCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	tk, err := c.client.authenticateAuthCode(ctx, c.tenantID, c.clientID, c.authCode, c.clientSecret, c.redirectURI, opts.Scopes)
	if err != nil {
		addGetTokenFailureLogs("Authorization Code Credential", err, true)
		return nil, err
	}
	logGetTokenSuccess(c, opts)
	return tk, nil
}

// AuthenticationPolicy implements the azcore.Credential interface on AuthorizationCodeCredential and calls the Bearer Token policy
// to get the bearer token.
func (c *AuthorizationCodeCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}

var _ azcore.TokenCredential = (*AuthorizationCodeCredential)(nil)
