// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
)

// RequestLogOptions configures the logging policy's behavior.
type RequestLogOptions struct {
	// placeholder for future configuration options
}

type requestLogPolicy struct {
	options RequestLogOptions
}

// NewRequestLogPolicy creates a RequestLogPolicy object configured using the specified options.
func NewRequestLogPolicy(o *RequestLogOptions) Policy {
	return &requestLogPolicy{}
}

// logPolicyOpValues is the struct containing the per-operation values
type logPolicyOpValues struct {
	try   int32
	start time.Time
}

func (p *requestLogPolicy) Do(req *Request) (*Response, error) {
	// Get the per-operation values. These are saved in the Message's map so that they persist across each retry calling into this policy object.
	var opValues logPolicyOpValues
	if req.OperationValue(&opValues); opValues.start.IsZero() {
		opValues.start = time.Now() // If this is the 1st try, record this operation's start time
	}
	opValues.try++ // The first try is #1 (not #0)
	req.SetOperationValue(opValues)

	// Log the outgoing request as informational
	if Log().Should(LogRequest) {
		b := &bytes.Buffer{}
		fmt.Fprintf(b, "==> OUTGOING REQUEST (Try=%d)\n", opValues.try)
		WriteRequestWithResponse(b, prepareRequestForLogging(req), nil, nil)
		Log().Write(LogRequest, b.String())
	}

	// Set the time for this particular retry operation and then Do the operation.
	tryStart := time.Now()
	response, err := req.Next() // Make the request
	tryEnd := time.Now()
	tryDuration := tryEnd.Sub(tryStart)
	opDuration := tryEnd.Sub(opValues.start)

	if Log().Should(LogResponse) {
		// We're going to log this; build the string to log
		b := &bytes.Buffer{}
		fmt.Fprintf(b, "==> REQUEST/RESPONSE (Try=%d/%v, OpTime=%v) -- ", opValues.try, tryDuration, opDuration)
		if err != nil { // This HTTP request did not get a response from the service
			fmt.Fprint(b, "REQUEST ERROR\n")
		} else {
			fmt.Fprint(b, "RESPONSE RECEIVED\n")
		}

		WriteRequestWithResponse(b, prepareRequestForLogging(req), response, err)
		if err != nil {
			// skip frames runtime.Callers() and runtime.StackTrace()
			b.WriteString(runtime.StackTrace(2, StackFrameCount))
		}
		Log().Write(LogResponse, b.String())
	}
	return response, err
}

// RedactSigQueryParam redacts the 'sig' query parameter in URL's raw query to protect secret.
func RedactSigQueryParam(rawQuery string) (bool, string) {
	rawQuery = strings.ToLower(rawQuery) // lowercase the string so we can look for ?sig= and &sig=
	sigFound := strings.Contains(rawQuery, "?sig=")
	if !sigFound {
		sigFound = strings.Contains(rawQuery, "&sig=")
		if !sigFound {
			return sigFound, rawQuery // [?|&]sig= not found; return same rawQuery passed in (no memory allocation)
		}
	}
	// [?|&]sig= found, redact its value
	values, _ := url.ParseQuery(rawQuery)
	for name := range values {
		if strings.EqualFold(name, "sig") {
			values[name] = []string{"REDACTED"}
		}
	}
	return sigFound, values.Encode()
}

func prepareRequestForLogging(req *Request) *Request {
	request := req
	if sigFound, rawQuery := RedactSigQueryParam(request.URL.RawQuery); sigFound {
		// Make copy so we don't destroy the query parameters we actually need to send in the request
		request = req.copy()
		request.URL.RawQuery = rawQuery
	}
	return request
}
