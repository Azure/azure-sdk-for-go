//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates_test

import (
	"context"
	"fmt"
	"hash/fnv"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/stretchr/testify/require"
)

const recordingDirectory = "sdk/security/keyvault/azcertificates/testdata"

var (
	certsToPurge = struct {
		mut   sync.Mutex
		names []string
	}{sync.Mutex{}, []string{}}

	credential azcore.TokenCredential
	vaultURL   string

	fakeVaultURL = fmt.Sprintf("https://%s.vault.azure.net/", recording.SanitizedValue)
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

	vaultURL = recording.GetEnvVariable("AZURE_KEYVAULT_URL", fakeVaultURL)

	if recording.GetRecordMode() != recording.LiveMode {
		err := recording.RemoveRegisteredSanitizers([]string{
			"AZSDK3430", // id in body
			"AZSDK3493", // name in body
		}, nil)
		if err != nil {
			panic(err)
		}
	}

	code := m.Run()
	if recording.GetRecordMode() != recording.PlaybackMode {
		// Purge test certs using a client whose requests aren't recorded. This
		// will be fast because the tests which created these certs requested their
		// deletion. Now, at the end of the run, Key Vault will have finished deleting
		// most of them...
		client, err := azcertificates.NewClient(vaultURL, credential, nil)
		if err != nil {
			panic(err)
		}
		for _, name := range certsToPurge.names {
			// ...but we need a retry loop for the others. Note this wouldn't benefit
			// from client-side parallelization because Key Vault's delete operations
			// are running in parallel. When the client waits on one deletion, it
			// effectively waits on all of them.
			for i := 0; i < 12; i++ {
				_, err := client.PurgeDeletedCertificate(context.Background(), name, nil)
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

func startTest(t *testing.T) *azcertificates.Client {
	err := recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &azcertificates.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	client, err := azcertificates.NewClient(vaultURL, credential, opts)
	require.NoError(t, err)
	return client
}

func getName(t *testing.T, prefix string) string {
	h := fnv.New32a()
	_, err := h.Write([]byte(t.Name()))
	require.NoError(t, err)
	return prefix + fmt.Sprint(h.Sum32())
}

func cleanUpCert(t *testing.T, client *azcertificates.Client, name string) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		return
	}
	if _, err := client.DeleteCertificate(context.Background(), name, nil); err == nil {
		certsToPurge.mut.Lock()
		defer certsToPurge.mut.Unlock()
		certsToPurge.names = append(certsToPurge.names, name)
	} else {
		t.Logf(`cleanUpCert failed for "%s": %v`, name, err)
	}
}
