//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets_test

import (
	"context"
	"fmt"
	"hash/fnv"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/stretchr/testify/require"
)

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
	vaultURL = strings.TrimSuffix(os.Getenv("AZURE_KEYVAULT_URL"), "/")
	if vaultURL == "" {
		if recording.GetRecordMode() != recording.PlaybackMode {
			panic("no value for AZURE_KEYVAULT_URL")
		}
		vaultURL = fakeVaultURL
	}
	err := recording.ResetProxy(nil)
	if err != nil {
		panic(err)
	}
	if recording.GetRecordMode() == recording.PlaybackMode {
		credential = &FakeCredential{}
	} else {
		tenantID := lookupEnvVar("AZSECRETS_TENANT_ID")
		clientID := lookupEnvVar("AZSECRETS_CLIENT_ID")
		secret := lookupEnvVar("AZSECRETS_CLIENT_SECRET")
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
		err = recording.AddHeaderRegexSanitizer("WWW-Authenticate", "https://local", `resource="(.*)"`, &recording.RecordingOptions{GroupForReplace: "1"})
		if err != nil {
			panic(err)
		}
		err = recording.AddBodyRegexSanitizer(fakeVaultURL, vaultURL, nil)
		if err != nil {
			panic(err)
		}
		defer func() {
			err := recording.ResetProxy(nil)
			if err != nil {
				panic(err)
			}
		}()
	}
	code := m.Run()
	if recording.GetRecordMode() != recording.PlaybackMode {
		// Purge test secrets using a client whose requests aren't recorded. This
		// will be fast because the tests which created these secrets requested their
		// deletion. Now, at the end of the run, Key Vault will have finished deleting
		// most of them...
		client := azsecrets.NewClient(vaultURL, credential, nil)
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
	os.Exit(code)
}

func startTest(t *testing.T) *azsecrets.Client {
	err := recording.Start(t, "sdk/keyvault/azsecrets/testdata", nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &azsecrets.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	return azsecrets.NewClient(vaultURL, credential, opts)
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
