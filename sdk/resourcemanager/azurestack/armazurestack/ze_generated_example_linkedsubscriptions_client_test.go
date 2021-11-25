//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armazurestack_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/azurestack/armazurestack"
)

// x-ms-original-file: specification/azurestack/resource-manager/Microsoft.AzureStack/preview/2020-06-01-preview/examples/LinkedSubscription/List.json
func ExampleLinkedSubscriptionsClient_ListByResourceGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armazurestack.NewLinkedSubscriptionsClient("<subscription-id>", cred, nil)
	pager := client.ListByResourceGroup("<resource-group>",
		nil)
	for pager.NextPage(ctx) {
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("LinkedSubscription.ID: %s\n", *v.ID)
		}
	}
}

// x-ms-original-file: specification/azurestack/resource-manager/Microsoft.AzureStack/preview/2020-06-01-preview/examples/LinkedSubscription/ListBySubscription.json
func ExampleLinkedSubscriptionsClient_ListBySubscription() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armazurestack.NewLinkedSubscriptionsClient("<subscription-id>", cred, nil)
	pager := client.ListBySubscription(nil)
	for pager.NextPage(ctx) {
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("LinkedSubscription.ID: %s\n", *v.ID)
		}
	}
}

// x-ms-original-file: specification/azurestack/resource-manager/Microsoft.AzureStack/preview/2020-06-01-preview/examples/LinkedSubscription/Get.json
func ExampleLinkedSubscriptionsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armazurestack.NewLinkedSubscriptionsClient("<subscription-id>", cred, nil)
	res, err := client.Get(ctx,
		"<resource-group>",
		"<linked-subscription-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("LinkedSubscription.ID: %s\n", *res.ID)
}

// x-ms-original-file: specification/azurestack/resource-manager/Microsoft.AzureStack/preview/2020-06-01-preview/examples/LinkedSubscription/Delete.json
func ExampleLinkedSubscriptionsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armazurestack.NewLinkedSubscriptionsClient("<subscription-id>", cred, nil)
	_, err = client.Delete(ctx,
		"<resource-group>",
		"<linked-subscription-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/azurestack/resource-manager/Microsoft.AzureStack/preview/2020-06-01-preview/examples/LinkedSubscription/Put.json
func ExampleLinkedSubscriptionsClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armazurestack.NewLinkedSubscriptionsClient("<subscription-id>", cred, nil)
	res, err := client.CreateOrUpdate(ctx,
		"<resource-group>",
		"<linked-subscription-name>",
		armazurestack.LinkedSubscriptionParameter{
			Location: armazurestack.LocationGlobal.ToPtr(),
			Properties: &armazurestack.LinkedSubscriptionParameterProperties{
				LinkedSubscriptionID:   to.StringPtr("<linked-subscription-id>"),
				RegistrationResourceID: to.StringPtr("<registration-resource-id>"),
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("LinkedSubscription.ID: %s\n", *res.ID)
}

// x-ms-original-file: specification/azurestack/resource-manager/Microsoft.AzureStack/preview/2020-06-01-preview/examples/LinkedSubscription/Patch.json
func ExampleLinkedSubscriptionsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armazurestack.NewLinkedSubscriptionsClient("<subscription-id>", cred, nil)
	res, err := client.Update(ctx,
		"<resource-group>",
		"<linked-subscription-name>",
		armazurestack.LinkedSubscriptionParameter{
			Location: armazurestack.LocationGlobal.ToPtr(),
			Properties: &armazurestack.LinkedSubscriptionParameterProperties{
				LinkedSubscriptionID:   to.StringPtr("<linked-subscription-id>"),
				RegistrationResourceID: to.StringPtr("<registration-resource-id>"),
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("LinkedSubscription.ID: %s\n", *res.ID)
}
