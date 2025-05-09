// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
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

func ExampleNewServiceClient_sovereignCloud() {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := accountName + ".table.core.windows.net"

	cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: cloud.AzureChina,
		},
	})
	if err != nil {
		panic(err)
	}

	client, err := aztables.NewServiceClient(serviceURL, cred, &aztables.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: cloud.AzureChina,
		},
	})
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

func ExampleServiceClient_GetAccountSASURL() {
	cred, err := aztables.NewSharedKeyCredential("myAccountName", "myAccountKey")
	if err != nil {
		panic(err)
	}
	service, err := aztables.NewServiceClientWithSharedKey("https://<myAccountName>.table.core.windows.net", cred, nil)
	if err != nil {
		panic(err)
	}

	resources := aztables.AccountSASResourceTypes{Service: true}
	permission := aztables.AccountSASPermissions{Read: true}
	start := time.Now()
	expiry := start.AddDate(1, 0, 0)
	sasURL, err := service.GetAccountSASURL(resources, permission, start, expiry)
	if err != nil {
		panic(err)
	}

	serviceURL := fmt.Sprintf("https://<myAccountName>.table.core.windows.net/?%s", sasURL)
	sasService, err := aztables.NewServiceClientWithNoCredential(serviceURL, nil)
	if err != nil {
		panic(err)
	}
	_ = sasService
}

type MyEntity struct {
	aztables.Entity
	Value int
}

func ExampleClient_SubmitTransaction() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net/%s", "myAccountName", "tableName")
	client, err := aztables.NewClient(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}

	batch := []aztables.TransactionAction{}

	baseEntity := MyEntity{
		Entity: aztables.Entity{
			PartitionKey: "myPartitionKey",
			RowKey:       "",
		},
	}
	for i := 0; i < 10; i++ {
		baseEntity.RowKey = fmt.Sprintf("rk-%d", i)
		baseEntity.Value = i
		marshalled, err := json.Marshal(baseEntity)
		if err != nil {
			panic(err)
		}
		batch = append(batch, aztables.TransactionAction{
			ActionType: aztables.TransactionTypeAdd,
			Entity:     marshalled,
		})
	}

	_, err = client.SubmitTransaction(context.TODO(), batch, nil)
	if err != nil {
		var httpErr *azcore.ResponseError
		if errors.As(err, &httpErr) {
			body, err := io.ReadAll(httpErr.RawResponse.Body)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(body)) // Do some parsing of the body
		} else {
			panic(err)
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
	_, err = service.CreateTable(context.TODO(), "fromServiceClient", nil)
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
	_, err = service.DeleteTable(context.TODO(), "fromServiceClient", nil)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_CreateTable() {
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
	_, err = client.CreateTable(context.TODO(), nil)
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
	_, err = client.Delete(context.TODO(), nil)
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

func ExampleClient_UpsertEntity() {
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

	_, err = client.AddEntity(context.TODO(), marshalled, nil)
	if err != nil {
		panic(err)
	}

	// Inserting an entity with int64s, binary, datetime, or guid types
	myAdvancedEntity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: "pk001",
			RowKey:       "rk002",
		},
		Properties: map[string]any{
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
	_, err = client.AddEntity(context.TODO(), marshalled, nil)
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
	_, err = client.DeleteEntity(context.TODO(), "pk001", "rk001", &aztables.DeleteEntityOptions{IfMatch: &anyETag})
	if err != nil {
		panic(err)
	}
}

func ExampleClient_NewListEntitiesPager() {
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

	// For more information about writing query strings, check out:
	//  - API Documentation: https://learn.microsoft.com/rest/api/storageservices/querying-tables-and-entities
	//  - README samples: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/data/aztables/README.md#writing-filters
	filter := fmt.Sprintf("PartitionKey eq '%s' or PartitionKey eq '%s'", "pk001", "pk002")
	pager := client.NewListEntitiesPager(&aztables.ListEntitiesOptions{Filter: &filter})

	pageCount := 1
	for pager.More() {
		response, err := pager.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		fmt.Printf("There are %d entities in page #%d\n", len(response.Entities), pageCount)

		for _, entity := range response.Entities {
			var myEntity aztables.EDMEntity
			err = json.Unmarshal(entity, &myEntity)
			if err != nil {
				panic(err)
			}

			sp := myEntity.Properties["String"].(string)
			dp := myEntity.Properties["Double"].(float64)
			dt := myEntity.Properties["DateTime"].(aztables.EDMDateTime)
			t1 := time.Time(dt)

			fmt.Printf("Received: %s, %s, %s, %.2f, %s", myEntity.PartitionKey, myEntity.RowKey, sp, dp, t1.String())
		}
		pageCount += 1
	}

	// To list all entities in a table, provide nil to Query()
	listPager := client.NewListEntitiesPager(nil)
	pageCount = 0
	for listPager.More() {
		response, err := listPager.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		fmt.Printf("There are %d entities in page #%d\n", len(response.Entities), pageCount)
		pageCount += 1

		for _, entity := range response.Entities {
			var myEntity aztables.EDMEntity
			err = json.Unmarshal(entity, &myEntity)
			if err != nil {
				panic(err)
			}

			sp := myEntity.Properties["String"].(string)
			dp := myEntity.Properties["Double"].(float64)
			dt := myEntity.Properties["DateTime"].(aztables.EDMDateTime)
			t1 := time.Time(dt)

			fmt.Printf("Received: %s, %s, %s, %.2f, %s", myEntity.PartitionKey, myEntity.RowKey, sp, dp, t1.String())
		}
	}
}

func ExampleServiceClient_NewListTablesPager() {
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
	pager := service.NewListTablesPager(&aztables.ListTablesOptions{Filter: &filter})

	pageCount := 1
	for pager.More() {
		response, err := pager.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		fmt.Printf("There are %d tables in page #%d\n", len(response.Tables), pageCount)
		for _, table := range response.Tables {
			fmt.Printf("\tTableName: %s\n", *table.Name)
		}
		pageCount += 1
	}
}

func ExampleServiceClient_SetProperties() {
	accountName, ok := os.LookupEnv("TABLES_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("TABLES_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := accountName + ".table.core.windows.net"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
	service, err := aztables.NewServiceClient(serviceURL, cred, nil)
	if err != nil {
		panic(err)
	}

	getResp, err := service.GetProperties(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	getResp.HourMetrics = &aztables.Metrics{
		Enabled: to.Ptr(true),
	}
	getResp.Logging = &aztables.Logging{
		Delete: to.Ptr(true),
		Read:   to.Ptr(true),
		Write:  to.Ptr(true),
	}
	getResp.Cors = append(getResp.Cors, &aztables.CorsRule{
		AllowedHeaders: to.Ptr("x-allowed-header"),
		AllowedMethods: to.Ptr("POST,GET"),
	})

	_, err = service.SetProperties(context.TODO(), getResp.ServiceProperties, nil)
	if err != nil {
		panic(err)
	}
}
