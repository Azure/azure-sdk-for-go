//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"errors"
	"io"
	"net/http"
	"reflect"
)

const (
	headerContentLength     = "Content-Length"
	headerContentType       = "Content-Type"
	headerOperationLocation = "Operation-Location"
	headerLocation          = "Location"
	headerRetryAfter        = "Retry-After"
	headerUserAgent         = "User-Agent"
)

// Policy represents an extensibility point for the Pipeline that can mutate the specified
// Request and react to the received Response.
type Policy interface {
	// Do applies the policy to the specified Request.  When implementing a Policy, mutate the
	// request before calling req.Next() to move on to the next policy, and respond to the result
	// before returning to the caller.
	Do(req *Request) (*http.Response, error)
}

// policyFunc is a type that implements the Policy interface.
// Use this type when implementing a stateless policy as a first-class function.
type policyFunc func(*Request) (*http.Response, error)

// Do implements the Policy interface on PolicyFunc.
func (pf policyFunc) Do(req *Request) (*http.Response, error) {
	return pf(req)
}

// Transporter represents an HTTP pipeline transport used to send HTTP requests and receive responses.
type Transporter interface {
	// Do sends the HTTP request and returns the HTTP response or error.
	Do(req *http.Request) (*http.Response, error)
}

// used to adapt a TransportPolicy to a Policy
type transportPolicy struct {
	trans Transporter
}

func (tp transportPolicy) Do(req *Request) (*http.Response, error) {
	resp, err := tp.trans.Do(req.Request)
	if err != nil {
		return nil, err
	} else if resp == nil {
		// there was no response and no error (rare but can happen)
		// this ensures the retry policy will retry the request
		return nil, errors.New("received nil response")
	}
	return resp, nil
}

// Pipeline represents a primitive for sending HTTP requests and receiving responses.
// Its behavior can be extended by specifying policies during construction.
type Pipeline struct {
	policies []Policy
}

// NewPipeline creates a new Pipeline object from the specified Transport and Policies.
// If no transport is provided then the default *http.Client transport will be used.
func NewPipeline(transport Transporter, policies ...Policy) Pipeline {
	if transport == nil {
		transport = defaultHTTPClient
	}
	// transport policy must always be the last in the slice
	policies = append(policies, policyFunc(httpHeaderPolicy), policyFunc(bodyDownloadPolicy), transportPolicy{trans: transport})
	return Pipeline{
		policies: policies,
	}
}

// Do is called for each and every HTTP request. It passes the request through all
// the Policy objects (which can transform the Request's URL/query parameters/headers)
// and ultimately sends the transformed HTTP request over the network.
func (p Pipeline) Do(req *Request) (*http.Response, error) {
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

// holds sentinel values used to send nulls
var nullables map[reflect.Type]interface{} = map[reflect.Type]interface{}{}

// NullValue is used to send an explicit 'null' within a request.
// This is typically used in JSON-MERGE-PATCH operations to delete a value.
func NullValue(v interface{}) interface{} {
	t := reflect.TypeOf(v)
	if k := t.Kind(); k != reflect.Ptr && k != reflect.Slice && k != reflect.Map {
		// t is not of pointer type, make it be of pointer type
		t = reflect.PtrTo(t)
	}
	v, found := nullables[t]
	if !found {
		var o reflect.Value
		if k := t.Kind(); k == reflect.Slice || k == reflect.Map {
			o = reflect.New(t) // *[]type / *map[]
			o = o.Elem()       // []type / map[]
		} else {
			o = reflect.New(t.Elem())
		}
		v = o.Interface()
		nullables[t] = v
	}
	// return the sentinel object
	return v
}

// IsNullValue returns true if the field contains a null sentinel value.
// This is used by custom marshallers to properly encode a null value.
func IsNullValue(v interface{}) bool {
	// see if our map has a sentinel object for this *T
	t := reflect.TypeOf(v)
	if k := t.Kind(); k != reflect.Ptr && k != reflect.Slice && k != reflect.Map {
		// v isn't a pointer type so it can never be a null
		return false
	}
	if o, found := nullables[t]; found {
		o1 := reflect.ValueOf(o)
		v1 := reflect.ValueOf(v)
		// we found it; return true if v points to the sentinel object.
		// NOTE: maps and slices can only be compared to nil, else you get
		// a runtime panic.  so we compare addresses instead.
		return o1.Pointer() == v1.Pointer()
	}
	// no sentinel object for this *t
	return false
}
