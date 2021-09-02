//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute"
)

func ExampleDisksClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewDisksClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	poller, err := client.BeginCreateOrUpdate(
		context.Background(),
		"<resource group name>",
		"<disk name>",
		armcompute.Disk{
			Resource: armcompute.Resource{
				Name:     to.StringPtr("<disk name>"),
				Location: to.StringPtr("<Azure location>"),
			},
			Properties: &armcompute.DiskProperties{
				CreationData: &armcompute.CreationData{
					CreateOption: armcompute.DiskCreateOptionEmpty.ToPtr(),
				},
				DiskSizeGB: to.Int32Ptr(64),
			},
		},
		nil,
	)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	resp, err := poller.PollUntilDone(context.Background(), 30*time.Second)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("disk ID: %v", *resp.Disk.ID)
}

func ExampleDisksClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewDisksClient(arm.NewDefaultConnection(cred, nil), "<subscription ID>")
	resp, err := client.Get(context.Background(), "<resource group name>", "<disk name>", nil)
	if err != nil {
		log.Fatalf("failed to obtain a response: %v", err)
	}
	log.Printf("disk ID: %s", *resp.Disk.ID)
}
