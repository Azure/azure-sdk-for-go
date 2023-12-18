# Release History

## 4.0.0-beta.3 (2023-11-24)
### Features Added

- New value `MigrationStateCleaningUp`, `MigrationStateValidationFailed` added to enum type `MigrationState`
- New value `MigrationSubStateCancelingRequestedDBMigrations`, `MigrationSubStateValidationInProgress` added to enum type `MigrationSubState`
- New value `ServerVersionSixteen` added to enum type `ServerVersion`
- New enum type `MigrationDbState` with values `MigrationDbStateCanceled`, `MigrationDbStateCanceling`, `MigrationDbStateFailed`, `MigrationDbStateInProgress`, `MigrationDbStateSucceeded`, `MigrationDbStateWaitingForCutoverTrigger`
- New enum type `MigrationOption` with values `MigrationOptionMigrate`, `MigrationOptionValidate`, `MigrationOptionValidateAndMigrate`
- New enum type `PrivateEndpointConnectionProvisioningState` with values `PrivateEndpointConnectionProvisioningStateCreating`, `PrivateEndpointConnectionProvisioningStateDeleting`, `PrivateEndpointConnectionProvisioningStateFailed`, `PrivateEndpointConnectionProvisioningStateSucceeded`
- New enum type `PrivateEndpointServiceConnectionStatus` with values `PrivateEndpointServiceConnectionStatusApproved`, `PrivateEndpointServiceConnectionStatusPending`, `PrivateEndpointServiceConnectionStatusRejected`
- New enum type `ReadReplicaPromoteMode` with values `ReadReplicaPromoteModeStandalone`, `ReadReplicaPromoteModeSwitchover`
- New enum type `ReplicationPromoteOption` with values `ReplicationPromoteOptionForced`, `ReplicationPromoteOptionPlanned`
- New enum type `ReplicationState` with values `ReplicationStateActive`, `ReplicationStateBroken`, `ReplicationStateCatchup`, `ReplicationStateProvisioning`, `ReplicationStateReconfiguring`, `ReplicationStateUpdating`
- New enum type `SSLMode` with values `SSLModePrefer`, `SSLModeRequire`, `SSLModeVerifyCA`, `SSLModeVerifyFull`
- New enum type `SourceType` with values `SourceTypeAWS`, `SourceTypeAzureVM`, `SourceTypeGCP`, `SourceTypeOnPremises`, `SourceTypePostgreSQLSingleServer`
- New enum type `StorageType` with values `StorageTypePremiumLRS`, `StorageTypePremiumV2LRS`
- New enum type `ThreatProtectionName` with values `ThreatProtectionNameDefault`
- New enum type `ThreatProtectionState` with values `ThreatProtectionStateDisabled`, `ThreatProtectionStateEnabled`
- New enum type `ValidationState` with values `ValidationStateFailed`, `ValidationStateSucceeded`, `ValidationStateWarning`
- New enum type `VirtualEndpointType` with values `VirtualEndpointTypeReadWrite`
- New function `*ClientFactory.NewPrivateEndpointConnectionClient() *PrivateEndpointConnectionClient`
- New function `*ClientFactory.NewPrivateEndpointConnectionsClient() *PrivateEndpointConnectionsClient`
- New function `*ClientFactory.NewPrivateLinkResourcesClient() *PrivateLinkResourcesClient`
- New function `*ClientFactory.NewQuotaUsagesClient() *QuotaUsagesClient`
- New function `*ClientFactory.NewServerThreatProtectionSettingsClient() *ServerThreatProtectionSettingsClient`
- New function `*ClientFactory.NewVirtualEndpointsClient() *VirtualEndpointsClient`
- New function `NewPrivateEndpointConnectionClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateEndpointConnectionClient, error)`
- New function `*PrivateEndpointConnectionClient.BeginDelete(context.Context, string, string, string, *PrivateEndpointConnectionClientBeginDeleteOptions) (*runtime.Poller[PrivateEndpointConnectionClientDeleteResponse], error)`
- New function `*PrivateEndpointConnectionClient.BeginUpdate(context.Context, string, string, string, PrivateEndpointConnection, *PrivateEndpointConnectionClientBeginUpdateOptions) (*runtime.Poller[PrivateEndpointConnectionClientUpdateResponse], error)`
- New function `NewPrivateEndpointConnectionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateEndpointConnectionsClient, error)`
- New function `*PrivateEndpointConnectionsClient.Get(context.Context, string, string, string, *PrivateEndpointConnectionsClientGetOptions) (PrivateEndpointConnectionsClientGetResponse, error)`
- New function `*PrivateEndpointConnectionsClient.NewListByServerPager(string, string, *PrivateEndpointConnectionsClientListByServerOptions) *runtime.Pager[PrivateEndpointConnectionsClientListByServerResponse]`
- New function `NewPrivateLinkResourcesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateLinkResourcesClient, error)`
- New function `*PrivateLinkResourcesClient.Get(context.Context, string, string, string, *PrivateLinkResourcesClientGetOptions) (PrivateLinkResourcesClientGetResponse, error)`
- New function `*PrivateLinkResourcesClient.NewListByServerPager(string, string, *PrivateLinkResourcesClientListByServerOptions) *runtime.Pager[PrivateLinkResourcesClientListByServerResponse]`
- New function `NewQuotaUsagesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*QuotaUsagesClient, error)`
- New function `*QuotaUsagesClient.NewListPager(string, *QuotaUsagesClientListOptions) *runtime.Pager[QuotaUsagesClientListResponse]`
- New function `NewServerThreatProtectionSettingsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ServerThreatProtectionSettingsClient, error)`
- New function `*ServerThreatProtectionSettingsClient.BeginCreateOrUpdate(context.Context, string, string, ThreatProtectionName, ServerThreatProtectionSettingsModel, *ServerThreatProtectionSettingsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ServerThreatProtectionSettingsClientCreateOrUpdateResponse], error)`
- New function `*ServerThreatProtectionSettingsClient.Get(context.Context, string, string, ThreatProtectionName, *ServerThreatProtectionSettingsClientGetOptions) (ServerThreatProtectionSettingsClientGetResponse, error)`
- New function `*ServerThreatProtectionSettingsClient.NewListByServerPager(string, string, *ServerThreatProtectionSettingsClientListByServerOptions) *runtime.Pager[ServerThreatProtectionSettingsClientListByServerResponse]`
- New function `NewVirtualEndpointsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VirtualEndpointsClient, error)`
- New function `*VirtualEndpointsClient.BeginCreate(context.Context, string, string, string, VirtualEndpointResource, *VirtualEndpointsClientBeginCreateOptions) (*runtime.Poller[VirtualEndpointsClientCreateResponse], error)`
- New function `*VirtualEndpointsClient.BeginDelete(context.Context, string, string, string, *VirtualEndpointsClientBeginDeleteOptions) (*runtime.Poller[VirtualEndpointsClientDeleteResponse], error)`
- New function `*VirtualEndpointsClient.Get(context.Context, string, string, string, *VirtualEndpointsClientGetOptions) (VirtualEndpointsClientGetResponse, error)`
- New function `*VirtualEndpointsClient.NewListByServerPager(string, string, *VirtualEndpointsClientListByServerOptions) *runtime.Pager[VirtualEndpointsClientListByServerResponse]`
- New function `*VirtualEndpointsClient.BeginUpdate(context.Context, string, string, string, VirtualEndpointResourceForPatch, *VirtualEndpointsClientBeginUpdateOptions) (*runtime.Poller[VirtualEndpointsClientUpdateResponse], error)`
- New struct `DbLevelValidationStatus`
- New struct `DbMigrationStatus`
- New struct `NameProperty`
- New struct `PrivateEndpoint`
- New struct `PrivateEndpointConnection`
- New struct `PrivateEndpointConnectionListResult`
- New struct `PrivateEndpointConnectionProperties`
- New struct `PrivateLinkResource`
- New struct `PrivateLinkResourceListResult`
- New struct `PrivateLinkResourceProperties`
- New struct `PrivateLinkServiceConnectionState`
- New struct `QuotaUsage`
- New struct `QuotaUsagesListResult`
- New struct `Replica`
- New struct `ServerThreatProtectionListResult`
- New struct `ServerThreatProtectionProperties`
- New struct `ServerThreatProtectionSettingsModel`
- New struct `ValidationDetails`
- New struct `ValidationMessage`
- New struct `ValidationSummaryItem`
- New struct `VirtualEndpointResource`
- New struct `VirtualEndpointResourceForPatch`
- New struct `VirtualEndpointResourceProperties`
- New struct `VirtualEndpointsListResult`
- New field `MigrationOption`, `SSLMode`, `SourceType` in struct `MigrationResourceProperties`
- New field `DbDetails`, `ValidationDetails` in struct `MigrationSubStateDetails`
- New field `PrivateEndpointConnections`, `Replica` in struct `ServerProperties`
- New field `Replica` in struct `ServerPropertiesForUpdate`
- New field `Throughput`, `Type` in struct `Storage`
- New field `MaximumStorageSizeMb`, `SupportedMaximumIops`, `SupportedMaximumThroughput`, `SupportedThroughput` in struct `StorageMbCapability`


## 3.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 4.0.0-beta.2 (2023-10-27)
### Breaking Changes

- Field `IopsTier` of struct `Storage` has been removed

### Features Added

- New field `Tier` in struct `Storage`


## 4.0.0-beta.1 (2023-05-26)
### Breaking Changes

- Type of `CapabilitiesListResult.Value` has been changed from `[]*CapabilityProperties` to `[]*FlexibleServerCapability`
- Type of `FastProvisioningEditionCapability.SupportedStorageGb` has been changed from `*int64` to `*int32`
- Type of `FlexibleServerEditionCapability.Status` has been changed from `*string` to `*CapabilityStatus`
- Type of `ServerVersionCapability.Status` has been changed from `*string` to `*CapabilityStatus`
- Type of `StorageEditionCapability.Status` has been changed from `*string` to `*CapabilityStatus`
- Type of `StorageTierCapability.Iops` has been changed from `*int64` to `*int32`
- Type of `StorageTierCapability.Status` has been changed from `*string` to `*CapabilityStatus`
- `IdentityTypeSystemAssigned` from enum `IdentityType` has been removed
- Struct `CapabilityProperties` has been removed
- Struct `HyperscaleNodeEditionCapability` has been removed
- Struct `NodeTypeCapability` has been removed
- Struct `StorageMBCapability` has been removed
- Struct `VcoreCapability` has been removed
- Field `SupportedServerVersions` of struct `FlexibleServerEditionCapability` has been removed
- Field `SupportedVcores` of struct `ServerVersionCapability` has been removed
- Field `SupportedStorageMB` of struct `StorageEditionCapability` has been removed
- Field `IsBaseline`, `TierName` of struct `StorageTierCapability` has been removed

### Features Added

- New value `CreateModeReviveDropped` added to enum type `CreateMode`
- New value `ServerVersionFifteen` added to enum type `ServerVersion`
- New enum type `AzureManagedDiskPerformanceTiers` with values `AzureManagedDiskPerformanceTiersP1`, `AzureManagedDiskPerformanceTiersP10`, `AzureManagedDiskPerformanceTiersP15`, `AzureManagedDiskPerformanceTiersP2`, `AzureManagedDiskPerformanceTiersP20`, `AzureManagedDiskPerformanceTiersP3`, `AzureManagedDiskPerformanceTiersP30`, `AzureManagedDiskPerformanceTiersP4`, `AzureManagedDiskPerformanceTiersP40`, `AzureManagedDiskPerformanceTiersP50`, `AzureManagedDiskPerformanceTiersP6`, `AzureManagedDiskPerformanceTiersP60`, `AzureManagedDiskPerformanceTiersP70`, `AzureManagedDiskPerformanceTiersP80`
- New enum type `CancelEnum` with values `CancelEnumFalse`, `CancelEnumTrue`
- New enum type `CapabilityStatus` with values `CapabilityStatusAvailable`, `CapabilityStatusDefault`, `CapabilityStatusDisabled`, `CapabilityStatusVisible`
- New enum type `ExecutionStatus` with values `ExecutionStatusCancelled`, `ExecutionStatusFailed`, `ExecutionStatusRunning`, `ExecutionStatusSucceeded`
- New enum type `FastProvisioningSupportedEnum` with values `FastProvisioningSupportedEnumDisabled`, `FastProvisioningSupportedEnumEnabled`
- New enum type `GeoBackupSupportedEnum` with values `GeoBackupSupportedEnumDisabled`, `GeoBackupSupportedEnumEnabled`
- New enum type `HaMode` with values `HaModeSameZone`, `HaModeZoneRedundant`
- New enum type `KeyStatusEnum` with values `KeyStatusEnumInvalid`, `KeyStatusEnumValid`
- New enum type `LogicalReplicationOnSourceDbEnum` with values `LogicalReplicationOnSourceDbEnumFalse`, `LogicalReplicationOnSourceDbEnumTrue`
- New enum type `MigrationListFilter` with values `MigrationListFilterActive`, `MigrationListFilterAll`
- New enum type `MigrationMode` with values `MigrationModeOffline`, `MigrationModeOnline`
- New enum type `MigrationNameAvailabilityReason` with values `MigrationNameAvailabilityReasonAlreadyExists`, `MigrationNameAvailabilityReasonInvalid`
- New enum type `MigrationState` with values `MigrationStateCanceled`, `MigrationStateFailed`, `MigrationStateInProgress`, `MigrationStateSucceeded`, `MigrationStateWaitingForUserAction`
- New enum type `MigrationSubState` with values `MigrationSubStateCompleted`, `MigrationSubStateCompletingMigration`, `MigrationSubStateMigratingData`, `MigrationSubStatePerformingPreRequisiteSteps`, `MigrationSubStateWaitingForCutoverTrigger`, `MigrationSubStateWaitingForDBsToMigrateSpecification`, `MigrationSubStateWaitingForDataMigrationScheduling`, `MigrationSubStateWaitingForDataMigrationWindow`, `MigrationSubStateWaitingForLogicalReplicationSetupRequestOnSourceDB`, `MigrationSubStateWaitingForTargetDBOverwriteConfirmation`
- New enum type `OnlineResizeSupportedEnum` with values `OnlineResizeSupportedEnumDisabled`, `OnlineResizeSupportedEnumEnabled`
- New enum type `OverwriteDbsInTargetEnum` with values `OverwriteDbsInTargetEnumFalse`, `OverwriteDbsInTargetEnumTrue`
- New enum type `RestrictedEnum` with values `RestrictedEnumDisabled`, `RestrictedEnumEnabled`
- New enum type `StartDataMigrationEnum` with values `StartDataMigrationEnumFalse`, `StartDataMigrationEnumTrue`
- New enum type `StorageAutoGrow` with values `StorageAutoGrowDisabled`, `StorageAutoGrowEnabled`
- New enum type `StorageAutoGrowthSupportedEnum` with values `StorageAutoGrowthSupportedEnumDisabled`, `StorageAutoGrowthSupportedEnumEnabled`
- New enum type `TriggerCutoverEnum` with values `TriggerCutoverEnumFalse`, `TriggerCutoverEnumTrue`
- New enum type `ZoneRedundantHaAndGeoBackupSupportedEnum` with values `ZoneRedundantHaAndGeoBackupSupportedEnumDisabled`, `ZoneRedundantHaAndGeoBackupSupportedEnumEnabled`
- New enum type `ZoneRedundantHaSupportedEnum` with values `ZoneRedundantHaSupportedEnumDisabled`, `ZoneRedundantHaSupportedEnumEnabled`
- New function `*ClientFactory.NewFlexibleServerClient() *FlexibleServerClient`
- New function `*ClientFactory.NewLogFilesClient() *LogFilesClient`
- New function `*ClientFactory.NewLtrBackupOperationsClient() *LtrBackupOperationsClient`
- New function `*ClientFactory.NewMigrationsClient() *MigrationsClient`
- New function `*ClientFactory.NewPostgreSQLManagementClient() *PostgreSQLManagementClient`
- New function `*ClientFactory.NewServerCapabilitiesClient() *ServerCapabilitiesClient`
- New function `NewFlexibleServerClient(string, azcore.TokenCredential, *arm.ClientOptions) (*FlexibleServerClient, error)`
- New function `*FlexibleServerClient.BeginStartLtrBackup(context.Context, string, string, LtrBackupRequest, *FlexibleServerClientBeginStartLtrBackupOptions) (*runtime.Poller[FlexibleServerClientStartLtrBackupResponse], error)`
- New function `*FlexibleServerClient.TriggerLtrPreBackup(context.Context, string, string, LtrPreBackupRequest, *FlexibleServerClientTriggerLtrPreBackupOptions) (FlexibleServerClientTriggerLtrPreBackupResponse, error)`
- New function `NewLogFilesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LogFilesClient, error)`
- New function `*LogFilesClient.NewListByServerPager(string, string, *LogFilesClientListByServerOptions) *runtime.Pager[LogFilesClientListByServerResponse]`
- New function `NewLtrBackupOperationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LtrBackupOperationsClient, error)`
- New function `*LtrBackupOperationsClient.Get(context.Context, string, string, string, *LtrBackupOperationsClientGetOptions) (LtrBackupOperationsClientGetResponse, error)`
- New function `*LtrBackupOperationsClient.NewListByServerPager(string, string, *LtrBackupOperationsClientListByServerOptions) *runtime.Pager[LtrBackupOperationsClientListByServerResponse]`
- New function `NewMigrationsClient(azcore.TokenCredential, *arm.ClientOptions) (*MigrationsClient, error)`
- New function `*MigrationsClient.Create(context.Context, string, string, string, string, MigrationResource, *MigrationsClientCreateOptions) (MigrationsClientCreateResponse, error)`
- New function `*MigrationsClient.Delete(context.Context, string, string, string, string, *MigrationsClientDeleteOptions) (MigrationsClientDeleteResponse, error)`
- New function `*MigrationsClient.Get(context.Context, string, string, string, string, *MigrationsClientGetOptions) (MigrationsClientGetResponse, error)`
- New function `*MigrationsClient.NewListByTargetServerPager(string, string, string, *MigrationsClientListByTargetServerOptions) *runtime.Pager[MigrationsClientListByTargetServerResponse]`
- New function `*MigrationsClient.Update(context.Context, string, string, string, string, MigrationResourceForPatch, *MigrationsClientUpdateOptions) (MigrationsClientUpdateResponse, error)`
- New function `PossibleStorageAutoGrowValues() []StorageAutoGrow`
- New function `NewPostgreSQLManagementClient(azcore.TokenCredential, *arm.ClientOptions) (*PostgreSQLManagementClient, error)`
- New function `*PostgreSQLManagementClient.CheckMigrationNameAvailability(context.Context, string, string, string, MigrationNameAvailabilityResource, *PostgreSQLManagementClientCheckMigrationNameAvailabilityOptions) (PostgreSQLManagementClientCheckMigrationNameAvailabilityResponse, error)`
- New function `NewServerCapabilitiesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ServerCapabilitiesClient, error)`
- New function `*ServerCapabilitiesClient.NewListPager(string, string, *ServerCapabilitiesClientListOptions) *runtime.Pager[ServerCapabilitiesClientListResponse]`
- New struct `AdminCredentials`
- New struct `BackupSettings`
- New struct `BackupStoreDetails`
- New struct `DbServerMetadata`
- New struct `FlexibleServerCapability`
- New struct `LogFile`
- New struct `LogFileListResult`
- New struct `LogFileProperties`
- New struct `LtrBackupOperationResponseProperties`
- New struct `LtrBackupRequest`
- New struct `LtrBackupResponse`
- New struct `LtrPreBackupRequest`
- New struct `LtrPreBackupResponse`
- New struct `LtrPreBackupResponseProperties`
- New struct `LtrServerBackupOperation`
- New struct `LtrServerBackupOperationList`
- New struct `MigrationNameAvailabilityResource`
- New struct `MigrationResource`
- New struct `MigrationResourceForPatch`
- New struct `MigrationResourceListResult`
- New struct `MigrationResourceProperties`
- New struct `MigrationResourcePropertiesForPatch`
- New struct `MigrationSecretParameters`
- New struct `MigrationStatus`
- New struct `MigrationSubStateDetails`
- New struct `ServerSKU`
- New struct `ServerSKUCapability`
- New struct `StorageMbCapability`
- New field `GeoBackupEncryptionKeyStatus`, `GeoBackupKeyURI`, `GeoBackupUserAssignedIdentityID`, `PrimaryEncryptionKeyStatus` in struct `DataEncryption`
- New field `Reason`, `ServerCount`, `Status`, `SupportedTier` in struct `FastProvisioningEditionCapability`
- New field `DefaultSKUName`, `Reason`, `SupportedServerSKUs` in struct `FlexibleServerEditionCapability`
- New field `Network` in struct `ServerPropertiesForUpdate`
- New field `Reason` in struct `ServerVersionCapability`
- New field `AutoGrow`, `Iops`, `IopsTier` in struct `Storage`
- New field `DefaultStorageSizeMb`, `Reason`, `SupportedStorageMb` in struct `StorageEditionCapability`
- New field `Reason` in struct `StorageTierCapability`
- New field `TenantID` in struct `UserAssignedIdentity`


## 3.0.0 (2023-04-28)
### Breaking Changes

- `ArmServerKeyTypeSystemAssigned` from enum `ArmServerKeyType` has been removed
- `ReplicationRoleGeoSyncReplica`, `ReplicationRoleSecondary`, `ReplicationRoleSyncReplica`, `ReplicationRoleWalReplica` from enum `ReplicationRole` has been removed

### Features Added

- New value `ArmServerKeyTypeSystemManaged` added to enum type `ArmServerKeyType`


## 2.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 2.1.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module

## 2.0.0 (2023-01-27)
### Breaking Changes

- Function `*CheckNameAvailabilityClient.Execute` parameter(s) have been changed from `(context.Context, NameAvailabilityRequest, *CheckNameAvailabilityClientExecuteOptions)` to `(context.Context, CheckNameAvailabilityRequest, *CheckNameAvailabilityClientExecuteOptions)`
- Function `*ConfigurationsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, string, Configuration, *ConfigurationsClientBeginUpdateOptions)` to `(context.Context, string, string, string, ConfigurationForUpdate, *ConfigurationsClientBeginUpdateOptions)`
- Type of `NameAvailability.Reason` has been changed from `*Reason` to `*CheckNameAvailabilityReason`
- Type alias `Reason` has been removed
- Struct `NameAvailabilityRequest` has been removed
- Field `Location` of struct `ServerForUpdate` has been removed

### Features Added

- New value `CreateModeGeoRestore`, `CreateModeReplica` added to type alias `CreateMode`
- New value `HighAvailabilityModeSameZone` added to type alias `HighAvailabilityMode`
- New type alias `ActiveDirectoryAuthEnum` with values `ActiveDirectoryAuthEnumDisabled`, `ActiveDirectoryAuthEnumEnabled`
- New type alias `ArmServerKeyType` with values `ArmServerKeyTypeAzureKeyVault`, `ArmServerKeyTypeSystemAssigned`
- New type alias `CheckNameAvailabilityReason` with values `CheckNameAvailabilityReasonAlreadyExists`, `CheckNameAvailabilityReasonInvalid`
- New type alias `IdentityType` with values `IdentityTypeNone`, `IdentityTypeSystemAssigned`, `IdentityTypeUserAssigned`
- New type alias `Origin` with values `OriginFull`
- New type alias `PasswordAuthEnum` with values `PasswordAuthEnumDisabled`, `PasswordAuthEnumEnabled`
- New type alias `PrincipalType` with values `PrincipalTypeGroup`, `PrincipalTypeServicePrincipal`, `PrincipalTypeUnknown`, `PrincipalTypeUser`
- New type alias `ReplicationRole` with values `ReplicationRoleAsyncReplica`, `ReplicationRoleGeoAsyncReplica`, `ReplicationRoleGeoSyncReplica`, `ReplicationRoleNone`, `ReplicationRolePrimary`, `ReplicationRoleSecondary`, `ReplicationRoleSyncReplica`, `ReplicationRoleWalReplica`
- New function `NewAdministratorsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AdministratorsClient, error)`
- New function `*AdministratorsClient.BeginCreate(context.Context, string, string, string, ActiveDirectoryAdministratorAdd, *AdministratorsClientBeginCreateOptions) (*runtime.Poller[AdministratorsClientCreateResponse], error)`
- New function `*AdministratorsClient.BeginDelete(context.Context, string, string, string, *AdministratorsClientBeginDeleteOptions) (*runtime.Poller[AdministratorsClientDeleteResponse], error)`
- New function `*AdministratorsClient.Get(context.Context, string, string, string, *AdministratorsClientGetOptions) (AdministratorsClientGetResponse, error)`
- New function `*AdministratorsClient.NewListByServerPager(string, string, *AdministratorsClientListByServerOptions) *runtime.Pager[AdministratorsClientListByServerResponse]`
- New function `NewBackupsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BackupsClient, error)`
- New function `*BackupsClient.Get(context.Context, string, string, string, *BackupsClientGetOptions) (BackupsClientGetResponse, error)`
- New function `*BackupsClient.NewListByServerPager(string, string, *BackupsClientListByServerOptions) *runtime.Pager[BackupsClientListByServerResponse]`
- New function `NewCheckNameAvailabilityWithLocationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CheckNameAvailabilityWithLocationClient, error)`
- New function `*CheckNameAvailabilityWithLocationClient.Execute(context.Context, string, CheckNameAvailabilityRequest, *CheckNameAvailabilityWithLocationClientExecuteOptions) (CheckNameAvailabilityWithLocationClientExecuteResponse, error)`
- New function `NewReplicasClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ReplicasClient, error)`
- New function `*ReplicasClient.NewListByServerPager(string, string, *ReplicasClientListByServerOptions) *runtime.Pager[ReplicasClientListByServerResponse]`
- New struct `ActiveDirectoryAdministrator`
- New struct `ActiveDirectoryAdministratorAdd`
- New struct `AdministratorListResult`
- New struct `AdministratorProperties`
- New struct `AdministratorPropertiesForAdd`
- New struct `AdministratorsClient`
- New struct `AdministratorsClientCreateResponse`
- New struct `AdministratorsClientDeleteResponse`
- New struct `AdministratorsClientListByServerResponse`
- New struct `AuthConfig`
- New struct `BackupsClient`
- New struct `BackupsClientListByServerResponse`
- New struct `CheckNameAvailabilityRequest`
- New struct `CheckNameAvailabilityWithLocationClient`
- New struct `ConfigurationForUpdate`
- New struct `DataEncryption`
- New struct `FastProvisioningEditionCapability`
- New struct `ReplicasClient`
- New struct `ReplicasClientListByServerResponse`
- New struct `ServerBackup`
- New struct `ServerBackupListResult`
- New struct `ServerBackupProperties`
- New struct `StorageTierCapability`
- New struct `UserAssignedIdentity`
- New struct `UserIdentity`
- New field `FastProvisioningSupported` in struct `CapabilityProperties`
- New field `SupportedFastProvisioningEditions` in struct `CapabilityProperties`
- New field `Identity` in struct `Server`
- New field `Identity` in struct `ServerForUpdate`
- New field `AuthConfig` in struct `ServerProperties`
- New field `DataEncryption` in struct `ServerProperties`
- New field `ReplicaCapacity` in struct `ServerProperties`
- New field `ReplicationRole` in struct `ServerProperties`
- New field `AuthConfig` in struct `ServerPropertiesForUpdate`
- New field `DataEncryption` in struct `ServerPropertiesForUpdate`
- New field `ReplicationRole` in struct `ServerPropertiesForUpdate`
- New field `Version` in struct `ServerPropertiesForUpdate`
- New field `SupportedVersionsToUpgrade` in struct `ServerVersionCapability`
- New field `SupportedUpgradableTierList` in struct `StorageMBCapability`


## 1.1.0 (2022-07-21)
### Features Added

- New const `ServerVersionFourteen`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).