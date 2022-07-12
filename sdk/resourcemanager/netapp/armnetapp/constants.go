//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armnetapp

const (
	moduleName    = "armnetapp"
	moduleVersion = "v2.0.0"
)

// ActiveDirectoryStatus - Status of the Active Directory
type ActiveDirectoryStatus string

const (
	// ActiveDirectoryStatusCreated - Active Directory created but not in use
	ActiveDirectoryStatusCreated ActiveDirectoryStatus = "Created"
	// ActiveDirectoryStatusDeleted - Active Directory Deleted
	ActiveDirectoryStatusDeleted ActiveDirectoryStatus = "Deleted"
	// ActiveDirectoryStatusError - Error with the Active Directory
	ActiveDirectoryStatusError ActiveDirectoryStatus = "Error"
	// ActiveDirectoryStatusInUse - Active Directory in use by SMB Volume
	ActiveDirectoryStatusInUse ActiveDirectoryStatus = "InUse"
	// ActiveDirectoryStatusUpdating - Active Directory Updating
	ActiveDirectoryStatusUpdating ActiveDirectoryStatus = "Updating"
)

// PossibleActiveDirectoryStatusValues returns the possible values for the ActiveDirectoryStatus const type.
func PossibleActiveDirectoryStatusValues() []ActiveDirectoryStatus {
	return []ActiveDirectoryStatus{
		ActiveDirectoryStatusCreated,
		ActiveDirectoryStatusDeleted,
		ActiveDirectoryStatusError,
		ActiveDirectoryStatusInUse,
		ActiveDirectoryStatusUpdating,
	}
}

// ApplicationType - Application Type
type ApplicationType string

const (
	ApplicationTypeSAPHANA ApplicationType = "SAP-HANA"
)

// PossibleApplicationTypeValues returns the possible values for the ApplicationType const type.
func PossibleApplicationTypeValues() []ApplicationType {
	return []ApplicationType{
		ApplicationTypeSAPHANA,
	}
}

// AvsDataStore - Specifies whether the volume is enabled for Azure VMware Solution (AVS) datastore purpose
type AvsDataStore string

const (
	// AvsDataStoreDisabled - avsDataStore is disabled
	AvsDataStoreDisabled AvsDataStore = "Disabled"
	// AvsDataStoreEnabled - avsDataStore is enabled
	AvsDataStoreEnabled AvsDataStore = "Enabled"
)

// PossibleAvsDataStoreValues returns the possible values for the AvsDataStore const type.
func PossibleAvsDataStoreValues() []AvsDataStore {
	return []AvsDataStore{
		AvsDataStoreDisabled,
		AvsDataStoreEnabled,
	}
}

// BackupType - Type of backup Manual or Scheduled
type BackupType string

const (
	// BackupTypeManual - Manual backup
	BackupTypeManual BackupType = "Manual"
	// BackupTypeScheduled - Scheduled backup
	BackupTypeScheduled BackupType = "Scheduled"
)

// PossibleBackupTypeValues returns the possible values for the BackupType const type.
func PossibleBackupTypeValues() []BackupType {
	return []BackupType{
		BackupTypeManual,
		BackupTypeScheduled,
	}
}

// CheckNameResourceTypes - Resource type used for verification.
type CheckNameResourceTypes string

const (
	CheckNameResourceTypesMicrosoftNetAppNetAppAccounts                              CheckNameResourceTypes = "Microsoft.NetApp/netAppAccounts"
	CheckNameResourceTypesMicrosoftNetAppNetAppAccountsCapacityPools                 CheckNameResourceTypes = "Microsoft.NetApp/netAppAccounts/capacityPools"
	CheckNameResourceTypesMicrosoftNetAppNetAppAccountsCapacityPoolsVolumes          CheckNameResourceTypes = "Microsoft.NetApp/netAppAccounts/capacityPools/volumes"
	CheckNameResourceTypesMicrosoftNetAppNetAppAccountsCapacityPoolsVolumesSnapshots CheckNameResourceTypes = "Microsoft.NetApp/netAppAccounts/capacityPools/volumes/snapshots"
)

// PossibleCheckNameResourceTypesValues returns the possible values for the CheckNameResourceTypes const type.
func PossibleCheckNameResourceTypesValues() []CheckNameResourceTypes {
	return []CheckNameResourceTypes{
		CheckNameResourceTypesMicrosoftNetAppNetAppAccounts,
		CheckNameResourceTypesMicrosoftNetAppNetAppAccountsCapacityPools,
		CheckNameResourceTypesMicrosoftNetAppNetAppAccountsCapacityPoolsVolumes,
		CheckNameResourceTypesMicrosoftNetAppNetAppAccountsCapacityPoolsVolumesSnapshots,
	}
}

// CheckQuotaNameResourceTypes - Resource type used for verification.
type CheckQuotaNameResourceTypes string

const (
	CheckQuotaNameResourceTypesMicrosoftNetAppNetAppAccounts                              CheckQuotaNameResourceTypes = "Microsoft.NetApp/netAppAccounts"
	CheckQuotaNameResourceTypesMicrosoftNetAppNetAppAccountsCapacityPools                 CheckQuotaNameResourceTypes = "Microsoft.NetApp/netAppAccounts/capacityPools"
	CheckQuotaNameResourceTypesMicrosoftNetAppNetAppAccountsCapacityPoolsVolumes          CheckQuotaNameResourceTypes = "Microsoft.NetApp/netAppAccounts/capacityPools/volumes"
	CheckQuotaNameResourceTypesMicrosoftNetAppNetAppAccountsCapacityPoolsVolumesSnapshots CheckQuotaNameResourceTypes = "Microsoft.NetApp/netAppAccounts/capacityPools/volumes/snapshots"
)

// PossibleCheckQuotaNameResourceTypesValues returns the possible values for the CheckQuotaNameResourceTypes const type.
func PossibleCheckQuotaNameResourceTypesValues() []CheckQuotaNameResourceTypes {
	return []CheckQuotaNameResourceTypes{
		CheckQuotaNameResourceTypesMicrosoftNetAppNetAppAccounts,
		CheckQuotaNameResourceTypesMicrosoftNetAppNetAppAccountsCapacityPools,
		CheckQuotaNameResourceTypesMicrosoftNetAppNetAppAccountsCapacityPoolsVolumes,
		CheckQuotaNameResourceTypesMicrosoftNetAppNetAppAccountsCapacityPoolsVolumesSnapshots,
	}
}

// ChownMode - This parameter specifies who is authorized to change the ownership of a file. restricted - Only root user can
// change the ownership of the file. unrestricted - Non-root users can change ownership of
// files that they own.
type ChownMode string

const (
	ChownModeRestricted   ChownMode = "Restricted"
	ChownModeUnrestricted ChownMode = "Unrestricted"
)

// PossibleChownModeValues returns the possible values for the ChownMode const type.
func PossibleChownModeValues() []ChownMode {
	return []ChownMode{
		ChownModeRestricted,
		ChownModeUnrestricted,
	}
}

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

// EnableSubvolumes - Flag indicating whether subvolume operations are enabled on the volume
type EnableSubvolumes string

const (
	// EnableSubvolumesDisabled - subvolumes are not enabled
	EnableSubvolumesDisabled EnableSubvolumes = "Disabled"
	// EnableSubvolumesEnabled - subvolumes are enabled
	EnableSubvolumesEnabled EnableSubvolumes = "Enabled"
)

// PossibleEnableSubvolumesValues returns the possible values for the EnableSubvolumes const type.
func PossibleEnableSubvolumesValues() []EnableSubvolumes {
	return []EnableSubvolumes{
		EnableSubvolumesDisabled,
		EnableSubvolumesEnabled,
	}
}

// EncryptionKeySource - Source of key used to encrypt data in volume. Possible values (case-insensitive) are: 'Microsoft.NetApp'
type EncryptionKeySource string

const (
	// EncryptionKeySourceMicrosoftNetApp - Microsoft-managed key encryption
	EncryptionKeySourceMicrosoftNetApp EncryptionKeySource = "Microsoft.NetApp"
)

// PossibleEncryptionKeySourceValues returns the possible values for the EncryptionKeySource const type.
func PossibleEncryptionKeySourceValues() []EncryptionKeySource {
	return []EncryptionKeySource{
		EncryptionKeySourceMicrosoftNetApp,
	}
}

// EncryptionType - Encryption type of the capacity pool, set encryption type for data at rest for this pool and all volumes
// in it. This value can only be set when creating new pool.
type EncryptionType string

const (
	// EncryptionTypeDouble - EncryptionType Double, volumes will use double encryption at rest
	EncryptionTypeDouble EncryptionType = "Double"
	// EncryptionTypeSingle - EncryptionType Single, volumes will use single encryption at rest
	EncryptionTypeSingle EncryptionType = "Single"
)

// PossibleEncryptionTypeValues returns the possible values for the EncryptionType const type.
func PossibleEncryptionTypeValues() []EncryptionType {
	return []EncryptionType{
		EncryptionTypeDouble,
		EncryptionTypeSingle,
	}
}

// EndpointType - Indicates whether the local volume is the source or destination for the Volume Replication
type EndpointType string

const (
	EndpointTypeDst EndpointType = "dst"
	EndpointTypeSrc EndpointType = "src"
)

// PossibleEndpointTypeValues returns the possible values for the EndpointType const type.
func PossibleEndpointTypeValues() []EndpointType {
	return []EndpointType{
		EndpointTypeDst,
		EndpointTypeSrc,
	}
}

// InAvailabilityReasonType - Invalid indicates the name provided does not match Azure App Service naming requirements. AlreadyExists
// indicates that the name is already in use and is therefore unavailable.
type InAvailabilityReasonType string

const (
	InAvailabilityReasonTypeAlreadyExists InAvailabilityReasonType = "AlreadyExists"
	InAvailabilityReasonTypeInvalid       InAvailabilityReasonType = "Invalid"
)

// PossibleInAvailabilityReasonTypeValues returns the possible values for the InAvailabilityReasonType const type.
func PossibleInAvailabilityReasonTypeValues() []InAvailabilityReasonType {
	return []InAvailabilityReasonType{
		InAvailabilityReasonTypeAlreadyExists,
		InAvailabilityReasonTypeInvalid,
	}
}

type MetricAggregationType string

const (
	MetricAggregationTypeAverage MetricAggregationType = "Average"
)

// PossibleMetricAggregationTypeValues returns the possible values for the MetricAggregationType const type.
func PossibleMetricAggregationTypeValues() []MetricAggregationType {
	return []MetricAggregationType{
		MetricAggregationTypeAverage,
	}
}

// MirrorState - The status of the replication
type MirrorState string

const (
	MirrorStateBroken        MirrorState = "Broken"
	MirrorStateMirrored      MirrorState = "Mirrored"
	MirrorStateUninitialized MirrorState = "Uninitialized"
)

// PossibleMirrorStateValues returns the possible values for the MirrorState const type.
func PossibleMirrorStateValues() []MirrorState {
	return []MirrorState{
		MirrorStateBroken,
		MirrorStateMirrored,
		MirrorStateUninitialized,
	}
}

// NetworkFeatures - Basic network, or Standard features available to the volume.
type NetworkFeatures string

const (
	// NetworkFeaturesBasic - Basic network feature.
	NetworkFeaturesBasic NetworkFeatures = "Basic"
	// NetworkFeaturesStandard - Standard network feature.
	NetworkFeaturesStandard NetworkFeatures = "Standard"
)

// PossibleNetworkFeaturesValues returns the possible values for the NetworkFeatures const type.
func PossibleNetworkFeaturesValues() []NetworkFeatures {
	return []NetworkFeatures{
		NetworkFeaturesBasic,
		NetworkFeaturesStandard,
	}
}

// ProvisioningState - Gets the status of the VolumeQuotaRule at the time the operation was called.
type ProvisioningState string

const (
	ProvisioningStateAccepted  ProvisioningState = "Accepted"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStatePatching  ProvisioningState = "Patching"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateMoving    ProvisioningState = "Moving"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

// PossibleProvisioningStateValues returns the possible values for the ProvisioningState const type.
func PossibleProvisioningStateValues() []ProvisioningState {
	return []ProvisioningState{
		ProvisioningStateAccepted,
		ProvisioningStateCreating,
		ProvisioningStatePatching,
		ProvisioningStateDeleting,
		ProvisioningStateMoving,
		ProvisioningStateFailed,
		ProvisioningStateSucceeded,
	}
}

// QosType - The qos type of the pool
type QosType string

const (
	// QosTypeAuto - qos type Auto
	QosTypeAuto QosType = "Auto"
	// QosTypeManual - qos type Manual
	QosTypeManual QosType = "Manual"
)

// PossibleQosTypeValues returns the possible values for the QosType const type.
func PossibleQosTypeValues() []QosType {
	return []QosType{
		QosTypeAuto,
		QosTypeManual,
	}
}

// RelationshipStatus - Status of the mirror relationship
type RelationshipStatus string

const (
	RelationshipStatusIdle         RelationshipStatus = "Idle"
	RelationshipStatusTransferring RelationshipStatus = "Transferring"
)

// PossibleRelationshipStatusValues returns the possible values for the RelationshipStatus const type.
func PossibleRelationshipStatusValues() []RelationshipStatus {
	return []RelationshipStatus{
		RelationshipStatusIdle,
		RelationshipStatusTransferring,
	}
}

// ReplicationSchedule - Schedule
type ReplicationSchedule string

const (
	ReplicationSchedule10Minutely ReplicationSchedule = "_10minutely"
	ReplicationScheduleDaily      ReplicationSchedule = "daily"
	ReplicationScheduleHourly     ReplicationSchedule = "hourly"
)

// PossibleReplicationScheduleValues returns the possible values for the ReplicationSchedule const type.
func PossibleReplicationScheduleValues() []ReplicationSchedule {
	return []ReplicationSchedule{
		ReplicationSchedule10Minutely,
		ReplicationScheduleDaily,
		ReplicationScheduleHourly,
	}
}

// SecurityStyle - The security style of volume, default unix, defaults to ntfs for dual protocol or CIFS protocol
type SecurityStyle string

const (
	SecurityStyleNtfs SecurityStyle = "ntfs"
	SecurityStyleUnix SecurityStyle = "unix"
)

// PossibleSecurityStyleValues returns the possible values for the SecurityStyle const type.
func PossibleSecurityStyleValues() []SecurityStyle {
	return []SecurityStyle{
		SecurityStyleNtfs,
		SecurityStyleUnix,
	}
}

// ServiceLevel - The service level of the file system
type ServiceLevel string

const (
	// ServiceLevelPremium - Premium service level
	ServiceLevelPremium ServiceLevel = "Premium"
	// ServiceLevelStandard - Standard service level
	ServiceLevelStandard ServiceLevel = "Standard"
	// ServiceLevelStandardZRS - Zone redundant storage service level
	ServiceLevelStandardZRS ServiceLevel = "StandardZRS"
	// ServiceLevelUltra - Ultra service level
	ServiceLevelUltra ServiceLevel = "Ultra"
)

// PossibleServiceLevelValues returns the possible values for the ServiceLevel const type.
func PossibleServiceLevelValues() []ServiceLevel {
	return []ServiceLevel{
		ServiceLevelPremium,
		ServiceLevelStandard,
		ServiceLevelStandardZRS,
		ServiceLevelUltra,
	}
}

// Type - Type of quota
type Type string

const (
	// TypeDefaultGroupQuota - Default group quota
	TypeDefaultGroupQuota Type = "DefaultGroupQuota"
	// TypeDefaultUserQuota - Default user quota
	TypeDefaultUserQuota Type = "DefaultUserQuota"
	// TypeIndividualGroupQuota - Individual group quota
	TypeIndividualGroupQuota Type = "IndividualGroupQuota"
	// TypeIndividualUserQuota - Individual user quota
	TypeIndividualUserQuota Type = "IndividualUserQuota"
)

// PossibleTypeValues returns the possible values for the Type const type.
func PossibleTypeValues() []Type {
	return []Type{
		TypeDefaultGroupQuota,
		TypeDefaultUserQuota,
		TypeIndividualGroupQuota,
		TypeIndividualUserQuota,
	}
}

// VolumeStorageToNetworkProximity - Provides storage to network proximity information for the volume.
type VolumeStorageToNetworkProximity string

const (
	// VolumeStorageToNetworkProximityDefault - Basic storage to network connectivity.
	VolumeStorageToNetworkProximityDefault VolumeStorageToNetworkProximity = "Default"
	// VolumeStorageToNetworkProximityT1 - Standard T1 storage to network connectivity.
	VolumeStorageToNetworkProximityT1 VolumeStorageToNetworkProximity = "T1"
	// VolumeStorageToNetworkProximityT2 - Standard T2 storage to network connectivity.
	VolumeStorageToNetworkProximityT2 VolumeStorageToNetworkProximity = "T2"
)

// PossibleVolumeStorageToNetworkProximityValues returns the possible values for the VolumeStorageToNetworkProximity const type.
func PossibleVolumeStorageToNetworkProximityValues() []VolumeStorageToNetworkProximity {
	return []VolumeStorageToNetworkProximity{
		VolumeStorageToNetworkProximityDefault,
		VolumeStorageToNetworkProximityT1,
		VolumeStorageToNetworkProximityT2,
	}
}
