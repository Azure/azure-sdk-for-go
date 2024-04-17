//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"fmt"
)

// Deprecated: the local recording API that uses this type is no longer supported. Call [SetDefaultMatcher] to
// configure the test proxy's matcher instead.
type RequestMatcher struct {
	context TestContext
	// IgnoredHeaders is a map acting as a hash set of the header names that will be ignored for matching.
	// Modifying the keys in the map will affect how headers are matched for recordings.
	IgnoredHeaders map[string]struct{}
	bodyMatcher    StringMatcher
	urlMatcher     StringMatcher
	methodMatcher  StringMatcher
}

// Deprecated: only deprecated methods use this type.
type StringMatcher func(reqVal string, recVal string) bool

const (
	recordingHeaderMissing = "Test recording headers do not match. Header '%s' is present in request but not in recording."
	requestHeaderMissing   = "Test recording headers do not match. Header '%s' is present in recording but not in request."
	headerValuesMismatch   = "Test recording header '%s' does not match. request: %s, recording: %s"
	methodMismatch         = "Test recording methods do not match. request: %s, recording: %s"
	urlMismatch            = "Test recording URLs do not match. request: %s, recording: %s"
	bodiesMismatch         = "Test recording bodies do not match.\nrequest: %s\nrecording: %s"
)

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
