//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package internal

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestParseID(t *testing.T) {
	examples := map[string]struct{ url, name, version *string }{
		"https://myvaultname.vault.azure.net/keys/key1053998307/b86c2e6ad9054f4abf69cc185b99aa60": {to.Ptr("https://myvaultname.vault.azure.net/"), to.Ptr("key1053998307"), to.Ptr("b86c2e6ad9054f4abf69cc185b99aa60")},
		"https://myvaultname.vault.azure.net/keys/key1053998307":                                  {to.Ptr("https://myvaultname.vault.azure.net/"), to.Ptr("key1053998307"), nil},
		"https://myvaultname.vault.azure.net/":                                                    {to.Ptr("https://myvaultname.vault.azure.net/"), nil, nil},
	}

	for url, result := range examples {
		url, name, version := ParseID(&url)
		if result.url == nil {
			require.Nil(t, url)
		} else {
			require.NotNil(t, url)
			require.Equal(t, *url, *result.url)
		}
		if result.name == nil {
			require.Nil(t, name)
		} else {
			require.NotNilf(t, name, "expected %s", *result.name)
			require.Equal(t, *name, *result.name)
		}
		if result.version == nil {
			require.Nil(t, version)
		} else {
			require.NotNil(t, version)
			require.Equal(t, *version, *result.version)
		}
	}
}
