// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs

import "context"

// ProcessorPartitionClient allows you to receive events, similar to a PartitionClient, with
// integration into a checkpoint store for tracking progress.
// NOTE: This type is instantiated using the Processor type, which handles dynamic load balancing.
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

// UpdateCheckpoint updates the checkpoint store.
// The next time a client loads, using the checkpoint store, it will start from the event _after_ the latestEvent.
func (p *ProcessorPartitionClient) UpdateCheckpoint(ctx context.Context, latestEvent *ReceivedEventData) error {
	return p.checkpointStore.UpdateCheckpoint(ctx, Checkpoint{
		CheckpointStoreAddress: CheckpointStoreAddress{
			ConsumerGroup:           p.consumerClientDetails.ConsumerGroup,
			EventHubName:            p.consumerClientDetails.EventHubName,
			FullyQualifiedNamespace: p.consumerClientDetails.FullyQualifiedNamespace,
			PartitionID:             p.partitionID,
		},
		CheckpointData: CheckpointData{
			SequenceNumber: &latestEvent.SequenceNumber,
			Offset:         latestEvent.Offset,
		},
	}, nil)
}

// PartitionID is the partition ID of the partition we're receiving from.
// This will not change.
func (p *ProcessorPartitionClient) PartitionID() string {
	return p.partitionID
}

// Close closes the partition client.
// This does not close the ConsumerClient that the Processor was started with.
func (c *ProcessorPartitionClient) Close(ctx context.Context) error {
	c.cleanupFn()

	if c.innerClient != nil {
		return c.innerClient.Close(ctx)
	}

	return nil
}
