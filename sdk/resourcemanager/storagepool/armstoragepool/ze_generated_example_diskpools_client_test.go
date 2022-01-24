//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armstoragepool_test

import (
	"context"
	"log"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagepool/armstoragepool"
)

// x-ms-original-file: specification/storagepool/resource-manager/Microsoft.StoragePool/stable/2021-08-01/examples/DiskPools_ListBySubscription.json
func ExampleDiskPoolsClient_ListBySubscription() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armstoragepool.NewDiskPoolsClient("<subscription-id>", cred, nil)
	pager := client.ListBySubscription(nil)
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

// x-ms-original-file: specification/storagepool/resource-manager/Microsoft.StoragePool/stable/2021-08-01/examples/DiskPools_ListByResourceGroup.json
func ExampleDiskPoolsClient_ListByResourceGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armstoragepool.NewDiskPoolsClient("<subscription-id>", cred, nil)
	pager := client.ListByResourceGroup("<resource-group-name>",
		nil)
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

// x-ms-original-file: specification/storagepool/resource-manager/Microsoft.StoragePool/stable/2021-08-01/examples/DiskPools_Put.json
func ExampleDiskPoolsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armstoragepool.NewDiskPoolsClient("<subscription-id>", cred, nil)
	poller, err := client.BeginCreateOrUpdate(ctx,
		"<resource-group-name>",
		"<disk-pool-name>",
		armstoragepool.DiskPoolCreate{
			Location: to.StringPtr("<location>"),
			Properties: &armstoragepool.DiskPoolCreateProperties{
				AvailabilityZones: []*string{
					to.StringPtr("1")},
				Disks: []*armstoragepool.Disk{
					{
						ID: to.StringPtr("<id>"),
					},
					{
						ID: to.StringPtr("<id>"),
					}},
				SubnetID: to.StringPtr("<subnet-id>"),
			},
			SKU: &armstoragepool.SKU{
				Name: to.StringPtr("<name>"),
				Tier: to.StringPtr("<tier>"),
			},
			Tags: map[string]*string{
				"key": to.StringPtr("value"),
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
	log.Printf("Response result: %#v\n", res.DiskPoolsClientCreateOrUpdateResult)
}

// x-ms-original-file: specification/storagepool/resource-manager/Microsoft.StoragePool/stable/2021-08-01/examples/DiskPools_Patch.json
func ExampleDiskPoolsClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armstoragepool.NewDiskPoolsClient("<subscription-id>", cred, nil)
	poller, err := client.BeginUpdate(ctx,
		"<resource-group-name>",
		"<disk-pool-name>",
		armstoragepool.DiskPoolUpdate{
			Properties: &armstoragepool.DiskPoolUpdateProperties{
				Disks: []*armstoragepool.Disk{
					{
						ID: to.StringPtr("<id>"),
					},
					{
						ID: to.StringPtr("<id>"),
					}},
			},
			SKU: &armstoragepool.SKU{
				Name: to.StringPtr("<name>"),
				Tier: to.StringPtr("<tier>"),
			},
			Tags: map[string]*string{
				"key": to.StringPtr("value"),
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
	log.Printf("Response result: %#v\n", res.DiskPoolsClientUpdateResult)
}

// x-ms-original-file: specification/storagepool/resource-manager/Microsoft.StoragePool/stable/2021-08-01/examples/DiskPools_Delete.json
func ExampleDiskPoolsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armstoragepool.NewDiskPoolsClient("<subscription-id>", cred, nil)
	poller, err := client.BeginDelete(ctx,
		"<resource-group-name>",
		"<disk-pool-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/storagepool/resource-manager/Microsoft.StoragePool/stable/2021-08-01/examples/DiskPools_Get.json
func ExampleDiskPoolsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armstoragepool.NewDiskPoolsClient("<subscription-id>", cred, nil)
	res, err := client.Get(ctx,
		"<resource-group-name>",
		"<disk-pool-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.DiskPoolsClientGetResult)
}

// x-ms-original-file: specification/storagepool/resource-manager/Microsoft.StoragePool/stable/2021-08-01/examples/DiskPools_GetOutboundNetworkDependencies.json
func ExampleDiskPoolsClient_ListOutboundNetworkDependenciesEndpoints() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armstoragepool.NewDiskPoolsClient("<subscription-id>", cred, nil)
	pager := client.ListOutboundNetworkDependenciesEndpoints("<resource-group-name>",
		"<disk-pool-name>",
		nil)
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

// x-ms-original-file: specification/storagepool/resource-manager/Microsoft.StoragePool/stable/2021-08-01/examples/DiskPools_Start.json
func ExampleDiskPoolsClient_BeginStart() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armstoragepool.NewDiskPoolsClient("<subscription-id>", cred, nil)
	poller, err := client.BeginStart(ctx,
		"<resource-group-name>",
		"<disk-pool-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/storagepool/resource-manager/Microsoft.StoragePool/stable/2021-08-01/examples/DiskPools_Deallocate.json
func ExampleDiskPoolsClient_BeginDeallocate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armstoragepool.NewDiskPoolsClient("<subscription-id>", cred, nil)
	poller, err := client.BeginDeallocate(ctx,
		"<resource-group-name>",
		"<disk-pool-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/storagepool/resource-manager/Microsoft.StoragePool/stable/2021-08-01/examples/DiskPools_Upgrade.json
func ExampleDiskPoolsClient_BeginUpgrade() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armstoragepool.NewDiskPoolsClient("<subscription-id>", cred, nil)
	poller, err := client.BeginUpgrade(ctx,
		"<resource-group-name>",
		"<disk-pool-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}
