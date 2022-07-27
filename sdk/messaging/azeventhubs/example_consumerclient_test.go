// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs_test

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
)

var consumerClient *azeventhubs.ConsumerClient
var consumerGroup string
var err error

func ExampleNewConsumerClient() {
	defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		panic(err)
	}

	consumerClient, err = azeventhubs.NewConsumerClient(consumerGroup, "<ex: myeventhubnamespace.servicebus.windows.net>", "eventhub-name", "partition id", defaultAzureCred, nil)

	if err != nil {
		panic(err)
	}
}

func ExampleNewConsumerClientFromConnectionString() {
	connectionString := "Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>;EntityPath=<entity path>"

	consumerClient, err = azeventhubs.NewConsumerClientFromConnectionString(consumerGroup, connectionString, "partition id", nil)

	if err != nil {
		panic(err)
	}
}

func ExampleNewConsumerClientForHubFromConnectionString() {
	connectionString := "Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>"

	consumerClient, err = azeventhubs.NewConsumerClientForHubFromConnectionString(consumerGroup, connectionString, "eventhub-name", "partition id", nil)

	if err != nil {
		panic(err)
	}
}

func ExampleConsumerClient_ReceiveEvents() {
	events, err := consumerClient.ReceiveEvents(context.TODO(), 100, nil)

	if err != nil {
		panic(err)
	}

	for _, evt := range events {
		fmt.Printf("Body: %s\n", string(evt.Body))
	}
}

func ExampleConsumerClient_GetEventHubProperties() {
	eventHubProps, err := consumerClient.GetEventHubProperties(context.TODO(), nil)

	if err != nil {
		panic(err)
	}

	for _, partitionID := range eventHubProps.PartitionIDs {
		fmt.Printf("Partition ID: %s\n", partitionID)
	}
}

func ExampleConsumerClient_GetPartitionProperties() {
	partitionProps, err := consumerClient.GetPartitionProperties(context.TODO(), "partition-id", nil)

	if err != nil {
		panic(err)
	}

	fmt.Printf("First sequence number for partition ID %s: %d\n", partitionProps.PartitionID, partitionProps.BeginningSequenceNumber)
	fmt.Printf("Last sequence number for partition ID %s: %d\n", partitionProps.PartitionID, partitionProps.LastEnqueuedSequenceNumber)
}
