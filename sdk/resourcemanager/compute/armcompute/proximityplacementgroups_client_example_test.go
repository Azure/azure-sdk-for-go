//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcompute_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/proximityPlacementGroupExamples/ProximityPlacementGroup_CreateOrUpdate.json
func ExampleProximityPlacementGroupsClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewProximityPlacementGroupsClient().CreateOrUpdate(ctx, "myResourceGroup", "myProximityPlacementGroup", armcompute.ProximityPlacementGroup{
		Location: to.Ptr("westus"),
		Properties: &armcompute.ProximityPlacementGroupProperties{
			Intent: &armcompute.ProximityPlacementGroupPropertiesIntent{
				VMSizes: []*string{
					to.Ptr("Basic_A0"),
					to.Ptr("Basic_A2")},
			},
			ProximityPlacementGroupType: to.Ptr(armcompute.ProximityPlacementGroupTypeStandard),
		},
		Zones: []*string{
			to.Ptr("1")},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ProximityPlacementGroup = armcompute.ProximityPlacementGroup{
	// 	Name: to.Ptr("myProximityPlacementGroup"),
	// 	Type: to.Ptr("Microsoft.Compute/proximityPlacementGroups"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/proximityPlacementGroups/myProximityPlacementGroup"),
	// 	Location: to.Ptr("westus"),
	// 	Properties: &armcompute.ProximityPlacementGroupProperties{
	// 		Intent: &armcompute.ProximityPlacementGroupPropertiesIntent{
	// 			VMSizes: []*string{
	// 				to.Ptr("Basic_A0"),
	// 				to.Ptr("Basic_A2")},
	// 			},
	// 			ProximityPlacementGroupType: to.Ptr(armcompute.ProximityPlacementGroupTypeStandard),
	// 		},
	// 		Zones: []*string{
	// 			to.Ptr("1")},
	// 		}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/proximityPlacementGroupExamples/ProximityPlacementGroup_Patch.json
func ExampleProximityPlacementGroupsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewProximityPlacementGroupsClient().Update(ctx, "myResourceGroup", "myProximityPlacementGroup", armcompute.ProximityPlacementGroupUpdate{
		Tags: map[string]*string{
			"additionalProp1": to.Ptr("string"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ProximityPlacementGroup = armcompute.ProximityPlacementGroup{
	// 	Name: to.Ptr("myProximityPlacementGroup"),
	// 	Type: to.Ptr("Microsoft.Compute/proximityPlacementGroups"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/proximityPlacementGroups/myProximityPlacementGroup"),
	// 	Location: to.Ptr("westus"),
	// 	Properties: &armcompute.ProximityPlacementGroupProperties{
	// 		ProximityPlacementGroupType: to.Ptr(armcompute.ProximityPlacementGroupTypeStandard),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/proximityPlacementGroupExamples/ProximityPlacementGroup_Delete.json
func ExampleProximityPlacementGroupsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewProximityPlacementGroupsClient().Delete(ctx, "myResourceGroup", "myProximityPlacementGroup", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/proximityPlacementGroupExamples/ProximityPlacementGroup_Get.json
func ExampleProximityPlacementGroupsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewProximityPlacementGroupsClient().Get(ctx, "myResourceGroup", "myProximityPlacementGroup", &armcompute.ProximityPlacementGroupsClientGetOptions{IncludeColocationStatus: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ProximityPlacementGroup = armcompute.ProximityPlacementGroup{
	// 	Name: to.Ptr("myProximityPlacementGroup"),
	// 	Type: to.Ptr("Microsoft.Compute/proximityPlacementGroups"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/proximityPlacementGroups/myProximityPlacementGroup"),
	// 	Location: to.Ptr("westus"),
	// 	Properties: &armcompute.ProximityPlacementGroupProperties{
	// 		AvailabilitySets: []*armcompute.SubResourceWithColocationStatus{
	// 			{
	// 				ID: to.Ptr("string"),
	// 		}},
	// 		Intent: &armcompute.ProximityPlacementGroupPropertiesIntent{
	// 			VMSizes: []*string{
	// 				to.Ptr("Basic_A0"),
	// 				to.Ptr("Basic_A2")},
	// 			},
	// 			ProximityPlacementGroupType: to.Ptr(armcompute.ProximityPlacementGroupTypeStandard),
	// 			VirtualMachineScaleSets: []*armcompute.SubResourceWithColocationStatus{
	// 				{
	// 					ID: to.Ptr("string"),
	// 			}},
	// 			VirtualMachines: []*armcompute.SubResourceWithColocationStatus{
	// 				{
	// 					ID: to.Ptr("string"),
	// 			}},
	// 		},
	// 		Zones: []*string{
	// 			to.Ptr("1")},
	// 		}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/proximityPlacementGroupExamples/ProximityPlacementGroup_ListBySubscription.json
func ExampleProximityPlacementGroupsClient_NewListBySubscriptionPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewProximityPlacementGroupsClient().NewListBySubscriptionPager(nil)
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
		// page.ProximityPlacementGroupListResult = armcompute.ProximityPlacementGroupListResult{
		// 	Value: []*armcompute.ProximityPlacementGroup{
		// 		{
		// 			Name: to.Ptr("myProximityPlacementGroup"),
		// 			Type: to.Ptr("Microsoft.Compute/proximityPlacementGroups"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/proximityPlacementGroups/myProximityPlacementGroup"),
		// 			Location: to.Ptr("westus"),
		// 			Properties: &armcompute.ProximityPlacementGroupProperties{
		// 				AvailabilitySets: []*armcompute.SubResourceWithColocationStatus{
		// 					{
		// 						ID: to.Ptr("string"),
		// 				}},
		// 				Intent: &armcompute.ProximityPlacementGroupPropertiesIntent{
		// 					VMSizes: []*string{
		// 						to.Ptr("Basic_A0"),
		// 						to.Ptr("Basic_A2")},
		// 					},
		// 					ProximityPlacementGroupType: to.Ptr(armcompute.ProximityPlacementGroupTypeStandard),
		// 					VirtualMachineScaleSets: []*armcompute.SubResourceWithColocationStatus{
		// 						{
		// 							ID: to.Ptr("string"),
		// 					}},
		// 					VirtualMachines: []*armcompute.SubResourceWithColocationStatus{
		// 						{
		// 							ID: to.Ptr("string"),
		// 					}},
		// 				},
		// 				Zones: []*string{
		// 					to.Ptr("1")},
		// 			}},
		// 		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/proximityPlacementGroupExamples/ProximityPlacementGroup_ListByResourceGroup.json
func ExampleProximityPlacementGroupsClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewProximityPlacementGroupsClient().NewListByResourceGroupPager("myResourceGroup", nil)
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
		// page.ProximityPlacementGroupListResult = armcompute.ProximityPlacementGroupListResult{
		// 	Value: []*armcompute.ProximityPlacementGroup{
		// 		{
		// 			Name: to.Ptr("myProximityPlacementGroup"),
		// 			Type: to.Ptr("Microsoft.Compute/proximityPlacementGroups"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/proximityPlacementGroups/myProximityPlacementGroup"),
		// 			Location: to.Ptr("westus"),
		// 			Properties: &armcompute.ProximityPlacementGroupProperties{
		// 				AvailabilitySets: []*armcompute.SubResourceWithColocationStatus{
		// 					{
		// 						ID: to.Ptr("string"),
		// 				}},
		// 				Intent: &armcompute.ProximityPlacementGroupPropertiesIntent{
		// 					VMSizes: []*string{
		// 						to.Ptr("Basic_A0"),
		// 						to.Ptr("Basic_A2")},
		// 					},
		// 					ProximityPlacementGroupType: to.Ptr(armcompute.ProximityPlacementGroupTypeStandard),
		// 					VirtualMachineScaleSets: []*armcompute.SubResourceWithColocationStatus{
		// 						{
		// 							ID: to.Ptr("string"),
		// 					}},
		// 					VirtualMachines: []*armcompute.SubResourceWithColocationStatus{
		// 						{
		// 							ID: to.Ptr("string"),
		// 					}},
		// 				},
		// 				Zones: []*string{
		// 					to.Ptr("1")},
		// 			}},
		// 		}
	}
}
