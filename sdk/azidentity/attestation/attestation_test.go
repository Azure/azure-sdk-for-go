// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package attestation_test

import (
	"context"
	"crypto/rand"
	"net/http"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	_ "github.com/Azure/azure-sdk-for-go/sdk/azidentity/attestation"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportedAttestation(t *testing.T) {
	clearManagedIdentityEnvironment(t)
	for _, tc := range []struct {
		name                string
		requestPoP, disable bool
	}{
		{name: "discover unsupported host", requestPoP: true},
		{name: "opt out", requestPoP: true, disable: true},
		{name: "ordinary bearer"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			wantPaths := []string{"/metadata/identity/oauth2/token"}
			if tc.requestPoP && !tc.disable {
				// A definitive v2 absence and Standard VM metadata avoid native operations on every platform.
				srv.AppendResponse(mock.WithStatusCode(http.StatusNotFound))
				srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
				srv.AppendResponse(mock.WithBody([]byte(`{"securityProfile":{"securityType":"Standard"}}`)))
				wantPaths = []string{
					"/metadata/identity/getplatformmetadata", "/metadata/identity/oauth2/token",
					"/metadata/instance/compute", "/metadata/identity/oauth2/token",
				}
			}
			srv.AppendResponse(mock.WithBody([]byte(`{"access_token":"bearer-token","expires_in":3600,"token_type":"Bearer"}`)))
			var paths []string
			cred, err := azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
				DisableMTLSPoP: tc.disable,
				ClientOptions: azcore.ClientOptions{
					Transport: srv,
					Retry:     policy.RetryOptions{MaxRetries: -1},
					PerCallPolicies: []policy.Policy{policyFunc(func(req *policy.Request) (*http.Response, error) {
						paths = append(paths, req.Raw().URL.Path)
						return req.Next()
					})},
				},
			})
			require.NoError(t, err)
			assert.Empty(t, paths, "constructing an activated credential must not perform discovery")
			tk, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
				Scopes: []string{"https://example.com/" + rand.Text()}, EnableMTLSPoP: tc.requestPoP,
			})
			require.NoError(t, err)
			assert.Equal(t, "bearer-token", tk.Token)
			assert.Nil(t, tk.BindingCertificate)
			assert.Equal(t, wantPaths, paths, "the blank import must activate capability discovery only for eligible requests")
		})
	}
}

func TestImportedAttestationWithKeyVault(t *testing.T) {
	clearManagedIdentityEnvironment(t)
	t.Setenv("AZURE_TOKEN_CREDENTIALS", "ManagedIdentityCredential")
	t.Setenv("AZURE_CLIENT_ID", "")
	require.NoError(t, os.Unsetenv("AZURE_CLIENT_ID"))
	for _, tc := range []struct {
		name         string
		dac, disable bool
	}{
		{name: "managed identity"},
		{name: "default credential", dac: true},
		{name: "managed identity opt out", disable: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			identity, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			wantPaths := []string{"/metadata/identity/oauth2/token"}
			if !tc.disable {
				identity.AppendResponse(mock.WithStatusCode(http.StatusNotFound))
				identity.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
				identity.AppendResponse(mock.WithBody([]byte(`{"securityProfile":{"securityType":"Standard"}}`)))
				wantPaths = []string{
					"/metadata/identity/getplatformmetadata", "/metadata/identity/oauth2/token",
					"/metadata/instance/compute", "/metadata/identity/oauth2/token",
				}
			}
			identity.AppendResponse(mock.WithBody([]byte(`{"access_token":"bearer-token","expires_in":3600,"token_type":"Bearer"}`)))
			var paths []string
			options := azcore.ClientOptions{
				Transport: identity,
				Retry:     policy.RetryOptions{MaxRetries: -1},
				PerCallPolicies: []policy.Policy{policyFunc(func(req *policy.Request) (*http.Response, error) {
					paths = append(paths, req.Raw().URL.Path)
					return req.Next()
				})},
			}
			var cred azcore.TokenCredential
			var err error
			if tc.dac {
				cred, err = azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{ClientOptions: options})
			} else {
				cred, err = azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
					ClientOptions: options, DisableMTLSPoP: tc.disable,
				})
			}
			require.NoError(t, err)
			assert.Empty(t, paths, "construction must not discover capabilities or load native dependencies")

			resource, closeResource := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer closeResource()
			resource.AppendResponse(mock.WithStatusCode(http.StatusUnauthorized), mock.WithHeader("WWW-Authenticate",
				`Bearer authorization="https://login.microsoftonline.com/tenant" resource="https://vault.azure.net/`+rand.Text()+`"`))
			resource.AppendResponse(mock.WithBody([]byte(`{"value":"test-value"}`)))
			client, err := azsecrets.NewClient("https://test.vault.azure.net", cred, &azsecrets.ClientOptions{
				ClientOptions:                        azcore.ClientOptions{Transport: resource},
				DisableChallengeResourceVerification: true,
			})
			require.NoError(t, err)
			response, err := client.GetSecret(context.Background(), "test-secret", "", nil)
			require.NoError(t, err)
			require.NotNil(t, response.Value)
			assert.Equal(t, "test-value", *response.Value)
			assert.Equal(t, wantPaths, paths, "default SDK signaling must select discovery, unless explicitly opted out")
		})
	}
}

func clearManagedIdentityEnvironment(t *testing.T) {
	t.Helper()
	for _, name := range []string{
		"IDENTITY_ENDPOINT", "IDENTITY_HEADER", "IDENTITY_SERVER_THUMBPRINT", "IMDS_ENDPOINT",
		"MSI_ENDPOINT", "MSI_SECRET", "AZURE_POD_IDENTITY_AUTHORITY_HOST", "MSAL_MI_DISABLE_IMDS_V2",
	} {
		t.Setenv(name, "")
	}
}

type policyFunc func(*policy.Request) (*http.Response, error)

func (f policyFunc) Do(req *policy.Request) (*http.Response, error) {
	return f(req)
}
