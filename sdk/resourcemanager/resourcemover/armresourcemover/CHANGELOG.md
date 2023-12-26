# Release History

## 1.3.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.2.0 (2023-10-27)
### Features Added

- New enum type `MoveType` with values `MoveTypeRegionToRegion`, `MoveTypeRegionToZone`
- New field `TargetResourceGroupName` in struct `AvailabilitySetResourceSettings`
- New field `TargetResourceGroupName` in struct `DiskEncryptionSetResourceSettings`
- New field `TargetResourceGroupName` in struct `KeyVaultResourceSettings`
- New field `TargetResourceGroupName` in struct `LoadBalancerResourceSettings`
- New field `MoveRegion`, `MoveType`, `Version` in struct `MoveCollectionProperties`
- New field `TargetResourceGroupName` in struct `NetworkInterfaceResourceSettings`
- New field `TargetResourceGroupName` in struct `NetworkSecurityGroupResourceSettings`
- New field `TargetResourceGroupName` in struct `PublicIPAddressResourceSettings`
- New field `TargetResourceGroupName` in struct `ResourceGroupResourceSettings`
- New field `TargetResourceGroupName` in struct `ResourceSettings`
- New field `TargetResourceGroupName` in struct `SQLDatabaseResourceSettings`
- New field `TargetResourceGroupName` in struct `SQLElasticPoolResourceSettings`
- New field `TargetResourceGroupName` in struct `SQLServerResourceSettings`
- New field `TargetResourceGroupName` in struct `VirtualMachineResourceSettings`
- New field `TargetResourceGroupName` in struct `VirtualNetworkResourceSettings`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcemover/armresourcemover` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).