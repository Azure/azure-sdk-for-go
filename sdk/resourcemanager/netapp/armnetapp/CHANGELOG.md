# Release History

## 7.4.0 (2025-02-28)
### Features Added

- New enum type `CoolAccessTieringPolicy` with values `CoolAccessTieringPolicyAuto`, `CoolAccessTieringPolicySnapshotOnly`
- New function `*AccountsClient.BeginChangeKeyVault(context.Context, string, string, *AccountsClientBeginChangeKeyVaultOptions) (*runtime.Poller[AccountsClientChangeKeyVaultResponse], error)`
- New function `*AccountsClient.BeginGetChangeKeyVaultInformation(context.Context, string, string, *AccountsClientBeginGetChangeKeyVaultInformationOptions) (*runtime.Poller[AccountsClientGetChangeKeyVaultInformationResponse], error)`
- New function `*AccountsClient.BeginTransitionToCmk(context.Context, string, string, *AccountsClientBeginTransitionToCmkOptions) (*runtime.Poller[AccountsClientTransitionToCmkResponse], error)`
- New struct `ChangeKeyVault`
- New struct `EncryptionTransitionRequest`
- New struct `KeyVaultPrivateEndpoint`
- New field `CoolAccessTieringPolicy` in struct `VolumePatchProperties`
- New field `CoolAccessTieringPolicy` in struct `VolumeProperties`


## 7.4.0-beta.1 (2024-11-18)
### Features Added

- New value `ServiceLevelFlexible` added to enum type `ServiceLevel`
- New enum type `AcceptGrowCapacityPoolForShortTermCloneSplit` with values `AcceptGrowCapacityPoolForShortTermCloneSplitAccepted`, `AcceptGrowCapacityPoolForShortTermCloneSplitDeclined`
- New enum type `ReplicationType` with values `ReplicationTypeCrossRegionReplication`, `ReplicationTypeCrossZoneReplication`
- New enum type `VolumeLanguage` with values `VolumeLanguageAr`, `VolumeLanguageArUTF8`, `VolumeLanguageC`, `VolumeLanguageCUTF8`, `VolumeLanguageCs`, `VolumeLanguageCsUTF8`, `VolumeLanguageDa`, `VolumeLanguageDaUTF8`, `VolumeLanguageDe`, `VolumeLanguageDeUTF8`, `VolumeLanguageEn`, `VolumeLanguageEnUTF8`, `VolumeLanguageEnUs`, `VolumeLanguageEnUsUTF8`, `VolumeLanguageEs`, `VolumeLanguageEsUTF8`, `VolumeLanguageFi`, `VolumeLanguageFiUTF8`, `VolumeLanguageFr`, `VolumeLanguageFrUTF8`, `VolumeLanguageHe`, `VolumeLanguageHeUTF8`, `VolumeLanguageHr`, `VolumeLanguageHrUTF8`, `VolumeLanguageHu`, `VolumeLanguageHuUTF8`, `VolumeLanguageIt`, `VolumeLanguageItUTF8`, `VolumeLanguageJa`, `VolumeLanguageJaJp932`, `VolumeLanguageJaJp932UTF8`, `VolumeLanguageJaJpPck`, `VolumeLanguageJaJpPckUTF8`, `VolumeLanguageJaJpPckV2`, `VolumeLanguageJaJpPckV2UTF8`, `VolumeLanguageJaUTF8`, `VolumeLanguageJaV1`, `VolumeLanguageJaV1UTF8`, `VolumeLanguageKo`, `VolumeLanguageKoUTF8`, `VolumeLanguageNl`, `VolumeLanguageNlUTF8`, `VolumeLanguageNo`, `VolumeLanguageNoUTF8`, `VolumeLanguagePl`, `VolumeLanguagePlUTF8`, `VolumeLanguagePt`, `VolumeLanguagePtUTF8`, `VolumeLanguageRo`, `VolumeLanguageRoUTF8`, `VolumeLanguageRu`, `VolumeLanguageRuUTF8`, `VolumeLanguageSk`, `VolumeLanguageSkUTF8`, `VolumeLanguageSl`, `VolumeLanguageSlUTF8`, `VolumeLanguageSv`, `VolumeLanguageSvUTF8`, `VolumeLanguageTr`, `VolumeLanguageTrUTF8`, `VolumeLanguageUTF8Mb4`, `VolumeLanguageZh`, `VolumeLanguageZhGbk`, `VolumeLanguageZhGbkUTF8`, `VolumeLanguageZhTw`, `VolumeLanguageZhTwBig5`, `VolumeLanguageZhTwBig5UTF8`, `VolumeLanguageZhTwUTF8`, `VolumeLanguageZhUTF8`
- New function `*AccountsClient.BeginChangeKeyVault(context.Context, string, string, *AccountsClientBeginChangeKeyVaultOptions) (*runtime.Poller[AccountsClientChangeKeyVaultResponse], error)`
- New function `*AccountsClient.BeginGetChangeKeyVaultInformation(context.Context, string, string, *AccountsClientBeginGetChangeKeyVaultInformationOptions) (*runtime.Poller[AccountsClientGetChangeKeyVaultInformationResponse], error)`
- New function `*AccountsClient.BeginTransitionToCmk(context.Context, string, string, *AccountsClientBeginTransitionToCmkOptions) (*runtime.Poller[AccountsClientTransitionToCmkResponse], error)`
- New function `*VolumesClient.BeginListQuotaReport(context.Context, string, string, string, string, *VolumesClientBeginListQuotaReportOptions) (*runtime.Poller[VolumesClientListQuotaReportResponse], error)`
- New function `*VolumesClient.BeginSplitCloneFromParent(context.Context, string, string, string, string, *VolumesClientBeginSplitCloneFromParentOptions) (*runtime.Poller[VolumesClientSplitCloneFromParentResponse], error)`
- New struct `ChangeKeyVault`
- New struct `DestinationReplication`
- New struct `EncryptionTransitionRequest`
- New struct `KeyVaultPrivateEndpoint`
- New struct `ListQuotaReportResponse`
- New struct `QuotaReport`
- New field `IsMultiAdEnabled`, `NfsV4IDDomain` in struct `AccountProperties`
- New field `IsLargeVolume` in struct `BackupProperties`
- New field `FederatedClientID` in struct `EncryptionIdentity`
- New field `CustomThroughputMibps` in struct `PoolPatchProperties`
- New field `CustomThroughputMibps` in struct `PoolProperties`
- New field `DestinationReplications` in struct `ReplicationObject`
- New field `AcceptGrowCapacityPoolForShortTermCloneSplit`, `InheritedSizeInBytes`, `Language` in struct `VolumeProperties`


## 7.3.0 (2024-10-23)
### Features Added

- New function `*VolumesClient.BeginAuthorizeExternalReplication(context.Context, string, string, string, string, *VolumesClientBeginAuthorizeExternalReplicationOptions) (*runtime.Poller[VolumesClientAuthorizeExternalReplicationResponse], error)`
- New function `*VolumesClient.BeginFinalizeExternalReplication(context.Context, string, string, string, string, *VolumesClientBeginFinalizeExternalReplicationOptions) (*runtime.Poller[VolumesClientFinalizeExternalReplicationResponse], error)`
- New function `*VolumesClient.BeginPeerExternalCluster(context.Context, string, string, string, string, PeerClusterForVolumeMigrationRequest, *VolumesClientBeginPeerExternalClusterOptions) (*runtime.Poller[VolumesClientPeerExternalClusterResponse], error)`
- New function `*VolumesClient.BeginPerformReplicationTransfer(context.Context, string, string, string, string, *VolumesClientBeginPerformReplicationTransferOptions) (*runtime.Poller[VolumesClientPerformReplicationTransferResponse], error)`
- New struct `ClusterPeerCommandResponse`
- New struct `PeerClusterForVolumeMigrationRequest`
- New struct `RemotePath`
- New struct `SvmPeerCommandResponse`
- New field `AvailabilityZone` in struct `FilePathAvailabilityRequest`
- New field `RemotePath` in struct `ReplicationObject`
- New field `EffectiveNetworkFeatures` in struct `VolumeProperties`


## 7.2.0 (2024-08-23)
### Features Added

- New field `ReplicationID` in struct `Replication`


## 7.1.0 (2024-07-26)
### Features Added

- New field `ProtocolTypes` in struct `VolumePatchProperties`


## 7.0.0 (2024-05-24)
### Breaking Changes

- Function `*BackupsClient.GetVolumeRestoreStatus` has been removed

### Features Added

- New enum type `BackupType` with values `BackupTypeManual`, `BackupTypeScheduled`
- New function `NewBackupVaultsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BackupVaultsClient, error)`
- New function `*BackupVaultsClient.BeginCreateOrUpdate(context.Context, string, string, string, BackupVault, *BackupVaultsClientBeginCreateOrUpdateOptions) (*runtime.Poller[BackupVaultsClientCreateOrUpdateResponse], error)`
- New function `*BackupVaultsClient.BeginDelete(context.Context, string, string, string, *BackupVaultsClientBeginDeleteOptions) (*runtime.Poller[BackupVaultsClientDeleteResponse], error)`
- New function `*BackupVaultsClient.Get(context.Context, string, string, string, *BackupVaultsClientGetOptions) (BackupVaultsClientGetResponse, error)`
- New function `*BackupVaultsClient.NewListByNetAppAccountPager(string, string, *BackupVaultsClientListByNetAppAccountOptions) *runtime.Pager[BackupVaultsClientListByNetAppAccountResponse]`
- New function `*BackupVaultsClient.BeginUpdate(context.Context, string, string, string, BackupVaultPatch, *BackupVaultsClientBeginUpdateOptions) (*runtime.Poller[BackupVaultsClientUpdateResponse], error)`
- New function `*BackupsClient.BeginCreate(context.Context, string, string, string, string, Backup, *BackupsClientBeginCreateOptions) (*runtime.Poller[BackupsClientCreateResponse], error)`
- New function `*BackupsClient.BeginDelete(context.Context, string, string, string, string, *BackupsClientBeginDeleteOptions) (*runtime.Poller[BackupsClientDeleteResponse], error)`
- New function `*BackupsClient.Get(context.Context, string, string, string, string, *BackupsClientGetOptions) (BackupsClientGetResponse, error)`
- New function `*BackupsClient.GetLatestStatus(context.Context, string, string, string, string, *BackupsClientGetLatestStatusOptions) (BackupsClientGetLatestStatusResponse, error)`
- New function `*BackupsClient.GetVolumeLatestRestoreStatus(context.Context, string, string, string, string, *BackupsClientGetVolumeLatestRestoreStatusOptions) (BackupsClientGetVolumeLatestRestoreStatusResponse, error)`
- New function `*BackupsClient.NewListByVaultPager(string, string, string, *BackupsClientListByVaultOptions) *runtime.Pager[BackupsClientListByVaultResponse]`
- New function `*BackupsClient.BeginUpdate(context.Context, string, string, string, string, BackupPatch, *BackupsClientBeginUpdateOptions) (*runtime.Poller[BackupsClientUpdateResponse], error)`
- New function `NewBackupsUnderAccountClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BackupsUnderAccountClient, error)`
- New function `*BackupsUnderAccountClient.BeginMigrateBackups(context.Context, string, string, BackupsMigrationRequest, *BackupsUnderAccountClientBeginMigrateBackupsOptions) (*runtime.Poller[BackupsUnderAccountClientMigrateBackupsResponse], error)`
- New function `NewBackupsUnderBackupVaultClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BackupsUnderBackupVaultClient, error)`
- New function `*BackupsUnderBackupVaultClient.BeginRestoreFiles(context.Context, string, string, string, string, BackupRestoreFiles, *BackupsUnderBackupVaultClientBeginRestoreFilesOptions) (*runtime.Poller[BackupsUnderBackupVaultClientRestoreFilesResponse], error)`
- New function `NewBackupsUnderVolumeClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BackupsUnderVolumeClient, error)`
- New function `*BackupsUnderVolumeClient.BeginMigrateBackups(context.Context, string, string, string, string, BackupsMigrationRequest, *BackupsUnderVolumeClientBeginMigrateBackupsOptions) (*runtime.Poller[BackupsUnderVolumeClientMigrateBackupsResponse], error)`
- New function `*ClientFactory.NewBackupVaultsClient() *BackupVaultsClient`
- New function `*ClientFactory.NewBackupsUnderAccountClient() *BackupsUnderAccountClient`
- New function `*ClientFactory.NewBackupsUnderBackupVaultClient() *BackupsUnderBackupVaultClient`
- New function `*ClientFactory.NewBackupsUnderVolumeClient() *BackupsUnderVolumeClient`
- New function `*ClientFactory.NewResourceRegionInfosClient() *ResourceRegionInfosClient`
- New function `NewResourceRegionInfosClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ResourceRegionInfosClient, error)`
- New function `*ResourceRegionInfosClient.Get(context.Context, string, *ResourceRegionInfosClientGetOptions) (ResourceRegionInfosClientGetResponse, error)`
- New function `*ResourceRegionInfosClient.NewListPager(string, *ResourceRegionInfosClientListOptions) *runtime.Pager[ResourceRegionInfosClientListResponse]`
- New struct `Backup`
- New struct `BackupPatch`
- New struct `BackupPatchProperties`
- New struct `BackupProperties`
- New struct `BackupRestoreFiles`
- New struct `BackupStatus`
- New struct `BackupVault`
- New struct `BackupVaultPatch`
- New struct `BackupVaultProperties`
- New struct `BackupVaultsList`
- New struct `BackupsList`
- New struct `BackupsMigrationRequest`
- New struct `RegionInfoResource`
- New struct `RegionInfosList`
- New struct `VolumeBackupProperties`
- New field `VolumeResourceID` in struct `VolumeBackups`
- New field `Backup` in struct `VolumePatchPropertiesDataProtection`
- New field `Backup` in struct `VolumePropertiesDataProtection`


## 6.0.0 (2024-03-22)
### Breaking Changes

- Field `DeploymentSpecID` of struct `VolumeGroupMetaData` has been removed

### Features Added

- New value `RelationshipStatusFailed`, `RelationshipStatusUnknown` added to enum type `RelationshipStatus`


## 6.0.0-beta.1 (2023-12-22)
### Breaking Changes

- Field `DeploymentSpecID` of struct `VolumeGroupMetaData` has been removed

### Features Added

- New enum type `BackupType` with values `BackupTypeManual`, `BackupTypeScheduled`
- New function `NewAccountBackupsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AccountBackupsClient, error)`
- New function `*AccountBackupsClient.BeginDelete(context.Context, string, string, string, *AccountBackupsClientBeginDeleteOptions) (*runtime.Poller[AccountBackupsClientDeleteResponse], error)`
- New function `*AccountBackupsClient.Get(context.Context, string, string, string, *AccountBackupsClientGetOptions) (AccountBackupsClientGetResponse, error)`
- New function `*AccountBackupsClient.NewListByNetAppAccountPager(string, string, *AccountBackupsClientListByNetAppAccountOptions) *runtime.Pager[AccountBackupsClientListByNetAppAccountResponse]`
- New function `*AccountsClient.BeginMigrateEncryptionKey(context.Context, string, string, *AccountsClientBeginMigrateEncryptionKeyOptions) (*runtime.Poller[AccountsClientMigrateEncryptionKeyResponse], error)`
- New function `NewBackupVaultsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BackupVaultsClient, error)`
- New function `*BackupVaultsClient.BeginCreateOrUpdate(context.Context, string, string, string, BackupVault, *BackupVaultsClientBeginCreateOrUpdateOptions) (*runtime.Poller[BackupVaultsClientCreateOrUpdateResponse], error)`
- New function `*BackupVaultsClient.BeginDelete(context.Context, string, string, string, *BackupVaultsClientBeginDeleteOptions) (*runtime.Poller[BackupVaultsClientDeleteResponse], error)`
- New function `*BackupVaultsClient.Get(context.Context, string, string, string, *BackupVaultsClientGetOptions) (BackupVaultsClientGetResponse, error)`
- New function `*BackupVaultsClient.NewListByNetAppAccountPager(string, string, *BackupVaultsClientListByNetAppAccountOptions) *runtime.Pager[BackupVaultsClientListByNetAppAccountResponse]`
- New function `*BackupVaultsClient.BeginUpdate(context.Context, string, string, string, BackupVaultPatch, *BackupVaultsClientBeginUpdateOptions) (*runtime.Poller[BackupVaultsClientUpdateResponse], error)`
- New function `*BackupsClient.BeginCreate(context.Context, string, string, string, string, Backup, *BackupsClientBeginCreateOptions) (*runtime.Poller[BackupsClientCreateResponse], error)`
- New function `*BackupsClient.BeginDelete(context.Context, string, string, string, string, *BackupsClientBeginDeleteOptions) (*runtime.Poller[BackupsClientDeleteResponse], error)`
- New function `*BackupsClient.Get(context.Context, string, string, string, string, *BackupsClientGetOptions) (BackupsClientGetResponse, error)`
- New function `*BackupsClient.GetLatestStatus(context.Context, string, string, string, string, *BackupsClientGetLatestStatusOptions) (BackupsClientGetLatestStatusResponse, error)`
- New function `*BackupsClient.NewListByVaultPager(string, string, string, *BackupsClientListByVaultOptions) *runtime.Pager[BackupsClientListByVaultResponse]`
- New function `*BackupsClient.BeginUpdate(context.Context, string, string, string, string, BackupPatch, *BackupsClientBeginUpdateOptions) (*runtime.Poller[BackupsClientUpdateResponse], error)`
- New function `NewBackupsUnderAccountClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BackupsUnderAccountClient, error)`
- New function `*BackupsUnderAccountClient.BeginMigrateBackups(context.Context, string, string, BackupsMigrationRequest, *BackupsUnderAccountClientBeginMigrateBackupsOptions) (*runtime.Poller[BackupsUnderAccountClientMigrateBackupsResponse], error)`
- New function `NewBackupsUnderBackupVaultClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BackupsUnderBackupVaultClient, error)`
- New function `*BackupsUnderBackupVaultClient.BeginRestoreFiles(context.Context, string, string, string, string, BackupRestoreFiles, *BackupsUnderBackupVaultClientBeginRestoreFilesOptions) (*runtime.Poller[BackupsUnderBackupVaultClientRestoreFilesResponse], error)`
- New function `NewBackupsUnderVolumeClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BackupsUnderVolumeClient, error)`
- New function `*BackupsUnderVolumeClient.BeginMigrateBackups(context.Context, string, string, string, string, BackupsMigrationRequest, *BackupsUnderVolumeClientBeginMigrateBackupsOptions) (*runtime.Poller[BackupsUnderVolumeClientMigrateBackupsResponse], error)`
- New function `*ClientFactory.NewAccountBackupsClient() *AccountBackupsClient`
- New function `*ClientFactory.NewBackupVaultsClient() *BackupVaultsClient`
- New function `*ClientFactory.NewBackupsUnderAccountClient() *BackupsUnderAccountClient`
- New function `*ClientFactory.NewBackupsUnderBackupVaultClient() *BackupsUnderBackupVaultClient`
- New function `*ClientFactory.NewBackupsUnderVolumeClient() *BackupsUnderVolumeClient`
- New function `*ClientFactory.NewResourceRegionInfosClient() *ResourceRegionInfosClient`
- New function `NewResourceRegionInfosClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ResourceRegionInfosClient, error)`
- New function `*ResourceRegionInfosClient.Get(context.Context, string, *ResourceRegionInfosClientGetOptions) (ResourceRegionInfosClientGetResponse, error)`
- New function `*ResourceRegionInfosClient.NewListPager(string, *ResourceRegionInfosClientListOptions) *runtime.Pager[ResourceRegionInfosClientListResponse]`
- New function `*VolumesClient.BeginSplitCloneFromParent(context.Context, string, string, string, string, *VolumesClientBeginSplitCloneFromParentOptions) (*runtime.Poller[VolumesClientSplitCloneFromParentResponse], error)`
- New struct `Backup`
- New struct `BackupPatch`
- New struct `BackupPatchProperties`
- New struct `BackupProperties`
- New struct `BackupRestoreFiles`
- New struct `BackupStatus`
- New struct `BackupVault`
- New struct `BackupVaultPatch`
- New struct `BackupVaultProperties`
- New struct `BackupVaultsList`
- New struct `BackupsList`
- New struct `BackupsMigrationRequest`
- New struct `EncryptionMigrationRequest`
- New struct `RegionInfoResource`
- New struct `RegionInfosList`
- New struct `RemotePath`
- New struct `VolumeBackupProperties`
- New field `IsMultiAdEnabled`, `NfsV4IDDomain` in struct `AccountProperties`
- New field `RemotePath` in struct `ReplicationObject`
- New field `Backup` in struct `VolumePatchPropertiesDataProtection`
- New field `InheritedSizeInBytes` in struct `VolumeProperties`
- New field `Backup` in struct `VolumePropertiesDataProtection`


## 5.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 5.0.0 (2023-10-27)
### Breaking Changes

- Enum `BackupType` has been removed
- Function `NewAccountBackupsClient` has been removed
- Function `*AccountBackupsClient.BeginDelete` has been removed
- Function `*AccountBackupsClient.Get` has been removed
- Function `*AccountBackupsClient.NewListPager` has been removed
- Function `*BackupsClient.BeginCreate` has been removed
- Function `*BackupsClient.BeginDelete` has been removed
- Function `*BackupsClient.Get` has been removed
- Function `*BackupsClient.GetStatus` has been removed
- Function `*BackupsClient.NewListPager` has been removed
- Function `*BackupsClient.BeginRestoreFiles` has been removed
- Function `*BackupsClient.BeginUpdate` has been removed
- Function `*ClientFactory.NewAccountBackupsClient` has been removed
- Struct `Backup` has been removed
- Struct `BackupPatch` has been removed
- Struct `BackupProperties` has been removed
- Struct `BackupRestoreFiles` has been removed
- Struct `BackupStatus` has been removed
- Struct `BackupsList` has been removed
- Struct `VolumeBackupProperties` has been removed
- Field `Backup` of struct `VolumePatchPropertiesDataProtection` has been removed
- Field `Backup` of struct `VolumePropertiesDataProtection` has been removed

### Features Added

- New value `ApplicationTypeORACLE` added to enum type `ApplicationType`
- New value `NetworkFeaturesBasicStandard`, `NetworkFeaturesStandardBasic` added to enum type `NetworkFeatures`
- New enum type `CoolAccessRetrievalPolicy` with values `CoolAccessRetrievalPolicyDefault`, `CoolAccessRetrievalPolicyNever`, `CoolAccessRetrievalPolicyOnRead`
- New enum type `NetworkSiblingSetProvisioningState` with values `NetworkSiblingSetProvisioningStateCanceled`, `NetworkSiblingSetProvisioningStateFailed`, `NetworkSiblingSetProvisioningStateSucceeded`, `NetworkSiblingSetProvisioningStateUpdating`
- New function `*ResourceClient.QueryNetworkSiblingSet(context.Context, string, QueryNetworkSiblingSetRequest, *ResourceClientQueryNetworkSiblingSetOptions) (ResourceClientQueryNetworkSiblingSetResponse, error)`
- New function `*ResourceClient.BeginUpdateNetworkSiblingSet(context.Context, string, UpdateNetworkSiblingSetRequest, *ResourceClientBeginUpdateNetworkSiblingSetOptions) (*runtime.Poller[ResourceClientUpdateNetworkSiblingSetResponse], error)`
- New function `*VolumesClient.BeginPopulateAvailabilityZone(context.Context, string, string, string, string, *VolumesClientBeginPopulateAvailabilityZoneOptions) (*runtime.Poller[VolumesClientPopulateAvailabilityZoneResponse], error)`
- New struct `NetworkSiblingSet`
- New struct `NicInfo`
- New struct `QueryNetworkSiblingSetRequest`
- New struct `UpdateNetworkSiblingSetRequest`
- New field `Zones` in struct `VolumeGroupVolumeProperties`
- New field `CoolAccessRetrievalPolicy`, `SmbAccessBasedEnumeration`, `SmbNonBrowsable` in struct `VolumePatchProperties`
- New field `CoolAccessRetrievalPolicy` in struct `VolumeProperties`


## 4.1.0 (2023-07-28)
### Features Added

- New value `RegionStorageToNetworkProximityAcrossT2`, `RegionStorageToNetworkProximityT1AndAcrossT2`, `RegionStorageToNetworkProximityT1AndT2AndAcrossT2`, `RegionStorageToNetworkProximityT2AndAcrossT2` added to enum type `RegionStorageToNetworkProximity`
- New value `VolumeStorageToNetworkProximityAcrossT2` added to enum type `VolumeStorageToNetworkProximity`
- New function `*VolumesClient.BeginListGetGroupIDListForLdapUser(context.Context, string, string, string, string, GetGroupIDListForLDAPUserRequest, *VolumesClientBeginListGetGroupIDListForLdapUserOptions) (*runtime.Poller[VolumesClientListGetGroupIDListForLdapUserResponse], error)`
- New struct `GetGroupIDListForLDAPUserRequest`
- New struct `GetGroupIDListForLDAPUserResponse`
- New field `Identity` in struct `AccountPatch`
- New field `SnapshotDirectoryVisible` in struct `VolumePatchProperties`
- New field `ActualThroughputMibps`, `OriginatingResourceID` in struct `VolumeProperties`


## 4.0.0 (2023-03-24)
### Breaking Changes

- Type of `Account.Identity` has been changed from `*Identity` to `*ManagedServiceIdentity`
- Type alias `IdentityType` has been removed
- Function `NewVaultsClient` has been removed
- Function `*VaultsClient.NewListPager` has been removed
- Struct `Identity` has been removed
- Struct `Vault` has been removed
- Struct `VaultList` has been removed
- Struct `VaultProperties` has been removed
- Struct `VaultsClient` has been removed
- Field `VaultID` of struct `VolumeBackupProperties` has been removed

### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New enum type `FileAccessLogs` with values `FileAccessLogsDisabled`, `FileAccessLogsEnabled`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New function `*BackupsClient.BeginRestoreFiles(context.Context, string, string, string, string, string, BackupRestoreFiles, *BackupsClientBeginRestoreFilesOptions) (*runtime.Poller[BackupsClientRestoreFilesResponse], error)`
- New function `*VolumesClient.BeginBreakFileLocks(context.Context, string, string, string, string, *VolumesClientBeginBreakFileLocksOptions) (*runtime.Poller[VolumesClientBreakFileLocksResponse], error)`
- New struct `BackupRestoreFiles`
- New struct `BreakFileLocksRequest`
- New struct `ManagedServiceIdentity`
- New struct `VolumeRelocationProperties`
- New field `PreferredServersForLdapClient` in struct `ActiveDirectory`
- New field `SystemData` in struct `Backup`
- New field `SystemData` in struct `Snapshot`
- New field `DataStoreResourceID` in struct `VolumeProperties`
- New field `FileAccessLogs` in struct `VolumeProperties`
- New field `IsLargeVolume` in struct `VolumeProperties`
- New field `ProvisionedAvailabilityZone` in struct `VolumeProperties`
- New field `VolumeRelocation` in struct `VolumePropertiesDataProtection`
- New field `Tags` in struct `VolumeQuotaRulePatch`


## 3.0.0 (2022-09-16)
### Breaking Changes

- Type of `AccountEncryption.KeySource` has been changed from `*string` to `*KeySource`
- Field `Location` of struct `Vault` has been removed

### Features Added

- New const `KeyVaultStatusCreated`
- New const `IdentityTypeNone`
- New const `IdentityTypeUserAssigned`
- New const `IdentityTypeSystemAssigned`
- New const `RegionStorageToNetworkProximityDefault`
- New const `KeyVaultStatusUpdating`
- New const `KeyVaultStatusError`
- New const `RegionStorageToNetworkProximityT1AndT2`
- New const `KeySourceMicrosoftKeyVault`
- New const `SmbNonBrowsableDisabled`
- New const `SmbNonBrowsableEnabled`
- New const `RegionStorageToNetworkProximityT1`
- New const `SmbAccessBasedEnumerationDisabled`
- New const `SmbAccessBasedEnumerationEnabled`
- New const `KeySourceMicrosoftNetApp`
- New const `RegionStorageToNetworkProximityT2`
- New const `KeyVaultStatusDeleted`
- New const `KeyVaultStatusInUse`
- New const `IdentityTypeSystemAssignedUserAssigned`
- New type alias `SmbAccessBasedEnumeration`
- New type alias `SmbNonBrowsable`
- New type alias `KeySource`
- New type alias `KeyVaultStatus`
- New type alias `IdentityType`
- New type alias `RegionStorageToNetworkProximity`
- New function `PossibleIdentityTypeValues() []IdentityType`
- New function `PossibleKeySourceValues() []KeySource`
- New function `*ResourceClient.QueryRegionInfo(context.Context, string, *ResourceClientQueryRegionInfoOptions) (ResourceClientQueryRegionInfoResponse, error)`
- New function `PossibleSmbAccessBasedEnumerationValues() []SmbAccessBasedEnumeration`
- New function `PossibleSmbNonBrowsableValues() []SmbNonBrowsable`
- New function `*AccountsClient.BeginRenewCredentials(context.Context, string, string, *AccountsClientBeginRenewCredentialsOptions) (*runtime.Poller[AccountsClientRenewCredentialsResponse], error)`
- New function `PossibleRegionStorageToNetworkProximityValues() []RegionStorageToNetworkProximity`
- New function `PossibleKeyVaultStatusValues() []KeyVaultStatus`
- New struct `AccountsClientBeginRenewCredentialsOptions`
- New struct `AccountsClientRenewCredentialsResponse`
- New struct `EncryptionIdentity`
- New struct `Identity`
- New struct `KeyVaultProperties`
- New struct `RegionInfo`
- New struct `RegionInfoAvailabilityZoneMappingsItem`
- New struct `RelocateVolumeRequest`
- New struct `ResourceClientQueryRegionInfoOptions`
- New struct `ResourceClientQueryRegionInfoResponse`
- New struct `UserAssignedIdentity`
- New field `Body` in struct `VolumesClientBeginRelocateOptions`
- New field `SmbAccessBasedEnumeration` in struct `VolumeProperties`
- New field `SmbNonBrowsable` in struct `VolumeProperties`
- New field `DeleteBaseSnapshot` in struct `VolumeProperties`
- New field `Identity` in struct `Account`
- New field `Identity` in struct `AccountEncryption`
- New field `KeyVaultProperties` in struct `AccountEncryption`
- New field `DisableShowmount` in struct `AccountProperties`


## 2.1.0 (2022-07-21)
### Features Added

- New const `EncryptionKeySourceMicrosoftKeyVault`
- New function `*VolumesClient.BeginReestablishReplication(context.Context, string, string, string, string, ReestablishReplicationRequest, *VolumesClientBeginReestablishReplicationOptions) (*runtime.Poller[VolumesClientReestablishReplicationResponse], error)`
- New struct `ReestablishReplicationRequest`
- New struct `VolumesClientBeginReestablishReplicationOptions`
- New struct `VolumesClientReestablishReplicationResponse`
- New field `CoolnessPeriod` in struct `VolumePatchProperties`
- New field `CoolAccess` in struct `VolumePatchProperties`
- New field `KeyVaultPrivateEndpointResourceID` in struct `VolumeProperties`
- New field `CoolAccess` in struct `PoolPatchProperties`


## 2.0.0 (2022-06-22)
### Breaking Changes

- Function `*BackupsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, string, string, string, *BackupsClientBeginUpdateOptions)` to `(context.Context, string, string, string, string, string, BackupPatch, *BackupsClientBeginUpdateOptions)`
- Type of `VolumeProperties.EncryptionKeySource` has been changed from `*string` to `*EncryptionKeySource`
- Field `Body` of struct `BackupsClientBeginUpdateOptions` has been removed
- Field `Tags` of struct `VolumeGroupDetails` has been removed
- Field `Tags` of struct `VolumeGroup` has been removed

### Features Added

- New const `ProvisioningStateDeleting`
- New const `TypeIndividualGroupQuota`
- New const `TypeDefaultUserQuota`
- New const `TypeIndividualUserQuota`
- New const `ProvisioningStateSucceeded`
- New const `ProvisioningStateAccepted`
- New const `ProvisioningStateCreating`
- New const `ProvisioningStateMoving`
- New const `ProvisioningStatePatching`
- New const `ProvisioningStateFailed`
- New const `EncryptionKeySourceMicrosoftNetApp`
- New const `TypeDefaultGroupQuota`
- New function `*VolumesClient.BeginRelocate(context.Context, string, string, string, string, *VolumesClientBeginRelocateOptions) (*runtime.Poller[VolumesClientRelocateResponse], error)`
- New function `*VolumeQuotaRulesClient.BeginDelete(context.Context, string, string, string, string, string, *VolumeQuotaRulesClientBeginDeleteOptions) (*runtime.Poller[VolumeQuotaRulesClientDeleteResponse], error)`
- New function `*VolumeQuotaRulesClient.Get(context.Context, string, string, string, string, string, *VolumeQuotaRulesClientGetOptions) (VolumeQuotaRulesClientGetResponse, error)`
- New function `*VolumesClient.NewListReplicationsPager(string, string, string, string, *VolumesClientListReplicationsOptions) *runtime.Pager[VolumesClientListReplicationsResponse]`
- New function `PossibleTypeValues() []Type`
- New function `*VolumesClient.BeginResetCifsPassword(context.Context, string, string, string, string, *VolumesClientBeginResetCifsPasswordOptions) (*runtime.Poller[VolumesClientResetCifsPasswordResponse], error)`
- New function `*VolumeQuotaRulesClient.BeginCreate(context.Context, string, string, string, string, string, VolumeQuotaRule, *VolumeQuotaRulesClientBeginCreateOptions) (*runtime.Poller[VolumeQuotaRulesClientCreateResponse], error)`
- New function `*VolumesClient.BeginFinalizeRelocation(context.Context, string, string, string, string, *VolumesClientBeginFinalizeRelocationOptions) (*runtime.Poller[VolumesClientFinalizeRelocationResponse], error)`
- New function `*VolumeQuotaRulesClient.NewListByVolumePager(string, string, string, string, *VolumeQuotaRulesClientListByVolumeOptions) *runtime.Pager[VolumeQuotaRulesClientListByVolumeResponse]`
- New function `PossibleProvisioningStateValues() []ProvisioningState`
- New function `*VolumesClient.BeginRevertRelocation(context.Context, string, string, string, string, *VolumesClientBeginRevertRelocationOptions) (*runtime.Poller[VolumesClientRevertRelocationResponse], error)`
- New function `*VolumeQuotaRulesClient.BeginUpdate(context.Context, string, string, string, string, string, VolumeQuotaRulePatch, *VolumeQuotaRulesClientBeginUpdateOptions) (*runtime.Poller[VolumeQuotaRulesClientUpdateResponse], error)`
- New function `PossibleEncryptionKeySourceValues() []EncryptionKeySource`
- New function `NewVolumeQuotaRulesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VolumeQuotaRulesClient, error)`
- New struct `ListReplications`
- New struct `Replication`
- New struct `TrackedResource`
- New struct `VolumeQuotaRule`
- New struct `VolumeQuotaRulePatch`
- New struct `VolumeQuotaRulesClient`
- New struct `VolumeQuotaRulesClientBeginCreateOptions`
- New struct `VolumeQuotaRulesClientBeginDeleteOptions`
- New struct `VolumeQuotaRulesClientBeginUpdateOptions`
- New struct `VolumeQuotaRulesClientCreateResponse`
- New struct `VolumeQuotaRulesClientDeleteResponse`
- New struct `VolumeQuotaRulesClientGetOptions`
- New struct `VolumeQuotaRulesClientGetResponse`
- New struct `VolumeQuotaRulesClientListByVolumeOptions`
- New struct `VolumeQuotaRulesClientListByVolumeResponse`
- New struct `VolumeQuotaRulesClientUpdateResponse`
- New struct `VolumeQuotaRulesList`
- New struct `VolumeQuotaRulesProperties`
- New struct `VolumeRelocationProperties`
- New struct `VolumesClientBeginFinalizeRelocationOptions`
- New struct `VolumesClientBeginRelocateOptions`
- New struct `VolumesClientBeginResetCifsPasswordOptions`
- New struct `VolumesClientBeginRevertRelocationOptions`
- New struct `VolumesClientFinalizeRelocationResponse`
- New struct `VolumesClientListReplicationsOptions`
- New struct `VolumesClientListReplicationsResponse`
- New struct `VolumesClientRelocateResponse`
- New struct `VolumesClientResetCifsPasswordResponse`
- New struct `VolumesClientRevertRelocationResponse`
- New field `SystemData` in struct `ProxyResource`
- New field `Zones` in struct `Volume`
- New field `SystemData` in struct `Resource`
- New field `Encrypted` in struct `VolumeProperties`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
