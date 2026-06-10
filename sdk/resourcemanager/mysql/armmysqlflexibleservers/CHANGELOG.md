# Release History

## 2.0.0 (2026-06-10)
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
- New function `NewAdvancedThreatProtectionSettingsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AdvancedThreatProtectionSettingsClient, error)`
- New function `*AdvancedThreatProtectionSettingsClient.Get(ctx context.Context, resourceGroupName string, serverName string, advancedThreatProtectionName AdvancedThreatProtectionName, options *AdvancedThreatProtectionSettingsClientGetOptions) (AdvancedThreatProtectionSettingsClientGetResponse, error)`
- New function `*AdvancedThreatProtectionSettingsClient.NewListPager(resourceGroupName string, serverName string, options *AdvancedThreatProtectionSettingsClientListOptions) *runtime.Pager[AdvancedThreatProtectionSettingsClientListResponse]`
- New function `*AdvancedThreatProtectionSettingsClient.BeginUpdate(ctx context.Context, resourceGroupName string, serverName string, advancedThreatProtectionName AdvancedThreatProtectionName, parameters AdvancedThreatProtectionForUpdate, options *AdvancedThreatProtectionSettingsClientBeginUpdateOptions) (*runtime.Poller[AdvancedThreatProtectionSettingsClientUpdateResponse], error)`
- New function `*AdvancedThreatProtectionSettingsClient.BeginUpdatePut(ctx context.Context, resourceGroupName string, serverName string, advancedThreatProtectionName AdvancedThreatProtectionName, parameters AdvancedThreatProtection, options *AdvancedThreatProtectionSettingsClientBeginUpdatePutOptions) (*runtime.Poller[AdvancedThreatProtectionSettingsClientUpdatePutResponse], error)`
- New function `NewAzureADAdministratorsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AzureADAdministratorsClient, error)`
- New function `*AzureADAdministratorsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, serverName string, administratorName AdministratorName, parameters AzureADAdministrator, options *AzureADAdministratorsClientBeginCreateOrUpdateOptions) (*runtime.Poller[AzureADAdministratorsClientCreateOrUpdateResponse], error)`
- New function `*AzureADAdministratorsClient.BeginDelete(ctx context.Context, resourceGroupName string, serverName string, administratorName AdministratorName, options *AzureADAdministratorsClientBeginDeleteOptions) (*runtime.Poller[AzureADAdministratorsClientDeleteResponse], error)`
- New function `*AzureADAdministratorsClient.Get(ctx context.Context, resourceGroupName string, serverName string, administratorName AdministratorName, options *AzureADAdministratorsClientGetOptions) (AzureADAdministratorsClientGetResponse, error)`
- New function `*AzureADAdministratorsClient.NewListByServerPager(resourceGroupName string, serverName string, options *AzureADAdministratorsClientListByServerOptions) *runtime.Pager[AzureADAdministratorsClientListByServerResponse]`
- New function `NewBackupAndExportClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*BackupAndExportClient, error)`
- New function `*BackupAndExportClient.BeginCreate(ctx context.Context, resourceGroupName string, serverName string, parameters BackupAndExportRequest, options *BackupAndExportClientBeginCreateOptions) (*runtime.Poller[BackupAndExportClientCreateResponse], error)`
- New function `*BackupAndExportClient.ValidateBackup(ctx context.Context, resourceGroupName string, serverName string, options *BackupAndExportClientValidateBackupOptions) (BackupAndExportClientValidateBackupResponse, error)`
- New function `*BackupAndExportResponseType.GetOperationProgressResponseType() *OperationProgressResponseType`
- New function `*BackupStoreDetails.GetBackupStoreDetails() *BackupStoreDetails`
- New function `*BackupsClient.Put(ctx context.Context, resourceGroupName string, serverName string, backupName string, options *BackupsClientPutOptions) (BackupsClientPutResponse, error)`
- New function `NewCheckNameAvailabilityWithoutLocationClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CheckNameAvailabilityWithoutLocationClient, error)`
- New function `*CheckNameAvailabilityWithoutLocationClient.Execute(ctx context.Context, nameAvailabilityRequest NameAvailabilityRequest, options *CheckNameAvailabilityWithoutLocationClientExecuteOptions) (CheckNameAvailabilityWithoutLocationClientExecuteResponse, error)`
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
- New function `*ConfigurationsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, serverName string, configurationName string, parameters Configuration, options *ConfigurationsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ConfigurationsClientCreateOrUpdateResponse], error)`
- New function `*FullBackupStoreDetails.GetBackupStoreDetails() *BackupStoreDetails`
- New function `*ImportFromStorageResponseType.GetOperationProgressResponseType() *OperationProgressResponseType`
- New function `NewLocationBasedCapabilitySetClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*LocationBasedCapabilitySetClient, error)`
- New function `*LocationBasedCapabilitySetClient.Get(ctx context.Context, locationName string, capabilitySetName string, options *LocationBasedCapabilitySetClientGetOptions) (LocationBasedCapabilitySetClientGetResponse, error)`
- New function `*LocationBasedCapabilitySetClient.NewListPager(locationName string, options *LocationBasedCapabilitySetClientListOptions) *runtime.Pager[LocationBasedCapabilitySetClientListResponse]`
- New function `NewLogFilesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*LogFilesClient, error)`
- New function `*LogFilesClient.NewListByServerPager(resourceGroupName string, serverName string, options *LogFilesClientListByServerOptions) *runtime.Pager[LogFilesClientListByServerResponse]`
- New function `NewLongRunningBackupClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*LongRunningBackupClient, error)`
- New function `*LongRunningBackupClient.BeginCreate(ctx context.Context, resourceGroupName string, serverName string, backupName string, parameters ServerBackupV2, options *LongRunningBackupClientBeginCreateOptions) (*runtime.Poller[LongRunningBackupClientCreateResponse], error)`
- New function `NewLongRunningBackupsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*LongRunningBackupsClient, error)`
- New function `*LongRunningBackupsClient.Get(ctx context.Context, resourceGroupName string, serverName string, backupName string, options *LongRunningBackupsClientGetOptions) (LongRunningBackupsClientGetResponse, error)`
- New function `*LongRunningBackupsClient.NewListPager(resourceGroupName string, serverName string, options *LongRunningBackupsClientListOptions) *runtime.Pager[LongRunningBackupsClientListResponse]`
- New function `NewMaintenancesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*MaintenancesClient, error)`
- New function `*MaintenancesClient.NewListPager(resourceGroupName string, serverName string, options *MaintenancesClientListOptions) *runtime.Pager[MaintenancesClientListResponse]`
- New function `*MaintenancesClient.Read(ctx context.Context, resourceGroupName string, serverName string, maintenanceName string, options *MaintenancesClientReadOptions) (MaintenancesClientReadResponse, error)`
- New function `*MaintenancesClient.BeginUpdate(ctx context.Context, resourceGroupName string, serverName string, maintenanceName string, parameters MaintenanceUpdate, options *MaintenancesClientBeginUpdateOptions) (*runtime.Poller[MaintenancesClientUpdateResponse], error)`
- New function `NewOperationProgressClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*OperationProgressClient, error)`
- New function `*OperationProgressClient.Get(ctx context.Context, locationName string, operationID string, options *OperationProgressClientGetOptions) (OperationProgressClientGetResponse, error)`
- New function `*OperationProgressResponseType.GetOperationProgressResponseType() *OperationProgressResponseType`
- New function `NewOperationResultsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*OperationResultsClient, error)`
- New function `*OperationResultsClient.Get(ctx context.Context, locationName string, operationID string, options *OperationResultsClientGetOptions) (OperationResultsClientGetResponse, error)`
- New function `NewPrivateEndpointConnectionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PrivateEndpointConnectionsClient, error)`
- New function `*PrivateEndpointConnectionsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, serverName string, privateEndpointConnectionName string, parameters PrivateEndpointConnection, options *PrivateEndpointConnectionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[PrivateEndpointConnectionsClientCreateOrUpdateResponse], error)`
- New function `*PrivateEndpointConnectionsClient.BeginDelete(ctx context.Context, resourceGroupName string, serverName string, privateEndpointConnectionName string, options *PrivateEndpointConnectionsClientBeginDeleteOptions) (*runtime.Poller[PrivateEndpointConnectionsClientDeleteResponse], error)`
- New function `*PrivateEndpointConnectionsClient.Get(ctx context.Context, resourceGroupName string, serverName string, privateEndpointConnectionName string, options *PrivateEndpointConnectionsClientGetOptions) (PrivateEndpointConnectionsClientGetResponse, error)`
- New function `*PrivateEndpointConnectionsClient.ListByServer(ctx context.Context, resourceGroupName string, serverName string, options *PrivateEndpointConnectionsClientListByServerOptions) (PrivateEndpointConnectionsClientListByServerResponse, error)`
- New function `NewPrivateLinkResourcesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PrivateLinkResourcesClient, error)`
- New function `*PrivateLinkResourcesClient.Get(ctx context.Context, resourceGroupName string, serverName string, groupName string, options *PrivateLinkResourcesClientGetOptions) (PrivateLinkResourcesClientGetResponse, error)`
- New function `*PrivateLinkResourcesClient.NewListByServerPager(resourceGroupName string, serverName string, options *PrivateLinkResourcesClientListByServerOptions) *runtime.Pager[PrivateLinkResourcesClientListByServerResponse]`
- New function `*ServersClient.BeginDetachVNet(ctx context.Context, resourceGroupName string, serverName string, parameters ServerDetachVNetParameter, options *ServersClientBeginDetachVNetOptions) (*runtime.Poller[ServersClientDetachVNetResponse], error)`
- New function `*ServersClient.BeginResetGtid(ctx context.Context, resourceGroupName string, serverName string, parameters ServerGtidSetParameter, options *ServersClientBeginResetGtidOptions) (*runtime.Poller[ServersClientResetGtidResponse], error)`
- New function `*ServersClient.ValidateEstimateHighAvailability(ctx context.Context, resourceGroupName string, serverName string, parameters HighAvailabilityValidationEstimation, options *ServersClientValidateEstimateHighAvailabilityOptions) (ServersClientValidateEstimateHighAvailabilityResponse, error)`
- New function `NewServersMigrationClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ServersMigrationClient, error)`
- New function `*ServersMigrationClient.BeginCutoverMigration(ctx context.Context, resourceGroupName string, serverName string, options *ServersMigrationClientBeginCutoverMigrationOptions) (*runtime.Poller[ServersMigrationClientCutoverMigrationResponse], error)`
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


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


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