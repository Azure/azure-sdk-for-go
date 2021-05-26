// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstorage_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/armstorage"
)

func ExampleBlobContainersClient_Create() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewBlobContainersClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.Create(
		context.Background(),
		"<resource group name>",
		"<storage account name>",
		"<container name>",
		armstorage.BlobContainer{}, nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("blob container ID: %v\n", *resp.BlobContainer.ID)
}

func ExampleBlobContainersClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewBlobContainersClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.Get(
		context.Background(),
		"<resource group name>",
		"<storage account name>",
		"<container name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("blob container ID: %v\n", *resp.BlobContainer.ID)
}

func ExampleBlobContainersClient_List() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewBlobContainersClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	pager := client.List("<resource group name>", "<storage account name>", nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.ListContainerItems.Value) == 0 {
			log.Fatal("missing payload")
		}
		for _, val := range resp.ListContainerItems.Value {
			log.Printf("container item: %v", *val.ID)
		}
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}
}

func ExampleBlobContainersClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewBlobContainersClient(armcore.NewDefaultConnection(cred, nil), "<subscription ID>")
	_, err = client.Delete(
		context.Background(),
		"<resource group name>",
		"<storage account name>",
		"<container name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
}
