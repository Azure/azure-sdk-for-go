// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armcontainerservicefleet_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservicefleet/armcontainerservicefleet/v2"
	"log"
)

// Generated from example definition: 2025-03-01/FleetMembers_Create.json
func ExampleFleetMembersClient_BeginCreate_createsAFleetMemberResourceWithALongRunningOperation() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservicefleet.NewClientFactory("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewFleetMembersClient().BeginCreate(ctx, "rg1", "fleet1", "member-1", armcontainerservicefleet.FleetMember{
		Properties: &armcontainerservicefleet.FleetMemberProperties{
			ClusterResourceID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/cluster-1"),
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
	// res = armcontainerservicefleet.FleetMembersClientCreateResponse{
	// 	FleetMember: &armcontainerservicefleet.FleetMember{
	// 		ID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/fleets/fleet-1/members/member-1"),
	// 		Name: to.Ptr("member-1"),
	// 		Type: to.Ptr("Microsoft.ContainerService/fleets/members"),
	// 		SystemData: &armcontainerservicefleet.SystemData{
	// 			CreatedBy: to.Ptr("someUser"),
	// 			CreatedByType: to.Ptr(armcontainerservicefleet.CreatedByTypeUser),
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-23T05:40:40.657Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("someOtherUser"),
	// 			LastModifiedByType: to.Ptr(armcontainerservicefleet.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-23T05:40:40.657Z"); return t}()),
	// 		},
	// 		ETag: to.Ptr("23ujdflewrj3="),
	// 		Properties: &armcontainerservicefleet.FleetMemberProperties{
	// 			ClusterResourceID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/cluster-1"),
	// 			ProvisioningState: to.Ptr(armcontainerservicefleet.FleetMemberProvisioningStateSucceeded),
	// 			Status: &armcontainerservicefleet.FleetMemberStatus{
	// 				LastOperationID: to.Ptr("operation-12345"),
	// 				LastOperationError: &armcontainerservicefleet.ErrorDetail{
	// 					Code: to.Ptr("None"),
	// 					Message: to.Ptr("No error"),
	// 				},
	// 			},
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2025-03-01/FleetMembers_Create_MaximumSet_Gen.json
func ExampleFleetMembersClient_BeginCreate_createsAFleetMemberResourceWithALongRunningOperationGeneratedByMaximumSetRule() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservicefleet.NewClientFactory("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewFleetMembersClient().BeginCreate(ctx, "rgfleets", "fleet1", "fleet1", armcontainerservicefleet.FleetMember{
		Properties: &armcontainerservicefleet.FleetMemberProperties{
			ClusterResourceID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/cluster-1"),
			Group:             to.Ptr("fleet1"),
		},
	}, &FleetMembersClientBeginCreateOptions{
		ifMatch:     to.Ptr("amkttadbw"),
		ifNoneMatch: to.Ptr("zoljoccbcg")})
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
	// res = armcontainerservicefleet.FleetMembersClientCreateResponse{
	// 	FleetMember: &armcontainerservicefleet.FleetMember{
	// 		Properties: &armcontainerservicefleet.FleetMemberProperties{
	// 			ClusterResourceID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/cluster-1"),
	// 			Group: to.Ptr("fleet1"),
	// 			ProvisioningState: to.Ptr(armcontainerservicefleet.FleetMemberProvisioningStateSucceeded),
	// 		},
	// 		ETag: to.Ptr("kd30rkdfo49="),
	// 		ID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/fleets/fleet-1/members/member-1"),
	// 		Name: to.Ptr("member-1"),
	// 		Type: to.Ptr("Microsoft.ContainerService/fleets/members"),
	// 		SystemData: &armcontainerservicefleet.SystemData{
	// 			CreatedBy: to.Ptr("someUser"),
	// 			CreatedByType: to.Ptr(armcontainerservicefleet.CreatedByTypeUser),
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-23T05:40:40.657Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("someOtherUser"),
	// 			LastModifiedByType: to.Ptr(armcontainerservicefleet.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-23T05:40:40.657Z"); return t}()),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2025-03-01/FleetMembers_Delete.json
func ExampleFleetMembersClient_BeginDelete_deletesAFleetMemberResourceAsynchronouslyWithALongRunningOperation() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservicefleet.NewClientFactory("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewFleetMembersClient().BeginDelete(ctx, "rg1", "fleet1", "member-1", nil)
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
	// res = armcontainerservicefleet.FleetMembersClientDeleteResponse{
	// }
}

// Generated from example definition: 2025-03-01/FleetMembers_Delete_MaximumSet_Gen.json
func ExampleFleetMembersClient_BeginDelete_deletesAFleetMemberResourceAsynchronouslyWithALongRunningOperationGeneratedByMaximumSetRule() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservicefleet.NewClientFactory("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewFleetMembersClient().BeginDelete(ctx, "rgfleets", "fleet1", "fleet1", &FleetMembersClientBeginDeleteOptions{
		ifMatch: to.Ptr("klroqfozx")})
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
	// res = armcontainerservicefleet.FleetMembersClientDeleteResponse{
	// }
}

// Generated from example definition: 2025-03-01/FleetMembers_Get.json
func ExampleFleetMembersClient_Get_getsAFleetMemberResource() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservicefleet.NewClientFactory("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewFleetMembersClient().Get(ctx, "rg1", "fleet1", "member-1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armcontainerservicefleet.FleetMembersClientGetResponse{
	// 	FleetMember: &armcontainerservicefleet.FleetMember{
	// 		ID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/fleets/fleet-1/members/member-1"),
	// 		Name: to.Ptr("member-1"),
	// 		Type: to.Ptr("Microsoft.ContainerService/fleets/members"),
	// 		SystemData: &armcontainerservicefleet.SystemData{
	// 			CreatedBy: to.Ptr("someUser"),
	// 			CreatedByType: to.Ptr(armcontainerservicefleet.CreatedByTypeUser),
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-23T05:40:40.657Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("someOtherUser"),
	// 			LastModifiedByType: to.Ptr(armcontainerservicefleet.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-23T05:40:40.657Z"); return t}()),
	// 		},
	// 		ETag: to.Ptr("kd30rkdfo49="),
	// 		Properties: &armcontainerservicefleet.FleetMemberProperties{
	// 			ClusterResourceID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/cluster-1"),
	// 			ProvisioningState: to.Ptr(armcontainerservicefleet.FleetMemberProvisioningStateSucceeded),
	// 			Status: &armcontainerservicefleet.FleetMemberStatus{
	// 				LastOperationID: to.Ptr("operation-12345"),
	// 				LastOperationError: &armcontainerservicefleet.ErrorDetail{
	// 					Code: to.Ptr("None"),
	// 					Message: to.Ptr("No error"),
	// 				},
	// 			},
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2025-03-01/FleetMembers_Get_MaximumSet_Gen.json
func ExampleFleetMembersClient_Get_getsAFleetMemberResourceGeneratedByMaximumSetRule() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservicefleet.NewClientFactory("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewFleetMembersClient().Get(ctx, "rgfleets", "fleet1", "fleet1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armcontainerservicefleet.FleetMembersClientGetResponse{
	// 	FleetMember: &armcontainerservicefleet.FleetMember{
	// 		Properties: &armcontainerservicefleet.FleetMemberProperties{
	// 			ClusterResourceID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/cluster-1"),
	// 			Group: to.Ptr("fleet1"),
	// 			ProvisioningState: to.Ptr(armcontainerservicefleet.FleetMemberProvisioningStateSucceeded),
	// 		},
	// 		ETag: to.Ptr("kd30rkdfo49="),
	// 		ID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/fleets/fleet-1/members/member-1"),
	// 		Name: to.Ptr("member-1"),
	// 		Type: to.Ptr("Microsoft.ContainerService/fleets/members"),
	// 		SystemData: &armcontainerservicefleet.SystemData{
	// 			CreatedBy: to.Ptr("someUser"),
	// 			CreatedByType: to.Ptr(armcontainerservicefleet.CreatedByTypeUser),
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-23T05:40:40.657Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("someOtherUser"),
	// 			LastModifiedByType: to.Ptr(armcontainerservicefleet.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-23T05:40:40.657Z"); return t}()),
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2025-03-01/FleetMembers_ListByFleet.json
func ExampleFleetMembersClient_NewListByFleetPager_listsTheMembersOfAFleet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservicefleet.NewClientFactory("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewFleetMembersClient().NewListByFleetPager("rg1", "fleet1", nil)
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
		// page = armcontainerservicefleet.FleetMembersClientListByFleetResponse{
		// 	FleetMemberListResult: armcontainerservicefleet.FleetMemberListResult{
		// 		Value: []*armcontainerservicefleet.FleetMember{
		// 			{
		// 				ID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/fleets/fleet-1/members/member-1"),
		// 				Name: to.Ptr("member-1"),
		// 				Type: to.Ptr("Microsoft.ContainerService/fleets/members"),
		// 				SystemData: &armcontainerservicefleet.SystemData{
		// 					CreatedBy: to.Ptr("someUser"),
		// 					CreatedByType: to.Ptr(armcontainerservicefleet.CreatedByTypeUser),
		// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-23T05:40:40.657Z"); return t}()),
		// 					LastModifiedBy: to.Ptr("someOtherUser"),
		// 					LastModifiedByType: to.Ptr(armcontainerservicefleet.CreatedByTypeUser),
		// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-23T05:40:40.657Z"); return t}()),
		// 				},
		// 				ETag: to.Ptr("kd30rkdfo49="),
		// 				Properties: &armcontainerservicefleet.FleetMemberProperties{
		// 					ClusterResourceID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/cluster-1"),
		// 					ProvisioningState: to.Ptr(armcontainerservicefleet.FleetMemberProvisioningStateSucceeded),
		// 					Status: &armcontainerservicefleet.FleetMemberStatus{
		// 						LastOperationID: to.Ptr("operation-12345"),
		// 						LastOperationError: &armcontainerservicefleet.ErrorDetail{
		// 							Code: to.Ptr("None"),
		// 							Message: to.Ptr("No error"),
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// }
	}
}

// Generated from example definition: 2025-03-01/FleetMembers_ListByFleet_MaximumSet_Gen.json
func ExampleFleetMembersClient_NewListByFleetPager_listsTheMembersOfAFleetGeneratedByMaximumSetRule() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservicefleet.NewClientFactory("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewFleetMembersClient().NewListByFleetPager("rgfleets", "fleet1", nil)
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
		// page = armcontainerservicefleet.FleetMembersClientListByFleetResponse{
		// 	FleetMemberListResult: armcontainerservicefleet.FleetMemberListResult{
		// 		Value: []*armcontainerservicefleet.FleetMember{
		// 			{
		// 				ID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/fleets/fleet-1/members/member-1"),
		// 				Name: to.Ptr("member-1"),
		// 				Type: to.Ptr("Microsoft.ContainerService/fleets/members"),
		// 				SystemData: &armcontainerservicefleet.SystemData{
		// 					CreatedBy: to.Ptr("someUser"),
		// 					CreatedByType: to.Ptr(armcontainerservicefleet.CreatedByTypeUser),
		// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-23T05:40:40.657Z"); return t}()),
		// 					LastModifiedBy: to.Ptr("someOtherUser"),
		// 					LastModifiedByType: to.Ptr(armcontainerservicefleet.CreatedByTypeUser),
		// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-23T05:40:40.657Z"); return t}()),
		// 				},
		// 				ETag: to.Ptr("kd30rkdfo49="),
		// 				Properties: &armcontainerservicefleet.FleetMemberProperties{
		// 					ClusterResourceID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/cluster-1"),
		// 					ProvisioningState: to.Ptr(armcontainerservicefleet.FleetMemberProvisioningStateSucceeded),
		// 					Group: to.Ptr("fleet1"),
		// 				},
		// 			},
		// 		},
		// 		NextLink: to.Ptr("https://microsoft.com/a"),
		// 	},
		// }
	}
}

// Generated from example definition: 2025-03-01/FleetMembers_Update.json
func ExampleFleetMembersClient_BeginUpdateAsync_updatesAFleetMemberResourceSynchronously() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservicefleet.NewClientFactory("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewFleetMembersClient().BeginUpdateAsync(ctx, "rg1", "fleet1", "member-1", armcontainerservicefleet.FleetMemberUpdate{
		Properties: &armcontainerservicefleet.FleetMemberUpdateProperties{
			Group: to.Ptr("staging"),
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
	// res = armcontainerservicefleet.FleetMembersClientUpdateAsyncResponse{
	// 	FleetMember: &armcontainerservicefleet.FleetMember{
	// 		ID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/fleets/fleet-1/members/member-1"),
	// 		Name: to.Ptr("member-1"),
	// 		Type: to.Ptr("Microsoft.ContainerService/fleets/members"),
	// 		SystemData: &armcontainerservicefleet.SystemData{
	// 			CreatedBy: to.Ptr("someUser"),
	// 			CreatedByType: to.Ptr(armcontainerservicefleet.CreatedByTypeUser),
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-23T05:40:40.657Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("someOtherUser"),
	// 			LastModifiedByType: to.Ptr(armcontainerservicefleet.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-23T05:40:40.657Z"); return t}()),
	// 		},
	// 		ETag: to.Ptr("23ujdflewrj3="),
	// 		Properties: &armcontainerservicefleet.FleetMemberProperties{
	// 			ClusterResourceID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/cluster-1"),
	// 			Group: to.Ptr("staging"),
	// 			ProvisioningState: to.Ptr(armcontainerservicefleet.FleetMemberProvisioningStateSucceeded),
	// 			Status: &armcontainerservicefleet.FleetMemberStatus{
	// 				LastOperationID: to.Ptr("operation-12345"),
	// 				LastOperationError: &armcontainerservicefleet.ErrorDetail{
	// 					Code: to.Ptr("None"),
	// 					Message: to.Ptr("No error"),
	// 				},
	// 			},
	// 		},
	// 	},
	// }
}

// Generated from example definition: 2025-03-01/FleetMembers_Update_MaximumSet_Gen.json
func ExampleFleetMembersClient_BeginUpdateAsync_updatesAFleetMemberResourceSynchronouslyGeneratedByMaximumSetRule() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservicefleet.NewClientFactory("00000000-0000-0000-0000-000000000000", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewFleetMembersClient().BeginUpdateAsync(ctx, "rgfleets", "fleet1", "fleet1", armcontainerservicefleet.FleetMemberUpdate{
		Properties: &armcontainerservicefleet.FleetMemberUpdateProperties{
			Group: to.Ptr("staging"),
		},
	}, &FleetMembersClientBeginUpdateAsyncOptions{
		ifMatch: to.Ptr("bjyjzzxvbs")})
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
	// res = armcontainerservicefleet.FleetMembersClientUpdateAsyncResponse{
	// 	FleetMember: &armcontainerservicefleet.FleetMember{
	// 		Properties: &armcontainerservicefleet.FleetMemberProperties{
	// 			ClusterResourceID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/cluster-1"),
	// 			Group: to.Ptr("fleet1"),
	// 			ProvisioningState: to.Ptr(armcontainerservicefleet.FleetMemberProvisioningStateSucceeded),
	// 		},
	// 		ETag: to.Ptr("kd30rkdfo49="),
	// 		ID: to.Ptr("/subscriptions/subid1/resourcegroups/rg1/providers/Microsoft.ContainerService/fleets/fleet-1/members/member-1"),
	// 		Name: to.Ptr("member-1"),
	// 		Type: to.Ptr("Microsoft.ContainerService/fleets/members"),
	// 		SystemData: &armcontainerservicefleet.SystemData{
	// 			CreatedBy: to.Ptr("someUser"),
	// 			CreatedByType: to.Ptr(armcontainerservicefleet.CreatedByTypeUser),
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-23T05:40:40.657Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("someOtherUser"),
	// 			LastModifiedByType: to.Ptr(armcontainerservicefleet.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-03-23T05:40:40.657Z"); return t}()),
	// 		},
	// 	},
	// }
}
