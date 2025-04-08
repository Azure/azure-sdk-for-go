//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcognitiveservices_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/cognitiveservices/resource-manager/Microsoft.CognitiveServices/stable/2024-10-01/examples/ListCommitmentTiers.json
func ExampleCommitmentTiersClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcognitiveservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewCommitmentTiersClient().NewListPager("location", nil)
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
		// page.CommitmentTierListResult = armcognitiveservices.CommitmentTierListResult{
		// 	Value: []*armcognitiveservices.CommitmentTier{
		// 		{
		// 			Cost: &armcognitiveservices.CommitmentCost{
		// 			},
		// 			HostingModel: to.Ptr(armcognitiveservices.HostingModelWeb),
		// 			Kind: to.Ptr("TextAnalytics"),
		// 			PlanType: to.Ptr("TA"),
		// 			Quota: &armcognitiveservices.CommitmentQuota{
		// 				Quantity: to.Ptr[int64](1000000),
		// 				Unit: to.Ptr("Transaction"),
		// 			},
		// 			SKUName: to.Ptr("S"),
		// 			Tier: to.Ptr("T1"),
		// 	}},
		// }
	}
}
