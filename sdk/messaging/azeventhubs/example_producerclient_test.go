// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
)

var producerClient *azeventhubs.ProducerClient

func ExampleNewProducerClient() {
	// `DefaultAzureCredential` tries several common credential types. For more credential types
	// see this link: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-credential-types.
	defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		panic(err)
	}

	producerClient, err = azeventhubs.NewProducerClient("<ex: myeventhubnamespace.servicebus.windows.net>", "eventhub-name", defaultAzureCred, nil)

	if err != nil {
		panic(err)
	}
}

func ExampleNewProducerClientFromConnectionString() {
	// if the connection string contains an EntityPath
	//
	connectionString := "Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>;EntityPath=<entity path>"
	producerClient, err = azeventhubs.NewProducerClientFromConnectionString(connectionString, "", nil)

	// or

	// if the connection string does not contain an EntityPath
	connectionString = "Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>"
	producerClient, err = azeventhubs.NewProducerClientFromConnectionString(connectionString, "eventhub-name", nil)

	if err != nil {
		panic(err)
	}
}

func ExampleProducerClient_SendEventDataBatch() {
	batch, err := producerClient.NewEventDataBatch(context.TODO(), nil)

	if err != nil {
		panic(err)
	}

	// See ExampleProducerClient_AddEventData for more information.
	err = batch.AddEventData(&azeventhubs.EventData{Body: []byte("hello")}, nil)

	if err != nil {
		panic(err)
	}

	err = producerClient.SendEventDataBatch(context.TODO(), batch, nil)

	if err != nil {
		panic(err)
	}
}

func ExampleEventDataBatch_AddEventData() {
	batch, err := producerClient.NewEventDataBatch(context.TODO(), nil)

	if err != nil {
		panic(err)
	}

	// can be called multiple times with new messages until you
	// receive an azeventhubs.ErrMessageTooLarge
	err = batch.AddEventData(&azeventhubs.EventData{
		Body: []byte("hello"),
	}, nil)

	if errors.Is(err, azeventhubs.ErrEventDataTooLarge) {
		// Message was too large to fit into this batch.
		//
		// At this point you'd usually just send the batch (using ProducerClient.SendEventDataBatch),
		// create a new one, and start filling up the batch again.
		//
		// If this is the _only_ message being added to the batch then it's too big in general, and
		// will need to be split or shrunk to fit.
		panic(err)
	} else if err != nil {
		panic(err)
	}

	err = producerClient.SendEventDataBatch(context.TODO(), batch, nil)

	if err != nil {
		panic(err)
	}
}

func ExampleEventDataBatch_AddEventData_rawAMQPMessages() {
	batch, err := producerClient.NewEventDataBatch(context.TODO(), nil)

	if err != nil {
		panic(err)
	}

	// This is functionally equivalent to EventDataBatch.AddEventData(), just with a more
	// advanced message format.
	// See ExampleEventDataBatch_AddEventData for more details.

	err = batch.AddAMQPAnnotatedMessage(&azeventhubs.AMQPAnnotatedMessage{
		Body: azeventhubs.AMQPAnnotatedMessageBody{
			Data: [][]byte{
				[]byte("hello"),
				[]byte("world"),
			},
		},
	}, nil)

	if err != nil {
		panic(err)
	}

	err = batch.AddAMQPAnnotatedMessage(&azeventhubs.AMQPAnnotatedMessage{
		Body: azeventhubs.AMQPAnnotatedMessageBody{
			Sequence: [][]any{
				// let the AMQP stack encode your strings (or other primitives) for you, no need
				// to convert them to bytes manually.
				{"hello", "world"},
				{"howdy", "world"},
			},
		},
	}, nil)

	if err != nil {
		panic(err)
	}

	err = batch.AddAMQPAnnotatedMessage(&azeventhubs.AMQPAnnotatedMessage{
		Body: azeventhubs.AMQPAnnotatedMessageBody{
			// let the AMQP stack encode your string (or other primitives) for you, no need
			// to convert them to bytes manually.
			Value: "hello world",
		},
	}, nil)

	if err != nil {
		panic(err)
	}

	err = producerClient.SendEventDataBatch(context.TODO(), batch, nil)

	if err != nil {
		panic(err)
	}
}

func ExampleProducerClient_GetEventHubProperties() {
	eventHubProps, err := producerClient.GetEventHubProperties(context.TODO(), nil)

	if err != nil {
		panic(err)
	}

	for _, partitionID := range eventHubProps.PartitionIDs {
		fmt.Printf("Partition ID: %s\n", partitionID)
	}
}

func ExampleProducerClient_GetPartitionProperties() {
	partitionProps, err := producerClient.GetPartitionProperties(context.TODO(), "partition-id", nil)

	if err != nil {
		panic(err)
	}

	fmt.Printf("First sequence number for partition ID %s: %d\n", partitionProps.PartitionID, partitionProps.BeginningSequenceNumber)
	fmt.Printf("Last sequence number for partition ID %s: %d\n", partitionProps.PartitionID, partitionProps.LastEnqueuedSequenceNumber)
}
