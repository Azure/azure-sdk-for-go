//go:build (js && wasm) || wasip1

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package azeventgrid

import (
	"context"
	"net"
)

// defaultTransportDialContext returns nil on WebAssembly targets, mirroring
// azcore's build-tagged helper (sdk/azcore/runtime/transport_default_dialer_wasm.go).
// A net.Dialer cannot be used on these platforms, so http.Transport must fall
// back to its own default dialing.
func defaultTransportDialContext(_ *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return nil
}
