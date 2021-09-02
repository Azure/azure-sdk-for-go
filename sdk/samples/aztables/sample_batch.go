package main

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/tables/aztable"
)

type MyEntity struct {
	aztable.Entity
	Price       float32
	Inventory   int32
	ProductName string
	OnSale      bool
}

func Sample_Batching() {
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
	client := aztable.NewTableClient("tableName", serviceURL, cred, nil)

	entity1 := MyEntity{
		Entity: aztable.Entity{
			PartitionKey: "pk001",
			RowKey:       "rk001",
		},
		Price:       3.99,
		Inventory:   10,
		ProductName: "Pens",
		OnSale:      false,
	}

	entity2 := MyEntity{
		Entity: aztable.Entity{
			PartitionKey: "pk001",
			RowKey:       "rk002",
		},
		Price:       19.99,
		Inventory:   15,
		ProductName: "Calculators",
		OnSale:      false,
	}

	entity3 := MyEntity{
		Entity: aztable.Entity{
			PartitionKey: "pk001",
			RowKey:       "rk003",
		},
		Price:       0.99,
		Inventory:   150,
		ProductName: "Pens",
		OnSale:      true,
	}

	entity4 := MyEntity{
		Entity: aztable.Entity{
			PartitionKey: "pk001",
			RowKey:       "rk004",
		},
		Price:       3.99,
		Inventory:   150,
		ProductName: "100ct Paper Clips",
		OnSale:      false,
	}

	batch := make([]aztable.TableTransactionAction, 4)
	batch[0] = aztable.TableTransactionAction{
		ActionType: aztable.Add,
		Entity:     entity1,
	}
	batch[1] = aztable.TableTransactionAction{
		ActionType: aztable.UpdateMerge,
		Entity:     entity2,
	}
	batch[2] = aztable.TableTransactionAction{
		ActionType: aztable.UpsertReplace,
		Entity:     entity3,
	}
	batch[3] = aztable.TableTransactionAction{
		ActionType: aztable.Delete,
		Entity:     entity4,
	}

	_, err := client.SubmitTransaction(context.Background(), batch, nil)
	check(err)
}
