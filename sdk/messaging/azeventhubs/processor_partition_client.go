// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs

import "context"

// ProcessorPartitionClient allows you to receive events, similar to a PartitionClient, with
// integration into a checkpoint store for tracking progress.
//
// This type is instantiated from the Processor type, which handles dynamic load balancing.
//
// See [example_processor_test.go] for an example of typical usage.
//
// NOTE: If you do NOT want to use dynamic load balancing, and would prefer to track state and ownership
// manually, use the [ConsumerClient] type instead.
//
// [example_processor_test.go]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azeventhubs/example_processor_test.go
type ProcessorPartitionClient struct {
	partitionID           string
	innerClient           *PartitionClient // *azeventhubs.PartitionClient
	checkpointStore       CheckpointStore
	cleanupFn             func()
	consumerClientDetails consumerClientDetails
}

// ReceiveEvents receives events until 'count' events have been received or the context has
// expired or been cancelled.
func (c *ProcessorPartitionClient) ReceiveEvents(ctx context.Context, count int, options *ReceiveEventsOptions) ([]*ReceivedEventData, error) {
	return c.innerClient.ReceiveEvents(ctx, count, options)
}

// UpdateCheckpoint updates the checkpoint store. This ensure that if the Processor is restarted it will
// start from after this point.
func (p *ProcessorPartitionClient) UpdateCheckpoint(ctx context.Context, latestEvent *ReceivedEventData) error {
	return p.checkpointStore.UpdateCheckpoint(ctx, Checkpoint{
		ConsumerGroup:           p.consumerClientDetails.ConsumerGroup,
		EventHubName:            p.consumerClientDetails.EventHubName,
		FullyQualifiedNamespace: p.consumerClientDetails.FullyQualifiedNamespace,
		PartitionID:             p.partitionID,
		SequenceNumber:          &latestEvent.SequenceNumber,
		Offset:                  latestEvent.Offset,
	}, nil)
}

// PartitionID is the partition ID of the partition we're receiving from.
// This will not change during the lifetime of this ProcessorPartitionClient.
func (p *ProcessorPartitionClient) PartitionID() string {
	return p.partitionID
}

// Close releases resources for the partition client.
// This does not close the ConsumerClient that the Processor was started with.
func (c *ProcessorPartitionClient) Close(ctx context.Context) error {
	c.cleanupFn()

	if c.innerClient != nil {
		return c.innerClient.Close(ctx)
	}

	return nil
}
