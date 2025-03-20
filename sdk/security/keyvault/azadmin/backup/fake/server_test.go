// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package fake_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/backup"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/backup/fake"
	"github.com/stretchr/testify/require"
)

var (
	storageURI = "https://test.blob.core.windows.net/backup/testing"
)

func getServer() fake.Server {
	return fake.Server{
		BeginFullBackup: func(ctx context.Context, azureStorageBlobContainerURI backup.SASTokenParameters, options *backup.BeginFullBackupOptions) (resp azfake.PollerResponder[backup.FullBackupResponse], errResp azfake.ErrorResponder) {
			resp.AddNonTerminalResponse(http.StatusAccepted, nil)
			resp.AddNonTerminalResponse(http.StatusAccepted, nil)
			backupResp := backup.FullBackupResponse{}
			backupResp.AzureStorageBlobContainerURI = to.Ptr(storageURI)
			backupResp.Status = to.Ptr("Succeeded")
			resp.SetTerminalResponse(http.StatusOK, backupResp, nil)
			return
		},
		BeginFullRestore: func(ctx context.Context, restoreBlobDetails backup.RestoreOperationParameters, options *backup.BeginFullRestoreOptions) (resp azfake.PollerResponder[backup.FullRestoreResponse], errResp azfake.ErrorResponder) {
			resp.AddNonTerminalResponse(http.StatusAccepted, nil)
			resp.AddNonTerminalResponse(http.StatusAccepted, nil)
			restoreResp := backup.FullRestoreResponse{}
			restoreResp.Status = to.Ptr("Succeeded")
			resp.SetTerminalResponse(http.StatusOK, restoreResp, nil)
			return
		},
		BeginPreFullBackup: func(ctx context.Context, preBackupOperationParameters backup.PreBackupOperationParameters, options *backup.BeginPreFullBackupOptions) (resp azfake.PollerResponder[backup.PreFullBackupResponse], errResp azfake.ErrorResponder) {
			resp.AddNonTerminalResponse(http.StatusAccepted, nil)
			resp.AddNonTerminalResponse(http.StatusAccepted, nil)
			backupResp := backup.PreFullBackupResponse{}
			backupResp.AzureStorageBlobContainerURI = to.Ptr(storageURI)
			backupResp.Status = to.Ptr("Succeeded")
			resp.SetTerminalResponse(http.StatusOK, backupResp, nil)
			return
		},
		BeginPreFullRestore: func(ctx context.Context, preRestoreOperationParameters backup.PreRestoreOperationParameters, options *backup.BeginPreFullRestoreOptions) (resp azfake.PollerResponder[backup.PreFullRestoreResponse], errResp azfake.ErrorResponder) {
			resp.AddNonTerminalResponse(http.StatusAccepted, nil)
			resp.AddNonTerminalResponse(http.StatusAccepted, nil)
			restoreResp := backup.PreFullRestoreResponse{}
			restoreResp.Status = to.Ptr("Succeeded")
			resp.SetTerminalResponse(http.StatusOK, restoreResp, nil)
			return
		},
		BeginSelectiveKeyRestore: func(ctx context.Context, keyName string, restoreBlobDetails backup.SelectiveKeyRestoreOperationParameters, options *backup.BeginSelectiveKeyRestoreOptions) (resp azfake.PollerResponder[backup.SelectiveKeyRestoreResponse], errResp azfake.ErrorResponder) {
			resp.AddNonTerminalResponse(http.StatusAccepted, nil)
			resp.AddNonTerminalResponse(http.StatusAccepted, nil)
			restoreResp := backup.SelectiveKeyRestoreResponse{}
			restoreResp.Status = to.Ptr("Succeeded")
			resp.SetTerminalResponse(http.StatusOK, restoreResp, nil)
			return
		},
	}
}

func TestServer(t *testing.T) {
	fakeServer := getServer()

	client, err := backup.NewClient("https://fake-vault.vault.azure.net", &azfake.TokenCredential{}, &backup.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: fake.NewServerTransport(&fakeServer),
		},
	})
	require.NoError(t, err)

	backupPoller, err := client.BeginFullBackup(context.Background(), backup.SASTokenParameters{StorageResourceURI: to.Ptr("testing")}, nil)
	require.NoError(t, err)
	backupResp, err := backupPoller.PollUntilDone(context.Background(), &runtime.PollUntilDoneOptions{
		Frequency: time.Second,
	})
	require.NoError(t, err)
	require.Equal(t, storageURI, *backupResp.AzureStorageBlobContainerURI)
	require.Equal(t, "Succeeded", *backupResp.Status)

	restorePoller, err := client.BeginFullRestore(context.Background(), backup.RestoreOperationParameters{FolderToRestore: to.Ptr("test"), SASTokenParameters: &backup.SASTokenParameters{StorageResourceURI: &storageURI, Token: to.Ptr("SASToken")}}, nil)
	require.NoError(t, err)
	restoreResp, err := restorePoller.PollUntilDone(context.Background(), &runtime.PollUntilDoneOptions{
		Frequency: time.Second,
	})
	require.NoError(t, err)
	require.Equal(t, "Succeeded", *restoreResp.Status)

	preBackupPoller, err := client.BeginPreFullBackup(context.Background(), backup.PreBackupOperationParameters{StorageResourceURI: to.Ptr("testing")}, nil)
	require.NoError(t, err)
	preBackupResp, err := preBackupPoller.PollUntilDone(context.Background(), &runtime.PollUntilDoneOptions{
		Frequency: time.Second,
	})
	require.NoError(t, err)
	require.Equal(t, storageURI, *preBackupResp.AzureStorageBlobContainerURI)
	require.Equal(t, "Succeeded", *preBackupResp.Status)

	preRestorePoller, err := client.BeginPreFullRestore(context.Background(), backup.PreRestoreOperationParameters{FolderToRestore: to.Ptr("test"), SASTokenParameters: &backup.SASTokenParameters{StorageResourceURI: &storageURI, Token: to.Ptr("SASToken")}}, nil)
	require.NoError(t, err)
	preRestoreResp, err := preRestorePoller.PollUntilDone(context.Background(), &runtime.PollUntilDoneOptions{
		Frequency: time.Second,
	})
	require.NoError(t, err)
	require.Equal(t, "Succeeded", *preRestoreResp.Status)

	selectiveRestorePoller, err := client.BeginSelectiveKeyRestore(context.Background(), "key-name", backup.SelectiveKeyRestoreOperationParameters{Folder: to.Ptr("test"), SASTokenParameters: &backup.SASTokenParameters{StorageResourceURI: &storageURI, Token: to.Ptr("SASToken")}}, nil)
	require.NoError(t, err)
	selectiveRestoreResp, err := selectiveRestorePoller.PollUntilDone(context.Background(), &runtime.PollUntilDoneOptions{
		Frequency: time.Second,
	})
	require.NoError(t, err)
	require.Equal(t, "Succeeded", *selectiveRestoreResp.Status)
}
