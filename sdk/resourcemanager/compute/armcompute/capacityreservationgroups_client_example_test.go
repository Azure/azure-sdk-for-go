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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/capacityReservationExamples/CapacityReservationGroup_CreateOrUpdate.json
func ExampleCapacityReservationGroupsClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewCapacityReservationGroupsClient().CreateOrUpdate(ctx, "myResourceGroup", "myCapacityReservationGroup", armcompute.CapacityReservationGroup{
		Location: to.Ptr("westus"),
		Tags: map[string]*string{
			"department": to.Ptr("finance"),
		},
		Properties: &armcompute.CapacityReservationGroupProperties{
			SharingProfile: &armcompute.ResourceSharingProfile{
				SubscriptionIDs: []*armcompute.SubResource{
					{
						ID: to.Ptr("/subscriptions/{subscription-id1}"),
					},
					{
						ID: to.Ptr("/subscriptions/{subscription-id2}"),
					}},
			},
		},
		Zones: []*string{
			to.Ptr("1"),
			to.Ptr("2")},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.CapacityReservationGroup = armcompute.CapacityReservationGroup{
	// 	Name: to.Ptr("myCapacityReservationGroup"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/CapacityReservationGroups/myCapacityReservationGroup"),
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 		"department": to.Ptr("finance"),
	// 		"owner": to.Ptr("myCompany"),
	// 	},
	// 	Properties: &armcompute.CapacityReservationGroupProperties{
	// 		SharingProfile: &armcompute.ResourceSharingProfile{
	// 			SubscriptionIDs: []*armcompute.SubResource{
	// 				{
	// 					ID: to.Ptr("/subscriptions/{subscription-id1}"),
	// 				},
	// 				{
	// 					ID: to.Ptr("/subscriptions/{subscription-id2}"),
	// 			}},
	// 		},
	// 	},
	// 	Zones: []*string{
	// 		to.Ptr("1"),
	// 		to.Ptr("2")},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/capacityReservationExamples/CapacityReservationGroup_Update_MaximumSet_Gen.json
func ExampleCapacityReservationGroupsClient_Update_capacityReservationGroupUpdateMaximumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewCapacityReservationGroupsClient().Update(ctx, "rgcompute", "aaaaaaaaaaaaaaaaaaaaaa", armcompute.CapacityReservationGroupUpdate{
		Tags: map[string]*string{
			"key5355": to.Ptr("aaa"),
		},
		Properties: &armcompute.CapacityReservationGroupProperties{
			InstanceView: &armcompute.CapacityReservationGroupInstanceView{},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.CapacityReservationGroup = armcompute.CapacityReservationGroup{
	// 	Name: to.Ptr("myCapacityReservationGroup"),
	// 	Type: to.Ptr("aaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/CapacityReservationGroups/myCapacityReservationGroup"),
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 	},
	// 	Properties: &armcompute.CapacityReservationGroupProperties{
	// 		CapacityReservations: []*armcompute.SubResourceReadOnly{
	// 			{
	// 				ID: to.Ptr("aaaa"),
	// 		}},
	// 		InstanceView: &armcompute.CapacityReservationGroupInstanceView{
	// 			CapacityReservations: []*armcompute.CapacityReservationInstanceViewWithName{
	// 				{
	// 					Statuses: []*armcompute.InstanceViewStatus{
	// 						{
	// 							Code: to.Ptr("aaaaaaaaaaaaaaaaaaaaaaa"),
	// 							DisplayStatus: to.Ptr("aaaaaa"),
	// 							Level: to.Ptr(armcompute.StatusLevelTypesInfo),
	// 							Message: to.Ptr("a"),
	// 							Time: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-11-30T12:58:26.522Z"); return t}()),
	// 					}},
	// 					UtilizationInfo: &armcompute.CapacityReservationUtilization{
	// 						VirtualMachinesAllocated: []*armcompute.SubResourceReadOnly{
	// 							{
	// 								ID: to.Ptr("aaaa"),
	// 						}},
	// 					},
	// 					Name: to.Ptr("aaaaaaaaaaaaaaaa"),
	// 			}},
	// 		},
	// 		VirtualMachinesAssociated: []*armcompute.SubResourceReadOnly{
	// 			{
	// 				ID: to.Ptr("aaaa"),
	// 		}},
	// 	},
	// 	Zones: []*string{
	// 		to.Ptr("1"),
	// 		to.Ptr("2")},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/capacityReservationExamples/CapacityReservationGroup_Update_MinimumSet_Gen.json
func ExampleCapacityReservationGroupsClient_Update_capacityReservationGroupUpdateMinimumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewCapacityReservationGroupsClient().Update(ctx, "rgcompute", "aaaaaaaaaaaaaaaaaaaaaa", armcompute.CapacityReservationGroupUpdate{}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.CapacityReservationGroup = armcompute.CapacityReservationGroup{
	// 	Location: to.Ptr("westus"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/capacityReservationExamples/CapacityReservationGroup_Delete_MaximumSet_Gen.json
func ExampleCapacityReservationGroupsClient_Delete_capacityReservationGroupDeleteMaximumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewCapacityReservationGroupsClient().Delete(ctx, "rgcompute", "a", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/capacityReservationExamples/CapacityReservationGroup_Delete_MinimumSet_Gen.json
func ExampleCapacityReservationGroupsClient_Delete_capacityReservationGroupDeleteMinimumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewCapacityReservationGroupsClient().Delete(ctx, "rgcompute", "aaaaaaaaaaaaaaaaaaaaaaaaaa", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/capacityReservationExamples/CapacityReservationGroup_Get.json
func ExampleCapacityReservationGroupsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewCapacityReservationGroupsClient().Get(ctx, "myResourceGroup", "myCapacityReservationGroup", &armcompute.CapacityReservationGroupsClientGetOptions{Expand: to.Ptr(armcompute.CapacityReservationGroupInstanceViewTypesInstanceView)})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.CapacityReservationGroup = armcompute.CapacityReservationGroup{
	// 	Name: to.Ptr("myCapacityReservationGroup"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/CapacityReservationGroups/myCapacityReservationGroup"),
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 		"{tagName}": to.Ptr("{tagValue}"),
	// 	},
	// 	Properties: &armcompute.CapacityReservationGroupProperties{
	// 		CapacityReservations: []*armcompute.SubResourceReadOnly{
	// 			{
	// 				ID: to.Ptr("subscriptions/{subscriptionId}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/capacityReservationGroups/myCapacityReservationGroup/capacityReservations/myCapacityReservation1"),
	// 			},
	// 			{
	// 				ID: to.Ptr("subscriptions/{subscriptionId}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/capacityReservationGroups/myCapacityReservationGroup/capacityReservations/myCapacityReservation2"),
	// 		}},
	// 		InstanceView: &armcompute.CapacityReservationGroupInstanceView{
	// 			CapacityReservations: []*armcompute.CapacityReservationInstanceViewWithName{
	// 				{
	// 					Statuses: []*armcompute.InstanceViewStatus{
	// 						{
	// 							Code: to.Ptr("ProvisioningState/succeeded"),
	// 							DisplayStatus: to.Ptr("Provisioning succeeded"),
	// 							Level: to.Ptr(armcompute.StatusLevelTypesInfo),
	// 					}},
	// 					UtilizationInfo: &armcompute.CapacityReservationUtilization{
	// 						CurrentCapacity: to.Ptr[int32](5),
	// 						VirtualMachinesAllocated: []*armcompute.SubResourceReadOnly{
	// 							{
	// 								ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM1"),
	// 							},
	// 							{
	// 								ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM2"),
	// 						}},
	// 					},
	// 					Name: to.Ptr("myCapacityReservation1"),
	// 				},
	// 				{
	// 					Statuses: []*armcompute.InstanceViewStatus{
	// 						{
	// 							Code: to.Ptr("ProvisioningState/succeeded"),
	// 							DisplayStatus: to.Ptr("Provisioning succeeded"),
	// 							Level: to.Ptr(armcompute.StatusLevelTypesInfo),
	// 					}},
	// 					UtilizationInfo: &armcompute.CapacityReservationUtilization{
	// 						CurrentCapacity: to.Ptr[int32](5),
	// 						VirtualMachinesAllocated: []*armcompute.SubResourceReadOnly{
	// 							{
	// 								ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM3"),
	// 							},
	// 							{
	// 								ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM4"),
	// 						}},
	// 					},
	// 					Name: to.Ptr("myCapacityReservation2"),
	// 			}},
	// 			SharedSubscriptionIDs: []*armcompute.SubResourceReadOnly{
	// 				{
	// 					ID: to.Ptr("/subscriptions/{subscription-id1}"),
	// 				},
	// 				{
	// 					ID: to.Ptr("/subscriptions/{subscription-id2}"),
	// 			}},
	// 		},
	// 		SharingProfile: &armcompute.ResourceSharingProfile{
	// 			SubscriptionIDs: []*armcompute.SubResource{
	// 				{
	// 					ID: to.Ptr("/subscriptions/{subscription-id1}"),
	// 				},
	// 				{
	// 					ID: to.Ptr("/subscriptions/{subscription-id2}"),
	// 			}},
	// 		},
	// 	},
	// 	Zones: []*string{
	// 		to.Ptr("3"),
	// 		to.Ptr("1")},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/capacityReservationExamples/CapacityReservationGroup_ListByResourceGroup.json
func ExampleCapacityReservationGroupsClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewCapacityReservationGroupsClient().NewListByResourceGroupPager("myResourceGroup", &armcompute.CapacityReservationGroupsClientListByResourceGroupOptions{Expand: to.Ptr(armcompute.ExpandTypesForGetCapacityReservationGroupsVirtualMachinesRef)})
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
		// page.CapacityReservationGroupListResult = armcompute.CapacityReservationGroupListResult{
		// 	Value: []*armcompute.CapacityReservationGroup{
		// 		{
		// 			Name: to.Ptr("{capacityReservationGroupName}"),
		// 			Type: to.Ptr("Microsoft.Compute/CapacityReservationGroups"),
		// 			ID: to.Ptr("/subscriptions/{subscriptionId}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/capacityReservationGroups/{capacityReservationGroupName}"),
		// 			Location: to.Ptr("West US"),
		// 			Properties: &armcompute.CapacityReservationGroupProperties{
		// 				CapacityReservations: []*armcompute.SubResourceReadOnly{
		// 					{
		// 						ID: to.Ptr("subscriptions/{subscriptionId}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/capacityReservationGroups/myCapacityReservationGroup/capacityReservations/myCapacityReservation1"),
		// 					},
		// 					{
		// 						ID: to.Ptr("subscriptions/{subscriptionId}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/capacityReservationGroups/myCapacityReservationGroup/capacityReservations/myCapacityReservation2"),
		// 				}},
		// 				VirtualMachinesAssociated: []*armcompute.SubResourceReadOnly{
		// 					{
		// 						ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM1"),
		// 					},
		// 					{
		// 						ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM2"),
		// 				}},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("{capacityReservationGroupName}"),
		// 			Type: to.Ptr("Microsoft.Compute/CapacityReservationGroups"),
		// 			ID: to.Ptr("/subscriptions/{subscriptionId}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/capacityReservationGroups/{capacityReservationGroupName}"),
		// 			Location: to.Ptr("West US"),
		// 			Properties: &armcompute.CapacityReservationGroupProperties{
		// 				CapacityReservations: []*armcompute.SubResourceReadOnly{
		// 					{
		// 						ID: to.Ptr("subscriptions/{subscriptionId}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/capacityReservationGroups/myCapacityReservationGroup/capacityReservations/myCapacityReservation3"),
		// 					},
		// 					{
		// 						ID: to.Ptr("subscriptions/{subscriptionId}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/capacityReservationGroups/myCapacityReservationGroup/capacityReservations/myCapacityReservation4"),
		// 				}},
		// 				VirtualMachinesAssociated: []*armcompute.SubResourceReadOnly{
		// 					{
		// 						ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM3"),
		// 				}},
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/capacityReservationExamples/CapacityReservationGroup_ListBySubscription.json
func ExampleCapacityReservationGroupsClient_NewListBySubscriptionPager_listCapacityReservationGroupsInSubscription() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewCapacityReservationGroupsClient().NewListBySubscriptionPager(&armcompute.CapacityReservationGroupsClientListBySubscriptionOptions{Expand: to.Ptr(armcompute.ExpandTypesForGetCapacityReservationGroupsVirtualMachinesRef),
		ResourceIDsOnly: nil,
	})
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
		// page.CapacityReservationGroupListResult = armcompute.CapacityReservationGroupListResult{
		// 	Value: []*armcompute.CapacityReservationGroup{
		// 		{
		// 			Name: to.Ptr("{capacityReservationGroupName}"),
		// 			Type: to.Ptr("Microsoft.Compute/CapacityReservationGroups"),
		// 			ID: to.Ptr("/subscriptions/{subscriptionId}/resourceGroups/myResourceGroup1/providers/Microsoft.Compute/capacityReservationGroups/{capacityReservationGroupName}"),
		// 			Location: to.Ptr("West US"),
		// 			Properties: &armcompute.CapacityReservationGroupProperties{
		// 				CapacityReservations: []*armcompute.SubResourceReadOnly{
		// 					{
		// 						ID: to.Ptr("subscriptions/{subscriptionId}/resourceGroups/myResourceGroup1/providers/Microsoft.Compute/capacityReservationGroups/myCapacityReservationGroup/capacityReservations/myCapacityReservation1"),
		// 					},
		// 					{
		// 						ID: to.Ptr("subscriptions/{subscriptionId}/resourceGroups/myResourceGroup1/providers/Microsoft.Compute/capacityReservationGroups/myCapacityReservationGroup/capacityReservations/myCapacityReservation2"),
		// 				}},
		// 				VirtualMachinesAssociated: []*armcompute.SubResourceReadOnly{
		// 					{
		// 						ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup1/providers/Microsoft.Compute/virtualMachines/myVM1"),
		// 					},
		// 					{
		// 						ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup1/providers/Microsoft.Compute/virtualMachines/myVM2"),
		// 				}},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("{capacityReservationGroupName}"),
		// 			Type: to.Ptr("Microsoft.Compute/CapacityReservationGroups"),
		// 			ID: to.Ptr("/subscriptions/{subscriptionId}/resourceGroups/myResourceGroup2/providers/Microsoft.Compute/capacityReservationGroups/{capacityReservationGroupName}"),
		// 			Location: to.Ptr("West US"),
		// 			Properties: &armcompute.CapacityReservationGroupProperties{
		// 				CapacityReservations: []*armcompute.SubResourceReadOnly{
		// 					{
		// 						ID: to.Ptr("subscriptions/{subscriptionId}/resourceGroups/myResourceGroup2/providers/Microsoft.Compute/capacityReservationGroups/myCapacityReservationGroup/capacityReservations/myCapacityReservation3"),
		// 					},
		// 					{
		// 						ID: to.Ptr("subscriptions/{subscriptionId}/resourceGroups/myResourceGroup2/providers/Microsoft.Compute/capacityReservationGroups/myCapacityReservationGroup/capacityReservations/myCapacityReservation4"),
		// 				}},
		// 				VirtualMachinesAssociated: []*armcompute.SubResourceReadOnly{
		// 					{
		// 						ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup2/providers/Microsoft.Compute/virtualMachines/myVM3"),
		// 				}},
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/069a65e8a6d1a6c0c58d9a9d97610b7103b6e8a5/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/capacityReservationExamples/CapacityReservationGroup_ListBySubscriptionWithResourceIdsQuery.json
func ExampleCapacityReservationGroupsClient_NewListBySubscriptionPager_listCapacityReservationGroupsWithResourceIdsOnlyInSubscription() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewCapacityReservationGroupsClient().NewListBySubscriptionPager(&armcompute.CapacityReservationGroupsClientListBySubscriptionOptions{Expand: nil,
		ResourceIDsOnly: to.Ptr(armcompute.ResourceIDOptionsForGetCapacityReservationGroupsAll),
	})
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
		// page.CapacityReservationGroupListResult = armcompute.CapacityReservationGroupListResult{
		// 	Value: []*armcompute.CapacityReservationGroup{
		// 		{
		// 			Type: to.Ptr("Microsoft.Compute/capacityReservationGroups"),
		// 			ID: to.Ptr("/subscriptions/{subscriptionId}/resourceGroups/myResourceGroup1/providers/Microsoft.Compute/capacityReservationGroups/{capacityReservationGroupName1}"),
		// 			Location: to.Ptr("SouthCentralUSSTG"),
		// 		},
		// 		{
		// 			Type: to.Ptr("Microsoft.Compute/capacityReservationGroups"),
		// 			ID: to.Ptr("/subscriptions/{subscriptionId2}/resourceGroups/myResourceGroup2/providers/Microsoft.Compute/capacityReservationGroups/{capacityReservationGroupName2}"),
		// 			Location: to.Ptr("SouthCentralUSSTG"),
		// 	}},
		// }
	}
}
