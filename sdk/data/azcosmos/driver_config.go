// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

// driverConfig is what a [Client] needs to open its driver resources. It exists so that the two
// driver implementations, the native binding and the stub, share one input type and client.go
// stays free of build tags.
type driverConfig struct {
	// endpoint is the account endpoint, already validated by newClient.
	endpoint string

	// accountKey is the master key, already validated by [NewKeyCredential].
	accountKey string

	// tokenCredential is set instead of accountKey for Microsoft Entra ID authentication.
	tokenCredential azcore.TokenCredential

	options ClientOptions
}

// usesTokenCredential reports whether the client was created with a token credential rather than
// an account key.
//
//nolint:unused // consumed by the driver-backed build in driver_native.go.
func (c driverConfig) usesTokenCredential() bool {
	return c.tokenCredential != nil
}
