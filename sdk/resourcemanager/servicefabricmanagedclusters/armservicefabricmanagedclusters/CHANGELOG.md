# Release History

## 0.6.0 (2025-11-10)
### Breaking Changes

- Operation `*ApplicationsClient.Update` has been changed to LRO, use `*ApplicationsClient.BeginUpdate` instead.
- Operation `*ManagedClustersClient.Update` has been changed to LRO, use `*ManagedClustersClient.BeginUpdate` instead.

### Features Added

- New enum type `HealthFilter` with values `HealthFilterAll`, `HealthFilterDefault`, `HealthFilterError`, `HealthFilterNone`, `HealthFilterOk`, `HealthFilterWarning`
- New enum type `RestartKind` with values `RestartKindSimultaneous`
- New enum type `RuntimeFailureAction` with values `RuntimeFailureActionManual`, `RuntimeFailureActionRollback`
- New enum type `RuntimeRollingUpgradeMode` with values `RuntimeRollingUpgradeModeMonitored`, `RuntimeRollingUpgradeModeUnmonitoredAuto`, `RuntimeRollingUpgradeModeUnmonitoredManual`
- New enum type `RuntimeUpgradeKind` with values `RuntimeUpgradeKindRolling`
- New function `*ApplicationsClient.BeginFetchHealth(context.Context, string, string, string, ApplicationFetchHealthRequest, *ApplicationsClientBeginFetchHealthOptions) (*runtime.Poller[ApplicationsClientFetchHealthResponse], error)`
- New function `*ApplicationsClient.BeginRestartDeployedCodePackage(context.Context, string, string, string, RestartDeployedCodePackageRequest, *ApplicationsClientBeginRestartDeployedCodePackageOptions) (*runtime.Poller[ApplicationsClientRestartDeployedCodePackageResponse], error)`
- New function `*ApplicationsClient.BeginUpdateUpgrade(context.Context, string, string, string, RuntimeUpdateApplicationUpgradeParameters, *ApplicationsClientBeginUpdateUpgradeOptions) (*runtime.Poller[ApplicationsClientUpdateUpgradeResponse], error)`
- New function `*ServicesClient.BeginRestartReplica(context.Context, string, string, string, string, RestartReplicaRequest, *ServicesClientBeginRestartReplicaOptions) (*runtime.Poller[ServicesClientRestartReplicaResponse], error)`
- New struct `ApplicationFetchHealthRequest`
- New struct `ApplicationUpdateParametersProperties`
- New struct `RestartDeployedCodePackageRequest`
- New struct `RestartReplicaRequest`
- New struct `RuntimeApplicationHealthPolicy`
- New struct `RuntimeRollingUpgradeUpdateMonitoringPolicy`
- New struct `RuntimeServiceTypeHealthPolicy`
- New struct `RuntimeUpdateApplicationUpgradeParameters`
- New field `Properties` in struct `ApplicationUpdateParameters`


## 0.5.0 (2025-08-12)
### Features Added

- New field `EnableOutboundOnlyNodeTypes` in struct `ManagedClusterProperties`
- New field `IsOutboundOnly` in struct `NodeTypeProperties`
- New field `NetworkIdentifier` in struct `ServiceEndpoint`


## 0.4.0 (2025-06-24)
### Breaking Changes

- Type of `ManagedClusterVersionDetails.SupportExpiryUTC` has been changed from `*string` to `*time.Time`
- Type of `SystemData.CreatedByType` has been changed from `*string` to `*CreatedByType`
- Type of `SystemData.LastModifiedByType` has been changed from `*string` to `*CreatedByType`

### Features Added

- New value `DiskTypePremiumV2LRS`, `DiskTypePremiumZRS`, `DiskTypeStandardSSDZRS` added to enum type `DiskType`
- New value `SecurityTypeConfidentialVM` added to enum type `SecurityType`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `FaultKind` with values `FaultKindZone`
- New enum type `FaultSimulationStatus` with values `FaultSimulationStatusActive`, `FaultSimulationStatusDone`, `FaultSimulationStatusStartFailed`, `FaultSimulationStatusStarting`, `FaultSimulationStatusStopFailed`, `FaultSimulationStatusStopping`
- New enum type `SecurityEncryptionType` with values `SecurityEncryptionTypeDiskWithVMGuestState`, `SecurityEncryptionTypeVMGuestStateOnly`
- New enum type `SfmcOperationStatus` with values `SfmcOperationStatusAborted`, `SfmcOperationStatusCanceled`, `SfmcOperationStatusCreated`, `SfmcOperationStatusFailed`, `SfmcOperationStatusStarted`, `SfmcOperationStatusSucceeded`
- New function `*FaultSimulationContent.GetFaultSimulationContent() *FaultSimulationContent`
- New function `*ManagedClustersClient.GetFaultSimulation(context.Context, string, string, FaultSimulationIDContent, *ManagedClustersClientGetFaultSimulationOptions) (ManagedClustersClientGetFaultSimulationResponse, error)`
- New function `*ManagedClustersClient.NewListFaultSimulationPager(string, string, *ManagedClustersClientListFaultSimulationOptions) *runtime.Pager[ManagedClustersClientListFaultSimulationResponse]`
- New function `*ManagedClustersClient.BeginStartFaultSimulation(context.Context, string, string, FaultSimulationContentWrapper, *ManagedClustersClientBeginStartFaultSimulationOptions) (*runtime.Poller[ManagedClustersClientStartFaultSimulationResponse], error)`
- New function `*ManagedClustersClient.BeginStopFaultSimulation(context.Context, string, string, FaultSimulationIDContent, *ManagedClustersClientBeginStopFaultSimulationOptions) (*runtime.Poller[ManagedClustersClientStopFaultSimulationResponse], error)`
- New function `*ZoneFaultSimulationContent.GetFaultSimulationContent() *FaultSimulationContent`
- New function `*NodeTypesClient.BeginDeallocate(context.Context, string, string, string, NodeTypeActionParameters, *NodeTypesClientBeginDeallocateOptions) (*runtime.Poller[NodeTypesClientDeallocateResponse], error)`
- New function `*NodeTypesClient.GetFaultSimulation(context.Context, string, string, string, FaultSimulationIDContent, *NodeTypesClientGetFaultSimulationOptions) (NodeTypesClientGetFaultSimulationResponse, error)`
- New function `*NodeTypesClient.NewListFaultSimulationPager(string, string, string, *NodeTypesClientListFaultSimulationOptions) *runtime.Pager[NodeTypesClientListFaultSimulationResponse]`
- New function `*NodeTypesClient.BeginRedeploy(context.Context, string, string, string, NodeTypeActionParameters, *NodeTypesClientBeginRedeployOptions) (*runtime.Poller[NodeTypesClientRedeployResponse], error)`
- New function `*NodeTypesClient.BeginStart(context.Context, string, string, string, NodeTypeActionParameters, *NodeTypesClientBeginStartOptions) (*runtime.Poller[NodeTypesClientStartResponse], error)`
- New function `*NodeTypesClient.BeginStartFaultSimulation(context.Context, string, string, string, FaultSimulationContentWrapper, *NodeTypesClientBeginStartFaultSimulationOptions) (*runtime.Poller[NodeTypesClientStartFaultSimulationResponse], error)`
- New function `*NodeTypesClient.BeginStopFaultSimulation(context.Context, string, string, string, FaultSimulationIDContent, *NodeTypesClientBeginStopFaultSimulationOptions) (*runtime.Poller[NodeTypesClientStopFaultSimulationResponse], error)`
- New struct `FaultSimulation`
- New struct `FaultSimulationConstraints`
- New struct `FaultSimulationContentWrapper`
- New struct `FaultSimulationDetails`
- New struct `FaultSimulationIDContent`
- New struct `FaultSimulationListResult`
- New struct `NodeTypeFaultSimulation`
- New struct `ZoneFaultSimulationContent`
- New field `VMImage` in struct `ManagedClusterProperties`
- New field `SecurityEncryptionType`, `ZoneBalance` in struct `NodeTypeProperties`


## 0.3.0 (2024-12-27)
### Breaking Changes

- Operation `*NodeTypesClient.Update` has been changed to LRO, use `*NodeTypesClient.BeginUpdate` instead.
- Field `CustomFqdn` of struct `ManagedClusterProperties` has been removed

### Features Added

- New field `AllocatedOutboundPorts` in struct `ManagedClusterProperties`


## 0.2.0 (2024-10-23)
### Features Added

- New enum type `AutoGeneratedDomainNameLabelScope` with values `AutoGeneratedDomainNameLabelScopeNoReuse`, `AutoGeneratedDomainNameLabelScopeResourceGroupReuse`, `AutoGeneratedDomainNameLabelScopeSubscriptionReuse`, `AutoGeneratedDomainNameLabelScopeTenantReuse`
- New struct `VMApplication`
- New field `AutoGeneratedDomainNameLabelScope`, `CustomFqdn` in struct `ManagedClusterProperties`
- New field `VMApplications` in struct `NodeTypeProperties`


## 0.1.0 (2024-07-29)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicefabricmanagedclusters/armservicefabricmanagedclusters` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).