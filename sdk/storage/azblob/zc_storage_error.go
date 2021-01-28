package azblob

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net"
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

// StorageError is the internal struct that replaces the generated StorageError.
// TL;DR: This implements xml.Unmarshaler, and when the original StorageError is substituted, this unmarshaler kicks in.
// This handles the description and details. defunkifyStorageError handles the response, cause, and service code.
type StorageError struct {
	cause error //
	response *http.Response
	description string

	serviceCode ServiceCodeType
	details     map[string]string
}

func handleError(err error) error {
	if err, ok := err.(*runtime.ResponseError); ok {
		return defunkifyStorageError(err)
	}

	return err
}

// defunkifyStorageError is a function that takes the "funky" *runtime.ResponseError and reduces it to a storageError.
func defunkifyStorageError(responseError *runtime.ResponseError) error {
	if err, ok := responseError.Unwrap().(*StorageError); ok {
		// errors.Unwrap(responseError.Unwrap())

		err.response = responseError.RawResponse()
		err.cause = nil

		if code, ok := err.details["Code"]; ok {
			err.serviceCode = ServiceCodeType(code)
			delete(err.details, "Code")
		}

		return err
	}

	return responseError
}

// // newStorageError creates an error object that implements the error interface.
// func newStorageError(cause error, response *http.Response, description string) error {
// 	return &StorageError{
// 		cause: cause,
// 		response: response,
// 		description: description,
//
// 		serviceCode: ServiceCodeType(response.Header.Get("x-ms-error-code")),
// 	}
// }

// ServiceCode returns service-error information. The caller may examine these values but should not modify any of them.
func (e *StorageError) ServiceCode() ServiceCodeType {
	return e.serviceCode
}

// Error implements the error interface's Error method to return a string representation of the error.
func (e *StorageError) Error() string {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "===== RESPONSE ERROR (ServiceCode=%s) =====\n", e.serviceCode)
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
	WriteRequestWithResponse(b, &azcore.Request{Request: e.response.Request}, e.response, e.cause)
	return b.String()
	// azcore.WriteRequestWithResponse(b, prepareRequestForLogging(req), e.response, nil)
	// return e.ErrorNode.Error(b.String())
}

func WriteRequestWithResponse(b *bytes.Buffer, request *azcore.Request, response *http.Response, err error) {
	// Write the request into the buffer.
	fmt.Fprint(b, "   "+request.Method+" "+request.URL.String()+"\n")
	writeHeader(b, request.Header)
	if response != nil {
		fmt.Fprintln(b, "   --------------------------------------------------------------------------------")
		fmt.Fprint(b, "   RESPONSE Status: "+response.Status+"\n")
		writeHeader(b, response.Header)
	}
	if err != nil {
		fmt.Fprintln(b, "   --------------------------------------------------------------------------------")
		fmt.Fprint(b, "   ERROR:\n"+err.Error()+"\n")
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

	if netError, ok := e.cause.(net.Error); ok {
		return netError.Temporary()
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
