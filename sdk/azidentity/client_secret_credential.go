// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/atomic"
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

// GetToken obtains a token from the Azure Active Directory service, using the specified client secret to authenticate.
// ctx: controlling the request lifetime.
// scopes: The list of scopes for which the token will have access.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *ClientSecretCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	if len(opts.Scopes) == 0 {
		return nil, &AuthenticationFailedError{Message: "You need to include valid scopes in order to request a token with this credential"}
	}
	return c.client.authenticate(ctx, c.tenantID, c.clientID, c.clientSecret, opts.Scopes)
}

// AuthenticationPolicy implements the azcore.Credential interface on ClientSecretCredential.
func (c *ClientSecretCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options.Scopes)
}

type bearerTokenPolicy struct {
	empty     chan bool    // indicates we don't have a token yet
	updating  atomic.Int64 // atomically set to 1 to indicate a refresh is in progress
	header    atomic.String
	expiresOn atomic.Time
	creds     azcore.TokenCredential // R/O
	scopes    []string               // R/O
}

func newBearerTokenPolicy(creds azcore.TokenCredential, scopes []string) *bearerTokenPolicy {
	bt := bearerTokenPolicy{creds: creds, scopes: scopes}
	// prime channel indicating there is no token
	bt.empty = make(chan bool, 1)
	bt.empty <- true
	return &bt
}

func (b *bearerTokenPolicy) Do(ctx context.Context, req *azcore.Request) (*azcore.Response, error) {
	if req.URL.Scheme != "https" {
		// HTTPS must be used, otherwise the tokens are at the risk of being exposed
		return nil, &AuthenticationFailedError{
			Message: "token credentials require a URL using the HTTPS protocol scheme",
		}
	}
	// check if the token is empty.  this will return true for
	// the first read.  all other reads will block until the
	// token has been obtained at which point it returns false
	if <-b.empty {
		tk, err := b.creds.GetToken(ctx, azcore.TokenRequestOptions{Scopes: b.scopes})
		if err != nil {
			// failed to get a token, let another go routine try
			b.empty <- true
			return nil, err
		}
		b.header.Store("Bearer " + tk.Token)
		b.expiresOn.Store(tk.ExpiresOn)
		// signal token has been initialized.  go routines
		// that read from b.empty will always get false.
		close(b.empty)
	}
	// create a "refresh window" before the token's real expiration date.
	// this allows callers to continue to use the old token while the
	// refresh is in progress.
	const window = 2 * time.Minute
	// check if the token has expired
	now := time.Now().UTC()
	exp := b.expiresOn.Load()
	if now.After(exp.Add(-window)) {
		// token has expired, set update flag and begin the refresh.
		// if no other go routine has initiated a refresh the calling
		// go routine will do it.
		if b.updating.CAS(0, 1) {
			tk, err := b.creds.GetToken(ctx, azcore.TokenRequestOptions{Scopes: b.scopes})
			if err != nil {
				// clear updating flag before returning so other
				// go routines can attempt to refresh
				b.updating.Store(0)
				return nil, err
			}
			b.header.Store("Bearer " + tk.Token)
			// set expiresOn last since refresh is predicated on it
			b.expiresOn.Store(tk.ExpiresOn)
			// signal the refresh is complete. this must happen after all shared
			// state has been updated.
			b.updating.Store(0)
		} // else { another go routine is refreshing the token, use previous token }
	}
	req.Request.Header.Set(azcore.HeaderAuthorization, b.header.Load())
	return req.Do(ctx)
}

var _ azcore.TokenCredential = (*ClientSecretCredential)(nil)

// TODO: rename credentialpolicyoptions
// TODO: wrap scopes in a token request type
