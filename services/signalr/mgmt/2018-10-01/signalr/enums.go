package signalr

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// FeatureFlags enumerates the values for feature flags.
type FeatureFlags string

const (
	// EnableConnectivityLogs ...
	EnableConnectivityLogs FeatureFlags = "EnableConnectivityLogs"
	// ServiceMode ...
	ServiceMode FeatureFlags = "ServiceMode"
)

// PossibleFeatureFlagsValues returns an array of possible values for the FeatureFlags const type.
func PossibleFeatureFlagsValues() []FeatureFlags {
	return []FeatureFlags{EnableConnectivityLogs, ServiceMode}
}

// KeyType enumerates the values for key type.
type KeyType string

const (
	// Primary ...
	Primary KeyType = "Primary"
	// Secondary ...
	Secondary KeyType = "Secondary"
)

// PossibleKeyTypeValues returns an array of possible values for the KeyType const type.
func PossibleKeyTypeValues() []KeyType {
	return []KeyType{Primary, Secondary}
}

// ProvisioningState enumerates the values for provisioning state.
type ProvisioningState string

const (
	// Canceled ...
	Canceled ProvisioningState = "Canceled"
	// Creating ...
	Creating ProvisioningState = "Creating"
	// Deleting ...
	Deleting ProvisioningState = "Deleting"
	// Failed ...
	Failed ProvisioningState = "Failed"
	// Moving ...
	Moving ProvisioningState = "Moving"
	// Running ...
	Running ProvisioningState = "Running"
	// Succeeded ...
	Succeeded ProvisioningState = "Succeeded"
	// Unknown ...
	Unknown ProvisioningState = "Unknown"
	// Updating ...
	Updating ProvisioningState = "Updating"
)

// PossibleProvisioningStateValues returns an array of possible values for the ProvisioningState const type.
func PossibleProvisioningStateValues() []ProvisioningState {
	return []ProvisioningState{Canceled, Creating, Deleting, Failed, Moving, Running, Succeeded, Unknown, Updating}
}

// SkuTier enumerates the values for sku tier.
type SkuTier string

const (
	// Basic ...
	Basic SkuTier = "Basic"
	// Free ...
	Free SkuTier = "Free"
	// Premium ...
	Premium SkuTier = "Premium"
	// Standard ...
	Standard SkuTier = "Standard"
)

// PossibleSkuTierValues returns an array of possible values for the SkuTier const type.
func PossibleSkuTierValues() []SkuTier {
	return []SkuTier{Basic, Free, Premium, Standard}
}
