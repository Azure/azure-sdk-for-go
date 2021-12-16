package servicefabric

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// ClusterState enumerates the values for cluster state.
type ClusterState string

const (
	// AutoScale ...
	AutoScale ClusterState = "AutoScale"
	// BaselineUpgrade ...
	BaselineUpgrade ClusterState = "BaselineUpgrade"
	// Deploying ...
	Deploying ClusterState = "Deploying"
	// EnforcingClusterVersion ...
	EnforcingClusterVersion ClusterState = "EnforcingClusterVersion"
	// Ready ...
	Ready ClusterState = "Ready"
	// UpdatingInfrastructure ...
	UpdatingInfrastructure ClusterState = "UpdatingInfrastructure"
	// UpdatingUserCertificate ...
	UpdatingUserCertificate ClusterState = "UpdatingUserCertificate"
	// UpdatingUserConfiguration ...
	UpdatingUserConfiguration ClusterState = "UpdatingUserConfiguration"
	// UpgradeServiceUnreachable ...
	UpgradeServiceUnreachable ClusterState = "UpgradeServiceUnreachable"
	// WaitingForNodes ...
	WaitingForNodes ClusterState = "WaitingForNodes"
)

// PossibleClusterStateValues returns an array of possible values for the ClusterState const type.
func PossibleClusterStateValues() []ClusterState {
	return []ClusterState{AutoScale, BaselineUpgrade, Deploying, EnforcingClusterVersion, Ready, UpdatingInfrastructure, UpdatingUserCertificate, UpdatingUserConfiguration, UpgradeServiceUnreachable, WaitingForNodes}
}

// DurabilityLevel enumerates the values for durability level.
type DurabilityLevel string

const (
	// Bronze ...
	Bronze DurabilityLevel = "Bronze"
	// Gold ...
	Gold DurabilityLevel = "Gold"
	// Silver ...
	Silver DurabilityLevel = "Silver"
)

// PossibleDurabilityLevelValues returns an array of possible values for the DurabilityLevel const type.
func PossibleDurabilityLevelValues() []DurabilityLevel {
	return []DurabilityLevel{Bronze, Gold, Silver}
}

// Environment enumerates the values for environment.
type Environment string

const (
	// Linux ...
	Linux Environment = "Linux"
	// Windows ...
	Windows Environment = "Windows"
)

// PossibleEnvironmentValues returns an array of possible values for the Environment const type.
func PossibleEnvironmentValues() []Environment {
	return []Environment{Linux, Windows}
}

// ProvisioningState enumerates the values for provisioning state.
type ProvisioningState string

const (
	// Canceled ...
	Canceled ProvisioningState = "Canceled"
	// Failed ...
	Failed ProvisioningState = "Failed"
	// Succeeded ...
	Succeeded ProvisioningState = "Succeeded"
	// Updating ...
	Updating ProvisioningState = "Updating"
)

// PossibleProvisioningStateValues returns an array of possible values for the ProvisioningState const type.
func PossibleProvisioningStateValues() []ProvisioningState {
	return []ProvisioningState{Canceled, Failed, Succeeded, Updating}
}

// ReliabilityLevel enumerates the values for reliability level.
type ReliabilityLevel string

const (
	// ReliabilityLevelBronze ...
	ReliabilityLevelBronze ReliabilityLevel = "Bronze"
	// ReliabilityLevelGold ...
	ReliabilityLevelGold ReliabilityLevel = "Gold"
	// ReliabilityLevelSilver ...
	ReliabilityLevelSilver ReliabilityLevel = "Silver"
)

// PossibleReliabilityLevelValues returns an array of possible values for the ReliabilityLevel const type.
func PossibleReliabilityLevelValues() []ReliabilityLevel {
	return []ReliabilityLevel{ReliabilityLevelBronze, ReliabilityLevelGold, ReliabilityLevelSilver}
}

// ReliabilityLevel1 enumerates the values for reliability level 1.
type ReliabilityLevel1 string

const (
	// ReliabilityLevel1Bronze ...
	ReliabilityLevel1Bronze ReliabilityLevel1 = "Bronze"
	// ReliabilityLevel1Gold ...
	ReliabilityLevel1Gold ReliabilityLevel1 = "Gold"
	// ReliabilityLevel1Platinum ...
	ReliabilityLevel1Platinum ReliabilityLevel1 = "Platinum"
	// ReliabilityLevel1Silver ...
	ReliabilityLevel1Silver ReliabilityLevel1 = "Silver"
)

// PossibleReliabilityLevel1Values returns an array of possible values for the ReliabilityLevel1 const type.
func PossibleReliabilityLevel1Values() []ReliabilityLevel1 {
	return []ReliabilityLevel1{ReliabilityLevel1Bronze, ReliabilityLevel1Gold, ReliabilityLevel1Platinum, ReliabilityLevel1Silver}
}

// UpgradeMode enumerates the values for upgrade mode.
type UpgradeMode string

const (
	// Automatic ...
	Automatic UpgradeMode = "Automatic"
	// Manual ...
	Manual UpgradeMode = "Manual"
)

// PossibleUpgradeModeValues returns an array of possible values for the UpgradeMode const type.
func PossibleUpgradeModeValues() []UpgradeMode {
	return []UpgradeMode{Automatic, Manual}
}

// UpgradeMode1 enumerates the values for upgrade mode 1.
type UpgradeMode1 string

const (
	// UpgradeMode1Automatic ...
	UpgradeMode1Automatic UpgradeMode1 = "Automatic"
	// UpgradeMode1Manual ...
	UpgradeMode1Manual UpgradeMode1 = "Manual"
)

// PossibleUpgradeMode1Values returns an array of possible values for the UpgradeMode1 const type.
func PossibleUpgradeMode1Values() []UpgradeMode1 {
	return []UpgradeMode1{UpgradeMode1Automatic, UpgradeMode1Manual}
}

// X509StoreName enumerates the values for x509 store name.
type X509StoreName string

const (
	// AddressBook ...
	AddressBook X509StoreName = "AddressBook"
	// AuthRoot ...
	AuthRoot X509StoreName = "AuthRoot"
	// CertificateAuthority ...
	CertificateAuthority X509StoreName = "CertificateAuthority"
	// Disallowed ...
	Disallowed X509StoreName = "Disallowed"
	// My ...
	My X509StoreName = "My"
	// Root ...
	Root X509StoreName = "Root"
	// TrustedPeople ...
	TrustedPeople X509StoreName = "TrustedPeople"
	// TrustedPublisher ...
	TrustedPublisher X509StoreName = "TrustedPublisher"
)

// PossibleX509StoreNameValues returns an array of possible values for the X509StoreName const type.
func PossibleX509StoreNameValues() []X509StoreName {
	return []X509StoreName{AddressBook, AuthRoot, CertificateAuthority, Disallowed, My, Root, TrustedPeople, TrustedPublisher}
}
