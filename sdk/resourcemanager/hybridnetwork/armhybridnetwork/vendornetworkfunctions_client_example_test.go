//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armhybridnetwork_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybridnetwork/armhybridnetwork/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/792db17291c758b2bfdbbc0d35d0e2f5b5a1bd05/specification/hybridnetwork/resource-manager/Microsoft.HybridNetwork/preview/2022-01-01-preview/examples/VendorNfGet.json
func ExampleVendorNetworkFunctionsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVendorNetworkFunctionsClient().Get(ctx, "eastus", "testVendor", "testServiceKey", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VendorNetworkFunction = armhybridnetwork.VendorNetworkFunction{
	// 	Name: to.Ptr("testServiceKey"),
	// 	Type: to.Ptr("Microsoft.HybridNetwork/locations/vendors/networkFunctions"),
	// 	ID: to.Ptr("/subscriptions/subid/providers/Microsoft.HybridNetwork/locations/eastus/vendors/testVendor/networkFunctions/testServiceKey"),
	// 	Properties: &armhybridnetwork.VendorNetworkFunctionPropertiesFormat{
	// 		NetworkFunctionVendorConfigurations: []*armhybridnetwork.NetworkFunctionVendorConfiguration{
	// 			{
	// 				NetworkInterfaces: []*armhybridnetwork.NetworkInterface{
	// 					{
	// 						IPConfigurations: []*armhybridnetwork.NetworkInterfaceIPConfiguration{
	// 							{
	// 								Gateway: to.Ptr(""),
	// 								IPAddress: to.Ptr(""),
	// 								IPAllocationMethod: to.Ptr(armhybridnetwork.IPAllocationMethodDynamic),
	// 								IPVersion: to.Ptr(armhybridnetwork.IPVersionIPv4),
	// 								Subnet: to.Ptr(""),
	// 						}},
	// 						MacAddress: to.Ptr(""),
	// 						NetworkInterfaceName: to.Ptr("nic1"),
	// 						VMSwitchType: to.Ptr(armhybridnetwork.VMSwitchTypeManagement),
	// 					},
	// 					{
	// 						IPConfigurations: []*armhybridnetwork.NetworkInterfaceIPConfiguration{
	// 							{
	// 								Gateway: to.Ptr(""),
	// 								IPAddress: to.Ptr(""),
	// 								IPAllocationMethod: to.Ptr(armhybridnetwork.IPAllocationMethodDynamic),
	// 								IPVersion: to.Ptr(armhybridnetwork.IPVersionIPv4),
	// 								Subnet: to.Ptr(""),
	// 						}},
	// 						MacAddress: to.Ptr("DC-97-F8-79-16-7D"),
	// 						NetworkInterfaceName: to.Ptr("nic2"),
	// 						VMSwitchType: to.Ptr(armhybridnetwork.VMSwitchTypeWan),
	// 				}},
	// 				OSProfile: &armhybridnetwork.OsProfile{
	// 					AdminUsername: to.Ptr("dummyuser"),
	// 					CustomData: to.Ptr("base-64 encoded string of custom data"),
	// 					LinuxConfiguration: &armhybridnetwork.LinuxConfiguration{
	// 						SSH: &armhybridnetwork.SSHConfiguration{
	// 							PublicKeys: []*armhybridnetwork.SSHPublicKey{
	// 								{
	// 									Path: to.Ptr("home/user/.ssh/authorized_keys"),
	// 									KeyData: to.Ptr("ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAgEAwrr66r8n6B8Y0zMF3dOpXEapIQD9DiYQ6D6/zwor9o39jSkHNiMMER/GETBbzP83LOcekm02aRjo55ArO7gPPVvCXbrirJu9pkm4AC4BBre5xSLS= user@constoso-DSH"),
	// 							}},
	// 						},
	// 					},
	// 				},
	// 				RoleName: to.Ptr("testRole"),
	// 				UserDataParameters: map[string]any{
	// 				},
	// 		}},
	// 		ProvisioningState: to.Ptr(armhybridnetwork.ProvisioningStateSucceeded),
	// 		SKUName: to.Ptr("testSku"),
	// 		SKUType: to.Ptr(armhybridnetwork.SKUTypeSDWAN),
	// 		VendorProvisioningState: to.Ptr(armhybridnetwork.VendorProvisioningStateProvisioning),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/792db17291c758b2bfdbbc0d35d0e2f5b5a1bd05/specification/hybridnetwork/resource-manager/Microsoft.HybridNetwork/preview/2022-01-01-preview/examples/VendorNfCreate.json
func ExampleVendorNetworkFunctionsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVendorNetworkFunctionsClient().BeginCreateOrUpdate(ctx, "eastus", "testVendor", "testServiceKey", armhybridnetwork.VendorNetworkFunction{
		Properties: &armhybridnetwork.VendorNetworkFunctionPropertiesFormat{
			NetworkFunctionVendorConfigurations: []*armhybridnetwork.NetworkFunctionVendorConfiguration{
				{
					NetworkInterfaces: []*armhybridnetwork.NetworkInterface{
						{
							IPConfigurations: []*armhybridnetwork.NetworkInterfaceIPConfiguration{
								{
									Gateway:            to.Ptr(""),
									IPAddress:          to.Ptr(""),
									IPAllocationMethod: to.Ptr(armhybridnetwork.IPAllocationMethodDynamic),
									IPVersion:          to.Ptr(armhybridnetwork.IPVersionIPv4),
									Subnet:             to.Ptr(""),
								}},
							MacAddress:           to.Ptr(""),
							NetworkInterfaceName: to.Ptr("nic1"),
							VMSwitchType:         to.Ptr(armhybridnetwork.VMSwitchTypeManagement),
						},
						{
							IPConfigurations: []*armhybridnetwork.NetworkInterfaceIPConfiguration{
								{
									Gateway:            to.Ptr(""),
									IPAddress:          to.Ptr(""),
									IPAllocationMethod: to.Ptr(armhybridnetwork.IPAllocationMethodDynamic),
									IPVersion:          to.Ptr(armhybridnetwork.IPVersionIPv4),
									Subnet:             to.Ptr(""),
								}},
							MacAddress:           to.Ptr("DC-97-F8-79-16-7D"),
							NetworkInterfaceName: to.Ptr("nic2"),
							VMSwitchType:         to.Ptr(armhybridnetwork.VMSwitchTypeWan),
						}},
					OSProfile: &armhybridnetwork.OsProfile{
						AdminUsername: to.Ptr("dummyuser"),
						CustomData:    to.Ptr("base-64 encoded string of custom data"),
						LinuxConfiguration: &armhybridnetwork.LinuxConfiguration{
							SSH: &armhybridnetwork.SSHConfiguration{
								PublicKeys: []*armhybridnetwork.SSHPublicKey{
									{
										Path:    to.Ptr("home/user/.ssh/authorized_keys"),
										KeyData: to.Ptr("ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAgEAwrr66r8n6B8Y0zMF3dOpXEapIQD9DiYQ6D6/zwor9o39jSkHNiMMER/GETBbzP83LOcekm02aRjo55ArO7gPPVvCXbrirJu9pkm4AC4BBre5xSLS= user@constoso-DSH"),
									}},
							},
						},
					},
					RoleName:           to.Ptr("testRole"),
					UserDataParameters: map[string]any{},
				}},
			SKUType:                 to.Ptr(armhybridnetwork.SKUTypeSDWAN),
			VendorProvisioningState: to.Ptr(armhybridnetwork.VendorProvisioningStateProvisioning),
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
	// res.VendorNetworkFunction = armhybridnetwork.VendorNetworkFunction{
	// 	Name: to.Ptr("testServiceKey"),
	// 	Type: to.Ptr("Microsoft.HybridNetwork/locations/vendors/networkFunctions"),
	// 	ID: to.Ptr("/subscriptions/subid/providers/Microsoft.HybridNetwork/locations/eastus/vendors/testVendor/networkFunctions/testServiceKey"),
	// 	Properties: &armhybridnetwork.VendorNetworkFunctionPropertiesFormat{
	// 		NetworkFunctionVendorConfigurations: []*armhybridnetwork.NetworkFunctionVendorConfiguration{
	// 			{
	// 				NetworkInterfaces: []*armhybridnetwork.NetworkInterface{
	// 					{
	// 						IPConfigurations: []*armhybridnetwork.NetworkInterfaceIPConfiguration{
	// 							{
	// 								Gateway: to.Ptr(""),
	// 								IPAddress: to.Ptr(""),
	// 								IPAllocationMethod: to.Ptr(armhybridnetwork.IPAllocationMethodDynamic),
	// 								IPVersion: to.Ptr(armhybridnetwork.IPVersionIPv4),
	// 								Subnet: to.Ptr(""),
	// 						}},
	// 						MacAddress: to.Ptr(""),
	// 						NetworkInterfaceName: to.Ptr("nic1"),
	// 						VMSwitchType: to.Ptr(armhybridnetwork.VMSwitchTypeManagement),
	// 					},
	// 					{
	// 						IPConfigurations: []*armhybridnetwork.NetworkInterfaceIPConfiguration{
	// 							{
	// 								Gateway: to.Ptr(""),
	// 								IPAddress: to.Ptr(""),
	// 								IPAllocationMethod: to.Ptr(armhybridnetwork.IPAllocationMethodDynamic),
	// 								IPVersion: to.Ptr(armhybridnetwork.IPVersionIPv4),
	// 								Subnet: to.Ptr(""),
	// 						}},
	// 						MacAddress: to.Ptr("DC-97-F8-79-16-7D"),
	// 						NetworkInterfaceName: to.Ptr("nic2"),
	// 						VMSwitchType: to.Ptr(armhybridnetwork.VMSwitchTypeWan),
	// 				}},
	// 				OSProfile: &armhybridnetwork.OsProfile{
	// 					AdminUsername: to.Ptr("dummyuser"),
	// 					CustomData: to.Ptr("base-64 encoded string of custom data"),
	// 					LinuxConfiguration: &armhybridnetwork.LinuxConfiguration{
	// 						SSH: &armhybridnetwork.SSHConfiguration{
	// 							PublicKeys: []*armhybridnetwork.SSHPublicKey{
	// 								{
	// 									Path: to.Ptr("home/user/.ssh/authorized_keys"),
	// 									KeyData: to.Ptr("ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAgEAwrr66r8n6B8Y0zMF3dOpXEapIQD9DiYQ6D6/zwor9o39jSkHNiMMER/GETBbzP83LOcekm02aRjo55ArO7gPPVvCXbrirJu9pkm4AC4BBre5xSLS= user@constoso-DSH"),
	// 							}},
	// 						},
	// 					},
	// 				},
	// 				RoleName: to.Ptr("testRole"),
	// 				UserDataParameters: map[string]any{
	// 				},
	// 		}},
	// 		ProvisioningState: to.Ptr(armhybridnetwork.ProvisioningStateSucceeded),
	// 		SKUName: to.Ptr("testSku"),
	// 		SKUType: to.Ptr(armhybridnetwork.SKUTypeSDWAN),
	// 		VendorProvisioningState: to.Ptr(armhybridnetwork.VendorProvisioningStateProvisioning),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/792db17291c758b2bfdbbc0d35d0e2f5b5a1bd05/specification/hybridnetwork/resource-manager/Microsoft.HybridNetwork/preview/2022-01-01-preview/examples/VendorNfListByVendor.json
func ExampleVendorNetworkFunctionsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhybridnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewVendorNetworkFunctionsClient().NewListPager("eastus", "testVendor", &armhybridnetwork.VendorNetworkFunctionsClientListOptions{Filter: nil})
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
		// page.VendorNetworkFunctionListResult = armhybridnetwork.VendorNetworkFunctionListResult{
		// 	Value: []*armhybridnetwork.VendorNetworkFunction{
		// 		{
		// 			Name: to.Ptr("TestServiceKey"),
		// 			Type: to.Ptr("Microsoft.HybridNetwork/locations/vendors/networkFunctions"),
		// 			ID: to.Ptr("/subscriptions/subid/providers/Microsoft.HybridNetwork/locations/eastus/vendors/testVendor/networkFunctions/testServiceKey"),
		// 			Properties: &armhybridnetwork.VendorNetworkFunctionPropertiesFormat{
		// 				NetworkFunctionVendorConfigurations: []*armhybridnetwork.NetworkFunctionVendorConfiguration{
		// 					{
		// 						NetworkInterfaces: []*armhybridnetwork.NetworkInterface{
		// 							{
		// 								IPConfigurations: []*armhybridnetwork.NetworkInterfaceIPConfiguration{
		// 									{
		// 										Gateway: to.Ptr(""),
		// 										IPAddress: to.Ptr(""),
		// 										IPAllocationMethod: to.Ptr(armhybridnetwork.IPAllocationMethodDynamic),
		// 										IPVersion: to.Ptr(armhybridnetwork.IPVersionIPv4),
		// 										Subnet: to.Ptr(""),
		// 								}},
		// 								MacAddress: to.Ptr(""),
		// 								NetworkInterfaceName: to.Ptr("nic1"),
		// 								VMSwitchType: to.Ptr(armhybridnetwork.VMSwitchTypeManagement),
		// 							},
		// 							{
		// 								IPConfigurations: []*armhybridnetwork.NetworkInterfaceIPConfiguration{
		// 									{
		// 										Gateway: to.Ptr(""),
		// 										IPAddress: to.Ptr(""),
		// 										IPAllocationMethod: to.Ptr(armhybridnetwork.IPAllocationMethodDynamic),
		// 										IPVersion: to.Ptr(armhybridnetwork.IPVersionIPv4),
		// 										Subnet: to.Ptr(""),
		// 								}},
		// 								MacAddress: to.Ptr("DC-97-F8-79-16-7D"),
		// 								NetworkInterfaceName: to.Ptr("nic2"),
		// 								VMSwitchType: to.Ptr(armhybridnetwork.VMSwitchTypeWan),
		// 						}},
		// 						OSProfile: &armhybridnetwork.OsProfile{
		// 							AdminUsername: to.Ptr("dummyuser"),
		// 							CustomData: to.Ptr("base-64 encoded string of custom data"),
		// 							LinuxConfiguration: &armhybridnetwork.LinuxConfiguration{
		// 								SSH: &armhybridnetwork.SSHConfiguration{
		// 									PublicKeys: []*armhybridnetwork.SSHPublicKey{
		// 										{
		// 											Path: to.Ptr("home/user/.ssh/authorized_keys"),
		// 											KeyData: to.Ptr("ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAgEAwrr66r8n6B8Y0zMF3dOpXEapIQD9DiYQ6D6/zwor9o39jSkHNiMMER/GETBbzP83LOcekm02aRjo55ArO7gPPVvCXbrirJu9pkm4AC4BBre5xSLS= user@constoso-DSH"),
		// 									}},
		// 								},
		// 							},
		// 						},
		// 						RoleName: to.Ptr("testRole"),
		// 						UserDataParameters: map[string]any{
		// 						},
		// 				}},
		// 				ProvisioningState: to.Ptr(armhybridnetwork.ProvisioningStateSucceeded),
		// 				SKUName: to.Ptr("testSku"),
		// 				SKUType: to.Ptr(armhybridnetwork.SKUTypeSDWAN),
		// 				VendorProvisioningState: to.Ptr(armhybridnetwork.VendorProvisioningStateProvisioning),
		// 			},
		// 	}},
		// }
	}
}
