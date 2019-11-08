// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"encoding/json"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	defaultAuthorityHost = "https://login.microsoftonline.com/"
)

// AuthenticationFailedError is a struct used to marshal responses when authentication has failed
type AuthenticationFailedError struct {
	Message       string `json:"error"`
	Description   string `json:"error_description"`
	Timestamp     string `json:"timestamp"`
	TraceID       string `json:"trace_id"`
	CorrelationID string `json:"correlation_id"`
	URI           string `json:"error_uri"`
	Response      *azcore.Response
}

func (e *AuthenticationFailedError) Error() string {
	return e.Message + ": " + e.Description
}

func newAuthenticationFailedError(resp *azcore.Response) error {
	var authFailed *AuthenticationFailedError
	err := json.Unmarshal(resp.Payload, authFailed)
	if err != nil {
		authFailed.Message = resp.Status
		authFailed.Description = "Failed to unmarshal response: " + err.Error()
	}
	authFailed.Response = resp
	return authFailed
}

// CredentialUnavailableError an error type returned when the conditions required to create a credential do not exist
type CredentialUnavailableError struct {
	Message string
}

func (e *CredentialUnavailableError) Error() string {
	return e.Message
}

// AggregateError a specific error type for the chained token credential
type AggregateError struct {
	Msg     string
	Err     error
	ErrList []error
}

func (e *AggregateError) Error() string {
	errString := ""
	for _, err := range e.ErrList {
		errString += "\n" + err.Error()
	}

	return errString
}

// IdentityClientOptions to configure requests made to Azure Identity Services
type IdentityClientOptions struct {
	PipelineOptions azcore.PipelineOptions
	AuthorityHost   *url.URL // The host of the Azure Active Directory authority. The default is https://login.microsoft.com
}

// NewIdentityClientOptions initializes an instance of IdentityClientOptions with default settings
func (c *IdentityClientOptions) setDefaultValues() *IdentityClientOptions {
	if c == nil {
		c = &IdentityClientOptions{}
		c.AuthorityHost, _ = url.Parse(defaultAuthorityHost)
	}
	return c
}

// NewDefaultPipeline creates a Pipeline using the specified pipeline options
func newDefaultPipeline(o azcore.PipelineOptions) azcore.Pipeline {
	if o.HTTPClient == nil {
		o.HTTPClient = azcore.DefaultHTTPClientTransport()
	}

	return azcore.NewPipeline(
		o.HTTPClient,
		azcore.NewTelemetryPolicy(o.Telemetry),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(o.Retry),
		azcore.NewRequestLogPolicy(o.LogOptions))
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
