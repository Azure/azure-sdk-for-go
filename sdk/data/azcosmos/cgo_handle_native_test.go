// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"runtime/cgo"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCGOHandleHelpersTolerateInvalidHandles(t *testing.T) {
	handle := cgo.NewHandle("value")
	value, valid := cgoHandleValue(handle)
	require.True(t, valid)
	require.Equal(t, "value", value)

	cgoHandleDelete(handle)
	require.NotPanics(t, func() {
		_, valid = cgoHandleValue(handle)
	})
	require.False(t, valid)
	require.NotPanics(t, func() { cgoHandleDelete(handle) })
}
