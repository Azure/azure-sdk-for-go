//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

const (
	fakeAttestationUrl = "https://fakeattestation"
	fakeVaultURL       = "https://fakekvurl.vault.azure.net"
	fakeMHSMURL        = "https://fakekvurl.managedhsm.azure.net"
)

var (
	enableHSM          bool
	testAttestationURL string
	testMHSMURL        string
	testVaultURL       string
	pathToPackage      = "sdk/keyvault/azkeys/testdata"
)

func TestMain(m *testing.M) {
	testAttestationURL = strings.TrimSuffix(recording.GetEnvVariable("AZURE_KEYVAULT_ATTESTATION_URL", fakeAttestationUrl), "/")
	testVaultURL = strings.TrimSuffix(recording.GetEnvVariable("AZURE_KEYVAULT_URL", fakeVaultURL), "/")
	testMHSMURL = strings.TrimSuffix(recording.GetEnvVariable("AZURE_MANAGEDHSM_URL", fakeMHSMURL), "/")
	enableHSM = testMHSMURL != fakeMHSMURL

	if testVaultURL == fakeVaultURL && recording.GetRecordMode() != recording.PlaybackMode {
		panic("no value for AZURE_KEYVAULT_URL")
	}
	err := recording.ResetProxy(nil)
	if err != nil {
		panic(err)
	}
	switch recording.GetRecordMode() {
	case recording.PlaybackMode:
		err := recording.SetDefaultMatcher(nil, &recording.SetDefaultMatcherOptions{
			ExcludedHeaders: []string{":path", ":authority", ":method", ":scheme"},
		})
		if err != nil {
			panic(err)
		}
	case recording.RecordingMode:
		defer func() {
			err := recording.ResetProxy(nil)
			if err != nil {
				panic(err)
			}
		}()
		err := recording.AddURISanitizer(fakeVaultURL, testVaultURL, nil)
		if err != nil {
			panic(err)
		}
		err = recording.AddURISanitizer(fakeAttestationUrl, testAttestationURL, nil)
		if err != nil {
			panic(err)
		}
		err = recording.AddRemoveHeaderSanitizer([]string{"Set-Cookie"}, nil)
		if err != nil {
			panic(err)
		}

		keyIDPaths := []string{"$.error.message", "$.key.kid", "$.recoveryId"}
		for _, path := range keyIDPaths {
			err = recording.AddBodyKeySanitizer(path, fakeVaultURL, testVaultURL, nil)
			if err != nil {
				panic(err)
			}
		}
		// these values aren't secret but we redact them anyway to avoid
		// alerts from automation scanning for JWTs or "token" values
		for _, attestation := range []string{"$.target", "$.token"} {
			err = recording.AddBodyKeySanitizer(attestation, "redacted", "", nil)
			if err != nil {
				panic(err)
			}
		}
		// we need to replace release policy data because it has the attestation service URL encoded
		// into it and therefore won't match in playback, when we don't have the URL used while recording
		realPolicyData := base64.StdEncoding.EncodeToString(getMarshalledReleasePolicy(testAttestationURL))
		fakePolicyData := base64.RawStdEncoding.EncodeToString(getMarshalledReleasePolicy(fakeAttestationUrl))
		err = recording.AddBodyKeySanitizer("$.release_policy.data", fakePolicyData, realPolicyData, nil)
		if err != nil {
			panic(err)
		}
		if enableHSM {
			err = recording.AddURISanitizer(fakeMHSMURL, testMHSMURL, nil)
			if err != nil {
				panic(err)
			}
			for _, path := range keyIDPaths {
				err = recording.AddBodyKeySanitizer(path, fakeMHSMURL, testMHSMURL, nil)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	os.Exit(m.Run())
}

func startTest(t *testing.T, testType string) {
	if recording.GetRecordMode() != recording.PlaybackMode && testType == HSMTEST && !enableHSM {
		t.Skip("set AZURE_MANAGEDHSM_URL to run this test")
	}
	err := recording.Start(t, pathToPackage, nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
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
	vaultUrl := testVaultURL
	if testType == HSMTEST {
		vaultUrl = testMHSMURL
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

func delay() *runtime.PollUntilDoneOptions {
	return &runtime.PollUntilDoneOptions{Frequency: time.Second}
}

func cleanUpKey(t *testing.T, client *Client, key string) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		return
	}

	resp, err := client.BeginDeleteKey(context.Background(), key, nil)
	if err != nil {
		return
	}

	_, err = resp.PollUntilDone(context.Background(), delay())
	require.NoError(t, err)

	_, err = client.PurgeDeletedKey(context.Background(), key, nil)
	require.NoError(t, err)
}

type FakeCredential struct{}

func NewFakeCredential(accountName, accountKey string) *FakeCredential {
	return &FakeCredential{}
}

func (f *FakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token:     "faketoken",
		ExpiresOn: time.Now().UTC().Add(time.Hour),
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

func validateKey(t *testing.T, key *Key) {
	require.NotNil(t, key)
	require.NotNil(t, key.Properties)
	validateProperties(t, key.Properties)
	require.NotNil(t, key.JSONWebKey)
	require.NotNil(t, key.ID)
	require.NotNil(t, key.Name)
}

func validateProperties(t *testing.T, props *Properties) {
	require.NotNil(t, props)
	if props.CreatedOn == nil {
		t.Fatalf("expected CreatedOn to be not nil")
	}
	if props.Enabled == nil {
		t.Fatalf("expected Enabled to be not nil")
	}
	if props.ID == nil {
		t.Fatalf("expected ID to be not nil")
	}
	if props.Name == nil {
		t.Fatalf("expected Name to be not nil")
	}
	if props.RecoverableDays == nil {
		t.Fatalf("expected RecoverableDays to be not nil")
	}
	if props.RecoveryLevel == nil {
		t.Fatalf("expected RecoveryLevel to be not nil")
	}
	if props.UpdatedOn == nil {
		t.Fatalf("expected UpdatedOn to be not nil")
	}
	if props.VaultURL == nil {
		t.Fatalf("expected VaultURL to be not nil")
	}
	if props.Version == nil {
		t.Fatalf("expected Version to be not nil")
	}
}

func getMarshalledReleasePolicy(attestationURL string) []byte {
	data, _ := json.Marshal(map[string]interface{}{
		"anyOf": []map[string]interface{}{
			{
				"anyOf": []map[string]interface{}{
					{
						"claim":  "sdk-test",
						"equals": "true",
					}},
				"authority": attestationURL,
			},
		},
		"version": "1.0.0",
	})
	return data
}
