//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package errorinfo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeError struct {
	Message string
}

func (fe *fakeError) Error() string {
	return fe.Message
}

func TestNonRetriableError(t *testing.T) {
	const dnr string = "Do Not Retry"

	// Create sample error.
	err := NonRetriableError(&fakeError{Message: dnr})
	// Check error message is correct
	require.Error(t, err, dnr)

	var e NonRetriable
	require.ErrorAs(t, err, &e)

	// Check Unwrap method on NonRetriable error type
	var fe *fakeError
	require.ErrorAs(t, err, &fe)
}
