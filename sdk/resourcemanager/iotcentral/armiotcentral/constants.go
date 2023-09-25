//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armiotcentral

const (
	moduleName = "armiotcentral"
	moduleVersion = "v1.1.1"
)

// AppSKU - The name of the SKU.
type AppSKU string

const (
	AppSKUST0 AppSKU = "ST0"
	AppSKUST1 AppSKU = "ST1"
	AppSKUST2 AppSKU = "ST2"
)

// PossibleAppSKUValues returns the possible values for the AppSKU const type.
func PossibleAppSKUValues() []AppSKU {
	return []AppSKU{	
		AppSKUST0,
		AppSKUST1,
		AppSKUST2,
	}
}

// AppState - The current state of the application.
type AppState string

const (
	AppStateCreated AppState = "created"
	AppStateSuspended AppState = "suspended"
)

// PossibleAppStateValues returns the possible values for the AppState const type.
func PossibleAppStateValues() []AppState {
	return []AppState{	
		AppStateCreated,
		AppStateSuspended,
	}
}

// SystemAssignedServiceIdentityType - Type of managed service identity (either system assigned, or none).
type SystemAssignedServiceIdentityType string

const (
	SystemAssignedServiceIdentityTypeNone SystemAssignedServiceIdentityType = "None"
	SystemAssignedServiceIdentityTypeSystemAssigned SystemAssignedServiceIdentityType = "SystemAssigned"
)

// PossibleSystemAssignedServiceIdentityTypeValues returns the possible values for the SystemAssignedServiceIdentityType const type.
func PossibleSystemAssignedServiceIdentityTypeValues() []SystemAssignedServiceIdentityType {
	return []SystemAssignedServiceIdentityType{	
		SystemAssignedServiceIdentityTypeNone,
		SystemAssignedServiceIdentityTypeSystemAssigned,
	}
}

