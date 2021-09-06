// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func Sample_QueryTables() {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY could not be found")
	}
	serviceURL := accountName + ".table.core.windows.net"

	cred, err := aztables.NewSharedKeyCredential(accountName, accountKey)
	check(err)
	service, err := aztables.NewServiceClient(serviceURL, cred, nil)
	check(err)

	myTable := "myTableName"
	filter := fmt.Sprintf("TableName ge '%v'", myTable)
	pager := service.ListTables(&aztables.ListTablesOptions{Filter: &filter})

	pageCount := 1
	for pager.NextPage(context.Background()) {
		response := pager.PageResponse()
		fmt.Printf("There are %d tables in page #%d\n", len(response.Tables), pageCount)
		for _, table := range response.Tables {
			fmt.Printf("\tTableName: %s\n", *table.TableName)
		}
		pageCount += 1
	}

	check(pager.Err())
}
