//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/generated"
)

// CryptoClient is a wrapper around generated.KeyVaultClient with additional data
// needed by crypto.Client.  This allows construction of a crypto.Client from the
// azkeys and crypto packages without having to use a bunch of type aliases.
type CryptoClient struct {
	kvClient   *generated.KeyVaultClient
	vaultURL   string
	keyName    string
	keyVersion string
}

// NewCryptoClient creates a new CryptoClient with the specified values.
func NewCryptoClient(vaultURL, keyName, keyVersion string, pl runtime.Pipeline) CryptoClient {
	return CryptoClient{
		kvClient:   generated.NewKeyVaultClient(pl),
		vaultURL:   vaultURL,
		keyName:    keyName,
		keyVersion: keyVersion,
	}
}

// Client returns the underlying generated client.
func Client(client CryptoClient) *generated.KeyVaultClient {
	return client.kvClient
}

// VaultURL returns the vault URL for this client.
func VaultURL(client CryptoClient) string {
	return client.vaultURL
}

// KeyName returns the key name for this client.
func KeyName(client CryptoClient) string {
	return client.keyName
}

// KeyVersion returns the key version for this client (can be empty).
func KeyVersion(client CryptoClient) string {
	return client.keyVersion
}
