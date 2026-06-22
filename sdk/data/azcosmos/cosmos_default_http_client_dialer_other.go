//go:build !wasm

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net"
)

// defaultCosmosTransportDialContext mirrors azcore's defaultTransportDialContext
// so the Cosmos default HTTP transport behaves consistently with azcore across
// build targets. On non-WASM platforms it returns the dialer's DialContext.
func defaultCosmosTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}
