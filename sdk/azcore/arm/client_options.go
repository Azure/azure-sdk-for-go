//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"

// ClientOptions contains configuration settings for a client's pipeline.
type ClientOptions struct {
	policy.ClientOptions

	// DisableRPRegistration disables the auto-RP registration policy. Defaults to false.
	DisableRPRegistration bool
}
