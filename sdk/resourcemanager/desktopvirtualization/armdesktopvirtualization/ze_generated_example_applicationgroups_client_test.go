//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdesktopvirtualization_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/desktopvirtualization/armdesktopvirtualization"
)

// x-ms-original-file: specification/desktopvirtualization/resource-manager/Microsoft.DesktopVirtualization/preview/2021-09-03-preview/examples/ApplicationGroup_Get.json
func ExampleApplicationGroupsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armdesktopvirtualization.NewApplicationGroupsClient("<subscription-id>", cred, nil)
	res, err := client.Get(ctx,
		"<resource-group-name>",
		"<application-group-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ApplicationGroup.ID: %s\n", *res.ID)
}

// x-ms-original-file: specification/desktopvirtualization/resource-manager/Microsoft.DesktopVirtualization/preview/2021-09-03-preview/examples/ApplicationGroup_Create.json
func ExampleApplicationGroupsClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armdesktopvirtualization.NewApplicationGroupsClient("<subscription-id>", cred, nil)
	res, err := client.CreateOrUpdate(ctx,
		"<resource-group-name>",
		"<application-group-name>",
		armdesktopvirtualization.ApplicationGroup{
			ResourceModelWithAllowedPropertySet: armdesktopvirtualization.ResourceModelWithAllowedPropertySet{
				Location: to.StringPtr("<location>"),
				Tags: map[string]*string{
					"tag1": to.StringPtr("value1"),
					"tag2": to.StringPtr("value2"),
				},
			},
			Properties: &armdesktopvirtualization.ApplicationGroupProperties{
				Description:          to.StringPtr("<description>"),
				ApplicationGroupType: armdesktopvirtualization.ApplicationGroupTypeRemoteApp.ToPtr(),
				FriendlyName:         to.StringPtr("<friendly-name>"),
				HostPoolArmPath:      to.StringPtr("<host-pool-arm-path>"),
				MigrationRequest: &armdesktopvirtualization.MigrationRequestProperties{
					MigrationPath: to.StringPtr("<migration-path>"),
					Operation:     armdesktopvirtualization.OperationStart.ToPtr(),
				},
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ApplicationGroup.ID: %s\n", *res.ID)
}

// x-ms-original-file: specification/desktopvirtualization/resource-manager/Microsoft.DesktopVirtualization/preview/2021-09-03-preview/examples/ApplicationGroup_Delete.json
func ExampleApplicationGroupsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armdesktopvirtualization.NewApplicationGroupsClient("<subscription-id>", cred, nil)
	_, err = client.Delete(ctx,
		"<resource-group-name>",
		"<application-group-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/desktopvirtualization/resource-manager/Microsoft.DesktopVirtualization/preview/2021-09-03-preview/examples/ApplicationGroup_Update.json
func ExampleApplicationGroupsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armdesktopvirtualization.NewApplicationGroupsClient("<subscription-id>", cred, nil)
	res, err := client.Update(ctx,
		"<resource-group-name>",
		"<application-group-name>",
		&armdesktopvirtualization.ApplicationGroupsUpdateOptions{ApplicationGroup: &armdesktopvirtualization.ApplicationGroupPatch{
			Properties: &armdesktopvirtualization.ApplicationGroupPatchProperties{
				Description:  to.StringPtr("<description>"),
				FriendlyName: to.StringPtr("<friendly-name>"),
			},
			Tags: map[string]*string{
				"tag1": to.StringPtr("value1"),
				"tag2": to.StringPtr("value2"),
			},
		},
		})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ApplicationGroup.ID: %s\n", *res.ID)
}

// x-ms-original-file: specification/desktopvirtualization/resource-manager/Microsoft.DesktopVirtualization/preview/2021-09-03-preview/examples/ApplicationGroup_ListByResourceGroup.json
func ExampleApplicationGroupsClient_ListByResourceGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armdesktopvirtualization.NewApplicationGroupsClient("<subscription-id>", cred, nil)
	pager := client.ListByResourceGroup("<resource-group-name>",
		&armdesktopvirtualization.ApplicationGroupsListByResourceGroupOptions{Filter: to.StringPtr("<filter>")})
	for pager.NextPage(ctx) {
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("ApplicationGroup.ID: %s\n", *v.ID)
		}
	}
}

// x-ms-original-file: specification/desktopvirtualization/resource-manager/Microsoft.DesktopVirtualization/preview/2021-09-03-preview/examples/ApplicationGroup_ListBySubscription.json
func ExampleApplicationGroupsClient_ListBySubscription() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armdesktopvirtualization.NewApplicationGroupsClient("<subscription-id>", cred, nil)
	pager := client.ListBySubscription(&armdesktopvirtualization.ApplicationGroupsListBySubscriptionOptions{Filter: to.StringPtr("<filter>")})
	for pager.NextPage(ctx) {
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("ApplicationGroup.ID: %s\n", *v.ID)
		}
	}
}
