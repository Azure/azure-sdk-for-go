// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

type InventoryEntity struct {
	aztables.Entity
	Price       float32
	Inventory   int32
	ProductName string
	OnSale      bool
}

func InsertEntity() {
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

	myEntity := InventoryEntity{
		Entity: aztables.Entity{
			PartitionKey: "pk001",
			RowKey:       "rk001",
		},
		Price:       3.99,
		Inventory:   20,
		ProductName: "Markers",
		OnSale:      false,
	}
	marshalled, err := json.Marshal(myEntity)
	check(err)

	_, err = client.AddEntity(context.Background(), marshalled, nil)
	check(err)

	// Inserting an entity with int64s, binary, datetime, or guid types
	myAdvancedEntity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: "pk001",
			RowKey:       "rk002",
		},
		Properties: map[string]interface{}{
			"Bool":     false,
			"Int32":    int32(1234),
			"Int64":    aztables.EDMInt64(123456789012),
			"Double":   1234.1234,
			"String":   "test",
			"Guid":     aztables.EDMGUID("4185404a-5818-48c3-b9be-f217df0dba6f"),
			"DateTime": aztables.EDMDateTime(time.Date(2013, time.August, 02, 17, 37, 43, 9004348, time.UTC)),
			"Binary":   aztables.EDMBinary("SomeBinary"),
		},
	}

	marshalled, err = json.Marshal(myAdvancedEntity)
	check(err)
	_, err = client.AddEntity(context.Background(), marshalled, nil)
	check(err)
}

func DeleteEntity() {
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

	anyETag := azcore.ETagAny
	_, err = client.DeleteEntity(context.Background(), "pk001", "rk001", &aztables.DeleteEntityOptions{IfMatch: &anyETag})
	check(err)
}
