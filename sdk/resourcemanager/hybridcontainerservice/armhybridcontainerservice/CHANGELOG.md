# Release History

## 0.2.0 (2023-03-24)
### Breaking Changes

- Struct `VirtualNetworksPropertiesInfraVnetProfileKubevirt` has been removed
- Field `Kubevirt` of struct `VirtualNetworksPropertiesInfraVnetProfile` has been removed

### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New function `*ProvisionedClustersClient.GetUpgradeProfile(context.Context, string, string, *ProvisionedClustersClientGetUpgradeProfileOptions) (ProvisionedClustersClientGetUpgradeProfileResponse, error)`
- New function `*ProvisionedClustersClient.BeginUpgradeNodeImageVersionForEntireCluster(context.Context, string, string, *ProvisionedClustersClientBeginUpgradeNodeImageVersionForEntireClusterOptions) (*runtime.Poller[ProvisionedClustersClientUpgradeNodeImageVersionForEntireClusterResponse], error)`
- New struct `ProvisionedClusterPoolUpgradeProfile`
- New struct `ProvisionedClusterPoolUpgradeProfileProperties`
- New struct `ProvisionedClusterUpgradeProfile`
- New struct `ProvisionedClusterUpgradeProfileProperties`
- New struct `VirtualNetworksPropertiesInfraVnetProfileNetworkCloud`
- New field `NetworkCloud` in struct `VirtualNetworksPropertiesInfraVnetProfile`


## 0.1.1 (2022-10-12)
### Other Changes
- Loosen Go version requirement.

## 0.1.0 (2022-09-13)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybridcontainerservice/armhybridcontainerservice` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.1.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).