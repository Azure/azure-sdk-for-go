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
	"time"

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

	key, err := createRandomName(t, "key")
	require.NoError(t, err)

	resp, err := client.CreateRSAKey(ctx, key, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Key)

	resp2, err := client.CreateRSAKey(ctx, key+"hsm", &CreateRSAKeyOptions{HardwareProtected: true})
	require.NoError(t, err)
	require.NotNil(t, resp2.Key)

	cleanUpKey(t, client, key)
	cleanUpKey(t, client, key+"hsm")
}

func TestCreateECKey(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	key, err := createRandomName(t, "key")
	require.NoError(t, err)

	resp, err := client.CreateECKey(ctx, key, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Key)

	cleanUpKey(t, client, key)
}

func TestCreateOCTKey(t *testing.T) {
	t.Skipf("OCT Key is HSM only")
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	key, err := createRandomName(t, "key")
	require.NoError(t, err)

	resp, err := client.CreateOCTKey(ctx, key, &CreateOCTKeyOptions{KeySize: to.Int32Ptr(256), HardwareProtected: true})
	require.NoError(t, err)
	require.NotNil(t, resp.Key)

	cleanUpKey(t, client, key)
}

func TestListKeys(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	for i := 0; i < 4; i++ {
		key, err := createRandomName(t, fmt.Sprintf("key-%d", i))
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

	for i := 0; i < 4; i++ {
		key, err := createRandomName(t, fmt.Sprintf("key-%d", i))
		require.NoError(t, err)
		cleanUpKey(t, client, key)
	}
}

func TestGetKey(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	key, err := createRandomName(t, "key")
	require.NoError(t, err)

	_, err = client.CreateKey(ctx, key, RSA, nil)
	require.NoError(t, err)

	resp, err := client.GetKey(ctx, key, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Key)
}

func TestDeleteKey(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	key, err := createRandomName(t, "key")
	require.NoError(t, err)
	defer cleanUpKey(t, client, key)

	_, err = client.CreateKey(ctx, key, RSA, nil)
	require.NoError(t, err)

	resp, err := client.BeginDeleteKey(ctx, key, nil)
	require.NoError(t, err)
	_, err = resp.PollUntilDone(ctx, 1*time.Second)
	require.Nil(t, err)

	_, err = client.GetKey(ctx, key, nil)
	require.Error(t, err)

	_, err = client.PurgeDeletedKey(ctx, key, nil)
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		_, err = client.GetDeletedKey(ctx, key, nil)
		if err != nil {
			break
		}
		require.NoError(t, err)
		recording.Sleep(time.Second * 2)
	}

	_, err = client.GetDeletedKey(ctx, key, nil)
	require.Error(t, err)
}

func TestBackupKey(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	key, err := createRandomName(t, "key")
	require.NoError(t, err)
	defer cleanUpKey(t, client, key)

	_, err = client.CreateKey(ctx, key, RSA, nil)
	require.NoError(t, err)

	resp, err := client.BackupKey(context.Background(), key, nil)
	require.NoError(t, err)
	require.Greater(t, len(resp.Value), 0)
}

func TestRecoverDeletedKey(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	key, err := createRandomName(t, "key")
	require.NoError(t, err)

	_, err = client.CreateRSAKey(context.Background(), key, nil)
	require.NoError(t, err)

	defer cleanUpKey(t, client, key)

	pollerResp, err := client.BeginDeleteKey(context.Background(), key, nil)
	require.NoError(t, err)

	_, err = pollerResp.PollUntilDone(context.Background(), delay())
	require.NoError(t, err)

	resp, err := client.BeginRecoverDeletedKey(context.Background(), key, nil)
	require.NoError(t, err)

	_, err = resp.PollUntilDone(context.Background(), delay())
	require.NoError(t, err)

	getResp, err := client.GetKey(context.Background(), key, nil)
	require.NoError(t, err)
	require.NotNil(t, getResp.Key)
}

func TestUpdateKeyProperties(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	key, err := createRandomName(t, "key")
	require.NoError(t, err)

	_, err = client.CreateRSAKey(context.Background(), key, nil)
	require.NoError(t, err)
	defer cleanUpKey(t, client, key)

	resp, err := client.UpdateKeyProperties(ctx, key, &UpdateKeyPropertiesOptions{
		Tags: map[string]*string{
			"Tag1": to.StringPtr("Val1"),
		},
		// KeyAttributes: &KeyAttributes{
		// 	Attributes: Attributes{
		// 		Expires: to.TimePtr(time.Now().AddDate(1, 0, 0)),
		// 	},
		// },
	})
	require.NoError(t, err)
	require.NotNil(t, resp.Attributes)
	require.Equal(t, *resp.Tags["Tag1"], "Val1")
	require.NotNil(t, resp.Attributes.Updated)
}
