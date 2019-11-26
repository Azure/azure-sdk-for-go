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
	// cond is used to synchronize token refresh.  the locker
	// must be locked when updating the following shared state.
	cond *sync.Cond

	// renewing indicates that the token is in the process of being refreshed
	renewing bool

	// header contains the authorization header value
	header string

	// expiresOn is when the token will expire
	expiresOn time.Time

	// the following fields are read-only
	creds  azcore.TokenCredential
	scopes []string
}

func newBearerTokenPolicy(creds azcore.TokenCredential, scopes []string) *bearerTokenPolicy {
	return &bearerTokenPolicy{
		cond:   sync.NewCond(&sync.Mutex{}),
		creds:  creds,
		scopes: scopes,
	}
}

func (b *bearerTokenPolicy) Do(ctx context.Context, req *azcore.Request) (*azcore.Response, error) {
	if req.URL.Scheme != "https" {
		// HTTPS must be used, otherwise the tokens are at the risk of being exposed
		return nil, &AuthenticationFailedError{
			Message: "token credentials require a URL using the HTTPS protocol scheme",
		}
	}
	// create a "refresh window" before the token's real expiration date.
	// this allows callers to continue to use the old token while the
	// refresh is in progress.
	const window = 2 * time.Minute
	now, getToken, header := time.Now(), false, ""
	// acquire exclusive lock
	b.cond.L.Lock()
	for {
		if b.expiresOn.IsZero() || b.expiresOn.Before(now) {
			// token was never obtained or has expired
			if !b.renewing {
				// another go routine isn't refreshing the token so this one will
				b.renewing = true
				getToken = true
				break
			}
			// getting here means this go routine will wait for the token to refresh
		} else if b.expiresOn.Add(-window).Before(now) {
			// token is within the expiration window
			if !b.renewing {
				// another go routine isn't refreshing the token so this one will
				b.renewing = true
				getToken = true
				break
			}
			// this go routine will use the existing token while another refreshes it
			header = b.header
			break
		} else {
			// token is not expiring yet so use it as-is
			header = b.header
			break
		}
		// wait for the token to refresh
		b.cond.Wait()
	}
	b.cond.L.Unlock()
	if getToken {
		// this go routine has been elected to refresh the token
		tk, err := b.creds.GetToken(ctx, b.scopes)
		if err != nil {
			return nil, err
		}
		header = "Bearer " + tk.Token
		// update shared state
		b.cond.L.Lock()
		b.renewing = false
		b.header = header
		b.expiresOn = tk.ExpiresOn
		// signal any waiters that the token has been refreshed
		b.cond.Broadcast()
		b.cond.L.Unlock()
	}
	req.Request.Header.Set(azcore.HeaderAuthorization, header)
	return req.Do(ctx)
}

var _ azcore.TokenCredential = (*ClientSecretCredential)(nil)
