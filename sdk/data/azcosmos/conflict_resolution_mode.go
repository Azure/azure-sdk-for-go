// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// ConflictResolutionMode defines the conflict resolution mode in the Azure Cosmos DB service.
type ConflictResolutionMode string

const (
	// Conflict resolution that uses the highest value of the conflicting documents property values.
	ConflictResolutionModeLastWriteWins ConflictResolutionMode = "LastWriterWins"
	// Custom conflict resolution mode that requires the definition of a stored procedure.
	ConflictResolutionModeCustom ConflictResolutionMode = "Custom"
)

// Returns a list of available consistency levels
func ConflictResolutionModeValues() []ConflictResolutionMode {
	return []ConflictResolutionMode{ConflictResolutionModeLastWriteWins, ConflictResolutionModeCustom}
}

func (c ConflictResolutionMode) ToPtr() *ConflictResolutionMode {
	return &c
}
