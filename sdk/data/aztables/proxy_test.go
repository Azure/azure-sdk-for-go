// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"fmt"
	"hash/fnv"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var pathToPackage = "sdk/data/aztables"

type recordingPolicy struct {
	options recording.RecordingOptions
	t       *testing.T
}

func (r recordingPolicy) Host() string {
	if r.options.UseHTTPS {
		return "localhost:5001"
	}
	return "localhost:5000"
}

func (r recordingPolicy) Scheme() string {
	if r.options.UseHTTPS {
		return "https"
	}
	return "http"
}

func NewRecordingPolicy(t *testing.T, o *recording.RecordingOptions) policy.Policy {
	if o == nil {
		o = &recording.RecordingOptions{UseHTTPS: true}
	}
	p := &recordingPolicy{options: *o, t: t}
	return p
}

func (p *recordingPolicy) Do(req *policy.Request) (resp *http.Response, err error) {
	if recording.GetRecordMode() != "live" {
		fmt.Println("BEFORE: ", req.Raw().URL)
		originalURLHost := req.Raw().URL.Host
		req.Raw().URL.Scheme = p.Scheme()
		req.Raw().URL.Host = p.Host()
		req.Raw().Host = p.Host()
		fmt.Println("AFTER: ", req.Raw().URL)

		req.Raw().Header.Set(recording.UpstreamURIHeader, fmt.Sprintf("%v://%v", p.Scheme(), originalURLHost))
		req.Raw().Header.Set(recording.ModeHeader, recording.GetRecordMode())
		req.Raw().Header.Set(recording.IDHeader, recording.GetRecordingId(p.t))
		fmt.Println(req.Raw().Header.Get(recording.UpstreamURIHeader))
	}
	return req.Next()
}

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
	p := NewRecordingPolicy(t, &recording.RecordingOptions{UseHTTPS: true})
	client, err := recording.GetHTTPClient(t)
	require.NoError(t, err)

	options := &ClientOptions{ClientOptions: azcore.ClientOptions{
		PerCallPolicies: []policy.Policy{p},
		Transport:       client,
	}}
	if !strings.HasSuffix(serviceURL, "/") && tableName != "" {
		serviceURL += "/"
	}
	serviceURL += tableName

	return NewClientWithSharedKey(serviceURL, &cred, options)
}

func createClientForRecordingWithNoCredential(t *testing.T, tableName string, serviceURL string) (*Client, error) {
	p := NewRecordingPolicy(t, &recording.RecordingOptions{UseHTTPS: true})
	client, err := recording.GetHTTPClient(t)
	require.NoError(t, err)

	options := &ClientOptions{ClientOptions: azcore.ClientOptions{
		PerCallPolicies: []policy.Policy{p},
		Transport:       client,
	}}
	if !strings.HasSuffix(serviceURL, "/") && tableName != "" {
		serviceURL += "/"
	}
	serviceURL += tableName

	return NewClientWithNoCredential(serviceURL, options)
}

func createServiceClientForRecording(t *testing.T, serviceURL string, cred SharedKeyCredential) (*ServiceClient, error) {
	p := NewRecordingPolicy(t, &recording.RecordingOptions{UseHTTPS: true})
	client, err := recording.GetHTTPClient(t)
	require.NoError(t, err)

	options := &ClientOptions{ClientOptions: azcore.ClientOptions{
		PerCallPolicies: []policy.Policy{p},
		Transport:       client,
	}}
	return NewServiceClientWithSharedKey(serviceURL, &cred, options)
}

func createServiceClientForRecordingWithNoCredential(t *testing.T, serviceURL string) (*ServiceClient, error) {
	p := NewRecordingPolicy(t, &recording.RecordingOptions{UseHTTPS: true})
	client, err := recording.GetHTTPClient(t)
	require.NoError(t, err)

	options := &ClientOptions{ClientOptions: azcore.ClientOptions{
		PerCallPolicies: []policy.Policy{p},
		Transport:       client,
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
		_, err = client.Create(context.Background(), nil)
		require.NoError(t, err)
	}

	return client, func() {
		_, err = client.Delete(context.Background(), nil)
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

func getAADCredential(t *testing.T) (azcore.TokenCredential, error) { //nolint
	if recording.GetRecordMode() == "playback" {
		cred := NewFakeCredential("fakestorageaccount", "fakeAccountKey")
		return cred, nil
	}

	return azidentity.NewDefaultAzureCredential(nil)
}

func getSharedKeyCredential(t *testing.T) (*SharedKeyCredential, error) {
	if recording.GetRecordMode() == "playback" {
		return NewSharedKeyCredential("accountName", "daaaaaaaaaabbbbbbbbbbcccccccccccccccccccdddddddddddddddddddeeeeeeeeeeefffffffffffggggg==")
	}

	accountName := recording.GetEnvVariable("TABLES_COSMOS_ACCOUNT_NAME", "fakestorageaccount")
	accountKey := recording.GetEnvVariable("TABLES_PRIMARY_COSMOS_ACCOUNT_KEY", "fakeAccountKey")

	return NewSharedKeyCredential(accountName, accountKey)
}

func createStorageClient(t *testing.T) (*Client, error) {
	var cred *SharedKeyCredential
	var err error
	accountName := recording.GetEnvVariable("TABLES_STORAGE_ACCOUNT_NAME", "fakestorageaccount")
	accountKey := recording.GetEnvVariable("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY", "fakestorageaccountkey")

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
	accountName := recording.GetEnvVariable("TABLES_COSMOS_ACCOUNT_NAME", "fakestorageaccount")
	if recording.GetRecordMode() == "playback" {
		accountName = "fakestorageaccount"
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
	accountName := recording.GetEnvVariable("TABLES_STORAGE_ACCOUNT_NAME", "fakestorageaccount")
	accountKey := recording.GetEnvVariable("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY", "fakestorageaccountkey")

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
	accountName := recording.GetEnvVariable("TABLES_COSMOS_ACCOUNT_NAME", "fakestorageaccount")
	if recording.GetRecordMode() == "playback" {
		accountName = "fakestorageaccount"
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
	pager := service.ListTables(nil)
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		for _, v := range resp.Tables {
			_, err := service.DeleteTable(ctx, *v.TableName, nil)
			if err != nil {
				return err
			}
		}
	}
	return pager.Err()
}
