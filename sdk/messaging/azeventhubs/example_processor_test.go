// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

// Example_consumingEventsUsingProcessor shows how to use the [Processor] type.
//
// The Processor type acts as a load balancer, ensuring that partitions are divided up amongst
// active Processor instances. You provide it with a [ConsumerClient] as well as a [CheckpointStore].
//
// You will loop, continually calling [Processor.NextPartitionClient] and using the [ProcessorPartitionClient]'s
// that are returned. This loop will run for the lifetime of your application, as ownership can change over
// time as new Processor instances are started, or die.
//
// As you process a partition using [ProcessorPartitionClient.ReceiveEvents] you will periodically
// call [ProcessorPartitionClient.UpdateCheckpoint], which stores your checkpoint information inside of
// the [CheckpointStore]. In the common case, this means your checkpoint information will be stored
// in Azure Blob storage.
//
// If you prefer to manually allocate partitions or to have more control over the process you can use
// the [ConsumerClient] type. See [example_consuming_events_test.go] for an example.
//
// [example_consuming_events_test.go]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azeventhubs/example_consuming_events_test.go
func Example_consumingEventsUsingProcessor() {
	// The Processor makes it simpler to do distributed consumption of an Event Hub.
	// It automatically coordinates with other Processor instances to ensure balanced
	// allocation of partitions and tracks status, durably, in a CheckpointStore.
	//
	// The built-in checkpoint store (available in the `azeventhubs/checkpoints` package) uses
	// Azure Blob storage.

	ehCS := os.Getenv("EVENTHUB_CONNECTION_STRING")
	eventHubName := os.Getenv("EVENTHUB_NAME")

	storageCS := os.Getenv("CHECKPOINTSTORE_STORAGE_CONNECTION_STRING")
	containerName := os.Getenv("CHECKPOINTSTORE_STORAGE_CONTAINER_NAME")

	// Create the checkpoint store
	//
	// NOTE: the Blob container must exist before the checkpoint store can be used.
	azBlobContainerClient, err := container.NewClientFromConnectionString(storageCS, containerName, nil)

	if err != nil {
		panic(err)
	}

	checkpointStore, err := checkpoints.NewBlobStore(azBlobContainerClient, nil)

	if err != nil {
		panic(err)
	}

	// Create a ConsumerClient
	//
	// The Processor (created below) will use this to create any PartitionClient instances, as ownership
	// is assigned.
	consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(ehCS, eventHubName, azeventhubs.DefaultConsumerGroup, nil)

	if err != nil {
		panic(err)
	}

	defer consumerClient.Close(context.TODO())

	// Create the Processor
	//
	// The Processor handles load balancing with other Processor instances, running in separate
	// processes or even on separate machines. Each one will use the checkpointStore to coordinate
	// state and ownership, dynamically.
	processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, nil)

	if err != nil {
		panic(err)
	}

	// This function will be launched in a goroutine and will loop and launch
	// processing of partitions, in parallel.
	dispatchPartitionClients := func() {
		// Our loop will terminate when:
		// - You cancel the passed in context
		// - The Processor.Run() call is cancelled or terminates.
		for {
			partitionClient := processor.NextPartitionClient(context.TODO())

			if partitionClient == nil {
				// this happens if Processor.Run terminates or if you cancel
				// the passed in context.
				break
			}

			// Each time you get a ProcessorPartitionClient, you are the
			// exclusive owner. The previous owner (if there was one) will get an
			// *azeventhubs.Error from their next call to PartitionClient.ReceiveEvents()
			// with a Code of CodeOwnershipLost.
			go func() {
				if err := processEvents(partitionClient); err != nil {
					panic(err)
				}
			}()
		}
	}

	go dispatchPartitionClients()

	// This context will control the lifetime of the Processor.Run call.
	// When it is cancelled the processor will stop running and also close
	// any ProcessorPartitionClient's it opened while running.
	processorCtx, processorCancel := context.WithCancel(context.TODO())
	defer processorCancel()

	// Run the load balancer. The dispatchPartitionClients goroutine, launched
	// above, will continually get new ProcessorPartitionClient's as partitions
	// are allocated.
	//
	// Stopping the processor is as simple as canceling the context that you passed
	// in to Run.
	if err := processor.Run(processorCtx); err != nil {
		panic(err)
	}
}

func processEvents(partitionClient *azeventhubs.ProcessorPartitionClient) error {
	// In other models of the Processor we have broken up the partition
	// lifecycle model.
	//
	// In Go, we model this as a function call, with a loop, using this structure:
	//
	// 1. [BEGIN] Initialize any partition specific resources.
	// 2. [CONTINUOUS] Run a loop, calling ReceiveEvents() and UpdateCheckpoint().
	// 3. [END] Close any resources associated with the processor.

	// [END] Do cleanup here, like shutting down database connections
	// or other resources used for processing this partition.
	defer closePartitionResources(partitionClient)

	// [BEGIN] Initialize any resources needed to process the partition
	if err := initializePartitionResources(partitionClient.PartitionID()); err != nil {
		return err
	}

	// [CONTINUOUS] loop until you lose ownership or your own criteria, checkpointing
	// as needed using UpdateCheckpoint.
	for {
		// Using a context with a timeout will allow ReceiveEvents() to return with events it
		// collected in a minute, or earlier if it actually gets all 100 events we requested.
		receiveCtx, receiveCtxCancel := context.WithTimeout(context.TODO(), time.Minute)
		events, err := partitionClient.ReceiveEvents(receiveCtx, 100, nil)
		receiveCtxCancel()

		// Timing out (context.DeadlineExceeded) is fine. We didn't receive our 100 events
		// but we might have received _some_ events.
		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			if eventHubError := (*azeventhubs.Error)(nil); errors.As(err, &eventHubError) && eventHubError.Code == exported.ErrorCodeOwnershipLost {
				// This means that the partition was "stolen" - this can happen as partitions are balanced between
				// consumers. We'll exit here and just let our "defer closePartitionResources" handle closing
				// resources, including the ProcessorPartitionClient.
				return nil
			}

			return err
		}

		fmt.Printf("Processing %d event(s)\n", len(events))

		for _, event := range events {
			// process the event in some way
			fmt.Printf("Event received with body %v\n", event.Body)
		}

		// it's possible to get zero events if the partition is empty, or if no new events have arrived
		// since your last receive.
		if len(events) != 0 {
			// Update the checkpoint with the last event received. If we lose ownership of this partition or
			// have to restart the next owner will start from this point.
			if err := partitionClient.UpdateCheckpoint(context.TODO(), events[len(events)-1]); err != nil {
				return err
			}
		}
	}
}

func initializePartitionResources(partitionID string) error {
	// initialize things that might be partition specific, like a
	// database connection.
	return nil
}

func closePartitionResources(partitionClient *azeventhubs.ProcessorPartitionClient) {
	// Each PartitionClient holds onto an external resource and should be closed if you're
	// not processing them anymore.
	defer partitionClient.Close(context.TODO())
}
