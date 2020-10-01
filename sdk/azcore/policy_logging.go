// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
)

// LogOptions configures the logging policy's behavior.
type LogOptions struct {
	// placeholder for future configuration options
}

// DefaultLogOptions returns an instance of LogOptions initialized with default values.
func DefaultLogOptions() LogOptions {
	return LogOptions{}
}

type logPolicy struct {
	options LogOptions
}

// NewLogPolicy creates a RequestLogPolicy object configured using the specified options.
// Pass nil to accept the default values; this is the same as passing the result
// from a call to DefaultLogOptions().
func NewLogPolicy(o *LogOptions) Policy {
	if o == nil {
		def := DefaultLogOptions()
		o = &def
	}
	return &logPolicy{options: *o}
}

// logPolicyOpValues is the struct containing the per-operation values
type logPolicyOpValues struct {
	try   int32
	start time.Time
}

func (p *logPolicy) Do(req *Request) (*Response, error) {
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
		writeRequestWithResponse(b, req, nil, nil)
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

		writeRequestWithResponse(b, req, response, err)
		if err != nil {
			// skip frames runtime.Callers() and runtime.StackTrace()
			b.WriteString(runtime.StackTrace(2, StackFrameCount))
		}
		Log().Write(LogResponse, b.String())
	}
	return response, err
}
