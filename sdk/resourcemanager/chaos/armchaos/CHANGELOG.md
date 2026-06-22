# Release History

## 2.1.0-beta.1 (2026-05-06)
### Features Added

- New value `ProvisioningStateRunning` added to enum type `ProvisioningState`
- New enum type `ActionDependencyType` with values `ActionDependencyTypeAction`
- New enum type `ActionKind` with values `ActionKindCancelable`, `ActionKindContinuous`, `ActionKindDiscrete`
- New enum type `ActionLifecycle` with values `ActionLifecycleAnyTerminal`, `ActionLifecycleFailure`, `ActionLifecycleRunning`, `ActionLifecycleSkipped`, `ActionLifecycleStart`, `ActionLifecycleSuccess`
- New enum type `ParameterType` with values `ParameterTypeArray`, `ParameterTypeBoolean`, `ParameterTypeNumber`, `ParameterTypeObject`, `ParameterTypeString`
- New enum type `PermissionsFixState` with values `PermissionsFixStateFailed`, `PermissionsFixStateInProgress`, `PermissionsFixStateNotStarted`, `PermissionsFixStatePartiallySucceeded`, `PermissionsFixStateSucceeded`, `PermissionsFixStateWhatIfCompleted`
- New enum type `PrivateEndpointServiceConnectionStatus` with values `PrivateEndpointServiceConnectionStatusApproved`, `PrivateEndpointServiceConnectionStatusPending`, `PrivateEndpointServiceConnectionStatusRejected`
- New enum type `PublicNetworkAccessOption` with values `PublicNetworkAccessOptionDisabled`, `PublicNetworkAccessOptionEnabled`
- New enum type `RecommendationStatus` with values `RecommendationStatusEvaluating`, `RecommendationStatusEvaluationCancelled`, `RecommendationStatusEvaluationFailed`, `RecommendationStatusNotApplicable`, `RecommendationStatusNotEvaluated`, `RecommendationStatusRecommended`
- New enum type `RoleAssignmentStatus` with values `RoleAssignmentStatusFailed`, `RoleAssignmentStatusPending`, `RoleAssignmentStatusSkipped`, `RoleAssignmentStatusSucceeded`
- New enum type `RunAfterBehavior` with values `RunAfterBehaviorAll`, `RunAfterBehaviorAny`, `RunAfterBehaviorAtLeastOne`
- New enum type `ScenarioRunState` with values `ScenarioRunStateCanceled`, `ScenarioRunStateCanceling`, `ScenarioRunStateCleaningUp`, `ScenarioRunStateFailed`, `ScenarioRunStateGenerating`, `ScenarioRunStatePreparing`, `ScenarioRunStateQueued`, `ScenarioRunStateResolving`, `ScenarioRunStateRunning`, `ScenarioRunStateStarting`, `ScenarioRunStateSucceeded`, `ScenarioRunStateValidating`, `ScenarioRunStateValidationSucceeded`
- New enum type `ScenarioSummaryState` with values `ScenarioSummaryStateCanceled`, `ScenarioSummaryStateCanceling`, `ScenarioSummaryStateFailed`, `ScenarioSummaryStateFailingOnError`, `ScenarioSummaryStatePending`, `ScenarioSummaryStateRunning`, `ScenarioSummaryStateSkipped`, `ScenarioSummaryStateStarting`, `ScenarioSummaryStateStopping`, `ScenarioSummaryStateSucceeded`
- New enum type `ScenarioValidationState` with values `ScenarioValidationStateAccepted`, `ScenarioValidationStateGenerating`, `ScenarioValidationStateNoResolvedResources`, `ScenarioValidationStateNotStarted`, `ScenarioValidationStateRequiresAttention`, `ScenarioValidationStateResolving`, `ScenarioValidationStateSucceeded`, `ScenarioValidationStateValidating`
- New enum type `WorkspaceEvaluationStatus` with values `WorkspaceEvaluationStatusCanceled`, `WorkspaceEvaluationStatusFailed`, `WorkspaceEvaluationStatusInProgress`, `WorkspaceEvaluationStatusPartiallySucceeded`, `WorkspaceEvaluationStatusPending`, `WorkspaceEvaluationStatusQueued`, `WorkspaceEvaluationStatusSucceeded`
- New enum type `ZoneResolutionMode` with values `ZoneResolutionModeLogical`, `ZoneResolutionModePhysical`
- New function `NewActionVersionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ActionVersionsClient, error)`
- New function `*ActionVersionsClient.Get(ctx context.Context, location string, actionName string, versionName string, options *ActionVersionsClientGetOptions) (ActionVersionsClientGetResponse, error)`
- New function `*ActionVersionsClient.NewListPager(location string, actionName string, options *ActionVersionsClientListOptions) *runtime.Pager[ActionVersionsClientListResponse]`
- New function `NewActionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ActionsClient, error)`
- New function `*ActionsClient.Get(ctx context.Context, location string, actionName string, options *ActionsClientGetOptions) (ActionsClientGetResponse, error)`
- New function `*ActionsClient.NewListPager(location string, options *ActionsClientListOptions) *runtime.Pager[ActionsClientListResponse]`
- New function `*ClientFactory.NewActionVersionsClient() *ActionVersionsClient`
- New function `*ClientFactory.NewActionsClient() *ActionsClient`
- New function `*ClientFactory.NewDiscoveredResourcesClient() *DiscoveredResourcesClient`
- New function `*ClientFactory.NewPrivateAccessesClient() *PrivateAccessesClient`
- New function `*ClientFactory.NewScenarioConfigurationsClient() *ScenarioConfigurationsClient`
- New function `*ClientFactory.NewScenarioRunsClient() *ScenarioRunsClient`
- New function `*ClientFactory.NewScenariosClient() *ScenariosClient`
- New function `*ClientFactory.NewWorkspacesClient() *WorkspacesClient`
- New function `NewDiscoveredResourcesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DiscoveredResourcesClient, error)`
- New function `*DiscoveredResourcesClient.Get(ctx context.Context, resourceGroupName string, workspaceName string, discoveredResourceName string, options *DiscoveredResourcesClientGetOptions) (DiscoveredResourcesClientGetResponse, error)`
- New function `*DiscoveredResourcesClient.NewListByWorkspacePager(resourceGroupName string, workspaceName string, options *DiscoveredResourcesClientListByWorkspaceOptions) *runtime.Pager[DiscoveredResourcesClientListByWorkspaceResponse]`
- New function `NewPrivateAccessesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PrivateAccessesClient, error)`
- New function `*PrivateAccessesClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, privateAccessName string, resource PrivateAccess, options *PrivateAccessesClientBeginCreateOrUpdateOptions) (*runtime.Poller[PrivateAccessesClientCreateOrUpdateResponse], error)`
- New function `*PrivateAccessesClient.BeginDelete(ctx context.Context, resourceGroupName string, privateAccessName string, options *PrivateAccessesClientBeginDeleteOptions) (*runtime.Poller[PrivateAccessesClientDeleteResponse], error)`
- New function `*PrivateAccessesClient.BeginDeleteAPrivateEndpointConnection(ctx context.Context, resourceGroupName string, privateAccessName string, privateEndpointConnectionName string, options *PrivateAccessesClientBeginDeleteAPrivateEndpointConnectionOptions) (*runtime.Poller[PrivateAccessesClientDeleteAPrivateEndpointConnectionResponse], error)`
- New function `*PrivateAccessesClient.Get(ctx context.Context, resourceGroupName string, privateAccessName string, options *PrivateAccessesClientGetOptions) (PrivateAccessesClientGetResponse, error)`
- New function `*PrivateAccessesClient.GetAPrivateEndpointConnection(ctx context.Context, resourceGroupName string, privateAccessName string, privateEndpointConnectionName string, options *PrivateAccessesClientGetAPrivateEndpointConnectionOptions) (PrivateAccessesClientGetAPrivateEndpointConnectionResponse, error)`
- New function `*PrivateAccessesClient.GetPrivateLinkResources(ctx context.Context, resourceGroupName string, privateAccessName string, options *PrivateAccessesClientGetPrivateLinkResourcesOptions) (PrivateAccessesClientGetPrivateLinkResourcesResponse, error)`
- New function `*PrivateAccessesClient.NewListAllPager(options *PrivateAccessesClientListAllOptions) *runtime.Pager[PrivateAccessesClientListAllResponse]`
- New function `*PrivateAccessesClient.NewListPager(resourceGroupName string, options *PrivateAccessesClientListOptions) *runtime.Pager[PrivateAccessesClientListResponse]`
- New function `*PrivateAccessesClient.NewListPrivateEndpointConnectionsPager(resourceGroupName string, privateAccessName string, options *PrivateAccessesClientListPrivateEndpointConnectionsOptions) *runtime.Pager[PrivateAccessesClientListPrivateEndpointConnectionsResponse]`
- New function `*PrivateAccessesClient.BeginUpdate(ctx context.Context, resourceGroupName string, privateAccessName string, properties PrivateAccessPatch, options *PrivateAccessesClientBeginUpdateOptions) (*runtime.Poller[PrivateAccessesClientUpdateResponse], error)`
- New function `NewScenarioConfigurationsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ScenarioConfigurationsClient, error)`
- New function `*ScenarioConfigurationsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, scenarioName string, scenarioConfigurationName string, resource ScenarioConfiguration, options *ScenarioConfigurationsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ScenarioConfigurationsClientCreateOrUpdateResponse], error)`
- New function `*ScenarioConfigurationsClient.BeginDelete(ctx context.Context, resourceGroupName string, workspaceName string, scenarioName string, scenarioConfigurationName string, options *ScenarioConfigurationsClientBeginDeleteOptions) (*runtime.Poller[ScenarioConfigurationsClientDeleteResponse], error)`
- New function `*ScenarioConfigurationsClient.Execute(ctx context.Context, resourceGroupName string, workspaceName string, scenarioName string, scenarioConfigurationName string, options *ScenarioConfigurationsClientExecuteOptions) (ScenarioConfigurationsClientExecuteResponse, error)`
- New function `*ScenarioConfigurationsClient.BeginFixResourcePermissions(ctx context.Context, resourceGroupName string, workspaceName string, scenarioName string, scenarioConfigurationName string, options *ScenarioConfigurationsClientBeginFixResourcePermissionsOptions) (*runtime.Poller[ScenarioConfigurationsClientFixResourcePermissionsResponse], error)`
- New function `*ScenarioConfigurationsClient.Get(ctx context.Context, resourceGroupName string, workspaceName string, scenarioName string, scenarioConfigurationName string, options *ScenarioConfigurationsClientGetOptions) (ScenarioConfigurationsClientGetResponse, error)`
- New function `*ScenarioConfigurationsClient.NewListAllPager(resourceGroupName string, workspaceName string, scenarioName string, options *ScenarioConfigurationsClientListAllOptions) *runtime.Pager[ScenarioConfigurationsClientListAllResponse]`
- New function `*ScenarioConfigurationsClient.BeginValidate(ctx context.Context, resourceGroupName string, workspaceName string, scenarioName string, scenarioConfigurationName string, options *ScenarioConfigurationsClientBeginValidateOptions) (*runtime.Poller[ScenarioConfigurationsClientValidateResponse], error)`
- New function `NewScenarioRunsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ScenarioRunsClient, error)`
- New function `*ScenarioRunsClient.Cancel(ctx context.Context, resourceGroupName string, workspaceName string, scenarioName string, runID string, options *ScenarioRunsClientCancelOptions) (ScenarioRunsClientCancelResponse, error)`
- New function `*ScenarioRunsClient.Get(ctx context.Context, resourceGroupName string, workspaceName string, scenarioName string, runID string, options *ScenarioRunsClientGetOptions) (ScenarioRunsClientGetResponse, error)`
- New function `*ScenarioRunsClient.NewListAllPager(resourceGroupName string, workspaceName string, scenarioName string, options *ScenarioRunsClientListAllOptions) *runtime.Pager[ScenarioRunsClientListAllResponse]`
- New function `NewScenariosClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ScenariosClient, error)`
- New function `*ScenariosClient.CreateOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, scenarioName string, resource Scenario, options *ScenariosClientCreateOrUpdateOptions) (ScenariosClientCreateOrUpdateResponse, error)`
- New function `*ScenariosClient.Delete(ctx context.Context, resourceGroupName string, workspaceName string, scenarioName string, options *ScenariosClientDeleteOptions) (ScenariosClientDeleteResponse, error)`
- New function `*ScenariosClient.Get(ctx context.Context, resourceGroupName string, workspaceName string, scenarioName string, options *ScenariosClientGetOptions) (ScenariosClientGetResponse, error)`
- New function `*ScenariosClient.NewListAllPager(resourceGroupName string, workspaceName string, options *ScenariosClientListAllOptions) *runtime.Pager[ScenariosClientListAllResponse]`
- New function `NewWorkspacesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*WorkspacesClient, error)`
- New function `*WorkspacesClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, resource Workspace, options *WorkspacesClientBeginCreateOrUpdateOptions) (*runtime.Poller[WorkspacesClientCreateOrUpdateResponse], error)`
- New function `*WorkspacesClient.BeginDelete(ctx context.Context, resourceGroupName string, workspaceName string, options *WorkspacesClientBeginDeleteOptions) (*runtime.Poller[WorkspacesClientDeleteResponse], error)`
- New function `*WorkspacesClient.Get(ctx context.Context, resourceGroupName string, workspaceName string, options *WorkspacesClientGetOptions) (WorkspacesClientGetResponse, error)`
- New function `*WorkspacesClient.NewListAllPager(options *WorkspacesClientListAllOptions) *runtime.Pager[WorkspacesClientListAllResponse]`
- New function `*WorkspacesClient.NewListPager(resourceGroupName string, options *WorkspacesClientListOptions) *runtime.Pager[WorkspacesClientListResponse]`
- New function `*WorkspacesClient.BeginRefreshRecommendations(ctx context.Context, resourceGroupName string, workspaceName string, options *WorkspacesClientBeginRefreshRecommendationsOptions) (*runtime.Poller[WorkspacesClientRefreshRecommendationsResponse], error)`
- New function `*WorkspacesClient.BeginUpdate(ctx context.Context, resourceGroupName string, workspaceName string, properties WorkspaceUpdate, options *WorkspacesClientBeginUpdateOptions) (*runtime.Poller[WorkspacesClientUpdateResponse], error)`
- New struct `Action`
- New struct `ActionDependency`
- New struct `ActionListResult`
- New struct `ActionProperties`
- New struct `ActionPropertiesParametersSchema`
- New struct `ActionSupportedTargetType`
- New struct `ActionVersion`
- New struct `ActionVersionListResult`
- New struct `ConfigurationExclusions`
- New struct `ConfigurationFilters`
- New struct `CustomerDataStorageProperties`
- New struct `DiscoveredResource`
- New struct `DiscoveredResourceListResult`
- New struct `DiscoveredResourceProperties`
- New struct `EntraIdentity`
- New struct `ExternalResource`
- New struct `FixResourcePermissionsRequest`
- New struct `OperationError`
- New struct `PermissionError`
- New struct `PermissionsFix`
- New struct `PermissionsFixProperties`
- New struct `PermissionsFixSummary`
- New struct `PhysicalToLogicalZoneMapping`
- New struct `PrivateAccess`
- New struct `PrivateAccessListResult`
- New struct `PrivateAccessPatch`
- New struct `PrivateAccessProperties`
- New struct `PrivateEndpoint`
- New struct `PrivateEndpointConnection`
- New struct `PrivateEndpointConnectionListResult`
- New struct `PrivateEndpointConnectionProperties`
- New struct `PrivateLinkResource`
- New struct `PrivateLinkResourceListResult`
- New struct `PrivateLinkResourceProperties`
- New struct `PrivateLinkServiceConnectionState`
- New struct `Recommendation`
- New struct `ResourceStateError`
- New struct `RoleAssignmentError`
- New struct `RoleAssignmentResult`
- New struct `RunAfter`
- New struct `Scenario`
- New struct `ScenarioAction`
- New struct `ScenarioConfiguration`
- New struct `ScenarioConfigurationListResult`
- New struct `ScenarioConfigurationProperties`
- New struct `ScenarioErrors`
- New struct `ScenarioEvaluationResultItem`
- New struct `ScenarioListResult`
- New struct `ScenarioParameter`
- New struct `ScenarioProperties`
- New struct `ScenarioRun`
- New struct `ScenarioRunListResult`
- New struct `ScenarioRunProperties`
- New struct `ScenarioRunResource`
- New struct `ScenarioRunSummaryAction`
- New struct `Validation`
- New struct `ValidationProperties`
- New struct `Workspace`
- New struct `WorkspaceEvaluation`
- New struct `WorkspaceEvaluationProperties`
- New struct `WorkspaceListResult`
- New struct `WorkspaceProperties`
- New struct `WorkspaceUpdate`
- New struct `ZoneResolutionInfo`
- New struct `ZoneResolutionMapping`
- New field `ProvisioningState` in struct `CapabilityProperties`
- New field `ProvisioningState` in struct `ExperimentExecutionDetailsProperties`
- New field `ProvisioningState` in struct `ExperimentExecutionProperties`
- New field `CustomerDataStorage` in struct `ExperimentProperties`


## 2.0.0 (2025-06-05)
### Breaking Changes

- Type of `ContinuousAction.Type` has been changed from `*string` to `*ExperimentActionType`
- Type of `DelayAction.Type` has been changed from `*string` to `*ExperimentActionType`
- Type of `DiscreteAction.Type` has been changed from `*string` to `*ExperimentActionType`
- Type of `Experiment.Identity` has been changed from `*ResourceIdentity` to `*ManagedServiceIdentity`
- Type of `ExperimentAction.Type` has been changed from `*string` to `*ExperimentActionType`
- Type of `ExperimentUpdate.Identity` has been changed from `*ResourceIdentity` to `*ManagedServiceIdentity`
- Enum `ResourceIdentityType` has been removed
- Function `*OperationsClient.NewListAllPager` has been removed
- Struct `OperationStatus` has been removed
- Struct `ResourceIdentity` has been removed
- Field `OperationStatus` of struct `OperationStatusesClientGetResponse` has been removed

### Features Added

- New enum type `ExperimentActionType` with values `ExperimentActionTypeContinuous`, `ExperimentActionTypeDelay`, `ExperimentActionTypeDiscrete`
- New enum type `ManagedServiceIdentityType` with values `ManagedServiceIdentityTypeNone`, `ManagedServiceIdentityTypeSystemAssigned`, `ManagedServiceIdentityTypeSystemAssignedUserAssigned`, `ManagedServiceIdentityTypeUserAssigned`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New struct `ManagedServiceIdentity`
- New struct `OperationStatusResult`
- New field `RequiredAzureRoleDefinitionIDs` in struct `CapabilityTypeProperties`
- New field `SystemData` in struct `ExperimentExecution`
- New anonymous field `OperationStatusResult` in struct `OperationStatusesClientGetResponse`


## 1.1.0 (2024-03-22)
### Features Added

- New field `Tags` in struct `ExperimentUpdate`


## 1.0.0 (2023-11-24)
### Breaking Changes

- Type of `ExperimentProperties.Selectors` has been changed from `[]SelectorClassification` to `[]TargetSelectorClassification`
- Type of `ExperimentProperties.Steps` has been changed from `[]*Step` to `[]*ExperimentStep`
- Function `*Action.GetAction` has been removed
- Function `*ContinuousAction.GetAction` has been removed
- Function `*DelayAction.GetAction` has been removed
- Function `*DiscreteAction.GetAction` has been removed
- Function `*ExperimentsClient.GetExecutionDetails` has been removed
- Function `*ExperimentsClient.GetStatus` has been removed
- Function `*ExperimentsClient.NewListAllStatusesPager` has been removed
- Function `*ExperimentsClient.NewListExecutionDetailsPager` has been removed
- Function `*Filter.GetFilter` has been removed
- Function `*ListSelector.GetSelector` has been removed
- Function `*QuerySelector.GetSelector` has been removed
- Function `*Selector.GetSelector` has been removed
- Function `*SimpleFilter.GetFilter` has been removed
- Operation `*ExperimentsClient.Cancel` has been changed to LRO, use `*ExperimentsClient.BeginCancel` instead.
- Operation `*ExperimentsClient.CreateOrUpdate` has been changed to LRO, use `*ExperimentsClient.BeginCreateOrUpdate` instead.
- Operation `*ExperimentsClient.Delete` has been changed to LRO, use `*ExperimentsClient.BeginDelete` instead.
- Operation `*ExperimentsClient.Start` has been changed to LRO, use `*ExperimentsClient.BeginStart` instead.
- Operation `*ExperimentsClient.Update` has been changed to LRO, use `*ExperimentsClient.BeginUpdate` instead.
- Struct `Branch` has been removed
- Struct `ExperimentCancelOperationResult` has been removed
- Struct `ExperimentExecutionDetailsListResult` has been removed
- Struct `ExperimentStartOperationResult` has been removed
- Struct `ExperimentStatus` has been removed
- Struct `ExperimentStatusListResult` has been removed
- Struct `ExperimentStatusProperties` has been removed
- Struct `ListSelector` has been removed
- Struct `QuerySelector` has been removed
- Struct `SimpleFilter` has been removed
- Struct `SimpleFilterParameters` has been removed
- Struct `Step` has been removed
- Field `CreatedDateTime`, `ExperimentID`, `LastActionDateTime`, `StartDateTime`, `StopDateTime` of struct `ExperimentExecutionDetailsProperties` has been removed
- Field `StartOnCreation` of struct `ExperimentProperties` has been removed

### Features Added

- Support for test fakes and OpenTelemetry trace spans.
- New enum type `ProvisioningState` with values `ProvisioningStateCanceled`, `ProvisioningStateCreating`, `ProvisioningStateDeleting`, `ProvisioningStateFailed`, `ProvisioningStateSucceeded`, `ProvisioningStateUpdating`
- New function `*ClientFactory.NewOperationStatusesClient() *OperationStatusesClient`
- New function `*ContinuousAction.GetExperimentAction() *ExperimentAction`
- New function `*DelayAction.GetExperimentAction() *ExperimentAction`
- New function `*DiscreteAction.GetExperimentAction() *ExperimentAction`
- New function `*ExperimentAction.GetExperimentAction() *ExperimentAction`
- New function `*ExperimentsClient.ExecutionDetails(context.Context, string, string, string, *ExperimentsClientExecutionDetailsOptions) (ExperimentsClientExecutionDetailsResponse, error)`
- New function `*ExperimentsClient.GetExecution(context.Context, string, string, string, *ExperimentsClientGetExecutionOptions) (ExperimentsClientGetExecutionResponse, error)`
- New function `*ExperimentsClient.NewListAllExecutionsPager(string, string, *ExperimentsClientListAllExecutionsOptions) *runtime.Pager[ExperimentsClientListAllExecutionsResponse]`
- New function `NewOperationStatusesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*OperationStatusesClient, error)`
- New function `*OperationStatusesClient.Get(context.Context, string, string, *OperationStatusesClientGetOptions) (OperationStatusesClientGetResponse, error)`
- New function `*TargetFilter.GetTargetFilter() *TargetFilter`
- New function `*TargetListSelector.GetTargetSelector() *TargetSelector`
- New function `*TargetQuerySelector.GetTargetSelector() *TargetSelector`
- New function `*TargetSelector.GetTargetSelector() *TargetSelector`
- New function `*TargetSimpleFilter.GetTargetFilter() *TargetFilter`
- New struct `ExperimentBranch`
- New struct `ExperimentExecution`
- New struct `ExperimentExecutionListResult`
- New struct `ExperimentExecutionProperties`
- New struct `ExperimentStep`
- New struct `OperationStatus`
- New struct `TargetListSelector`
- New struct `TargetQuerySelector`
- New struct `TargetSimpleFilter`
- New struct `TargetSimpleFilterParameters`
- New field `LastActionAt`, `StartedAt`, `StoppedAt` in struct `ExperimentExecutionDetailsProperties`
- New field `ProvisioningState` in struct `ExperimentProperties`


## 0.7.0 (2023-08-25)
### Breaking Changes

- Type of `ExperimentProperties.Selectors` has been changed from `[]*Selector` to `[]SelectorClassification`
- Type of `TargetReference.Type` has been changed from `*string` to `*TargetReferenceType`
- `SelectorTypePercent`, `SelectorTypeRandom`, `SelectorTypeTag` from enum `SelectorType` has been removed
- Operation `*ExperimentsClient.BeginCancel` has been changed to non-LRO, use `*ExperimentsClient.Cancel` instead.
- Operation `*ExperimentsClient.BeginCreateOrUpdate` has been changed to non-LRO, use `*ExperimentsClient.CreateOrUpdate` instead.
- Field `Targets` of struct `Selector` has been removed

### Features Added

- New value `ResourceIdentityTypeUserAssigned` added to enum type `ResourceIdentityType`
- New value `SelectorTypeQuery` added to enum type `SelectorType`
- New enum type `FilterType` with values `FilterTypeSimple`
- New enum type `TargetReferenceType` with values `TargetReferenceTypeChaosTarget`
- New function `*ExperimentsClient.Update(context.Context, string, string, ExperimentUpdate, *ExperimentsClientUpdateOptions) (ExperimentsClientUpdateResponse, error)`
- New function `*Filter.GetFilter() *Filter`
- New function `*ListSelector.GetSelector() *Selector`
- New function `*QuerySelector.GetSelector() *Selector`
- New function `*Selector.GetSelector() *Selector`
- New function `*SimpleFilter.GetFilter() *Filter`
- New struct `CapabilityTypePropertiesRuntimeProperties`
- New struct `ExperimentUpdate`
- New struct `ListSelector`
- New struct `QuerySelector`
- New struct `SimpleFilter`
- New struct `SimpleFilterParameters`
- New struct `UserAssignedIdentity`
- New field `AzureRbacActions`, `AzureRbacDataActions`, `Kind`, `RuntimeProperties` in struct `CapabilityTypeProperties`
- New field `UserAssignedIdentities` in struct `ResourceIdentity`


## 0.6.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 0.6.0 (2023-03-28)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.5.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/chaos/armchaos` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
