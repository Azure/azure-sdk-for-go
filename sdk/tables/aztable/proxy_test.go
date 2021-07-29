package aztable

import (
	"context"
	"fmt"
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

func createTableServiceClientForRecording(t *testing.T, serviceURL string) (*TableServiceClient, error) {
	policy := recording.NewRecordingPolicy(&recording.RecordingOptions{UseHTTPS: true})
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)
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

func createStorageTableClient(t *testing.T) (*TableClient, error) {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		accountName = "fakestorageaccount"
		fmt.Println("STORAGE KEY")
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
		fmt.Println("COSMOS KEY")
		accountKey = "fakekey"
	}
	serviceURL := cosmosURI(accountName, "cosmos.azure.com")
	cred, err := NewSharedKeyCredential(accountName, accountKey)
	require.NoError(t, err)
	return createTableClientForRecording(t, "createPseudoRandomName", serviceURL, cred)
}

// func Test_TestProxyPolicy(t *testing.T) {
// 	require := require.New(t)
// 	err := recording.StartRecording(t, nil)
// 	require.NoError(err)
// 	defer recording.StopRecording(t, nil)

// 	client, err := createTableClientForRecording(t, "testproxy", "https://seankaneprim.table.core.windows.net")
// 	require.NoError(err)

// 	_, err = client.Create(ctx)
// 	require.NoError(err)

// 	_, err = client.Delete(ctx)
// 	require.NoError(err)
// }
