//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armappplatform

const (
	moduleName    = "armappplatform"
	moduleVersion = "v0.3.0"
)

// APIPortalProvisioningState - State of the API portal.
type APIPortalProvisioningState string

const (
	APIPortalProvisioningStateCreating  APIPortalProvisioningState = "Creating"
	APIPortalProvisioningStateDeleting  APIPortalProvisioningState = "Deleting"
	APIPortalProvisioningStateFailed    APIPortalProvisioningState = "Failed"
	APIPortalProvisioningStateSucceeded APIPortalProvisioningState = "Succeeded"
	APIPortalProvisioningStateUpdating  APIPortalProvisioningState = "Updating"
)

// PossibleAPIPortalProvisioningStateValues returns the possible values for the APIPortalProvisioningState const type.
func PossibleAPIPortalProvisioningStateValues() []APIPortalProvisioningState {
	return []APIPortalProvisioningState{
		APIPortalProvisioningStateCreating,
		APIPortalProvisioningStateDeleting,
		APIPortalProvisioningStateFailed,
		APIPortalProvisioningStateSucceeded,
		APIPortalProvisioningStateUpdating,
	}
}

// ToPtr returns a *APIPortalProvisioningState pointing to the current value.
func (c APIPortalProvisioningState) ToPtr() *APIPortalProvisioningState {
	return &c
}

// AppResourceProvisioningState - Provisioning state of the App
type AppResourceProvisioningState string

const (
	AppResourceProvisioningStateCreating  AppResourceProvisioningState = "Creating"
	AppResourceProvisioningStateDeleting  AppResourceProvisioningState = "Deleting"
	AppResourceProvisioningStateFailed    AppResourceProvisioningState = "Failed"
	AppResourceProvisioningStateSucceeded AppResourceProvisioningState = "Succeeded"
	AppResourceProvisioningStateUpdating  AppResourceProvisioningState = "Updating"
)

// PossibleAppResourceProvisioningStateValues returns the possible values for the AppResourceProvisioningState const type.
func PossibleAppResourceProvisioningStateValues() []AppResourceProvisioningState {
	return []AppResourceProvisioningState{
		AppResourceProvisioningStateCreating,
		AppResourceProvisioningStateDeleting,
		AppResourceProvisioningStateFailed,
		AppResourceProvisioningStateSucceeded,
		AppResourceProvisioningStateUpdating,
	}
}

// ToPtr returns a *AppResourceProvisioningState pointing to the current value.
func (c AppResourceProvisioningState) ToPtr() *AppResourceProvisioningState {
	return &c
}

// BindingType - Buildpack Binding Type
type BindingType string

const (
	BindingTypeApacheSkyWalking    BindingType = "ApacheSkyWalking"
	BindingTypeAppDynamics         BindingType = "AppDynamics"
	BindingTypeApplicationInsights BindingType = "ApplicationInsights"
	BindingTypeDynatrace           BindingType = "Dynatrace"
	BindingTypeElasticAPM          BindingType = "ElasticAPM"
	BindingTypeNewRelic            BindingType = "NewRelic"
)

// PossibleBindingTypeValues returns the possible values for the BindingType const type.
func PossibleBindingTypeValues() []BindingType {
	return []BindingType{
		BindingTypeApacheSkyWalking,
		BindingTypeAppDynamics,
		BindingTypeApplicationInsights,
		BindingTypeDynatrace,
		BindingTypeElasticAPM,
		BindingTypeNewRelic,
	}
}

// ToPtr returns a *BindingType pointing to the current value.
func (c BindingType) ToPtr() *BindingType {
	return &c
}

// BuildProvisioningState - Provisioning state of the KPack build result
type BuildProvisioningState string

const (
	BuildProvisioningStateCreating  BuildProvisioningState = "Creating"
	BuildProvisioningStateDeleting  BuildProvisioningState = "Deleting"
	BuildProvisioningStateFailed    BuildProvisioningState = "Failed"
	BuildProvisioningStateSucceeded BuildProvisioningState = "Succeeded"
	BuildProvisioningStateUpdating  BuildProvisioningState = "Updating"
)

// PossibleBuildProvisioningStateValues returns the possible values for the BuildProvisioningState const type.
func PossibleBuildProvisioningStateValues() []BuildProvisioningState {
	return []BuildProvisioningState{
		BuildProvisioningStateCreating,
		BuildProvisioningStateDeleting,
		BuildProvisioningStateFailed,
		BuildProvisioningStateSucceeded,
		BuildProvisioningStateUpdating,
	}
}

// ToPtr returns a *BuildProvisioningState pointing to the current value.
func (c BuildProvisioningState) ToPtr() *BuildProvisioningState {
	return &c
}

// BuildResultProvisioningState - Provisioning state of the KPack build result
type BuildResultProvisioningState string

const (
	BuildResultProvisioningStateBuilding  BuildResultProvisioningState = "Building"
	BuildResultProvisioningStateDeleting  BuildResultProvisioningState = "Deleting"
	BuildResultProvisioningStateFailed    BuildResultProvisioningState = "Failed"
	BuildResultProvisioningStateQueuing   BuildResultProvisioningState = "Queuing"
	BuildResultProvisioningStateSucceeded BuildResultProvisioningState = "Succeeded"
)

// PossibleBuildResultProvisioningStateValues returns the possible values for the BuildResultProvisioningState const type.
func PossibleBuildResultProvisioningStateValues() []BuildResultProvisioningState {
	return []BuildResultProvisioningState{
		BuildResultProvisioningStateBuilding,
		BuildResultProvisioningStateDeleting,
		BuildResultProvisioningStateFailed,
		BuildResultProvisioningStateQueuing,
		BuildResultProvisioningStateSucceeded,
	}
}

// ToPtr returns a *BuildResultProvisioningState pointing to the current value.
func (c BuildResultProvisioningState) ToPtr() *BuildResultProvisioningState {
	return &c
}

// BuildServiceProvisioningState - Provisioning state of the KPack build result
type BuildServiceProvisioningState string

const (
	BuildServiceProvisioningStateCreating  BuildServiceProvisioningState = "Creating"
	BuildServiceProvisioningStateDeleting  BuildServiceProvisioningState = "Deleting"
	BuildServiceProvisioningStateFailed    BuildServiceProvisioningState = "Failed"
	BuildServiceProvisioningStateSucceeded BuildServiceProvisioningState = "Succeeded"
	BuildServiceProvisioningStateUpdating  BuildServiceProvisioningState = "Updating"
)

// PossibleBuildServiceProvisioningStateValues returns the possible values for the BuildServiceProvisioningState const type.
func PossibleBuildServiceProvisioningStateValues() []BuildServiceProvisioningState {
	return []BuildServiceProvisioningState{
		BuildServiceProvisioningStateCreating,
		BuildServiceProvisioningStateDeleting,
		BuildServiceProvisioningStateFailed,
		BuildServiceProvisioningStateSucceeded,
		BuildServiceProvisioningStateUpdating,
	}
}

// ToPtr returns a *BuildServiceProvisioningState pointing to the current value.
func (c BuildServiceProvisioningState) ToPtr() *BuildServiceProvisioningState {
	return &c
}

// BuilderProvisioningState - Builder provision status.
type BuilderProvisioningState string

const (
	BuilderProvisioningStateCreating  BuilderProvisioningState = "Creating"
	BuilderProvisioningStateDeleting  BuilderProvisioningState = "Deleting"
	BuilderProvisioningStateFailed    BuilderProvisioningState = "Failed"
	BuilderProvisioningStateSucceeded BuilderProvisioningState = "Succeeded"
	BuilderProvisioningStateUpdating  BuilderProvisioningState = "Updating"
)

// PossibleBuilderProvisioningStateValues returns the possible values for the BuilderProvisioningState const type.
func PossibleBuilderProvisioningStateValues() []BuilderProvisioningState {
	return []BuilderProvisioningState{
		BuilderProvisioningStateCreating,
		BuilderProvisioningStateDeleting,
		BuilderProvisioningStateFailed,
		BuilderProvisioningStateSucceeded,
		BuilderProvisioningStateUpdating,
	}
}

// ToPtr returns a *BuilderProvisioningState pointing to the current value.
func (c BuilderProvisioningState) ToPtr() *BuilderProvisioningState {
	return &c
}

// BuildpackBindingProvisioningState - State of the Buildpack Binding.
type BuildpackBindingProvisioningState string

const (
	BuildpackBindingProvisioningStateCreating  BuildpackBindingProvisioningState = "Creating"
	BuildpackBindingProvisioningStateDeleting  BuildpackBindingProvisioningState = "Deleting"
	BuildpackBindingProvisioningStateFailed    BuildpackBindingProvisioningState = "Failed"
	BuildpackBindingProvisioningStateSucceeded BuildpackBindingProvisioningState = "Succeeded"
	BuildpackBindingProvisioningStateUpdating  BuildpackBindingProvisioningState = "Updating"
)

// PossibleBuildpackBindingProvisioningStateValues returns the possible values for the BuildpackBindingProvisioningState const type.
func PossibleBuildpackBindingProvisioningStateValues() []BuildpackBindingProvisioningState {
	return []BuildpackBindingProvisioningState{
		BuildpackBindingProvisioningStateCreating,
		BuildpackBindingProvisioningStateDeleting,
		BuildpackBindingProvisioningStateFailed,
		BuildpackBindingProvisioningStateSucceeded,
		BuildpackBindingProvisioningStateUpdating,
	}
}

// ToPtr returns a *BuildpackBindingProvisioningState pointing to the current value.
func (c BuildpackBindingProvisioningState) ToPtr() *BuildpackBindingProvisioningState {
	return &c
}

// ConfigServerState - State of the config server.
type ConfigServerState string

const (
	ConfigServerStateDeleted      ConfigServerState = "Deleted"
	ConfigServerStateFailed       ConfigServerState = "Failed"
	ConfigServerStateNotAvailable ConfigServerState = "NotAvailable"
	ConfigServerStateSucceeded    ConfigServerState = "Succeeded"
	ConfigServerStateUpdating     ConfigServerState = "Updating"
)

// PossibleConfigServerStateValues returns the possible values for the ConfigServerState const type.
func PossibleConfigServerStateValues() []ConfigServerState {
	return []ConfigServerState{
		ConfigServerStateDeleted,
		ConfigServerStateFailed,
		ConfigServerStateNotAvailable,
		ConfigServerStateSucceeded,
		ConfigServerStateUpdating,
	}
}

// ToPtr returns a *ConfigServerState pointing to the current value.
func (c ConfigServerState) ToPtr() *ConfigServerState {
	return &c
}

// ConfigurationServiceProvisioningState - State of the Application Configuration Service.
type ConfigurationServiceProvisioningState string

const (
	ConfigurationServiceProvisioningStateCreating  ConfigurationServiceProvisioningState = "Creating"
	ConfigurationServiceProvisioningStateDeleting  ConfigurationServiceProvisioningState = "Deleting"
	ConfigurationServiceProvisioningStateFailed    ConfigurationServiceProvisioningState = "Failed"
	ConfigurationServiceProvisioningStateSucceeded ConfigurationServiceProvisioningState = "Succeeded"
	ConfigurationServiceProvisioningStateUpdating  ConfigurationServiceProvisioningState = "Updating"
)

// PossibleConfigurationServiceProvisioningStateValues returns the possible values for the ConfigurationServiceProvisioningState const type.
func PossibleConfigurationServiceProvisioningStateValues() []ConfigurationServiceProvisioningState {
	return []ConfigurationServiceProvisioningState{
		ConfigurationServiceProvisioningStateCreating,
		ConfigurationServiceProvisioningStateDeleting,
		ConfigurationServiceProvisioningStateFailed,
		ConfigurationServiceProvisioningStateSucceeded,
		ConfigurationServiceProvisioningStateUpdating,
	}
}

// ToPtr returns a *ConfigurationServiceProvisioningState pointing to the current value.
func (c ConfigurationServiceProvisioningState) ToPtr() *ConfigurationServiceProvisioningState {
	return &c
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

// ToPtr returns a *CreatedByType pointing to the current value.
func (c CreatedByType) ToPtr() *CreatedByType {
	return &c
}

// CustomPersistentDiskPropertiesType - The type of the underlying resource to mount as a persistent disk.
type CustomPersistentDiskPropertiesType string

const (
	CustomPersistentDiskPropertiesTypeAzureFileVolume CustomPersistentDiskPropertiesType = "AzureFileVolume"
)

// PossibleCustomPersistentDiskPropertiesTypeValues returns the possible values for the CustomPersistentDiskPropertiesType const type.
func PossibleCustomPersistentDiskPropertiesTypeValues() []CustomPersistentDiskPropertiesType {
	return []CustomPersistentDiskPropertiesType{
		CustomPersistentDiskPropertiesTypeAzureFileVolume,
	}
}

// ToPtr returns a *CustomPersistentDiskPropertiesType pointing to the current value.
func (c CustomPersistentDiskPropertiesType) ToPtr() *CustomPersistentDiskPropertiesType {
	return &c
}

// DeploymentResourceProvisioningState - Provisioning state of the Deployment
type DeploymentResourceProvisioningState string

const (
	DeploymentResourceProvisioningStateCreating  DeploymentResourceProvisioningState = "Creating"
	DeploymentResourceProvisioningStateFailed    DeploymentResourceProvisioningState = "Failed"
	DeploymentResourceProvisioningStateSucceeded DeploymentResourceProvisioningState = "Succeeded"
	DeploymentResourceProvisioningStateUpdating  DeploymentResourceProvisioningState = "Updating"
)

// PossibleDeploymentResourceProvisioningStateValues returns the possible values for the DeploymentResourceProvisioningState const type.
func PossibleDeploymentResourceProvisioningStateValues() []DeploymentResourceProvisioningState {
	return []DeploymentResourceProvisioningState{
		DeploymentResourceProvisioningStateCreating,
		DeploymentResourceProvisioningStateFailed,
		DeploymentResourceProvisioningStateSucceeded,
		DeploymentResourceProvisioningStateUpdating,
	}
}

// ToPtr returns a *DeploymentResourceProvisioningState pointing to the current value.
func (c DeploymentResourceProvisioningState) ToPtr() *DeploymentResourceProvisioningState {
	return &c
}

// DeploymentResourceStatus - Status of the Deployment
type DeploymentResourceStatus string

const (
	DeploymentResourceStatusRunning DeploymentResourceStatus = "Running"
	DeploymentResourceStatusStopped DeploymentResourceStatus = "Stopped"
)

// PossibleDeploymentResourceStatusValues returns the possible values for the DeploymentResourceStatus const type.
func PossibleDeploymentResourceStatusValues() []DeploymentResourceStatus {
	return []DeploymentResourceStatus{
		DeploymentResourceStatusRunning,
		DeploymentResourceStatusStopped,
	}
}

// ToPtr returns a *DeploymentResourceStatus pointing to the current value.
func (c DeploymentResourceStatus) ToPtr() *DeploymentResourceStatus {
	return &c
}

// GatewayProvisioningState - State of the Spring Cloud Gateway.
type GatewayProvisioningState string

const (
	GatewayProvisioningStateCreating  GatewayProvisioningState = "Creating"
	GatewayProvisioningStateDeleting  GatewayProvisioningState = "Deleting"
	GatewayProvisioningStateFailed    GatewayProvisioningState = "Failed"
	GatewayProvisioningStateSucceeded GatewayProvisioningState = "Succeeded"
	GatewayProvisioningStateUpdating  GatewayProvisioningState = "Updating"
)

// PossibleGatewayProvisioningStateValues returns the possible values for the GatewayProvisioningState const type.
func PossibleGatewayProvisioningStateValues() []GatewayProvisioningState {
	return []GatewayProvisioningState{
		GatewayProvisioningStateCreating,
		GatewayProvisioningStateDeleting,
		GatewayProvisioningStateFailed,
		GatewayProvisioningStateSucceeded,
		GatewayProvisioningStateUpdating,
	}
}

// ToPtr returns a *GatewayProvisioningState pointing to the current value.
func (c GatewayProvisioningState) ToPtr() *GatewayProvisioningState {
	return &c
}

// KPackBuildStageProvisioningState - The provisioning state of this build stage resource.
type KPackBuildStageProvisioningState string

const (
	KPackBuildStageProvisioningStateFailed     KPackBuildStageProvisioningState = "Failed"
	KPackBuildStageProvisioningStateNotStarted KPackBuildStageProvisioningState = "NotStarted"
	KPackBuildStageProvisioningStateRunning    KPackBuildStageProvisioningState = "Running"
	KPackBuildStageProvisioningStateSucceeded  KPackBuildStageProvisioningState = "Succeeded"
)

// PossibleKPackBuildStageProvisioningStateValues returns the possible values for the KPackBuildStageProvisioningState const type.
func PossibleKPackBuildStageProvisioningStateValues() []KPackBuildStageProvisioningState {
	return []KPackBuildStageProvisioningState{
		KPackBuildStageProvisioningStateFailed,
		KPackBuildStageProvisioningStateNotStarted,
		KPackBuildStageProvisioningStateRunning,
		KPackBuildStageProvisioningStateSucceeded,
	}
}

// ToPtr returns a *KPackBuildStageProvisioningState pointing to the current value.
func (c KPackBuildStageProvisioningState) ToPtr() *KPackBuildStageProvisioningState {
	return &c
}

// LastModifiedByType - The type of identity that last modified the resource.
type LastModifiedByType string

const (
	LastModifiedByTypeApplication     LastModifiedByType = "Application"
	LastModifiedByTypeKey             LastModifiedByType = "Key"
	LastModifiedByTypeManagedIdentity LastModifiedByType = "ManagedIdentity"
	LastModifiedByTypeUser            LastModifiedByType = "User"
)

// PossibleLastModifiedByTypeValues returns the possible values for the LastModifiedByType const type.
func PossibleLastModifiedByTypeValues() []LastModifiedByType {
	return []LastModifiedByType{
		LastModifiedByTypeApplication,
		LastModifiedByTypeKey,
		LastModifiedByTypeManagedIdentity,
		LastModifiedByTypeUser,
	}
}

// ToPtr returns a *LastModifiedByType pointing to the current value.
func (c LastModifiedByType) ToPtr() *LastModifiedByType {
	return &c
}

// ManagedIdentityType - Type of the managed identity
type ManagedIdentityType string

const (
	ManagedIdentityTypeNone                       ManagedIdentityType = "None"
	ManagedIdentityTypeSystemAssigned             ManagedIdentityType = "SystemAssigned"
	ManagedIdentityTypeSystemAssignedUserAssigned ManagedIdentityType = "SystemAssigned,UserAssigned"
	ManagedIdentityTypeUserAssigned               ManagedIdentityType = "UserAssigned"
)

// PossibleManagedIdentityTypeValues returns the possible values for the ManagedIdentityType const type.
func PossibleManagedIdentityTypeValues() []ManagedIdentityType {
	return []ManagedIdentityType{
		ManagedIdentityTypeNone,
		ManagedIdentityTypeSystemAssigned,
		ManagedIdentityTypeSystemAssignedUserAssigned,
		ManagedIdentityTypeUserAssigned,
	}
}

// ToPtr returns a *ManagedIdentityType pointing to the current value.
func (c ManagedIdentityType) ToPtr() *ManagedIdentityType {
	return &c
}

// MonitoringSettingState - State of the Monitoring Setting.
type MonitoringSettingState string

const (
	MonitoringSettingStateFailed       MonitoringSettingState = "Failed"
	MonitoringSettingStateNotAvailable MonitoringSettingState = "NotAvailable"
	MonitoringSettingStateSucceeded    MonitoringSettingState = "Succeeded"
	MonitoringSettingStateUpdating     MonitoringSettingState = "Updating"
)

// PossibleMonitoringSettingStateValues returns the possible values for the MonitoringSettingState const type.
func PossibleMonitoringSettingStateValues() []MonitoringSettingState {
	return []MonitoringSettingState{
		MonitoringSettingStateFailed,
		MonitoringSettingStateNotAvailable,
		MonitoringSettingStateSucceeded,
		MonitoringSettingStateUpdating,
	}
}

// ToPtr returns a *MonitoringSettingState pointing to the current value.
func (c MonitoringSettingState) ToPtr() *MonitoringSettingState {
	return &c
}

// PowerState - Power state of the Service
type PowerState string

const (
	PowerStateRunning PowerState = "Running"
	PowerStateStopped PowerState = "Stopped"
)

// PossiblePowerStateValues returns the possible values for the PowerState const type.
func PossiblePowerStateValues() []PowerState {
	return []PowerState{
		PowerStateRunning,
		PowerStateStopped,
	}
}

// ToPtr returns a *PowerState pointing to the current value.
func (c PowerState) ToPtr() *PowerState {
	return &c
}

// ProvisioningState - Provisioning state of the Service
type ProvisioningState string

const (
	ProvisioningStateCreating   ProvisioningState = "Creating"
	ProvisioningStateDeleted    ProvisioningState = "Deleted"
	ProvisioningStateDeleting   ProvisioningState = "Deleting"
	ProvisioningStateFailed     ProvisioningState = "Failed"
	ProvisioningStateMoveFailed ProvisioningState = "MoveFailed"
	ProvisioningStateMoved      ProvisioningState = "Moved"
	ProvisioningStateMoving     ProvisioningState = "Moving"
	ProvisioningStateSucceeded  ProvisioningState = "Succeeded"
	ProvisioningStateUpdating   ProvisioningState = "Updating"
)

// PossibleProvisioningStateValues returns the possible values for the ProvisioningState const type.
func PossibleProvisioningStateValues() []ProvisioningState {
	return []ProvisioningState{
		ProvisioningStateCreating,
		ProvisioningStateDeleted,
		ProvisioningStateDeleting,
		ProvisioningStateFailed,
		ProvisioningStateMoveFailed,
		ProvisioningStateMoved,
		ProvisioningStateMoving,
		ProvisioningStateSucceeded,
		ProvisioningStateUpdating,
	}
}

// ToPtr returns a *ProvisioningState pointing to the current value.
func (c ProvisioningState) ToPtr() *ProvisioningState {
	return &c
}

// ResourceSKURestrictionsReasonCode - Gets the reason for restriction. Possible values include: 'QuotaId', 'NotAvailableForSubscription'
type ResourceSKURestrictionsReasonCode string

const (
	ResourceSKURestrictionsReasonCodeNotAvailableForSubscription ResourceSKURestrictionsReasonCode = "NotAvailableForSubscription"
	ResourceSKURestrictionsReasonCodeQuotaID                     ResourceSKURestrictionsReasonCode = "QuotaId"
)

// PossibleResourceSKURestrictionsReasonCodeValues returns the possible values for the ResourceSKURestrictionsReasonCode const type.
func PossibleResourceSKURestrictionsReasonCodeValues() []ResourceSKURestrictionsReasonCode {
	return []ResourceSKURestrictionsReasonCode{
		ResourceSKURestrictionsReasonCodeNotAvailableForSubscription,
		ResourceSKURestrictionsReasonCodeQuotaID,
	}
}

// ToPtr returns a *ResourceSKURestrictionsReasonCode pointing to the current value.
func (c ResourceSKURestrictionsReasonCode) ToPtr() *ResourceSKURestrictionsReasonCode {
	return &c
}

// ResourceSKURestrictionsType - Gets the type of restrictions. Possible values include: 'Location', 'Zone'
type ResourceSKURestrictionsType string

const (
	ResourceSKURestrictionsTypeLocation ResourceSKURestrictionsType = "Location"
	ResourceSKURestrictionsTypeZone     ResourceSKURestrictionsType = "Zone"
)

// PossibleResourceSKURestrictionsTypeValues returns the possible values for the ResourceSKURestrictionsType const type.
func PossibleResourceSKURestrictionsTypeValues() []ResourceSKURestrictionsType {
	return []ResourceSKURestrictionsType{
		ResourceSKURestrictionsTypeLocation,
		ResourceSKURestrictionsTypeZone,
	}
}

// ToPtr returns a *ResourceSKURestrictionsType pointing to the current value.
func (c ResourceSKURestrictionsType) ToPtr() *ResourceSKURestrictionsType {
	return &c
}

// SKUScaleType - Gets or sets the type of the scale.
type SKUScaleType string

const (
	SKUScaleTypeAutomatic SKUScaleType = "Automatic"
	SKUScaleTypeManual    SKUScaleType = "Manual"
	SKUScaleTypeNone      SKUScaleType = "None"
)

// PossibleSKUScaleTypeValues returns the possible values for the SKUScaleType const type.
func PossibleSKUScaleTypeValues() []SKUScaleType {
	return []SKUScaleType{
		SKUScaleTypeAutomatic,
		SKUScaleTypeManual,
		SKUScaleTypeNone,
	}
}

// ToPtr returns a *SKUScaleType pointing to the current value.
func (c SKUScaleType) ToPtr() *SKUScaleType {
	return &c
}

// ServiceRegistryProvisioningState - State of the Service Registry.
type ServiceRegistryProvisioningState string

const (
	ServiceRegistryProvisioningStateCreating  ServiceRegistryProvisioningState = "Creating"
	ServiceRegistryProvisioningStateDeleting  ServiceRegistryProvisioningState = "Deleting"
	ServiceRegistryProvisioningStateFailed    ServiceRegistryProvisioningState = "Failed"
	ServiceRegistryProvisioningStateSucceeded ServiceRegistryProvisioningState = "Succeeded"
	ServiceRegistryProvisioningStateUpdating  ServiceRegistryProvisioningState = "Updating"
)

// PossibleServiceRegistryProvisioningStateValues returns the possible values for the ServiceRegistryProvisioningState const type.
func PossibleServiceRegistryProvisioningStateValues() []ServiceRegistryProvisioningState {
	return []ServiceRegistryProvisioningState{
		ServiceRegistryProvisioningStateCreating,
		ServiceRegistryProvisioningStateDeleting,
		ServiceRegistryProvisioningStateFailed,
		ServiceRegistryProvisioningStateSucceeded,
		ServiceRegistryProvisioningStateUpdating,
	}
}

// ToPtr returns a *ServiceRegistryProvisioningState pointing to the current value.
func (c ServiceRegistryProvisioningState) ToPtr() *ServiceRegistryProvisioningState {
	return &c
}

// StoragePropertiesStorageType - The type of the storage.
type StoragePropertiesStorageType string

const (
	StoragePropertiesStorageTypeStorageAccount StoragePropertiesStorageType = "StorageAccount"
)

// PossibleStoragePropertiesStorageTypeValues returns the possible values for the StoragePropertiesStorageType const type.
func PossibleStoragePropertiesStorageTypeValues() []StoragePropertiesStorageType {
	return []StoragePropertiesStorageType{
		StoragePropertiesStorageTypeStorageAccount,
	}
}

// ToPtr returns a *StoragePropertiesStorageType pointing to the current value.
func (c StoragePropertiesStorageType) ToPtr() *StoragePropertiesStorageType {
	return &c
}

// SupportedRuntimePlatform - The platform of this runtime version (possible values: "Java" or ".NET").
type SupportedRuntimePlatform string

const (
	SupportedRuntimePlatformJava    SupportedRuntimePlatform = "Java"
	SupportedRuntimePlatformNETCore SupportedRuntimePlatform = ".NET Core"
)

// PossibleSupportedRuntimePlatformValues returns the possible values for the SupportedRuntimePlatform const type.
func PossibleSupportedRuntimePlatformValues() []SupportedRuntimePlatform {
	return []SupportedRuntimePlatform{
		SupportedRuntimePlatformJava,
		SupportedRuntimePlatformNETCore,
	}
}

// ToPtr returns a *SupportedRuntimePlatform pointing to the current value.
func (c SupportedRuntimePlatform) ToPtr() *SupportedRuntimePlatform {
	return &c
}

// SupportedRuntimeValue - The raw value which could be passed to deployment CRUD operations.
type SupportedRuntimeValue string

const (
	SupportedRuntimeValueJava11    SupportedRuntimeValue = "Java_11"
	SupportedRuntimeValueJava17    SupportedRuntimeValue = "Java_17"
	SupportedRuntimeValueJava8     SupportedRuntimeValue = "Java_8"
	SupportedRuntimeValueNetCore31 SupportedRuntimeValue = "NetCore_31"
)

// PossibleSupportedRuntimeValueValues returns the possible values for the SupportedRuntimeValue const type.
func PossibleSupportedRuntimeValueValues() []SupportedRuntimeValue {
	return []SupportedRuntimeValue{
		SupportedRuntimeValueJava11,
		SupportedRuntimeValueJava17,
		SupportedRuntimeValueJava8,
		SupportedRuntimeValueNetCore31,
	}
}

// ToPtr returns a *SupportedRuntimeValue pointing to the current value.
func (c SupportedRuntimeValue) ToPtr() *SupportedRuntimeValue {
	return &c
}

// TestKeyType - Type of the test key
type TestKeyType string

const (
	TestKeyTypePrimary   TestKeyType = "Primary"
	TestKeyTypeSecondary TestKeyType = "Secondary"
)

// PossibleTestKeyTypeValues returns the possible values for the TestKeyType const type.
func PossibleTestKeyTypeValues() []TestKeyType {
	return []TestKeyType{
		TestKeyTypePrimary,
		TestKeyTypeSecondary,
	}
}

// ToPtr returns a *TestKeyType pointing to the current value.
func (c TestKeyType) ToPtr() *TestKeyType {
	return &c
}

// TrafficDirection - The direction of required traffic
type TrafficDirection string

const (
	TrafficDirectionInbound  TrafficDirection = "Inbound"
	TrafficDirectionOutbound TrafficDirection = "Outbound"
)

// PossibleTrafficDirectionValues returns the possible values for the TrafficDirection const type.
func PossibleTrafficDirectionValues() []TrafficDirection {
	return []TrafficDirection{
		TrafficDirectionInbound,
		TrafficDirectionOutbound,
	}
}

// ToPtr returns a *TrafficDirection pointing to the current value.
func (c TrafficDirection) ToPtr() *TrafficDirection {
	return &c
}
