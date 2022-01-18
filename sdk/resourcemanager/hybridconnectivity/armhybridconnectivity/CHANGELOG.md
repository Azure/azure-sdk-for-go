# Release History

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*EndpointsClient.Update` parameter(s) have been changed from `(context.Context, string, string, EndpointResource, *EndpointsUpdateOptions)` to `(context.Context, string, string, EndpointResource, *EndpointsClientUpdateOptions)`
- Function `*EndpointsClient.Update` return value(s) have been changed from `(EndpointsUpdateResponse, error)` to `(EndpointsClientUpdateResponse, error)`
- Function `*EndpointsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, EndpointResource, *EndpointsCreateOrUpdateOptions)` to `(context.Context, string, string, EndpointResource, *EndpointsClientCreateOrUpdateOptions)`
- Function `*EndpointsClient.CreateOrUpdate` return value(s) have been changed from `(EndpointsCreateOrUpdateResponse, error)` to `(EndpointsClientCreateOrUpdateResponse, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*EndpointsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *EndpointsGetOptions)` to `(context.Context, string, string, *EndpointsClientGetOptions)`
- Function `*EndpointsClient.Get` return value(s) have been changed from `(EndpointsGetResponse, error)` to `(EndpointsClientGetResponse, error)`
- Function `*EndpointsClient.ListCredentials` parameter(s) have been changed from `(context.Context, string, string, *EndpointsListCredentialsOptions)` to `(context.Context, string, string, *EndpointsClientListCredentialsOptions)`
- Function `*EndpointsClient.ListCredentials` return value(s) have been changed from `(EndpointsListCredentialsResponse, error)` to `(EndpointsClientListCredentialsResponse, error)`
- Function `*EndpointsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *EndpointsDeleteOptions)` to `(context.Context, string, string, *EndpointsClientDeleteOptions)`
- Function `*EndpointsClient.Delete` return value(s) have been changed from `(EndpointsDeleteResponse, error)` to `(EndpointsClientDeleteResponse, error)`
- Function `*EndpointsClient.List` parameter(s) have been changed from `(string, *EndpointsListOptions)` to `(string, *EndpointsClientListOptions)`
- Function `*EndpointsClient.List` return value(s) have been changed from `(*EndpointsListPager)` to `(*EndpointsClientListPager)`
- Function `ErrorResponse.Error` has been removed
- Function `*OperationsListPager.PageResponse` has been removed
- Function `*EndpointsListPager.PageResponse` has been removed
- Function `Resource.MarshalJSON` has been removed
- Function `*EndpointsListPager.Err` has been removed
- Function `*EndpointsListPager.NextPage` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Function `*OperationsListPager.Err` has been removed
- Struct `EndpointsCreateOrUpdateOptions` has been removed
- Struct `EndpointsCreateOrUpdateResponse` has been removed
- Struct `EndpointsCreateOrUpdateResult` has been removed
- Struct `EndpointsDeleteOptions` has been removed
- Struct `EndpointsDeleteResponse` has been removed
- Struct `EndpointsGetOptions` has been removed
- Struct `EndpointsGetResponse` has been removed
- Struct `EndpointsGetResult` has been removed
- Struct `EndpointsListCredentialsOptions` has been removed
- Struct `EndpointsListCredentialsResponse` has been removed
- Struct `EndpointsListCredentialsResult` has been removed
- Struct `EndpointsListOptions` has been removed
- Struct `EndpointsListPager` has been removed
- Struct `EndpointsListResponse` has been removed
- Struct `EndpointsListResult` has been removed
- Struct `EndpointsUpdateOptions` has been removed
- Struct `EndpointsUpdateResponse` has been removed
- Struct `EndpointsUpdateResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `Resource` of struct `ProxyResource` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed
- Field `ProxyResource` of struct `EndpointResource` has been removed

### Features Added

- New function `*EndpointsClientListPager.PageResponse() EndpointsClientListResponse`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*EndpointsClientListPager.Err() error`
- New function `*OperationsClientListPager.Err() error`
- New function `*EndpointsClientListPager.NextPage(context.Context) bool`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New struct `EndpointsClientCreateOrUpdateOptions`
- New struct `EndpointsClientCreateOrUpdateResponse`
- New struct `EndpointsClientCreateOrUpdateResult`
- New struct `EndpointsClientDeleteOptions`
- New struct `EndpointsClientDeleteResponse`
- New struct `EndpointsClientGetOptions`
- New struct `EndpointsClientGetResponse`
- New struct `EndpointsClientGetResult`
- New struct `EndpointsClientListCredentialsOptions`
- New struct `EndpointsClientListCredentialsResponse`
- New struct `EndpointsClientListCredentialsResult`
- New struct `EndpointsClientListOptions`
- New struct `EndpointsClientListPager`
- New struct `EndpointsClientListResponse`
- New struct `EndpointsClientListResult`
- New struct `EndpointsClientUpdateOptions`
- New struct `EndpointsClientUpdateResponse`
- New struct `EndpointsClientUpdateResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New field `ID` in struct `ProxyResource`
- New field `Name` in struct `ProxyResource`
- New field `Type` in struct `ProxyResource`
- New field `Error` in struct `ErrorResponse`
- New field `Name` in struct `EndpointResource`
- New field `Type` in struct `EndpointResource`
- New field `ID` in struct `EndpointResource`


## 0.1.0 (2021-12-07)

- Init release.
