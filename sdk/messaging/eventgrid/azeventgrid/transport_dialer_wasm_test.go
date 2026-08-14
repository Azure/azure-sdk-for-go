//go:build (js && wasm) || wasip1

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package azeventgrid

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultTransportDialContext(t *testing.T) {
	// On WASM targets the helper returns nil so http.Transport uses its own
	// default dialing (a net.Dialer cannot be used there).
	require.Nil(t, defaultTransportDialContext(&net.Dialer{}))
}
