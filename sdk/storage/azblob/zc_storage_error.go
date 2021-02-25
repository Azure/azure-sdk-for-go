// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
)

// // StorageError identifies a responder-generated network or response parsing error.
// type StorageError interface {
// 	// ResponseError implements error's Error(), net.Error's Temporary() and Timeout() methods & Response().
// 	// ResponseError
//
// 	// ServiceCode returns a service error code. Your code can use this to make error recovery decisions.
// 	ServiceCode() ServiceCodeType
// }

// InternalError is an internal error type that all errors get wrapped in.
type InternalError struct {
	cause error
}

func (e InternalError) Error() string {
	if (errors.Is(e.cause, StorageError{})) {
		return e.cause.Error()
	}

	return fmt.Sprintf("===== INTERNAL ERROR =====\n%s", e.cause.Error())
}

func (e InternalError) Is(err error) bool {
	_, ok := err.(InternalError)

	return ok
}

// StorageError is the internal struct that replaces the generated StorageError.
// TL;DR: This implements xml.Unmarshaler, and when the original StorageError is substituted, this unmarshaler kicks in.
// This handles the description and details. defunkifyStorageError handles the response, cause, and service code.
type StorageError struct {
	response *http.Response
	description string

	ServiceCode ServiceCodeType
	details     map[string]string
}

func handleError(err error) error {
	if err, ok := err.(*runtime.ResponseError); ok {
		return InternalError{defunkifyStorageError(err) }
	}

	if err != nil {
		return InternalError{err}
	}

	return nil
}

// defunkifyStorageError is a function that takes the "funky" *runtime.ResponseError and reduces it to a storageError.
func defunkifyStorageError(responseError *runtime.ResponseError) error {
	if err, ok := responseError.Unwrap().(*StorageError); ok {
		// errors.Unwrap(responseError.Unwrap())

		err.response = responseError.RawResponse()

		err.ServiceCode = ServiceCodeType(responseError.RawResponse().Header.Get("x-ms-error-code"))

		if code, ok := err.details["Code"]; ok {
			err.ServiceCode = ServiceCodeType(code)
			delete(err.details, "Code")
		}

		return err
	}

	return InternalError{
		cause: responseError,
	}
}

// ServiceCode returns service-error information. The caller may examine these values but should not modify any of them.
func (e *StorageError) ServiceCode() ServiceCodeType {
	return e.serviceCode
}

// ServiceCode returns service-error information. The caller may examine these values but should not modify any of them.
func (e *StorageError) StatusCode() int {
	return e.response.StatusCode
}

// Error implements the error interface's Error method to return a string representation of the error.
func (e StorageError) Error() string {
	b := &bytes.Buffer{}

	if e.response != nil {
		fmt.Fprintf(b, "===== RESPONSE ERROR (ServiceCode=%s) =====\n", e.ServiceCode)
		fmt.Fprintf(b, "Description=%s, Details: ", e.description)
		if len(e.details) == 0 {
			b.WriteString("(none)\n")
		} else {
			b.WriteRune('\n')
			keys := make([]string, 0, len(e.details))
			// Alphabetize the details
			for k := range e.details {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Fprintf(b, "   %s: %+v\n", k, e.details[k])
			}
		}
		// req := azcore.Request{Request: e.response.Request}.Copy() // Make a copy of the response's request
		writeRequestWithResponse(b, &azcore.Request{Request: e.response.Request}, e.response)
	}

	return b.String()
	// azcore.writeRequestWithResponse(b, prepareRequestForLogging(req), e.response, nil)
	// return e.ErrorNode.Error(b.String())
}

func (e StorageError) Is(err error) bool {
	_, ok := err.(StorageError)

	return ok
}

func (e StorageError) Response() *http.Response {
	return e.response
}

func writeRequestWithResponse(b *bytes.Buffer, request *azcore.Request, response *http.Response) {
	// Write the request into the buffer.
	fmt.Fprint(b, "   "+request.Method+" "+request.URL.String()+"\n")
	writeHeader(b, request.Header)
	if response != nil {
		fmt.Fprintln(b, "   --------------------------------------------------------------------------------")
		fmt.Fprint(b, "   RESPONSE Status: "+response.Status+"\n")
		writeHeader(b, response.Header)
	}
}

// formatHeaders appends an HTTP request's or response's header into a Buffer.
func writeHeader(b *bytes.Buffer, header map[string][]string) {
	if len(header) == 0 {
		b.WriteString("   (no headers)\n")
		return
	}
	keys := make([]string, 0, len(header))
	// Alphabetize the headers
	for k := range header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		// Redact the value of any Authorization header to prevent security information from persisting in logs
		value := interface{}("REDACTED")
		if !strings.EqualFold(k, "Authorization") {
			value = header[k]
		}
		fmt.Fprintf(b, "   %s: %+v\n", k, value)
	}
}

// Temporary returns true if the error occurred due to a temporary condition (including an HTTP status of 500 or 503).
func (e *StorageError) Temporary() bool {
	if e.response != nil {
		if (e.response.StatusCode == http.StatusInternalServerError) || (e.response.StatusCode == http.StatusServiceUnavailable) || (e.response.StatusCode == http.StatusBadGateway) {
			return true
		}
	}
  
	return false
}

// UnmarshalXML performs custom unmarshalling of XML-formatted Azure storage request errors.
func (e *StorageError) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	tokName := ""
	var t xml.Token
	for t, err = d.Token(); err == nil; t, err = d.Token() {
		switch tt := t.(type) {
		case xml.StartElement:
			tokName = tt.Name.Local
			break
		case xml.CharData:
			switch tokName {
			case "Message":
				e.description = string(tt)
			default:
				if e.details == nil {
					e.details = map[string]string{}
				}
				e.details[tokName] = string(tt)
			}
		}
	}
	return nil
}
