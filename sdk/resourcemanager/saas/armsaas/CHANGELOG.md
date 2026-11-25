# Release History

## 1.0.0 (2026-11-25)
### Breaking Changes

- Function `*SubscriptionLevelClient.BeginCreateOrUpdate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, resourceName string, parameters ResourceCreation, options *SubscriptionLevelClientBeginCreateOrUpdateOptions)` to `(ctx context.Context, resourceGroupName string, resourceName string, resource Resource, options *SubscriptionLevelClientBeginCreateOrUpdateOptions)`
- Function `*SubscriptionLevelClient.BeginUpdate` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, resourceName string, parameters ResourceCreation, options *SubscriptionLevelClientBeginUpdateOptions)` to `(ctx context.Context, resourceGroupName string, resourceName string, properties ResourceUpdate, options *SubscriptionLevelClientBeginUpdateOptions)`
- Type of `ResourceProperties.Quantity` has been changed from `*float32` to `*int64`
- Enum `PaymentChannelType` has been removed
- Enum `SaasAppStatus` has been removed
- Function `NewApplicationsClient` has been removed
- Function `*ApplicationsClient.NewListPager` has been removed
- Function `NewClient` has been removed
- Function `*Client.BeginCreateResource` has been removed
- Function `*Client.BeginDelete` has been removed
- Function `*Client.GetResource` has been removed
- Function `*Client.BeginUpdateResource` has been removed
- Function `*ClientFactory.NewApplicationsClient` has been removed
- Function `*ClientFactory.NewClient` has been removed
- Function `*ClientFactory.NewOperationClient` has been removed
- Function `*ClientFactory.NewResourcesClient` has been removed
- Function `NewOperationClient` has been removed
- Function `*OperationClient.BeginGet` has been removed
- Function `NewResourcesClient` has been removed
- Function `*ResourcesClient.ListAccessToken` has been removed
- Function `*ResourcesClient.NewListPager` has been removed
- Function `*SubscriptionLevelClient.NewListByAzureSubscriptionPager` has been removed
- Function `*SubscriptionLevelClient.BeginMoveResources` has been removed
- Function `*SubscriptionLevelClient.BeginUpdateToUnsubscribed` has been removed
- Function `*SubscriptionLevelClient.ValidateMoveResources` has been removed
- Struct `App` has been removed
- Struct `AppOperation` has been removed
- Struct `AppOperationDisplay` has been removed
- Struct `AppOperationsResponseWithContinuation` has been removed
- Struct `AppPlan` has been removed
- Struct `AppProperties` has been removed
- Struct `AppResponseWithContinuation` has been removed
- Struct `CreationProperties` has been removed
- Struct `DeleteOptions` has been removed
- Struct `ErrorResponse` has been removed
- Struct `MoveResource` has been removed
- Struct `Properties` has been removed
- Struct `PropertiesTerm` has been removed
- Struct `ResourceCreation` has been removed
- Struct `ResourceResponseWithContinuation` has been removed
- Struct `Result` has been removed
- Field `AppOperationsResponseWithContinuation` of struct `OperationsClientListResponse` has been removed
- Field `Created`, `LastModified`, `OfferID`, `PaymentChannelMetadata`, `PaymentChannelType`, `PublisherID`, `PublisherTestEnvironment`, `SKUID`, `SaasResourceName`, `SaasSessionID`, `SaasSubscriptionID`, `Term`, `TermID` of struct `ResourceProperties` has been removed
- Field `ResourceResponseWithContinuation` of struct `SubscriptionLevelClientListByResourceGroupResponse` has been removed

### Features Added

- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `CheckNameAvailabilityReason` with values `CheckNameAvailabilityReasonAlreadyExists`, `CheckNameAvailabilityReasonInvalid`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `DeletionSource` with values `DeletionSourceAzureSubscriptionDeletion`, `DeletionSourceCustomerCancellation`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `PatchIntent` with values `PatchIntentReset`, `PatchIntentRetain`, `PatchIntentScheduleChange`
- New enum type `SaasProvisioningState` with values `SaasProvisioningStateAccepted`, `SaasProvisioningStateCanceled`, `SaasProvisioningStateFailed`, `SaasProvisioningStateInProgress`, `SaasProvisioningStateNotStarted`, `SaasProvisioningStateSucceeded`
- New enum type `SubscriptionStatusReason` with values `SubscriptionStatusReasonAzureSubscriptionDeletion`, `SubscriptionStatusReasonNone`, `SubscriptionStatusReasonSuspendedInternal`
- New function `*ClientFactory.NewOperationStatusesClient() *OperationStatusesClient`
- New function `*ClientFactory.NewSubscriptionLevelGlobalClient() *SubscriptionLevelGlobalClient`
- New function `*ClientFactory.NewTenantLevelClient() *TenantLevelClient`
- New function `NewOperationStatusesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*OperationStatusesClient, error)`
- New function `*OperationStatusesClient.GetSubscriptionOperationStatus(ctx context.Context, operationID string, options *OperationStatusesClientGetSubscriptionOperationStatusOptions) (OperationStatusesClientGetSubscriptionOperationStatusResponse, error)`
- New function `*OperationStatusesClient.GetTenantOperationStatus(ctx context.Context, operationID string, options *OperationStatusesClientGetTenantOperationStatusOptions) (OperationStatusesClientGetTenantOperationStatusResponse, error)`
- New function `*SubscriptionLevelClient.BeginUnsubscribe(ctx context.Context, resourceGroupName string, resourceName string, body UnsubscribeOptions, options *SubscriptionLevelClientBeginUnsubscribeOptions) (*runtime.Poller[SubscriptionLevelClientUnsubscribeResponse], error)`
- New function `*SubscriptionLevelClient.ValidateUpdate(ctx context.Context, resourceGroupName string, resourceName string, body ResourceUpdate, options *SubscriptionLevelClientValidateUpdateOptions) (SubscriptionLevelClientValidateUpdateResponse, error)`
- New function `NewSubscriptionLevelGlobalClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SubscriptionLevelGlobalClient, error)`
- New function `*SubscriptionLevelGlobalClient.CheckNameAvailability(ctx context.Context, body CheckNameAvailabilityRequest, options *SubscriptionLevelGlobalClientCheckNameAvailabilityOptions) (SubscriptionLevelGlobalClientCheckNameAvailabilityResponse, error)`
- New function `*SubscriptionLevelGlobalClient.NewListByAzureSubscriptionPager(options *SubscriptionLevelGlobalClientListByAzureSubscriptionOptions) *runtime.Pager[SubscriptionLevelGlobalClientListByAzureSubscriptionResponse]`
- New function `*SubscriptionLevelGlobalClient.ValidatePurchase(ctx context.Context, resourceGroupName string, body PurchaseValidationRequest, options *SubscriptionLevelGlobalClientValidatePurchaseOptions) (SubscriptionLevelGlobalClientValidatePurchaseResponse, error)`
- New function `*SubscriptionLevelGlobalClient.ValidatePurchaseBySubscription(ctx context.Context, body PurchaseValidationRequest, options *SubscriptionLevelGlobalClientValidatePurchaseBySubscriptionOptions) (SubscriptionLevelGlobalClientValidatePurchaseBySubscriptionResponse, error)`
- New function `NewTenantLevelClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*TenantLevelClient, error)`
- New function `*TenantLevelClient.BeginDeleteResource(ctx context.Context, resourceID string, options *TenantLevelClientBeginDeleteResourceOptions) (*runtime.Poller[TenantLevelClientDeleteResourceResponse], error)`
- New function `*TenantLevelClient.GetResource(ctx context.Context, resourceID string, options *TenantLevelClientGetResourceOptions) (TenantLevelClientGetResourceResponse, error)`
- New function `*TenantLevelClient.ListAccessToken(ctx context.Context, resourceID string, options *TenantLevelClientListAccessTokenOptions) (TenantLevelClientListAccessTokenResponse, error)`
- New function `*TenantLevelClient.NewListPager(options *TenantLevelClientListOptions) *runtime.Pager[TenantLevelClientListResponse]`
- New function `*TenantLevelClient.BeginUpdateResource(ctx context.Context, resourceID string, properties TenantResourceUpdate, options *TenantLevelClientBeginUpdateResourceOptions) (*runtime.Poller[TenantLevelClientUpdateResourceResponse], error)`
- New struct `CheckNameAvailabilityRequest`
- New struct `CheckNameAvailabilityResponse`
- New struct `Operation`
- New struct `OperationDisplay`
- New struct `OperationListResult`
- New struct `OperationStatusResult`
- New struct `PurchaseValidationProperties`
- New struct `PurchaseValidationRequest`
- New struct `ResourceListResult`
- New struct `ResourceUpdate`
- New struct `ResourceUpdateProperties`
- New struct `ResourceValidationError`
- New struct `ResourceValidationResult`
- New struct `SystemData`
- New struct `TenantResource`
- New struct `TenantResourceListResult`
- New struct `TenantResourceUpdate`
- New struct `TenantResourceUpdateProperties`
- New struct `UnsubscribeOptions`
- New field `ExpiryTime`, `OfferURIWithToken` in struct `AccessTokenResult`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New field `Location`, `SystemData` in struct `Resource`
- New field `AssetID`, `BillingPeriod`, `Description`, `FutureStartDate`, `Intent`, `OfferName`, `PlanName`, `PrivateOfferID`, `ProductCode`, `ProvisioningState`, `PublisherName`, `PurchaserEmail`, `QuoteID`, `RiskPropertyBagHeader`, `StatusReason`, `Storefront`, `TermEndDate`, `TermStartDate`, `TermUnit` in struct `ResourceProperties`
- New field `ExpandCustomerData` in struct `SubscriptionLevelClientGetOptions`
- New anonymous field `ResourceListResult` in struct `SubscriptionLevelClientListByResourceGroupResponse`


## 0.7.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.6.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 0.6.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.5.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/saas/armsaas` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).