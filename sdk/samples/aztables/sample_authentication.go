// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func AuthenticateWithSharedKey() *aztables.ServiceClient {
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
	client, err := aztables.NewServiceClient(serviceURL, cred, nil)
	check(err)
	return client
}

func AuthenticateWithTokenCredential() *aztables.ServiceClient {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := accountName + ".table.core.windows.net"

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	check(err)
	client, err := aztables.NewServiceClient(serviceURL, cred, nil)
	check(err)
	return client
}

func main() {
	AuthenticateWithTokenCredential()
	AuthenticateWithSharedKey()

	// sample_batch.go
	Sample_Batching()

	// sample_authentication.go
	CreateTableClient()
	CreateServiceClient()

	// sample_create_delete_table.go
	CreateDeleteFromServiceClient()
	CreateDeleteFromTableClient()

	// sample_insert_delete_entities.go
	InsertEntity()
	DeleteEntity()

	// sample_query_table.go
	Sample_QueryTable()

	// sample_query_tables.go
	Sample_QueryTables()

	// sample_update_entities.go
	UpdateEntities()
	MergeEntities()
	UpsertEntity()
}
