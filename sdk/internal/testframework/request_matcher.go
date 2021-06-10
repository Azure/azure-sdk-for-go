// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package testframework

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/dnaeon/go-vcr/cassette"
)

type RequestMatcher struct {
	context        TestContext
	IgnoredHeaders map[string]*string
	bodyMatcher    StringMatcher
	urlMatcher     StringMatcher
	methodMatcher  StringMatcher
}

type StringMatcher func(reqVal string, recVal string) bool
type matcherWrapper func(matcher StringMatcher, testContext TestContext) bool

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
var bodiesMismatch = "Test recording bodies do not match.\nrequest: %s\nrecording: %s"

func DefaultMatcher(testContext TestContext) *RequestMatcher {
	// The default sanitizer sanitizes the Authorization header
	matcher := &RequestMatcher{
		context:        testContext,
		IgnoredHeaders: ignoredHeaders,
	}
	matcher.SetBodyMatcher(defaultStringMatcher)
	matcher.SetURLMatcher(defaultStringMatcher)
	matcher.SetMethodMatcher(defaultStringMatcher)

	return matcher
}

func (m *RequestMatcher) SetBodyMatcher(matcher StringMatcher) {
	m.setMatcher(matcher, bodiesMismatch)
}

func (m *RequestMatcher) SetURLMatcher(matcher StringMatcher) {
	m.setMatcher(matcher, urlMismatch)
}

func (m *RequestMatcher) SetMethodMatcher(matcher StringMatcher) {
	m.setMatcher(matcher, methodMismatch)
}

func (m *RequestMatcher) setMatcher(matcher StringMatcher, message string) {
	m.bodyMatcher = func(reqVal string, recVal string) bool {
		isMatch := matcher(reqVal, recVal)
		if !isMatch {
			m.context.Log(fmt.Sprintf(message, recVal, recVal))
		}
		return isMatch
	}
}

func defaultStringMatcher(s1 string, s2 string) bool {
	return s1 == s2
}

func getBody(r *http.Request) string {
	body := bytes.Buffer{}
	if r.Body != nil {
		_, err := body.ReadFrom(r.Body)
		if err != nil {
			return "could not parse body: " + err.Error()
		}
		r.Body = ioutil.NopCloser(&body)
	}
	return body.String()
}

func getUrl(r *http.Request) string {
	return r.URL.String()
}

func getMethod(r *http.Request) string {
	return r.Method
}

func (m *RequestMatcher) compareBodies(r *http.Request, recordedBody string) bool {
	body := getBody(r)
	return m.bodyMatcher(body, recordedBody)
}

func (m *RequestMatcher) compareURLs(r *http.Request, recordedUrl string) bool {
	body := getUrl(r)
	return m.urlMatcher(body, recordedUrl)
}

func (m *RequestMatcher) compareMethods(r *http.Request, recordedMethod string) bool {
	body := getMethod(r)
	return m.urlMatcher(body, recordedMethod)
}

func (m *RequestMatcher) compareHeaders(r *http.Request, i cassette.Request) bool {
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
				m.context.Log(fmt.Sprintf(headerValuesMismatch, key, requestHeader, recordedHeader))
				return false
			}

		} else {
			// header not found
			m.context.Log(fmt.Sprintf(recordingHeaderMissing, key))
			return false
		}
	}
	if len(unVisitedCassetteKeys) > 0 {
		// headers exist in the recording that do not exist in the request
		for headerName := range unVisitedCassetteKeys {
			m.context.Log(fmt.Sprintf(requestHeaderMissing, headerName))
		}
		return false
	}
	return true
}
