// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && windows && amd64

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWindowsAMD64SelectsNativeDriver(t *testing.T) {
	require.True(t, driverAvailable)
	require.NoError(t, verifyDriverVersion())
}
