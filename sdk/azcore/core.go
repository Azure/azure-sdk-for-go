// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"io"
	"net/http"
)

// Policy represents an extensibility point for the Pipeline that can mutate the specified
// Request and react to the received Response.
type Policy interface {
	// Do applies the policy to the specified Request.  When implementing a Policy, mutate the
	// request before calling req.Next() to move on to the next policy, and respond to the result
	// before returning to the caller.
	Do(req *Request) (*Response, error)
}

// PolicyFunc is a type that implements the Policy interface.
// Use this type when implementing a stateless policy as a first-class function.
type PolicyFunc func(*Request) (*Response, error)

// Do implements the Policy interface on PolicyFunc.
func (pf PolicyFunc) Do(req *Request) (*Response, error) {
	return pf(req)
}

// Transport represents an HTTP pipeline transport used to send HTTP requests and receive responses.
type Transport interface {
	// Do sends the HTTP request and returns the HTTP response or error.
	Do(req *http.Request) (*http.Response, error)
}

// TransportFunc is a type that implements the Transport interface.
// Use this type when implementing a stateless transport as a first-class function.
type TransportFunc func(*http.Request) (*http.Response, error)

// Do implements the Transport interface on TransportFunc.
func (tf TransportFunc) Do(req *http.Request) (*http.Response, error) {
	return tf(req)
}

// used to adapt a TransportPolicy to a Policy
type transportPolicy struct {
	trans Transport
}

func (tp transportPolicy) Do(req *Request) (*Response, error) {
	resp, err := tp.trans.Do(req.Request)
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

// NewPipeline creates a new Pipeline object from the specified Transport and Policies.
// If no transport is provided then the default *http.Client transport will be used.
func NewPipeline(transport Transport, policies ...Policy) Pipeline {
	if transport == nil {
		transport = defaultHTTPClient
	}
	// transport policy must always be the last in the slice
	policies = append(policies, PolicyFunc(httpHeaderPolicy), PolicyFunc(bodyDownloadPolicy), transportPolicy{trans: transport})
	return Pipeline{
		policies: policies,
	}
}

// Do is called for each and every HTTP request. It passes the request through all
// the Policy objects (which can transform the Request's URL/query parameters/headers)
// and ultimately sends the transformed HTTP request over the network.
func (p Pipeline) Do(req *Request) (*Response, error) {
	if err := req.valid(); err != nil {
		return nil, err
	}
	req.policies = p.policies
	return req.Next()
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
