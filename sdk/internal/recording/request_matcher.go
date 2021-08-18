//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/dnaeon/go-vcr/cassette"
)

type RequestMatcher struct {
	ignoredHeaders map[string]*string
}

type StringMatcher func(reqVal string, recVal string) bool

var ignoredHeaders = map[string]struct{}{
	"Date":                   {},
	"X-Ms-Date":              {},
	"x-ms-date":              {},
	"x-ms-client-request-id": {},
	"User-Agent":             {},
	"Request-Id":             {},
	"traceparent":            {},
	"Authorization":          {},
}

const (
	recordingHeaderMissing = "Test recording headers do not match. Header '%s' is present in request but not in recording."
	requestHeaderMissing   = "Test recording headers do not match. Header '%s' is present in recording but not in request."
	headerValuesMismatch   = "Test recording header '%s' does not match. request: %s, recording: %s"
	methodMismatch         = "Test recording methods do not match. request: %s, recording: %s"
	urlMismatch            = "Test recording URLs do not match. request: %s, recording: %s"
	bodiesMismatch         = "Test recording bodies do not match.\nrequest: %s\nrecording: %s"
)

// defaultMatcher returns a new RequestMatcher configured with the default matching behavior.
func defaultMatcher(testContext TestContext) *RequestMatcher {
	// The default sanitizer sanitizes the Authorization header
	matcher := &RequestMatcher{
		context:        testContext,
		IgnoredHeaders: ignoredHeaders,
	}
	matcher.SetBodyMatcher(func(req string, rec string) bool {
		return defaultStringMatcher(req, rec)
	})
	matcher.SetURLMatcher(func(req string, rec string) bool {
		return defaultStringMatcher(req, rec)
	})
	matcher.SetMethodMatcher(func(req string, rec string) bool {
		return defaultStringMatcher(req, rec)
	})

	return matcher
}

// SetBodyMatcher replaces the default matching behavior with a custom StringMatcher that compares the string value of the request body payload with the string value of the recorded body payload.
func (m *RequestMatcher) SetBodyMatcher(matcher StringMatcher) {
	m.bodyMatcher = func(reqVal string, recVal string) bool {
		isMatch := matcher(reqVal, recVal)
		if !isMatch {
			m.context.Log(fmt.Sprintf(bodiesMismatch, recVal, recVal))
		}
		return isMatch
	}
}

// SetURLMatcher replaces the default matching behavior with a custom StringMatcher that compares the string value of the request URL with the string value of the recorded URL
func (m *RequestMatcher) SetURLMatcher(matcher StringMatcher) {
	m.urlMatcher = func(reqVal string, recVal string) bool {
		isMatch := matcher(reqVal, recVal)
		if !isMatch {
			m.context.Log(fmt.Sprintf(urlMismatch, recVal, recVal))
		}
		return isMatch
	}
}

// SetMethodMatcher replaces the default matching behavior with a custom StringMatcher that compares the string value of the request method with the string value of the recorded method
func (m *RequestMatcher) SetMethodMatcher(matcher StringMatcher) {
	m.methodMatcher = func(reqVal string, recVal string) bool {
		isMatch := matcher(reqVal, recVal)
		if !isMatch {
			m.context.Log(fmt.Sprintf(methodMismatch, recVal, recVal))
		}
		return isMatch
	}
}

var recordingHeaderMissing = "Test recording headers do not match. Header '%s' is present in request but not in recording."
var requestHeaderMissing = "Test recording headers do not match. Header '%s' is present in recording but not in request."
var headerValuesMismatch = "Test recording header '%s' does not match. request: %s, recording: %s"
var methodMismatch = "Test recording methods do not match. request: %s, recording: %s"
var urlMismatch = "Test recording URLs do not match. request: %s, recording: %s"
var bodiesMismatch = "Test recording bodies do not match.\nrequest: %s\nrecording: %s"

func compareBodies(r *http.Request, i cassette.Request, c TestContext) bool {
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
		c.Log(fmt.Sprintf(bodiesMismatch, body.String(), i.Body))
	}
	return bodiesMatch
}

func compareURLs(r *http.Request, i cassette.Request, c TestContext) bool {
	if r.URL.String() != i.URL {
		c.Log(fmt.Sprintf(urlMismatch, r.URL.String(), i.URL))
		return false
	}
	return true
}

func compareMethods(r *http.Request, i cassette.Request, c TestContext) bool {
	if r.Method != i.Method {
		c.Log(fmt.Sprintf(methodMismatch, r.Method, i.Method))
		return false
	}
	return true
}

func compareHeaders(r *http.Request, i cassette.Request, c TestContext) bool {
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
				c.Log(fmt.Sprintf(headerValuesMismatch, key, requestHeader, recordedHeader))
				return false
			}

		} else {
			// header not found
			c.Log(fmt.Sprintf(recordingHeaderMissing, key))
			return false
		}
	}
	if len(unVisitedCassetteKeys) > 0 {
		// headers exist in the recording that do not exist in the request
		for headerName := range unVisitedCassetteKeys {
			c.Log(fmt.Sprintf(requestHeaderMissing, headerName))
		}
		return false
	}
	return true
}
