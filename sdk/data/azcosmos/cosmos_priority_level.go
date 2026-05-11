// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// PriorityLevel defines the priority level for a request.
type PriorityLevel string

const (
	PriorityLevelHigh PriorityLevel = "High"
	PriorityLevelLow  PriorityLevel = "Low"
)

// PriorityLevelValues returns a list of available priority levels.
func PriorityLevelValues() []PriorityLevel {
	return []PriorityLevel{PriorityLevelHigh, PriorityLevelLow}
}

// ToPtr returns a *PriorityLevel.
func (p PriorityLevel) ToPtr() *PriorityLevel {
	return &p
}
