// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ContainerProperties represents the properties of a container.
type ContainerProperties struct {
	// ID contains the unique id of the container.
	ID string `json:"id"`
	// ETag contains the entity etag of the container.
	ETag *azcore.ETag `json:"_etag,omitempty"`
	// SelfLink contains the self-link of the container.
	SelfLink string `json:"_self,omitempty"`
	// ResourceID contains the resource id of the container.
	ResourceID string `json:"_rid,omitempty"`
	// LastModified contains the last modified time of the container.
	LastModified int64 `json:"_ts,omitempty"`
	// DefaultTimeToLive contains the default time to live in seconds for items in the container.
	// For more information see https://docs.microsoft.com/azure/cosmos-db/time-to-live#time-to-live-configurations
	DefaultTimeToLive *int32 `json:"defaultTtl,omitempty"`
	// AnalyticalStoreTimeToLiveInSeconds contains the default time to live in seconds for analytical store in the container.
	// For more information see https://docs.microsoft.com/azure/cosmos-db/analytical-store-introduction#analytical-ttl
	AnalyticalStoreTimeToLiveInSeconds *int32 `json:"analyticalStorageTtl,omitempty"`
	// PartitionKeyDefinition contains the partition key definition of the container.
	PartitionKeyDefinition PartitionKeyDefinition `json:"partitionKey,omitempty"`
	// IndexingPolicy contains the indexing definition of the container.
	IndexingPolicy *IndexingPolicy `json:"indexingPolicy,omitempty"`
	// UniqueKeyPolicy contains the unique key policy of the container.
	UniqueKeyPolicy *UniqueKeyPolicy `json:"uniqueKeyPolicy,omitempty"`
	// ConflictResolutionPolicy contains the conflict resolution policy of the container.
	ConflictResolutionPolicy *ConflictResolutionPolicy `json:"conflictResolutionPolicy,omitempty"`
}
