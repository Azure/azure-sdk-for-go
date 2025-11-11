# Release History

## 5.0.0 (2025-11-11)
### Breaking Changes

- Function `*ConfigurationsClient.BeginPut` parameter(s) have been changed from `(context.Context, string, string, string, Configuration, *ConfigurationsClientBeginPutOptions)` to `(context.Context, string, string, string, ConfigurationForUpdate, *ConfigurationsClientBeginPutOptions)`
- Function `NewMigrationsClient` parameter(s) have been changed from `(azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*MigrationsClient.Create` parameter(s) have been changed from `(context.Context, string, string, string, string, MigrationResource, *MigrationsClientCreateOptions)` to `(context.Context, string, string, string, Migration, *MigrationsClientCreateOptions)`
- Function `*MigrationsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, string, *MigrationsClientGetOptions)` to `(context.Context, string, string, string, *MigrationsClientGetOptions)`
- Function `*MigrationsClient.NewListByTargetServerPager` parameter(s) have been changed from `(string, string, string, *MigrationsClientListByTargetServerOptions)` to `(string, string, *MigrationsClientListByTargetServerOptions)`
- Function `*MigrationsClient.Update` parameter(s) have been changed from `(context.Context, string, string, string, string, MigrationResourceForPatch, *MigrationsClientUpdateOptions)` to `(context.Context, string, string, string, MigrationResourceForPatch, *MigrationsClientUpdateOptions)`
- Function `*ServerThreatProtectionSettingsClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, ThreatProtectionName, ServerThreatProtectionSettingsModel, *ServerThreatProtectionSettingsClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, ThreatProtectionName, AdvancedThreatProtectionSettingsModel, *ServerThreatProtectionSettingsClientBeginCreateOrUpdateOptions)`
- Function `*ServersClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, ServerForUpdate, *ServersClientBeginUpdateOptions)` to `(context.Context, string, string, ServerForPatch, *ServersClientBeginUpdateOptions)`
- Function `*VirtualEndpointsClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, string, VirtualEndpointResource, *VirtualEndpointsClientBeginCreateOptions)` to `(context.Context, string, string, string, VirtualEndpoint, *VirtualEndpointsClientBeginCreateOptions)`
- Type of `AuthConfig.ActiveDirectoryAuth` has been changed from `*ActiveDirectoryAuthEnum` to `*MicrosoftEntraAuth`
- Type of `AuthConfig.PasswordAuth` has been changed from `*PasswordAuthEnum` to `*PasswordBasedAuth`
- Type of `Backup.GeoRedundantBackup` has been changed from `*GeoRedundantBackupEnum` to `*GeographicallyRedundantBackup`
- Type of `DataEncryption.GeoBackupEncryptionKeyStatus` has been changed from `*KeyStatusEnum` to `*EncryptionKeyStatus`
- Type of `DataEncryption.PrimaryEncryptionKeyStatus` has been changed from `*KeyStatusEnum` to `*EncryptionKeyStatus`
- Type of `DataEncryption.Type` has been changed from `*ArmServerKeyType` to `*DataEncryptionType`
- Type of `HighAvailability.State` has been changed from `*ServerHAState` to `*HighAvailabilityState`
- Type of `LtrPreBackupResponse.Properties` has been changed from `*LtrPreBackupResponseProperties` to `*BackupsLongTermRetentionResponseProperties`
- Type of `LtrServerBackupOperationList.Value` has been changed from `[]*LtrServerBackupOperation` to `[]*BackupsLongTermRetentionOperation`
- Type of `MigrationResourceForPatch.Properties` has been changed from `*MigrationResourcePropertiesForPatch` to `*MigrationPropertiesForPatch`
- Type of `MigrationStatus.CurrentSubStateDetails` has been changed from `*MigrationSubStateDetails` to `*MigrationSubstateDetails`
- Type of `Replica.PromoteOption` has been changed from `*ReplicationPromoteOption` to `*ReadReplicaPromoteOption`
- Type of `ServerProperties.Version` has been changed from `*ServerVersion` to `*PostgresMajorVersion`
- Type of `ServerSKUCapability.SupportedHaMode` has been changed from `[]*HaMode` to `[]*HighAvailabilityMode`
- Type of `Storage.Tier` has been changed from `*AzureManagedDiskPerformanceTiers` to `*AzureManagedDiskPerformanceTier`
- `HighAvailabilityModeDisabled` from enum `HighAvailabilityMode` has been removed
- Enum `ActiveDirectoryAuthEnum` has been removed
- Enum `ArmServerKeyType` has been removed
- Enum `AzureManagedDiskPerformanceTiers` has been removed
- Enum `CancelEnum` has been removed
- Enum `CreateModeForUpdate` has been removed
- Enum `FastProvisioningSupportedEnum` has been removed
- Enum `GeoBackupSupportedEnum` has been removed
- Enum `GeoRedundantBackupEnum` has been removed
- Enum `HaMode` has been removed
- Enum `KeyStatusEnum` has been removed
- Enum `LogicalReplicationOnSourceDbEnum` has been removed
- Enum `MigrateRolesEnum` has been removed
- Enum `MigrationDbState` has been removed
- Enum `MigrationSubState` has been removed
- Enum `OnlineResizeSupportedEnum` has been removed
- Enum `Origin` has been removed
- Enum `OverwriteDbsInTargetEnum` has been removed
- Enum `PasswordAuthEnum` has been removed
- Enum `ReplicationPromoteOption` has been removed
- Enum `RestrictedEnum` has been removed
- Enum `ServerHAState` has been removed
- Enum `ServerVersion` has been removed
- Enum `StartDataMigrationEnum` has been removed
- Enum `StorageAutoGrowthSupportedEnum` has been removed
- Enum `TriggerCutoverEnum` has been removed
- Enum `ZoneRedundantHaAndGeoBackupSupportedEnum` has been removed
- Enum `ZoneRedundantHaSupportedEnum` has been removed
- Function `NewAdministratorsClient` has been removed
- Function `*AdministratorsClient.BeginCreate` has been removed
- Function `*AdministratorsClient.BeginDelete` has been removed
- Function `*AdministratorsClient.Get` has been removed
- Function `*AdministratorsClient.NewListByServerPager` has been removed
- Function `NewBackupsClient` has been removed
- Function `*BackupsClient.BeginCreate` has been removed
- Function `*BackupsClient.BeginDelete` has been removed
- Function `*BackupsClient.Get` has been removed
- Function `*BackupsClient.NewListByServerPager` has been removed
- Function `NewCheckNameAvailabilityClient` has been removed
- Function `*CheckNameAvailabilityClient.Execute` has been removed
- Function `NewCheckNameAvailabilityWithLocationClient` has been removed
- Function `*CheckNameAvailabilityWithLocationClient.Execute` has been removed
- Function `*ClientFactory.NewAdministratorsClient` has been removed
- Function `*ClientFactory.NewBackupsClient` has been removed
- Function `*ClientFactory.NewCheckNameAvailabilityClient` has been removed
- Function `*ClientFactory.NewCheckNameAvailabilityWithLocationClient` has been removed
- Function `*ClientFactory.NewFlexibleServerClient` has been removed
- Function `*ClientFactory.NewGetPrivateDNSZoneSuffixClient` has been removed
- Function `*ClientFactory.NewLocationBasedCapabilitiesClient` has been removed
- Function `*ClientFactory.NewLogFilesClient` has been removed
- Function `*ClientFactory.NewLtrBackupOperationsClient` has been removed
- Function `*ClientFactory.NewPostgreSQLServerManagementClient` has been removed
- Function `*ClientFactory.NewPrivateEndpointConnectionClient` has been removed
- Function `*ClientFactory.NewServerCapabilitiesClient` has been removed
- Function `NewFlexibleServerClient` has been removed
- Function `*FlexibleServerClient.BeginStartLtrBackup` has been removed
- Function `*FlexibleServerClient.TriggerLtrPreBackup` has been removed
- Function `NewGetPrivateDNSZoneSuffixClient` has been removed
- Function `*GetPrivateDNSZoneSuffixClient.Execute` has been removed
- Function `NewLocationBasedCapabilitiesClient` has been removed
- Function `*LocationBasedCapabilitiesClient.NewExecutePager` has been removed
- Function `NewLogFilesClient` has been removed
- Function `*LogFilesClient.NewListByServerPager` has been removed
- Function `NewLtrBackupOperationsClient` has been removed
- Function `*LtrBackupOperationsClient.Get` has been removed
- Function `*LtrBackupOperationsClient.NewListByServerPager` has been removed
- Function `*MigrationsClient.Delete` has been removed
- Function `NewPostgreSQLServerManagementClient` has been removed
- Function `*PostgreSQLServerManagementClient.CheckMigrationNameAvailability` has been removed
- Function `NewPrivateEndpointConnectionClient` has been removed
- Function `*PrivateEndpointConnectionClient.BeginDelete` has been removed
- Function `*PrivateEndpointConnectionClient.BeginUpdate` has been removed
- Function `NewServerCapabilitiesClient` has been removed
- Function `*ServerCapabilitiesClient.NewListPager` has been removed
- Function `*ServerThreatProtectionSettingsClient.Get` has been removed
- Function `*ServerThreatProtectionSettingsClient.NewListByServerPager` has been removed
- Function `*ServersClient.BeginCreate` has been removed
- Function `*ServersClient.NewListPager` has been removed
- Function `*VirtualNetworkSubnetUsageClient.Execute` has been removed
- Operation `*OperationsClient.List` has supported pagination, use `*OperationsClient.NewListPager` instead.
- Struct `ActiveDirectoryAdministrator` has been removed
- Struct `ActiveDirectoryAdministratorAdd` has been removed
- Struct `AdministratorListResult` has been removed
- Struct `AdministratorProperties` has been removed
- Struct `AdministratorPropertiesForAdd` has been removed
- Struct `CapabilitiesListResult` has been removed
- Struct `ConfigurationListResult` has been removed
- Struct `DatabaseListResult` has been removed
- Struct `DbMigrationStatus` has been removed
- Struct `FirewallRuleListResult` has been removed
- Struct `FlexibleServerCapability` has been removed
- Struct `FlexibleServerEditionCapability` has been removed
- Struct `LogFile` has been removed
- Struct `LogFileListResult` has been removed
- Struct `LogFileProperties` has been removed
- Struct `LtrBackupRequest` has been removed
- Struct `LtrBackupResponse` has been removed
- Struct `LtrPreBackupResponseProperties` has been removed
- Struct `LtrServerBackupOperation` has been removed
- Struct `MigrationNameAvailabilityResource` has been removed
- Struct `MigrationResource` has been removed
- Struct `MigrationResourceListResult` has been removed
- Struct `MigrationResourceProperties` has been removed
- Struct `MigrationResourcePropertiesForPatch` has been removed
- Struct `MigrationSubStateDetails` has been removed
- Struct `NameAvailability` has been removed
- Struct `OperationListResult` has been removed
- Struct `PrivateEndpointConnectionListResult` has been removed
- Struct `PrivateLinkResourceListResult` has been removed
- Struct `ServerBackup` has been removed
- Struct `ServerBackupListResult` has been removed
- Struct `ServerBackupProperties` has been removed
- Struct `ServerForUpdate` has been removed
- Struct `ServerListResult` has been removed
- Struct `ServerPropertiesForUpdate` has been removed
- Struct `ServerThreatProtectionListResult` has been removed
- Struct `ServerThreatProtectionProperties` has been removed
- Struct `ServerThreatProtectionSettingsModel` has been removed
- Struct `VirtualEndpointResource` has been removed
- Struct `VirtualEndpointsListResult` has been removed
- Struct `VirtualNetworkSubnetUsageResult` has been removed
- Field `ConfigurationListResult` of struct `ConfigurationsClientListByServerResponse` has been removed
- Field `Configuration` of struct `ConfigurationsClientPutResponse` has been removed
- Field `Configuration` of struct `ConfigurationsClientUpdateResponse` has been removed
- Field `Database` of struct `DatabasesClientCreateResponse` has been removed
- Field `DatabaseListResult` of struct `DatabasesClientListByServerResponse` has been removed
- Field `FirewallRule` of struct `FirewallRulesClientCreateOrUpdateResponse` has been removed
- Field `FirewallRuleListResult` of struct `FirewallRulesClientListByServerResponse` has been removed
- Field `MigrationResource` of struct `MigrationsClientCreateResponse` has been removed
- Field `MigrationResource` of struct `MigrationsClientGetResponse` has been removed
- Field `MigrationResourceListResult` of struct `MigrationsClientListByTargetServerResponse` has been removed
- Field `MigrationResource` of struct `MigrationsClientUpdateResponse` has been removed
- Field `PrivateEndpointConnectionListResult` of struct `PrivateEndpointConnectionsClientListByServerResponse` has been removed
- Field `PrivateLinkResourceListResult` of struct `PrivateLinkResourcesClientListByServerResponse` has been removed
- Field `ServerListResult` of struct `ReplicasClientListByServerResponse` has been removed
- Field `ServerThreatProtectionSettingsModel` of struct `ServerThreatProtectionSettingsClientCreateOrUpdateResponse` has been removed
- Field `ServerListResult` of struct `ServersClientListByResourceGroupResponse` has been removed
- Field `Server` of struct `ServersClientUpdateResponse` has been removed
- Field `VirtualEndpointResource` of struct `VirtualEndpointsClientCreateResponse` has been removed
- Field `VirtualEndpointResource` of struct `VirtualEndpointsClientGetResponse` has been removed
- Field `VirtualEndpointsListResult` of struct `VirtualEndpointsClientListByServerResponse` has been removed
- Field `VirtualEndpointResource` of struct `VirtualEndpointsClientUpdateResponse` has been removed

### Features Added

- New value `ConfigurationDataTypeSet`, `ConfigurationDataTypeString` added to enum type `ConfigurationDataType`
- New value `ServerStateInaccessible`, `ServerStateProvisioning`, `ServerStateRestarting` added to enum type `ServerState`
- New value `SourceTypeApsaraDBRDS`, `SourceTypeCrunchyPostgreSQL`, `SourceTypeDigitalOceanDroplets`, `SourceTypeDigitalOceanPostgreSQL`, `SourceTypeEDBOracleServer`, `SourceTypeEDBPostgreSQL`, `SourceTypeHerokuPostgreSQL`, `SourceTypeHuaweiCompute`, `SourceTypeHuaweiRDS`, `SourceTypePostgreSQLCosmosDB`, `SourceTypePostgreSQLFlexibleServer`, `SourceTypeSupabasePostgreSQL` added to enum type `SourceType`
- New value `StorageTypeUltraSSDLRS` added to enum type `StorageType`
- New enum type `AzureManagedDiskPerformanceTier` with values `AzureManagedDiskPerformanceTierP1`, `AzureManagedDiskPerformanceTierP10`, `AzureManagedDiskPerformanceTierP15`, `AzureManagedDiskPerformanceTierP2`, `AzureManagedDiskPerformanceTierP20`, `AzureManagedDiskPerformanceTierP3`, `AzureManagedDiskPerformanceTierP30`, `AzureManagedDiskPerformanceTierP4`, `AzureManagedDiskPerformanceTierP40`, `AzureManagedDiskPerformanceTierP50`, `AzureManagedDiskPerformanceTierP6`, `AzureManagedDiskPerformanceTierP60`, `AzureManagedDiskPerformanceTierP70`, `AzureManagedDiskPerformanceTierP80`
- New enum type `BackupType` with values `BackupTypeCustomerOnDemand`, `BackupTypeFull`
- New enum type `Cancel` with values `CancelFalse`, `CancelTrue`
- New enum type `CreateModeForPatch` with values `CreateModeForPatchDefault`, `CreateModeForPatchUpdate`
- New enum type `DataEncryptionType` with values `DataEncryptionTypeAzureKeyVault`, `DataEncryptionTypeSystemManaged`
- New enum type `EncryptionKeyStatus` with values `EncryptionKeyStatusInvalid`, `EncryptionKeyStatusValid`
- New enum type `FastProvisioningSupport` with values `FastProvisioningSupportDisabled`, `FastProvisioningSupportEnabled`
- New enum type `FeatureStatus` with values `FeatureStatusDisabled`, `FeatureStatusEnabled`
- New enum type `GeographicallyRedundantBackup` with values `GeographicallyRedundantBackupDisabled`, `GeographicallyRedundantBackupEnabled`
- New enum type `GeographicallyRedundantBackupSupport` with values `GeographicallyRedundantBackupSupportDisabled`, `GeographicallyRedundantBackupSupportEnabled`
- New enum type `HighAvailabilityState` with values `HighAvailabilityStateCreatingStandby`, `HighAvailabilityStateFailingOver`, `HighAvailabilityStateHealthy`, `HighAvailabilityStateNotEnabled`, `HighAvailabilityStateRemovingStandby`, `HighAvailabilityStateReplicatingData`
- New enum type `LocationRestricted` with values `LocationRestrictedDisabled`, `LocationRestrictedEnabled`
- New enum type `LogicalReplicationOnSourceServer` with values `LogicalReplicationOnSourceServerFalse`, `LogicalReplicationOnSourceServerTrue`
- New enum type `MicrosoftEntraAuth` with values `MicrosoftEntraAuthDisabled`, `MicrosoftEntraAuthEnabled`
- New enum type `MigrateRolesAndPermissions` with values `MigrateRolesAndPermissionsFalse`, `MigrateRolesAndPermissionsTrue`
- New enum type `MigrationDatabaseState` with values `MigrationDatabaseStateCanceled`, `MigrationDatabaseStateCanceling`, `MigrationDatabaseStateFailed`, `MigrationDatabaseStateInProgress`, `MigrationDatabaseStateSucceeded`, `MigrationDatabaseStateWaitingForCutoverTrigger`
- New enum type `MigrationSubstate` with values `MigrationSubstateCancelingRequestedDBMigrations`, `MigrationSubstateCompleted`, `MigrationSubstateCompletingMigration`, `MigrationSubstateMigratingData`, `MigrationSubstatePerformingPreRequisiteSteps`, `MigrationSubstateValidationInProgress`, `MigrationSubstateWaitingForCutoverTrigger`, `MigrationSubstateWaitingForDBsToMigrateSpecification`, `MigrationSubstateWaitingForDataMigrationScheduling`, `MigrationSubstateWaitingForDataMigrationWindow`, `MigrationSubstateWaitingForLogicalReplicationSetupRequestOnSourceDB`, `MigrationSubstateWaitingForTargetDBOverwriteConfirmation`
- New enum type `OnlineStorageResizeSupport` with values `OnlineStorageResizeSupportDisabled`, `OnlineStorageResizeSupportEnabled`
- New enum type `OverwriteDatabasesOnTargetServer` with values `OverwriteDatabasesOnTargetServerFalse`, `OverwriteDatabasesOnTargetServerTrue`
- New enum type `PasswordBasedAuth` with values `PasswordBasedAuthDisabled`, `PasswordBasedAuthEnabled`
- New enum type `PostgresMajorVersion` with values `PostgresMajorVersionEighteen`, `PostgresMajorVersionEleven`, `PostgresMajorVersionFifteen`, `PostgresMajorVersionFourteen`, `PostgresMajorVersionSeventeen`, `PostgresMajorVersionSixteen`, `PostgresMajorVersionThirteen`, `PostgresMajorVersionTwelve`
- New enum type `ReadReplicaPromoteOption` with values `ReadReplicaPromoteOptionForced`, `ReadReplicaPromoteOptionPlanned`
- New enum type `RecommendationType` with values `RecommendationTypeAnalyzeTable`, `RecommendationTypeCreateIndex`, `RecommendationTypeDropIndex`, `RecommendationTypeReIndex`
- New enum type `StartDataMigration` with values `StartDataMigrationFalse`, `StartDataMigrationTrue`
- New enum type `StorageAutoGrowthSupport` with values `StorageAutoGrowthSupportDisabled`, `StorageAutoGrowthSupportEnabled`
- New enum type `TriggerCutover` with values `TriggerCutoverFalse`, `TriggerCutoverTrue`
- New enum type `TuningOption` with values `TuningOptionIndex`, `TuningOptionTable`
- New enum type `ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupport` with values `ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupportDisabled`, `ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupportEnabled`
- New enum type `ZoneRedundantHighAvailabilitySupport` with values `ZoneRedundantHighAvailabilitySupportDisabled`, `ZoneRedundantHighAvailabilitySupportEnabled`
- New function `NewAdministratorsMicrosoftEntraClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AdministratorsMicrosoftEntraClient, error)`
- New function `*AdministratorsMicrosoftEntraClient.BeginCreateOrUpdate(context.Context, string, string, string, AdministratorMicrosoftEntraAdd, *AdministratorsMicrosoftEntraClientBeginCreateOrUpdateOptions) (*runtime.Poller[AdministratorsMicrosoftEntraClientCreateOrUpdateResponse], error)`
- New function `*AdministratorsMicrosoftEntraClient.BeginDelete(context.Context, string, string, string, *AdministratorsMicrosoftEntraClientBeginDeleteOptions) (*runtime.Poller[AdministratorsMicrosoftEntraClientDeleteResponse], error)`
- New function `*AdministratorsMicrosoftEntraClient.Get(context.Context, string, string, string, *AdministratorsMicrosoftEntraClientGetOptions) (AdministratorsMicrosoftEntraClientGetResponse, error)`
- New function `*AdministratorsMicrosoftEntraClient.NewListByServerPager(string, string, *AdministratorsMicrosoftEntraClientListByServerOptions) *runtime.Pager[AdministratorsMicrosoftEntraClientListByServerResponse]`
- New function `NewAdvancedThreatProtectionSettingsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AdvancedThreatProtectionSettingsClient, error)`
- New function `*AdvancedThreatProtectionSettingsClient.Get(context.Context, string, string, ThreatProtectionName, *AdvancedThreatProtectionSettingsClientGetOptions) (AdvancedThreatProtectionSettingsClientGetResponse, error)`
- New function `*AdvancedThreatProtectionSettingsClient.NewListByServerPager(string, string, *AdvancedThreatProtectionSettingsClientListByServerOptions) *runtime.Pager[AdvancedThreatProtectionSettingsClientListByServerResponse]`
- New function `NewBackupsAutomaticAndOnDemandClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BackupsAutomaticAndOnDemandClient, error)`
- New function `*BackupsAutomaticAndOnDemandClient.BeginCreate(context.Context, string, string, string, *BackupsAutomaticAndOnDemandClientBeginCreateOptions) (*runtime.Poller[BackupsAutomaticAndOnDemandClientCreateResponse], error)`
- New function `*BackupsAutomaticAndOnDemandClient.BeginDelete(context.Context, string, string, string, *BackupsAutomaticAndOnDemandClientBeginDeleteOptions) (*runtime.Poller[BackupsAutomaticAndOnDemandClientDeleteResponse], error)`
- New function `*BackupsAutomaticAndOnDemandClient.Get(context.Context, string, string, string, *BackupsAutomaticAndOnDemandClientGetOptions) (BackupsAutomaticAndOnDemandClientGetResponse, error)`
- New function `*BackupsAutomaticAndOnDemandClient.NewListByServerPager(string, string, *BackupsAutomaticAndOnDemandClientListByServerOptions) *runtime.Pager[BackupsAutomaticAndOnDemandClientListByServerResponse]`
- New function `NewBackupsLongTermRetentionClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BackupsLongTermRetentionClient, error)`
- New function `*BackupsLongTermRetentionClient.CheckPrerequisites(context.Context, string, string, LtrPreBackupRequest, *BackupsLongTermRetentionClientCheckPrerequisitesOptions) (BackupsLongTermRetentionClientCheckPrerequisitesResponse, error)`
- New function `*BackupsLongTermRetentionClient.Get(context.Context, string, string, string, *BackupsLongTermRetentionClientGetOptions) (BackupsLongTermRetentionClientGetResponse, error)`
- New function `*BackupsLongTermRetentionClient.NewListByServerPager(string, string, *BackupsLongTermRetentionClientListByServerOptions) *runtime.Pager[BackupsLongTermRetentionClientListByServerResponse]`
- New function `*BackupsLongTermRetentionClient.BeginStart(context.Context, string, string, BackupsLongTermRetentionRequest, *BackupsLongTermRetentionClientBeginStartOptions) (*runtime.Poller[BackupsLongTermRetentionClientStartResponse], error)`
- New function `NewCapabilitiesByLocationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CapabilitiesByLocationClient, error)`
- New function `*CapabilitiesByLocationClient.NewListPager(string, *CapabilitiesByLocationClientListOptions) *runtime.Pager[CapabilitiesByLocationClientListResponse]`
- New function `NewCapabilitiesByServerClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CapabilitiesByServerClient, error)`
- New function `*CapabilitiesByServerClient.NewListPager(string, string, *CapabilitiesByServerClientListOptions) *runtime.Pager[CapabilitiesByServerClientListResponse]`
- New function `NewCapturedLogsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CapturedLogsClient, error)`
- New function `*CapturedLogsClient.NewListByServerPager(string, string, *CapturedLogsClientListByServerOptions) *runtime.Pager[CapturedLogsClientListByServerResponse]`
- New function `*ClientFactory.NewAdministratorsMicrosoftEntraClient() *AdministratorsMicrosoftEntraClient`
- New function `*ClientFactory.NewAdvancedThreatProtectionSettingsClient() *AdvancedThreatProtectionSettingsClient`
- New function `*ClientFactory.NewBackupsAutomaticAndOnDemandClient() *BackupsAutomaticAndOnDemandClient`
- New function `*ClientFactory.NewBackupsLongTermRetentionClient() *BackupsLongTermRetentionClient`
- New function `*ClientFactory.NewCapabilitiesByLocationClient() *CapabilitiesByLocationClient`
- New function `*ClientFactory.NewCapabilitiesByServerClient() *CapabilitiesByServerClient`
- New function `*ClientFactory.NewCapturedLogsClient() *CapturedLogsClient`
- New function `*ClientFactory.NewNameAvailabilityClient() *NameAvailabilityClient`
- New function `*ClientFactory.NewPrivateDNSZoneSuffixClient() *PrivateDNSZoneSuffixClient`
- New function `*ClientFactory.NewQuotaUsagesClient() *QuotaUsagesClient`
- New function `*ClientFactory.NewTuningOptionsClient() *TuningOptionsClient`
- New function `*MigrationsClient.Cancel(context.Context, string, string, string, *MigrationsClientCancelOptions) (MigrationsClientCancelResponse, error)`
- New function `*MigrationsClient.CheckNameAvailability(context.Context, string, string, MigrationNameAvailability, *MigrationsClientCheckNameAvailabilityOptions) (MigrationsClientCheckNameAvailabilityResponse, error)`
- New function `PossibleGeographicallyRedundantBackupValues() []GeographicallyRedundantBackup`
- New function `NewPrivateDNSZoneSuffixClient(azcore.TokenCredential, *arm.ClientOptions) (*PrivateDNSZoneSuffixClient, error)`
- New function `*PrivateDNSZoneSuffixClient.Get(context.Context, *PrivateDNSZoneSuffixClientGetOptions) (PrivateDNSZoneSuffixClientGetResponse, error)`
- New function `*PrivateEndpointConnectionsClient.BeginDelete(context.Context, string, string, string, *PrivateEndpointConnectionsClientBeginDeleteOptions) (*runtime.Poller[PrivateEndpointConnectionsClientDeleteResponse], error)`
- New function `*PrivateEndpointConnectionsClient.BeginUpdate(context.Context, string, string, string, PrivateEndpointConnection, *PrivateEndpointConnectionsClientBeginUpdateOptions) (*runtime.Poller[PrivateEndpointConnectionsClientUpdateResponse], error)`
- New function `NewQuotaUsagesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*QuotaUsagesClient, error)`
- New function `*QuotaUsagesClient.NewListPager(string, *QuotaUsagesClientListOptions) *runtime.Pager[QuotaUsagesClientListResponse]`
- New function `*ServersClient.BeginCreateOrUpdate(context.Context, string, string, Server, *ServersClientBeginCreateOrUpdateOptions) (*runtime.Poller[ServersClientCreateOrUpdateResponse], error)`
- New function `*ServersClient.NewListBySubscriptionPager(*ServersClientListBySubscriptionOptions) *runtime.Pager[ServersClientListBySubscriptionResponse]`
- New function `NewTuningOptionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*TuningOptionsClient, error)`
- New function `*TuningOptionsClient.Get(context.Context, string, string, TuningOption, *TuningOptionsClientGetOptions) (TuningOptionsClientGetResponse, error)`
- New function `*TuningOptionsClient.NewListByServerPager(string, string, *TuningOptionsClientListByServerOptions) *runtime.Pager[TuningOptionsClientListByServerResponse]`
- New function `*TuningOptionsClient.NewListRecommendationsPager(string, string, TuningOption, *TuningOptionsClientListRecommendationsOptions) *runtime.Pager[TuningOptionsClientListRecommendationsResponse]`
- New function `*VirtualNetworkSubnetUsageClient.List(context.Context, string, VirtualNetworkSubnetUsageParameter, *VirtualNetworkSubnetUsageClientListOptions) (VirtualNetworkSubnetUsageClientListResponse, error)`
- New function `NewNameAvailabilityClient(string, azcore.TokenCredential, *arm.ClientOptions) (*NameAvailabilityClient, error)`
- New function `*NameAvailabilityClient.CheckGlobally(context.Context, CheckNameAvailabilityRequest, *NameAvailabilityClientCheckGloballyOptions) (NameAvailabilityClientCheckGloballyResponse, error)`
- New function `*NameAvailabilityClient.CheckWithLocation(context.Context, string, CheckNameAvailabilityRequest, *NameAvailabilityClientCheckWithLocationOptions) (NameAvailabilityClientCheckWithLocationResponse, error)`
- New struct `AdminCredentialsForPatch`
- New struct `AdministratorMicrosoftEntra`
- New struct `AdministratorMicrosoftEntraAdd`
- New struct `AdministratorMicrosoftEntraList`
- New struct `AdministratorMicrosoftEntraProperties`
- New struct `AdministratorMicrosoftEntraPropertiesForAdd`
- New struct `AdvancedThreatProtectionSettingsList`
- New struct `AdvancedThreatProtectionSettingsModel`
- New struct `AdvancedThreatProtectionSettingsProperties`
- New struct `AuthConfigForPatch`
- New struct `BackupAutomaticAndOnDemand`
- New struct `BackupAutomaticAndOnDemandList`
- New struct `BackupAutomaticAndOnDemandProperties`
- New struct `BackupForPatch`
- New struct `BackupsLongTermRetentionOperation`
- New struct `BackupsLongTermRetentionRequest`
- New struct `BackupsLongTermRetentionResponse`
- New struct `BackupsLongTermRetentionResponseProperties`
- New struct `Capability`
- New struct `CapabilityList`
- New struct `CapturedLog`
- New struct `CapturedLogList`
- New struct `CapturedLogProperties`
- New struct `Cluster`
- New struct `ConfigurationList`
- New struct `DatabaseList`
- New struct `DatabaseMigrationState`
- New struct `FirewallRuleList`
- New struct `HighAvailabilityForPatch`
- New struct `ImpactRecord`
- New struct `MaintenanceWindowForPatch`
- New struct `Migration`
- New struct `MigrationList`
- New struct `MigrationNameAvailability`
- New struct `MigrationProperties`
- New struct `MigrationPropertiesForPatch`
- New struct `MigrationSecretParametersForPatch`
- New struct `MigrationSubstateDetails`
- New struct `NameAvailabilityModel`
- New struct `NameProperty`
- New struct `ObjectRecommendation`
- New struct `ObjectRecommendationDetails`
- New struct `ObjectRecommendationList`
- New struct `ObjectRecommendationProperties`
- New struct `ObjectRecommendationPropertiesAnalyzedWorkload`
- New struct `ObjectRecommendationPropertiesImplementationDetails`
- New struct `OperationList`
- New struct `PrivateEndpointConnectionList`
- New struct `PrivateLinkResourceList`
- New struct `QuotaUsage`
- New struct `QuotaUsageList`
- New struct `SKUForPatch`
- New struct `ServerEditionCapability`
- New struct `ServerForPatch`
- New struct `ServerList`
- New struct `ServerPropertiesForPatch`
- New struct `SupportedFeature`
- New struct `TuningOptions`
- New struct `TuningOptionsList`
- New struct `VirtualEndpoint`
- New struct `VirtualEndpointsList`
- New struct `VirtualNetworkSubnetUsageModel`
- New anonymous field `ConfigurationList` in struct `ConfigurationsClientListByServerResponse`
- New anonymous field `DatabaseList` in struct `DatabasesClientListByServerResponse`
- New anonymous field `FirewallRuleList` in struct `FirewallRulesClientListByServerResponse`
- New anonymous field `Migration` in struct `MigrationsClientCreateResponse`
- New anonymous field `Migration` in struct `MigrationsClientGetResponse`
- New anonymous field `MigrationList` in struct `MigrationsClientListByTargetServerResponse`
- New anonymous field `Migration` in struct `MigrationsClientUpdateResponse`
- New anonymous field `PrivateEndpointConnectionList` in struct `PrivateEndpointConnectionsClientListByServerResponse`
- New anonymous field `PrivateLinkResourceList` in struct `PrivateLinkResourcesClientListByServerResponse`
- New anonymous field `ServerList` in struct `ReplicasClientListByServerResponse`
- New field `Cluster` in struct `ServerProperties`
- New field `SecurityProfile`, `SupportedFeatures` in struct `ServerSKUCapability`
- New field `SupportedFeatures` in struct `ServerVersionCapability`
- New anonymous field `ServerList` in struct `ServersClientListByResourceGroupResponse`
- New anonymous field `VirtualEndpoint` in struct `VirtualEndpointsClientGetResponse`
- New anonymous field `VirtualEndpointsList` in struct `VirtualEndpointsClientListByServerResponse`


## 5.0.0-beta.1 (2025-05-23)
### Breaking Changes

- Function `*ClientFactory.NewPostgreSQLServerManagementClient` has been removed
- Function `NewPostgreSQLServerManagementClient` has been removed
- Function `*PostgreSQLServerManagementClient.CheckMigrationNameAvailability` has been removed
- Operation `*OperationsClient.List` has supported pagination, use `*OperationsClient.NewListPager` instead.

### Features Added

- New value `ServerStateInaccessible`, `ServerStateProvisioning`, `ServerStateRestarting` added to enum type `ServerState`
- New value `ServerVersionSeventeen` added to enum type `ServerVersion`
- New value `SourceTypeApsaraDBRDS`, `SourceTypeCrunchyPostgreSQL`, `SourceTypeDigitalOceanDroplets`, `SourceTypeDigitalOceanPostgreSQL`, `SourceTypeEDBOracleServer`, `SourceTypeEDBPostgreSQL`, `SourceTypeHerokuPostgreSQL`, `SourceTypeHuaweiCompute`, `SourceTypeHuaweiRDS`, `SourceTypePostgreSQLCosmosDB`, `SourceTypePostgreSQLFlexibleServer`, `SourceTypeSupabasePostgreSQL` added to enum type `SourceType`
- New value `StorageTypeUltraSSDLRS` added to enum type `StorageType`
- New enum type `RecommendationType` with values `RecommendationTypeCreateIndex`, `RecommendationTypeDropIndex`
- New enum type `RecommendationTypeEnum` with values `RecommendationTypeEnumCreateIndex`, `RecommendationTypeEnumDropIndex`, `RecommendationTypeEnumReIndex`
- New enum type `SupportedFeatureStatusEnum` with values `SupportedFeatureStatusEnumDisabled`, `SupportedFeatureStatusEnumEnabled`
- New enum type `TuningOptionEnum` with values `TuningOptionEnumConfiguration`, `TuningOptionEnumIndex`
- New function `*ClientFactory.NewPostgreSQLManagementClient() *PostgreSQLManagementClient`
- New function `*ClientFactory.NewQuotaUsagesClient() *QuotaUsagesClient`
- New function `*ClientFactory.NewTuningConfigurationClient() *TuningConfigurationClient`
- New function `*ClientFactory.NewTuningIndexClient() *TuningIndexClient`
- New function `*ClientFactory.NewTuningOptionsClient() *TuningOptionsClient`
- New function `NewPostgreSQLManagementClient(azcore.TokenCredential, *arm.ClientOptions) (*PostgreSQLManagementClient, error)`
- New function `*PostgreSQLManagementClient.CheckMigrationNameAvailability(context.Context, string, string, string, MigrationNameAvailabilityResource, *PostgreSQLManagementClientCheckMigrationNameAvailabilityOptions) (PostgreSQLManagementClientCheckMigrationNameAvailabilityResponse, error)`
- New function `NewQuotaUsagesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*QuotaUsagesClient, error)`
- New function `*QuotaUsagesClient.NewListPager(string, *QuotaUsagesClientListOptions) *runtime.Pager[QuotaUsagesClientListResponse]`
- New function `NewTuningConfigurationClient(string, azcore.TokenCredential, *arm.ClientOptions) (*TuningConfigurationClient, error)`
- New function `*TuningConfigurationClient.BeginDisable(context.Context, string, string, TuningOptionEnum, *TuningConfigurationClientBeginDisableOptions) (*runtime.Poller[TuningConfigurationClientDisableResponse], error)`
- New function `*TuningConfigurationClient.BeginEnable(context.Context, string, string, TuningOptionEnum, *TuningConfigurationClientBeginEnableOptions) (*runtime.Poller[TuningConfigurationClientEnableResponse], error)`
- New function `*TuningConfigurationClient.NewListSessionDetailsPager(string, string, TuningOptionEnum, string, *TuningConfigurationClientListSessionDetailsOptions) *runtime.Pager[TuningConfigurationClientListSessionDetailsResponse]`
- New function `*TuningConfigurationClient.NewListSessionsPager(string, string, TuningOptionEnum, *TuningConfigurationClientListSessionsOptions) *runtime.Pager[TuningConfigurationClientListSessionsResponse]`
- New function `*TuningConfigurationClient.BeginStartSession(context.Context, string, string, TuningOptionEnum, ConfigTuningRequestParameter, *TuningConfigurationClientBeginStartSessionOptions) (*runtime.Poller[TuningConfigurationClientStartSessionResponse], error)`
- New function `*TuningConfigurationClient.BeginStopSession(context.Context, string, string, TuningOptionEnum, *TuningConfigurationClientBeginStopSessionOptions) (*runtime.Poller[TuningConfigurationClientStopSessionResponse], error)`
- New function `NewTuningIndexClient(string, azcore.TokenCredential, *arm.ClientOptions) (*TuningIndexClient, error)`
- New function `*TuningIndexClient.NewListRecommendationsPager(string, string, TuningOptionEnum, *TuningIndexClientListRecommendationsOptions) *runtime.Pager[TuningIndexClientListRecommendationsResponse]`
- New function `NewTuningOptionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*TuningOptionsClient, error)`
- New function `*TuningOptionsClient.Get(context.Context, string, string, TuningOptionEnum, *TuningOptionsClientGetOptions) (TuningOptionsClientGetResponse, error)`
- New function `*TuningOptionsClient.NewListByServerPager(string, string, *TuningOptionsClientListByServerOptions) *runtime.Pager[TuningOptionsClientListByServerResponse]`
- New struct `Cluster`
- New struct `ConfigTuningRequestParameter`
- New struct `ImpactRecord`
- New struct `IndexRecommendationDetails`
- New struct `IndexRecommendationListResult`
- New struct `IndexRecommendationResource`
- New struct `IndexRecommendationResourceProperties`
- New struct `IndexRecommendationResourcePropertiesAnalyzedWorkload`
- New struct `IndexRecommendationResourcePropertiesImplementationDetails`
- New struct `NameProperty`
- New struct `QuotaUsage`
- New struct `QuotaUsagesListResult`
- New struct `SessionDetailsListResult`
- New struct `SessionDetailsResource`
- New struct `SessionResource`
- New struct `SessionsListResult`
- New struct `SupportedFeature`
- New struct `TuningOptionsListResult`
- New struct `TuningOptionsResource`
- New field `SupportedFeatures` in struct `FlexibleServerCapability`
- New field `Cluster` in struct `ServerProperties`
- New field `Cluster` in struct `ServerPropertiesForUpdate`
- New field `SecurityProfile`, `SupportedFeatures` in struct `ServerSKUCapability`
- New field `SupportedFeatures` in struct `ServerVersionCapability`


## 4.1.0 (2025-03-17)
### Features Added

- New value `IdentityTypeSystemAssignedUserAssigned` added to enum type `IdentityType`
- New field `PrincipalID` in struct `UserAssignedIdentity`


## 4.0.0 (2025-01-04)
### Breaking Changes

- Type of `CapabilitiesListResult.Value` has been changed from `[]*CapabilityProperties` to `[]*FlexibleServerCapability`
- Type of `FastProvisioningEditionCapability.SupportedStorageGb` has been changed from `*int64` to `*int32`
- Type of `FlexibleServerEditionCapability.Status` has been changed from `*string` to `*CapabilityStatus`
- Type of `ServerVersionCapability.Status` has been changed from `*string` to `*CapabilityStatus`
- Type of `StorageEditionCapability.Status` has been changed from `*string` to `*CapabilityStatus`
- Type of `StorageTierCapability.Iops` has been changed from `*int64` to `*int32`
- Type of `StorageTierCapability.Status` has been changed from `*string` to `*CapabilityStatus`
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
- New value `OriginCustomerOnDemand` added to enum type `Origin`
- New value `ServerVersionFifteen`, `ServerVersionSixteen` added to enum type `ServerVersion`
- New enum type `AzureManagedDiskPerformanceTiers` with values `AzureManagedDiskPerformanceTiersP1`, `AzureManagedDiskPerformanceTiersP10`, `AzureManagedDiskPerformanceTiersP15`, `AzureManagedDiskPerformanceTiersP2`, `AzureManagedDiskPerformanceTiersP20`, `AzureManagedDiskPerformanceTiersP3`, `AzureManagedDiskPerformanceTiersP30`, `AzureManagedDiskPerformanceTiersP4`, `AzureManagedDiskPerformanceTiersP40`, `AzureManagedDiskPerformanceTiersP50`, `AzureManagedDiskPerformanceTiersP6`, `AzureManagedDiskPerformanceTiersP60`, `AzureManagedDiskPerformanceTiersP70`, `AzureManagedDiskPerformanceTiersP80`
- New enum type `CancelEnum` with values `CancelEnumFalse`, `CancelEnumTrue`
- New enum type `CapabilityStatus` with values `CapabilityStatusAvailable`, `CapabilityStatusDefault`, `CapabilityStatusDisabled`, `CapabilityStatusVisible`
- New enum type `ExecutionStatus` with values `ExecutionStatusCancelled`, `ExecutionStatusFailed`, `ExecutionStatusRunning`, `ExecutionStatusSucceeded`
- New enum type `FastProvisioningSupportedEnum` with values `FastProvisioningSupportedEnumDisabled`, `FastProvisioningSupportedEnumEnabled`
- New enum type `GeoBackupSupportedEnum` with values `GeoBackupSupportedEnumDisabled`, `GeoBackupSupportedEnumEnabled`
- New enum type `HaMode` with values `HaModeSameZone`, `HaModeZoneRedundant`
- New enum type `KeyStatusEnum` with values `KeyStatusEnumInvalid`, `KeyStatusEnumValid`
- New enum type `LogicalReplicationOnSourceDbEnum` with values `LogicalReplicationOnSourceDbEnumFalse`, `LogicalReplicationOnSourceDbEnumTrue`
- New enum type `MigrateRolesEnum` with values `MigrateRolesEnumFalse`, `MigrateRolesEnumTrue`
- New enum type `MigrationDbState` with values `MigrationDbStateCanceled`, `MigrationDbStateCanceling`, `MigrationDbStateFailed`, `MigrationDbStateInProgress`, `MigrationDbStateSucceeded`, `MigrationDbStateWaitingForCutoverTrigger`
- New enum type `MigrationListFilter` with values `MigrationListFilterActive`, `MigrationListFilterAll`
- New enum type `MigrationMode` with values `MigrationModeOffline`, `MigrationModeOnline`
- New enum type `MigrationNameAvailabilityReason` with values `MigrationNameAvailabilityReasonAlreadyExists`, `MigrationNameAvailabilityReasonInvalid`
- New enum type `MigrationOption` with values `MigrationOptionMigrate`, `MigrationOptionValidate`, `MigrationOptionValidateAndMigrate`
- New enum type `MigrationState` with values `MigrationStateCanceled`, `MigrationStateCleaningUp`, `MigrationStateFailed`, `MigrationStateInProgress`, `MigrationStateSucceeded`, `MigrationStateValidationFailed`, `MigrationStateWaitingForUserAction`
- New enum type `MigrationSubState` with values `MigrationSubStateCancelingRequestedDBMigrations`, `MigrationSubStateCompleted`, `MigrationSubStateCompletingMigration`, `MigrationSubStateMigratingData`, `MigrationSubStatePerformingPreRequisiteSteps`, `MigrationSubStateValidationInProgress`, `MigrationSubStateWaitingForCutoverTrigger`, `MigrationSubStateWaitingForDBsToMigrateSpecification`, `MigrationSubStateWaitingForDataMigrationScheduling`, `MigrationSubStateWaitingForDataMigrationWindow`, `MigrationSubStateWaitingForLogicalReplicationSetupRequestOnSourceDB`, `MigrationSubStateWaitingForTargetDBOverwriteConfirmation`
- New enum type `OnlineResizeSupportedEnum` with values `OnlineResizeSupportedEnumDisabled`, `OnlineResizeSupportedEnumEnabled`
- New enum type `OverwriteDbsInTargetEnum` with values `OverwriteDbsInTargetEnumFalse`, `OverwriteDbsInTargetEnumTrue`
- New enum type `PrivateEndpointConnectionProvisioningState` with values `PrivateEndpointConnectionProvisioningStateCreating`, `PrivateEndpointConnectionProvisioningStateDeleting`, `PrivateEndpointConnectionProvisioningStateFailed`, `PrivateEndpointConnectionProvisioningStateSucceeded`
- New enum type `PrivateEndpointServiceConnectionStatus` with values `PrivateEndpointServiceConnectionStatusApproved`, `PrivateEndpointServiceConnectionStatusPending`, `PrivateEndpointServiceConnectionStatusRejected`
- New enum type `ReadReplicaPromoteMode` with values `ReadReplicaPromoteModeStandalone`, `ReadReplicaPromoteModeSwitchover`
- New enum type `ReplicationPromoteOption` with values `ReplicationPromoteOptionForced`, `ReplicationPromoteOptionPlanned`
- New enum type `ReplicationState` with values `ReplicationStateActive`, `ReplicationStateBroken`, `ReplicationStateCatchup`, `ReplicationStateProvisioning`, `ReplicationStateReconfiguring`, `ReplicationStateUpdating`
- New enum type `RestrictedEnum` with values `RestrictedEnumDisabled`, `RestrictedEnumEnabled`
- New enum type `SSLMode` with values `SSLModePrefer`, `SSLModeRequire`, `SSLModeVerifyCA`, `SSLModeVerifyFull`
- New enum type `SourceType` with values `SourceTypeAWS`, `SourceTypeAWSAURORA`, `SourceTypeAWSEC2`, `SourceTypeAWSRDS`, `SourceTypeAzureVM`, `SourceTypeEDB`, `SourceTypeGCP`, `SourceTypeGCPAlloyDB`, `SourceTypeGCPCloudSQL`, `SourceTypeGCPCompute`, `SourceTypeOnPremises`, `SourceTypePostgreSQLSingleServer`
- New enum type `StartDataMigrationEnum` with values `StartDataMigrationEnumFalse`, `StartDataMigrationEnumTrue`
- New enum type `StorageAutoGrow` with values `StorageAutoGrowDisabled`, `StorageAutoGrowEnabled`
- New enum type `StorageAutoGrowthSupportedEnum` with values `StorageAutoGrowthSupportedEnumDisabled`, `StorageAutoGrowthSupportedEnumEnabled`
- New enum type `StorageType` with values `StorageTypePremiumLRS`, `StorageTypePremiumV2LRS`
- New enum type `ThreatProtectionName` with values `ThreatProtectionNameDefault`
- New enum type `ThreatProtectionState` with values `ThreatProtectionStateDisabled`, `ThreatProtectionStateEnabled`
- New enum type `TriggerCutoverEnum` with values `TriggerCutoverEnumFalse`, `TriggerCutoverEnumTrue`
- New enum type `ValidationState` with values `ValidationStateFailed`, `ValidationStateSucceeded`, `ValidationStateWarning`
- New enum type `VirtualEndpointType` with values `VirtualEndpointTypeReadWrite`
- New enum type `ZoneRedundantHaAndGeoBackupSupportedEnum` with values `ZoneRedundantHaAndGeoBackupSupportedEnumDisabled`, `ZoneRedundantHaAndGeoBackupSupportedEnumEnabled`
- New enum type `ZoneRedundantHaSupportedEnum` with values `ZoneRedundantHaSupportedEnumDisabled`, `ZoneRedundantHaSupportedEnumEnabled`
- New function `*BackupsClient.BeginCreate(context.Context, string, string, string, *BackupsClientBeginCreateOptions) (*runtime.Poller[BackupsClientCreateResponse], error)`
- New function `*BackupsClient.BeginDelete(context.Context, string, string, string, *BackupsClientBeginDeleteOptions) (*runtime.Poller[BackupsClientDeleteResponse], error)`
- New function `*ClientFactory.NewFlexibleServerClient() *FlexibleServerClient`
- New function `*ClientFactory.NewLogFilesClient() *LogFilesClient`
- New function `*ClientFactory.NewLtrBackupOperationsClient() *LtrBackupOperationsClient`
- New function `*ClientFactory.NewMigrationsClient() *MigrationsClient`
- New function `*ClientFactory.NewPostgreSQLServerManagementClient() *PostgreSQLServerManagementClient`
- New function `*ClientFactory.NewPrivateEndpointConnectionClient() *PrivateEndpointConnectionClient`
- New function `*ClientFactory.NewPrivateEndpointConnectionsClient() *PrivateEndpointConnectionsClient`
- New function `*ClientFactory.NewPrivateLinkResourcesClient() *PrivateLinkResourcesClient`
- New function `*ClientFactory.NewServerCapabilitiesClient() *ServerCapabilitiesClient`
- New function `*ClientFactory.NewServerThreatProtectionSettingsClient() *ServerThreatProtectionSettingsClient`
- New function `*ClientFactory.NewVirtualEndpointsClient() *VirtualEndpointsClient`
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
- New function `NewPostgreSQLServerManagementClient(azcore.TokenCredential, *arm.ClientOptions) (*PostgreSQLServerManagementClient, error)`
- New function `*PostgreSQLServerManagementClient.CheckMigrationNameAvailability(context.Context, string, string, string, MigrationNameAvailabilityResource, *PostgreSQLServerManagementClientCheckMigrationNameAvailabilityOptions) (PostgreSQLServerManagementClientCheckMigrationNameAvailabilityResponse, error)`
- New function `NewPrivateEndpointConnectionClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateEndpointConnectionClient, error)`
- New function `*PrivateEndpointConnectionClient.BeginDelete(context.Context, string, string, string, *PrivateEndpointConnectionClientBeginDeleteOptions) (*runtime.Poller[PrivateEndpointConnectionClientDeleteResponse], error)`
- New function `*PrivateEndpointConnectionClient.BeginUpdate(context.Context, string, string, string, PrivateEndpointConnection, *PrivateEndpointConnectionClientBeginUpdateOptions) (*runtime.Poller[PrivateEndpointConnectionClientUpdateResponse], error)`
- New function `NewPrivateEndpointConnectionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateEndpointConnectionsClient, error)`
- New function `*PrivateEndpointConnectionsClient.Get(context.Context, string, string, string, *PrivateEndpointConnectionsClientGetOptions) (PrivateEndpointConnectionsClientGetResponse, error)`
- New function `*PrivateEndpointConnectionsClient.NewListByServerPager(string, string, *PrivateEndpointConnectionsClientListByServerOptions) *runtime.Pager[PrivateEndpointConnectionsClientListByServerResponse]`
- New function `NewPrivateLinkResourcesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateLinkResourcesClient, error)`
- New function `*PrivateLinkResourcesClient.Get(context.Context, string, string, string, *PrivateLinkResourcesClientGetOptions) (PrivateLinkResourcesClientGetResponse, error)`
- New function `*PrivateLinkResourcesClient.NewListByServerPager(string, string, *PrivateLinkResourcesClientListByServerOptions) *runtime.Pager[PrivateLinkResourcesClientListByServerResponse]`
- New function `NewServerCapabilitiesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ServerCapabilitiesClient, error)`
- New function `*ServerCapabilitiesClient.NewListPager(string, string, *ServerCapabilitiesClientListOptions) *runtime.Pager[ServerCapabilitiesClientListResponse]`
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
- New struct `AdminCredentials`
- New struct `BackupSettings`
- New struct `BackupStoreDetails`
- New struct `DbLevelValidationStatus`
- New struct `DbMigrationStatus`
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
- New struct `PrivateEndpoint`
- New struct `PrivateEndpointConnection`
- New struct `PrivateEndpointConnectionListResult`
- New struct `PrivateEndpointConnectionProperties`
- New struct `PrivateLinkResource`
- New struct `PrivateLinkResourceListResult`
- New struct `PrivateLinkResourceProperties`
- New struct `PrivateLinkServiceConnectionState`
- New struct `Replica`
- New struct `ServerSKU`
- New struct `ServerSKUCapability`
- New struct `ServerThreatProtectionListResult`
- New struct `ServerThreatProtectionProperties`
- New struct `ServerThreatProtectionSettingsModel`
- New struct `StorageMbCapability`
- New struct `ValidationDetails`
- New struct `ValidationMessage`
- New struct `ValidationSummaryItem`
- New struct `VirtualEndpointResource`
- New struct `VirtualEndpointResourceForPatch`
- New struct `VirtualEndpointResourceProperties`
- New struct `VirtualEndpointsListResult`
- New field `GeoBackupEncryptionKeyStatus`, `GeoBackupKeyURI`, `GeoBackupUserAssignedIdentityID`, `PrimaryEncryptionKeyStatus` in struct `DataEncryption`
- New field `Reason`, `ServerCount`, `Status`, `SupportedTier` in struct `FastProvisioningEditionCapability`
- New field `DefaultSKUName`, `Reason`, `SupportedServerSKUs` in struct `FlexibleServerEditionCapability`
- New field `PrivateEndpointConnections`, `Replica` in struct `ServerProperties`
- New field `AdministratorLogin`, `Network`, `Replica` in struct `ServerPropertiesForUpdate`
- New field `Reason` in struct `ServerVersionCapability`
- New field `AutoGrow`, `Iops`, `Throughput`, `Tier`, `Type` in struct `Storage`
- New field `DefaultStorageSizeMb`, `Reason`, `SupportedStorageMb` in struct `StorageEditionCapability`
- New field `Reason` in struct `StorageTierCapability`
- New field `TenantID` in struct `UserAssignedIdentity`


## 4.0.0-beta.5 (2024-04-26)
### Features Added

- New value `SourceTypeAWSAURORA`, `SourceTypeAWSEC2`, `SourceTypeAWSRDS`, `SourceTypeEDB`, `SourceTypeGCPAlloyDB`, `SourceTypeGCPCloudSQL`, `SourceTypeGCPCompute` added to enum type `SourceType`
- New enum type `MigrateRolesEnum` with values `MigrateRolesEnumFalse`, `MigrateRolesEnumTrue`
- New field `MigrateRoles`, `MigrationInstanceResourceID` in struct `MigrationResourceProperties`
- New field `MigrateRoles` in struct `MigrationResourcePropertiesForPatch`


## 4.0.0-beta.4 (2023-12-22)
### Other Changes

- Operation `ServerThreatProtectionSettingsClient.BeginCreateOrUpdate` increase `202` response.


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
