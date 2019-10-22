// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import "context"

// ManagedIdentityCredential Attempts authentication using a managed identity that has been assigned to the deployment environment.This authentication type works in Azure VMs,
// App Service and Azure Functions applications, as well as inside of Azure Cloud Shell. More information about configuring managed identities can be found here:
// https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview
type ManagedIdentityCredential struct {
	TokenCredential
	clientID string
	client   ManagedIdentityClient
}

// NewManagedIdentityCredential creates an instance of the ManagedIdentityCredential capable of authenticating a resource with a managed identity.
// The client id to authenticate for a user assigned managed identity.  More information on user assigned managed identities cam be found here:
// https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview#how-a-user-assigned-managed-identity-works-with-an-azure-vm
// Options that allow to configure the management of the requests sent to the Azure Active Directory service.
func NewManagedIdentityCredential(clientID string, options *IdentityClientOptions) (ManagedIdentityCredential, error) {
	if options == nil {
		options, _ = newIdentityClientOptions()
	}
	client, err := NewManagedIdentityClient(options)
	if err != nil {
		return ManagedIdentityCredential{}, err
	}

	return ManagedIdentityCredential{clientID: clientID, client: *client}, nil
}

// GetToken obtains an AccessToken from the Managed Identity service if available.
// Scopes: The list of scopes for which the token will have access.
// Returns an AccessToken which can be used to authenticate service client calls, or a default AccessToken if no managed identity is available.
// TODO: check params this func is expecting and returned token
func (c ManagedIdentityCredential) GetToken(ctx context.Context, scopes []string) (*AccessToken, error) {
	return c.client.Authenticate(ctx, c.clientID, scopes)
}
