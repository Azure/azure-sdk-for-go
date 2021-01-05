// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/pkg/browser"
)

// InteractiveBrowserCredentialOptions can be used when providing additional credential information, such as a client secret.
// Also use these options to modify the default pipeline behavior through the TokenCredentialOptions.
// Call DefaultInteractiveBrowserCredentialOptions() to create an instance populated with default values.
type InteractiveBrowserCredentialOptions struct {
	// The Azure Active Directory tenant (directory) ID of the service principal.
	TenantID string
	// The client (application) ID of the service principal.
	ClientID string
	// The client secret that was generated for the App Registration used to authenticate the client. Only applies for web apps.
	ClientSecret string
	// The redirect URL used to request the authorization code. Must be the same URL that is configured for the App Registration.
	RedirectURL string
	// The localhost port for the local server that will be used to redirect back. If left with a zero value, a random port
	// will be selected.
	Port int
	// The host of the Azure Active Directory authority. The default is AzurePublicCloud.
	// Leave empty to allow overriding the value from the AZURE_AUTHORITY_HOST environment variable.
	AuthorityHost string
	// HTTPClient sets the transport for making HTTP requests
	// Leave this as nil to use the default HTTP transport
	HTTPClient azcore.Transport
	// Retry configures the built-in retry policy behavior
	Retry azcore.RetryOptions
	// Telemetry configures the built-in telemetry policy behavior
	Telemetry azcore.TelemetryOptions
	// Logging configures the built-in logging policy behavior.
	Logging azcore.LogOptions
}

// DefaultInteractiveBrowserCredentialOptions returns an instance of InteractiveBrowserCredentialOptions initialized with default values.
func DefaultInteractiveBrowserCredentialOptions() InteractiveBrowserCredentialOptions {
	return InteractiveBrowserCredentialOptions{
		TenantID:  organizationsTenantID,
		ClientID:  developerSignOnClientID,
		Retry:     azcore.DefaultRetryOptions(),
		Telemetry: azcore.DefaultTelemetryOptions(),
		Logging:   azcore.DefaultLogOptions(),
	}
}

// InteractiveBrowserCredential enables authentication to Azure Active Directory using an interactive browser to log in.
type InteractiveBrowserCredential struct {
	client *aadIdentityClient
	// options contains data necessary to authenticate through an interactive browser window
	options InteractiveBrowserCredentialOptions
}

// NewInteractiveBrowserCredential constructs a new InteractiveBrowserCredential with the details needed to authenticate against Azure Active Directory through an interactive browser window.
// options: allow to configure the management of the requests sent to Azure Active Directory, pass in nil for default behavior.
func NewInteractiveBrowserCredential(options *InteractiveBrowserCredentialOptions) (*InteractiveBrowserCredential, error) {
	if options == nil {
		temp := DefaultInteractiveBrowserCredentialOptions()
		options = &temp
	}
	if !validTenantID(options.TenantID) {
		return nil, &CredentialUnavailableError{credentialType: "Interactive Browser Credential", message: tenantIDValidationErr}
	}
	authorityHost, err := setAuthorityHost(options.AuthorityHost)
	if err != nil {
		return nil, err
	}
	c, err := newAADIdentityClient(authorityHost, pipelineOptions{HTTPClient: options.HTTPClient, Retry: options.Retry, Telemetry: options.Telemetry, Logging: options.Logging})
	if err != nil {
		return nil, err
	}
	return &InteractiveBrowserCredential{options: *options, client: c}, nil
}

// GetToken obtains a token from Azure Active Directory using an interactive browser to authenticate.
// ctx: Context used to control the request lifetime.
// opts: TokenRequestOptions contains the list of scopes for which the token will have access.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *InteractiveBrowserCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	tk, err := c.client.authenticateInteractiveBrowser(ctx, c.options.TenantID, c.options.ClientID, c.options.ClientSecret, c.options.RedirectURL, c.options.Port, opts.Scopes)
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

// authCodeReceiver is used to allow for testing without opening an interactive browser window. Allows mocking a response authorization code and redirect URI.
var authCodeReceiver = func(authorityHost string, tenantID string, clientID string, redirectURI string, port int, scopes []string) (*interactiveConfig, error) {
	return interactiveBrowserLogin(authorityHost, tenantID, clientID, redirectURI, port, scopes)
}

// interactiveBrowserLogin opens an interactive browser with the specified tenant and client IDs provided then returns the authorization code
// received or an error.
func interactiveBrowserLogin(authorityHost string, tenantID string, clientID string, redirectURL string, port int, scopes []string) (*interactiveConfig, error) {
	const authPathFormat = "%s/oauth2/v2.0/authorize?response_type=code&response_mode=query&client_id=%s&redirect_uri=%s&state=%s&scope=%s&prompt=select_account"
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
	if redirectURL == "" {
		redirectURL = rs.Start(state, port)
	}
	defer rs.Stop()
	u, err := url.Parse(authorityHost)
	if err != nil {
		return nil, err
	}
	authPath := fmt.Sprintf(authPathFormat, tenantID, clientID, redirectURL, state, strings.Join(scopes, " "))
	u.Path = azcore.JoinPaths(u.Path, authPath)
	// open browser window so user can select credentials
	err = browser.OpenURL(u.String())
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
