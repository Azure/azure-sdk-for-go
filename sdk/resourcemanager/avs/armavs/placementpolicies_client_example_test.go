// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armavs_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/avs/armavs/v2"
	"log"
)

// Generated from example definition: 2024-09-01/PlacementPolicies_CreateOrUpdate.json
func ExamplePlacementPoliciesClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armavs.NewClientFactory("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewPlacementPoliciesClient().BeginCreateOrUpdate(ctx, "group1", "cloud1", "cluster1", "policy1", armavs.PlacementPolicy{
		Properties: &armavs.VMHostPlacementPolicyProperties{
			Type: to.Ptr(armavs.PlacementPolicyTypeVMHost),
			VMMembers: []*string{
				to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-128"),
				to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-256"),
			},
			HostMembers: []*string{
				to.Ptr("fakehost22.nyc1.kubernetes.center"),
				to.Ptr("fakehost23.nyc1.kubernetes.center"),
				to.Ptr("fakehost24.nyc1.kubernetes.center"),
			},
			AffinityType:           to.Ptr(armavs.AffinityTypeAntiAffinity),
			AffinityStrength:       to.Ptr(armavs.AffinityStrengthMust),
			AzureHybridBenefitType: to.Ptr(armavs.AzureHybridBenefitTypeSQLHost),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armavs.PlacementPoliciesClientCreateOrUpdateResponse{
	// 	PlacementPolicy: &armavs.PlacementPolicy{
	// 		ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/placementPolicies/policy1"),
	// 		Name: to.Ptr("policy1"),
	// 		Type: to.Ptr("Microsoft.AVS/privateClouds/clusters/placementPolicies"),
	// 		Properties: &armavs.VMHostPlacementPolicyProperties{
	// 			DisplayName: to.Ptr("policy1"),
	// 			Type: to.Ptr(armavs.PlacementPolicyTypeVMHost),
	// 			State: to.Ptr(armavs.PlacementPolicyStateEnabled),
	// 			VMMembers: []*string{
	// 				to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-128"),
	// 				to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-256"),
	// 			},
	// 			HostMembers: []*string{
	// 				to.Ptr("fakehost22.nyc1.kubernetes.center"),
	// 				to.Ptr("fakehost23.nyc1.kubernetes.center"),
	// 				to.Ptr("fakehost24.nyc1.kubernetes.center"),
	// 			},
	// 			AffinityType: to.Ptr(armavs.AffinityTypeAntiAffinity),
	// 			AffinityStrength: to.Ptr(armavs.AffinityStrengthMust),
	// 			AzureHybridBenefitType: to.Ptr(armavs.AzureHybridBenefitTypeSQLHost),
	// 			ProvisioningState: to.Ptr(armavs.PlacementPolicyProvisioningStateSucceeded),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2024-09-01/PlacementPolicies_Delete.json
func ExamplePlacementPoliciesClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armavs.NewClientFactory("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewPlacementPoliciesClient().BeginDelete(ctx, "group1", "cloud1", "cluster1", "policy1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armavs.PlacementPoliciesClientDeleteResponse{
	// }
}

// Generated from example definition: 2024-09-01/PlacementPolicies_Get.json
func ExamplePlacementPoliciesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armavs.NewClientFactory("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewPlacementPoliciesClient().Get(ctx, "group1", "cloud1", "cluster1", "policy1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armavs.PlacementPoliciesClientGetResponse{
	// 	PlacementPolicy: &armavs.PlacementPolicy{
	// 		ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/placementPolicies/policy1"),
	// 		Name: to.Ptr("policy1"),
	// 		Type: to.Ptr("Microsoft.AVS/privateClouds/clusters/placementPolicies"),
	// 		Properties: &armavs.VMHostPlacementPolicyProperties{
	// 			DisplayName: to.Ptr("policy1"),
	// 			Type: to.Ptr(armavs.PlacementPolicyTypeVMHost),
	// 			State: to.Ptr(armavs.PlacementPolicyStateEnabled),
	// 			VMMembers: []*string{
	// 				to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-128"),
	// 				to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-256"),
	// 			},
	// 			HostMembers: []*string{
	// 				to.Ptr("fakehost22.nyc1.kubernetes.center"),
	// 				to.Ptr("fakehost23.nyc1.kubernetes.center"),
	// 				to.Ptr("fakehost24.nyc1.kubernetes.center"),
	// 			},
	// 			AffinityType: to.Ptr(armavs.AffinityTypeAntiAffinity),
	// 			AffinityStrength: to.Ptr(armavs.AffinityStrengthMust),
	// 			AzureHybridBenefitType: to.Ptr(armavs.AzureHybridBenefitTypeSQLHost),
	// 			ProvisioningState: to.Ptr(armavs.PlacementPolicyProvisioningStateSucceeded),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2024-09-01/PlacementPolicies_List.json
func ExamplePlacementPoliciesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armavs.NewClientFactory("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewPlacementPoliciesClient().NewListPager("group1", "cloud1", "cluster1", nil)
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
		// page = armavs.PlacementPoliciesClientListResponse{
		// 	PlacementPoliciesList: armavs.PlacementPoliciesList{
		// 		Value: []*armavs.PlacementPolicy{
		// 			{
		// 				ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/placementPolicies/policy1"),
		// 				Name: to.Ptr("policy1"),
		// 				Type: to.Ptr("Microsoft.AVS/privateClouds/clusters/placementPolicies"),
		// 				Properties: &armavs.VMHostPlacementPolicyProperties{
		// 					DisplayName: to.Ptr("policy1"),
		// 					Type: to.Ptr(armavs.PlacementPolicyTypeVMHost),
		// 					State: to.Ptr(armavs.PlacementPolicyStateEnabled),
		// 					VMMembers: []*string{
		// 						to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-128"),
		// 						to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-256"),
		// 					},
		// 					HostMembers: []*string{
		// 						to.Ptr("fakehost22.nyc1.kubernetes.center"),
		// 						to.Ptr("fakehost23.nyc1.kubernetes.center"),
		// 						to.Ptr("fakehost24.nyc1.kubernetes.center"),
		// 					},
		// 					AffinityType: to.Ptr(armavs.AffinityTypeAntiAffinity),
		// 					AffinityStrength: to.Ptr(armavs.AffinityStrengthMust),
		// 					AzureHybridBenefitType: to.Ptr(armavs.AzureHybridBenefitTypeSQLHost),
		// 					ProvisioningState: to.Ptr(armavs.PlacementPolicyProvisioningStateSucceeded),
		// 				},
		// 			},
		// 			{
		// 				ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/placementPolicies/policy2"),
		// 				Name: to.Ptr("policy2"),
		// 				Type: to.Ptr("Microsoft.AVS/privateClouds/clusters/placementPolicies"),
		// 				Properties: &armavs.VMPlacementPolicyProperties{
		// 					DisplayName: to.Ptr("policy2"),
		// 					Type: to.Ptr(armavs.PlacementPolicyTypeVMVM),
		// 					State: to.Ptr(armavs.PlacementPolicyStateEnabled),
		// 					VMMembers: []*string{
		// 						to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-128"),
		// 						to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-256"),
		// 					},
		// 					AffinityType: to.Ptr(armavs.AffinityTypeAffinity),
		// 					ProvisioningState: to.Ptr(armavs.PlacementPolicyProvisioningStateSucceeded),
		// 				},
		// 			},
		// 		},
		// 	},
		// }
	}
}

// Generated from example definition: 2024-09-01/PlacementPolicies_Update.json
func ExamplePlacementPoliciesClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armavs.NewClientFactory("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewPlacementPoliciesClient().BeginUpdate(ctx, "group1", "cloud1", "cluster1", "policy1", armavs.PlacementPolicyUpdate{
		Properties: &armavs.PlacementPolicyUpdateProperties{
			State: to.Ptr(armavs.PlacementPolicyStateDisabled),
			VMMembers: []*string{
				to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-128"),
				to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-256"),
			},
			HostMembers: []*string{
				to.Ptr("fakehost22.nyc1.kubernetes.center"),
				to.Ptr("fakehost23.nyc1.kubernetes.center"),
				to.Ptr("fakehost24.nyc1.kubernetes.center"),
			},
			AffinityStrength:       to.Ptr(armavs.AffinityStrengthMust),
			AzureHybridBenefitType: to.Ptr(armavs.AzureHybridBenefitTypeSQLHost),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armavs.PlacementPoliciesClientUpdateResponse{
	// 	PlacementPolicy: &armavs.PlacementPolicy{
	// 		ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/placementPolicies/policy1"),
	// 		Name: to.Ptr("policy1"),
	// 		Type: to.Ptr("Microsoft.AVS/privateClouds/clusters/placementPolicies"),
	// 		Properties: &armavs.VMHostPlacementPolicyProperties{
	// 			DisplayName: to.Ptr("policy1"),
	// 			Type: to.Ptr(armavs.PlacementPolicyTypeVMHost),
	// 			State: to.Ptr(armavs.PlacementPolicyStateDisabled),
	// 			VMMembers: []*string{
	// 				to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-128"),
	// 				to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-256"),
	// 			},
	// 			HostMembers: []*string{
	// 				to.Ptr("fakehost22.nyc1.kubernetes.center"),
	// 				to.Ptr("fakehost23.nyc1.kubernetes.center"),
	// 				to.Ptr("fakehost24.nyc1.kubernetes.center"),
	// 			},
	// 			AffinityType: to.Ptr(armavs.AffinityTypeAntiAffinity),
	// 			AffinityStrength: to.Ptr(armavs.AffinityStrengthMust),
	// 			AzureHybridBenefitType: to.Ptr(armavs.AzureHybridBenefitTypeSQLHost),
	// 			ProvisioningState: to.Ptr(armavs.PlacementPolicyProvisioningStateSucceeded),
	// 		},
	// 	},
	// }
}
