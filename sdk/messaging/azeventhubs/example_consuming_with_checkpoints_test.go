// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// Shows how to use the [Processor] type, using a [ConsumerClient] and [CheckpointStore].
//
// The Processor type acts as a load balancer, ensuring that partitions are divided up evenly
// amongst active Processor instances. It also allows storing (and restoring) checkpoints of progress.
//
// NOTE: If you want to manually allocate partitions or to have more control over the process you can use
// the [ConsumerClient]. See [example_consuming_events_test.go] for an example.
//
// [example_consuming_events_test.go]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azeventhubs/example_consuming_events_test.go
func Example_consumingEventsWithCheckpoints() {
	// The Processor makes it simpler to do distributed consumption of an Event Hub.
	// It automatically coordinates with other Processor instances to ensure balanced
	// allocation of partitions and tracks status, durably, in a CheckpointStore.
	//
	// The built-in checkpoint store (available in the `azeventhubs/checkpoints` package) uses
	// Azure Blob storage.

	eventHubNamespace := os.Getenv("EVENTHUB_NAMESPACE")
	eventHubName := os.Getenv("EVENTHUB_NAME")

	storageEndpoint := os.Getenv("CHECKPOINTSTORE_STORAGE_ENDPOINT")
	storageContainerName := os.Getenv("CHECKPOINTSTORE_STORAGE_CONTAINER_NAME")

	if eventHubName == "" || eventHubNamespace == "" || storageEndpoint == "" || storageContainerName == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	consumerClient, checkpointStore, err := createClientsForExample(eventHubNamespace, eventHubName, storageEndpoint, storageContainerName)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	defer func() { _ = consumerClient.Close(context.TODO()) }()

	// Create the Processor
	//
	// The Processor handles load balancing with other Processor instances, running in separate
	// processes or even on separate machines. Each one will use the checkpointStore to coordinate
	// state and ownership, dynamically.
	processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	// Run in the background, launching goroutines to process each partition
	go dispatchPartitionClients(processor)

	// Run the load balancer. The dispatchPartitionClients goroutine (launched above)
	// will receive and dispatch ProcessorPartitionClients as partitions are claimed.
	//
	// Stopping the processor is as simple as canceling the context that you passed
	// in to Run.
	processorCtx, processorCancel := context.WithCancel(context.TODO())
	defer processorCancel()

	if err := processor.Run(processorCtx); err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}
}

func dispatchPartitionClients(processor *azeventhubs.Processor) {
	for {
		processorPartitionClient := processor.NextPartitionClient(context.TODO())

		if processorPartitionClient == nil {
			// Processor has stopped
			break
		}

		go func() {
			if err := processEventsForPartition(processorPartitionClient); err != nil {
				// TODO: Update the following line with your application specific error handling logic
				log.Fatalf("ERROR: %s", err)
			}
		}()
	}
}

// processEventsForPartition shows the typical pattern for processing a partition.
func processEventsForPartition(partitionClient *azeventhubs.ProcessorPartitionClient) error {
	// 1. [BEGIN] Initialize any partition specific resources for your application.
	// 2. [CONTINUOUS] Loop, calling ReceiveEvents() and UpdateCheckpoint().
	// 3. [END] Cleanup any resources.

	defer func() {
		// 3/3 [END] Do cleanup here, like shutting down database clients
		// or other resources used for processing this partition.
		shutdownPartitionResources(partitionClient)
	}()

	// 1/3 [BEGIN] Initialize any partition specific resources for your application.
	if err := initializePartitionResources(partitionClient.PartitionID()); err != nil {
		return err
	}

	// 2/3 [CONTINUOUS] Receive events, checkpointing as needed using UpdateCheckpoint.
	log.Printf("Starting to receive for partition %s", partitionClient.PartitionID())
	for {
		// Wait up to a minute for 100 events, otherwise returns whatever we collected during that time.
		receiveCtx, cancelReceive := context.WithTimeout(context.TODO(), time.Minute)
		events, err := partitionClient.ReceiveEvents(receiveCtx, 100, nil)
		cancelReceive()

		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			var eventHubError *azeventhubs.Error

			if errors.As(err, &eventHubError) && eventHubError.Code == azeventhubs.ErrorCodeOwnershipLost {
				return nil
			}

			return err
		}

		if len(events) == 0 {
			continue
		}

		log.Printf("Received %d event(s)", len(events))

		for _, event := range events {
			log.Printf("Event received with body %v", event.Body)
		}

		// Updates the checkpoint with the latest event received. If processing needs to restart
		// it will restart from this point, automatically.
		if err := partitionClient.UpdateCheckpoint(context.TODO(), events[len(events)-1], nil); err != nil {
			return err
		}
	}
}

func initializePartitionResources(partitionID string) error {
	// initialize things that might be partition specific, like a
	// database connection.
	log.Printf("Initializing partition related resources for partition %s", partitionID)
	return nil
}

func shutdownPartitionResources(partitionClient *azeventhubs.ProcessorPartitionClient) {
	// Each PartitionClient holds onto an external resource and should be closed if you're
	// not processing them anymore.
	defer func() { _ = partitionClient.Close(context.TODO()) }()

	log.Printf("Shutting down partition related resources for partition %s", partitionClient.PartitionID())
}

func createClientsForExample(eventHubNamespace, eventHubName, storageServiceURL, storageContainerName string) (*azeventhubs.ConsumerClient, azeventhubs.CheckpointStore, error) {
	defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		return nil, nil, err
	}

	// NOTE: the storageContainerName must exist before the checkpoint store can be used.
	blobClient, err := azblob.NewClient(storageServiceURL, defaultAzureCred, nil)

	if err != nil {
		return nil, nil, err
	}

	azBlobContainerClient := blobClient.ServiceClient().NewContainerClient(storageContainerName)

	checkpointStore, err := checkpoints.NewBlobStore(azBlobContainerClient, nil)

	if err != nil {
		return nil, nil, err
	}

	consumerClient, err := azeventhubs.NewConsumerClient(eventHubNamespace, eventHubName, azeventhubs.DefaultConsumerGroup, defaultAzureCred, nil)

	if err != nil {
		return nil, nil, err
	}

	return consumerClient, checkpointStore, nil
}
