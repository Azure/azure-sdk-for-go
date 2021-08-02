package aztable

import (
	"context"
	"fmt"
	"hash/fnv"
	"os"
	"strings"
	"testing"
	"time"

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
	}
	accountKey, ok := os.LookupEnv("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY")
	if !ok {
		t.Log("STORAGE KEY")
		accountKey = "fakekey"
	}
	serviceURL := storageURI(accountName, "core.windows.net")
	cred, err := NewSharedKeyCredential(accountName, accountKey)
	require.NoError(t, err)
	tableName, err := createRandomName(t, "tableName")
	require.NoError(t, err)
	return createTableClientForRecording(t, tableName, serviceURL, cred)
}

func createCosmosTableClient(t *testing.T) (*TableClient, error) {
	accountName, ok := os.LookupEnv("TABLES_COSMOS_ACCOUNT_NAME")
	if !ok {
		accountName = "fakecosmosaccount"
	}
	accountKey, ok := os.LookupEnv("TABLES_PRIMARY_COSMOS_ACCOUNT_KEY")
	if !ok {
		t.Log("COSMOS KEY")
		accountKey = "fakekey"
	}
	serviceURL := cosmosURI(accountName, "cosmos.azure.com")
	cred, err := createSharedKey(accountName, accountKey)
	require.NoError(t, err)
	tableName, err := createRandomName(t, "tableName")
	require.NoError(t, err)
	return createTableClientForRecording(t, tableName, serviceURL, cred)
}

func createStorageServiceClient(t *testing.T) (*TableServiceClient, error) {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		accountName = "fakestorageaccount"
	}
	serviceURL := storageURI(accountName, "core.windows.net")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)
	return createTableServiceClientForRecording(t, serviceURL, cred)
}

func createCosmosServiceClient(t *testing.T) (*TableServiceClient, error) {
	accountName := getEnvVariable("TABLES_COSMOS_ACCOUNT_NAME", "fakestorageaccount")
	accountKey := getEnvVariable("TABLES_PRIMARY_COSMOS_ACCOUNT_KEY", "fakekey")
	serviceURL := cosmosURI(accountName, "cosmos.azure.com")
	cred, err := createSharedKey(accountName, accountKey)
	require.NoError(t, err)
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
			_, err := service.DeleteTable(ctx, *v.TableName)
			if err != nil {
				return err
			}
		}
	}
	return pager.Err()
}

// This looks up an environment variable and if it is not found, returns the recordedValue
func getEnvVariable(varName string, recordedValue string) string {
	val, ok := os.LookupEnv(varName)
	if !ok {
		return recordedValue
	}
	return val
}

func createSharedKey(accountName, accountKey string) (azcore.Credential, error) {
	if os.Getenv("AZURE_RECORD_MODE") == "record" {
		return NewSharedKeyCredential(accountName, accountKey)
	}

	return NewFakeCredential(accountName, accountKey), nil
}

type FakeCredential struct {
	accountName string
	accountKey  string
}

func NewFakeCredential(accountName, accountKey string) *FakeCredential {
	return &FakeCredential{
		accountName: accountName,
		accountKey:  accountKey,
	}
}

func (f *FakeCredential) AuthenticationPolicy(azcore.AuthenticationPolicyOptions) azcore.Policy {
	return azcore.PolicyFunc(func(req *azcore.Request) (*azcore.Response, error) {
		// Do nothing, authentication headers are stripped from recordings
		authHeader := strings.Join([]string{"Authorization ", f.accountName, ":", f.accountKey}, "")
		req.Request.Header.Set(azcore.HeaderAuthorization, authHeader)
		return req.Next()
	})
}

func testSleep() {
	if os.Getenv("AZURE_RECORD_MODE") == "record" {
		time.Sleep(45 * time.Second)
	}
}
