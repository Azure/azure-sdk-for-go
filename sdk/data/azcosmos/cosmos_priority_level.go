// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// PriorityLevel defines the priority level for a Cosmos DB request.
// When the total number of RU/s consumed exceeds the provisioned capacity,
// low-priority requests are throttled before high-priority requests.
// Valid values are PriorityLevelHigh and PriorityLevelLow.
// For more information, see https://aka.ms/CosmosDB/PriorityBasedExecution
type PriorityLevel string

const (
	// PriorityLevelHigh is the default priority level. High-priority requests are served before low-priority requests.
	PriorityLevelHigh PriorityLevel = "High"
	// PriorityLevelLow marks a request as low priority. These requests are throttled first when over RU/s budget.
	PriorityLevelLow PriorityLevel = "Low"
)

// PriorityLevelValues returns a list of available priority levels.
func PriorityLevelValues() []PriorityLevel {
	return []PriorityLevel{PriorityLevelHigh, PriorityLevelLow}
}

// ToPtr returns a *PriorityLevel.
func (p PriorityLevel) ToPtr() *PriorityLevel {
	return &p
}
