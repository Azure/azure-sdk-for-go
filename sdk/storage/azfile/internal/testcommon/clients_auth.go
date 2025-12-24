// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// Contains common helpers for TESTS ONLY
package testcommon

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

type TestAccountType string

const (
	TestAccountDefault    TestAccountType = ""
	TestAccountSecondary  TestAccountType = "SECONDARY_"
	TestAccountPremium    TestAccountType = "PREMIUM_"
	TestAccountSoftDelete TestAccountType = "SOFT_DELETE_"
)

const (
	DefaultEndpointSuffix       = "core.windows.net/"
	DefaultFileEndpointSuffix   = "file.core.windows.net/"
	AccountNameEnvVar           = "AZURE_STORAGE_ACCOUNT_NAME"
	AccountKeyEnvVar            = "AZURE_STORAGE_ACCOUNT_KEY"
	DefaultEndpointSuffixEnvVar = "AZURE_STORAGE_ENDPOINT_SUFFIX"
	EncryptionScopeEnvVar       = "AZURE_STORAGE_ENCRYPTION_SCOPE"
	PremiumAccountNameEnvVar    = "FILE_STORAGE_ACCOUNT_NAME"
	PremiumAccountKeyEnvVar     = "FILE_STORAGE_ACCOUNT_KEY"
)

const (
	FakeStorageAccount = "fakestorage"
	FakeStorageURL     = "https://fakestorage.file.core.windows.net"
	FakeToken          = "faketoken"
)

const (
	ISO8601                  = "2006-01-02T15:04:05.0000000Z07:00"
	FilePermissionFormatSddl = "sddl"
	FilePermissionBinary     = "Binary"
)

var (
	SampleSDDL   = `O:S-1-5-32-548G:S-1-5-21-397955417-626881126-188441444-512D:(A;;RPWPCCDCLCSWRCWDWOGA;;;S-1-0-0)`
	SampleBinary = `AQAUhGwAAACIAAAAAAAAABQAAAACAFgAAwAAAAAAFAD/AR8AAQEAAAAAAAUSAAAAAAAYAP8BHwABAgAAAAAABSAAAAAgAgAAAAAkAKkAEgABBQAAAAAABRUAAABZUbgXZnJdJWRjOwuMmS4AAQUAAAAAAAUVAAAAoGXPfnhLm1/nfIdwr/1IAQEFAAAAAAAFFQAAAKBlz354S5tf53yHcAECAAA=`
)

var BasicMetadata = map[string]*string{
	"foo": to.Ptr("foovalue"),
	"bar": to.Ptr("barvalue"),
}

func SetClientOptions(t *testing.T, opts *azcore.ClientOptions) {
	opts.Logging.AllowedHeaders = append(opts.Logging.AllowedHeaders, "X-Request-Mismatch", "X-Request-Mismatch-Error")

	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts.Transport = transport
}

func GetServiceClient(t *testing.T, accountType TestAccountType, options *service.ClientOptions) (*service.Client, error) {
	if options == nil {
		options = &service.ClientOptions{}
	}

	SetClientOptions(t, &options.ClientOptions)

	cred, err := GetGenericSharedKeyCredential(accountType)
	if err != nil {
		return nil, err
	}

	serviceClient, err := service.NewClientWithSharedKeyCredential("https://"+cred.AccountName()+".file.core.windows.net/", cred, options)

	return serviceClient, err
}

func GetServiceClientNoCredential(t *testing.T, sasUrl string, options *service.ClientOptions) (*service.Client, error) {
	if options == nil {
		options = &service.ClientOptions{}
	}

	SetClientOptions(t, &options.ClientOptions)

	serviceClient, err := service.NewClientWithNoCredential(sasUrl, options)

	return serviceClient, err
}

func GetGenericAccountInfo(accountType TestAccountType) (string, string) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		return FakeStorageAccount, "ZmFrZQ=="
	}
	accountNameEnvVar := string(accountType) + AccountNameEnvVar
	accountKeyEnvVar := string(accountType) + AccountKeyEnvVar
	if accountType == TestAccountPremium {
		accountNameEnvVar = string(accountType) + PremiumAccountNameEnvVar
		accountKeyEnvVar = string(accountType) + PremiumAccountKeyEnvVar
	}
	accountName, _ := GetRequiredEnv(accountNameEnvVar)
	accountKey, _ := GetRequiredEnv(accountKeyEnvVar)
	return accountName, accountKey
}

func GetGenericSharedKeyCredential(accountType TestAccountType) (*service.SharedKeyCredential, error) {
	accountName, accountKey := GetGenericAccountInfo(accountType)
	if accountName == "" || accountKey == "" {
		if accountType == TestAccountPremium {
			return nil, errors.New(string(accountType) + PremiumAccountNameEnvVar + " and/or " + string(accountType) + PremiumAccountKeyEnvVar + " environment variables not specified.")
		} else {
			return nil, errors.New(string(accountType) + AccountNameEnvVar + " and/or " + string(accountType) + AccountKeyEnvVar + " environment variables not specified.")
		}
	}
	return service.NewSharedKeyCredential(accountName, accountKey)
}

func GetGenericConnectionString(accountType TestAccountType) (*string, error) {
	accountName, accountKey := GetGenericAccountInfo(accountType)
	if accountName == "" || accountKey == "" {
		if accountType == TestAccountPremium {
			return nil, errors.New(string(accountType) + PremiumAccountNameEnvVar + " and/or " + string(accountType) + PremiumAccountKeyEnvVar + " environment variables not specified.")
		} else {
			return nil, errors.New(string(accountType) + AccountNameEnvVar + " and/or " + string(accountType) + AccountKeyEnvVar + " environment variables not specified.")
		}
	}
	connectionString := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=core.windows.net/",
		accountName, accountKey)
	return &connectionString, nil
}

func GetGenericTokenCredential() (azcore.TokenCredential, error) {
	return credential.New(nil)
}

func GetServiceClientFromConnectionString(t *testing.T, accountType TestAccountType, options *service.ClientOptions) (*service.Client, error) {
	if options == nil {
		options = &service.ClientOptions{}
	}
	SetClientOptions(t, &options.ClientOptions)

	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	options.Transport = transport

	cred, err := GetGenericConnectionString(accountType)
	if err != nil {
		return nil, err
	}
	svcClient, err := service.NewClientFromConnectionString(*cred, options)
	return svcClient, err
}

func GetShareClient(shareName string, s *service.Client) *share.Client {
	return s.NewShareClient(shareName)
}

func CreateNewShare(ctx context.Context, _require *require.Assertions, shareName string, svcClient *service.Client) *share.Client {
	shareClient := GetShareClient(shareName, svcClient)
	_, err := shareClient.Create(ctx, nil)
	_require.NoError(err)
	return shareClient
}

func DeleteShare(ctx context.Context, _require *require.Assertions, shareClient *share.Client) {
	_, err := shareClient.Delete(ctx, nil)
	_require.NoError(err)
}

func GetDirectoryClient(dirName string, s *share.Client) *directory.Client {
	return s.NewDirectoryClient(dirName)
}

func CreateNewDirectory(ctx context.Context, _require *require.Assertions, dirName string, shareClient *share.Client) *directory.Client {
	dirClient := GetDirectoryClient(dirName, shareClient)
	_, err := dirClient.Create(ctx, nil)
	_require.NoError(err)
	return dirClient
}

func DeleteDirectory(ctx context.Context, _require *require.Assertions, dirClient *directory.Client) {
	_, err := dirClient.Delete(ctx, nil)
	_require.NoError(err)
}

func GetFileClientFromShare(fileName string, shareClient *share.Client) *file.Client {
	return shareClient.NewRootDirectoryClient().NewFileClient(fileName)
}

func CreateNewFileFromShare(ctx context.Context, _require *require.Assertions, fileName string, fileSize int64, shareClient *share.Client) *file.Client {
	fClient := GetFileClientFromShare(fileName, shareClient)

	_, err := fClient.Create(ctx, fileSize, nil)
	_require.NoError(err)

	return fClient
}

func CreateNewFileFromShareWithData(ctx context.Context, _require *require.Assertions, fileName string, shareClient *share.Client) *file.Client {
	fClient := GetFileClientFromShare(fileName, shareClient)

	_, err := fClient.Create(ctx, int64(len(FileDefaultData)), nil)
	_require.NoError(err)

	_, err = fClient.UploadRange(ctx, 0, streaming.NopCloser(strings.NewReader(FileDefaultData)), nil)
	_require.NoError(err)

	return fClient
}

func DeleteFile(ctx context.Context, _require *require.Assertions, fileClient *file.Client) {
	_, err := fileClient.Delete(ctx, nil)
	_require.NoError(err)
}
