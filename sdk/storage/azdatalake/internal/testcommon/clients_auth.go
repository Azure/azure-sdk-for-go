//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package testcommon

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/service"
	"github.com/stretchr/testify/require"
)

const (
	DefaultEndpointSuffix         = "core.windows.net/"
	DefaultBlobEndpointSuffix     = "blob.core.windows.net/"
	AccountNameEnvVar             = "AZURE_STORAGE_ACCOUNT_NAME"
	AccountKeyEnvVar              = "AZURE_STORAGE_ACCOUNT_KEY"
	DefaultEndpointSuffixEnvVar   = "AZURE_STORAGE_ENDPOINT_SUFFIX"
	DataLakeEncryptionScopeEnvVar = "DATALAKE_AZURE_STORAGE_ENCRYPTION_SCOPE"
	SubscriptionID                = "SUBSCRIPTION_ID"
	ResourceGroupName             = "RESOURCE_GROUP_NAME"
)

const (
	FakeStorageAccount = "fakestorage"
	FakeBlobStorageURL = "https://fakestorage.blob.core.windows.net"
	FakeDFSStorageURL  = "https://fakestorage.dfs.core.windows.net"
	FakeToken          = "faketoken"
)

var BasicMetadata = map[string]*string{"Foo": to.Ptr("bar")}

var (
	DatalakeContentType        = "my_type"
	DatalakeContentDisposition = "my_disposition"
	DatalakeCacheControl       = "control"
	DatalakeContentLanguage    = "my_language"
	DatalakeContentEncoding    = "my_encoding"
)

var (
	testEncryptedKey        = "MDEyMzQ1NjcwMTIzNDU2NzAxMjM0NTY3MDEyMzQ1Njc="
	testEncryptedHash       = "3QFFFpRA5+XANHqwwbT4yXDmrT/2JaLt/FKHjzhOdoE="
	testEncryptionAlgorithm = file.EncryptionAlgorithmTypeAES256
	TestEncryptionContext   = "test_encryption_context"
	TestCPKByValue          = file.CPKInfo{
		EncryptionKey:       &testEncryptedKey,
		EncryptionKeySHA256: &testEncryptedHash,
		EncryptionAlgorithm: &testEncryptionAlgorithm,
	}
	TestEncryptionScope = "datalaketestencryptionscope"
	TestCPKScopeInfo    = container.CPKScopeInfo{
		DefaultEncryptionScope:         &TestEncryptionScope,
		PreventEncryptionScopeOverride: to.Ptr(false),
	}
)

var BasicHeaders = file.HTTPHeaders{
	ContentType:        &DatalakeContentType,
	ContentDisposition: &DatalakeContentDisposition,
	CacheControl:       &DatalakeCacheControl,
	ContentMD5:         nil,
	ContentLanguage:    &DatalakeContentLanguage,
	ContentEncoding:    &DatalakeContentEncoding,
}

type TestAccountType string

const (
	TestAccountDefault    TestAccountType = ""
	TestAccountSecondary  TestAccountType = "SECONDARY_"
	TestAccountPremium    TestAccountType = "PREMIUM_"
	TestAccountSoftDelete TestAccountType = "SOFT_DELETE_"
	TestAccountDatalake   TestAccountType = "DATALAKE_"
	TestAccountImmutable  TestAccountType = "IMMUTABLE_"
)

func SetClientOptions(t *testing.T, opts *azcore.ClientOptions) {
	opts.Logging.AllowedHeaders = append(opts.Logging.AllowedHeaders, "X-Request-Mismatch", "X-Request-Mismatch-Error")

	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts.Transport = transport
}

func GetGenericAccountInfo(accountType TestAccountType) (string, string) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		return FakeStorageAccount, "ZmFrZQ=="
	}
	accountNameEnvVar := string(accountType) + AccountNameEnvVar
	accountKeyEnvVar := string(accountType) + AccountKeyEnvVar
	accountName, _ := GetRequiredEnv(accountNameEnvVar)
	accountKey, _ := GetRequiredEnv(accountKeyEnvVar)
	return accountName, accountKey
}

func GetGenericSharedKeyCredential(accountType TestAccountType) (*azdatalake.SharedKeyCredential, error) {
	accountName, accountKey := GetGenericAccountInfo(accountType)
	if accountName == "" || accountKey == "" {
		return nil, errors.New(string(accountType) + AccountNameEnvVar + " and/or " + string(accountType) + AccountKeyEnvVar + " environment variables not specified.")
	}
	return azdatalake.NewSharedKeyCredential(accountName, accountKey)
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

	serviceClient, err := service.NewClientWithSharedKeyCredential("https://"+cred.AccountName()+".dfs.core.windows.net/", cred, options)

	return serviceClient, err
}

func GetFileSystemClient(fsName string, t *testing.T, accountType TestAccountType, options *filesystem.ClientOptions) (*filesystem.Client, error) {
	if options == nil {
		options = &filesystem.ClientOptions{}
	}

	SetClientOptions(t, &options.ClientOptions)

	cred, err := GetGenericSharedKeyCredential(accountType)
	if err != nil {
		return nil, err
	}

	filesystemClient, err := filesystem.NewClientWithSharedKeyCredential("https://"+cred.AccountName()+".dfs.core.windows.net/"+fsName, cred, options)

	return filesystemClient, err
}

func GetFileClient(fsName, fName string, t *testing.T, accountType TestAccountType, options *file.ClientOptions) (*file.Client, error) {
	if options == nil {
		options = &file.ClientOptions{}
	}

	SetClientOptions(t, &options.ClientOptions)

	cred, err := GetGenericSharedKeyCredential(accountType)
	if err != nil {
		return nil, err
	}

	fileClient, err := file.NewClientWithSharedKeyCredential("https://"+cred.AccountName()+".dfs.core.windows.net/"+fsName+"/"+fName, cred, options)

	return fileClient, err
}

func CreateNewFile(ctx context.Context, _require *require.Assertions, fileName string, filesystemClient *filesystem.Client) *file.Client {
	fileClient := filesystemClient.NewFileClient(fileName)
	_, err := fileClient.Create(ctx, nil)
	_require.NoError(err)
	return fileClient
}

func CreateNewDir(ctx context.Context, _require *require.Assertions, dirName string, filesystemClient *filesystem.Client) *directory.Client {
	dirClient := filesystemClient.NewDirectoryClient(dirName)
	_, err := dirClient.Create(ctx, nil)
	_require.NoError(err)
	return dirClient
}

func GetDirClient(fsName, dirName string, t *testing.T, accountType TestAccountType, options *directory.ClientOptions) (*directory.Client, error) {
	if options == nil {
		options = &directory.ClientOptions{}
	}

	SetClientOptions(t, &options.ClientOptions)

	cred, err := GetGenericSharedKeyCredential(accountType)
	if err != nil {
		return nil, err
	}

	dirClient, err := directory.NewClientWithSharedKeyCredential("https://"+cred.AccountName()+".dfs.core.windows.net/"+fsName+"/"+dirName, cred, options)

	return dirClient, err
}

func ServiceGetFileSystemClient(filesystemName string, s *service.Client) *filesystem.Client {
	return s.NewFileSystemClient(filesystemName)
}

func DeleteFileSystem(ctx context.Context, _require *require.Assertions, filesystemClient *filesystem.Client) {
	_, err := filesystemClient.Delete(ctx, nil)
	_require.NoError(err)
}

func DeleteFile(ctx context.Context, _require *require.Assertions, fileClient *file.Client) {
	_, err := fileClient.Delete(ctx, nil)
	_require.NoError(err)
}

func DeleteDir(ctx context.Context, _require *require.Assertions, dirClient *directory.Client) {
	_, err := dirClient.Delete(ctx, nil)
	_require.NoError(err)
}

func GetGenericConnectionString(accountType TestAccountType) (*string, error) {
	accountName, accountKey := GetGenericAccountInfo(accountType)
	if accountName == "" || accountKey == "" {
		return nil, errors.New(string(accountType) + AccountNameEnvVar + " and/or " + string(accountType) + AccountKeyEnvVar + " environment variables not specified.")
	}
	connectionString := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=core.windows.net/",
		accountName, accountKey)
	return &connectionString, nil
}

func CreateNewFileSystem(ctx context.Context, _require *require.Assertions, filesystemName string, serviceClient *service.Client) *filesystem.Client {
	fsClient := ServiceGetFileSystemClient(filesystemName, serviceClient)

	_, err := fsClient.Create(ctx, nil)
	_require.NoError(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)
	return fsClient
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

func GetServiceClientNoCredential(t *testing.T, sasUrl string, options *service.ClientOptions) (*service.Client, error) {
	if options == nil {
		options = &service.ClientOptions{}
	}

	SetClientOptions(t, &options.ClientOptions)

	serviceClient, err := service.NewClientWithNoCredential(sasUrl, options)

	return serviceClient, err
}

func GetGenericTokenCredential() (azcore.TokenCredential, error) {
	return credential.New(nil)
}

func GetUserDelegationSAS(svcClient *service.Client, filePath string, permissions sas.FilePermissions) (string, error) {
	// Set current and past time and create key
	now := time.Now().UTC().Add(-10 * time.Second)
	expiry := now.Add(2 * time.Hour)
	info := service.KeyInfo{
		Start:  to.Ptr(now.UTC().Format(sas.TimeFormat)),
		Expiry: to.Ptr(expiry.UTC().Format(sas.TimeFormat)),
	}

	udc, err := svcClient.GetUserDelegationCredential(context.Background(), info, nil)
	if err != nil {
		return "", err
	}

	// Create Blob Signature Values with desired permissions and sign with user delegation credential
	sasQueryParams, err := sas.DatalakeSignatureValues{
		Protocol:    sas.ProtocolHTTPS,
		StartTime:   time.Now().UTC().Add(time.Second * -10),
		ExpiryTime:  time.Now().UTC().Add(15 * time.Minute),
		Permissions: permissions.String(),
		FilePath:    filePath,
	}.SignWithUserDelegation(udc)
	if err != nil {
		return "", err
	}

	return sasQueryParams.Encode(), nil
}
