// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	defaultAuthorityHost = "https://login.microsoftonline.com/"
)

var (
	// StatusCodesForRetry is the default set of HTTP status code for which the policy will retry.
	StatusCodesForRetry = [6]int{
		http.StatusRequestTimeout,      // 408
		http.StatusTooManyRequests,     // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout,      // 504
	}

	successStatusCodes = [2]int{
		http.StatusOK,      // 200
		http.StatusCreated, // 201
	}
)

var (
	defaultAuthorityHostURL   *url.URL
	defaultIdentityClientOpts *IdentityClientOptions
)

func init() {
	var err error
	defaultAuthorityHostURL, err = url.Parse(defaultAuthorityHost)
	if err != nil {
		// TODO fix this
		panic(err)
	}

	defaultIdentityClientOpts = &IdentityClientOptions{AuthorityHost: defaultAuthorityHostURL}
}

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
	authFailed := &AuthenticationFailedError{}
	if resp == nil {
		authFailed.Message = "Something unexpected happened when attempting to authenticate"
		return authFailed
	}

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
	CredentialType string
	Message        string
}

func (e *CredentialUnavailableError) Error() string {
	return e.CredentialType + ": " + e.Message
}

// ChainedCredentialError ..
type ChainedCredentialError struct {
	ErrorList []*CredentialUnavailableError
	Message   string
}

func (e *ChainedCredentialError) Error() string {
	if len(e.ErrorList) > 0 {
		msg := ""
		for _, err := range e.ErrorList {
			msg += err.Error() + "\n"
		}
		e.Message = msg
	}
	return e.Message
}

// IdentityClientOptions to configure requests made to Azure Identity Services
type IdentityClientOptions struct {
	PipelineOptions azcore.PipelineOptions
	AuthorityHost   *url.URL // The host of the Azure Active Directory authority. The default is https://login.microsoft.com
}

// TODO singleton default options?
// TODO this is unnecessary if we keep the functionality in the init
// NewIdentityClientOptions initializes an instance of IdentityClientOptions with default settings
func (c *IdentityClientOptions) setDefaultValues() *IdentityClientOptions {
	// Invert logic if c!= nil
	if c == nil {
		c = &IdentityClientOptions{}
		c.AuthorityHost = defaultAuthorityHostURL
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

// hasStatusCode returns true if the Response's status code is one of the specified values.
func hasStatusCode(r *azcore.Response, statusCodes ...int) bool {
	if r == nil {
		return false
	}
	for _, sc := range statusCodes {
		if r.StatusCode == sc {
			return true
		}
	}
	return false
}

// // MSIPipelineOptions is used to configure a request policy pipeline's retry policy and logging.
// type MSIPipelineOptions struct {
// 	// Retry configures the built-in retry policy behavior.
// 	Retry RetryOptions

// 	// Telemetry configures the built-in telemetry policy behavior.
// 	Telemetry azcore.TelemetryOptions

// 	// HTTPClient sets the transport for making HTTP requests.
// 	// Leave this as nil to use the default HTTP transport.
// 	HTTPClient azcore.Transport

// 	// LogOptions configures the built-in request logging policy behavior.
// 	LogOptions azcore.RequestLogOptions
// }

// NewDefaultMSIPipeline creates a Pipeline using the specified pipeline options
func newDefaultMSIPipeline(o azcore.PipelineOptions) azcore.Pipeline {
	if o.HTTPClient == nil {
		o.HTTPClient = azcore.DefaultHTTPClientTransport()
	}

	return azcore.NewPipeline(
		o.HTTPClient,
		azcore.NewTelemetryPolicy(o.Telemetry),
		azcore.NewUniqueRequestIDPolicy(),
		NewMSIRetryPolicy(o.Retry),
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
