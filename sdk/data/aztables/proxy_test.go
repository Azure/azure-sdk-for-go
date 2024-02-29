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
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

const recordingDirectory = "sdk/data/aztables/testdata"

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		proxy, err := recording.StartTestProxy(recordingDirectory, nil)
		if err != nil {
			panic(err)
		}

		defer func() {
			err := recording.StopTestProxy(proxy)
			if err != nil {
				panic(err)
			}
		}()
	}

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
	return m.Run()
}

const tableNamePrefix = "tableName"

type FakeCredential struct {
	accountName string
	accountKey  string
}

func (f *FakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "***", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func NewFakeCredential(accountName, accountKey string) *FakeCredential {
	return &FakeCredential{
		accountName: accountName,
		accountKey:  accountKey,
	}
}

func createClientForRecording(t *testing.T, tableName string, serviceURL string, cred SharedKeyCredential, tp tracing.Provider) (*Client, error) {
	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{ClientOptions: azcore.ClientOptions{
		TracingProvider: tp,
		Transport:       client,
	}}
	if !strings.HasSuffix(serviceURL, "/") && tableName != "" {
		serviceURL += "/"
	}
	serviceURL += tableName

	return NewClientWithSharedKey(serviceURL, &cred, options)
}

func createClientForRecordingWithNoCredential(t *testing.T, tableName string, serviceURL string, tp tracing.Provider) (*Client, error) {
	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{ClientOptions: azcore.ClientOptions{
		TracingProvider: tp,
		Transport:       client,
	}}
	if !strings.HasSuffix(serviceURL, "/") && tableName != "" {
		serviceURL += "/"
	}
	serviceURL += tableName

	return NewClientWithNoCredential(serviceURL, options)
}

func createServiceClientForRecording(t *testing.T, serviceURL string, cred SharedKeyCredential, tp tracing.Provider) (*ServiceClient, error) {
	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{ClientOptions: azcore.ClientOptions{
		TracingProvider: tp,
		Transport:       client,
	}}
	return NewServiceClientWithSharedKey(serviceURL, &cred, options)
}

func createServiceClientForRecordingWithNoCredential(t *testing.T, serviceURL string, tp tracing.Provider) (*ServiceClient, error) {
	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{ClientOptions: azcore.ClientOptions{
		TracingProvider: tp,
		Transport:       client,
	}}
	return NewServiceClientWithNoCredential(serviceURL, options)
}

func initClientTest(t *testing.T, service string, createTable bool, tp tracing.Provider) (*Client, func()) {
	var client *Client
	var err error
	if service == string(storageEndpoint) {
		client, err = createStorageClient(t, tp)
		require.NoError(t, err)
	} else if service == string(cosmosEndpoint) {
		client, err = createCosmosClient(t, tp)
		require.NoError(t, err)
	}

	err = recording.Start(t, recordingDirectory, nil)
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

func initServiceTest(t *testing.T, service string, tp tracing.Provider) (*ServiceClient, func()) {
	var client *ServiceClient
	var err error
	if service == string(storageEndpoint) {
		client, err = createStorageServiceClient(t, tp)
		require.NoError(t, err)
	} else if service == string(cosmosEndpoint) {
		client, err = createCosmosServiceClient(t, tp)
		require.NoError(t, err)
	}

	err = recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)

	return client, func() {
		err = recording.Stop(t, nil)
		require.NoError(t, err)
	}
}

func getSharedKeyCredential() (*SharedKeyCredential, error) {
	if recording.GetRecordMode() == "playback" {
		return NewSharedKeyCredential("accountName", "daaaaaaaaaabbbbbbbbbbcccccccccccccccccccdddddddddddddddddddeeeeeeeeeeefffffffffffggggg==")
	}

	accountName := recording.GetEnvVariable("TABLES_COSMOS_ACCOUNT_NAME", "fakeaccount")
	accountKey := recording.GetEnvVariable("TABLES_PRIMARY_COSMOS_ACCOUNT_KEY", "fakeAccountKey")

	return NewSharedKeyCredential(accountName, accountKey)
}

func createStorageClient(t *testing.T, tp tracing.Provider) (*Client, error) {
	var cred *SharedKeyCredential
	var err error
	accountName := recording.GetEnvVariable("TABLES_STORAGE_ACCOUNT_NAME", "fakeaccount")
	accountKey := recording.GetEnvVariable("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY", "fakeaccountkey")

	if recording.GetRecordMode() == "playback" {
		cred, err = getSharedKeyCredential()
		require.NoError(t, err)
	} else {
		cred, err = NewSharedKeyCredential(accountName, accountKey)
		require.NoError(t, err)
	}

	serviceURL := storageURI(accountName)

	tableName, err := createRandomName(t, tableNamePrefix)
	require.NoError(t, err)

	return createClientForRecording(t, tableName, serviceURL, *cred, tp)
}

func createCosmosClient(t *testing.T, tp tracing.Provider) (*Client, error) {
	var cred *SharedKeyCredential
	accountName := recording.GetEnvVariable("TABLES_COSMOS_ACCOUNT_NAME", "fakeaccount")
	if recording.GetRecordMode() == "playback" {
		accountName = "fakeaccount"
	}

	cred, err := getSharedKeyCredential()
	require.NoError(t, err)

	serviceURL := cosmosURI(accountName)

	tableName, err := createRandomName(t, tableNamePrefix)
	require.NoError(t, err)

	return createClientForRecording(t, tableName, serviceURL, *cred, tp)
}

func createStorageServiceClient(t *testing.T, tp tracing.Provider) (*ServiceClient, error) {
	var cred *SharedKeyCredential
	var err error
	accountName := recording.GetEnvVariable("TABLES_STORAGE_ACCOUNT_NAME", "fakeaccount")
	accountKey := recording.GetEnvVariable("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY", "fakeaccountkey")

	if recording.GetRecordMode() == "playback" {
		cred, err = getSharedKeyCredential()
		require.NoError(t, err)
	} else {
		cred, err = NewSharedKeyCredential(accountName, accountKey)
		require.NoError(t, err)
	}

	serviceURL := storageURI(accountName)

	return createServiceClientForRecording(t, serviceURL, *cred, tp)
}

func createCosmosServiceClient(t *testing.T, tp tracing.Provider) (*ServiceClient, error) {
	var cred *SharedKeyCredential
	accountName := recording.GetEnvVariable("TABLES_COSMOS_ACCOUNT_NAME", "fakeaccount")
	if recording.GetRecordMode() == "playback" {
		accountName = "fakeaccount"
	}

	cred, err := getSharedKeyCredential()
	require.NoError(t, err)

	serviceURL := cosmosURI(accountName)

	return createServiceClientForRecording(t, serviceURL, *cred, tp)
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
