# Release History

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*SubscriptionFeatureRegistrationsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, *SubscriptionFeatureRegistrationsCreateOrUpdateOptions)` to `(context.Context, string, string, *SubscriptionFeatureRegistrationsClientCreateOrUpdateOptions)`
- Function `*SubscriptionFeatureRegistrationsClient.CreateOrUpdate` return value(s) have been changed from `(SubscriptionFeatureRegistrationsCreateOrUpdateResponse, error)` to `(SubscriptionFeatureRegistrationsClientCreateOrUpdateResponse, error)`
- Function `*SubscriptionFeatureRegistrationsClient.ListBySubscription` parameter(s) have been changed from `(string, *SubscriptionFeatureRegistrationsListBySubscriptionOptions)` to `(string, *SubscriptionFeatureRegistrationsClientListBySubscriptionOptions)`
- Function `*SubscriptionFeatureRegistrationsClient.ListBySubscription` return value(s) have been changed from `(*SubscriptionFeatureRegistrationsListBySubscriptionPager)` to `(*SubscriptionFeatureRegistrationsClientListBySubscriptionPager)`
- Function `*SubscriptionFeatureRegistrationsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *SubscriptionFeatureRegistrationsGetOptions)` to `(context.Context, string, string, *SubscriptionFeatureRegistrationsClientGetOptions)`
- Function `*SubscriptionFeatureRegistrationsClient.Get` return value(s) have been changed from `(SubscriptionFeatureRegistrationsGetResponse, error)` to `(SubscriptionFeatureRegistrationsClientGetResponse, error)`
- Function `*SubscriptionFeatureRegistrationsClient.ListAllBySubscription` parameter(s) have been changed from `(*SubscriptionFeatureRegistrationsListAllBySubscriptionOptions)` to `(*SubscriptionFeatureRegistrationsClientListAllBySubscriptionOptions)`
- Function `*SubscriptionFeatureRegistrationsClient.ListAllBySubscription` return value(s) have been changed from `(*SubscriptionFeatureRegistrationsListAllBySubscriptionPager)` to `(*SubscriptionFeatureRegistrationsClientListAllBySubscriptionPager)`
- Function `*SubscriptionFeatureRegistrationsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *SubscriptionFeatureRegistrationsDeleteOptions)` to `(context.Context, string, string, *SubscriptionFeatureRegistrationsClientDeleteOptions)`
- Function `*SubscriptionFeatureRegistrationsClient.Delete` return value(s) have been changed from `(SubscriptionFeatureRegistrationsDeleteResponse, error)` to `(SubscriptionFeatureRegistrationsClientDeleteResponse, error)`
- Function `*FeaturesClient.ListAll` has been removed
- Function `*SubscriptionFeatureRegistrationsListBySubscriptionPager.NextPage` has been removed
- Function `*FeaturesClient.List` has been removed
- Function `*FeaturesListPager.Err` has been removed
- Function `*FeaturesClient.Register` has been removed
- Function `*FeaturesClient.Get` has been removed
- Function `*FeaturesListPager.PageResponse` has been removed
- Function `*FeaturesClient.Unregister` has been removed
- Function `*FeaturesListAllPager.PageResponse` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `*FeaturesListAllPager.Err` has been removed
- Function `*FeaturesListPager.NextPage` has been removed
- Function `*SubscriptionFeatureRegistrationsListAllBySubscriptionPager.Err` has been removed
- Function `*SubscriptionFeatureRegistrationsListBySubscriptionPager.PageResponse` has been removed
- Function `NewFeaturesClient` has been removed
- Function `*FeaturesListAllPager.NextPage` has been removed
- Function `*SubscriptionFeatureRegistrationsListBySubscriptionPager.Err` has been removed
- Function `*SubscriptionFeatureRegistrationsListAllBySubscriptionPager.PageResponse` has been removed
- Function `*SubscriptionFeatureRegistrationsListAllBySubscriptionPager.NextPage` has been removed
- Struct `FeaturesClient` has been removed
- Struct `FeaturesGetOptions` has been removed
- Struct `FeaturesGetResponse` has been removed
- Struct `FeaturesGetResult` has been removed
- Struct `FeaturesListAllOptions` has been removed
- Struct `FeaturesListAllPager` has been removed
- Struct `FeaturesListAllResponse` has been removed
- Struct `FeaturesListAllResult` has been removed
- Struct `FeaturesListOptions` has been removed
- Struct `FeaturesListPager` has been removed
- Struct `FeaturesListResponse` has been removed
- Struct `FeaturesListResult` has been removed
- Struct `FeaturesRegisterOptions` has been removed
- Struct `FeaturesRegisterResponse` has been removed
- Struct `FeaturesRegisterResult` has been removed
- Struct `FeaturesUnregisterOptions` has been removed
- Struct `FeaturesUnregisterResponse` has been removed
- Struct `FeaturesUnregisterResult` has been removed
- Struct `SubscriptionFeatureRegistrationsCreateOrUpdateOptions` has been removed
- Struct `SubscriptionFeatureRegistrationsCreateOrUpdateResponse` has been removed
- Struct `SubscriptionFeatureRegistrationsCreateOrUpdateResult` has been removed
- Struct `SubscriptionFeatureRegistrationsDeleteOptions` has been removed
- Struct `SubscriptionFeatureRegistrationsDeleteResponse` has been removed
- Struct `SubscriptionFeatureRegistrationsGetOptions` has been removed
- Struct `SubscriptionFeatureRegistrationsGetResponse` has been removed
- Struct `SubscriptionFeatureRegistrationsGetResult` has been removed
- Struct `SubscriptionFeatureRegistrationsListAllBySubscriptionOptions` has been removed
- Struct `SubscriptionFeatureRegistrationsListAllBySubscriptionPager` has been removed
- Struct `SubscriptionFeatureRegistrationsListAllBySubscriptionResponse` has been removed
- Struct `SubscriptionFeatureRegistrationsListAllBySubscriptionResult` has been removed
- Struct `SubscriptionFeatureRegistrationsListBySubscriptionOptions` has been removed
- Struct `SubscriptionFeatureRegistrationsListBySubscriptionPager` has been removed
- Struct `SubscriptionFeatureRegistrationsListBySubscriptionResponse` has been removed
- Struct `SubscriptionFeatureRegistrationsListBySubscriptionResult` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed
- Field `ProxyResource` of struct `SubscriptionFeatureRegistration` has been removed

### Features Added

- New function `*ClientListPager.PageResponse() ClientListResponse`
- New function `*Client.List(string, *ClientListOptions) *ClientListPager`
- New function `*ClientListAllPager.NextPage(context.Context) bool`
- New function `*SubscriptionFeatureRegistrationsClientListAllBySubscriptionPager.Err() error`
- New function `*SubscriptionFeatureRegistrationsClientListAllBySubscriptionPager.NextPage(context.Context) bool`
- New function `*SubscriptionFeatureRegistrationsClientListBySubscriptionPager.NextPage(context.Context) bool`
- New function `*SubscriptionFeatureRegistrationsClientListBySubscriptionPager.PageResponse() SubscriptionFeatureRegistrationsClientListBySubscriptionResponse`
- New function `*ClientListPager.Err() error`
- New function `*Client.Get(context.Context, string, string, *ClientGetOptions) (ClientGetResponse, error)`
- New function `*ClientListAllPager.PageResponse() ClientListAllResponse`
- New function `*SubscriptionFeatureRegistrationsClientListBySubscriptionPager.Err() error`
- New function `*Client.Unregister(context.Context, string, string, *ClientUnregisterOptions) (ClientUnregisterResponse, error)`
- New function `*Client.ListAll(*ClientListAllOptions) *ClientListAllPager`
- New function `*Client.Register(context.Context, string, string, *ClientRegisterOptions) (ClientRegisterResponse, error)`
- New function `NewClient(string, azcore.TokenCredential, *arm.ClientOptions) *Client`
- New function `*SubscriptionFeatureRegistrationsClientListAllBySubscriptionPager.PageResponse() SubscriptionFeatureRegistrationsClientListAllBySubscriptionResponse`
- New function `*ClientListAllPager.Err() error`
- New function `*ClientListPager.NextPage(context.Context) bool`
- New struct `Client`
- New struct `ClientGetOptions`
- New struct `ClientGetResponse`
- New struct `ClientGetResult`
- New struct `ClientListAllOptions`
- New struct `ClientListAllPager`
- New struct `ClientListAllResponse`
- New struct `ClientListAllResult`
- New struct `ClientListOptions`
- New struct `ClientListPager`
- New struct `ClientListResponse`
- New struct `ClientListResult`
- New struct `ClientRegisterOptions`
- New struct `ClientRegisterResponse`
- New struct `ClientRegisterResult`
- New struct `ClientUnregisterOptions`
- New struct `ClientUnregisterResponse`
- New struct `ClientUnregisterResult`
- New struct `SubscriptionFeatureRegistrationsClientCreateOrUpdateOptions`
- New struct `SubscriptionFeatureRegistrationsClientCreateOrUpdateResponse`
- New struct `SubscriptionFeatureRegistrationsClientCreateOrUpdateResult`
- New struct `SubscriptionFeatureRegistrationsClientDeleteOptions`
- New struct `SubscriptionFeatureRegistrationsClientDeleteResponse`
- New struct `SubscriptionFeatureRegistrationsClientGetOptions`
- New struct `SubscriptionFeatureRegistrationsClientGetResponse`
- New struct `SubscriptionFeatureRegistrationsClientGetResult`
- New struct `SubscriptionFeatureRegistrationsClientListAllBySubscriptionOptions`
- New struct `SubscriptionFeatureRegistrationsClientListAllBySubscriptionPager`
- New struct `SubscriptionFeatureRegistrationsClientListAllBySubscriptionResponse`
- New struct `SubscriptionFeatureRegistrationsClientListAllBySubscriptionResult`
- New struct `SubscriptionFeatureRegistrationsClientListBySubscriptionOptions`
- New struct `SubscriptionFeatureRegistrationsClientListBySubscriptionPager`
- New struct `SubscriptionFeatureRegistrationsClientListBySubscriptionResponse`
- New struct `SubscriptionFeatureRegistrationsClientListBySubscriptionResult`
- New field `Error` in struct `ErrorResponse`
- New field `ID` in struct `SubscriptionFeatureRegistration`
- New field `Name` in struct `SubscriptionFeatureRegistration`
- New field `Type` in struct `SubscriptionFeatureRegistration`


## 0.1.1 (2021-12-13)

### Other Changes

- Fix the go minimum version to `1.16`

## 0.1.0 (2021-11-16)

- Initial preview release.
