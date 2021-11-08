//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcrypto

import (
	"context"
	"crypto/sha256"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func TestMain(m *testing.M) {
	// Initialize
	if recording.GetRecordMode() == "record" {
		vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
		err := recording.AddURISanitizer(fakeKvURL, vaultUrl, nil)
		if err != nil {
			panic(err)
		}
	}

	// Run
	exitVal := m.Run()

	// cleanup

	os.Exit(exitVal)
}

func TestNewClient(t *testing.T) {
	for _, key := range []string{"https://mykeyvault.vault.azure.net/keys/keyabcdef/1234567890", "https://mykeyvault.vault.azure.net/keys/keyabcdef"} {
		cred := getCredential(t)

		client, err := NewClient(key, cred, nil)
		require.NoError(t, err)

		require.Equal(t, client.vaultURL, "https://mykeyvault.vault.azure.net/")
		require.Equal(t, "keyabcdef", client.keyID)
		if strings.Contains(key, "1234567890") {
			require.Equal(t, client.keyVersion, "1234567890")
		} else {
			require.Equal(t, client.keyVersion, "")
		}
	}
}

func TestClient_Decrypt(t *testing.T) {
	stop := startTest(t)
	defer stop()

	keyName, err := createRandomName(t, "key")
	require.NoError(t, err)

	keyClient, err := createKeyClient(t)
	require.NoError(t, err)
	resp, err := keyClient.CreateRSAKey(ctx, keyName, nil)
	require.NoError(t, err)
	key := resp.Key

	cryptoClient, err := createClient(t, *key.ID)
	require.NoError(t, err)

	encryptResponse, err := cryptoClient.Encrypt(ctx, AlgorithmRSAOAEP, []byte("plaintext"), nil)
	require.NoError(t, err)
	require.NotNil(t, encryptResponse)

	decryptResponse, err := cryptoClient.Decrypt(ctx, AlgorithmRSAOAEP, encryptResponse.Result, nil)
	require.NoError(t, err)
	require.Equal(t, decryptResponse.Result, []byte("plaintext"))
}

func TestClient_WrapUnwrap(t *testing.T) {
	stop := startTest(t)
	defer stop()

	keyName, err := createRandomName(t, "key")
	require.NoError(t, err)

	keyClient, err := createKeyClient(t)
	require.NoError(t, err)
	resp, err := keyClient.CreateRSAKey(ctx, keyName, nil)
	require.NoError(t, err)
	key := resp.Key

	cryptoClient, err := createClient(t, *key.ID)
	require.NoError(t, err)

	keyBytes := []byte("5063e6aaa845f150200547944fd199679c98ed6f99da0a0b2dafeaf1f4684496fd532c1c229968cb9dee44957fcef7ccef59ceda0b362e56bcd78fd3faee5781c623c0bb22b35beabde0664fd30e0e824aba3dd1b0afffc4a3d955ede20cf6a854d52cfd")

	// Wrap
	wrapResp, err := cryptoClient.WrapKey(ctx, RSAOAEP, keyBytes, nil)
	require.NoError(t, err)

	// Unwrap
	unwrapResp, err := cryptoClient.UnwrapKey(ctx, RSAOAEP, wrapResp.Result, nil)
	require.NoError(t, err)
	require.Equal(t, keyBytes, unwrapResp.Result)

}

func TestClient_SignVerify(t *testing.T) {
	stop := startTest(t)
	defer stop()

	keyName, err := createRandomName(t, "key")
	require.NoError(t, err)

	keyClient, err := createKeyClient(t)
	require.NoError(t, err)
	resp, err := keyClient.CreateRSAKey(ctx, keyName, nil)
	require.NoError(t, err)
	key := resp.Key

	cryptoClient, err := createClient(t, *key.ID)
	require.NoError(t, err)

	hasher := sha256.New()
	_, err = hasher.Write([]byte("plaintext"))
	require.NoError(t, err)
	digest := hasher.Sum(nil)

	signResponse, err := cryptoClient.Sign(ctx, RS256, digest, nil)
	require.NoError(t, err)

	verifyResponse, err := cryptoClient.Verify(ctx, RS256, digest, signResponse.Result, nil)
	require.NoError(t, err)
	require.True(t, *verifyResponse.IsValid)
}
