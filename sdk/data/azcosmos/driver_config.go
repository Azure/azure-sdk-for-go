// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// driverConfig is what a [Client] needs to open its driver resources. It exists so that the two
// driver implementations, the native binding and the stub, share one input type and client.go
// stays free of build tags.
type driverConfig struct {
	// endpoint is the account endpoint, already validated by newClient.
	endpoint string

	// accountKey is the master key, already validated by [NewKeyCredential]. It is empty when the
	// caller supplied a token credential, which the driver cannot yet accept; see openDriver.
	accountKey string

	options ClientOptions
}

// usesTokenCredential reports whether the client was created with a token credential rather than
// an account key.
//
//nolint:unused // consumed by the driver-backed build in driver_native.go.
func (c driverConfig) usesTokenCredential() bool {
	return c.accountKey == ""
}
