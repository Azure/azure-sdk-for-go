//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	ContentTypeAppJSON = "application/json"
	ContentTypeAppXML  = "application/xml"
)

const (
	HeaderAzureAsync        = "Azure-AsyncOperation"
	HeaderContentLength     = "Content-Length"
	HeaderContentType       = "Content-Type"
	HeaderLocation          = "Location"
	HeaderOperationLocation = "Operation-Location"
	HeaderRetryAfter        = "Retry-After"
	HeaderUserAgent         = "User-Agent"
)

const (
	DefaultMaxRetries = 3
)

const (
	// Module is the name of the calling module used in telemetry data.
	Module = "azcore"

	// Version is the semantic version (see http://semver.org) of this module.
	Version = "v0.19.0"
)

// CtxWithHTTPHeaderKey is used as a context key for adding/retrieving http.Header.
type CtxWithHTTPHeaderKey struct{}

// CtxWithRetryOptionsKey is used as a context key for adding/retrieving RetryOptions.
type CtxWithRetryOptionsKey struct{}

type nopCloser struct {
	io.ReadSeeker
}

func (n nopCloser) Close() error {
	return nil
}

// NopCloser returns a ReadSeekCloser with a no-op close method wrapping the provided io.ReadSeeker.
func NopCloser(rs io.ReadSeeker) io.ReadSeekCloser {
	return nopCloser{rs}
}

// BodyDownloadPolicyOpValues is the struct containing the per-operation values
type BodyDownloadPolicyOpValues struct {
	Skip bool
}

// PolicyFunc is a type that implements the Policy interface.
// Use this type when implementing a stateless policy as a first-class function.
type PolicyFunc func(*Request) (*http.Response, error)

// Do implements the Policy interface on PolicyFunc.
func (pf PolicyFunc) Do(req *Request) (*http.Response, error) {
	return pf(req)
}

func NewResponseError(inner error, resp *http.Response) error {
	return &ResponseError{inner: inner, resp: resp}
}

type ResponseError struct {
	inner error
	resp  *http.Response
}

// Error implements the error interface for type ResponseError.
func (e *ResponseError) Error() string {
	return e.inner.Error()
}

// Unwrap returns the inner error.
func (e *ResponseError) Unwrap() error {
	return e.inner
}

// RawResponse returns the HTTP response associated with this error.
func (e *ResponseError) RawResponse() *http.Response {
	return e.resp
}

// NonRetriable indicates this error is non-transient.
func (e *ResponseError) NonRetriable() {
	// marker method
}

// Delay waits for the duration to elapse or the context to be cancelled.
func Delay(ctx context.Context, delay time.Duration) error {
	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ErrNoBody is returned if the response didn't contain a body.
var ErrNoBody = errors.New("the response did not contain a body")

// GetJSON reads the response body into a raw JSON object.
// It returns ErrNoBody if there was no content.
func GetJSON(resp *http.Response) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, ErrNoBody
	}
	// put the body back so it's available to others
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	// unmarshall the body to get the value
	var jsonBody map[string]interface{}
	if err = json.Unmarshal(body, &jsonBody); err != nil {
		return nil, err
	}
	return jsonBody, nil
}
