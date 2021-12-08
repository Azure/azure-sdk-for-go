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

var pathToPackage = "sdk/keyvault/azkeys/crypto/testdata"

const fakeKvURL = "https://fakekvurl.vault.azure.net/"

func startTest(t *testing.T) func() {
	err := recording.Start(t, pathToPackage, nil)
	require.NoError(t, err)
	return func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	}
}

func createRandomName(t *testing.T, prefix string) (string, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(t.Name()))
	return prefix + fmt.Sprint(h.Sum32()), err
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

func createClient(t *testing.T, key string) (*Client, error) {
	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{
		azcore.ClientOptions{
			Transport:       client,
		},
	}

	cred := getCredential(t)

	return NewClient(key, cred, options)
}

func createKeyClient(t *testing.T) (*azkeys.Client, error) {
	vaultUrl := recording.GetEnvVariable("AZURE_KEYVAULT_URL", fakeKvURL)
	if recording.GetRecordMode() == "playback" {
		vaultUrl = fakeKvURL
	}

	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &azkeys.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport:       client,
		},
	}

	cred := getCredential(t)

	return azkeys.NewClient(vaultUrl, cred, options)
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
