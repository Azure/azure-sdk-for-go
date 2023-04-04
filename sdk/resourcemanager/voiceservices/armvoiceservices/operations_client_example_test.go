//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armvoiceservices_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/voiceservices/armvoiceservices"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/aa85f59e259c4b12197b57b221067c40fa2fe3f1/specification/voiceservices/resource-manager/Microsoft.VoiceServices/stable/2023-01-31/examples/Operations_List.json
func ExampleOperationsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armvoiceservices.NewClientFactory("<subscription-id>", cred, nil)
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
		// page.OperationListResult = armvoiceservices.OperationListResult{
		// 	Value: []*armvoiceservices.Operation{
		// 		{
		// 			Name: to.Ptr("Microsoft.VoiceService/communicationsGateways/write"),
		// 			Display: &armvoiceservices.OperationDisplay{
		// 				Description: to.Ptr("Write communicationsGateways resource"),
		// 				Operation: to.Ptr("write"),
		// 				Provider: to.Ptr("Microsoft.VoiceService"),
		// 				Resource: to.Ptr("communicationsGateways"),
		// 			},
		// 	}},
		// }
	}
}
