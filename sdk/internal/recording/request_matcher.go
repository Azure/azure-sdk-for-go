//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/dnaeon/go-vcr/cassette"
)

type RequestMatcher struct {
	context TestContext
	// IgnoredHeaders is a map acting as a hash set of the header names that will be ignored for matching.
	// Modifying the keys in the map will affect how headers are matched for recordings.
	IgnoredHeaders map[string]struct{}
	bodyMatcher    StringMatcher
	urlMatcher     StringMatcher
	methodMatcher  StringMatcher
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
		r.Body = io.NopCloser(&body)
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
	url := getUrl(r)
	return m.urlMatcher(url, recordedUrl)
}

func (m *RequestMatcher) compareMethods(r *http.Request, recordedMethod string) bool {
	method := getMethod(r)
	return m.methodMatcher(method, recordedMethod)
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
