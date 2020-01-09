// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	deviceCodeGrantType = "urn:ietf:params:oauth:grant-type:device_code"
)

// DeviceCodeResult is used to store device code related information to help the user login and allow the device code flow to continue
// to request a token to authenticate a user
type DeviceCodeResult struct {
	UserCode        string `json:"user_code"`        // User code returned by the service
	DeviceCode      string `json:"device_code"`      // Device code returned by the service
	VerificationURL string `json:"verification_uri"` // Verification URL where the user must navigate to authenticate using the device code and credentials.
	Interval        int64  `json:"interval"`         // Polling interval time to check for completion of authentication flow.
	Message         string `json:"message"`          // User friendly text response that can be used for display purpose.
}

// DeviceCodeCredential authenticates a user using the device code flow, and provides access tokens for that user account.
// For more information on the device code authentication flow see https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-device-code
type DeviceCodeCredential struct {
	client       *aadIdentityClient
	tenantID     string       // Gets the Azure Active Directory tenant (directory) Id of the service principal
	clientID     string       // Gets the client (application) ID of the service principal
	callback     func(string) // Sends the user a message with a verification URL and device code to sign in to the login server
	refreshToken string       // Gets the refresh token sent from the service and will be used to retreive new access tokens after the initial request for a token
}

// NewDeviceCodeCredential constructs a new DeviceCodeCredential with the details needed to authenticate against Azure Active Directory with a device code.
// tenantID: The Azure Active Directory tenant (directory) Id of the service principal.
// clientID: The client (application) ID of the service principal.
// callback: The callback function used to send the login message back to the user
// options: allow to configure the management of the requests sent to the Azure Active Directory service.
func NewDeviceCodeCredential(tenantID string, clientID string, callback func(string), options *TokenCredentialOptions) *DeviceCodeCredential {
	return &DeviceCodeCredential{tenantID: tenantID, clientID: clientID, callback: callback, client: newAADIdentityClient(options)}
}

// GetToken obtains a token from the Azure Active Directory service, following the device code authentication
// flow. This function first requests a device code and requests that the user login to continue authenticating.
// This function will keep polling the service for a token meanwhile the user logs.
// - scopes: The list of scopes for which the token will have access.
// - ctx: controlling the request lifetime.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *DeviceCodeCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	if len(c.refreshToken) != 0 {
		tk, err := c.client.refreshAccessToken(ctx, c.tenantID, c.clientID, "", c.refreshToken, opts.Scopes)
		if tk == nil { // this would only happen if there was an unmarshal error or a failure to parse ExpiresOn
			return nil, err
		}
		// assign new refresh token to the credential for future use
		c.refreshToken = tk.refreshToken
		// passing the access token and/or error back up
		return tk.token, err
	}
	// if there is no refreshToken, then begin the Device Code flow from the beginning
	// make initial request to the device code endpoint for a device code and instructions for authentication
	dc, err := c.client.requestNewDeviceCode(ctx, c.tenantID, c.clientID, opts.Scopes)
	if err != nil {
		return nil, err
	}
	// send authentication flow instructions back to the user to log in and authorize the device
	c.callback(dc.Message)
	// poll the token endpoint until a valid access token is received or until authentication fails
	for {
		tk, err := c.client.authenticateDeviceCode(ctx, c.tenantID, c.clientID, dc.DeviceCode, opts.Scopes)
		// if there is no error, save the refresh token and return the token credential
		if err == nil {
			c.refreshToken = tk.refreshToken
			return tk.token, err
		}
		// if there is an error, check for an AADAuthenticationFailedError in order to check the status for token retrieval
		var authRespErr *AADAuthenticationFailedError
		// if the error is not an AADAuthenticationFailedError, then fail here since something unexpected occurred
		if !errors.As(err, &authRespErr) {
			return nil, err
		}
		switch authRespErr.Message {
		// wait for the interval specified from the initial device code endpoint and then poll for the token again
		case "authorization_pending":
			time.Sleep(time.Duration(dc.Interval) * time.Second)
		default:
			// any other error should be returned
			return nil, err
		}
	}
}

// AuthenticationPolicy implements the azcore.Credential interface on ClientSecretCredential.
func (c *DeviceCodeCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}
