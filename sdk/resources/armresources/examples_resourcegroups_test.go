// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armresources_test

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func ExampleResourceGroupsClient_CheckExistence() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armresources.NewResourceGroupsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
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
	client := armresources.NewResourceGroupsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
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
	client := armresources.NewResourceGroupsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
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
	client := armresources.NewResourceGroupsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
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
	client := armresources.NewResourceGroupsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
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
	client := armresources.NewResourceGroupsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	rg, err := client.Get(context.Background(), "<resource group name>", nil)
	if err != nil {
		log.Fatalf("failed to get resource group: %v", err)
	}
	log.Printf("resource group name: %s\n", *rg.ResourceGroup.Name)
}

func ExampleResourceGroupsClient_ResumeDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armresources.NewResourceGroupsClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	pollerResp, err := client.BeginDelete(context.Background(), "<resource group name>", nil)
	if err != nil {
		log.Fatalf("failed to get resource group: %v", err)
	}
	tk, err := pollerResp.Poller.ResumeToken()
	if err != nil {
		log.Fatalf("failed to get resource group poller resume token: %v", err)
	}
	rgPoller, err := client.ResumeDelete(context.Background(), tk)
	if err != nil {
		log.Fatalf("failed to get resource group: %v", err)
	}
	for {
		_, err := rgPoller.Poller.Poll(context.Background())
		if err != nil {
			log.Fatalf("failed to poll for status: %v", err)
		}
		if rgPoller.Poller.Done() {
			break
		}
		time.Sleep(5 * time.Second)
	}
	_, err = rgPoller.Poller.FinalResponse(context.Background())
	if err != nil {
		log.Fatalf("failed to get final poller response: %v", err)
	}
}
