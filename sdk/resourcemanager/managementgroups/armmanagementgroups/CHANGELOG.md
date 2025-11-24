# Release History

## 2.0.0 (2025-11-24)
### Breaking Changes

- Function `NewClientFactory` parameter(s) have been changed from `(azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewManagementGroupSubscriptionsClient` parameter(s) have been changed from `(azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `*ManagementGroupSubscriptionsClient.Create` parameter(s) have been changed from `(context.Context, string, string, *ManagementGroupSubscriptionsClientCreateOptions)` to `(context.Context, string, *ManagementGroupSubscriptionsClientCreateOptions)`
- Function `*ManagementGroupSubscriptionsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *ManagementGroupSubscriptionsClientDeleteOptions)` to `(context.Context, string, *ManagementGroupSubscriptionsClientDeleteOptions)`
- Function `*ManagementGroupSubscriptionsClient.GetSubscription` parameter(s) have been changed from `(context.Context, string, string, *ManagementGroupSubscriptionsClientGetSubscriptionOptions)` to `(context.Context, string, *ManagementGroupSubscriptionsClientGetSubscriptionOptions)`
- Type of `Operation.Display` has been changed from `*OperationDisplayProperties` to `*OperationDisplay`
- Function `NewAPIClient` has been removed
- Function `*APIClient.CheckNameAvailability` has been removed
- Function `*APIClient.StartTenantBackfill` has been removed
- Function `*APIClient.TenantBackfillStatus` has been removed
- Function `*ClientFactory.NewAPIClient` has been removed
- Operation `*HierarchySettingsClient.List` has supported pagination, use `*HierarchySettingsClient.NewListPager` instead.
- Struct `CreateManagementGroupChildInfo` has been removed
- Struct `EntityHierarchyItem` has been removed
- Struct `EntityHierarchyItemProperties` has been removed
- Struct `ErrorDetails` has been removed
- Struct `ErrorResponse` has been removed
- Struct `OperationDisplayProperties` has been removed
- Struct `OperationResults` has been removed
- Field `Count` of struct `EntityListResult` has been removed

### Features Added

- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New function `*ClientFactory.NewManagementClient() *ManagementClient`
- New function `NewManagementClient(azcore.TokenCredential, *arm.ClientOptions) (*ManagementClient, error)`
- New function `*ManagementClient.CheckNameAvailability(context.Context, CheckNameAvailabilityRequest, *ManagementClientCheckNameAvailabilityOptions) (ManagementClientCheckNameAvailabilityResponse, error)`
- New function `*ManagementClient.NewClient() *Client`
- New function `*ManagementClient.NewEntitiesClient() *EntitiesClient`
- New function `*ManagementClient.NewHierarchySettingsClient() *HierarchySettingsClient`
- New function `*ManagementClient.NewManagementGroupSubscriptionsClient() *ManagementGroupSubscriptionsClient`
- New function `*ManagementClient.NewOperationsClient() *OperationsClient`
- New function `*ManagementClient.StartTenantBackfill(context.Context, *ManagementClientStartTenantBackfillOptions) (ManagementClientStartTenantBackfillResponse, error)`
- New function `*ManagementClient.TenantBackfillStatus(context.Context, *ManagementClientTenantBackfillStatusOptions) (ManagementClientTenantBackfillStatusResponse, error)`
- New struct `OperationDisplay`
- New struct `SystemData`
- New field `SystemData` in struct `HierarchySettings`
- New field `SystemData` in struct `ManagementGroup`
- New field `ActionType`, `IsDataAction`, `Origin` in struct `Operation`
- New field `SystemData` in struct `SubscriptionUnderManagementGroup`


## 1.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 1.1.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 1.0.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/managementgroups/armmanagementgroups` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).