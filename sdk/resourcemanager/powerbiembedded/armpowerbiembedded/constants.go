//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armpowerbiembedded

const (
	moduleName    = "armpowerbiembedded"
	moduleVersion = "v1.1.0"
)

// AccessKeyName - Key name
type AccessKeyName string

const (
	AccessKeyNameKey1 AccessKeyName = "key1"
	AccessKeyNameKey2 AccessKeyName = "key2"
)

// PossibleAccessKeyNameValues returns the possible values for the AccessKeyName const type.
func PossibleAccessKeyNameValues() []AccessKeyName {
	return []AccessKeyName{
		AccessKeyNameKey1,
		AccessKeyNameKey2,
	}
}

// AzureSKUName - SKU name
type AzureSKUName string

const (
	AzureSKUNameS1 AzureSKUName = "S1"
)

// PossibleAzureSKUNameValues returns the possible values for the AzureSKUName const type.
func PossibleAzureSKUNameValues() []AzureSKUName {
	return []AzureSKUName{
		AzureSKUNameS1,
	}
}

// AzureSKUTier - SKU tier
type AzureSKUTier string

const (
	AzureSKUTierStandard AzureSKUTier = "Standard"
)

// PossibleAzureSKUTierValues returns the possible values for the AzureSKUTier const type.
func PossibleAzureSKUTierValues() []AzureSKUTier {
	return []AzureSKUTier{
		AzureSKUTierStandard,
	}
}

// CheckNameReason - Reason why the workspace collection name cannot be used.
type CheckNameReason string

const (
	CheckNameReasonInvalid     CheckNameReason = "Invalid"
	CheckNameReasonUnavailable CheckNameReason = "Unavailable"
)

// PossibleCheckNameReasonValues returns the possible values for the CheckNameReason const type.
func PossibleCheckNameReasonValues() []CheckNameReason {
	return []CheckNameReason{
		CheckNameReasonInvalid,
		CheckNameReasonUnavailable,
	}
}
