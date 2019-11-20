// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"runtime"
	"sync/atomic"
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
func (c *ClientSecretCredential) GetToken(ctx context.Context, scopes []string) (*azcore.AccessToken, error) {
	return c.client.authenticate(ctx, c.tenantID, c.clientID, c.clientSecret, scopes)
}

// AuthenticationPolicy implements the azcore.Credential interface on ClientSecretCredential.
func (c *ClientSecretCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options.Scopes)
}

type bearerTokenPolicy struct {
	updating  int64                  // atomically set to 1 to indicate a refresh is in progress
	header    atomic.Value           // string
	expiresOn atomic.Value           // time.Time
	creds     azcore.TokenCredential // R/O
	scopes    []string               // R/O
}

func newBearerTokenPolicy(creds azcore.TokenCredential, scopes []string) *bearerTokenPolicy {
	bt := bearerTokenPolicy{creds: creds, scopes: scopes}
	// set initial values to their zero-value
	bt.header.Store("")
	bt.expiresOn.Store(time.Time{})
	return &bt
}

func (b *bearerTokenPolicy) Do(ctx context.Context, req *azcore.Request) (*azcore.Response, error) {
	// create a "refresh window" before the token's real expiration date.
	// this allows callers to continue to use the old token while the
	// refresh is in progress.
	const window = 2 * time.Minute
	var bt string
	// check if the token has expired
	now := time.Now().UTC()
	exp := b.expiresOn.Load().(time.Time)
	if now.After(exp.Add(-window)) {
		// token has expired, set update flag and begin the refresh.
		// if no other go routine has initiated a refresh the calling
		// go routine will do it.
		if atomic.CompareAndSwapInt64(&b.updating, 0, 1) {
			tk, err := b.creds.GetToken(ctx, b.scopes)
			if err != nil {
				// clear updating flag before returning
				atomic.StoreInt64(&b.updating, 0)
				return nil, err
			}
			// we must set header before expiresOn.  this is to prevent a race
			// when setting header for the very first time.  if expiresOn is set
			// first it's possible to get an empty header.
			b.header.Store("Bearer " + tk.Token)
			b.expiresOn.Store(tk.ExpiresOn)
			// signal the refresh is complete. this must happen after all shared
			// state has been updated.
			atomic.StoreInt64(&b.updating, 0)
		} // else { another go routine is refreshing the token }
		for {
			bt = b.header.Load().(string)
			// for the very first call to Do() the header will be
			// the empty string.  in this case we need to spin-wait
			// until header contains a value.  subsequent calls to
			// Do() will never spin-wait.
			if bt != "" {
				break
			}
			runtime.Gosched()
		}
	} else {
		// token is still valid
		bt = b.header.Load().(string)
	}
	req.Request.Header.Set(azcore.HeaderAuthorization, bt)
	return req.Do(ctx)
}

var _ azcore.TokenCredential = (*ClientSecretCredential)(nil)
