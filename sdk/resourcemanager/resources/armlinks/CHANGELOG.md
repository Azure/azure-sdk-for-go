# Release History

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
