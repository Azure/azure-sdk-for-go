//go:build !wasm

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package azeventgrid

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultTransportDialContext(t *testing.T) {
	// On non-WASM targets the helper returns the dialer's DialContext.
	require.NotNil(t, defaultTransportDialContext(&net.Dialer{}))
}
