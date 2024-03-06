// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// PartitionKeyDefinition represents a partition key definition in the Azure Cosmos DB database service.
// A partition key definition defines the path for the partition key property.
type PartitionKeyDefinition struct {
	// Paths returns the list of partition key paths of the container.
	Paths []string `json:"paths"`
	// Version returns the version of the hash partitioning of the container.
	Version int `json:"version,omitempty"`
}
