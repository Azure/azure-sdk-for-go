//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

func TestMain(m *testing.M) {
	// Initialize
	if recording.GetRecordMode() == "record" {
		vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
		err := recording.AddUriSanitizer("https://fakekvurl.vault.azure.net/", vaultUrl, nil)
		if err != nil {
			panic(err)
		}
	}

	// Run
	exitVal := m.Run()

	// cleanup

	os.Exit(exitVal)
}

func startTest(t *testing.T) func() {
	err := recording.StartRecording(t, pathToPackage, nil)
	require.NoError(t, err)
	return func() {
		err := recording.StopRecording(t, nil)
		require.NoError(t, err)
	}
}

func TestCreateKeyRSA(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	key, err := createRandomName(t, "secret")
	require.NoError(t, err)

	resp, err := client.CreateKey(ctx, key, RSA, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Key)

	resp, err = client.CreateKey(ctx, key, RSAHSM, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Key)
}

func TestCreateECKey(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	key, err := createRandomName(t, "secret")
	require.NoError(t, err)

	resp, err := client.CreateECKey(ctx, key, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Key)
}

func TestCreateOCTKey(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	key, err := createRandomName(t, "secret")
	require.NoError(t, err)

	// _, err = client.CreateKey(ctx, key, JSONWebKeyTypeOct, nil)
	// require.NoError(t, err)

	resp, err := client.CreateOCTKey(ctx, key, &CreateOCTKeyOptions{KeySize: to.Int32Ptr(256)})
	require.NoError(t, err)
	require.NotNil(t, resp.Key)
}

func TestListKeys(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	for i := 0; i < 4; i++ {
		key, err := createRandomName(t, fmt.Sprintf("secret-%d", i))
		require.NoError(t, err)

		_, err = client.CreateKey(ctx, key, RSA, nil)
		require.NoError(t, err)
	}

	pager := client.ListKeys(nil)
	count := 0
	for pager.NextPage(ctx) {
		count += len(pager.PageResponse().Keys)
		for _, key := range pager.PageResponse().Keys {
			require.NotNil(t, key)
		}
	}

	require.NoError(t, pager.Err())
	require.GreaterOrEqual(t, count, 4)
}

func TestGetKey(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	key, err := createRandomName(t, "secret")
	require.NoError(t, err)

	_, err = client.CreateKey(ctx, key, RSA, nil)
	require.NoError(t, err)

	resp, err := client.GetKey(ctx, key, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Key)
}
