//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armnetwork_test

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

func ExampleNetworkInterfacesClient_BeginCreateOrUpdate() {
	// Replace the "subnet" variable with a call to armnetwork.SubnetsClient.Get to retreive the subnet
	// instance that will be assigned to the network interface.
	var subnet *armnetwork.Subnet
	// Replace the "ipAddress" variable with a call to armnetwork.PublicIPAddressesClient.Get to retreive the
	// public IP address instance that will be assigned to the network interface.
	var ipAddress *armnetwork.PublicIPAddress
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkInterfacesClient("<subscription ID>", cred, nil)
	poller, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		"<NIC name>",
		armnetwork.NetworkInterface{
			Resource: armnetwork.Resource{
				Name:     to.StringPtr("<NIC name>"),
				Location: to.StringPtr("<Azure location>"),
			},
			Properties: &armnetwork.NetworkInterfacePropertiesFormat{
				IPConfigurations: []*armnetwork.NetworkInterfaceIPConfiguration{
					{
						Name: to.StringPtr("<NIC name>"),
						Properties: &armnetwork.NetworkInterfaceIPConfigurationPropertiesFormat{
							Subnet:                    subnet,
							PrivateIPAllocationMethod: armnetwork.IPAllocationMethodDynamic.ToPtr(),
							PublicIPAddress:           ipAddress,
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
	log.Printf("NIC ID: %v", *resp.NetworkInterface.ID)
}

func ExampleNetworkInterfacesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkInterfacesClient("<subscription ID>", cred, nil)
	resp, err := client.Get(context.Background(), "<resource group name>", "<NIC name>", nil)
	if err != nil {
		log.Fatalf("failed to get resource: %v", err)
	}
	log.Printf("NIC ID: %v", *resp.NetworkInterface.ID)
}

func ExampleNetworkInterfacesClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkInterfacesClient("<subscription ID>", cred, nil)
	poller, err := client.BeginDelete(context.Background(), "<resource group name>", "<NIC name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to delete NIC: %v", err)
	}
}
