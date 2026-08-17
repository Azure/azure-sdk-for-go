//go:build !wasm

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package azeventgrid

import (
	"context"
	"net"
)

// defaultTransportDialContext returns the dialer's DialContext. It mirrors
// azcore's build-tagged helper (sdk/azcore/runtime/transport_default_dialer_other.go)
// so that newDefaultTransport matches azcore's default transport on every
// platform.
func defaultTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}
