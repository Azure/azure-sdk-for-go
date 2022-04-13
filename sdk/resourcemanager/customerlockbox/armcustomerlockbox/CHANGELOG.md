# Release History

## 0.3.0 (2022-04-11)
### Breaking Changes

- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `NewRequestsClient` return value(s) have been changed from `(*RequestsClient)` to `(*RequestsClient, error)`
- Function `NewGetClient` return value(s) have been changed from `(*GetClient)` to `(*GetClient, error)`
- Function `NewPostClient` return value(s) have been changed from `(*PostClient)` to `(*PostClient, error)`
- Function `*RequestsClient.List` return value(s) have been changed from `(*RequestsClientListPager)` to `(*runtime.Pager[RequestsClientListResponse])`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `Status.ToPtr` has been removed
- Function `*OperationsClientListPager.NextPage` has been removed
- Function `*OperationsClientListPager.Err` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Function `*RequestsClientListPager.Err` has been removed
- Function `*RequestsClientListPager.PageResponse` has been removed
- Function `*RequestsClientListPager.NextPage` has been removed
- Struct `GetClientTenantOptedInResult` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Struct `RequestsClientGetResult` has been removed
- Struct `RequestsClientListPager` has been removed
- Struct `RequestsClientListResult` has been removed
- Struct `RequestsClientUpdateStatusResult` has been removed
- Field `GetClientTenantOptedInResult` of struct `GetClientTenantOptedInResponse` has been removed
- Field `RawResponse` of struct `GetClientTenantOptedInResponse` has been removed
- Field `RawResponse` of struct `PostClientDisableLockboxResponse` has been removed
- Field `RawResponse` of struct `PostClientEnableLockboxResponse` has been removed
- Field `RequestsClientListResult` of struct `RequestsClientListResponse` has been removed
- Field `RawResponse` of struct `RequestsClientListResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `RequestsClientUpdateStatusResult` of struct `RequestsClientUpdateStatusResponse` has been removed
- Field `RawResponse` of struct `RequestsClientUpdateStatusResponse` has been removed
- Field `RequestsClientGetResult` of struct `RequestsClientGetResponse` has been removed
- Field `RawResponse` of struct `RequestsClientGetResponse` has been removed

### Features Added

- New anonymous field `LockboxRequestResponse` in struct `RequestsClientGetResponse`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New anonymous field `Approval` in struct `RequestsClientUpdateStatusResponse`
- New anonymous field `RequestListResult` in struct `RequestsClientListResponse`
- New anonymous field `TenantOptInResponse` in struct `GetClientTenantOptedInResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*RequestsClient.List` parameter(s) have been changed from `(string, *RequestsListOptions)` to `(string, *RequestsClientListOptions)`
- Function `*RequestsClient.List` return value(s) have been changed from `(*RequestsListPager)` to `(*RequestsClientListPager)`
- Function `*RequestsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *RequestsGetOptions)` to `(context.Context, string, string, *RequestsClientGetOptions)`
- Function `*RequestsClient.Get` return value(s) have been changed from `(RequestsGetResponse, error)` to `(RequestsClientGetResponse, error)`
- Function `*GetClient.TenantOptedIn` parameter(s) have been changed from `(context.Context, string, *GetTenantOptedInOptions)` to `(context.Context, string, *GetClientTenantOptedInOptions)`
- Function `*GetClient.TenantOptedIn` return value(s) have been changed from `(GetTenantOptedInResponse, error)` to `(GetClientTenantOptedInResponse, error)`
- Function `*PostClient.EnableLockbox` parameter(s) have been changed from `(context.Context, *PostEnableLockboxOptions)` to `(context.Context, *PostClientEnableLockboxOptions)`
- Function `*PostClient.EnableLockbox` return value(s) have been changed from `(PostEnableLockboxResponse, error)` to `(PostClientEnableLockboxResponse, error)`
- Function `*PostClient.DisableLockbox` parameter(s) have been changed from `(context.Context, *PostDisableLockboxOptions)` to `(context.Context, *PostClientDisableLockboxOptions)`
- Function `*PostClient.DisableLockbox` return value(s) have been changed from `(PostDisableLockboxResponse, error)` to `(PostClientDisableLockboxResponse, error)`
- Function `*RequestsClient.UpdateStatus` parameter(s) have been changed from `(context.Context, string, string, Approval, *RequestsUpdateStatusOptions)` to `(context.Context, string, string, Approval, *RequestsClientUpdateStatusOptions)`
- Function `*RequestsClient.UpdateStatus` return value(s) have been changed from `(RequestsUpdateStatusResponse, error)` to `(RequestsClientUpdateStatusResponse, error)`
- Function `*OperationsListPager.PageResponse` has been removed
- Function `*RequestsListPager.Err` has been removed
- Function `*RequestsListPager.NextPage` has been removed
- Function `*RequestsListPager.PageResponse` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `*OperationsListPager.Err` has been removed
- Struct `GetTenantOptedInOptions` has been removed
- Struct `GetTenantOptedInResponse` has been removed
- Struct `GetTenantOptedInResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Struct `PostDisableLockboxOptions` has been removed
- Struct `PostDisableLockboxResponse` has been removed
- Struct `PostEnableLockboxOptions` has been removed
- Struct `PostEnableLockboxResponse` has been removed
- Struct `RequestsGetOptions` has been removed
- Struct `RequestsGetResponse` has been removed
- Struct `RequestsGetResult` has been removed
- Struct `RequestsListOptions` has been removed
- Struct `RequestsListPager` has been removed
- Struct `RequestsListResponse` has been removed
- Struct `RequestsListResult` has been removed
- Struct `RequestsUpdateStatusOptions` has been removed
- Struct `RequestsUpdateStatusResponse` has been removed
- Struct `RequestsUpdateStatusResult` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed

### Features Added

- New function `*RequestsClientListPager.NextPage(context.Context) bool`
- New function `*RequestsClientListPager.Err() error`
- New function `*RequestsClientListPager.PageResponse() RequestsClientListResponse`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*OperationsClientListPager.Err() error`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New struct `GetClientTenantOptedInOptions`
- New struct `GetClientTenantOptedInResponse`
- New struct `GetClientTenantOptedInResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New struct `PostClientDisableLockboxOptions`
- New struct `PostClientDisableLockboxResponse`
- New struct `PostClientEnableLockboxOptions`
- New struct `PostClientEnableLockboxResponse`
- New struct `RequestsClientGetOptions`
- New struct `RequestsClientGetResponse`
- New struct `RequestsClientGetResult`
- New struct `RequestsClientListOptions`
- New struct `RequestsClientListPager`
- New struct `RequestsClientListResponse`
- New struct `RequestsClientListResult`
- New struct `RequestsClientUpdateStatusOptions`
- New struct `RequestsClientUpdateStatusResponse`
- New struct `RequestsClientUpdateStatusResult`
- New field `Error` in struct `ErrorResponse`


## 0.1.0 (2021-12-07)

- Initial preview release.
