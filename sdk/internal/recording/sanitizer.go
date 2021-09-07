//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
)

type Sanitizer struct {
	recorder          *recorder.Recorder
	headersToSanitize []string
	urlSanitizer      StringSanitizer
	bodySanitizer     StringSanitizer
}

// StringSanitizer is a func that will modify the string pointed to by the parameter into a sanitized value.
type StringSanitizer func(*string)

// SanitizedValue is the default placeholder value to be used for sanitized strings.
const SanitizedValue string = "sanitized"

// SanitizedBase64Value is the default placeholder value to be used for sanitized base-64 encoded strings.
const SanitizedBase64Value string = "Kg=="

var sanitizedValueSlice = []string{SanitizedValue}

// defaultSanitizer returns a new RecordingSanitizer with the default sanitizing behavior.
// To customize sanitization, call AddSanitizedHeaders, AddBodySanitizer, or AddUrlSanitizer.
func defaultSanitizer(recorder *recorder.Recorder) *Sanitizer {
	// The default sanitizer sanitizes the Authorization header
	s := &Sanitizer{headersToSanitize: []string{"Authorization"}, recorder: recorder, urlSanitizer: DefaultStringSanitizer, bodySanitizer: DefaultStringSanitizer}
	recorder.AddSaveFilter(s.applySaveFilter)

	return s
}

// AddSanitizedHeaders adds the supplied header names to the list of headers to be sanitized on request and response recordings.
func (s *Sanitizer) AddSanitizedHeaders(headers ...string) {
	s.headersToSanitize = append(s.headersToSanitize, headers...)
}

// AddBodysanitizer configures the supplied StringSanitizer to sanitize recording request and response bodies
func (s *Sanitizer) AddBodysanitizer(sanitizer StringSanitizer) {
	s.bodySanitizer = sanitizer
}

// AddUriSanitizer configures the supplied StringSanitizer to sanitize recording request and response URLs
func (s *Sanitizer) AddUrlSanitizer(sanitizer StringSanitizer) {
	s.urlSanitizer = sanitizer
}

func (s *Sanitizer) sanitizeHeaders(header http.Header) {
	for _, headerName := range s.headersToSanitize {
		if _, ok := header[headerName]; ok {
			header[headerName] = sanitizedValueSlice
		}
	}
}

func (s *Sanitizer) sanitizeBodies(body *string) {
	s.bodySanitizer(body)
}

func (s *Sanitizer) sanitizeURL(url *string) {
	s.urlSanitizer(url)
}

func (s *Sanitizer) applySaveFilter(i *cassette.Interaction) error {
	s.sanitizeHeaders(i.Request.Headers)
	s.sanitizeHeaders(i.Response.Headers)
	s.sanitizeURL(&i.Request.URL)
	if len(i.Request.Body) > 0 {
		s.sanitizeBodies(&i.Request.Body)
	}
	if len(i.Response.Body) > 0 {
		s.sanitizeBodies(&i.Response.Body)
	}
	return nil
}

func DefaultStringSanitizer(s *string) {}

func AddBodyKeySanitizer(jsonPath, replacementValue, regex, groupForReplace string, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.HostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "BodyKeySanitizer")
	bodyContent := map[string]string{
		"jsonPath": jsonPath,
	}
	if replacementValue != "" {
		bodyContent["value"] = replacementValue
	}
	if regex != "" {
		bodyContent["regex"] = regex
	}
	if groupForReplace != "" {
		bodyContent["groupForReplace"] = groupForReplace
	}

	marshalled, err := json.Marshal(bodyContent)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	_, err = client.Do(req)
	return err
}

// value: the substitution value
// regex: the regex to match on request/response entries
// groupForReplace: If your regex has multiple groups, the named group which to replace. If your regex does not, make this an empty string
func AddBodyRegexSanitizer(value, regex, groupForReplace string, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.HostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "BodyRegexSanitizer")
	bodyContent := map[string]string{
		"value": value,
	}
	if regex != "" {
		bodyContent["regex"] = regex
	}
	if groupForReplace != "" {
		bodyContent["groupForReplace"] = groupForReplace
	}

	marshalled, err := json.Marshal(bodyContent)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	_, err = client.Do(req)
	return err
}

// value: the substitution value
// regex: the regex to match on request/response entries
// groupForReplace: If your regex has multiple groups, the named group which to replace. If your regex does not, make this an empty string
func AddContinuationSanitizer(key, method string, resetAfterFirst bool, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.HostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "ContinuationSanitizer")
	bodyContent := map[string]string{
		"key":             key,
		"method":          method,
		"resetAfterFirst": fmt.Sprintf("%v", resetAfterFirst),
	}

	marshalled, err := json.Marshal(bodyContent)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	_, err = client.Do(req)
	return err
}

func AddGeneralRegexSanitizer(value, regex, groupForReplace string, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.HostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "GeneralRegexSanitizer")
	bodyContent := map[string]string{
		"value":           value,
		"regex":           regex,
		"groupForReplace": groupForReplace,
	}

	marshalled, err := json.Marshal(bodyContent)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	_, err = client.Do(req)
	return err
}

func AddHeaderRegexSanitizer(key, replacementValue, regex, groupForReplace string, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.HostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "HeaderRegexSanitizer")
	bodyContent := map[string]string{
		"key": key,
	}
	if replacementValue != "" {
		bodyContent["value"] = replacementValue
	}
	if regex != "" {
		bodyContent["regex"] = regex
	}
	if groupForReplace != "" {
		bodyContent["groupForReplace"] = groupForReplace
	}

	marshalled, err := json.Marshal(bodyContent)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	_, err = client.Do(req)
	return err
}

func AddOAuthResponseSanitizer(options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.HostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "OAuthResponseSanitizer")
	_, err = client.Do(req)
	return err
}

func AddRemoveHeaderSanitizer(headersForRemoval []string, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.HostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "RemoveHeaderSanitizer")
	bodyContent := map[string]string{
		"headersForRemoval": strings.Join(headersForRemoval, ","),
	}

	marshalled, err := json.Marshal(bodyContent)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	_, err = client.Do(req)
	return err
}

func AddReplaceRequestSubscriptionIdSanitizer(value string, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.HostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "ReplaceRequestSubscriptionId")
	bodyContent := map[string]string{
		"value": value,
	}

	marshalled, err := json.Marshal(bodyContent)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	_, err = client.Do(req)
	return err
}

func AddUriSanitizer(replacement, regex string, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%v/Admin/AddSanitizer", options.HostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "UriRegexSanitizer")
	bodyContent := map[string]string{
		"value": replacement,
		"regex": regex,
	}
	marshalled, err := json.Marshal(bodyContent)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	_, err = client.Do(req)
	return err
}
