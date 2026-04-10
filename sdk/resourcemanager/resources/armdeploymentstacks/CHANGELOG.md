# Release History

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