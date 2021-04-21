// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armresources_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/arm/resources/2020-06-01/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func ExampleResourcesClient_GetByID() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armresources.NewResourcesClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.GetByID(context.Background(), "<resource ID>", "<api version>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("resource ID: %v\n", resp.GenericResource.ID)
}

func ExampleResourcesClient_ListByResourceGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armresources.NewResourcesClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	page := client.ListByResourceGroup("<resource group name>", nil)
	for page.NextPage(context.Background()) {
		resp := page.PageResponse()
		if len(*resp.ResourceListResult.Value) == 0 {
			log.Fatal("missing payload")
		}
		for _, val := range *resp.ResourceListResult.Value {
			log.Printf("resource: %v", *val.ID)
		}
	}
	if err := page.Err(); err != nil {
		log.Fatal(err)
	}

}
