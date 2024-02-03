//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armnetwork_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v5"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4be63be2cabdebeb6974b22d89ed6fdb44392541/specification/network/resource-manager/Microsoft.Network/stable/2023-09-01/examples/LoadBalancerNetworkInterfaceListSimple.json
func ExampleLoadBalancerNetworkInterfacesClient_NewListPager_loadBalancerNetworkInterfaceListSimple() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewLoadBalancerNetworkInterfacesClient().NewListPager("testrg", "lb", nil)
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
		// page.InterfaceListResult = armnetwork.InterfaceListResult{
		// 	Value: []*armnetwork.Interface{
		// 		{
		// 			Name: to.Ptr("mynic"),
		// 			Type: to.Ptr("Microsoft.Network/networkInterfaces"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/networkInterfaces/mynic"),
		// 			Location: to.Ptr("westus"),
		// 			Etag: to.Ptr("W/\\\"00000000-0000-0000-0000-000000000000\\\""),
		// 			Properties: &armnetwork.InterfacePropertiesFormat{
		// 				DNSSettings: &armnetwork.InterfaceDNSSettings{
		// 					AppliedDNSServers: []*string{
		// 					},
		// 					DNSServers: []*string{
		// 					},
		// 				},
		// 				EnableAcceleratedNetworking: to.Ptr(false),
		// 				EnableIPForwarding: to.Ptr(false),
		// 				IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
		// 					{
		// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/networkInterfaces/mynic/ipConfigurations/ipconfig1"),
		// 						Name: to.Ptr("ipconfig1"),
		// 						Etag: to.Ptr("W/\\\"00000000-0000-0000-0000-000000000000\\\""),
		// 						Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
		// 							LoadBalancerBackendAddressPools: []*armnetwork.BackendAddressPool{
		// 								{
		// 									ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/loadBalancers/lb/backendAddressPools/bepool1"),
		// 							}},
		// 							LoadBalancerInboundNatRules: []*armnetwork.InboundNatRule{
		// 								{
		// 									ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/loadBalancers/lb/inboundNatRules/inbound1"),
		// 							}},
		// 							PrivateIPAddress: to.Ptr("10.0.1.4"),
		// 							PrivateIPAddressVersion: to.Ptr(armnetwork.IPVersionIPv4),
		// 							PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
		// 							ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 							Subnet: &armnetwork.Subnet{
		// 								ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/myVirtualNetwork/subnets/frontendSubnet"),
		// 							},
		// 						},
		// 				}},
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4be63be2cabdebeb6974b22d89ed6fdb44392541/specification/network/resource-manager/Microsoft.Network/stable/2023-09-01/examples/LoadBalancerNetworkInterfaceListVmss.json
func ExampleLoadBalancerNetworkInterfacesClient_NewListPager_loadBalancerNetworkInterfaceListVmss() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewLoadBalancerNetworkInterfacesClient().NewListPager("testrg", "lb", nil)
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
		// page.InterfaceListResult = armnetwork.InterfaceListResult{
		// 	Value: []*armnetwork.Interface{
		// 		{
		// 			Name: to.Ptr("vmss1Nic"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Compute/virtualMachineScaleSets/vmss1/virtualMachines/0/networkInterfaces/vmss1Nic"),
		// 			Etag: to.Ptr("W/\\\"00000000-0000-0000-0000-000000000000\\\""),
		// 			Properties: &armnetwork.InterfacePropertiesFormat{
		// 				DNSSettings: &armnetwork.InterfaceDNSSettings{
		// 					AppliedDNSServers: []*string{
		// 					},
		// 					DNSServers: []*string{
		// 					},
		// 					InternalDomainNameSuffix: to.Ptr("aaaaaaaaaaaaaaaaaaaaaaaaaa.dx.internal.cloudapp.net"),
		// 				},
		// 				EnableAcceleratedNetworking: to.Ptr(false),
		// 				EnableIPForwarding: to.Ptr(false),
		// 				IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
		// 					{
		// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Compute/virtualMachineScaleSets/vmss1/virtualMachines/0/networkInterfaces/vmss1Nic/ipConfigurations/vmss1IpConfig"),
		// 						Name: to.Ptr("vmss1IpConfig"),
		// 						Etag: to.Ptr("W/\\\"00000000-0000-0000-0000-000000000000\\\""),
		// 						Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
		// 							LoadBalancerBackendAddressPools: []*armnetwork.BackendAddressPool{
		// 								{
		// 									ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/loadBalancers/lb/backendAddressPools/bepool"),
		// 							}},
		// 							LoadBalancerInboundNatRules: []*armnetwork.InboundNatRule{
		// 								{
		// 									ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/loadBalancers/lb/inboundNatRules/natpool.0"),
		// 							}},
		// 							Primary: to.Ptr(true),
		// 							PrivateIPAddress: to.Ptr("10.0.0.4"),
		// 							PrivateIPAddressVersion: to.Ptr(armnetwork.IPVersionIPv4),
		// 							PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
		// 							ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 							Subnet: &armnetwork.Subnet{
		// 								ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/vmss1Vnet/subnets/default"),
		// 							},
		// 						},
		// 				}},
		// 				MacAddress: to.Ptr("00-00-00-00-00-00"),
		// 				Primary: to.Ptr(true),
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 				VirtualMachine: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Compute/virtualMachineScaleSets/vmss1/virtualMachines/0"),
		// 				},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("vmss1Nic"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Compute/virtualMachineScaleSets/vmss1/virtualMachines/1/networkInterfaces/vmss1Nic"),
		// 			Etag: to.Ptr("W/\\\"00000000-0000-0000-0000-000000000000\\\""),
		// 			Properties: &armnetwork.InterfacePropertiesFormat{
		// 				DNSSettings: &armnetwork.InterfaceDNSSettings{
		// 					AppliedDNSServers: []*string{
		// 					},
		// 					DNSServers: []*string{
		// 					},
		// 					InternalDomainNameSuffix: to.Ptr("aaaaaaaaaaaaaaaaaaaaaaaaaa.dx.internal.cloudapp.net"),
		// 				},
		// 				EnableAcceleratedNetworking: to.Ptr(false),
		// 				EnableIPForwarding: to.Ptr(false),
		// 				IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
		// 					{
		// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Compute/virtualMachineScaleSets/vmss1/virtualMachines/1/networkInterfaces/vmss1Nic/ipConfigurations/vmss1IpConfig"),
		// 						Name: to.Ptr("vmss1IpConfig"),
		// 						Etag: to.Ptr("W/\\\"00000000-0000-0000-0000-000000000000\\\""),
		// 						Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
		// 							LoadBalancerBackendAddressPools: []*armnetwork.BackendAddressPool{
		// 								{
		// 									ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/loadBalancers/lb/backendAddressPools/bepool"),
		// 							}},
		// 							LoadBalancerInboundNatRules: []*armnetwork.InboundNatRule{
		// 								{
		// 									ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/loadBalancers/lb/inboundNatRules/natpool.1"),
		// 							}},
		// 							Primary: to.Ptr(true),
		// 							PrivateIPAddress: to.Ptr("10.0.0.5"),
		// 							PrivateIPAddressVersion: to.Ptr(armnetwork.IPVersionIPv4),
		// 							PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
		// 							ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 							Subnet: &armnetwork.Subnet{
		// 								ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/virtualNetworks/vmss1Vnet/subnets/default"),
		// 							},
		// 						},
		// 				}},
		// 				MacAddress: to.Ptr("00-00-00-00-00-00"),
		// 				Primary: to.Ptr(true),
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 				VirtualMachine: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Compute/virtualMachineScaleSets/vmss1/virtualMachines/1"),
		// 				},
		// 			},
		// 	}},
		// }
	}
}
