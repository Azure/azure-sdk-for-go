// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// LogClassification is used to group entries.  Each group can be toggled on or off.
type LogClassification = log.Classification

const (
	// LogRequest entries contain information about HTTP requests.
	// This includes information like the URL, query parameters, and headers.
	LogRequest = log.Request

	// LogResponse entries contain information about HTTP responses.
	// This includes information like the HTTP status code, headers, and request URL.
	LogResponse = log.Response

	// LogRetryPolicy entries contain information specific to the retry policy in use.
	LogRetryPolicy = log.RetryPolicy

	// LogLongRunningOperation entries contain information specific to long-running operations.
	// This includes information like polling location, operation state and sleep intervals.
	LogLongRunningOperation = log.LongRunningOperation
)

// LogSetClassifications is used to control which classifications are written to
// the log.  By default all log classifications are writen.
func LogSetClassifications(cls ...LogClassification) {
	log.SetClassifications(cls...)
}

// SetListener will set the Logger to write to the specified Listener.
func LogSetListener(lst func(log.Classification, string)) {
	log.SetListener(lst)
}

// for testing purposes
func resetClassifications() {
	log.TestResetClassifications()
}
