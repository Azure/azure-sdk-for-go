// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"fmt"
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

// UnmarshalJSON implements the json.Unmarshaler interface
func (pkd *PartitionKeyDefinition) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", pkd, err)
	}

	for key, value := range rawMsg {
		switch key {
		case "kind":
			if err := json.Unmarshal(value, &pkd.Kind); err != nil {
				return fmt.Errorf("unmarshalling type %T: %v", pkd, err)
			}
		case "paths":
			if err := json.Unmarshal(value, &pkd.Paths); err != nil {
				return fmt.Errorf("unmarshalling type %T: %v", pkd, err)
			}
		case "version":
			if err := json.Unmarshal(value, &pkd.Version); err != nil {
				return fmt.Errorf("unmarshalling type %T: %v", pkd, err)
			}
		}
	}

	if pkd.Kind == "" && len(pkd.Paths) == 1 {
		pkd.Kind = PartitionKeyKindHash
	} else if pkd.Kind == "" && len(pkd.Paths) > 1 {
		pkd.Kind = PartitionKeyKindMultiHash
	}

	return nil
}
