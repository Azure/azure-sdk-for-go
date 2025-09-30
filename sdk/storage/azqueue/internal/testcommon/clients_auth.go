//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// Contains common helpers for TESTS ONLY

package testcommon

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/v2"
	"github.com/stretchr/testify/require"
)

type TestAccountType string

const (
	TestAccountDefault   TestAccountType = ""
	TestAccountSecondary TestAccountType = "SECONDARY_"
	TestAccountPremium   TestAccountType = "PREMIUM_"
)

const (
	DefaultEndpointSuffix       = "core.windows.net/"
	DefaultQueueEndpointSuffix  = "queue.core.windows.net/"
	AccountNameEnvVar           = "AZURE_STORAGE_ACCOUNT_NAME"
	AccountKeyEnvVar            = "AZURE_STORAGE_ACCOUNT_KEY"
	DefaultEndpointSuffixEnvVar = "AZURE_STORAGE_ENDPOINT_SUFFIX"
)

const (
	FakeStorageAccount = "fakestorage"
	FakeStorageURL     = "https://fakestorage.queue.core.windows.net"
)

var BasicMetadata = map[string]*string{"Foo": to.Ptr("bar")}

func setClientOptions(t *testing.T, opts *azcore.ClientOptions) {
	opts.Logging.AllowedHeaders = append(opts.Logging.AllowedHeaders, "X-Request-Mismatch", "X-Request-Mismatch-Error")

	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts.Transport = transport
}

func GetServiceClient(t *testing.T, accountType TestAccountType, options *azqueue.ClientOptions) (*azqueue.ServiceClient, error) {
	if options == nil {
		options = &azqueue.ClientOptions{}
	}

	setClientOptions(t, &options.ClientOptions)

	cred, err := GetGenericCredential(accountType)
	if err != nil {
		return nil, err
	}

	serviceClient, err := azqueue.NewServiceClientWithSharedKeyCredential("https://"+cred.AccountName()+".queue.core.windows.net/", cred, options)

	return serviceClient, err
}

func GetAccountInfo(accountType TestAccountType) (string, string) {
	accountNameEnvVar := string(accountType) + AccountNameEnvVar
	accountKeyEnvVar := string(accountType) + AccountKeyEnvVar
	accountName, _ := GetRequiredEnv(accountNameEnvVar)
	accountKey, _ := GetRequiredEnv(accountKeyEnvVar)
	return accountName, accountKey
}

func GetGenericCredential(accountType TestAccountType) (*azqueue.SharedKeyCredential, error) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		return azqueue.NewSharedKeyCredential(FakeStorageAccount, "ZmFrZQ==")
	}

	accountName, accountKey := GetAccountInfo(accountType)
	if accountName == "" || accountKey == "" {
		return nil, errors.New(string(accountType) + AccountNameEnvVar + " and/or " + string(accountType) + AccountKeyEnvVar + " environment variables not specified.")
	}
	return azqueue.NewSharedKeyCredential(accountName, accountKey)
}

func GetConnectionString(accountType TestAccountType) string {
	accountName, accountKey := GetAccountInfo(accountType)
	connectionString := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=core.windows.net/",
		accountName, accountKey)
	return connectionString
}

func GetServiceClientFromConnectionString(t *testing.T, accountType TestAccountType, options *azqueue.ClientOptions) (*azqueue.ServiceClient, error) {
	if options == nil {
		options = &azqueue.ClientOptions{}
	}

	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	options.Transport = transport

	if recording.GetRecordMode() == recording.PlaybackMode {
		return azqueue.NewServiceClientWithNoCredential(FakeStorageURL, options)
	}

	connectionString := GetConnectionString(accountType)
	svcClient, err := azqueue.NewServiceClientFromConnectionString(connectionString, options)
	return svcClient, err
}

func GetQueueClient(queueName string, serviceClient *azqueue.ServiceClient) *azqueue.QueueClient {
	return serviceClient.NewQueueClient(queueName)
}

func CreateNewQueue(ctx context.Context, _require *require.Assertions, queueName string, serviceClient *azqueue.ServiceClient) *azqueue.QueueClient {
	queueClient := GetQueueClient(queueName, serviceClient)

	_, err := queueClient.Create(ctx, nil)
	_require.NoError(err)
	return queueClient
}

func DeleteQueue(ctx context.Context, _require *require.Assertions, queueClient *azqueue.QueueClient) {
	_, err := queueClient.Delete(ctx, nil)
	_require.NoError(err)
}
