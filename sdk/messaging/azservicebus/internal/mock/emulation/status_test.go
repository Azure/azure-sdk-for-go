// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package emulation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStatus(t *testing.T) {
	top := NewStatus(nil)

	top.CloseWithError(errors.New("Hello World"))
	top.CloseWithError(errors.New("Hello World 2"))

	require.EqualError(t, top.Err(), "Hello World")

	child := NewStatus(top)

	select {
	case <-child.Done():
		require.EqualError(t, child.Err(), "Hello World")
	default:
		require.Fail(t, "Should have been closed already")
	}
}
