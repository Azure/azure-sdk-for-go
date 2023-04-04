//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armwindowsiot

const (
	moduleName    = "armwindowsiot"
	moduleVersion = "v1.1.0"
)

// ServiceNameUnavailabilityReason - The reason for unavailability.
type ServiceNameUnavailabilityReason string

const (
	ServiceNameUnavailabilityReasonInvalid       ServiceNameUnavailabilityReason = "Invalid"
	ServiceNameUnavailabilityReasonAlreadyExists ServiceNameUnavailabilityReason = "AlreadyExists"
)

// PossibleServiceNameUnavailabilityReasonValues returns the possible values for the ServiceNameUnavailabilityReason const type.
func PossibleServiceNameUnavailabilityReasonValues() []ServiceNameUnavailabilityReason {
	return []ServiceNameUnavailabilityReason{
		ServiceNameUnavailabilityReasonInvalid,
		ServiceNameUnavailabilityReasonAlreadyExists,
	}
}
