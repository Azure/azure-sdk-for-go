// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armnetwork_test

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func ExampleSubnetsOperations_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewSubnetsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		"<virtual network name>",
		"<subnet name>",
		armnetwork.Subnet{
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix: to.StringPtr("10.0.0.0/16"),
			},
		},
		nil,
	)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	resp, err := poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("subnet ID: %v", *resp.Subnet.ID)
}

func ExampleSubnetsOperations_BeginCreateOrUpdate_withNetworkSecurityGroup() {
	// Replace the "nsg" variable with a call to armnetwork.NetworkSecurityGroups.Get to retreive the
	//network security group instance that will be assigned to the subnet.
	var nsg *armnetwork.NetworkSecurityGroup
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewSubnetsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		"<virtual network name>",
		"<subnet name>",
		armnetwork.Subnet{
			Properties: &armnetwork.SubnetPropertiesFormat{
				AddressPrefix:        to.StringPtr("10.0.0.0/16"),
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
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("subnet ID: %v", *resp.Subnet.ID)
}

func ExampleSubnetsOperations_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewSubnetsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.Get(context.Background(), "<resource group name>", "<virtual network name>", "<subnet name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("subnet ID: %v", *resp.Subnet.ID)
}

func ExampleSubnetsOperations_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewSubnetsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.BeginDelete(context.Background(), "<resource group name>", "<virtual network name>", "<subnet name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = resp.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to delete resource: %v", err)
	}
}
