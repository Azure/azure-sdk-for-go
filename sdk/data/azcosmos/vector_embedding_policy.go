// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// VectorEmbeddingPolicy represents the vector embedding policy for a container.
// This policy defines the vector embedding configurations that specify how vectors
// are stored and searched within the container.
type VectorEmbeddingPolicy struct {
	// VectorEmbeddings contains the list of vector embedding definitions for the container.
	VectorEmbeddings []VectorEmbedding `json:"vectorEmbeddings"`
}

// VectorEmbedding represents a single vector embedding definition within a container.
type VectorEmbedding struct {
	// Path contains the JSON path to the vector property in the document.
	// Example: "/vector1" or "/embeddings/textVector"
	Path string `json:"path"`

	// DataType specifies the data type of the vector elements.
	// Supported values: "float32" (default), "int8", "uint8"
	DataType VectorDataType `json:"dataType"`

	// DistanceFunction specifies the metric used to compute distance/similarity.
	// Supported values: "cosine", "dotproduct", "euclidean"
	DistanceFunction VectorDistanceFunction `json:"distanceFunction"`

	// Dimensions specifies the dimensionality or length of each vector in the path.
	// All vectors in a path should have the same number of dimensions.
	// Default: 1536
	Dimensions int32 `json:"dimensions"`
}

// VectorDataType represents the supported data types for vector elements.
type VectorDataType string

const (
	// VectorDataTypeFloat32 represents 32-bit floating point numbers (default).
	VectorDataTypeFloat32 VectorDataType = "float32"

	// VectorDataTypeFloat16 represents 16-bit floating point numbers.
	VectorDataTypeFloat16 VectorDataType = "float16"

	// VectorDataTypeInt8 represents 8-bit signed integers.
	VectorDataTypeInt8 VectorDataType = "int8"

	// VectorDataTypeUint8 represents 8-bit unsigned integers.
	VectorDataTypeUint8 VectorDataType = "uint8"
)

// VectorDistanceFunction represents the supported distance functions for vector similarity.
type VectorDistanceFunction string

const (
	// VectorDistanceFunctionCosine uses cosine similarity.
	// Values range from -1 (least similar) to +1 (most similar).
	VectorDistanceFunctionCosine VectorDistanceFunction = "cosine"

	// VectorDistanceFunctionDotProduct uses dot product similarity.
	// Values range from -inf (least similar) to +inf (most similar).
	VectorDistanceFunctionDotProduct VectorDistanceFunction = "dotproduct"

	// VectorDistanceFunctionEuclidean uses Euclidean distance.
	// Values range from 0 (most similar) to +inf (least similar).
	VectorDistanceFunctionEuclidean VectorDistanceFunction = "euclidean"
)
