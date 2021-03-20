// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package testframework

import (
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"net/http"
)

type RecordingSanitizer struct {
	recorder          *recorder.Recorder
	headersToSanitize map[string]*string
	urlSanitizer      StringSanitizer
}

type StringSanitizer func(*string) error

const SanitizedValue string = "sanitized"
const SanitizedBase64Value string = "Kg=="

var sanitizedValueSlice = []string{SanitizedValue}

func DefaultSanitizer(recorder *recorder.Recorder) *RecordingSanitizer {
	// The default sanitizer sanitizes the Authorization header
	s := &RecordingSanitizer{headersToSanitize: map[string]*string{"Authorization": nil}, recorder: recorder, urlSanitizer: DefaultStringSanitizer}
	recorder.AddSaveFilter(s.applySaveFilter)

	return s
}

func (s *RecordingSanitizer) AddSanitizedHeaders(headers ...string) {
	for _, headerName := range headers {
		s.headersToSanitize[headerName] = nil
	}
}

func (s *RecordingSanitizer) sanitizeHeaders(header http.Header) {
	for headerName := range s.headersToSanitize {
		if _, ok := header[headerName]; ok {
			header[headerName] = sanitizedValueSlice
		}
	}
}

func (s *RecordingSanitizer) AddUrlSanitizer(sanitizer StringSanitizer) {
	s.urlSanitizer = sanitizer
}

func (s *RecordingSanitizer) sanitizeURL(url *string) error {
	return s.urlSanitizer(url)
}

func (s *RecordingSanitizer) applySaveFilter(i *cassette.Interaction) error {
	s.sanitizeHeaders(i.Request.Headers)
	return s.sanitizeURL(&i.Request.URL)
}

func DefaultStringSanitizer(s *string) error {
	return nil
}
