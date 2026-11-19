// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azsecrets_test

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tokenCredentialFunc func(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error)

func (c tokenCredentialFunc) GetToken(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return c(ctx, tro)
}

func TestClientMTLSPoP(t *testing.T) {
	for _, mode := range []string{"bound", "bearer", "tenant fallback", "tenant transition"} {
		t.Run(mode, func(t *testing.T) {
			var boundRequests atomic.Int32
			server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") == "" {
					assert.Empty(t, r.TLS.PeerCertificates)
					w.Header().Set("WWW-Authenticate", `Bearer authorization="https://login.microsoftonline.com/tenant", resource="https://vault.azure.net"`)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				if mode == "bound" || (mode == "tenant transition" && r.Header.Get("Authorization") == "mtls_pop bound") {
					assert.Equal(t, "mtls_pop bound", r.Header.Get("Authorization"))
					assert.Equal(t, "true", r.Header.Get("x-ms-tokenboundauth"))
					assert.Len(t, r.TLS.PeerCertificates, 1)
					if mode == "tenant transition" && boundRequests.Add(1)%2 == 0 {
						w.Header().Set("WWW-Authenticate", `Bearer authorization="https://login.microsoftonline.com/tenant2", resource="https://vault.azure.net"`)
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
				} else {
					assert.Equal(t, "Bearer bearer", r.Header.Get("Authorization"))
					assert.Empty(t, r.Header.Get("x-ms-tokenboundauth"))
					assert.Empty(t, r.TLS.PeerCertificates)
				}
				w.Header().Set("Content-Type", "application/json")
				_, err := io.WriteString(w, `{"value":"test value"}`)
				assert.NoError(t, err)
			}))
			server.TLS = &tls.Config{ClientAuth: tls.RequestClientCert}
			server.StartTLS()
			defer server.Close()
			var modes []bool
			leaf := tokenCredentialFunc(func(_ context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
				modes = append(modes, tro.EnableMTLSPoP)
				if mode != "tenant transition" {
					require.Equal(t, "tenant", tro.TenantID)
				}
				require.Equal(t, []string{"https://vault.azure.net/.default"}, tro.Scopes)
				if (mode == "tenant fallback" || (mode == "tenant transition" && tro.TenantID == "tenant2")) && tro.EnableMTLSPoP {
					return azcore.AccessToken{}, &azidentity.AuthenticationFailedError{RawResponse: &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       io.NopCloser(strings.NewReader(`{"error":"invalid_request","error_description":"AADSTS3921996: tenant not eligible"}`)),
					}}
				}
				tk := azcore.AccessToken{Token: "bearer", ExpiresOn: time.Now().Add(time.Hour)}
				if mode == "bound" || (mode == "tenant transition" && tro.TenantID == "tenant") {
					require.True(t, tro.EnableMTLSPoP)
					tk.Token, tk.BindingCertificate = "bound", &server.TLS.Certificates[0]
				}
				return tk, nil
			})
			chain, err := azidentity.NewChainedTokenCredential([]azcore.TokenCredential{
				tokenCredentialFunc(func(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
					return azcore.AccessToken{}, azidentity.NewCredentialUnavailableError("not configured")
				}), leaf,
			}, nil)
			require.NoError(t, err)
			for clientIndex := range 2 {
				client, err := azsecrets.NewClient(server.URL, chain, &azsecrets.ClientOptions{
					ClientOptions:                        azcore.ClientOptions{Transport: server.Client()},
					DisableChallengeResourceVerification: true,
				})
				require.NoError(t, err)
				reads := 2 - clientIndex
				if mode == "tenant transition" {
					reads++
				}
				for range reads {
					res, err := client.GetSecret(context.Background(), "test", "", nil)
					require.NoError(t, err)
					require.NotNil(t, res.Value)
					require.Equal(t, "test value", *res.Value)
				}
			}
			if mode == "tenant transition" {
				require.Equal(t, []bool{true, true, false, false, true, true, false}, modes,
					"a late tenant denial must clear the prior binding certificate and authorization headers")
			} else if mode == "tenant fallback" {
				require.Equal(t, []bool{true, false, false, true, false}, modes,
					"fresh bearer acquisition and tenant memory must be isolated to each SDK client")
			} else {
				require.Equal(t, []bool{true, true}, modes, "each client must signal PoP automatically and cache the actual returned token")
			}
		})
	}
}
