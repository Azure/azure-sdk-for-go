//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armhardwaresecuritymodules_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hardwaresecuritymodules/armhardwaresecuritymodules/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/41e4538ed7bb3ceac3c1322c9455a0812ed110ac/specification/hardwaresecuritymodules/resource-manager/Microsoft.HardwareSecurityModules/preview/2023-12-10-preview/examples/CloudHsmClusterPrivateLinkResource_ListByCloudHsmCluster_MaximumSet_Gen.json
func ExampleCloudHsmClusterPrivateLinkResourcesClient_ListByCloudHsmCluster() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhardwaresecuritymodules.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewCloudHsmClusterPrivateLinkResourcesClient().ListByCloudHsmCluster(ctx, "rgcloudhsm", "chsm1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.PrivateLinkResourceListResult = armhardwaresecuritymodules.PrivateLinkResourceListResult{
	// 	Value: []*armhardwaresecuritymodules.PrivateLinkResource{
	// 		{
	// 			Name: to.Ptr("sample-pls"),
	// 			Type: to.Ptr("Microsoft.HardwareSecurityModules/cloudHsmClusters/privateLinkResources"),
	// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rgcloudhsm/providers/Microsoft.HardwareSecurityModules/cloudHsmClusters/chsm1/privateLinkResources/sample-pls"),
	// 			Properties: &armhardwaresecuritymodules.PrivateLinkResourceProperties{
	// 				GroupID: to.Ptr("cloudHsm"),
	// 				RequiredMembers: []*string{
	// 					to.Ptr("hsm1"),
	// 					to.Ptr("hsm2"),
	// 					to.Ptr("hsm3")},
	// 					RequiredZoneNames: []*string{
	// 						to.Ptr("privatelink.cloudhsm.azure-int.net")},
	// 					},
	// 			}},
	// 		}
}
