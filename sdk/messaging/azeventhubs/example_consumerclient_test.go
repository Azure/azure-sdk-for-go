// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
)

var consumerClient *azeventhubs.ConsumerClient
var err error

func ExampleNewConsumerClient() {
	// `DefaultAzureCredential` tries several common credential types. For more credential types
	// see this link: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-credential-types.
	defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		panic(err)
	}

	consumerClient, err = azeventhubs.NewConsumerClient("<ex: myeventhubnamespace.servicebus.windows.net>", "eventhub-name", azeventhubs.DefaultConsumerGroup, defaultAzureCred, nil)

	if err != nil {
		panic(err)
	}
}

func ExampleNewConsumerClientFromConnectionString() {
	// if the connection string contains an EntityPath
	//
	connectionString := "Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>;EntityPath=<entity path>"
	consumerClient, err = azeventhubs.NewConsumerClientFromConnectionString(connectionString, "", azeventhubs.DefaultConsumerGroup, nil)

	// or

	// if the connection string does not contain an EntityPath
	connectionString = "Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>"
	consumerClient, err = azeventhubs.NewConsumerClientFromConnectionString(connectionString, "eventhub-name", azeventhubs.DefaultConsumerGroup, nil)

	if err != nil {
		panic(err)
	}
}

func ExampleConsumerClient_NewPartitionClient_receiveEvents() {
	const partitionID = "0"

	partitionClient, err := consumerClient.NewPartitionClient(partitionID, nil)

	if err != nil {
		panic(err)
	}

	defer partitionClient.Close(context.TODO())

	// Using a context with a timeout will allow ReceiveEvents() to return with events it
	// collected in a minute, or earlier if it actually gets all 100 events we requested.
	receiveCtx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()
	events, err := partitionClient.ReceiveEvents(receiveCtx, 100, nil)

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

func ExampleConsumerClient_NewPartitionClient_configuringPrefetch() {
	const partitionID = "0"

	// Prefetching configures the Event Hubs client to continually cache events, up to the configured size
	// in PartitionClientOptions.Prefetch. PartitionClient.ReceiveEvents will read from the cache first,
	// which can improve throughput in situations where you might normally be forced to request and wait
	// for more events.

	// By default, prefetch is enabled.
	partitionClient, err := consumerClient.NewPartitionClient(partitionID, nil)

	if err != nil {
		panic(err)
	}

	defer partitionClient.Close(context.TODO())

	// You can configure the prefetch buffer size as well. The default is 300.
	partitionClientWithCustomPrefetch, err := consumerClient.NewPartitionClient(partitionID, &azeventhubs.PartitionClientOptions{
		Prefetch: 301,
	})

	if err != nil {
		panic(err)
	}

	defer partitionClientWithCustomPrefetch.Close(context.TODO())

	// And prefetch can be disabled if you prefer to manually control the flow of events. Excess
	// events (that arrive after your ReceiveEvents() call has completed) will still be
	// buffered internally, but they will not be automatically replenished.
	partitionClientWithPrefetchDisabled, err := consumerClient.NewPartitionClient(partitionID, &azeventhubs.PartitionClientOptions{
		Prefetch: -1,
	})

	if err != nil {
		panic(err)
	}

	defer partitionClientWithPrefetchDisabled.Close(context.TODO())

	// Using a context with a timeout will allow ReceiveEvents() to return with events it
	// collected in a minute, or earlier if it actually gets all 100 events we requested.
	receiveCtx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()
	events, err := partitionClient.ReceiveEvents(receiveCtx, 100, nil)

	if err != nil {
		panic(err)
	}

	for _, evt := range events {
		fmt.Printf("Body: %s\n", string(evt.Body))
	}
}

func ExampleNewConsumerClient_usingCustomEndpoint() {
	// `DefaultAzureCredential` tries several common credential types. For more credential types
	// see this link: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-credential-types.
	defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		panic(err)
	}

	consumerClient, err = azeventhubs.NewConsumerClient("<ex: myeventhubnamespace.servicebus.windows.net>", "eventhub-name", azeventhubs.DefaultConsumerGroup, defaultAzureCred, &azeventhubs.ConsumerClientOptions{
		// A custom endpoint can be used when you need to connect to a TCP proxy.
		CustomEndpoint: "<address/hostname of TCP proxy>",
	})

	if err != nil {
		panic(err)
	}
}

func ExampleNewConsumerClient_configuringRetries() {
	// `DefaultAzureCredential` tries several common credential types. For more credential types
	// see this link: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-credential-types.
	defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		panic(err)
	}

	consumerClient, err = azeventhubs.NewConsumerClient("<ex: myeventhubnamespace.servicebus.windows.net>", "eventhub-name", azeventhubs.DefaultConsumerGroup, defaultAzureCred, &azeventhubs.ConsumerClientOptions{
		RetryOptions: azeventhubs.RetryOptions{
			// NOTE: these are the default values.
			MaxRetries:    3,
			RetryDelay:    time.Second,
			MaxRetryDelay: 120 * time.Second,
		},
	})

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		fmt.Printf("ERROR: %s\n", err)
		return
	}
}
