//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armservicefabric_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicefabric/armservicefabric/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4a04d8730b543e6355ba7e1539799d2b054aae97/specification/servicefabric/resource-manager/Microsoft.ServiceFabric/stable/2021-06-01/examples/ListOperations.json
func ExampleOperationsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armservicefabric.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewOperationsClient().NewListPager(nil)
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
		// page.OperationListResult = armservicefabric.OperationListResult{
		// 	Value: []*armservicefabric.OperationResult{
		// 		{
		// 			Name: to.Ptr("aaaaaaaaaaaaaaaaaaaaaaaaa"),
		// 			Display: &armservicefabric.AvailableOperationDisplay{
		// 				Description: to.Ptr("aaaaaaaaaaaaaaa"),
		// 				Operation: to.Ptr("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		// 				Provider: to.Ptr("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		// 				Resource: to.Ptr("aaaaaaaaaaaaaaaaaa"),
		// 			},
		// 			IsDataAction: to.Ptr(true),
		// 			Origin: to.Ptr("aaaaaaaaaaaaaa"),
		// 	}},
		// }
	}
}
