// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

// RequestLogOptions configures the retry policy's behavior.
type RequestLogOptions struct {
	// LogWarningIfTryOverThreshold logs a warning if a tried operation takes longer than the specified
	// duration (-1=no logging; 0=default threshold).
	LogWarningIfTryOverThreshold time.Duration
}

func (o RequestLogOptions) defaults() RequestLogOptions {
	if o.LogWarningIfTryOverThreshold == 0 {
		// It would be good to relate this to https://azure.microsoft.com/en-us/support/legal/sla/storage/v1_2/
		// But this monitors the time to get the HTTP response; NOT the time to download the response body.
		o.LogWarningIfTryOverThreshold = 3 * time.Second // Default to 3 seconds
	}
	return o
}

type requestLogPolicy struct {
	options RequestLogOptions
}

// NewRequestLogPolicy creates a RequestLogPolicy object configured using the specified options.
func NewRequestLogPolicy(o RequestLogOptions) Policy {
	o = o.defaults() // Force defaults to be calculated
	return &requestLogPolicy{options: o}
}

// logPolicyOpValues is the struct containing the per-operation values
type logPolicyOpValues struct {
	try   int32
	start time.Time
}

func (p *requestLogPolicy) Do(ctx context.Context, req *Request) (*Response, error) {
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
	response, err := req.Next(ctx) // Make the request
	tryEnd := time.Now()
	tryDuration := tryEnd.Sub(tryStart)
	opDuration := tryEnd.Sub(opValues.start)

	logClass := LogResponse // Default logging information

	// If the response took too long, we'll upgrade to warning.
	if p.options.LogWarningIfTryOverThreshold > 0 && tryDuration > p.options.LogWarningIfTryOverThreshold {
		// Log a warning if the try duration exceeded the specified threshold
		logClass = LogSlowResponse
	}

	if err == nil { // We got a response from the service
		sc := response.StatusCode
		if ((sc >= 400 && sc <= 499) && sc != http.StatusNotFound && sc != http.StatusConflict && sc != http.StatusPreconditionFailed && sc != http.StatusRequestedRangeNotSatisfiable) || (sc >= 500 && sc <= 599) {
			logClass = LogError // Promote to Error any 4xx (except those listed is an error) or any 5xx
		} else {
			// For other status codes, we leave the level as is.
		}
	} else { // This error did not get an HTTP response from the service; upgrade the severity to Error
		logClass = LogError
	}

	if Log().Should(logClass) {
		// We're going to log this; build the string to log
		b := &bytes.Buffer{}
		slow := ""
		if p.options.LogWarningIfTryOverThreshold > 0 && tryDuration > p.options.LogWarningIfTryOverThreshold {
			slow = fmt.Sprintf("[SLOW >%v]", p.options.LogWarningIfTryOverThreshold)
		}
		fmt.Fprintf(b, "==> REQUEST/RESPONSE (Try=%d/%v%s, OpTime=%v) -- ", opValues.try, tryDuration, slow, opDuration)
		if err != nil { // This HTTP request did not get a response from the service
			fmt.Fprint(b, "REQUEST ERROR\n")
		} else {
			if logClass == LogError {
				fmt.Fprint(b, "RESPONSE STATUS CODE ERROR\n")
			} else {
				fmt.Fprint(b, "RESPONSE SUCCESSFULLY RECEIVED\n")
			}
		}

		WriteRequestWithResponse(b, prepareRequestForLogging(req), response, err)
		if logClass == LogError {
			b.Write(stack()) // For errors (or lower levels), we append the stack trace (an expensive operation)
		}
		Log().Write(logClass, b.String())
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

func stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}
