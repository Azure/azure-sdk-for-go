// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	defaultAuthorityHost = "https://login.microsoftonline.com/"
)

// IdentityClientOptions to configure requests made to Azure Identity Services
type IdentityClientOptions struct {
	pipelineOptions azcore.PipelineOptions
	AuthorityHost   url.URL // The host of the Azure Active Directory authority. The default is https://login.microsoft.com
}

// // ClientOptions is used to ...
// type ClientOptions struct {
// 	Transport azcore.Pipeline
// }

// newIdentityClientOptions function initializes the default Authority Host to authenticate against
func newIdentityClientOptions() (*IdentityClientOptions, error) {
	var c *IdentityClientOptions = &IdentityClientOptions{}
	res, err := url.Parse(defaultAuthorityHost)
	if err != nil {
		return nil, fmt.Errorf("Parse: %w", err)
	}
	c.AuthorityHost = *res
	return c, nil
}

// TokenCredential interface serves as an anonymous field for all other credentials to use the GetToken method
// CP: check this
type TokenCredential interface {
	GetToken(ctx context.Context, scopes []string) (*AccessToken, error)
}

// // ServicePrincipalCertificateSecret implements ServicePrincipalSecret for generic RSA cert auth with signed JWTs.
// type ServicePrincipalCertificateSecret struct {
// 	Certificate *x509.Certificate
// 	PrivateKey  *rsa.PrivateKey
// }

// NewDefaultPipeline creates a Pipeline using the specified credentials and options.
func NewDefaultPipeline(o azcore.PipelineOptions) azcore.Pipeline {
	if o.HTTPClient == nil {
		o.HTTPClient = azcore.DefaultHTTPClientPolicy()
	}
	// Closest to API goes first; closest to the wire goes last
	// NOTE: The credential's policy factory must appear close to the wire so it can sign any
	// changes made by other factories (like UniqueRequestIDPolicyFactory)
	return azcore.NewPipeline(
		// azcore.NewTelemetryPolicy(o.Telemetry),
		// azcore.NewUniqueRequestIDPolicy(),
		// azcore.NewRetryPolicy(o.Retry),
		// azcore.NewBodyDownloadPolicy(),
		// azcore.NewRequestLogPolicy(o.RequestLog),
		o.HTTPClient)
}

/*********************************************************************
func newDeviceCodeInfo(dc DeviceCodeResult) DeviceCodeInfo {
	return DeviceCodeInfo{UserCode: dc.UserCode, DeviceCode: dc.DeviceCode, VerificationUrl: dc.VerificationUri,
		ExpiresOn: dc.ExpiresOn, Internal: dc.Interval, Message: dc.Message, ClientId: dc.ClientId, Scopes: dc.Scopes}
}

// Details of the device code to present to a user to allow them to authenticate through the device code authentication flow.
type DeviceCodeInfo struct {
	// JMR: Make all these private and add public getter methods?
	UserCode        string        // User code returned by the service
	DeviceCode      string        // Device code returned by the service
	VerificationUrl string        // Verification URL where the user must navigate to authenticate using the device code and credentials. JMR: URL?
	ExpiresOn       time.Duration // Time when the device code will expire.
	Interval        int64         // Polling interval time to check for completion of authentication flow.
	Message         string        // User friendly text response that can be used for display purpose.
	ClientId        string        // Identifier of the client requesting device code.
	Scopes          []string      // List of the scopes that would be held by token. JMR: Should be readonly
}
**************************************************************/

/*******************************************************
const defaultSuffix = "/.default"

func scopesToResource(scope string) string {
	if strings.EndsWith(scope, defaultSuffix) {
		return scope[:len(scope)-len(defaultSuffix)]
	}
	return scope
}

func resourceToScope(resource string) string {
	resource += defaultSuffix
	return resource
}
**********************************************************/
