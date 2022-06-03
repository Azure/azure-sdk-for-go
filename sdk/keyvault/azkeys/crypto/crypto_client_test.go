//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConstructor(t *testing.T) {
	client, err := NewClient("https://fakekvurl.vault.azure.net/keys/key89075156/0b29f1d3760f4407aeb996868c9a02a7", &FakeCredential{}, nil)
	require.NoError(t, err)
	require.NotNil(t, client.client())
	require.Equal(t, client.keyID(), "key89075156")
	require.Equal(t, client.keyVersion(), "0b29f1d3760f4407aeb996868c9a02a7")

	client, err = NewClient("https://fakekvurl.vault.azure.net/keys/key89075156", &FakeCredential{}, nil)
	require.NoError(t, err)
	require.NotNil(t, client.client())
	require.Equal(t, client.keyID(), "key89075156")
	require.Equal(t, client.keyVersion(), "")

	_, err = NewClient("https://fakekvurl.vault.azure.net/", &FakeCredential{}, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "URL is not for a specific key, expect path to start with '/keys/', received")
}
