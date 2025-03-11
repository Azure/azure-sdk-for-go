//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcontainerservice_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v6"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/913304d94f4cd9fb66aa3c72f6ed897b12b38b90/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/preview/2024-10-02-preview/examples/NodeImageVersions_List.json
func ExampleClient_NewListNodeImageVersionsPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewClient().NewListNodeImageVersionsPager("location1", nil)
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
		// page.NodeImageVersionsListResult = armcontainerservice.NodeImageVersionsListResult{
		// 	Value: []*armcontainerservice.NodeImageVersion{
		// 		{
		// 			FullName: to.Ptr("AKSCBLMariner-V1-202308.28.0"),
		// 			OS: to.Ptr("AKSCBLMariner"),
		// 			SKU: to.Ptr("V1"),
		// 			Version: to.Ptr("202308.28.0"),
		// 		},
		// 		{
		// 			FullName: to.Ptr("AKSUbuntu-2204gen2minimalcontainerd-202401.12.0"),
		// 			OS: to.Ptr("AKSUbuntu"),
		// 			SKU: to.Ptr("2204gen2minimalcontainerd"),
		// 			Version: to.Ptr("202401.12.0"),
		// 	}},
		// }
	}
}
