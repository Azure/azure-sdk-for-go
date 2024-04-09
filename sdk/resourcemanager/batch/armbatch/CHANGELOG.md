# Release History

## 2.3.0 (2024-03-22)
### Features Added

- New enum type `UpgradeMode` with values `UpgradeModeAutomatic`, `UpgradeModeManual`, `UpgradeModeRolling`
- New struct `AutomaticOSUpgradePolicy`
- New struct `RollingUpgradePolicy`
- New struct `UpgradePolicy`
- New field `UpgradePolicy` in struct `PoolProperties`
- New field `BatchSupportEndOfLife` in struct `SupportedSKU`


## 2.2.0 (2023-12-22)
### Features Added

- New value `StorageAccountTypeStandardSSDLRS` added to enum type `StorageAccountType`
- New struct `ManagedDisk`
- New struct `SecurityProfile`
- New struct `ServiceArtifactReference`
- New struct `UefiSettings`
- New field `Caching`, `DiskSizeGB`, `ManagedDisk`, `WriteAcceleratorEnabled` in struct `OSDisk`
- New field `ResourceTags` in struct `PoolProperties`
- New field `SecurityProfile`, `ServiceArtifactReference` in struct `VirtualMachineConfiguration`


## 2.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0 (2023-07-28)
### Breaking Changes

- Type of `ContainerConfiguration.Type` has been changed from `*string` to `*ContainerType`

### Features Added

- New enum type `ContainerType` with values `ContainerTypeCriCompatible`, `ContainerTypeDockerCompatible`
- New field `EnableAcceleratedNetworking` in struct `NetworkConfiguration`
- New field `EnableAutomaticUpgrade` in struct `VMExtension`


## 1.2.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 1.2.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.1.0 (2022-11-10)
### Features Added

- New const `PrivateEndpointConnectionProvisioningStateCreating`
- New const `NodeCommunicationModeDefault`
- New const `EndpointAccessDefaultActionAllow`
- New const `NodeCommunicationModeClassic`
- New const `PrivateEndpointConnectionProvisioningStateCancelled`
- New const `EndpointAccessDefaultActionDeny`
- New const `PrivateEndpointConnectionProvisioningStateDeleting`
- New const `NodeCommunicationModeSimplified`
- New type alias `EndpointAccessDefaultAction`
- New type alias `NodeCommunicationMode`
- New function `*PrivateEndpointConnectionClient.BeginDelete(context.Context, string, string, string, *PrivateEndpointConnectionClientBeginDeleteOptions) (*runtime.Poller[PrivateEndpointConnectionClientDeleteResponse], error)`
- New function `PossibleEndpointAccessDefaultActionValues() []EndpointAccessDefaultAction`
- New function `PossibleNodeCommunicationModeValues() []NodeCommunicationMode`
- New struct `EndpointAccessProfile`
- New struct `IPRule`
- New struct `NetworkProfile`
- New struct `PrivateEndpointConnectionClientBeginDeleteOptions`
- New struct `PrivateEndpointConnectionClientDeleteResponse`
- New field `NetworkProfile` in struct `AccountUpdateProperties`
- New field `PublicNetworkAccess` in struct `AccountUpdateProperties`
- New field `NodeManagementEndpoint` in struct `AccountProperties`
- New field `NetworkProfile` in struct `AccountProperties`
- New field `GroupIDs` in struct `PrivateEndpointConnectionProperties`
- New field `NetworkProfile` in struct `AccountCreateProperties`
- New field `TargetNodeCommunicationMode` in struct `PoolProperties`
- New field `CurrentNodeCommunicationMode` in struct `PoolProperties`


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/batch/armbatch` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).