// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
)

// ClientSecretCredential enables authentication to Azure Active Directory using a client secret that was generated for an App Registration.  More information on how
// to configure a client secret can be found here:
// https://docs.microsoft.com/en-us/azure/active-directory/develop/quickstart-configure-app-access-web-apis#add-credentials-to-your-web-application
type ClientSecretCredential struct {
	TokenCredential

	client aadIdentityClient
	// Gets the Azure Active Directory tenant (directory) Id of the service principal
	TenantID string

	// Gets the client (application) ID of the service principal
	ClientID string

	// Gets the client secret that was generated for the App Registration used to authenticate the client.
	ClientSecret string
}

// NewClientSecretCredential constructs a new ClientSecretCredential with the details needed to authenticate against Azure Active Directory with a client secret.
// tenantID: The Azure Active Directory tenant (directory) Id of the service principal.
// clientID: The client (application) ID of the service principal.
// clientSecret: A client secret that was generated for the App Registration used to authenticate the client.
// options: allow to configure the management of the requests sent to the Azure Active Directory service.
func NewClientSecretCredential(tenantID string, clientID string, clientSecret string, options *IdentityClientOptions) ClientSecretCredential {
	if options == nil {
		options, _ = newIdentityClientOptions()
	}
	var client aadIdentityClient = newAADIdentityClient(options)

	return ClientSecretCredential{TenantID: tenantID, ClientID: clientID, ClientSecret: clientSecret, client: client}
}

// GetToken obtains a token from the Azure Active Directory service, using the specified client secret to authenticate.
// scopes: The list of scopes for which the token will have access.
// ctx: controlling the request lifetime.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c ClientSecretCredential) GetToken(ctx context.Context, scopes []string) (*AccessToken, error) {
	return c.client.Authenticate(ctx, c.TenantID, c.ClientID, c.ClientSecret, scopes)
}
