//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armlabservices_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/labservices/armlabservices"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/labservices/resource-manager/Microsoft.LabServices/preview/2021-11-15-preview/examples/Labs/listLabs.json
func ExampleLabsClient_NewListBySubscriptionPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armlabservices.NewLabsClient("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListBySubscriptionPager(&armlabservices.LabsClientListBySubscriptionOptions{Filter: nil})
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/labservices/resource-manager/Microsoft.LabServices/preview/2021-11-15-preview/examples/Labs/listResourceGroupLabs.json
func ExampleLabsClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armlabservices.NewLabsClient("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListByResourceGroupPager("testrg123",
		nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/labservices/resource-manager/Microsoft.LabServices/preview/2021-11-15-preview/examples/Labs/getLab.json
func ExampleLabsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armlabservices.NewLabsClient("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"testrg123",
		"testlab",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/labservices/resource-manager/Microsoft.LabServices/preview/2021-11-15-preview/examples/Labs/putLab.json
func ExampleLabsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armlabservices.NewLabsClient("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginCreateOrUpdate(ctx,
		"testrg123",
		"testlab",
		armlabservices.Lab{
			Location: to.Ptr("westus"),
			Properties: &armlabservices.LabProperties{
				Description: to.Ptr("This is a test lab."),
				AutoShutdownProfile: &armlabservices.AutoShutdownProfile{
					DisconnectDelay:          to.Ptr("00:05"),
					IdleDelay:                to.Ptr("01:00"),
					NoConnectDelay:           to.Ptr("01:00"),
					ShutdownOnDisconnect:     to.Ptr(armlabservices.EnableStateEnabled),
					ShutdownOnIdle:           to.Ptr(armlabservices.ShutdownOnIdleModeUserAbsence),
					ShutdownWhenNotConnected: to.Ptr(armlabservices.EnableStateEnabled),
				},
				ConnectionProfile: &armlabservices.ConnectionProfile{
					ClientRdpAccess: to.Ptr(armlabservices.ConnectionTypePublic),
					ClientSSHAccess: to.Ptr(armlabservices.ConnectionTypePublic),
					WebRdpAccess:    to.Ptr(armlabservices.ConnectionTypeNone),
					WebSSHAccess:    to.Ptr(armlabservices.ConnectionTypeNone),
				},
				LabPlanID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.LabServices/labPlans/testlabplan"),
				SecurityProfile: &armlabservices.SecurityProfile{
					OpenAccess: to.Ptr(armlabservices.EnableStateDisabled),
				},
				Title: to.Ptr("Test Lab"),
				VirtualMachineProfile: &armlabservices.VirtualMachineProfile{
					AdditionalCapabilities: &armlabservices.VirtualMachineAdditionalCapabilities{
						InstallGpuDrivers: to.Ptr(armlabservices.EnableStateDisabled),
					},
					AdminUser: &armlabservices.Credentials{
						Username: to.Ptr("test-user"),
					},
					CreateOption: to.Ptr(armlabservices.CreateOptionTemplateVM),
					ImageReference: &armlabservices.ImageReference{
						Offer:     to.Ptr("WindowsServer"),
						Publisher: to.Ptr("Microsoft"),
						SKU:       to.Ptr("2019-Datacenter"),
						Version:   to.Ptr("2019.0.20190410"),
					},
					SKU: &armlabservices.SKU{
						Name: to.Ptr("Medium"),
					},
					UsageQuota:        to.Ptr("10:00"),
					UseSharedPassword: to.Ptr(armlabservices.EnableStateDisabled),
				},
				NetworkProfile: &armlabservices.LabNetworkProfile{
					SubnetID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.Network/virtualNetworks/test-vnet/subnets/default"),
				},
				State: to.Ptr(armlabservices.LabStateDraft),
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/labservices/resource-manager/Microsoft.LabServices/preview/2021-11-15-preview/examples/Labs/patchLab.json
func ExampleLabsClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armlabservices.NewLabsClient("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginUpdate(ctx,
		"testrg123",
		"testlab",
		armlabservices.LabUpdate{
			Properties: &armlabservices.LabUpdateProperties{
				SecurityProfile: &armlabservices.SecurityProfile{
					OpenAccess: to.Ptr(armlabservices.EnableStateEnabled),
				},
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/labservices/resource-manager/Microsoft.LabServices/preview/2021-11-15-preview/examples/Labs/deleteLab.json
func ExampleLabsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armlabservices.NewLabsClient("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginDelete(ctx,
		"testrg123",
		"testlab",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/labservices/resource-manager/Microsoft.LabServices/preview/2021-11-15-preview/examples/Labs/publishLab.json
func ExampleLabsClient_BeginPublish() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armlabservices.NewLabsClient("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginPublish(ctx,
		"testrg123",
		"testlab",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/labservices/resource-manager/Microsoft.LabServices/preview/2021-11-15-preview/examples/Labs/syncLab.json
func ExampleLabsClient_BeginSyncGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armlabservices.NewLabsClient("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginSyncGroup(ctx,
		"testrg123",
		"testlab",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
