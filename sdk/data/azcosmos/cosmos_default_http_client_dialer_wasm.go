//go:build (js && wasm) || wasip1

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net"
)

// defaultCosmosTransportDialContext mirrors azcore's defaultTransportDialContext
// so the Cosmos default HTTP transport behaves consistently with azcore across
// build targets. On WASM/wasip1 it returns nil, which lets the runtime use the
// platform-specific HTTP transport instead of a net.Dialer-based DialContext.
func defaultCosmosTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return nil
}
