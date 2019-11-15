// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ManagedIdentityCredential attempts authentication using a managed identity that has been assigned to the deployment environment. This authentication type works in Azure VMs,
// App Service and Azure Functions applications, as well as inside of Azure Cloud Shell. More information about configuring managed identities can be found here:
// https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview
type ManagedIdentityCredential struct {
	clientID string
	client   *managedIdentityClient
}

// NewManagedIdentityCredential creates an instance of the ManagedIdentityCredential capable of authenticating a resource with a managed identity.
// - clientID: The client id to authenticate for a user assigned managed identity.  More information on user assigned managed identities cam be found here:
// https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview#how-a-user-assigned-managed-identity-works-with-an-azure-vm
// - options: Options that allow to configure the management of the requests sent to the Azure Active Directory service.
func NewManagedIdentityCredential(clientID string, options *IdentityClientOptions) (*ManagedIdentityCredential, error) {
	client, err := newManagedIdentityClient(options)
	if err != nil {
		return nil, err
	}

	return &ManagedIdentityCredential{clientID: clientID, client: client}, nil
}

// GetToken obtains an AccessToken from the Managed Identity service if available.
// - scopes: The list of scopes for which the token will have access.
// Returns an AccessToken which can be used to authenticate service client calls, or a default AccessToken if no managed identity is available.
func (c *ManagedIdentityCredential) GetToken(ctx context.Context, scopes []string) (*azcore.AccessToken, error) {
	return c.client.authenticate(ctx, c.clientID, scopes)
}

// AuthenticationPolicy implements the azcore.Credential interface on ManagedIdentityCredential.
func (c *ManagedIdentityCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options.Scopes)
}
