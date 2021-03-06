package storageimportexport

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// DriveState enumerates the values for drive state.
type DriveState string

const (
	// Completed ...
	Completed DriveState = "Completed"
	// CompletedMoreInfo ...
	CompletedMoreInfo DriveState = "CompletedMoreInfo"
	// NeverReceived ...
	NeverReceived DriveState = "NeverReceived"
	// Received ...
	Received DriveState = "Received"
	// ShippedBack ...
	ShippedBack DriveState = "ShippedBack"
	// Specified ...
	Specified DriveState = "Specified"
	// Transferring ...
	Transferring DriveState = "Transferring"
)

// PossibleDriveStateValues returns an array of possible values for the DriveState const type.
func PossibleDriveStateValues() []DriveState {
	return []DriveState{Completed, CompletedMoreInfo, NeverReceived, Received, ShippedBack, Specified, Transferring}
}

// KekType enumerates the values for kek type.
type KekType string

const (
	// CustomerManaged ...
	CustomerManaged KekType = "CustomerManaged"
	// MicrosoftManaged ...
	MicrosoftManaged KekType = "MicrosoftManaged"
)

// PossibleKekTypeValues returns an array of possible values for the KekType const type.
func PossibleKekTypeValues() []KekType {
	return []KekType{CustomerManaged, MicrosoftManaged}
}

// Type enumerates the values for type.
type Type string

const (
	// None ...
	None Type = "None"
	// SystemAssigned ...
	SystemAssigned Type = "SystemAssigned"
	// UserAssigned ...
	UserAssigned Type = "UserAssigned"
)

// PossibleTypeValues returns an array of possible values for the Type const type.
func PossibleTypeValues() []Type {
	return []Type{None, SystemAssigned, UserAssigned}
}
