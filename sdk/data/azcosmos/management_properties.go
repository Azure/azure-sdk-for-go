// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

// DatabaseProperties represents the properties of a database.
type DatabaseProperties struct {
	// ID is the unique id of the database. Required on [Client.CreateDatabase]; server-assigned
	// fields below are populated by the service and ignored on create.
	ID string `json:"id"`
	// ETag is the entity tag of the database (`_etag`).
	ETag *azcore.ETag `json:"_etag,omitempty"`
	// SelfLink is the self-link of the database (`_self`).
	SelfLink string `json:"_self,omitempty"`
	// ResourceID is the resource id of the database (`_rid`).
	ResourceID string `json:"_rid,omitempty"`
}

// PartitionKeyKind represents the type of partition key used by a container.
type PartitionKeyKind string

const (
	// PartitionKeyKindHash selects a single-path hash partition key.
	PartitionKeyKindHash PartitionKeyKind = "Hash"
	// PartitionKeyKindMultiHash selects a hierarchical (multi-path) partition key.
	PartitionKeyKindMultiHash PartitionKeyKind = "MultiHash"
)

// PartitionKeyDefinition defines the partition key path(s) of a container.
type PartitionKeyDefinition struct {
	// Kind is the partition key kind. Empty infers [PartitionKeyKindHash] for a single path and
	// [PartitionKeyKindMultiHash] for more than one, matching the service default.
	Kind PartitionKeyKind `json:"kind,omitempty"`
	// Paths is the list of partition key paths, for example "/tenantId".
	Paths []string `json:"paths"`
	// Version is the partition key hash version. Zero uses the service default.
	Version int `json:"version,omitempty"`
}

// IndexingMode selects how a container's indexing policy is applied.
type IndexingMode string

const (
	// IndexingModeConsistent updates the index synchronously with every write.
	IndexingModeConsistent IndexingMode = "consistent"
	// IndexingModeNone disables indexing.
	IndexingModeNone IndexingMode = "none"
)

// IncludedPath is a JSON path included in a container's indexing policy.
type IncludedPath struct {
	Path string `json:"path"`
}

// ExcludedPath is a JSON path excluded from a container's indexing policy.
type ExcludedPath struct {
	Path string `json:"path"`
}

// CompositeIndexOrder is the sort order of a [CompositeIndex] path.
type CompositeIndexOrder string

const (
	CompositeIndexAscending  CompositeIndexOrder = "ascending"
	CompositeIndexDescending CompositeIndexOrder = "descending"
)

// CompositeIndex is one path within a composite index, used to serve an ORDER BY query over
// multiple properties.
type CompositeIndex struct {
	Path  string              `json:"path"`
	Order CompositeIndexOrder `json:"order"`
}

// IndexingPolicy configures how a container indexes items.
//
// This mirrors the subset of the service's indexing policy schema needed to create a container
// with consistent indexing, included/excluded paths and composite indexes. Spatial, vector and
// full-text indexing are not represented here.
type IndexingPolicy struct {
	// Automatic selects whether the policy applies automatically on write.
	Automatic bool `json:"automatic"`
	// IndexingMode selects the indexing mode. Empty uses the service default
	// ([IndexingModeConsistent]).
	IndexingMode IndexingMode `json:"indexingMode,omitempty"`
	// IncludedPaths lists the paths to index. Empty indexes every path.
	IncludedPaths []IncludedPath `json:"includedPaths,omitempty"`
	// ExcludedPaths lists the paths to exclude from indexing.
	ExcludedPaths []ExcludedPath `json:"excludedPaths,omitempty"`
	// CompositeIndexes lists the composite indexes, each an ordered slice of paths.
	CompositeIndexes [][]CompositeIndex `json:"compositeIndexes,omitempty"`
}

// UniqueKey is a set of paths whose combined values must be unique across the items of a
// container.
type UniqueKey struct {
	Paths []string `json:"paths"`
}

// UniqueKeyPolicy configures a container's unique key constraints. Unique keys can only be set
// when the container is created.
type UniqueKeyPolicy struct {
	UniqueKeys []UniqueKey `json:"uniqueKeys"`
}

// ContainerProperties represents the properties of a container.
type ContainerProperties struct {
	// ID is the unique id of the container. Required on [DatabaseClient.CreateContainer].
	ID string `json:"id"`
	// ETag is the entity tag of the container (`_etag`).
	ETag *azcore.ETag `json:"_etag,omitempty"`
	// SelfLink is the self-link of the container (`_self`).
	SelfLink string `json:"_self,omitempty"`
	// ResourceID is the resource id of the container (`_rid`).
	ResourceID string `json:"_rid,omitempty"`
	// DefaultTimeToLive is the default time to live, in seconds, for items in the container. Nil
	// disables TTL expiration.
	// https://learn.microsoft.com/azure/cosmos-db/time-to-live
	DefaultTimeToLive *int32 `json:"defaultTtl,omitempty"`
	// PartitionKeyDefinition is the partition key path(s) of the container. Required on create and
	// immutable afterward.
	PartitionKeyDefinition PartitionKeyDefinition `json:"partitionKey"`
	// IndexingPolicy configures indexing. Nil uses the service default policy.
	IndexingPolicy *IndexingPolicy `json:"indexingPolicy,omitempty"`
	// UniqueKeyPolicy configures unique key constraints. Can only be set on create.
	UniqueKeyPolicy *UniqueKeyPolicy `json:"uniqueKeyPolicy,omitempty"`
}
