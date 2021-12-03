// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"math/rand"
	"net/url"
	"path"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/pkg/browser"
)

// InteractiveBrowserCredentialOptions contains optional parameters for InteractiveBrowserCredential.
type InteractiveBrowserCredentialOptions struct {
	azcore.ClientOptions

	// TenantID is the Azure Active Directory tenant the credential authenticates in. Defaults to the
	// "organizations" tenant, which can authenticate work and school accounts.
	TenantID string
	// ClientID is the ID of the application users will authenticate to.
	// Defaults to the ID of an Azure development application.
	ClientID string
	// RedirectURL will be supported in a future version but presently doesn't work: https://github.com/Azure/azure-sdk-for-go/issues/15632.
	// Applications which have "http://localhost" registered as a redirect URL need not set this option.
	RedirectURL string
	// AuthorityHost is the base URL of an Azure Active Directory authority. Defaults
	// to the value of environment variable AZURE_AUTHORITY_HOST, if set, or AzurePublicCloud.
	AuthorityHost AuthorityHost
}

func (o *InteractiveBrowserCredentialOptions) init() {
	if o.TenantID == "" {
		o.TenantID = organizationsTenantID
	}
	if o.ClientID == "" {
		o.ClientID = developerSignOnClientID
	}
}

// InteractiveBrowserCredential opens a browser to interactively authenticate a user.
type InteractiveBrowserCredential struct {
	client  *aadIdentityClient
	options InteractiveBrowserCredentialOptions
}

// NewInteractiveBrowserCredential constructs a new InteractiveBrowserCredential.
// options: Optional configuration.
func NewInteractiveBrowserCredential(options *InteractiveBrowserCredentialOptions) (*InteractiveBrowserCredential, error) {
	cp := InteractiveBrowserCredentialOptions{}
	if options != nil {
		cp = *options
	}
	cp.init()
	if !validTenantID(cp.TenantID) {
		return nil, errors.New(tenantIDValidationErr)
	}
	authorityHost, err := setAuthorityHost(cp.AuthorityHost)
	if err != nil {
		return nil, err
	}
	c, err := newAADIdentityClient(authorityHost, &cp.ClientOptions)
	if err != nil {
		return nil, err
	}
	return &InteractiveBrowserCredential{options: cp, client: c}, nil
}

// GetToken obtains a token from Azure Active Directory. This method is called automatically by Azure SDK clients.
// ctx: Context used to control the request lifetime.
// opts: Options for the token request, in particular the desired scope of the access token.
func (c *InteractiveBrowserCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	tk, err := c.client.authenticateInteractiveBrowser(ctx, &c.options, opts.Scopes)
	if err != nil {
		addGetTokenFailureLogs("Interactive Browser Credential", err, true)
		return nil, err
	}
	logGetTokenSuccess(c, opts)
	return tk, nil
}

var _ azcore.TokenCredential = (*InteractiveBrowserCredential)(nil)

// authCodeReceiver is used to allow for testing without opening an interactive browser window. Allows mocking a response authorization code and redirect URI.
var authCodeReceiver = func(ctx context.Context, authorityHost string, opts *InteractiveBrowserCredentialOptions, scopes []string) (*interactiveConfig, error) {
	return interactiveBrowserLogin(ctx, authorityHost, opts, scopes)
}

// interactiveBrowserLogin opens an interactive browser with the specified tenant and client IDs provided then returns the authorization code
// received or an error.
func interactiveBrowserLogin(ctx context.Context, authorityHost string, opts *InteractiveBrowserCredentialOptions, scopes []string) (*interactiveConfig, error) {
	// start local redirect server so login can call us back
	rs := newServer()
	state, err := uuid.New()
	if err != nil {
		return nil, err
	}
	redirectURL := rs.Start(state.String())
	defer rs.Stop()
	u, err := url.Parse(authorityHost)
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Add("response_type", "code")
	values.Add("response_mode", "query")
	values.Add("client_id", opts.ClientID)
	values.Add("redirect_uri", redirectURL)
	values.Add("state", state.String())
	values.Add("scope", strings.Join(scopes, " "))
	values.Add("prompt", "select_account")
	cv := ""
	// the code verifier is a random 32-byte sequence that's been base-64 encoded without padding.
	// it's used to prevent MitM attacks during auth code flow, see https://tools.ietf.org/html/rfc7636
	b := make([]byte, 32) // nolint:gosimple
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	cv = base64.RawURLEncoding.EncodeToString(b)
	// for PKCE, create a hash of the code verifier
	cvh := sha256.Sum256([]byte(cv))
	values.Add("code_challenge", base64.RawURLEncoding.EncodeToString(cvh[:]))
	values.Add("code_challenge_method", "S256")
	u.Path = runtime.JoinPaths(u.Path, opts.TenantID, path.Join(oauthPath(opts.TenantID), "/authorize"))
	u.RawQuery = values.Encode()
	// open browser window so user can select credentials
	if err = browser.OpenURL(u.String()); err != nil {
		return nil, err
	}
	// now wait until the logic calls us back
	if err := rs.WaitForCallback(ctx); err != nil {
		return nil, err
	}

	authCode, err := rs.AuthorizationCode()
	if err != nil {
		return nil, err
	}
	return &interactiveConfig{
		authCode:     authCode,
		codeVerifier: cv,
		redirectURI:  redirectURL,
	}, nil
}
