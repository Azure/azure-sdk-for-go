package testcommon

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/service"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	DefaultEndpointSuffix       = "core.windows.net/"
	DefaultBlobEndpointSuffix   = "blob.core.windows.net/"
	AccountNameEnvVar           = "AZURE_STORAGE_ACCOUNT_NAME"
	AccountKeyEnvVar            = "AZURE_STORAGE_ACCOUNT_KEY"
	DefaultEndpointSuffixEnvVar = "AZURE_STORAGE_ENDPOINT_SUFFIX"
	SubscriptionID              = "SUBSCRIPTION_ID"
	ResourceGroupName           = "RESOURCE_GROUP_NAME"
)

const (
	FakeStorageAccount = "fakestorage"
	FakeBlobStorageURL = "https://fakestorage.blob.core.windows.net"
	FakeDFSStorageURL  = "https://fakestorage.dfs.core.windows.net"
	FakeToken          = "faketoken"
)

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

func GetFilesystemClient(fsName string, t *testing.T, accountType TestAccountType, options *filesystem.ClientOptions) (*filesystem.Client, error) {
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

func DeleteFilesystem(ctx context.Context, _require *require.Assertions, filesystemClient *filesystem.Client) {
	_, err := filesystemClient.Delete(ctx, nil)
	_require.Nil(err)
}
