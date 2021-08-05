// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"
	"fmt"
	"hash/fnv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var AADAuthenticationScope = "https://storage.azure.com/.default"

func createTableClientForRecording(t *testing.T, tableName string, serviceURL string, cred azcore.Credential) (*TableClient, error) {
	policy := recording.NewRecordingPolicy(&recording.RecordingOptions{UseHTTPS: true})
	client, err := recording.GetHTTPClient()
	require.NoError(t, err)
	options := &TableClientOptions{
		Scopes:         []string{AADAuthenticationScope},
		PerCallOptions: []azcore.Policy{policy},
		HTTPClient:     client,
	}
	return NewTableClient(tableName, serviceURL, cred, options)
}

func createTableServiceClientForRecording(t *testing.T, serviceURL string, cred azcore.Credential) (*TableServiceClient, error) {
	policy := recording.NewRecordingPolicy(&recording.RecordingOptions{UseHTTPS: true})
	client, err := recording.GetHTTPClient()
	require.NoError(t, err)
	options := &TableClientOptions{
		Scopes:         []string{AADAuthenticationScope},
		PerCallOptions: []azcore.Policy{policy},
		HTTPClient:     client,
	}
	return NewTableServiceClient(serviceURL, cred, options)
}

func initClientTest(t *testing.T, service string, createTable bool) (*TableClient, func()) {
	var client *TableClient
	var err error
	if service == string(StorageEndpoint) {
		client, err = createStorageTableClient(t)
		require.NoError(t, err)
	} else if service == string(CosmosEndpoint) {
		client, err = createCosmosTableClient(t)
		require.NoError(t, err)
	}

	err = recording.StartRecording(t, nil)
	require.NoError(t, err)

	if createTable {
		_, err = client.Create(context.Background())
		require.NoError(t, err)
	}

	return client, func() {
		_, err = client.Delete(context.Background(), nil)
		require.NoError(t, err)
		err = recording.StopRecording(t, nil)
		require.NoError(t, err)
	}
}

func initServiceTest(t *testing.T, service string) (*TableServiceClient, func()) {
	var client *TableServiceClient
	var err error
	if service == string(StorageEndpoint) {
		client, err = createStorageServiceClient(t)
		require.NoError(t, err)
	} else if service == string(CosmosEndpoint) {
		client, err = createCosmosServiceClient(t)
		require.NoError(t, err)
	}

	err = recording.StartRecording(t, nil)
	require.NoError(t, err)

	return client, func() {
		err = recording.StopRecording(t, nil)
		require.NoError(t, err)
	}
}

func getAADCredential(t *testing.T) (azcore.Credential, error) {
	if recording.InPlayback() {
		return recording.NewFakeCredential("fakestorageaccount", "fakeAccountKey"), nil
	}

	accountName := recording.GetEnvVariable(t, "TABLES_STORAGE_ACCOUNT_NAME", "fakestorageaccount")

	err := recording.AddUriSanitizer("fakestorageaccount", accountName, nil)
	require.NoError(t, err)

	return azidentity.NewDefaultAzureCredential(nil)
}

func getSharedKeyCredential(t *testing.T) (azcore.Credential, error) {
	if recording.InPlayback() {
		return recording.NewFakeCredential("fakestorageaccount", "fakeAccountKey"), nil
	}

	accountName := recording.GetEnvVariable(t, "TABLES_COSMOS_ACCOUNT_NAME", "fakestorageaccount")
	accountKey := recording.GetEnvVariable(t, "TABLES_PRIMARY_COSMOS_ACCOUNT_KEY", "fakeAccountKey")

	err := recording.AddUriSanitizer("fakestorageaccount", accountName, nil)
	require.NoError(t, err)

	return NewSharedKeyCredential(accountName, accountKey)
}

func createStorageTableClient(t *testing.T) (*TableClient, error) {
	var cred azcore.Credential
	accountName := "fakestorageaccount"

	cred, err := getAADCredential(t)
	require.NoError(t, err)

	serviceURL := storageURI(accountName, "core.windows.net")

	tableName, err := createRandomName(t, "tableName")
	require.NoError(t, err)

	return createTableClientForRecording(t, tableName, serviceURL, cred)
}

func createCosmosTableClient(t *testing.T) (*TableClient, error) {
	var cred azcore.Credential
	accountName := "fakestorageaccount"

	cred, err := getSharedKeyCredential(t)
	require.NoError(t, err)

	serviceURL := cosmosURI(accountName, "cosmos.azure.com")

	tableName, err := createRandomName(t, "tableName")
	require.NoError(t, err)

	return createTableClientForRecording(t, tableName, serviceURL, cred)
}

func createStorageServiceClient(t *testing.T) (*TableServiceClient, error) {
	var cred azcore.Credential
	accountName := "fakestorageaccount"

	cred, err := getAADCredential(t)
	require.NoError(t, err)

	serviceURL := storageURI(accountName, "core.windows.net")

	return createTableServiceClientForRecording(t, serviceURL, cred)
}

func createCosmosServiceClient(t *testing.T) (*TableServiceClient, error) {
	var cred azcore.Credential
	accountName := "fakestorageaccount"

	cred, err := getSharedKeyCredential(t)
	require.NoError(t, err)

	serviceURL := cosmosURI(accountName, "cosmos.azure.com")

	return createTableServiceClientForRecording(t, serviceURL, cred)
}

func createRandomName(t *testing.T, prefix string) (string, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(t.Name()))
	return prefix + fmt.Sprint(h.Sum32()), err
}

func clearAllTables(service *TableServiceClient) error {
	pager := service.ListTables(nil)
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		for _, v := range resp.TableQueryResponse.Value {
			_, err := service.DeleteTable(ctx, *v.TableName, nil)
			if err != nil {
				return err
			}
		}
	}
	return pager.Err()
}
