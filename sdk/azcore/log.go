// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/logger"
)

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

func SetClassifications(cls ...logger.LogClassification) {
	logger.Log().SetClassifications(cls)
}

func SetListener(lst logger.Listener) {
	logger.Log().SetListener(lst)
}
