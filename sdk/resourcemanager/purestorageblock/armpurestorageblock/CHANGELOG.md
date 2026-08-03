# Release History

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