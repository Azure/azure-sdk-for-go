// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos_test

import (
	"context"
	"log"

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
	if !response.Success {
		if index, found := response.FailedOperationIndex(); found {
			log.Printf("operation %d failed with HTTP %d", index, response.OperationResults[index].StatusCode)
		}
		closeClient()
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: the batch did not commit")
	}
	log.Printf("committed %d operations, charge %.2f RU", len(response.OperationResults), response.RequestCharge)
	log.Printf("read result: %s", response.OperationResults[1].ResourceBody)
	closeClient()
}
