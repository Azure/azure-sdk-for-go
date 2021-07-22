package main

import (
	"context"
	"encoding/json"
	"os"
	"time"
)

// "github.com/Azure/azure-sdk-for-go/sdk/tables/aztable"

type MyEntity struct {
	aztable.Entity
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
	serviceURL := accountName + ".table.core.windows.net"

	cred := aztable.SharedKeyCredential(accountName, accountKey)
	client := aztable.NewTableClient("myTable", serviceURL, cred, nil)

	myEntity := MyEntity{
		aztable.Entity: aztable.Entity{
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
	myAdvancedEntity := aztable.EdmEntity{
		Entity: aztable.Entity{
			PartitionKey: "pk001",
			RowKey:       "rk002",
		},
		Properties: map[string]interface{}{
			"Bool":     false,
			"Int32":    int32(1234),
			"Int64":    aztable.EdmInt64(123456789012),
			"Double":   1234.1234,
			"String":   "test",
			"Guid":     aztable.EdmGuid("4185404a-5818-48c3-b9be-f217df0dba6f"),
			"DateTime": aztable.EdmDateTime(time.Date(2013, time.August, 02, 17, 37, 43, 9004348, time.UTC)),
			"Binary":   aztable.EdmBinary("SomeBinary"),
		},
	}

	marshalled, err = json.Marshal(myAdvancedEntity)
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
	serviceURL := accountName + ".table.core.windows.net"

	cred := aztable.SharedKeyCredential(accountName, accountKey)
	client := aztable.NewTableClient("myTable", serviceURL, cred, nil)

	_, err := client.DeleteEntity(context.Background(), "pk001", "rk001", "*")
	check(err)
}

func main() {
	InsertEntity()
	DeleteEntity()
}
