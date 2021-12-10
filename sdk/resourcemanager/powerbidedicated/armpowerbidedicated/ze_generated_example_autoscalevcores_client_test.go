//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armpowerbidedicated_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/powerbidedicated/armpowerbidedicated"
)

// x-ms-original-file: specification/powerbidedicated/resource-manager/Microsoft.PowerBIdedicated/stable/2021-01-01/examples/getAutoScaleVCore.json
func ExampleAutoScaleVCoresClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpowerbidedicated.NewAutoScaleVCoresClient("<subscription-id>", cred, nil)
	res, err := client.Get(ctx,
		"<resource-group-name>",
		"<vcore-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("AutoScaleVCore.ID: %s\n", *res.ID)
}

// x-ms-original-file: specification/powerbidedicated/resource-manager/Microsoft.PowerBIdedicated/stable/2021-01-01/examples/createAutoScaleVCore.json
func ExampleAutoScaleVCoresClient_Create() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpowerbidedicated.NewAutoScaleVCoresClient("<subscription-id>", cred, nil)
	res, err := client.Create(ctx,
		"<resource-group-name>",
		"<vcore-name>",
		armpowerbidedicated.AutoScaleVCore{
			Resource: armpowerbidedicated.Resource{
				Location: to.StringPtr("<location>"),
				Tags: map[string]*string{
					"testKey": to.StringPtr("testValue"),
				},
			},
			Properties: &armpowerbidedicated.AutoScaleVCoreProperties{
				AutoScaleVCoreMutableProperties: armpowerbidedicated.AutoScaleVCoreMutableProperties{
					CapacityLimit: to.Int32Ptr(10),
				},
				CapacityObjectID: to.StringPtr("<capacity-object-id>"),
			},
			SKU: &armpowerbidedicated.AutoScaleVCoreSKU{
				Name:     to.StringPtr("<name>"),
				Capacity: to.Int32Ptr(0),
				Tier:     armpowerbidedicated.VCoreSKUTierAutoScale.ToPtr(),
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("AutoScaleVCore.ID: %s\n", *res.ID)
}

// x-ms-original-file: specification/powerbidedicated/resource-manager/Microsoft.PowerBIdedicated/stable/2021-01-01/examples/deleteAutoScaleVCore.json
func ExampleAutoScaleVCoresClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpowerbidedicated.NewAutoScaleVCoresClient("<subscription-id>", cred, nil)
	_, err = client.Delete(ctx,
		"<resource-group-name>",
		"<vcore-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/powerbidedicated/resource-manager/Microsoft.PowerBIdedicated/stable/2021-01-01/examples/updateAutoScaleVCore.json
func ExampleAutoScaleVCoresClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpowerbidedicated.NewAutoScaleVCoresClient("<subscription-id>", cred, nil)
	res, err := client.Update(ctx,
		"<resource-group-name>",
		"<vcore-name>",
		armpowerbidedicated.AutoScaleVCoreUpdateParameters{
			Properties: &armpowerbidedicated.AutoScaleVCoreMutableProperties{
				CapacityLimit: to.Int32Ptr(20),
			},
			SKU: &armpowerbidedicated.AutoScaleVCoreSKU{
				Name:     to.StringPtr("<name>"),
				Capacity: to.Int32Ptr(0),
				Tier:     armpowerbidedicated.VCoreSKUTierAutoScale.ToPtr(),
			},
			Tags: map[string]*string{
				"testKey": to.StringPtr("testValue"),
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("AutoScaleVCore.ID: %s\n", *res.ID)
}

// x-ms-original-file: specification/powerbidedicated/resource-manager/Microsoft.PowerBIdedicated/stable/2021-01-01/examples/listAutoScaleVCoresInResourceGroup.json
func ExampleAutoScaleVCoresClient_ListByResourceGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpowerbidedicated.NewAutoScaleVCoresClient("<subscription-id>", cred, nil)
	_, err = client.ListByResourceGroup(ctx,
		"<resource-group-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/powerbidedicated/resource-manager/Microsoft.PowerBIdedicated/stable/2021-01-01/examples/listAutoScaleVCoresInSubscription.json
func ExampleAutoScaleVCoresClient_ListBySubscription() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armpowerbidedicated.NewAutoScaleVCoresClient("<subscription-id>", cred, nil)
	_, err = client.ListBySubscription(ctx,
		nil)
	if err != nil {
		log.Fatal(err)
	}
}
