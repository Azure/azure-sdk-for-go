//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsql_test

import (
	"context"
	"log"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
)

// x-ms-original-file: specification/sql/resource-manager/Microsoft.Sql/preview/2020-11-01-preview/examples/SyncMemberGet.json
func ExampleSyncMembersClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armsql.NewSyncMembersClient("<subscription-id>", cred, nil)
	res, err := client.Get(ctx,
		"<resource-group-name>",
		"<server-name>",
		"<database-name>",
		"<sync-group-name>",
		"<sync-member-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response result: %#v\n", res.SyncMembersClientGetResult)
}

// x-ms-original-file: specification/sql/resource-manager/Microsoft.Sql/preview/2020-11-01-preview/examples/SyncMemberCreate.json
func ExampleSyncMembersClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armsql.NewSyncMembersClient("<subscription-id>", cred, nil)
	poller, err := client.BeginCreateOrUpdate(ctx,
		"<resource-group-name>",
		"<server-name>",
		"<database-name>",
		"<sync-group-name>",
		"<sync-member-name>",
		armsql.SyncMember{
			Properties: &armsql.SyncMemberProperties{
				DatabaseName:                      to.StringPtr("<database-name>"),
				DatabaseType:                      armsql.SyncMemberDbType("AzureSqlDatabase").ToPtr(),
				ServerName:                        to.StringPtr("<server-name>"),
				SyncDirection:                     armsql.SyncDirection("Bidirectional").ToPtr(),
				SyncMemberAzureDatabaseResourceID: to.StringPtr("<sync-member-azure-database-resource-id>"),
				UsePrivateLinkConnection:          to.BoolPtr(true),
				UserName:                          to.StringPtr("<user-name>"),
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
	log.Printf("Response result: %#v\n", res.SyncMembersClientCreateOrUpdateResult)
}

// x-ms-original-file: specification/sql/resource-manager/Microsoft.Sql/preview/2020-11-01-preview/examples/SyncMemberDelete.json
func ExampleSyncMembersClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armsql.NewSyncMembersClient("<subscription-id>", cred, nil)
	poller, err := client.BeginDelete(ctx,
		"<resource-group-name>",
		"<server-name>",
		"<database-name>",
		"<sync-group-name>",
		"<sync-member-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/sql/resource-manager/Microsoft.Sql/preview/2020-11-01-preview/examples/SyncMemberPatch.json
func ExampleSyncMembersClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armsql.NewSyncMembersClient("<subscription-id>", cred, nil)
	poller, err := client.BeginUpdate(ctx,
		"<resource-group-name>",
		"<server-name>",
		"<database-name>",
		"<sync-group-name>",
		"<sync-member-name>",
		armsql.SyncMember{
			Properties: &armsql.SyncMemberProperties{
				DatabaseName:                      to.StringPtr("<database-name>"),
				DatabaseType:                      armsql.SyncMemberDbType("AzureSqlDatabase").ToPtr(),
				ServerName:                        to.StringPtr("<server-name>"),
				SyncDirection:                     armsql.SyncDirection("Bidirectional").ToPtr(),
				SyncMemberAzureDatabaseResourceID: to.StringPtr("<sync-member-azure-database-resource-id>"),
				UsePrivateLinkConnection:          to.BoolPtr(true),
				UserName:                          to.StringPtr("<user-name>"),
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
	log.Printf("Response result: %#v\n", res.SyncMembersClientUpdateResult)
}

// x-ms-original-file: specification/sql/resource-manager/Microsoft.Sql/preview/2020-11-01-preview/examples/SyncMemberListBySyncGroup.json
func ExampleSyncMembersClient_ListBySyncGroup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armsql.NewSyncMembersClient("<subscription-id>", cred, nil)
	pager := client.ListBySyncGroup("<resource-group-name>",
		"<server-name>",
		"<database-name>",
		"<sync-group-name>",
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

// x-ms-original-file: specification/sql/resource-manager/Microsoft.Sql/preview/2020-11-01-preview/examples/SyncMemberGetSchema.json
func ExampleSyncMembersClient_ListMemberSchemas() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armsql.NewSyncMembersClient("<subscription-id>", cred, nil)
	pager := client.ListMemberSchemas("<resource-group-name>",
		"<server-name>",
		"<database-name>",
		"<sync-group-name>",
		"<sync-member-name>",
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

// x-ms-original-file: specification/sql/resource-manager/Microsoft.Sql/preview/2020-11-01-preview/examples/SyncMemberRefreshSchema.json
func ExampleSyncMembersClient_BeginRefreshMemberSchema() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armsql.NewSyncMembersClient("<subscription-id>", cred, nil)
	poller, err := client.BeginRefreshMemberSchema(ctx,
		"<resource-group-name>",
		"<server-name>",
		"<database-name>",
		"<sync-group-name>",
		"<sync-member-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}
