//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

// NewResponseError creates a new *ResponseError from the provided HTTP response.
func NewResponseError(resp *http.Response) error {
	respErr := &ResponseError{
		StatusCode:  resp.StatusCode,
		RawResponse: resp,
	}

	// prefer the error code in the response header
	if ec := resp.Header.Get("x-ms-error-code"); ec != "" {
		respErr.ErrorCode = ec
	}

	// write the request method and URL with response status code
	msg := &bytes.Buffer{}
	fmt.Fprintf(msg, "%s %s://%s%s\n", resp.Request.Method, resp.Request.URL.Scheme, resp.Request.URL.Host, resp.Request.URL.Path)
	fmt.Fprintln(msg, "--------------------------------------------------------------------------------")
	fmt.Fprintf(msg, "RESPONSE %d: %s\n", resp.StatusCode, resp.Status)

	body, err := Payload(resp)
	if err != nil {
		body = []byte(fmt.Sprintf("Error reading response body: %v", err))
		goto Epilog
	}

	// if we didn't get x-ms-error-code, check in the response body
	if respErr.ErrorCode == "" && len(body) > 0 {
		if code := extractErrorCodeJSON(body); code != "" {
			respErr.ErrorCode = code
		} else if code := extractErrorCodeXML(body); code != "" {
			respErr.ErrorCode = code
		}
	}

Epilog:
	if respErr.ErrorCode != "" {
		fmt.Fprintf(msg, "ERROR CODE: %s\n", respErr.ErrorCode)
	} else {
		fmt.Fprintln(msg, "ERROR CODE UNAVAILABLE")
	}

	fmt.Fprintln(msg, "--------------------------------------------------------------------------------")
	if len(body) > 0 {
		if err := json.Indent(msg, body, "", "  "); err != nil {
			// failed to pretty-print so just dump it verbatim
			fmt.Fprint(msg, string(body))
		}
		// the standard library doesn't have a pretty-printer for XML
		fmt.Fprintln(msg)
	} else {
		fmt.Fprintln(msg, "Response contained no body")
	}
	fmt.Fprintln(msg, "--------------------------------------------------------------------------------")

	respErr.msg = msg.String()
	return respErr
}

func extractErrorCodeJSON(body []byte) string {
	var rawObj map[string]interface{}
	if err := json.Unmarshal(body, &rawObj); err != nil {
		// not a JSON object
		return ""
	}

	// check if this is a wrapped error, i.e. { "error": { ... } }
	// if so then unwrap it
	if wrapped, ok := rawObj["error"]; ok {
		unwrapped, ok := wrapped.(map[string]interface{})
		if !ok {
			return ""
		}
		rawObj = unwrapped
	}

	// now check for the error code
	code, ok := rawObj["code"]
	if !ok {
		return ""
	}
	codeStr, ok := code.(string)
	if !ok {
		return ""
	}
	return codeStr
}

func extractErrorCodeXML(body []byte) string {
	// regular expression is much easier than dealing with the XML parser
	rx := regexp.MustCompile(`<[c|C]ode>\s*(\w+)\s*<\/[c|C]ode>`)
	res := rx.FindStringSubmatch(string(body))
	if len(res) != 2 {
		return ""
	}
	// first submatch is the entire thing, second one is the captured error code
	return res[1]
}

// ResponseError is returned when a request is made to a service and
// the service returns a non-success HTTP status code.
// Use errors.As() to access this type in the error chain.
type ResponseError struct {
	msg string

	// ErrorCode is the error code returned by the resource provider if available.
	ErrorCode string

	// StatusCode is the HTTP status code as defined in https://pkg.go.dev/net/http#pkg-constants.
	StatusCode int

	// RawResponse is the underlying HTTP response.
	RawResponse *http.Response
}

// Error implements the error interface for type ResponseError.
// Note that the message contents are not contractual and can change over time.
func (e *ResponseError) Error() string {
	return e.msg
}
