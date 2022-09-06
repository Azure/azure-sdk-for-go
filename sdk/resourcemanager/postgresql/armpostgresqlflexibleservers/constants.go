//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armpostgresqlflexibleservers

const (
	moduleName    = "armpostgresqlflexibleservers"
	moduleVersion = "v1.1.0"
)

// ConfigurationDataType - Data type of the configuration.
type ConfigurationDataType string

const (
	ConfigurationDataTypeBoolean     ConfigurationDataType = "Boolean"
	ConfigurationDataTypeEnumeration ConfigurationDataType = "Enumeration"
	ConfigurationDataTypeInteger     ConfigurationDataType = "Integer"
	ConfigurationDataTypeNumeric     ConfigurationDataType = "Numeric"
)

// PossibleConfigurationDataTypeValues returns the possible values for the ConfigurationDataType const type.
func PossibleConfigurationDataTypeValues() []ConfigurationDataType {
	return []ConfigurationDataType{
		ConfigurationDataTypeBoolean,
		ConfigurationDataTypeEnumeration,
		ConfigurationDataTypeInteger,
		ConfigurationDataTypeNumeric,
	}
}

// CreateMode - The mode to create a new PostgreSQL server.
type CreateMode string

const (
	CreateModeCreate             CreateMode = "Create"
	CreateModeDefault            CreateMode = "Default"
	CreateModePointInTimeRestore CreateMode = "PointInTimeRestore"
	CreateModeUpdate             CreateMode = "Update"
)

// PossibleCreateModeValues returns the possible values for the CreateMode const type.
func PossibleCreateModeValues() []CreateMode {
	return []CreateMode{
		CreateModeCreate,
		CreateModeDefault,
		CreateModePointInTimeRestore,
		CreateModeUpdate,
	}
}

// CreateModeForUpdate - The mode to update a new PostgreSQL server.
type CreateModeForUpdate string

const (
	CreateModeForUpdateDefault CreateModeForUpdate = "Default"
	CreateModeForUpdateUpdate  CreateModeForUpdate = "Update"
)

// PossibleCreateModeForUpdateValues returns the possible values for the CreateModeForUpdate const type.
func PossibleCreateModeForUpdateValues() []CreateModeForUpdate {
	return []CreateModeForUpdate{
		CreateModeForUpdateDefault,
		CreateModeForUpdateUpdate,
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

// FailoverMode - Failover mode.
type FailoverMode string

const (
	FailoverModeForcedFailover    FailoverMode = "ForcedFailover"
	FailoverModeForcedSwitchover  FailoverMode = "ForcedSwitchover"
	FailoverModePlannedFailover   FailoverMode = "PlannedFailover"
	FailoverModePlannedSwitchover FailoverMode = "PlannedSwitchover"
)

// PossibleFailoverModeValues returns the possible values for the FailoverMode const type.
func PossibleFailoverModeValues() []FailoverMode {
	return []FailoverMode{
		FailoverModeForcedFailover,
		FailoverModeForcedSwitchover,
		FailoverModePlannedFailover,
		FailoverModePlannedSwitchover,
	}
}

// GeoRedundantBackupEnum - A value indicating whether Geo-Redundant backup is enabled on the server.
type GeoRedundantBackupEnum string

const (
	GeoRedundantBackupEnumDisabled GeoRedundantBackupEnum = "Disabled"
	GeoRedundantBackupEnumEnabled  GeoRedundantBackupEnum = "Enabled"
)

// PossibleGeoRedundantBackupEnumValues returns the possible values for the GeoRedundantBackupEnum const type.
func PossibleGeoRedundantBackupEnumValues() []GeoRedundantBackupEnum {
	return []GeoRedundantBackupEnum{
		GeoRedundantBackupEnumDisabled,
		GeoRedundantBackupEnumEnabled,
	}
}

// HighAvailabilityMode - The HA mode for the server.
type HighAvailabilityMode string

const (
	HighAvailabilityModeDisabled      HighAvailabilityMode = "Disabled"
	HighAvailabilityModeZoneRedundant HighAvailabilityMode = "ZoneRedundant"
)

// PossibleHighAvailabilityModeValues returns the possible values for the HighAvailabilityMode const type.
func PossibleHighAvailabilityModeValues() []HighAvailabilityMode {
	return []HighAvailabilityMode{
		HighAvailabilityModeDisabled,
		HighAvailabilityModeZoneRedundant,
	}
}

// OperationOrigin - The intended executor of the operation.
type OperationOrigin string

const (
	OperationOriginNotSpecified OperationOrigin = "NotSpecified"
	OperationOriginSystem       OperationOrigin = "system"
	OperationOriginUser         OperationOrigin = "user"
)

// PossibleOperationOriginValues returns the possible values for the OperationOrigin const type.
func PossibleOperationOriginValues() []OperationOrigin {
	return []OperationOrigin{
		OperationOriginNotSpecified,
		OperationOriginSystem,
		OperationOriginUser,
	}
}

// Reason - The name availability reason.
type Reason string

const (
	ReasonAlreadyExists Reason = "AlreadyExists"
	ReasonInvalid       Reason = "Invalid"
)

// PossibleReasonValues returns the possible values for the Reason const type.
func PossibleReasonValues() []Reason {
	return []Reason{
		ReasonAlreadyExists,
		ReasonInvalid,
	}
}

// SKUTier - The tier of the particular SKU, e.g. Burstable.
type SKUTier string

const (
	SKUTierBurstable       SKUTier = "Burstable"
	SKUTierGeneralPurpose  SKUTier = "GeneralPurpose"
	SKUTierMemoryOptimized SKUTier = "MemoryOptimized"
)

// PossibleSKUTierValues returns the possible values for the SKUTier const type.
func PossibleSKUTierValues() []SKUTier {
	return []SKUTier{
		SKUTierBurstable,
		SKUTierGeneralPurpose,
		SKUTierMemoryOptimized,
	}
}

// ServerHAState - A state of a HA server that is visible to user.
type ServerHAState string

const (
	ServerHAStateCreatingStandby ServerHAState = "CreatingStandby"
	ServerHAStateFailingOver     ServerHAState = "FailingOver"
	ServerHAStateHealthy         ServerHAState = "Healthy"
	ServerHAStateNotEnabled      ServerHAState = "NotEnabled"
	ServerHAStateRemovingStandby ServerHAState = "RemovingStandby"
	ServerHAStateReplicatingData ServerHAState = "ReplicatingData"
)

// PossibleServerHAStateValues returns the possible values for the ServerHAState const type.
func PossibleServerHAStateValues() []ServerHAState {
	return []ServerHAState{
		ServerHAStateCreatingStandby,
		ServerHAStateFailingOver,
		ServerHAStateHealthy,
		ServerHAStateNotEnabled,
		ServerHAStateRemovingStandby,
		ServerHAStateReplicatingData,
	}
}

// ServerPublicNetworkAccessState - public network access is enabled or not
type ServerPublicNetworkAccessState string

const (
	ServerPublicNetworkAccessStateDisabled ServerPublicNetworkAccessState = "Disabled"
	ServerPublicNetworkAccessStateEnabled  ServerPublicNetworkAccessState = "Enabled"
)

// PossibleServerPublicNetworkAccessStateValues returns the possible values for the ServerPublicNetworkAccessState const type.
func PossibleServerPublicNetworkAccessStateValues() []ServerPublicNetworkAccessState {
	return []ServerPublicNetworkAccessState{
		ServerPublicNetworkAccessStateDisabled,
		ServerPublicNetworkAccessStateEnabled,
	}
}

// ServerState - A state of a server that is visible to user.
type ServerState string

const (
	ServerStateDisabled ServerState = "Disabled"
	ServerStateDropping ServerState = "Dropping"
	ServerStateReady    ServerState = "Ready"
	ServerStateStarting ServerState = "Starting"
	ServerStateStopped  ServerState = "Stopped"
	ServerStateStopping ServerState = "Stopping"
	ServerStateUpdating ServerState = "Updating"
)

// PossibleServerStateValues returns the possible values for the ServerState const type.
func PossibleServerStateValues() []ServerState {
	return []ServerState{
		ServerStateDisabled,
		ServerStateDropping,
		ServerStateReady,
		ServerStateStarting,
		ServerStateStopped,
		ServerStateStopping,
		ServerStateUpdating,
	}
}

// ServerVersion - The version of a server.
type ServerVersion string

const (
	ServerVersionEleven   ServerVersion = "11"
	ServerVersionFourteen ServerVersion = "14"
	ServerVersionThirteen ServerVersion = "13"
	ServerVersionTwelve   ServerVersion = "12"
)

// PossibleServerVersionValues returns the possible values for the ServerVersion const type.
func PossibleServerVersionValues() []ServerVersion {
	return []ServerVersion{
		ServerVersionEleven,
		ServerVersionFourteen,
		ServerVersionThirteen,
		ServerVersionTwelve,
	}
}
