# Release History

## 3.0.0 (2026-04-03)
### Breaking Changes

- Enum `ActionType` has been removed
- Enum `ApplicationArtifactName` has been removed
- Enum `ApplicationArtifactType` has been removed
- Enum `ApplicationDefinitionArtifactName` has been removed
- Enum `ApplicationLockLevel` has been removed
- Enum `ApplicationManagementMode` has been removed
- Enum `CreatedByType` has been removed
- Enum `DeploymentMode` has been removed
- Enum `JitApprovalMode` has been removed
- Enum `JitApproverType` has been removed
- Enum `JitRequestState` has been removed
- Enum `JitSchedulingType` has been removed
- Enum `Origin` has been removed
- Enum `ProvisioningState` has been removed
- Enum `ResourceIdentityType` has been removed
- Enum `Status` has been removed
- Enum `Substatus` has been removed
- Function `NewApplicationClient` has been removed
- Function `*ApplicationClient.NewListOperationsPager` has been removed
- Function `NewApplicationDefinitionsClient` has been removed
- Function `*ApplicationDefinitionsClient.CreateOrUpdate` has been removed
- Function `*ApplicationDefinitionsClient.CreateOrUpdateByID` has been removed
- Function `*ApplicationDefinitionsClient.Delete` has been removed
- Function `*ApplicationDefinitionsClient.DeleteByID` has been removed
- Function `*ApplicationDefinitionsClient.Get` has been removed
- Function `*ApplicationDefinitionsClient.GetByID` has been removed
- Function `*ApplicationDefinitionsClient.NewListByResourceGroupPager` has been removed
- Function `*ApplicationDefinitionsClient.NewListBySubscriptionPager` has been removed
- Function `*ApplicationDefinitionsClient.Update` has been removed
- Function `*ApplicationDefinitionsClient.UpdateByID` has been removed
- Function `NewApplicationsClient` has been removed
- Function `*ApplicationsClient.BeginCreateOrUpdate` has been removed
- Function `*ApplicationsClient.BeginCreateOrUpdateByID` has been removed
- Function `*ApplicationsClient.BeginDelete` has been removed
- Function `*ApplicationsClient.BeginDeleteByID` has been removed
- Function `*ApplicationsClient.Get` has been removed
- Function `*ApplicationsClient.GetByID` has been removed
- Function `*ApplicationsClient.ListAllowedUpgradePlans` has been removed
- Function `*ApplicationsClient.NewListByResourceGroupPager` has been removed
- Function `*ApplicationsClient.NewListBySubscriptionPager` has been removed
- Function `*ApplicationsClient.ListTokens` has been removed
- Function `*ApplicationsClient.BeginRefreshPermissions` has been removed
- Function `*ApplicationsClient.BeginUpdate` has been removed
- Function `*ApplicationsClient.BeginUpdateAccess` has been removed
- Function `*ApplicationsClient.BeginUpdateByID` has been removed
- Function `NewClientFactory` has been removed
- Function `*ClientFactory.NewApplicationClient` has been removed
- Function `*ClientFactory.NewApplicationDefinitionsClient` has been removed
- Function `*ClientFactory.NewApplicationsClient` has been removed
- Function `*ClientFactory.NewJitRequestsClient` has been removed
- Function `NewJitRequestsClient` has been removed
- Function `*JitRequestsClient.BeginCreateOrUpdate` has been removed
- Function `*JitRequestsClient.Delete` has been removed
- Function `*JitRequestsClient.Get` has been removed
- Function `*JitRequestsClient.ListByResourceGroup` has been removed
- Function `*JitRequestsClient.ListBySubscription` has been removed
- Function `*JitRequestsClient.Update` has been removed
- Struct `AllowedUpgradePlansResult` has been removed
- Struct `Application` has been removed
- Struct `ApplicationArtifact` has been removed
- Struct `ApplicationAuthorization` has been removed
- Struct `ApplicationBillingDetailsDefinition` has been removed
- Struct `ApplicationClientDetails` has been removed
- Struct `ApplicationDefinition` has been removed
- Struct `ApplicationDefinitionArtifact` has been removed
- Struct `ApplicationDefinitionListResult` has been removed
- Struct `ApplicationDefinitionPatchable` has been removed
- Struct `ApplicationDefinitionProperties` has been removed
- Struct `ApplicationDeploymentPolicy` has been removed
- Struct `ApplicationJitAccessPolicy` has been removed
- Struct `ApplicationListResult` has been removed
- Struct `ApplicationManagementPolicy` has been removed
- Struct `ApplicationNotificationEndpoint` has been removed
- Struct `ApplicationNotificationPolicy` has been removed
- Struct `ApplicationPackageContact` has been removed
- Struct `ApplicationPackageLockingPolicyDefinition` has been removed
- Struct `ApplicationPackageSupportUrls` has been removed
- Struct `ApplicationPatchable` has been removed
- Struct `ApplicationPolicy` has been removed
- Struct `ApplicationProperties` has been removed
- Struct `ClientFactory` has been removed
- Struct `ErrorAdditionalInfo` has been removed
- Struct `ErrorDetail` has been removed
- Struct `ErrorResponse` has been removed
- Struct `GenericResource` has been removed
- Struct `Identity` has been removed
- Struct `JitApproverDefinition` has been removed
- Struct `JitAuthorizationPolicies` has been removed
- Struct `JitRequestDefinition` has been removed
- Struct `JitRequestDefinitionListResult` has been removed
- Struct `JitRequestMetadata` has been removed
- Struct `JitRequestPatchable` has been removed
- Struct `JitRequestProperties` has been removed
- Struct `JitSchedulingPolicy` has been removed
- Struct `ListTokenRequest` has been removed
- Struct `ManagedIdentityToken` has been removed
- Struct `ManagedIdentityTokenResult` has been removed
- Struct `Operation` has been removed
- Struct `OperationDisplay` has been removed
- Struct `OperationListResult` has been removed
- Struct `Plan` has been removed
- Struct `PlanPatchable` has been removed
- Struct `Resource` has been removed
- Struct `SKU` has been removed
- Struct `SystemData` has been removed
- Struct `UpdateAccessDefinition` has been removed
- Struct `UserAssignedResourceIdentity` has been removed


## 2.1.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 2.0.0 (2023-09-22)
### Breaking Changes

- Operation `*ApplicationDefinitionsClient.BeginCreateOrUpdate` has been changed to non-LRO, use `*ApplicationDefinitionsClient.CreateOrUpdate` instead.
- Operation `*ApplicationDefinitionsClient.BeginDelete` has been changed to non-LRO, use `*ApplicationDefinitionsClient.Delete` instead.
- Operation `*ApplicationsClient.Update` has been changed to LRO, use `*ApplicationsClient.BeginUpdate` instead.
- Struct `ApplicationPropertiesPatchable` has been removed
- Field `ProvisioningState` of struct `ApplicationDefinitionProperties` has been removed

### Features Added

- Function `ApplicationsClient.Get` no longer throws an exception when the response is `http.StatusNotFound`
- New enum type `Status` with values `StatusElevate`, `StatusNotSpecified`, `StatusRemove`
- New enum type `Substatus` with values `SubstatusApproved`, `SubstatusDenied`, `SubstatusExpired`, `SubstatusFailed`, `SubstatusNotSpecified`, `SubstatusTimeout`
- New function `*ApplicationDefinitionsClient.CreateOrUpdateByID(context.Context, string, string, ApplicationDefinition, *ApplicationDefinitionsClientCreateOrUpdateByIDOptions) (ApplicationDefinitionsClientCreateOrUpdateByIDResponse, error)`
- New function `*ApplicationDefinitionsClient.DeleteByID(context.Context, string, string, *ApplicationDefinitionsClientDeleteByIDOptions) (ApplicationDefinitionsClientDeleteByIDResponse, error)`
- New function `*ApplicationDefinitionsClient.GetByID(context.Context, string, string, *ApplicationDefinitionsClientGetByIDOptions) (ApplicationDefinitionsClientGetByIDResponse, error)`
- New function `*ApplicationDefinitionsClient.UpdateByID(context.Context, string, string, ApplicationDefinitionPatchable, *ApplicationDefinitionsClientUpdateByIDOptions) (ApplicationDefinitionsClientUpdateByIDResponse, error)`
- New function `*ApplicationsClient.BeginCreateOrUpdateByID(context.Context, string, Application, *ApplicationsClientBeginCreateOrUpdateByIDOptions) (*runtime.Poller[ApplicationsClientCreateOrUpdateByIDResponse], error)`
- New function `*ApplicationsClient.BeginDeleteByID(context.Context, string, *ApplicationsClientBeginDeleteByIDOptions) (*runtime.Poller[ApplicationsClientDeleteByIDResponse], error)`
- New function `*ApplicationsClient.GetByID(context.Context, string, *ApplicationsClientGetByIDOptions) (ApplicationsClientGetByIDResponse, error)`
- New function `*ApplicationsClient.ListTokens(context.Context, string, string, ListTokenRequest, *ApplicationsClientListTokensOptions) (ApplicationsClientListTokensResponse, error)`
- New function `*ApplicationsClient.BeginUpdateAccess(context.Context, string, string, UpdateAccessDefinition, *ApplicationsClientBeginUpdateAccessOptions) (*runtime.Poller[ApplicationsClientUpdateAccessResponse], error)`
- New function `*ApplicationsClient.BeginUpdateByID(context.Context, string, *ApplicationsClientBeginUpdateByIDOptions) (*runtime.Poller[ApplicationsClientUpdateByIDResponse], error)`
- New struct `AllowedUpgradePlansResult`
- New struct `JitRequestMetadata`
- New struct `ListTokenRequest`
- New struct `ManagedIdentityToken`
- New struct `ManagedIdentityTokenResult`
- New struct `UpdateAccessDefinition`
- New anonymous field `AllowedUpgradePlansResult` in struct `ApplicationsClientListAllowedUpgradePlansResponse`


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/solutions/armmanagedapplications` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).