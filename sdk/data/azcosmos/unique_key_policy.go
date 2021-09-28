// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// UniqueKeyPolicy represents a unique key policy for a container.
// For more information see https://docs.microsoft.com/azure/cosmos-db/unique-keys
type UniqueKeyPolicy struct {
	// Automatic defines if the indexing policy is automatic or manual.
	UniqueKeys []UniqueKey `json:"uniqueKeys"`
}

// UniqueKey represents a unique key for a container.
// For more information see https://docs.microsoft.com/azure/cosmos-db/unique-keys
type UniqueKey struct {
	// Paths define a sets of paths which must be unique for each document.
	Paths []string `json:"paths"`
}
