//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package recording

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Deprecated: the local recording API that uses this type is no longer supported. Use [Start] and [Stop] to make
// recordings with the test proxy instead, and configure its sanitizers with functions such as [AddBodyKeySanitizer]
// and [AddRemoveHeaderSanitizer].
type Sanitizer struct{}

// StringSanitizer is a func that will modify the string pointed to by the parameter into a sanitized value.
//
// Deprecated: the local sanitizer API that uses this type is no longer supported. Use test proxy sanitizer
// configuration functions such as [AddBodyKeySanitizer] and [AddRemoveHeaderSanitizer] instead.
type StringSanitizer func(*string)

// SanitizedValue is the default placeholder value to be used for sanitized strings.
const SanitizedValue string = "Sanitized"

// SanitizedBase64Value is the default placeholder value to be used for sanitized base-64 encoded strings.
const SanitizedBase64Value string = "Kg=="

// AddSanitizedHeaders adds the supplied header names to the list of headers to be sanitized on request and response recordings.
func (s *Sanitizer) AddSanitizedHeaders(headers ...string) {
	panic(errUnsupportedAPI)
}

// AddBodysanitizer configures the supplied StringSanitizer to sanitize recording request and response bodies
func (s *Sanitizer) AddBodysanitizer(sanitizer StringSanitizer) {
	panic(errUnsupportedAPI)
}

// AddUriSanitizer configures the supplied StringSanitizer to sanitize recording request and response URLs
func (s *Sanitizer) AddUrlSanitizer(sanitizer StringSanitizer) {
	panic(errUnsupportedAPI)
}

// Deprecated: the local sanitizer API that uses this function is no longer supported.
func DefaultStringSanitizer(s *string) {}

func handleProxyResponse(resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	return fmt.Errorf("%s responded %s: %s", resp.Request.URL, resp.Status, body)
}

// handleTestLevelSanitizer sets the "x-recording-id" header if options.TestInstance is not nil
func handleTestLevelSanitizer(req *http.Request, options *RecordingOptions) {
	if options == nil || options.TestInstance == nil {
		return
	}

	if recordingID := GetRecordingId(options.TestInstance); recordingID != "" {
		req.Header.Set(IDHeader, recordingID)
	}
}

// AddBodyKeySanitizer adds a sanitizer for JSON Bodies. jsonPath is the path to the key, value
// is the value to replace with, and regex is the string to match in the body. If your regex includes a group
// options.GroupForReplace specifies which group to replace
func AddBodyKeySanitizer(jsonPath, value, regex string, options *RecordingOptions) error {
	if recordMode == LiveMode {
		return nil
	}
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.baseURL())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "BodyKeySanitizer")
	handleTestLevelSanitizer(req, options)

	marshalled, err := json.MarshalIndent(struct {
		JSONPath        string `json:"jsonPath"`
		Value           string `json:"value,omitempty"`
		Regex           string `json:"regex,omitempty"`
		GroupForReplace string `json:"groupForReplace,omitempty"`
	}{
		JSONPath:        jsonPath,
		Value:           value,
		Regex:           regex,
		GroupForReplace: options.GroupForReplace,
	}, "", "")
	if err != nil {
		return err
	}

	req.Body = io.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	return handleProxyResponse(client.Do(req))
}

// AddBodyRegexSanitizer offers regex replace within a returned JSON body. value is the
// substitution value, regex can be a simple regex or a substitution operation if
// options.GroupForReplace is set.
func AddBodyRegexSanitizer(value, regex string, options *RecordingOptions) error {
	if recordMode == LiveMode {
		return nil
	}
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.baseURL())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "BodyRegexSanitizer")
	handleTestLevelSanitizer(req, options)

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

	req.Body = io.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	return handleProxyResponse(client.Do(req))
}

// AddContinuationSanitizer is used to anonymize private keys in response/request pairs.
// key: the name of the header whos value will be replaced from response -> next request
// method: the method by which the value of the targeted key will be replaced. Defaults to GUID replacement
// resetAfterFirt: Do we need multiple pairs replaced? Or do we want to replace each value with the same value.
func AddContinuationSanitizer(key, method string, resetAfterFirst bool, options *RecordingOptions) error {
	if recordMode == LiveMode {
		return nil
	}
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.baseURL())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "ContinuationSanitizer")
	handleTestLevelSanitizer(req, options)

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

	req.Body = io.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	return handleProxyResponse(client.Do(req))
}

// AddGeneralRegexSanitizer adds a general regex across request/response Body, Headers, and URI.
// value is the substitution value, regex can be defined as a simple regex replace or a substition
// operation if options.GroupForReplace specifies which group to replace.
func AddGeneralRegexSanitizer(value, regex string, options *RecordingOptions) error {
	if recordMode == LiveMode {
		return nil
	}
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.baseURL())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "GeneralRegexSanitizer")
	handleTestLevelSanitizer(req, options)

	marshalled, err := json.MarshalIndent(struct {
		Value           string `json:"value"`
		Regex           string `json:"regex"`
		GroupForReplace string `json:"groupForReplace,omitempty"`
	}{
		Value:           value,
		Regex:           regex,
		GroupForReplace: options.GroupForReplace,
	}, "", "")
	if err != nil {
		return err
	}

	req.Body = io.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	return handleProxyResponse(client.Do(req))
}

// AddHeaderRegexSanitizer can be used to replace a key with a specific value: set regex to ""
// OR can be used to do a simple regex replace operation by setting key, value, and regex.
// OR To do a substitution operation if options.GroupForReplace is set.
// key is the name of the header to operate against. value is the substitution or whole new header
// value. regex can be defined as a simple regex replace or a substitution operation if
// options.GroupForReplace is set.
func AddHeaderRegexSanitizer(key, value, regex string, options *RecordingOptions) error {
	if recordMode == LiveMode {
		return nil
	}
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.baseURL())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "HeaderRegexSanitizer")
	handleTestLevelSanitizer(req, options)

	marshalled, err := json.MarshalIndent(struct {
		Key             string `json:"key"`
		Value           string `json:"value,omitempty"`
		Regex           string `json:"regex,omitempty"`
		GroupForReplace string `json:"groupForReplace,omitempty"`
	}{
		Key:             key,
		Value:           value,
		Regex:           regex,
		GroupForReplace: options.GroupForReplace,
	}, "", "")
	if err != nil {
		return err
	}

	req.Body = io.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	return handleProxyResponse(client.Do(req))
}

// AddOAuthResponseSanitizer cleans all request/response pairs taht match an oauth regex in their URI
func AddOAuthResponseSanitizer(options *RecordingOptions) error {
	if recordMode == LiveMode {
		return nil
	}
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.baseURL())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "OAuthResponseSanitizer")
	handleTestLevelSanitizer(req, options)

	return handleProxyResponse(client.Do(req))
}

// AddRemoveHeaderSanitizer removes a list of headers from request/responses.
func AddRemoveHeaderSanitizer(headersForRemoval []string, options *RecordingOptions) error {
	if recordMode == LiveMode {
		return nil
	}
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.baseURL())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "RemoveHeaderSanitizer")
	handleTestLevelSanitizer(req, options)

	if options.TestInstance != nil {
		recordingID := GetRecordingId(options.TestInstance)
		if recordingID == "" {
			return fmt.Errorf("did not find a recording ID for test with name '%s'. Did you make sure to call Start?", options.TestInstance.Name())
		}
		req.Header.Set(IDHeader, recordingID)
	}

	marshalled, err := json.MarshalIndent(struct {
		HeadersForRemoval string `json:"headersForRemoval"`
	}{
		HeadersForRemoval: strings.Join(headersForRemoval, ""),
	}, "", "")
	if err != nil {
		return err
	}

	req.Body = io.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	return handleProxyResponse(client.Do(req))
}

// AddURISanitizer sanitizes URIs via regex. value is the substition value, regex is
// either a simple regex or a substitution operation if options.GroupForReplace is defined.
func AddURISanitizer(value, regex string, options *RecordingOptions) error {
	if recordMode == LiveMode {
		return nil
	}
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%v/Admin/AddSanitizer", options.baseURL())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "UriRegexSanitizer")
	handleTestLevelSanitizer(req, options)

	marshalled, err := json.MarshalIndent(struct {
		Value string `json:"value"`
		Regex string `json:"regex"`
	}{
		Value: value,
		Regex: regex,
	}, "", "")
	if err != nil {
		return err
	}

	req.Body = io.NopCloser(bytes.NewReader(marshalled))
	req.ContentLength = int64(len(marshalled))
	return handleProxyResponse(client.Do(req))
}

// AddURISubscriptionIDSanitizer replaces real subscriptionIDs within a URI with a default
// or configured fake value. To use the default value set value to "", otherwise value specifies the replacement value.
func AddURISubscriptionIDSanitizer(value string, options *RecordingOptions) error {
	if recordMode == LiveMode {
		return nil
	}
	if options == nil {
		options = defaultOptions()
	}
	url := fmt.Sprintf("%s/Admin/AddSanitizer", options.baseURL())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("x-abstraction-identifier", "UriSubscriptionIdSanitizer")
	handleTestLevelSanitizer(req, options)

	if value != "" {
		marshalled, err := json.MarshalIndent(struct {
			Value string `json:"value,omitempty"`
		}{
			Value: value,
		}, "", "")
		if err != nil {
			return err
		}

		req.Body = io.NopCloser(bytes.NewReader(marshalled))
		req.ContentLength = int64(len(marshalled))
	}
	return handleProxyResponse(client.Do(req))
}

// RemoveRegisteredSanitizers selectively disables sanitizers that are enabled by default such as "AZSDK1001"
func RemoveRegisteredSanitizers(sanitizerIDs []string, options *RecordingOptions) error {
	if recordMode == LiveMode {
		return nil
	}
	if options == nil {
		options = defaultOptions()
	}
	b, err := json.Marshal(struct {
		Sanitizers []string `json:"Sanitizers"`
	}{
		Sanitizers: sanitizerIDs,
	})
	if err != nil {
		return err
	}
	url := options.baseURL() + "/Admin/RemoveSanitizers"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	handleTestLevelSanitizer(req, options)
	return handleProxyResponse(client.Do(req))
}

// ResetProxy restores the proxy's default sanitizers, matchers, and transforms
func ResetProxy(options *RecordingOptions) error {
	if recordMode == LiveMode {
		return nil
	}
	if options == nil {
		options = defaultOptions()
	}

	url := fmt.Sprintf("%v/Admin/Reset", options.baseURL())
	req, err := http.NewRequest("POST", url, nil)

	if options.TestInstance != nil {
		recordingID := GetRecordingId(options.TestInstance)
		if recordingID == "" {
			return fmt.Errorf("did not find a recording ID for test with name '%s'. Did you make sure to call Start?", options.TestInstance.Name())
		}
		req.Header.Set(IDHeader, recordingID)
	}

	if err != nil {
		return err
	}
	return handleProxyResponse(client.Do(req))
}
