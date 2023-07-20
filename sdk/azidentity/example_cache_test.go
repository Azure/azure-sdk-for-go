//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	// importing the cache module registers the cache implementation for the current platform
	_ "github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
)

// Credentials, excepting those that authenticate via external tools like [AzureCLICredential],
// cache authentication data in memory by default. Most of these credentials also support optional
// persistent caching. This example shows how to enable and configure that for a credential. It
// shows only [InteractiveBrowserCredential], however all credentials that support persistent caching have
// the same [TokenCachePersistenceOptions] API.
func Example_persistentCache() {
	cred, err := azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		// Non-nil TokenCachePersistenceOptions enables persistent caching with default options.
		// See TokenCachePersistenceOptions documentation for details of the supported options.
		TokenCachePersistenceOptions: &azidentity.TokenCachePersistenceOptions{},
	})
	if err != nil {
		// TODO: handle error
	}
	// TODO: use credential
	_ = cred
}
