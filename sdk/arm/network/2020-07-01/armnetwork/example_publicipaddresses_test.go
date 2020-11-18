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

func ExamplePublicIPAddressesOperations_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewPublicIPAddressesClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		"<IP name>",
		armnetwork.PublicIPAddress{
			Resource: armnetwork.Resource{
				Name:     to.StringPtr("<IP name>"),
				Location: to.StringPtr("<Azure location>"),
			},
			Properties: &armnetwork.PublicIPAddressPropertiesFormat{
				PublicIPAddressVersion:   armnetwork.IPVersionIPv4.ToPtr(),
				PublicIPAllocationMethod: armnetwork.IPAllocationMethodStatic.ToPtr(),
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
	log.Printf("public IP address ID: %v", *resp.PublicIPAddress.ID)
}

func ExamplePublicIPAddressesOperations_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewPublicIPAddressesClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.Get(context.Background(), "<resource group name>", "<IP name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("public IP address ID: %v", *resp.PublicIPAddress.ID)
}

func ExamplePublicIPAddressesOperations_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewPublicIPAddressesClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.BeginDelete(context.Background(), "<resource group name>", "<IP name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = resp.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to delete resource: %v", err)
	}
}
