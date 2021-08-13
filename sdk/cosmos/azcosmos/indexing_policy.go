// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// IndexingPolicy represents an indexing policy for a container.
// For more information see https://docs.microsoft.com/azure/cosmos-db/index-policy
type IndexingPolicy struct {
	// Automatic defines if the indexing policy is automatic or manual.
	Automatic bool `json:"automatic"`
	// IndexingMode for the container.
	IndexingMode IndexingMode `json:"indexingMode,omitempty"`
	// Paths to be indexed.
	IncludedPaths []IncludedPath `json:"includedPaths,omitempty"`
	// Paths to be excluded.
	ExcludedPaths []ExcludedPath `json:"excludedPaths,omitempty"`
	// Spatial indexes.
	SpatialIndexes []SpatialIndex `json:"spatialIndexes,omitempty"`
	// Spatial indexes.
	CompositeIndexes [][]CompositeIndex `json:"compositeIndexes,omitempty"`
}

// IncludedPath represents a json path to be included in indexing.
type IncludedPath struct {
	// Path to be included.
	Path string `json:"path"`
}

// ExcludedPath represents a json path to be excluded from indexing.
type ExcludedPath struct {
	// Path to be excluded.
	Path string `json:"path"`
}

// SpatialIndex represents a spatial index.
type SpatialIndex struct {
	// Path for the index.
	Path string `json:"path"`
	// SpatialType of the spatial index.
	SpatialTypes []SpatialType `json:"types"`
}

type CompositeIndex struct {
	// Path for the index.
	Path string `json:"path"`
	// Order represents the order of the composite index.
	// For example if you want to run the query "SELECT * FROM c ORDER BY c.age asc, c.height desc",
	// then you need to make the order for "/age" "ascending" and the order for "/height" "descending".
	Order CompositeIndexOrder `json:"order"`
}
