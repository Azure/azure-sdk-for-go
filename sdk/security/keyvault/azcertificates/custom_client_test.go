// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

import (
	"context"
	"encoding/json"
	"io"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestNewPlatformManaged(t *testing.T) {
	metadata := map[string]any{
		"issuer":       "internal-ca",
		"rotationDays": 90,
	}

	pm := NewPlatformManaged("tls-server", metadata)

	require.NotNil(t, pm.CertificateUsage)
	require.Equal(t, "tls-server", *pm.CertificateUsage)
	require.Equal(t, metadata, pm.Metadata)
}

func TestNewPlatformManagedNilMetadata(t *testing.T) {
	pm := NewPlatformManaged("tls-server", nil)

	require.NotNil(t, pm.CertificateUsage)
	require.Equal(t, "tls-server", *pm.CertificateUsage)
	require.Nil(t, pm.Metadata)

	data, err := json.Marshal(pm)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"certificateUsage": "tls-server"
	}`, string(data))
}

func TestPlatformManagedSerde(t *testing.T) {
	policy := CertificatePolicy{
		IssuerParameters: &IssuerParameters{Name: to.Ptr("Self")},
		PlatformManaged: NewPlatformManaged("tls-server", map[string]any{
			"issuer": "internal-ca",
			"nested": map[string]any{
				"enabled": true,
			},
			"usages": []any{"server", "client"},
		}),
	}

	data, err := json.Marshal(policy)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"issuer": {"name": "Self"},
		"platformManaged": {
			"certificateUsage": "tls-server",
			"metadata": {
				"issuer": "internal-ca",
				"nested": {"enabled": true},
				"usages": ["server", "client"]
			}
		}
	}`, string(data))

	var roundTrip CertificatePolicy
	err = json.Unmarshal(data, &roundTrip)
	require.NoError(t, err)
	require.NotNil(t, roundTrip.PlatformManaged)
	require.NotNil(t, roundTrip.PlatformManaged.CertificateUsage)
	require.Equal(t, "tls-server", *roundTrip.PlatformManaged.CertificateUsage)
	require.Equal(t, "internal-ca", roundTrip.PlatformManaged.Metadata["issuer"])
	require.Equal(t, map[string]any{"enabled": true}, roundTrip.PlatformManaged.Metadata["nested"])
	require.Equal(t, []any{"server", "client"}, roundTrip.PlatformManaged.Metadata["usages"])
}

func TestCreateCertificateRequestIncludesPlatformManaged(t *testing.T) {
	client := &Client{vaultBaseUrl: "https://fakevault.vault.azure.net"}
	parameters := CreateCertificateParameters{
		CertificatePolicy: &CertificatePolicy{
			IssuerParameters: &IssuerParameters{Name: to.Ptr("Self")},
			PlatformManaged: NewPlatformManaged("tls-server", map[string]any{
				"issuer":       "internal-ca",
				"rotationDays": 90,
			}),
		},
	}

	req, err := client.createCertificateCreateRequest(context.Background(), "cert-name", parameters, nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, req.Raw().Body.Close())
	}()
	require.Equal(t, version20260301Preview, req.Raw().URL.Query().Get("api-version"))

	body, err := io.ReadAll(req.Raw().Body)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"policy": {
			"issuer": {"name": "Self"},
			"platformManaged": {
				"certificateUsage": "tls-server",
				"metadata": {
					"issuer": "internal-ca",
					"rotationDays": 90
				}
			}
		}
	}`, string(body))
}

func TestUpdateCertificatePolicyRequestIncludesPlatformManaged(t *testing.T) {
	client := &Client{vaultBaseUrl: "https://fakevault.vault.azure.net"}
	policy := CertificatePolicy{
		IssuerParameters: &IssuerParameters{Name: to.Ptr("Self")},
		PlatformManaged: NewPlatformManaged("tls-client", map[string]any{
			"issuer": "internal-ca",
			"renewal": map[string]any{
				"enabled": true,
			},
		}),
	}

	req, err := client.updateCertificatePolicyCreateRequest(context.Background(), "cert-name", policy, nil)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, req.Raw().Body.Close())
	}()
	require.Equal(t, version20260301Preview, req.Raw().URL.Query().Get("api-version"))

	body, err := io.ReadAll(req.Raw().Body)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"issuer": {"name": "Self"},
		"platformManaged": {
			"certificateUsage": "tls-client",
			"metadata": {
				"issuer": "internal-ca",
				"renewal": {
					"enabled": true
				}
			}
		}
	}`, string(body))
}
