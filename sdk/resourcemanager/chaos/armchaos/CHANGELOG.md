# Release History

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
