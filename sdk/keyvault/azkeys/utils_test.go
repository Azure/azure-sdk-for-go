//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"context"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var pathToPackage = "sdk/keyvault/azkeys/testdata"

const fakeKvURL = "https://fakekvurl.vault.azure.net/"
const fakeKvMHSMURL = "https://fakekvurl.managedhsm.azure.net/"

var enableHSM = true

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

		mhsmURL, ok := os.LookupEnv("AZURE_MANAGEDHSM_URL")
		if !ok {
			fmt.Println("Did not find managed HSM url, skipping those tests")
			enableHSM = false
		} else {
			err = recording.AddURISanitizer(fakeKvMHSMURL, mhsmURL, nil)
			if err != nil {
				panic(err)
			}

			err = recording.AddBodyKeySanitizer("$.key.kid", mhsmURL, fakeKvMHSMURL, nil)
			if err != nil {
				panic(err)
			}
		}
	} else if recording.GetRecordMode() == "live" {
		_, ok := os.LookupEnv("AZURE_MANAGEDHSM_URL")
		if !ok {
			fmt.Println("Did not find managed HSM url, skipping those tests")
			enableHSM = false
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

func skipHSM(t *testing.T, testType string) {
	if testType == HSMTEST && !enableHSM {
		if recording.GetRecordMode() != recording.PlaybackMode {
			t.Log("Skipping HSM Test")
			t.Skip()
		}
	}
}

func alwaysSkipHSM(t *testing.T, testType string) {
	if testType == HSMTEST {
		t.Log("Skipping HSM Test")
		t.Skip()
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

func createClient(t *testing.T, testType string) (*Client, error) {
	vaultUrl := recording.GetEnvVariable("AZURE_KEYVAULT_URL", fakeKvURL)
	// var credOptions *azidentity.ClientSecretCredentialOptions
	if testType == HSMTEST {
		vaultUrl = recording.GetEnvVariable("AZURE_MANAGEDHSM_URL", fakeKvMHSMURL)
	}

	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{
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

	return NewClient(vaultUrl, cred, options)
}

func delay() time.Duration {
	if recording.GetRecordMode() == "playback" {
		return 1 * time.Microsecond
	}
	return 250 * time.Millisecond
}

func cleanUpKey(t *testing.T, client *Client, key string) {
	resp, err := client.BeginDeleteKey(context.Background(), key, nil)
	if err != nil {
		return
	}

	_, err = resp.PollUntilDone(context.Background(), delay())
	require.NoError(t, err)

	_, err = client.PurgeDeletedKey(context.Background(), key, nil)
	require.NoError(t, err)
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

func toBytes(s string, t *testing.T) []byte {
	if len(s)%2 == 1 {
		s = fmt.Sprintf("0%s", s)
	}
	ret, err := hex.DecodeString(s)
	require.NoError(t, err)
	return ret
}
