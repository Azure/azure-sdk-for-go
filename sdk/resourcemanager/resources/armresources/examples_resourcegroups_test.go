//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armresources_test

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

func ExampleResourceGroupsClient_CheckExistence() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armresources.NewResourceGroupsClient("<subscription ID>", cred, nil)
	resp, err := client.CheckExistence(context.Background(), "<resource group name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("resource group exists: %v\n", resp.Success)
}

func ExampleResourceGroupsClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armresources.NewResourceGroupsClient("<subscription ID>", cred, nil)
	resp, err := client.CreateOrUpdate(context.Background(), "<resource group name>", armresources.ResourceGroup{
		Location: to.StringPtr("<location>"),
	}, nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("resource group ID: %s\n", *resp.ResourceGroup.ID)
}

func ExampleResourceGroupsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armresources.NewResourceGroupsClient("<subscription ID>", cred, nil)
	resp, err := client.Update(
		context.Background(),
		"<resource group name>",
		armresources.ResourceGroupPatchable{
			Tags: map[string]*string{
				"exampleTag": to.StringPtr("exampleTagValue"),
			}},
		nil)
	if err != nil {
		log.Fatalf("failed to update resource group: %v", err)
	}
	log.Printf("updated resource group: %v", *resp.ResourceGroup.ID)
}

func ExampleResourceGroupsClient_List() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armresources.NewResourceGroupsClient("<subscription ID>", cred, nil)
	pager := client.List(nil)
	for pager.NextPage(context.Background()) {
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, rg := range pager.PageResponse().ResourceGroupListResult.Value {
			log.Printf("resource group ID: %s\n", *rg.ID)
		}
	}
}

func ExampleResourceGroupsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armresources.NewResourceGroupsClient("<subscription ID>", cred, nil)
	poller, err := client.BeginDelete(context.Background(), "<resource group name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	_, err = poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to delete resource group: %v", err)
	}
}

func ExampleResourceGroupsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armresources.NewResourceGroupsClient("<subscription ID>", cred, nil)
	rg, err := client.Get(context.Background(), "<resource group name>", nil)
	if err != nil {
		log.Fatalf("failed to get resource group: %v", err)
	}
	log.Printf("resource group name: %s\n", *rg.ResourceGroup.Name)
}
