//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// Contains common helpers for TESTS ONLY
package testcommon

import (
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/service"
	"github.com/stretchr/testify/require"
	"testing"
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

var BasicMetadata = map[string]string{"Foo": "bar"}

func setClientOptions(t *testing.T, opts *azcore.ClientOptions) {
	opts.Logging.AllowedHeaders = append(opts.Logging.AllowedHeaders, "X-Request-Mismatch", "X-Request-Mismatch-Error")

	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts.Transport = transport
}

func GetServiceClient(t *testing.T, accountType TestAccountType, options *service.ClientOptions) (*service.Client, error) {
	if options == nil {
		options = &service.ClientOptions{}
	}

	setClientOptions(t, &options.ClientOptions)

	cred, err := GetGenericCredential(accountType)
	if err != nil {
		return nil, err
	}

	serviceClient, err := service.NewClientWithSharedKeyCredential("https://"+cred.AccountName()+".queue.core.windows.net/", cred, options)

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

func GetServiceClientFromConnectionString(t *testing.T, accountType TestAccountType, options *service.ClientOptions) (*service.Client, error) {
	if options == nil {
		options = &service.ClientOptions{}
	}

	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	options.Transport = transport

	if recording.GetRecordMode() == recording.PlaybackMode {
		return service.NewClientWithNoCredential(FakeStorageURL, options)
	}

	connectionString := GetConnectionString(accountType)
	svcClient, err := service.NewClientFromConnectionString(connectionString, options)
	return svcClient, err
}

// TODO: GetQueueClient()
// TODO: CreateNewQueue()
// TODO: DeleteQueue()
