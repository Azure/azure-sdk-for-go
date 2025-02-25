//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armquota_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/quota/armquota"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/8691e5081766c7ad602a9e55de841d07bed5196a/specification/quota/resource-manager/Microsoft.Quota/stable/2025-03-01/examples/SubscriptionQuotaAllocationRequests/PatchSubscriptionQuotaAllocationRequest-Compute.json
func ExampleGroupQuotaSubscriptionAllocationRequestClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armquota.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGroupQuotaSubscriptionAllocationRequestClient().BeginUpdate(ctx, "E7EC67B3-7657-4966-BFFC-41EFD36BAA09", "groupquota1", "Microsoft.Compute", "westus", armquota.SubscriptionQuotaAllocationsList{
		Properties: &armquota.SubscriptionQuotaAllocationsListProperties{
			Value: []*armquota.SubscriptionQuotaAllocations{
				{
					Properties: &armquota.SubscriptionQuotaAllocationsProperties{
						Limit:        to.Ptr[int64](110),
						ResourceName: to.Ptr("standardddv4family"),
					},
				},
				{
					Properties: &armquota.SubscriptionQuotaAllocationsProperties{
						Limit:        to.Ptr[int64](110),
						ResourceName: to.Ptr("standardav2family"),
					},
				}},
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
	// res.SubscriptionQuotaAllocationsList = armquota.SubscriptionQuotaAllocationsList{
	// 	Name: to.Ptr("westus"),
	// 	Type: to.Ptr("Microsoft.Quota/groupQuotas/quotaAllocations"),
	// 	ID: to.Ptr("/providers/Microsoft.Management/managementGroups/E7EC67B3-7657-4966-BFFC-41EFD36BAA09/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Quota/groupQuotas/groupquota1/resourceProviders/Microsoft.Compute/quotaAllocations/westus"),
	// 	Properties: &armquota.SubscriptionQuotaAllocationsListProperties{
	// 		ProvisioningState: to.Ptr(armquota.RequestStateSucceeded),
	// 		Value: []*armquota.SubscriptionQuotaAllocations{
	// 			{
	// 				Properties: &armquota.SubscriptionQuotaAllocationsProperties{
	// 					Name: &armquota.SubscriptionQuotaDetailsName{
	// 						LocalizedValue: to.Ptr("standard DDv4 Family vCPUs"),
	// 						Value: to.Ptr("standardddv4family"),
	// 					},
	// 					Limit: to.Ptr[int64](25),
	// 					ResourceName: to.Ptr("standardddv4family"),
	// 					ShareableQuota: to.Ptr[int64](15),
	// 				},
	// 			},
	// 			{
	// 				Properties: &armquota.SubscriptionQuotaAllocationsProperties{
	// 					Name: &armquota.SubscriptionQuotaDetailsName{
	// 						LocalizedValue: to.Ptr("standard Av2 Family vCPUs"),
	// 						Value: to.Ptr("standardav2family"),
	// 					},
	// 					Limit: to.Ptr[int64](30),
	// 					ResourceName: to.Ptr("standardav2family"),
	// 					ShareableQuota: to.Ptr[int64](0),
	// 				},
	// 		}},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/8691e5081766c7ad602a9e55de841d07bed5196a/specification/quota/resource-manager/Microsoft.Quota/stable/2025-03-01/examples/SubscriptionQuotaAllocationRequests/SubscriptionQuotaAllocationRequests_Get-Compute.json
func ExampleGroupQuotaSubscriptionAllocationRequestClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armquota.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewGroupQuotaSubscriptionAllocationRequestClient().Get(ctx, "E7EC67B3-7657-4966-BFFC-41EFD36BAA09", "groupquota1", "Microsoft.Compute", "AE000000-0000-0000-0000-00000000000A", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.AllocationRequestStatus = armquota.AllocationRequestStatus{
	// 	Name: to.Ptr("AE000000-0000-0000-0000-00000000000A"),
	// 	Type: to.Ptr("Microsoft.Quota/groupQuotas/quotaAllocationRequests"),
	// 	ID: to.Ptr("/providers/Microsoft.Management/managementGroups/E7EC67B3-7657-4966-BFFC-41EFD36BAA09/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Quota/groupQuotas/groupquota1/quotaAllocationRequests/AE000000-0000-0000-0000-00000000000A"),
	// 	Properties: &armquota.AllocationRequestStatusProperties{
	// 		FaultCode: to.Ptr("ContactSupport"),
	// 		ProvisioningState: to.Ptr(armquota.RequestStateSucceeded),
	// 		RequestSubmitTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-11-17T01:06:02.191Z"); return t}()),
	// 		RequestedResource: &armquota.AllocationRequestBase{
	// 			Properties: &armquota.AllocationRequestBaseProperties{
	// 				Name: &armquota.AllocationRequestBasePropertiesName{
	// 					LocalizedValue: to.Ptr("standard Av2 Family vCPUs"),
	// 					Value: to.Ptr("standardav2family"),
	// 				},
	// 				Limit: to.Ptr[int64](75),
	// 				Region: to.Ptr("westus"),
	// 			},
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/8691e5081766c7ad602a9e55de841d07bed5196a/specification/quota/resource-manager/Microsoft.Quota/stable/2025-03-01/examples/SubscriptionQuotaAllocationRequests/SubscriptionQuotaAllocationRequests_List-Compute.json
func ExampleGroupQuotaSubscriptionAllocationRequestClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armquota.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewGroupQuotaSubscriptionAllocationRequestClient().NewListPager("E7EC67B3-7657-4966-BFFC-41EFD36BAA09", "groupquota1", "Microsoft.Compute", "location eq westus", nil)
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
		// page.AllocationRequestStatusList = armquota.AllocationRequestStatusList{
		// 	Value: []*armquota.AllocationRequestStatus{
		// 		{
		// 			Name: to.Ptr("AE000000-0000-0000-0000-00000000000A"),
		// 			Type: to.Ptr("Microsoft.Quota/groupQuotas/quotaAllocationRequests"),
		// 			ID: to.Ptr("/providers/Microsoft.Management/managementGroups/E7EC67B3-7657-4966-BFFC-41EFD36BAA09/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Quota/groupQuotas/groupquota1/resourceProviders/Microsoft.Compute/quotaAllocationRequests/AE000000-0000-0000-0000-00000000000A"),
		// 			Properties: &armquota.AllocationRequestStatusProperties{
		// 				FaultCode: to.Ptr("ContactSupport"),
		// 				ProvisioningState: to.Ptr(armquota.RequestStateSucceeded),
		// 				RequestSubmitTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-03-20T06:18:59.913Z"); return t}()),
		// 				RequestedResource: &armquota.AllocationRequestBase{
		// 					Properties: &armquota.AllocationRequestBaseProperties{
		// 						Name: &armquota.AllocationRequestBasePropertiesName{
		// 							LocalizedValue: to.Ptr("standard Av2 Family vCPUs"),
		// 							Value: to.Ptr("standardav2family"),
		// 						},
		// 						Limit: to.Ptr[int64](75),
		// 						Region: to.Ptr("westus"),
		// 					},
		// 				},
		// 			},
		// 	}},
		// }
	}
}
