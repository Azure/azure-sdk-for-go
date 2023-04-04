//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armstoragecache_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagecache/armstoragecache/v3"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e2749bb2cbee0b4c447a9d6c1d7cbce3d415abd4/specification/storagecache/resource-manager/Microsoft.StorageCache/preview/2023-03-01-preview/examples/UsageModels_List.json
func ExampleUsageModelsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstoragecache.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewUsageModelsClient().NewListPager(nil)
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
		// page.UsageModelsResult = armstoragecache.UsageModelsResult{
		// 	Value: []*armstoragecache.UsageModel{
		// 		{
		// 			Display: &armstoragecache.UsageModelDisplay{
		// 				Description: to.Ptr("Read only, with a default verification timer of 30 seconds. Verification timer has a minimum value of 1 and maximum value of 31536000. Write-back timer must have value of 0 or be null."),
		// 			},
		// 			ModelName: to.Ptr("READ_ONLY"),
		// 			TargetType: to.Ptr("nfs3"),
		// 		},
		// 		{
		// 			Display: &armstoragecache.UsageModelDisplay{
		// 				Description: to.Ptr("Read-write, with a default verification timer of 8 hours and default write-back timer of 1 hour. Verification timer and write-back timer have a minimum value of 1 and maximum value of 31536000."),
		// 			},
		// 			ModelName: to.Ptr("READ_WRITE"),
		// 			TargetType: to.Ptr("nfs3"),
		// 	}},
		// }
	}
}
