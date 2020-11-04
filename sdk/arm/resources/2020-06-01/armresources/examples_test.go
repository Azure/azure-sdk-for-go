// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armresources_test

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/arm/resources/2020-06-01/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func ExampleResourceGroupsOperations_CheckExistence() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armresources.NewResourceGroupsClient(armresources.NewDefaultClient(cred, nil), "<subscription ID>")
	resp, err := client.CheckExistence(context.Background(), "<resource group name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("resource group exists: %v\n", resp.Success)
}

func ExampleResourceGroupsOperations_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armresources.NewResourceGroupsClient(armresources.NewDefaultClient(cred, nil), "<subscription ID>")
	resp, err := client.CreateOrUpdate(context.Background(), "<resource group name>", armresources.ResourceGroup{
		Location: to.StringPtr("<location>"),
	}, nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("resource group ID: %s\n", *resp.ResourceGroup.ID)
}

func ExampleResourceGroupsOperations_List() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armresources.NewResourceGroupsClient(armresources.NewDefaultClient(cred, nil), "<subscription ID>")
	pager := client.List(nil)
	for pager.NextPage(context.Background()) {
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, rg := range *pager.PageResponse().ResourceGroupListResult.Value {
			log.Printf("resource group ID: %s\n", *rg.ID)
		}
	}
}

func ExampleResourceGroupsOperations_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armresources.NewResourceGroupsClient(armresources.NewDefaultClient(cred, nil), "<subscription ID>")
	poller, err := client.BeginDelete(context.Background(), "<resource group name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to delete resource group: %v", err)
	}
}
