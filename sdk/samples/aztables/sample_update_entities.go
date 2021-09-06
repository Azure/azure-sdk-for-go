// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func UpdateEntities() {
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

	// 1. First add an entity
	myEntity := MyEntity{
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

	// 2. Update the entity object
	myEntity.OnSale = true
	myEntity.Price = 3.49

	// 3. Update the entity on the service
	marshalled, err = json.Marshal(myEntity)
	check(err)
	_, err = client.UpdateEntity(context.Background(), marshalled, &aztables.UpdateEntityOptions{UpdateMode: aztables.MergeEntity})
	check(err)
}

func MergeEntities() {
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

	marshalled, err := json.Marshal(myAdvancedEntity)
	check(err)
	_, err = client.AddEntity(context.Background(), marshalled, nil)
	check(err)

	// Delete properties
	delete(myAdvancedEntity.Properties, "Guid")
	delete(myAdvancedEntity.Properties, "String")
	marshalled, err = json.Marshal(myAdvancedEntity)
	check(err)
	_, err = client.UpdateEntity(context.Background(), marshalled, &aztables.UpdateEntityOptions{UpdateMode: aztables.ReplaceEntity})
	check(err)
}

func UpsertEntity() {
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

	// 1. First insert an entity
	myEntity := MyEntity{
		Entity: aztables.Entity{
			PartitionKey: "pk001",
			RowKey:       "rk001",
		},
		Price:       1.99,
		Inventory:   35,
		ProductName: "Pens",
		OnSale:      true,
	}
	marshalled, err := json.Marshal(myEntity)
	check(err)

	_, err = client.InsertEntity(context.Background(), marshalled, &aztables.InsertEntityOptions{UpdateMode: aztables.MergeEntity})
	check(err)

	// 2. Update the entity object
	myEntity.OnSale = false

	// 3. Insert the entity on the service. MergeEntity will merge entity properties, whereas ReplaceEntity will replace entity properties
	marshalled, err = json.Marshal(myEntity)
	check(err)
	_, err = client.InsertEntity(context.Background(), marshalled, &aztables.InsertEntityOptions{UpdateMode: aztables.MergeEntity})
	check(err)
}
