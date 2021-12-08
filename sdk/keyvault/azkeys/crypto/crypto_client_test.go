//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package crypto

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	for _, key := range []string{"https://mykeyvault.vault.azure.net/keys/keyabcdef/1234567890", "https://mykeyvault.vault.azure.net/keys/keyabcdef"} {
		cred := getCredential(t)

		client, err := NewClient(key, cred, nil)
		require.NoError(t, err)

		require.Equal(t, client.vaultURL, "https://mykeyvault.vault.azure.net/")
		require.Equal(t, "keyabcdef", client.keyID)
		if strings.Contains(key, "1234567890") {
			require.Equal(t, client.keyVersion, "1234567890")
		} else {
			require.Equal(t, client.keyVersion, "")
		}
	}
}
