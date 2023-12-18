// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestSASServiceClient(t *testing.T) {
	recording.LiveOnly(t)
	accountName := os.Getenv("TABLES_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY")
	cred, err := NewSharedKeyCredential(accountName, accountKey)
	require.NoError(t, err)

	serviceClient, err := NewServiceClientWithSharedKey(fmt.Sprintf("https://%s.table.core.windows.net/", accountName), cred, nil)
	require.NoError(t, err)

	tableName, err := createRandomName(t, tableNamePrefix)
	require.NoError(t, err)

	delete := func() {
		_, err := serviceClient.DeleteTable(ctx, tableName, nil)
		require.NoError(t, err)
	}
	defer delete()

	_, err = serviceClient.CreateTable(ctx, tableName, nil)
	require.NoError(t, err)

	resources := AccountSASResourceTypes{
		Object:    true,
		Service:   true,
		Container: true,
	}
	permissions := AccountSASPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Update: true,
		Delete: true,
	}

	start := time.Now().Add(-1 * time.Hour).UTC()
	expiry := time.Now().Add(24 * time.Hour).UTC()

	sasUrl, err := serviceClient.GetAccountSASURL(resources, permissions, start, expiry)
	require.NoError(t, err)

	err = recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)
	svcClient, err := createServiceClientForRecordingWithNoCredential(t, sasUrl, tracing.Provider{})
	require.NoError(t, err)
	defer require.NoError(t, recording.Stop(t, nil))

	_, err = svcClient.CreateTable(ctx, tableName+"002", nil)
	require.NoError(t, err)

	_, err = svcClient.DeleteTable(ctx, tableName+"002", nil)
	require.NoError(t, err)
}

func TestSASClient(t *testing.T) {
	recording.LiveOnly(t)
	accountName := os.Getenv("TABLES_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY")
	cred, err := NewSharedKeyCredential(accountName, accountKey)
	require.NoError(t, err)

	serviceClient, err := NewServiceClientWithSharedKey(fmt.Sprintf("https://%s.table.core.windows.net/", accountName), cred, nil)
	require.NoError(t, err)

	tableName, err := createRandomName(t, tableNamePrefix)
	require.NoError(t, err)

	delete := func() {
		_, err := serviceClient.DeleteTable(ctx, tableName, nil)
		require.NoError(t, err)
	}
	defer delete()

	_, err = serviceClient.CreateTable(ctx, tableName, nil)
	require.NoError(t, err)

	permissions := SASPermissions{
		Read: true,
		Add:  true,
	}

	start := time.Now().Add(-1 * time.Hour).UTC()
	expiry := time.Now().Add(24 * time.Hour).UTC()

	c := serviceClient.NewClient(tableName)
	sasUrl, err := c.GetTableSASURL(permissions, start, expiry)
	require.NoError(t, err)

	err = recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)
	client, err := createClientForRecordingWithNoCredential(t, "", sasUrl, tracing.Provider{})
	require.NoError(t, err)
	defer require.NoError(t, recording.Stop(t, nil))

	entity := map[string]string{
		"PartitionKey": "pk001",
		"RowKey":       "rk001",
		"Value":        "5",
	}
	marshalled, err := json.Marshal(entity)
	require.NoError(t, err)

	_, err = client.AddEntity(ctx, marshalled, nil)
	require.NoError(t, err)
}

func TestSASClientReadOnly(t *testing.T) {
	recording.LiveOnly(t)
	accountName := os.Getenv("TABLES_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY")
	cred, err := NewSharedKeyCredential(accountName, accountKey)
	require.NoError(t, err)

	serviceClient, err := NewServiceClientWithSharedKey(fmt.Sprintf("https://%s.table.core.windows.net/", accountName), cred, nil)
	require.NoError(t, err)

	tableName, err := createRandomName(t, tableNamePrefix)
	require.NoError(t, err)

	delete := func() {
		_, err := serviceClient.DeleteTable(ctx, tableName, nil)
		require.NoError(t, err)
	}
	defer delete()

	_, err = serviceClient.CreateTable(ctx, tableName, nil)
	require.NoError(t, err)

	client := serviceClient.NewClient(tableName)
	err = insertNEntities("pk001", 4, client)
	require.NoError(t, err)

	permissions := SASPermissions{
		Read: true,
	}

	start := time.Now().Add(-1 * time.Hour).UTC()
	expiry := time.Now().Add(24 * time.Hour).UTC()

	c := serviceClient.NewClient(tableName)
	sasUrl, err := c.GetTableSASURL(permissions, start, expiry)
	require.NoError(t, err)

	err = recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)
	client, err = createClientForRecordingWithNoCredential(t, "", sasUrl, tracing.Provider{})
	require.NoError(t, err)
	defer require.NoError(t, recording.Stop(t, nil))

	entity := map[string]string{
		"PartitionKey": "pk001",
		"RowKey":       "rk001",
		"Value":        "5",
	}
	marshalled, err := json.Marshal(entity)
	require.NoError(t, err)

	// Failure on a read
	_, err = client.AddEntity(ctx, marshalled, nil)
	require.Error(t, err)
	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, "AuthorizationPermissionMismatch", httpErr.ErrorCode)

	// Success on a list
	pager := client.NewListEntitiesPager(nil)
	count := 0
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		require.NoError(t, err)
		count += len(resp.Entities)
	}

	require.Equal(t, 4, count)
}

func TestSASCosmosClientReadOnly(t *testing.T) {
	recording.LiveOnly(t)
	accountName := os.Getenv("TABLES_COSMOS_ACCOUNT_NAME")
	accountKey := os.Getenv("TABLES_PRIMARY_COSMOS_ACCOUNT_KEY")
	cred, err := NewSharedKeyCredential(accountName, accountKey)
	require.NoError(t, err)

	serviceClient, err := NewServiceClientWithSharedKey(fmt.Sprintf("https://%s.table.cosmos.azure.com/", accountName), cred, nil)
	require.NoError(t, err)

	tableName, err := createRandomName(t, tableNamePrefix)
	require.NoError(t, err)

	delete := func() {
		_, err := serviceClient.DeleteTable(ctx, tableName, nil)
		require.NoError(t, err)
	}
	defer delete()

	_, err = serviceClient.CreateTable(ctx, tableName, nil)
	require.NoError(t, err)

	client := serviceClient.NewClient(tableName)
	err = insertNEntities("pk001", 4, client)
	require.NoError(t, err)

	permissions := SASPermissions{
		Read: true,
	}

	start := time.Now().Add(-1 * time.Hour).UTC()
	expiry := time.Now().Add(24 * time.Hour).UTC()

	c := serviceClient.NewClient(tableName)
	sasUrl, err := c.GetTableSASURL(permissions, start, expiry)
	require.NoError(t, err)

	err = recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)
	client, err = createClientForRecordingWithNoCredential(t, "", sasUrl, tracing.Provider{})
	require.NoError(t, err)
	defer require.NoError(t, recording.Stop(t, nil))

	entity := map[string]string{
		"PartitionKey": "pk001",
		"RowKey":       "rk001",
		"Value":        "5",
	}
	marshalled, err := json.Marshal(entity)
	require.NoError(t, err)

	// Failure on a read
	_, err = client.AddEntity(ctx, marshalled, nil)
	require.Error(t, err)
	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, "Forbidden", httpErr.ErrorCode)

	// Success on a list
	pager := client.NewListEntitiesPager(nil)
	count := 0
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		require.NoError(t, err)
		count += len(resp.Entities)
	}

	require.Equal(t, 4, count)
}
