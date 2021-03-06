package consumption

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// Datagrain enumerates the values for datagrain.
type Datagrain string

const (
	// DailyGrain Daily grain of data
	DailyGrain Datagrain = "daily"
	// MonthlyGrain Monthly grain of data
	MonthlyGrain Datagrain = "monthly"
)

// PossibleDatagrainValues returns an array of possible values for the Datagrain const type.
func PossibleDatagrainValues() []Datagrain {
	return []Datagrain{DailyGrain, MonthlyGrain}
}
