// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

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
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// 1. Set up session level sanitizers
	switch recording.GetRecordMode() {
	case recording.PlaybackMode:
		err := recording.SetDefaultMatcher(nil, &recording.SetDefaultMatcherOptions{
			ExcludedHeaders: []string{":path", ":auth", ":method", ":scheme"},
		})
		if err != nil {
			panic(err)
		}
	case recording.RecordingMode:
		for _, val := range []string{"TABLES_COSMOS_ACCOUNT_NAME", "TABLES_STORAGE_ACCOUNT_NAME"} {
			account, ok := os.LookupEnv(val)
			if !ok {
				fmt.Printf("Could not find environment variable: %s", val)
				os.Exit(1)
			}

			err := recording.AddGeneralRegexSanitizer("fakeaccount", account, nil)
			if err != nil {
				panic(err)
			}
		}

	}
	// Run tests
	exitVal := m.Run()

	// 3. Reset
	// TODO: Add after sanitizer PR
	if recording.GetRecordMode() != "live" {
		err := recording.ResetProxy(nil)
		if err != nil {
			panic(err)
		}
	}

	// 4. Error out if applicable
	os.Exit(exitVal)
}

var pathToPackage = "sdk/data/aztables/testdata"

const tableNamePrefix = "tableName"

type FakeCredential struct {
	accountName string
	accountKey  string
}

func (f *FakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	return &azcore.AccessToken{Token: "***", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func NewFakeCredential(accountName, accountKey string) *FakeCredential {
	return &FakeCredential{
		accountName: accountName,
		accountKey:  accountKey,
	}
}

func createClientForRecording(t *testing.T, tableName string, serviceURL string, cred SharedKeyCredential) (*Client, error) {
	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{ClientOptions: azcore.ClientOptions{
		Transport: client,
	}}
	if !strings.HasSuffix(serviceURL, "/") && tableName != "" {
		serviceURL += "/"
	}
	serviceURL += tableName

	return NewClientWithSharedKey(serviceURL, &cred, options)
}

func createClientForRecordingWithNoCredential(t *testing.T, tableName string, serviceURL string) (*Client, error) {
	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{ClientOptions: azcore.ClientOptions{
		Transport: client,
	}}
	if !strings.HasSuffix(serviceURL, "/") && tableName != "" {
		serviceURL += "/"
	}
	serviceURL += tableName

	return NewClientWithNoCredential(serviceURL, options)
}

func createServiceClientForRecording(t *testing.T, serviceURL string, cred SharedKeyCredential) (*ServiceClient, error) {
	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{ClientOptions: azcore.ClientOptions{
		Transport: client,
	}}
	return NewServiceClientWithSharedKey(serviceURL, &cred, options)
}

func createServiceClientForRecordingWithNoCredential(t *testing.T, serviceURL string) (*ServiceClient, error) {
	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{ClientOptions: azcore.ClientOptions{
		Transport: client,
	}}
	return NewServiceClientWithNoCredential(serviceURL, options)
}

func initClientTest(t *testing.T, service string, createTable bool) (*Client, func()) {
	var client *Client
	var err error
	if service == string(storageEndpoint) {
		client, err = createStorageClient(t)
		require.NoError(t, err)
	} else if service == string(cosmosEndpoint) {
		client, err = createCosmosClient(t)
		require.NoError(t, err)
	}

	err = recording.Start(t, pathToPackage, nil)
	require.NoError(t, err)

	if createTable {
		_, err = client.CreateTable(ctx, nil)
		require.NoError(t, err)
	}

	return client, func() {
		_, err = client.Delete(ctx, nil)
		require.NoError(t, err)
		err = recording.Stop(t, nil)
		require.NoError(t, err)
	}
}

func initServiceTest(t *testing.T, service string) (*ServiceClient, func()) {
	var client *ServiceClient
	var err error
	if service == string(storageEndpoint) {
		client, err = createStorageServiceClient(t)
		require.NoError(t, err)
	} else if service == string(cosmosEndpoint) {
		client, err = createCosmosServiceClient(t)
		require.NoError(t, err)
	}

	err = recording.Start(t, pathToPackage, nil)
	require.NoError(t, err)

	return client, func() {
		err = recording.Stop(t, nil)
		require.NoError(t, err)
	}
}

func getSharedKeyCredential(t *testing.T) (*SharedKeyCredential, error) {
	if recording.GetRecordMode() == "playback" {
		return NewSharedKeyCredential("accountName", "daaaaaaaaaabbbbbbbbbbcccccccccccccccccccdddddddddddddddddddeeeeeeeeeeefffffffffffggggg==")
	}

	accountName := recording.GetEnvVariable("TABLES_COSMOS_ACCOUNT_NAME", "fakeaccount")
	accountKey := recording.GetEnvVariable("TABLES_PRIMARY_COSMOS_ACCOUNT_KEY", "fakeAccountKey")

	return NewSharedKeyCredential(accountName, accountKey)
}

func createStorageClient(t *testing.T) (*Client, error) {
	var cred *SharedKeyCredential
	var err error
	accountName := recording.GetEnvVariable("TABLES_STORAGE_ACCOUNT_NAME", "fakeaccount")
	accountKey := recording.GetEnvVariable("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY", "fakeaccountkey")

	if recording.GetRecordMode() == "playback" {
		cred, err = getSharedKeyCredential(t)
		require.NoError(t, err)
	} else {
		cred, err = NewSharedKeyCredential(accountName, accountKey)
		require.NoError(t, err)
	}

	serviceURL := storageURI(accountName)

	tableName, err := createRandomName(t, tableNamePrefix)
	require.NoError(t, err)

	return createClientForRecording(t, tableName, serviceURL, *cred)
}

func createCosmosClient(t *testing.T) (*Client, error) {
	var cred *SharedKeyCredential
	accountName := recording.GetEnvVariable("TABLES_COSMOS_ACCOUNT_NAME", "fakeaccount")
	if recording.GetRecordMode() == "playback" {
		accountName = "fakeaccount"
	}

	cred, err := getSharedKeyCredential(t)
	require.NoError(t, err)

	serviceURL := cosmosURI(accountName)

	tableName, err := createRandomName(t, tableNamePrefix)
	require.NoError(t, err)

	return createClientForRecording(t, tableName, serviceURL, *cred)
}

func createStorageServiceClient(t *testing.T) (*ServiceClient, error) {
	var cred *SharedKeyCredential
	var err error
	accountName := recording.GetEnvVariable("TABLES_STORAGE_ACCOUNT_NAME", "fakeaccount")
	accountKey := recording.GetEnvVariable("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY", "fakeaccountkey")

	if recording.GetRecordMode() == "playback" {
		cred, err = getSharedKeyCredential(t)
		require.NoError(t, err)
	} else {
		cred, err = NewSharedKeyCredential(accountName, accountKey)
		require.NoError(t, err)
	}

	serviceURL := storageURI(accountName)

	return createServiceClientForRecording(t, serviceURL, *cred)
}

func createCosmosServiceClient(t *testing.T) (*ServiceClient, error) {
	var cred *SharedKeyCredential
	accountName := recording.GetEnvVariable("TABLES_COSMOS_ACCOUNT_NAME", "fakeaccount")
	if recording.GetRecordMode() == "playback" {
		accountName = "fakeaccount"
	}

	cred, err := getSharedKeyCredential(t)
	require.NoError(t, err)

	serviceURL := cosmosURI(accountName)

	return createServiceClientForRecording(t, serviceURL, *cred)
}

func createRandomName(t *testing.T, prefix string) (string, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(t.Name()))
	return prefix + fmt.Sprint(h.Sum32()), err
}

func clearAllTables(service *ServiceClient) error {
	pager := service.NewListTablesPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, v := range resp.Tables {
			_, err := service.DeleteTable(ctx, *v.Name, nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
