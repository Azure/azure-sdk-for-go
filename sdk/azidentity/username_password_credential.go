// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

// ------------------------ NOTES:
// CP: To test this I need to create my own tenant in my own personal azure account create a user with auth permissions and an app registration.

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// UsernamePasswordCredential enables authentication to Azure Active Directory using a user's  username and password. If the user has MFA enabled this
// credential will fail to get a token returning an AuthenticationFailureError. Also, this credential requires a high degree of trust and is not
// recommended outside of prototyping when more secure credentials can be used.
type UsernamePasswordCredential struct {
	azcore.TokenCredential

	client *aadIdentityClient

	TenantID string // Gets the Azure Active Directory tenant (directory) Id of the service principal

	ClientID string // Gets the client (application) ID of the service principal

	Username string // Gets the user account's user name

	Password string // Gets the user account's password
}

// NewUsernamePasswordCredential constructs a new UsernamePasswordCredential with the details needed to authenticate against Azure Active Directory with
// a simple username and password.
// - tenantID: The Azure Active Directory tenant (directory) Id of the service principal.
// - clientID: The client (application) ID of the service principal.
// - username: A user's account username
// - password: A user's account password
// - options: The options configure the management of the requests sent to the Azure Active Directory service.
func NewUsernamePasswordCredential(tenantID string, clientID string, username string, password string, options *IdentityClientOptions) (*UsernamePasswordCredential, error) {
	return &UsernamePasswordCredential{TenantID: tenantID, ClientID: clientID, Username: username, Password: password, client: newAADIdentityClient(options)}, nil
}

// GetToken obtains a token from the Azure Active Directory service, using the specified username and password.
// - scopes: The list of scopes for which the token will have access.
// - ctx: controlling the request lifetime.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c UsernamePasswordCredential) GetToken(ctx context.Context, scopes []string) (*azcore.AccessToken, error) {
	return c.client.authenticateUsernamePassword(ctx, c.TenantID, c.ClientID, c.Username, c.Password, scopes)
}
