# Release History

## 0.3.0 (2022-01-13)
### Breaking Changes

- Function `*UserAssignedIdentitiesClient.Update` parameter(s) have been changed from `(context.Context, string, string, IdentityUpdate, *UserAssignedIdentitiesUpdateOptions)` to `(context.Context, string, string, IdentityUpdate, *UserAssignedIdentitiesClientUpdateOptions)`
- Function `*UserAssignedIdentitiesClient.Update` return value(s) have been changed from `(UserAssignedIdentitiesUpdateResponse, error)` to `(UserAssignedIdentitiesClientUpdateResponse, error)`
- Function `*UserAssignedIdentitiesClient.Get` parameter(s) have been changed from `(context.Context, string, string, *UserAssignedIdentitiesGetOptions)` to `(context.Context, string, string, *UserAssignedIdentitiesClientGetOptions)`
- Function `*UserAssignedIdentitiesClient.Get` return value(s) have been changed from `(UserAssignedIdentitiesGetResponse, error)` to `(UserAssignedIdentitiesClientGetResponse, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*UserAssignedIdentitiesClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *UserAssignedIdentitiesDeleteOptions)` to `(context.Context, string, string, *UserAssignedIdentitiesClientDeleteOptions)`
- Function `*UserAssignedIdentitiesClient.Delete` return value(s) have been changed from `(UserAssignedIdentitiesDeleteResponse, error)` to `(UserAssignedIdentitiesClientDeleteResponse, error)`
- Function `*UserAssignedIdentitiesClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, Identity, *UserAssignedIdentitiesCreateOrUpdateOptions)` to `(context.Context, string, string, Identity, *UserAssignedIdentitiesClientCreateOrUpdateOptions)`
- Function `*UserAssignedIdentitiesClient.CreateOrUpdate` return value(s) have been changed from `(UserAssignedIdentitiesCreateOrUpdateResponse, error)` to `(UserAssignedIdentitiesClientCreateOrUpdateResponse, error)`
- Function `*SystemAssignedIdentitiesClient.GetByScope` parameter(s) have been changed from `(context.Context, string, *SystemAssignedIdentitiesGetByScopeOptions)` to `(context.Context, string, *SystemAssignedIdentitiesClientGetByScopeOptions)`
- Function `*SystemAssignedIdentitiesClient.GetByScope` return value(s) have been changed from `(SystemAssignedIdentitiesGetByScopeResponse, error)` to `(SystemAssignedIdentitiesClientGetByScopeResponse, error)`
- Function `*UserAssignedIdentitiesClient.ListBySubscription` parameter(s) have been changed from `(*UserAssignedIdentitiesListBySubscriptionOptions)` to `(*UserAssignedIdentitiesClientListBySubscriptionOptions)`
- Function `*UserAssignedIdentitiesClient.ListBySubscription` return value(s) have been changed from `(*UserAssignedIdentitiesListBySubscriptionPager)` to `(*UserAssignedIdentitiesClientListBySubscriptionPager)`
- Function `*UserAssignedIdentitiesClient.ListByResourceGroup` parameter(s) have been changed from `(string, *UserAssignedIdentitiesListByResourceGroupOptions)` to `(string, *UserAssignedIdentitiesClientListByResourceGroupOptions)`
- Function `*UserAssignedIdentitiesClient.ListByResourceGroup` return value(s) have been changed from `(*UserAssignedIdentitiesListByResourceGroupPager)` to `(*UserAssignedIdentitiesClientListByResourceGroupPager)`
- Function `*OperationsListPager.PageResponse` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Function `*UserAssignedIdentitiesListByResourceGroupPager.NextPage` has been removed
- Function `*UserAssignedIdentitiesListByResourceGroupPager.PageResponse` has been removed
- Function `*UserAssignedIdentitiesListBySubscriptionPager.Err` has been removed
- Function `Resource.MarshalJSON` has been removed
- Function `*UserAssignedIdentitiesListBySubscriptionPager.PageResponse` has been removed
- Function `CloudError.Error` has been removed
- Function `*UserAssignedIdentitiesListBySubscriptionPager.NextPage` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `*UserAssignedIdentitiesListByResourceGroupPager.Err` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Struct `SystemAssignedIdentitiesGetByScopeOptions` has been removed
- Struct `SystemAssignedIdentitiesGetByScopeResponse` has been removed
- Struct `SystemAssignedIdentitiesGetByScopeResult` has been removed
- Struct `UserAssignedIdentitiesCreateOrUpdateOptions` has been removed
- Struct `UserAssignedIdentitiesCreateOrUpdateResponse` has been removed
- Struct `UserAssignedIdentitiesCreateOrUpdateResult` has been removed
- Struct `UserAssignedIdentitiesDeleteOptions` has been removed
- Struct `UserAssignedIdentitiesDeleteResponse` has been removed
- Struct `UserAssignedIdentitiesGetOptions` has been removed
- Struct `UserAssignedIdentitiesGetResponse` has been removed
- Struct `UserAssignedIdentitiesGetResult` has been removed
- Struct `UserAssignedIdentitiesListByResourceGroupOptions` has been removed
- Struct `UserAssignedIdentitiesListByResourceGroupPager` has been removed
- Struct `UserAssignedIdentitiesListByResourceGroupResponse` has been removed
- Struct `UserAssignedIdentitiesListByResourceGroupResult` has been removed
- Struct `UserAssignedIdentitiesListBySubscriptionOptions` has been removed
- Struct `UserAssignedIdentitiesListBySubscriptionPager` has been removed
- Struct `UserAssignedIdentitiesListBySubscriptionResponse` has been removed
- Struct `UserAssignedIdentitiesListBySubscriptionResult` has been removed
- Struct `UserAssignedIdentitiesUpdateOptions` has been removed
- Struct `UserAssignedIdentitiesUpdateResponse` has been removed
- Struct `UserAssignedIdentitiesUpdateResult` has been removed
- Field `ProxyResource` of struct `SystemAssignedIdentity` has been removed
- Field `Resource` of struct `ProxyResource` has been removed
- Field `Resource` of struct `IdentityUpdate` has been removed
- Field `TrackedResource` of struct `Identity` has been removed
- Field `InnerError` of struct `CloudError` has been removed
- Field `Resource` of struct `TrackedResource` has been removed

### Features Added

- New function `*OperationsClientListPager.Err() error`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*UserAssignedIdentitiesClientListBySubscriptionPager.Err() error`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*UserAssignedIdentitiesClientListBySubscriptionPager.NextPage(context.Context) bool`
- New function `*UserAssignedIdentitiesClientListByResourceGroupPager.PageResponse() UserAssignedIdentitiesClientListByResourceGroupResponse`
- New function `*UserAssignedIdentitiesClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*UserAssignedIdentitiesClientListBySubscriptionPager.PageResponse() UserAssignedIdentitiesClientListBySubscriptionResponse`
- New function `*UserAssignedIdentitiesClientListByResourceGroupPager.Err() error`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New struct `SystemAssignedIdentitiesClientGetByScopeOptions`
- New struct `SystemAssignedIdentitiesClientGetByScopeResponse`
- New struct `SystemAssignedIdentitiesClientGetByScopeResult`
- New struct `UserAssignedIdentitiesClientCreateOrUpdateOptions`
- New struct `UserAssignedIdentitiesClientCreateOrUpdateResponse`
- New struct `UserAssignedIdentitiesClientCreateOrUpdateResult`
- New struct `UserAssignedIdentitiesClientDeleteOptions`
- New struct `UserAssignedIdentitiesClientDeleteResponse`
- New struct `UserAssignedIdentitiesClientGetOptions`
- New struct `UserAssignedIdentitiesClientGetResponse`
- New struct `UserAssignedIdentitiesClientGetResult`
- New struct `UserAssignedIdentitiesClientListByResourceGroupOptions`
- New struct `UserAssignedIdentitiesClientListByResourceGroupPager`
- New struct `UserAssignedIdentitiesClientListByResourceGroupResponse`
- New struct `UserAssignedIdentitiesClientListByResourceGroupResult`
- New struct `UserAssignedIdentitiesClientListBySubscriptionOptions`
- New struct `UserAssignedIdentitiesClientListBySubscriptionPager`
- New struct `UserAssignedIdentitiesClientListBySubscriptionResponse`
- New struct `UserAssignedIdentitiesClientListBySubscriptionResult`
- New struct `UserAssignedIdentitiesClientUpdateOptions`
- New struct `UserAssignedIdentitiesClientUpdateResponse`
- New struct `UserAssignedIdentitiesClientUpdateResult`
- New field `ID` in struct `TrackedResource`
- New field `Name` in struct `TrackedResource`
- New field `Type` in struct `TrackedResource`
- New field `Error` in struct `CloudError`
- New field `Location` in struct `Identity`
- New field `Tags` in struct `Identity`
- New field `ID` in struct `Identity`
- New field `Name` in struct `Identity`
- New field `Type` in struct `Identity`
- New field `ID` in struct `ProxyResource`
- New field `Name` in struct `ProxyResource`
- New field `Type` in struct `ProxyResource`
- New field `ID` in struct `SystemAssignedIdentity`
- New field `Name` in struct `SystemAssignedIdentity`
- New field `Type` in struct `SystemAssignedIdentity`
- New field `ID` in struct `IdentityUpdate`
- New field `Name` in struct `IdentityUpdate`
- New field `Type` in struct `IdentityUpdate`


## 0.2.0 (2021-10-29)

### Breaking Changes

- `arm.Connection` has been removed in `github.com/Azure/azure-sdk-for-go/sdk/azcore/v0.20.0`
- The parameters of `NewXXXClient` has been changed from `(con *arm.Connection, subscriptionID string)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`

## 0.1.0 (2021-10-15)

- Initial preview release.
