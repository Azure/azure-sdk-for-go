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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/28a163a3b110ac1d7002449ee2d04553c725ec0d/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/preview/2024-10-02-preview/examples/TrustedAccessRoles_List.json
func ExampleTrustedAccessRolesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewTrustedAccessRolesClient().NewListPager("westus2", nil)
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
		// page.TrustedAccessRoleListResult = armcontainerservice.TrustedAccessRoleListResult{
		// 	Value: []*armcontainerservice.TrustedAccessRole{
		// 		{
		// 			Name: to.Ptr("reader"),
		// 			Rules: []*armcontainerservice.TrustedAccessRoleRule{
		// 				{
		// 					APIGroups: []*string{
		// 						to.Ptr("")},
		// 						NonResourceURLs: []*string{
		// 						},
		// 						ResourceNames: []*string{
		// 						},
		// 						Resources: []*string{
		// 							to.Ptr("pods")},
		// 							Verbs: []*string{
		// 								to.Ptr("get")},
		// 						}},
		// 						SourceResourceType: to.Ptr("Microsoft.MachineLearningServices/workspaces"),
		// 				}},
		// 			}
	}
}
