//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package errorinfo

import (
	"errors"
	"testing"
)

func TestNonRetriableError(t *testing.T) {
	const dnr string = "Do Not Retry"

	// Create sample error.
	err := NonRetriableError(errors.New(dnr))

	// Check error message is correct
	if err.Error() != dnr {
		t.Fatalf("Expected error message to be '%q' but got '%q'.", dnr, err.Error())
	}

	var e *nonRetriableError
	if !errors.As(err, &e) {
		t.Fatalf("Expected error to be of type nonRetriableError")
	}
}
