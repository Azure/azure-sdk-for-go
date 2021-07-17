// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// Consistency levels supported by the Azure Cosmos DB service.
type ConsistencyLevel int

const (
	ConsistencyLevelStrong           ConsistencyLevel = iota
	ConsistencyLevelBoundedStaleness ConsistencyLevel = iota
	ConsistencyLevelSession          ConsistencyLevel = iota
	ConsistencyLevelEventual         ConsistencyLevel = iota
	ConsistencyLevelConsistentPrefix ConsistencyLevel = iota
)

// Returns a list of available consistency levels
func ConsistencyLevelValues() []ConsistencyLevel {
	return []ConsistencyLevel{ConsistencyLevelStrong, ConsistencyLevelBoundedStaleness, ConsistencyLevelSession, ConsistencyLevelEventual, ConsistencyLevelConsistentPrefix}
}

func (c ConsistencyLevel) ToPtr() *ConsistencyLevel {
	return &c
}
