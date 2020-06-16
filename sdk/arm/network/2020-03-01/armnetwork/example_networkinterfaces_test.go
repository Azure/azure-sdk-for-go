// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armnetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

const (
	nicName = "sampleNIC"
)

func getNetworkInterfacesClient() NetworkInterfacesOperations {
	nicClient, err := NewDefaultClient(getCredential(), nil)
	if err != nil {
		panic(err)
	}
	return nicClient.NetworkInterfacesOperations(subscriptionID)
}

func getSubnetClient() SubnetsOperations {
	nicClient, err := NewDefaultClient(getCredential(), nil)
	if err != nil {
		panic(err)
	}
	return nicClient.SubnetsOperations(subscriptionID)
}

// Environment variables that are required to run this example:
// AZURE_SUBSCRIPTION_ID, AZURE_TENANT_ID, AZURE_CLIENT_ID, AZURE_CLIENT_SECRET,
// AZURE_RESOURCE_GROUP, AZURE_LOCATION, AZURE_IP.
// See example_config.go for information about the environment variables.
func ExampleNetworkInterfacesOperations_BeginCreateOrUpdate() {
	ip := getIP()
	subnet := getSubnet()
	client := getNetworkInterfacesClient()
	nic, err := client.BeginCreateOrUpdate(
		context.Background(),
		resourceGroupName,
		nicName,
		NetworkInterface{
			Resource: Resource{
				Name:     to.StringPtr(nicName),
				Location: to.StringPtr(location),
			},
			Properties: &NetworkInterfacePropertiesFormat{
				IPConfigurations: &[]NetworkInterfaceIPConfiguration{
					{
						Name: to.StringPtr(nicName),
						Properties: &NetworkInterfaceIPConfigurationPropertiesFormat{
							Subnet:                    subnet,
							PrivateIPAllocationMethod: IPAllocationMethodDynamic.ToPtr(),
							PublicIPAddress:           ip,
						},
					},
				},
			},
		})
	if err != nil {
		panic(err)
	}
	result, err := nic.PollUntilDone(context.Background(), 1*time.Second)
	if err != nil {
		panic(err)
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
// See example_config.go for information about the environment variables.
func ExampleNetworkInterfacesOperations_Get() {
	client := getNetworkInterfacesClient()
	nic, err := client.Get(context.Background(), resourceGroupName, nicName, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(*nic.NetworkInterface.Name)
	// Output:
	// sampleNIC
}

// Environment variables that are required to run this example:
// AZURE_SUBSCRIPTION_ID, AZURE_TENANT_ID, AZURE_CLIENT_ID, AZURE_CLIENT_SECRET,
// AZURE_RESOURCE_GROUP.
// See example_config.go for information about the environment variables.
func ExampleNetworkInterfacesOperations_BeginDelete() {
	client := getNetworkInterfacesClient()
	nic, err := client.BeginDelete(context.Background(), resourceGroupName, nicName)
	if err != nil {
		panic(err)
	}
	res, err := nic.PollUntilDone(context.Background(), 5*time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.StatusCode)
	// Output:
	// 200
}

func getIP() *PublicIPAddress {
	client := getPublicIPClient()
	resp, err := client.Get(context.Background(), resourceGroupName, ip, nil)
	if err != nil {
		panic(err)
	}
	return resp.PublicIPAddress
}

func getSubnet() *Subnet {
	client := getSubnetClient()
	resp, err := client.Get(context.Background(), resourceGroupName, vnetName, subnetName, nil)
	if err != nil {
		panic(err)
	}
	return resp.Subnet
}
