//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azsecrets

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Initialize
	if recording.GetRecordMode() == "record" {
		vaultUrl := os.Getenv("AZURE_KEYVAULT_URL")
		err := recording.AddURISanitizer("https://fakekvurl.vault.azure.net/", vaultUrl, nil)
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
	err := recording.Start(t, pathToPackage, nil)
	require.NoError(t, err)
	return func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	}
}

func TestSetGetSecret(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	secret, err := createRandomName(t, "secret")
	require.NoError(t, err)
	value, err := createRandomName(t, "value")
	require.NoError(t, err)

	defer cleanUpSecret(t, client, secret)

	_, err = client.SetSecret(context.Background(), secret, value, nil)
	require.NoError(t, err)

	getResp, err := client.GetSecret(context.Background(), secret, nil)
	require.NoError(t, err)
	require.Equal(t, *getResp.Value, value)
}

func TestListSecretVersions(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	secret, err := createRandomName(t, "secret")
	require.NoError(t, err)
	value, err := createRandomName(t, "value")
	require.NoError(t, err)

	_, err = client.SetSecret(context.Background(), secret, value, nil)
	require.NoError(t, err)
	_, err = client.SetSecret(context.Background(), secret, value+"1", nil)
	require.NoError(t, err)
	_, err = client.SetSecret(context.Background(), secret, value+"2", nil)
	require.NoError(t, err)
	defer cleanUpSecret(t, client, secret)

	count := 0
	pager := client.ListSecretVersions(secret, nil)
	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()
		count += len(page.Secrets)
	}
	require.GreaterOrEqual(t, count, 3)
	require.NoError(t, pager.Err())
}

func TestListSecrets(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	_, err = client.SetSecret(context.Background(), "secret1", "value", nil)
	require.NoError(t, err)
	_, err = client.SetSecret(context.Background(), "secret2", "value", nil)
	require.NoError(t, err)
	_, err = client.SetSecret(context.Background(), "secret3", "value", nil)
	require.NoError(t, err)
	_, err = client.SetSecret(context.Background(), "secret4", "value", nil)
	require.NoError(t, err)

	defer cleanUpSecret(t, client, "secret1")
	defer cleanUpSecret(t, client, "secret2")
	defer cleanUpSecret(t, client, "secret3")
	defer cleanUpSecret(t, client, "secret4")

	count := 0
	pager := client.ListSecrets(nil)
	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()
		count += len(page.Secrets)
	}
	require.Equal(t, count, 4)
	require.NoError(t, pager.Err())
}

func TestListDeletedSecrets(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	secret1, err := createRandomName(t, "secret1")
	require.NoError(t, err)
	value1, err := createRandomName(t, "value1")
	require.NoError(t, err)
	secret2, err := createRandomName(t, "secret2")
	require.NoError(t, err)
	value2, err := createRandomName(t, "value2")
	require.NoError(t, err)

	// 1. Create 2 secrets
	_, err = client.SetSecret(context.Background(), secret1, value1, nil)
	require.NoError(t, err)

	_, err = client.SetSecret(context.Background(), secret2, value2, nil)
	require.NoError(t, err)

	// 2. Delete both secrets
	resp, err := client.BeginDeleteSecret(context.Background(), secret1, nil)
	require.NoError(t, err)
	_, err = resp.PollUntilDone(context.Background(), delay())
	require.NoError(t, err)

	resp, err = client.BeginDeleteSecret(context.Background(), secret2, nil)
	require.NoError(t, err)
	_, err = resp.PollUntilDone(context.Background(), delay())
	require.NoError(t, err)

	f := func() {
		_, err := client.PurgeDeletedSecret(context.Background(), secret1, nil)
		require.NoError(t, err)
		_, err = client.PurgeDeletedSecret(context.Background(), secret2, nil)
		require.NoError(t, err)
	}
	defer f()

	// Make sure both secrets show up in deleted secrets
	deletedSecrets := map[string]bool{
		secret1: false,
		secret2: false,
	}
	count := 0
	pager := client.ListDeletedSecrets(nil)
	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()
		count += len(page.DeletedSecrets)
		for _, secret := range page.DeletedSecrets {
			for deleted := range deletedSecrets {
				if strings.Contains(*secret.ID, deleted) {
					deletedSecrets[deleted] = true
					break
				}
			}
		}
	}

	for _, deleted := range deletedSecrets {
		require.True(t, deleted)
	}
}

func TestDeleteSecret(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	secret, err := createRandomName(t, "secret1")
	require.NoError(t, err)
	value, err := createRandomName(t, "value1")
	require.NoError(t, err)

	_, err = client.SetSecret(context.Background(), secret, value, nil)
	require.NoError(t, err)

	resp, err := client.BeginDeleteSecret(context.Background(), secret, nil)
	require.NoError(t, err)

	_, err = resp.PollUntilDone(context.Background(), delay())
	require.NoError(t, err)

	_, err = client.GetDeletedSecret(context.Background(), secret, nil)
	require.NoError(t, err)

	_, err = client.PurgeDeletedSecret(context.Background(), secret, nil)
	require.NoError(t, err)

	_, err = client.GetSecret(context.Background(), secret, nil)
	require.Error(t, err)
}

func TestPurgeDeletedSecret(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	secret, err := createRandomName(t, "secret1")
	require.NoError(t, err)
	value, err := createRandomName(t, "value1")
	require.NoError(t, err)

	_, err = client.SetSecret(context.Background(), secret, value, nil)
	require.NoError(t, err)

	resp, err := client.BeginDeleteSecret(context.Background(), secret, nil)
	require.NoError(t, err)

	_, err = resp.PollUntilDone(context.Background(), delay())
	require.NoError(t, err)

	_, err = client.PurgeDeletedSecret(context.Background(), secret, nil)
	require.NoError(t, err)

	pager := client.ListDeletedSecrets(nil)
	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()
		for _, secret := range page.DeletedSecrets {
			require.NotEqual(t, *secret.ID, secret)
		}
	}
}

func TestUpdateSecretProperties(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	secret, err := createRandomName(t, "secret")
	require.NoError(t, err)
	value, err := createRandomName(t, "value")
	require.NoError(t, err)

	_, err = client.SetSecret(context.Background(), secret, value, nil)
	require.NoError(t, err)

	defer cleanUpSecret(t, client, secret)

	getResp, err := client.GetSecret(context.Background(), secret, nil)
	require.NoError(t, err)
	require.Equal(t, *getResp.Value, value)

	params := Properties{
		ContentType: to.StringPtr("password"),
		Tags: map[string]string{
			"Tag1": "TagVal1",
		},
	}

	_, err = client.UpdateSecretProperties(context.Background(), secret, params, nil)
	require.NoError(t, err)

	getResp, err = client.GetSecret(context.Background(), secret, nil)
	require.NoError(t, err)
	require.Equal(t, *getResp.Value, value)
	require.Equal(t, getResp.Tags["Tag1"], "TagVal1")
	require.Equal(t, *getResp.ContentType, "password")
}

func TestBeginRecoverDeletedSecret(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	secret, err := createRandomName(t, "secret")
	require.NoError(t, err)
	value, err := createRandomName(t, "value")
	require.NoError(t, err)

	_, err = client.SetSecret(context.Background(), secret, value, nil)
	require.NoError(t, err)

	defer cleanUpSecret(t, client, secret)

	pollerResp, err := client.BeginDeleteSecret(context.Background(), secret, nil)
	require.NoError(t, err)

	_, err = pollerResp.PollUntilDone(context.Background(), delay())
	require.NoError(t, err)

	resp, err := client.BeginRecoverDeletedSecret(context.Background(), secret, nil)
	require.NoError(t, err)

	_, err = resp.PollUntilDone(context.Background(), delay())
	require.NoError(t, err)

	_, err = client.SetSecret(context.Background(), secret, value, nil)
	require.NoError(t, err)

	getResp, err := client.GetSecret(context.Background(), secret, nil)
	require.NoError(t, err)
	require.Equal(t, *getResp.Value, value)
}

func TestBackupSecret(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t)
	require.NoError(t, err)

	secret, err := createRandomName(t, "secrets")
	require.NoError(t, err)
	value, err := createRandomName(t, "value")
	require.NoError(t, err)

	_, err = client.SetSecret(context.Background(), secret, value, nil)
	require.NoError(t, err)

	defer cleanUpSecret(t, client, secret)

	backupResp, err := client.BackupSecret(context.Background(), secret, nil)
	require.NoError(t, err)
	require.Greater(t, len(backupResp.Value), 0)

	respPoller, err := client.BeginDeleteSecret(context.Background(), secret, nil)
	require.NoError(t, err)
	_, err = respPoller.PollUntilDone(context.Background(), delay())
	require.NoError(t, err)

	_, err = client.PurgeDeletedSecret(context.Background(), secret, nil)
	require.NoError(t, err)

	_, err = client.GetSecret(context.Background(), secret, nil)
	var httpErr azcore.HTTPResponse
	require.True(t, errors.As(err, &httpErr))
	require.Equal(t, httpErr.RawResponse().StatusCode, http.StatusNotFound)

	_, err = client.GetDeletedSecret(context.Background(), secret, nil)
	require.True(t, errors.As(err, &httpErr))
	require.Equal(t, httpErr.RawResponse().StatusCode, http.StatusNotFound)

	// Poll this operation manually
	var restoreResp RestoreSecretBackupResponse
	var i int
	for i = 0; i < 10; i++ {
		restoreResp, err = client.RestoreSecretBackup(context.Background(), backupResp.Value, nil)
		if err == nil {
			break
		}
		time.Sleep(delay())
	}
	require.NoError(t, err)
	require.Contains(t, *restoreResp.ID, secret)

	// Now the Secret should be Get-able
	_, err = client.GetSecret(context.Background(), secret, nil)
	require.NoError(t, err)
}
