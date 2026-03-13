# Release History

## 2.0.0-beta.1 (2026-03-13)
### Breaking Changes

- Function `NewClient` parameter(s) have been changed from `(credential azcore.TokenCredential, options *arm.ClientOptions)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`
- Function `*Client.AcceptOwnershipStatus` parameter(s) have been changed from `(ctx context.Context, subscriptionID string, options *ClientAcceptOwnershipStatusOptions)` to `(ctx context.Context, options *ClientAcceptOwnershipStatusOptions)`
- Function `*Client.BeginAcceptOwnership` parameter(s) have been changed from `(ctx context.Context, subscriptionID string, body AcceptOwnershipRequest, options *ClientBeginAcceptOwnershipOptions)` to `(ctx context.Context, body AcceptOwnershipRequest, options *ClientBeginAcceptOwnershipOptions)`
- Function `NewClientFactory` parameter(s) have been changed from `(credential azcore.TokenCredential, options *arm.ClientOptions)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`
- Function `NewSubscriptionsClient` parameter(s) have been changed from `(credential azcore.TokenCredential, options *arm.ClientOptions)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`
- Enum `SpendingLimit` has been removed
- Enum `SubscriptionState` has been removed
- Function `*ClientFactory.NewTenantsClient` has been removed
- Function `*SubscriptionsClient.Get` has been removed
- Function `*SubscriptionsClient.NewListLocationsPager` has been removed
- Function `*SubscriptionsClient.NewListPager` has been removed
- Function `NewTenantsClient` has been removed
- Function `*TenantsClient.NewListPager` has been removed
- Operation `*AliasClient.List` has supported pagination, use `*AliasClient.NewListPager` instead.
- Struct `ErrorResponse` has been removed
- Struct `ErrorResponseBody` has been removed
- Struct `ListResult` has been removed
- Struct `Location` has been removed
- Struct `LocationListResult` has been removed
- Struct `Policies` has been removed
- Struct `Subscription` has been removed
- Struct `TenantIDDescription` has been removed
- Struct `TenantListResult` has been removed

### Features Added

- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `ChangeDirectoryOperationStatus` with values `ChangeDirectoryOperationStatusCompleted`, `ChangeDirectoryOperationStatusInProgress`, `ChangeDirectoryOperationStatusInitialized`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `Provisioning` with values `ProvisioningAccepted`, `ProvisioningPending`, `ProvisioningSucceeded`
- New function `*ClientFactory.NewOperationClient() *OperationClient`
- New function `NewOperationClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*OperationClient, error)`
- New function `*OperationClient.Get(ctx context.Context, operationID string, options *OperationClientGetOptions) (OperationClientGetResponse, error)`
- New function `*SubscriptionsClient.AcceptTargetDirectory(ctx context.Context, options *SubscriptionsClientAcceptTargetDirectoryOptions) (SubscriptionsClientAcceptTargetDirectoryResponse, error)`
- New function `*SubscriptionsClient.DeleteTargetDirectory(ctx context.Context, options *SubscriptionsClientDeleteTargetDirectoryOptions) (SubscriptionsClientDeleteTargetDirectoryResponse, error)`
- New function `*SubscriptionsClient.GetTargetDirectory(ctx context.Context, options *SubscriptionsClientGetTargetDirectoryOptions) (SubscriptionsClientGetTargetDirectoryResponse, error)`
- New function `*SubscriptionsClient.NewListTargetDirectoryPager(options *SubscriptionsClientListTargetDirectoryOptions) *runtime.Pager[SubscriptionsClientListTargetDirectoryResponse]`
- New function `*SubscriptionsClient.PutTargetDirectory(ctx context.Context, body TargetDirectoryRequest, options *SubscriptionsClientPutTargetDirectoryOptions) (SubscriptionsClientPutTargetDirectoryResponse, error)`
- New function `*SubscriptionsClient.TargetDirectoryStatus(ctx context.Context, options *SubscriptionsClientTargetDirectoryStatusOptions) (SubscriptionsClientTargetDirectoryStatusResponse, error)`
- New struct `AcceptOwnershipFinalResult`
- New struct `CreationResult`
- New struct `TargetDirectoryListResult`
- New struct `TargetDirectoryRequest`
- New struct `TargetDirectoryRequestProperties`
- New struct `TargetDirectoryResult`
- New struct `TargetDirectoryResultProperties`
- New field `ProvisioningState` in struct `AcceptOwnershipStatusResponse`
- New anonymous field `AcceptOwnershipFinalResult` in struct `ClientAcceptOwnershipResponse`
- New field `ActionType`, `Origin` in struct `Operation`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).