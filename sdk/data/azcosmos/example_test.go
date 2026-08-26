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

// exampleContainer builds the client and container the operation examples below work against, so
// that each of those can show the operation rather than repeating the setup. ExampleNewClient
// shows the construction it stands in for.
//
// The returned function closes the client. It is called rather than deferred because a deferred
// call would not run if a later log.Fatalf fired.
func exampleContainer() (*azcosmos.ContainerClient, func()) {
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

	return container, func() {
		if err := client.Close(); err != nil {
			// TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}
	}
}

// A client is safe for concurrent use and holds the caches that make requests cheap, so create one
// per account and keep it for the lifetime of the application rather than one per operation.
//
// Close is called directly rather than deferred here only because log.Fatalf would skip a deferred
// call. In an application that returns errors rather than exiting, defer it.
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

// An account key authenticates too, though Entra ID is preferable where it is available.
//
// Routing is worth setting: naming the region the application runs in lets the SDK order the
// account's regions by proximity to it, rather than leaving the order to the account.
func ExampleNewClientWithKey() {
	cred, err := azcosmos.NewKeyCredential("myAccountKey")
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	client, err := azcosmos.NewClientWithKey("https://myaccount.documents.azure.com", cred, &azcosmos.ClientOptions{
		Routing: azcosmos.ProximityTo(azcosmos.RegionEastUS),
	})
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	if err := client.Close(); err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
}

// When the regions to prefer are known ahead of time, give them explicitly. The order is a
// preference, not a restriction: once it is exhausted the client may still use other regions.
func ExamplePreferredRegions() {
	options := &azcosmos.ClientOptions{
		Routing: azcosmos.PreferredRegions(azcosmos.RegionEastUS, azcosmos.RegionWestUS),
	}

	log.Printf("client will prefer East US, then West US: %+v", options.Routing)
}

func ExampleContainerClient_CreateItem() {
	container, closeClient := exampleContainer()

	item, err := json.Marshal(map[string]string{
		"id":           "item-1",
		"categoryName": "gear-surf-surfboards",
	})
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	// The partition key value must match the value of the partition key path in the item, and the
	// id must match the item's own id property.
	pk := azcosmos.NewPartitionKeyString("gear-surf-surfboards")

	response, err := container.CreateItem(context.TODO(), pk, "item-1", item, nil)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	log.Printf("created the item, charged %v RU", response.RequestCharge)

	closeClient()
}

func ExampleContainerClient_ReadItem() {
	container, closeClient := exampleContainer()
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

	closeClient()
}

// Reads can relax the account's consistency level, and can carry a session token captured from a
// write elsewhere so that the read is guaranteed to observe it.
func ExampleContainerClient_ReadItem_sessionConsistency() {
	container, closeClient := exampleContainer()

	options := &azcosmos.ReadItemOptions{
		Operation: azcosmos.OperationOptions{
			ConsistencyStrategy: azcosmos.ReadConsistencyStrategySession,
		},
		// Taken from the ItemResponse of the write this read needs to observe.
		SessionToken: "0:-1#42",
	}

	pk := azcosmos.NewPartitionKeyString("gear-surf-surfboards")

	if _, err := container.ReadItem(context.TODO(), pk, "item-1", options); err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	closeClient()
}

// A missing item is reported as an *azcosmos.Error with the CodeNotFound code, which is usually
// something to handle rather than to fail on.
func ExampleError() {
	container, closeClient := exampleContainer()
	pk := azcosmos.NewPartitionKeyString("gear-surf-surfboards")

	_, err := container.ReadItem(context.TODO(), pk, "item-1", nil)

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

	closeClient()
}

// A container with a hierarchical partition key takes one component per level of its definition,
// in the order the paths are declared.
func ExamplePartitionKey_hierarchical() {
	pk := azcosmos.NewPartitionKeyString("Contoso").
		AppendString("Redmond").
		AppendNumber(98052)

	log.Printf("built a partition key with %d components", pk.Len())
}
