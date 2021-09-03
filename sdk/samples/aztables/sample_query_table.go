// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func Sample_QueryTable() {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net/%s", accountName, "myTable")

	cred, err := aztables.NewSharedKeyCredential(accountName, accountKey)
	check(err)
	client, err := aztables.NewClient(serviceURL, cred, nil)
	check(err)

	filter := fmt.Sprintf("PartitionKey eq '%v' or PartitionKey eq '%v'", "pk001", "pk002")
	pager := client.List(&aztables.ListEntitiesOptions{Filter: &filter})

	pageCount := 1
	for pager.NextPage(context.Background()) {
		response := pager.PageResponse()
		fmt.Printf("There are %d entities in page #%d\n", len(response.Entities), pageCount)
		pageCount += 1
	}
	check(pager.Err())

	// To list all entities in a table, provide nil to Query()
	listPager := client.List(nil)
	pageCount = 1
	for listPager.NextPage(context.Background()) {
		response := listPager.PageResponse()
		fmt.Printf("There are %d entities in page #%d\n", len(response.Entities), pageCount)
		pageCount += 1
	}
	check(pager.Err())
}
