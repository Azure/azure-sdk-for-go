//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcosmos_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v3"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/011ecc5633300a5eefe43dde748f269d39e96458/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/stable/2025-04-15/examples/CosmosDBRestorableSqlContainerList.json
func ExampleRestorableSQLContainersClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcosmos.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewRestorableSQLContainersClient().NewListPager("WestUS", "98a570f2-63db-4117-91f0-366327b7b353", &armcosmos.RestorableSQLContainersClientListOptions{RestorableSQLDatabaseRid: to.Ptr("3fu-hg=="),
		StartTime: nil,
		EndTime:   nil,
	})
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
		// page.RestorableSQLContainersListResult = armcosmos.RestorableSQLContainersListResult{
		// 	Value: []*armcosmos.RestorableSQLContainerGetResult{
		// 		{
		// 			Name: to.Ptr("79609a98-3394-41f8-911f-cfab0c075c86"),
		// 			Type: to.Ptr("Microsoft.DocumentDB/locations/restorableDatabaseAccounts/restorableSqlContainers"),
		// 			ID: to.Ptr("/subscriptions/subid/providers/Microsoft.DocumentDb/locations/westus/restorableDatabaseAccounts/98a570f2-63db-4117-91f0-366327b7b353/restorableSqlContainers/79609a98-3394-41f8-911f-cfab0c075c86"),
		// 			Properties: &armcosmos.RestorableSQLContainerProperties{
		// 				Resource: &armcosmos.RestorableSQLContainerPropertiesResource{
		// 					Rid: to.Ptr("zAyAPQAAAA=="),
		// 					CanUndelete: to.Ptr("invalid"),
		// 					Container: &armcosmos.RestorableSQLContainerPropertiesResourceContainer{
		// 						Etag: to.Ptr("\"00003e00-0000-0700-0000-5f85338a0000\""),
		// 						Rid: to.Ptr("V18LoLrv-qA="),
		// 						ConflictResolutionPolicy: &armcosmos.ConflictResolutionPolicy{
		// 							ConflictResolutionPath: to.Ptr("/_ts"),
		// 							ConflictResolutionProcedure: to.Ptr(""),
		// 							Mode: to.Ptr(armcosmos.ConflictResolutionModeLastWriterWins),
		// 						},
		// 						ID: to.Ptr("Container1"),
		// 						IndexingPolicy: &armcosmos.IndexingPolicy{
		// 							Automatic: to.Ptr(true),
		// 							ExcludedPaths: []*armcosmos.ExcludedPath{
		// 								{
		// 									Path: to.Ptr("/\"_etag\"/?"),
		// 							}},
		// 							IncludedPaths: []*armcosmos.IncludedPath{
		// 								{
		// 									Path: to.Ptr("/*"),
		// 								},
		// 								{
		// 									Path: to.Ptr("/\"_ts\"/?"),
		// 							}},
		// 							IndexingMode: to.Ptr(armcosmos.IndexingModeConsistent),
		// 						},
		// 						Self: to.Ptr("dbs/V18LoA==/colls/V18LoLrv-qA=/"),
		// 					},
		// 					EventTimestamp: to.Ptr("2020-10-13T04:56:42Z"),
		// 					OperationType: to.Ptr(armcosmos.OperationTypeCreate),
		// 					OwnerID: to.Ptr("Container1"),
		// 					OwnerResourceID: to.Ptr("V18LoLrv-qA="),
		// 				},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("e85298a1-c631-4726-825e-a7ca092e9098"),
		// 			Type: to.Ptr("Microsoft.DocumentDB/locations/restorableDatabaseAccounts/restorableSqlContainers"),
		// 			ID: to.Ptr("/subscriptions/subid/providers/Microsoft.DocumentDb/locations/westus/restorableDatabaseAccounts/98a570f2-63db-4117-91f0-366327b7b353/restorableSqlContainers/e85298a1-c631-4726-825e-a7ca092e9098"),
		// 			Properties: &armcosmos.RestorableSQLContainerProperties{
		// 				Resource: &armcosmos.RestorableSQLContainerPropertiesResource{
		// 					Rid: to.Ptr("PrArcgAAAA=="),
		// 					CanUndelete: to.Ptr("invalid"),
		// 					Container: &armcosmos.RestorableSQLContainerPropertiesResourceContainer{
		// 						Etag: to.Ptr("\"00004400-0000-0700-0000-5f85351f0000\""),
		// 						Rid: to.Ptr("V18LoLrv-qA="),
		// 						ConflictResolutionPolicy: &armcosmos.ConflictResolutionPolicy{
		// 							ConflictResolutionPath: to.Ptr("/_ts"),
		// 							ConflictResolutionProcedure: to.Ptr(""),
		// 							Mode: to.Ptr(armcosmos.ConflictResolutionModeLastWriterWins),
		// 						},
		// 						DefaultTTL: to.Ptr[int32](12345),
		// 						ID: to.Ptr("Container1"),
		// 						IndexingPolicy: &armcosmos.IndexingPolicy{
		// 							Automatic: to.Ptr(true),
		// 							ExcludedPaths: []*armcosmos.ExcludedPath{
		// 								{
		// 									Path: to.Ptr("/\"_etag\"/?"),
		// 							}},
		// 							IncludedPaths: []*armcosmos.IncludedPath{
		// 								{
		// 									Path: to.Ptr("/*"),
		// 								},
		// 								{
		// 									Path: to.Ptr("/\"_ts\"/?"),
		// 							}},
		// 							IndexingMode: to.Ptr(armcosmos.IndexingModeConsistent),
		// 						},
		// 						Self: to.Ptr("dbs/V18LoA==/colls/V18LoLrv-qA=/"),
		// 					},
		// 					EventTimestamp: to.Ptr("2020-10-13T05:03:27Z"),
		// 					OperationType: to.Ptr(armcosmos.OperationTypeReplace),
		// 					OwnerID: to.Ptr("Container1"),
		// 					OwnerResourceID: to.Ptr("V18LoLrv-qA="),
		// 				},
		// 			},
		// 	}},
		// }
	}
}
