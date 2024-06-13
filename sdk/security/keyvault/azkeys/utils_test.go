//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys_test

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/internal"
	"github.com/stretchr/testify/require"
)

const recordingDirectory = "sdk/security/keyvault/azkeys/testdata"

var (
	fakeVaultURL       = fmt.Sprintf("https://%s.vault.azure.net/", recording.SanitizedValue)
	fakeHsmURL         = fmt.Sprintf("https://%s.managedhsm.azure.net/", recording.SanitizedValue)
	fakeAttestationURL = fmt.Sprintf("https://%s.azurewebsites.net/", recording.SanitizedValue)

	keysToPurge = struct {
		mut   sync.Mutex
		names map[string][]string // maps vault URL to key names
	}{sync.Mutex{}, map[string][]string{}}

	credential     azcore.TokenCredential
	enableHSM      bool
	attestationURL string
	mhsmURL        string
	vaultURL       string
)

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	var proxy *recording.TestProxyInstance
	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		var err error
		proxy, err = recording.StartTestProxy(recordingDirectory, nil)
		if err != nil {
			panic(err)
		}

		defer func() {
			err := recording.StopTestProxy(proxy)
			if err != nil {
				panic(err)
			}
		}()
	}

	var err error
	credential, err = azcred.New(nil)
	if err != nil {
		panic(err)
	}

	attestationURL = recording.GetEnvVariable("AZURE_KEYVAULT_ATTESTATION_URL", fakeAttestationURL)
	mhsmURL = recording.GetEnvVariable("AZURE_MANAGEDHSM_URL", fakeHsmURL)
	vaultURL = recording.GetEnvVariable("AZURE_KEYVAULT_URL", fakeVaultURL)
	enableHSM = mhsmURL != fakeHsmURL

	if recording.GetRecordMode() == recording.RecordingMode {
		for _, path := range []string{"$.error.message", "$.key.kid", "$.recoveryId"} {
			err := recording.AddBodyKeySanitizer(path, fakeVaultURL, vaultURL, nil)
			if err != nil {
				panic(err)
			}
			err = recording.AddBodyKeySanitizer(path, fakeHsmURL, mhsmURL, nil)
			if err != nil {
				panic(err)
			}
		}
		// these values aren't secret but we redact them anyway to avoid
		// alerts from automation scanning for JWTs or "token" values
		for _, attestation := range []string{"$.target", "$.token"} {
			err := recording.AddBodyKeySanitizer(attestation, recording.SanitizedValue, "", nil)
			if err != nil {
				panic(err)
			}
		}
		// we need to replace release policy data because it has the attestation service URL encoded
		// into it and therefore won't match in playback, when we don't have the URL used while recording
		fakePolicyData := base64.RawStdEncoding.EncodeToString(getMarshalledReleasePolicy(fakeAttestationURL))
		err := recording.AddBodyKeySanitizer("$.release_policy.data", fakePolicyData, "", nil)
		if err != nil {
			panic(err)
		}
		err = recording.AddRemoveHeaderSanitizer([]string{"Set-Cookie"}, nil)
		if err != nil {
			panic(err)
		}
	}
	code := m.Run()
	if recording.GetRecordMode() != recording.PlaybackMode {
		// Purge test keys using a client whose requests aren't recorded. This
		// will be fast because the tests which created these keys requested their
		// deletion. Now, at the end of the run, Key Vault will have finished deleting
		// most of them...
		for URL, names := range keysToPurge.names {
			client, err := azkeys.NewClient(URL, credential, nil)
			if err != nil {
				panic(err)
			}
			for _, name := range names {
				// ...but we need a retry loop for the others. Note this wouldn't benefit
				// from client-side parallelization because Key Vault's delete operations
				// are running in parallel. When the client waits on one deletion, it
				// effectively waits on all of them.
				for i := 0; i < 12; i++ {
					_, err := client.PurgeDeletedKey(context.Background(), name, nil)
					if err == nil {
						break
					}
					if i < 11 {
						recording.Sleep(10 * time.Second)
					}
				}
			}
		}
	}

	return code
}

func startTest(t *testing.T, mhsmTest bool) *azkeys.Client {
	if recording.GetRecordMode() != recording.PlaybackMode && mhsmTest && !enableHSM {
		t.Skip("set AZURE_MANAGEDHSM_URL to run this test")
	}
	err := recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	URL := vaultURL
	if mhsmTest {
		URL = mhsmURL
	}
	opts := &azkeys.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	client, err := azkeys.NewClient(URL, credential, opts)
	require.NoError(t, err)
	return client
}

func createRandomName(t *testing.T, prefix string) string {
	h := fnv.New32a()
	_, err := h.Write([]byte(t.Name()))
	require.NoError(t, err)
	return prefix + fmt.Sprint(h.Sum32())
}

func cleanUpKey(t *testing.T, client *azkeys.Client, id *azkeys.ID) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		return
	}
	URL, name, _ := internal.ParseID((*string)(id))
	if _, err := client.DeleteKey(context.Background(), *name, nil); err == nil {
		keysToPurge.mut.Lock()
		defer keysToPurge.mut.Unlock()
		keysToPurge.names[*URL] = append(keysToPurge.names[*URL], *name)
	} else {
		t.Logf(`cleanUpKey failed for "%s": %v`, *name, err)
	}
}

func toBytes(s string, t *testing.T) []byte {
	if len(s)%2 == 1 {
		s = fmt.Sprintf("0%s", s)
	}
	ret, err := hex.DecodeString(s)
	require.NoError(t, err)
	return ret
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

type serdeModel interface {
	json.Marshaler
	json.Unmarshaler
}

func testSerde[T serdeModel](t *testing.T, model T) {
	data, err := model.MarshalJSON()
	require.NoError(t, err)
	err = model.UnmarshalJSON(data)
	require.NoError(t, err)

	// testing unmarshal error scenarios
	err = model.UnmarshalJSON(nil)
	require.Error(t, err)

	m := regexp.MustCompile(":.*$")
	modifiedData := m.ReplaceAllString(string(data), `:["test", "test1", "test2"]}`)
	if modifiedData != "{}" {
		data3 := []byte(modifiedData)
		err = model.UnmarshalJSON(data3)
		require.Error(t, err)
	}
}
