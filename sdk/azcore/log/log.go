//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package log provides functionality for configuring logging facilities.
package log

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// Classification is used to group entries.  Each group can be toggled on or off.
type Classification = log.Classification

const (
	// Request entries contain information about HTTP requests.
	// This includes information like the URL, query parameters, and headers.
	Request = log.Request

	// Response entries contain information about HTTP responses.
	// This includes information like the HTTP status code, headers, and request URL.
	Response = log.Response

	// RetryPolicy entries contain information specific to the retry policy in use.
	RetryPolicy = log.RetryPolicy

	// LongRunningOperation entries contain information specific to long-running operations.
	// This includes information like polling location, operation state and sleep intervals.
	LongRunningOperation = log.LongRunningOperation
)

// SetClassifications is used to control which classifications are written to
// the log.  By default all log classifications are writen.
func SetClassifications(cls ...Classification) {
	log.SetClassifications(cls...)
}

// SetListener will set the Logger to write to the specified Listener.
func SetListener(lst func(log.Classification, string)) {
	log.SetListener(lst)
}

// for testing purposes
func resetClassifications() {
	log.TestResetClassifications()
}
