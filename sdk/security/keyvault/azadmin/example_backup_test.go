// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azadmin_test

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin"
)

var backupClient *azadmin.BackupClient

func ExampleNewBackupClient() {
	vaultURL := "https://<TODO: your vault name>.managedhsm.azure.net/"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azadmin.NewBackupClient(vaultURL, cred, nil)
	if err != nil {
		// TODO: handle error
	}

	_ = client
}

func ExampleBackupClient_BeginFullBackup() {
	storageParameters := azadmin.SASTokenParameter{
		StorageResourceURI: to.Ptr("https://<storage-account>.blob.core.windows.net/<container>"),
		Token:              to.Ptr("SAS_TOKEN"),
	}
	backupPoller, err := backupClient.BeginFullBackup(context.Background(), storageParameters, nil)
	if err != nil {
		// TODO: handle error
	}
	backupResults, err := backupPoller.PollUntilDone(context.Background(), nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Printf("Status of backup: %s", *backupResults.Status)
}
func ExampleBackupClient_BeginFullRestore() {

	// FolderToRestore can be extracted from blob url returned from the backup operation

	restoreOperationParameters := azadmin.RestoreOperationParameters{
		FolderToRestore: to.Ptr("FOLDER_NAME"),
		SasTokenParameters: &azadmin.SASTokenParameter{
			StorageResourceURI: to.Ptr("https://<storage-account>.blob.core.windows.net/<container>"),
			Token:              to.Ptr("SAS_TOKEN"),
		},
	}
	restorePoller, err := backupClient.BeginFullRestore(context.Background(), restoreOperationParameters, nil)
	if err != nil {
		// TODO: handle error
	}

	// Poll for the results
	restoreResults, err := restorePoller.PollUntilDone(context.Background(), nil)
	if err != nil {
		// TODO: handle error
	}

	fmt.Printf("Status of restore: %s", *restoreResults.Status)
}

func ExampleBackupClient_BeginSelectiveKeyRestore() {

	restoreOperationParameters := azadmin.SelectiveKeyRestoreOperationParameters{
		Folder: to.Ptr("FOLDER_NAME"),
		SasTokenParameters: &azadmin.SASTokenParameter{
			StorageResourceURI: to.Ptr("https://<storage-account>.blob.core.windows.net/<container>"),
			Token:              to.Ptr("SAS_TOKEN"),
		},
	}
	selectivePoller, err := backupClient.BeginSelectiveKeyRestore(context.Background(), "keyName", restoreOperationParameters, nil)
	if err != nil {
		// TODO: handle error
	}
	selectiveResults, err := selectivePoller.PollUntilDone(context.Background(), nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Printf("Status of the selective restore: %s", *selectiveResults.Status)
}
