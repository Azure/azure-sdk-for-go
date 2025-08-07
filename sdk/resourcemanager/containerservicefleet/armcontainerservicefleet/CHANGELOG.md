# Release History

## 3.0.0-beta.1 (2025-08-07)
### Breaking Changes

- Struct `ErrorAdditionalInfoInfo` has been removed

### Features Added

- Type of `ErrorAdditionalInfo.Info` has been changed from `*ErrorAdditionalInfoInfo` to `any`
- New value `UpdateStatePending` added to enum type `UpdateState`
- New value `UpgradeChannelTargetKubernetesVersion` added to enum type `UpgradeChannel`
- New enum type `GateProvisioningState` with values `GateProvisioningStateCanceled`, `GateProvisioningStateFailed`, `GateProvisioningStateSucceeded`
- New enum type `GateState` with values `GateStateCompleted`, `GateStatePending`, `GateStateSkipped`
- New enum type `GateType` with values `GateTypeApproval`
- New enum type `Timing` with values `TimingAfter`, `TimingBefore`
- New function `*ClientFactory.NewGatesClient() *GatesClient`
- New function `NewGatesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*GatesClient, error)`
- New function `*GatesClient.Get(context.Context, string, string, string, *GatesClientGetOptions) (GatesClientGetResponse, error)`
- New function `*GatesClient.NewListByFleetPager(string, string, *GatesClientListByFleetOptions) *runtime.Pager[GatesClientListByFleetResponse]`
- New function `*GatesClient.BeginUpdate(context.Context, string, string, string, GatePatch, *GatesClientBeginUpdateOptions) (*runtime.Poller[GatesClientUpdateResponse], error)`
- New struct `Gate`
- New struct `GateConfiguration`
- New struct `GateListResult`
- New struct `GatePatch`
- New struct `GatePatchProperties`
- New struct `GateProperties`
- New struct `GateTarget`
- New struct `UpdateRunGateStatus`
- New struct `UpdateRunGateTargetProperties`
- New field `LongTermSupport`, `TargetKubernetesVersion` in struct `AutoUpgradeProfileProperties`
- New field `Labels` in struct `FleetMemberProperties`
- New field `Labels` in struct `FleetMemberUpdateProperties`
- New field `AfterGates`, `BeforeGates` in struct `UpdateGroup`
- New field `AfterGates`, `BeforeGates` in struct `UpdateGroupStatus`
- New field `AfterGates`, `BeforeGates` in struct `UpdateStage`
- New field `AfterGates`, `BeforeGates` in struct `UpdateStageStatus`


## 2.0.0 (2025-04-15)
### Breaking Changes

- Type of `ErrorAdditionalInfo.Info` has been changed from `any` to `*ErrorAdditionalInfoInfo`
- `ManagedServiceIdentityTypeSystemAssignedUserAssigned` from enum `ManagedServiceIdentityType` has been removed
- Function `*FleetMembersClient.BeginUpdate` has been removed
- Function `*FleetsClient.BeginCreateOrUpdate` has been removed
- Function `*FleetsClient.BeginUpdate` has been removed

### Features Added

- New value `ManagedServiceIdentityTypeSystemAndUserAssigned` added to enum type `ManagedServiceIdentityType`
- New value `NodeImageSelectionTypeCustom` added to enum type `NodeImageSelectionType`
- New enum type `AutoUpgradeLastTriggerStatus` with values `AutoUpgradeLastTriggerStatusFailed`, `AutoUpgradeLastTriggerStatusSucceeded`
- New enum type `AutoUpgradeNodeImageSelectionType` with values `AutoUpgradeNodeImageSelectionTypeConsistent`, `AutoUpgradeNodeImageSelectionTypeLatest`
- New enum type `AutoUpgradeProfileProvisioningState` with values `AutoUpgradeProfileProvisioningStateCanceled`, `AutoUpgradeProfileProvisioningStateFailed`, `AutoUpgradeProfileProvisioningStateSucceeded`
- New enum type `UpgradeChannel` with values `UpgradeChannelNodeImage`, `UpgradeChannelRapid`, `UpgradeChannelStable`
- New function `NewAutoUpgradeProfileOperationsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AutoUpgradeProfileOperationsClient, error)`
- New function `*AutoUpgradeProfileOperationsClient.BeginGenerateUpdateRun(context.Context, string, string, string, *AutoUpgradeProfileOperationsClientBeginGenerateUpdateRunOptions) (*runtime.Poller[AutoUpgradeProfileOperationsClientGenerateUpdateRunResponse], error)`
- New function `NewAutoUpgradeProfilesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AutoUpgradeProfilesClient, error)`
- New function `*AutoUpgradeProfilesClient.BeginCreateOrUpdate(context.Context, string, string, string, AutoUpgradeProfile, *AutoUpgradeProfilesClientBeginCreateOrUpdateOptions) (*runtime.Poller[AutoUpgradeProfilesClientCreateOrUpdateResponse], error)`
- New function `*AutoUpgradeProfilesClient.BeginDelete(context.Context, string, string, string, *AutoUpgradeProfilesClientBeginDeleteOptions) (*runtime.Poller[AutoUpgradeProfilesClientDeleteResponse], error)`
- New function `*AutoUpgradeProfilesClient.Get(context.Context, string, string, string, *AutoUpgradeProfilesClientGetOptions) (AutoUpgradeProfilesClientGetResponse, error)`
- New function `*AutoUpgradeProfilesClient.NewListByFleetPager(string, string, *AutoUpgradeProfilesClientListByFleetOptions) *runtime.Pager[AutoUpgradeProfilesClientListByFleetResponse]`
- New function `*ClientFactory.NewAutoUpgradeProfileOperationsClient() *AutoUpgradeProfileOperationsClient`
- New function `*ClientFactory.NewAutoUpgradeProfilesClient() *AutoUpgradeProfilesClient`
- New function `*FleetMembersClient.BeginUpdateAsync(context.Context, string, string, string, FleetMemberUpdate, *FleetMembersClientBeginUpdateAsyncOptions) (*runtime.Poller[FleetMembersClientUpdateAsyncResponse], error)`
- New function `*FleetsClient.BeginCreate(context.Context, string, string, Fleet, *FleetsClientBeginCreateOptions) (*runtime.Poller[FleetsClientCreateResponse], error)`
- New function `*FleetsClient.BeginUpdateAsync(context.Context, string, string, FleetPatch, *FleetsClientBeginUpdateAsyncOptions) (*runtime.Poller[FleetsClientUpdateAsyncResponse], error)`
- New struct `AutoUpgradeNodeImageSelection`
- New struct `AutoUpgradeProfile`
- New struct `AutoUpgradeProfileListResult`
- New struct `AutoUpgradeProfileProperties`
- New struct `AutoUpgradeProfileStatus`
- New struct `ErrorAdditionalInfoInfo`
- New struct `FleetMemberStatus`
- New struct `FleetStatus`
- New struct `GenerateResponse`
- New field `EnableVnetIntegration`, `SubnetID` in struct `APIServerAccessProfile`
- New field `Status` in struct `FleetMemberProperties`
- New field `Status` in struct `FleetProperties`
- New field `CustomNodeImageVersions` in struct `NodeImageSelection`
- New field `AutoUpgradeProfileID` in struct `UpdateRunProperties`


## 1.3.0-beta.1 (2024-10-25)
### Features Added

- New value `NodeImageSelectionTypeCustom` added to enum type `NodeImageSelectionType`
- New enum type `AutoUpgradeNodeImageSelectionType` with values `AutoUpgradeNodeImageSelectionTypeConsistent`, `AutoUpgradeNodeImageSelectionTypeLatest`
- New enum type `AutoUpgradeProfileProvisioningState` with values `AutoUpgradeProfileProvisioningStateCanceled`, `AutoUpgradeProfileProvisioningStateFailed`, `AutoUpgradeProfileProvisioningStateSucceeded`
- New enum type `UpgradeChannel` with values `UpgradeChannelNodeImage`, `UpgradeChannelRapid`, `UpgradeChannelStable`
- New function `NewAutoUpgradeProfilesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*AutoUpgradeProfilesClient, error)`
- New function `*AutoUpgradeProfilesClient.BeginCreateOrUpdate(context.Context, string, string, string, AutoUpgradeProfile, *AutoUpgradeProfilesClientBeginCreateOrUpdateOptions) (*runtime.Poller[AutoUpgradeProfilesClientCreateOrUpdateResponse], error)`
- New function `*AutoUpgradeProfilesClient.BeginDelete(context.Context, string, string, string, *AutoUpgradeProfilesClientBeginDeleteOptions) (*runtime.Poller[AutoUpgradeProfilesClientDeleteResponse], error)`
- New function `*AutoUpgradeProfilesClient.Get(context.Context, string, string, string, *AutoUpgradeProfilesClientGetOptions) (AutoUpgradeProfilesClientGetResponse, error)`
- New function `*AutoUpgradeProfilesClient.NewListByFleetPager(string, string, *AutoUpgradeProfilesClientListByFleetOptions) *runtime.Pager[AutoUpgradeProfilesClientListByFleetResponse]`
- New function `*ClientFactory.NewAutoUpgradeProfilesClient() *AutoUpgradeProfilesClient`
- New struct `AutoUpgradeNodeImageSelection`
- New struct `AutoUpgradeProfile`
- New struct `AutoUpgradeProfileListResult`
- New struct `AutoUpgradeProfileProperties`
- New field `EnableVnetIntegration`, `SubnetID` in struct `APIServerAccessProfile`
- New field `CustomNodeImageVersions` in struct `NodeImageSelection`


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