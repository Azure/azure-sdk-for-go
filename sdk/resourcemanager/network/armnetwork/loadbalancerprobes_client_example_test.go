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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/LoadBalancerProbeList.json
func ExampleLoadBalancerProbesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewLoadBalancerProbesClient().NewListPager("testrg", "lb", nil)
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
		// page.LoadBalancerProbeListResult = armnetwork.LoadBalancerProbeListResult{
		// 	Value: []*armnetwork.Probe{
		// 		{
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/loadBalancers/lb/probes/prlb"),
		// 			Name: to.Ptr("prlb"),
		// 			Type: to.Ptr("Microsoft.Network/loadBalancers/probes"),
		// 			Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
		// 			Properties: &armnetwork.ProbePropertiesFormat{
		// 				IntervalInSeconds: to.Ptr[int32](15),
		// 				LoadBalancingRules: []*armnetwork.SubResource{
		// 					{
		// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/loadBalancers/lb/loadBalancingRules/rulelb"),
		// 				}},
		// 				NumberOfProbes: to.Ptr[int32](2),
		// 				Port: to.Ptr[int32](80),
		// 				ProbeThreshold: to.Ptr[int32](1),
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				RequestPath: to.Ptr("healthcheck.aspx"),
		// 				Protocol: to.Ptr(armnetwork.ProbeProtocolHTTP),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/LoadBalancerProbeGet.json
func ExampleLoadBalancerProbesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewLoadBalancerProbesClient().Get(ctx, "testrg", "lb", "probe1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Probe = armnetwork.Probe{
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/loadBalancers/lb/probes/probe1"),
	// 	Name: to.Ptr("probe1"),
	// 	Type: to.Ptr("Microsoft.Network/loadBalancers/probes"),
	// 	Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
	// 	Properties: &armnetwork.ProbePropertiesFormat{
	// 		IntervalInSeconds: to.Ptr[int32](15),
	// 		LoadBalancingRules: []*armnetwork.SubResource{
	// 			{
	// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/testrg/providers/Microsoft.Network/loadBalancers/lb/loadBalancingRules/rulelb"),
	// 		}},
	// 		NumberOfProbes: to.Ptr[int32](2),
	// 		Port: to.Ptr[int32](80),
	// 		ProbeThreshold: to.Ptr[int32](1),
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		RequestPath: to.Ptr("healthcheck.aspx"),
	// 		Protocol: to.Ptr(armnetwork.ProbeProtocolHTTP),
	// 	},
	// }
}
