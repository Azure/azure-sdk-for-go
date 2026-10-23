# Release History

## 1.0.0 (2026-10-23)
### Breaking Changes

- Type of `ExternalAuthProperties.Status` has been changed from `*ResourceStatus` to `*ExternalAuthResourceStatus`
- Type of `HcpOpenShiftClusterProperties.Status` has been changed from `*ResourceStatus` to `*ClusterResourceStatus`
- Type of `NodePoolProperties.Status` has been changed from `*ResourceStatus` to `*NodePoolResourceStatus`
- Struct `ResourceStatus` has been removed

### Features Added

- New enum type `KmsKeyVaultType` with values `KmsKeyVaultTypeKeyVault`, `KmsKeyVaultTypeManagedHSM`
- New struct `ClusterActiveVersion`
- New struct `ClusterResourceStatus`
- New struct `ContainerRegistryProfile`
- New struct `ExternalAuthResourceStatus`
- New struct `NodePoolActiveVersion`
- New struct `NodePoolResourceStatus`
- New field `KeyVaultType` in struct `KmsEncryptionProfile`
- New field `ContainerRegistry` in struct `PlatformProfile`


## 0.1.0 (2026-08-19)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redhatopenshifthcp/armredhatopenshifthcp` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).