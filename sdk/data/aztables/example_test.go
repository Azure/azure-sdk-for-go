// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func ExampleNewSharedKeyCredential() {
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
	if err != nil {
		panic(err)
	}
	client, err := aztables.NewServiceClientWithSharedKey(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(client)
}

func ExampleNewServiceClient() {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := accountName + ".table.core.windows.net"

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
	client, err := aztables.NewServiceClient(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(client)
}

func ExampleNewServiceClientWithSharedKey() {
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
	if err != nil {
		panic(err)
	}
	client, err := aztables.NewServiceClientWithSharedKey(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(client)
}

func ExampleNewServiceClientWithNoCredential() {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	sharedAccessSignature, ok := os.LookupEnv("TABLES_SHARED_ACCESS_SIGNATURE")
	if !ok {
		panic("TABLES_SHARED_ACCESS_SIGNATURE could not be found")
	}
	serviceURL := fmt.Sprintf("%s.table.core.windows.net/?%s", accountName, sharedAccessSignature)

	client, err := aztables.NewServiceClientWithNoCredential(serviceURL, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(client)
}

type MyEntity struct {
	aztables.Entity
	Price       float32
	Inventory   int32
	ProductName string
	OnSale      bool
}

func ExampleClient_SubmitTransaction() {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net/%s", accountName, "tableName")

	cred, err := aztables.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		panic(err)
	}
	client, err := aztables.NewClientWithSharedKey(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}

	batch := []aztables.TransactionAction{}

	entity1 := MyEntity{
		Entity: aztables.Entity{
			PartitionKey: "pk001",
			RowKey:       "rk001",
		},
		Price:       3.99,
		Inventory:   10,
		ProductName: "Pens",
		OnSale:      false,
	}
	marshalled, err := json.Marshal(entity1)
	if err != nil {
		panic(err)
	}
	batch = append(batch, aztables.TransactionAction{
		ActionType: aztables.Add,
		Entity:     marshalled,
	})

	entity2 := MyEntity{
		Entity: aztables.Entity{
			PartitionKey: "pk001",
			RowKey:       "rk002",
		},
		Price:       19.99,
		Inventory:   15,
		ProductName: "Calculators",
		OnSale:      false,
	}
	marshalled, err = json.Marshal(entity2)
	if err != nil {
		panic(err)
	}
	batch = append(batch, aztables.TransactionAction{
		ActionType: aztables.UpdateMerge,
		Entity:     marshalled,
	})

	entity3 := MyEntity{
		Entity: aztables.Entity{
			PartitionKey: "pk001",
			RowKey:       "rk003",
		},
		Price:       0.99,
		Inventory:   150,
		ProductName: "Pens",
		OnSale:      true,
	}
	marshalled, err = json.Marshal(entity3)
	if err != nil {
		panic(err)
	}
	batch = append(batch, aztables.TransactionAction{
		ActionType: aztables.InsertReplace,
		Entity:     marshalled,
	})

	entity4 := MyEntity{
		Entity: aztables.Entity{
			PartitionKey: "pk001",
			RowKey:       "rk004",
		},
		Price:       3.99,
		Inventory:   150,
		ProductName: "100ct Paper Clips",
		OnSale:      false,
	}
	marshalled, err = json.Marshal(entity4)
	if err != nil {
		panic(err)
	}
	batch = append(batch, aztables.TransactionAction{
		ActionType: aztables.Delete,
		Entity:     marshalled,
	})

	resp, err := client.SubmitTransaction(context.Background(), batch, nil)
	if err != nil {
		panic(err)
	}

	for _, subResp := range *resp.TransactionResponses {
		if subResp.StatusCode != http.StatusAccepted {
			fmt.Println(subResp.Body)
		}
	}
}

func ExampleServiceClient_CreateTable() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net", accountName)

	service, err := aztables.NewServiceClient(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}

	// Create a table
	_, err = service.CreateTable(context.Background(), "fromServiceClient", nil)
	if err != nil {
		panic(err)
	}
}

func ExampleServiceClient_DeleteTable() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net", accountName)

	service, err := aztables.NewServiceClient(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}

	// Delete a table
	_, err = service.DeleteTable(context.Background(), "fromServiceClient", nil)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_Create() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net/%s", accountName, "fromTableClient")
	client, err := aztables.NewClient(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}

	// Create a table
	_, err = client.Create(context.Background(), nil)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net/%s", accountName, "fromTableClient")
	client, err := aztables.NewClient(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}

	// Delete a table
	_, err = client.Delete(context.Background(), nil)
	if err != nil {
		panic(err)
	}
}

func ExampleNewClient() {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net/%s", accountName, "myTableName")

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
	client, err := aztables.NewClient(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(client)
}

func ExampleNewClientWithSharedKey() {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("TABLES_PRIMARY_STORAGE_ACCOUNT_KEY could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net/%s", accountName, "myTableName")

	cred, err := aztables.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		panic(err)
	}
	client, err := aztables.NewClientWithSharedKey(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(client)
}

type InventoryEntity struct {
	aztables.Entity
	Price       float32
	Inventory   int32
	ProductName string
	OnSale      bool
}

func ExampleClient_InsertEntity() {
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
	if err != nil {
		panic(err)
	}
	client, err := aztables.NewClientWithSharedKey(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}

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
	if err != nil {
		panic(err)
	}

	_, err = client.AddEntity(context.Background(), marshalled, nil)
	if err != nil {
		panic(err)
	}

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
	if err != nil {
		panic(err)
	}
	_, err = client.AddEntity(context.Background(), marshalled, nil)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_DeleteEntity() {
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
	if err != nil {
		panic(err)
	}
	client, err := aztables.NewClientWithSharedKey(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}

	anyETag := azcore.ETagAny
	_, err = client.DeleteEntity(context.Background(), "pk001", "rk001", &aztables.DeleteEntityOptions{IfMatch: &anyETag})
	if err != nil {
		panic(err)
	}
}

func ExampleClient_List() {
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
	if err != nil {
		panic(err)
	}
	client, err := aztables.NewClientWithSharedKey(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}

	filter := fmt.Sprintf("PartitionKey eq '%v' or PartitionKey eq '%v'", "pk001", "pk002")
	pager := client.List(&aztables.ListEntitiesOptions{Filter: &filter})

	pageCount := 1
	for pager.NextPage(context.Background()) {
		response := pager.PageResponse()
		fmt.Printf("There are %d entities in page #%d\n", len(response.Entities), pageCount)
		pageCount += 1
	}
	if err := pager.Err(); err != nil {
		panic(err)
	}

	// To list all entities in a table, provide nil to Query()
	listPager := client.List(nil)
	pageCount = 1
	for listPager.NextPage(context.Background()) {
		response := listPager.PageResponse()
		fmt.Printf("There are %d entities in page #%d\n", len(response.Entities), pageCount)
		pageCount += 1
	}
	if err := pager.Err(); err != nil {
		panic(err)
	}
}

func ExampleServiceClient_ListTables() {
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
	if err != nil {
		panic(err)
	}
	service, err := aztables.NewServiceClientWithSharedKey(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}

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
	if err := pager.Err(); err != nil {
		panic(err)
	}
}
