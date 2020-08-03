package servicefabric

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
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

// DefaultMoveCost enumerates the values for default move cost.
type DefaultMoveCost string

const (
	// High ...
	High DefaultMoveCost = "High"
	// Low ...
	Low DefaultMoveCost = "Low"
	// Medium ...
	Medium DefaultMoveCost = "Medium"
	// Zero ...
	Zero DefaultMoveCost = "Zero"
)

// PossibleDefaultMoveCostValues returns an array of possible values for the DefaultMoveCost const type.
func PossibleDefaultMoveCostValues() []DefaultMoveCost {
	return []DefaultMoveCost{High, Low, Medium, Zero}
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

// PartitionScheme enumerates the values for partition scheme.
type PartitionScheme string

const (
	// PartitionSchemeNamed ...
	PartitionSchemeNamed PartitionScheme = "Named"
	// PartitionSchemePartitionSchemeDescription ...
	PartitionSchemePartitionSchemeDescription PartitionScheme = "PartitionSchemeDescription"
	// PartitionSchemeSingleton ...
	PartitionSchemeSingleton PartitionScheme = "Singleton"
	// PartitionSchemeUniformInt64Range ...
	PartitionSchemeUniformInt64Range PartitionScheme = "UniformInt64Range"
)

// PossiblePartitionSchemeValues returns an array of possible values for the PartitionScheme const type.
func PossiblePartitionSchemeValues() []PartitionScheme {
	return []PartitionScheme{PartitionSchemeNamed, PartitionSchemePartitionSchemeDescription, PartitionSchemeSingleton, PartitionSchemeUniformInt64Range}
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
	// ReliabilityLevelNone ...
	ReliabilityLevelNone ReliabilityLevel = "None"
	// ReliabilityLevelPlatinum ...
	ReliabilityLevelPlatinum ReliabilityLevel = "Platinum"
	// ReliabilityLevelSilver ...
	ReliabilityLevelSilver ReliabilityLevel = "Silver"
)

// PossibleReliabilityLevelValues returns an array of possible values for the ReliabilityLevel const type.
func PossibleReliabilityLevelValues() []ReliabilityLevel {
	return []ReliabilityLevel{ReliabilityLevelBronze, ReliabilityLevelGold, ReliabilityLevelNone, ReliabilityLevelPlatinum, ReliabilityLevelSilver}
}

// ReliabilityLevel1 enumerates the values for reliability level 1.
type ReliabilityLevel1 string

const (
	// ReliabilityLevel1Bronze ...
	ReliabilityLevel1Bronze ReliabilityLevel1 = "Bronze"
	// ReliabilityLevel1Gold ...
	ReliabilityLevel1Gold ReliabilityLevel1 = "Gold"
	// ReliabilityLevel1Silver ...
	ReliabilityLevel1Silver ReliabilityLevel1 = "Silver"
)

// PossibleReliabilityLevel1Values returns an array of possible values for the ReliabilityLevel1 const type.
func PossibleReliabilityLevel1Values() []ReliabilityLevel1 {
	return []ReliabilityLevel1{ReliabilityLevel1Bronze, ReliabilityLevel1Gold, ReliabilityLevel1Silver}
}

// Scheme enumerates the values for scheme.
type Scheme string

const (
	// Affinity ...
	Affinity Scheme = "Affinity"
	// AlignedAffinity ...
	AlignedAffinity Scheme = "AlignedAffinity"
	// Invalid ...
	Invalid Scheme = "Invalid"
	// NonAlignedAffinity ...
	NonAlignedAffinity Scheme = "NonAlignedAffinity"
)

// PossibleSchemeValues returns an array of possible values for the Scheme const type.
func PossibleSchemeValues() []Scheme {
	return []Scheme{Affinity, AlignedAffinity, Invalid, NonAlignedAffinity}
}

// ServiceKind enumerates the values for service kind.
type ServiceKind string

const (
	// ServiceKindServiceProperties ...
	ServiceKindServiceProperties ServiceKind = "ServiceProperties"
	// ServiceKindStateful ...
	ServiceKindStateful ServiceKind = "Stateful"
	// ServiceKindStateless ...
	ServiceKindStateless ServiceKind = "Stateless"
)

// PossibleServiceKindValues returns an array of possible values for the ServiceKind const type.
func PossibleServiceKindValues() []ServiceKind {
	return []ServiceKind{ServiceKindServiceProperties, ServiceKindStateful, ServiceKindStateless}
}

// ServiceKindBasicServiceUpdateProperties enumerates the values for service kind basic service update
// properties.
type ServiceKindBasicServiceUpdateProperties string

const (
	// ServiceKindBasicServiceUpdatePropertiesServiceKindServiceUpdateProperties ...
	ServiceKindBasicServiceUpdatePropertiesServiceKindServiceUpdateProperties ServiceKindBasicServiceUpdateProperties = "ServiceUpdateProperties"
	// ServiceKindBasicServiceUpdatePropertiesServiceKindStateful ...
	ServiceKindBasicServiceUpdatePropertiesServiceKindStateful ServiceKindBasicServiceUpdateProperties = "Stateful"
	// ServiceKindBasicServiceUpdatePropertiesServiceKindStateless ...
	ServiceKindBasicServiceUpdatePropertiesServiceKindStateless ServiceKindBasicServiceUpdateProperties = "Stateless"
)

// PossibleServiceKindBasicServiceUpdatePropertiesValues returns an array of possible values for the ServiceKindBasicServiceUpdateProperties const type.
func PossibleServiceKindBasicServiceUpdatePropertiesValues() []ServiceKindBasicServiceUpdateProperties {
	return []ServiceKindBasicServiceUpdateProperties{ServiceKindBasicServiceUpdatePropertiesServiceKindServiceUpdateProperties, ServiceKindBasicServiceUpdatePropertiesServiceKindStateful, ServiceKindBasicServiceUpdatePropertiesServiceKindStateless}
}

// Type enumerates the values for type.
type Type string

const (
	// TypeServicePlacementPolicyDescription ...
	TypeServicePlacementPolicyDescription Type = "ServicePlacementPolicyDescription"
)

// PossibleTypeValues returns an array of possible values for the Type const type.
func PossibleTypeValues() []Type {
	return []Type{TypeServicePlacementPolicyDescription}
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

// Weight enumerates the values for weight.
type Weight string

const (
	// WeightHigh ...
	WeightHigh Weight = "High"
	// WeightLow ...
	WeightLow Weight = "Low"
	// WeightMedium ...
	WeightMedium Weight = "Medium"
	// WeightZero ...
	WeightZero Weight = "Zero"
)

// PossibleWeightValues returns an array of possible values for the Weight const type.
func PossibleWeightValues() []Weight {
	return []Weight{WeightHigh, WeightLow, WeightMedium, WeightZero}
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
