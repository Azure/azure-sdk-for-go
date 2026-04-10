# Release History

## 1.2.0 (2025-10-30)
### Features Added

- New value `ProvisioningStatesDeleted`, `ProvisioningStatesRestoring` added to enum type `ProvisioningStates`
- New enum type `AutoScalePolicyEnforcement` with values `AutoScalePolicyEnforcementDisabled`, `AutoScalePolicyEnforcementEnabled`, `AutoScalePolicyEnforcementNone`
- New function `*VolumesClient.BeginPreBackup(context.Context, string, string, string, VolumeNameList, *VolumesClientBeginPreBackupOptions) (*runtime.Poller[VolumesClientPreBackupResponse], error)`
- New function `*VolumesClient.BeginPreRestore(context.Context, string, string, string, DiskSnapshotList, *VolumesClientBeginPreRestoreOptions) (*runtime.Poller[VolumesClientPreRestoreResponse], error)`
- New struct `AutoScaleProperties`
- New struct `DiskSnapshotList`
- New struct `PreValidationResponse`
- New struct `ScaleUpProperties`
- New struct `VolumeNameList`
- New field `AutoScaleProperties` in struct `Properties`
- New field `AutoScaleProperties` in struct `UpdateProperties`


## 1.2.0-beta.2 (2025-04-24)
### Features Added

- New value `ProvisioningStatesDeleted`, `ProvisioningStatesRestoring`, `ProvisioningStatesSoftDeleting` added to enum type `ProvisioningStates`
- New enum type `DeleteType` with values `DeleteTypePermanent`
- New enum type `PolicyState` with values `PolicyStateDisabled`, `PolicyStateEnabled`
- New enum type `XMSAccessSoftDeletedResources` with values `XMSAccessSoftDeletedResourcesFalse`, `XMSAccessSoftDeletedResourcesTrue`
- New function `*ClientFactory.NewManagementClient() *ManagementClient`
- New function `NewManagementClient(string, azcore.TokenCredential, *arm.ClientOptions) (*ManagementClient, error)`
- New function `*ManagementClient.BeginRestoreVolume(context.Context, string, string, string, string, *ManagementClientBeginRestoreVolumeOptions) (*runtime.Poller[ManagementClientRestoreVolumeResponse], error)`
- New function `*VolumesClient.BeginPreBackup(context.Context, string, string, string, VolumeNameList, *VolumesClientBeginPreBackupOptions) (*runtime.Poller[VolumesClientPreBackupResponse], error)`
- New function `*VolumesClient.BeginPreRestore(context.Context, string, string, string, DiskSnapshotList, *VolumesClientBeginPreRestoreOptions) (*runtime.Poller[VolumesClientPreRestoreResponse], error)`
- New struct `DeleteRetentionPolicy`
- New struct `DiskSnapshotList`
- New struct `PreValidationResponse`
- New struct `VolumeNameList`
- New field `DeleteRetentionPolicy` in struct `VolumeGroupProperties`
- New field `DeleteRetentionPolicy` in struct `VolumeGroupUpdateProperties`
- New field `XMSAccessSoftDeletedResources` in struct `VolumeGroupsClientListByElasticSanOptions`
- New field `DeleteType` in struct `VolumesClientBeginDeleteOptions`
- New field `XMSAccessSoftDeletedResources` in struct `VolumesClientListByVolumeGroupOptions`


## 1.2.0-beta.1 (2024-10-23)
### Features Added

- New enum type `AutoScalePolicyEnforcement` with values `AutoScalePolicyEnforcementDisabled`, `AutoScalePolicyEnforcementEnabled`, `AutoScalePolicyEnforcementNone`
- New struct `AutoScaleProperties`
- New struct `ScaleUpProperties`
- New field `AutoScaleProperties` in struct `Properties`
- New field `AutoScaleProperties` in struct `UpdateProperties`


## 1.1.0 (2024-08-22)
### Features Added

- New field `EnforceDataIntegrityCheckForIscsi` in struct `VolumeGroupProperties`
- New field `EnforceDataIntegrityCheckForIscsi` in struct `VolumeGroupUpdateProperties`


## 1.0.0 (2024-01-26)
### Other Changes

- Release stable version.


## 0.5.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.4.0 (2023-10-27)
### Breaking Changes

- Type of `SourceCreationData.CreateSource` has been changed from `*string` to `*VolumeCreateOption`
- Type of `VirtualNetworkRule.Action` has been changed from `*string` to `*Action`
- Enum `State` has been removed
- Field `SourceURI` of struct `SourceCreationData` has been removed
- Field `State` of struct `VirtualNetworkRule` has been removed

### Features Added

- New value `EncryptionTypeEncryptionAtRestWithCustomerManagedKey` added to enum type `EncryptionType`
- New enum type `Action` with values `ActionAllow`
- New enum type `IdentityType` with values `IdentityTypeNone`, `IdentityTypeSystemAssigned`, `IdentityTypeUserAssigned`
- New enum type `PublicNetworkAccess` with values `PublicNetworkAccessDisabled`, `PublicNetworkAccessEnabled`
- New enum type `VolumeCreateOption` with values `VolumeCreateOptionDisk`, `VolumeCreateOptionDiskRestorePoint`, `VolumeCreateOptionDiskSnapshot`, `VolumeCreateOptionNone`, `VolumeCreateOptionVolumeSnapshot`
- New enum type `XMSDeleteSnapshots` with values `XMSDeleteSnapshotsFalse`, `XMSDeleteSnapshotsTrue`
- New enum type `XMSForceDelete` with values `XMSForceDeleteFalse`, `XMSForceDeleteTrue`
- New function `*ClientFactory.NewVolumeSnapshotsClient() *VolumeSnapshotsClient`
- New function `NewVolumeSnapshotsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*VolumeSnapshotsClient, error)`
- New function `*VolumeSnapshotsClient.BeginCreate(context.Context, string, string, string, string, Snapshot, *VolumeSnapshotsClientBeginCreateOptions) (*runtime.Poller[VolumeSnapshotsClientCreateResponse], error)`
- New function `*VolumeSnapshotsClient.BeginDelete(context.Context, string, string, string, string, *VolumeSnapshotsClientBeginDeleteOptions) (*runtime.Poller[VolumeSnapshotsClientDeleteResponse], error)`
- New function `*VolumeSnapshotsClient.Get(context.Context, string, string, string, string, *VolumeSnapshotsClientGetOptions) (VolumeSnapshotsClientGetResponse, error)`
- New function `*VolumeSnapshotsClient.NewListByVolumeGroupPager(string, string, string, *VolumeSnapshotsClientListByVolumeGroupOptions) *runtime.Pager[VolumeSnapshotsClientListByVolumeGroupResponse]`
- New struct `EncryptionIdentity`
- New struct `EncryptionProperties`
- New struct `Identity`
- New struct `KeyVaultProperties`
- New struct `ManagedByInfo`
- New struct `Snapshot`
- New struct `SnapshotCreationData`
- New struct `SnapshotList`
- New struct `SnapshotProperties`
- New struct `UserAssignedIdentity`
- New field `PublicNetworkAccess` in struct `Properties`
- New field `SourceID` in struct `SourceCreationData`
- New field `PublicNetworkAccess` in struct `UpdateProperties`
- New field `Identity` in struct `VolumeGroup`
- New field `EncryptionProperties` in struct `VolumeGroupProperties`
- New field `Identity` in struct `VolumeGroupUpdate`
- New field `EncryptionProperties` in struct `VolumeGroupUpdateProperties`
- New field `ManagedBy`, `ProvisioningState` in struct `VolumeProperties`
- New field `ManagedBy` in struct `VolumeUpdateProperties`
- New field `XMSDeleteSnapshots`, `XMSForceDelete` in struct `VolumesClientBeginDeleteOptions`


## 0.3.0 (2023-07-28)
### Breaking Changes

- Type of `OperationListResult.Value` has been changed from `[]*RPOperation` to `[]*Operation`
- Struct `RPOperation` has been removed
- Field `Tags` of struct `Volume` has been removed
- Field `Tags` of struct `VolumeGroup` has been removed
- Field `Tags` of struct `VolumeGroupUpdate` has been removed
- Field `Tags` of struct `VolumeUpdate` has been removed

### Features Added

- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `PrivateEndpointServiceConnectionStatus` with values `PrivateEndpointServiceConnectionStatusApproved`, `PrivateEndpointServiceConnectionStatusFailed`, `PrivateEndpointServiceConnectionStatusPending`, `PrivateEndpointServiceConnectionStatusRejected`
- New function `*ClientFactory.NewPrivateEndpointConnectionsClient() *PrivateEndpointConnectionsClient`
- New function `*ClientFactory.NewPrivateLinkResourcesClient() *PrivateLinkResourcesClient`
- New function `NewPrivateEndpointConnectionsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateEndpointConnectionsClient, error)`
- New function `*PrivateEndpointConnectionsClient.BeginCreate(context.Context, string, string, string, PrivateEndpointConnection, *PrivateEndpointConnectionsClientBeginCreateOptions) (*runtime.Poller[PrivateEndpointConnectionsClientCreateResponse], error)`
- New function `*PrivateEndpointConnectionsClient.BeginDelete(context.Context, string, string, string, *PrivateEndpointConnectionsClientBeginDeleteOptions) (*runtime.Poller[PrivateEndpointConnectionsClientDeleteResponse], error)`
- New function `*PrivateEndpointConnectionsClient.Get(context.Context, string, string, string, *PrivateEndpointConnectionsClientGetOptions) (PrivateEndpointConnectionsClientGetResponse, error)`
- New function `*PrivateEndpointConnectionsClient.NewListPager(string, string, *PrivateEndpointConnectionsClientListOptions) *runtime.Pager[PrivateEndpointConnectionsClientListResponse]`
- New function `NewPrivateLinkResourcesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PrivateLinkResourcesClient, error)`
- New function `*PrivateLinkResourcesClient.ListByElasticSan(context.Context, string, string, *PrivateLinkResourcesClientListByElasticSanOptions) (PrivateLinkResourcesClientListByElasticSanResponse, error)`
- New struct `Operation`
- New struct `PrivateEndpoint`
- New struct `PrivateEndpointConnection`
- New struct `PrivateEndpointConnectionListResult`
- New struct `PrivateEndpointConnectionProperties`
- New struct `PrivateLinkResource`
- New struct `PrivateLinkResourceListResult`
- New struct `PrivateLinkResourceProperties`
- New struct `PrivateLinkServiceConnectionState`
- New field `PrivateEndpointConnections` in struct `Properties`
- New field `NextLink` in struct `SKUInformationList`
- New field `PrivateEndpointConnections` in struct `VolumeGroupProperties`


## 0.2.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.1.0 (2022-10-21)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/elasticsan/armelasticsan` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).