// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"fmt"

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

func SetClassifications(cls ...LogClassification) {
	input := make([]logger.LogClassification, 0)
	for _, l := range cls {
		input = append(input, logger.LogClassification(l))
	}
	logger.Log().SetClassifications(input...)
}

type Listener func(LogClassification, string)

func transform(lst Listener) logger.Listener {
	return func(l logger.LogClassification, msg string) {
		azcoreL := LogClassification(l)
		lst(azcoreL, msg)
	}
}

func SetListener(lst Listener) {
	if lst == nil {
		fmt.Println("nil listener")
		logger.Log().SetListener(nil)
	} else {
		fmt.Println("Not nil listener")
		logger.Log().SetListener(transform(lst))
	}
}

// for testing purposes
func resetClassifications() {
	logger.Log().SetClassifications([]logger.LogClassification{}...)
}
