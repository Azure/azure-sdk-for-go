# Release History

## 2.0.0 (2023-03-24)
### Breaking Changes

- Struct `CloudError` has been removed
- Struct `CloudErrorBody` has been removed

### Features Added

- New function `NewClientFactory(string, azcore.TokenCredential, *arm.ClientOptions) (*ClientFactory, error)`
- New function `*ClientFactory.NewMoveCollectionsClient() *MoveCollectionsClient`
- New function `*ClientFactory.NewMoveResourcesClient() *MoveResourcesClient`
- New function `*ClientFactory.NewOperationsDiscoveryClient() *OperationsDiscoveryClient`
- New function `*ClientFactory.NewUnresolvedDependenciesClient() *UnresolvedDependenciesClient`
- New struct `ClientFactory`
- New field `TargetResourceGroupName` in struct `AvailabilitySetResourceSettings`
- New field `TargetResourceGroupName` in struct `DiskEncryptionSetResourceSettings`
- New field `TargetResourceGroupName` in struct `KeyVaultResourceSettings`
- New field `TargetResourceGroupName` in struct `LoadBalancerResourceSettings`
- New field `Version` in struct `MoveCollectionProperties`
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


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcemover/armresourcemover` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).