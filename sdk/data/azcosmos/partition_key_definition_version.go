// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// Version of the hash partitioning
type PartitionKeyDefinitionVersion int

const (
	// Original version of hash partitioning.
	PartitionKeyDefinitionVersion1 PartitionKeyDefinitionVersion = 1
	// Enhanced version of hash partitioning - offers better distribution of long partition keys and uses less storage.
	PartitionKeyDefinitionVersion2 PartitionKeyDefinitionVersion = 2
)

// Returns a list of available consistency levels
func PartitionKeyDefinitionVersionValues() []PartitionKeyDefinitionVersion {
	return []PartitionKeyDefinitionVersion{PartitionKeyDefinitionVersion1, PartitionKeyDefinitionVersion2}
}

func (c PartitionKeyDefinitionVersion) ToPtr() *PartitionKeyDefinitionVersion {
	return &c
}
