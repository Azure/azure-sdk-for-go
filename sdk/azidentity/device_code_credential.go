// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	deviceCodeGrantType = "urn:ietf:params:oauth:grant-type:device_code"
)

// DeviceCodeResult is used to store device code related information to help the user login and allow the device code flow to continue
// to request a token to authenticate a user
type DeviceCodeResult struct {
	UserCode        string        `json:"user_code"`        // User code returned by the service
	DeviceCode      string        `json:"device_code"`      // Device code returned by the service
	VerificationURL string        `json:"verification_uri"` // Verification URL where the user must navigate to authenticate using the device code and credentials. JMR: URL?
	ExpiresOn       time.Duration // Time when the device code will expire.
	Interval        int64         `json:"interval"` // Polling interval time to check for completion of authentication flow.
	Message         string        `json:"message"`  // User friendly text response that can be used for display purpose.
	ClientID        string        // Identifier of the client requesting device code.
	Scopes          []string      // List of the scopes that would be held by token. JMR: Should be readonly
}

// DeviceCodeCredential authenticates a user using the device code flow, and provides access tokens for that user account.
// For more information on the device code authentication flow see https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-device-code
type DeviceCodeCredential struct {
	client   *aadIdentityClient
	tenantID string       // Gets the Azure Active Directory tenant (directory) Id of the service principal
	clientID string       // Gets the client (application) ID of the service principal
	callback func(string) // Sends the user a message with a verification URL and device code to sign in to the login server
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
	dc, err := c.client.requestNewDeviceCode(ctx, c.tenantID, c.clientID, opts.Scopes)
	if err != nil {
		return nil, err
	}
	c.callback(dc.Message)
	for {
		tk, err := c.client.authenticateDeviceCode(ctx, c.tenantID, c.clientID, dc.DeviceCode, opts.Scopes)
		if err == nil {
			return tk, err
		}
		var authFailed *AuthenticationFailedError
		if !errors.As(err, &authFailed) {
			return nil, err
		}
		var authRespErr *AuthenticationResponseError
		if !errors.As(authFailed.Unwrap(), &authRespErr) {
			return nil, err
		}
		switch authRespErr.Message {
		case "authorization_pending":
			time.Sleep(time.Duration(dc.Interval) * time.Second)
		default:
			// Any other error should be returned
			return nil, err
		}
	}
}

// AuthenticationPolicy implements the azcore.Credential interface on ClientSecretCredential.
func (c *DeviceCodeCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}

func createDeviceCodeResult(res *azcore.Response) (*DeviceCodeResult, error) {
	value := &DeviceCodeResult{}
	if err := json.Unmarshal(res.Payload, &value); err != nil {
		return nil, fmt.Errorf("DeviceCodeResult: %w", err)
	}
	return value, nil
}
