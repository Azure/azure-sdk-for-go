// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"errors"
	"io"
	"net/http"
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

// Do is called for each and every HTTP request. It passes the Context and request through
// all the Policy objects (which can transform the Request's URL/query parameters/headers)
// and ultimately sends the transformed HTTP request over the network.
func (p Pipeline) Do(ctx context.Context, req *Request) (*Response, error) {
	req.policies = p.policies
	return req.Next(ctx)
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

// IterationDone is returned by an iterator's Next method when iteration is complete.
var IterationDone = errors.New("no more items in iterator")

// Retrier provides methods describing if an error should be considered as transient.
type Retrier interface {
	// IsNotRetriable returns true for error types that are not retriable.
	IsNotRetriable() bool
}
