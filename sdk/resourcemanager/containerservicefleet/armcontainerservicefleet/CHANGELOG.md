# Release History

## 1.2.0 (2024-05-24)
### Features Added

- New value `ManagedClusterUpgradeTypeControlPlaneOnly` added to enum type `ManagedClusterUpgradeType`
- New enum type `TargetType` with values `TargetTypeAfterStageWait`, `TargetTypeGroup`, `TargetTypeMember`, `TargetTypeStage`
- New function `*UpdateRunsClient.BeginSkip(context.Context, string, string, string, SkipProperties, *UpdateRunsClientBeginSkipOptions) (*runtime.Poller[UpdateRunsClientSkipResponse], error)`
- New struct `APIServerAccessProfile`
- New struct `AgentProfile`
- New struct `FleetHubProfile`
- New struct `SkipProperties`
- New struct `SkipTarget`
- New field `HubProfile` in struct `FleetProperties`


## 1.2.0-beta.1 (2024-04-26)
### Features Added

- New value `ManagedClusterUpgradeTypeControlPlaneOnly` added to enum type `ManagedClusterUpgradeType`
- New enum type `TargetType` with values `TargetTypeAfterStageWait`, `TargetTypeGroup`, `TargetTypeMember`, `TargetTypeStage`
- New function `*UpdateRunsClient.BeginSkip(context.Context, string, string, string, SkipProperties, *UpdateRunsClientBeginSkipOptions) (*runtime.Poller[UpdateRunsClientSkipResponse], error)`
- New struct `APIServerAccessProfile`
- New struct `AgentProfile`
- New struct `FleetHubProfile`
- New struct `SkipProperties`
- New struct `SkipTarget`
- New field `HubProfile` in struct `FleetProperties`


## 1.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.0.0 (2023-10-27)
### Breaking Changes

- Struct `APIServerAccessProfile` has been removed
- Struct `AgentProfile` has been removed
- Struct `FleetHubProfile` has been removed
- Field `HubProfile` of struct `FleetProperties` has been removed


## 0.3.0 (2023-10-27)
### Features Added

- New enum type `FleetUpdateStrategyProvisioningState` with values `FleetUpdateStrategyProvisioningStateCanceled`, `FleetUpdateStrategyProvisioningStateFailed`, `FleetUpdateStrategyProvisioningStateSucceeded`
- New function `*ClientFactory.NewFleetUpdateStrategiesClient() *FleetUpdateStrategiesClient`
- New function `NewFleetUpdateStrategiesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*FleetUpdateStrategiesClient, error)`
- New function `*FleetUpdateStrategiesClient.BeginCreateOrUpdate(context.Context, string, string, string, FleetUpdateStrategy, *FleetUpdateStrategiesClientBeginCreateOrUpdateOptions) (*runtime.Poller[FleetUpdateStrategiesClientCreateOrUpdateResponse], error)`
- New function `*FleetUpdateStrategiesClient.BeginDelete(context.Context, string, string, string, *FleetUpdateStrategiesClientBeginDeleteOptions) (*runtime.Poller[FleetUpdateStrategiesClientDeleteResponse], error)`
- New function `*FleetUpdateStrategiesClient.Get(context.Context, string, string, string, *FleetUpdateStrategiesClientGetOptions) (FleetUpdateStrategiesClientGetResponse, error)`
- New function `*FleetUpdateStrategiesClient.NewListByFleetPager(string, string, *FleetUpdateStrategiesClientListByFleetOptions) *runtime.Pager[FleetUpdateStrategiesClientListByFleetResponse]`
- New struct `FleetUpdateStrategy`
- New struct `FleetUpdateStrategyListResult`
- New struct `FleetUpdateStrategyProperties`
- New field `VMSize` in struct `AgentProfile`
- New field `PortalFqdn` in struct `FleetHubProfile`
- New field `UpdateStrategyID` in struct `UpdateRunProperties`


## 0.2.0 (2023-09-22)
### Breaking Changes

- Operation `*FleetMembersClient.Update` has been changed to LRO, use `*FleetMembersClient.BeginUpdate` instead.
- Operation `*FleetsClient.Update` has been changed to LRO, use `*FleetsClient.BeginUpdate` instead.

### Features Added

- New value `UpdateStateSkipped` added to enum type `UpdateState`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New enum type `NodeImageSelectionType` with values `NodeImageSelectionTypeConsistent`, `NodeImageSelectionTypeLatest`
- New struct `APIServerAccessProfile`
- New struct `AgentProfile`
- New struct `ManagedServiceIdentity`
- New struct `NodeImageSelection`
- New struct `NodeImageSelectionStatus`
- New struct `NodeImageVersion`
- New struct `UserAssignedIdentity`
- New field `Identity` in struct `Fleet`
- New field `APIServerAccessProfile`, `AgentProfile` in struct `FleetHubProfile`
- New field `Identity` in struct `FleetPatch`
- New field `NodeImageSelection` in struct `ManagedClusterUpdate`
- New field `Message` in struct `MemberUpdateStatus`
- New field `NodeImageSelection` in struct `UpdateRunStatus`


## 0.1.0 (2023-06-23)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservicefleet/armcontainerservicefleet` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).