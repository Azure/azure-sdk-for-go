//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armappconfiguration_test

import (
	"context"
	"log"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appconfiguration/armappconfiguration"
)

// x-ms-original-file: specification/appconfiguration/resource-manager/Microsoft.AppConfiguration/preview/2021-03-01-preview/examples/ConfigurationStoresList.json
func ExampleConfigurationStoresClient_List() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armappconfiguration.NewConfigurationStoresClient("<subscription-id>", cred, nil)
	pager := client.List(&armappconfiguration.ConfigurationStoresClientListOptions{SkipToken: nil})
	for {
		nextResult := pager.NextPage(ctx)
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		if !nextResult {
			break
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("Pager result: %#v\n", v)
		}
	}
}

// x-ms-original-file: specification/appconfiguration/resource-manager/Microsoft.AppConfiguration/preview/2021-03-01-preview/examples/ConfigurationStoresListByResourceGroup.json
func ExampleConfigurationStoresClient_ListByResourceGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armappconfiguration.NewConfigurationStoresClient("<subscription-id>", cred, nil)
	pager := client.ListByResourceGroup("<resource-group-name>",
		&armappconfiguration.ConfigurationStoresClientListByResourceGroupOptions{SkipToken: nil})
	for {
		nextResult := pager.NextPage(ctx)
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		if !nextResult {
			break
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("Pager result: %#v\n", v)
		}
	}
}

// x-ms-original-file: specification/appconfiguration/resource-manager/Microsoft.AppConfiguration/preview/2021-03-01-preview/examples/ConfigurationStoresGet.json
func ExampleConfigurationStoresClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armappconfiguration.NewConfigurationStoresClient("<subscription-id>", cred, nil)
	res, err := client.Get(ctx,
		"<resource-group-name>",
		"<config-store-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.ConfigurationStoresClientGetResult)
}

// x-ms-original-file: specification/appconfiguration/resource-manager/Microsoft.AppConfiguration/preview/2021-03-01-preview/examples/ConfigurationStoresCreate.json
func ExampleConfigurationStoresClient_BeginCreate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armappconfiguration.NewConfigurationStoresClient("<subscription-id>", cred, nil)
	poller, err := client.BeginCreate(ctx,
		"<resource-group-name>",
		"<config-store-name>",
		armappconfiguration.ConfigurationStore{
			Location: to.StringPtr("<location>"),
			Tags: map[string]*string{
				"myTag": to.StringPtr("myTagValue"),
			},
			SKU: &armappconfiguration.SKU{
				Name: to.StringPtr("<name>"),
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.ConfigurationStoresClientCreateResult)
}

// x-ms-original-file: specification/appconfiguration/resource-manager/Microsoft.AppConfiguration/preview/2021-03-01-preview/examples/ConfigurationStoresDelete.json
func ExampleConfigurationStoresClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armappconfiguration.NewConfigurationStoresClient("<subscription-id>", cred, nil)
	poller, err := client.BeginDelete(ctx,
		"<resource-group-name>",
		"<config-store-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/appconfiguration/resource-manager/Microsoft.AppConfiguration/preview/2021-03-01-preview/examples/ConfigurationStoresUpdate.json
func ExampleConfigurationStoresClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armappconfiguration.NewConfigurationStoresClient("<subscription-id>", cred, nil)
	poller, err := client.BeginUpdate(ctx,
		"<resource-group-name>",
		"<config-store-name>",
		armappconfiguration.ConfigurationStoreUpdateParameters{
			SKU: &armappconfiguration.SKU{
				Name: to.StringPtr("<name>"),
			},
			Tags: map[string]*string{
				"Category": to.StringPtr("Marketing"),
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.ConfigurationStoresClientUpdateResult)
}

// x-ms-original-file: specification/appconfiguration/resource-manager/Microsoft.AppConfiguration/preview/2021-03-01-preview/examples/ConfigurationStoresListKeys.json
func ExampleConfigurationStoresClient_ListKeys() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armappconfiguration.NewConfigurationStoresClient("<subscription-id>", cred, nil)
	pager := client.ListKeys("<resource-group-name>",
		"<config-store-name>",
		&armappconfiguration.ConfigurationStoresClientListKeysOptions{SkipToken: nil})
	for {
		nextResult := pager.NextPage(ctx)
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		if !nextResult {
			break
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("Pager result: %#v\n", v)
		}
	}
}

// x-ms-original-file: specification/appconfiguration/resource-manager/Microsoft.AppConfiguration/preview/2021-03-01-preview/examples/ConfigurationStoresRegenerateKey.json
func ExampleConfigurationStoresClient_RegenerateKey() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armappconfiguration.NewConfigurationStoresClient("<subscription-id>", cred, nil)
	res, err := client.RegenerateKey(ctx,
		"<resource-group-name>",
		"<config-store-name>",
		armappconfiguration.RegenerateKeyParameters{
			ID: to.StringPtr("<id>"),
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.ConfigurationStoresClientRegenerateKeyResult)
}
