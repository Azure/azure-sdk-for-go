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
	// Composite indexes.
	CompositeIndexes [][]CompositeIndex `json:"compositeIndexes,omitempty"`
	// Vector indexes for vector search capabilities.
	VectorIndexes []VectorIndex `json:"vectorIndexes,omitempty"`
	// Full text indexes for full-text search capabilities.
	FullTextIndexes []FullTextIndex `json:"fullTextIndexes,omitempty"`
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

// CompositeIndex is used when queries have an ORDER BY clause with two or more properties
type CompositeIndex struct {
	// Path for the index.
	Path string `json:"path"`
	// Order represents the order of the composite index.
	// For example if you want to run the query "SELECT * FROM c ORDER BY c.age asc, c.height desc",
	// then you need to make the order for "/age" "ascending" and the order for "/height" "descending".
	Order CompositeIndexOrder `json:"order"`
}

// VectorIndex represents a vector index for efficient vector search operations.
type VectorIndex struct {
	// Path to the vector property in the document.
	Path string `json:"path"`
	// Type of vector index algorithm to use.
	Type VectorIndexType `json:"type"`
}

// VectorIndexType represents the supported vector index algorithms in Azure Cosmos DB.
type VectorIndexType string

const (
	// VectorIndexTypeFlat uses a flat (brute-force) index that provides 100% accuracy.
	// Suitable for smaller datasets and has a limitation of 505 dimensions.
	VectorIndexTypeFlat VectorIndexType = "flat"

	// VectorIndexTypeQuantizedFlat uses a quantized flat index that compresses vectors
	// before storing on the index. Provides high accuracy with better performance than flat.
	// Supports up to 4,096 dimensions and is recommended for up to ~50,000 vectors per partition.
	VectorIndexTypeQuantizedFlat VectorIndexType = "quantizedFlat"

	// VectorIndexTypeDiskANN uses DiskANN algorithm for high-performance vector search.
	// Provides the best performance for large datasets with more than 50,000 vectors per partition.
	// Supports up to 4,096 dimensions.
	VectorIndexTypeDiskANN VectorIndexType = "diskANN"
)

// FullTextIndex represents a full-text index for efficient text search operations.
type FullTextIndex struct {
	// Path to the text property in the document.
	Path string `json:"path"`
}
