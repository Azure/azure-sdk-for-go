// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ClientCertificateCredential enables authentication of a service principal to Azure Active Directory using a certificate that is assigned to its App Registration. More information
// on how to configure certificate authentication can be found here:
// https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-certificate-credentials#register-your-certificate-with-azure-ad
type ClientCertificateCredential struct {
	client            *aadIdentityClient
	tenantID          string // The Azure Active Directory tenant (directory) ID of the service principal
	clientID          string // The client (application) ID of the service principal
	clientCertificate string // Path to the client certificate generated for the App Registration used to authenticate the client
}

// NewClientCertificateCredential creates an instance of ClientCertificateCredential with the details needed to authenticate against Azure Active Directory with the specified certificate.
// tenantID: The Azure Active Directory tenant (directory) ID of the service principal.
// clientID: The client (application) ID of the service principal.
// clientCertificate: The path to the client certificate that was generated for the App Registration used to authenticate the client.
// options: configure the management of the requests sent to Azure Active Directory.
func NewClientCertificateCredential(tenantID string, clientID string, clientCertificate string, options *TokenCredentialOptions) (*ClientCertificateCredential, error) {
	_, err := os.Stat(clientCertificate)
	if err != nil {
		credErr := &CredentialUnavailableError{CredentialType: "Client Certificate Credential", Message: "Certificate file not found in path: " + clientCertificate}
		azcore.Log().Write(azcore.LogError, logCredentialError(credErr.CredentialType, credErr))
		return nil, credErr
	}
	c, err := newAADIdentityClient(options)
	if err != nil {
		return nil, err
	}
	return &ClientCertificateCredential{tenantID: tenantID, clientID: clientID, clientCertificate: clientCertificate, client: c}, nil
}

// GetToken obtains a token from Azure Active Directory, using the certificate in the file path.
// scopes: The list of scopes for which the token will have access.
// ctx: controlling the request lifetime.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *ClientCertificateCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	tk, err := c.client.authenticateCertificate(ctx, c.tenantID, c.clientID, c.clientCertificate, opts.Scopes)
	if err != nil {
		addGetTokenFailureLogs("Client Certificate Credential", err)
		return nil, err
	}
	azcore.Log().Write(LogCredential, logGetTokenSuccess(c, opts))
	return tk, nil
}

// AuthenticationPolicy implements the azcore.Credential interface on ClientSecretCredential.
func (c *ClientCertificateCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}
