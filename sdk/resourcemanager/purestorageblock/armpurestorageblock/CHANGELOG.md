# Release History

## 1.1.0-beta.2 (2026-08-21)
### Features Added

- New enum type `PlatformConsoleAuthType` with values `PlatformConsoleAuthTypeSSH`
- New enum type `PlatformConsoleRole` with values `PlatformConsoleRoleArrayAdmin`, `PlatformConsoleRoleReadOnly`, `PlatformConsoleRoleStorageAdmin`
- New enum type `VolumeGroupSourceType` with values `VolumeGroupSourceTypeNone`, `VolumeGroupSourceTypeRecoverableVolumeGroup`, `VolumeGroupSourceTypeSnapshot`, `VolumeGroupSourceTypeVolumeGroup`
- New enum type `VolumeSourceType` with values `VolumeSourceTypeNone`, `VolumeSourceTypeRecoverableVolume`, `VolumeSourceTypeSerialNumber`, `VolumeSourceTypeSnapshot`, `VolumeSourceTypeVolume`
- New function `*ClientFactory.NewRecoverableVolumeGroupsClient() *RecoverableVolumeGroupsClient`
- New function `*ClientFactory.NewSaaSOperationGroupClient() *SaaSOperationGroupClient`
- New function `*ClientFactory.NewVolumeGroupSnapshotsClient() *VolumeGroupSnapshotsClient`
- New function `*PlatformConsoleAuthConfig.GetPlatformConsoleAuthConfig() *PlatformConsoleAuthConfig`
- New function `*PlatformConsoleAuthResult.GetPlatformConsoleAuthResult() *PlatformConsoleAuthResult`
- New function `NewRecoverableVolumeGroupsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*RecoverableVolumeGroupsClient, error)`
- New function `*RecoverableVolumeGroupsClient.BeginDelete(ctx context.Context, resourceGroupName string, storagePoolName string, recoverableVolumeGroupName string, options *RecoverableVolumeGroupsClientBeginDeleteOptions) (*runtime.Poller[RecoverableVolumeGroupsClientDeleteResponse], error)`
- New function `*RecoverableVolumeGroupsClient.Get(ctx context.Context, resourceGroupName string, storagePoolName string, recoverableVolumeGroupName string, options *RecoverableVolumeGroupsClientGetOptions) (RecoverableVolumeGroupsClientGetResponse, error)`
- New function `*RecoverableVolumeGroupsClient.NewListByStoragePoolPager(resourceGroupName string, storagePoolName string, options *RecoverableVolumeGroupsClientListByStoragePoolOptions) *runtime.Pager[RecoverableVolumeGroupsClientListByStoragePoolResponse]`
- New function `*ReservationsClient.LatestLinkedSaaS(ctx context.Context, resourceGroupName string, reservationName string, options *ReservationsClientLatestLinkedSaaSOptions) (ReservationsClientLatestLinkedSaaSResponse, error)`
- New function `*ReservationsClient.BeginLinkSaaS(ctx context.Context, resourceGroupName string, reservationName string, body LinkSaaSRequest, options *ReservationsClientBeginLinkSaaSOptions) (*runtime.Poller[ReservationsClientLinkSaaSResponse], error)`
- New function `*SSHPlatformConsoleAuthConfig.GetPlatformConsoleAuthConfig() *PlatformConsoleAuthConfig`
- New function `*SSHPlatformConsoleAuthResult.GetPlatformConsoleAuthResult() *PlatformConsoleAuthResult`
- New function `NewSaaSOperationGroupClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SaaSOperationGroupClient, error)`
- New function `*SaaSOperationGroupClient.BeginActivateResource(ctx context.Context, body ActivateSaaSRequest, options *SaaSOperationGroupClientBeginActivateResourceOptions) (*runtime.Poller[SaaSOperationGroupClientActivateResourceResponse], error)`
- New function `*StoragePoolsClient.ConfigurePlatformConsoleAuth(ctx context.Context, resourceGroupName string, storagePoolName string, config PlatformConsoleAuthConfigClassification, options *StoragePoolsClientConfigurePlatformConsoleAuthOptions) (StoragePoolsClientConfigurePlatformConsoleAuthResponse, error)`
- New function `*StoragePoolsClient.ListPlatformConsoleActivationCode(ctx context.Context, resourceGroupName string, storagePoolName string, options *StoragePoolsClientListPlatformConsoleActivationCodeOptions) (StoragePoolsClientListPlatformConsoleActivationCodeResponse, error)`
- New function `NewVolumeGroupSnapshotsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VolumeGroupSnapshotsClient, error)`
- New function `*VolumeGroupSnapshotsClient.BeginCreate(ctx context.Context, resourceGroupName string, storagePoolName string, volumeGroupName string, snapshotName string, resource VolumeGroupSnapshot, options *VolumeGroupSnapshotsClientBeginCreateOptions) (*runtime.Poller[VolumeGroupSnapshotsClientCreateResponse], error)`
- New function `*VolumeGroupSnapshotsClient.BeginDelete(ctx context.Context, resourceGroupName string, storagePoolName string, volumeGroupName string, snapshotName string, options *VolumeGroupSnapshotsClientBeginDeleteOptions) (*runtime.Poller[VolumeGroupSnapshotsClientDeleteResponse], error)`
- New function `*VolumeGroupSnapshotsClient.Get(ctx context.Context, resourceGroupName string, storagePoolName string, volumeGroupName string, snapshotName string, options *VolumeGroupSnapshotsClientGetOptions) (VolumeGroupSnapshotsClientGetResponse, error)`
- New function `*VolumeGroupSnapshotsClient.NewListByVolumeGroupPager(resourceGroupName string, storagePoolName string, volumeGroupName string, options *VolumeGroupSnapshotsClientListByVolumeGroupOptions) *runtime.Pager[VolumeGroupSnapshotsClientListByVolumeGroupResponse]`
- New function `*VolumeGroupSnapshotsClient.ListSnapshots(ctx context.Context, resourceGroupName string, storagePoolName string, volumeGroupName string, properties VolumeGroupSnapshotListRequest, options *VolumeGroupSnapshotsClientListSnapshotsOptions) (VolumeGroupSnapshotsClientListSnapshotsResponse, error)`
- New function `*VolumeGroupsClient.BeginOverwrite(ctx context.Context, resourceGroupName string, storagePoolName string, volumeGroupName string, body VolumeGroupOverwriteRequest, options *VolumeGroupsClientBeginOverwriteOptions) (*runtime.Poller[VolumeGroupsClientOverwriteResponse], error)`
- New function `*VolumesClient.BeginOverwrite(ctx context.Context, resourceGroupName string, storagePoolName string, volumeGroupName string, volumeName string, body VolumeOverwriteRequest, options *VolumesClientBeginOverwriteOptions) (*runtime.Poller[VolumesClientOverwriteResponse], error)`
- New struct `ActivateSaaSRequest`
- New struct `DestroyedStateProperties`
- New struct `LatestLinkedSaaSResponse`
- New struct `LinkSaaSRequest`
- New struct `PlatformConsoleAccessSettings`
- New struct `PlatformConsoleActivationCode`
- New struct `PlatformConsoleSettings`
- New struct `PlatformConsoleSubnet`
- New struct `RecoverableVolumeGroup`
- New struct `RecoverableVolumeGroupListResult`
- New struct `RecoverableVolumeGroupProperties`
- New struct `SSHPlatformConsoleAuthConfig`
- New struct `SSHPlatformConsoleAuthResult`
- New struct `SaaSResourceDetailsResponse`
- New struct `VolumeGroupOverwriteRequest`
- New struct `VolumeGroupSnapshot`
- New struct `VolumeGroupSnapshotListRequest`
- New struct `VolumeGroupSnapshotListResult`
- New struct `VolumeGroupSnapshotPostListResult`
- New struct `VolumeGroupSnapshotProperties`
- New struct `VolumeOverwriteRequest`
- New struct `VolumeSnapshotInfo`
- New struct `VolumeSnapshotSource`
- New field `SoftDeletion`, `SourceRecoverableVolumeResourceID`, `SourceSerialNumber`, `SourceType`, `SourceVolumeSnapshot` in struct `AzureVolumeProperties`
- New field `SaaSResourceID` in struct `MarketplaceDetails`
- New field `PlatformConsoleSettings` in struct `StoragePoolProperties`
- New field `PlatformConsoleSettings` in struct `StoragePoolUpdateProperties`
- New field `SourceRecoverableVolumeGroupResourceID`, `SourceSnapshotResourceID`, `SourceType` in struct `VolumeGroupProperties`


## 1.1.0-beta.1 (2026-05-14)
### Features Added

- New function `*ClientFactory.NewVolumeGroupsClient() *VolumeGroupsClient`
- New function `*ClientFactory.NewVolumesClient() *VolumesClient`
- New function `NewVolumeGroupsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VolumeGroupsClient, error)`
- New function `*VolumeGroupsClient.BeginCreate(ctx context.Context, resourceGroupName string, storagePoolName string, volumeGroupName string, resource VolumeGroup, options *VolumeGroupsClientBeginCreateOptions) (*runtime.Poller[VolumeGroupsClientCreateResponse], error)`
- New function `*VolumeGroupsClient.BeginDelete(ctx context.Context, resourceGroupName string, storagePoolName string, volumeGroupName string, options *VolumeGroupsClientBeginDeleteOptions) (*runtime.Poller[VolumeGroupsClientDeleteResponse], error)`
- New function `*VolumeGroupsClient.Get(ctx context.Context, resourceGroupName string, storagePoolName string, volumeGroupName string, options *VolumeGroupsClientGetOptions) (VolumeGroupsClientGetResponse, error)`
- New function `*VolumeGroupsClient.GetStatus(ctx context.Context, resourceGroupName string, storagePoolName string, volumeGroupName string, options *VolumeGroupsClientGetStatusOptions) (VolumeGroupsClientGetStatusResponse, error)`
- New function `*VolumeGroupsClient.NewListByStoragePoolPager(resourceGroupName string, storagePoolName string, options *VolumeGroupsClientListByStoragePoolOptions) *runtime.Pager[VolumeGroupsClientListByStoragePoolResponse]`
- New function `*VolumeGroupsClient.ListConnectionParameters(ctx context.Context, resourceGroupName string, storagePoolName string, volumeGroupName string, options *VolumeGroupsClientListConnectionParametersOptions) (VolumeGroupsClientListConnectionParametersResponse, error)`
- New function `*VolumeGroupsClient.BeginUpdate(ctx context.Context, resourceGroupName string, storagePoolName string, volumeGroupName string, properties VolumeGroupUpdate, options *VolumeGroupsClientBeginUpdateOptions) (*runtime.Poller[VolumeGroupsClientUpdateResponse], error)`
- New function `NewVolumesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VolumesClient, error)`
- New function `*VolumesClient.BeginCreate(ctx context.Context, resourceGroupName string, storagePoolName string, volumeGroupName string, volumeName string, resource Volume, options *VolumesClientBeginCreateOptions) (*runtime.Poller[VolumesClientCreateResponse], error)`
- New function `*VolumesClient.BeginDelete(ctx context.Context, resourceGroupName string, storagePoolName string, volumeGroupName string, volumeName string, options *VolumesClientBeginDeleteOptions) (*runtime.Poller[VolumesClientDeleteResponse], error)`
- New function `*VolumesClient.Get(ctx context.Context, resourceGroupName string, storagePoolName string, volumeGroupName string, volumeName string, options *VolumesClientGetOptions) (VolumesClientGetResponse, error)`
- New function `*VolumesClient.NewListByVolumeGroupPager(resourceGroupName string, storagePoolName string, volumeGroupName string, options *VolumesClientListByVolumeGroupOptions) *runtime.Pager[VolumesClientListByVolumeGroupResponse]`
- New function `*VolumesClient.BeginUpdate(ctx context.Context, resourceGroupName string, storagePoolName string, volumeGroupName string, volumeName string, properties VolumeUpdate, options *VolumesClientBeginUpdateOptions) (*runtime.Poller[VolumesClientUpdateResponse], error)`
- New struct `AzureVolumeProperties`
- New struct `ConnectionParametersResponse`
- New struct `IscsiConnectionParameters`
- New struct `IscsiEndpoint`
- New struct `PerformanceParameters`
- New struct `ProtectionParameters`
- New struct `Volume`
- New struct `VolumeGroup`
- New struct `VolumeGroupListResult`
- New struct `VolumeGroupProperties`
- New struct `VolumeGroupStatus`
- New struct `VolumeGroupUpdate`
- New struct `VolumeGroupUpdateProperties`
- New struct `VolumeListResult`
- New struct `VolumeUpdate`
- New struct `VolumeUpdateProperties`


## 1.0.0 (2025-07-01)
### Other Changes

* Updated to use API version 2024-11-01

## 0.1.0 (2025-05-27)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/purestorageblock/armpurestorageblock` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).