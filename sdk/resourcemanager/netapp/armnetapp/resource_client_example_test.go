//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armnetapp_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp/v7"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/44319b51c6f952fdc9543d3dc4fdd9959350d102/specification/netapp/resource-manager/Microsoft.NetApp/stable/2025-03-01/examples/CheckNameAvailability.json
func ExampleResourceClient_CheckNameAvailability() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewResourceClient().CheckNameAvailability(ctx, "eastus", armnetapp.ResourceNameAvailabilityRequest{
		Name:          to.Ptr("accName"),
		Type:          to.Ptr(armnetapp.CheckNameResourceTypesMicrosoftNetAppNetAppAccounts),
		ResourceGroup: to.Ptr("myRG"),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.CheckAvailabilityResponse = armnetapp.CheckAvailabilityResponse{
	// 	IsAvailable: to.Ptr(true),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/44319b51c6f952fdc9543d3dc4fdd9959350d102/specification/netapp/resource-manager/Microsoft.NetApp/stable/2025-03-01/examples/CheckFilePathAvailability.json
func ExampleResourceClient_CheckFilePathAvailability() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewResourceClient().CheckFilePathAvailability(ctx, "eastus", armnetapp.FilePathAvailabilityRequest{
		Name:     to.Ptr("my-exact-filepth"),
		SubnetID: to.Ptr("/subscriptions/9760acf5-4638-11e7-9bdb-020073ca7778/resourceGroups/myRP/providers/Microsoft.Network/virtualNetworks/testvnet3/subnets/testsubnet3"),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.CheckAvailabilityResponse = armnetapp.CheckAvailabilityResponse{
	// 	IsAvailable: to.Ptr(true),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/44319b51c6f952fdc9543d3dc4fdd9959350d102/specification/netapp/resource-manager/Microsoft.NetApp/stable/2025-03-01/examples/CheckQuotaAvailability.json
func ExampleResourceClient_CheckQuotaAvailability() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewResourceClient().CheckQuotaAvailability(ctx, "eastus", armnetapp.QuotaAvailabilityRequest{
		Name:          to.Ptr("resource1"),
		Type:          to.Ptr(armnetapp.CheckQuotaNameResourceTypesMicrosoftNetAppNetAppAccounts),
		ResourceGroup: to.Ptr("myRG"),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.CheckAvailabilityResponse = armnetapp.CheckAvailabilityResponse{
	// 	IsAvailable: to.Ptr(true),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/44319b51c6f952fdc9543d3dc4fdd9959350d102/specification/netapp/resource-manager/Microsoft.NetApp/stable/2025-03-01/examples/RegionInfo.json
func ExampleResourceClient_QueryRegionInfo() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewResourceClient().QueryRegionInfo(ctx, "eastus", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.RegionInfo = armnetapp.RegionInfo{
	// 	AvailabilityZoneMappings: []*armnetapp.RegionInfoAvailabilityZoneMappingsItem{
	// 		{
	// 			AvailabilityZone: to.Ptr("1"),
	// 			IsAvailable: to.Ptr(true),
	// 	}},
	// 	StorageToNetworkProximity: to.Ptr(armnetapp.RegionStorageToNetworkProximityT2),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/44319b51c6f952fdc9543d3dc4fdd9959350d102/specification/netapp/resource-manager/Microsoft.NetApp/stable/2025-03-01/examples/NetworkSiblingSet_Query.json
func ExampleResourceClient_QueryNetworkSiblingSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewResourceClient().QueryNetworkSiblingSet(ctx, "eastus", armnetapp.QueryNetworkSiblingSetRequest{
		NetworkSiblingSetID: to.Ptr("9760acf5-4638-11e7-9bdb-020073ca3333"),
		SubnetID:            to.Ptr("/subscriptions/9760acf5-4638-11e7-9bdb-020073ca7778/resourceGroups/myRP/providers/Microsoft.Network/virtualNetworks/testVnet/subnets/testSubnet"),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.NetworkSiblingSet = armnetapp.NetworkSiblingSet{
	// 	NetworkFeatures: to.Ptr(armnetapp.NetworkFeaturesStandard),
	// 	NetworkSiblingSetID: to.Ptr("9760acf5-4638-11e7-9bdb-020073ca3333"),
	// 	NetworkSiblingSetStateID: to.Ptr("12345_44420.8001578125"),
	// 	NicInfoList: []*armnetapp.NicInfo{
	// 		{
	// 			IPAddress: to.Ptr("1.2.3.4"),
	// 			VolumeResourceIDs: []*string{
	// 				to.Ptr("/subscriptions/9760acf5-4638-11e7-9bdb-020073ca7778/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume10"),
	// 				to.Ptr("/subscriptions/9760acf5-4638-11e7-9bdb-020073ca7778/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume11")},
	// 			},
	// 			{
	// 				IPAddress: to.Ptr("1.2.3.5"),
	// 				VolumeResourceIDs: []*string{
	// 					to.Ptr("/subscriptions/9760acf5-4638-11e7-9bdb-020073ca7778/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account2/capacityPools/pool2/volumes/volume20"),
	// 					to.Ptr("/subscriptions/9760acf5-4638-11e7-9bdb-020073ca7778/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account2/capacityPools/pool2/volumes/volume21")},
	// 				},
	// 				{
	// 					IPAddress: to.Ptr("1.2.3.9"),
	// 					VolumeResourceIDs: []*string{
	// 					},
	// 			}},
	// 			ProvisioningState: to.Ptr(armnetapp.NetworkSiblingSetProvisioningStateSucceeded),
	// 			SubnetID: to.Ptr("/subscriptions/9760acf5-4638-11e7-9bdb-020073ca7778/resourceGroups/myRP/providers/Microsoft.Network/virtualNetworks/testVnet/subnets/testSubnet"),
	// 		}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/44319b51c6f952fdc9543d3dc4fdd9959350d102/specification/netapp/resource-manager/Microsoft.NetApp/stable/2025-03-01/examples/NetworkSiblingSet_Update.json
func ExampleResourceClient_BeginUpdateNetworkSiblingSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewResourceClient().BeginUpdateNetworkSiblingSet(ctx, "eastus", armnetapp.UpdateNetworkSiblingSetRequest{
		NetworkFeatures:          to.Ptr(armnetapp.NetworkFeaturesStandard),
		NetworkSiblingSetID:      to.Ptr("9760acf5-4638-11e7-9bdb-020073ca3333"),
		NetworkSiblingSetStateID: to.Ptr("12345_44420.8001578125"),
		SubnetID:                 to.Ptr("/subscriptions/9760acf5-4638-11e7-9bdb-020073ca7778/resourceGroups/myRP/providers/Microsoft.Network/virtualNetworks/testVnet/subnets/testSubnet"),
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
	// res.NetworkSiblingSet = armnetapp.NetworkSiblingSet{
	// 	NetworkFeatures: to.Ptr(armnetapp.NetworkFeaturesStandard),
	// 	NetworkSiblingSetID: to.Ptr("9760acf5-4638-11e7-9bdb-020073ca3333"),
	// 	NetworkSiblingSetStateID: to.Ptr("12345_44420.8001578125"),
	// 	NicInfoList: []*armnetapp.NicInfo{
	// 		{
	// 			IPAddress: to.Ptr("1.2.3.4"),
	// 			VolumeResourceIDs: []*string{
	// 				to.Ptr("/subscriptions/9760acf5-4638-11e7-9bdb-020073ca7778/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume10"),
	// 				to.Ptr("/subscriptions/9760acf5-4638-11e7-9bdb-020073ca7778/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume11")},
	// 			},
	// 			{
	// 				IPAddress: to.Ptr("1.2.3.5"),
	// 				VolumeResourceIDs: []*string{
	// 					to.Ptr("/subscriptions/9760acf5-4638-11e7-9bdb-020073ca7778/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account2/capacityPools/pool2/volumes/volume20"),
	// 					to.Ptr("/subscriptions/9760acf5-4638-11e7-9bdb-020073ca7778/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account2/capacityPools/pool2/volumes/volume21")},
	// 				},
	// 				{
	// 					IPAddress: to.Ptr("1.2.3.9"),
	// 					VolumeResourceIDs: []*string{
	// 					},
	// 			}},
	// 			ProvisioningState: to.Ptr(armnetapp.NetworkSiblingSetProvisioningStateSucceeded),
	// 			SubnetID: to.Ptr("/subscriptions/9760acf5-4638-11e7-9bdb-020073ca7778/resourceGroups/myRP/providers/Microsoft.Network/virtualNetworks/testVnet/subnets/testSubnet"),
	// 		}
}
