// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const (
	deviceCodeGrantType = "urn:ietf:params:oauth:grant-type:device_code"
)

// DeviceCodeCredentialOptions contains optional parameters for DeviceCodeCredential.
type DeviceCodeCredentialOptions struct {
	azcore.ClientOptions

	// TenantID is the Azure Active Directory tenant the credential authenticates in. Defaults to the
	// "organizations" tenant, which can authenticate work and school accounts. Required for single-tenant
	// applications.
	TenantID string
	// ClientID is the ID of the application users will authenticate to.
	// Defaults to the ID of an Azure development application.
	ClientID string
	// UserPrompt controls how the credential presents authentication instructions. The credential calls
	// this function with authentication details when it receives a device code. By default, the credential
	// prints these details to stdout.
	UserPrompt func(context.Context, DeviceCodeMessage) error
	// AuthorityHost is the base URL of an Azure Active Directory authority. Defaults
	// to the value of environment variable AZURE_AUTHORITY_HOST, if set, or AzurePublicCloud.
	AuthorityHost AuthorityHost
}

func (o *DeviceCodeCredentialOptions) init() {
	if o.TenantID == "" {
		o.TenantID = organizationsTenantID
	}
	if o.ClientID == "" {
		o.ClientID = developerSignOnClientID
	}
	if o.UserPrompt == nil {
		o.UserPrompt = func(ctx context.Context, dc DeviceCodeMessage) error {
			fmt.Println(dc.Message)
			return nil
		}
	}
}

// DeviceCodeMessage contains the information a user needs to complete authentication.
type DeviceCodeMessage struct {
	// UserCode is the user code returned by the service.
	UserCode string `json:"user_code"`
	// VerificationURL is the URL at which the user must authenticate.
	VerificationURL string `json:"verification_uri"`
	// Message is user instruction from Azure Active Directory.
	Message string `json:"message"`
}

// DeviceCodeCredential acquires tokens for a user via the device code flow, which has the
// user browse to an Azure Active Directory URL, enter a code, and authenticate. It's useful
// for authenticating a user in an environment without a web browser, such as an SSH session.
// If a web browser is available, InteractiveBrowserCredential is more convenient because it
// automatically opens a browser to the login page.
type DeviceCodeCredential struct {
	client       *aadIdentityClient
	tenantID     string
	clientID     string
	userPrompt   func(context.Context, DeviceCodeMessage) error
	refreshToken string
}

// NewDeviceCodeCredential creates a DeviceCodeCredential.
// options: Optional configuration.
func NewDeviceCodeCredential(options *DeviceCodeCredentialOptions) (*DeviceCodeCredential, error) {
	cp := DeviceCodeCredentialOptions{}
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
	return &DeviceCodeCredential{tenantID: cp.TenantID, clientID: cp.ClientID, userPrompt: cp.UserPrompt, client: c}, nil
}

// GetToken obtains a token from Azure Active Directory. It will begin the device code flow and poll until the user completes authentication.
// This method is called automatically by Azure SDK clients.
// ctx: Context used to control the request lifetime.
// opts: Options for the token request, in particular the desired scope of the access token.
func (c *DeviceCodeCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	// ensure we request offline_access
	scopes := make([]string, len(opts.Scopes))
	copy(scopes, opts.Scopes)
	for i, scope := range scopes {
		if scope == "offline_access" {
			break
		} else if i == len(scopes)-1 {
			scopes = append(scopes, "offline_access")
		}
	}
	if len(c.refreshToken) != 0 {
		tk, err := c.client.refreshAccessToken(ctx, c.tenantID, c.clientID, "", c.refreshToken, scopes)
		if err != nil {
			addGetTokenFailureLogs("Device Code Credential", err, true)
			return nil, err
		}
		c.refreshToken = tk.refreshToken
		logGetTokenSuccess(c, opts)
		return tk.token, nil
	}
	dc, err := c.client.requestNewDeviceCode(ctx, c.tenantID, c.clientID, scopes)
	if err != nil {
		authErr := newAuthenticationFailedError(err, nil)
		addGetTokenFailureLogs("Device Code Credential", authErr, true)
		return nil, authErr
	}
	// send authentication flow instructions back to the user to log in and authorize the device
	err = c.userPrompt(ctx, DeviceCodeMessage{
		UserCode:        dc.UserCode,
		VerificationURL: dc.VerificationURL,
		Message:         dc.Message,
	})
	if err != nil {
		return nil, err
	}
	// poll the token endpoint until a valid access token is received or until authentication fails
	for {
		tk, err := c.client.authenticateDeviceCode(ctx, c.tenantID, c.clientID, dc.DeviceCode, scopes)
		// if there is no error, save the refresh token and return the token credential
		if err == nil {
			c.refreshToken = tk.refreshToken
			logGetTokenSuccess(c, opts)
			return tk.token, err
		}
		// if there is an error, check for an AADAuthenticationFailedError in order to check the status for token retrieval
		// if the error is not an AADAuthenticationFailedError, then fail here since something unexpected occurred
		var authFailed AuthenticationFailedError
		if errors.As(err, &authFailed) && strings.Contains(authFailed.Error(), "authorization_pending") {
			// wait for the interval specified from the initial device code endpoint and then poll for the token again
			time.Sleep(time.Duration(dc.Interval) * time.Second)
		} else {
			addGetTokenFailureLogs("Device Code Credential", err, true)
			return nil, err
		}
	}
}

// deviceCodeResult is used to store device code related information to help the user login and allow the device code flow to continue
// to request a token to authenticate a user
type deviceCodeResult struct {
	UserCode        string `json:"user_code"`        // User code returned by the service.
	DeviceCode      string `json:"device_code"`      // Device code returned by the service.
	VerificationURL string `json:"verification_uri"` // Verification URL where the user must navigate to authenticate using the device code and credentials.
	Interval        int64  `json:"interval"`         // Polling interval time to check for completion of authentication flow.
	Message         string `json:"message"`          // User friendly text response that can be used for display purposes.
}

var _ azcore.TokenCredential = (*DeviceCodeCredential)(nil)
