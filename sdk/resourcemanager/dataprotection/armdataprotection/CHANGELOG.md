# Release History

## 2.0.0-beta.1 (2023-01-13)
### Breaking Changes

- Type of `ResourceGuard.ProvisioningState` has been changed from `*ProvisioningState` to `*ResourceGuardProvisioningState`

### Features Added

- New value `StorageSettingTypesZoneRedundant` added to type alias `StorageSettingTypes`
- New type alias `ExistingResourcePolicy` with values `ExistingResourcePolicyPatch`, `ExistingResourcePolicySkip`
- New type alias `ImmutabilityState` with values `ImmutabilityStateDisabled`, `ImmutabilityStateLocked`, `ImmutabilityStateUnlocked`
- New type alias `PersistentVolumeRestoreMode` with values `PersistentVolumeRestoreModeRestoreWithVolumeData`, `PersistentVolumeRestoreModeRestoreWithoutVolumeData`
- New type alias `ResourceGuardProvisioningState` with values `ResourceGuardProvisioningStateFailed`, `ResourceGuardProvisioningStateProvisioning`, `ResourceGuardProvisioningStateSucceeded`, `ResourceGuardProvisioningStateUnknown`, `ResourceGuardProvisioningStateUpdating`
- New type alias `SoftDeleteState` with values `SoftDeleteStateAlwaysOn`, `SoftDeleteStateOff`, `SoftDeleteStateOn`
- New function `*BackupDatasourceParameters.GetBackupDatasourceParameters() *BackupDatasourceParameters`
- New function `NewBackupInstancesExtensionRoutingClient(azcore.TokenCredential, *arm.ClientOptions) (*BackupInstancesExtensionRoutingClient, error)`
- New function `*BackupInstancesExtensionRoutingClient.NewListPager(string, *BackupInstancesExtensionRoutingClientListOptions) *runtime.Pager[BackupInstancesExtensionRoutingClientListResponse]`
- New function `*BlobBackupDatasourceParameters.GetBackupDatasourceParameters() *BackupDatasourceParameters`
- New function `NewDeletedBackupInstancesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DeletedBackupInstancesClient, error)`
- New function `*DeletedBackupInstancesClient.Get(context.Context, string, string, string, *DeletedBackupInstancesClientGetOptions) (DeletedBackupInstancesClientGetResponse, error)`
- New function `*DeletedBackupInstancesClient.NewListPager(string, string, *DeletedBackupInstancesClientListOptions) *runtime.Pager[DeletedBackupInstancesClientListResponse]`
- New function `*DeletedBackupInstancesClient.BeginUndelete(context.Context, string, string, string, *DeletedBackupInstancesClientBeginUndeleteOptions) (*runtime.Poller[DeletedBackupInstancesClientUndeleteResponse], error)`
- New function `NewDppResourceGuardProxyClient(string, azcore.TokenCredential, *arm.ClientOptions) (*DppResourceGuardProxyClient, error)`
- New function `*DppResourceGuardProxyClient.Delete(context.Context, string, string, string, *DppResourceGuardProxyClientDeleteOptions) (DppResourceGuardProxyClientDeleteResponse, error)`
- New function `*DppResourceGuardProxyClient.Get(context.Context, string, string, string, *DppResourceGuardProxyClientGetOptions) (DppResourceGuardProxyClientGetResponse, error)`
- New function `*DppResourceGuardProxyClient.NewListPager(string, string, *DppResourceGuardProxyClientListOptions) *runtime.Pager[DppResourceGuardProxyClientListResponse]`
- New function `*DppResourceGuardProxyClient.Put(context.Context, string, string, string, ResourceGuardProxyBaseResource, *DppResourceGuardProxyClientPutOptions) (DppResourceGuardProxyClientPutResponse, error)`
- New function `*DppResourceGuardProxyClient.UnlockDelete(context.Context, string, string, string, UnlockDeleteRequest, *DppResourceGuardProxyClientUnlockDeleteOptions) (DppResourceGuardProxyClientUnlockDeleteResponse, error)`
- New function `*ItemPathBasedRestoreCriteria.GetItemLevelRestoreCriteria() *ItemLevelRestoreCriteria`
- New function `*KubernetesClusterBackupDatasourceParameters.GetBackupDatasourceParameters() *BackupDatasourceParameters`
- New function `*KubernetesClusterRestoreCriteria.GetItemLevelRestoreCriteria() *ItemLevelRestoreCriteria`
- New struct `BackupInstancesExtensionRoutingClient`
- New struct `BackupInstancesExtensionRoutingClientListResponse`
- New struct `BlobBackupDatasourceParameters`
- New struct `DeletedBackupInstance`
- New struct `DeletedBackupInstanceResource`
- New struct `DeletedBackupInstanceResourceList`
- New struct `DeletedBackupInstancesClient`
- New struct `DeletedBackupInstancesClientListResponse`
- New struct `DeletedBackupInstancesClientUndeleteResponse`
- New struct `DeletionInfo`
- New struct `DppProxyResource`
- New struct `DppResourceGuardProxyClient`
- New struct `DppResourceGuardProxyClientListResponse`
- New struct `ImmutabilitySettings`
- New struct `ItemPathBasedRestoreCriteria`
- New struct `KubernetesClusterBackupDatasourceParameters`
- New struct `KubernetesClusterRestoreCriteria`
- New struct `ResourceGuardOperationDetail`
- New struct `ResourceGuardProxyBase`
- New struct `ResourceGuardProxyBaseResource`
- New struct `ResourceGuardProxyBaseResourceList`
- New struct `SecuritySettings`
- New struct `SoftDeleteSettings`
- New struct `UnlockDeleteRequest`
- New struct `UnlockDeleteResponse`
- New field `Tags` in struct `BackupInstanceResource`
- New field `IsVaultProtectedByResourceGuard` in struct `BackupVault`
- New field `SecuritySettings` in struct `BackupVault`
- New field `SecuritySettings` in struct `PatchBackupVaultInput`
- New field `BackupDatasourceParametersList` in struct `PolicyParameters`


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).