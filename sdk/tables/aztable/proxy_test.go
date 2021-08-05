// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var AADAuthenticationScope = "https://storage.azure.com/.default"

var localCertFile = "C:/github.com/azure-sdk-tools/tools/test-proxy/docker/dev_certificate/dotnet-devcert.crt"

func getRootCas(filePath *string) (*x509.CertPool, error) {
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		rootCAs = x509.NewCertPool()
	}

	cert, err := ioutil.ReadFile(*filePath)
	if err != nil {
		fmt.Println("error opening cert file")
		return nil, err
	}

	if ok := rootCAs.AppendCertsFromPEM(cert); !ok {
		fmt.Println("No certs appended, using system certs only")
	}

	return rootCAs, nil
}

func getHTTPClient() (*http.Client, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig.InsecureSkipVerify = true

	rootCAs, err := getRootCas(&localCertFile)
	if err != nil {
		return nil, err
	}
	transport.TLSClientConfig.RootCAs = rootCAs
	transport.TLSClientConfig.MinVersion = tls.VersionTLS12
	defaultHttpClient := &http.Client{
		Transport: transport,
	}
	return defaultHttpClient, nil
}

func createTableClientForRecording(t *testing.T, tableName string, serviceURL string, cred azcore.Credential) (*TableClient, error) {
	policy := recording.NewRecordingPolicy(&recording.RecordingOptions{UseHTTPS: true})
	client, err := getHTTPClient()
	require.NoError(t, err)
	options := &TableClientOptions{
		Scopes:         []string{AADAuthenticationScope},
		PerCallOptions: []azcore.Policy{policy},
		HTTPClient:     client,
	}
	fmt.Println("service url: ", serviceURL)
	return NewTableClient(tableName, serviceURL, cred, options)
}

func createTableServiceClientForRecording(t *testing.T, serviceURL string, cred azcore.Credential) (*TableServiceClient, error) {
	policy := recording.NewRecordingPolicy(&recording.RecordingOptions{UseHTTPS: true})
	client, err := getHTTPClient()
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

func createStorageTableClient(t *testing.T) (*TableClient, error) {
	var cred azcore.Credential
	accountName := "fakestorageaccount"
	var ok bool
	if recording.InPlayback() {
		fmt.Println("IN PLAYBACK")
		cred = NewFakeCredential(accountName, "fakeAccountKey")
	} else {
		accountName, ok = os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
		if !ok {
			accountName = "fakestorageaccount"
		}

		err := recording.AddUriSanitizer("fakestorageaccount", accountName, nil)
		require.NoError(t, err)

		cred, err = azidentity.NewDefaultAzureCredential(nil)
		require.NoError(t, err)
	}
	fmt.Println("ACCOUNTNAME: ", accountName)

	serviceURL := storageURI(accountName, "core.windows.net")

	tableName, err := createRandomName(t, "tableName")
	require.NoError(t, err)

	return createTableClientForRecording(t, tableName, serviceURL, cred)
}

func createCosmosTableClient(t *testing.T) (*TableClient, error) {
	var cred azcore.Credential
	accountName := "fakestorageaccount"
	var ok bool
	if recording.InPlayback() {
		cred = NewFakeCredential(accountName, "fakeAccountKey")
	} else {
		accountName, ok = os.LookupEnv("TABLES_COSMOS_ACCOUNT_NAME")
		if !ok {
			accountName = "fakecosmosaccount"
		}
		accountKey, ok := os.LookupEnv("TABLES_PRIMARY_COSMOS_ACCOUNT_KEY")
		if !ok {
			t.Log("COSMOS KEY")
			accountKey = "fakekey"
		}

		err := recording.AddUriSanitizer("fakestorageaccount", accountName, nil)
		require.NoError(t, err)
		cred, err = NewSharedKeyCredential(accountName, accountKey)
		require.NoError(t, err)
	}
	fmt.Println("ACCOUNTNAME: ", accountName)

	serviceURL := cosmosURI(accountName, "cosmos.azure.com")

	tableName, err := createRandomName(t, "tableName")
	require.NoError(t, err)

	return createTableClientForRecording(t, tableName, serviceURL, cred)
}

func createStorageServiceClient(t *testing.T) (*TableServiceClient, error) {
	var cred azcore.Credential
	accountName := "fakestorageaccount"
	if recording.InPlayback() {
		cred = NewFakeCredential(accountName, "fakeAccountKey")
	} else {
		accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
		if !ok {
			accountName = "fakestorageaccount"
		}

		err := recording.AddUriSanitizer("fakestorageaccount", accountName, nil)
		require.NoError(t, err)

		cred, err = azidentity.NewDefaultAzureCredential(nil)
		require.NoError(t, err)
	}

	// accountName := recording.GetEnvVariable("TABLES_STORAGE_ACCOUNT_NAME", "fakestorageaccount")

	// err := recording.AddUriSanitizer("fakestorageaccount", accountName, nil)
	// require.NoError(t, err)

	serviceURL := storageURI(accountName, "core.windows.net")
	// cred, err := azidentity.NewDefaultAzureCredential(nil)
	// require.NoError(t, err)
	return createTableServiceClientForRecording(t, serviceURL, cred)
}

func createCosmosServiceClient(t *testing.T) (*TableServiceClient, error) {
	var cred azcore.Credential
	accountName := "fakestorageaccount"
	if recording.InPlayback() {
		cred = NewFakeCredential(accountName, "fakeAccountKey")
	} else {
		accountName, ok := os.LookupEnv("TABLES_COSMOS_ACCOUNT_NAME")
		if !ok {
			accountName = "fakecosmosaccount"
		}
		accountKey, ok := os.LookupEnv("TABLES_PRIMARY_COSMOS_ACCOUNT_KEY")
		if !ok {
			t.Log("COSMOS KEY")
			accountKey = "fakekey"
		}

		err := recording.AddUriSanitizer("fakestorageaccount", accountName, nil)
		require.NoError(t, err)
		cred, err = NewSharedKeyCredential(accountName, accountKey)
		require.NoError(t, err)
	}
	// accountName := recording.GetEnvVariable("TABLES_COSMOS_ACCOUNT_NAME", "fakestorageaccount")
	// accountKey := recording.GetEnvVariable("TABLES_PRIMARY_COSMOS_ACCOUNT_KEY", "fakekey")

	// err := recording.AddUriSanitizer("fakestorageaccount", accountName, nil)
	// require.NoError(t, err)

	serviceURL := cosmosURI(accountName, "cosmos.azure.com")
	// cred, err := createSharedKey(accountName, accountKey)
	// require.NoError(t, err)
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

func createAADCredential() (azcore.Credential, error) {
	if recording.GetRecordMode() == recording.ModeRecording {
		return azidentity.NewDefaultAzureCredential(nil)
	}
	return NewFakeCredential("accountName", "accountKey"), nil
}

func createSharedKey(accountName, accountKey string) (azcore.Credential, error) {
	if recording.GetRecordMode() == recording.ModeRecording {
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
		authHeader := strings.Join([]string{"Authorization ", f.accountName, ":", f.accountKey}, "")
		req.Request.Header.Set(azcore.HeaderAuthorization, authHeader)
		return req.Next()
	})
}
