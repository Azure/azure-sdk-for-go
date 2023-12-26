//go:build go1.18 && windows
// +build go1.18,windows

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

func TestCaching(t *testing.T) {
	for _, test := range []struct {
		ctor func(azidentity.TokenCachePersistenceOptions) (azcore.TokenCredential, error)
		name string
	}{
		{
			func(tcpo azidentity.TokenCachePersistenceOptions) (azcore.TokenCredential, error) {
				opts := azidentity.ClientSecretCredentialOptions{
					ClientOptions:                policy.ClientOptions{Transport: &mockSTS{}},
					TokenCachePersistenceOptions: &tcpo,
				}
				return azidentity.NewClientSecretCredential("tenantID", "clientID", "secret", &opts)
			},
			"confidential",
		},
		{
			func(tcpo azidentity.TokenCachePersistenceOptions) (azcore.TokenCredential, error) {
				opts := azidentity.DeviceCodeCredentialOptions{
					ClientOptions:                policy.ClientOptions{Transport: &mockSTS{}},
					TokenCachePersistenceOptions: &tcpo,
				}
				return azidentity.NewDeviceCodeCredential(&opts)
			},
			"public",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			tcpo := azidentity.TokenCachePersistenceOptions{
				Name: strings.ReplaceAll(t.Name(), string(filepath.Separator), "_"),
			}
			if a, e := storage(azidentity.TokenCachePersistenceOptions{Name: tcpo.Name + ".nocae"}); e == nil {
				defer func() { a.Delete(ctx) }()
			}
			cred, err := test.ctor(tcpo)
			require.NoError(t, err)
			tro := policy.TokenRequestOptions{Scopes: []string{"scope"}}
			tk, err := cred.GetToken(ctx, tro)
			require.NoError(t, err)

			cred2, err := test.ctor(tcpo)
			require.NoError(t, err)
			tk2, err := cred2.GetToken(ctx, tro)
			require.NoError(t, err)
			require.Equal(t, tk.Token, tk2.Token)
		})
	}
}
