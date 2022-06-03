//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/stretchr/testify/require"
)

func TestCryptoClient(t *testing.T) {
	const (
		vaultURL   = "vault"
		keyName    = "name"
		keyVersion = "version"
	)
	pl := runtime.NewPipeline("test", "v0.0.1", runtime.PipelineOptions{}, nil)
	client := NewCryptoClient(vaultURL, keyName, keyVersion, pl)
	require.Equal(t, vaultURL, VaultURL(client))
	require.Equal(t, keyName, KeyName(client))
	require.Equal(t, keyVersion, KeyVersion(client))
	require.NotZero(t, Client(client))
}
