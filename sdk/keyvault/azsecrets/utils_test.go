//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets

import (
	"context"
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
	fakeVaultURL  = "https://fakekvurl.vault.azure.net"
	pathToPackage = "sdk/keyvault/azsecrets/testdata"
)

var liveVaultURL string

func TestMain(m *testing.M) {
	liveVaultURL = strings.TrimSuffix(os.Getenv("AZURE_KEYVAULT_URL"), "/")
	if liveVaultURL == "" && recording.GetRecordMode() != recording.PlaybackMode {
		panic("no value for AZURE_KEYVAULT_URL")
	}
	err := recording.ResetProxy(nil)
	if err != nil {
		panic(err)
	}
	if recording.GetRecordMode() == recording.RecordingMode {
		err := recording.AddURISanitizer(fakeVaultURL, liveVaultURL, nil)
		if err != nil {
			panic(err)
		}
		err = recording.AddBodyRegexSanitizer(fakeVaultURL, liveVaultURL, nil)
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
	os.Exit(m.Run())
}

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

func createClient(t *testing.T) (*Client, error) {
	vaultURL := liveVaultURL
	if vaultURL == "" {
		vaultURL = fakeVaultURL
	}

	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{
		azcore.ClientOptions{
			Transport: client,
		},
	}

	var cred azcore.TokenCredential
	if recording.GetRecordMode() != "playback" {
		tenantId := lookupEnvVar("AZSECRETS_TENANT_ID")
		clientId := lookupEnvVar("AZSECRETS_CLIENT_ID")
		clientSecret := lookupEnvVar("AZSECRETS_CLIENT_SECRET")
		cred, err = azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
		require.NoError(t, err)
	} else {
		cred = NewFakeCredential()
	}

	return NewClient(vaultURL, cred, options)
}

func getPollingOptions() *runtime.PollUntilDoneOptions {
	freq := time.Second
	if recording.GetRecordMode() == recording.RecordingMode {
		freq = time.Minute
	}
	return &runtime.PollUntilDoneOptions{Frequency: freq}
}

func cleanUpSecret(t *testing.T, client *Client, secret string) {
	resp, err := client.BeginDeleteSecret(context.Background(), secret, nil)
	require.NoError(t, err)

	_, err = resp.PollUntilDone(context.Background(), getPollingOptions())
	require.NoError(t, err)

	_, err = client.PurgeDeletedSecret(context.Background(), secret, nil)
	require.NoError(t, err)
}

type FakeCredential struct{}

func NewFakeCredential() *FakeCredential {
	return &FakeCredential{}
}

func (f *FakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "faketoken", ExpiresOn: time.Now().Add(time.Hour).UTC()}, nil
}
