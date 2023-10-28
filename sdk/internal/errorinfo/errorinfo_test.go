//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package errorinfo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNonRetriableError(t *testing.T) {
	//Create sample error.
	err := NonRetriableError(errors.New("Do Not Retry"))

	// Check error message is correct
	errMsg := "Do Not Retry"
	if err.Error() != errMsg {
		t.Fatalf("Expected error message to be '%q' but got '%q'.", errMsg, err.Error())
	}

	var nonRetriableErr interface{ NonRetriable() }
	require.True(t, errors.As(&nonRetriableError{}, &nonRetriableErr))
}
