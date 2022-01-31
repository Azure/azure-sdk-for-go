//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package crypto

import (
	"context"
	"fmt"
	"hash/fnv"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
	"github.com/stretchr/testify/require"
)

const fakeKvURL = "https://fakekvurl.vault.azure.net/"

var pathToPackage = "sdk/keyvault/azkeys/testdata"

func TestMain(m *testing.M) {
	// Initialize
	if recording.GetRecordMode() == "record" {
		err := recording.ResetProxy(nil)
		if err != nil {
			panic(err)
		}

		vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
		err = recording.AddURISanitizer(fakeKvURL, vaultUrl, nil)
		if err != nil {
			panic(err)
		}

		err = recording.AddBodyKeySanitizer("$.key.kid", fakeKvURL, vaultUrl, nil)
		if err != nil {
			panic(err)
		}

		err = recording.AddBodyKeySanitizer("$.recoveryId", fakeKvURL, vaultUrl, nil)
		if err != nil {
			panic(err)
		}

		tenantID := os.Getenv("AZKEYS_TENANT_ID")
		err = recording.AddHeaderRegexSanitizer("WWW-Authenticate", "00000000-0000-0000-0000-000000000000", tenantID, nil)
		if err != nil {
			panic(err)
		}

	}

	// Run tests
	exitVal := m.Run()

	// 3. Reset
	if recording.GetRecordMode() != "live" {
		err := recording.ResetProxy(nil)
		if err != nil {
			panic(err)
		}
	}

	// 4. Error out if applicable
	os.Exit(exitVal)
}

func startTest(t *testing.T) func() {
	err := recording.Start(t, pathToPackage, nil)
	require.NoError(t, err)
	return func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	}
}

func lookupEnvVar(s string) string {
	ret, ok := os.LookupEnv(s)
	if !ok {
		panic(fmt.Sprintf("Could not find env var: '%s'", s))
	}
	return ret
}

func getCredential(t *testing.T) azcore.TokenCredential {
	if recording.GetRecordMode() != "playback" {
		tenantId := lookupEnvVar("AZKEYS_TENANT_ID")
		clientId := lookupEnvVar("AZKEYS_CLIENT_ID")
		clientSecret := lookupEnvVar("AZKEYS_CLIENT_SECRET")
		cred, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
		require.NoError(t, err)
		return cred
	} else {
		return NewFakeCredential("fake", "fake")
	}
}

type FakeCredential struct {
	accountName string
	accountKey  string
}

func NewFakeCredential(accountName, accountKey string) *FakeCredential {
	return &FakeCredential{
		accountName: accountName,
		accountKey:  accountKey,
	}
}

func (f *FakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	return &azcore.AccessToken{
		Token:     "faketoken",
		ExpiresOn: time.Date(2040, time.January, 1, 1, 1, 1, 1, time.UTC),
	}, nil
}

func createCryptoClient(t *testing.T, key string) (*Client, error) {
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: transport,
		},
	}

	cred := getCredential(t)

	return NewClient(key, cred, options)
}

func createRandomName(t *testing.T, prefix string) (string, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(t.Name()))
	return prefix + fmt.Sprint(h.Sum32()), err
}

func createClient(t *testing.T) (*azkeys.Client, error) {
	vaultUrl := recording.GetEnvVariable("AZURE_KEYVAULT_URL", fakeKvURL)

	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &azkeys.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: transport,
		},
	}

	var cred azcore.TokenCredential
	if recording.GetRecordMode() != "playback" {
		tenantId := lookupEnvVar("AZKEYS_TENANT_ID")
		clientId := lookupEnvVar("AZKEYS_CLIENT_ID")
		clientSecret := lookupEnvVar("AZKEYS_CLIENT_SECRET")
		cred, err = azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
		require.NoError(t, err)
	} else {
		cred = NewFakeCredential("fake", "fake")
	}

	return azkeys.NewClient(vaultUrl, cred, options)
}
