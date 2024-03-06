// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// ConflictResolutionPolicy represents a conflict resolution policy for a container.
// For more information see https://docs.microsoft.com/azure/cosmos-db/unique-keys
type ConflictResolutionPolicy struct {
	// Conflict resolution mode. By default, the conflict resolution mode is LastWriteWins.
	Mode ConflictResolutionMode `json:"mode"`
	// The path which is present in each item in the container to be used on LastWriteWins conflict resolution.
	// It must be an integer value.
	ResolutionPath string `json:"conflictResolutionPath,omitempty"`
	// The stored procedure path on Custom conflict.
	// The path should be the full path to the procedure
	ResolutionProcedure string `json:"conflictResolutionProcedure,omitempty"`
}
