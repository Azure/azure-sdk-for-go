# Release History

## 8.1.0-beta.1 (2025-12-03)
### Features Added

- New enum type `ActualRansomwareProtectionState` with values `ActualRansomwareProtectionStateDisabled`, `ActualRansomwareProtectionStateEnabled`, `ActualRansomwareProtectionStateLearning`, `ActualRansomwareProtectionStatePaused`
- New enum type `BreakthroughMode` with values `BreakthroughModeDisabled`, `BreakthroughModeEnabled`
- New enum type `BucketPatchPermissions` with values `BucketPatchPermissionsReadOnly`, `BucketPatchPermissionsReadWrite`
- New enum type `BucketPermissions` with values `BucketPermissionsReadOnly`, `BucketPermissionsReadWrite`
- New enum type `CacheLifeCycleState` with values `CacheLifeCycleStateClusterPeeringOfferSent`, `CacheLifeCycleStateCreating`, `CacheLifeCycleStateFailed`, `CacheLifeCycleStateSucceeded`, `CacheLifeCycleStateVserverPeeringOfferSent`
- New enum type `CacheProvisioningState` with values `CacheProvisioningStateCanceled`, `CacheProvisioningStateCreating`, `CacheProvisioningStateDeleting`, `CacheProvisioningStateFailed`, `CacheProvisioningStateSucceeded`, `CacheProvisioningStateUpdating`
- New enum type `CheckElasticResourceAvailabilityReason` with values `CheckElasticResourceAvailabilityReasonAlreadyExists`, `CheckElasticResourceAvailabilityReasonInvalid`
- New enum type `CheckElasticResourceAvailabilityStatus` with values `CheckElasticResourceAvailabilityStatusFalse`, `CheckElasticResourceAvailabilityStatusTrue`
- New enum type `CifsChangeNotifyState` with values `CifsChangeNotifyStateDisabled`, `CifsChangeNotifyStateEnabled`
- New enum type `CredentialsStatus` with values `CredentialsStatusActive`, `CredentialsStatusCredentialsExpired`, `CredentialsStatusNoCredentialsSet`
- New enum type `DayOfWeek` with values `DayOfWeekFriday`, `DayOfWeekMonday`, `DayOfWeekSaturday`, `DayOfWeekSunday`, `DayOfWeekThursday`, `DayOfWeekTuesday`, `DayOfWeekWednesday`
- New enum type `DesiredRansomwareProtectionState` with values `DesiredRansomwareProtectionStateDisabled`, `DesiredRansomwareProtectionStateEnabled`
- New enum type `ElasticBackupPolicyState` with values `ElasticBackupPolicyStateDisabled`, `ElasticBackupPolicyStateEnabled`
- New enum type `ElasticBackupType` with values `ElasticBackupTypeManual`, `ElasticBackupTypeScheduled`
- New enum type `ElasticKeyVaultStatus` with values `ElasticKeyVaultStatusCreated`, `ElasticKeyVaultStatusDeleted`, `ElasticKeyVaultStatusError`, `ElasticKeyVaultStatusInUse`, `ElasticKeyVaultStatusUpdating`
- New enum type `ElasticNfsv3Access` with values `ElasticNfsv3AccessDisabled`, `ElasticNfsv3AccessEnabled`
- New enum type `ElasticNfsv4Access` with values `ElasticNfsv4AccessDisabled`, `ElasticNfsv4AccessEnabled`
- New enum type `ElasticPoolEncryptionKeySource` with values `ElasticPoolEncryptionKeySourceKeyVault`, `ElasticPoolEncryptionKeySourceNetApp`
- New enum type `ElasticProtocolType` with values `ElasticProtocolTypeNFSv3`, `ElasticProtocolTypeNFSv4`, `ElasticProtocolTypeSMB`
- New enum type `ElasticResourceAvailabilityStatus` with values `ElasticResourceAvailabilityStatusOffline`, `ElasticResourceAvailabilityStatusOnline`
- New enum type `ElasticRootAccess` with values `ElasticRootAccessDisabled`, `ElasticRootAccessEnabled`
- New enum type `ElasticServiceLevel` with values `ElasticServiceLevelZoneRedundant`
- New enum type `ElasticSmbEncryption` with values `ElasticSmbEncryptionDisabled`, `ElasticSmbEncryptionEnabled`
- New enum type `ElasticUnixAccessRule` with values `ElasticUnixAccessRuleNoAccess`, `ElasticUnixAccessRuleReadOnly`, `ElasticUnixAccessRuleReadWrite`
- New enum type `ElasticVolumePolicyEnforcement` with values `ElasticVolumePolicyEnforcementEnforced`, `ElasticVolumePolicyEnforcementNotEnforced`
- New enum type `ElasticVolumeRestorationState` with values `ElasticVolumeRestorationStateFailed`, `ElasticVolumeRestorationStateRestored`, `ElasticVolumeRestorationStateRestoring`
- New enum type `EnableWriteBackState` with values `EnableWriteBackStateDisabled`, `EnableWriteBackStateEnabled`
- New enum type `EncryptionState` with values `EncryptionStateDisabled`, `EncryptionStateEnabled`
- New enum type `ExternalReplicationSetupStatus` with values `ExternalReplicationSetupStatusClusterPeerPending`, `ExternalReplicationSetupStatusClusterPeerRequired`, `ExternalReplicationSetupStatusNoActionRequired`, `ExternalReplicationSetupStatusReplicationCreateRequired`, `ExternalReplicationSetupStatusVServerPeerRequired`
- New enum type `GlobalFileLockingState` with values `GlobalFileLockingStateDisabled`, `GlobalFileLockingStateEnabled`
- New enum type `KerberosState` with values `KerberosStateDisabled`, `KerberosStateEnabled`
- New enum type `LargeVolumeType` with values `LargeVolumeTypeExtraLargeVolume7Dot2PiB`, `LargeVolumeTypeLargeVolume`
- New enum type `LdapServerType` with values `LdapServerTypeActiveDirectory`, `LdapServerTypeOpenLDAP`
- New enum type `LdapState` with values `LdapStateDisabled`, `LdapStateEnabled`
- New enum type `PolicyStatus` with values `PolicyStatusDisabled`, `PolicyStatusEnabled`
- New enum type `ProtocolTypes` with values `ProtocolTypesNFSv3`, `ProtocolTypesNFSv4`, `ProtocolTypesSMB`
- New enum type `RansomwareReportSeverity` with values `RansomwareReportSeverityHigh`, `RansomwareReportSeverityLow`, `RansomwareReportSeverityModerate`, `RansomwareReportSeverityNone`
- New enum type `RansomwareReportState` with values `RansomwareReportStateActive`, `RansomwareReportStateResolved`
- New enum type `RansomwareSuspectResolution` with values `RansomwareSuspectResolutionFalsePositive`, `RansomwareSuspectResolutionPotentialThreat`
- New enum type `SmbEncryptionState` with values `SmbEncryptionStateDisabled`, `SmbEncryptionStateEnabled`
- New enum type `SnapshotDirectoryVisibility` with values `SnapshotDirectoryVisibilityHidden`, `SnapshotDirectoryVisibilityVisible`
- New enum type `SnapshotUsage` with values `SnapshotUsageCreateNewSnapshot`, `SnapshotUsageUseExistingSnapshot`
- New enum type `VolumeLanguage` with values `VolumeLanguageAr`, `VolumeLanguageArUTF8`, `VolumeLanguageC`, `VolumeLanguageCUTF8`, `VolumeLanguageCs`, `VolumeLanguageCsUTF8`, `VolumeLanguageDa`, `VolumeLanguageDaUTF8`, `VolumeLanguageDe`, `VolumeLanguageDeUTF8`, `VolumeLanguageEn`, `VolumeLanguageEnUTF8`, `VolumeLanguageEnUs`, `VolumeLanguageEnUsUTF8`, `VolumeLanguageEs`, `VolumeLanguageEsUTF8`, `VolumeLanguageFi`, `VolumeLanguageFiUTF8`, `VolumeLanguageFr`, `VolumeLanguageFrUTF8`, `VolumeLanguageHe`, `VolumeLanguageHeUTF8`, `VolumeLanguageHr`, `VolumeLanguageHrUTF8`, `VolumeLanguageHu`, `VolumeLanguageHuUTF8`, `VolumeLanguageIt`, `VolumeLanguageItUTF8`, `VolumeLanguageJa`, `VolumeLanguageJaJp932`, `VolumeLanguageJaJp932UTF8`, `VolumeLanguageJaJpPck`, `VolumeLanguageJaJpPckUTF8`, `VolumeLanguageJaJpPckV2`, `VolumeLanguageJaJpPckV2UTF8`, `VolumeLanguageJaUTF8`, `VolumeLanguageJaV1`, `VolumeLanguageJaV1UTF8`, `VolumeLanguageKo`, `VolumeLanguageKoUTF8`, `VolumeLanguageNl`, `VolumeLanguageNlUTF8`, `VolumeLanguageNo`, `VolumeLanguageNoUTF8`, `VolumeLanguagePl`, `VolumeLanguagePlUTF8`, `VolumeLanguagePt`, `VolumeLanguagePtUTF8`, `VolumeLanguageRo`, `VolumeLanguageRoUTF8`, `VolumeLanguageRu`, `VolumeLanguageRuUTF8`, `VolumeLanguageSk`, `VolumeLanguageSkUTF8`, `VolumeLanguageSl`, `VolumeLanguageSlUTF8`, `VolumeLanguageSv`, `VolumeLanguageSvUTF8`, `VolumeLanguageTr`, `VolumeLanguageTrUTF8`, `VolumeLanguageUTF8Mb4`, `VolumeLanguageZh`, `VolumeLanguageZhGbk`, `VolumeLanguageZhGbkUTF8`, `VolumeLanguageZhTw`, `VolumeLanguageZhTwBig5`, `VolumeLanguageZhTwBig5UTF8`, `VolumeLanguageZhTwUTF8`, `VolumeLanguageZhUTF8`
- New enum type `VolumeSize` with values `VolumeSizeLarge`, `VolumeSizeRegular`
- New function `NewActiveDirectoryConfigsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ActiveDirectoryConfigsClient, error)`
- New function `*ActiveDirectoryConfigsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, activeDirectoryConfigName string, body ActiveDirectoryConfig, options *ActiveDirectoryConfigsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ActiveDirectoryConfigsClientCreateOrUpdateResponse], error)`
- New function `*ActiveDirectoryConfigsClient.BeginDelete(ctx context.Context, resourceGroupName string, activeDirectoryConfigName string, options *ActiveDirectoryConfigsClientBeginDeleteOptions) (*runtime.Poller[ActiveDirectoryConfigsClientDeleteResponse], error)`
- New function `*ActiveDirectoryConfigsClient.Get(ctx context.Context, resourceGroupName string, activeDirectoryConfigName string, options *ActiveDirectoryConfigsClientGetOptions) (ActiveDirectoryConfigsClientGetResponse, error)`
- New function `*ActiveDirectoryConfigsClient.NewListByResourceGroupPager(resourceGroupName string, options *ActiveDirectoryConfigsClientListByResourceGroupOptions) *runtime.Pager[ActiveDirectoryConfigsClientListByResourceGroupResponse]`
- New function `*ActiveDirectoryConfigsClient.NewListBySubscriptionPager(options *ActiveDirectoryConfigsClientListBySubscriptionOptions) *runtime.Pager[ActiveDirectoryConfigsClientListBySubscriptionResponse]`
- New function `*ActiveDirectoryConfigsClient.BeginUpdate(ctx context.Context, resourceGroupName string, activeDirectoryConfigName string, body ActiveDirectoryConfigUpdate, options *ActiveDirectoryConfigsClientBeginUpdateOptions) (*runtime.Poller[ActiveDirectoryConfigsClientUpdateResponse], error)`
- New function `NewBucketsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*BucketsClient, error)`
- New function `*BucketsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, bucketName string, body Bucket, options *BucketsClientBeginCreateOrUpdateOptions) (*runtime.Poller[BucketsClientCreateOrUpdateResponse], error)`
- New function `*BucketsClient.BeginDelete(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, bucketName string, options *BucketsClientBeginDeleteOptions) (*runtime.Poller[BucketsClientDeleteResponse], error)`
- New function `*BucketsClient.GenerateCredentials(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, bucketName string, body BucketCredentialsExpiry, options *BucketsClientGenerateCredentialsOptions) (BucketsClientGenerateCredentialsResponse, error)`
- New function `*BucketsClient.Get(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, bucketName string, options *BucketsClientGetOptions) (BucketsClientGetResponse, error)`
- New function `*BucketsClient.NewListPager(resourceGroupName string, accountName string, poolName string, volumeName string, options *BucketsClientListOptions) *runtime.Pager[BucketsClientListResponse]`
- New function `*BucketsClient.BeginUpdate(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, bucketName string, body BucketPatch, options *BucketsClientBeginUpdateOptions) (*runtime.Poller[BucketsClientUpdateResponse], error)`
- New function `NewCachesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CachesClient, error)`
- New function `*CachesClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, accountName string, poolName string, cacheName string, body Cache, options *CachesClientBeginCreateOrUpdateOptions) (*runtime.Poller[CachesClientCreateOrUpdateResponse], error)`
- New function `*CachesClient.BeginDelete(ctx context.Context, resourceGroupName string, accountName string, poolName string, cacheName string, options *CachesClientBeginDeleteOptions) (*runtime.Poller[CachesClientDeleteResponse], error)`
- New function `*CachesClient.Get(ctx context.Context, resourceGroupName string, accountName string, poolName string, cacheName string, options *CachesClientGetOptions) (CachesClientGetResponse, error)`
- New function `*CachesClient.NewListByCapacityPoolsPager(resourceGroupName string, accountName string, poolName string, options *CachesClientListByCapacityPoolsOptions) *runtime.Pager[CachesClientListByCapacityPoolsResponse]`
- New function `*CachesClient.ListPeeringPassphrases(ctx context.Context, resourceGroupName string, accountName string, poolName string, cacheName string, options *CachesClientListPeeringPassphrasesOptions) (CachesClientListPeeringPassphrasesResponse, error)`
- New function `*CachesClient.BeginPoolChange(ctx context.Context, resourceGroupName string, accountName string, poolName string, cacheName string, body PoolChangeRequest, options *CachesClientBeginPoolChangeOptions) (*runtime.Poller[CachesClientPoolChangeResponse], error)`
- New function `*CachesClient.BeginUpdate(ctx context.Context, resourceGroupName string, accountName string, poolName string, cacheName string, body CacheUpdate, options *CachesClientBeginUpdateOptions) (*runtime.Poller[CachesClientUpdateResponse], error)`
- New function `*ClientFactory.NewActiveDirectoryConfigsClient() *ActiveDirectoryConfigsClient`
- New function `*ClientFactory.NewBucketsClient() *BucketsClient`
- New function `*ClientFactory.NewCachesClient() *CachesClient`
- New function `*ClientFactory.NewElasticAccountsClient() *ElasticAccountsClient`
- New function `*ClientFactory.NewElasticBackupPoliciesClient() *ElasticBackupPoliciesClient`
- New function `*ClientFactory.NewElasticBackupVaultsClient() *ElasticBackupVaultsClient`
- New function `*ClientFactory.NewElasticBackupsClient() *ElasticBackupsClient`
- New function `*ClientFactory.NewElasticCapacityPoolsClient() *ElasticCapacityPoolsClient`
- New function `*ClientFactory.NewElasticSnapshotPoliciesClient() *ElasticSnapshotPoliciesClient`
- New function `*ClientFactory.NewElasticSnapshotsClient() *ElasticSnapshotsClient`
- New function `*ClientFactory.NewElasticVolumesClient() *ElasticVolumesClient`
- New function `*ClientFactory.NewRansomwareReportsClient() *RansomwareReportsClient`
- New function `NewElasticAccountsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ElasticAccountsClient, error)`
- New function `*ElasticAccountsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, accountName string, body ElasticAccount, options *ElasticAccountsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ElasticAccountsClientCreateOrUpdateResponse], error)`
- New function `*ElasticAccountsClient.BeginDelete(ctx context.Context, resourceGroupName string, accountName string, options *ElasticAccountsClientBeginDeleteOptions) (*runtime.Poller[ElasticAccountsClientDeleteResponse], error)`
- New function `*ElasticAccountsClient.Get(ctx context.Context, resourceGroupName string, accountName string, options *ElasticAccountsClientGetOptions) (ElasticAccountsClientGetResponse, error)`
- New function `*ElasticAccountsClient.NewListByResourceGroupPager(resourceGroupName string, options *ElasticAccountsClientListByResourceGroupOptions) *runtime.Pager[ElasticAccountsClientListByResourceGroupResponse]`
- New function `*ElasticAccountsClient.NewListBySubscriptionPager(options *ElasticAccountsClientListBySubscriptionOptions) *runtime.Pager[ElasticAccountsClientListBySubscriptionResponse]`
- New function `*ElasticAccountsClient.BeginUpdate(ctx context.Context, resourceGroupName string, accountName string, body ElasticAccountUpdate, options *ElasticAccountsClientBeginUpdateOptions) (*runtime.Poller[ElasticAccountsClientUpdateResponse], error)`
- New function `NewElasticBackupPoliciesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ElasticBackupPoliciesClient, error)`
- New function `*ElasticBackupPoliciesClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, accountName string, backupPolicyName string, body ElasticBackupPolicy, options *ElasticBackupPoliciesClientBeginCreateOrUpdateOptions) (*runtime.Poller[ElasticBackupPoliciesClientCreateOrUpdateResponse], error)`
- New function `*ElasticBackupPoliciesClient.BeginDelete(ctx context.Context, resourceGroupName string, accountName string, backupPolicyName string, options *ElasticBackupPoliciesClientBeginDeleteOptions) (*runtime.Poller[ElasticBackupPoliciesClientDeleteResponse], error)`
- New function `*ElasticBackupPoliciesClient.Get(ctx context.Context, resourceGroupName string, accountName string, backupPolicyName string, options *ElasticBackupPoliciesClientGetOptions) (ElasticBackupPoliciesClientGetResponse, error)`
- New function `*ElasticBackupPoliciesClient.NewListByElasticAccountPager(resourceGroupName string, accountName string, options *ElasticBackupPoliciesClientListByElasticAccountOptions) *runtime.Pager[ElasticBackupPoliciesClientListByElasticAccountResponse]`
- New function `*ElasticBackupPoliciesClient.BeginUpdate(ctx context.Context, resourceGroupName string, accountName string, backupPolicyName string, body ElasticBackupPolicyUpdate, options *ElasticBackupPoliciesClientBeginUpdateOptions) (*runtime.Poller[ElasticBackupPoliciesClientUpdateResponse], error)`
- New function `NewElasticBackupVaultsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ElasticBackupVaultsClient, error)`
- New function `*ElasticBackupVaultsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, accountName string, backupVaultName string, body ElasticBackupVault, options *ElasticBackupVaultsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ElasticBackupVaultsClientCreateOrUpdateResponse], error)`
- New function `*ElasticBackupVaultsClient.BeginDelete(ctx context.Context, resourceGroupName string, accountName string, backupVaultName string, options *ElasticBackupVaultsClientBeginDeleteOptions) (*runtime.Poller[ElasticBackupVaultsClientDeleteResponse], error)`
- New function `*ElasticBackupVaultsClient.Get(ctx context.Context, resourceGroupName string, accountName string, backupVaultName string, options *ElasticBackupVaultsClientGetOptions) (ElasticBackupVaultsClientGetResponse, error)`
- New function `*ElasticBackupVaultsClient.NewListByElasticAccountPager(resourceGroupName string, accountName string, options *ElasticBackupVaultsClientListByElasticAccountOptions) *runtime.Pager[ElasticBackupVaultsClientListByElasticAccountResponse]`
- New function `*ElasticBackupVaultsClient.BeginUpdate(ctx context.Context, resourceGroupName string, accountName string, backupVaultName string, body ElasticBackupVaultUpdate, options *ElasticBackupVaultsClientBeginUpdateOptions) (*runtime.Poller[ElasticBackupVaultsClientUpdateResponse], error)`
- New function `NewElasticBackupsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ElasticBackupsClient, error)`
- New function `*ElasticBackupsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, accountName string, backupVaultName string, backupName string, body ElasticBackup, options *ElasticBackupsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ElasticBackupsClientCreateOrUpdateResponse], error)`
- New function `*ElasticBackupsClient.BeginDelete(ctx context.Context, resourceGroupName string, accountName string, backupVaultName string, backupName string, options *ElasticBackupsClientBeginDeleteOptions) (*runtime.Poller[ElasticBackupsClientDeleteResponse], error)`
- New function `*ElasticBackupsClient.Get(ctx context.Context, resourceGroupName string, accountName string, backupVaultName string, backupName string, options *ElasticBackupsClientGetOptions) (ElasticBackupsClientGetResponse, error)`
- New function `*ElasticBackupsClient.NewListByVaultPager(resourceGroupName string, accountName string, backupVaultName string, options *ElasticBackupsClientListByVaultOptions) *runtime.Pager[ElasticBackupsClientListByVaultResponse]`
- New function `*ElasticBackupsClient.BeginUpdate(ctx context.Context, resourceGroupName string, accountName string, backupVaultName string, backupName string, body ElasticBackup, options *ElasticBackupsClientBeginUpdateOptions) (*runtime.Poller[ElasticBackupsClientUpdateResponse], error)`
- New function `NewElasticCapacityPoolsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ElasticCapacityPoolsClient, error)`
- New function `*ElasticCapacityPoolsClient.BeginChangeZone(ctx context.Context, resourceGroupName string, accountName string, poolName string, body ChangeZoneRequest, options *ElasticCapacityPoolsClientBeginChangeZoneOptions) (*runtime.Poller[ElasticCapacityPoolsClientChangeZoneResponse], error)`
- New function `*ElasticCapacityPoolsClient.CheckVolumeFilePathAvailability(ctx context.Context, resourceGroupName string, accountName string, poolName string, body CheckElasticVolumeFilePathAvailabilityRequest, options *ElasticCapacityPoolsClientCheckVolumeFilePathAvailabilityOptions) (ElasticCapacityPoolsClientCheckVolumeFilePathAvailabilityResponse, error)`
- New function `*ElasticCapacityPoolsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, accountName string, poolName string, body ElasticCapacityPool, options *ElasticCapacityPoolsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ElasticCapacityPoolsClientCreateOrUpdateResponse], error)`
- New function `*ElasticCapacityPoolsClient.BeginDelete(ctx context.Context, resourceGroupName string, accountName string, poolName string, options *ElasticCapacityPoolsClientBeginDeleteOptions) (*runtime.Poller[ElasticCapacityPoolsClientDeleteResponse], error)`
- New function `*ElasticCapacityPoolsClient.Get(ctx context.Context, resourceGroupName string, accountName string, poolName string, options *ElasticCapacityPoolsClientGetOptions) (ElasticCapacityPoolsClientGetResponse, error)`
- New function `*ElasticCapacityPoolsClient.NewListByElasticAccountPager(resourceGroupName string, accountName string, options *ElasticCapacityPoolsClientListByElasticAccountOptions) *runtime.Pager[ElasticCapacityPoolsClientListByElasticAccountResponse]`
- New function `*ElasticCapacityPoolsClient.BeginUpdate(ctx context.Context, resourceGroupName string, accountName string, poolName string, body ElasticCapacityPoolUpdate, options *ElasticCapacityPoolsClientBeginUpdateOptions) (*runtime.Poller[ElasticCapacityPoolsClientUpdateResponse], error)`
- New function `NewElasticSnapshotPoliciesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ElasticSnapshotPoliciesClient, error)`
- New function `*ElasticSnapshotPoliciesClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, accountName string, snapshotPolicyName string, body ElasticSnapshotPolicy, options *ElasticSnapshotPoliciesClientBeginCreateOrUpdateOptions) (*runtime.Poller[ElasticSnapshotPoliciesClientCreateOrUpdateResponse], error)`
- New function `*ElasticSnapshotPoliciesClient.BeginDelete(ctx context.Context, resourceGroupName string, accountName string, snapshotPolicyName string, options *ElasticSnapshotPoliciesClientBeginDeleteOptions) (*runtime.Poller[ElasticSnapshotPoliciesClientDeleteResponse], error)`
- New function `*ElasticSnapshotPoliciesClient.Get(ctx context.Context, resourceGroupName string, accountName string, snapshotPolicyName string, options *ElasticSnapshotPoliciesClientGetOptions) (ElasticSnapshotPoliciesClientGetResponse, error)`
- New function `*ElasticSnapshotPoliciesClient.NewListByElasticAccountPager(resourceGroupName string, accountName string, options *ElasticSnapshotPoliciesClientListByElasticAccountOptions) *runtime.Pager[ElasticSnapshotPoliciesClientListByElasticAccountResponse]`
- New function `*ElasticSnapshotPoliciesClient.NewListElasticVolumesPager(resourceGroupName string, accountName string, snapshotPolicyName string, options *ElasticSnapshotPoliciesClientListElasticVolumesOptions) *runtime.Pager[ElasticSnapshotPoliciesClientListElasticVolumesResponse]`
- New function `*ElasticSnapshotPoliciesClient.BeginUpdate(ctx context.Context, resourceGroupName string, accountName string, snapshotPolicyName string, body ElasticSnapshotPolicyUpdate, options *ElasticSnapshotPoliciesClientBeginUpdateOptions) (*runtime.Poller[ElasticSnapshotPoliciesClientUpdateResponse], error)`
- New function `NewElasticSnapshotsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ElasticSnapshotsClient, error)`
- New function `*ElasticSnapshotsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, snapshotName string, body ElasticSnapshot, options *ElasticSnapshotsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ElasticSnapshotsClientCreateOrUpdateResponse], error)`
- New function `*ElasticSnapshotsClient.BeginDelete(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, snapshotName string, options *ElasticSnapshotsClientBeginDeleteOptions) (*runtime.Poller[ElasticSnapshotsClientDeleteResponse], error)`
- New function `*ElasticSnapshotsClient.Get(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, snapshotName string, options *ElasticSnapshotsClientGetOptions) (ElasticSnapshotsClientGetResponse, error)`
- New function `*ElasticSnapshotsClient.NewListByElasticVolumePager(resourceGroupName string, accountName string, poolName string, volumeName string, options *ElasticSnapshotsClientListByElasticVolumeOptions) *runtime.Pager[ElasticSnapshotsClientListByElasticVolumeResponse]`
- New function `NewElasticVolumesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ElasticVolumesClient, error)`
- New function `*ElasticVolumesClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, body ElasticVolume, options *ElasticVolumesClientBeginCreateOrUpdateOptions) (*runtime.Poller[ElasticVolumesClientCreateOrUpdateResponse], error)`
- New function `*ElasticVolumesClient.BeginDelete(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, options *ElasticVolumesClientBeginDeleteOptions) (*runtime.Poller[ElasticVolumesClientDeleteResponse], error)`
- New function `*ElasticVolumesClient.Get(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, options *ElasticVolumesClientGetOptions) (ElasticVolumesClientGetResponse, error)`
- New function `*ElasticVolumesClient.NewListByElasticPoolPager(resourceGroupName string, accountName string, poolName string, options *ElasticVolumesClientListByElasticPoolOptions) *runtime.Pager[ElasticVolumesClientListByElasticPoolResponse]`
- New function `*ElasticVolumesClient.BeginRevert(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, body ElasticVolumeRevert, options *ElasticVolumesClientBeginRevertOptions) (*runtime.Poller[ElasticVolumesClientRevertResponse], error)`
- New function `*ElasticVolumesClient.BeginUpdate(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, body ElasticVolumeUpdate, options *ElasticVolumesClientBeginUpdateOptions) (*runtime.Poller[ElasticVolumesClientUpdateResponse], error)`
- New function `NewRansomwareReportsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*RansomwareReportsClient, error)`
- New function `*RansomwareReportsClient.BeginClearSuspects(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, ransomwareReportName string, body RansomwareSuspectsClearRequest, options *RansomwareReportsClientBeginClearSuspectsOptions) (*runtime.Poller[RansomwareReportsClientClearSuspectsResponse], error)`
- New function `*RansomwareReportsClient.Get(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, ransomwareReportName string, options *RansomwareReportsClientGetOptions) (RansomwareReportsClientGetResponse, error)`
- New function `*RansomwareReportsClient.NewListPager(resourceGroupName string, accountName string, poolName string, volumeName string, options *RansomwareReportsClientListOptions) *runtime.Pager[RansomwareReportsClientListResponse]`
- New function `*VolumesClient.BeginListQuotaReport(ctx context.Context, resourceGroupName string, accountName string, poolName string, volumeName string, options *VolumesClientBeginListQuotaReportOptions) (*runtime.Poller[VolumesClientListQuotaReportResponse], error)`
- New struct `ActiveDirectoryConfig`
- New struct `ActiveDirectoryConfigListResult`
- New struct `ActiveDirectoryConfigProperties`
- New struct `ActiveDirectoryConfigUpdate`
- New struct `ActiveDirectoryConfigUpdateProperties`
- New struct `Bucket`
- New struct `BucketCredentialsExpiry`
- New struct `BucketGenerateCredentials`
- New struct `BucketList`
- New struct `BucketPatch`
- New struct `BucketPatchProperties`
- New struct `BucketProperties`
- New struct `BucketServerPatchProperties`
- New struct `BucketServerProperties`
- New struct `Cache`
- New struct `CacheList`
- New struct `CacheMountTargetProperties`
- New struct `CacheProperties`
- New struct `CachePropertiesExportPolicy`
- New struct `CacheUpdate`
- New struct `CacheUpdateProperties`
- New struct `ChangeZoneRequest`
- New struct `CheckElasticResourceAvailabilityResponse`
- New struct `CheckElasticVolumeFilePathAvailabilityRequest`
- New struct `CifsUser`
- New struct `ElasticAccount`
- New struct `ElasticAccountListResult`
- New struct `ElasticAccountProperties`
- New struct `ElasticAccountUpdate`
- New struct `ElasticAccountUpdateProperties`
- New struct `ElasticBackup`
- New struct `ElasticBackupListResult`
- New struct `ElasticBackupPolicy`
- New struct `ElasticBackupPolicyListResult`
- New struct `ElasticBackupPolicyProperties`
- New struct `ElasticBackupPolicyUpdate`
- New struct `ElasticBackupPolicyUpdateProperties`
- New struct `ElasticBackupProperties`
- New struct `ElasticBackupVault`
- New struct `ElasticBackupVaultListResult`
- New struct `ElasticBackupVaultProperties`
- New struct `ElasticBackupVaultUpdate`
- New struct `ElasticCapacityPool`
- New struct `ElasticCapacityPoolListResult`
- New struct `ElasticCapacityPoolProperties`
- New struct `ElasticCapacityPoolUpdate`
- New struct `ElasticCapacityPoolUpdateProperties`
- New struct `ElasticEncryption`
- New struct `ElasticEncryptionConfiguration`
- New struct `ElasticEncryptionIdentity`
- New struct `ElasticExportPolicy`
- New struct `ElasticExportPolicyRule`
- New struct `ElasticKeyVaultProperties`
- New struct `ElasticMountTargetProperties`
- New struct `ElasticSmbPatchProperties`
- New struct `ElasticSmbProperties`
- New struct `ElasticSnapshot`
- New struct `ElasticSnapshotListResult`
- New struct `ElasticSnapshotPolicy`
- New struct `ElasticSnapshotPolicyDailySchedule`
- New struct `ElasticSnapshotPolicyHourlySchedule`
- New struct `ElasticSnapshotPolicyListResult`
- New struct `ElasticSnapshotPolicyMonthlySchedule`
- New struct `ElasticSnapshotPolicyProperties`
- New struct `ElasticSnapshotPolicyUpdate`
- New struct `ElasticSnapshotPolicyUpdateProperties`
- New struct `ElasticSnapshotPolicyVolumeList`
- New struct `ElasticSnapshotPolicyWeeklySchedule`
- New struct `ElasticSnapshotProperties`
- New struct `ElasticVolume`
- New struct `ElasticVolumeBackupProperties`
- New struct `ElasticVolumeDataProtectionPatchProperties`
- New struct `ElasticVolumeDataProtectionProperties`
- New struct `ElasticVolumeListResult`
- New struct `ElasticVolumeProperties`
- New struct `ElasticVolumeRevert`
- New struct `ElasticVolumeSnapshotProperties`
- New struct `ElasticVolumeUpdate`
- New struct `ElasticVolumeUpdateProperties`
- New struct `FileSystemUser`
- New struct `LdapConfiguration`
- New struct `ListQuotaReportResponse`
- New struct `NfsUser`
- New struct `OriginClusterInformation`
- New struct `PeeringPassphrases`
- New struct `QuotaReport`
- New struct `RansomwareProtectionPatchSettings`
- New struct `RansomwareProtectionSettings`
- New struct `RansomwareReport`
- New struct `RansomwareReportProperties`
- New struct `RansomwareReportsList`
- New struct `RansomwareSuspects`
- New struct `RansomwareSuspectsClearRequest`
- New struct `SecretPassword`
- New struct `SecretPasswordIdentity`
- New struct `SecretPasswordKeyVaultProperties`
- New struct `SmbSettings`
- New struct `SuspectFile`
- New field `LdapConfiguration` in struct `AccountProperties`
- New field `ExternalReplicationSetupInfo`, `ExternalReplicationSetupStatus`, `MirrorState`, `RelationshipStatus` in struct `ReplicationObject`
- New field `RansomwareProtection` in struct `VolumePatchPropertiesDataProtection`
- New field `BreakthroughMode`, `Language`, `LargeVolumeType`, `LdapServerType` in struct `VolumeProperties`
- New field `RansomwareProtection` in struct `VolumePropertiesDataProtection`


## 8.0.0 (2025-11-14)
### Breaking Changes

- Function `*SnapshotsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, string, string, string, any, *SnapshotsClientBeginUpdateOptions)` to `(context.Context, string, string, string, string, string, SnapshotPatch, *SnapshotsClientBeginUpdateOptions)`
- Type of `BackupStatus.RelationshipStatus` has been changed from `*RelationshipStatus` to `*VolumeBackupRelationshipStatus`
- Type of `PoolPatchProperties.CustomThroughputMibps` has been changed from `*float32` to `*int32`
- Type of `PoolProperties.CustomThroughputMibps` has been changed from `*float32` to `*int32`
- Type of `ReplicationStatus.RelationshipStatus` has been changed from `*RelationshipStatus` to `*VolumeReplicationRelationshipStatus`
- Type of `RestoreStatus.RelationshipStatus` has been changed from `*RelationshipStatus` to `*VolumeRestoreRelationshipStatus`
- Enum `RelationshipStatus` has been removed

### Features Added

- New value `CheckNameResourceTypesMicrosoftNetAppNetAppAccountsBackupVaultsBackups`, `CheckNameResourceTypesMicrosoftNetAppNetAppAccountsCapacityPoolsVolumesBackups` added to enum type `CheckNameResourceTypes`
- New value `CheckQuotaNameResourceTypesMicrosoftNetAppNetAppAccountsBackupVaultsBackups`, `CheckQuotaNameResourceTypesMicrosoftNetAppNetAppAccountsCapacityPoolsVolumesBackups` added to enum type `CheckQuotaNameResourceTypes`
- New value `ProvisioningStateUpdating` added to enum type `ProvisioningState`
- New enum type `Exclude` with values `ExcludeDeleted`, `ExcludeNone`
- New enum type `ReplicationMirrorState` with values `ReplicationMirrorStateBroken`, `ReplicationMirrorStateMirrored`, `ReplicationMirrorStateUninitialized`
- New enum type `VolumeBackupRelationshipStatus` with values `VolumeBackupRelationshipStatusFailed`, `VolumeBackupRelationshipStatusIdle`, `VolumeBackupRelationshipStatusTransferring`, `VolumeBackupRelationshipStatusUnknown`
- New enum type `VolumeReplicationRelationshipStatus` with values `VolumeReplicationRelationshipStatusIdle`, `VolumeReplicationRelationshipStatusTransferring`
- New enum type `VolumeRestoreRelationshipStatus` with values `VolumeRestoreRelationshipStatusFailed`, `VolumeRestoreRelationshipStatusIdle`, `VolumeRestoreRelationshipStatusTransferring`, `VolumeRestoreRelationshipStatusUnknown`
- New function `*ClientFactory.NewResourceQuotaLimitsAccountClient() *ResourceQuotaLimitsAccountClient`
- New function `NewResourceQuotaLimitsAccountClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ResourceQuotaLimitsAccountClient, error)`
- New function `*ResourceQuotaLimitsAccountClient.Get(context.Context, string, string, string, *ResourceQuotaLimitsAccountClientGetOptions) (ResourceQuotaLimitsAccountClientGetResponse, error)`
- New function `*ResourceQuotaLimitsAccountClient.NewListPager(string, string, *ResourceQuotaLimitsAccountClientListOptions) *runtime.Pager[ResourceQuotaLimitsAccountClientListResponse]`
- New struct `ListReplicationsRequest`
- New struct `SnapshotPatch`
- New field `NextLink` in struct `BackupPoliciesList`
- New field `NextLink` in struct `ListReplications`
- New field `MirrorState`, `ReplicationCreationTime`, `ReplicationDeletionTime` in struct `Replication`
- New field `NextLink` in struct `SnapshotPoliciesList`
- New field `NextLink` in struct `SnapshotPolicyVolumeList`
- New field `NextLink` in struct `SnapshotsList`
- New field `Usage` in struct `SubscriptionQuotaItemProperties`
- New field `SystemData` in struct `VolumeGroupDetails`
- New field `NextLink` in struct `VolumeGroupList`
- New field `NextLink` in struct `VolumeQuotaRulesList`
- New field `Body` in struct `VolumesClientListReplicationsOptions`
- New anonymous field `Volume` in struct `VolumesClientSplitCloneFromParentResponse`


## 8.0.0-beta.2 (2025-10-20)
### Breaking Changes

- Function `*SnapshotsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, string, string, string, any, *SnapshotsClientBeginUpdateOptions)` to `(context.Context, string, string, string, string, string, SnapshotPatch, *SnapshotsClientBeginUpdateOptions)`
- Type of `BackupStatus.RelationshipStatus` has been changed from `*RelationshipStatus` to `*VolumeBackupRelationshipStatus`
- Type of `BucketPatchProperties.ProvisioningState` has been changed from `*NetappProvisioningState` to `*ProvisioningState`
- Type of `BucketProperties.ProvisioningState` has been changed from `*NetappProvisioningState` to `*ProvisioningState`
- Type of `PoolPatchProperties.CustomThroughputMibps` has been changed from `*float32` to `*int32`
- Type of `PoolProperties.CustomThroughputMibps` has been changed from `*float32` to `*int32`
- Type of `ReplicationStatus.RelationshipStatus` has been changed from `*RelationshipStatus` to `*VolumeReplicationRelationshipStatus`
- Type of `RestoreStatus.RelationshipStatus` has been changed from `*RelationshipStatus` to `*VolumeRestoreRelationshipStatus`
- Enum `NetappProvisioningState` has been removed
- Enum `RelationshipStatus` has been removed
- Field `NextLink` of struct `ListQuotaReportResponse` has been removed

### Features Added

- New enum type `BucketPatchPermissions` with values `BucketPatchPermissionsReadOnly`, `BucketPatchPermissionsReadWrite`
- New enum type `BucketPermissions` with values `BucketPermissionsReadOnly`, `BucketPermissionsReadWrite`
- New enum type `VolumeBackupRelationshipStatus` with values `VolumeBackupRelationshipStatusFailed`, `VolumeBackupRelationshipStatusIdle`, `VolumeBackupRelationshipStatusTransferring`, `VolumeBackupRelationshipStatusUnknown`
- New enum type `VolumeReplicationRelationshipStatus` with values `VolumeReplicationRelationshipStatusIdle`, `VolumeReplicationRelationshipStatusTransferring`
- New enum type `VolumeRestoreRelationshipStatus` with values `VolumeRestoreRelationshipStatusFailed`, `VolumeRestoreRelationshipStatusIdle`, `VolumeRestoreRelationshipStatusTransferring`, `VolumeRestoreRelationshipStatusUnknown`
- New struct `SnapshotPatch`
- New field `NextLink` in struct `BackupPoliciesList`
- New field `Permissions` in struct `BucketPatchProperties`
- New field `Permissions` in struct `BucketProperties`
- New field `NextLink` in struct `ListReplications`
- New field `MirrorState`, `RelationshipStatus` in struct `ReplicationObject`
- New field `NextLink` in struct `SnapshotPoliciesList`
- New field `NextLink` in struct `SnapshotPolicyVolumeList`
- New field `NextLink` in struct `SnapshotsList`
- New field `SystemData` in struct `VolumeGroupDetails`
- New field `NextLink` in struct `VolumeGroupList`
- New field `NextLink` in struct `VolumeQuotaRulesList`
- New anonymous field `Volume` in struct `VolumesClientSplitCloneFromParentResponse`


## 7.7.0 (2025-08-13)
### Features Added

- New value `ServiceLevelFlexible` added to enum type `ServiceLevel`
- New enum type `AcceptGrowCapacityPoolForShortTermCloneSplit` with values `AcceptGrowCapacityPoolForShortTermCloneSplitAccepted`, `AcceptGrowCapacityPoolForShortTermCloneSplitDeclined`
- New function `*VolumesClient.BeginSplitCloneFromParent(context.Context, string, string, string, string, *VolumesClientBeginSplitCloneFromParentOptions) (*runtime.Poller[VolumesClientSplitCloneFromParentResponse], error)`
- New field `CustomThroughputMibps` in struct `PoolPatchProperties`
- New field `CustomThroughputMibps` in struct `PoolProperties`
- New field `AcceptGrowCapacityPoolForShortTermCloneSplit`, `InheritedSizeInBytes` in struct `VolumeProperties`


## 7.6.0 (2025-07-25)
### Features Added

- New field `NextLink` in struct `SubscriptionQuotaItemList`


## 8.0.0-beta.1 (2025-05-23)
### Breaking Changes

- Struct `SubscriptionQuotaItem` has been removed
- Struct `SubscriptionQuotaItemList` has been removed
- Struct `SubscriptionQuotaItemProperties` has been removed
- Field `SubscriptionQuotaItem` of struct `ResourceQuotaLimitsClientGetResponse` has been removed
- Field `SubscriptionQuotaItemList` of struct `ResourceQuotaLimitsClientListResponse` has been removed

### Features Added

- New value `ServiceLevelFlexible` added to enum type `ServiceLevel`
- New enum type `AcceptGrowCapacityPoolForShortTermCloneSplit` with values `AcceptGrowCapacityPoolForShortTermCloneSplitAccepted`, `AcceptGrowCapacityPoolForShortTermCloneSplitDeclined`
- New enum type `CredentialsStatus` with values `CredentialsStatusActive`, `CredentialsStatusCredentialsExpired`, `CredentialsStatusNoCredentialsSet`
- New enum type `ExternalReplicationSetupStatus` with values `ExternalReplicationSetupStatusClusterPeerPending`, `ExternalReplicationSetupStatusClusterPeerRequired`, `ExternalReplicationSetupStatusNoActionRequired`, `ExternalReplicationSetupStatusReplicationCreateRequired`, `ExternalReplicationSetupStatusVServerPeerRequired`
- New enum type `LdapServerType` with values `LdapServerTypeActiveDirectory`, `LdapServerTypeOpenLDAP`
- New enum type `NetappProvisioningState` with values `NetappProvisioningStateAccepted`, `NetappProvisioningStateCanceled`, `NetappProvisioningStateDeleting`, `NetappProvisioningStateFailed`, `NetappProvisioningStateProvisioning`, `NetappProvisioningStateSucceeded`, `NetappProvisioningStateUpdating`
- New enum type `VolumeLanguage` with values `VolumeLanguageAr`, `VolumeLanguageArUTF8`, `VolumeLanguageC`, `VolumeLanguageCUTF8`, `VolumeLanguageCs`, `VolumeLanguageCsUTF8`, `VolumeLanguageDa`, `VolumeLanguageDaUTF8`, `VolumeLanguageDe`, `VolumeLanguageDeUTF8`, `VolumeLanguageEn`, `VolumeLanguageEnUTF8`, `VolumeLanguageEnUs`, `VolumeLanguageEnUsUTF8`, `VolumeLanguageEs`, `VolumeLanguageEsUTF8`, `VolumeLanguageFi`, `VolumeLanguageFiUTF8`, `VolumeLanguageFr`, `VolumeLanguageFrUTF8`, `VolumeLanguageHe`, `VolumeLanguageHeUTF8`, `VolumeLanguageHr`, `VolumeLanguageHrUTF8`, `VolumeLanguageHu`, `VolumeLanguageHuUTF8`, `VolumeLanguageIt`, `VolumeLanguageItUTF8`, `VolumeLanguageJa`, `VolumeLanguageJaJp932`, `VolumeLanguageJaJp932UTF8`, `VolumeLanguageJaJpPck`, `VolumeLanguageJaJpPckUTF8`, `VolumeLanguageJaJpPckV2`, `VolumeLanguageJaJpPckV2UTF8`, `VolumeLanguageJaUTF8`, `VolumeLanguageJaV1`, `VolumeLanguageJaV1UTF8`, `VolumeLanguageKo`, `VolumeLanguageKoUTF8`, `VolumeLanguageNl`, `VolumeLanguageNlUTF8`, `VolumeLanguageNo`, `VolumeLanguageNoUTF8`, `VolumeLanguagePl`, `VolumeLanguagePlUTF8`, `VolumeLanguagePt`, `VolumeLanguagePtUTF8`, `VolumeLanguageRo`, `VolumeLanguageRoUTF8`, `VolumeLanguageRu`, `VolumeLanguageRuUTF8`, `VolumeLanguageSk`, `VolumeLanguageSkUTF8`, `VolumeLanguageSl`, `VolumeLanguageSlUTF8`, `VolumeLanguageSv`, `VolumeLanguageSvUTF8`, `VolumeLanguageTr`, `VolumeLanguageTrUTF8`, `VolumeLanguageUTF8Mb4`, `VolumeLanguageZh`, `VolumeLanguageZhGbk`, `VolumeLanguageZhGbkUTF8`, `VolumeLanguageZhTw`, `VolumeLanguageZhTwBig5`, `VolumeLanguageZhTwBig5UTF8`, `VolumeLanguageZhTwUTF8`, `VolumeLanguageZhUTF8`
- New function `NewBucketsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BucketsClient, error)`
- New function `*BucketsClient.BeginCreateOrUpdate(context.Context, string, string, string, string, string, Bucket, *BucketsClientBeginCreateOrUpdateOptions) (*runtime.Poller[BucketsClientCreateOrUpdateResponse], error)`
- New function `*BucketsClient.BeginDelete(context.Context, string, string, string, string, string, *BucketsClientBeginDeleteOptions) (*runtime.Poller[BucketsClientDeleteResponse], error)`
- New function `*BucketsClient.GenerateCredentials(context.Context, string, string, string, string, string, BucketCredentialsExpiry, *BucketsClientGenerateCredentialsOptions) (BucketsClientGenerateCredentialsResponse, error)`
- New function `*BucketsClient.Get(context.Context, string, string, string, string, string, *BucketsClientGetOptions) (BucketsClientGetResponse, error)`
- New function `*BucketsClient.NewListPager(string, string, string, string, *BucketsClientListOptions) *runtime.Pager[BucketsClientListResponse]`
- New function `*BucketsClient.BeginUpdate(context.Context, string, string, string, string, string, BucketPatch, *BucketsClientBeginUpdateOptions) (*runtime.Poller[BucketsClientUpdateResponse], error)`
- New function `*ClientFactory.NewBucketsClient() *BucketsClient`
- New function `*ClientFactory.NewResourceQuotaLimitsAccountClient() *ResourceQuotaLimitsAccountClient`
- New function `NewResourceQuotaLimitsAccountClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ResourceQuotaLimitsAccountClient, error)`
- New function `*ResourceQuotaLimitsAccountClient.Get(context.Context, string, string, string, *ResourceQuotaLimitsAccountClientGetOptions) (ResourceQuotaLimitsAccountClientGetResponse, error)`
- New function `*ResourceQuotaLimitsAccountClient.NewListPager(string, string, *ResourceQuotaLimitsAccountClientListOptions) *runtime.Pager[ResourceQuotaLimitsAccountClientListResponse]`
- New function `*VolumesClient.BeginListQuotaReport(context.Context, string, string, string, string, *VolumesClientBeginListQuotaReportOptions) (*runtime.Poller[VolumesClientListQuotaReportResponse], error)`
- New function `*VolumesClient.BeginSplitCloneFromParent(context.Context, string, string, string, string, *VolumesClientBeginSplitCloneFromParentOptions) (*runtime.Poller[VolumesClientSplitCloneFromParentResponse], error)`
- New struct `Bucket`
- New struct `BucketCredentialsExpiry`
- New struct `BucketGenerateCredentials`
- New struct `BucketList`
- New struct `BucketPatch`
- New struct `BucketPatchProperties`
- New struct `BucketProperties`
- New struct `BucketServerPatchProperties`
- New struct `BucketServerProperties`
- New struct `CifsUser`
- New struct `FileSystemUser`
- New struct `LdapConfiguration`
- New struct `ListQuotaReportResponse`
- New struct `NfsUser`
- New struct `QuotaItem`
- New struct `QuotaItemList`
- New struct `QuotaItemProperties`
- New struct `QuotaReport`
- New field `LdapConfiguration` in struct `AccountProperties`
- New field `CustomThroughputMibps` in struct `PoolPatchProperties`
- New field `CustomThroughputMibps` in struct `PoolProperties`
- New field `ExternalReplicationSetupInfo`, `ExternalReplicationSetupStatus` in struct `ReplicationObject`
- New anonymous field `QuotaItem` in struct `ResourceQuotaLimitsClientGetResponse`
- New anonymous field `QuotaItemList` in struct `ResourceQuotaLimitsClientListResponse`
- New field `AcceptGrowCapacityPoolForShortTermCloneSplit`, `InheritedSizeInBytes`, `Language`, `LdapServerType` in struct `VolumeProperties`


## 7.5.0 (2025-04-25)
### Features Added

- New enum type `MultiAdStatus` with values `MultiAdStatusDisabled`, `MultiAdStatusEnabled`
- New enum type `ReplicationType` with values `ReplicationTypeCrossRegionReplication`, `ReplicationTypeCrossZoneReplication`
- New function `*ClientFactory.NewResourceUsagesClient() *ResourceUsagesClient`
- New function `NewResourceUsagesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ResourceUsagesClient, error)`
- New function `*ResourceUsagesClient.Get(context.Context, string, string, *ResourceUsagesClientGetOptions) (ResourceUsagesClientGetResponse, error)`
- New function `*ResourceUsagesClient.NewListPager(string, *ResourceUsagesClientListOptions) *runtime.Pager[ResourceUsagesClientListResponse]`
- New struct `DestinationReplication`
- New struct `UsageName`
- New struct `UsageProperties`
- New struct `UsageResult`
- New struct `UsagesListResult`
- New field `MultiAdStatus`, `NfsV4IDDomain` in struct `AccountProperties`
- New field `CompletionDate`, `IsLargeVolume`, `SnapshotCreationDate` in struct `BackupProperties`
- New field `FederatedClientID` in struct `EncryptionIdentity`
- New field `NextLink` in struct `OperationListResult`
- New field `DestinationReplications` in struct `ReplicationObject`


## 7.4.0 (2025-02-12)
### Features Added

- New enum type `CoolAccessTieringPolicy` with values `CoolAccessTieringPolicyAuto`, `CoolAccessTieringPolicySnapshotOnly`
- New function `*AccountsClient.BeginChangeKeyVault(context.Context, string, string, *AccountsClientBeginChangeKeyVaultOptions) (*runtime.Poller[AccountsClientChangeKeyVaultResponse], error)`
- New function `*AccountsClient.BeginGetChangeKeyVaultInformation(context.Context, string, string, *AccountsClientBeginGetChangeKeyVaultInformationOptions) (*runtime.Poller[AccountsClientGetChangeKeyVaultInformationResponse], error)`
- New function `*AccountsClient.BeginTransitionToCmk(context.Context, string, string, *AccountsClientBeginTransitionToCmkOptions) (*runtime.Poller[AccountsClientTransitionToCmkResponse], error)`
- New struct `ChangeKeyVault`
- New struct `EncryptionTransitionRequest`
- New struct `GetKeyVaultStatusResponse`
- New struct `GetKeyVaultStatusResponseProperties`
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
