//go:build go1.18 && (darwin || linux || windows)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// the test in this file must be defined in azidentity_test because it imports azidentity/cache

package azidentity_test

import (
	"context"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func TestPersistentCacheLive(t *testing.T) {
	if recording.GetRecordMode() != recording.LiveMode {
		t.Skip("this test runs only in live mode")
	}
	if runtime.GOOS == "darwin" && os.Getenv("AZIDENTITY_RUN_MANUAL_TESTS") == "" {
		t.Skip("set AZIDENTITY_RUN_MANUAL_TESTS to run this test on macOS")
	}
	armURL := os.Getenv("RESOURCE_MANAGER_URL")
	if armURL == "" {
		t.Skip("set RESOURCE_MANAGER_URL to run this test")
	}
	tro := policy.TokenRequestOptions{Scopes: []string{armURL + "/.default"}}
	for _, test := range []struct {
		credential func(*testing.T, azidentity.AuthenticationRecord, azidentity.Cache) (azcore.TokenCredential, error)
		name       string
	}{
		{
			credential: func(t *testing.T, _ azidentity.AuthenticationRecord, c azidentity.Cache) (azcore.TokenCredential, error) {
				t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22879")
				clientID := os.Getenv("IDENTITY_SP_CLIENT_ID")
				secret := os.Getenv("IDENTITY_SP_CLIENT_SECRET")
				tenantID := os.Getenv("IDENTITY_SP_TENANT_ID")
				if clientID == "" || secret == "" || tenantID == "" {
					t.Skip("set IDENTITY_SP_* with service principal configuration to run this test")
				}
				return azidentity.NewClientSecretCredential(tenantID, clientID, secret,
					&azidentity.ClientSecretCredentialOptions{Cache: c},
				)
			},
			name: "confidential",
		},
		{
			credential: func(t *testing.T, r azidentity.AuthenticationRecord, c azidentity.Cache) (azcore.TokenCredential, error) {
				clientID := "04b07795-8ddb-461a-bbee-02f9e1bf7b46"
				password := os.Getenv("AZURE_IDENTITY_TEST_PASSWORD")
				tenantID := os.Getenv("AZURE_IDENTITY_TEST_TENANTID")
				username := os.Getenv("AZURE_IDENTITY_TEST_USERNAME")
				if password == "" || tenantID == "" || username == "" {
					t.Skip("set AZURE_IDENTITY_TEST_* with user configuration to run this test")
				}
				return azidentity.NewUsernamePasswordCredential(tenantID, clientID, username, password,
					&azidentity.UsernamePasswordCredentialOptions{
						AuthenticationRecord: r,
						Cache:                c,
					},
				)
			},
			name: "public",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			c, err := cache.New(&cache.Options{Name: strings.ReplaceAll(t.Name(), "/", "_")})
			require.NoError(t, err)

			rec := azidentity.AuthenticationRecord{}
			cred, err := test.credential(t, rec, c)
			require.NoError(t, err)
			if test.name == "public" {
				type authenticater interface {
					Authenticate(context.Context, *policy.TokenRequestOptions) (azidentity.AuthenticationRecord, error)
				}
				a, ok := cred.(authenticater)
				require.True(t, ok, "test bug: public credential must implement Authenticate")
				rec, err = a.Authenticate(ctx, &tro)
				require.NoError(t, err)
			}
			tk, err := cred.GetToken(ctx, tro)
			require.NoError(t, err)

			cred2, err := test.credential(t, rec, c)
			require.NoError(t, err)
			tk2, err := cred2.GetToken(ctx, tro)
			require.NoError(t, err)
			// require.Equal is more to the point but prints a value i.e. logs a token when expected != actual
			require.True(t, tk.Token == tk2.Token, "expected a cached token")

			caeTRO := tro
			caeTRO.EnableCAE = true
			tk3, err := cred.GetToken(ctx, caeTRO)
			require.NoError(t, err)
			require.False(t, tk.Token == tk3.Token, "expected a new token because the cached one isn't a CAE token")

			tk4, err := cred2.GetToken(ctx, caeTRO)
			require.NoError(t, err)
			require.True(t, tk3.Token == tk4.Token, "expected a cached token")
			require.False(t, tk.Token == tk4.Token, "expected a CAE token")
		})
	}
}
