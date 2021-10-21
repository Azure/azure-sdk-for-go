// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
)

// AuthorityHost is the base URL for Azure Active Directory
type AuthorityHost string

const (
	// AzureChina is a global constant to use in order to access the Azure China cloud.
	AzureChina AuthorityHost = "https://login.chinacloudapi.cn/"
	// AzureGermany is a global constant to use in order to access the Azure Germany cloud.
	AzureGermany AuthorityHost = "https://login.microsoftonline.de/"
	// AzureGovernment is a global constant to use in order to access the Azure Government cloud.
	AzureGovernment AuthorityHost = "https://login.microsoftonline.us/"
	// AzurePublicCloud is a global constant to use in order to access the Azure public cloud.
	AzurePublicCloud AuthorityHost = "https://login.microsoftonline.com"
)

// defaultSuffix is the default AADv2 scope
const defaultSuffix = "/.default"

const (
	headerUserAgent   = "User-Agent"
	headerURLEncoded  = "application/x-www-form-urlencoded"
	headerMetadata    = "Metadata"
	headerContentType = "Content-Type"
)

const tenantIDValidationErr = "Invalid tenantID provided. You can locate your tenantID by following the instructions listed here: https://docs.microsoft.com/partner-center/find-ids-and-domain-names."

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

// AuthenticationFailedError is returned when the authentication request has failed.
type AuthenticationFailedError struct {
	inner error
	msg   string
	resp  *http.Response
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
	return e.msg
}

// RawResponse returns the HTTP response motivating the error, if available
func (e *AuthenticationFailedError) RawResponse() *http.Response {
	return e.resp
}

var _ azcore.HTTPResponse = (*AuthenticationFailedError)(nil)
var _ errorinfo.NonRetriable = (*AuthenticationFailedError)(nil)

// CredentialUnavailableError is the error type returned when the conditions required to
// create a credential do not exist or are unavailable.
type CredentialUnavailableError struct {
	// CredentialType holds the name of the credential that is unavailable
	credentialType string
	// Message contains the reason why the credential is unavailable
	message string
}

func (e *CredentialUnavailableError) Error() string {
	return e.credentialType + ": " + e.message
}

// NonRetriable indicates that this error should not be retried.
func (e *CredentialUnavailableError) NonRetriable() {
	// marker method
}

var _ errorinfo.NonRetriable = (*CredentialUnavailableError)(nil)

// setAuthorityHost initializes the authority host for credentials.
func setAuthorityHost(authorityHost AuthorityHost) (string, error) {
	host := string(authorityHost)
	if host == "" {
		host = string(AzurePublicCloud)
		if envAuthorityHost := os.Getenv("AZURE_AUTHORITY_HOST"); envAuthorityHost != "" {
			host = envAuthorityHost
		}
	}
	u, err := url.Parse(host)
	if err != nil {
		return "", err
	}
	if u.Scheme != "https" {
		return "", errors.New("cannot use an authority host without https")
	}
	return host, nil
}

// newDefaultMSIPipeline creates a pipeline using the specified pipeline options needed
// for a Managed Identity, such as a MSI specific retry policy.
func newDefaultMSIPipeline(o ManagedIdentityCredentialOptions) runtime.Pipeline {
	cp := o.ClientOptions
	v := reflect.ValueOf(cp.Retry)
	// if the caller hasn't set retry options, set defaults appropriate for IMDS
	if !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface()) {
		cp.Retry = policy.RetryOptions{
			MaxRetries:    5,
			MaxRetryDelay: 1 * time.Minute,
			RetryDelay:    2 * time.Second,
			TryTimeout:    1 * time.Minute,
			StatusCodes: []int{
				// The following status codes are a subset of those found in azcore.StatusCodesForRetry, these are the only ones specifically needed for IMDS
				http.StatusRequestTimeout,      // 408
				http.StatusTooManyRequests,     // 429
				http.StatusInternalServerError, // 500
				http.StatusBadGateway,          // 502
				http.StatusGatewayTimeout,      // 504
				http.StatusNotFound,            // 404
				http.StatusGone,                // 410
				// all remaining 5xx
				http.StatusNotImplemented,                // 501
				http.StatusHTTPVersionNotSupported,       // 505
				http.StatusVariantAlsoNegotiates,         // 506
				http.StatusInsufficientStorage,           // 507
				http.StatusLoopDetected,                  // 508
				http.StatusNotExtended,                   // 510
				http.StatusNetworkAuthenticationRequired, // 511
			},
		}
	}
	return runtime.NewPipeline(component, version, nil, nil, &cp)
}

// validTenantID return true is it receives a valid tenantID, returns false otherwise
func validTenantID(tenantID string) bool {
	match, err := regexp.MatchString("^[0-9a-zA-Z-.]+$", tenantID)
	if err != nil {
		return false
	}
	return match
}

// tokenEndpoint takes a given path and appends "/token" to the end of the path
func tokenEndpoint(p string) string {
	return path.Join(p, "/token")
}

// oauthPath returns the oauth path for AAD or ADFS
func oauthPath(tenantID string) string {
	if tenantID == "adfs" {
		return "/oauth2"
	}
	return "/oauth2/v2.0"
}
