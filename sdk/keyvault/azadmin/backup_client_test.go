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

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azadmin"
	"github.com/stretchr/testify/require"
)

func TestBackupRestore(t *testing.T) {
	client := startBackupTest(t)

	storageResourceUri := os.Getenv("BLOB_RESOURCE_URI")
	token := os.Getenv("SAS_TOKEN")
	sasToken := azadmin.SASTokenParameter{
		StorageResourceURI: &storageResourceUri,
		Token:              &token,
	}
	testSerde(t, &sasToken)

	// backup the vault
	backupPoller, err := client.BeginFullBackup(context.Background(), sasToken, nil)
	require.NoError(t, err)
	backupResults, err := backupPoller.PollUntilDone(context.Background(), nil)
	require.NoError(t, err)
	require.Nil(t, backupResults.Error)
	require.Equal(t, "Succeeded", *backupResults.Status)
	require.Contains(t, *backupResults.AzureStorageBlobContainerURI, storageResourceUri)
	testSerde(t, &backupResults)

	// restore the backup
	s := *backupResults.AzureStorageBlobContainerURI
	folderURI := s[strings.LastIndex(s, "/")+1:]
	restorePoller, err := client.BeginFullRestore(context.Background(), azadmin.RestoreOperationParameters{
		FolderToRestore:    &folderURI,
		SasTokenParameters: &sasToken,
	}, nil)
	require.NoError(t, err)
	restoreResults, err := restorePoller.PollUntilDone(context.Background(), nil)
	require.NoError(t, err)
	require.Nil(t, restoreResults.Error)
	require.Equal(t, "Succeeded", restoreResults.Status)
	require.NotNil(t, restoreResults.StartTime)
	require.NotNil(t, restoreResults.EndTime)
	require.NotNil(t, restoreResults.JobID)
}

func TestBackupRestoreWithResumeToken(t *testing.T) {
}

func TestBeginSelectiveKeyRestoreOperation(t *testing.T) {
}
