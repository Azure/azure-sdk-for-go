# Release History

## 0.4.1 (2022-04-18)
### Other Changes


## 0.4.0 (2022-04-18)
### Breaking Changes

- Function `*ManagementLocksClient.ListAtResourceLevel` has been removed
- Function `*ManagementLocksClient.ListAtSubscriptionLevel` has been removed
- Function `*ManagementLocksClient.ListByScope` has been removed
- Function `*ManagementLocksClient.ListAtResourceGroupLevel` has been removed
- Function `*AuthorizationOperationsClient.List` has been removed

### Features Added

- New function `*ManagementLocksClient.NewListAtResourceLevelPager(string, string, string, string, string, *ManagementLocksClientListAtResourceLevelOptions) *runtime.Pager[ManagementLocksClientListAtResourceLevelResponse]`
- New function `*AuthorizationOperationsClient.NewListPager(*AuthorizationOperationsClientListOptions) *runtime.Pager[AuthorizationOperationsClientListResponse]`
- New function `*ManagementLocksClient.NewListAtSubscriptionLevelPager(*ManagementLocksClientListAtSubscriptionLevelOptions) *runtime.Pager[ManagementLocksClientListAtSubscriptionLevelResponse]`
- New function `*ManagementLocksClient.NewListAtResourceGroupLevelPager(string, *ManagementLocksClientListAtResourceGroupLevelOptions) *runtime.Pager[ManagementLocksClientListAtResourceGroupLevelResponse]`
- New function `*ManagementLocksClient.NewListByScopePager(string, *ManagementLocksClientListByScopeOptions) *runtime.Pager[ManagementLocksClientListByScopeResponse]`


## 0.3.0 (2022-04-14)
### Breaking Changes

- Function `*ManagementLocksClient.ListAtResourceGroupLevel` return value(s) have been changed from `(*ManagementLocksClientListAtResourceGroupLevelPager)` to `(*runtime.Pager[ManagementLocksClientListAtResourceGroupLevelResponse])`
- Function `NewAuthorizationOperationsClient` return value(s) have been changed from `(*AuthorizationOperationsClient)` to `(*AuthorizationOperationsClient, error)`
- Function `*ManagementLocksClient.ListAtResourceLevel` return value(s) have been changed from `(*ManagementLocksClientListAtResourceLevelPager)` to `(*runtime.Pager[ManagementLocksClientListAtResourceLevelResponse])`
- Function `*ManagementLocksClient.ListAtSubscriptionLevel` return value(s) have been changed from `(*ManagementLocksClientListAtSubscriptionLevelPager)` to `(*runtime.Pager[ManagementLocksClientListAtSubscriptionLevelResponse])`
- Function `*AuthorizationOperationsClient.List` return value(s) have been changed from `(*AuthorizationOperationsClientListPager)` to `(*runtime.Pager[AuthorizationOperationsClientListResponse])`
- Function `NewManagementLocksClient` return value(s) have been changed from `(*ManagementLocksClient)` to `(*ManagementLocksClient, error)`
- Function `*ManagementLocksClient.ListByScope` return value(s) have been changed from `(*ManagementLocksClientListByScopePager)` to `(*runtime.Pager[ManagementLocksClientListByScopeResponse])`
- Type of `ErrorAdditionalInfo.Info` has been changed from `map[string]interface{}` to `interface{}`
- Function `*ManagementLocksClientListAtSubscriptionLevelPager.NextPage` has been removed
- Function `*ManagementLocksClientListByScopePager.Err` has been removed
- Function `*AuthorizationOperationsClientListPager.Err` has been removed
- Function `*AuthorizationOperationsClientListPager.PageResponse` has been removed
- Function `CreatedByType.ToPtr` has been removed
- Function `*ManagementLocksClientListAtResourceLevelPager.NextPage` has been removed
- Function `*ManagementLocksClientListAtResourceGroupLevelPager.Err` has been removed
- Function `*ManagementLocksClientListAtSubscriptionLevelPager.PageResponse` has been removed
- Function `*ManagementLocksClientListAtResourceLevelPager.PageResponse` has been removed
- Function `*AuthorizationOperationsClientListPager.NextPage` has been removed
- Function `*ManagementLocksClientListAtSubscriptionLevelPager.Err` has been removed
- Function `*ManagementLocksClientListAtResourceGroupLevelPager.NextPage` has been removed
- Function `*ManagementLocksClientListByScopePager.PageResponse` has been removed
- Function `LockLevel.ToPtr` has been removed
- Function `*ManagementLocksClientListByScopePager.NextPage` has been removed
- Function `*ManagementLocksClientListAtResourceLevelPager.Err` has been removed
- Function `*ManagementLocksClientListAtResourceGroupLevelPager.PageResponse` has been removed
- Struct `AuthorizationOperationsClientListPager` has been removed
- Struct `AuthorizationOperationsClientListResult` has been removed
- Struct `ManagementLocksClientCreateOrUpdateAtResourceGroupLevelResult` has been removed
- Struct `ManagementLocksClientCreateOrUpdateAtResourceLevelResult` has been removed
- Struct `ManagementLocksClientCreateOrUpdateAtSubscriptionLevelResult` has been removed
- Struct `ManagementLocksClientCreateOrUpdateByScopeResult` has been removed
- Struct `ManagementLocksClientGetAtResourceGroupLevelResult` has been removed
- Struct `ManagementLocksClientGetAtResourceLevelResult` has been removed
- Struct `ManagementLocksClientGetAtSubscriptionLevelResult` has been removed
- Struct `ManagementLocksClientGetByScopeResult` has been removed
- Struct `ManagementLocksClientListAtResourceGroupLevelPager` has been removed
- Struct `ManagementLocksClientListAtResourceGroupLevelResult` has been removed
- Struct `ManagementLocksClientListAtResourceLevelPager` has been removed
- Struct `ManagementLocksClientListAtResourceLevelResult` has been removed
- Struct `ManagementLocksClientListAtSubscriptionLevelPager` has been removed
- Struct `ManagementLocksClientListAtSubscriptionLevelResult` has been removed
- Struct `ManagementLocksClientListByScopePager` has been removed
- Struct `ManagementLocksClientListByScopeResult` has been removed
- Field `ManagementLocksClientCreateOrUpdateAtResourceLevelResult` of struct `ManagementLocksClientCreateOrUpdateAtResourceLevelResponse` has been removed
- Field `RawResponse` of struct `ManagementLocksClientCreateOrUpdateAtResourceLevelResponse` has been removed
- Field `RawResponse` of struct `ManagementLocksClientDeleteAtResourceGroupLevelResponse` has been removed
- Field `ManagementLocksClientCreateOrUpdateByScopeResult` of struct `ManagementLocksClientCreateOrUpdateByScopeResponse` has been removed
- Field `RawResponse` of struct `ManagementLocksClientCreateOrUpdateByScopeResponse` has been removed
- Field `ManagementLocksClientListAtResourceGroupLevelResult` of struct `ManagementLocksClientListAtResourceGroupLevelResponse` has been removed
- Field `RawResponse` of struct `ManagementLocksClientListAtResourceGroupLevelResponse` has been removed
- Field `ManagementLocksClientCreateOrUpdateAtSubscriptionLevelResult` of struct `ManagementLocksClientCreateOrUpdateAtSubscriptionLevelResponse` has been removed
- Field `RawResponse` of struct `ManagementLocksClientCreateOrUpdateAtSubscriptionLevelResponse` has been removed
- Field `ManagementLocksClientGetAtResourceLevelResult` of struct `ManagementLocksClientGetAtResourceLevelResponse` has been removed
- Field `RawResponse` of struct `ManagementLocksClientGetAtResourceLevelResponse` has been removed
- Field `ManagementLocksClientGetAtSubscriptionLevelResult` of struct `ManagementLocksClientGetAtSubscriptionLevelResponse` has been removed
- Field `RawResponse` of struct `ManagementLocksClientGetAtSubscriptionLevelResponse` has been removed
- Field `ManagementLocksClientGetByScopeResult` of struct `ManagementLocksClientGetByScopeResponse` has been removed
- Field `RawResponse` of struct `ManagementLocksClientGetByScopeResponse` has been removed
- Field `ManagementLocksClientGetAtResourceGroupLevelResult` of struct `ManagementLocksClientGetAtResourceGroupLevelResponse` has been removed
- Field `RawResponse` of struct `ManagementLocksClientGetAtResourceGroupLevelResponse` has been removed
- Field `ManagementLocksClientListByScopeResult` of struct `ManagementLocksClientListByScopeResponse` has been removed
- Field `RawResponse` of struct `ManagementLocksClientListByScopeResponse` has been removed
- Field `RawResponse` of struct `ManagementLocksClientDeleteAtResourceLevelResponse` has been removed
- Field `AuthorizationOperationsClientListResult` of struct `AuthorizationOperationsClientListResponse` has been removed
- Field `RawResponse` of struct `AuthorizationOperationsClientListResponse` has been removed
- Field `ManagementLocksClientListAtSubscriptionLevelResult` of struct `ManagementLocksClientListAtSubscriptionLevelResponse` has been removed
- Field `RawResponse` of struct `ManagementLocksClientListAtSubscriptionLevelResponse` has been removed
- Field `ManagementLocksClientListAtResourceLevelResult` of struct `ManagementLocksClientListAtResourceLevelResponse` has been removed
- Field `RawResponse` of struct `ManagementLocksClientListAtResourceLevelResponse` has been removed
- Field `ManagementLocksClientCreateOrUpdateAtResourceGroupLevelResult` of struct `ManagementLocksClientCreateOrUpdateAtResourceGroupLevelResponse` has been removed
- Field `RawResponse` of struct `ManagementLocksClientCreateOrUpdateAtResourceGroupLevelResponse` has been removed
- Field `RawResponse` of struct `ManagementLocksClientDeleteAtSubscriptionLevelResponse` has been removed
- Field `RawResponse` of struct `ManagementLocksClientDeleteByScopeResponse` has been removed

### Features Added

- New anonymous field `ManagementLockObject` in struct `ManagementLocksClientCreateOrUpdateAtResourceLevelResponse`
- New anonymous field `ManagementLockObject` in struct `ManagementLocksClientGetAtResourceLevelResponse`
- New anonymous field `ManagementLockListResult` in struct `ManagementLocksClientListAtResourceGroupLevelResponse`
- New anonymous field `ManagementLockObject` in struct `ManagementLocksClientCreateOrUpdateAtResourceGroupLevelResponse`
- New anonymous field `ManagementLockObject` in struct `ManagementLocksClientGetAtSubscriptionLevelResponse`
- New anonymous field `ManagementLockObject` in struct `ManagementLocksClientCreateOrUpdateAtSubscriptionLevelResponse`
- New anonymous field `ManagementLockListResult` in struct `ManagementLocksClientListByScopeResponse`
- New anonymous field `ManagementLockListResult` in struct `ManagementLocksClientListAtResourceLevelResponse`
- New anonymous field `ManagementLockObject` in struct `ManagementLocksClientGetAtResourceGroupLevelResponse`
- New anonymous field `ManagementLockListResult` in struct `ManagementLocksClientListAtSubscriptionLevelResponse`
- New anonymous field `ManagementLockObject` in struct `ManagementLocksClientCreateOrUpdateByScopeResponse`
- New anonymous field `ManagementLockObject` in struct `ManagementLocksClientGetByScopeResponse`
- New anonymous field `OperationListResult` in struct `AuthorizationOperationsClientListResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*ManagementLocksClient.ListAtSubscriptionLevel` parameter(s) have been changed from `(*ManagementLocksListAtSubscriptionLevelOptions)` to `(*ManagementLocksClientListAtSubscriptionLevelOptions)`
- Function `*ManagementLocksClient.ListAtSubscriptionLevel` return value(s) have been changed from `(*ManagementLocksListAtSubscriptionLevelPager)` to `(*ManagementLocksClientListAtSubscriptionLevelPager)`
- Function `*ManagementLocksClient.ListAtResourceGroupLevel` parameter(s) have been changed from `(string, *ManagementLocksListAtResourceGroupLevelOptions)` to `(string, *ManagementLocksClientListAtResourceGroupLevelOptions)`
- Function `*ManagementLocksClient.ListAtResourceGroupLevel` return value(s) have been changed from `(*ManagementLocksListAtResourceGroupLevelPager)` to `(*ManagementLocksClientListAtResourceGroupLevelPager)`
- Function `*ManagementLocksClient.DeleteByScope` parameter(s) have been changed from `(context.Context, string, string, *ManagementLocksDeleteByScopeOptions)` to `(context.Context, string, string, *ManagementLocksClientDeleteByScopeOptions)`
- Function `*ManagementLocksClient.DeleteByScope` return value(s) have been changed from `(ManagementLocksDeleteByScopeResponse, error)` to `(ManagementLocksClientDeleteByScopeResponse, error)`
- Function `*ManagementLocksClient.GetAtResourceGroupLevel` parameter(s) have been changed from `(context.Context, string, string, *ManagementLocksGetAtResourceGroupLevelOptions)` to `(context.Context, string, string, *ManagementLocksClientGetAtResourceGroupLevelOptions)`
- Function `*ManagementLocksClient.GetAtResourceGroupLevel` return value(s) have been changed from `(ManagementLocksGetAtResourceGroupLevelResponse, error)` to `(ManagementLocksClientGetAtResourceGroupLevelResponse, error)`
- Function `*ManagementLocksClient.ListByScope` parameter(s) have been changed from `(string, *ManagementLocksListByScopeOptions)` to `(string, *ManagementLocksClientListByScopeOptions)`
- Function `*ManagementLocksClient.ListByScope` return value(s) have been changed from `(*ManagementLocksListByScopePager)` to `(*ManagementLocksClientListByScopePager)`
- Function `*AuthorizationOperationsClient.List` parameter(s) have been changed from `(*AuthorizationOperationsListOptions)` to `(*AuthorizationOperationsClientListOptions)`
- Function `*AuthorizationOperationsClient.List` return value(s) have been changed from `(*AuthorizationOperationsListPager)` to `(*AuthorizationOperationsClientListPager)`
- Function `*ManagementLocksClient.GetAtResourceLevel` parameter(s) have been changed from `(context.Context, string, string, string, string, string, string, *ManagementLocksGetAtResourceLevelOptions)` to `(context.Context, string, string, string, string, string, string, *ManagementLocksClientGetAtResourceLevelOptions)`
- Function `*ManagementLocksClient.GetAtResourceLevel` return value(s) have been changed from `(ManagementLocksGetAtResourceLevelResponse, error)` to `(ManagementLocksClientGetAtResourceLevelResponse, error)`
- Function `*ManagementLocksClient.DeleteAtResourceGroupLevel` parameter(s) have been changed from `(context.Context, string, string, *ManagementLocksDeleteAtResourceGroupLevelOptions)` to `(context.Context, string, string, *ManagementLocksClientDeleteAtResourceGroupLevelOptions)`
- Function `*ManagementLocksClient.DeleteAtResourceGroupLevel` return value(s) have been changed from `(ManagementLocksDeleteAtResourceGroupLevelResponse, error)` to `(ManagementLocksClientDeleteAtResourceGroupLevelResponse, error)`
- Function `*ManagementLocksClient.DeleteAtSubscriptionLevel` parameter(s) have been changed from `(context.Context, string, *ManagementLocksDeleteAtSubscriptionLevelOptions)` to `(context.Context, string, *ManagementLocksClientDeleteAtSubscriptionLevelOptions)`
- Function `*ManagementLocksClient.DeleteAtSubscriptionLevel` return value(s) have been changed from `(ManagementLocksDeleteAtSubscriptionLevelResponse, error)` to `(ManagementLocksClientDeleteAtSubscriptionLevelResponse, error)`
- Function `*ManagementLocksClient.GetByScope` parameter(s) have been changed from `(context.Context, string, string, *ManagementLocksGetByScopeOptions)` to `(context.Context, string, string, *ManagementLocksClientGetByScopeOptions)`
- Function `*ManagementLocksClient.GetByScope` return value(s) have been changed from `(ManagementLocksGetByScopeResponse, error)` to `(ManagementLocksClientGetByScopeResponse, error)`
- Function `*ManagementLocksClient.CreateOrUpdateByScope` parameter(s) have been changed from `(context.Context, string, string, ManagementLockObject, *ManagementLocksCreateOrUpdateByScopeOptions)` to `(context.Context, string, string, ManagementLockObject, *ManagementLocksClientCreateOrUpdateByScopeOptions)`
- Function `*ManagementLocksClient.CreateOrUpdateByScope` return value(s) have been changed from `(ManagementLocksCreateOrUpdateByScopeResponse, error)` to `(ManagementLocksClientCreateOrUpdateByScopeResponse, error)`
- Function `*ManagementLocksClient.CreateOrUpdateAtSubscriptionLevel` parameter(s) have been changed from `(context.Context, string, ManagementLockObject, *ManagementLocksCreateOrUpdateAtSubscriptionLevelOptions)` to `(context.Context, string, ManagementLockObject, *ManagementLocksClientCreateOrUpdateAtSubscriptionLevelOptions)`
- Function `*ManagementLocksClient.CreateOrUpdateAtSubscriptionLevel` return value(s) have been changed from `(ManagementLocksCreateOrUpdateAtSubscriptionLevelResponse, error)` to `(ManagementLocksClientCreateOrUpdateAtSubscriptionLevelResponse, error)`
- Function `*ManagementLocksClient.ListAtResourceLevel` parameter(s) have been changed from `(string, string, string, string, string, *ManagementLocksListAtResourceLevelOptions)` to `(string, string, string, string, string, *ManagementLocksClientListAtResourceLevelOptions)`
- Function `*ManagementLocksClient.ListAtResourceLevel` return value(s) have been changed from `(*ManagementLocksListAtResourceLevelPager)` to `(*ManagementLocksClientListAtResourceLevelPager)`
- Function `*ManagementLocksClient.CreateOrUpdateAtResourceLevel` parameter(s) have been changed from `(context.Context, string, string, string, string, string, string, ManagementLockObject, *ManagementLocksCreateOrUpdateAtResourceLevelOptions)` to `(context.Context, string, string, string, string, string, string, ManagementLockObject, *ManagementLocksClientCreateOrUpdateAtResourceLevelOptions)`
- Function `*ManagementLocksClient.CreateOrUpdateAtResourceLevel` return value(s) have been changed from `(ManagementLocksCreateOrUpdateAtResourceLevelResponse, error)` to `(ManagementLocksClientCreateOrUpdateAtResourceLevelResponse, error)`
- Function `*ManagementLocksClient.CreateOrUpdateAtResourceGroupLevel` parameter(s) have been changed from `(context.Context, string, string, ManagementLockObject, *ManagementLocksCreateOrUpdateAtResourceGroupLevelOptions)` to `(context.Context, string, string, ManagementLockObject, *ManagementLocksClientCreateOrUpdateAtResourceGroupLevelOptions)`
- Function `*ManagementLocksClient.CreateOrUpdateAtResourceGroupLevel` return value(s) have been changed from `(ManagementLocksCreateOrUpdateAtResourceGroupLevelResponse, error)` to `(ManagementLocksClientCreateOrUpdateAtResourceGroupLevelResponse, error)`
- Function `*ManagementLocksClient.DeleteAtResourceLevel` parameter(s) have been changed from `(context.Context, string, string, string, string, string, string, *ManagementLocksDeleteAtResourceLevelOptions)` to `(context.Context, string, string, string, string, string, string, *ManagementLocksClientDeleteAtResourceLevelOptions)`
- Function `*ManagementLocksClient.DeleteAtResourceLevel` return value(s) have been changed from `(ManagementLocksDeleteAtResourceLevelResponse, error)` to `(ManagementLocksClientDeleteAtResourceLevelResponse, error)`
- Function `*ManagementLocksClient.GetAtSubscriptionLevel` parameter(s) have been changed from `(context.Context, string, *ManagementLocksGetAtSubscriptionLevelOptions)` to `(context.Context, string, *ManagementLocksClientGetAtSubscriptionLevelOptions)`
- Function `*ManagementLocksClient.GetAtSubscriptionLevel` return value(s) have been changed from `(ManagementLocksGetAtSubscriptionLevelResponse, error)` to `(ManagementLocksClientGetAtSubscriptionLevelResponse, error)`
- Function `*AuthorizationOperationsListPager.NextPage` has been removed
- Function `*ManagementLocksListByScopePager.Err` has been removed
- Function `*ManagementLocksListAtResourceGroupLevelPager.Err` has been removed
- Function `*AuthorizationOperationsListPager.Err` has been removed
- Function `*ManagementLocksListByScopePager.PageResponse` has been removed
- Function `*ManagementLocksListAtSubscriptionLevelPager.PageResponse` has been removed
- Function `*ManagementLocksListAtResourceLevelPager.NextPage` has been removed
- Function `*ManagementLocksListAtResourceGroupLevelPager.NextPage` has been removed
- Function `*ManagementLocksListAtResourceLevelPager.Err` has been removed
- Function `*ManagementLocksListAtSubscriptionLevelPager.Err` has been removed
- Function `*ManagementLocksListAtSubscriptionLevelPager.NextPage` has been removed
- Function `*ManagementLocksListAtResourceGroupLevelPager.PageResponse` has been removed
- Function `*ManagementLocksListAtResourceLevelPager.PageResponse` has been removed
- Function `*ManagementLocksListByScopePager.NextPage` has been removed
- Function `*AuthorizationOperationsListPager.PageResponse` has been removed
- Function `ErrorResponse.Error` has been removed
- Struct `AuthorizationOperationsListOptions` has been removed
- Struct `AuthorizationOperationsListPager` has been removed
- Struct `AuthorizationOperationsListResponse` has been removed
- Struct `AuthorizationOperationsListResult` has been removed
- Struct `ManagementLocksCreateOrUpdateAtResourceGroupLevelOptions` has been removed
- Struct `ManagementLocksCreateOrUpdateAtResourceGroupLevelResponse` has been removed
- Struct `ManagementLocksCreateOrUpdateAtResourceGroupLevelResult` has been removed
- Struct `ManagementLocksCreateOrUpdateAtResourceLevelOptions` has been removed
- Struct `ManagementLocksCreateOrUpdateAtResourceLevelResponse` has been removed
- Struct `ManagementLocksCreateOrUpdateAtResourceLevelResult` has been removed
- Struct `ManagementLocksCreateOrUpdateAtSubscriptionLevelOptions` has been removed
- Struct `ManagementLocksCreateOrUpdateAtSubscriptionLevelResponse` has been removed
- Struct `ManagementLocksCreateOrUpdateAtSubscriptionLevelResult` has been removed
- Struct `ManagementLocksCreateOrUpdateByScopeOptions` has been removed
- Struct `ManagementLocksCreateOrUpdateByScopeResponse` has been removed
- Struct `ManagementLocksCreateOrUpdateByScopeResult` has been removed
- Struct `ManagementLocksDeleteAtResourceGroupLevelOptions` has been removed
- Struct `ManagementLocksDeleteAtResourceGroupLevelResponse` has been removed
- Struct `ManagementLocksDeleteAtResourceLevelOptions` has been removed
- Struct `ManagementLocksDeleteAtResourceLevelResponse` has been removed
- Struct `ManagementLocksDeleteAtSubscriptionLevelOptions` has been removed
- Struct `ManagementLocksDeleteAtSubscriptionLevelResponse` has been removed
- Struct `ManagementLocksDeleteByScopeOptions` has been removed
- Struct `ManagementLocksDeleteByScopeResponse` has been removed
- Struct `ManagementLocksGetAtResourceGroupLevelOptions` has been removed
- Struct `ManagementLocksGetAtResourceGroupLevelResponse` has been removed
- Struct `ManagementLocksGetAtResourceGroupLevelResult` has been removed
- Struct `ManagementLocksGetAtResourceLevelOptions` has been removed
- Struct `ManagementLocksGetAtResourceLevelResponse` has been removed
- Struct `ManagementLocksGetAtResourceLevelResult` has been removed
- Struct `ManagementLocksGetAtSubscriptionLevelOptions` has been removed
- Struct `ManagementLocksGetAtSubscriptionLevelResponse` has been removed
- Struct `ManagementLocksGetAtSubscriptionLevelResult` has been removed
- Struct `ManagementLocksGetByScopeOptions` has been removed
- Struct `ManagementLocksGetByScopeResponse` has been removed
- Struct `ManagementLocksGetByScopeResult` has been removed
- Struct `ManagementLocksListAtResourceGroupLevelOptions` has been removed
- Struct `ManagementLocksListAtResourceGroupLevelPager` has been removed
- Struct `ManagementLocksListAtResourceGroupLevelResponse` has been removed
- Struct `ManagementLocksListAtResourceGroupLevelResult` has been removed
- Struct `ManagementLocksListAtResourceLevelOptions` has been removed
- Struct `ManagementLocksListAtResourceLevelPager` has been removed
- Struct `ManagementLocksListAtResourceLevelResponse` has been removed
- Struct `ManagementLocksListAtResourceLevelResult` has been removed
- Struct `ManagementLocksListAtSubscriptionLevelOptions` has been removed
- Struct `ManagementLocksListAtSubscriptionLevelPager` has been removed
- Struct `ManagementLocksListAtSubscriptionLevelResponse` has been removed
- Struct `ManagementLocksListAtSubscriptionLevelResult` has been removed
- Struct `ManagementLocksListByScopeOptions` has been removed
- Struct `ManagementLocksListByScopePager` has been removed
- Struct `ManagementLocksListByScopeResponse` has been removed
- Struct `ManagementLocksListByScopeResult` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed

### Features Added

- New function `*ManagementLocksClientListAtResourceGroupLevelPager.PageResponse() ManagementLocksClientListAtResourceGroupLevelResponse`
- New function `*ManagementLocksClientListAtResourceGroupLevelPager.NextPage(context.Context) bool`
- New function `*ManagementLocksClientListAtSubscriptionLevelPager.Err() error`
- New function `*AuthorizationOperationsClientListPager.NextPage(context.Context) bool`
- New function `*AuthorizationOperationsClientListPager.PageResponse() AuthorizationOperationsClientListResponse`
- New function `*ManagementLocksClientListByScopePager.PageResponse() ManagementLocksClientListByScopeResponse`
- New function `*ManagementLocksClientListAtResourceLevelPager.Err() error`
- New function `*AuthorizationOperationsClientListPager.Err() error`
- New function `*ManagementLocksClientListAtSubscriptionLevelPager.NextPage(context.Context) bool`
- New function `*ManagementLocksClientListAtResourceLevelPager.NextPage(context.Context) bool`
- New function `*ManagementLocksClientListAtResourceLevelPager.PageResponse() ManagementLocksClientListAtResourceLevelResponse`
- New function `*ManagementLocksClientListByScopePager.NextPage(context.Context) bool`
- New function `*ManagementLocksClientListAtSubscriptionLevelPager.PageResponse() ManagementLocksClientListAtSubscriptionLevelResponse`
- New function `*ManagementLocksClientListByScopePager.Err() error`
- New function `*ManagementLocksClientListAtResourceGroupLevelPager.Err() error`
- New struct `AuthorizationOperationsClientListOptions`
- New struct `AuthorizationOperationsClientListPager`
- New struct `AuthorizationOperationsClientListResponse`
- New struct `AuthorizationOperationsClientListResult`
- New struct `ManagementLocksClientCreateOrUpdateAtResourceGroupLevelOptions`
- New struct `ManagementLocksClientCreateOrUpdateAtResourceGroupLevelResponse`
- New struct `ManagementLocksClientCreateOrUpdateAtResourceGroupLevelResult`
- New struct `ManagementLocksClientCreateOrUpdateAtResourceLevelOptions`
- New struct `ManagementLocksClientCreateOrUpdateAtResourceLevelResponse`
- New struct `ManagementLocksClientCreateOrUpdateAtResourceLevelResult`
- New struct `ManagementLocksClientCreateOrUpdateAtSubscriptionLevelOptions`
- New struct `ManagementLocksClientCreateOrUpdateAtSubscriptionLevelResponse`
- New struct `ManagementLocksClientCreateOrUpdateAtSubscriptionLevelResult`
- New struct `ManagementLocksClientCreateOrUpdateByScopeOptions`
- New struct `ManagementLocksClientCreateOrUpdateByScopeResponse`
- New struct `ManagementLocksClientCreateOrUpdateByScopeResult`
- New struct `ManagementLocksClientDeleteAtResourceGroupLevelOptions`
- New struct `ManagementLocksClientDeleteAtResourceGroupLevelResponse`
- New struct `ManagementLocksClientDeleteAtResourceLevelOptions`
- New struct `ManagementLocksClientDeleteAtResourceLevelResponse`
- New struct `ManagementLocksClientDeleteAtSubscriptionLevelOptions`
- New struct `ManagementLocksClientDeleteAtSubscriptionLevelResponse`
- New struct `ManagementLocksClientDeleteByScopeOptions`
- New struct `ManagementLocksClientDeleteByScopeResponse`
- New struct `ManagementLocksClientGetAtResourceGroupLevelOptions`
- New struct `ManagementLocksClientGetAtResourceGroupLevelResponse`
- New struct `ManagementLocksClientGetAtResourceGroupLevelResult`
- New struct `ManagementLocksClientGetAtResourceLevelOptions`
- New struct `ManagementLocksClientGetAtResourceLevelResponse`
- New struct `ManagementLocksClientGetAtResourceLevelResult`
- New struct `ManagementLocksClientGetAtSubscriptionLevelOptions`
- New struct `ManagementLocksClientGetAtSubscriptionLevelResponse`
- New struct `ManagementLocksClientGetAtSubscriptionLevelResult`
- New struct `ManagementLocksClientGetByScopeOptions`
- New struct `ManagementLocksClientGetByScopeResponse`
- New struct `ManagementLocksClientGetByScopeResult`
- New struct `ManagementLocksClientListAtResourceGroupLevelOptions`
- New struct `ManagementLocksClientListAtResourceGroupLevelPager`
- New struct `ManagementLocksClientListAtResourceGroupLevelResponse`
- New struct `ManagementLocksClientListAtResourceGroupLevelResult`
- New struct `ManagementLocksClientListAtResourceLevelOptions`
- New struct `ManagementLocksClientListAtResourceLevelPager`
- New struct `ManagementLocksClientListAtResourceLevelResponse`
- New struct `ManagementLocksClientListAtResourceLevelResult`
- New struct `ManagementLocksClientListAtSubscriptionLevelOptions`
- New struct `ManagementLocksClientListAtSubscriptionLevelPager`
- New struct `ManagementLocksClientListAtSubscriptionLevelResponse`
- New struct `ManagementLocksClientListAtSubscriptionLevelResult`
- New struct `ManagementLocksClientListByScopeOptions`
- New struct `ManagementLocksClientListByScopePager`
- New struct `ManagementLocksClientListByScopeResponse`
- New struct `ManagementLocksClientListByScopeResult`
- New field `Error` in struct `ErrorResponse`


## 0.1.1 (2021-12-13)

### Other Changes

- Fix the go minimum version to `1.16`

## 0.1.0 (2021-11-16)

- Initial preview release.
