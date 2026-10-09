// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos_test

import (
	"context"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
)

func ExampleContainerClient_ExecuteTransactionalBatch() {
	container, closeClient := exampleContainer()
	// The container in this example has the hierarchical paths /tenant and /region.
	key := azcosmos.NewPartitionKeyString("Contoso").AppendString("west")
	batch := azcosmos.NewTransactionalBatch(key)
	err := batch.CreateItem([]byte(`{"id":"order-1","tenant":"Contoso","region":"west","total":42}`), nil)
	if err != nil {
		closeClient()
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	if err := batch.ReadItem("order-1", nil); err != nil {
		closeClient()
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	response, err := container.ExecuteTransactionalBatch(context.TODO(), batch, nil)
	if err != nil {
		closeClient()
		// An execution error is not proof that writes did not commit.
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	log.Printf("batch HTTP %d, charge %.2f RU", response.StatusCode, response.RequestCharge)
	for i, result := range response.OperationResults {
		log.Printf("operation %d: HTTP %d", i, result.StatusCode)
	}
	log.Printf("read result: %s", response.OperationResults[1].ResourceBody)
	closeClient()
}

func ExampleTransactionalBatch_ReadItem() {
	container, closeClient := exampleContainer()
	key := azcosmos.NewPartitionKeyString("Contoso").AppendString("west")
	cachedETag := azcore.ETag(`"cached-etag"`)
	batch := azcosmos.NewTransactionalBatch(key)
	if err := batch.ReadItem("order-1", &azcosmos.TransactionalBatchReadItemOptions{IfNoneMatchETag: &cachedETag}); err != nil {
		closeClient()
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	response, err := container.ExecuteTransactionalBatch(context.TODO(), batch, nil)
	if err != nil {
		closeClient()
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	result := response.OperationResults[0]
	switch result.StatusCode {
	case http.StatusNotModified:
		log.Print("the cached item is unchanged")
	case http.StatusOK:
		log.Printf("updated item: %s, ETag: %s", result.ResourceBody, result.ETag)
	default:
		log.Printf("read HTTP %d", result.StatusCode)
	}
	closeClient()
}
