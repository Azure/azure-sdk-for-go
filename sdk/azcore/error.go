//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"net/http"
)

var (
	// StackFrameCount contains the number of stack frames to include when a trace is being collected.
	StackFrameCount = 32
)

// HTTPResponse provides access to an HTTP response when available.
// Errors returned from failed API calls will implement this interface.
// Use errors.As() to access this interface in the error chain.
// If there was no HTTP response then this interface will be omitted
// from any error in the chain.
type HTTPResponse interface {
	RawResponse() *http.Response
}

// NonRetriableError represents a non-transient error.  This works in
// conjunction with the retry policy, indicating that the error condition
// is idempotent, so no retries will be attempted.
// Use errors.As() to access this interface in the error chain.
type NonRetriableError interface {
	error
	NonRetriable()
}

// NewResponseError wraps the specified error with an error that provides access to an HTTP response.
// If an HTTP request returns a non-successful status code, wrap the response and the associated error
// in this error type so that callers can access the underlying *http.Response as required.
// DO NOT wrap failed HTTP requests that returned an error and no response with this type.
func NewResponseError(inner error, resp *http.Response) error {
	return &responseError{inner: inner, resp: resp}
}

type responseError struct {
	inner error
	resp  *http.Response
}

// Error implements the error interface for type ResponseError.
func (e *responseError) Error() string {
	return e.inner.Error()
}

// Unwrap returns the inner error.
func (e *responseError) Unwrap() error {
	return e.inner
}

// RawResponse returns the HTTP response associated with this error.
func (e *responseError) RawResponse() *http.Response {
	return e.resp
}

// NonRetriable indicates this error is non-transient.
func (e *responseError) NonRetriable() {
	// marker method
}

var _ HTTPResponse = (*responseError)(nil)
var _ NonRetriableError = (*responseError)(nil)
