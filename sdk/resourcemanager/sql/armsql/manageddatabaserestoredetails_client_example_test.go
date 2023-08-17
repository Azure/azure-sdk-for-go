//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armsql_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/fe78d8f1e7bd86c778c7e1cafd52cb0e9fec67ef/specification/sql/resource-manager/Microsoft.Sql/preview/2022-05-01-preview/examples/ManagedDatabaseRestoreDetails.json
func ExampleManagedDatabaseRestoreDetailsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsql.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewManagedDatabaseRestoreDetailsClient().Get(ctx, "Default-SQL-SouthEastAsia", "managedInstance", "testdb", armsql.RestoreDetailsNameDefault, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ManagedDatabaseRestoreDetailsResult = armsql.ManagedDatabaseRestoreDetailsResult{
	// 	Name: to.Ptr("current"),
	// 	Type: to.Ptr("Microsoft.Sql/managedInstances/databases/restoreDetails"),
	// 	ID: to.Ptr("/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/Default-SQL-SouthEastAsia/providers/Microsoft.Sql/managedInstances/managedInstance/databases/testdb/restoreDetails/current"),
	// 	Properties: &armsql.ManagedDatabaseRestoreDetailsProperties{
	// 		Type: to.Ptr("LRSRestore"),
	// 		CurrentBackupType: to.Ptr("Log"),
	// 		CurrentRestorePlanSizeMB: to.Ptr[int32](47),
	// 		CurrentRestoredSizeMB: to.Ptr[int32](25),
	// 		CurrentRestoringFileName: to.Ptr("RestoreDetailsFullBlownExampleLog10.bak"),
	// 		DiffBackupSets: []*armsql.ManagedDatabaseRestoreDetailsBackupSetProperties{
	// 			{
	// 				BackupSizeMB: to.Ptr[int32](0),
	// 				FirstStripeName: to.Ptr("RestoreDetailsFullBlownExampleDiff2.bak"),
	// 				NumberOfStripes: to.Ptr[int32](1),
	// 				Status: to.Ptr("Skipped"),
	// 			},
	// 			{
	// 				BackupSizeMB: to.Ptr[int32](1),
	// 				FirstStripeName: to.Ptr("RestoreDetailsFullBlownExampleDiff3_1.bak"),
	// 				NumberOfStripes: to.Ptr[int32](3),
	// 				RestoreFinishedTimestampUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-09T12:20:21.3667454Z"); return t}()),
	// 				RestoreStartedTimestampUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-09T12:19:40.5455092Z"); return t}()),
	// 				Status: to.Ptr("Restored"),
	// 		}},
	// 		FullBackupSets: []*armsql.ManagedDatabaseRestoreDetailsBackupSetProperties{
	// 			{
	// 				BackupSizeMB: to.Ptr[int32](2),
	// 				FirstStripeName: to.Ptr("RestoreDetailsFullBlownExampleFull2.bak"),
	// 				NumberOfStripes: to.Ptr[int32](1),
	// 				Status: to.Ptr("Skipped"),
	// 			},
	// 			{
	// 				BackupSizeMB: to.Ptr[int32](3),
	// 				FirstStripeName: to.Ptr("RestoreDetailsFullBlownExampleFull3_1.bak"),
	// 				NumberOfStripes: to.Ptr[int32](3),
	// 				RestoreFinishedTimestampUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-09T12:18:41.3785089Z"); return t}()),
	// 				RestoreStartedTimestampUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-09T12:15:40.3143263Z"); return t}()),
	// 				Status: to.Ptr("Restored"),
	// 		}},
	// 		LastRestoredFileName: to.Ptr("RestoreDetailsFullBlownExampleLog9_1.bak"),
	// 		LastRestoredFileTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-09T12:27:20.4936305Z"); return t}()),
	// 		LastUploadedFileName: to.Ptr("RestoreDetailsFullBlownExampleLog11.bak"),
	// 		LastUploadedFileTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-01T07:54:21Z"); return t}()),
	// 		LogBackupSets: []*armsql.ManagedDatabaseRestoreDetailsBackupSetProperties{
	// 			{
	// 				BackupSizeMB: to.Ptr[int32](0),
	// 				FirstStripeName: to.Ptr("RestoreDetailsFullBlownExampleLog2.bak"),
	// 				NumberOfStripes: to.Ptr[int32](1),
	// 				Status: to.Ptr("Skipped"),
	// 			},
	// 			{
	// 				BackupSizeMB: to.Ptr[int32](8),
	// 				FirstStripeName: to.Ptr("RestoreDetailsFullBlownExampleLog3.bak"),
	// 				NumberOfStripes: to.Ptr[int32](1),
	// 				Status: to.Ptr("Skipped"),
	// 			},
	// 			{
	// 				BackupSizeMB: to.Ptr[int32](11),
	// 				FirstStripeName: to.Ptr("RestoreDetailsFullBlownExampleLog4.bak"),
	// 				NumberOfStripes: to.Ptr[int32](1),
	// 				Status: to.Ptr("Skipped"),
	// 			},
	// 			{
	// 				BackupSizeMB: to.Ptr[int32](7),
	// 				FirstStripeName: to.Ptr("RestoreDetailsFullBlownExampleLog5.bak"),
	// 				NumberOfStripes: to.Ptr[int32](1),
	// 				RestoreFinishedTimestampUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-09T12:21:54.7557851Z"); return t}()),
	// 				RestoreStartedTimestampUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-09T12:21:01.7717453Z"); return t}()),
	// 				Status: to.Ptr("Restored"),
	// 			},
	// 			{
	// 				BackupSizeMB: to.Ptr[int32](3),
	// 				FirstStripeName: to.Ptr("RestoreDetailsFullBlownExampleLog6.bak"),
	// 				NumberOfStripes: to.Ptr[int32](1),
	// 				RestoreFinishedTimestampUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-09T12:23:03.709407Z"); return t}()),
	// 				RestoreStartedTimestampUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-09T12:22:41.8784062Z"); return t}()),
	// 				Status: to.Ptr("Restored"),
	// 			},
	// 			{
	// 				BackupSizeMB: to.Ptr[int32](4),
	// 				FirstStripeName: to.Ptr("RestoreDetailsFullBlownExampleLog7.bak"),
	// 				NumberOfStripes: to.Ptr[int32](1),
	// 				RestoreFinishedTimestampUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-09T12:23:52.9274047Z"); return t}()),
	// 				RestoreStartedTimestampUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-09T12:23:36.6264066Z"); return t}()),
	// 				Status: to.Ptr("Restored"),
	// 			},
	// 			{
	// 				BackupSizeMB: to.Ptr[int32](3),
	// 				FirstStripeName: to.Ptr("RestoreDetailsFullBlownExampleLog8.bak"),
	// 				NumberOfStripes: to.Ptr[int32](1),
	// 				RestoreFinishedTimestampUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-09T12:24:37.9954063Z"); return t}()),
	// 				RestoreStartedTimestampUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-09T12:24:25.899407Z"); return t}()),
	// 				Status: to.Ptr("Restored"),
	// 			},
	// 			{
	// 				BackupSizeMB: to.Ptr[int32](4),
	// 				FirstStripeName: to.Ptr("RestoreDetailsFullBlownExampleLog9_1.bak"),
	// 				NumberOfStripes: to.Ptr[int32](4),
	// 				RestoreFinishedTimestampUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-09T12:25:27.808409Z"); return t}()),
	// 				RestoreStartedTimestampUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-09T12:25:10.8804065Z"); return t}()),
	// 				Status: to.Ptr("Restored"),
	// 			},
	// 			{
	// 				BackupSizeMB: to.Ptr[int32](15),
	// 				FirstStripeName: to.Ptr("RestoreDetailsFullBlownExampleLog10.bak"),
	// 				NumberOfStripes: to.Ptr[int32](1),
	// 				RestoreStartedTimestampUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-09T12:26:00.7813103Z"); return t}()),
	// 				Status: to.Ptr("Restoring"),
	// 			},
	// 			{
	// 				BackupSizeMB: to.Ptr[int32](7),
	// 				FirstStripeName: to.Ptr("RestoreDetailsFullBlownExampleLog11.bak"),
	// 				NumberOfStripes: to.Ptr[int32](1),
	// 				Status: to.Ptr("Queued"),
	// 		}},
	// 		NumberOfFilesDetected: to.Ptr[int32](25),
	// 		NumberOfFilesQueued: to.Ptr[int32](1),
	// 		NumberOfFilesRestored: to.Ptr[int32](14),
	// 		NumberOfFilesRestoring: to.Ptr[int32](1),
	// 		NumberOfFilesSkipped: to.Ptr[int32](8),
	// 		NumberOfFilesUnrestorable: to.Ptr[int32](1),
	// 		PercentCompleted: to.Ptr[int32](53),
	// 		Status: to.Ptr("Restoring"),
	// 		UnrestorableFiles: []*armsql.ManagedDatabaseRestoreDetailsUnrestorableFileProperties{
	// 			{
	// 				Name: to.Ptr("ImageFile.JPG"),
	// 		}},
	// 	},
	// }
}
