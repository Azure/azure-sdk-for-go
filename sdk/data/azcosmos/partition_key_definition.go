// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
)

// PartitionKeyKind represents the type of the partition key that is used in an Azure Cosmos DB container.
type PartitionKeyKind string

const (
	PartitionKeyKindHash      PartitionKeyKind = "Hash"
	PartitionKeyKindMultiHash PartitionKeyKind = "MultiHash"
)

// PartitionKeyDefinition represents a partition key definition in the Azure Cosmos DB database service.
// A partition key definition defines the path for the partition key property.
type PartitionKeyDefinition struct {
	// Kind returns the kind of partition key definition.
	Kind PartitionKeyKind `json:"kind"`
	// Paths returns the list of partition key paths of the container.
	Paths []string `json:"paths"`
	// Version returns the version of the hash partitioning of the container.
	Version int `json:"version,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface
// If the Kind is not set, it will be inferred based on the number of paths.
func (pkd PartitionKeyDefinition) MarshalJSON() ([]byte, error) {
	var paths_length = len(pkd.Paths)

	var kind PartitionKeyKind
	if pkd.Kind != "" {
		kind = pkd.Kind
	} else if pkd.Kind == "" && paths_length == 1 {
		kind = PartitionKeyKindHash
	} else if pkd.Kind == "" && paths_length > 1 {
		kind = PartitionKeyKindMultiHash
	}

	return json.Marshal(struct {
		Kind    PartitionKeyKind `json:"kind"`
		Paths   []string         `json:"paths"`
		Version int              `json:"version,omitempty"`
	}{
		Kind:    kind,
		Paths:   pkd.Paths,
		Version: pkd.Version,
	})
}
