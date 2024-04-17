//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

// Deprecated: the local recording API that uses this type is no longer supported. Call [SetDefaultMatcher] to
// configure the test proxy's matcher instead.
type RequestMatcher struct {
	// IgnoredHeaders is a map acting as a hash set of the header names that will be ignored for matching.
	// Modifying the keys in the map will affect how headers are matched for recordings.
	IgnoredHeaders map[string]struct{}
}

// Deprecated: only deprecated methods use this type.
type StringMatcher func(reqVal string, recVal string) bool

// SetBodyMatcher replaces the default matching behavior with a custom StringMatcher that compares the string value of the request body payload with the string value of the recorded body payload.
func (m *RequestMatcher) SetBodyMatcher(matcher StringMatcher) {
	panic(errUnsupportedAPI)
}

// SetURLMatcher replaces the default matching behavior with a custom StringMatcher that compares the string value of the request URL with the string value of the recorded URL
func (m *RequestMatcher) SetURLMatcher(matcher StringMatcher) {
	panic(errUnsupportedAPI)
}

// SetMethodMatcher replaces the default matching behavior with a custom StringMatcher that compares the string value of the request method with the string value of the recorded method
func (m *RequestMatcher) SetMethodMatcher(matcher StringMatcher) {
	panic(errUnsupportedAPI)
}
