# Release History

## 4.0.0 (2025-10-13)
### Breaking Changes

- Struct `AzureBackupFindRestorableTimeRangesRequestResource` has been removed
- Struct `ClientDiscoveryDisplay` has been removed
- Struct `ClientDiscoveryForLogSpecification` has been removed
- Struct `ClientDiscoveryForProperties` has been removed
- Struct `ClientDiscoveryForServiceSpecification` has been removed
- Struct `ClientDiscoveryResponse` has been removed
- Struct `ClientDiscoveryValueForSingleAPI` has been removed
- Struct `DppBaseTrackedResource` has been removed
- Struct `DppProxyResource` has been removed
- Struct `DppResource` has been removed
- Struct `DppResourceList` has been removed
- Struct `DppTrackedResource` has been removed
- Struct `DppTrackedResourceList` has been removed
- Struct `DppWorkerRequest` has been removed
- Struct `ErrorDetail` has been removed
- Struct `ErrorResponse` has been removed
- Struct `RecoveryPointsFilters` has been removed
- Field `ClientDiscoveryResponse` of struct `OperationsClientListResponse` has been removed

### Features Added

- New enum type `AKSVolumeTypes` with values `AKSVolumeTypesAzureDisk`, `AKSVolumeTypesAzureFileShareSMB`
- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New function `*AdlsBlobBackupDatasourceParameters.GetBackupDatasourceParameters() *BackupDatasourceParameters`
- New function `*AdlsBlobBackupDatasourceParameters.GetBlobBackupDatasourceParameters() *BlobBackupDatasourceParameters`
- New function `*BackupInstancesClient.BeginValidateForModifyBackup(context.Context, string, string, string, ValidateForModifyBackupRequest, *BackupInstancesClientBeginValidateForModifyBackupOptions) (*runtime.Poller[BackupInstancesClientValidateForModifyBackupResponse], error)`
- New function `*BlobBackupDatasourceParameters.GetBlobBackupDatasourceParameters() *BlobBackupDatasourceParameters`
- New struct `AdlsBlobBackupDatasourceParameters`
- New struct `Operation`
- New struct `OperationDisplay`
- New struct `OperationListResult`
- New struct `ValidateForModifyBackupRequest`
- New field `SystemData` in struct `DppBaseResource`
- New field `RenameTo` in struct `ItemPathBasedRestoreCriteria`
- New field `IncludedVolumeTypes` in struct `KubernetesClusterBackupDatasourceParameters`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`


## 3.1.0 (2024-07-26)
### Features Added

- New enum type `BCDRSecurityLevel` with values `BCDRSecurityLevelExcellent`, `BCDRSecurityLevelFair`, `BCDRSecurityLevelGood`, `BCDRSecurityLevelNotSupported`, `BCDRSecurityLevelPoor`
- New enum type `EncryptionState` with values `EncryptionStateDisabled`, `EncryptionStateEnabled`, `EncryptionStateInconsistent`
- New enum type `IdentityType` with values `IdentityTypeSystemAssigned`, `IdentityTypeUserAssigned`
- New enum type `InfrastructureEncryptionState` with values `InfrastructureEncryptionStateDisabled`, `InfrastructureEncryptionStateEnabled`
- New function `NewBackupInstancesExtensionRoutingClient(azcore.TokenCredential, *arm.ClientOptions) (*BackupInstancesExtensionRoutingClient, error)`
- New function `*BackupInstancesExtensionRoutingClient.NewListPager(string, *BackupInstancesExtensionRoutingClientListOptions) *runtime.Pager[BackupInstancesExtensionRoutingClientListResponse]`
- New function `*ClientFactory.NewBackupInstancesExtensionRoutingClient() *BackupInstancesExtensionRoutingClient`
- New struct `CmkKekIdentity`
- New struct `CmkKeyVaultProperties`
- New struct `EncryptionSettings`
- New struct `StopProtectionRequest`
- New struct `SuspendBackupRequest`
- New field `ResourceGuardOperationRequests` in struct `AzureBackupRecoveryPointBasedRestoreRequest`
- New field `ResourceGuardOperationRequests` in struct `AzureBackupRecoveryTimeBasedRestoreRequest`
- New field `ResourceGuardOperationRequests` in struct `AzureBackupRestoreRequest`
- New field `ResourceGuardOperationRequests` in struct `AzureBackupRestoreWithRehydrationRequest`
- New field `ResourceGuardOperationRequests` in struct `BackupInstance`
- New field `XMSAuthorizationAuxiliary` in struct `BackupInstancesClientBeginCreateOrUpdateOptions`
- New field `XMSAuthorizationAuxiliary` in struct `BackupInstancesClientBeginDeleteOptions`
- New field `Parameters`, `XMSAuthorizationAuxiliary` in struct `BackupInstancesClientBeginStopProtectionOptions`
- New field `Parameters`, `XMSAuthorizationAuxiliary` in struct `BackupInstancesClientBeginSuspendBackupsOptions`
- New field `XMSAuthorizationAuxiliary` in struct `BackupInstancesClientBeginTriggerRestoreOptions`
- New field `BcdrSecurityLevel`, `ResourceGuardOperationRequests` in struct `BackupVault`
- New field `XMSAuthorizationAuxiliary` in struct `BackupVaultsClientBeginCreateOrUpdateOptions`
- New field `XMSAuthorizationAuxiliary` in struct `BackupVaultsClientBeginUpdateOptions`
- New field `ResourceGuardOperationRequests` in struct `DeletedBackupInstance`
- New field `XMSAuthorizationAuxiliary` in struct `DppResourceGuardProxyClientUnlockDeleteOptions`
- New field `ResourceModifierReference` in struct `KubernetesClusterRestoreCriteria`
- New field `ResourceModifierReference` in struct `KubernetesClusterVaultTierRestoreCriteria`
- New field `ResourceGuardOperationRequests` in struct `PatchBackupVaultInput`
- New field `EncryptionSettings` in struct `SecuritySettings`


## 3.0.0 (2023-12-22)
### Breaking Changes

- Type of `BaseResourceProperties.ObjectType` has been changed from `*string` to `*ResourcePropertiesObjectType`

### Features Added

- New enum type `RecoveryPointCompletionState` with values `RecoveryPointCompletionStateCompleted`, `RecoveryPointCompletionStatePartial`
- New enum type `ResourcePropertiesObjectType` with values `ResourcePropertiesObjectTypeDefaultResourceProperties`
- New function `*BackupInstancesClient.BeginTriggerCrossRegionRestore(context.Context, string, string, CrossRegionRestoreRequestObject, *BackupInstancesClientBeginTriggerCrossRegionRestoreOptions) (*runtime.Poller[BackupInstancesClientTriggerCrossRegionRestoreResponse], error)`
- New function `*BackupInstancesClient.BeginValidateCrossRegionRestore(context.Context, string, string, ValidateCrossRegionRestoreRequestObject, *BackupInstancesClientBeginValidateCrossRegionRestoreOptions) (*runtime.Poller[BackupInstancesClientValidateCrossRegionRestoreResponse], error)`
- New function `*ClientFactory.NewFetchCrossRegionRestoreJobClient() *FetchCrossRegionRestoreJobClient`
- New function `*ClientFactory.NewFetchCrossRegionRestoreJobsClient() *FetchCrossRegionRestoreJobsClient`
- New function `*ClientFactory.NewFetchSecondaryRecoveryPointsClient() *FetchSecondaryRecoveryPointsClient`
- New function `*DefaultResourceProperties.GetBaseResourceProperties() *BaseResourceProperties`
- New function `NewFetchCrossRegionRestoreJobClient(string, azcore.TokenCredential, *arm.ClientOptions) (*FetchCrossRegionRestoreJobClient, error)`
- New function `*FetchCrossRegionRestoreJobClient.Get(context.Context, string, string, CrossRegionRestoreJobRequest, *FetchCrossRegionRestoreJobClientGetOptions) (FetchCrossRegionRestoreJobClientGetResponse, error)`
- New function `NewFetchCrossRegionRestoreJobsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*FetchCrossRegionRestoreJobsClient, error)`
- New function `*FetchCrossRegionRestoreJobsClient.NewListPager(string, string, CrossRegionRestoreJobsRequest, *FetchCrossRegionRestoreJobsClientListOptions) *runtime.Pager[FetchCrossRegionRestoreJobsClientListResponse]`
- New function `NewFetchSecondaryRecoveryPointsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*FetchSecondaryRecoveryPointsClient, error)`
- New function `*FetchSecondaryRecoveryPointsClient.NewListPager(string, string, FetchSecondaryRPsRequestParameters, *FetchSecondaryRecoveryPointsClientListOptions) *runtime.Pager[FetchSecondaryRecoveryPointsClientListResponse]`
- New function `*KubernetesClusterVaultTierRestoreCriteria.GetItemLevelRestoreCriteria() *ItemLevelRestoreCriteria`
- New struct `CrossRegionRestoreDetails`
- New struct `CrossRegionRestoreJobRequest`
- New struct `CrossRegionRestoreJobsRequest`
- New struct `CrossRegionRestoreRequestObject`
- New struct `DefaultResourceProperties`
- New struct `ErrorDetail`
- New struct `ErrorResponse`
- New struct `FetchSecondaryRPsRequestParameters`
- New struct `KubernetesClusterVaultTierRestoreCriteria`
- New struct `UserFacingWarningDetail`
- New struct `ValidateCrossRegionRestoreRequestObject`
- New field `RecoveryPointState` in struct `AzureBackupDiscreteRecoveryPoint`
- New field `ReplicatedRegions` in struct `BackupVault`
- New field `WarningDetails` in struct `JobExtendedInfo`


## 2.4.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.3.0 (2023-07-28)
### Features Added

- New enum type `CrossRegionRestoreState` with values `CrossRegionRestoreStateDisabled`, `CrossRegionRestoreStateEnabled`
- New enum type `SecureScoreLevel` with values `SecureScoreLevelAdequate`, `SecureScoreLevelMaximum`, `SecureScoreLevelMinimum`, `SecureScoreLevelNone`, `SecureScoreLevelNotSupported`
- New function `*BaseResourceProperties.GetBaseResourceProperties() *BaseResourceProperties`
- New struct `CrossRegionRestoreSettings`
- New struct `IdentityDetails`
- New struct `NamespacedNameResource`
- New struct `UserAssignedIdentity`
- New field `RehydrationPriority` in struct `AzureBackupJob`
- New field `IdentityDetails` in struct `AzureBackupRecoveryPointBasedRestoreRequest`
- New field `IdentityDetails` in struct `AzureBackupRecoveryTimeBasedRestoreRequest`
- New field `IdentityDetails` in struct `AzureBackupRestoreRequest`
- New field `IdentityDetails` in struct `AzureBackupRestoreWithRehydrationRequest`
- New field `IdentityDetails` in struct `BackupInstance`
- New field `SecureScore` in struct `BackupVault`
- New field `ResourceProperties` in struct `Datasource`
- New field `ResourceProperties` in struct `DatasourceSet`
- New field `IdentityDetails` in struct `DeletedBackupInstance`
- New field `UserAssignedIdentities` in struct `DppIdentityDetails`
- New field `CrossRegionRestoreSettings` in struct `FeatureSettings`
- New field `BackupHookReferences` in struct `KubernetesClusterBackupDatasourceParameters`
- New field `RestoreHookReferences` in struct `KubernetesClusterRestoreCriteria`


## 2.2.0 (2023-06-23)
### Features Added

- New function `*ClientFactory.NewDppResourceGuardProxyClient() *DppResourceGuardProxyClient`
- New function `NewDppResourceGuardProxyClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DppResourceGuardProxyClient, error)`
- New function `*DppResourceGuardProxyClient.CreateOrUpdate(context.Context, string, string, string, ResourceGuardProxyBaseResource, *DppResourceGuardProxyClientCreateOrUpdateOptions) (DppResourceGuardProxyClientCreateOrUpdateResponse, error)`
- New function `*DppResourceGuardProxyClient.Delete(context.Context, string, string, string, *DppResourceGuardProxyClientDeleteOptions) (DppResourceGuardProxyClientDeleteResponse, error)`
- New function `*DppResourceGuardProxyClient.Get(context.Context, string, string, string, *DppResourceGuardProxyClientGetOptions) (DppResourceGuardProxyClientGetResponse, error)`
- New function `*DppResourceGuardProxyClient.NewListPager(string, string, *DppResourceGuardProxyClientListOptions) *runtime.Pager[DppResourceGuardProxyClientListResponse]`
- New function `*DppResourceGuardProxyClient.UnlockDelete(context.Context, string, string, string, UnlockDeleteRequest, *DppResourceGuardProxyClientUnlockDeleteOptions) (DppResourceGuardProxyClientUnlockDeleteResponse, error)`
- New struct `ResourceGuardOperationDetail`
- New struct `ResourceGuardProxyBase`
- New struct `ResourceGuardProxyBaseResource`
- New struct `ResourceGuardProxyBaseResourceList`
- New struct `UnlockDeleteRequest`
- New struct `UnlockDeleteResponse`


## 2.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 2.1.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 2.0.0 (2023-02-24)
### Breaking Changes

- Function `*ResourceGuardsClient.Patch` parameter(s) have been changed from `(context.Context, string, string, PatchResourceRequestInput, *ResourceGuardsClientPatchOptions)` to `(context.Context, string, string, PatchResourceGuardInput, *ResourceGuardsClientPatchOptions)`
- Const `StorageSettingStoreTypesSnapshotStore` from type alias `StorageSettingStoreTypes` has been removed
- Operation `*BackupVaultsClient.Delete` has been changed to LRO, use `*BackupVaultsClient.BeginDelete` instead.
- Field `Identity` of struct `ResourceGuardResource` has been removed

### Features Added

- New value `SourceDataStoreTypeOperationalStore` added to type alias `SourceDataStoreType`
- New value `StorageSettingStoreTypesOperationalStore` added to type alias `StorageSettingStoreTypes`
- New value `StorageSettingTypesZoneRedundant` added to type alias `StorageSettingTypes`
- New type alias `CrossSubscriptionRestoreState` with values `CrossSubscriptionRestoreStateDisabled`, `CrossSubscriptionRestoreStateEnabled`, `CrossSubscriptionRestoreStatePermanentlyDisabled`
- New type alias `ExistingResourcePolicy` with values `ExistingResourcePolicyPatch`, `ExistingResourcePolicySkip`
- New type alias `ImmutabilityState` with values `ImmutabilityStateDisabled`, `ImmutabilityStateLocked`, `ImmutabilityStateUnlocked`
- New type alias `PersistentVolumeRestoreMode` with values `PersistentVolumeRestoreModeRestoreWithVolumeData`, `PersistentVolumeRestoreModeRestoreWithoutVolumeData`
- New type alias `SoftDeleteState` with values `SoftDeleteStateAlwaysOn`, `SoftDeleteStateOff`, `SoftDeleteStateOn`
- New function `*BackupDatasourceParameters.GetBackupDatasourceParameters() *BackupDatasourceParameters`
- New function `*BlobBackupDatasourceParameters.GetBackupDatasourceParameters() *BackupDatasourceParameters`
- New function `NewDeletedBackupInstancesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DeletedBackupInstancesClient, error)`
- New function `*DeletedBackupInstancesClient.Get(context.Context, string, string, string, *DeletedBackupInstancesClientGetOptions) (DeletedBackupInstancesClientGetResponse, error)`
- New function `*DeletedBackupInstancesClient.NewListPager(string, string, *DeletedBackupInstancesClientListOptions) *runtime.Pager[DeletedBackupInstancesClientListResponse]`
- New function `*DeletedBackupInstancesClient.BeginUndelete(context.Context, string, string, string, *DeletedBackupInstancesClientBeginUndeleteOptions) (*runtime.Poller[DeletedBackupInstancesClientUndeleteResponse], error)`
- New function `*ItemPathBasedRestoreCriteria.GetItemLevelRestoreCriteria() *ItemLevelRestoreCriteria`
- New function `*KubernetesClusterBackupDatasourceParameters.GetBackupDatasourceParameters() *BackupDatasourceParameters`
- New function `*KubernetesClusterRestoreCriteria.GetItemLevelRestoreCriteria() *ItemLevelRestoreCriteria`
- New struct `BlobBackupDatasourceParameters`
- New struct `CrossSubscriptionRestoreSettings`
- New struct `DeletedBackupInstance`
- New struct `DeletedBackupInstanceResource`
- New struct `DeletedBackupInstanceResourceList`
- New struct `DeletedBackupInstancesClient`
- New struct `DeletedBackupInstancesClientListResponse`
- New struct `DeletedBackupInstancesClientUndeleteResponse`
- New struct `DeletionInfo`
- New struct `DppBaseTrackedResource`
- New struct `DppProxyResource`
- New struct `FeatureSettings`
- New struct `ImmutabilitySettings`
- New struct `ItemPathBasedRestoreCriteria`
- New struct `KubernetesClusterBackupDatasourceParameters`
- New struct `KubernetesClusterRestoreCriteria`
- New struct `PatchResourceGuardInput`
- New struct `SecuritySettings`
- New struct `SoftDeleteSettings`
- New field `ExpiryTime` in struct `AzureBackupDiscreteRecoveryPoint`
- New field `Tags` in struct `BackupInstanceResource`
- New field `FeatureSettings` in struct `BackupVault`
- New field `IsVaultProtectedByResourceGuard` in struct `BackupVault`
- New field `SecuritySettings` in struct `BackupVault`
- New field `FeatureSettings` in struct `PatchBackupVaultInput`
- New field `SecuritySettings` in struct `PatchBackupVaultInput`
- New field `BackupDatasourceParametersList` in struct `PolicyParameters`
- New field `TargetResourceArmID` in struct `TargetDetails`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).