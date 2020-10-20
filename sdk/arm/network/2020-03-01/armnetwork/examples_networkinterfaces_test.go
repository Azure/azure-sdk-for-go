// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armnetwork_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-03-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

// This sample requires the following environment variables to be set correctly in order to run:
// AZURE_LOCATION: Azure location where the resource will be created.
// AZURE_RESOURCE_GROUP: Azure resource group to retrieve and create resources.
// AZURE_SUBSCRIPTION_ID: The subscription ID where the resource group exists.
// AZURE_TENANT_ID: The Azure Active Directory tenant (directory) ID of the service principal.
// AZURE_CLIENT_ID: The client (application) ID of the service principal.
// AZURE_CLIENT_SECRET: A client secret that was generated for the App Registration used to authenticate the client.
// AZURE_IP: The name of the public IP addresses created in Azure.
// AZURE_VNET: The Virtual Network that will be used in the examples.
// AZURE_SUBNET: The subnet that exists in the Virtual Network.

const (
	nicName = "sampleNIC"
)

func getNetworkInterfacesClient() armnetwork.NetworkInterfacesOperations {
	return armnetwork.NewNetworkInterfacesClient(
		armnetwork.NewDefaultClient(getCredential(), nil),
		subscriptionID)
}

func getSubnetClient() armnetwork.SubnetsOperations {
	return armnetwork.NewSubnetsClient(
		armnetwork.NewDefaultClient(getCredential(), nil),
		subscriptionID)
}

// Environment variables that are required to run this example:
// AZURE_SUBSCRIPTION_ID, AZURE_TENANT_ID, AZURE_CLIENT_ID, AZURE_CLIENT_SECRET,
// AZURE_RESOURCE_GROUP, AZURE_LOCATION, AZURE_IP, AZURE_VNET, AZURE_SUBNET.
func ExampleNetworkInterfacesOperations_BeginCreateOrUpdate() {
	ip := getIP()
	subnet := getSubnet()
	client := getNetworkInterfacesClient()
	nic, err := client.BeginCreateOrUpdate(
		context.Background(),
		resourceGroupName,
		nicName,
		armnetwork.NetworkInterface{
			Resource: armnetwork.Resource{
				Name:     to.StringPtr(nicName),
				Location: to.StringPtr(location),
			},
			Properties: &armnetwork.NetworkInterfacePropertiesFormat{
				IPConfigurations: &[]armnetwork.NetworkInterfaceIPConfiguration{
					{
						Name: to.StringPtr(nicName),
						Properties: &armnetwork.NetworkInterfaceIPConfigurationPropertiesFormat{
							Subnet:                    subnet,
							PrivateIPAllocationMethod: armnetwork.IPAllocationMethodDynamic.ToPtr(),
							PublicIPAddress:           ip,
						},
					},
				},
			},
		},
		nil)
	if err != nil {
		panic(err.Error())
	}
	result, err := nic.PollUntilDone(context.Background(), 1*time.Second)
	if err != nil {
		panic(err.Error())
	}
	if result != nil {
		fmt.Println(*result.NetworkInterface.Name)
	}
	// Output:
	// sampleNIC
}

// Environment variables that are required to run this example:
// AZURE_SUBSCRIPTION_ID, AZURE_TENANT_ID, AZURE_CLIENT_ID, AZURE_CLIENT_SECRET,
// AZURE_RESOURCE_GROUP.
func ExampleNetworkInterfacesOperations_Get() {
	client := getNetworkInterfacesClient()
	nic, err := client.Get(context.Background(), resourceGroupName, nicName, nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(*nic.NetworkInterface.Name)
	// Output:
	// sampleNIC
}

// Environment variables that are required to run this example:
// AZURE_SUBSCRIPTION_ID, AZURE_TENANT_ID, AZURE_CLIENT_ID, AZURE_CLIENT_SECRET,
// AZURE_RESOURCE_GROUP.
func ExampleNetworkInterfacesOperations_BeginDelete() {
	client := getNetworkInterfacesClient()
	nic, err := client.BeginDelete(context.Background(), resourceGroupName, nicName, nil)
	if err != nil {
		panic(err.Error())
	}
	res, err := nic.PollUntilDone(context.Background(), 5*time.Second)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res.StatusCode)
	// Output:
	// 200
}

func getIP() *armnetwork.PublicIPAddress {
	client := getPublicIPClient()
	resp, err := client.Get(context.Background(), resourceGroupName, ip, nil)
	if err != nil {
		panic(err.Error())
	}
	return resp.PublicIPAddress
}

func getSubnet() *armnetwork.Subnet {
	client := getSubnetClient()
	resp, err := client.Get(context.Background(), resourceGroupName, vnetName, subnetName, nil)
	if err != nil {
		panic(err.Error())
	}
	return resp.Subnet
}
