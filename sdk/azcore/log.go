// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/logger"
)

// LogClassification is used to group entries.  Each group can be toggled on or off
type LogClassification logger.LogClassification

const (
	// LogRequest entries contain information about HTTP requests.
	// This includes information like the URL, query parameters, and headers.
	LogRequest LogClassification = "Request"

	// LogResponse entries contain information about HTTP responses.
	// This includes information like the HTTP status code, headers, and request URL.
	LogResponse LogClassification = "Response"

	// LogRetryPolicy entries contain information specific to the retry policy in use.
	LogRetryPolicy LogClassification = "RetryPolicy"

	// LogLongRunningOperation entries contain information specific to long-running operations.
	// This includes information like polling location, operation state and sleep intervals.
	LogLongRunningOperation LogClassification = "LongRunningOperation"
)

// SetClassifications is used to control which classifications are written to
// the log.  By default all log classifications are writen.
func SetClassifications(cls ...LogClassification) {
	input := make([]logger.LogClassification, 0)
	for _, l := range cls {
		input = append(input, logger.LogClassification(l))
	}
	logger.Log().SetClassifications(input...)
}

// Listener is the function signature invoked when writing log entries.
// A Listener is required to perform its own synchronization if it's expected to be called
// from multiple Go routines
type Listener func(LogClassification, string)

// transform to convert the azcore.Listener type into a usable one for internal.logger module
func transform(lst Listener) logger.Listener {
	return func(l logger.LogClassification, msg string) {
		azcoreL := LogClassification(l)
		lst(azcoreL, msg)
	}
}

// SetListener will set the Logger to write to the specified Listener.
func SetListener(lst Listener) {
	if lst == nil {
		logger.Log().SetListener(nil)
	} else {
		logger.Log().SetListener(transform(lst))
	}
}

// for testing purposes
func resetClassifications() {
	logger.Log().SetClassifications([]logger.LogClassification{}...)
}
