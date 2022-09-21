package documentdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// AnalyticalStorageSchemaType enumerates the values for analytical storage schema type.
type AnalyticalStorageSchemaType string

const (
	// FullFidelity ...
	FullFidelity AnalyticalStorageSchemaType = "FullFidelity"
	// WellDefined ...
	WellDefined AnalyticalStorageSchemaType = "WellDefined"
)

// PossibleAnalyticalStorageSchemaTypeValues returns an array of possible values for the AnalyticalStorageSchemaType const type.
func PossibleAnalyticalStorageSchemaTypeValues() []AnalyticalStorageSchemaType {
	return []AnalyticalStorageSchemaType{FullFidelity, WellDefined}
}

// APIType enumerates the values for api type.
type APIType string

const (
	// Cassandra ...
	Cassandra APIType = "Cassandra"
	// Gremlin ...
	Gremlin APIType = "Gremlin"
	// GremlinV2 ...
	GremlinV2 APIType = "GremlinV2"
	// MongoDB ...
	MongoDB APIType = "MongoDB"
	// SQL ...
	SQL APIType = "Sql"
	// Table ...
	Table APIType = "Table"
)

// PossibleAPITypeValues returns an array of possible values for the APIType const type.
func PossibleAPITypeValues() []APIType {
	return []APIType{Cassandra, Gremlin, GremlinV2, MongoDB, SQL, Table}
}

// AuthenticationMethod enumerates the values for authentication method.
type AuthenticationMethod string

const (
	// AuthenticationMethodCassandra ...
	AuthenticationMethodCassandra AuthenticationMethod = "Cassandra"
	// AuthenticationMethodNone ...
	AuthenticationMethodNone AuthenticationMethod = "None"
)

// PossibleAuthenticationMethodValues returns an array of possible values for the AuthenticationMethod const type.
func PossibleAuthenticationMethodValues() []AuthenticationMethod {
	return []AuthenticationMethod{AuthenticationMethodCassandra, AuthenticationMethodNone}
}

// BackupPolicyMigrationStatus enumerates the values for backup policy migration status.
type BackupPolicyMigrationStatus string

const (
	// Completed ...
	Completed BackupPolicyMigrationStatus = "Completed"
	// Failed ...
	Failed BackupPolicyMigrationStatus = "Failed"
	// InProgress ...
	InProgress BackupPolicyMigrationStatus = "InProgress"
	// Invalid ...
	Invalid BackupPolicyMigrationStatus = "Invalid"
)

// PossibleBackupPolicyMigrationStatusValues returns an array of possible values for the BackupPolicyMigrationStatus const type.
func PossibleBackupPolicyMigrationStatusValues() []BackupPolicyMigrationStatus {
	return []BackupPolicyMigrationStatus{Completed, Failed, InProgress, Invalid}
}

// BackupPolicyType enumerates the values for backup policy type.
type BackupPolicyType string

const (
	// Continuous ...
	Continuous BackupPolicyType = "Continuous"
	// Periodic ...
	Periodic BackupPolicyType = "Periodic"
)

// PossibleBackupPolicyTypeValues returns an array of possible values for the BackupPolicyType const type.
func PossibleBackupPolicyTypeValues() []BackupPolicyType {
	return []BackupPolicyType{Continuous, Periodic}
}

// BackupStorageRedundancy enumerates the values for backup storage redundancy.
type BackupStorageRedundancy string

const (
	// Geo ...
	Geo BackupStorageRedundancy = "Geo"
	// Local ...
	Local BackupStorageRedundancy = "Local"
	// Zone ...
	Zone BackupStorageRedundancy = "Zone"
)

// PossibleBackupStorageRedundancyValues returns an array of possible values for the BackupStorageRedundancy const type.
func PossibleBackupStorageRedundancyValues() []BackupStorageRedundancy {
	return []BackupStorageRedundancy{Geo, Local, Zone}
}

// CompositePathSortOrder enumerates the values for composite path sort order.
type CompositePathSortOrder string

const (
	// Ascending ...
	Ascending CompositePathSortOrder = "ascending"
	// Descending ...
	Descending CompositePathSortOrder = "descending"
)

// PossibleCompositePathSortOrderValues returns an array of possible values for the CompositePathSortOrder const type.
func PossibleCompositePathSortOrderValues() []CompositePathSortOrder {
	return []CompositePathSortOrder{Ascending, Descending}
}

// ConflictResolutionMode enumerates the values for conflict resolution mode.
type ConflictResolutionMode string

const (
	// Custom ...
	Custom ConflictResolutionMode = "Custom"
	// LastWriterWins ...
	LastWriterWins ConflictResolutionMode = "LastWriterWins"
)

// PossibleConflictResolutionModeValues returns an array of possible values for the ConflictResolutionMode const type.
func PossibleConflictResolutionModeValues() []ConflictResolutionMode {
	return []ConflictResolutionMode{Custom, LastWriterWins}
}

// ConnectionState enumerates the values for connection state.
type ConnectionState string

const (
	// DatacenterToDatacenterNetworkError ...
	DatacenterToDatacenterNetworkError ConnectionState = "DatacenterToDatacenterNetworkError"
	// InternalError ...
	InternalError ConnectionState = "InternalError"
	// InternalOperatorToDataCenterCertificateError ...
	InternalOperatorToDataCenterCertificateError ConnectionState = "InternalOperatorToDataCenterCertificateError"
	// OK ...
	OK ConnectionState = "OK"
	// OperatorToDataCenterNetworkError ...
	OperatorToDataCenterNetworkError ConnectionState = "OperatorToDataCenterNetworkError"
	// Unknown ...
	Unknown ConnectionState = "Unknown"
)

// PossibleConnectionStateValues returns an array of possible values for the ConnectionState const type.
func PossibleConnectionStateValues() []ConnectionState {
	return []ConnectionState{DatacenterToDatacenterNetworkError, InternalError, InternalOperatorToDataCenterCertificateError, OK, OperatorToDataCenterNetworkError, Unknown}
}

// ConnectorOffer enumerates the values for connector offer.
type ConnectorOffer string

const (
	// Small ...
	Small ConnectorOffer = "Small"
)

// PossibleConnectorOfferValues returns an array of possible values for the ConnectorOffer const type.
func PossibleConnectorOfferValues() []ConnectorOffer {
	return []ConnectorOffer{Small}
}

// CreatedByType enumerates the values for created by type.
type CreatedByType string

const (
	// Application ...
	Application CreatedByType = "Application"
	// Key ...
	Key CreatedByType = "Key"
	// ManagedIdentity ...
	ManagedIdentity CreatedByType = "ManagedIdentity"
	// User ...
	User CreatedByType = "User"
)

// PossibleCreatedByTypeValues returns an array of possible values for the CreatedByType const type.
func PossibleCreatedByTypeValues() []CreatedByType {
	return []CreatedByType{Application, Key, ManagedIdentity, User}
}

// CreateMode enumerates the values for create mode.
type CreateMode string

const (
	// Default ...
	Default CreateMode = "Default"
	// Restore ...
	Restore CreateMode = "Restore"
)

// PossibleCreateModeValues returns an array of possible values for the CreateMode const type.
func PossibleCreateModeValues() []CreateMode {
	return []CreateMode{Default, Restore}
}

// DatabaseAccountKind enumerates the values for database account kind.
type DatabaseAccountKind string

const (
	// DatabaseAccountKindGlobalDocumentDB ...
	DatabaseAccountKindGlobalDocumentDB DatabaseAccountKind = "GlobalDocumentDB"
	// DatabaseAccountKindMongoDB ...
	DatabaseAccountKindMongoDB DatabaseAccountKind = "MongoDB"
	// DatabaseAccountKindParse ...
	DatabaseAccountKindParse DatabaseAccountKind = "Parse"
)

// PossibleDatabaseAccountKindValues returns an array of possible values for the DatabaseAccountKind const type.
func PossibleDatabaseAccountKindValues() []DatabaseAccountKind {
	return []DatabaseAccountKind{DatabaseAccountKindGlobalDocumentDB, DatabaseAccountKindMongoDB, DatabaseAccountKindParse}
}

// DatabaseAccountOfferType enumerates the values for database account offer type.
type DatabaseAccountOfferType string

const (
	// Standard ...
	Standard DatabaseAccountOfferType = "Standard"
)

// PossibleDatabaseAccountOfferTypeValues returns an array of possible values for the DatabaseAccountOfferType const type.
func PossibleDatabaseAccountOfferTypeValues() []DatabaseAccountOfferType {
	return []DatabaseAccountOfferType{Standard}
}

// DataType enumerates the values for data type.
type DataType string

const (
	// LineString ...
	LineString DataType = "LineString"
	// MultiPolygon ...
	MultiPolygon DataType = "MultiPolygon"
	// Number ...
	Number DataType = "Number"
	// Point ...
	Point DataType = "Point"
	// Polygon ...
	Polygon DataType = "Polygon"
	// String ...
	String DataType = "String"
)

// PossibleDataTypeValues returns an array of possible values for the DataType const type.
func PossibleDataTypeValues() []DataType {
	return []DataType{LineString, MultiPolygon, Number, Point, Polygon, String}
}

// DefaultConsistencyLevel enumerates the values for default consistency level.
type DefaultConsistencyLevel string

const (
	// BoundedStaleness ...
	BoundedStaleness DefaultConsistencyLevel = "BoundedStaleness"
	// ConsistentPrefix ...
	ConsistentPrefix DefaultConsistencyLevel = "ConsistentPrefix"
	// Eventual ...
	Eventual DefaultConsistencyLevel = "Eventual"
	// Session ...
	Session DefaultConsistencyLevel = "Session"
	// Strong ...
	Strong DefaultConsistencyLevel = "Strong"
)

// PossibleDefaultConsistencyLevelValues returns an array of possible values for the DefaultConsistencyLevel const type.
func PossibleDefaultConsistencyLevelValues() []DefaultConsistencyLevel {
	return []DefaultConsistencyLevel{BoundedStaleness, ConsistentPrefix, Eventual, Session, Strong}
}

// IndexingMode enumerates the values for indexing mode.
type IndexingMode string

const (
	// Consistent ...
	Consistent IndexingMode = "consistent"
	// Lazy ...
	Lazy IndexingMode = "lazy"
	// None ...
	None IndexingMode = "none"
)

// PossibleIndexingModeValues returns an array of possible values for the IndexingMode const type.
func PossibleIndexingModeValues() []IndexingMode {
	return []IndexingMode{Consistent, Lazy, None}
}

// IndexKind enumerates the values for index kind.
type IndexKind string

const (
	// Hash ...
	Hash IndexKind = "Hash"
	// Range ...
	Range IndexKind = "Range"
	// Spatial ...
	Spatial IndexKind = "Spatial"
)

// PossibleIndexKindValues returns an array of possible values for the IndexKind const type.
func PossibleIndexKindValues() []IndexKind {
	return []IndexKind{Hash, Range, Spatial}
}

// KeyKind enumerates the values for key kind.
type KeyKind string

const (
	// Primary ...
	Primary KeyKind = "primary"
	// PrimaryReadonly ...
	PrimaryReadonly KeyKind = "primaryReadonly"
	// Secondary ...
	Secondary KeyKind = "secondary"
	// SecondaryReadonly ...
	SecondaryReadonly KeyKind = "secondaryReadonly"
)

// PossibleKeyKindValues returns an array of possible values for the KeyKind const type.
func PossibleKeyKindValues() []KeyKind {
	return []KeyKind{Primary, PrimaryReadonly, Secondary, SecondaryReadonly}
}

// ManagedCassandraProvisioningState enumerates the values for managed cassandra provisioning state.
type ManagedCassandraProvisioningState string

const (
	// ManagedCassandraProvisioningStateCanceled ...
	ManagedCassandraProvisioningStateCanceled ManagedCassandraProvisioningState = "Canceled"
	// ManagedCassandraProvisioningStateCreating ...
	ManagedCassandraProvisioningStateCreating ManagedCassandraProvisioningState = "Creating"
	// ManagedCassandraProvisioningStateDeleting ...
	ManagedCassandraProvisioningStateDeleting ManagedCassandraProvisioningState = "Deleting"
	// ManagedCassandraProvisioningStateFailed ...
	ManagedCassandraProvisioningStateFailed ManagedCassandraProvisioningState = "Failed"
	// ManagedCassandraProvisioningStateSucceeded ...
	ManagedCassandraProvisioningStateSucceeded ManagedCassandraProvisioningState = "Succeeded"
	// ManagedCassandraProvisioningStateUpdating ...
	ManagedCassandraProvisioningStateUpdating ManagedCassandraProvisioningState = "Updating"
)

// PossibleManagedCassandraProvisioningStateValues returns an array of possible values for the ManagedCassandraProvisioningState const type.
func PossibleManagedCassandraProvisioningStateValues() []ManagedCassandraProvisioningState {
	return []ManagedCassandraProvisioningState{ManagedCassandraProvisioningStateCanceled, ManagedCassandraProvisioningStateCreating, ManagedCassandraProvisioningStateDeleting, ManagedCassandraProvisioningStateFailed, ManagedCassandraProvisioningStateSucceeded, ManagedCassandraProvisioningStateUpdating}
}

// ManagedCassandraResourceIdentityType enumerates the values for managed cassandra resource identity type.
type ManagedCassandraResourceIdentityType string

const (
	// ManagedCassandraResourceIdentityTypeNone ...
	ManagedCassandraResourceIdentityTypeNone ManagedCassandraResourceIdentityType = "None"
	// ManagedCassandraResourceIdentityTypeSystemAssigned ...
	ManagedCassandraResourceIdentityTypeSystemAssigned ManagedCassandraResourceIdentityType = "SystemAssigned"
)

// PossibleManagedCassandraResourceIdentityTypeValues returns an array of possible values for the ManagedCassandraResourceIdentityType const type.
func PossibleManagedCassandraResourceIdentityTypeValues() []ManagedCassandraResourceIdentityType {
	return []ManagedCassandraResourceIdentityType{ManagedCassandraResourceIdentityTypeNone, ManagedCassandraResourceIdentityTypeSystemAssigned}
}

// MongoRoleDefinitionType enumerates the values for mongo role definition type.
type MongoRoleDefinitionType string

const (
	// BuiltInRole ...
	BuiltInRole MongoRoleDefinitionType = "BuiltInRole"
	// CustomRole ...
	CustomRole MongoRoleDefinitionType = "CustomRole"
)

// PossibleMongoRoleDefinitionTypeValues returns an array of possible values for the MongoRoleDefinitionType const type.
func PossibleMongoRoleDefinitionTypeValues() []MongoRoleDefinitionType {
	return []MongoRoleDefinitionType{BuiltInRole, CustomRole}
}

// NetworkACLBypass enumerates the values for network acl bypass.
type NetworkACLBypass string

const (
	// NetworkACLBypassAzureServices ...
	NetworkACLBypassAzureServices NetworkACLBypass = "AzureServices"
	// NetworkACLBypassNone ...
	NetworkACLBypassNone NetworkACLBypass = "None"
)

// PossibleNetworkACLBypassValues returns an array of possible values for the NetworkACLBypass const type.
func PossibleNetworkACLBypassValues() []NetworkACLBypass {
	return []NetworkACLBypass{NetworkACLBypassAzureServices, NetworkACLBypassNone}
}

// NodeState enumerates the values for node state.
type NodeState string

const (
	// Joining ...
	Joining NodeState = "Joining"
	// Leaving ...
	Leaving NodeState = "Leaving"
	// Moving ...
	Moving NodeState = "Moving"
	// Normal ...
	Normal NodeState = "Normal"
	// Stopped ...
	Stopped NodeState = "Stopped"
)

// PossibleNodeStateValues returns an array of possible values for the NodeState const type.
func PossibleNodeStateValues() []NodeState {
	return []NodeState{Joining, Leaving, Moving, Normal, Stopped}
}

// NodeStatus enumerates the values for node status.
type NodeStatus string

const (
	// Down ...
	Down NodeStatus = "Down"
	// Up ...
	Up NodeStatus = "Up"
)

// PossibleNodeStatusValues returns an array of possible values for the NodeStatus const type.
func PossibleNodeStatusValues() []NodeStatus {
	return []NodeStatus{Down, Up}
}

// OperationType enumerates the values for operation type.
type OperationType string

const (
	// Create ...
	Create OperationType = "Create"
	// Delete ...
	Delete OperationType = "Delete"
	// Replace ...
	Replace OperationType = "Replace"
	// SystemOperation ...
	SystemOperation OperationType = "SystemOperation"
)

// PossibleOperationTypeValues returns an array of possible values for the OperationType const type.
func PossibleOperationTypeValues() []OperationType {
	return []OperationType{Create, Delete, Replace, SystemOperation}
}

// PartitionKind enumerates the values for partition kind.
type PartitionKind string

const (
	// PartitionKindHash ...
	PartitionKindHash PartitionKind = "Hash"
	// PartitionKindMultiHash ...
	PartitionKindMultiHash PartitionKind = "MultiHash"
	// PartitionKindRange ...
	PartitionKindRange PartitionKind = "Range"
)

// PossiblePartitionKindValues returns an array of possible values for the PartitionKind const type.
func PossiblePartitionKindValues() []PartitionKind {
	return []PartitionKind{PartitionKindHash, PartitionKindMultiHash, PartitionKindRange}
}

// PrimaryAggregationType enumerates the values for primary aggregation type.
type PrimaryAggregationType string

const (
	// PrimaryAggregationTypeAverage ...
	PrimaryAggregationTypeAverage PrimaryAggregationType = "Average"
	// PrimaryAggregationTypeLast ...
	PrimaryAggregationTypeLast PrimaryAggregationType = "Last"
	// PrimaryAggregationTypeMaximum ...
	PrimaryAggregationTypeMaximum PrimaryAggregationType = "Maximum"
	// PrimaryAggregationTypeMinimum ...
	PrimaryAggregationTypeMinimum PrimaryAggregationType = "Minimum"
	// PrimaryAggregationTypeNone ...
	PrimaryAggregationTypeNone PrimaryAggregationType = "None"
	// PrimaryAggregationTypeTotal ...
	PrimaryAggregationTypeTotal PrimaryAggregationType = "Total"
)

// PossiblePrimaryAggregationTypeValues returns an array of possible values for the PrimaryAggregationType const type.
func PossiblePrimaryAggregationTypeValues() []PrimaryAggregationType {
	return []PrimaryAggregationType{PrimaryAggregationTypeAverage, PrimaryAggregationTypeLast, PrimaryAggregationTypeMaximum, PrimaryAggregationTypeMinimum, PrimaryAggregationTypeNone, PrimaryAggregationTypeTotal}
}

// PublicNetworkAccess enumerates the values for public network access.
type PublicNetworkAccess string

const (
	// Disabled ...
	Disabled PublicNetworkAccess = "Disabled"
	// Enabled ...
	Enabled PublicNetworkAccess = "Enabled"
)

// PossiblePublicNetworkAccessValues returns an array of possible values for the PublicNetworkAccess const type.
func PossiblePublicNetworkAccessValues() []PublicNetworkAccess {
	return []PublicNetworkAccess{Disabled, Enabled}
}

// ResourceIdentityType enumerates the values for resource identity type.
type ResourceIdentityType string

const (
	// ResourceIdentityTypeNone ...
	ResourceIdentityTypeNone ResourceIdentityType = "None"
	// ResourceIdentityTypeSystemAssigned ...
	ResourceIdentityTypeSystemAssigned ResourceIdentityType = "SystemAssigned"
	// ResourceIdentityTypeSystemAssignedUserAssigned ...
	ResourceIdentityTypeSystemAssignedUserAssigned ResourceIdentityType = "SystemAssigned,UserAssigned"
	// ResourceIdentityTypeUserAssigned ...
	ResourceIdentityTypeUserAssigned ResourceIdentityType = "UserAssigned"
)

// PossibleResourceIdentityTypeValues returns an array of possible values for the ResourceIdentityType const type.
func PossibleResourceIdentityTypeValues() []ResourceIdentityType {
	return []ResourceIdentityType{ResourceIdentityTypeNone, ResourceIdentityTypeSystemAssigned, ResourceIdentityTypeSystemAssignedUserAssigned, ResourceIdentityTypeUserAssigned}
}

// RestoreMode enumerates the values for restore mode.
type RestoreMode string

const (
	// PointInTime ...
	PointInTime RestoreMode = "PointInTime"
)

// PossibleRestoreModeValues returns an array of possible values for the RestoreMode const type.
func PossibleRestoreModeValues() []RestoreMode {
	return []RestoreMode{PointInTime}
}

// RoleDefinitionType enumerates the values for role definition type.
type RoleDefinitionType string

const (
	// RoleDefinitionTypeBuiltInRole ...
	RoleDefinitionTypeBuiltInRole RoleDefinitionType = "BuiltInRole"
	// RoleDefinitionTypeCustomRole ...
	RoleDefinitionTypeCustomRole RoleDefinitionType = "CustomRole"
)

// PossibleRoleDefinitionTypeValues returns an array of possible values for the RoleDefinitionType const type.
func PossibleRoleDefinitionTypeValues() []RoleDefinitionType {
	return []RoleDefinitionType{RoleDefinitionTypeBuiltInRole, RoleDefinitionTypeCustomRole}
}

// ServerVersion enumerates the values for server version.
type ServerVersion string

const (
	// FourFullStopTwo ...
	FourFullStopTwo ServerVersion = "4.2"
	// FourFullStopZero ...
	FourFullStopZero ServerVersion = "4.0"
	// ThreeFullStopSix ...
	ThreeFullStopSix ServerVersion = "3.6"
	// ThreeFullStopTwo ...
	ThreeFullStopTwo ServerVersion = "3.2"
)

// PossibleServerVersionValues returns an array of possible values for the ServerVersion const type.
func PossibleServerVersionValues() []ServerVersion {
	return []ServerVersion{FourFullStopTwo, FourFullStopZero, ThreeFullStopSix, ThreeFullStopTwo}
}

// ServiceSize enumerates the values for service size.
type ServiceSize string

const (
	// CosmosD16s ...
	CosmosD16s ServiceSize = "Cosmos.D16s"
	// CosmosD4s ...
	CosmosD4s ServiceSize = "Cosmos.D4s"
	// CosmosD8s ...
	CosmosD8s ServiceSize = "Cosmos.D8s"
)

// PossibleServiceSizeValues returns an array of possible values for the ServiceSize const type.
func PossibleServiceSizeValues() []ServiceSize {
	return []ServiceSize{CosmosD16s, CosmosD4s, CosmosD8s}
}

// ServiceStatus enumerates the values for service status.
type ServiceStatus string

const (
	// ServiceStatusCreating ...
	ServiceStatusCreating ServiceStatus = "Creating"
	// ServiceStatusDeleting ...
	ServiceStatusDeleting ServiceStatus = "Deleting"
	// ServiceStatusError ...
	ServiceStatusError ServiceStatus = "Error"
	// ServiceStatusRunning ...
	ServiceStatusRunning ServiceStatus = "Running"
	// ServiceStatusStopped ...
	ServiceStatusStopped ServiceStatus = "Stopped"
	// ServiceStatusUpdating ...
	ServiceStatusUpdating ServiceStatus = "Updating"
)

// PossibleServiceStatusValues returns an array of possible values for the ServiceStatus const type.
func PossibleServiceStatusValues() []ServiceStatus {
	return []ServiceStatus{ServiceStatusCreating, ServiceStatusDeleting, ServiceStatusError, ServiceStatusRunning, ServiceStatusStopped, ServiceStatusUpdating}
}

// ServiceType enumerates the values for service type.
type ServiceType string

const (
	// DataTransfer ...
	DataTransfer ServiceType = "DataTransfer"
	// GraphAPICompute ...
	GraphAPICompute ServiceType = "GraphAPICompute"
	// MaterializedViewsBuilder ...
	MaterializedViewsBuilder ServiceType = "MaterializedViewsBuilder"
	// SQLDedicatedGateway ...
	SQLDedicatedGateway ServiceType = "SqlDedicatedGateway"
)

// PossibleServiceTypeValues returns an array of possible values for the ServiceType const type.
func PossibleServiceTypeValues() []ServiceType {
	return []ServiceType{DataTransfer, GraphAPICompute, MaterializedViewsBuilder, SQLDedicatedGateway}
}

// ServiceTypeBasicServiceResourceProperties enumerates the values for service type basic service resource
// properties.
type ServiceTypeBasicServiceResourceProperties string

const (
	// ServiceTypeDataTransfer ...
	ServiceTypeDataTransfer ServiceTypeBasicServiceResourceProperties = "DataTransfer"
	// ServiceTypeGraphAPICompute ...
	ServiceTypeGraphAPICompute ServiceTypeBasicServiceResourceProperties = "GraphAPICompute"
	// ServiceTypeMaterializedViewsBuilder ...
	ServiceTypeMaterializedViewsBuilder ServiceTypeBasicServiceResourceProperties = "MaterializedViewsBuilder"
	// ServiceTypeServiceResourceProperties ...
	ServiceTypeServiceResourceProperties ServiceTypeBasicServiceResourceProperties = "ServiceResourceProperties"
	// ServiceTypeSQLDedicatedGateway ...
	ServiceTypeSQLDedicatedGateway ServiceTypeBasicServiceResourceProperties = "SqlDedicatedGateway"
)

// PossibleServiceTypeBasicServiceResourcePropertiesValues returns an array of possible values for the ServiceTypeBasicServiceResourceProperties const type.
func PossibleServiceTypeBasicServiceResourcePropertiesValues() []ServiceTypeBasicServiceResourceProperties {
	return []ServiceTypeBasicServiceResourceProperties{ServiceTypeDataTransfer, ServiceTypeGraphAPICompute, ServiceTypeMaterializedViewsBuilder, ServiceTypeServiceResourceProperties, ServiceTypeSQLDedicatedGateway}
}

// SpatialType enumerates the values for spatial type.
type SpatialType string

const (
	// SpatialTypeLineString ...
	SpatialTypeLineString SpatialType = "LineString"
	// SpatialTypeMultiPolygon ...
	SpatialTypeMultiPolygon SpatialType = "MultiPolygon"
	// SpatialTypePoint ...
	SpatialTypePoint SpatialType = "Point"
	// SpatialTypePolygon ...
	SpatialTypePolygon SpatialType = "Polygon"
)

// PossibleSpatialTypeValues returns an array of possible values for the SpatialType const type.
func PossibleSpatialTypeValues() []SpatialType {
	return []SpatialType{SpatialTypeLineString, SpatialTypeMultiPolygon, SpatialTypePoint, SpatialTypePolygon}
}

// TriggerOperation enumerates the values for trigger operation.
type TriggerOperation string

const (
	// TriggerOperationAll ...
	TriggerOperationAll TriggerOperation = "All"
	// TriggerOperationCreate ...
	TriggerOperationCreate TriggerOperation = "Create"
	// TriggerOperationDelete ...
	TriggerOperationDelete TriggerOperation = "Delete"
	// TriggerOperationReplace ...
	TriggerOperationReplace TriggerOperation = "Replace"
	// TriggerOperationUpdate ...
	TriggerOperationUpdate TriggerOperation = "Update"
)

// PossibleTriggerOperationValues returns an array of possible values for the TriggerOperation const type.
func PossibleTriggerOperationValues() []TriggerOperation {
	return []TriggerOperation{TriggerOperationAll, TriggerOperationCreate, TriggerOperationDelete, TriggerOperationReplace, TriggerOperationUpdate}
}

// TriggerType enumerates the values for trigger type.
type TriggerType string

const (
	// Post ...
	Post TriggerType = "Post"
	// Pre ...
	Pre TriggerType = "Pre"
)

// PossibleTriggerTypeValues returns an array of possible values for the TriggerType const type.
func PossibleTriggerTypeValues() []TriggerType {
	return []TriggerType{Post, Pre}
}

// Type enumerates the values for type.
type Type string

const (
	// TypeBackupPolicy ...
	TypeBackupPolicy Type = "BackupPolicy"
	// TypeContinuous ...
	TypeContinuous Type = "Continuous"
	// TypePeriodic ...
	TypePeriodic Type = "Periodic"
)

// PossibleTypeValues returns an array of possible values for the Type const type.
func PossibleTypeValues() []Type {
	return []Type{TypeBackupPolicy, TypeContinuous, TypePeriodic}
}

// UnitType enumerates the values for unit type.
type UnitType string

const (
	// Bytes ...
	Bytes UnitType = "Bytes"
	// BytesPerSecond ...
	BytesPerSecond UnitType = "BytesPerSecond"
	// Count ...
	Count UnitType = "Count"
	// CountPerSecond ...
	CountPerSecond UnitType = "CountPerSecond"
	// Milliseconds ...
	Milliseconds UnitType = "Milliseconds"
	// Percent ...
	Percent UnitType = "Percent"
	// Seconds ...
	Seconds UnitType = "Seconds"
)

// PossibleUnitTypeValues returns an array of possible values for the UnitType const type.
func PossibleUnitTypeValues() []UnitType {
	return []UnitType{Bytes, BytesPerSecond, Count, CountPerSecond, Milliseconds, Percent, Seconds}
}
