# Release History

## 1.0.0 (2024-09-26)
### Breaking Changes

- Type of `StandbyContainerGroupPoolResourceUpdateProperties.ContainerGroupProperties` has been changed from `*ContainerGroupPropertiesUpdate` to `*ContainerGroupProperties`
- Type of `StandbyContainerGroupPoolResourceUpdateProperties.ElasticityProfile` has been changed from `*StandbyContainerGroupPoolElasticityProfileUpdate` to `*StandbyContainerGroupPoolElasticityProfile`
- Type of `StandbyVirtualMachinePoolResourceUpdateProperties.ElasticityProfile` has been changed from `*StandbyVirtualMachinePoolElasticityProfileUpdate` to `*StandbyVirtualMachinePoolElasticityProfile`
- Struct `ContainerGroupProfileUpdate` has been removed
- Struct `ContainerGroupPropertiesUpdate` has been removed
- Struct `StandbyContainerGroupPoolElasticityProfileUpdate` has been removed
- Struct `StandbyVirtualMachinePoolElasticityProfileUpdate` has been removed

### Features Added

- New function `*ClientFactory.NewStandbyContainerGroupPoolRuntimeViewsClient() *StandbyContainerGroupPoolRuntimeViewsClient`
- New function `*ClientFactory.NewStandbyVirtualMachinePoolRuntimeViewsClient() *StandbyVirtualMachinePoolRuntimeViewsClient`
- New function `NewStandbyContainerGroupPoolRuntimeViewsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*StandbyContainerGroupPoolRuntimeViewsClient, error)`
- New function `*StandbyContainerGroupPoolRuntimeViewsClient.Get(context.Context, string, string, string, *StandbyContainerGroupPoolRuntimeViewsClientGetOptions) (StandbyContainerGroupPoolRuntimeViewsClientGetResponse, error)`
- New function `*StandbyContainerGroupPoolRuntimeViewsClient.NewListByStandbyPoolPager(string, string, *StandbyContainerGroupPoolRuntimeViewsClientListByStandbyPoolOptions) *runtime.Pager[StandbyContainerGroupPoolRuntimeViewsClientListByStandbyPoolResponse]`
- New function `NewStandbyVirtualMachinePoolRuntimeViewsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*StandbyVirtualMachinePoolRuntimeViewsClient, error)`
- New function `*StandbyVirtualMachinePoolRuntimeViewsClient.Get(context.Context, string, string, string, *StandbyVirtualMachinePoolRuntimeViewsClientGetOptions) (StandbyVirtualMachinePoolRuntimeViewsClientGetResponse, error)`
- New function `*StandbyVirtualMachinePoolRuntimeViewsClient.NewListByStandbyPoolPager(string, string, *StandbyVirtualMachinePoolRuntimeViewsClientListByStandbyPoolOptions) *runtime.Pager[StandbyVirtualMachinePoolRuntimeViewsClientListByStandbyPoolResponse]`
- New struct `ContainerGroupInstanceCountSummary`
- New struct `PoolResourceStateCount`
- New struct `StandbyContainerGroupPoolRuntimeViewResource`
- New struct `StandbyContainerGroupPoolRuntimeViewResourceListResult`
- New struct `StandbyContainerGroupPoolRuntimeViewResourceProperties`
- New struct `StandbyVirtualMachinePoolRuntimeViewResource`
- New struct `StandbyVirtualMachinePoolRuntimeViewResourceListResult`
- New struct `StandbyVirtualMachinePoolRuntimeViewResourceProperties`
- New struct `VirtualMachineInstanceCountSummary`
- New field `MinReadyCapacity` in struct `StandbyVirtualMachinePoolElasticityProfile`


## 0.1.0 (2024-04-26)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/standbypool/armstandbypool` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).