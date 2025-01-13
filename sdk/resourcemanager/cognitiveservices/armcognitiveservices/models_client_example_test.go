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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/cognitiveservices/resource-manager/Microsoft.CognitiveServices/stable/2024-10-01/examples/ListLocationModels.json
func ExampleModelsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcognitiveservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewModelsClient().NewListPager("WestUS", nil)
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
		// page.ModelListResult = armcognitiveservices.ModelListResult{
		// 	Value: []*armcognitiveservices.Model{
		// 		{
		// 			Kind: to.Ptr("OpenAI"),
		// 			Model: &armcognitiveservices.AccountModel{
		// 				Name: to.Ptr("ada"),
		// 				Format: to.Ptr("OpenAI"),
		// 				Version: to.Ptr("1"),
		// 				Capabilities: map[string]*string{
		// 					"FineTuneTokensMaxValue": to.Ptr("37500000"),
		// 					"completion": to.Ptr("true"),
		// 					"fineTune": to.Ptr("true"),
		// 					"inference": to.Ptr("false"),
		// 					"scaleType": to.Ptr("Manual"),
		// 					"search": to.Ptr("true"),
		// 				},
		// 				Deprecation: &armcognitiveservices.ModelDeprecationInfo{
		// 					FineTune: to.Ptr("2024-01-01T00:00:00Z"),
		// 					Inference: to.Ptr("2024-01-01T00:00:00Z"),
		// 				},
		// 				FinetuneCapabilities: map[string]*string{
		// 					"FineTuneTokensMaxValue": to.Ptr("37500000"),
		// 					"completion": to.Ptr("true"),
		// 					"fineTune": to.Ptr("true"),
		// 					"inference": to.Ptr("false"),
		// 					"scaleType": to.Ptr("Manual,Standard"),
		// 					"search": to.Ptr("true"),
		// 				},
		// 				LifecycleStatus: to.Ptr(armcognitiveservices.ModelLifecycleStatusPreview),
		// 				MaxCapacity: to.Ptr[int32](3),
		// 				SKUs: []*armcognitiveservices.ModelSKU{
		// 					{
		// 						Name: to.Ptr("provisioned"),
		// 						Capacity: &armcognitiveservices.CapacityConfig{
		// 							Default: to.Ptr[int32](100),
		// 							Maximum: to.Ptr[int32](1000),
		// 							Minimum: to.Ptr[int32](100),
		// 							Step: to.Ptr[int32](100),
		// 						},
		// 						UsageName: to.Ptr("OpenAI.Provisioned.Class1"),
		// 				}},
		// 				SystemData: &armcognitiveservices.SystemData{
		// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-10-07T00:00:00.000Z"); return t}()),
		// 					CreatedBy: to.Ptr("Microsoft"),
		// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-10-07T00:00:00.000Z"); return t}()),
		// 					LastModifiedBy: to.Ptr("Microsoft"),
		// 				},
		// 			},
		// 			SKUName: to.Ptr("S0"),
		// 	}},
		// }
	}
}
