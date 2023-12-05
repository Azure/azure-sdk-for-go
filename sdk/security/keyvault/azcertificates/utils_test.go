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
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/internal"
	"github.com/stretchr/testify/require"
)

const recordingDirectory = "sdk/security/keyvault/azcertificates/testdata"
const fakeVaultURL = "https://fakevault.local"

var (
	certsToPurge = struct {
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
		proxy, err := recording.StartTestProxy(recordingDirectory, nil)
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

	vaultURL = internal.GetEnvVar("AZURE_KEYVAULT_URL", fakeVaultURL)

	credential = internal.GetCredential("AZADMIN")

	if recording.GetRecordMode() == recording.RecordingMode {
		opts := proxy.Options
		opts.GroupForReplace = "1"
		err := recording.AddHeaderRegexSanitizer("WWW-Authenticate", "https://local", `resource="(.*)"`, opts)
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
	internal.StartRecording(t, recordingDirectory)
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

// pollCertOperation polls a certificate operation for up to 2 minutes, stopping when it completes.
// It fails the test if a poll fails or the operation doesn't complete successfully in the allotted time.
func pollCertOperation(t *testing.T, client *azcertificates.Client, name string) {
	var err error
	var op azcertificates.GetCertificateOperationResponse
	for i := 0; i < 24; i++ {
		op, err = client.GetCertificateOperation(ctx, name, nil)
		require.NoError(t, err)
		require.NotNil(t, op.Status)
		switch s := *op.Status; s {
		case "completed":
			return
		case "cancelled":
			t.Fatal("cert creation cancelled")
		case "inProgress":
			// sleep and continue
		default:
			t.Fatalf(`unexpected status "%s"`, s)
		}
		if i < 23 {
			recording.Sleep(5 * time.Second)
		} else {
			t.Fatal("cert creation didn't complete in time")
		}
	}
}
