// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	contentTypeAppJSON = "application/json"
	contentTypeAppXML  = "application/xml"
)

// Request is an abstraction over the creation of an HTTP request as it passes through the pipeline.
type Request struct {
	*http.Request
	policies []Policy
	values   opValues
}

type opValues map[reflect.Type]interface{}

// Set adds/changes a value
func (ov opValues) set(value interface{}) {
	ov[reflect.TypeOf(value)] = value
}

// Get looks for a value set by SetValue first
func (ov opValues) get(value interface{}) bool {
	v, ok := ov[reflect.ValueOf(value).Elem().Type()]
	if ok {
		reflect.ValueOf(value).Elem().Set(reflect.ValueOf(v))
	}
	return ok
}

// NewRequest creates a new Request with the specified input.
func NewRequest(httpMethod string, endpoint url.URL) *Request {
	// removeEmptyPort strips the empty port in ":port" to ""
	// as mandated by RFC 3986 Section 6.2.3.
	// adapted from removeEmptyPort() in net/http.go
	if strings.LastIndex(endpoint.Host, ":") > strings.LastIndex(endpoint.Host, "]") {
		endpoint.Host = strings.TrimSuffix(endpoint.Host, ":")
	}
	return &Request{
		Request: &http.Request{
			Method:     httpMethod,
			URL:        &endpoint,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     http.Header{},
			Host:       endpoint.Host,
		},
	}
}

// Next calls the next policy in the pipeline.
// If there are no more policies, nil and ErrNoMorePolicies are returned.
// This method is intended to be called from pipeline policies.
// To send a request through a pipeline call Pipeline.Do().
func (req *Request) Next(ctx context.Context) (*Response, error) {
	if len(req.policies) == 0 {
		return nil, ErrNoMorePolicies
	}
	nextPolicy := req.policies[0]
	nextReq := *req
	nextReq.policies = nextReq.policies[1:]
	return nextPolicy.Do(ctx, &nextReq)
}

// MarshalAsJSON calls json.Marshal() to get the JSON encoding of v then calls SetBody.
// If json.Marshal fails a MarshalError is returned.  Any error from SetBody is returned.
func (req *Request) MarshalAsJSON(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("error marshalling type %s: %w", reflect.TypeOf(v).Name(), err)
	}
	req.Header.Set(HeaderContentType, contentTypeAppJSON)
	return req.SetBody(NopCloser(bytes.NewReader(b)))
}

// MarshalAsXML calls xml.Marshal() to get the XML encoding of v then calls SetBody.
// If xml.Marshal fails a MarshalError is returned.  Any error from SetBody is returned.
func (req *Request) MarshalAsXML(v interface{}) error {
	b, err := xml.Marshal(v)
	if err != nil {
		return fmt.Errorf("error marshalling type %s: %w", reflect.TypeOf(v).Name(), err)
	}
	req.Header.Set(HeaderContentType, contentTypeAppXML)
	return req.SetBody(NopCloser(bytes.NewReader(b)))
}

// SetOperationValue adds/changes a mutable key/value associated with a single operation.
func (req *Request) SetOperationValue(value interface{}) {
	if req.values == nil {
		req.values = opValues{}
	}
	req.values.set(value)
}

// OperationValue looks for a value set by SetOperationValue().
func (req *Request) OperationValue(value interface{}) bool {
	if req.values == nil {
		return false
	}
	return req.values.get(value)
}

// SetBody sets the specified ReadSeekCloser as the HTTP request body.
func (req *Request) SetBody(body ReadSeekCloser) error {
	// Set the body and content length.
	size, err := body.Seek(0, io.SeekEnd) // Seek to the end to get the stream's size
	if err != nil {
		return err
	}
	if size == 0 {
		body.Close()
		return nil
	}
	_, err = body.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	req.Request.Body = body
	req.Request.ContentLength = size
	return nil
}

// SkipBodyDownload will disable automatic downloading of the response body.
func (req *Request) SkipBodyDownload() {
	req.SetOperationValue(bodyDownloadPolicyOpValues{skip: true})
}

// returns true if auto-body download policy is enabled
func (req *Request) bodyDownloadEnabled() bool {
	var opValues bodyDownloadPolicyOpValues
	req.OperationValue(&opValues)
	return !opValues.skip
}

// RewindBody seeks the request's Body stream back to the beginning so it can be resent when retrying an operation.
func (req *Request) RewindBody() error {
	if req.Body != nil {
		// Reset the stream back to the beginning
		_, err := req.Body.(io.Seeker).Seek(0, io.SeekStart)
		return err
	}
	return nil
}

// Close closes the request body.
func (req *Request) Close() error {
	if req.Body == nil {
		return nil
	}
	return req.Body.Close()
}

func (req *Request) copy() *Request {
	clonedURL := *req.URL
	// Copy the values and immutable references
	return &Request{
		Request: &http.Request{
			Method:        req.Method,
			URL:           &clonedURL,
			Proto:         req.Proto,
			ProtoMajor:    req.ProtoMajor,
			ProtoMinor:    req.ProtoMinor,
			Header:        req.Header.Clone(),
			Host:          req.URL.Host,
			Body:          req.Body, // shallow copy
			ContentLength: req.ContentLength,
			GetBody:       req.GetBody,
		},
	}
}
