package windowsiot

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// ServiceNameUnavailabilityReason enumerates the values for service name unavailability reason.
type ServiceNameUnavailabilityReason string

const (
	// AlreadyExists ...
	AlreadyExists ServiceNameUnavailabilityReason = "AlreadyExists"
	// Invalid ...
	Invalid ServiceNameUnavailabilityReason = "Invalid"
)

// PossibleServiceNameUnavailabilityReasonValues returns an array of possible values for the ServiceNameUnavailabilityReason const type.
func PossibleServiceNameUnavailabilityReasonValues() []ServiceNameUnavailabilityReason {
	return []ServiceNameUnavailabilityReason{AlreadyExists, Invalid}
}
