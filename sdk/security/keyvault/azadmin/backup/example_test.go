// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package backup_test

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/backup"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
)

var client *backup.Client
var keyClient *azkeys.Client

func ExampleNewClient() {
	vaultURL := "https://<TODO: your vault name>.managedhsm.azure.net/"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := backup.NewClient(vaultURL, cred, nil)
	if err != nil {
		// TODO: handle error
	}

	_ = client
}

func ExampleClient_BeginFullBackup() {
	storageParameters := backup.SASTokenParameters{
		StorageResourceURI: to.Ptr("https://<storage-account>.blob.core.windows.net/<container>"),
		Token:              to.Ptr("<your SAS token>"),
	}
	backupPoller, err := client.BeginFullBackup(context.Background(), storageParameters, nil)
	if err != nil {
		// TODO: handle error
	}
	backupResults, err := backupPoller.PollUntilDone(context.Background(), nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Printf("Status of backup: %s", *backupResults.Status)
}
func ExampleClient_BeginFullRestore() {
	// first, backup the managed HSM to a blob storage container
	storageParameters := backup.SASTokenParameters{
		StorageResourceURI: to.Ptr("https://<storage-account>.blob.core.windows.net/<container>"),
		Token:              to.Ptr("<your SAS token>"),
	}
	backupPoller, err := client.BeginFullBackup(context.Background(), storageParameters, nil)
	if err != nil {
		// TODO: handle error
	}
	backupResults, err := backupPoller.PollUntilDone(context.Background(), nil)
	if err != nil {
		// TODO: handle error
	}

	// FolderToRestore is the folder in the blob container your managed HSM was uploaded to
	// FolderToRestore can be extracted from the returned backupResults.AzureStorageBlobContainerURI
	s := *backupResults.AzureStorageBlobContainerURI
	folderName := s[strings.LastIndex(s, "/")+1:]

	// begin the restore operation
	restoreOperationParameters := backup.RestoreOperationParameters{
		FolderToRestore: to.Ptr(folderName),
		SASTokenParameters: &backup.SASTokenParameters{
			StorageResourceURI: to.Ptr("https://<storage-account>.blob.core.windows.net/<container>"),
			Token:              to.Ptr("<your SAS token>"),
		},
	}
	restorePoller, err := client.BeginFullRestore(context.Background(), restoreOperationParameters, nil)
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

func ExampleClient_BeginSelectiveKeyRestore() {
	// first, create a key to backup
	params := azkeys.CreateKeyParameters{
		KeySize: to.Ptr(int32(2048)),
		Kty:     to.Ptr(azkeys.KeyTypeRSA),
	}
	_, err := keyClient.CreateKey(context.TODO(), "<key-name>", params, nil)
	if err != nil {
		// TODO: handle error
	}

	// backup the vault
	storageParameters := backup.SASTokenParameters{
		StorageResourceURI: to.Ptr("https://<storage-account>.blob.core.windows.net/<container>"),
		Token:              to.Ptr("<your SAS token>"),
	}
	backupPoller, err := client.BeginFullBackup(context.Background(), storageParameters, nil)
	if err != nil {
		// TODO: handle error
	}
	backupResults, err := backupPoller.PollUntilDone(context.Background(), nil)
	if err != nil {
		// TODO: handle error
	}

	// extract the folder name where the vault was backed up
	s := *backupResults.AzureStorageBlobContainerURI
	folderName := s[strings.LastIndex(s, "/")+1:]

	// restore the key
	restoreOperationParameters := backup.SelectiveKeyRestoreOperationParameters{
		Folder: to.Ptr(folderName),
		SASTokenParameters: &backup.SASTokenParameters{
			StorageResourceURI: to.Ptr("https://<storage-account>.blob.core.windows.net/<container>"),
			Token:              to.Ptr("<your SAS token>"),
		},
	}
	selectivePoller, err := client.BeginSelectiveKeyRestore(context.Background(), "<key-name>", restoreOperationParameters, nil)
	if err != nil {
		// TODO: handle error
	}
	selectiveResults, err := selectivePoller.PollUntilDone(context.Background(), nil)
	if err != nil {
		// TODO: handle error
	}
	fmt.Printf("Status of the selective restore: %s", *selectiveResults.Status)
}
