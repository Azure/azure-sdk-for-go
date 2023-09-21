//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const credNameDeviceCode = "DeviceCodeCredential"

// DeviceCodeCredentialOptions contains optional parameters for DeviceCodeCredential.
type DeviceCodeCredentialOptions struct {
	azcore.ClientOptions

	// AdditionallyAllowedTenants specifies additional tenants for which the credential may acquire
	// tokens. Add the wildcard value "*" to allow the credential to acquire tokens for any tenant.
	AdditionallyAllowedTenants []string

	// AuthenticationRecord returned by a call to a credential's Authenticate method. Set this option
	// to enable the credential to use data from a previous authentication.
	AuthenticationRecord AuthenticationRecord

	// ClientID is the ID of the application users will authenticate to.
	// Defaults to the ID of an Azure development application.
	ClientID string

	// DisableAutomaticAuthentication prevents the credential from automatically prompting the user to authenticate.
	// When this option is true, [DeviceCodeCredential.GetToken] will return [ErrAuthenticationRequired] when user
	// interaction is necessary to acquire a token.
	DisableAutomaticAuthentication bool

	// DisableInstanceDiscovery should be set true only by applications authenticating in disconnected clouds, or
	// private clouds such as Azure Stack. It determines whether the credential requests Azure AD instance metadata
	// from https://login.microsoft.com before authenticating. Setting this to true will skip this request, making
	// the application responsible for ensuring the configured authority is valid and trustworthy.
	DisableInstanceDiscovery bool

	// TenantID is the Azure Active Directory tenant the credential authenticates in. Defaults to the
	// "organizations" tenant, which can authenticate work and school accounts. Required for single-tenant
	// applications.
	TenantID string

	// TokenCachePersistenceOptions enables persistent token caching when not nil.
	TokenCachePersistenceOptions *TokenCachePersistenceOptions

	// UserPrompt controls how the credential presents authentication instructions. The credential calls
	// this function with authentication details when it receives a device code. By default, the credential
	// prints these details to stdout.
	UserPrompt func(context.Context, DeviceCodeMessage) error
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
// If a web browser is available, [InteractiveBrowserCredential] is more convenient because it
// automatically opens a browser to the login page.
type DeviceCodeCredential struct {
	client *publicClient
}

// NewDeviceCodeCredential creates a DeviceCodeCredential. Pass nil to accept default options.
func NewDeviceCodeCredential(options *DeviceCodeCredentialOptions) (*DeviceCodeCredential, error) {
	cp := DeviceCodeCredentialOptions{}
	if options != nil {
		cp = *options
	}
	cp.init()
	msalOpts := publicClientOptions{
		AdditionallyAllowedTenants:     cp.AdditionallyAllowedTenants,
		ClientOptions:                  cp.ClientOptions,
		DeviceCodePrompt:               cp.UserPrompt,
		DisableAutomaticAuthentication: cp.DisableAutomaticAuthentication,
		DisableInstanceDiscovery:       cp.DisableInstanceDiscovery,
		TokenCachePersistenceOptions:   cp.TokenCachePersistenceOptions,
	}
	c, err := newPublicClient(cp.TenantID, cp.ClientID, credNameDeviceCode, msalOpts)
	if err != nil {
		return nil, err
	}
	c.name = credNameDeviceCode
	return &DeviceCodeCredential{client: c}, nil
}

// Authenticate a user via the device code flow. Subsequent calls to GetToken will automatically use the returned AuthenticationRecord.
func (c *DeviceCodeCredential) Authenticate(ctx context.Context, opts *policy.TokenRequestOptions) (AuthenticationRecord, error) {
	return c.client.Authenticate(ctx, opts)
}

// GetToken requests an access token from Azure Active Directory. It will begin the device code flow and poll until the user completes authentication.
// This method is called automatically by Azure SDK clients.
func (c *DeviceCodeCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return c.client.GetToken(ctx, opts)
}

var _ azcore.TokenCredential = (*DeviceCodeCredential)(nil)
