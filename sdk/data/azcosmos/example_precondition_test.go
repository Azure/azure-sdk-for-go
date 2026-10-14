// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
)

func ExampleIfMatch() {
	container, closeClient := exampleContainer()
	// Use the ETag from a prior item response.
	condition, err := azcosmos.IfMatch(azcore.ETag(`"previous-etag"`))
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	_, err = container.ReplaceItem(context.TODO(), azcosmos.NewPartitionKeyString("partition"), "item",
		[]byte(`{"id":"item","pk":"partition","quantity":2}`),
		&azcosmos.ReplaceItemOptions{Precondition: condition})
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	closeClient()
}

func ExampleIfNoneMatch() {
	container, closeClient := exampleContainer()
	condition, err := azcosmos.IfNoneMatch(azcore.ETag(`"cached-etag"`))
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	response, err := container.ReadItem(context.TODO(), azcosmos.NewPartitionKeyString("partition"), "item",
		&azcosmos.ReadItemOptions{Precondition: condition})
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	log.Printf("conditional read returned %d bytes", len(response.Value))
	closeClient()
}
