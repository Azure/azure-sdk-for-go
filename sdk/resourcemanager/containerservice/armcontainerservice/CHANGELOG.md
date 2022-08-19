# Release History

## 2.0.0-beta.3 (2022-08-19)
### Breaking Changes

- Function `*FleetsClient.Update` parameter(s) have been changed from `(context.Context, string, string, FleetPatch, *FleetsClientUpdateOptions)` to `(context.Context, string, string, *FleetsClientUpdateOptions)`

### Features Added

- New const `UpdateModeOff`
- New const `UpdateModeAuto`
- New const `ControlledValuesRequestsAndLimits`
- New const `ControlledValuesRequestsOnly`
- New const `PublicNetworkAccessSecuredByPerimeter`
- New const `UpdateModeRecreate`
- New const `UpdateModeInitial`
- New type alias `UpdateMode`
- New type alias `ControlledValues`
- New function `*ManagedClustersClient.AbortLatestOperation(context.Context, string, string, *ManagedClustersClientAbortLatestOperationOptions) (ManagedClustersClientAbortLatestOperationResponse, error)`
- New function `PossibleUpdateModeValues() []UpdateMode`
- New function `*AgentPoolsClient.AbortLatestOperation(context.Context, string, string, string, *AgentPoolsClientAbortLatestOperationOptions) (AgentPoolsClientAbortLatestOperationResponse, error)`
- New function `PossibleControlledValuesValues() []ControlledValues`
- New struct `AgentPoolsClientAbortLatestOperationOptions`
- New struct `AgentPoolsClientAbortLatestOperationResponse`
- New struct `ManagedClusterAzureMonitorProfile`
- New struct `ManagedClusterAzureMonitorProfileKubeStateMetrics`
- New struct `ManagedClusterAzureMonitorProfileMetrics`
- New struct `ManagedClusterSecurityProfileImageCleaner`
- New struct `ManagedClusterWorkloadAutoScalerProfileVerticalPodAutoscaler`
- New struct `ManagedClustersClientAbortLatestOperationOptions`
- New struct `ManagedClustersClientAbortLatestOperationResponse`
- New field `VerticalPodAutoscaler` in struct `ManagedClusterWorkloadAutoScalerProfile`
- New field `AzureMonitorProfile` in struct `ManagedClusterProperties`
- New field `ImageCleaner` in struct `ManagedClusterSecurityProfile`
- New field `Parameters` in struct `FleetsClientUpdateOptions`


## 2.0.0 (2022-07-22)
### Breaking Changes

- Struct `ManagedClusterSecurityProfileAzureDefender` has been removed
- Field `AzureDefender` of struct `ManagedClusterSecurityProfile` has been removed

### Features Added

- New const `KeyVaultNetworkAccessTypesPrivate`
- New const `NetworkPluginNone`
- New const `KeyVaultNetworkAccessTypesPublic`
- New function `PossibleKeyVaultNetworkAccessTypesValues() []KeyVaultNetworkAccessTypes`
- New struct `AzureKeyVaultKms`
- New struct `ManagedClusterSecurityProfileDefender`
- New struct `ManagedClusterSecurityProfileDefenderSecurityMonitoring`
- New field `HostGroupID` in struct `ManagedClusterAgentPoolProfileProperties`
- New field `HostGroupID` in struct `ManagedClusterAgentPoolProfile`
- New field `AzureKeyVaultKms` in struct `ManagedClusterSecurityProfile`
- New field `Defender` in struct `ManagedClusterSecurityProfile`


## 1.0.0 (2022-05-16)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).