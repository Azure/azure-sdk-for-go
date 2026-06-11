// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys/fake"
	"github.com/stretchr/testify/require"
)

// These tests verify the client wires requests and responses correctly using the generated fake server.
func TestSecureWrapKey_Fake(t *testing.T) {
	const keyName, keyVersion = "key", "version"
	wrappedKey := []byte("wrapped-aes-key")

	server := fake.Server{
		SecureWrapKey: func(ctx context.Context, name string, version string, parameters azkeys.SecureKeyWrapOperationParameters, options *azkeys.SecureWrapKeyOptions) (resp azfake.Responder[azkeys.SecureWrapKeyResponse], errResp azfake.ErrorResponder) {
			require.NotNil(t, parameters.Algorithm)
			require.Equal(t, azkeys.JSONWebKeyWrapAlgorithmRSAOAEP256, *parameters.Algorithm)
			resp.SetResponse(http.StatusOK, azkeys.SecureWrapKeyResponse{
				SecureKeyOperationResult: azkeys.SecureKeyOperationResult{
					Algorithm: parameters.Algorithm,
					Kid:       to.Ptr("https://fake.vault.azure.net/keys/" + keyName + "/" + keyVersion),
					Value:     wrappedKey,
				},
			}, nil)
			return
		},
	}

	client, err := azkeys.NewClient(vaultURL, &azcred.Fake{}, &azkeys.ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: fake.NewServerTransport(&server)},
	})
	require.NoError(t, err)

	params := azkeys.SecureKeyWrapOperationParameters{
		Algorithm: to.Ptr(azkeys.JSONWebKeyWrapAlgorithmRSAOAEP256),
	}
	resp, err := client.SecureWrapKey(context.Background(), keyName, keyVersion, params, nil)
	require.NoError(t, err)
	require.Equal(t, wrappedKey, resp.Value)
	require.Equal(t, azkeys.JSONWebKeyWrapAlgorithmRSAOAEP256, *resp.Algorithm)
}

func TestSecureUnwrapKey_Fake(t *testing.T) {
	const keyName, keyVersion = "key", "version"
	unwrappedKey := []byte("unwrapped-aes-key")
	wrappedKey := []byte("wrapped-aes-key")

	server := fake.Server{
		SecureUnwrapKey: func(ctx context.Context, name string, version string, parameters azkeys.SecureKeyUnWrapOperationParameters, options *azkeys.SecureUnwrapKeyOptions) (resp azfake.Responder[azkeys.SecureUnwrapKeyResponse], errResp azfake.ErrorResponder) {
			require.NotNil(t, parameters.Algorithm)
			require.Equal(t, azkeys.JSONWebKeyWrapAlgorithmRSAOAEP256, *parameters.Algorithm)
			require.Equal(t, wrappedKey, parameters.Value)
			require.Equal(t, "attestation-token", *parameters.TargetAttestationToken)
			resp.SetResponse(http.StatusOK, azkeys.SecureUnwrapKeyResponse{
				SecureKeyOperationResult: azkeys.SecureKeyOperationResult{
					Algorithm: parameters.Algorithm,
					Kid:       to.Ptr("https://fake.vault.azure.net/keys/" + keyName + "/" + keyVersion),
					Value:     unwrappedKey,
				},
			}, nil)
			return
		},
	}

	client, err := azkeys.NewClient(vaultURL, &azcred.Fake{}, &azkeys.ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: fake.NewServerTransport(&server)},
	})
	require.NoError(t, err)

	params := azkeys.SecureKeyUnWrapOperationParameters{
		Algorithm:              to.Ptr(azkeys.JSONWebKeyWrapAlgorithmRSAOAEP256),
		TargetAttestationToken: to.Ptr("attestation-token"),
		Value:                  wrappedKey,
	}
	resp, err := client.SecureUnwrapKey(context.Background(), keyName, keyVersion, params, nil)
	require.NoError(t, err)
	require.Equal(t, unwrappedKey, resp.Value)
	require.Equal(t, azkeys.JSONWebKeyWrapAlgorithmRSAOAEP256, *resp.Algorithm)
}
