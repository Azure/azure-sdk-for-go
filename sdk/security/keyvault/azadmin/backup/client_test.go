//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package backup_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/backup"
	"github.com/stretchr/testify/require"
)

const fakeHsmURL = "https://fakehsm.managedhsm.azure.net/"

var (
	credential azcore.TokenCredential
	hsmURL     string
)

func TestMain(m *testing.M) {
	if recording.GetRecordMode() != recording.PlaybackMode {
		hsmURL = os.Getenv("AZURE_MANAGEDHSM_URL")
	}
	if hsmURL == "" {
		if recording.GetRecordMode() != recording.PlaybackMode {
			panic("no value for AZURE_MANAGEDHSM_URL")
		}
		hsmURL = fakeHsmURL
	}

	err := recording.ResetProxy(nil)
	if err != nil {
		panic(err)
	}
	if recording.GetRecordMode() == recording.PlaybackMode {
		credential = &FakeCredential{}
	} else {
		tenantID := lookupEnvVar("KEYVAULT_TENANT_ID")
		clientID := lookupEnvVar("KEYVAULT_CLIENT_ID")
		secret := lookupEnvVar("KEYVAULT_CLIENT_SECRET")
		credential, err = azidentity.NewClientSecretCredential(tenantID, clientID, secret, nil)
		if err != nil {
			panic(err)
		}
	}
	if recording.GetRecordMode() == recording.RecordingMode {
		err := recording.AddGeneralRegexSanitizer(fakeHsmURL, hsmURL, nil)
		if err != nil {
			panic(err)
		}
		defer func() {
			err := recording.ResetProxy(nil)
			if err != nil {
				panic(err)
			}
		}()
	}
	code := m.Run()
	os.Exit(code)
}

func startRecording(t *testing.T) {
	err := recording.Start(t, "sdk/security/keyvault/azadmin/testdata", nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
}

func startBackupTest(t *testing.T) (*backup.Client, backup.SASTokenParameter) {
	startRecording(t)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &backup.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	client, err := backup.NewClient(hsmURL, credential, opts)
	require.NoError(t, err)

	storageResourceUri := "https://" + os.Getenv("BLOB_STORAGE_ACCOUNT_NAME") + ".blob." + os.Getenv("KEYVAULT_STORAGE_ENDPOINT_SUFFIX") + "/" + os.Getenv("BLOB_CONTAINER_NAME")
	token := os.Getenv("BLOB_STORAGE_SAS_TOKEN")
	sasToken := backup.SASTokenParameter{
		StorageResourceURI: &storageResourceUri,
		Token:              &token,
	}

	return client, sasToken
}

func lookupEnvVar(s string) string {
	ret, ok := os.LookupEnv(s)
	if !ok {
		panic(fmt.Sprintf("Could not find env var: '%s'", s))
	}
	return ret
}

type FakeCredential struct{}

func (f *FakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "faketoken", ExpiresOn: time.Now().Add(time.Hour).UTC()}, nil
}

type serdeModel interface {
	json.Marshaler
	json.Unmarshaler
}

func testSerde[T serdeModel](t *testing.T, model T) {
	data, err := model.MarshalJSON()
	require.NoError(t, err)
	err = model.UnmarshalJSON(data)
	require.NoError(t, err)

	// testing unmarshal error scenarios
	var data2 []byte
	err = model.UnmarshalJSON(data2)
	require.Error(t, err)

	m := regexp.MustCompile(":.*$")
	modifiedData := m.ReplaceAllString(string(data), ":false}")
	if modifiedData != "{}" {
		data3 := []byte(modifiedData)
		err = model.UnmarshalJSON(data3)
		require.Error(t, err)
	}
}

func TestBackupRestore(t *testing.T) {
	client, sasToken := startBackupTest(t)

	testSerde(t, &sasToken)

	// backup the vault
	backupPoller, err := client.BeginFullBackup(context.Background(), sasToken, nil)
	require.NoError(t, err)
	backupResults, err := backupPoller.PollUntilDone(context.Background(), nil)
	require.NoError(t, err)
	require.Nil(t, backupResults.Error)
	require.Equal(t, "Succeeded", *backupResults.Status)
	require.Contains(t, *backupResults.AzureStorageBlobContainerURI, os.Getenv("BLOB_RESOURCE_URI"))
	testSerde(t, &backupResults)

	// restore the backup
	s := *backupResults.AzureStorageBlobContainerURI
	folderURI := s[strings.LastIndex(s, "/")+1:]
	restoreOperationParameters := backup.RestoreOperationParameters{
		FolderToRestore:    &folderURI,
		SasTokenParameters: &sasToken,
	}
	testSerde(t, &restoreOperationParameters)
	restorePoller, err := client.BeginFullRestore(context.Background(), restoreOperationParameters, nil)
	require.NoError(t, err)
	restoreResults, err := restorePoller.PollUntilDone(context.Background(), nil)
	require.NoError(t, err)
	require.Nil(t, restoreResults.Error)
	require.Equal(t, "Succeeded", *restoreResults.Status)
	require.NotNil(t, restoreResults.StartTime)
	require.NotNil(t, restoreResults.EndTime)
	require.NotNil(t, restoreResults.JobID)
	testSerde(t, &restoreResults)

	// additional waiting to avoid conflicts with resources in other tests
	if recording.GetRecordMode() != recording.PlaybackMode {
		time.Sleep(60 * time.Second)
	}
}

func TestBackupRestoreWithResumeToken(t *testing.T) {
	client, sasToken := startBackupTest(t)

	// backup the vault
	backupPoller, err := client.BeginFullBackup(context.Background(), sasToken, nil)
	require.NoError(t, err)

	// create a new poller from a continuation token
	token, err := backupPoller.ResumeToken()
	require.NoError(t, err)
	newBackupPoller, err := client.BeginFullBackup(context.Background(), sasToken, &backup.ClientBeginFullBackupOptions{ResumeToken: token})
	require.NoError(t, err)
	backupResults, err := newBackupPoller.PollUntilDone(context.Background(), nil)
	require.NoError(t, err)
	require.Nil(t, backupResults.Error)
	require.Equal(t, "Succeeded", *backupResults.Status)
	require.Contains(t, *backupResults.AzureStorageBlobContainerURI, os.Getenv("BLOB_RESOURCE_URI"))
	testSerde(t, &backupResults)

	// restore the backup
	s := *backupResults.AzureStorageBlobContainerURI
	folderURI := s[strings.LastIndex(s, "/")+1:]
	restoreOperationParameters := backup.RestoreOperationParameters{
		FolderToRestore:    &folderURI,
		SasTokenParameters: &sasToken,
	}
	testSerde(t, &restoreOperationParameters)
	restorePoller, err := client.BeginFullRestore(context.Background(), restoreOperationParameters, nil)
	require.NoError(t, err)

	// create a new poller from a continuation token
	restoreToken, err := restorePoller.ResumeToken()
	require.NoError(t, err)
	newRestorePoller, err := client.BeginFullRestore(context.Background(), restoreOperationParameters, &backup.ClientBeginFullRestoreOptions{ResumeToken: restoreToken})
	require.NoError(t, err)
	restoreResults, err := newRestorePoller.PollUntilDone(context.Background(), nil)
	require.NoError(t, err)
	require.Nil(t, restoreResults.Error)
	require.Equal(t, "Succeeded", *restoreResults.Status)
	require.NotNil(t, restoreResults.StartTime)
	require.NotNil(t, restoreResults.EndTime)
	require.NotNil(t, restoreResults.JobID)
	testSerde(t, &restoreResults)

	// additional waiting to avoid conflicts with resources in other tests
	if recording.GetRecordMode() != recording.PlaybackMode {
		time.Sleep(60 * time.Second)
	}
}

func TestBeginSelectiveKeyRestoreOperation(t *testing.T) {
	backupClient, sasToken := startBackupTest(t)

	// create a key to selectively restore
	cred := credential
	keyClient, err := azkeys.NewClient(hsmURL, cred, nil)
	require.NoError(t, err)
	params := azkeys.CreateKeyParameters{
		KeySize: to.Ptr(int32(2048)),
		Kty:     to.Ptr(azkeys.JSONWebKeyTypeRSA),
	}
	_, err = keyClient.CreateKey(context.TODO(), "selective-restore-test-key", params, nil)
	require.NoError(t, err)

	// backup the vault
	backupPoller, err := backupClient.BeginFullBackup(context.Background(), sasToken, nil)
	require.NoError(t, err)
	backupResults, err := backupPoller.PollUntilDone(context.Background(), nil)
	require.NoError(t, err)

	// restore the key
	s := *backupResults.AzureStorageBlobContainerURI
	folderURI := s[strings.LastIndex(s, "/")+1:]
	restoreOperationParameters := backup.SelectiveKeyRestoreOperationParameters{
		Folder:             &folderURI,
		SasTokenParameters: &sasToken,
	}
	testSerde(t, &restoreOperationParameters)
	selectivePoller, err := backupClient.BeginSelectiveKeyRestore(context.Background(), "selective-restore-test-key", restoreOperationParameters, nil)
	require.NoError(t, err)
	selectiveResults, err := selectivePoller.PollUntilDone(context.Background(), nil)
	require.NoError(t, err)
	testSerde(t, &selectiveResults)

	// additional waiting to avoid conflicts with resources in other tests
	if recording.GetRecordMode() != recording.PlaybackMode {
		time.Sleep(60 * time.Second)
	}
}
