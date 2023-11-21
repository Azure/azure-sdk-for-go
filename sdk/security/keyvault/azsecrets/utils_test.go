//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/stretchr/testify/require"
)

const recordingDirectory = "sdk/security/keyvault/azsecrets/testdata"
const fakeVaultURL = "https://fakevault.local"

var (
	secretsToPurge = struct {
		mut   sync.Mutex
		names []string
	}{sync.Mutex{}, []string{}}

	credential azcore.TokenCredential
	vaultURL   string
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

	vaultURL = fakeVaultURL
	if recording.GetRecordMode() != recording.PlaybackMode {
		if u, ok := os.LookupEnv("AZURE_KEYVAULT_URL"); ok && u != "" {
			vaultURL = strings.TrimSuffix(u, "/")
		} else {
			panic("no value for AZURE_KEYVAULT_URL")
		}
	}
	if recording.GetRecordMode() == recording.PlaybackMode {
		credential = &FakeCredential{}
	} else {
		tenantID := lookupEnvVar("AZSECRETS_TENANT_ID")
		clientID := lookupEnvVar("AZSECRETS_CLIENT_ID")
		secret := lookupEnvVar("AZSECRETS_CLIENT_SECRET")
		var err error
		credential, err = azidentity.NewClientSecretCredential(tenantID, clientID, secret, nil)
		if err != nil {
			panic(err)
		}
	}
	if recording.GetRecordMode() == recording.RecordingMode {
		err := recording.AddURISanitizer(fakeVaultURL, vaultURL, nil)
		if err != nil {
			panic(err)
		}
		opts := proxy.Options
		opts.GroupForReplace = "1"
		err = recording.AddHeaderRegexSanitizer("WWW-Authenticate", "https://local", `resource="(.*)"`, opts)
		if err != nil {
			panic(err)
		}
		err = recording.AddBodyRegexSanitizer(fakeVaultURL, vaultURL, nil)
		if err != nil {
			panic(err)
		}
	}
	code := m.Run()
	if recording.GetRecordMode() != recording.PlaybackMode {
		// Purge test secrets using a client whose requests aren't recorded. This
		// will be fast because the tests which created these secrets requested their
		// deletion. Now, at the end of the run, Key Vault will have finished deleting
		// most of them...
		client, err := azsecrets.NewClient(vaultURL, credential, nil)
		if err != nil {
			panic(err)
		}
		for _, name := range secretsToPurge.names {
			// ...but we need a retry loop for the others. Note this wouldn't benefit
			// from client-side parallelization because Key Vault's delete operations
			// are running in parallel. When the client waits on one deletion, it
			// effectively waits on all of them.
			for i := 0; i < 12; i++ {
				_, err := client.PurgeDeletedSecret(context.Background(), name, nil)
				if err == nil {
					break
				}
				if i < 11 {
					recording.Sleep(10 * time.Second)
				}
			}
		}
	}

	return code
}

func startTest(t *testing.T) *azsecrets.Client {
	err := recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &azsecrets.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	client, err := azsecrets.NewClient(vaultURL, credential, opts)
	require.NoError(t, err)
	return client
}

func createRandomName(t *testing.T, prefix string) string {
	h := fnv.New32a()
	_, err := h.Write([]byte(t.Name()))
	require.NoError(t, err)
	return prefix + fmt.Sprint(h.Sum32())
}

func lookupEnvVar(s string) string {
	ret, ok := os.LookupEnv(s)
	if !ok {
		panic(fmt.Sprintf("Could not find env var: '%s'", s))
	}
	return ret
}

func cleanUpSecret(t *testing.T, client *azsecrets.Client, name string) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		return
	}
	if _, err := client.DeleteSecret(context.Background(), name, nil); err == nil {
		secretsToPurge.mut.Lock()
		defer secretsToPurge.mut.Unlock()
		secretsToPurge.names = append(secretsToPurge.names, name)
	} else {
		t.Logf(`cleanUpSecret failed for "%s": %v`, name, err)
	}
}

type FakeCredential struct{}

func (f *FakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "faketoken", ExpiresOn: time.Now().Add(time.Hour).UTC()}, nil
}

// pollStatus calls a function until it stops returning a response error with the given status code.
// If this takes more than 2 minutes, it fails the test.
func pollStatus(t *testing.T, expectedStatus int, fn func() error) {
	var err error
	for i := 0; i < 12; i++ {
		err = fn()
		var respErr *azcore.ResponseError
		if !(errors.As(err, &respErr) && respErr.StatusCode == expectedStatus) {
			break
		}
		if i < 11 {
			recording.Sleep(10 * time.Second)
		}
	}
	require.NoError(t, err)
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
	var data2 []byte
	err = model.UnmarshalJSON(data2)
	require.Error(t, err)

	m := regexp.MustCompile(":.*$")
	modifiedData := m.ReplaceAllString(string(data), ":false}")
	if modifiedData != "{}" {
		data3 := []byte(modifiedData)
		err = model.UnmarshalJSON(data3)
		require.Error(t, err)
	}
}
