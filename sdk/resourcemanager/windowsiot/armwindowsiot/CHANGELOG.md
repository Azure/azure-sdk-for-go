# Release History

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*ServicesClient.List` parameter(s) have been changed from `(*ServicesListOptions)` to `(*ServicesClientListOptions)`
- Function `*ServicesClient.List` return value(s) have been changed from `(*ServicesListPager)` to `(*ServicesClientListPager)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*ServicesClient.CheckDeviceServiceNameAvailability` parameter(s) have been changed from `(context.Context, DeviceServiceCheckNameAvailabilityParameters, *ServicesCheckDeviceServiceNameAvailabilityOptions)` to `(context.Context, DeviceServiceCheckNameAvailabilityParameters, *ServicesClientCheckDeviceServiceNameAvailabilityOptions)`
- Function `*ServicesClient.CheckDeviceServiceNameAvailability` return value(s) have been changed from `(ServicesCheckDeviceServiceNameAvailabilityResponse, error)` to `(ServicesClientCheckDeviceServiceNameAvailabilityResponse, error)`
- Function `*ServicesClient.Get` parameter(s) have been changed from `(context.Context, string, string, *ServicesGetOptions)` to `(context.Context, string, string, *ServicesClientGetOptions)`
- Function `*ServicesClient.Get` return value(s) have been changed from `(ServicesGetResponse, error)` to `(ServicesClientGetResponse, error)`
- Function `*ServicesClient.Update` parameter(s) have been changed from `(context.Context, string, string, DeviceService, *ServicesUpdateOptions)` to `(context.Context, string, string, DeviceService, *ServicesClientUpdateOptions)`
- Function `*ServicesClient.Update` return value(s) have been changed from `(ServicesUpdateResponse, error)` to `(ServicesClientUpdateResponse, error)`
- Function `*ServicesClient.ListByResourceGroup` parameter(s) have been changed from `(string, *ServicesListByResourceGroupOptions)` to `(string, *ServicesClientListByResourceGroupOptions)`
- Function `*ServicesClient.ListByResourceGroup` return value(s) have been changed from `(*ServicesListByResourceGroupPager)` to `(*ServicesClientListByResourceGroupPager)`
- Function `*ServicesClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, DeviceService, *ServicesCreateOrUpdateOptions)` to `(context.Context, string, string, DeviceService, *ServicesClientCreateOrUpdateOptions)`
- Function `*ServicesClient.CreateOrUpdate` return value(s) have been changed from `(ServicesCreateOrUpdateResponse, error)` to `(ServicesClientCreateOrUpdateResponse, error)`
- Function `*ServicesClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *ServicesDeleteOptions)` to `(context.Context, string, string, *ServicesClientDeleteOptions)`
- Function `*ServicesClient.Delete` return value(s) have been changed from `(ServicesDeleteResponse, error)` to `(ServicesClientDeleteResponse, error)`
- Function `*OperationsListPager.PageResponse` has been removed
- Function `*ServicesListByResourceGroupPager.NextPage` has been removed
- Function `Resource.MarshalJSON` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `*ServicesListByResourceGroupPager.Err` has been removed
- Function `*ServicesListPager.NextPage` has been removed
- Function `*ServicesListPager.PageResponse` has been removed
- Function `ErrorDetails.Error` has been removed
- Function `*ServicesListPager.Err` has been removed
- Function `*ServicesListByResourceGroupPager.PageResponse` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Struct `ServicesCheckDeviceServiceNameAvailabilityOptions` has been removed
- Struct `ServicesCheckDeviceServiceNameAvailabilityResponse` has been removed
- Struct `ServicesCheckDeviceServiceNameAvailabilityResult` has been removed
- Struct `ServicesCreateOrUpdateOptions` has been removed
- Struct `ServicesCreateOrUpdateResponse` has been removed
- Struct `ServicesCreateOrUpdateResult` has been removed
- Struct `ServicesDeleteOptions` has been removed
- Struct `ServicesDeleteResponse` has been removed
- Struct `ServicesDeleteResult` has been removed
- Struct `ServicesGetOptions` has been removed
- Struct `ServicesGetResponse` has been removed
- Struct `ServicesGetResult` has been removed
- Struct `ServicesListByResourceGroupOptions` has been removed
- Struct `ServicesListByResourceGroupPager` has been removed
- Struct `ServicesListByResourceGroupResponse` has been removed
- Struct `ServicesListByResourceGroupResult` has been removed
- Struct `ServicesListOptions` has been removed
- Struct `ServicesListPager` has been removed
- Struct `ServicesListResponse` has been removed
- Struct `ServicesListResult` has been removed
- Struct `ServicesUpdateOptions` has been removed
- Struct `ServicesUpdateResponse` has been removed
- Struct `ServicesUpdateResult` has been removed
- Field `Resource` of struct `TrackedResource` has been removed
- Field `TrackedResource` of struct `DeviceService` has been removed
- Field `Resource` of struct `ProxyResource` has been removed
- Field `InnerError` of struct `ErrorDetails` has been removed

### Features Added

- New function `*ServicesClientListPager.PageResponse() ServicesClientListResponse`
- New function `*ServicesClientListPager.Err() error`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*OperationsClientListPager.Err() error`
- New function `*ServicesClientListByResourceGroupPager.Err() error`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*ServicesClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*ServicesClientListPager.NextPage(context.Context) bool`
- New function `*ServicesClientListByResourceGroupPager.PageResponse() ServicesClientListByResourceGroupResponse`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New struct `ServicesClientCheckDeviceServiceNameAvailabilityOptions`
- New struct `ServicesClientCheckDeviceServiceNameAvailabilityResponse`
- New struct `ServicesClientCheckDeviceServiceNameAvailabilityResult`
- New struct `ServicesClientCreateOrUpdateOptions`
- New struct `ServicesClientCreateOrUpdateResponse`
- New struct `ServicesClientCreateOrUpdateResult`
- New struct `ServicesClientDeleteOptions`
- New struct `ServicesClientDeleteResponse`
- New struct `ServicesClientDeleteResult`
- New struct `ServicesClientGetOptions`
- New struct `ServicesClientGetResponse`
- New struct `ServicesClientGetResult`
- New struct `ServicesClientListByResourceGroupOptions`
- New struct `ServicesClientListByResourceGroupPager`
- New struct `ServicesClientListByResourceGroupResponse`
- New struct `ServicesClientListByResourceGroupResult`
- New struct `ServicesClientListOptions`
- New struct `ServicesClientListPager`
- New struct `ServicesClientListResponse`
- New struct `ServicesClientListResult`
- New struct `ServicesClientUpdateOptions`
- New struct `ServicesClientUpdateResponse`
- New struct `ServicesClientUpdateResult`
- New field `Error` in struct `ErrorDetails`
- New field `ID` in struct `ProxyResource`
- New field `Name` in struct `ProxyResource`
- New field `Type` in struct `ProxyResource`
- New field `ID` in struct `TrackedResource`
- New field `Name` in struct `TrackedResource`
- New field `Type` in struct `TrackedResource`
- New field `Tags` in struct `DeviceService`
- New field `ID` in struct `DeviceService`
- New field `Name` in struct `DeviceService`
- New field `Type` in struct `DeviceService`
- New field `Location` in struct `DeviceService`


## 0.1.0 (2021-12-22)

- Init release.
