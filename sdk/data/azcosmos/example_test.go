// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// cSpell:ignore Contoso upserted

package azcosmos_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
)

type catalogItem struct {
	ID           string `json:"id"`
	CategoryName string `json:"categoryName"`
	Quantity     int    `json:"quantity"`
}

func ExampleContainerClient_NewQueryItemsPager() {
	container, closeClient := exampleContainer()
	query, err := azcosmos.NewQuery("SELECT * FROM c WHERE c.price >= @minimum").
		WithParameter("@minimum", 10)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	scope := azcosmos.NewFeedScopeForPartitionKey(azcosmos.NewPartitionKeyString("catalog"))
	pager := container.NewQueryItemsPager[catalogItem](query, scope, &azcosmos.QueryOptions{
		Feed: azcosmos.FeedOptions{PageSizeHint: 25},
	})
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			closeErr := pager.Close()
			closeClient()
			// TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", errors.Join(err, closeErr))
		}
		log.Printf("received %d items, charge %.2f RU", len(page.Items), page.RequestCharge)
		// To pause, call pager.ContinuationToken(ctx) before exhaustion, then Close.
	}
	if err := pager.Close(); err != nil {
		closeClient()
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	closeClient()
}

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
	if err := client.Initialize(context.TODO()); err != nil {
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
	if err := client.Initialize(context.TODO()); err != nil {
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
// Routing is worth setting when the preferred region order is known.
func ExampleNewClientWithKey() {
	cred, err := azcosmos.NewKeyCredential("bXktYWNjb3VudC1rZXk=")
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	client, err := azcosmos.NewClientWithKey("https://myaccount.documents.azure.com", cred, &azcosmos.ClientOptions{
		Routing: azcosmos.PreferredRegions(azcosmos.RegionEastUS),
	})

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	if err := client.Initialize(context.TODO()); err != nil {
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

	item := catalogItem{ID: "item-1", CategoryName: "gear-surf-surfboards"}

	// The partition key value must match the value of the partition key path in the item, and the
	// id must match the item's own id property.
	pk := azcosmos.NewPartitionKeyString("gear-surf-surfboards")

	response, err := container.CreateItem(context.TODO(), pk, "item-1", item, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
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

	item, err := response.IntoModel[catalogItem]()
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	log.Printf("read item %v, charged %v RU", item.ID, response.RequestCharge)

	closeClient()
}

func ExampleContainerClient_ReplaceItem() {
	container, closeClient := exampleContainer()
	pk := azcosmos.NewPartitionKeyString("gear-surf-surfboards")

	item := catalogItem{ID: "item-1", CategoryName: "gear-surf-surfboards", Quantity: 3}

	response, err := container.ReplaceItem(context.TODO(), pk, "item-1", item, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	log.Printf("replaced the item, charged %v RU", response.RequestCharge)

	closeClient()
}

func ExampleContainerClient_UpsertItem() {
	container, closeClient := exampleContainer()
	pk := azcosmos.NewPartitionKeyString("gear-surf-surfboards")

	item := catalogItem{ID: "item-1", CategoryName: "gear-surf-surfboards", Quantity: 3}

	response, err := container.UpsertItem(context.TODO(), pk, "item-1", item, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	log.Printf("upserted the item, charged %v RU", response.RequestCharge)

	closeClient()
}

func ExampleContainerClient_DeleteItem() {
	container, closeClient := exampleContainer()
	pk := azcosmos.NewPartitionKeyString("gear-surf-surfboards")

	response, err := container.DeleteItem(context.TODO(), pk, "item-1", nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	log.Printf("deleted the item, charged %v RU", response.RequestCharge)

	closeClient()
}

func ExampleContainerClient_PatchItem() {
	container, closeClient := exampleContainer()
	pk := azcosmos.NewPartitionKeyString("gear-surf-surfboards")

	var operations azcosmos.PatchOperations
	if err := operations.AppendSet("/quantity", 5); err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	if err := operations.AppendIncrement("/inventoryCount", 1); err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	response, err := container.PatchItem(context.TODO(), pk, "item-1", operations, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	log.Printf("patched the item, charged %v RU", response.RequestCharge)
	if len(response.Value) != 0 {
		item, err := response.IntoModel[catalogItem]()
		if err != nil {
			// TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}
		log.Printf("updated quantity: %d", item.Quantity)
	}

	closeClient()
}

type numberPreservingSerializer struct{}

func (numberPreservingSerializer) Marshal(value any) ([]byte, error) {
	return json.Marshal(value)
}

func (numberPreservingSerializer) Unmarshal(data []byte, value any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(value); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return errors.New("unexpected trailing JSON value")
		}
		return err
	}
	return nil
}

func ExampleClientOptions_serializer() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	client, err := azcosmos.NewClient("https://myaccount.documents.azure.com", cred, &azcosmos.ClientOptions{
		Serializer: numberPreservingSerializer{},
	})
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	container, err := client.NewContainer("myDatabase", "myContainer")
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	pager := container.NewQueryItemsPager[map[string]any](
		azcosmos.NewQuery("SELECT c.quantity FROM c"),
		azcosmos.NewFeedScopeForPartitionKey(azcosmos.NewPartitionKeyString("catalog")), nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			// TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}
		for _, projection := range page.Items {
			log.Printf("quantity without float64 rounding: %v", projection["quantity"])
		}
	}
	if err := pager.Close(); err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	if err := client.Close(); err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
}

func ExampleContainerClient_NewQueryItemsPager_scalar() {
	container, closeClient := exampleContainer()
	pager := container.NewQueryItemsPager[*int](
		azcosmos.NewQuery("SELECT VALUE c.quantity FROM c"),
		azcosmos.NewFeedScopeForPartitionKey(azcosmos.NewPartitionKeyString("catalog")), nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			// TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}
		for _, quantity := range page.Items {
			if quantity != nil {
				log.Printf("quantity: %d", *quantity)
			}
		}
	}
	if err := pager.Close(); err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
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
