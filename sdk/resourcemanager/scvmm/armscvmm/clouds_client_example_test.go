//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armscvmm_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/scvmm/armscvmm"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/scvmm/resource-manager/Microsoft.ScVmm/stable/2023-10-07/examples/Clouds_ListBySubscription_MaximumSet_Gen.json
func ExampleCloudsClient_NewListBySubscriptionPager_cloudsListBySubscriptionMaximumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armscvmm.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewCloudsClient().NewListBySubscriptionPager(nil)
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
		// page.CloudListResult = armscvmm.CloudListResult{
		// 	Value: []*armscvmm.Cloud{
		// 		{
		// 			Name: to.Ptr("wwcwalpiufsfbnydxpr"),
		// 			Type: to.Ptr("qnaaimszbuokldohwrdfuiitpy"),
		// 			ID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ScVmm/clouds/cloudResourceName"),
		// 			SystemData: &armscvmm.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-29T22:28:00.094Z"); return t}()),
		// 				CreatedBy: to.Ptr("p"),
		// 				CreatedByType: to.Ptr(armscvmm.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-29T22:28:00.095Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("goxcwpyyqlxndquly"),
		// 				LastModifiedByType: to.Ptr(armscvmm.CreatedByTypeUser),
		// 			},
		// 			Location: to.Ptr("khwsdmaxfhmbu"),
		// 			Tags: map[string]*string{
		// 				"key4295": to.Ptr("wngosgcbdifaxdobufuuqxtho"),
		// 			},
		// 			ExtendedLocation: &armscvmm.ExtendedLocation{
		// 				Name: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ExtendedLocation/customLocations/customLocationName"),
		// 				Type: to.Ptr("customLocation"),
		// 			},
		// 			Properties: &armscvmm.CloudProperties{
		// 				CloudCapacity: &armscvmm.CloudCapacity{
		// 					CPUCount: to.Ptr[int64](4),
		// 					MemoryMB: to.Ptr[int64](19),
		// 					VMCount: to.Ptr[int64](28),
		// 				},
		// 				CloudName: to.Ptr("menarjsplhcqvnkjdwieroir"),
		// 				InventoryItemID: to.Ptr("qjd"),
		// 				ProvisioningState: to.Ptr(armscvmm.ResourceProvisioningStateSucceeded),
		// 				StorageQosPolicies: []*armscvmm.StorageQosPolicy{
		// 					{
		// 						Name: to.Ptr("hvqcentnbwcunxhzfavyewhwlo"),
		// 						BandwidthLimit: to.Ptr[int64](26),
		// 						ID: to.Ptr("oclhgkydaw"),
		// 						IopsMaximum: to.Ptr[int64](6),
		// 						IopsMinimum: to.Ptr[int64](25),
		// 						PolicyID: to.Ptr("lvcylbmxrqjgarvhfny"),
		// 				}},
		// 				UUID: to.Ptr("12345678-1234-1234-1234-12345678abcd"),
		// 				VmmServerID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ScVmm/vmmServers/vmmServerName"),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/scvmm/resource-manager/Microsoft.ScVmm/stable/2023-10-07/examples/Clouds_ListBySubscription_MinimumSet_Gen.json
func ExampleCloudsClient_NewListBySubscriptionPager_cloudsListBySubscriptionMinimumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armscvmm.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewCloudsClient().NewListBySubscriptionPager(nil)
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
		// page.CloudListResult = armscvmm.CloudListResult{
		// 	Value: []*armscvmm.Cloud{
		// 		{
		// 			ID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ScVmm/clouds/cloudResourceName"),
		// 			Location: to.Ptr("khwsdmaxfhmbu"),
		// 			ExtendedLocation: &armscvmm.ExtendedLocation{
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/scvmm/resource-manager/Microsoft.ScVmm/stable/2023-10-07/examples/Clouds_ListByResourceGroup_MaximumSet_Gen.json
func ExampleCloudsClient_NewListByResourceGroupPager_cloudsListByResourceGroupMaximumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armscvmm.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewCloudsClient().NewListByResourceGroupPager("rgscvmm", nil)
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
		// page.CloudListResult = armscvmm.CloudListResult{
		// 	Value: []*armscvmm.Cloud{
		// 		{
		// 			Name: to.Ptr("wwcwalpiufsfbnydxpr"),
		// 			Type: to.Ptr("qnaaimszbuokldohwrdfuiitpy"),
		// 			ID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ScVmm/clouds/cloudResourceName"),
		// 			SystemData: &armscvmm.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-29T22:28:00.094Z"); return t}()),
		// 				CreatedBy: to.Ptr("p"),
		// 				CreatedByType: to.Ptr(armscvmm.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-29T22:28:00.095Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("goxcwpyyqlxndquly"),
		// 				LastModifiedByType: to.Ptr(armscvmm.CreatedByTypeUser),
		// 			},
		// 			Location: to.Ptr("khwsdmaxfhmbu"),
		// 			Tags: map[string]*string{
		// 				"key4295": to.Ptr("wngosgcbdifaxdobufuuqxtho"),
		// 			},
		// 			ExtendedLocation: &armscvmm.ExtendedLocation{
		// 				Name: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ExtendedLocation/customLocations/customLocationName"),
		// 				Type: to.Ptr("customLocation"),
		// 			},
		// 			Properties: &armscvmm.CloudProperties{
		// 				CloudCapacity: &armscvmm.CloudCapacity{
		// 					CPUCount: to.Ptr[int64](4),
		// 					MemoryMB: to.Ptr[int64](19),
		// 					VMCount: to.Ptr[int64](28),
		// 				},
		// 				CloudName: to.Ptr("menarjsplhcqvnkjdwieroir"),
		// 				InventoryItemID: to.Ptr("qjd"),
		// 				ProvisioningState: to.Ptr(armscvmm.ResourceProvisioningStateSucceeded),
		// 				StorageQosPolicies: []*armscvmm.StorageQosPolicy{
		// 					{
		// 						Name: to.Ptr("hvqcentnbwcunxhzfavyewhwlo"),
		// 						BandwidthLimit: to.Ptr[int64](26),
		// 						ID: to.Ptr("oclhgkydaw"),
		// 						IopsMaximum: to.Ptr[int64](6),
		// 						IopsMinimum: to.Ptr[int64](25),
		// 						PolicyID: to.Ptr("lvcylbmxrqjgarvhfny"),
		// 				}},
		// 				UUID: to.Ptr("12345678-1234-1234-1234-12345678abcd"),
		// 				VmmServerID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ScVmm/vmmServers/vmmServerName"),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/scvmm/resource-manager/Microsoft.ScVmm/stable/2023-10-07/examples/Clouds_ListByResourceGroup_MinimumSet_Gen.json
func ExampleCloudsClient_NewListByResourceGroupPager_cloudsListByResourceGroupMinimumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armscvmm.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewCloudsClient().NewListByResourceGroupPager("rgscvmm", nil)
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
		// page.CloudListResult = armscvmm.CloudListResult{
		// 	Value: []*armscvmm.Cloud{
		// 		{
		// 			ID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ScVmm/clouds/cloudResourceName"),
		// 			Location: to.Ptr("khwsdmaxfhmbu"),
		// 			ExtendedLocation: &armscvmm.ExtendedLocation{
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/scvmm/resource-manager/Microsoft.ScVmm/stable/2023-10-07/examples/Clouds_Get_MaximumSet_Gen.json
func ExampleCloudsClient_Get_cloudsGetMaximumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armscvmm.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewCloudsClient().Get(ctx, "rgscvmm", "_", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Cloud = armscvmm.Cloud{
	// 	Name: to.Ptr("wwcwalpiufsfbnydxpr"),
	// 	Type: to.Ptr("qnaaimszbuokldohwrdfuiitpy"),
	// 	ID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ScVmm/clouds/cloudResourceName"),
	// 	SystemData: &armscvmm.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-29T22:28:00.094Z"); return t}()),
	// 		CreatedBy: to.Ptr("p"),
	// 		CreatedByType: to.Ptr(armscvmm.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-29T22:28:00.095Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("goxcwpyyqlxndquly"),
	// 		LastModifiedByType: to.Ptr(armscvmm.CreatedByTypeUser),
	// 	},
	// 	Location: to.Ptr("khwsdmaxfhmbu"),
	// 	Tags: map[string]*string{
	// 		"key4295": to.Ptr("wngosgcbdifaxdobufuuqxtho"),
	// 	},
	// 	ExtendedLocation: &armscvmm.ExtendedLocation{
	// 		Name: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ExtendedLocation/customLocations/customLocationName"),
	// 		Type: to.Ptr("customLocation"),
	// 	},
	// 	Properties: &armscvmm.CloudProperties{
	// 		CloudCapacity: &armscvmm.CloudCapacity{
	// 			CPUCount: to.Ptr[int64](4),
	// 			MemoryMB: to.Ptr[int64](19),
	// 			VMCount: to.Ptr[int64](28),
	// 		},
	// 		CloudName: to.Ptr("menarjsplhcqvnkjdwieroir"),
	// 		InventoryItemID: to.Ptr("qjd"),
	// 		ProvisioningState: to.Ptr(armscvmm.ResourceProvisioningStateSucceeded),
	// 		StorageQosPolicies: []*armscvmm.StorageQosPolicy{
	// 			{
	// 				Name: to.Ptr("hvqcentnbwcunxhzfavyewhwlo"),
	// 				BandwidthLimit: to.Ptr[int64](26),
	// 				ID: to.Ptr("oclhgkydaw"),
	// 				IopsMaximum: to.Ptr[int64](6),
	// 				IopsMinimum: to.Ptr[int64](25),
	// 				PolicyID: to.Ptr("lvcylbmxrqjgarvhfny"),
	// 		}},
	// 		UUID: to.Ptr("12345678-1234-1234-1234-12345678abcd"),
	// 		VmmServerID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ScVmm/vmmServers/vmmServerName"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/scvmm/resource-manager/Microsoft.ScVmm/stable/2023-10-07/examples/Clouds_Get_MinimumSet_Gen.json
func ExampleCloudsClient_Get_cloudsGetMinimumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armscvmm.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewCloudsClient().Get(ctx, "rgscvmm", "i", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Cloud = armscvmm.Cloud{
	// 	ID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ScVmm/clouds/cloudResourceName"),
	// 	Location: to.Ptr("khwsdmaxfhmbu"),
	// 	ExtendedLocation: &armscvmm.ExtendedLocation{
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/scvmm/resource-manager/Microsoft.ScVmm/stable/2023-10-07/examples/Clouds_CreateOrUpdate_MaximumSet_Gen.json
func ExampleCloudsClient_BeginCreateOrUpdate_cloudsCreateOrUpdateMaximumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armscvmm.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewCloudsClient().BeginCreateOrUpdate(ctx, "rgscvmm", "2", armscvmm.Cloud{
		Location: to.Ptr("khwsdmaxfhmbu"),
		Tags: map[string]*string{
			"key4295": to.Ptr("wngosgcbdifaxdobufuuqxtho"),
		},
		ExtendedLocation: &armscvmm.ExtendedLocation{
			Name: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ExtendedLocation/customLocations/customLocationName"),
			Type: to.Ptr("customLocation"),
		},
		Properties: &armscvmm.CloudProperties{
			CloudCapacity:   &armscvmm.CloudCapacity{},
			InventoryItemID: to.Ptr("qjd"),
			UUID:            to.Ptr("12345678-1234-1234-1234-12345678abcd"),
			VmmServerID:     to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ScVmm/vmmServers/vmmServerName"),
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
	// res.Cloud = armscvmm.Cloud{
	// 	Name: to.Ptr("wwcwalpiufsfbnydxpr"),
	// 	Type: to.Ptr("qnaaimszbuokldohwrdfuiitpy"),
	// 	ID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ScVmm/clouds/cloudResourceName"),
	// 	SystemData: &armscvmm.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-29T22:28:00.094Z"); return t}()),
	// 		CreatedBy: to.Ptr("p"),
	// 		CreatedByType: to.Ptr(armscvmm.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-29T22:28:00.095Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("goxcwpyyqlxndquly"),
	// 		LastModifiedByType: to.Ptr(armscvmm.CreatedByTypeUser),
	// 	},
	// 	Location: to.Ptr("khwsdmaxfhmbu"),
	// 	Tags: map[string]*string{
	// 		"key4295": to.Ptr("wngosgcbdifaxdobufuuqxtho"),
	// 	},
	// 	ExtendedLocation: &armscvmm.ExtendedLocation{
	// 		Name: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ExtendedLocation/customLocations/customLocationName"),
	// 		Type: to.Ptr("customLocation"),
	// 	},
	// 	Properties: &armscvmm.CloudProperties{
	// 		CloudCapacity: &armscvmm.CloudCapacity{
	// 			CPUCount: to.Ptr[int64](4),
	// 			MemoryMB: to.Ptr[int64](19),
	// 			VMCount: to.Ptr[int64](28),
	// 		},
	// 		CloudName: to.Ptr("menarjsplhcqvnkjdwieroir"),
	// 		InventoryItemID: to.Ptr("qjd"),
	// 		ProvisioningState: to.Ptr(armscvmm.ResourceProvisioningStateSucceeded),
	// 		StorageQosPolicies: []*armscvmm.StorageQosPolicy{
	// 			{
	// 				Name: to.Ptr("hvqcentnbwcunxhzfavyewhwlo"),
	// 				BandwidthLimit: to.Ptr[int64](26),
	// 				ID: to.Ptr("oclhgkydaw"),
	// 				IopsMaximum: to.Ptr[int64](6),
	// 				IopsMinimum: to.Ptr[int64](25),
	// 				PolicyID: to.Ptr("lvcylbmxrqjgarvhfny"),
	// 		}},
	// 		UUID: to.Ptr("12345678-1234-1234-1234-12345678abcd"),
	// 		VmmServerID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ScVmm/vmmServers/vmmServerName"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/scvmm/resource-manager/Microsoft.ScVmm/stable/2023-10-07/examples/Clouds_CreateOrUpdate_MinimumSet_Gen.json
func ExampleCloudsClient_BeginCreateOrUpdate_cloudsCreateOrUpdateMinimumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armscvmm.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewCloudsClient().BeginCreateOrUpdate(ctx, "rgscvmm", "-", armscvmm.Cloud{
		Location:         to.Ptr("khwsdmaxfhmbu"),
		ExtendedLocation: &armscvmm.ExtendedLocation{},
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
	// res.Cloud = armscvmm.Cloud{
	// 	ID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ScVmm/clouds/cloudResourceName"),
	// 	Location: to.Ptr("khwsdmaxfhmbu"),
	// 	ExtendedLocation: &armscvmm.ExtendedLocation{
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/scvmm/resource-manager/Microsoft.ScVmm/stable/2023-10-07/examples/Clouds_Update_MaximumSet_Gen.json
func ExampleCloudsClient_BeginUpdate_cloudsUpdateMaximumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armscvmm.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewCloudsClient().BeginUpdate(ctx, "rgscvmm", "P", armscvmm.CloudTagsUpdate{
		Tags: map[string]*string{
			"key5266": to.Ptr("hjpcnwmpnixsolrxnbl"),
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
	// res.Cloud = armscvmm.Cloud{
	// 	Name: to.Ptr("wwcwalpiufsfbnydxpr"),
	// 	Type: to.Ptr("qnaaimszbuokldohwrdfuiitpy"),
	// 	ID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ScVmm/clouds/cloudResourceName"),
	// 	SystemData: &armscvmm.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-29T22:28:00.094Z"); return t}()),
	// 		CreatedBy: to.Ptr("p"),
	// 		CreatedByType: to.Ptr(armscvmm.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-29T22:28:00.095Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("goxcwpyyqlxndquly"),
	// 		LastModifiedByType: to.Ptr(armscvmm.CreatedByTypeUser),
	// 	},
	// 	Location: to.Ptr("khwsdmaxfhmbu"),
	// 	Tags: map[string]*string{
	// 		"key4295": to.Ptr("wngosgcbdifaxdobufuuqxtho"),
	// 	},
	// 	ExtendedLocation: &armscvmm.ExtendedLocation{
	// 		Name: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ExtendedLocation/customLocations/customLocationName"),
	// 		Type: to.Ptr("customLocation"),
	// 	},
	// 	Properties: &armscvmm.CloudProperties{
	// 		CloudCapacity: &armscvmm.CloudCapacity{
	// 			CPUCount: to.Ptr[int64](4),
	// 			MemoryMB: to.Ptr[int64](19),
	// 			VMCount: to.Ptr[int64](28),
	// 		},
	// 		CloudName: to.Ptr("menarjsplhcqvnkjdwieroir"),
	// 		InventoryItemID: to.Ptr("qjd"),
	// 		ProvisioningState: to.Ptr(armscvmm.ResourceProvisioningStateSucceeded),
	// 		StorageQosPolicies: []*armscvmm.StorageQosPolicy{
	// 			{
	// 				Name: to.Ptr("hvqcentnbwcunxhzfavyewhwlo"),
	// 				BandwidthLimit: to.Ptr[int64](26),
	// 				ID: to.Ptr("oclhgkydaw"),
	// 				IopsMaximum: to.Ptr[int64](6),
	// 				IopsMinimum: to.Ptr[int64](25),
	// 				PolicyID: to.Ptr("lvcylbmxrqjgarvhfny"),
	// 		}},
	// 		UUID: to.Ptr("12345678-1234-1234-1234-12345678abcd"),
	// 		VmmServerID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.ScVmm/vmmServers/vmmServerName"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/scvmm/resource-manager/Microsoft.ScVmm/stable/2023-10-07/examples/Clouds_Update_MinimumSet_Gen.json
func ExampleCloudsClient_BeginUpdate_cloudsUpdateMinimumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armscvmm.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewCloudsClient().BeginUpdate(ctx, "rgscvmm", "_", armscvmm.CloudTagsUpdate{}, nil)
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
	// res.Cloud = armscvmm.Cloud{
	// 	Location: to.Ptr("khwsdmaxfhmbu"),
	// 	ExtendedLocation: &armscvmm.ExtendedLocation{
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/scvmm/resource-manager/Microsoft.ScVmm/stable/2023-10-07/examples/Clouds_Delete_MaximumSet_Gen.json
func ExampleCloudsClient_BeginDelete_cloudsDeleteMaximumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armscvmm.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewCloudsClient().BeginDelete(ctx, "rgscvmm", "-", &armscvmm.CloudsClientBeginDeleteOptions{Force: to.Ptr(armscvmm.ForceDeleteTrue)})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/scvmm/resource-manager/Microsoft.ScVmm/stable/2023-10-07/examples/Clouds_Delete_MinimumSet_Gen.json
func ExampleCloudsClient_BeginDelete_cloudsDeleteMinimumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armscvmm.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewCloudsClient().BeginDelete(ctx, "rgscvmm", "1", &armscvmm.CloudsClientBeginDeleteOptions{Force: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
