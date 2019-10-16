// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// Request is an abstraction over the creation of an HTTP request as it passes through the pipeline.
type Request struct {
	*http.Request
	policies []Policy
	qp       url.Values
	values   opValues
}

type opValues map[reflect.Type]interface{}

// Set adds/changes a value
func (ov opValues) set(value interface{}) {
	ov[reflect.TypeOf(value)] = value
}

// Get looks for a value set by SetValue first
func (ov opValues) get(value interface{}) bool {
	v, ok := ov[reflect.TypeOf(value).Elem()]
	if ok {
		reflect.Indirect(reflect.ValueOf(value)).Set(reflect.ValueOf(v))
	}
	return ok
}

// NewRequest ...
func (p Pipeline) NewRequest(method string, URL url.URL) *Request {
	// removeEmptyPort strips the empty port in ":port" to ""
	// as mandated by RFC 3986 Section 6.2.3.
	// adapted from removeEmptyPort() in net/http.go
	if strings.LastIndex(URL.Host, ":") > strings.LastIndex(URL.Host, "]") {
		URL.Host = strings.TrimSuffix(URL.Host, ":")
	}
	return &Request{
		Request: &http.Request{
			Method:     method,
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

// Do is called for each and every HTTP request. It passes the Context through all the Policy objects
// (which can transform the Request's URL/query parameters/headers) and ultimately sends the
// transformed HTTP request over the network.
func (req *Request) Do(ctx context.Context) (*Response, error) {
	nextPolicy := req.policies[0]
	nextMsg := req
	nextMsg.policies = nextMsg.policies[1:]
	// encode any pending query params
	if nextMsg.qp != nil {
		nextMsg.Request.URL.RawQuery = nextMsg.qp.Encode()
		nextMsg.qp = nil
	}
	return nextPolicy.Do(ctx, nextMsg)
}

// MarshalAsXML calls xml.Marshal() to get the XML encoding of v then calls SetBody.
// If xml.Marshal fails a MarshalError is returned.  Any error from SetBody is returned.
func (req *Request) MarshalAsXML(v interface{}) error {
	b, err := xml.Marshal(v)
	if err != nil {
		return fmt.Errorf("error marshalling type %s: %w", reflect.TypeOf(v).Name(), err)
	}
	req.Header.Set("Content-Type", "application/xml")
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
	return req.values.get(&value)
}

// SetQueryParam sets the key to value.
func (req *Request) SetQueryParam(key, value string) {
	if req.qp == nil {
		req.qp = req.Request.URL.Query()
	}
	req.qp.Set(key, value)
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
	body = &retryableRequestBody{body: body}
	req.Request.Body = body
	req.Request.ContentLength = size
	req.Request.GetBody = func() (io.ReadCloser, error) {
		_, err := body.Seek(0, io.SeekStart) // Seek back to the beginning of the stream
		return body, err
	}
	return nil
}

// SkipBodyDownload will disable automatic downloading of the response body via the DownloadResponseBody policy.
func (req *Request) SkipBodyDownload() {
	req.SetOperationValue(bodyDownloadPolicyOpValues{skip: true})
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

// ********** The following type/methods implement the retryableRequestBody (a ReadSeekCloser)

// This struct is used when sending a body to the network
type retryableRequestBody struct {
	body io.ReadSeeker // Seeking is required to support retries
}

// Read reads a block of data from an inner stream and reports progress
func (b *retryableRequestBody) Read(p []byte) (n int, err error) {
	return b.body.Read(p)
}

func (b *retryableRequestBody) Seek(offset int64, whence int) (offsetFromStart int64, err error) {
	return b.body.Seek(offset, whence)
}

func (b *retryableRequestBody) Close() error {
	// We don't want the underlying transport to close the request body on transient failures so this is a nop.
	// The pipeline closes the request body upon success.
	return nil
}

func (b *retryableRequestBody) realClose() error {
	if c, ok := b.body.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
