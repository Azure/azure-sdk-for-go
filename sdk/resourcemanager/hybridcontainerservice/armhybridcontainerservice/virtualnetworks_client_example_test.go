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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/a60468a0c5e2beb054680ae488fb9f92699f0a0d/specification/hybridaks/resource-manager/Microsoft.HybridContainerService/preview/2022-09-01-preview/examples/GetVirtualNetwork.json
func ExampleVirtualNetworksClient_Retrieve() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhybridcontainerservice.NewVirtualNetworksClient("a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Retrieve(ctx, "test-arcappliance-resgrp", "test-vnet-static", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VirtualNetworks = armhybridcontainerservice.VirtualNetworks{
	// 	Name: to.Ptr("test-vnet-static"),
	// 	Type: to.Ptr("microsoft.hybridcontainerservice/virtualnetworks"),
	// 	ID: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourceGroups/test-arcappliance-resgrp/providers/Microsoft.HybridContainerService/virtualNetworks/test-vnet-static"),
	// 	Location: to.Ptr("westus"),
	// 	ExtendedLocation: &armhybridcontainerservice.VirtualNetworksExtendedLocation{
	// 		Name: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourcegroups/test-arcappliance-resgrp/providers/microsoft.extendedlocation/customlocations/testcustomlocation"),
	// 		Type: to.Ptr("CustomLocation"),
	// 	},
	// 	Properties: &armhybridcontainerservice.VirtualNetworksProperties{
	// 		DNSServers: []*string{
	// 			to.Ptr("192.168.0.1")},
	// 			Gateway: to.Ptr("192.168.0.1"),
	// 			InfraVnetProfile: &armhybridcontainerservice.VirtualNetworksPropertiesInfraVnetProfile{
	// 				Hci: &armhybridcontainerservice.VirtualNetworksPropertiesInfraVnetProfileHci{
	// 					MocGroup: to.Ptr("target-group"),
	// 					MocLocation: to.Ptr("MocLocation"),
	// 					MocVnetName: to.Ptr("test-vnet"),
	// 				},
	// 			},
	// 			IPAddressPrefix: to.Ptr("192.168.0.0/16"),
	// 			ProvisioningState: to.Ptr(armhybridcontainerservice.ProvisioningStateSucceeded),
	// 			VipPool: []*armhybridcontainerservice.VirtualNetworksPropertiesVipPoolItem{
	// 				{
	// 					EndIP: to.Ptr("192.168.0.50"),
	// 					StartIP: to.Ptr("192.168.0.10"),
	// 			}},
	// 			VmipPool: []*armhybridcontainerservice.VirtualNetworksPropertiesVmipPoolItem{
	// 				{
	// 					EndIP: to.Ptr("192.168.0.130"),
	// 					StartIP: to.Ptr("192.168.0.110"),
	// 			}},
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/a60468a0c5e2beb054680ae488fb9f92699f0a0d/specification/hybridaks/resource-manager/Microsoft.HybridContainerService/preview/2022-09-01-preview/examples/PutVirtualNetwork.json
func ExampleVirtualNetworksClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhybridcontainerservice.NewVirtualNetworksClient("a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginCreateOrUpdate(ctx, "test-arcappliance-resgrp", "test-vnet-static", armhybridcontainerservice.VirtualNetworks{
		Location: to.Ptr("westus"),
		ExtendedLocation: &armhybridcontainerservice.VirtualNetworksExtendedLocation{
			Name: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourcegroups/test-arcappliance-resgrp/providers/microsoft.extendedlocation/customlocations/testcustomlocation"),
			Type: to.Ptr("CustomLocation"),
		},
		Properties: &armhybridcontainerservice.VirtualNetworksProperties{
			InfraVnetProfile: &armhybridcontainerservice.VirtualNetworksPropertiesInfraVnetProfile{
				Hci: &armhybridcontainerservice.VirtualNetworksPropertiesInfraVnetProfileHci{
					MocGroup:    to.Ptr("target-group"),
					MocLocation: to.Ptr("MocLocation"),
					MocVnetName: to.Ptr("test-vnet"),
				},
			},
			VipPool: []*armhybridcontainerservice.VirtualNetworksPropertiesVipPoolItem{
				{
					EndIP:   to.Ptr("192.168.0.50"),
					StartIP: to.Ptr("192.168.0.10"),
				}},
			VmipPool: []*armhybridcontainerservice.VirtualNetworksPropertiesVmipPoolItem{
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
	// res.VirtualNetworks = armhybridcontainerservice.VirtualNetworks{
	// 	Name: to.Ptr("test-vnet-static"),
	// 	Type: to.Ptr("microsoft.hybridcontainerservice/virtualnetworks"),
	// 	ID: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourceGroups/test-arcappliance-resgrp/providers/Microsoft.HybridContainerService/virtualNetworks/test-vnet-static"),
	// 	Location: to.Ptr("westus"),
	// 	ExtendedLocation: &armhybridcontainerservice.VirtualNetworksExtendedLocation{
	// 		Name: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourcegroups/test-arcappliance-resgrp/providers/microsoft.extendedlocation/customlocations/testcustomlocation"),
	// 		Type: to.Ptr("CustomLocation"),
	// 	},
	// 	Properties: &armhybridcontainerservice.VirtualNetworksProperties{
	// 		DNSServers: []*string{
	// 			to.Ptr("192.168.0.1")},
	// 			Gateway: to.Ptr("192.168.0.1"),
	// 			InfraVnetProfile: &armhybridcontainerservice.VirtualNetworksPropertiesInfraVnetProfile{
	// 				Hci: &armhybridcontainerservice.VirtualNetworksPropertiesInfraVnetProfileHci{
	// 					MocGroup: to.Ptr("target-group"),
	// 					MocLocation: to.Ptr("MocLocation"),
	// 					MocVnetName: to.Ptr("test-vnet"),
	// 				},
	// 			},
	// 			IPAddressPrefix: to.Ptr("192.168.0.0/16"),
	// 			ProvisioningState: to.Ptr(armhybridcontainerservice.ProvisioningStateSucceeded),
	// 			VipPool: []*armhybridcontainerservice.VirtualNetworksPropertiesVipPoolItem{
	// 				{
	// 					EndIP: to.Ptr("192.168.0.50"),
	// 					StartIP: to.Ptr("192.168.0.10"),
	// 			}},
	// 			VmipPool: []*armhybridcontainerservice.VirtualNetworksPropertiesVmipPoolItem{
	// 				{
	// 					EndIP: to.Ptr("192.168.0.130"),
	// 					StartIP: to.Ptr("192.168.0.110"),
	// 			}},
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/a60468a0c5e2beb054680ae488fb9f92699f0a0d/specification/hybridaks/resource-manager/Microsoft.HybridContainerService/preview/2022-09-01-preview/examples/DeleteVirtualNetwork.json
func ExampleVirtualNetworksClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhybridcontainerservice.NewVirtualNetworksClient("a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Delete(ctx, "test-arcappliance-resgrp", "test-vnet-static", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/a60468a0c5e2beb054680ae488fb9f92699f0a0d/specification/hybridaks/resource-manager/Microsoft.HybridContainerService/preview/2022-09-01-preview/examples/UpdateVirtualNetwork.json
func ExampleVirtualNetworksClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhybridcontainerservice.NewVirtualNetworksClient("a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginUpdate(ctx, "test-arcappliance-resgrp", "test-vnet-static", armhybridcontainerservice.VirtualNetworksPatch{
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
	// res.VirtualNetworks = armhybridcontainerservice.VirtualNetworks{
	// 	Name: to.Ptr("test-vnet-static"),
	// 	Type: to.Ptr("microsoft.hybridcontainerservice/virtualnetworks"),
	// 	ID: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourceGroups/test-arcappliance-resgrp/providers/Microsoft.HybridContainerService/virtualNetworks/test-vnet-static"),
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 		"additionalProperties": to.Ptr("sample"),
	// 	},
	// 	ExtendedLocation: &armhybridcontainerservice.VirtualNetworksExtendedLocation{
	// 		Name: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourcegroups/test-arcappliance-resgrp/providers/microsoft.extendedlocation/customlocations/testcustomlocation"),
	// 		Type: to.Ptr("CustomLocation"),
	// 	},
	// 	Properties: &armhybridcontainerservice.VirtualNetworksProperties{
	// 		DNSServers: []*string{
	// 			to.Ptr("192.168.0.1")},
	// 			Gateway: to.Ptr("192.168.0.1"),
	// 			InfraVnetProfile: &armhybridcontainerservice.VirtualNetworksPropertiesInfraVnetProfile{
	// 				Hci: &armhybridcontainerservice.VirtualNetworksPropertiesInfraVnetProfileHci{
	// 					MocGroup: to.Ptr("target-group"),
	// 					MocLocation: to.Ptr("MocLocation"),
	// 					MocVnetName: to.Ptr("test-vnet"),
	// 				},
	// 			},
	// 			IPAddressPrefix: to.Ptr("192.168.0.0/16"),
	// 			ProvisioningState: to.Ptr(armhybridcontainerservice.ProvisioningStateSucceeded),
	// 			VipPool: []*armhybridcontainerservice.VirtualNetworksPropertiesVipPoolItem{
	// 				{
	// 					EndIP: to.Ptr("192.168.0.50"),
	// 					StartIP: to.Ptr("192.168.0.10"),
	// 			}},
	// 			VmipPool: []*armhybridcontainerservice.VirtualNetworksPropertiesVmipPoolItem{
	// 				{
	// 					EndIP: to.Ptr("192.168.0.130"),
	// 					StartIP: to.Ptr("192.168.0.110"),
	// 			}},
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/a60468a0c5e2beb054680ae488fb9f92699f0a0d/specification/hybridaks/resource-manager/Microsoft.HybridContainerService/preview/2022-09-01-preview/examples/ListVirtualNetworkByResourceGroup.json
func ExampleVirtualNetworksClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhybridcontainerservice.NewVirtualNetworksClient("a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListByResourceGroupPager("test-arcappliance-resgrp", nil)
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
		// 	Value: []*armhybridcontainerservice.VirtualNetworks{
		// 		{
		// 			Name: to.Ptr("test-vnet-static"),
		// 			Type: to.Ptr("microsoft.hybridcontainerservice/virtualnetworks"),
		// 			ID: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourceGroups/test-arcappliance-resgrp/providers/Microsoft.HybridContainerService/virtualNetworks/test-vnet-static"),
		// 			Location: to.Ptr("westus"),
		// 			ExtendedLocation: &armhybridcontainerservice.VirtualNetworksExtendedLocation{
		// 				Name: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourcegroups/test-arcappliance-resgrp/providers/microsoft.extendedlocation/customlocations/testcustomlocation"),
		// 				Type: to.Ptr("CustomLocation"),
		// 			},
		// 			Properties: &armhybridcontainerservice.VirtualNetworksProperties{
		// 				DNSServers: []*string{
		// 					to.Ptr("192.168.0.1")},
		// 					Gateway: to.Ptr("192.168.0.1"),
		// 					InfraVnetProfile: &armhybridcontainerservice.VirtualNetworksPropertiesInfraVnetProfile{
		// 						Hci: &armhybridcontainerservice.VirtualNetworksPropertiesInfraVnetProfileHci{
		// 							MocGroup: to.Ptr("target-group"),
		// 							MocLocation: to.Ptr("MocLocation"),
		// 							MocVnetName: to.Ptr("test-vnet"),
		// 						},
		// 					},
		// 					IPAddressPrefix: to.Ptr("192.168.0.0/16"),
		// 					ProvisioningState: to.Ptr(armhybridcontainerservice.ProvisioningStateSucceeded),
		// 					VipPool: []*armhybridcontainerservice.VirtualNetworksPropertiesVipPoolItem{
		// 						{
		// 							EndIP: to.Ptr("192.168.0.50"),
		// 							StartIP: to.Ptr("192.168.0.10"),
		// 					}},
		// 					VmipPool: []*armhybridcontainerservice.VirtualNetworksPropertiesVmipPoolItem{
		// 						{
		// 							EndIP: to.Ptr("192.168.0.130"),
		// 							StartIP: to.Ptr("192.168.0.110"),
		// 					}},
		// 				},
		// 		}},
		// 	}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/a60468a0c5e2beb054680ae488fb9f92699f0a0d/specification/hybridaks/resource-manager/Microsoft.HybridContainerService/preview/2022-09-01-preview/examples/ListVirtualNetworkBySubscription.json
func ExampleVirtualNetworksClient_NewListBySubscriptionPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armhybridcontainerservice.NewVirtualNetworksClient("a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListBySubscriptionPager(nil)
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
		// 	Value: []*armhybridcontainerservice.VirtualNetworks{
		// 		{
		// 			Name: to.Ptr("test-vnet-static"),
		// 			Type: to.Ptr("microsoft.hybridcontainerservice/virtualnetworks"),
		// 			ID: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourceGroups/test-arcappliance-resgrp/providers/Microsoft.HybridContainerService/virtualNetworks/test-vnet-static"),
		// 			Location: to.Ptr("westus"),
		// 			ExtendedLocation: &armhybridcontainerservice.VirtualNetworksExtendedLocation{
		// 				Name: to.Ptr("/subscriptions/a3e42606-29b1-4d7d-b1d9-9ff6b9d3c71b/resourcegroups/test-arcappliance-resgrp/providers/microsoft.extendedlocation/customlocations/testcustomlocation"),
		// 				Type: to.Ptr("CustomLocation"),
		// 			},
		// 			Properties: &armhybridcontainerservice.VirtualNetworksProperties{
		// 				DNSServers: []*string{
		// 					to.Ptr("192.168.0.1")},
		// 					Gateway: to.Ptr("192.168.0.1"),
		// 					InfraVnetProfile: &armhybridcontainerservice.VirtualNetworksPropertiesInfraVnetProfile{
		// 						Hci: &armhybridcontainerservice.VirtualNetworksPropertiesInfraVnetProfileHci{
		// 							MocGroup: to.Ptr("target-group"),
		// 							MocLocation: to.Ptr("MocLocation"),
		// 							MocVnetName: to.Ptr("test-vnet"),
		// 						},
		// 					},
		// 					IPAddressPrefix: to.Ptr("192.168.0.0/16"),
		// 					ProvisioningState: to.Ptr(armhybridcontainerservice.ProvisioningStateSucceeded),
		// 					VipPool: []*armhybridcontainerservice.VirtualNetworksPropertiesVipPoolItem{
		// 						{
		// 							EndIP: to.Ptr("192.168.0.50"),
		// 							StartIP: to.Ptr("192.168.0.10"),
		// 					}},
		// 					VmipPool: []*armhybridcontainerservice.VirtualNetworksPropertiesVmipPoolItem{
		// 						{
		// 							EndIP: to.Ptr("192.168.0.130"),
		// 							StartIP: to.Ptr("192.168.0.110"),
		// 					}},
		// 				},
		// 		}},
		// 	}
	}
}
