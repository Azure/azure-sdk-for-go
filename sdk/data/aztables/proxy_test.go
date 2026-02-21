// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"fmt"
	"hash/fnv"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/stretchr/testify/require"
)

const recordingDirectory = "sdk/data/aztables/testdata"
const fakeAccount = recording.SanitizedValue

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	if recording.GetRecordMode() != recording.LiveMode {
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
	for _, v := range []string{"TABLES_COSMOS_ACCOUNT_NAME", "TABLES_STORAGE_ACCOUNT_NAME"} {
		account := recording.GetEnvVariable(v, recording.SanitizedValue)
		if account != recording.SanitizedValue {
			err := recording.AddGeneralRegexSanitizer(recording.SanitizedValue, account, nil)
			if err != nil {
				panic(err)
			}
		} else if recording.GetRecordMode() != recording.PlaybackMode {
			panic("no value for " + v)
		}
	}
	err := recording.AddGeneralRegexSanitizer("batch_00000000-0000-0000-0000-000000000000", "batch_[0-9A-Fa-f]{8}[-]([0-9A-Fa-f]{4}[-]?){3}[0-9a-fA-F]{12}", nil)
	if err != nil {
		panic(err)
	}
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

func createClientForRecording(t *testing.T, tableName string, serviceURL string, tp tracing.Provider) (*Client, error) {
	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	tokenCredential, err := credential.New(nil)
	require.NoError(t, err)

	options := &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Retry: policy.RetryOptions{
				StatusCodes: statusCodesForRetry(),
			},
			TracingProvider: tp,
			Transport:       client,
		},
	}
	if !strings.HasSuffix(serviceURL, "/") && tableName != "" {
		serviceURL += "/"
	}
	serviceURL += tableName

	return NewClient(serviceURL, tokenCredential, options)
}

func createClientForRecordingForSharedKey(t *testing.T, tableName string, serviceURL string, cred SharedKeyCredential, tp tracing.Provider) (*Client, error) {
	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Retry: policy.RetryOptions{
				StatusCodes: statusCodesForRetry(),
			},
			TracingProvider: tp,
			Transport:       client,
		},
	}
	if !strings.HasSuffix(serviceURL, "/") && tableName != "" {
		serviceURL += "/"
	}
	serviceURL += tableName

	return NewClientWithSharedKey(serviceURL, &cred, options)
}

func createClientForRecordingWithNoCredential(t *testing.T, tableName string, serviceURL string, tp tracing.Provider) (*Client, error) {
	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Retry: policy.RetryOptions{
				StatusCodes: statusCodesForRetry(),
			},
			TracingProvider: tp,
			Transport:       client,
		},
	}
	if !strings.HasSuffix(serviceURL, "/") && tableName != "" {
		serviceURL += "/"
	}
	serviceURL += tableName

	return NewClientWithNoCredential(serviceURL, options)
}

func createServiceClientForRecording(t *testing.T, serviceURL string, tp tracing.Provider) (*ServiceClient, error) {
	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	tokenCredential, err := credential.New(nil)
	require.NoError(t, err)

	options := &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Retry: policy.RetryOptions{
				StatusCodes: statusCodesForRetry(),
			},
			TracingProvider: tp,
			Transport:       client,
		},
	}
	return NewServiceClient(serviceURL, tokenCredential, options)
}

func createServiceClientForRecordingForSharedKey(t *testing.T, serviceURL string, cred SharedKeyCredential, tp tracing.Provider) (*ServiceClient, error) {
	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Retry: policy.RetryOptions{
				StatusCodes: statusCodesForRetry(),
			},
			TracingProvider: tp,
			Transport:       client,
		},
	}
	return NewServiceClientWithSharedKey(serviceURL, &cred, options)
}

func createServiceClientForRecordingWithNoCredential(t *testing.T, serviceURL string, tp tracing.Provider) (*ServiceClient, error) {
	client, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Retry: policy.RetryOptions{
				StatusCodes: statusCodesForRetry(),
			},
			TracingProvider: tp,
			Transport:       client,
		},
	}
	return NewServiceClientWithNoCredential(serviceURL, options)
}

func initClientTest(t *testing.T, service endpointType, createTable bool, tp tracing.Provider) *Client {
	var client *Client
	var err error

	switch service {
	case storageEndpoint:
		client, err = createStorageClient(t, tp, &testClientOptions{UseSharedKey: true})
	case storageTokenCredentialEndpoint:
		client, err = createStorageClient(t, tp, &testClientOptions{UseSharedKey: false})
	case cosmosEndpoint:
		client, err = createCosmosClient(t, tp, &testClientOptions{UseSharedKey: true})
	case cosmosTokenCredentialEndpoint:
		client, err = createCosmosClient(t, tp, &testClientOptions{UseSharedKey: false})
	default:
		require.FailNowf(t, "Invalid client test option", "%s", string(service))
	}

	require.NoError(t, err)

	err = recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)

	if createTable {
		_, err = client.CreateTable(ctx, nil)
		require.NoError(t, err, "failed to create table %s", client.name)
	}

	t.Cleanup(func() {
		_, err = client.Delete(ctx, nil)
		require.NoError(t, err, "failed to delete table %s", client.name)
		err = recording.Stop(t, nil)
		require.NoError(t, err)
	})

	return client
}

func initServiceTest(t *testing.T, service endpointType, tp tracing.Provider) *ServiceClient {
	var client *ServiceClient
	var err error
	switch service {
	case storageEndpoint:
		client, err = createStorageServiceClient(t, tp, &testClientOptions{UseSharedKey: true})
	case storageTokenCredentialEndpoint:
		client, err = createStorageServiceClient(t, tp, &testClientOptions{UseSharedKey: false})
	case cosmosEndpoint:
		client, err = createCosmosServiceClient(t, tp, &testClientOptions{UseSharedKey: true})
	case cosmosTokenCredentialEndpoint:
		client, err = createCosmosServiceClient(t, tp, &testClientOptions{UseSharedKey: false})
	default:
		require.FailNowf(t, "Invalid service test option", "%s", string(service))
	}
	require.NoError(t, err)

	err = recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)

	t.Cleanup(func() {
		err = recording.Stop(t, nil)
		require.NoError(t, err)
	})

	return client
}

func getSharedKeyCredential() (*SharedKeyCredential, error) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		return NewSharedKeyCredential("accountName", "daaaaaaaaaabbbbbbbbbbcccccccccccccccccccdddddddddddddddddddeeeeeeeeeeefffffffffffggggg==")
	}

	accountName := recording.GetEnvVariable("TABLES_COSMOS_ACCOUNT_NAME", fakeAccount)
	accountKey := recording.GetEnvVariable("TABLES_PRIMARY_COSMOS_ACCOUNT_KEY", "fakeAccountKey")

	return NewSharedKeyCredential(accountName, accountKey)
}

func createStorageClient(t *testing.T, tp tracing.Provider, options *testClientOptions) (*Client, error) {
	if options == nil {
		options = &testClientOptions{}
	}

	var err error
	accountName := recording.GetEnvVariable("TABLES_STORAGE_ACCOUNT_NAME", fakeAccount)
	accountKey := recording.GetEnvVariable("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY", "fakeaccountkey")

	serviceURL := storageURI(accountName)

	tableName, err := createRandomName(t, tableNamePrefix)
	require.NoError(t, err)

	if options.UseSharedKey {
		var cred *SharedKeyCredential

		if recording.GetRecordMode() == recording.PlaybackMode {
			cred, err = getSharedKeyCredential()
			require.NoError(t, err)
		} else {
			cred, err = NewSharedKeyCredential(accountName, accountKey)
			require.NoError(t, err)
		}

		return createClientForRecordingForSharedKey(t, tableName, serviceURL, *cred, tp)
	}

	return createClientForRecording(t, tableName, serviceURL, tp)
}

type testClientOptions struct {
	UseSharedKey bool
}

func createCosmosClient(t *testing.T, tp tracing.Provider, options *testClientOptions) (*Client, error) {
	if options == nil {
		options = &testClientOptions{}
	}

	accountName := recording.GetEnvVariable("TABLES_COSMOS_ACCOUNT_NAME", fakeAccount)
	if recording.GetRecordMode() == recording.PlaybackMode {
		accountName = fakeAccount
	}

	serviceURL := cosmosURI(accountName)

	tableName, err := createRandomName(t, tableNamePrefix)
	require.NoError(t, err)

	if options.UseSharedKey {
		cred, err := getSharedKeyCredential()
		require.NoError(t, err)
		return createClientForRecordingForSharedKey(t, tableName, serviceURL, *cred, tp)
	}

	return createClientForRecording(t, tableName, serviceURL, tp)
}

func createStorageServiceClient(t *testing.T, tp tracing.Provider, options *testClientOptions) (*ServiceClient, error) {
	if options == nil {
		options = &testClientOptions{}
	}

	accountName := recording.GetEnvVariable("TABLES_STORAGE_ACCOUNT_NAME", fakeAccount)
	accountKey := recording.GetEnvVariable("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY", "fakeaccountkey")
	serviceURL := storageURI(accountName)

	if options.UseSharedKey {
		var cred *SharedKeyCredential
		var err error

		if recording.GetRecordMode() == recording.PlaybackMode {
			cred, err = getSharedKeyCredential()
			require.NoError(t, err)
		} else {
			cred, err = NewSharedKeyCredential(accountName, accountKey)
			require.NoError(t, err)
		}

		return createServiceClientForRecordingForSharedKey(t, serviceURL, *cred, tp)
	}

	return createServiceClientForRecording(t, serviceURL, tp)
}

func createCosmosServiceClient(t *testing.T, tp tracing.Provider, options *testClientOptions) (*ServiceClient, error) {
	if options == nil {
		options = &testClientOptions{}
	}

	accountName := recording.GetEnvVariable("TABLES_COSMOS_ACCOUNT_NAME", fakeAccount)

	if recording.GetRecordMode() == recording.PlaybackMode {
		accountName = fakeAccount
	}

	serviceURL := cosmosURI(accountName)

	if options.UseSharedKey {
		var cred *SharedKeyCredential

		cred, err := getSharedKeyCredential()
		require.NoError(t, err)

		return createServiceClientForRecordingForSharedKey(t, serviceURL, *cred, tp)
	}

	return createServiceClientForRecording(t, serviceURL, tp)
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

func statusCodesForRetry() []int {
	// we add 403 to the standard list of status
	// codes as we see transient live test failures
	// due to 403s
	return []int{
		http.StatusForbidden,           // 403
		http.StatusRequestTimeout,      // 408
		http.StatusTooManyRequests,     // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout,      // 504
	}
}
