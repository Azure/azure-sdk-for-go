# Release History

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*AzureBareMetalInstancesClient.Get` parameter(s) have been changed from `(context.Context, string, string, *AzureBareMetalInstancesGetOptions)` to `(context.Context, string, string, *AzureBareMetalInstancesClientGetOptions)`
- Function `*AzureBareMetalInstancesClient.Get` return value(s) have been changed from `(AzureBareMetalInstancesGetResponse, error)` to `(AzureBareMetalInstancesClientGetResponse, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(context.Context, *OperationsListOptions)` to `(context.Context, *OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(OperationsListResponse, error)` to `(OperationsClientListResponse, error)`
- Function `*AzureBareMetalInstancesClient.ListByResourceGroup` parameter(s) have been changed from `(string, *AzureBareMetalInstancesListByResourceGroupOptions)` to `(string, *AzureBareMetalInstancesClientListByResourceGroupOptions)`
- Function `*AzureBareMetalInstancesClient.ListByResourceGroup` return value(s) have been changed from `(*AzureBareMetalInstancesListByResourceGroupPager)` to `(*AzureBareMetalInstancesClientListByResourceGroupPager)`
- Function `*AzureBareMetalInstancesClient.ListBySubscription` parameter(s) have been changed from `(*AzureBareMetalInstancesListBySubscriptionOptions)` to `(*AzureBareMetalInstancesClientListBySubscriptionOptions)`
- Function `*AzureBareMetalInstancesClient.ListBySubscription` return value(s) have been changed from `(*AzureBareMetalInstancesListBySubscriptionPager)` to `(*AzureBareMetalInstancesClientListBySubscriptionPager)`
- Function `*AzureBareMetalInstancesClient.Update` parameter(s) have been changed from `(context.Context, string, string, Tags, *AzureBareMetalInstancesUpdateOptions)` to `(context.Context, string, string, Tags, *AzureBareMetalInstancesClientUpdateOptions)`
- Function `*AzureBareMetalInstancesClient.Update` return value(s) have been changed from `(AzureBareMetalInstancesUpdateResponse, error)` to `(AzureBareMetalInstancesClientUpdateResponse, error)`
- Function `*AzureBareMetalInstancesListByResourceGroupPager.Err` has been removed
- Function `*AzureBareMetalInstancesListByResourceGroupPager.NextPage` has been removed
- Function `*AzureBareMetalInstancesListBySubscriptionPager.NextPage` has been removed
- Function `*AzureBareMetalInstancesListByResourceGroupPager.PageResponse` has been removed
- Function `*AzureBareMetalInstancesListBySubscriptionPager.Err` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `*AzureBareMetalInstancesListBySubscriptionPager.PageResponse` has been removed
- Function `Resource.MarshalJSON` has been removed
- Struct `AzureBareMetalInstancesGetOptions` has been removed
- Struct `AzureBareMetalInstancesGetResponse` has been removed
- Struct `AzureBareMetalInstancesGetResult` has been removed
- Struct `AzureBareMetalInstancesListByResourceGroupOptions` has been removed
- Struct `AzureBareMetalInstancesListByResourceGroupPager` has been removed
- Struct `AzureBareMetalInstancesListByResourceGroupResponse` has been removed
- Struct `AzureBareMetalInstancesListByResourceGroupResult` has been removed
- Struct `AzureBareMetalInstancesListBySubscriptionOptions` has been removed
- Struct `AzureBareMetalInstancesListBySubscriptionPager` has been removed
- Struct `AzureBareMetalInstancesListBySubscriptionResponse` has been removed
- Struct `AzureBareMetalInstancesListBySubscriptionResult` has been removed
- Struct `AzureBareMetalInstancesUpdateOptions` has been removed
- Struct `AzureBareMetalInstancesUpdateResponse` has been removed
- Struct `AzureBareMetalInstancesUpdateResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `Resource` of struct `TrackedResource` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed
- Field `TrackedResource` of struct `AzureBareMetalInstance` has been removed

### Features Added

- New function `*AzureBareMetalInstancesClientListBySubscriptionPager.PageResponse() AzureBareMetalInstancesClientListBySubscriptionResponse`
- New function `*AzureBareMetalInstancesClientListByResourceGroupPager.Err() error`
- New function `*AzureBareMetalInstancesClientListBySubscriptionPager.Err() error`
- New function `*AzureBareMetalInstancesClientListByResourceGroupPager.PageResponse() AzureBareMetalInstancesClientListByResourceGroupResponse`
- New function `*AzureBareMetalInstancesClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*AzureBareMetalInstancesClientListBySubscriptionPager.NextPage(context.Context) bool`
- New struct `AzureBareMetalInstancesClientGetOptions`
- New struct `AzureBareMetalInstancesClientGetResponse`
- New struct `AzureBareMetalInstancesClientGetResult`
- New struct `AzureBareMetalInstancesClientListByResourceGroupOptions`
- New struct `AzureBareMetalInstancesClientListByResourceGroupPager`
- New struct `AzureBareMetalInstancesClientListByResourceGroupResponse`
- New struct `AzureBareMetalInstancesClientListByResourceGroupResult`
- New struct `AzureBareMetalInstancesClientListBySubscriptionOptions`
- New struct `AzureBareMetalInstancesClientListBySubscriptionPager`
- New struct `AzureBareMetalInstancesClientListBySubscriptionResponse`
- New struct `AzureBareMetalInstancesClientListBySubscriptionResult`
- New struct `AzureBareMetalInstancesClientUpdateOptions`
- New struct `AzureBareMetalInstancesClientUpdateResponse`
- New struct `AzureBareMetalInstancesClientUpdateResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New field `ID` in struct `TrackedResource`
- New field `Name` in struct `TrackedResource`
- New field `Type` in struct `TrackedResource`
- New field `ID` in struct `AzureBareMetalInstance`
- New field `Name` in struct `AzureBareMetalInstance`
- New field `Type` in struct `AzureBareMetalInstance`
- New field `Location` in struct `AzureBareMetalInstance`
- New field `Tags` in struct `AzureBareMetalInstance`
- New field `Error` in struct `ErrorResponse`


## 0.1.0 (2021-12-01)

- Initial preview release.
