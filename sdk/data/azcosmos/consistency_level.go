// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// Consistency levels supported by the Azure Cosmos DB service.
type ConsistencyLevel string

const (
	ConsistencyLevelStrong           ConsistencyLevel = "Strong"
	ConsistencyLevelBoundedStaleness ConsistencyLevel = "BoundedStaleness"
	ConsistencyLevelSession          ConsistencyLevel = "Session"
	ConsistencyLevelEventual         ConsistencyLevel = "Eventual"
	ConsistencyLevelConsistentPrefix ConsistencyLevel = "ConsistentPrefix"
)

// Returns a list of available consistency levels
func ConsistencyLevelValues() []ConsistencyLevel {
	return []ConsistencyLevel{ConsistencyLevelStrong, ConsistencyLevelBoundedStaleness, ConsistencyLevelSession, ConsistencyLevelEventual, ConsistencyLevelConsistentPrefix}
}

func (c ConsistencyLevel) ToPtr() *ConsistencyLevel {
	return &c
}
