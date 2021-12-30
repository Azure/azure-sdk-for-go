// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	msal "github.com/AzureAD/microsoft-authentication-library-for-go/apps/errors"
)

// AuthenticationFailedError indicates an authentication request has failed.
type AuthenticationFailedError interface {
	errorinfo.NonRetriable
	RawResponse() *http.Response
	authenticationFailed()
}

type authenticationFailedError struct {
	error
	resp *http.Response
}

func newAuthenticationFailedError(err error, resp *http.Response) AuthenticationFailedError {
	if resp == nil {
		var e msal.CallErr
		if errors.As(err, &e) {
			return authenticationFailedError{err, e.Resp}
		}
	}
	return authenticationFailedError{err, resp}
}

// NonRetriable indicates that this error should not be retried.
func (authenticationFailedError) NonRetriable() {
	// marker method
}

// AuthenticationFailed indicates that an authentication attempt failed
func (authenticationFailedError) authenticationFailed() {
	// marker method
}

// RawResponse returns the HTTP response motivating the error, if available.
func (e authenticationFailedError) RawResponse() *http.Response {
	return e.resp
}

var _ AuthenticationFailedError = (*authenticationFailedError)(nil)
var _ errorinfo.NonRetriable = (*authenticationFailedError)(nil)

// CredentialUnavailableError indicates a credential can't attempt authenticate
// because it lacks required data or state.
type CredentialUnavailableError interface {
	errorinfo.NonRetriable
	credentialUnavailable()
}

type credentialUnavailableError struct {
	credType string
	message  string
}

func newCredentialUnavailableError(credType, message string) CredentialUnavailableError {
	return credentialUnavailableError{credType: credType, message: message}
}

func (e credentialUnavailableError) Error() string {
	return e.credType + ": " + e.message
}

// NonRetriable indicates that this error should not be retried.
func (e credentialUnavailableError) NonRetriable() {
	// marker method
}

// CredentialUnavailable indicates that the credential didn't attempt to authenticate
func (e credentialUnavailableError) credentialUnavailable() {
	// marker method
}

var _ CredentialUnavailableError = (*credentialUnavailableError)(nil)
var _ errorinfo.NonRetriable = (*credentialUnavailableError)(nil)
