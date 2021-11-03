// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func ExampleNewClientWithKey() {
	endpoint, ok := os.LookupEnv("AZURE_COSMOS_ENDPOINT")
	if !ok {
		panic("AZURE_COSMOS_ENDPOINT could not be found")
	}

	key, ok := os.LookupEnv("AZURE_COSMOS_KEY")
	if !ok {
		panic("AZURE_COSMOS_KEY could not be found")
	}

	// Create new Cosmos DB client.
	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		panic(err)
	}

	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(client)
}

func ExampleClient_CreateDatabase() {
	endpoint, ok := os.LookupEnv("AZURE_COSMOS_ENDPOINT")
	if !ok {
		panic("AZURE_COSMOS_ENDPOINT could not be found")
	}

	key, ok := os.LookupEnv("AZURE_COSMOS_KEY")
	if !ok {
		panic("AZURE_COSMOS_KEY could not be found")
	}

	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		panic(err)
	}

	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		panic(err)
	}

	databaseProperties := azcosmos.DatabaseProperties{ID: "databaseName"}
	databaseResponse, err := client.CreateDatabase(context.Background(), databaseProperties, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Database created. ActivityId %s", databaseResponse.ActivityID)
}

func ExampleDatabaseClient_CreateContainer() {
	endpoint, ok := os.LookupEnv("AZURE_COSMOS_ENDPOINT")
	if !ok {
		panic("AZURE_COSMOS_ENDPOINT could not be found")
	}

	key, ok := os.LookupEnv("AZURE_COSMOS_KEY")
	if !ok {
		panic("AZURE_COSMOS_KEY could not be found")
	}

	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		panic(err)
	}

	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		panic(err)
	}

	database, err := client.NewDatabase("databaseName")
	if err != nil {
		panic(err)
	}

	properties := azcosmos.ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: []string{"/myPartitionKey"},
		},
	}

	throughput := azcosmos.NewManualThroughputProperties(400)

	resp, err := database.CreateContainer(context.Background(), properties, &azcosmos.CreateContainerOptions{ThroughputProperties: &throughput})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Container created. ActivityId %s", resp.ActivityID)
}

func ExampleContainerClient_ReplaceThroughput() {
	endpoint, ok := os.LookupEnv("AZURE_COSMOS_ENDPOINT")
	if !ok {
		panic("AZURE_COSMOS_ENDPOINT could not be found")
	}

	key, ok := os.LookupEnv("AZURE_COSMOS_KEY")
	if !ok {
		panic("AZURE_COSMOS_KEY could not be found")
	}

	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		panic(err)
	}

	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		panic(err)
	}

	container, err := client.NewContainer("databaseName", "aContainer")
	if err != nil {
		panic(err)
	}

	throughputResponse, err := container.ReadThroughput(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	manualThroughput, hasManual := throughputResponse.ThroughputProperties.ManualThroughput()
	if !hasManual {
		panic("Expected to have manual throughput")
	}
	fmt.Printf("Container is provisioned with %v RU/s", manualThroughput)

	// Replace manual throughput
	newScale := azcosmos.NewManualThroughputProperties(500)
	replaceThroughputResponse, err := container.ReplaceThroughput(context.Background(), newScale, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Throughput updated. ActivityId %s", replaceThroughputResponse.ActivityID)
}

func ExampleContainerClient_Replace() {
	endpoint, ok := os.LookupEnv("AZURE_COSMOS_ENDPOINT")
	if !ok {
		panic("AZURE_COSMOS_ENDPOINT could not be found")
	}

	key, ok := os.LookupEnv("AZURE_COSMOS_KEY")
	if !ok {
		panic("AZURE_COSMOS_KEY could not be found")
	}

	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		panic(err)
	}

	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		panic(err)
	}

	container, err := client.NewContainer("databaseName", "aContainer")
	if err != nil {
		panic(err)
	}

	containerResponse, err := container.Read(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	// Changing the indexing policy
	containerResponse.ContainerProperties.IndexingPolicy = &azcosmos.IndexingPolicy{
		IncludedPaths: []azcosmos.IncludedPath{},
		ExcludedPaths: []azcosmos.ExcludedPath{},
		Automatic:     false,
		IndexingMode:  azcosmos.IndexingModeNone,
	}

	// Replace container properties
	replaceResponse, err := container.Replace(context.Background(), *containerResponse.ContainerProperties, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Container updated. ActivityId %s", replaceResponse.ActivityID)
}

func ExampleContainerClient_CreateItem() {
	endpoint, ok := os.LookupEnv("AZURE_COSMOS_ENDPOINT")
	if !ok {
		panic("AZURE_COSMOS_ENDPOINT could not be found")
	}

	key, ok := os.LookupEnv("AZURE_COSMOS_KEY")
	if !ok {
		panic("AZURE_COSMOS_KEY could not be found")
	}

	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		panic(err)
	}

	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		panic(err)
	}

	container, err := client.NewContainer("databaseName", "aContainer")
	if err != nil {
		panic(err)
	}

	pk := azcosmos.NewPartitionKeyString("newPartitionKey")

	item := map[string]string{
		"id":             "anId",
		"value":          "2",
		"myPartitionKey": "newPartitionKey",
	}

	marshalled, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}

	itemResponse, err := container.CreateItem(context.Background(), pk, marshalled, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Item created. ActivityId %s consuming %v RU", itemResponse.ActivityID, itemResponse.RequestCharge)
}

func ExampleContainerClient_ReadItem() {
	endpoint, ok := os.LookupEnv("AZURE_COSMOS_ENDPOINT")
	if !ok {
		panic("AZURE_COSMOS_ENDPOINT could not be found")
	}

	key, ok := os.LookupEnv("AZURE_COSMOS_KEY")
	if !ok {
		panic("AZURE_COSMOS_KEY could not be found")
	}

	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		panic(err)
	}

	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		panic(err)
	}

	container, err := client.NewContainer("databaseName", "aContainer")
	if err != nil {
		panic(err)
	}

	pk := azcosmos.NewPartitionKeyString("newPartitionKey")

	id := "anId"
	itemResponse, err := container.ReadItem(context.Background(), pk, id, nil)
	if err != nil {
		panic(err)
	}

	var itemResponseBody map[string]string
	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Item read. ActivityId %s consuming %v RU", itemResponse.ActivityID, itemResponse.RequestCharge)
}

func ExampleContainerClient_ReplaceItem() {
	endpoint, ok := os.LookupEnv("AZURE_COSMOS_ENDPOINT")
	if !ok {
		panic("AZURE_COSMOS_ENDPOINT could not be found")
	}

	key, ok := os.LookupEnv("AZURE_COSMOS_KEY")
	if !ok {
		panic("AZURE_COSMOS_KEY could not be found")
	}

	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		panic(err)
	}

	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		panic(err)
	}

	container, err := client.NewContainer("databaseName", "aContainer")
	if err != nil {
		panic(err)
	}

	pk := azcosmos.NewPartitionKeyString("newPartitionKey")

	id := "anId"
	itemResponse, err := container.ReadItem(context.Background(), pk, id, nil)
	if err != nil {
		panic(err)
	}

	var itemResponseBody map[string]string
	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
	if err != nil {
		panic(err)
	}

	// Modify some property
	itemResponseBody["value"] = "newValue"
	marshalledReplace, err := json.Marshal(itemResponseBody)
	if err != nil {
		panic(err)
	}

	itemResponse, err = container.ReplaceItem(context.Background(), pk, id, marshalledReplace, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Item replaced. ActivityId %s consuming %v RU", itemResponse.ActivityID, itemResponse.RequestCharge)
}

func ExampleContainerClient_DeleteItem() {
	endpoint, ok := os.LookupEnv("AZURE_COSMOS_ENDPOINT")
	if !ok {
		panic("AZURE_COSMOS_ENDPOINT could not be found")
	}

	key, ok := os.LookupEnv("AZURE_COSMOS_KEY")
	if !ok {
		panic("AZURE_COSMOS_KEY could not be found")
	}

	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		panic(err)
	}

	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		panic(err)
	}

	container, err := client.NewContainer("databaseName", "aContainer")
	if err != nil {
		panic(err)
	}

	pk := azcosmos.NewPartitionKeyString("newPartitionKey")

	id := "anId"
	itemResponse, err := container.DeleteItem(context.Background(), pk, id, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Item deleted. ActivityId %s consuming %v RU", itemResponse.ActivityID, itemResponse.RequestCharge)
}

func ExampleContainerClient_ReadItem_sessionConsistency() {
	endpoint, ok := os.LookupEnv("AZURE_COSMOS_ENDPOINT")
	if !ok {
		panic("AZURE_COSMOS_ENDPOINT could not be found")
	}

	key, ok := os.LookupEnv("AZURE_COSMOS_KEY")
	if !ok {
		panic("AZURE_COSMOS_KEY could not be found")
	}

	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		panic(err)
	}

	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		panic(err)
	}

	container, err := client.NewContainer("databaseName", "aContainer")
	if err != nil {
		panic(err)
	}

	pk := azcosmos.NewPartitionKeyString("newPartitionKey")
	id := "anId"
	item := map[string]string{
		"id":             "anId",
		"value":          "2",
		"myPartitionKey": "newPartitionKey",
	}

	marshalled, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}

	itemResponse, err := container.CreateItem(context.Background(), pk, marshalled, nil)
	if err != nil {
		panic(err)
	}

	itemSessionToken := itemResponse.SessionToken
	fmt.Printf("Create response contained session %s", itemSessionToken)

	// In another client, maintain the session by passing the session token
	itemResponse, err = container.ReadItem(context.Background(), pk, id, &azcosmos.ItemOptions{SessionToken: itemSessionToken})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Item read. ActivityId %s consuming %v RU", itemResponse.ActivityID, itemResponse.RequestCharge)
}

// Azure Cosmos DB supports optimistic concurrency control to prevent lost updates or deletes and detection of conflicting operations.
// Check the item response status code. If an error is emitted and the response code is 412 then retry operation.
func ExampleContainerClient_ReplaceItem_optimisticConcurrency() {
	endpoint, ok := os.LookupEnv("AZURE_COSMOS_ENDPOINT")
	if !ok {
		panic("AZURE_COSMOS_ENDPOINT could not be found")
	}

	key, ok := os.LookupEnv("AZURE_COSMOS_KEY")
	if !ok {
		panic("AZURE_COSMOS_KEY could not be found")
	}

	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		panic(err)
	}

	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		panic(err)
	}

	container, err := client.NewContainer("databaseName", "aContainer")
	if err != nil {
		panic(err)
	}

	pk := azcosmos.NewPartitionKeyString("newPartitionKey")
	id := "anId"

	numberRetry := 3 // Defining a limit on retries
	err = retryOptimisticConcurrency(numberRetry, 10*time.Millisecond, func() (bool, error) {
		itemResponse, err := container.ReadItem(context.Background(), pk, id, nil)
		if err != nil {
			panic(err)
		}

		var itemResponseBody map[string]string
		err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
		if err != nil {
			panic(err)
		}

		// Change a value in the item response body.
		itemResponseBody["value"] = "newValue"

		marshalledReplace, err := json.Marshal(itemResponseBody)
		if err != nil {
			panic(err)
		}

		// Replace with Etag
		etag := itemResponse.ETag
		itemResponse, err = container.ReplaceItem(context.Background(), pk, id, marshalledReplace, &azcosmos.ItemOptions{IfMatchEtag: &etag})
		var httpErr azcore.HTTPResponse

		return (errors.As(err, &httpErr) && itemResponse.RawResponse.StatusCode == 412), err
	})
	if err != nil {
		panic(err)
	}
}

func retryOptimisticConcurrency(retryAttempts int, wait time.Duration, retry func() (bool, error)) (result error) {
	for i := 0; ; i++ {
		retryResult, err := retry()
		if err != nil {
			break
		}

		if !(retryResult) {
			break
		}

		if i >= (retryAttempts - 1) {
			break
		}

		fmt.Printf("retrying after error: %v", err)

		time.Sleep(wait)
	}
	return fmt.Errorf("Cosmos DB retry attempts %d, error: %s", retryAttempts, result)
}
