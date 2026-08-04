// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos_test

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
)

// A client is safe for concurrent use and holds the caches that make requests cheap, so create one
// per account and keep it for the lifetime of the application rather than one per operation.
//
// Close is called directly rather than deferred in these examples only because log.Fatalf would
// skip a deferred call. In an application that returns errors rather than exiting, defer it.
func ExampleNewClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	client, err := azcosmos.NewClient("https://myaccount.documents.azure.com", cred, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	container, err := client.NewContainer("myDatabase", "myContainer")
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	log.Printf("ready to operate on container %s", container.ID())

	if err := client.Close(); err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
}

func ExampleNewClientFromConnectionString() {
	client, err := azcosmos.NewClientFromConnectionString(
		"AccountEndpoint=https://myaccount.documents.azure.com;AccountKey=myAccountKey;", nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	if err := client.Close(); err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
}

func ExampleContainerClient_CreateItem() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	client, err := azcosmos.NewClient("https://myaccount.documents.azure.com", cred, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	container, err := client.NewContainer("myDatabase", "myContainer")
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	item, err := json.Marshal(map[string]string{
		"id":           "item-1",
		"categoryName": "gear-surf-surfboards",
	})
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// The partition key value must match the value of the partition key path in the item.
	pk := azcosmos.NewPartitionKeyString("gear-surf-surfboards")

	response, err := container.CreateItem(context.TODO(), pk, item, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	log.Printf("created the item, charged %v RU", response.RequestCharge)
}

func ExampleContainerClient_ReadItem() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	client, err := azcosmos.NewClient("https://myaccount.documents.azure.com", cred, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	container, err := client.NewContainer("myDatabase", "myContainer")
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	pk := azcosmos.NewPartitionKeyString("gear-surf-surfboards")

	response, err := container.ReadItem(context.TODO(), pk, "item-1", nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	var item map[string]any
	if err := json.Unmarshal(response.Value, &item); err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	log.Printf("read item %v, charged %v RU", item["id"], response.RequestCharge)
}

// Reads can relax the account's consistency level, and can carry a session token captured from a
// write elsewhere so that the read is guaranteed to observe it.
func ExampleContainerClient_ReadItem_sessionConsistency() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	client, err := azcosmos.NewClient("https://myaccount.documents.azure.com", cred, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	container, err := client.NewContainer("myDatabase", "myContainer")
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	options := &azcosmos.ReadItemOptions{
		ConsistencyStrategy: azcosmos.ReadConsistencyStrategySession,
		// Taken from the ItemResponse of the write this read needs to observe.
		SessionToken: "0:-1#42",
	}

	pk := azcosmos.NewPartitionKeyString("gear-surf-surfboards")

	if _, err := container.ReadItem(context.TODO(), pk, "item-1", options); err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
}

// A missing item is reported as an *azcosmos.Error with the CodeNotFound code, which is usually
// something to handle rather than to fail on.
func ExampleError() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	client, err := azcosmos.NewClient("https://myaccount.documents.azure.com", cred, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	container, err := client.NewContainer("myDatabase", "myContainer")
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	pk := azcosmos.NewPartitionKeyString("gear-surf-surfboards")

	_, err = container.ReadItem(context.TODO(), pk, "item-1", nil)

	var cosmosErr *azcosmos.Error
	switch {
	case err == nil:
		log.Println("read the item")
	case errors.As(err, &cosmosErr) && cosmosErr.Code == azcosmos.CodeNotFound:
		log.Println("no such item")
	default:
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
}

// A container with a hierarchical partition key takes one component per level of its definition,
// in the order the paths are declared.
func ExamplePartitionKey_hierarchical() {
	pk := azcosmos.NewPartitionKeyString("Contoso").
		AppendString("Redmond").
		AppendNumber(98052)

	log.Printf("built a partition key with %d components", pk.Len())
}
