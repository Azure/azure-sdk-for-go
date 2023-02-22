//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azadmin_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin"
	"github.com/stretchr/testify/require"
)

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
	restoreOperationParameters := azadmin.RestoreOperationParameters{
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
	newBackupPoller, err := client.BeginFullBackup(context.Background(), sasToken, &azadmin.BackupClientBeginFullBackupOptions{ResumeToken: token})
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
	restoreOperationParameters := azadmin.RestoreOperationParameters{
		FolderToRestore:    &folderURI,
		SasTokenParameters: &sasToken,
	}
	testSerde(t, &restoreOperationParameters)
	restorePoller, err := client.BeginFullRestore(context.Background(), restoreOperationParameters, nil)
	require.NoError(t, err)

	// create a new poller from a continuation token
	restoreToken, err := restorePoller.ResumeToken()
	require.NoError(t, err)
	newRestorePoller, err := client.BeginFullRestore(context.Background(), restoreOperationParameters, &azadmin.BackupClientBeginFullRestoreOptions{ResumeToken: restoreToken})
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
	restoreOperationParameters := azadmin.SelectiveKeyRestoreOperationParameters{
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
