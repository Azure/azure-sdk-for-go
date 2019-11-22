// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"sync"
	"time"

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
func NewClientSecretCredential(tenantID string, clientID string, clientSecret string, options *IdentityClientOptions) *ClientSecretCredential {
	return &ClientSecretCredential{tenantID: tenantID, clientID: clientID, clientSecret: clientSecret, client: newAADIdentityClient(options)}
}

// TODO: make sure guid params are always in the same order

// GetToken obtains a token from the Azure Active Directory service, using the specified client secret to authenticate.
// ctx: controlling the request lifetime.
// scopes: The list of scopes for which the token will have access.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *ClientSecretCredential) GetToken(ctx context.Context, opts *azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	return c.client.authenticate(ctx, c.tenantID, c.clientID, c.clientSecret, opts.Scopes)
}

// AuthenticationPolicy implements the azcore.Credential interface on ClientSecretCredential.
func (c *ClientSecretCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options.Scopes)
}

type bearerTokenPolicy struct {
	// take lock when manipulating header/expiresOn fields
	lock      sync.RWMutex
	header    string
	expiresOn time.Time
	creds     azcore.TokenCredential // R/O
	scopes    []string               // R/O
}

func newBearerTokenPolicy(creds azcore.TokenCredential, scopes []string) *bearerTokenPolicy {
	// set the token as expired so first call to Do() refreshes it
	return &bearerTokenPolicy{creds: creds, scopes: scopes, expiresOn: time.Now().UTC()}
}

func (b *bearerTokenPolicy) Do(ctx context.Context, req *azcore.Request) (*azcore.Response, error) {
	var bt string
	// take read lock and check if the token has expired
	b.lock.RLock()
	now := time.Now().UTC()
	if now.Equal(b.expiresOn) || now.After(b.expiresOn) {
		// token has expired, take the write lock then check again
		b.lock.RUnlock()
		// don't defer Unlock(), we want to release it ASAP
		b.lock.Lock()
		if now.Equal(b.expiresOn) || now.After(b.expiresOn) {
			// token has expired, get a new one and update shared state
			tk, err := b.creds.GetToken(ctx, &azcore.TokenRequestOptions{Scopes: b.scopes})
			if err != nil {
				b.lock.Unlock()
				return nil, err
			}
			b.expiresOn = tk.ExpiresOn
			b.header = "Bearer " + tk.Token
		} // else { another go routine already refreshed the token }
		bt = b.header
		b.lock.Unlock()
	} else {
		// token is still valid
		bt = b.header
		b.lock.RUnlock()
	}
	// no locks are to be held at this point
	req.Request.Header.Set(azcore.HeaderAuthorization, bt)
	return req.Do(ctx)
}

var _ azcore.TokenCredential = (*ClientSecretCredential)(nil)

// TODO: rename credentialpolicyoptions
// TODO: wrap scopes in a token request type
