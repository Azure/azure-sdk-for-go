//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armstoragesync_test

import (
	"context"
	"log"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagesync/armstoragesync"
)

// x-ms-original-file: specification/storagesync/resource-manager/Microsoft.StorageSync/stable/2020-09-01/examples/RegisteredServers_ListByStorageSyncService.json
func ExampleRegisteredServersClient_ListByStorageSyncService() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armstoragesync.NewRegisteredServersClient("<subscription-id>", cred, nil)
	res, err := client.ListByStorageSyncService(ctx,
		"<resource-group-name>",
		"<storage-sync-service-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RegisteredServersClientListByStorageSyncServiceResult)
}

// x-ms-original-file: specification/storagesync/resource-manager/Microsoft.StorageSync/stable/2020-09-01/examples/RegisteredServers_Get.json
func ExampleRegisteredServersClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armstoragesync.NewRegisteredServersClient("<subscription-id>", cred, nil)
	res, err := client.Get(ctx,
		"<resource-group-name>",
		"<storage-sync-service-name>",
		"<server-id>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.RegisteredServersClientGetResult)
}

// x-ms-original-file: specification/storagesync/resource-manager/Microsoft.StorageSync/stable/2020-09-01/examples/RegisteredServers_Create.json
func ExampleRegisteredServersClient_BeginCreate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armstoragesync.NewRegisteredServersClient("<subscription-id>", cred, nil)
	poller, err := client.BeginCreate(ctx,
		"<resource-group-name>",
		"<storage-sync-service-name>",
		"<server-id>",
		armstoragesync.RegisteredServerCreateParameters{
			Properties: &armstoragesync.RegisteredServerCreateParametersProperties{
				AgentVersion:      to.StringPtr("<agent-version>"),
				FriendlyName:      to.StringPtr("<friendly-name>"),
				ServerCertificate: to.StringPtr("<server-certificate>"),
				ServerID:          to.StringPtr("<server-id>"),
				ServerOSVersion:   to.StringPtr("<server-osversion>"),
				ServerRole:        to.StringPtr("<server-role>"),
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
	log.Printf("Response result: %#v\n", res.RegisteredServersClientCreateResult)
}

// x-ms-original-file: specification/storagesync/resource-manager/Microsoft.StorageSync/stable/2020-09-01/examples/RegisteredServers_Delete.json
func ExampleRegisteredServersClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armstoragesync.NewRegisteredServersClient("<subscription-id>", cred, nil)
	poller, err := client.BeginDelete(ctx,
		"<resource-group-name>",
		"<storage-sync-service-name>",
		"<server-id>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/storagesync/resource-manager/Microsoft.StorageSync/stable/2020-09-01/examples/RegisteredServers_TriggerRollover.json
func ExampleRegisteredServersClient_BeginTriggerRollover() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armstoragesync.NewRegisteredServersClient("<subscription-id>", cred, nil)
	poller, err := client.BeginTriggerRollover(ctx,
		"<resource-group-name>",
		"<storage-sync-service-name>",
		"<server-id>",
		armstoragesync.TriggerRolloverRequest{
			ServerCertificate: to.StringPtr("<server-certificate>"),
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}
