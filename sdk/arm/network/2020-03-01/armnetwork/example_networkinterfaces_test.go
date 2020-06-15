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
					NetworkInterfaceIPConfiguration{
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
}

func ExampleNetworkInterfacesOperations_Get() {
	client := getNetworkInterfacesClient()
	nic, err := client.Get(context.Background(), resourceGroupName, nicName, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(*nic.NetworkInterface.Name)
}

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
}

func getIP() *PublicIPAddress {
	client := getPublicIPClient()
	resp, err := client.Get(context.Background(), resourceGroupName, ipName, nil)
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
