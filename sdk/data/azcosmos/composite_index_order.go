// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// These are the ordering values available for composite indexes in the Azure Cosmos DB database service.
// For more information see https://docs.microsoft.com/azure/cosmos-db/index-policy
type CompositeIndexOrder string

const (
	// Ascending sort order for composite paths.
	CompositeIndexAscending CompositeIndexOrder = "ascending"
	// Descending sort order for composite paths.
	CompositeIndexDescending CompositeIndexOrder = "descending"
)

// Returns a list of available consistency levels
func CompositeIndexOrderValues() []CompositeIndexOrder {
	return []CompositeIndexOrder{CompositeIndexAscending, CompositeIndexDescending}
}

func (c CompositeIndexOrder) ToPtr() *CompositeIndexOrder {
	return &c
}
