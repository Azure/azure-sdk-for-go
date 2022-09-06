//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armpostgresqlflexibleservers

import "time"

// Backup properties of a server
type Backup struct {
	// Backup retention days for the server.
	BackupRetentionDays *int32 `json:"backupRetentionDays,omitempty"`

	// A value indicating whether Geo-Redundant backup is enabled on the server.
	GeoRedundantBackup *GeoRedundantBackupEnum `json:"geoRedundantBackup,omitempty"`

	// READ-ONLY; The earliest restore point time (ISO8601 format) for server.
	EarliestRestoreDate *time.Time `json:"earliestRestoreDate,omitempty" azure:"ro"`
}

// CapabilitiesListResult - location capability
type CapabilitiesListResult struct {
	// READ-ONLY; Link to retrieve next page of results.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A list of supported capabilities.
	Value []*CapabilityProperties `json:"value,omitempty" azure:"ro"`
}

// CapabilityProperties - Location capabilities.
type CapabilityProperties struct {
	// READ-ONLY; A value indicating whether a new server in this region can have geo-backups to paired region.
	GeoBackupSupported *bool `json:"geoBackupSupported,omitempty" azure:"ro"`

	// READ-ONLY; The status
	Status *string `json:"status,omitempty" azure:"ro"`

	// READ-ONLY
	SupportedFlexibleServerEditions []*FlexibleServerEditionCapability `json:"supportedFlexibleServerEditions,omitempty" azure:"ro"`

	// READ-ONLY; Supported high availability mode
	SupportedHAMode []*string `json:"supportedHAMode,omitempty" azure:"ro"`

	// READ-ONLY
	SupportedHyperscaleNodeEditions []*HyperscaleNodeEditionCapability `json:"supportedHyperscaleNodeEditions,omitempty" azure:"ro"`

	// READ-ONLY; zone name
	Zone *string `json:"zone,omitempty" azure:"ro"`

	// READ-ONLY; A value indicating whether a new server in this region can have geo-backups to paired region.
	ZoneRedundantHaAndGeoBackupSupported *bool `json:"zoneRedundantHaAndGeoBackupSupported,omitempty" azure:"ro"`

	// READ-ONLY; A value indicating whether a new server in this region can support multi zone HA.
	ZoneRedundantHaSupported *bool `json:"zoneRedundantHaSupported,omitempty" azure:"ro"`
}

// CheckNameAvailabilityClientExecuteOptions contains the optional parameters for the CheckNameAvailabilityClient.Execute
// method.
type CheckNameAvailabilityClientExecuteOptions struct {
	// placeholder for future optional parameters
}

// CloudError - An error response from the Batch service.
type CloudError struct {
	// Common error response for all Azure Resource Manager APIs to return error details for failed operations. (This also follows
	// the OData error response format.)
	Error *ErrorResponse `json:"error,omitempty"`
}

// Configuration - Represents a Configuration.
type Configuration struct {
	// The properties of a configuration.
	Properties *ConfigurationProperties `json:"properties,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The system metadata relating to this resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// ConfigurationListResult - A list of server configurations.
type ConfigurationListResult struct {
	// The link used to get the next page of operations.
	NextLink *string `json:"nextLink,omitempty"`

	// The list of server configurations.
	Value []*Configuration `json:"value,omitempty"`
}

// ConfigurationProperties - The properties of a configuration.
type ConfigurationProperties struct {
	// Source of the configuration.
	Source *string `json:"source,omitempty"`

	// Value of the configuration.
	Value *string `json:"value,omitempty"`

	// READ-ONLY; Allowed values of the configuration.
	AllowedValues *string `json:"allowedValues,omitempty" azure:"ro"`

	// READ-ONLY; Data type of the configuration.
	DataType *ConfigurationDataType `json:"dataType,omitempty" azure:"ro"`

	// READ-ONLY; Default value of the configuration.
	DefaultValue *string `json:"defaultValue,omitempty" azure:"ro"`

	// READ-ONLY; Description of the configuration.
	Description *string `json:"description,omitempty" azure:"ro"`

	// READ-ONLY; Configuration documentation link.
	DocumentationLink *string `json:"documentationLink,omitempty" azure:"ro"`

	// READ-ONLY; Configuration is pending restart or not.
	IsConfigPendingRestart *bool `json:"isConfigPendingRestart,omitempty" azure:"ro"`

	// READ-ONLY; Configuration dynamic or static.
	IsDynamicConfig *bool `json:"isDynamicConfig,omitempty" azure:"ro"`

	// READ-ONLY; Configuration read-only or not.
	IsReadOnly *bool `json:"isReadOnly,omitempty" azure:"ro"`

	// READ-ONLY; Configuration unit.
	Unit *string `json:"unit,omitempty" azure:"ro"`
}

// ConfigurationsClientBeginPutOptions contains the optional parameters for the ConfigurationsClient.BeginPut method.
type ConfigurationsClientBeginPutOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ConfigurationsClientBeginUpdateOptions contains the optional parameters for the ConfigurationsClient.BeginUpdate method.
type ConfigurationsClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ConfigurationsClientGetOptions contains the optional parameters for the ConfigurationsClient.Get method.
type ConfigurationsClientGetOptions struct {
	// placeholder for future optional parameters
}

// ConfigurationsClientListByServerOptions contains the optional parameters for the ConfigurationsClient.ListByServer method.
type ConfigurationsClientListByServerOptions struct {
	// placeholder for future optional parameters
}

// Database - Represents a Database.
type Database struct {
	// The properties of a database.
	Properties *DatabaseProperties `json:"properties,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The system metadata relating to this resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// DatabaseListResult - A List of databases.
type DatabaseListResult struct {
	// The link used to get the next page of databases.
	NextLink *string `json:"nextLink,omitempty"`

	// The list of databases housed in a server
	Value []*Database `json:"value,omitempty"`
}

// DatabaseProperties - The properties of a database.
type DatabaseProperties struct {
	// The charset of the database.
	Charset *string `json:"charset,omitempty"`

	// The collation of the database.
	Collation *string `json:"collation,omitempty"`
}

// DatabasesClientBeginCreateOptions contains the optional parameters for the DatabasesClient.BeginCreate method.
type DatabasesClientBeginCreateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// DatabasesClientBeginDeleteOptions contains the optional parameters for the DatabasesClient.BeginDelete method.
type DatabasesClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// DatabasesClientGetOptions contains the optional parameters for the DatabasesClient.Get method.
type DatabasesClientGetOptions struct {
	// placeholder for future optional parameters
}

// DatabasesClientListByServerOptions contains the optional parameters for the DatabasesClient.ListByServer method.
type DatabasesClientListByServerOptions struct {
	// placeholder for future optional parameters
}

// DelegatedSubnetUsage - Delegated subnet usage data.
type DelegatedSubnetUsage struct {
	// READ-ONLY; name of the subnet
	SubnetName *string `json:"subnetName,omitempty" azure:"ro"`

	// READ-ONLY; Number of used delegated subnets
	Usage *int64 `json:"usage,omitempty" azure:"ro"`
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info interface{} `json:"info,omitempty" azure:"ro"`

	// READ-ONLY; The additional info type.
	Type *string `json:"type,omitempty" azure:"ro"`
}

// ErrorResponse - Common error response for all Azure Resource Manager APIs to return error details for failed operations.
// (This also follows the OData error response format.)
type ErrorResponse struct {
	// READ-ONLY; The error additional info.
	AdditionalInfo []*ErrorAdditionalInfo `json:"additionalInfo,omitempty" azure:"ro"`

	// READ-ONLY; The error code.
	Code *string `json:"code,omitempty" azure:"ro"`

	// READ-ONLY; The error details.
	Details []*ErrorResponse `json:"details,omitempty" azure:"ro"`

	// READ-ONLY; The error message.
	Message *string `json:"message,omitempty" azure:"ro"`

	// READ-ONLY; The error target.
	Target *string `json:"target,omitempty" azure:"ro"`
}

// FirewallRule - Represents a server firewall rule.
type FirewallRule struct {
	// REQUIRED; The properties of a firewall rule.
	Properties *FirewallRuleProperties `json:"properties,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The system metadata relating to this resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// FirewallRuleListResult - A list of firewall rules.
type FirewallRuleListResult struct {
	// The link used to get the next page of operations.
	NextLink *string `json:"nextLink,omitempty"`

	// The list of firewall rules in a server.
	Value []*FirewallRule `json:"value,omitempty"`
}

// FirewallRuleProperties - The properties of a server firewall rule.
type FirewallRuleProperties struct {
	// REQUIRED; The end IP address of the server firewall rule. Must be IPv4 format.
	EndIPAddress *string `json:"endIpAddress,omitempty"`

	// REQUIRED; The start IP address of the server firewall rule. Must be IPv4 format.
	StartIPAddress *string `json:"startIpAddress,omitempty"`
}

// FirewallRulesClientBeginCreateOrUpdateOptions contains the optional parameters for the FirewallRulesClient.BeginCreateOrUpdate
// method.
type FirewallRulesClientBeginCreateOrUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// FirewallRulesClientBeginDeleteOptions contains the optional parameters for the FirewallRulesClient.BeginDelete method.
type FirewallRulesClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// FirewallRulesClientGetOptions contains the optional parameters for the FirewallRulesClient.Get method.
type FirewallRulesClientGetOptions struct {
	// placeholder for future optional parameters
}

// FirewallRulesClientListByServerOptions contains the optional parameters for the FirewallRulesClient.ListByServer method.
type FirewallRulesClientListByServerOptions struct {
	// placeholder for future optional parameters
}

// FlexibleServerEditionCapability - Flexible server edition capabilities.
type FlexibleServerEditionCapability struct {
	// READ-ONLY; Server edition name
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The status
	Status *string `json:"status,omitempty" azure:"ro"`

	// READ-ONLY; The list of server versions supported by this server edition.
	SupportedServerVersions []*ServerVersionCapability `json:"supportedServerVersions,omitempty" azure:"ro"`

	// READ-ONLY; The list of editions supported by this server edition.
	SupportedStorageEditions []*StorageEditionCapability `json:"supportedStorageEditions,omitempty" azure:"ro"`
}

// GetPrivateDNSZoneSuffixClientExecuteOptions contains the optional parameters for the GetPrivateDNSZoneSuffixClient.Execute
// method.
type GetPrivateDNSZoneSuffixClientExecuteOptions struct {
	// placeholder for future optional parameters
}

// HighAvailability - High availability properties of a server
type HighAvailability struct {
	// The HA mode for the server.
	Mode *HighAvailabilityMode `json:"mode,omitempty"`

	// availability zone information of the standby.
	StandbyAvailabilityZone *string `json:"standbyAvailabilityZone,omitempty"`

	// READ-ONLY; A state of a HA server that is visible to user.
	State *ServerHAState `json:"state,omitempty" azure:"ro"`
}

// HyperscaleNodeEditionCapability - Hyperscale node edition capabilities.
type HyperscaleNodeEditionCapability struct {
	// READ-ONLY; Server edition name
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The status
	Status *string `json:"status,omitempty" azure:"ro"`

	// READ-ONLY; The list of Node Types supported by this server edition.
	SupportedNodeTypes []*NodeTypeCapability `json:"supportedNodeTypes,omitempty" azure:"ro"`

	// READ-ONLY; The list of server versions supported by this server edition.
	SupportedServerVersions []*ServerVersionCapability `json:"supportedServerVersions,omitempty" azure:"ro"`

	// READ-ONLY; The list of editions supported by this server edition.
	SupportedStorageEditions []*StorageEditionCapability `json:"supportedStorageEditions,omitempty" azure:"ro"`
}

// LocationBasedCapabilitiesClientExecuteOptions contains the optional parameters for the LocationBasedCapabilitiesClient.Execute
// method.
type LocationBasedCapabilitiesClientExecuteOptions struct {
	// placeholder for future optional parameters
}

// MaintenanceWindow - Maintenance window properties of a server.
type MaintenanceWindow struct {
	// indicates whether custom window is enabled or disabled
	CustomWindow *string `json:"customWindow,omitempty"`

	// day of week for maintenance window
	DayOfWeek *int32 `json:"dayOfWeek,omitempty"`

	// start hour for maintenance window
	StartHour *int32 `json:"startHour,omitempty"`

	// start minute for maintenance window
	StartMinute *int32 `json:"startMinute,omitempty"`
}

// NameAvailability - Represents a resource name availability.
type NameAvailability struct {
	// READ-ONLY; Error Message.
	Message *string `json:"message,omitempty" azure:"ro"`

	// READ-ONLY; name of the PostgreSQL server.
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; Indicates whether the resource name is available.
	NameAvailable *bool `json:"nameAvailable,omitempty" azure:"ro"`

	// READ-ONLY; The name availability reason.
	Reason *Reason `json:"reason,omitempty" azure:"ro"`

	// READ-ONLY; type of the server
	Type *string `json:"type,omitempty" azure:"ro"`
}

// NameAvailabilityRequest - Request from client to check resource name availability.
type NameAvailabilityRequest struct {
	// REQUIRED; Resource name to verify.
	Name *string `json:"name,omitempty"`

	// Resource type used for verification.
	Type *string `json:"type,omitempty"`
}

// Network properties of a server
type Network struct {
	// delegated subnet arm resource id.
	DelegatedSubnetResourceID *string `json:"delegatedSubnetResourceId,omitempty"`

	// private dns zone arm resource id.
	PrivateDNSZoneArmResourceID *string `json:"privateDnsZoneArmResourceId,omitempty"`

	// READ-ONLY; public network access is enabled or not
	PublicNetworkAccess *ServerPublicNetworkAccessState `json:"publicNetworkAccess,omitempty" azure:"ro"`
}

// NodeTypeCapability - node type capability
type NodeTypeCapability struct {
	// READ-ONLY; note type name
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; note type
	NodeType *string `json:"nodeType,omitempty" azure:"ro"`

	// READ-ONLY; The status
	Status *string `json:"status,omitempty" azure:"ro"`
}

// Operation - REST API operation definition.
type Operation struct {
	// Indicates whether the operation is a data action
	IsDataAction *bool `json:"isDataAction,omitempty"`

	// READ-ONLY; The localized display information for this particular operation or action.
	Display *OperationDisplay `json:"display,omitempty" azure:"ro"`

	// READ-ONLY; The name of the operation being performed on this particular object.
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The intended executor of the operation.
	Origin *OperationOrigin `json:"origin,omitempty" azure:"ro"`

	// READ-ONLY; Additional descriptions for the operation.
	Properties map[string]interface{} `json:"properties,omitempty" azure:"ro"`
}

// OperationDisplay - Display metadata associated with the operation.
type OperationDisplay struct {
	// READ-ONLY; Operation description.
	Description *string `json:"description,omitempty" azure:"ro"`

	// READ-ONLY; Localized friendly name for the operation.
	Operation *string `json:"operation,omitempty" azure:"ro"`

	// READ-ONLY; Operation resource provider name.
	Provider *string `json:"provider,omitempty" azure:"ro"`

	// READ-ONLY; Resource on which the operation is performed.
	Resource *string `json:"resource,omitempty" azure:"ro"`
}

// OperationListResult - A list of resource provider operations.
type OperationListResult struct {
	// URL client should use to fetch the next page (per server side paging). It's null for now, added for future use.
	NextLink *string `json:"nextLink,omitempty"`

	// Collection of available operation details
	Value []*Operation `json:"value,omitempty"`
}

// OperationsClientListOptions contains the optional parameters for the OperationsClient.List method.
type OperationsClientListOptions struct {
	// placeholder for future optional parameters
}

// ProxyResource - The resource model definition for a Azure Resource Manager proxy resource. It will not have tags and a
// location
type ProxyResource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// Resource - Common fields that are returned in the response for all Azure Resource Manager resources
type Resource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// RestartParameter - Represents server restart parameters.
type RestartParameter struct {
	// Failover mode.
	FailoverMode *FailoverMode `json:"failoverMode,omitempty"`

	// Indicates whether to restart the server with failover.
	RestartWithFailover *bool `json:"restartWithFailover,omitempty"`
}

// SKU - Sku information related properties of a server.
type SKU struct {
	// REQUIRED; The name of the sku, typically, tier + family + cores, e.g. StandardD4sv3.
	Name *string `json:"name,omitempty"`

	// REQUIRED; The tier of the particular SKU, e.g. Burstable.
	Tier *SKUTier `json:"tier,omitempty"`
}

// Server - Represents a server.
type Server struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// Properties of the server.
	Properties *ServerProperties `json:"properties,omitempty"`

	// The SKU (pricing tier) of the server.
	SKU *SKU `json:"sku,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The system metadata relating to this resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// ServerForUpdate - Represents a server to be updated.
type ServerForUpdate struct {
	// The location the resource resides in.
	Location *string `json:"location,omitempty"`

	// Properties of the server.
	Properties *ServerPropertiesForUpdate `json:"properties,omitempty"`

	// The SKU (pricing tier) of the server.
	SKU *SKU `json:"sku,omitempty"`

	// Application-specific metadata in the form of key-value pairs.
	Tags map[string]*string `json:"tags,omitempty"`
}

// ServerListResult - A list of servers.
type ServerListResult struct {
	// The link used to get the next page of operations.
	NextLink *string `json:"nextLink,omitempty"`

	// The list of flexible servers
	Value []*Server `json:"value,omitempty"`
}

// ServerProperties - The properties of a server.
type ServerProperties struct {
	// The administrator's login name of a server. Can only be specified when the server is being created (and is required for
	// creation).
	AdministratorLogin *string `json:"administratorLogin,omitempty"`

	// The administrator login password (required for server creation).
	AdministratorLoginPassword *string `json:"administratorLoginPassword,omitempty"`

	// availability zone information of the server.
	AvailabilityZone *string `json:"availabilityZone,omitempty"`

	// Backup properties of a server.
	Backup *Backup `json:"backup,omitempty"`

	// The mode to create a new PostgreSQL server.
	CreateMode *CreateMode `json:"createMode,omitempty"`

	// High availability properties of a server.
	HighAvailability *HighAvailability `json:"highAvailability,omitempty"`

	// Maintenance window properties of a server.
	MaintenanceWindow *MaintenanceWindow `json:"maintenanceWindow,omitempty"`

	// Network properties of a server.
	Network *Network `json:"network,omitempty"`

	// Restore point creation time (ISO8601 format), specifying the time to restore from. It's required when 'createMode' is 'PointInTimeRestore'.
	PointInTimeUTC *time.Time `json:"pointInTimeUTC,omitempty"`

	// The source server resource ID to restore from. It's required when 'createMode' is 'PointInTimeRestore'.
	SourceServerResourceID *string `json:"sourceServerResourceId,omitempty"`

	// Storage properties of a server.
	Storage *Storage `json:"storage,omitempty"`

	// PostgreSQL Server version.
	Version *ServerVersion `json:"version,omitempty"`

	// READ-ONLY; The fully qualified domain name of a server.
	FullyQualifiedDomainName *string `json:"fullyQualifiedDomainName,omitempty" azure:"ro"`

	// READ-ONLY; The minor version of the server.
	MinorVersion *string `json:"minorVersion,omitempty" azure:"ro"`

	// READ-ONLY; A state of a server that is visible to user.
	State *ServerState `json:"state,omitempty" azure:"ro"`
}

type ServerPropertiesForUpdate struct {
	// The password of the administrator login.
	AdministratorLoginPassword *string `json:"administratorLoginPassword,omitempty"`

	// Backup properties of a server.
	Backup *Backup `json:"backup,omitempty"`

	// The mode to update a new PostgreSQL server.
	CreateMode *CreateModeForUpdate `json:"createMode,omitempty"`

	// High availability properties of a server.
	HighAvailability *HighAvailability `json:"highAvailability,omitempty"`

	// Maintenance window properties of a server.
	MaintenanceWindow *MaintenanceWindow `json:"maintenanceWindow,omitempty"`

	// Storage properties of a server.
	Storage *Storage `json:"storage,omitempty"`
}

// ServerVersionCapability - Server version capabilities.
type ServerVersionCapability struct {
	// READ-ONLY; server version
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The status
	Status *string `json:"status,omitempty" azure:"ro"`

	// READ-ONLY
	SupportedVcores []*VcoreCapability `json:"supportedVcores,omitempty" azure:"ro"`
}

// ServersClientBeginCreateOptions contains the optional parameters for the ServersClient.BeginCreate method.
type ServersClientBeginCreateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ServersClientBeginDeleteOptions contains the optional parameters for the ServersClient.BeginDelete method.
type ServersClientBeginDeleteOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ServersClientBeginRestartOptions contains the optional parameters for the ServersClient.BeginRestart method.
type ServersClientBeginRestartOptions struct {
	// The parameters for restarting a server.
	Parameters *RestartParameter
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ServersClientBeginStartOptions contains the optional parameters for the ServersClient.BeginStart method.
type ServersClientBeginStartOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ServersClientBeginStopOptions contains the optional parameters for the ServersClient.BeginStop method.
type ServersClientBeginStopOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ServersClientBeginUpdateOptions contains the optional parameters for the ServersClient.BeginUpdate method.
type ServersClientBeginUpdateOptions struct {
	// Resumes the LRO from the provided token.
	ResumeToken string
}

// ServersClientGetOptions contains the optional parameters for the ServersClient.Get method.
type ServersClientGetOptions struct {
	// placeholder for future optional parameters
}

// ServersClientListByResourceGroupOptions contains the optional parameters for the ServersClient.ListByResourceGroup method.
type ServersClientListByResourceGroupOptions struct {
	// placeholder for future optional parameters
}

// ServersClientListOptions contains the optional parameters for the ServersClient.List method.
type ServersClientListOptions struct {
	// placeholder for future optional parameters
}

// Storage properties of a server
type Storage struct {
	// Max storage allowed for a server.
	StorageSizeGB *int32 `json:"storageSizeGB,omitempty"`
}

// StorageEditionCapability - storage edition capability
type StorageEditionCapability struct {
	// READ-ONLY; storage edition name
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The status
	Status *string `json:"status,omitempty" azure:"ro"`

	// READ-ONLY
	SupportedStorageMB []*StorageMBCapability `json:"supportedStorageMB,omitempty" azure:"ro"`
}

// StorageMBCapability - storage size in MB capability
type StorageMBCapability struct {
	// READ-ONLY; storage MB name
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The status
	Status *string `json:"status,omitempty" azure:"ro"`

	// READ-ONLY; storage size in MB
	StorageSizeMB *int64 `json:"storageSizeMB,omitempty" azure:"ro"`

	// READ-ONLY; supported IOPS
	SupportedIops *int64 `json:"supportedIops,omitempty" azure:"ro"`
}

// SystemData - Metadata pertaining to creation and last modification of the resource.
type SystemData struct {
	// The timestamp of resource creation (UTC).
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	// The identity that created the resource.
	CreatedBy *string `json:"createdBy,omitempty"`

	// The type of identity that created the resource.
	CreatedByType *CreatedByType `json:"createdByType,omitempty"`

	// The timestamp of resource last modification (UTC)
	LastModifiedAt *time.Time `json:"lastModifiedAt,omitempty"`

	// The identity that last modified the resource.
	LastModifiedBy *string `json:"lastModifiedBy,omitempty"`

	// The type of identity that last modified the resource.
	LastModifiedByType *CreatedByType `json:"lastModifiedByType,omitempty"`
}

// TrackedResource - The resource model definition for an Azure Resource Manager tracked top level resource which has 'tags'
// and a 'location'
type TrackedResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// VcoreCapability - Vcores capability
type VcoreCapability struct {
	// READ-ONLY; vCore name
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The status
	Status *string `json:"status,omitempty" azure:"ro"`

	// READ-ONLY; supported IOPS
	SupportedIops *int64 `json:"supportedIops,omitempty" azure:"ro"`

	// READ-ONLY; supported memory per vCore in MB
	SupportedMemoryPerVcoreMB *int64 `json:"supportedMemoryPerVcoreMB,omitempty" azure:"ro"`

	// READ-ONLY; supported vCores
	VCores *int64 `json:"vCores,omitempty" azure:"ro"`
}

// VirtualNetworkSubnetUsageClientExecuteOptions contains the optional parameters for the VirtualNetworkSubnetUsageClient.Execute
// method.
type VirtualNetworkSubnetUsageClientExecuteOptions struct {
	// placeholder for future optional parameters
}

// VirtualNetworkSubnetUsageParameter - Virtual network subnet usage parameter
type VirtualNetworkSubnetUsageParameter struct {
	// Virtual network resource id.
	VirtualNetworkArmResourceID *string `json:"virtualNetworkArmResourceId,omitempty"`
}

// VirtualNetworkSubnetUsageResult - Virtual network subnet usage data.
type VirtualNetworkSubnetUsageResult struct {
	// READ-ONLY
	DelegatedSubnetsUsage []*DelegatedSubnetUsage `json:"delegatedSubnetsUsage,omitempty" azure:"ro"`

	// READ-ONLY; The location the resource resides in.
	Location *string `json:"location,omitempty" azure:"ro"`

	// READ-ONLY; The subscription ID.
	SubscriptionID *string `json:"subscriptionId,omitempty" azure:"ro"`
}
