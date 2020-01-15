// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Response represents the response from an HTTP request.
type Response struct {
	*http.Response

	// Payload contains the contents of the HTTP response body if available.
	Payload []byte
}

// CheckStatusCode returns a RequestError if the Response's status code isn't one of the specified values.
func (r *Response) CheckStatusCode(statusCodes ...int) error {
	if !r.HasStatusCode(statusCodes...) {
		msg := r.Status
		if len(r.Payload) > 0 {
			msg = string(r.Payload)
		}
		return newRequestError(msg, r)
	}
	return nil
}

// HasStatusCode returns true if the Response's status code is one of the specified values.
func (r *Response) HasStatusCode(statusCodes ...int) bool {
	if r == nil {
		return false
	}
	for _, sc := range statusCodes {
		if r.StatusCode == sc {
			return true
		}
	}
	return false
}

// UnmarshalAsXML calls xml.Unmarshal() to unmarshal the received payload into the value pointed to by v.
// If no payload was received a RequestError is returned.  If xml.Unmarshal fails a UnmarshalError is returned.
func (r *Response) UnmarshalAsXML(v interface{}) error {
	if len(r.Payload) == 0 {
		return newRequestError("missing payload", r)
	}
	r.removeBOM()
	err := xml.Unmarshal(r.Payload, v)
	if err != nil {
		err = fmt.Errorf("unmarshalling type %s: %w", reflect.TypeOf(v).Elem().Name(), err)
	}
	return err
}

// Drain reads the response body to completion then closes it.  The bytes read are discarded.
func (r *Response) Drain() {
	if r != nil && r.Body != nil {
		io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
	}
}

// removeBOM removes any byte-order mark prefix from the payload if present.
func (r *Response) removeBOM() {
	// UTF8
	r.Payload = bytes.TrimPrefix(r.Payload, []byte("\xef\xbb\xbf"))
}

// RetryAfter returns (non-zero, true) if the response contains a Retry-After header value
func (r *Response) RetryAfter() (time.Duration, bool) {
	if r == nil {
		return 0, false
	}
	if retryAfter, _ := strconv.Atoi(r.Header.Get("Retry-After")); retryAfter > 0 {
		return time.Duration(retryAfter) * time.Second, true
	}
	return 0, false
}

// WriteRequestWithResponse appends a formatted HTTP request into a Buffer. If request and/or err are
// not nil, then these are also written into the Buffer.
func WriteRequestWithResponse(b *bytes.Buffer, request *Request, response *Response, err error) {
	// Write the request into the buffer.
	fmt.Fprint(b, "   "+request.Method+" "+request.URL.String()+"\n")
	writeHeader(b, request.Header)
	if response != nil {
		fmt.Fprintln(b, "   --------------------------------------------------------------------------------")
		fmt.Fprint(b, "   RESPONSE Status: "+response.Status+"\n")
		writeHeader(b, response.Header)
	}
	if err != nil {
		fmt.Fprintln(b, "   --------------------------------------------------------------------------------")
		fmt.Fprint(b, "   ERROR:\n"+err.Error()+"\n")
	}
}

// formatHeaders appends an HTTP request's or response's header into a Buffer.
func writeHeader(b *bytes.Buffer, header http.Header) {
	if len(header) == 0 {
		b.WriteString("   (no headers)\n")
		return
	}
	keys := make([]string, 0, len(header))
	// Alphabetize the headers
	for k := range header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		// Redact the value of any Authorization header to prevent security information from persisting in logs
		value := interface{}("REDACTED")
		if !strings.EqualFold(k, "Authorization") {
			value = header[k]
		}
		fmt.Fprintf(b, "   %s: %+v\n", k, value)
	}
}
