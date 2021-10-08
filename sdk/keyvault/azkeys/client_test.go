//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"context"
	"os"
	"testing"

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

	resp, err := client.CreateKey(ctx, key, JSONWebKeyTypeRSA, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Key)

	resp, err = client.CreateKey(ctx, key, JSONWebKeyTypeRSAHSM, nil)
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

	resp, err = client.CreateECKey(ctx, key, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Key)
}
