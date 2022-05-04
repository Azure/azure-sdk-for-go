//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmixedreality

const (
	moduleName    = "armmixedreality"
	moduleVersion = "v0.4.0"
)

// CreatedByType - The type of identity that created the resource.
type CreatedByType string

const (
	CreatedByTypeApplication     CreatedByType = "Application"
	CreatedByTypeKey             CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser            CreatedByType = "User"
)

// PossibleCreatedByTypeValues returns the possible values for the CreatedByType const type.
func PossibleCreatedByTypeValues() []CreatedByType {
	return []CreatedByType{
		CreatedByTypeApplication,
		CreatedByTypeKey,
		CreatedByTypeManagedIdentity,
		CreatedByTypeUser,
	}
}

// NameUnavailableReason - reason of name unavailable.
type NameUnavailableReason string

const (
	NameUnavailableReasonAlreadyExists NameUnavailableReason = "AlreadyExists"
	NameUnavailableReasonInvalid       NameUnavailableReason = "Invalid"
)

// PossibleNameUnavailableReasonValues returns the possible values for the NameUnavailableReason const type.
func PossibleNameUnavailableReasonValues() []NameUnavailableReason {
	return []NameUnavailableReason{
		NameUnavailableReasonAlreadyExists,
		NameUnavailableReasonInvalid,
	}
}

// SKUTier - This field is required to be implemented by the Resource Provider if the service has more than one tier, but
// is not required on a PUT.
type SKUTier string

const (
	SKUTierFree     SKUTier = "Free"
	SKUTierBasic    SKUTier = "Basic"
	SKUTierStandard SKUTier = "Standard"
	SKUTierPremium  SKUTier = "Premium"
)

// PossibleSKUTierValues returns the possible values for the SKUTier const type.
func PossibleSKUTierValues() []SKUTier {
	return []SKUTier{
		SKUTierFree,
		SKUTierBasic,
		SKUTierStandard,
		SKUTierPremium,
	}
}

// Serial - Serial of key to be regenerated
type Serial int32

const (
	// SerialPrimary - The Primary Key
	SerialPrimary Serial = 1
	// SerialSecondary - The Secondary Key
	SerialSecondary Serial = 2
)

// PossibleSerialValues returns the possible values for the Serial const type.
func PossibleSerialValues() []Serial {
	return []Serial{
		SerialPrimary,
		SerialSecondary,
	}
}
