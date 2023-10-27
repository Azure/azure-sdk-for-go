//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"testing"
)

func TestNonRetriableError(t *testing.T) {
	//Create sample error.
	err := NonRetriableError(NewSampleError("Do Not Retry"))

	// Check if the error is of type NonRetriableError
	if _, ok := err.(*nonRetriableError); !ok {
		t.Fatal("Expected error to be of type NonRetriableError")
	}

	// Check error message is correct
	errMsg := "Do Not Retry"
	if err.Error() != errMsg {
		t.Fatalf("Expected error message to be '%q' but got '%q'.", errMsg, err.Error())
	}
}

func NewSampleError(message string) error {
	return &sampleError{message}
}

type sampleError struct {
	message string
}

func (e *sampleError) Error() string {
	return e.message
}
