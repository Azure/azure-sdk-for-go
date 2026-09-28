# Release History

## 2.0.0-beta.1 (2026-09-28)
### Breaking Changes

- Type of `ApplicationArtifact.Name` has been changed from `*string` to `*ApplicationArtifactName`
- Type of `ApplicationDefinitionProperties.Artifacts` has been changed from `[]*ApplicationArtifact` to `[]*ApplicationDefinitionArtifact`
- Type of `ApplicationDefinitionProperties.Authorizations` has been changed from `[]*ApplicationProviderAuthorization` to `[]*ApplicationAuthorization`
- Type of `ApplicationDefinitionProperties.IsEnabled` has been changed from `*string` to `*bool`
- Type of `ApplicationPatchable.Properties` has been changed from `*ApplicationPropertiesPatchable` to `*ApplicationProperties`
- Type of `Identity.Type` has been changed from `*string` to `*ResourceIdentityType`
- Struct `ApplicationProviderAuthorization` has been removed
- Field `Identity` of struct `ApplicationDefinition` has been removed
- Field `ErrorCode`, `ErrorMessage`, `HTTPStatus` of struct `ErrorResponse` has been removed
- Field `Identity` of struct `GenericResource` has been removed

### Features Added

- New value `ApplicationArtifactTypeNotSpecified` added to enum type `ApplicationArtifactType`
- New value `ProvisioningStateNotSpecified` added to enum type `ProvisioningState`
- New enum type `ApplicationArtifactName` with values `ApplicationArtifactNameAuthorizations`, `ApplicationArtifactNameCustomRoleDefinition`, `ApplicationArtifactNameNotSpecified`, `ApplicationArtifactNameViewDefinition`
- New enum type `ApplicationDefinitionArtifactName` with values `ApplicationDefinitionArtifactNameApplicationResourceTemplate`, `ApplicationDefinitionArtifactNameCreateUIDefinition`, `ApplicationDefinitionArtifactNameMainTemplateParameters`, `ApplicationDefinitionArtifactNameNotSpecified`
- New enum type `ApplicationManagementMode` with values `ApplicationManagementModeManaged`, `ApplicationManagementModeNotSpecified`, `ApplicationManagementModeUnmanaged`
- New enum type `DeploymentMode` with values `DeploymentModeComplete`, `DeploymentModeIncremental`, `DeploymentModeNotSpecified`
- New enum type `JitApprovalMode` with values `JitApprovalModeAutoApprove`, `JitApprovalModeManualApprove`, `JitApprovalModeNotSpecified`
- New enum type `JitApproverType` with values `JitApproverTypeGroup`, `JitApproverTypeUser`
- New enum type `JitRequestState` with values `JitRequestStateApproved`, `JitRequestStateCanceled`, `JitRequestStateDenied`, `JitRequestStateExpired`, `JitRequestStateFailed`, `JitRequestStateNotSpecified`, `JitRequestStatePending`, `JitRequestStateTimeout`
- New enum type `JitSchedulingType` with values `JitSchedulingTypeNotSpecified`, `JitSchedulingTypeOnce`, `JitSchedulingTypeRecurring`
- New enum type `ResourceIdentityType` with values `ResourceIdentityTypeNone`, `ResourceIdentityTypeSystemAssigned`, `ResourceIdentityTypeSystemAssignedUserAssigned`, `ResourceIdentityTypeUserAssigned`
- New function `*ApplicationsClient.BeginRefreshPermissions(ctx context.Context, resourceGroupName string, applicationName string, options *ApplicationsClientBeginRefreshPermissionsOptions) (*runtime.Poller[ApplicationsClientRefreshPermissionsResponse], error)`
- New function `*ClientFactory.NewJitRequestsClient() *JitRequestsClient`
- New function `NewJitRequestsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*JitRequestsClient, error)`
- New function `*JitRequestsClient.BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, jitRequestName string, parameters JitRequestDefinition, options *JitRequestsClientBeginCreateOrUpdateOptions) (*runtime.Poller[JitRequestsClientCreateOrUpdateResponse], error)`
- New function `*JitRequestsClient.Delete(ctx context.Context, resourceGroupName string, jitRequestName string, options *JitRequestsClientDeleteOptions) (JitRequestsClientDeleteResponse, error)`
- New function `*JitRequestsClient.Get(ctx context.Context, resourceGroupName string, jitRequestName string, options *JitRequestsClientGetOptions) (JitRequestsClientGetResponse, error)`
- New function `*JitRequestsClient.ListByResourceGroup(ctx context.Context, resourceGroupName string, options *JitRequestsClientListByResourceGroupOptions) (JitRequestsClientListByResourceGroupResponse, error)`
- New function `*JitRequestsClient.ListBySubscription(ctx context.Context, options *JitRequestsClientListBySubscriptionOptions) (JitRequestsClientListBySubscriptionResponse, error)`
- New function `*JitRequestsClient.Update(ctx context.Context, resourceGroupName string, jitRequestName string, parameters JitRequestPatchable, options *JitRequestsClientUpdateOptions) (JitRequestsClientUpdateResponse, error)`
- New struct `ApplicationAuthorization`
- New struct `ApplicationBillingDetailsDefinition`
- New struct `ApplicationClientDetails`
- New struct `ApplicationDefinitionArtifact`
- New struct `ApplicationDeploymentPolicy`
- New struct `ApplicationJitAccessPolicy`
- New struct `ApplicationManagementPolicy`
- New struct `ApplicationNotificationEndpoint`
- New struct `ApplicationNotificationPolicy`
- New struct `ApplicationPackageContact`
- New struct `ApplicationPackageLockingPolicyDefinition`
- New struct `ApplicationPackageSupportUrls`
- New struct `ApplicationPolicy`
- New struct `ErrorAdditionalInfo`
- New struct `JitApproverDefinition`
- New struct `JitAuthorizationPolicies`
- New struct `JitRequestDefinition`
- New struct `JitRequestDefinitionListResult`
- New struct `JitRequestPatchable`
- New struct `JitRequestProperties`
- New struct `JitSchedulingPolicy`
- New struct `UserAssignedResourceIdentity`
- New field `DeploymentPolicy`, `LockingPolicy`, `ManagementPolicy`, `NotificationPolicy`, `Policies`, `StorageAccountID` in struct `ApplicationDefinitionProperties`
- New field `Artifacts`, `Authorizations`, `BillingDetails`, `CreatedBy`, `CustomerSupport`, `JitAccessPolicy`, `ManagementMode`, `PublisherTenantID`, `SupportUrls`, `UpdatedBy` in struct `ApplicationProperties`
- New field `AdditionalInfo`, `Code`, `Details`, `Message`, `Target` in struct `ErrorResponse`
- New field `UserAssignedIdentities` in struct `Identity`
- New field `IsDataAction` in struct `Operation`
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