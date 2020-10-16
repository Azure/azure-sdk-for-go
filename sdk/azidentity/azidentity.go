// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	// AzureChina is a global constant to use in order to access the Azure China cloud.
	AzureChina = "https://login.chinacloudapi.cn/"
	// AzureGermany is a global constant to use in order to access the Azure Germany cloud.
	AzureGermany = "https://login.microsoftonline.de/"
	// AzureGovernment is a global constant to use in order to access the Azure Government cloud.
	AzureGovernment = "https://login.microsoftonline.us/"
	// AzurePublicCloud is a global constant to use in order to access the Azure public cloud.
	AzurePublicCloud = "https://login.microsoftonline.com/"
	// defaultSuffix is a suffix the signals that a string is in scope format
	defaultSuffix = "/.default"
)

var (
	successStatusCodes = [2]int{
		http.StatusOK,      // 200
		http.StatusCreated, // 201
	}
)

type tokenResponse struct {
	token        *azcore.AccessToken
	refreshToken string
}

// AADAuthenticationFailedError is used to unmarshal error responses received from Azure Active Directory.
type AADAuthenticationFailedError struct {
	Message       string `json:"error"`
	Description   string `json:"error_description"`
	Timestamp     string `json:"timestamp"`
	TraceID       string `json:"trace_id"`
	CorrelationID string `json:"correlation_id"`
	URI           string `json:"error_uri"`
	Response      *azcore.Response
}

func (e *AADAuthenticationFailedError) Error() string {
	msg := e.Message
	if len(e.Description) > 0 {
		msg += " " + e.Description
	}
	return msg
}

// AuthenticationFailedError is returned when the authentication request has failed.
type AuthenticationFailedError struct {
	inner error
	msg   string
}

// Unwrap method on AuthenticationFailedError provides access to the inner error if available.
func (e *AuthenticationFailedError) Unwrap() error {
	return e.inner
}

// NonRetriable indicates that this error should not be retried.
func (e *AuthenticationFailedError) NonRetriable() {
	// marker method
}

func (e *AuthenticationFailedError) Error() string {
	if len(e.msg) == 0 {
		e.msg = e.inner.Error()
	}
	return e.msg
}

var _ azcore.NonRetriableError = (*AuthenticationFailedError)(nil)

func newAADAuthenticationFailedError(resp *azcore.Response) error {
	authFailed := &AADAuthenticationFailedError{Response: resp}
	err := resp.UnmarshalAsJSON(authFailed)
	if err != nil {
		authFailed.Message = resp.Status
		authFailed.Description = "Failed to unmarshal response: " + err.Error()
	}
	return authFailed
}

// CredentialUnavailableError is the error type returned when the conditions required to
// create a credential do not exist or are unavailable.
type CredentialUnavailableError struct {
	// CredentialType holds the name of the credential that is unavailable
	CredentialType string
	// Message contains the reason why the credential is unavailable
	Message string
}

func (e *CredentialUnavailableError) Error() string {
	return e.CredentialType + ": " + e.Message
}

// NonRetriable indicates that this error should not be retried.
func (e *CredentialUnavailableError) NonRetriable() {
	// marker method
}

var _ azcore.NonRetriableError = (*CredentialUnavailableError)(nil)

// TokenCredentialOptions are used to configure how requests are made to Azure Active Directory.
type TokenCredentialOptions struct {
	// The host of the Azure Active Directory authority. The default is https://login.microsoft.com
	AuthorityHost string

	// HTTPClient sets the transport for making HTTP requests
	// Leave this as nil to use the default HTTP transport
	HTTPClient azcore.Transport

	// Retry configures the built-in retry policy behavior
	Retry *azcore.RetryOptions

	// Telemetry configures the built-in telemetry policy behavior
	Telemetry azcore.TelemetryOptions
}

// setDefaultValues initializes an instance of TokenCredentialOptions with default settings.
func (c *TokenCredentialOptions) setDefaultValues() (*TokenCredentialOptions, error) {
	authorityHost := AzurePublicCloud
	if envAuthorityHost := os.Getenv("AZURE_AUTHORITY_HOST"); envAuthorityHost != "" {
		authorityHost = envAuthorityHost
	}

	if c == nil {
		c = &TokenCredentialOptions{AuthorityHost: authorityHost}
	}

	if c.AuthorityHost == "" {
		c.AuthorityHost = authorityHost
	}

	s, err := url.Parse(c.AuthorityHost)
	if err != nil {
		return nil, err
	}

	if s.Scheme != "https" {
		return nil, errors.New("cannot use an authority host without https")
	}
	return c, nil
}

// newDefaultPipeline creates a pipeline using the specified pipeline options.
func newDefaultPipeline(o TokenCredentialOptions) azcore.Pipeline {
	return azcore.NewPipeline(
		o.HTTPClient,
		azcore.NewTelemetryPolicy(&o.Telemetry),
		azcore.NewRetryPolicy(o.Retry),
		azcore.NewLogPolicy(nil))
}

// newDefaultMSIPipeline creates a pipeline using the specified pipeline options needed
// for a Managed Identity, such as a MSI specific retry policy.
func newDefaultMSIPipeline(o ManagedIdentityCredentialOptions) azcore.Pipeline {
	var statusCodes []int
	// retry policy for MSI is not end-user configurable
	retryOpts := azcore.RetryOptions{
		MaxRetries:    5,
		MaxRetryDelay: 1 * time.Minute,
		RetryDelay:    2 * time.Second,
		TryTimeout:    1 * time.Minute,
		StatusCodes: append(statusCodes,
			// The following status codes are a subset of those found in azcore.StatusCodesForRetry, these are the only ones specifically needed for MSI scenarios
			http.StatusRequestTimeout,      // 408
			http.StatusTooManyRequests,     // 429
			http.StatusInternalServerError, // 500
			http.StatusBadGateway,          // 502
			http.StatusGatewayTimeout,      // 504
			http.StatusNotFound,            //404
			http.StatusGone,                //410
			// all remaining 5xx
			http.StatusNotImplemented,                 // 501
			http.StatusHTTPVersionNotSupported,        // 505
			http.StatusVariantAlsoNegotiates,          // 506
			http.StatusInsufficientStorage,            // 507
			http.StatusLoopDetected,                   // 508
			http.StatusNotExtended,                    // 510
			http.StatusNetworkAuthenticationRequired), // 511
	}

	return azcore.NewPipeline(
		o.HTTPClient,
		azcore.NewTelemetryPolicy(&o.Telemetry),
		azcore.NewRetryPolicy(&retryOpts),
		azcore.NewLogPolicy(nil))
}
