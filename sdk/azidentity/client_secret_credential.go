// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ClientSecretCredential enables authentication to Azure Active Directory using a client secret that was generated for an App Registration.  More information on how
// to configure a client secret can be found here:
// https://docs.microsoft.com/en-us/azure/active-directory/develop/quickstart-configure-app-access-web-apis#add-credentials-to-your-web-application
type ClientSecretCredential struct {
	client       *aadIdentityClient
	tenantID     string // Gets the Azure Active Directory tenant (directory) Id of the service principal
	clientID     string // Gets the client (application) ID of the service principal
	clientSecret string // Gets the client secret that was generated for the App Registration used to authenticate the client.
}

// NewClientSecretCredential constructs a new ClientSecretCredential with the details needed to authenticate against Azure Active Directory with a client secret.
// tenantID: The Azure Active Directory tenant (directory) Id of the service principal.
// clientID: The client (application) ID of the service principal.
// clientSecret: A client secret that was generated for the App Registration used to authenticate the client.
// options: allow to configure the management of the requests sent to the Azure Active Directory service.
func NewClientSecretCredential(tenantID string, clientID string, clientSecret string, options *TokenCredentialOptions) (*ClientSecretCredential, error) {
	c, err := newAADIdentityClient(options)
	if err != nil {
		return nil, err
	}
	c.options.AuthorityHost.Path = path.Join(c.options.AuthorityHost.Path, tenantID+tokenEndpoint)
	return &ClientSecretCredential{tenantID: tenantID, clientID: clientID, clientSecret: clientSecret, client: c}, nil
}

// GetToken obtains a token from the Azure Active Directory service, using the specified client secret to authenticate.
// ctx: controlling the request lifetime.
// scopes: The list of scopes for which the token will have access.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *ClientSecretCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	tk, err := c.client.authenticate(ctx, c.tenantID, c.clientID, c.clientSecret, opts.Scopes)
	log := azcore.Log()
	if err != nil {
		msg := fmt.Sprintf("Azure Identity => ERROR in GetToken() call for %T: %s", c, err.Error())
		log.Write(CredentialClassification, msg)
	} else {
		msg := fmt.Sprintf("Azure Identity => GetToken() result for %T: SUCCESS", c)
		log.Write(CredentialClassification, msg)

	}
	if log.Should(CredentialClassificationVerbose) {
		vmsg := fmt.Sprintf("Azure Identity => Scopes: [%s]", strings.Join(opts.Scopes, ", "))
		log.Write(CredentialClassificationVerbose, vmsg)
	}
	return tk, err
}

// AuthenticationPolicy implements the azcore.Credential interface on ClientSecretCredential.
func (c *ClientSecretCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}

var _ azcore.TokenCredential = (*ClientSecretCredential)(nil)
