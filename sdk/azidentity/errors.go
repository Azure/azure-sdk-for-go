// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
)

// AuthenticationFailedError indicates an authentication request has failed.
type AuthenticationFailedError interface {
	azcore.HTTPResponse
	errorinfo.NonRetriable
	AuthenticationFailed()
}

type authenticationFailedError struct {
	error
	resp *http.Response
}

func newAuthenticationFailedError(err error, resp *http.Response) AuthenticationFailedError {
	return authenticationFailedError{err, resp}
}

// NonRetriable indicates that this error should not be retried.
func (authenticationFailedError) NonRetriable() {
	// marker method
}

// AuthenticationFailed indicates that an authentication attempt failed
func (authenticationFailedError) AuthenticationFailed() {
	// marker method
}

// RawResponse returns the HTTP response motivating the error, if available.
func (e authenticationFailedError) RawResponse() *http.Response {
	return e.resp
}

var _ AuthenticationFailedError = (*authenticationFailedError)(nil)
var _ azcore.HTTPResponse = (*authenticationFailedError)(nil)
var _ errorinfo.NonRetriable = (*authenticationFailedError)(nil)

// CredentialUnavailableError indicates a credential can't attempt authenticate
// because it lacks required data or state.
type CredentialUnavailableError interface {
	errorinfo.NonRetriable
	CredentialUnavailable()
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
func (e credentialUnavailableError) CredentialUnavailable() {
	// marker method
}

var _ CredentialUnavailableError = (*credentialUnavailableError)(nil)
var _ errorinfo.NonRetriable = (*credentialUnavailableError)(nil)
