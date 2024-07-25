//go:build go1.18 && (linux || windows)
// +build go1.18
// +build linux windows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cache

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func TestCache(t *testing.T) {
	before := cacheDir
	t.Cleanup(func() { cacheDir = before })
	cacheDir = func() (string, error) { return t.TempDir(), nil }
	for _, test := range []struct {
		credential func(azidentity.Cache) (azcore.TokenCredential, error)
		name       string
	}{
		{
			func(c azidentity.Cache) (azcore.TokenCredential, error) {
				opts := azidentity.ClientSecretCredentialOptions{
					Cache:         c,
					ClientOptions: policy.ClientOptions{Transport: &mockSTS{}},
				}
				return azidentity.NewClientSecretCredential("tenantID", "clientID", "secret", &opts)
			},
			"confidential",
		},
		{
			func(c azidentity.Cache) (azcore.TokenCredential, error) {
				opts := azidentity.DeviceCodeCredentialOptions{
					Cache:         c,
					ClientOptions: policy.ClientOptions{Transport: &mockSTS{}},
				}
				return azidentity.NewDeviceCodeCredential(&opts)
			},
			"public",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			name := strings.ReplaceAll(t.Name(), string(filepath.Separator), "_")
			cache, err := New(&Options{Name: name})
			require.NoError(t, err)
			cred, err := test.credential(cache)
			require.NoError(t, err)
			tro := policy.TokenRequestOptions{Scopes: []string{"scope"}}
			tk, err := cred.GetToken(ctx, tro)
			require.NoError(t, err)

			cred2, err := test.credential(cache)
			require.NoError(t, err)
			tk2, err := cred2.GetToken(ctx, tro)
			require.NoError(t, err)
			require.Equal(t, tk.Token, tk2.Token)
		})
	}
}
