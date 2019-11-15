// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Policy represents an extensibility point for the Pipeline that can mutate the specified
// Request and react to the received Response.
type Policy interface {
	// Do applies the policy to the specified Request.  When implementing a Policy, mutate the
	// request before calling req.Do() to move on to the next policy, and respond to the result
	// before returning to the caller.
	Do(ctx context.Context, req *Request) (*Response, error)
}

// PolicyFunc is a type that implements the Policy interface.
// Use this type when implementing a stateless policy as a first-class function.
type PolicyFunc func(context.Context, *Request) (*Response, error)

// Do implements the Policy interface on PolicyFunc.
func (pf PolicyFunc) Do(ctx context.Context, req *Request) (*Response, error) {
	return pf(ctx, req)
}

// Transport represents an HTTP pipeline transport used to send HTTP requests and receive responses.
type Transport interface {
	// Do sends the HTTP request and returns the HTTP response or error.
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// transportFunc is a type that implements the Transport interface.
// Use this type when implementing a stateless transport as a first-class function.
type transportFunc func(context.Context, *http.Request) (*http.Response, error)

// Do implements the Transport interface on transportFunc.
func (tf transportFunc) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return tf(ctx, req)
}

// used to adapt a TransportPolicy to a Policy
type transportPolicy struct {
	trans Transport
}

func (tp transportPolicy) Do(ctx context.Context, req *Request) (*Response, error) {
	resp, err := tp.trans.Do(ctx, req.Request)
	if err != nil {
		return nil, err
	}
	return &Response{Response: resp}, nil
}

// Pipeline represents a primitive for sending HTTP requests and receiving responses.
// Its behavior can be extended by specifying policies during construction.
type Pipeline struct {
	policies []Policy
}

// NewPipeline creates a new goroutine-safe Pipeline object from the specified Policies.
// If no transport is provided then the default HTTP transport will be used.
func NewPipeline(transport Transport, policies ...Policy) Pipeline {
	if transport == nil {
		transport = DefaultHTTPClientTransport()
	}
	// transport policy must always be the last in the slice
	policies = append(policies, newBodyDownloadPolicy(), transportPolicy{trans: transport})
	return Pipeline{
		policies: policies,
	}
}

// NewRequest creates a new Request associated with this pipeline.
func (p Pipeline) NewRequest(httpMethod string, URL url.URL) *Request {
	// removeEmptyPort strips the empty port in ":port" to ""
	// as mandated by RFC 3986 Section 6.2.3.
	// adapted from removeEmptyPort() in net/http.go
	if strings.LastIndex(URL.Host, ":") > strings.LastIndex(URL.Host, "]") {
		URL.Host = strings.TrimSuffix(URL.Host, ":")
	}
	return &Request{
		Request: &http.Request{
			Method:     httpMethod,
			URL:        &URL,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     http.Header{},
			Host:       URL.Host,
		},
		policies: p.policies,
	}
}

// PipelineOptions is used to configure a request policy pipeline's retry policy and logging.
type PipelineOptions struct {
	// Retry configures the built-in retry policy behavior.
	Retry RetryOptions

	// Telemetry configures the built-in telemetry policy behavior.
	Telemetry TelemetryOptions

	// HTTPClient sets the transport for making HTTP requests.
	// Leave this as nil to use the default HTTP transport.
	HTTPClient Transport

	// LogOptions configures the built-in request logging policy behavior.
	LogOptions RequestLogOptions
}

// ReadSeekCloser is the interface that groups the io.ReadCloser and io.Seeker interfaces.
type ReadSeekCloser interface {
	io.ReadCloser
	io.Seeker
}

type nopCloser struct {
	io.ReadSeeker
}

func (n nopCloser) Close() error {
	return nil
}

// NopCloser returns a ReadSeekCloser with a no-op close method wrapping the provided io.ReadSeeker.
func NopCloser(rs io.ReadSeeker) ReadSeekCloser {
	return nopCloser{rs}
}
