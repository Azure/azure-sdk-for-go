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

// AddURISanitizer configures the supplied StringSanitizer to sanitize recording request and response URLs
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

func handleProxyResponse(resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusOK {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("there was an error communicating with the test proxy: %s", body)
}

func AddBodyKeySanitizer(jsonPath, replacementValue, regex string, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.hostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "BodyKeySanitizer")

	marshalled, err := json.MarshalIndent(struct {
		JSONPath        string `json:"jsonPath"`
		Value           string `json:"value,omitempty"`
		Regex           string `json:"regex,omitempty"`
		GroupForReplace string `json:"groupForReplace,omitempty"`
	}{
		JSONPath:        jsonPath,
		Value:           replacementValue,
		Regex:           regex,
		GroupForReplace: options.GroupForReplace,
	}, "", "")
	if err != nil {
		return err
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	return handleProxyResponse(client.Do(req))
}

// value: the substitution value
// regex: the regex to match on request/response entries
func AddBodyRegexSanitizer(value, regex string, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.hostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "BodyRegexSanitizer")

	marshalled, err := json.MarshalIndent(struct {
		Value           string `json:"value"`
		Regex           string `json:"regex,omitempty"`
		GroupForReplace string `json:"groupForReplace,omitempty"`
	}{
		Value:           value,
		Regex:           regex,
		GroupForReplace: options.GroupForReplace,
	}, "", "")
	if err != nil {
		return err
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	return handleProxyResponse(client.Do(req))
}

// key: the name of the header whos value will be replaced from response -> next request
// method: the method by which the value of the targeted key will be replaced. Defaults to GUID replacement
// resetAfterFirt: Do we need multiple pairs replaced? Or do we want to replace each value with the same value.
func AddContinuationSanitizer(key, method string, resetAfterFirst bool, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.hostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "ContinuationSanitizer")

	marshalled, err := json.MarshalIndent(struct {
		Key             string `json:"key"`
		Method          string `json:"method"`
		ResetAfterFirst string `json:"resetAfterFirst"`
	}{
		Key:             key,
		Method:          method,
		ResetAfterFirst: fmt.Sprintf("%v", resetAfterFirst),
	}, "", "")
	if err != nil {
		return err
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	return handleProxyResponse(client.Do(req))
}

func AddGeneralRegexSanitizer(value, regex string, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.hostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "GeneralRegexSanitizer")

	marshalled, err := json.MarshalIndent(struct {
		Value           string `json:"value"`
		Regex           string `json:"regex"`
		GroupForReplace string `json:"groupForReplace"`
	}{
		Value:           value,
		Regex:           regex,
		GroupForReplace: options.GroupForReplace,
	}, "", "")
	if err != nil {
		return err
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	return handleProxyResponse(client.Do(req))
}

func AddHeaderRegexSanitizer(key, replacementValue, regex string, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.hostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "HeaderRegexSanitizer")

	marshalled, err := json.MarshalIndent(struct {
		Key             string `json:"key"`
		Value           string `json:"value,omitempty"`
		Regex           string `json:"regex,omitempty"`
		GroupForReplace string `json:"groupForReplace,omitempty"`
	}{
		Key:             key,
		Value:           replacementValue,
		Regex:           regex,
		GroupForReplace: options.GroupForReplace,
	}, "", "")
	if err != nil {
		return err
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	return handleProxyResponse(client.Do(req))
}

func AddOAuthResponseSanitizer(options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.hostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "OAuthResponseSanitizer")
	return handleProxyResponse(client.Do(req))
}

func AddRemoveHeaderSanitizer(headersForRemoval []string, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.hostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "RemoveHeaderSanitizer")

	marshalled, err := json.MarshalIndent(struct {
		HeadersForRemoval string `json:"headersForRemoval"`
	}{
		HeadersForRemoval: strings.Join(headersForRemoval, ""),
	}, "", "")
	if err != nil {
		return err
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	return handleProxyResponse(client.Do(req))
}

func AddURISanitizer(replacement, regex string, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%v/Admin/AddSanitizer", options.hostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "UriRegexSanitizer")

	marshalled, err := json.MarshalIndent(struct {
		Value string `json:"value"`
		Regex string `json:"regex"`
	}{
		Value: replacement,
		Regex: regex,
	}, "", "")
	if err != nil {
		return err
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	return handleProxyResponse(client.Do(req))
}

func AddURISubscriptionIDSanitizer(value string, options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.hostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "UriSubscriptionIdSanitizer")

	if value != "" {
		marshalled, err := json.MarshalIndent(struct {
			Value string `json:"value,omitempty"`
		}{
			Value: value,
		}, "", "")
		if err != nil {
			return err
		}
		
		req.Body = ioutil.NopCloser(bytes.NewReader(marshalled))
		req.ContentLength = int64(len(marshalled))
	}
	return handleProxyResponse(client.Do(req))
}

func ResetSanitizers(options *RecordingOptions) error {
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%v/Admin/Reset", options.hostScheme())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	return handleProxyResponse(client.Do(req))
}
