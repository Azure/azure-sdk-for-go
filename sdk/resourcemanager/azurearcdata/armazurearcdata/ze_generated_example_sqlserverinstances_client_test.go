//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armazurearcdata_test

import (
	"context"
	"log"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/azurearcdata/armazurearcdata"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/azurearcdata/resource-manager/Microsoft.AzureArcData/preview/2022-03-01-preview/examples/ListSubscriptionSqlServerInstance.json
func ExampleSQLServerInstancesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armazurearcdata.NewSQLServerInstancesClient("00000000-1111-2222-3333-444444444444", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListPager(nil)
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/azurearcdata/resource-manager/Microsoft.AzureArcData/preview/2022-03-01-preview/examples/ListByResourceGroupSqlServerInstance.json
func ExampleSQLServerInstancesClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armazurearcdata.NewSQLServerInstancesClient("00000000-1111-2222-3333-444444444444", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListByResourceGroupPager("testrg",
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/azurearcdata/resource-manager/Microsoft.AzureArcData/preview/2022-03-01-preview/examples/GetSqlServerInstance.json
func ExampleSQLServerInstancesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armazurearcdata.NewSQLServerInstancesClient("00000000-1111-2222-3333-444444444444", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"testrg",
		"testsqlServerInstance",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/azurearcdata/resource-manager/Microsoft.AzureArcData/preview/2022-03-01-preview/examples/CreateOrUpdateSqlServerInstance.json
func ExampleSQLServerInstancesClient_BeginCreate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armazurearcdata.NewSQLServerInstancesClient("00000000-1111-2222-3333-444444444444", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginCreate(ctx,
		"testrg",
		"testsqlServerInstance",
		armazurearcdata.SQLServerInstance{
			Location: to.Ptr("northeurope"),
			Tags: map[string]*string{
				"mytag": to.Ptr("myval"),
			},
			Properties: &armazurearcdata.SQLServerInstanceProperties{
				AzureDefenderStatus:            to.Ptr(armazurearcdata.DefenderStatusProtected),
				AzureDefenderStatusLastUpdated: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-01-02T17:18:19.1234567Z"); return t }()),
				Collation:                      to.Ptr("collation"),
				ContainerResourceID:            to.Ptr("Resource id of hosting Arc Machine"),
				CurrentVersion:                 to.Ptr("2012"),
				Edition:                        to.Ptr(armazurearcdata.EditionTypeDeveloper),
				HostType:                       to.Ptr(armazurearcdata.HostTypePhysicalServer),
				InstanceName:                   to.Ptr("name of instance"),
				LicenseType:                    to.Ptr(armazurearcdata.ArcSQLServerLicenseTypeFree),
				PatchLevel:                     to.Ptr("patchlevel"),
				ProductID:                      to.Ptr("sql id"),
				Status:                         to.Ptr(armazurearcdata.ConnectionStatusRegistered),
				TCPDynamicPorts:                to.Ptr("1433"),
				TCPStaticPorts:                 to.Ptr("1433"),
				VCore:                          to.Ptr("4"),
				Version:                        to.Ptr(armazurearcdata.SQLVersionSQLServer2012),
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/azurearcdata/resource-manager/Microsoft.AzureArcData/preview/2022-03-01-preview/examples/DeleteSqlServerInstance.json
func ExampleSQLServerInstancesClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armazurearcdata.NewSQLServerInstancesClient("00000000-1111-2222-3333-444444444444", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginDelete(ctx,
		"testrg",
		"testsqlServerInstance",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/azurearcdata/resource-manager/Microsoft.AzureArcData/preview/2022-03-01-preview/examples/UpdateSqlServerInstance.json
func ExampleSQLServerInstancesClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armazurearcdata.NewSQLServerInstancesClient("00000000-1111-2222-3333-444444444444", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Update(ctx,
		"testrg",
		"testsqlServerInstance",
		armazurearcdata.SQLServerInstanceUpdate{
			Tags: map[string]*string{
				"mytag": to.Ptr("myval"),
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}
