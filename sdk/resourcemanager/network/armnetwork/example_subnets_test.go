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

func ExampleSubnetsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewSubnetsClient("<subscription ID>", cred, nil)
	poller, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		"<virtual network name>",
		"<subnet name>",
		armnetwork.Subnet{
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: to.StringPtr("10.0.0.0/16"), // NOTE: the allowed address range can change based on the network the subnet is created in
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
	log.Printf("subnet ID: %v", *resp.Subnet.ID)
}

func ExampleSubnetsClient_BeginCreateOrUpdate_withNetworkSecurityGroup() {
	// Replace the "nsg" variable with a call to armnetwork.NetworkSecurityGroups.Get to retreive the
	//network security group instance that will be assigned to the subnet.
	var nsg *armnetwork.NetworkSecurityGroup
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewSubnetsClient("<subscription ID>", cred, nil)
	poller, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		"<virtual network name>",
		"<subnet name>",
		armnetwork.Subnet{
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix:        to.StringPtr("10.0.0.0/16"), // NOTE: the allowed address range can change based on the network the subnet is created in
				NetworkSecurityGroup: nsg,
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
	log.Printf("subnet ID: %v", *resp.Subnet.ID)
}

func ExampleSubnetsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewSubnetsClient("<subscription ID>", cred, nil)
	resp, err := client.Get(context.Background(), "<resource group name>", "<virtual network name>", "<subnet name>", nil)
	if err != nil {
		log.Fatalf("failed to get resource: %v", err)
	}
	log.Printf("subnet ID: %v", *resp.Subnet.ID)
}

func ExampleSubnetsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewSubnetsClient("<subscription ID>", cred, nil)
	resp, err := client.BeginDelete(context.Background(), "<resource group name>", "<virtual network name>", "<subnet name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = resp.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to delete resource: %v", err)
	}
}
