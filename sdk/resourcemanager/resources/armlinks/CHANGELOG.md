# Release History

## 0.4.1 (2022-04-18)
### Other Changes


## 0.4.0 (2022-04-18)
### Breaking Changes

- Function `*ResourceLinksClient.ListAtSubscription` has been removed
- Function `*OperationsClient.List` has been removed
- Function `*ResourceLinksClient.ListAtSourceScope` has been removed

### Features Added

- New function `*ResourceLinksClient.NewListAtSourceScopePager(string, *ResourceLinksClientListAtSourceScopeOptions) *runtime.Pager[ResourceLinksClientListAtSourceScopeResponse]`
- New function `*ResourceLinksClient.NewListAtSubscriptionPager(*ResourceLinksClientListAtSubscriptionOptions) *runtime.Pager[ResourceLinksClientListAtSubscriptionResponse]`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`


## 0.3.0 (2022-04-14)
### Breaking Changes

- Function `NewResourceLinksClient` return value(s) have been changed from `(*ResourceLinksClient)` to `(*ResourceLinksClient, error)`
- Function `*ResourceLinksClient.ListAtSourceScope` return value(s) have been changed from `(*ResourceLinksClientListAtSourceScopePager)` to `(*runtime.Pager[ResourceLinksClientListAtSourceScopeResponse])`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `*ResourceLinksClient.ListAtSubscription` return value(s) have been changed from `(*ResourceLinksClientListAtSubscriptionPager)` to `(*runtime.Pager[ResourceLinksClientListAtSubscriptionResponse])`
- Type of `ResourceLink.Type` has been changed from `map[string]interface{}` to `interface{}`
- Function `*ResourceLinksClientListAtSourceScopePager.NextPage` has been removed
- Function `*ResourceLinksClientListAtSubscriptionPager.Err` has been removed
- Function `*ResourceLinksClientListAtSubscriptionPager.PageResponse` has been removed
- Function `*OperationsClientListPager.Err` has been removed
- Function `*ResourceLinksClientListAtSourceScopePager.Err` has been removed
- Function `*ResourceLinksClientListAtSourceScopePager.PageResponse` has been removed
- Function `*ResourceLinksClientListAtSubscriptionPager.NextPage` has been removed
- Function `*OperationsClientListPager.NextPage` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Struct `ResourceLinksClientCreateOrUpdateResult` has been removed
- Struct `ResourceLinksClientGetResult` has been removed
- Struct `ResourceLinksClientListAtSourceScopePager` has been removed
- Struct `ResourceLinksClientListAtSourceScopeResult` has been removed
- Struct `ResourceLinksClientListAtSubscriptionPager` has been removed
- Struct `ResourceLinksClientListAtSubscriptionResult` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `ResourceLinksClientCreateOrUpdateResult` of struct `ResourceLinksClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `ResourceLinksClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `ResourceLinksClientDeleteResponse` has been removed
- Field `ResourceLinksClientListAtSubscriptionResult` of struct `ResourceLinksClientListAtSubscriptionResponse` has been removed
- Field `RawResponse` of struct `ResourceLinksClientListAtSubscriptionResponse` has been removed
- Field `ResourceLinksClientGetResult` of struct `ResourceLinksClientGetResponse` has been removed
- Field `RawResponse` of struct `ResourceLinksClientGetResponse` has been removed
- Field `ResourceLinksClientListAtSourceScopeResult` of struct `ResourceLinksClientListAtSourceScopeResponse` has been removed
- Field `RawResponse` of struct `ResourceLinksClientListAtSourceScopeResponse` has been removed

### Features Added

- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New anonymous field `ResourceLinkResult` in struct `ResourceLinksClientListAtSourceScopeResponse`
- New anonymous field `ResourceLink` in struct `ResourceLinksClientCreateOrUpdateResponse`
- New anonymous field `ResourceLinkResult` in struct `ResourceLinksClientListAtSubscriptionResponse`
- New anonymous field `ResourceLink` in struct `ResourceLinksClientGetResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*ResourceLinksClient.ListAtSourceScope` parameter(s) have been changed from `(string, *ResourceLinksListAtSourceScopeOptions)` to `(string, *ResourceLinksClientListAtSourceScopeOptions)`
- Function `*ResourceLinksClient.ListAtSourceScope` return value(s) have been changed from `(*ResourceLinksListAtSourceScopePager)` to `(*ResourceLinksClientListAtSourceScopePager)`
- Function `*ResourceLinksClient.Get` parameter(s) have been changed from `(context.Context, string, *ResourceLinksGetOptions)` to `(context.Context, string, *ResourceLinksClientGetOptions)`
- Function `*ResourceLinksClient.Get` return value(s) have been changed from `(ResourceLinksGetResponse, error)` to `(ResourceLinksClientGetResponse, error)`
- Function `*ResourceLinksClient.Delete` parameter(s) have been changed from `(context.Context, string, *ResourceLinksDeleteOptions)` to `(context.Context, string, *ResourceLinksClientDeleteOptions)`
- Function `*ResourceLinksClient.Delete` return value(s) have been changed from `(ResourceLinksDeleteResponse, error)` to `(ResourceLinksClientDeleteResponse, error)`
- Function `*ResourceLinksClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, ResourceLink, *ResourceLinksCreateOrUpdateOptions)` to `(context.Context, string, ResourceLink, *ResourceLinksClientCreateOrUpdateOptions)`
- Function `*ResourceLinksClient.CreateOrUpdate` return value(s) have been changed from `(ResourceLinksCreateOrUpdateResponse, error)` to `(ResourceLinksClientCreateOrUpdateResponse, error)`
- Function `*ResourceLinksClient.ListAtSubscription` parameter(s) have been changed from `(*ResourceLinksListAtSubscriptionOptions)` to `(*ResourceLinksClientListAtSubscriptionOptions)`
- Function `*ResourceLinksClient.ListAtSubscription` return value(s) have been changed from `(*ResourceLinksListAtSubscriptionPager)` to `(*ResourceLinksClientListAtSubscriptionPager)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*OperationsListPager.NextPage` has been removed
- Function `*ResourceLinksListAtSourceScopePager.NextPage` has been removed
- Function `*ResourceLinksListAtSubscriptionPager.PageResponse` has been removed
- Function `*ResourceLinksListAtSubscriptionPager.NextPage` has been removed
- Function `*ResourceLinksListAtSourceScopePager.Err` has been removed
- Function `*ResourceLinksListAtSubscriptionPager.Err` has been removed
- Function `*ResourceLinksListAtSourceScopePager.PageResponse` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `*OperationsListPager.PageResponse` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Struct `ResourceLinksCreateOrUpdateOptions` has been removed
- Struct `ResourceLinksCreateOrUpdateResponse` has been removed
- Struct `ResourceLinksCreateOrUpdateResult` has been removed
- Struct `ResourceLinksDeleteOptions` has been removed
- Struct `ResourceLinksDeleteResponse` has been removed
- Struct `ResourceLinksGetOptions` has been removed
- Struct `ResourceLinksGetResponse` has been removed
- Struct `ResourceLinksGetResult` has been removed
- Struct `ResourceLinksListAtSourceScopeOptions` has been removed
- Struct `ResourceLinksListAtSourceScopePager` has been removed
- Struct `ResourceLinksListAtSourceScopeResponse` has been removed
- Struct `ResourceLinksListAtSourceScopeResult` has been removed
- Struct `ResourceLinksListAtSubscriptionOptions` has been removed
- Struct `ResourceLinksListAtSubscriptionPager` has been removed
- Struct `ResourceLinksListAtSubscriptionResponse` has been removed
- Struct `ResourceLinksListAtSubscriptionResult` has been removed

### Features Added

- New function `*OperationsClientListPager.Err() error`
- New function `*ResourceLinksClientListAtSubscriptionPager.NextPage(context.Context) bool`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*ResourceLinksClientListAtSourceScopePager.Err() error`
- New function `*ResourceLinksClientListAtSourceScopePager.PageResponse() ResourceLinksClientListAtSourceScopeResponse`
- New function `*ResourceLinksClientListAtSourceScopePager.NextPage(context.Context) bool`
- New function `*ResourceLinksClientListAtSubscriptionPager.Err() error`
- New function `*ResourceLinksClientListAtSubscriptionPager.PageResponse() ResourceLinksClientListAtSubscriptionResponse`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New struct `ResourceLinksClientCreateOrUpdateOptions`
- New struct `ResourceLinksClientCreateOrUpdateResponse`
- New struct `ResourceLinksClientCreateOrUpdateResult`
- New struct `ResourceLinksClientDeleteOptions`
- New struct `ResourceLinksClientDeleteResponse`
- New struct `ResourceLinksClientGetOptions`
- New struct `ResourceLinksClientGetResponse`
- New struct `ResourceLinksClientGetResult`
- New struct `ResourceLinksClientListAtSourceScopeOptions`
- New struct `ResourceLinksClientListAtSourceScopePager`
- New struct `ResourceLinksClientListAtSourceScopeResponse`
- New struct `ResourceLinksClientListAtSourceScopeResult`
- New struct `ResourceLinksClientListAtSubscriptionOptions`
- New struct `ResourceLinksClientListAtSubscriptionPager`
- New struct `ResourceLinksClientListAtSubscriptionResponse`
- New struct `ResourceLinksClientListAtSubscriptionResult`


## 0.1.1 (2021-12-13)

### Other Changes

- Fix the go minimum version to `1.16`

## 0.1.0 (2021-11-16)

- Initial preview release.
