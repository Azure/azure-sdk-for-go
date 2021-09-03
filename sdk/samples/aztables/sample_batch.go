// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

type MyEntity struct {
	aztables.Entity
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
	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net/%s", accountName, "tableName")

	cred, err := aztables.NewSharedKeyCredential(accountName, accountKey)
	check(err)
	client, err := aztables.NewClient(serviceURL, cred, nil)
	check(err)

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
	check(err)
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
	check(err)
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
	check(err)
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
	check(err)
	batch = append(batch, aztables.TransactionAction{
		ActionType: aztables.Delete,
		Entity:     marshalled,
	})

	resp, err := client.SubmitTransaction(context.Background(), batch, nil)
	check(err)

	for _, subResp := range *resp.TransactionResponses {
		if subResp.StatusCode != http.StatusAccepted {
			fmt.Println(subResp.Body)
		}
	}
}
