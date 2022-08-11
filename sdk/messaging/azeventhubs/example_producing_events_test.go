// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs_test

import (
	"context"
	"errors"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
)

func Example_producing() {
	eventHubNamespace := os.Getenv("EVENTHUB_NAMESPACE") // <ex: myeventhubnamespace.servicebus.windows.net>
	eventHubName := os.Getenv("EVENTHUB_NAME")

	defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		panic(err)
	}

	// Can also use a connection string:
	//
	// producerClient, err = azeventhubs.NewProducerClientFromConnectionString(connectionString, eventHubName, nil)
	//
	producerClient, err = azeventhubs.NewProducerClient(eventHubNamespace, eventHubName, defaultAzureCred, nil)

	if err != nil {
		panic(err)
	}

	defer producerClient.Close(context.TODO())

	// Other examples:
	//
	// sending a batch to a specific partition:
	// batch, err := producerClient.NewEventDataBatch(context.TODO(), &azeventhubs.NewEventDataBatchOptions{
	// 	PartitionID: to.Ptr("0"),
	// })
	//
	// targeting a batch using a partition key
	// batch, err := producerClient.NewEventDataBatch(context.TODO(), &azeventhubs.NewEventDataBatchOptions{
	// 	PartitionKey: to.Ptr("partition key"),
	// })
	batch, err := producerClient.NewEventDataBatch(context.TODO(), nil)

	if err != nil {
		panic(err)
	}

	eventData := &azeventhubs.EventData{
		Body: []byte("hello"),
	}

	err = batch.AddEventData(eventData, nil)

	if errors.Is(err, azeventhubs.ErrEventDataTooLarge) {
		// EventData is too large for this batch.
		//
		// If the batch is empty and this happens then the event will never be sendable at it's current
		// size as it exceeds what the link allows.
		//
		// Otherwise, it's simplest to send this batch and create a new one, starting with this event.
		panic(err)
	} else if err != nil {
		panic(err)
	}

	err = producerClient.SendEventBatch(context.TODO(), batch, nil)

	if err != nil {
		panic(err)
	}
}
