//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armavs_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/avs/armavs"
)

// x-ms-original-file: specification/vmware/resource-manager/Microsoft.AVS/stable/2021-12-01/examples/HcxEnterpriseSites_List.json
func ExampleHcxEnterpriseSitesClient_List() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armavs.NewHcxEnterpriseSitesClient("<subscription-id>", cred, nil)
	pager := client.List("<resource-group-name>",
		"<private-cloud-name>",
		nil)
	for pager.NextPage(ctx) {
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("HcxEnterpriseSite.ID: %s\n", *v.ID)
		}
	}
}

// x-ms-original-file: specification/vmware/resource-manager/Microsoft.AVS/stable/2021-12-01/examples/HcxEnterpriseSites_Get.json
func ExampleHcxEnterpriseSitesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armavs.NewHcxEnterpriseSitesClient("<subscription-id>", cred, nil)
	res, err := client.Get(ctx,
		"<resource-group-name>",
		"<private-cloud-name>",
		"<hcx-enterprise-site-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("HcxEnterpriseSite.ID: %s\n", *res.ID)
}

// x-ms-original-file: specification/vmware/resource-manager/Microsoft.AVS/stable/2021-12-01/examples/HcxEnterpriseSites_CreateOrUpdate.json
func ExampleHcxEnterpriseSitesClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armavs.NewHcxEnterpriseSitesClient("<subscription-id>", cred, nil)
	res, err := client.CreateOrUpdate(ctx,
		"<resource-group-name>",
		"<private-cloud-name>",
		"<hcx-enterprise-site-name>",
		armavs.HcxEnterpriseSite{},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("HcxEnterpriseSite.ID: %s\n", *res.ID)
}

// x-ms-original-file: specification/vmware/resource-manager/Microsoft.AVS/stable/2021-12-01/examples/HcxEnterpriseSites_Delete.json
func ExampleHcxEnterpriseSitesClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armavs.NewHcxEnterpriseSitesClient("<subscription-id>", cred, nil)
	_, err = client.Delete(ctx,
		"<resource-group-name>",
		"<private-cloud-name>",
		"<hcx-enterprise-site-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
}
