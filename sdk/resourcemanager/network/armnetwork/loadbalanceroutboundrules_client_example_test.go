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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4be63be2cabdebeb6974b22d89ed6fdb44392541/specification/network/resource-manager/Microsoft.Network/stable/2023-09-01/examples/LoadBalancerOutboundRuleList.json
func ExampleLoadBalancerOutboundRulesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewLoadBalancerOutboundRulesClient().NewListPager("testrg", "lb1", nil)
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
		// page.LoadBalancerOutboundRuleListResult = armnetwork.LoadBalancerOutboundRuleListResult{
		// 	Value: []*armnetwork.OutboundRule{
		// 		{
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/loadBalancers/lb1/outboundRules/rule1"),
		// 			Name: to.Ptr("rule1"),
		// 			Type: to.Ptr("Microsoft.Network/loadBalancers/outboundRules"),
		// 			Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
		// 			Properties: &armnetwork.OutboundRulePropertiesFormat{
		// 				AllocatedOutboundPorts: to.Ptr[int32](64),
		// 				BackendAddressPool: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/loadBalancers/lb1/backendAddressPools/bepool1"),
		// 				},
		// 				EnableTCPReset: to.Ptr(true),
		// 				FrontendIPConfigurations: []*armnetwork.SubResource{
		// 					{
		// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/loadBalancers/lb1/frontendIPConfigurations/lbfrontend"),
		// 				}},
		// 				IdleTimeoutInMinutes: to.Ptr[int32](15),
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				Protocol: to.Ptr(armnetwork.LoadBalancerOutboundRuleProtocolTCP),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4be63be2cabdebeb6974b22d89ed6fdb44392541/specification/network/resource-manager/Microsoft.Network/stable/2023-09-01/examples/LoadBalancerOutboundRuleGet.json
func ExampleLoadBalancerOutboundRulesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewLoadBalancerOutboundRulesClient().Get(ctx, "testrg", "lb1", "rule1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.OutboundRule = armnetwork.OutboundRule{
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/loadBalancers/lb1/outboundRules/rule1"),
	// 	Name: to.Ptr("rule1"),
	// 	Type: to.Ptr("Microsoft.Network/loadBalancers/outboundRules"),
	// 	Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
	// 	Properties: &armnetwork.OutboundRulePropertiesFormat{
	// 		AllocatedOutboundPorts: to.Ptr[int32](64),
	// 		BackendAddressPool: &armnetwork.SubResource{
	// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/loadBalancers/lb1/backendAddressPools/bepool1"),
	// 		},
	// 		EnableTCPReset: to.Ptr(true),
	// 		FrontendIPConfigurations: []*armnetwork.SubResource{
	// 			{
	// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/loadBalancers/lb1/frontendIPConfigurations/lbfrontend"),
	// 		}},
	// 		IdleTimeoutInMinutes: to.Ptr[int32](15),
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		Protocol: to.Ptr(armnetwork.LoadBalancerOutboundRuleProtocolTCP),
	// 	},
	// }
}
