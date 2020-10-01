// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/pkg/browser"
)

// InteractiveBrowserCredential enables authentication to Azure Active Directory using an interactive browser to log in.
type InteractiveBrowserCredential struct {
	client       *aadIdentityClient
	tenantID     string // Gets the Azure Active Directory tenant (directory) ID of the service principal
	clientID     string // Gets the client (application) ID of the service principal
	clientSecret string // Gets the client secret that was generated for the App Registration used to authenticate the client.
}

// NewInteractiveBrowserCredential constructs a new InteractiveBrowserCredential with the details needed to authenticate against Azure Active Directory through an interactive browser window.
// tenantID: The Azure Active Directory tenant (directory) ID of the service principal.
// clientID: The client (application) ID of the service principal.
// clientSecret: Gets the client secret that was generated for the App Registration used to authenticate the client.
// options: allow to configure the management of the requests sent to Azure Active Directory.
func NewInteractiveBrowserCredential(tenantID string, clientID string, clientSecret string, options *TokenCredentialOptions) (*InteractiveBrowserCredential, error) {
	c, err := newAADIdentityClient(options)
	if err != nil {
		return nil, err
	}
	return &InteractiveBrowserCredential{tenantID: tenantID, clientID: clientID, clientSecret: clientSecret, client: c}, nil
}

// GetToken obtains a token from Azure Active Directory using an interactive browser to authenticate.
// ctx: Context used to control the request lifetime.
// opts: TokenRequestOptions contains the list of scopes for which the token will have access.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *InteractiveBrowserCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	tk, err := c.client.authenticateInteractiveBrowser(ctx, c.tenantID, c.clientID, c.clientSecret, opts.Scopes)
	if err != nil {
		addGetTokenFailureLogs("Interactive Browser Credential", err, true)
		return nil, err
	}
	logGetTokenSuccess(c, opts)
	return tk, nil
}

// AuthenticationPolicy implements the azcore.Credential interface on InteractiveBrowserCredential and calls the Bearer Token policy
// to get the bearer token.
func (c *InteractiveBrowserCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}

var _ azcore.TokenCredential = (*InteractiveBrowserCredential)(nil)

var authCodeReceiver = func(tenantID string, clientID string, scopes []string) (*interactiveConfig, error) {
	return interactiveBrowserLogin(tenantID, clientID, scopes)
}

// interactiveBrowserLogin opens an interactive browser with the specified tenant and client IDs provided then returns the authorization code
// received or an error.
func interactiveBrowserLogin(tenantID string, clientID string, scopes []string) (*interactiveConfig, error) {
	const authURLFormat = "https://login.microsoftonline.com/%s/oauth2/v2.0/authorize?response_type=code&response_mode=query&client_id=%s&redirect_uri=%s&state=%s&scope=%s&prompt=select_account"
	if tenantID == "" {
		tenantID = "common"
	}
	if clientID == "" {
		clientID = "04b07795-8ddb-461a-bbee-02f9e1bf7b46"
	}
	state := func() string {
		rand.Seed(time.Now().Unix())
		// generate a 20-char random alpha-numeric string
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		buff := make([]byte, 20)
		for i := range buff {
			buff[i] = charset[rand.Intn(len(charset))]
		}
		return string(buff)
	}()
	// start local redirect server so login can call us back
	rs := newServer()
	redirectURL := rs.Start(state)
	defer rs.Stop()
	authURL := fmt.Sprintf(authURLFormat, tenantID, clientID, redirectURL, state, strings.Join(scopes, " "))
	fmt.Println(authURL)
	// open browser window so user can select credentials
	err := browser.OpenURL(authURL)
	if err != nil {
		return nil, err
	}
	// now wait until the logic calls us back
	rs.WaitForCallback()

	authCode, err := rs.AuthorizationCode()
	if err != nil {
		return nil, err
	}
	return &interactiveConfig{
		authCode:    authCode,
		redirectURI: redirectURL,
	}, nil
}
