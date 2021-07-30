package aztable

import (
	"context"
	"fmt"
	"hash/fnv"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var AADAuthenticationScope = "https://storage.azure.com/.default"

func createTableClientForRecording(t *testing.T, tableName string, serviceURL string, cred azcore.Credential) (*TableClient, error) {
	policy := recording.NewRecordingPolicy(&recording.RecordingOptions{UseHTTPS: true})
	options := &TableClientOptions{
		Scopes:         []string{AADAuthenticationScope},
		PerCallOptions: []azcore.Policy{policy},
	}
	return NewTableClient(tableName, serviceURL, cred, options)
}

func createTableServiceClientForRecording(t *testing.T, serviceURL string, cred azcore.Credential) (*TableServiceClient, error) {
	policy := recording.NewRecordingPolicy(&recording.RecordingOptions{UseHTTPS: true})
	options := &TableClientOptions{
		Scopes:         []string{AADAuthenticationScope},
		PerCallOptions: []azcore.Policy{policy},
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
		_, err = client.Delete(context.Background())
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

func createStorageTableClient(t *testing.T) (*TableClient, error) {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		accountName = "fakestorageaccount"
		t.Log("STORAGE KEY")
	}
	serviceURL := storageURI(accountName, "core.windows.net")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)
	return createTableClientForRecording(t, "createPseudoRandomName", serviceURL, cred)
}

func createCosmosTableClient(t *testing.T) (*TableClient, error) {
	accountName, ok := os.LookupEnv("TABLES_COSMOS_ACCOUNT_NAME")
	if !ok {
		accountName = "fakestorageaccount"
	}
	accountKey, ok := os.LookupEnv("TABLES_PRIMARY_COSMOS_ACCOUNT_KEY")
	if !ok {
		t.Log("COSMOS KEY")
		accountKey = "fakekey"
	}
	serviceURL := cosmosURI(accountName, "cosmos.azure.com")
	cred, err := NewSharedKeyCredential(accountName, accountKey)
	require.NoError(t, err)
	return createTableClientForRecording(t, "createPseudoRandomName", serviceURL, cred)
}

func createStorageServiceClient(t *testing.T) (*TableServiceClient, error) {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		accountName = "fakestorageaccount"
		t.Log("STORAGE KEY")
	}
	serviceURL := storageURI(accountName, "core.windows.net")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)
	return createTableServiceClientForRecording(t, serviceURL, cred)
}

func createCosmosServiceClient(t *testing.T) (*TableServiceClient, error) {
	accountName, ok := os.LookupEnv("TABLES_COSMOS_ACCOUNT_NAME")
	if !ok {
		t.Log("No cosmos account name provided.")
		accountName = "fakestorageaccount"
	}
	accountKey, ok := os.LookupEnv("TABLES_PRIMARY_COSMOS_ACCOUNT_KEY")
	if !ok {
		t.Log("No key provided for cosmos")
		accountKey = "fakekey"
	}
	serviceURL := cosmosURI(accountName, "cosmos.azure.com")
	cred, err := NewSharedKeyCredential(accountName, accountKey)
	require.NoError(t, err)
	return createTableServiceClientForRecording(t, serviceURL, cred)
}

func createRandomName(t *testing.T, prefix string) (string, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(t.Name()))
	return prefix + fmt.Sprint(h.Sum32()), err
}

func clearAllTables2(service *TableServiceClient) error {
	pager := service.ListTables(nil)
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		for _, v := range resp.TableQueryResponse.Value {
			_, err := service.DeleteTable(ctx, *v.TableName)
			if err != nil {
				return err
			}
		}
	}
	return pager.Err()
}
