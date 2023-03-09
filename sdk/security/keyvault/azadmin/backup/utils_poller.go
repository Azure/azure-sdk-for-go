//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package backup

// this file contains handwritten additions to the generated code
// code to support the custom poller handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// the well-known set of LRO status/provisioning state values.
const (
	statusSucceeded  = "Succeeded"
	statusCanceled   = "Canceled"
	statusFailed     = "Failed"
	statusInProgress = "InProgress"
)

// isTerminalState returns true if the LRO's state is terminal.
func isTerminalState(s string) bool {
	return strings.EqualFold(s, statusSucceeded) || strings.EqualFold(s, statusFailed) || strings.EqualFold(s, statusCanceled)
}

// failed returns true if the LRO's state is terminal failure.
func failed(s string) bool {
	return strings.EqualFold(s, statusFailed) || strings.EqualFold(s, statusCanceled)
}

// returns true if the LRO response contains a valid HTTP status code
func statusCodeValid(resp *http.Response) bool {
	return hasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusCreated, http.StatusNoContent)
}

// isValidURL verifies that the URL is valid and absolute.
func isValidURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.IsAbs()
}

// errNoBody is returned if the response didn't contain a body.
var errNoBody = errors.New("the response did not contain a body")

// getJSON reads the response body into a raw JSON object.
// It returns ErrNoBody if there was no content.
func getJSON(resp *http.Response) (map[string]any, error) {
	body, err := payload(resp)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, errNoBody
	}
	// unmarshall the body to get the value
	var jsonBody map[string]any
	if err = json.Unmarshal(body, &jsonBody); err != nil {
		return nil, err
	}
	return jsonBody, nil
}

// provisioningState returns the provisioning state from the response or the empty string.
func provisioningState(jsonBody map[string]any) string {
	jsonProps, ok := jsonBody["properties"]
	if !ok {
		return ""
	}
	props, ok := jsonProps.(map[string]any)
	if !ok {
		return ""
	}
	rawPs, ok := props["provisioningState"]
	if !ok {
		return ""
	}
	ps, ok := rawPs.(string)
	if !ok {
		return ""
	}
	return ps
}

// status returns the status from the response or the empty string.
func status(jsonBody map[string]any) string {
	rawStatus, ok := jsonBody["status"]
	if !ok {
		return ""
	}
	status, ok := rawStatus.(string)
	if !ok {
		return ""
	}
	return status
}

// getStatus returns the LRO's status from the response body.
// Typically used for Azure-AsyncOperation flows.
// If there is no status in the response body the empty string is returned.
func getStatus(resp *http.Response) (string, error) {
	jsonBody, err := getJSON(resp)
	if err != nil {
		return "", err
	}
	return status(jsonBody), nil
}

// getProvisioningState returns the LRO's state from the response body.
// If there is no state in the response body the empty string is returned.
func getProvisioningState(resp *http.Response) (string, error) {
	jsonBody, err := getJSON(resp)
	if err != nil {
		return "", err
	}
	return provisioningState(jsonBody), nil
}

// hasStatusCode returns true if the Response's status code is one of the specified values.
// Exported as runtime.HasStatusCode().
func hasStatusCode(resp *http.Response, statusCodes ...int) bool {
	if resp == nil {
		return false
	}
	for _, sc := range statusCodes {
		if resp.StatusCode == sc {
			return true
		}
	}
	return false
}

// payload reads and returns the response body or an error.
// On a successful read, the response body is cached.
// Subsequent reads will access the cached value.
// Exported as runtime.Payload().
func payload(resp *http.Response) ([]byte, error) {
	// r.Body won't be a nopClosingBytesReader if downloading was skipped
	if buf, ok := resp.Body.(*nopClosingBytesReader); ok {
		return buf.Bytes(), nil
	}
	bytesBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	resp.Body = &nopClosingBytesReader{s: bytesBody}
	return bytesBody, nil
}

// nopClosingBytesReader is an io.ReadSeekCloser around a byte slice.
// It also provides direct access to the byte slice to avoid rereading.
type nopClosingBytesReader struct {
	s []byte
	i int64
}

// Bytes returns the underlying byte slice.
func (r *nopClosingBytesReader) Bytes() []byte {
	return r.s
}

// Close implements the io.Closer interface.
func (*nopClosingBytesReader) Close() error {
	return nil
}

// Read implements the io.Reader interface.
func (r *nopClosingBytesReader) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return
}

// Set replaces the existing byte slice with the specified byte slice and resets the reader.
func (r *nopClosingBytesReader) Set(b []byte) {
	r.s = b
	r.i = 0
}

// Seek implements the io.Seeker interface.
func (r *nopClosingBytesReader) Seek(offset int64, whence int) (int64, error) {
	var i int64
	switch whence {
	case io.SeekStart:
		i = offset
	case io.SeekCurrent:
		i = r.i + offset
	case io.SeekEnd:
		i = int64(len(r.s)) + offset
	default:
		return 0, errors.New("nopClosingBytesReader: invalid whence")
	}
	if i < 0 {
		return 0, errors.New("nopClosingBytesReader: negative position")
	}
	r.i = i
	return i, nil
}
