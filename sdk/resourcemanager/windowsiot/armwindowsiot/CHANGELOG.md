# Release History

## 0.4.0 (2022-04-18)
### Breaking Changes

- Function `*OperationsClient.List` has been removed
- Function `*ServicesClient.List` has been removed
- Function `*ServicesClient.ListByResourceGroup` has been removed

### Features Added

- New function `*ServicesClient.NewListByResourceGroupPager(string, *ServicesClientListByResourceGroupOptions) *runtime.Pager[ServicesClientListByResourceGroupResponse]`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New function `*ServicesClient.NewListPager(*ServicesClientListOptions) *runtime.Pager[ServicesClientListResponse]`


## 0.3.0 (2022-04-13)
### Breaking Changes

- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `*ServicesClient.ListByResourceGroup` return value(s) have been changed from `(*ServicesClientListByResourceGroupPager)` to `(*runtime.Pager[ServicesClientListByResourceGroupResponse])`
- Function `*ServicesClient.List` return value(s) have been changed from `(*ServicesClientListPager)` to `(*runtime.Pager[ServicesClientListResponse])`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `NewServicesClient` return value(s) have been changed from `(*ServicesClient)` to `(*ServicesClient, error)`
- Function `*OperationsClientListPager.NextPage` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Function `*ServicesClientListByResourceGroupPager.NextPage` has been removed
- Function `ServiceNameUnavailabilityReason.ToPtr` has been removed
- Function `*OperationsClientListPager.Err` has been removed
- Function `*ServicesClientListPager.PageResponse` has been removed
- Function `*ServicesClientListByResourceGroupPager.Err` has been removed
- Function `*ServicesClientListPager.Err` has been removed
- Function `*ServicesClientListPager.NextPage` has been removed
- Function `*ServicesClientListByResourceGroupPager.PageResponse` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Struct `ServicesClientCheckDeviceServiceNameAvailabilityResult` has been removed
- Struct `ServicesClientCreateOrUpdateResult` has been removed
- Struct `ServicesClientDeleteResult` has been removed
- Struct `ServicesClientGetResult` has been removed
- Struct `ServicesClientListByResourceGroupPager` has been removed
- Struct `ServicesClientListByResourceGroupResult` has been removed
- Struct `ServicesClientListPager` has been removed
- Struct `ServicesClientListResult` has been removed
- Struct `ServicesClientUpdateResult` has been removed
- Field `ServicesClientListResult` of struct `ServicesClientListResponse` has been removed
- Field `RawResponse` of struct `ServicesClientListResponse` has been removed
- Field `ServicesClientCheckDeviceServiceNameAvailabilityResult` of struct `ServicesClientCheckDeviceServiceNameAvailabilityResponse` has been removed
- Field `RawResponse` of struct `ServicesClientCheckDeviceServiceNameAvailabilityResponse` has been removed
- Field `ServicesClientCreateOrUpdateResult` of struct `ServicesClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `ServicesClientCreateOrUpdateResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `ServicesClientUpdateResult` of struct `ServicesClientUpdateResponse` has been removed
- Field `RawResponse` of struct `ServicesClientUpdateResponse` has been removed
- Field `ServicesClientGetResult` of struct `ServicesClientGetResponse` has been removed
- Field `RawResponse` of struct `ServicesClientGetResponse` has been removed
- Field `ServicesClientDeleteResult` of struct `ServicesClientDeleteResponse` has been removed
- Field `RawResponse` of struct `ServicesClientDeleteResponse` has been removed
- Field `ServicesClientListByResourceGroupResult` of struct `ServicesClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `ServicesClientListByResourceGroupResponse` has been removed

### Features Added

- New anonymous field `DeviceService` in struct `ServicesClientDeleteResponse`
- New anonymous field `DeviceServiceDescriptionListResult` in struct `ServicesClientListByResourceGroupResponse`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New anonymous field `DeviceServiceDescriptionListResult` in struct `ServicesClientListResponse`
- New anonymous field `DeviceService` in struct `ServicesClientCreateOrUpdateResponse`
- New anonymous field `DeviceService` in struct `ServicesClientGetResponse`
- New anonymous field `DeviceServiceNameAvailabilityInfo` in struct `ServicesClientCheckDeviceServiceNameAvailabilityResponse`
- New anonymous field `DeviceService` in struct `ServicesClientUpdateResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

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
