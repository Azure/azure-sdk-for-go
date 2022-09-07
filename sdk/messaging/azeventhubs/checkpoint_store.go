// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs

import (
	"context"
	"time"
)

// CheckpointStore is used by multiple consumers to coordinate progress and ownership for partitions.
type CheckpointStore interface {
	// ClaimOwnership attempts to claim ownership of the partitions in partitionOwnership and returns
	// the actual partitions that were claimed.
	ClaimOwnership(ctx context.Context, partitionOwnership []Ownership, options *ClaimOwnershipOptions) ([]Ownership, error)

	// ListCheckpoints lists all the available checkpoints.
	ListCheckpoints(ctx context.Context, fullyQualifiedNamespace string, eventHubName string, consumerGroup string, options *ListCheckpointsOptions) ([]Checkpoint, error)

	// ListOwnership lists all ownerships.
	ListOwnership(ctx context.Context, fullyQualifiedNamespace string, eventHubName string, consumerGroup string, options *ListOwnershipOptions) ([]Ownership, error)

	// UpdateCheckpoint updates a specific checkpoint with a sequence and offset.
	UpdateCheckpoint(ctx context.Context, checkpoint Checkpoint, options *UpdateCheckpointOptions) error
}

// Ownership tracks which consumer owns a particular partition.
type Ownership struct {
	CheckpointStoreAddress
	OwnershipData
}

// OwnershipData is data specific to ownership, like the
// current owner, by ID, and the last time the ownership was
// updated, which is used to calculate if ownership has expired.
type OwnershipData struct {
	OwnerID          string
	LastModifiedTime time.Time
	ETag             string
}

// Checkpoint tracks the last succesfully processed event in a partition.
type Checkpoint struct {
	CheckpointStoreAddress
	CheckpointData
}

// CheckpointStoreAddress contains the properties needed to uniquely address checkpoints
// or ownership.
type CheckpointStoreAddress struct {
	ConsumerGroup           string
	EventHubName            string
	FullyQualifiedNamespace string
	PartitionID             string
}

// CheckpointData tracks latest offset and sequence number that have been
// processed by the client.
type CheckpointData struct {
	Offset         *int64
	SequenceNumber *int64
}

// ListCheckpointsOptions contains optional parameters for the ListCheckpoints function
type ListCheckpointsOptions struct {
	// For future expansion
}

// ListOwnershipOptions contains optional parameters for the ListOwnership function
type ListOwnershipOptions struct {
	// For future expansion
}

// UpdateCheckpointOptions contains optional parameters for the UpdateCheckpoint function
type UpdateCheckpointOptions struct {
	// For future expansion
}

// ClaimOwnershipOptions contains optional parameters for the ClaimOwnership function
type ClaimOwnershipOptions struct {
	// For future expansion
}
