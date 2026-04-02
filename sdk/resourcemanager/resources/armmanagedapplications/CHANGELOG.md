# Release History

## 2.0.0-beta.1 (2026-04-02)
### Breaking Changes

- Type of `ApplicationArtifact.Name` has been changed from `*string` to `*ApplicationArtifactName`
- Type of `ApplicationDefinitionProperties.Artifacts` has been changed from `[]*ApplicationArtifact` to `[]*ApplicationDefinitionArtifact`
- Type of `ApplicationDefinitionProperties.Authorizations` has been changed from `[]*ApplicationProviderAuthorization` to `[]*ApplicationAuthorization`
- Type of `ApplicationDefinitionProperties.IsEnabled` has been changed from `*string` to `*bool`
- Type of `ApplicationPatchable.Properties` has been changed from `*ApplicationPropertiesPatchable` to `*ApplicationProperties`
- Type of `Identity.Type` has been changed from `*string` to `*ResourceIdentityType`
- `ProvisioningStateCreated`, `ProvisioningStateCreating`, `ProvisioningStateReady` from enum `ProvisioningState` has been removed
- Function `NewApplicationClient` has been removed
- Function `*ApplicationClient.NewListOperationsPager` has been removed
- Function `*ClientFactory.NewApplicationClient` has been removed
- Operation `*ApplicationDefinitionsClient.BeginCreateOrUpdate` has been changed to non-LRO, use `*ApplicationDefinitionsClient.CreateOrUpdate` instead.
- Operation `*ApplicationDefinitionsClient.BeginCreateOrUpdateByID` has been changed to non-LRO, use `*ApplicationDefinitionsClient.CreateOrUpdateByID` instead.
- Operation `*ApplicationDefinitionsClient.BeginDelete` has been changed to non-LRO, use `*ApplicationDefinitionsClient.Delete` instead.
- Operation `*ApplicationDefinitionsClient.BeginDeleteByID` has been changed to non-LRO, use `*ApplicationDefinitionsClient.DeleteByID` instead.
- Operation `*ApplicationsClient.Update` has been changed to LRO, use `*ApplicationsClient.BeginUpdate` instead.
- Operation `*ApplicationsClient.UpdateByID` has been changed to LRO, use `*ApplicationsClient.BeginUpdateByID` instead.
- Struct `ApplicationPropertiesPatchable` has been removed
- Struct `ApplicationProviderAuthorization` has been removed
- Struct `ErrorResponse` has been removed
- Struct `GenericResource` has been removed
- Struct `Resource` has been removed
- Field `Identity` of struct `ApplicationDefinition` has been removed

### Features Added

- New value `ApplicationArtifactTypeNotSpecified` added to enum type `ApplicationArtifactType`
- New value `ProvisioningStateNotSpecified` added to enum type `ProvisioningState`
- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `ApplicationArtifactName` with values `ApplicationArtifactNameAuthorizations`, `ApplicationArtifactNameCustomRoleDefinition`, `ApplicationArtifactNameNotSpecified`, `ApplicationArtifactNameViewDefinition`
- New enum type `ApplicationDefinitionArtifactName` with values `ApplicationDefinitionArtifactNameApplicationResourceTemplate`, `ApplicationDefinitionArtifactNameCreateUIDefinition`, `ApplicationDefinitionArtifactNameMainTemplateParameters`, `ApplicationDefinitionArtifactNameNotSpecified`
- New enum type `ApplicationManagementMode` with values `ApplicationManagementModeManaged`, `ApplicationManagementModeNotSpecified`, `ApplicationManagementModeUnmanaged`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `DeploymentMode` with values `DeploymentModeComplete`, `DeploymentModeIncremental`, `DeploymentModeNotSpecified`
- New enum type `JitApprovalMode` with values `JitApprovalModeAutoApprove`, `JitApprovalModeManualApprove`, `JitApprovalModeNotSpecified`
- New enum type `JitApproverType` with values `JitApproverTypeGroup`, `JitApproverTypeUser`
- New enum type `JitRequestState` with values `JitRequestStateApproved`, `JitRequestStateCanceled`, `JitRequestStateDenied`, `JitRequestStateExpired`, `JitRequestStateFailed`, `JitRequestStateNotSpecified`, `JitRequestStatePending`, `JitRequestStateTimeout`
- New enum type `JitSchedulingType` with values `JitSchedulingTypeNotSpecified`, `JitSchedulingTypeOnce`, `JitSchedulingTypeRecurring`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `ResourceIdentityType` with values `ResourceIdentityTypeNone`, `ResourceIdentityTypeSystemAssigned`, `ResourceIdentityTypeSystemAssignedUserAssigned`, `ResourceIdentityTypeUserAssigned`
- New enum type `Status` with values `StatusElevate`, `StatusNotSpecified`, `StatusRemove`
- New enum type `Substatus` with values `SubstatusApproved`, `SubstatusDenied`, `SubstatusExpired`, `SubstatusFailed`, `SubstatusNotSpecified`, `SubstatusTimeout`
- New function `*ApplicationDefinitionsClient.NewListBySubscriptionPager(options *ApplicationDefinitionsClientListBySubscriptionOptions) *runtime.Pager[ApplicationDefinitionsClientListBySubscriptionResponse]`
- New function `*ApplicationDefinitionsClient.Update(ctx context.Context, resourceGroupName string, applicationDefinitionName string, parameters ApplicationDefinitionPatchable, options *ApplicationDefinitionsClientUpdateOptions) (ApplicationDefinitionsClientUpdateResponse, error)`
- New function `*ApplicationDefinitionsClient.UpdateByID(ctx context.Context, resourceGroupName string, applicationDefinitionName string, parameters ApplicationDefinitionPatchable, options *ApplicationDefinitionsClientUpdateByIDOptions) (ApplicationDefinitionsClientUpdateByIDResponse, error)`
- New function `*ApplicationsClient.ListAllowedUpgradePlans(ctx context.Context, resourceGroupName string, applicationName string, options *ApplicationsClientListAllowedUpgradePlansOptions) (ApplicationsClientListAllowedUpgradePlansResponse, error)`
- New function `*ApplicationsClient.ListTokens(ctx context.Context, resourceGroupName string, applicationName string, parameters ListTokenRequest, options *ApplicationsClientListTokensOptions) (ApplicationsClientListTokensResponse, error)`
- New function `*ApplicationsClient.BeginRefreshPermissions(ctx context.Context, resourceGroupName string, applicationName string, options *ApplicationsClientBeginRefreshPermissionsOptions) (*runtime.Poller[ApplicationsClientRefreshPermissionsResponse], error)`
- New function `*ApplicationsClient.BeginUpdateAccess(ctx context.Context, resourceGroupName string, applicationName string, parameters UpdateAccessDefinition, options *ApplicationsClientBeginUpdateAccessOptions) (*runtime.Poller[ApplicationsClientUpdateAccessResponse], error)`
- New function `*ClientFactory.NewJitRequestsClient() *JitRequestsClient`
- New function `*ClientFactory.NewSolutionsClient() *SolutionsClient`
- New function `NewJitRequestsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*JitRequestsClient, error)`
- New function `*JitRequestsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, jitRequestName string, parameters JitRequestDefinition, options *JitRequestsClientBeginCreateOrUpdateOptions) (*runtime.Poller[JitRequestsClientCreateOrUpdateResponse], error)`
- New function `*JitRequestsClient.Delete(ctx context.Context, resourceGroupName string, jitRequestName string, options *JitRequestsClientDeleteOptions) (JitRequestsClientDeleteResponse, error)`
- New function `*JitRequestsClient.Get(ctx context.Context, resourceGroupName string, jitRequestName string, options *JitRequestsClientGetOptions) (JitRequestsClientGetResponse, error)`
- New function `*JitRequestsClient.ListByResourceGroup(ctx context.Context, resourceGroupName string, options *JitRequestsClientListByResourceGroupOptions) (JitRequestsClientListByResourceGroupResponse, error)`
- New function `*JitRequestsClient.ListBySubscription(ctx context.Context, options *JitRequestsClientListBySubscriptionOptions) (JitRequestsClientListBySubscriptionResponse, error)`
- New function `*JitRequestsClient.Update(ctx context.Context, resourceGroupName string, jitRequestName string, parameters JitRequestPatchable, options *JitRequestsClientUpdateOptions) (JitRequestsClientUpdateResponse, error)`
- New function `NewSolutionsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*SolutionsClient, error)`
- New function `*SolutionsClient.NewListOperationsPager(options *SolutionsClientListOperationsOptions) *runtime.Pager[SolutionsClientListOperationsResponse]`
- New function `*SolutionsClient.NewApplicationDefinitionsClient() *ApplicationDefinitionsClient`
- New function `*SolutionsClient.NewApplicationsClient() *ApplicationsClient`
- New function `*SolutionsClient.NewJitRequestsClient() *JitRequestsClient`
- New function `*SolutionsClient.PortalRegistryPackage(ctx context.Context, parameters RegistryPackagePlan, options *SolutionsClientPortalRegistryPackageOptions) (SolutionsClientPortalRegistryPackageResponse, error)`
- New struct `AllowedUpgradePlansResult`
- New struct `ApplicationAuthorization`
- New struct `ApplicationBillingDetailsDefinition`
- New struct `ApplicationClientDetails`
- New struct `ApplicationDefinitionArtifact`
- New struct `ApplicationDefinitionPatchable`
- New struct `ApplicationDeploymentPolicy`
- New struct `ApplicationJitAccessPolicy`
- New struct `ApplicationManagementPolicy`
- New struct `ApplicationNotificationEndpoint`
- New struct `ApplicationNotificationPolicy`
- New struct `ApplicationPackageContact`
- New struct `ApplicationPackageLockingPolicyDefinition`
- New struct `ApplicationPackageSupportUrls`
- New struct `ApplicationPolicy`
- New struct `JitApproverDefinition`
- New struct `JitAuthorizationPolicies`
- New struct `JitRequestDefinition`
- New struct `JitRequestDefinitionListResult`
- New struct `JitRequestMetadata`
- New struct `JitRequestPatchable`
- New struct `JitRequestProperties`
- New struct `JitSchedulingPolicy`
- New struct `ListTokenRequest`
- New struct `ManagedIdentityToken`
- New struct `ManagedIdentityTokenResult`
- New struct `RegistryPackage`
- New struct `RegistryPackageLinks`
- New struct `RegistryPackagePlan`
- New struct `SystemData`
- New struct `UpdateAccessDefinition`
- New struct `UserAssignedResourceIdentity`
- New field `SystemData` in struct `Application`
- New field `SystemData` in struct `ApplicationDefinition`
- New field `DeploymentPolicy`, `LockingPolicy`, `ManagementPolicy`, `NotificationPolicy`, `Policies`, `StorageAccountID` in struct `ApplicationDefinitionProperties`
- New field `SystemData` in struct `ApplicationPatchable`
- New field `Artifacts`, `Authorizations`, `BillingDetails`, `CreatedBy`, `CustomerSupport`, `JitAccessPolicy`, `ManagementMode`, `PublisherTenantID`, `SupportUrls`, `UpdatedBy` in struct `ApplicationProperties`
- New field `UserAssignedIdentities` in struct `Identity`
- New field `ActionType`, `IsDataAction`, `Origin` in struct `Operation`
- New field `Description` in struct `OperationDisplay`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 1.1.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-16)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armmanagedapplications` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).