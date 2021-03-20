// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package testframework

import (
	"bytes"
<<<<<<< HEAD
<<<<<<< HEAD
	"fmt"
=======
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	"fmt"
>>>>>>> 656c2801d (eliminate dependency on go-check)
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/dnaeon/go-vcr/cassette"
<<<<<<< HEAD
<<<<<<< HEAD
=======
	chk "gopkg.in/check.v1"
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
>>>>>>> 656c2801d (eliminate dependency on go-check)
)

type RequestMatcher struct {
	ignoredHeaders map[string]*string
}

var ignoredHeaders = map[string]*string{
	"Date":                   nil,
	"X-Ms-Date":              nil,
	"x-ms-date":              nil,
	"x-ms-client-request-id": nil,
	"User-Agent":             nil,
	"Request-Id":             nil,
	"traceparent":            nil,
	"Authorization":          nil,
}

var recordingHeaderMissing = "Test recording headers do not match. Header '%s' is present in request but not in recording."
var requestHeaderMissing = "Test recording headers do not match. Header '%s' is present in recording but not in request."
var headerValuesMismatch = "Test recording header '%s' does not match. request: %s, recording: %s"
var methodMismatch = "Test recording methods do not match. request: %s, recording: %s"
var urlMismatch = "Test recording URLs do not match. request: %s, recording: %s"
<<<<<<< HEAD
<<<<<<< HEAD
var bodiesMismatch = "Test recording bodies do not match.\nrequest: %s\nrecording: %s"

func compareBodies(r *http.Request, i cassette.Request, c TestContext) bool {
=======
=======
var bodiesMismatch = "Test recording bodies do not match.\nrequest: %s\nrecording: %s"
>>>>>>> c983287ea (refactor tests to testify)

<<<<<<< HEAD
func compareBodies(r *http.Request, i cassette.Request, c *chk.C) bool {
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
func compareBodies(r *http.Request, i cassette.Request, c TestContext) bool {
>>>>>>> 656c2801d (eliminate dependency on go-check)
	body := bytes.Buffer{}
	if r.Body != nil {
		_, err := body.ReadFrom(r.Body)
		if err != nil {
			return false
		}
		r.Body = ioutil.NopCloser(&body)
	}
	bodiesMatch := body.String() == i.Body
	if !bodiesMatch {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		c.Log(fmt.Sprintf(bodiesMismatch, body.String(), i.Body))
=======
		c.Logf("Test recording bodies do not match.\nrequest:   %s\nrecording: %s", body.String(), i.Body)
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
		c.Log(fmt.Sprintf("Test recording bodies do not match.\nrequest: %s\nrecording: %s", body.String(), i.Body))
>>>>>>> 656c2801d (eliminate dependency on go-check)
=======
		c.Log(fmt.Sprintf(bodiesMismatch, body.String(), i.Body))
>>>>>>> c983287ea (refactor tests to testify)
	}
	return bodiesMatch
}

<<<<<<< HEAD
<<<<<<< HEAD
func compareURLs(r *http.Request, i cassette.Request, c TestContext) bool {
	if r.URL.String() != i.URL {
		c.Log(fmt.Sprintf(urlMismatch, r.URL.String(), i.URL))
=======
func compareURLs(r *http.Request, i cassette.Request, c *chk.C) bool {
	if r.URL.String() != i.URL {
		c.Logf(urlMismatch, r.URL.String(), i.URL)
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
func compareURLs(r *http.Request, i cassette.Request, c TestContext) bool {
	if r.URL.String() != i.URL {
		c.Log(fmt.Sprintf(urlMismatch, r.URL.String(), i.URL))
>>>>>>> 656c2801d (eliminate dependency on go-check)
		return false
	}
	return true
}

<<<<<<< HEAD
<<<<<<< HEAD
func compareMethods(r *http.Request, i cassette.Request, c TestContext) bool {
	if r.Method != i.Method {
		c.Log(fmt.Sprintf(methodMismatch, r.Method, i.Method))
=======
func compareMethods(r *http.Request, i cassette.Request, c *chk.C) bool {
	if r.Method != i.Method {
		c.Logf(methodMismatch, r.Method, i.Method)
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
func compareMethods(r *http.Request, i cassette.Request, c TestContext) bool {
	if r.Method != i.Method {
		c.Log(fmt.Sprintf(methodMismatch, r.Method, i.Method))
>>>>>>> 656c2801d (eliminate dependency on go-check)
		return false
	}
	return true
}

<<<<<<< HEAD
<<<<<<< HEAD
func compareHeaders(r *http.Request, i cassette.Request, c TestContext) bool {
=======
func compareHeaders(r *http.Request, i cassette.Request, c *chk.C) bool {
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
func compareHeaders(r *http.Request, i cassette.Request, c TestContext) bool {
>>>>>>> 656c2801d (eliminate dependency on go-check)
	unVisitedCassetteKeys := make(map[string]*string, len(i.Headers))
	// clone the cassette keys to track which we have seen
	for k := range i.Headers {
		if _, ignore := ignoredHeaders[k]; ignore {
			// don't copy ignored headers
			continue
		}
		unVisitedCassetteKeys[k] = nil
	}
	//iterate through all the request headers to compare them to cassette headers
	for key, requestHeader := range r.Header {
		if _, ignore := ignoredHeaders[key]; ignore {
			// this is an ignorable header
			continue
		}
		delete(unVisitedCassetteKeys, key)
		if recordedHeader, foundMatch := i.Headers[key]; foundMatch {
			headersMatch := reflect.DeepEqual(requestHeader, recordedHeader)
			if !headersMatch {
				// headers don't match
<<<<<<< HEAD
<<<<<<< HEAD
				c.Log(fmt.Sprintf(headerValuesMismatch, key, requestHeader, recordedHeader))
=======
				c.Logf(headerValuesMismatch, key, requestHeader, recordedHeader)
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
				c.Log(fmt.Sprintf(headerValuesMismatch, key, requestHeader, recordedHeader))
>>>>>>> 656c2801d (eliminate dependency on go-check)
				return false
			}

		} else {
			// header not found
<<<<<<< HEAD
<<<<<<< HEAD
			c.Log(fmt.Sprintf(recordingHeaderMissing, key))
=======
			c.Logf(recordingHeaderMissing, key)
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
			c.Log(fmt.Sprintf(recordingHeaderMissing, key))
>>>>>>> 656c2801d (eliminate dependency on go-check)
			return false
		}
	}
	if len(unVisitedCassetteKeys) > 0 {
		// headers exist in the recording that do not exist in the request
		for headerName := range unVisitedCassetteKeys {
<<<<<<< HEAD
<<<<<<< HEAD
			c.Log(fmt.Sprintf(requestHeaderMissing, headerName))
=======
			c.Logf(requestHeaderMissing, headerName)
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
			c.Log(fmt.Sprintf(requestHeaderMissing, headerName))
>>>>>>> 656c2801d (eliminate dependency on go-check)
		}
		return false
	}
	return true
}
