# Release History

## 0.5.0 (2022-04-20)
### Features Added

- New function `*EndpointsClient.ListManagedProxyDetails(context.Context, string, string, ManagedProxyRequest, *EndpointsClientListManagedProxyDetailsOptions) (EndpointsClientListManagedProxyDetailsResponse, error)`
- New struct `AADProfileProperties`
- New struct `EndpointsClientListManagedProxyDetailsOptions`
- New struct `EndpointsClientListManagedProxyDetailsResponse`
- New struct `IngressGatewayResource`
- New struct `IngressProfileProperties`
- New struct `ManagedProxyRequest`
- New struct `ManagedProxyResource`


## 0.4.0 (2022-04-15)
### Breaking Changes

- Function `*EndpointsClient.List` has been removed
- Function `*OperationsClient.List` has been removed

### Features Added

- New function `*EndpointsClient.NewListPager(string, *EndpointsClientListOptions) *runtime.Pager[EndpointsClientListResponse]`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`


## 0.3.0 (2022-04-11)
### Breaking Changes

- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `*EndpointsClient.List` return value(s) have been changed from `(*EndpointsClientListPager)` to `(*runtime.Pager[EndpointsClientListResponse])`
- Function `NewEndpointsClient` return value(s) have been changed from `(*EndpointsClient)` to `(*EndpointsClient, error)`
- Type of `ErrorAdditionalInfo.Info` has been changed from `map[string]interface{}` to `interface{}`
- Function `*OperationsClientListPager.Err` has been removed
- Function `CreatedByType.ToPtr` has been removed
- Function `*EndpointsClientListPager.NextPage` has been removed
- Function `Origin.ToPtr` has been removed
- Function `Type.ToPtr` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Function `*EndpointsClientListPager.PageResponse` has been removed
- Function `ActionType.ToPtr` has been removed
- Function `*EndpointsClientListPager.Err` has been removed
- Function `*OperationsClientListPager.NextPage` has been removed
- Struct `EndpointsClientCreateOrUpdateResult` has been removed
- Struct `EndpointsClientGetResult` has been removed
- Struct `EndpointsClientListCredentialsResult` has been removed
- Struct `EndpointsClientListPager` has been removed
- Struct `EndpointsClientListResult` has been removed
- Struct `EndpointsClientUpdateResult` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Field `RawResponse` of struct `EndpointsClientDeleteResponse` has been removed
- Field `EndpointsClientListCredentialsResult` of struct `EndpointsClientListCredentialsResponse` has been removed
- Field `RawResponse` of struct `EndpointsClientListCredentialsResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `EndpointsClientUpdateResult` of struct `EndpointsClientUpdateResponse` has been removed
- Field `RawResponse` of struct `EndpointsClientUpdateResponse` has been removed
- Field `EndpointsClientCreateOrUpdateResult` of struct `EndpointsClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `EndpointsClientCreateOrUpdateResponse` has been removed
- Field `EndpointsClientListResult` of struct `EndpointsClientListResponse` has been removed
- Field `RawResponse` of struct `EndpointsClientListResponse` has been removed
- Field `EndpointsClientGetResult` of struct `EndpointsClientGetResponse` has been removed
- Field `RawResponse` of struct `EndpointsClientGetResponse` has been removed

### Features Added

- New anonymous field `EndpointAccessResource` in struct `EndpointsClientListCredentialsResponse`
- New anonymous field `EndpointResource` in struct `EndpointsClientUpdateResponse`
- New anonymous field `EndpointResource` in struct `EndpointsClientGetResponse`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New anonymous field `EndpointsList` in struct `EndpointsClientListResponse`
- New anonymous field `EndpointResource` in struct `EndpointsClientCreateOrUpdateResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

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
