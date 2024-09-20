//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package backup_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/backup"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/stretchr/testify/require"
)

func TestBackupRestore(t *testing.T) {
	client, sasToken := startBackupTest(t)
	testSerde(t, &sasToken)

	// if in live mode, wait 30 seconds before polling again
	// in playback mode, we don't need that delay, so shortening the retry frequency
	var frequency time.Duration
	if recording.GetRecordMode() == recording.PlaybackMode {
		frequency = time.Second
	}

	// backup the vault
	backupPoller, err := client.BeginFullBackup(context.Background(), sasToken, nil)
	require.NoError(t, err)
	backupResults, err := backupPoller.PollUntilDone(context.Background(), &runtime.PollUntilDoneOptions{
		Frequency: frequency,
	})
	require.NoError(t, err)
	require.Nil(t, backupResults.Error)
	require.Equal(t, "Succeeded", *backupResults.Status)
	require.Contains(t, *backupResults.AzureStorageBlobContainerURI, blobURL)
	testSerde(t, &backupResults)

	// restore the backup
	s := *backupResults.AzureStorageBlobContainerURI
	folderURI := s[strings.LastIndex(s, "/")+1:]
	restoreOperationParameters := backup.RestoreOperationParameters{
		FolderToRestore:    &folderURI,
		SASTokenParameters: &sasToken,
	}
	testSerde(t, &restoreOperationParameters)
	restorePoller, err := client.BeginFullRestore(context.Background(), restoreOperationParameters, nil)
	require.NoError(t, err)
	restoreResults, err := restorePoller.PollUntilDone(context.Background(), &runtime.PollUntilDoneOptions{
		Frequency: frequency,
	})
	require.NoError(t, err)
	require.Nil(t, restoreResults.Error)
	require.Equal(t, "Succeeded", *restoreResults.Status)
	require.NotNil(t, restoreResults.StartTime)
	require.NotNil(t, restoreResults.EndTime)
	require.NotNil(t, restoreResults.JobID)
	testSerde(t, &restoreResults)

	// additional waiting to avoid conflicts with resources in other tests
	if recording.GetRecordMode() != recording.PlaybackMode {
		time.Sleep(61 * time.Second)
	}
}

func TestBackupRestoreWithResumeToken(t *testing.T) {
	client, sasToken := startBackupTest(t)

	// if in live mode, wait 30 seconds before polling again
	// in playback mode, we don't need that delay, so shortening the retry frequency
	var frequency time.Duration
	if recording.GetRecordMode() == recording.PlaybackMode {
		frequency = time.Second
	}

	// backup the vault
	backupPoller, err := client.BeginFullBackup(context.Background(), sasToken, nil)
	require.NoError(t, err)

	// create a new poller from a continuation token
	token, err := backupPoller.ResumeToken()
	require.NoError(t, err)
	newBackupPoller, err := client.BeginFullBackup(context.Background(), sasToken, &backup.BeginFullBackupOptions{ResumeToken: token})
	require.NoError(t, err)
	backupResults, err := newBackupPoller.PollUntilDone(context.Background(), &runtime.PollUntilDoneOptions{
		Frequency: frequency,
	})
	require.NoError(t, err)
	require.Nil(t, backupResults.Error)
	require.Equal(t, "Succeeded", *backupResults.Status)
	require.Contains(t, *backupResults.AzureStorageBlobContainerURI, blobURL)
	testSerde(t, &backupResults)

	// restore the backup
	s := *backupResults.AzureStorageBlobContainerURI
	folderURI := s[strings.LastIndex(s, "/")+1:]
	restoreOperationParameters := backup.RestoreOperationParameters{
		FolderToRestore:    &folderURI,
		SASTokenParameters: &sasToken,
	}
	testSerde(t, &restoreOperationParameters)
	restorePoller, err := client.BeginFullRestore(context.Background(), restoreOperationParameters, nil)
	require.NoError(t, err)

	// create a new poller from a continuation token
	restoreToken, err := restorePoller.ResumeToken()
	require.NoError(t, err)
	newRestorePoller, err := client.BeginFullRestore(context.Background(), restoreOperationParameters, &backup.BeginFullRestoreOptions{ResumeToken: restoreToken})
	require.NoError(t, err)
	restoreResults, err := newRestorePoller.PollUntilDone(context.Background(), &runtime.PollUntilDoneOptions{
		Frequency: frequency,
	})
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

	// if in live mode, wait 30 seconds before polling again
	// in playback mode, we don't need that delay, so shortening the retry frequency
	var frequency time.Duration
	if recording.GetRecordMode() == recording.PlaybackMode {
		frequency = time.Second
	}

	// create a key to selectively restore
	if recording.GetRecordMode() != recording.PlaybackMode {
		cred := credential
		keyClient, err := azkeys.NewClient(hsmURL, cred, nil)
		require.NoError(t, err)
		params := azkeys.CreateKeyParameters{
			KeySize: to.Ptr(int32(2048)),
			Kty:     to.Ptr(azkeys.KeyTypeRSA),
		}
		_, err = keyClient.CreateKey(context.TODO(), "selective-restore-test-key", params, nil)
		require.NoError(t, err)
	}

	// backup the vault
	backupPoller, err := backupClient.BeginFullBackup(context.Background(), sasToken, nil)
	require.NoError(t, err)
	backupResults, err := backupPoller.PollUntilDone(context.Background(), &runtime.PollUntilDoneOptions{
		Frequency: frequency,
	})
	require.NoError(t, err)

	// restore the key
	s := *backupResults.AzureStorageBlobContainerURI
	folderURI := s[strings.LastIndex(s, "/")+1:]
	restoreOperationParameters := backup.SelectiveKeyRestoreOperationParameters{
		Folder:             &folderURI,
		SASTokenParameters: &sasToken,
	}
	testSerde(t, &restoreOperationParameters)
	selectivePoller, err := backupClient.BeginSelectiveKeyRestore(context.Background(), "selective-restore-test-key", restoreOperationParameters, nil)
	require.NoError(t, err)
	selectiveResults, err := selectivePoller.PollUntilDone(context.Background(), &runtime.PollUntilDoneOptions{
		Frequency: frequency,
	})
	require.NoError(t, err)
	testSerde(t, &selectiveResults)

	// additional waiting to avoid conflicts with resources in other tests
	if recording.GetRecordMode() != recording.PlaybackMode {
		time.Sleep(60 * time.Second)
	}
}
