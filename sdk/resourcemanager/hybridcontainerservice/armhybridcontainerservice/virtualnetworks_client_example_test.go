//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armhybridcontainerservice_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybridcontainerservice/armhybridcontainerservice"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/41e4538ed7bb3ceac3c1322c9455a0812ed110ac/specification/hybridaks/resource-manager/Microsoft.HybridContainerService/stable/2024-01-01/examples/GetVirtualNetwork.json
func ExampleVirtualNetworksClient_Retrieve() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVirtualNetworksClient().Retrieve(ctx, "test-arcappliance-resgrp", "test-vnet-static", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VirtualNetwork = armhybridcontainerservice.VirtualNetwork{
	// 	Name: to.Ptr("test-vnet-static"),
	// 	Type: to.Ptr("microsoft.hybridcontainerservice/virtualnetworks"),
	// 	ID: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourceGroups/test-arcappliance-resgrp/providers/Microsoft.HybridContainerService/virtualNetworks/test-vnet-static"),
	// 	Location: to.Ptr("westus"),
	// 	ExtendedLocation: &armhybridcontainerservice.VirtualNetworkExtendedLocation{
	// 		Name: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourcegroups/test-arcappliance-resgrp/providers/microsoft.extendedlocation/customlocations/testcustomlocation"),
	// 		Type: to.Ptr(armhybridcontainerservice.ExtendedLocationTypesCustomLocation),
	// 	},
	// 	Properties: &armhybridcontainerservice.VirtualNetworkProperties{
	// 		DNSServers: []*string{
	// 			to.Ptr("192.168.0.1")},
	// 			Gateway: to.Ptr("192.168.0.1"),
	// 			InfraVnetProfile: &armhybridcontainerservice.VirtualNetworkPropertiesInfraVnetProfile{
	// 				Hci: &armhybridcontainerservice.VirtualNetworkPropertiesInfraVnetProfileHci{
	// 					MocGroup: to.Ptr("target-group"),
	// 					MocLocation: to.Ptr("MocLocation"),
	// 					MocVnetName: to.Ptr("vnet1"),
	// 				},
	// 			},
	// 			IPAddressPrefix: to.Ptr("192.168.0.0/16"),
	// 			ProvisioningState: to.Ptr(armhybridcontainerservice.ProvisioningStateSucceeded),
	// 			VipPool: []*armhybridcontainerservice.VirtualNetworkPropertiesVipPoolItem{
	// 				{
	// 					EndIP: to.Ptr("192.168.0.50"),
	// 					StartIP: to.Ptr("192.168.0.10"),
	// 			}},
	// 			VlanID: to.Ptr[int32](10),
	// 			VmipPool: []*armhybridcontainerservice.VirtualNetworkPropertiesVmipPoolItem{
	// 				{
	// 					EndIP: to.Ptr("192.168.0.130"),
	// 					StartIP: to.Ptr("192.168.0.110"),
	// 			}},
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/41e4538ed7bb3ceac3c1322c9455a0812ed110ac/specification/hybridaks/resource-manager/Microsoft.HybridContainerService/stable/2024-01-01/examples/PutVirtualNetwork.json
func ExampleVirtualNetworksClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualNetworksClient().BeginCreateOrUpdate(ctx, "test-arcappliance-resgrp", "test-vnet-static", armhybridcontainerservice.VirtualNetwork{
		Location: to.Ptr("westus"),
		ExtendedLocation: &armhybridcontainerservice.VirtualNetworkExtendedLocation{
			Name: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourcegroups/test-arcappliance-resgrp/providers/microsoft.extendedlocation/customlocations/testcustomlocation"),
			Type: to.Ptr(armhybridcontainerservice.ExtendedLocationTypesCustomLocation),
		},
		Properties: &armhybridcontainerservice.VirtualNetworkProperties{
			DNSServers: []*string{
				to.Ptr("192.168.0.1")},
			Gateway: to.Ptr("192.168.0.1"),
			InfraVnetProfile: &armhybridcontainerservice.VirtualNetworkPropertiesInfraVnetProfile{
				Hci: &armhybridcontainerservice.VirtualNetworkPropertiesInfraVnetProfileHci{
					MocGroup:    to.Ptr("target-group"),
					MocLocation: to.Ptr("MocLocation"),
					MocVnetName: to.Ptr("vnet1"),
				},
			},
			IPAddressPrefix: to.Ptr("192.168.0.0/16"),
			VipPool: []*armhybridcontainerservice.VirtualNetworkPropertiesVipPoolItem{
				{
					EndIP:   to.Ptr("192.168.0.50"),
					StartIP: to.Ptr("192.168.0.10"),
				}},
			VlanID: to.Ptr[int32](10),
			VmipPool: []*armhybridcontainerservice.VirtualNetworkPropertiesVmipPoolItem{
				{
					EndIP:   to.Ptr("192.168.0.130"),
					StartIP: to.Ptr("192.168.0.110"),
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
	// res.VirtualNetwork = armhybridcontainerservice.VirtualNetwork{
	// 	Name: to.Ptr("test-vnet-static"),
	// 	Type: to.Ptr("microsoft.hybridcontainerservice/virtualnetworks"),
	// 	ID: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourceGroups/test-arcappliance-resgrp/providers/Microsoft.HybridContainerService/virtualNetworks/test-vnet-static"),
	// 	Location: to.Ptr("westus"),
	// 	ExtendedLocation: &armhybridcontainerservice.VirtualNetworkExtendedLocation{
	// 		Name: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourcegroups/test-arcappliance-resgrp/providers/microsoft.extendedlocation/customlocations/testcustomlocation"),
	// 		Type: to.Ptr(armhybridcontainerservice.ExtendedLocationTypesCustomLocation),
	// 	},
	// 	Properties: &armhybridcontainerservice.VirtualNetworkProperties{
	// 		DNSServers: []*string{
	// 			to.Ptr("192.168.0.1")},
	// 			Gateway: to.Ptr("192.168.0.1"),
	// 			InfraVnetProfile: &armhybridcontainerservice.VirtualNetworkPropertiesInfraVnetProfile{
	// 				Hci: &armhybridcontainerservice.VirtualNetworkPropertiesInfraVnetProfileHci{
	// 					MocGroup: to.Ptr("target-group"),
	// 					MocLocation: to.Ptr("MocLocation"),
	// 					MocVnetName: to.Ptr("vnet1"),
	// 				},
	// 			},
	// 			IPAddressPrefix: to.Ptr("192.168.0.0/16"),
	// 			ProvisioningState: to.Ptr(armhybridcontainerservice.ProvisioningStateSucceeded),
	// 			VipPool: []*armhybridcontainerservice.VirtualNetworkPropertiesVipPoolItem{
	// 				{
	// 					EndIP: to.Ptr("192.168.0.50"),
	// 					StartIP: to.Ptr("192.168.0.10"),
	// 			}},
	// 			VlanID: to.Ptr[int32](10),
	// 			VmipPool: []*armhybridcontainerservice.VirtualNetworkPropertiesVmipPoolItem{
	// 				{
	// 					EndIP: to.Ptr("192.168.0.130"),
	// 					StartIP: to.Ptr("192.168.0.110"),
	// 			}},
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/41e4538ed7bb3ceac3c1322c9455a0812ed110ac/specification/hybridaks/resource-manager/Microsoft.HybridContainerService/stable/2024-01-01/examples/DeleteVirtualNetwork.json
func ExampleVirtualNetworksClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualNetworksClient().BeginDelete(ctx, "test-arcappliance-resgrp", "test-vnet-static", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/41e4538ed7bb3ceac3c1322c9455a0812ed110ac/specification/hybridaks/resource-manager/Microsoft.HybridContainerService/stable/2024-01-01/examples/UpdateVirtualNetwork.json
func ExampleVirtualNetworksClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualNetworksClient().BeginUpdate(ctx, "test-arcappliance-resgrp", "test-vnet-static", armhybridcontainerservice.VirtualNetworksPatch{
		Tags: map[string]*string{
			"additionalProperties": to.Ptr("sample"),
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
	// res.VirtualNetwork = armhybridcontainerservice.VirtualNetwork{
	// 	Name: to.Ptr("test-vnet-static"),
	// 	Type: to.Ptr("microsoft.hybridcontainerservice/virtualnetworks"),
	// 	ID: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourceGroups/test-arcappliance-resgrp/providers/Microsoft.HybridContainerService/virtualNetworks/test-vnet-static"),
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 		"additionalProperties": to.Ptr("sample"),
	// 	},
	// 	ExtendedLocation: &armhybridcontainerservice.VirtualNetworkExtendedLocation{
	// 		Name: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourcegroups/test-arcappliance-resgrp/providers/microsoft.extendedlocation/customlocations/testcustomlocation"),
	// 		Type: to.Ptr(armhybridcontainerservice.ExtendedLocationTypesCustomLocation),
	// 	},
	// 	Properties: &armhybridcontainerservice.VirtualNetworkProperties{
	// 		DNSServers: []*string{
	// 			to.Ptr("192.168.0.1")},
	// 			Gateway: to.Ptr("192.168.0.1"),
	// 			InfraVnetProfile: &armhybridcontainerservice.VirtualNetworkPropertiesInfraVnetProfile{
	// 				Hci: &armhybridcontainerservice.VirtualNetworkPropertiesInfraVnetProfileHci{
	// 					MocGroup: to.Ptr("target-group"),
	// 					MocLocation: to.Ptr("MocLocation"),
	// 					MocVnetName: to.Ptr("vnet1"),
	// 				},
	// 			},
	// 			IPAddressPrefix: to.Ptr("192.168.0.0/16"),
	// 			ProvisioningState: to.Ptr(armhybridcontainerservice.ProvisioningStateSucceeded),
	// 			VipPool: []*armhybridcontainerservice.VirtualNetworkPropertiesVipPoolItem{
	// 				{
	// 					EndIP: to.Ptr("192.168.0.50"),
	// 					StartIP: to.Ptr("192.168.0.10"),
	// 			}},
	// 			VlanID: to.Ptr[int32](10),
	// 			VmipPool: []*armhybridcontainerservice.VirtualNetworkPropertiesVmipPoolItem{
	// 				{
	// 					EndIP: to.Ptr("192.168.0.130"),
	// 					StartIP: to.Ptr("192.168.0.110"),
	// 			}},
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/41e4538ed7bb3ceac3c1322c9455a0812ed110ac/specification/hybridaks/resource-manager/Microsoft.HybridContainerService/stable/2024-01-01/examples/ListVirtualNetworkByResourceGroup.json
func ExampleVirtualNetworksClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewVirtualNetworksClient().NewListByResourceGroupPager("test-arcappliance-resgrp", nil)
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
		// page.VirtualNetworksListResult = armhybridcontainerservice.VirtualNetworksListResult{
		// 	Value: []*armhybridcontainerservice.VirtualNetwork{
		// 		{
		// 			Name: to.Ptr("test-vnet-static"),
		// 			Type: to.Ptr("microsoft.hybridcontainerservice/virtualnetworks"),
		// 			ID: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourceGroups/test-arcappliance-resgrp/providers/Microsoft.HybridContainerService/virtualNetworks/test-vnet-static"),
		// 			Location: to.Ptr("westus"),
		// 			ExtendedLocation: &armhybridcontainerservice.VirtualNetworkExtendedLocation{
		// 				Name: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourcegroups/test-arcappliance-resgrp/providers/microsoft.extendedlocation/customlocations/testcustomlocation"),
		// 				Type: to.Ptr(armhybridcontainerservice.ExtendedLocationTypesCustomLocation),
		// 			},
		// 			Properties: &armhybridcontainerservice.VirtualNetworkProperties{
		// 				DNSServers: []*string{
		// 					to.Ptr("192.168.0.1")},
		// 					Gateway: to.Ptr("192.168.0.1"),
		// 					InfraVnetProfile: &armhybridcontainerservice.VirtualNetworkPropertiesInfraVnetProfile{
		// 						Hci: &armhybridcontainerservice.VirtualNetworkPropertiesInfraVnetProfileHci{
		// 							MocGroup: to.Ptr("target-group"),
		// 							MocLocation: to.Ptr("MocLocation"),
		// 							MocVnetName: to.Ptr("vnet1"),
		// 						},
		// 					},
		// 					IPAddressPrefix: to.Ptr("192.168.0.0/16"),
		// 					ProvisioningState: to.Ptr(armhybridcontainerservice.ProvisioningStateSucceeded),
		// 					VipPool: []*armhybridcontainerservice.VirtualNetworkPropertiesVipPoolItem{
		// 						{
		// 							EndIP: to.Ptr("192.168.0.50"),
		// 							StartIP: to.Ptr("192.168.0.10"),
		// 					}},
		// 					VlanID: to.Ptr[int32](10),
		// 					VmipPool: []*armhybridcontainerservice.VirtualNetworkPropertiesVmipPoolItem{
		// 						{
		// 							EndIP: to.Ptr("192.168.0.130"),
		// 							StartIP: to.Ptr("192.168.0.110"),
		// 					}},
		// 				},
		// 		}},
		// 	}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/41e4538ed7bb3ceac3c1322c9455a0812ed110ac/specification/hybridaks/resource-manager/Microsoft.HybridContainerService/stable/2024-01-01/examples/ListVirtualNetworkBySubscription.json
func ExampleVirtualNetworksClient_NewListBySubscriptionPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewVirtualNetworksClient().NewListBySubscriptionPager(nil)
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
		// page.VirtualNetworksListResult = armhybridcontainerservice.VirtualNetworksListResult{
		// 	Value: []*armhybridcontainerservice.VirtualNetwork{
		// 		{
		// 			Name: to.Ptr("test-vnet-static"),
		// 			Type: to.Ptr("microsoft.hybridcontainerservice/virtualnetworks"),
		// 			ID: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourceGroups/test-arcappliance-resgrp/providers/Microsoft.HybridContainerService/virtualNetworks/test-vnet-static"),
		// 			Location: to.Ptr("westus"),
		// 			ExtendedLocation: &armhybridcontainerservice.VirtualNetworkExtendedLocation{
		// 				Name: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourcegroups/test-arcappliance-resgrp/providers/microsoft.extendedlocation/customlocations/testcustomlocation"),
		// 				Type: to.Ptr(armhybridcontainerservice.ExtendedLocationTypesCustomLocation),
		// 			},
		// 			Properties: &armhybridcontainerservice.VirtualNetworkProperties{
		// 				DNSServers: []*string{
		// 					to.Ptr("192.168.0.1")},
		// 					Gateway: to.Ptr("192.168.0.1"),
		// 					InfraVnetProfile: &armhybridcontainerservice.VirtualNetworkPropertiesInfraVnetProfile{
		// 						Hci: &armhybridcontainerservice.VirtualNetworkPropertiesInfraVnetProfileHci{
		// 							MocGroup: to.Ptr("target-group"),
		// 							MocLocation: to.Ptr("MocLocation"),
		// 							MocVnetName: to.Ptr("vnet1"),
		// 						},
		// 					},
		// 					IPAddressPrefix: to.Ptr("192.168.0.0/16"),
		// 					ProvisioningState: to.Ptr(armhybridcontainerservice.ProvisioningStateSucceeded),
		// 					VipPool: []*armhybridcontainerservice.VirtualNetworkPropertiesVipPoolItem{
		// 						{
		// 							EndIP: to.Ptr("192.168.0.50"),
		// 							StartIP: to.Ptr("192.168.0.10"),
		// 					}},
		// 					VlanID: to.Ptr[int32](10),
		// 					VmipPool: []*armhybridcontainerservice.VirtualNetworkPropertiesVmipPoolItem{
		// 						{
		// 							EndIP: to.Ptr("192.168.0.130"),
		// 							StartIP: to.Ptr("192.168.0.110"),
		// 					}},
		// 				},
		// 		}},
		// 	}
	}
}
