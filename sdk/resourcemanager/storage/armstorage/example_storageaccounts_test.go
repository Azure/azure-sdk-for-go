//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstorage_test

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

func ExampleStorageAccountsClient_BeginCreate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewStorageAccountsClient("<subscription ID>", cred, nil)
	poller, err := client.BeginCreate(
		context.Background(),
		"<resource group name>",
		"<storage account name>",
		armstorage.StorageAccountCreateParameters{
			SKU: &armstorage.SKU{
				Name: armstorage.SKUNameStandardLRS.ToPtr(),
				Tier: armstorage.SKUTierStandard.ToPtr(),
			},
			Kind:     armstorage.KindBlobStorage.ToPtr(),
			Location: to.StringPtr("<Azure location>"),
			Properties: &armstorage.StorageAccountPropertiesCreateParameters{
				AccessTier: armstorage.AccessTierCool.ToPtr(),
			},
		}, nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	resp, err := poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to create storage account: %v", err)
	}
	log.Printf("storage account ID: %v\n", *resp.StorageAccount.ID)
}

func ExampleStorageAccountsClient_List() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewStorageAccountsClient("<subscription ID>", cred, nil)
	pager := client.List(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.StorageAccountListResult.Value) == 0 {
			log.Fatal("missing payload")
		}
		for _, val := range resp.StorageAccountListResult.Value {
			log.Printf("storage account: %v", *val.ID)
		}
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}
}

func ExampleStorageAccountsClient_ListByResourceGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewStorageAccountsClient("<subscription ID>", cred, nil)
	pager := client.ListByResourceGroup("<resource group name>", nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		if len(resp.StorageAccountListResult.Value) == 0 {
			log.Fatal("missing payload")
		}
		for _, val := range resp.StorageAccountListResult.Value {
			log.Printf("storage account: %v", *val.ID)
		}
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}
}

func ExampleStorageAccountsClient_CheckNameAvailability() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewStorageAccountsClient("<subscription ID>", cred, nil)
	resp, err := client.CheckNameAvailability(
		context.Background(),
		armstorage.StorageAccountCheckNameAvailabilityParameters{
			Name: to.StringPtr("<storage account name>"),
			Type: to.StringPtr("Microsoft.Storage/storageAccounts"),
		},
		nil)
	if err != nil {
		log.Fatalf("failed to delete account: %v", err)
	}
	log.Printf("name availability: %v", *resp.CheckNameAvailabilityResult.NameAvailable)
}

func ExampleStorageAccountsClient_ListKeys() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewStorageAccountsClient("<subscription ID>", cred, nil)
	resp, err := client.ListKeys(context.Background(), "<resource group name>", "<storage account name>", nil)
	if err != nil {
		log.Fatalf("failed to delete account: %v", err)
	}
	for _, k := range resp.StorageAccountListKeysResult.Keys {
		log.Printf("account key: %v", *k.KeyName)
	}
}

func ExampleStorageAccountsClient_GetProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewStorageAccountsClient("<subscription ID>", cred, nil)
	resp, err := client.GetProperties(context.Background(), "<resource group name>", "<storage account name>", nil)
	if err != nil {
		log.Fatalf("failed to delete account: %v", err)
	}
	log.Printf("storage account properties Access Tier: %v", *resp.StorageAccount.Properties.AccessTier)
}

func ExampleStorageAccountsClient_RegenerateKey() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewStorageAccountsClient("<subscription ID>", cred, nil)
	resp, err := client.RegenerateKey(context.Background(), "<resource group name>", "<storage account name>", armstorage.StorageAccountRegenerateKeyParameters{KeyName: to.StringPtr("<key name>")}, nil)
	if err != nil {
		log.Fatalf("failed to delete account: %v", err)
	}
	for _, k := range resp.StorageAccountListKeysResult.Keys {
		log.Printf("key: %v, value: %v", *k.KeyName, *k.Value)
	}
}

func ExampleStorageAccountsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewStorageAccountsClient("<subscription ID>", cred, nil)
	resp, err := client.Update(
		context.Background(),
		"<resource group name>",
		"<storage account name>",
		armstorage.StorageAccountUpdateParameters{
			Tags: map[string]*string{
				"who rocks": to.StringPtr("golang"),
				"where":     to.StringPtr("on azure")}}, nil)
	if err != nil {
		log.Fatalf("failed to delete account: %v", err)
	}
	log.Printf("storage account ID: %v", *resp.StorageAccount.ID)
}
func ExampleStorageAccountsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armstorage.NewStorageAccountsClient("<subscription ID>", cred, nil)
	_, err = client.Delete(context.Background(), "<resource group name>", "<storage account name>", nil)
	if err != nil {
		log.Fatalf("failed to delete account: %v", err)
	}
}
