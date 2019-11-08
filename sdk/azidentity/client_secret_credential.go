// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"

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
func NewClientSecretCredential(tenantID string, clientID string, clientSecret string, options *IdentityClientOptions) (*ClientSecretCredential, error) {
	return &ClientSecretCredential{tenantID: tenantID, clientID: clientID, clientSecret: clientSecret, client: newAADIdentityClient(options)}, nil
}

// NewClientSecretCredentialWithPipeline constructs a new ClientSecretCredential with the details needed to authenticate against Azure Active Directory with a client secret.
// tenantID: The Azure Active Directory tenant (directory) Id of the service principal.
// clientID: The client (application) ID of the service principal.
// clientSecret: A client secret that was generated for the App Registration used to authenticate the client.
// options: allow to configure the management of the requests sent to the Azure Active Directory service.
// pipeline: Custom pipeline to be used for API requests.
func NewClientSecretCredentialWithPipeline(tenantID string, clientID string, clientSecret string, options *IdentityClientOptions, pipeline azcore.Pipeline) (*ClientSecretCredential, error) {
	return &ClientSecretCredential{tenantID: tenantID, clientID: clientID, clientSecret: clientSecret, client: newAADIdentityClientWithPipeline(options, pipeline)}, nil
}

// GetToken obtains a token from the Azure Active Directory service, using the specified client secret to authenticate.
// ctx: controlling the request lifetime.
// scopes: The list of scopes for which the token will have access.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c ClientSecretCredential) GetToken(ctx context.Context, scopes []string) (*azcore.AccessToken, error) {
	return c.client.authenticate(ctx, c.tenantID, c.clientID, c.clientSecret, scopes)
}

// Policy implements the azcore.Credential interface on ClientSecretCredential.
func (c ClientSecretCredential) Policy(options azcore.CredentialPolicyOptions) azcore.Policy {
	return azcore.PolicyFunc(func(ctx context.Context, req *azcore.Request) (*azcore.Response, error) {
		return applyToken(ctx, c, req, options)
	})
}

func applyToken(ctx context.Context, token azcore.TokenCredential, req *azcore.Request, options azcore.CredentialPolicyOptions) (*azcore.Response, error) {
	tk, err := token.GetToken(ctx, options.Scopes)
	if err != nil {
		return nil, err
	}
	req.Request.Header.Set(azcore.HeaderAuthorization, "Bearer "+tk.Token)
	return req.Do(ctx)
}
