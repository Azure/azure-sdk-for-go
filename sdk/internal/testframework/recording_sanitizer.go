// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package testframework

import (
<<<<<<< HEAD
<<<<<<< HEAD
	"net/http"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
=======
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"net/http"
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	"net/http"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
>>>>>>> a450961e3 (fb)
)

type RecordingSanitizer struct {
	recorder          *recorder.Recorder
	headersToSanitize map[string]*string
	urlSanitizer      StringSanitizer
<<<<<<< HEAD
<<<<<<< HEAD
	bodySanitizer     StringSanitizer
}

type StringSanitizer func(*string)
=======
}

type StringSanitizer func(*string) error
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	bodySanitizer     StringSanitizer
}

type StringSanitizer func(*string)
>>>>>>> a450961e3 (fb)

const SanitizedValue string = "sanitized"
const SanitizedBase64Value string = "Kg=="

var sanitizedValueSlice = []string{SanitizedValue}

func DefaultSanitizer(recorder *recorder.Recorder) *RecordingSanitizer {
	// The default sanitizer sanitizes the Authorization header
	s := &RecordingSanitizer{headersToSanitize: map[string]*string{"Authorization": nil}, recorder: recorder, urlSanitizer: DefaultStringSanitizer}
	recorder.AddSaveFilter(s.applySaveFilter)

	return s
}

<<<<<<< HEAD
<<<<<<< HEAD
// AddSanitizedHeaders adds the supplied header names to the list of headers to be sanitized on request and response recordings.
=======
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
// AddSanitizedHeaders adds the supplied header names to the list of headers to be sanitized on request and response recordings.
>>>>>>> a450961e3 (fb)
func (s *RecordingSanitizer) AddSanitizedHeaders(headers ...string) {
	for _, headerName := range headers {
		s.headersToSanitize[headerName] = nil
	}
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> a450961e3 (fb)
// AddBodysanitizer configures the supplied StringSanitizer to sanitize recording request and response bodies
func (s *RecordingSanitizer) AddBodysanitizer(sanitizer StringSanitizer) {
	s.bodySanitizer = sanitizer
}

// AddUriSanitizer configures the supplied StringSanitizer to sanitize recording request and response URLs
func (s *RecordingSanitizer) AddUrlSanitizer(sanitizer StringSanitizer) {
	s.urlSanitizer = sanitizer
}

<<<<<<< HEAD
=======
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
>>>>>>> a450961e3 (fb)
func (s *RecordingSanitizer) sanitizeHeaders(header http.Header) {
	for headerName := range s.headersToSanitize {
		if _, ok := header[headerName]; ok {
			header[headerName] = sanitizedValueSlice
		}
	}
}

<<<<<<< HEAD
<<<<<<< HEAD
func (s *RecordingSanitizer) sanitizeBodies(body *string) {
	s.bodySanitizer(body)
}

func (s *RecordingSanitizer) sanitizeURL(url *string) {
	s.urlSanitizer(url)
=======
func (s *RecordingSanitizer) AddUrlSanitizer(sanitizer StringSanitizer) {
	s.urlSanitizer = sanitizer
}

func (s *RecordingSanitizer) sanitizeURL(url *string) error {
	return s.urlSanitizer(url)
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
func (s *RecordingSanitizer) sanitizeBodies(body *string) {
	s.bodySanitizer(body)
}

func (s *RecordingSanitizer) sanitizeURL(url *string) {
	s.urlSanitizer(url)
>>>>>>> a450961e3 (fb)
}

func (s *RecordingSanitizer) applySaveFilter(i *cassette.Interaction) error {
	s.sanitizeHeaders(i.Request.Headers)
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> a450961e3 (fb)
	s.sanitizeHeaders(i.Response.Headers)
	s.sanitizeURL(&i.Request.URL)
	if len(i.Request.Body) > 0 {
		s.sanitizeBodies(&i.Request.Body)
	}
	if len(i.Response.Body) > 0 {
		s.sanitizeBodies(&i.Response.Body)
	}
<<<<<<< HEAD
	return nil
}

func DefaultStringSanitizer(s *string) {}
=======
	return s.sanitizeURL(&i.Request.URL)
}

func DefaultStringSanitizer(s *string) error {
	return nil
}
>>>>>>> a90f565e4 (testFramework initial implementation with aztables tests)
=======
	return nil
}

func DefaultStringSanitizer(s *string) {}
>>>>>>> a450961e3 (fb)
