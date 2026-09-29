// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build !cgo || !((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuntimeUnavailableDiagnostic(t *testing.T) {
	runtime, err := NewRuntime(nil)
	require.Nil(t, runtime)
	require.ErrorContains(t, err, "requires CGO_ENABLED=1")
	var zero Runtime
	require.NoError(t, zero.Close())
	require.NoError(t, zero.Close())
	require.Error(t, zero.SetOperationOptions(OperationOptions{}))
}
