# Release History

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