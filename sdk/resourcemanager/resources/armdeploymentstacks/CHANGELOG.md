# Release History

## 2.0.0 (2026-01-27)
### Breaking Changes

- Type of `ClientBeginDeleteAtManagementGroupOptions.UnmanageActionManagementGroups` has been changed from `*UnmanageActionManagementGroupMode` to `*DeploymentStacksDeleteDetachEnum`
- Type of `ClientBeginDeleteAtManagementGroupOptions.UnmanageActionResourceGroups` has been changed from `*UnmanageActionResourceGroupMode` to `*DeploymentStacksDeleteDetachEnum`
- Type of `ClientBeginDeleteAtManagementGroupOptions.UnmanageActionResources` has been changed from `*UnmanageActionResourceMode` to `*DeploymentStacksDeleteDetachEnum`
- Type of `ClientBeginDeleteAtResourceGroupOptions.UnmanageActionManagementGroups` has been changed from `*UnmanageActionManagementGroupMode` to `*DeploymentStacksDeleteDetachEnum`
- Type of `ClientBeginDeleteAtResourceGroupOptions.UnmanageActionResourceGroups` has been changed from `*UnmanageActionResourceGroupMode` to `*DeploymentStacksDeleteDetachEnum`
- Type of `ClientBeginDeleteAtResourceGroupOptions.UnmanageActionResources` has been changed from `*UnmanageActionResourceMode` to `*DeploymentStacksDeleteDetachEnum`
- Type of `ClientBeginDeleteAtSubscriptionOptions.UnmanageActionManagementGroups` has been changed from `*UnmanageActionManagementGroupMode` to `*DeploymentStacksDeleteDetachEnum`
- Type of `ClientBeginDeleteAtSubscriptionOptions.UnmanageActionResourceGroups` has been changed from `*UnmanageActionResourceGroupMode` to `*DeploymentStacksDeleteDetachEnum`
- Type of `ClientBeginDeleteAtSubscriptionOptions.UnmanageActionResources` has been changed from `*UnmanageActionResourceMode` to `*DeploymentStacksDeleteDetachEnum`
- Enum `UnmanageActionManagementGroupMode` has been removed
- Enum `UnmanageActionResourceGroupMode` has been removed
- Enum `UnmanageActionResourceMode` has been removed

### Features Added

- New value `DenyStatusModeUnknown` added to enum type `DenyStatusMode`
- New value `DeploymentStackProvisioningStateInitializing`, `DeploymentStackProvisioningStateRunning` added to enum type `DeploymentStackProvisioningState`
- New enum type `DeploymentStacksDiagnosticLevel` with values `DeploymentStacksDiagnosticLevelError`, `DeploymentStacksDiagnosticLevelInfo`, `DeploymentStacksDiagnosticLevelWarning`
- New enum type `DeploymentStacksManagementStatus` with values `DeploymentStacksManagementStatusManaged`, `DeploymentStacksManagementStatusUnknown`, `DeploymentStacksManagementStatusUnmanaged`
- New enum type `DeploymentStacksResourcesWithoutDeleteSupportEnum` with values `DeploymentStacksResourcesWithoutDeleteSupportEnumDetach`, `DeploymentStacksResourcesWithoutDeleteSupportEnumFail`
- New enum type `DeploymentStacksWhatIfChangeCertainty` with values `DeploymentStacksWhatIfChangeCertaintyDefinite`, `DeploymentStacksWhatIfChangeCertaintyPotential`
- New enum type `DeploymentStacksWhatIfChangeType` with values `DeploymentStacksWhatIfChangeTypeCreate`, `DeploymentStacksWhatIfChangeTypeDelete`, `DeploymentStacksWhatIfChangeTypeDetach`, `DeploymentStacksWhatIfChangeTypeModify`, `DeploymentStacksWhatIfChangeTypeNoChange`, `DeploymentStacksWhatIfChangeTypeUnsupported`
- New enum type `DeploymentStacksWhatIfPropertyChangeType` with values `DeploymentStacksWhatIfPropertyChangeTypeArray`, `DeploymentStacksWhatIfPropertyChangeTypeCreate`, `DeploymentStacksWhatIfPropertyChangeTypeDelete`, `DeploymentStacksWhatIfPropertyChangeTypeModify`, `DeploymentStacksWhatIfPropertyChangeTypeNoEffect`
- New enum type `ValidationLevel` with values `ValidationLevelProvider`, `ValidationLevelProviderNoRbac`, `ValidationLevelTemplate`
- New function `*ClientFactory.NewWhatIfResultsAtManagementGroupClient() *WhatIfResultsAtManagementGroupClient`
- New function `*ClientFactory.NewWhatIfResultsAtResourceGroupClient() *WhatIfResultsAtResourceGroupClient`
- New function `*ClientFactory.NewWhatIfResultsAtSubscriptionClient() *WhatIfResultsAtSubscriptionClient`
- New function `NewWhatIfResultsAtManagementGroupClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*WhatIfResultsAtManagementGroupClient, error)`
- New function `*WhatIfResultsAtManagementGroupClient.BeginCreateOrUpdate(ctx context.Context, managementGroupID string, deploymentStacksWhatIfResultName string, resource WhatIfResult, options *WhatIfResultsAtManagementGroupClientBeginCreateOrUpdateOptions) (*runtime.Poller[WhatIfResultsAtManagementGroupClientCreateOrUpdateResponse], error)`
- New function `*WhatIfResultsAtManagementGroupClient.Delete(ctx context.Context, managementGroupID string, deploymentStacksWhatIfResultName string, options *WhatIfResultsAtManagementGroupClientDeleteOptions) (WhatIfResultsAtManagementGroupClientDeleteResponse, error)`
- New function `*WhatIfResultsAtManagementGroupClient.Get(ctx context.Context, managementGroupID string, deploymentStacksWhatIfResultName string, options *WhatIfResultsAtManagementGroupClientGetOptions) (WhatIfResultsAtManagementGroupClientGetResponse, error)`
- New function `*WhatIfResultsAtManagementGroupClient.NewListPager(managementGroupID string, options *WhatIfResultsAtManagementGroupClientListOptions) *runtime.Pager[WhatIfResultsAtManagementGroupClientListResponse]`
- New function `*WhatIfResultsAtManagementGroupClient.BeginWhatIf(ctx context.Context, managementGroupID string, deploymentStacksWhatIfResultName string, options *WhatIfResultsAtManagementGroupClientBeginWhatIfOptions) (*runtime.Poller[WhatIfResultsAtManagementGroupClientWhatIfResponse], error)`
- New function `NewWhatIfResultsAtResourceGroupClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*WhatIfResultsAtResourceGroupClient, error)`
- New function `*WhatIfResultsAtResourceGroupClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, deploymentStacksWhatIfResultName string, resource WhatIfResult, options *WhatIfResultsAtResourceGroupClientBeginCreateOrUpdateOptions) (*runtime.Poller[WhatIfResultsAtResourceGroupClientCreateOrUpdateResponse], error)`
- New function `*WhatIfResultsAtResourceGroupClient.Delete(ctx context.Context, resourceGroupName string, deploymentStacksWhatIfResultName string, options *WhatIfResultsAtResourceGroupClientDeleteOptions) (WhatIfResultsAtResourceGroupClientDeleteResponse, error)`
- New function `*WhatIfResultsAtResourceGroupClient.Get(ctx context.Context, resourceGroupName string, deploymentStacksWhatIfResultName string, options *WhatIfResultsAtResourceGroupClientGetOptions) (WhatIfResultsAtResourceGroupClientGetResponse, error)`
- New function `*WhatIfResultsAtResourceGroupClient.NewListPager(resourceGroupName string, options *WhatIfResultsAtResourceGroupClientListOptions) *runtime.Pager[WhatIfResultsAtResourceGroupClientListResponse]`
- New function `*WhatIfResultsAtResourceGroupClient.BeginWhatIf(ctx context.Context, resourceGroupName string, deploymentStacksWhatIfResultName string, options *WhatIfResultsAtResourceGroupClientBeginWhatIfOptions) (*runtime.Poller[WhatIfResultsAtResourceGroupClientWhatIfResponse], error)`
- New function `NewWhatIfResultsAtSubscriptionClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*WhatIfResultsAtSubscriptionClient, error)`
- New function `*WhatIfResultsAtSubscriptionClient.BeginCreateOrUpdate(ctx context.Context, deploymentStacksWhatIfResultName string, resource WhatIfResult, options *WhatIfResultsAtSubscriptionClientBeginCreateOrUpdateOptions) (*runtime.Poller[WhatIfResultsAtSubscriptionClientCreateOrUpdateResponse], error)`
- New function `*WhatIfResultsAtSubscriptionClient.Delete(ctx context.Context, deploymentStacksWhatIfResultName string, options *WhatIfResultsAtSubscriptionClientDeleteOptions) (WhatIfResultsAtSubscriptionClientDeleteResponse, error)`
- New function `*WhatIfResultsAtSubscriptionClient.Get(ctx context.Context, deploymentStacksWhatIfResultName string, options *WhatIfResultsAtSubscriptionClientGetOptions) (WhatIfResultsAtSubscriptionClientGetResponse, error)`
- New function `*WhatIfResultsAtSubscriptionClient.NewListPager(options *WhatIfResultsAtSubscriptionClientListOptions) *runtime.Pager[WhatIfResultsAtSubscriptionClientListResponse]`
- New function `*WhatIfResultsAtSubscriptionClient.BeginWhatIf(ctx context.Context, deploymentStacksWhatIfResultName string, options *WhatIfResultsAtSubscriptionClientBeginWhatIfOptions) (*runtime.Poller[WhatIfResultsAtSubscriptionClientWhatIfResponse], error)`
- New struct `ChangeBase`
- New struct `ChangeBaseDenyStatusMode`
- New struct `ChangeBaseDeploymentStacksManagementStatus`
- New struct `ChangeDeltaDenySettings`
- New struct `ChangeDeltaRecord`
- New struct `DeploymentExtension`
- New struct `DeploymentExtensionConfig`
- New struct `DeploymentExtensionConfigItem`
- New struct `DeploymentExternalInput`
- New struct `DeploymentExternalInputDefinition`
- New struct `Diagnostic`
- New struct `WhatIfChange`
- New struct `WhatIfPropertyChange`
- New struct `WhatIfResourceChange`
- New struct `WhatIfResult`
- New struct `WhatIfResultListResult`
- New struct `WhatIfResultProperties`
- New field `ResourcesWithoutDeleteSupport` in struct `ActionOnUnmanage`
- New field `UnmanageActionResourcesWithoutDeleteSupport` in struct `ClientBeginDeleteAtManagementGroupOptions`
- New field `UnmanageActionResourcesWithoutDeleteSupport` in struct `ClientBeginDeleteAtResourceGroupOptions`
- New field `UnmanageActionResourcesWithoutDeleteSupport` in struct `ClientBeginDeleteAtSubscriptionOptions`
- New field `Expression` in struct `DeploymentParameter`
- New field `DeploymentExtensions`, `ExtensionConfigs`, `ExternalInputDefinitions`, `ExternalInputs`, `ValidationLevel` in struct `DeploymentStackProperties`
- New field `DeploymentExtensions`, `ValidationLevel` in struct `DeploymentStackValidateProperties`
- New field `APIVersion`, `Extension`, `Identifiers`, `Type` in struct `ManagedResourceReference`
- New field `APIVersion`, `Extension`, `Identifiers`, `Type` in struct `ResourceReference`
- New field `APIVersion`, `Extension`, `Identifiers`, `Type` in struct `ResourceReferenceExtended`


## 1.0.1 (2025-07-23)
### Other Changes

- Adopt latest code gen optimization.

## 1.0.0 (2024-06-21)
### Breaking Changes

- Type of `DeploymentStackProperties.ActionOnUnmanage` has been changed from `*DeploymentStackPropertiesActionOnUnmanage` to `*ActionOnUnmanage`
- Type of `DeploymentStackProperties.Error` has been changed from `*ErrorResponse` to `*ErrorDetail`
- Type of `DeploymentStackProperties.Parameters` has been changed from `any` to `map[string]*DeploymentParameter`
- Type of `ResourceReferenceExtended.Error` has been changed from `*ErrorResponse` to `*ErrorDetail`
- `DeploymentStackProvisioningStateLocking` from enum `DeploymentStackProvisioningState` has been removed
- `ResourceStatusModeNone` from enum `ResourceStatusMode` has been removed
- Struct `DeploymentStackPropertiesActionOnUnmanage` has been removed
- Struct `ErrorResponse` has been removed

### Features Added

- New value `DeploymentStackProvisioningStateUpdatingDenyAssignments` added to enum type `DeploymentStackProvisioningState`
- New function `*Client.BeginValidateStackAtManagementGroup(context.Context, string, string, DeploymentStack, *ClientBeginValidateStackAtManagementGroupOptions) (*runtime.Poller[ClientValidateStackAtManagementGroupResponse], error)`
- New function `*Client.BeginValidateStackAtResourceGroup(context.Context, string, string, DeploymentStack, *ClientBeginValidateStackAtResourceGroupOptions) (*runtime.Poller[ClientValidateStackAtResourceGroupResponse], error)`
- New function `*Client.BeginValidateStackAtSubscription(context.Context, string, DeploymentStack, *ClientBeginValidateStackAtSubscriptionOptions) (*runtime.Poller[ClientValidateStackAtSubscriptionResponse], error)`
- New struct `ActionOnUnmanage`
- New struct `DeploymentParameter`
- New struct `DeploymentStackValidateProperties`
- New struct `DeploymentStackValidateResult`
- New struct `KeyVaultParameterReference`
- New struct `KeyVaultReference`
- New field `BypassStackOutOfSyncError` in struct `ClientBeginDeleteAtManagementGroupOptions`
- New field `BypassStackOutOfSyncError`, `UnmanageActionManagementGroups` in struct `ClientBeginDeleteAtResourceGroupOptions`
- New field `BypassStackOutOfSyncError`, `UnmanageActionManagementGroups` in struct `ClientBeginDeleteAtSubscriptionOptions`
- New field `BypassStackOutOfSyncError`, `CorrelationID` in struct `DeploymentStackProperties`


## 0.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.1.0 (2023-08-25)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armdeploymentstacks` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).