//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armsearch_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/search/armsearch"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/f6f50c6388fd5836fa142384641b8353a99874ef/specification/search/resource-manager/Microsoft.Search/preview/2024-06-01-preview/examples/GetQuotaUsagesList.json
func ExampleUsagesClient_NewListBySubscriptionPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsearch.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewUsagesClient().NewListBySubscriptionPager("westus", &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
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
		// page.QuotaUsagesListResult = armsearch.QuotaUsagesListResult{
		// 	Value: []*armsearch.QuotaUsageResult{
		// 		{
		// 			Name: &armsearch.QuotaUsageResultName{
		// 				LocalizedValue: to.Ptr("F - Free"),
		// 				Value: to.Ptr("free"),
		// 			},
		// 			CurrentValue: to.Ptr[int32](8),
		// 			ID: to.Ptr("/subscriptions/{subscriptionId}/providers/Microsoft.Search/locations/{location}/usages/free"),
		// 			Limit: to.Ptr[int32](16),
		// 			Unit: to.Ptr("Count"),
		// 		},
		// 		{
		// 			Name: &armsearch.QuotaUsageResultName{
		// 				LocalizedValue: to.Ptr("B - Basic"),
		// 				Value: to.Ptr("basic"),
		// 			},
		// 			CurrentValue: to.Ptr[int32](8),
		// 			ID: to.Ptr("/subscriptions/{subscriptionId}/providers/Microsoft.Search/locations/{location}/usages/basic"),
		// 			Limit: to.Ptr[int32](16),
		// 			Unit: to.Ptr("Count"),
		// 		},
		// 		{
		// 			Name: &armsearch.QuotaUsageResultName{
		// 				LocalizedValue: to.Ptr("S - Standard"),
		// 				Value: to.Ptr("standard"),
		// 			},
		// 			CurrentValue: to.Ptr[int32](8),
		// 			ID: to.Ptr("/subscriptions/{subscriptionId}/providers/Microsoft.Search/locations/{location}/usages/standard"),
		// 			Limit: to.Ptr[int32](16),
		// 			Unit: to.Ptr("Count"),
		// 		},
		// 		{
		// 			Name: &armsearch.QuotaUsageResultName{
		// 				LocalizedValue: to.Ptr("S2 - Standard2"),
		// 				Value: to.Ptr("standard2"),
		// 			},
		// 			CurrentValue: to.Ptr[int32](8),
		// 			ID: to.Ptr("/subscriptions/{subscriptionId}/providers/Microsoft.Search/locations/{location}/usages/standard2"),
		// 			Limit: to.Ptr[int32](16),
		// 			Unit: to.Ptr("Count"),
		// 		},
		// 		{
		// 			Name: &armsearch.QuotaUsageResultName{
		// 				LocalizedValue: to.Ptr("S3 - Standard3"),
		// 				Value: to.Ptr("standard3"),
		// 			},
		// 			CurrentValue: to.Ptr[int32](8),
		// 			ID: to.Ptr("/subscriptions/{subscriptionId}/providers/Microsoft.Search/locations/{location}/usages/standard3"),
		// 			Limit: to.Ptr[int32](16),
		// 			Unit: to.Ptr("Count"),
		// 		},
		// 		{
		// 			Name: &armsearch.QuotaUsageResultName{
		// 				LocalizedValue: to.Ptr("L1 - Storage Optimized"),
		// 				Value: to.Ptr("storageOptimizedL1"),
		// 			},
		// 			CurrentValue: to.Ptr[int32](8),
		// 			ID: to.Ptr("/subscriptions/{subscriptionId}/providers/Microsoft.Search/locations/{location}/usages/storageOptimizedL1"),
		// 			Limit: to.Ptr[int32](16),
		// 			Unit: to.Ptr("Count"),
		// 		},
		// 		{
		// 			Name: &armsearch.QuotaUsageResultName{
		// 				LocalizedValue: to.Ptr("L2 - Storage Optimized"),
		// 				Value: to.Ptr("storageOptimizedL2"),
		// 			},
		// 			CurrentValue: to.Ptr[int32](8),
		// 			ID: to.Ptr("/subscriptions/{subscriptionId}/providers/Microsoft.Search/locations/{location}/usages/storageOptimizedL2"),
		// 			Limit: to.Ptr[int32](16),
		// 			Unit: to.Ptr("Count"),
		// 	}},
		// }
	}
}
