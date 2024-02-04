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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/45f5b5a166c75a878d0f5404e74bd1855ff48894/specification/network/resource-manager/Microsoft.Network/stable/2023-09-01/examples/DefaultSecurityRuleList.json
func ExampleDefaultSecurityRulesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewDefaultSecurityRulesClient().NewListPager("testrg", "nsg1", nil)
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
		// page.SecurityRuleListResult = armnetwork.SecurityRuleListResult{
		// 	Value: []*armnetwork.SecurityRule{
		// 		{
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/networkSecurityGroups/nsg1/defaultSecurityRules/AllowVnetInBound"),
		// 			Name: to.Ptr("AllowVnetInBound"),
		// 			Properties: &armnetwork.SecurityRulePropertiesFormat{
		// 				Description: to.Ptr("Allow inbound traffic from all VMs in VNET"),
		// 				Access: to.Ptr(armnetwork.SecurityRuleAccessAllow),
		// 				DestinationAddressPrefix: to.Ptr("VirtualNetwork"),
		// 				DestinationAddressPrefixes: []*string{
		// 				},
		// 				DestinationPortRange: to.Ptr("*"),
		// 				DestinationPortRanges: []*string{
		// 				},
		// 				Direction: to.Ptr(armnetwork.SecurityRuleDirectionInbound),
		// 				Priority: to.Ptr[int32](65000),
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				SourceAddressPrefix: to.Ptr("VirtualNetwork"),
		// 				SourceAddressPrefixes: []*string{
		// 				},
		// 				SourcePortRange: to.Ptr("*"),
		// 				SourcePortRanges: []*string{
		// 				},
		// 				Protocol: to.Ptr(armnetwork.SecurityRuleProtocolAsterisk),
		// 			},
		// 		},
		// 		{
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/networkSecurityGroups/nsg1/defaultSecurityRules/AllowAzureLoadBalancerInBound"),
		// 			Name: to.Ptr("AllowAzureLoadBalancerInBound"),
		// 			Properties: &armnetwork.SecurityRulePropertiesFormat{
		// 				Description: to.Ptr("Allow inbound traffic from azure load balancer"),
		// 				Access: to.Ptr(armnetwork.SecurityRuleAccessAllow),
		// 				DestinationAddressPrefix: to.Ptr("*"),
		// 				DestinationAddressPrefixes: []*string{
		// 				},
		// 				DestinationPortRange: to.Ptr("*"),
		// 				DestinationPortRanges: []*string{
		// 				},
		// 				Direction: to.Ptr(armnetwork.SecurityRuleDirectionInbound),
		// 				Priority: to.Ptr[int32](65001),
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				SourceAddressPrefix: to.Ptr("AzureLoadBalancer"),
		// 				SourceAddressPrefixes: []*string{
		// 				},
		// 				SourcePortRange: to.Ptr("*"),
		// 				SourcePortRanges: []*string{
		// 				},
		// 				Protocol: to.Ptr(armnetwork.SecurityRuleProtocolAsterisk),
		// 			},
		// 		},
		// 		{
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/networkSecurityGroups/nsg1/defaultSecurityRules/DenyAllInBound"),
		// 			Name: to.Ptr("DenyAllInBound"),
		// 			Properties: &armnetwork.SecurityRulePropertiesFormat{
		// 				Description: to.Ptr("Deny all inbound traffic"),
		// 				Access: to.Ptr(armnetwork.SecurityRuleAccessDeny),
		// 				DestinationAddressPrefix: to.Ptr("*"),
		// 				DestinationAddressPrefixes: []*string{
		// 				},
		// 				DestinationPortRange: to.Ptr("*"),
		// 				DestinationPortRanges: []*string{
		// 				},
		// 				Direction: to.Ptr(armnetwork.SecurityRuleDirectionInbound),
		// 				Priority: to.Ptr[int32](65500),
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				SourceAddressPrefix: to.Ptr("*"),
		// 				SourceAddressPrefixes: []*string{
		// 				},
		// 				SourcePortRange: to.Ptr("*"),
		// 				SourcePortRanges: []*string{
		// 				},
		// 				Protocol: to.Ptr(armnetwork.SecurityRuleProtocolAsterisk),
		// 			},
		// 		},
		// 		{
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/networkSecurityGroups/nsg1/defaultSecurityRules/AllowVnetOutBound"),
		// 			Name: to.Ptr("AllowVnetOutBound"),
		// 			Properties: &armnetwork.SecurityRulePropertiesFormat{
		// 				Description: to.Ptr("Allow outbound traffic from all VMs to all VMs in VNET"),
		// 				Access: to.Ptr(armnetwork.SecurityRuleAccessAllow),
		// 				DestinationAddressPrefix: to.Ptr("VirtualNetwork"),
		// 				DestinationAddressPrefixes: []*string{
		// 				},
		// 				DestinationPortRange: to.Ptr("*"),
		// 				DestinationPortRanges: []*string{
		// 				},
		// 				Direction: to.Ptr(armnetwork.SecurityRuleDirectionOutbound),
		// 				Priority: to.Ptr[int32](65000),
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				SourceAddressPrefix: to.Ptr("VirtualNetwork"),
		// 				SourceAddressPrefixes: []*string{
		// 				},
		// 				SourcePortRange: to.Ptr("*"),
		// 				SourcePortRanges: []*string{
		// 				},
		// 				Protocol: to.Ptr(armnetwork.SecurityRuleProtocolAsterisk),
		// 			},
		// 		},
		// 		{
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/networkSecurityGroups/nsg1/defaultSecurityRules/AllowInternetOutBound"),
		// 			Name: to.Ptr("AllowInternetOutBound"),
		// 			Properties: &armnetwork.SecurityRulePropertiesFormat{
		// 				Description: to.Ptr("Allow outbound traffic from all VMs to Internet"),
		// 				Access: to.Ptr(armnetwork.SecurityRuleAccessAllow),
		// 				DestinationAddressPrefix: to.Ptr("Internet"),
		// 				DestinationAddressPrefixes: []*string{
		// 				},
		// 				DestinationPortRange: to.Ptr("*"),
		// 				DestinationPortRanges: []*string{
		// 				},
		// 				Direction: to.Ptr(armnetwork.SecurityRuleDirectionOutbound),
		// 				Priority: to.Ptr[int32](65001),
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				SourceAddressPrefix: to.Ptr("*"),
		// 				SourceAddressPrefixes: []*string{
		// 				},
		// 				SourcePortRange: to.Ptr("*"),
		// 				SourcePortRanges: []*string{
		// 				},
		// 				Protocol: to.Ptr(armnetwork.SecurityRuleProtocolAsterisk),
		// 			},
		// 		},
		// 		{
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/networkSecurityGroups/nsg1/defaultSecurityRules/DenyAllOutBound"),
		// 			Name: to.Ptr("DenyAllOutBound"),
		// 			Properties: &armnetwork.SecurityRulePropertiesFormat{
		// 				Description: to.Ptr("Deny all outbound traffic"),
		// 				Access: to.Ptr(armnetwork.SecurityRuleAccessDeny),
		// 				DestinationAddressPrefix: to.Ptr("*"),
		// 				DestinationAddressPrefixes: []*string{
		// 				},
		// 				DestinationPortRange: to.Ptr("*"),
		// 				DestinationPortRanges: []*string{
		// 				},
		// 				Direction: to.Ptr(armnetwork.SecurityRuleDirectionOutbound),
		// 				Priority: to.Ptr[int32](65500),
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				SourceAddressPrefix: to.Ptr("*"),
		// 				SourceAddressPrefixes: []*string{
		// 				},
		// 				SourcePortRange: to.Ptr("*"),
		// 				SourcePortRanges: []*string{
		// 				},
		// 				Protocol: to.Ptr(armnetwork.SecurityRuleProtocolAsterisk),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/45f5b5a166c75a878d0f5404e74bd1855ff48894/specification/network/resource-manager/Microsoft.Network/stable/2023-09-01/examples/DefaultSecurityRuleGet.json
func ExampleDefaultSecurityRulesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewDefaultSecurityRulesClient().Get(ctx, "testrg", "nsg1", "AllowVnetInBound", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SecurityRule = armnetwork.SecurityRule{
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/networkSecurityGroups/nsg1/defaultSecurityRules/AllowVnetInBound"),
	// 	Name: to.Ptr("AllowVnetInBound"),
	// 	Properties: &armnetwork.SecurityRulePropertiesFormat{
	// 		Description: to.Ptr("Allow inbound traffic from all VMs in VNET"),
	// 		Access: to.Ptr(armnetwork.SecurityRuleAccessAllow),
	// 		DestinationAddressPrefix: to.Ptr("VirtualNetwork"),
	// 		DestinationAddressPrefixes: []*string{
	// 		},
	// 		DestinationPortRange: to.Ptr("*"),
	// 		DestinationPortRanges: []*string{
	// 		},
	// 		Direction: to.Ptr(armnetwork.SecurityRuleDirectionInbound),
	// 		Priority: to.Ptr[int32](65000),
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		SourceAddressPrefix: to.Ptr("VirtualNetwork"),
	// 		SourceAddressPrefixes: []*string{
	// 		},
	// 		SourcePortRange: to.Ptr("*"),
	// 		SourcePortRanges: []*string{
	// 		},
	// 		Protocol: to.Ptr(armnetwork.SecurityRuleProtocolAsterisk),
	// 	},
	// }
}
