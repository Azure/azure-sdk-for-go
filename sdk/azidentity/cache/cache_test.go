//go:build linux || windows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cache

import (
	"context"
	"fmt"
	"io"
	"net/http"
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

func TestSeparation(t *testing.T) {
	before := cacheDir
	defer func() { cacheDir = before }()
	td := t.TempDir()
	cacheDir = func() (string, error) { return td, nil }

	names := []string{}
	for i := 0; i < 3; i++ {
		names = append(names, fmt.Sprintf("%s-%d", t.Name(), i))
	}
	test := func(t *testing.T, enableCAE bool) {
		tro := policy.TokenRequestOptions{EnableCAE: enableCAE, Scopes: []string{"scope"}}
		// create a cache for each name, containing a unique token for that name
		for _, name := range names {
			expected := name
			if enableCAE {
				expected += "-CAE"
			}
			c, err := New(&Options{Name: name})
			require.NoError(t, err)
			require.NotNil(t, c)
			o := azidentity.ClientSecretCredentialOptions{
				Cache: c,
				ClientOptions: policy.ClientOptions{
					Transport: &mockSTS{
						tokenRequestCallback: func(*http.Request) *http.Response {
							body := fmt.Sprintf(`{"access_token":%q,"expires_in":3600}`, expected)
							return &http.Response{
								StatusCode: 200,
								Header:     http.Header{"Content-Type": []string{"application/json"}},
								Body:       io.NopCloser(strings.NewReader(body)),
							}
						},
					},
				},
			}
			cred, err := azidentity.NewClientSecretCredential("tenantID", "clientID", "secret", &o)
			require.NoError(t, err)
			actual, err := cred.GetToken(ctx, tro)
			require.NoError(t, err)
			require.Equal(t, expected, actual.Token)
		}
		// verify the caches contain the expected tokens
		for _, name := range names {
			expected := name
			if tro.EnableCAE {
				expected += "-CAE"
			}
			c, err := New(&Options{Name: name})
			require.NoError(t, err)
			require.NotNil(t, c)
			o := azidentity.ClientSecretCredentialOptions{
				Cache: c,
				ClientOptions: policy.ClientOptions{
					Transport: &mockSTS{
						tokenRequestCallback: func(*http.Request) *http.Response {
							t.Error("credential should have found a cached token")
							return nil
						},
					},
				},
			}
			cred, err := azidentity.NewClientSecretCredential("tenantID", "clientID", "secret", &o)
			require.NoError(t, err)
			actual, err := cred.GetToken(ctx, tro)
			require.NoError(t, err)
			require.Equal(t, expected, actual.Token)
		}
	}

	// caches having different names shouldn't share data
	test(t, false)
	// caches having the same name should separate CAE and non-CAE tokens
	test(t, true)
	test(t, false)
}
