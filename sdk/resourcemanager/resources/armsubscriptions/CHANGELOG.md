# Release History

## 0.4.1 (2022-04-18)
### Other Changes


## 0.4.0 (2022-04-18)
### Breaking Changes

- Function `*TenantsClient.List` has been removed
- Function `*Client.List` has been removed
- Function `*Client.ListLocations` has been removed

### Features Added

- New function `*TenantsClient.NewListPager(*TenantsClientListOptions) *runtime.Pager[TenantsClientListResponse]`
- New function `*Client.NewListPager(*ClientListOptions) *runtime.Pager[ClientListResponse]`
- New function `*Client.NewListLocationsPager(string, *ClientListLocationsOptions) *runtime.Pager[ClientListLocationsResponse]`


## 0.3.0 (2022-04-14)
### Breaking Changes

- Function `*TenantsClient.List` return value(s) have been changed from `(*TenantsClientListPager)` to `(*runtime.Pager[TenantsClientListResponse])`
- Function `NewClient` return value(s) have been changed from `(*Client)` to `(*Client, error)`
- Function `*Client.ListLocations` parameter(s) have been changed from `(context.Context, string, *ClientListLocationsOptions)` to `(string, *ClientListLocationsOptions)`
- Function `*Client.ListLocations` return value(s) have been changed from `(ClientListLocationsResponse, error)` to `(*runtime.Pager[ClientListLocationsResponse])`
- Function `*Client.List` return value(s) have been changed from `(*ClientListPager)` to `(*runtime.Pager[ClientListResponse])`
- Function `NewSubscriptionClient` return value(s) have been changed from `(*SubscriptionClient)` to `(*SubscriptionClient, error)`
- Function `NewTenantsClient` return value(s) have been changed from `(*TenantsClient)` to `(*TenantsClient, error)`
- Type of `ErrorAdditionalInfo.Info` has been changed from `map[string]interface{}` to `interface{}`
- Function `ResourceNameStatus.ToPtr` has been removed
- Function `RegionType.ToPtr` has been removed
- Function `*ClientListPager.Err` has been removed
- Function `RegionCategory.ToPtr` has been removed
- Function `SubscriptionState.ToPtr` has been removed
- Function `TenantCategory.ToPtr` has been removed
- Function `*TenantsClientListPager.NextPage` has been removed
- Function `*TenantsClientListPager.PageResponse` has been removed
- Function `*ClientListPager.NextPage` has been removed
- Function `*TenantsClientListPager.Err` has been removed
- Function `SpendingLimit.ToPtr` has been removed
- Function `*ClientListPager.PageResponse` has been removed
- Function `LocationType.ToPtr` has been removed
- Struct `ClientCheckZonePeersResult` has been removed
- Struct `ClientGetResult` has been removed
- Struct `ClientListLocationsResult` has been removed
- Struct `ClientListPager` has been removed
- Struct `ClientListResult` has been removed
- Struct `SubscriptionClientCheckResourceNameResult` has been removed
- Struct `TenantsClientListPager` has been removed
- Struct `TenantsClientListResult` has been removed
- Field `ClientListResult` of struct `ClientListResponse` has been removed
- Field `RawResponse` of struct `ClientListResponse` has been removed
- Field `ClientCheckZonePeersResult` of struct `ClientCheckZonePeersResponse` has been removed
- Field `RawResponse` of struct `ClientCheckZonePeersResponse` has been removed
- Field `ClientGetResult` of struct `ClientGetResponse` has been removed
- Field `RawResponse` of struct `ClientGetResponse` has been removed
- Field `TenantsClientListResult` of struct `TenantsClientListResponse` has been removed
- Field `RawResponse` of struct `TenantsClientListResponse` has been removed
- Field `ClientListLocationsResult` of struct `ClientListLocationsResponse` has been removed
- Field `RawResponse` of struct `ClientListLocationsResponse` has been removed
- Field `SubscriptionClientCheckResourceNameResult` of struct `SubscriptionClientCheckResourceNameResponse` has been removed
- Field `RawResponse` of struct `SubscriptionClientCheckResourceNameResponse` has been removed

### Features Added

- New anonymous field `Subscription` in struct `ClientGetResponse`
- New anonymous field `CheckZonePeersResult` in struct `ClientCheckZonePeersResponse`
- New anonymous field `TenantListResult` in struct `TenantsClientListResponse`
- New anonymous field `SubscriptionListResult` in struct `ClientListResponse`
- New anonymous field `LocationListResult` in struct `ClientListLocationsResponse`
- New anonymous field `CheckResourceNameResult` in struct `SubscriptionClientCheckResourceNameResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-17)
### Breaking Changes

- Function `*TenantsClient.List` parameter(s) have been changed from `(*TenantsListOptions)` to `(*TenantsClientListOptions)`
- Function `*TenantsClient.List` return value(s) have been changed from `(*TenantsListPager)` to `(*TenantsClientListPager)`
- Function `*TenantsListPager.Err` has been removed
- Function `*TenantsListPager.PageResponse` has been removed
- Function `*TenantsListPager.NextPage` has been removed
- Function `*SubscriptionsClient.List` has been removed
- Function `*SubscriptionsClient.ListLocations` has been removed
- Function `NewSubscriptionsClient` has been removed
- Function `*SubscriptionsListPager.PageResponse` has been removed
- Function `*SubscriptionsListPager.NextPage` has been removed
- Function `*SubscriptionsListPager.Err` has been removed
- Function `*SubscriptionsClient.Get` has been removed
- Function `CloudError.Error` has been removed
- Struct `SubscriptionsClient` has been removed
- Struct `SubscriptionsGetOptions` has been removed
- Struct `SubscriptionsGetResponse` has been removed
- Struct `SubscriptionsGetResult` has been removed
- Struct `SubscriptionsListLocationsOptions` has been removed
- Struct `SubscriptionsListLocationsResponse` has been removed
- Struct `SubscriptionsListLocationsResult` has been removed
- Struct `SubscriptionsListOptions` has been removed
- Struct `SubscriptionsListPager` has been removed
- Struct `SubscriptionsListResponse` has been removed
- Struct `SubscriptionsListResult` has been removed
- Struct `TenantsListOptions` has been removed
- Struct `TenantsListPager` has been removed
- Struct `TenantsListResponse` has been removed
- Struct `TenantsListResult` has been removed
- Field `InnerError` of struct `CloudError` has been removed

### Features Added

- New function `ErrorDetail.MarshalJSON() ([]byte, error)`
- New function `*Client.ListLocations(context.Context, string, *ClientListLocationsOptions) (ClientListLocationsResponse, error)`
- New function `NewClient(azcore.TokenCredential, *arm.ClientOptions) *Client`
- New function `*TenantsClientListPager.PageResponse() TenantsClientListResponse`
- New function `CheckZonePeersResult.MarshalJSON() ([]byte, error)`
- New function `*ClientListPager.Err() error`
- New function `*TenantsClientListPager.NextPage(context.Context) bool`
- New function `*Client.CheckZonePeers(context.Context, string, CheckZonePeersRequest, *ClientCheckZonePeersOptions) (ClientCheckZonePeersResponse, error)`
- New function `CheckZonePeersRequest.MarshalJSON() ([]byte, error)`
- New function `*Client.Get(context.Context, string, *ClientGetOptions) (ClientGetResponse, error)`
- New function `*ClientListPager.PageResponse() ClientListResponse`
- New function `AvailabilityZonePeers.MarshalJSON() ([]byte, error)`
- New function `*Client.List(*ClientListOptions) *ClientListPager`
- New function `*ClientListPager.NextPage(context.Context) bool`
- New function `*TenantsClientListPager.Err() error`
- New struct `AvailabilityZonePeers`
- New struct `CheckZonePeersRequest`
- New struct `CheckZonePeersResult`
- New struct `Client`
- New struct `ClientCheckZonePeersOptions`
- New struct `ClientCheckZonePeersResponse`
- New struct `ClientCheckZonePeersResult`
- New struct `ClientGetOptions`
- New struct `ClientGetResponse`
- New struct `ClientGetResult`
- New struct `ClientListLocationsOptions`
- New struct `ClientListLocationsResponse`
- New struct `ClientListLocationsResult`
- New struct `ClientListOptions`
- New struct `ClientListPager`
- New struct `ClientListResponse`
- New struct `ClientListResult`
- New struct `ErrorDetail`
- New struct `ErrorResponseAutoGenerated`
- New struct `Peers`
- New struct `TenantsClientListOptions`
- New struct `TenantsClientListPager`
- New struct `TenantsClientListResponse`
- New struct `TenantsClientListResult`
- New field `Error` in struct `CloudError`


## 0.1.1 (2021-12-13)

### Other Changes

- Fix the go minimum version to `1.16`

## 0.1.0 (2021-11-16)

- Initial preview release.
