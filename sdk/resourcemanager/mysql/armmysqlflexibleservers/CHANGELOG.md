# Release History

## 2.0.0 (2025-10-13)
### Breaking Changes

- Type of `Operation.Origin` has been changed from `*string` to `*Origin`
- Type of `Server.Identity` has been changed from `*Identity` to `*MySQLServerIdentity`
- Type of `Server.SKU` has been changed from `*SKU` to `*MySQLServerSKU`
- Type of `ServerForUpdate.Identity` has been changed from `*Identity` to `*MySQLServerIdentity`
- Type of `ServerForUpdate.SKU` has been changed from `*SKU` to `*MySQLServerSKU`
- Enum `SKUTier` has been removed
- Struct `ErrorResponse` has been removed
- Struct `Identity` has been removed
- Struct `ProxyResource` has been removed
- Struct `Resource` has been removed
- Struct `SKU` has been removed
- Struct `TrackedResource` has been removed

### Features Added

- New enum type `AdministratorName` with values `AdministratorNameActiveDirectory`
- New enum type `AdministratorType` with values `AdministratorTypeActiveDirectory`
- New enum type `AdvancedThreatProtectionName` with values `AdvancedThreatProtectionNameDefault`
- New enum type `AdvancedThreatProtectionProvisioningState` with values `AdvancedThreatProtectionProvisioningStateCanceled`, `AdvancedThreatProtectionProvisioningStateFailed`, `AdvancedThreatProtectionProvisioningStateSucceeded`, `AdvancedThreatProtectionProvisioningStateUpdating`
- New enum type `AdvancedThreatProtectionState` with values `AdvancedThreatProtectionStateDisabled`, `AdvancedThreatProtectionStateEnabled`
- New enum type `BackupFormat` with values `BackupFormatCollatedFormat`, `BackupFormatRaw`
- New enum type `BackupType` with values `BackupTypeFULL`
- New enum type `BatchOfMaintenance` with values `BatchOfMaintenanceBatch1`, `BatchOfMaintenanceBatch2`, `BatchOfMaintenanceDefault`
- New enum type `ImportSourceStorageType` with values `ImportSourceStorageTypeAzureBlob`
- New enum type `MaintenanceProvisioningState` with values `MaintenanceProvisioningStateCreating`, `MaintenanceProvisioningStateDeleting`, `MaintenanceProvisioningStateFailed`, `MaintenanceProvisioningStateSucceeded`
- New enum type `MaintenanceState` with values `MaintenanceStateCanceled`, `MaintenanceStateCompleted`, `MaintenanceStateInPreparation`, `MaintenanceStateProcessing`, `MaintenanceStateReScheduled`, `MaintenanceStateScheduled`
- New enum type `MaintenanceType` with values `MaintenanceTypeHotFixes`, `MaintenanceTypeMinorVersionUpgrade`, `MaintenanceTypeRoutineMaintenance`, `MaintenanceTypeSecurityPatches`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeUserAssigned`
- New enum type `ObjectType` with values `ObjectTypeBackupAndExportResponse`, `ObjectTypeImportFromStorageResponse`
- New enum type `OperationStatus` with values `OperationStatusCancelInProgress`, `OperationStatusCanceled`, `OperationStatusFailed`, `OperationStatusInProgress`, `OperationStatusPending`, `OperationStatusSucceeded`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `PatchStrategy` with values `PatchStrategyRegular`, `PatchStrategyVirtualCanary`
- New enum type `PrivateEndpointConnectionProvisioningState` with values `PrivateEndpointConnectionProvisioningStateCreating`, `PrivateEndpointConnectionProvisioningStateDeleting`, `PrivateEndpointConnectionProvisioningStateFailed`, `PrivateEndpointConnectionProvisioningStateSucceeded`
- New enum type `PrivateEndpointServiceConnectionStatus` with values `PrivateEndpointServiceConnectionStatusApproved`, `PrivateEndpointServiceConnectionStatusPending`, `PrivateEndpointServiceConnectionStatusRejected`
- New enum type `ProvisioningState` with values `ProvisioningStateCanceled`, `ProvisioningStateCreating`, `ProvisioningStateDeleting`, `ProvisioningStateFailed`, `ProvisioningStateSucceeded`
- New enum type `ResetAllToDefault` with values `ResetAllToDefaultFalse`, `ResetAllToDefaultTrue`
- New enum type `ServerSKUTier` with values `ServerSKUTierBurstable`, `ServerSKUTierGeneralPurpose`, `ServerSKUTierMemoryOptimized`
- New enum type `StorageRedundancyEnum` with values `StorageRedundancyEnumLocalRedundancy`, `StorageRedundancyEnumZoneRedundancy`
- New function `NewAdvancedThreatProtectionSettingsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AdvancedThreatProtectionSettingsClient, error)`
- New function `*AdvancedThreatProtectionSettingsClient.Get(context.Context, string, string, AdvancedThreatProtectionName, *AdvancedThreatProtectionSettingsClientGetOptions) (AdvancedThreatProtectionSettingsClientGetResponse, error)`
- New function `*AdvancedThreatProtectionSettingsClient.NewListPager(string, string, *AdvancedThreatProtectionSettingsClientListOptions) *runtime.Pager[AdvancedThreatProtectionSettingsClientListResponse]`
- New function `*AdvancedThreatProtectionSettingsClient.BeginUpdate(context.Context, string, string, AdvancedThreatProtectionName, AdvancedThreatProtectionForUpdate, *AdvancedThreatProtectionSettingsClientBeginUpdateOptions) (*runtime.Poller[AdvancedThreatProtectionSettingsClientUpdateResponse], error)`
- New function `*AdvancedThreatProtectionSettingsClient.BeginUpdatePut(context.Context, string, string, AdvancedThreatProtectionName, AdvancedThreatProtection, *AdvancedThreatProtectionSettingsClientBeginUpdatePutOptions) (*runtime.Poller[AdvancedThreatProtectionSettingsClientUpdatePutResponse], error)`
- New function `NewAzureADAdministratorsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AzureADAdministratorsClient, error)`
- New function `*AzureADAdministratorsClient.BeginCreateOrUpdate(context.Context, string, string, AdministratorName, AzureADAdministrator, *AzureADAdministratorsClientBeginCreateOrUpdateOptions) (*runtime.Poller[AzureADAdministratorsClientCreateOrUpdateResponse], error)`
- New function `*AzureADAdministratorsClient.BeginDelete(context.Context, string, string, AdministratorName, *AzureADAdministratorsClientBeginDeleteOptions) (*runtime.Poller[AzureADAdministratorsClientDeleteResponse], error)`
- New function `*AzureADAdministratorsClient.Get(context.Context, string, string, AdministratorName, *AzureADAdministratorsClientGetOptions) (AzureADAdministratorsClientGetResponse, error)`
- New function `*AzureADAdministratorsClient.NewListByServerPager(string, string, *AzureADAdministratorsClientListByServerOptions) *runtime.Pager[AzureADAdministratorsClientListByServerResponse]`
- New function `NewBackupAndExportClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BackupAndExportClient, error)`
- New function `*BackupAndExportClient.BeginCreate(context.Context, string, string, BackupAndExportRequest, *BackupAndExportClientBeginCreateOptions) (*runtime.Poller[BackupAndExportClientCreateResponse], error)`
- New function `*BackupAndExportClient.ValidateBackup(context.Context, string, string, *BackupAndExportClientValidateBackupOptions) (BackupAndExportClientValidateBackupResponse, error)`
- New function `*BackupAndExportResponseType.GetOperationProgressResponseType() *OperationProgressResponseType`
- New function `*BackupStoreDetails.GetBackupStoreDetails() *BackupStoreDetails`
- New function `*BackupsClient.Put(context.Context, string, string, string, *BackupsClientPutOptions) (BackupsClientPutResponse, error)`
- New function `NewCheckNameAvailabilityWithoutLocationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CheckNameAvailabilityWithoutLocationClient, error)`
- New function `*CheckNameAvailabilityWithoutLocationClient.Execute(context.Context, NameAvailabilityRequest, *CheckNameAvailabilityWithoutLocationClientExecuteOptions) (CheckNameAvailabilityWithoutLocationClientExecuteResponse, error)`
- New function `*ClientFactory.NewAdvancedThreatProtectionSettingsClient() *AdvancedThreatProtectionSettingsClient`
- New function `*ClientFactory.NewAzureADAdministratorsClient() *AzureADAdministratorsClient`
- New function `*ClientFactory.NewBackupAndExportClient() *BackupAndExportClient`
- New function `*ClientFactory.NewCheckNameAvailabilityWithoutLocationClient() *CheckNameAvailabilityWithoutLocationClient`
- New function `*ClientFactory.NewLocationBasedCapabilitySetClient() *LocationBasedCapabilitySetClient`
- New function `*ClientFactory.NewLogFilesClient() *LogFilesClient`
- New function `*ClientFactory.NewLongRunningBackupClient() *LongRunningBackupClient`
- New function `*ClientFactory.NewLongRunningBackupsClient() *LongRunningBackupsClient`
- New function `*ClientFactory.NewMaintenancesClient() *MaintenancesClient`
- New function `*ClientFactory.NewOperationProgressClient() *OperationProgressClient`
- New function `*ClientFactory.NewOperationResultsClient() *OperationResultsClient`
- New function `*ClientFactory.NewPrivateEndpointConnectionsClient() *PrivateEndpointConnectionsClient`
- New function `*ClientFactory.NewPrivateLinkResourcesClient() *PrivateLinkResourcesClient`
- New function `*ClientFactory.NewServersMigrationClient() *ServersMigrationClient`
- New function `*ConfigurationsClient.BeginCreateOrUpdate(context.Context, string, string, string, Configuration, *ConfigurationsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ConfigurationsClientCreateOrUpdateResponse], error)`
- New function `*FullBackupStoreDetails.GetBackupStoreDetails() *BackupStoreDetails`
- New function `*ImportFromStorageResponseType.GetOperationProgressResponseType() *OperationProgressResponseType`
- New function `NewLocationBasedCapabilitySetClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LocationBasedCapabilitySetClient, error)`
- New function `*LocationBasedCapabilitySetClient.Get(context.Context, string, string, *LocationBasedCapabilitySetClientGetOptions) (LocationBasedCapabilitySetClientGetResponse, error)`
- New function `*LocationBasedCapabilitySetClient.NewListPager(string, *LocationBasedCapabilitySetClientListOptions) *runtime.Pager[LocationBasedCapabilitySetClientListResponse]`
- New function `NewLogFilesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LogFilesClient, error)`
- New function `*LogFilesClient.NewListByServerPager(string, string, *LogFilesClientListByServerOptions) *runtime.Pager[LogFilesClientListByServerResponse]`
- New function `NewLongRunningBackupClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LongRunningBackupClient, error)`
- New function `*LongRunningBackupClient.BeginCreate(context.Context, string, string, string, ServerBackupV2, *LongRunningBackupClientBeginCreateOptions) (*runtime.Poller[LongRunningBackupClientCreateResponse], error)`
- New function `NewLongRunningBackupsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LongRunningBackupsClient, error)`
- New function `*LongRunningBackupsClient.Get(context.Context, string, string, string, *LongRunningBackupsClientGetOptions) (LongRunningBackupsClientGetResponse, error)`
- New function `*LongRunningBackupsClient.NewListPager(string, string, *LongRunningBackupsClientListOptions) *runtime.Pager[LongRunningBackupsClientListResponse]`
- New function `NewMaintenancesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MaintenancesClient, error)`
- New function `*MaintenancesClient.NewListPager(string, string, *MaintenancesClientListOptions) *runtime.Pager[MaintenancesClientListResponse]`
- New function `*MaintenancesClient.Read(context.Context, string, string, string, *MaintenancesClientReadOptions) (MaintenancesClientReadResponse, error)`
- New function `*MaintenancesClient.BeginUpdate(context.Context, string, string, string, MaintenanceUpdate, *MaintenancesClientBeginUpdateOptions) (*runtime.Poller[MaintenancesClientUpdateResponse], error)`
- New function `NewOperationProgressClient(string, azcore.TokenCredential, *arm.ClientOptions) (*OperationProgressClient, error)`
- New function `*OperationProgressClient.Get(context.Context, string, string, *OperationProgressClientGetOptions) (OperationProgressClientGetResponse, error)`
- New function `*OperationProgressResponseType.GetOperationProgressResponseType() *OperationProgressResponseType`
- New function `NewOperationResultsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*OperationResultsClient, error)`
- New function `*OperationResultsClient.Get(context.Context, string, string, *OperationResultsClientGetOptions) (OperationResultsClientGetResponse, error)`
- New function `NewPrivateEndpointConnectionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateEndpointConnectionsClient, error)`
- New function `*PrivateEndpointConnectionsClient.BeginCreateOrUpdate(context.Context, string, string, string, PrivateEndpointConnection, *PrivateEndpointConnectionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[PrivateEndpointConnectionsClientCreateOrUpdateResponse], error)`
- New function `*PrivateEndpointConnectionsClient.BeginDelete(context.Context, string, string, string, *PrivateEndpointConnectionsClientBeginDeleteOptions) (*runtime.Poller[PrivateEndpointConnectionsClientDeleteResponse], error)`
- New function `*PrivateEndpointConnectionsClient.Get(context.Context, string, string, string, *PrivateEndpointConnectionsClientGetOptions) (PrivateEndpointConnectionsClientGetResponse, error)`
- New function `*PrivateEndpointConnectionsClient.ListByServer(context.Context, string, string, *PrivateEndpointConnectionsClientListByServerOptions) (PrivateEndpointConnectionsClientListByServerResponse, error)`
- New function `NewPrivateLinkResourcesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateLinkResourcesClient, error)`
- New function `*PrivateLinkResourcesClient.Get(context.Context, string, string, string, *PrivateLinkResourcesClientGetOptions) (PrivateLinkResourcesClientGetResponse, error)`
- New function `*PrivateLinkResourcesClient.NewListByServerPager(string, string, *PrivateLinkResourcesClientListByServerOptions) *runtime.Pager[PrivateLinkResourcesClientListByServerResponse]`
- New function `*ServersClient.BeginDetachVNet(context.Context, string, string, ServerDetachVNetParameter, *ServersClientBeginDetachVNetOptions) (*runtime.Poller[ServersClientDetachVNetResponse], error)`
- New function `*ServersClient.BeginResetGtid(context.Context, string, string, ServerGtidSetParameter, *ServersClientBeginResetGtidOptions) (*runtime.Poller[ServersClientResetGtidResponse], error)`
- New function `*ServersClient.ValidateEstimateHighAvailability(context.Context, string, string, HighAvailabilityValidationEstimation, *ServersClientValidateEstimateHighAvailabilityOptions) (ServersClientValidateEstimateHighAvailabilityResponse, error)`
- New function `NewServersMigrationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ServersMigrationClient, error)`
- New function `*ServersMigrationClient.BeginCutoverMigration(context.Context, string, string, *ServersMigrationClientBeginCutoverMigrationOptions) (*runtime.Poller[ServersMigrationClientCutoverMigrationResponse], error)`
- New struct `AdministratorListResult`
- New struct `AdministratorProperties`
- New struct `AdvancedThreatProtection`
- New struct `AdvancedThreatProtectionForUpdate`
- New struct `AdvancedThreatProtectionListResult`
- New struct `AdvancedThreatProtectionProperties`
- New struct `AdvancedThreatProtectionUpdateProperties`
- New struct `AzureADAdministrator`
- New struct `BackupAndExportRequest`
- New struct `BackupAndExportResponse`
- New struct `BackupAndExportResponseProperties`
- New struct `BackupAndExportResponseType`
- New struct `BackupSettings`
- New struct `Capability`
- New struct `CapabilityPropertiesV2`
- New struct `CapabilitySetsList`
- New struct `ErrorDetail`
- New struct `FeatureProperty`
- New struct `FullBackupStoreDetails`
- New struct `HighAvailabilityValidationEstimation`
- New struct `ImportFromStorageResponseType`
- New struct `ImportSourceProperties`
- New struct `LogFile`
- New struct `LogFileListResult`
- New struct `LogFileProperties`
- New struct `Maintenance`
- New struct `MaintenanceListResult`
- New struct `MaintenancePolicy`
- New struct `MaintenanceProperties`
- New struct `MaintenancePropertiesForUpdate`
- New struct `MaintenanceUpdate`
- New struct `MySQLServerIdentity`
- New struct `MySQLServerSKU`
- New struct `OperationProgressResult`
- New struct `OperationStatusExtendedResult`
- New struct `OperationStatusResult`
- New struct `PrivateEndpoint`
- New struct `PrivateEndpointConnection`
- New struct `PrivateEndpointConnectionListResult`
- New struct `PrivateEndpointConnectionProperties`
- New struct `PrivateLinkResource`
- New struct `PrivateLinkResourceListResult`
- New struct `PrivateLinkResourceProperties`
- New struct `PrivateLinkServiceConnectionState`
- New struct `SKUCapabilityV2`
- New struct `ServerBackupPropertiesV2`
- New struct `ServerBackupV2`
- New struct `ServerBackupV2ListResult`
- New struct `ServerDetachVNetParameter`
- New struct `ServerEditionCapabilityV2`
- New struct `ServerGtidSetParameter`
- New struct `ServerVersionCapabilityV2`
- New struct `ValidateBackupResponse`
- New struct `ValidateBackupResponseProperties`
- New field `BackupIntervalHours` in struct `Backup`
- New field `ResetAllToDefault` in struct `ConfigurationListForBatchUpdate`
- New field `CurrentValue`, `DocumentationLink` in struct `ConfigurationProperties`
- New field `Keyword`, `Page`, `PageSize`, `Tags` in struct `ConfigurationsClientListByServerOptions`
- New field `BatchOfMaintenance` in struct `MaintenanceWindow`
- New field `DatabasePort`, `FullVersion`, `ImportSourceProperties`, `MaintenancePolicy`, `PrivateEndpointConnections` in struct `ServerProperties`
- New field `MaintenancePolicy`, `Network`, `Version` in struct `ServerPropertiesForUpdate`
- New field `AutoIoScaling`, `LogOnDisk`, `StorageRedundancy` in struct `Storage`
- New field `MaxBackupIntervalHours`, `MinBackupIntervalHours` in struct `StorageEditionCapability`
- New field `Location`, `SubscriptionID` in struct `VirtualNetworkSubnetUsageResult`


## 2.0.0-beta.4 (2025-02-27)
### Features Added

- New enum type `PatchStrategy` with values `PatchStrategyRegular`, `PatchStrategyVirtualCanary`
- New enum type `StorageRedundancyEnum` with values `StorageRedundancyEnumLocalRedundancy`, `StorageRedundancyEnumZoneRedundancy`
- New function `*ServersClient.BeginDetachVNet(context.Context, string, string, ServerDetachVNetParameter, *ServersClientBeginDetachVNetOptions) (*runtime.Poller[ServersClientDetachVNetResponse], error)`
- New struct `FeatureProperty`
- New struct `MaintenancePolicy`
- New struct `ServerDetachVNetParameter`
- New field `SupportedFeatures` in struct `CapabilityPropertiesV2`
- New field `DatabasePort`, `FullVersion`, `MaintenancePolicy` in struct `ServerProperties`
- New field `MaintenancePolicy` in struct `ServerPropertiesForUpdate`
- New field `StorageRedundancy` in struct `Storage`


## 2.0.0-beta.3 (2024-04-30)
### Breaking Changes

- Type of `BackupAndExportResponse.Error` has been changed from `*ErrorResponse` to `*ErrorDetail`
- Type of `Server.Identity` has been changed from `*Identity` to `*MySQLServerIdentity`
- Type of `Server.SKU` has been changed from `*SKU` to `*MySQLServerSKU`
- Type of `ServerForUpdate.Identity` has been changed from `*Identity` to `*MySQLServerIdentity`
- Type of `ServerForUpdate.SKU` has been changed from `*SKU` to `*MySQLServerSKU`
- `BackupFormatNone` from enum `BackupFormat` has been removed
- Enum `SKUTier` has been removed
- Struct `Identity` has been removed
- Struct `SKU` has been removed
- Field `AdditionalInfo`, `Code`, `Details`, `Message`, `Target` of struct `ErrorResponse` has been removed

### Features Added

- New value `BackupFormatRaw` added to enum type `BackupFormat`
- New enum type `AdvancedThreatProtectionName` with values `AdvancedThreatProtectionNameDefault`
- New enum type `AdvancedThreatProtectionProvisioningState` with values `AdvancedThreatProtectionProvisioningStateCanceled`, `AdvancedThreatProtectionProvisioningStateFailed`, `AdvancedThreatProtectionProvisioningStateSucceeded`, `AdvancedThreatProtectionProvisioningStateUpdating`
- New enum type `AdvancedThreatProtectionState` with values `AdvancedThreatProtectionStateDisabled`, `AdvancedThreatProtectionStateEnabled`
- New enum type `BackupType` with values `BackupTypeFULL`
- New enum type `ImportSourceStorageType` with values `ImportSourceStorageTypeAzureBlob`
- New enum type `MaintenanceProvisioningState` with values `MaintenanceProvisioningStateCreating`, `MaintenanceProvisioningStateDeleting`, `MaintenanceProvisioningStateFailed`, `MaintenanceProvisioningStateSucceeded`
- New enum type `MaintenanceState` with values `MaintenanceStateCanceled`, `MaintenanceStateCompleted`, `MaintenanceStateInPreparation`, `MaintenanceStateProcessing`, `MaintenanceStateReScheduled`, `MaintenanceStateScheduled`
- New enum type `MaintenanceType` with values `MaintenanceTypeHotFixes`, `MaintenanceTypeMinorVersionUpgrade`, `MaintenanceTypeRoutineMaintenance`, `MaintenanceTypeSecurityPatches`
- New enum type `ObjectType` with values `ObjectTypeBackupAndExportResponse`, `ObjectTypeImportFromStorageResponse`
- New enum type `PrivateEndpointConnectionProvisioningState` with values `PrivateEndpointConnectionProvisioningStateCreating`, `PrivateEndpointConnectionProvisioningStateDeleting`, `PrivateEndpointConnectionProvisioningStateFailed`, `PrivateEndpointConnectionProvisioningStateSucceeded`
- New enum type `PrivateEndpointServiceConnectionStatus` with values `PrivateEndpointServiceConnectionStatusApproved`, `PrivateEndpointServiceConnectionStatusPending`, `PrivateEndpointServiceConnectionStatusRejected`
- New enum type `ProvisioningState` with values `ProvisioningStateCanceled`, `ProvisioningStateCreating`, `ProvisioningStateDeleting`, `ProvisioningStateFailed`, `ProvisioningStateSucceeded`
- New enum type `ServerSKUTier` with values `ServerSKUTierBurstable`, `ServerSKUTierGeneralPurpose`, `ServerSKUTierMemoryOptimized`
- New function `NewAdvancedThreatProtectionSettingsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AdvancedThreatProtectionSettingsClient, error)`
- New function `*AdvancedThreatProtectionSettingsClient.Get(context.Context, string, string, AdvancedThreatProtectionName, *AdvancedThreatProtectionSettingsClientGetOptions) (AdvancedThreatProtectionSettingsClientGetResponse, error)`
- New function `*AdvancedThreatProtectionSettingsClient.NewListPager(string, string, *AdvancedThreatProtectionSettingsClientListOptions) *runtime.Pager[AdvancedThreatProtectionSettingsClientListResponse]`
- New function `*AdvancedThreatProtectionSettingsClient.BeginUpdate(context.Context, string, string, AdvancedThreatProtectionName, AdvancedThreatProtectionForUpdate, *AdvancedThreatProtectionSettingsClientBeginUpdateOptions) (*runtime.Poller[AdvancedThreatProtectionSettingsClientUpdateResponse], error)`
- New function `*AdvancedThreatProtectionSettingsClient.BeginUpdatePut(context.Context, string, string, AdvancedThreatProtectionName, AdvancedThreatProtection, *AdvancedThreatProtectionSettingsClientBeginUpdatePutOptions) (*runtime.Poller[AdvancedThreatProtectionSettingsClientUpdatePutResponse], error)`
- New function `*BackupAndExportResponseType.GetOperationProgressResponseType() *OperationProgressResponseType`
- New function `*ClientFactory.NewAdvancedThreatProtectionSettingsClient() *AdvancedThreatProtectionSettingsClient`
- New function `*ClientFactory.NewLocationBasedCapabilitySetClient() *LocationBasedCapabilitySetClient`
- New function `*ClientFactory.NewLongRunningBackupClient() *LongRunningBackupClient`
- New function `*ClientFactory.NewLongRunningBackupsClient() *LongRunningBackupsClient`
- New function `*ClientFactory.NewMaintenancesClient() *MaintenancesClient`
- New function `*ClientFactory.NewOperationProgressClient() *OperationProgressClient`
- New function `*ClientFactory.NewOperationResultsClient() *OperationResultsClient`
- New function `*ClientFactory.NewServersMigrationClient() *ServersMigrationClient`
- New function `*ImportFromStorageResponseType.GetOperationProgressResponseType() *OperationProgressResponseType`
- New function `NewLocationBasedCapabilitySetClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LocationBasedCapabilitySetClient, error)`
- New function `*LocationBasedCapabilitySetClient.Get(context.Context, string, string, *LocationBasedCapabilitySetClientGetOptions) (LocationBasedCapabilitySetClientGetResponse, error)`
- New function `*LocationBasedCapabilitySetClient.NewListPager(string, *LocationBasedCapabilitySetClientListOptions) *runtime.Pager[LocationBasedCapabilitySetClientListResponse]`
- New function `NewLongRunningBackupClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LongRunningBackupClient, error)`
- New function `*LongRunningBackupClient.BeginCreate(context.Context, string, string, string, *LongRunningBackupClientBeginCreateOptions) (*runtime.Poller[LongRunningBackupClientCreateResponse], error)`
- New function `NewLongRunningBackupsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LongRunningBackupsClient, error)`
- New function `*LongRunningBackupsClient.Get(context.Context, string, string, string, *LongRunningBackupsClientGetOptions) (LongRunningBackupsClientGetResponse, error)`
- New function `*LongRunningBackupsClient.NewListPager(string, string, *LongRunningBackupsClientListOptions) *runtime.Pager[LongRunningBackupsClientListResponse]`
- New function `NewMaintenancesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*MaintenancesClient, error)`
- New function `*MaintenancesClient.NewListPager(string, string, *MaintenancesClientListOptions) *runtime.Pager[MaintenancesClientListResponse]`
- New function `*MaintenancesClient.Read(context.Context, string, string, string, *MaintenancesClientReadOptions) (MaintenancesClientReadResponse, error)`
- New function `*MaintenancesClient.BeginUpdate(context.Context, string, string, string, *MaintenancesClientBeginUpdateOptions) (*runtime.Poller[MaintenancesClientUpdateResponse], error)`
- New function `NewOperationProgressClient(string, azcore.TokenCredential, *arm.ClientOptions) (*OperationProgressClient, error)`
- New function `*OperationProgressClient.Get(context.Context, string, string, *OperationProgressClientGetOptions) (OperationProgressClientGetResponse, error)`
- New function `*OperationProgressResponseType.GetOperationProgressResponseType() *OperationProgressResponseType`
- New function `NewOperationResultsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*OperationResultsClient, error)`
- New function `*OperationResultsClient.Get(context.Context, string, string, *OperationResultsClientGetOptions) (OperationResultsClientGetResponse, error)`
- New function `*ServersClient.ValidateEstimateHighAvailability(context.Context, string, string, HighAvailabilityValidationEstimation, *ServersClientValidateEstimateHighAvailabilityOptions) (ServersClientValidateEstimateHighAvailabilityResponse, error)`
- New function `NewServersMigrationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ServersMigrationClient, error)`
- New function `*ServersMigrationClient.BeginCutoverMigration(context.Context, string, string, *ServersMigrationClientBeginCutoverMigrationOptions) (*runtime.Poller[ServersMigrationClientCutoverMigrationResponse], error)`
- New struct `AdvancedThreatProtection`
- New struct `AdvancedThreatProtectionForUpdate`
- New struct `AdvancedThreatProtectionListResult`
- New struct `AdvancedThreatProtectionProperties`
- New struct `AdvancedThreatProtectionUpdateProperties`
- New struct `BackupAndExportResponseType`
- New struct `Capability`
- New struct `CapabilityPropertiesV2`
- New struct `CapabilitySetsList`
- New struct `ErrorDetail`
- New struct `HighAvailabilityValidationEstimation`
- New struct `ImportFromStorageResponseType`
- New struct `ImportSourceProperties`
- New struct `Maintenance`
- New struct `MaintenanceListResult`
- New struct `MaintenanceProperties`
- New struct `MaintenancePropertiesForUpdate`
- New struct `MaintenanceUpdate`
- New struct `MySQLServerIdentity`
- New struct `MySQLServerSKU`
- New struct `OperationProgressResult`
- New struct `OperationStatusExtendedResult`
- New struct `OperationStatusResult`
- New struct `PrivateEndpoint`
- New struct `PrivateEndpointConnection`
- New struct `PrivateEndpointConnectionProperties`
- New struct `PrivateLinkServiceConnectionState`
- New struct `Provisioning`
- New struct `ProvisioningStateProperties`
- New struct `SKUCapabilityV2`
- New struct `ServerBackupPropertiesV2`
- New struct `ServerBackupV2`
- New struct `ServerBackupV2ListResult`
- New struct `ServerEditionCapabilityV2`
- New struct `ServerVersionCapabilityV2`
- New field `BackupIntervalHours` in struct `Backup`
- New field `SystemData` in struct `BackupAndExportResponse`
- New field `Error` in struct `ErrorResponse`
- New field `SystemData` in struct `ProxyResource`
- New field `SystemData` in struct `Resource`
- New field `ImportSourceProperties`, `PrivateEndpointConnections` in struct `ServerProperties`
- New field `Network` in struct `ServerPropertiesForUpdate`
- New field `MaxBackupIntervalHours`, `MinBackupIntervalHours` in struct `StorageEditionCapability`
- New field `SystemData` in struct `TrackedResource`


## 2.0.0-beta.2 (2023-11-30)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0-beta.1 (2023-05-26)
### Breaking Changes

- Type of `Identity.Type` has been changed from `*string` to `*ManagedServiceIdentityType`

### Features Added

- New enum type `AdministratorName` with values `AdministratorNameActiveDirectory`
- New enum type `AdministratorType` with values `AdministratorTypeActiveDirectory`
- New enum type `BackupFormat` with values `BackupFormatCollatedFormat`, `BackupFormatNone`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeUserAssigned`
- New enum type `OperationStatus` with values `OperationStatusCancelInProgress`, `OperationStatusCanceled`, `OperationStatusFailed`, `OperationStatusInProgress`, `OperationStatusPending`, `OperationStatusSucceeded`
- New enum type `ResetAllToDefault` with values `ResetAllToDefaultFalse`, `ResetAllToDefaultTrue`
- New function `NewAzureADAdministratorsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AzureADAdministratorsClient, error)`
- New function `*AzureADAdministratorsClient.BeginCreateOrUpdate(context.Context, string, string, AdministratorName, AzureADAdministrator, *AzureADAdministratorsClientBeginCreateOrUpdateOptions) (*runtime.Poller[AzureADAdministratorsClientCreateOrUpdateResponse], error)`
- New function `*AzureADAdministratorsClient.BeginDelete(context.Context, string, string, AdministratorName, *AzureADAdministratorsClientBeginDeleteOptions) (*runtime.Poller[AzureADAdministratorsClientDeleteResponse], error)`
- New function `*AzureADAdministratorsClient.Get(context.Context, string, string, AdministratorName, *AzureADAdministratorsClientGetOptions) (AzureADAdministratorsClientGetResponse, error)`
- New function `*AzureADAdministratorsClient.NewListByServerPager(string, string, *AzureADAdministratorsClientListByServerOptions) *runtime.Pager[AzureADAdministratorsClientListByServerResponse]`
- New function `NewBackupAndExportClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BackupAndExportClient, error)`
- New function `*BackupAndExportClient.BeginCreate(context.Context, string, string, BackupAndExportRequest, *BackupAndExportClientBeginCreateOptions) (*runtime.Poller[BackupAndExportClientCreateResponse], error)`
- New function `*BackupAndExportClient.ValidateBackup(context.Context, string, string, *BackupAndExportClientValidateBackupOptions) (BackupAndExportClientValidateBackupResponse, error)`
- New function `*BackupStoreDetails.GetBackupStoreDetails() *BackupStoreDetails`
- New function `*BackupsClient.Put(context.Context, string, string, string, *BackupsClientPutOptions) (BackupsClientPutResponse, error)`
- New function `NewCheckNameAvailabilityWithoutLocationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CheckNameAvailabilityWithoutLocationClient, error)`
- New function `*CheckNameAvailabilityWithoutLocationClient.Execute(context.Context, NameAvailabilityRequest, *CheckNameAvailabilityWithoutLocationClientExecuteOptions) (CheckNameAvailabilityWithoutLocationClientExecuteResponse, error)`
- New function `*ClientFactory.NewAzureADAdministratorsClient() *AzureADAdministratorsClient`
- New function `*ClientFactory.NewBackupAndExportClient() *BackupAndExportClient`
- New function `*ClientFactory.NewCheckNameAvailabilityWithoutLocationClient() *CheckNameAvailabilityWithoutLocationClient`
- New function `*ClientFactory.NewLogFilesClient() *LogFilesClient`
- New function `*ConfigurationsClient.BeginCreateOrUpdate(context.Context, string, string, string, Configuration, *ConfigurationsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ConfigurationsClientCreateOrUpdateResponse], error)`
- New function `*FullBackupStoreDetails.GetBackupStoreDetails() *BackupStoreDetails`
- New function `NewLogFilesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*LogFilesClient, error)`
- New function `*LogFilesClient.NewListByServerPager(string, string, *LogFilesClientListByServerOptions) *runtime.Pager[LogFilesClientListByServerResponse]`
- New function `*ServersClient.BeginResetGtid(context.Context, string, string, ServerGtidSetParameter, *ServersClientBeginResetGtidOptions) (*runtime.Poller[ServersClientResetGtidResponse], error)`
- New struct `AdministratorListResult`
- New struct `AdministratorProperties`
- New struct `AzureADAdministrator`
- New struct `BackupAndExportRequest`
- New struct `BackupAndExportResponse`
- New struct `BackupAndExportResponseProperties`
- New struct `BackupRequestBase`
- New struct `BackupSettings`
- New struct `FullBackupStoreDetails`
- New struct `LogFile`
- New struct `LogFileListResult`
- New struct `LogFileProperties`
- New struct `ServerGtidSetParameter`
- New struct `ValidateBackupResponse`
- New struct `ValidateBackupResponseProperties`
- New field `ResetAllToDefault` in struct `ConfigurationListForBatchUpdate`
- New field `CurrentValue`, `DocumentationLink` in struct `ConfigurationProperties`
- New field `Keyword`, `Page`, `PageSize`, `Tags` in struct `ConfigurationsClientListByServerOptions`
- New field `Version` in struct `ServerPropertiesForUpdate`
- New field `AutoIoScaling`, `LogOnDisk` in struct `Storage`
- New field `Location`, `SubscriptionID` in struct `VirtualNetworkSubnetUsageResult`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mysql/armmysqlflexibleservers` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).