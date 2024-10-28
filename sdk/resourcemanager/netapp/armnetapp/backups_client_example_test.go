//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armnetapp_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp/v7"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/366aaa13cdd218b9adac716680e49473673410c8/specification/netapp/resource-manager/Microsoft.NetApp/stable/2024-07-01/examples/Volumes_LatestBackupStatus.json
func ExampleBackupsClient_GetLatestStatus() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewBackupsClient().GetLatestStatus(ctx, "myRG", "account1", "pool1", "volume1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.BackupStatus = armnetapp.BackupStatus{
	// 	ErrorMessage: to.Ptr(""),
	// 	Healthy: to.Ptr(true),
	// 	LastTransferSize: to.Ptr[int64](100000),
	// 	LastTransferType: to.Ptr(""),
	// 	MirrorState: to.Ptr(armnetapp.MirrorStateMirrored),
	// 	RelationshipStatus: to.Ptr(armnetapp.RelationshipStatusIdle),
	// 	TotalTransferBytes: to.Ptr[int64](100000),
	// 	UnhealthyReason: to.Ptr(""),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/366aaa13cdd218b9adac716680e49473673410c8/specification/netapp/resource-manager/Microsoft.NetApp/stable/2024-07-01/examples/Volumes_LatestRestoreStatus.json
func ExampleBackupsClient_GetVolumeLatestRestoreStatus() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewBackupsClient().GetVolumeLatestRestoreStatus(ctx, "myRG", "account1", "pool1", "volume1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.RestoreStatus = armnetapp.RestoreStatus{
	// 	ErrorMessage: to.Ptr(""),
	// 	Healthy: to.Ptr(true),
	// 	MirrorState: to.Ptr(armnetapp.MirrorStateUninitialized),
	// 	RelationshipStatus: to.Ptr(armnetapp.RelationshipStatusIdle),
	// 	TotalTransferBytes: to.Ptr[int64](100000),
	// 	UnhealthyReason: to.Ptr(""),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/366aaa13cdd218b9adac716680e49473673410c8/specification/netapp/resource-manager/Microsoft.NetApp/stable/2024-07-01/examples/BackupsUnderBackupVault_List.json
func ExampleBackupsClient_NewListByVaultPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewBackupsClient().NewListByVaultPager("myRG", "account1", "backupVault1", &armnetapp.BackupsClientListByVaultOptions{Filter: nil})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range page.Value {
			// You could use page here. We use blank identifier for just demo purposes.
			_ = v
		}
		// If the HTTP response code is 200 as defined in example definition, your page structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
		// page.BackupsList = armnetapp.BackupsList{
		// 	Value: []*armnetapp.Backup{
		// 		{
		// 			Name: to.Ptr("account1/backupVault1/backup1"),
		// 			Type: to.Ptr("Microsoft.NetApp/netAppAccounts/backupVaults/backups"),
		// 			ID: to.Ptr("/subscriptions/D633CC2E-722B-4AE1-B636-BBD9E4C60ED9/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/backupVaults/backupVault1/backups/backup1"),
		// 			Properties: &armnetapp.BackupProperties{
		// 				BackupPolicyResourceID: to.Ptr("/subscriptions/D633CC2E-722B-4AE1-B636-BBD9E4C60ED9/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/backupPolicies/policy1"),
		// 				BackupType: to.Ptr(armnetapp.BackupTypeManual),
		// 				CreationDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-08-15T13:23:33.000Z"); return t}()),
		// 				Label: to.Ptr("myLabel"),
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				Size: to.Ptr[int64](10011),
		// 				SnapshotName: to.Ptr("backup1"),
		// 				VolumeResourceID: to.Ptr("/subscriptions/D633CC2E-722B-4AE1-B636-BBD9E4C60ED9/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPool/pool1/volumes/volume1"),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/366aaa13cdd218b9adac716680e49473673410c8/specification/netapp/resource-manager/Microsoft.NetApp/stable/2024-07-01/examples/BackupsUnderBackupVault_Get.json
func ExampleBackupsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewBackupsClient().Get(ctx, "myRG", "account1", "backupVault1", "backup1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Backup = armnetapp.Backup{
	// 	Name: to.Ptr("account1/backupVault1/backup1"),
	// 	Type: to.Ptr("Microsoft.NetApp/netAppAccounts/backupVaults/backups"),
	// 	ID: to.Ptr("/subscriptions/D633CC2E-722B-4AE1-B636-BBD9E4C60ED9/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/backupVaults/backupVault1/backups/backup1"),
	// 	Properties: &armnetapp.BackupProperties{
	// 		BackupPolicyResourceID: to.Ptr("/subscriptions/D633CC2E-722B-4AE1-B636-BBD9E4C60ED9/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/backupPolicies/policy1"),
	// 		BackupType: to.Ptr(armnetapp.BackupTypeManual),
	// 		CreationDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-08-15T13:23:33.000Z"); return t}()),
	// 		Label: to.Ptr("myLabel"),
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		Size: to.Ptr[int64](10011),
	// 		SnapshotName: to.Ptr("backup1"),
	// 		VolumeResourceID: to.Ptr("/subscriptions/D633CC2E-722B-4AE1-B636-BBD9E4C60ED9/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPool/pool1/volumes/volume1"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/366aaa13cdd218b9adac716680e49473673410c8/specification/netapp/resource-manager/Microsoft.NetApp/stable/2024-07-01/examples/BackupsUnderBackupVault_Create.json
func ExampleBackupsClient_BeginCreate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewBackupsClient().BeginCreate(ctx, "myRG", "account1", "backupVault1", "backup1", armnetapp.Backup{
		Properties: &armnetapp.BackupProperties{
			Label:            to.Ptr("myLabel"),
			VolumeResourceID: to.Ptr("/subscriptions/D633CC2E-722B-4AE1-B636-BBD9E4C60ED9/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPool/pool1/volumes/volume1"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Backup = armnetapp.Backup{
	// 	Name: to.Ptr("account1/backupVault1/backup1"),
	// 	Type: to.Ptr("Microsoft.NetApp/netAppAccounts/backupVaults/backups"),
	// 	ID: to.Ptr("/subscriptions/D633CC2E-722B-4AE1-B636-BBD9E4C60ED9/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/backupVaults/backupVault1/backups/backup1"),
	// 	Properties: &armnetapp.BackupProperties{
	// 		BackupType: to.Ptr(armnetapp.BackupTypeManual),
	// 		CreationDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-08-15T13:23:33.000Z"); return t}()),
	// 		Label: to.Ptr("myLabel"),
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		Size: to.Ptr[int64](10011),
	// 		SnapshotName: to.Ptr("backup1"),
	// 		VolumeResourceID: to.Ptr("/subscriptions/D633CC2E-722B-4AE1-B636-BBD9E4C60ED9/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPool/pool1/volumes/volume1"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/366aaa13cdd218b9adac716680e49473673410c8/specification/netapp/resource-manager/Microsoft.NetApp/stable/2024-07-01/examples/BackupsUnderBackupVault_Update.json
func ExampleBackupsClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewBackupsClient().BeginUpdate(ctx, "myRG", "account1", "backupVault1", "backup1", armnetapp.BackupPatch{}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Backup = armnetapp.Backup{
	// 	Name: to.Ptr("account1/backupVault1/backup1"),
	// 	Type: to.Ptr("Microsoft.NetApp/netAppAccounts/backupVaults/backups"),
	// 	ID: to.Ptr("/subscriptions/D633CC2E-722B-4AE1-B636-BBD9E4C60ED9/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/backupVaults/backupVault1/backups/backup1"),
	// 	Properties: &armnetapp.BackupProperties{
	// 		BackupType: to.Ptr(armnetapp.BackupTypeManual),
	// 		CreationDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-08-15T13:23:33.000Z"); return t}()),
	// 		Label: to.Ptr("myLabel"),
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		Size: to.Ptr[int64](10011),
	// 		SnapshotName: to.Ptr("backup1"),
	// 		VolumeResourceID: to.Ptr("/subscriptions/D633CC2E-722B-4AE1-B636-BBD9E4C60ED9/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPool/pool1/volumes/volume1"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/366aaa13cdd218b9adac716680e49473673410c8/specification/netapp/resource-manager/Microsoft.NetApp/stable/2024-07-01/examples/BackupsUnderBackupVault_Delete.json
func ExampleBackupsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewBackupsClient().BeginDelete(ctx, "resourceGroup", "account1", "backupVault1", "backup1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
