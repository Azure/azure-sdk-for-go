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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/b43042075540b8d67cce7d3d9f70b9b9f5a359da/specification/search/resource-manager/Microsoft.Search/preview/2025-02-01-preview/examples/GetQuotaUsage.json
func ExampleManagementClient_UsageBySubscriptionSKU() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsearch.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewManagementClient().UsageBySubscriptionSKU(ctx, "westus", "free", &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.QuotaUsageResult = armsearch.QuotaUsageResult{
	// 	Name: &armsearch.QuotaUsageResultName{
	// 		LocalizedValue: to.Ptr("F - Free"),
	// 		Value: to.Ptr("free"),
	// 	},
	// 	CurrentValue: to.Ptr[int32](8),
	// 	ID: to.Ptr("/subscriptions/{subscriptionId}/providers/Microsoft.Search/locations/{location}/usages/{skuName}"),
	// 	Limit: to.Ptr[int32](16),
	// 	Unit: to.Ptr("Count"),
	// }
}
