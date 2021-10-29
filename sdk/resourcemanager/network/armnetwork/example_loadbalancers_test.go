//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armnetwork_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

func ExampleLoadBalancersClient_BeginCreateOrUpdate() {
	// populate the following variables with your specific information
	resourceGroupName := "<resource group name>"
	location := "<Azure location>"
	loadBalancerName := "<load balancer name>"
	frontEndIPConfigName := "<frontend IP config name>"
	backEndAddressPoolName := "<backend address pool name>"
	probeName := "<probe name>"
	idPrefix := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers", "<subscription ID>", resourceGroupName)
	// Replace the "ipAddress" variable with a call to armnetwork.PublicIPAddressesClient.Get to retreive the
	// public IP address instance that will be assigned to the load balancer.
	var ipAddress *armnetwork.PublicIPAddress

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewLoadBalancersClient("<subscription ID>", cred, nil)
	poller, err := client.BeginCreateOrUpdate(
		context.Background(),
		resourceGroupName,
		loadBalancerName,
		armnetwork.LoadBalancer{
			Resource: armnetwork.Resource{
				Name:     to.StringPtr(loadBalancerName),
				Location: to.StringPtr(location),
			},
			Properties: &armnetwork.LoadBalancerPropertiesFormat{
				FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
					{
						Name: to.StringPtr(frontEndIPConfigName),
						Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
							PrivateIPAllocationMethod: armnetwork.IPAllocationMethodDynamic.ToPtr(),
							PublicIPAddress:           ipAddress,
						},
					},
				},
				BackendAddressPools: []*armnetwork.BackendAddressPool{
					{
						Name: &backEndAddressPoolName,
					},
				},
				Probes: []*armnetwork.Probe{
					{
						Name: &probeName,
						Properties: &armnetwork.ProbePropertiesFormat{
							Protocol:          armnetwork.ProbeProtocolHTTP.ToPtr(),
							Port:              to.Int32Ptr(80),
							IntervalInSeconds: to.Int32Ptr(15),
							NumberOfProbes:    to.Int32Ptr(4),
							RequestPath:       to.StringPtr("healthprobe.aspx"),
						},
					},
				},
				LoadBalancingRules: []*armnetwork.LoadBalancingRule{
					{
						Name: to.StringPtr("lbRule"),
						Properties: &armnetwork.LoadBalancingRulePropertiesFormat{
							Protocol:             armnetwork.TransportProtocolTCP.ToPtr(),
							FrontendPort:         to.Int32Ptr(80),
							BackendPort:          to.Int32Ptr(80),
							IdleTimeoutInMinutes: to.Int32Ptr(4),
							EnableFloatingIP:     to.BoolPtr(false),
							LoadDistribution:     armnetwork.LoadDistributionDefault.ToPtr(),
							FrontendIPConfiguration: &armnetwork.SubResource{
								ID: to.StringPtr(fmt.Sprintf("/%s/%s/frontendIPConfigurations/%s", idPrefix, loadBalancerName, frontEndIPConfigName)),
							},
							BackendAddressPool: &armnetwork.SubResource{
								ID: to.StringPtr(fmt.Sprintf("/%s/%s/backendAddressPools/%s", idPrefix, loadBalancerName, backEndAddressPoolName)),
							},
							Probe: &armnetwork.SubResource{
								ID: to.StringPtr(fmt.Sprintf("/%s/%s/probes/%s", idPrefix, loadBalancerName, probeName)),
							},
						},
					},
				},
				InboundNatRules: []*armnetwork.InboundNatRule{
					{
						Name: to.StringPtr("natRule1"),
						Properties: &armnetwork.InboundNatRulePropertiesFormat{
							Protocol:             armnetwork.TransportProtocolTCP.ToPtr(),
							FrontendPort:         to.Int32Ptr(21),
							BackendPort:          to.Int32Ptr(22),
							EnableFloatingIP:     to.BoolPtr(false),
							IdleTimeoutInMinutes: to.Int32Ptr(4),
							FrontendIPConfiguration: &armnetwork.SubResource{
								ID: to.StringPtr(fmt.Sprintf("/%s/%s/frontendIPConfigurations/%s", idPrefix, loadBalancerName, frontEndIPConfigName)),
							},
						},
					},
					{
						Name: to.StringPtr("natRule2"),
						Properties: &armnetwork.InboundNatRulePropertiesFormat{
							Protocol:             armnetwork.TransportProtocolTCP.ToPtr(),
							FrontendPort:         to.Int32Ptr(23),
							BackendPort:          to.Int32Ptr(22),
							EnableFloatingIP:     to.BoolPtr(false),
							IdleTimeoutInMinutes: to.Int32Ptr(4),
							FrontendIPConfiguration: &armnetwork.SubResource{
								ID: to.StringPtr(fmt.Sprintf("/%s/%s/frontendIPConfigurations/%s", idPrefix, loadBalancerName, frontEndIPConfigName)),
							},
						},
					},
				},
			},
		},
		nil,
	)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	resp, err := poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}
	log.Printf("load balancer ID: %v", *resp.LoadBalancer.ID)
}

func ExampleLoadBalancersClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewLoadBalancersClient("<subscription ID>", cred, nil)
	resp, err := client.Get(context.Background(), "<resource group name>", "<load balancer name>", nil)
	if err != nil {
		log.Fatalf("failed to get resource: %v", err)
	}
	log.Printf("load balancer ID: %v", *resp.LoadBalancer.ID)
}

func ExampleLoadBalancersClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewLoadBalancersClient("<subscription ID>", cred, nil)
	resp, err := client.BeginDelete(context.Background(), "<resource group name>", "<load balancer name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = resp.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to delete resource: %v", err)
	}
}
