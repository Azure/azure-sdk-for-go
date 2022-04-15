# Release History

## 0.4.0 (2022-04-15)
### Breaking Changes

- Function `*CustomLocationsClient.ListEnabledResourceTypes` has been removed
- Function `*CustomLocationsClient.ListBySubscription` has been removed
- Function `*CustomLocationsClient.ListOperations` has been removed
- Function `*CustomLocationsClient.ListByResourceGroup` has been removed

### Features Added

- New function `*CustomLocationsClient.NewListBySubscriptionPager(*CustomLocationsClientListBySubscriptionOptions) *runtime.Pager[CustomLocationsClientListBySubscriptionResponse]`
- New function `*CustomLocationsClient.NewListEnabledResourceTypesPager(string, string, *CustomLocationsClientListEnabledResourceTypesOptions) *runtime.Pager[CustomLocationsClientListEnabledResourceTypesResponse]`
- New function `*CustomLocationsClient.NewListByResourceGroupPager(string, *CustomLocationsClientListByResourceGroupOptions) *runtime.Pager[CustomLocationsClientListByResourceGroupResponse]`
- New function `*CustomLocationsClient.NewListOperationsPager(*CustomLocationsClientListOperationsOptions) *runtime.Pager[CustomLocationsClientListOperationsResponse]`


## 0.3.0 (2022-04-11)
### Breaking Changes

- Function `*CustomLocationsClient.ListEnabledResourceTypes` return value(s) have been changed from `(*CustomLocationsClientListEnabledResourceTypesPager)` to `(*runtime.Pager[CustomLocationsClientListEnabledResourceTypesResponse])`
- Function `*CustomLocationsClient.BeginCreateOrUpdate` return value(s) have been changed from `(CustomLocationsClientCreateOrUpdatePollerResponse, error)` to `(*armruntime.Poller[CustomLocationsClientCreateOrUpdateResponse], error)`
- Function `*CustomLocationsClient.ListBySubscription` return value(s) have been changed from `(*CustomLocationsClientListBySubscriptionPager)` to `(*runtime.Pager[CustomLocationsClientListBySubscriptionResponse])`
- Function `*CustomLocationsClient.ListOperations` return value(s) have been changed from `(*CustomLocationsClientListOperationsPager)` to `(*runtime.Pager[CustomLocationsClientListOperationsResponse])`
- Function `NewCustomLocationsClient` return value(s) have been changed from `(*CustomLocationsClient)` to `(*CustomLocationsClient, error)`
- Function `*CustomLocationsClient.BeginDelete` return value(s) have been changed from `(CustomLocationsClientDeletePollerResponse, error)` to `(*armruntime.Poller[CustomLocationsClientDeleteResponse], error)`
- Function `*CustomLocationsClient.ListByResourceGroup` return value(s) have been changed from `(*CustomLocationsClientListByResourceGroupPager)` to `(*runtime.Pager[CustomLocationsClientListByResourceGroupResponse])`
- Type of `ErrorAdditionalInfo.Info` has been changed from `map[string]interface{}` to `interface{}`
- Function `*CustomLocationsClientDeletePoller.Done` has been removed
- Function `*CustomLocationsClientListBySubscriptionPager.PageResponse` has been removed
- Function `*CustomLocationsClientListEnabledResourceTypesPager.PageResponse` has been removed
- Function `*CustomLocationsClientListBySubscriptionPager.Err` has been removed
- Function `*CustomLocationsClientListOperationsPager.Err` has been removed
- Function `*CustomLocationsClientListByResourceGroupPager.NextPage` has been removed
- Function `*CustomLocationsClientCreateOrUpdatePoller.Poll` has been removed
- Function `*CustomLocationsClientCreateOrUpdatePollerResponse.Resume` has been removed
- Function `*CustomLocationsClientDeletePollerResponse.Resume` has been removed
- Function `*CustomLocationsClientCreateOrUpdatePoller.FinalResponse` has been removed
- Function `CustomLocationsClientCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `*CustomLocationsClientListEnabledResourceTypesPager.NextPage` has been removed
- Function `*CustomLocationsClientCreateOrUpdatePoller.Done` has been removed
- Function `*CustomLocationsClientCreateOrUpdatePoller.ResumeToken` has been removed
- Function `CustomLocationsClientDeletePollerResponse.PollUntilDone` has been removed
- Function `*CustomLocationsClientDeletePoller.FinalResponse` has been removed
- Function `*CustomLocationsClientDeletePoller.ResumeToken` has been removed
- Function `*CustomLocationsClientListByResourceGroupPager.PageResponse` has been removed
- Function `*CustomLocationsClientListOperationsPager.NextPage` has been removed
- Function `*CustomLocationsClientListEnabledResourceTypesPager.Err` has been removed
- Function `*CustomLocationsClientListByResourceGroupPager.Err` has been removed
- Function `*CustomLocationsClientListOperationsPager.PageResponse` has been removed
- Function `*CustomLocationsClientDeletePoller.Poll` has been removed
- Function `HostType.ToPtr` has been removed
- Function `CreatedByType.ToPtr` has been removed
- Function `ResourceIdentityType.ToPtr` has been removed
- Function `*CustomLocationsClientListBySubscriptionPager.NextPage` has been removed
- Struct `CustomLocationsClientCreateOrUpdatePoller` has been removed
- Struct `CustomLocationsClientCreateOrUpdatePollerResponse` has been removed
- Struct `CustomLocationsClientCreateOrUpdateResult` has been removed
- Struct `CustomLocationsClientDeletePoller` has been removed
- Struct `CustomLocationsClientDeletePollerResponse` has been removed
- Struct `CustomLocationsClientGetResult` has been removed
- Struct `CustomLocationsClientListByResourceGroupPager` has been removed
- Struct `CustomLocationsClientListByResourceGroupResult` has been removed
- Struct `CustomLocationsClientListBySubscriptionPager` has been removed
- Struct `CustomLocationsClientListBySubscriptionResult` has been removed
- Struct `CustomLocationsClientListEnabledResourceTypesPager` has been removed
- Struct `CustomLocationsClientListEnabledResourceTypesResult` has been removed
- Struct `CustomLocationsClientListOperationsPager` has been removed
- Struct `CustomLocationsClientListOperationsResult` has been removed
- Struct `CustomLocationsClientUpdateResult` has been removed
- Field `CustomLocationsClientListBySubscriptionResult` of struct `CustomLocationsClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `CustomLocationsClientListBySubscriptionResponse` has been removed
- Field `CustomLocationsClientCreateOrUpdateResult` of struct `CustomLocationsClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `CustomLocationsClientCreateOrUpdateResponse` has been removed
- Field `CustomLocationsClientListEnabledResourceTypesResult` of struct `CustomLocationsClientListEnabledResourceTypesResponse` has been removed
- Field `RawResponse` of struct `CustomLocationsClientListEnabledResourceTypesResponse` has been removed
- Field `CustomLocationsClientListByResourceGroupResult` of struct `CustomLocationsClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `CustomLocationsClientListByResourceGroupResponse` has been removed
- Field `CustomLocationsClientGetResult` of struct `CustomLocationsClientGetResponse` has been removed
- Field `RawResponse` of struct `CustomLocationsClientGetResponse` has been removed
- Field `RawResponse` of struct `CustomLocationsClientDeleteResponse` has been removed
- Field `CustomLocationsClientListOperationsResult` of struct `CustomLocationsClientListOperationsResponse` has been removed
- Field `RawResponse` of struct `CustomLocationsClientListOperationsResponse` has been removed
- Field `CustomLocationsClientUpdateResult` of struct `CustomLocationsClientUpdateResponse` has been removed
- Field `RawResponse` of struct `CustomLocationsClientUpdateResponse` has been removed

### Features Added

- New anonymous field `CustomLocation` in struct `CustomLocationsClientUpdateResponse`
- New anonymous field `CustomLocationListResult` in struct `CustomLocationsClientListBySubscriptionResponse`
- New anonymous field `CustomLocation` in struct `CustomLocationsClientCreateOrUpdateResponse`
- New anonymous field `CustomLocation` in struct `CustomLocationsClientGetResponse`
- New anonymous field `EnabledResourceTypesListResult` in struct `CustomLocationsClientListEnabledResourceTypesResponse`
- New anonymous field `CustomLocationOperationsList` in struct `CustomLocationsClientListOperationsResponse`
- New anonymous field `CustomLocationListResult` in struct `CustomLocationsClientListByResourceGroupResponse`
- New field `ResumeToken` in struct `CustomLocationsClientBeginDeleteOptions`
- New field `ResumeToken` in struct `CustomLocationsClientBeginCreateOrUpdateOptions`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*CustomLocationsClient.Update` parameter(s) have been changed from `(context.Context, string, string, PatchableCustomLocations, *CustomLocationsUpdateOptions)` to `(context.Context, string, string, PatchableCustomLocations, *CustomLocationsClientUpdateOptions)`
- Function `*CustomLocationsClient.Update` return value(s) have been changed from `(CustomLocationsUpdateResponse, error)` to `(CustomLocationsClientUpdateResponse, error)`
- Function `*CustomLocationsClient.ListByResourceGroup` parameter(s) have been changed from `(string, *CustomLocationsListByResourceGroupOptions)` to `(string, *CustomLocationsClientListByResourceGroupOptions)`
- Function `*CustomLocationsClient.ListByResourceGroup` return value(s) have been changed from `(*CustomLocationsListByResourceGroupPager)` to `(*CustomLocationsClientListByResourceGroupPager)`
- Function `*CustomLocationsClient.ListEnabledResourceTypes` parameter(s) have been changed from `(string, string, *CustomLocationsListEnabledResourceTypesOptions)` to `(string, string, *CustomLocationsClientListEnabledResourceTypesOptions)`
- Function `*CustomLocationsClient.ListEnabledResourceTypes` return value(s) have been changed from `(*CustomLocationsListEnabledResourceTypesPager)` to `(*CustomLocationsClientListEnabledResourceTypesPager)`
- Function `*CustomLocationsClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, CustomLocation, *CustomLocationsBeginCreateOrUpdateOptions)` to `(context.Context, string, string, CustomLocation, *CustomLocationsClientBeginCreateOrUpdateOptions)`
- Function `*CustomLocationsClient.BeginCreateOrUpdate` return value(s) have been changed from `(CustomLocationsCreateOrUpdatePollerResponse, error)` to `(CustomLocationsClientCreateOrUpdatePollerResponse, error)`
- Function `*CustomLocationsClient.ListOperations` parameter(s) have been changed from `(*CustomLocationsListOperationsOptions)` to `(*CustomLocationsClientListOperationsOptions)`
- Function `*CustomLocationsClient.ListOperations` return value(s) have been changed from `(*CustomLocationsListOperationsPager)` to `(*CustomLocationsClientListOperationsPager)`
- Function `*CustomLocationsClient.ListBySubscription` parameter(s) have been changed from `(*CustomLocationsListBySubscriptionOptions)` to `(*CustomLocationsClientListBySubscriptionOptions)`
- Function `*CustomLocationsClient.ListBySubscription` return value(s) have been changed from `(*CustomLocationsListBySubscriptionPager)` to `(*CustomLocationsClientListBySubscriptionPager)`
- Function `*CustomLocationsClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *CustomLocationsBeginDeleteOptions)` to `(context.Context, string, string, *CustomLocationsClientBeginDeleteOptions)`
- Function `*CustomLocationsClient.BeginDelete` return value(s) have been changed from `(CustomLocationsDeletePollerResponse, error)` to `(CustomLocationsClientDeletePollerResponse, error)`
- Function `*CustomLocationsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *CustomLocationsGetOptions)` to `(context.Context, string, string, *CustomLocationsClientGetOptions)`
- Function `*CustomLocationsClient.Get` return value(s) have been changed from `(CustomLocationsGetResponse, error)` to `(CustomLocationsClientGetResponse, error)`
- Function `*CustomLocationsListByResourceGroupPager.NextPage` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `*CustomLocationsDeletePoller.ResumeToken` has been removed
- Function `*CustomLocationsListEnabledResourceTypesPager.Err` has been removed
- Function `*CustomLocationsCreateOrUpdatePoller.ResumeToken` has been removed
- Function `CustomLocationsDeletePollerResponse.PollUntilDone` has been removed
- Function `*CustomLocationsCreateOrUpdatePollerResponse.Resume` has been removed
- Function `*CustomLocationsListByResourceGroupPager.Err` has been removed
- Function `EnabledResourceType.MarshalJSON` has been removed
- Function `*CustomLocationsListBySubscriptionPager.PageResponse` has been removed
- Function `*CustomLocationsCreateOrUpdatePoller.FinalResponse` has been removed
- Function `CustomLocationsCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `Resource.MarshalJSON` has been removed
- Function `*CustomLocationsDeletePollerResponse.Resume` has been removed
- Function `*CustomLocationsListOperationsPager.Err` has been removed
- Function `*CustomLocationsListOperationsPager.NextPage` has been removed
- Function `*CustomLocationsListOperationsPager.PageResponse` has been removed
- Function `*CustomLocationsListEnabledResourceTypesPager.PageResponse` has been removed
- Function `*CustomLocationsCreateOrUpdatePoller.Poll` has been removed
- Function `*CustomLocationsDeletePoller.Done` has been removed
- Function `*CustomLocationsListByResourceGroupPager.PageResponse` has been removed
- Function `*CustomLocationsCreateOrUpdatePoller.Done` has been removed
- Function `*CustomLocationsListBySubscriptionPager.Err` has been removed
- Function `*CustomLocationsListEnabledResourceTypesPager.NextPage` has been removed
- Function `*CustomLocationsDeletePoller.Poll` has been removed
- Function `*CustomLocationsDeletePoller.FinalResponse` has been removed
- Function `*CustomLocationsListBySubscriptionPager.NextPage` has been removed
- Struct `CustomLocationsBeginCreateOrUpdateOptions` has been removed
- Struct `CustomLocationsBeginDeleteOptions` has been removed
- Struct `CustomLocationsCreateOrUpdatePoller` has been removed
- Struct `CustomLocationsCreateOrUpdatePollerResponse` has been removed
- Struct `CustomLocationsCreateOrUpdateResponse` has been removed
- Struct `CustomLocationsCreateOrUpdateResult` has been removed
- Struct `CustomLocationsDeletePoller` has been removed
- Struct `CustomLocationsDeletePollerResponse` has been removed
- Struct `CustomLocationsDeleteResponse` has been removed
- Struct `CustomLocationsGetOptions` has been removed
- Struct `CustomLocationsGetResponse` has been removed
- Struct `CustomLocationsGetResult` has been removed
- Struct `CustomLocationsListByResourceGroupOptions` has been removed
- Struct `CustomLocationsListByResourceGroupPager` has been removed
- Struct `CustomLocationsListByResourceGroupResponse` has been removed
- Struct `CustomLocationsListByResourceGroupResult` has been removed
- Struct `CustomLocationsListBySubscriptionOptions` has been removed
- Struct `CustomLocationsListBySubscriptionPager` has been removed
- Struct `CustomLocationsListBySubscriptionResponse` has been removed
- Struct `CustomLocationsListBySubscriptionResult` has been removed
- Struct `CustomLocationsListEnabledResourceTypesOptions` has been removed
- Struct `CustomLocationsListEnabledResourceTypesPager` has been removed
- Struct `CustomLocationsListEnabledResourceTypesResponse` has been removed
- Struct `CustomLocationsListEnabledResourceTypesResult` has been removed
- Struct `CustomLocationsListOperationsOptions` has been removed
- Struct `CustomLocationsListOperationsPager` has been removed
- Struct `CustomLocationsListOperationsResponse` has been removed
- Struct `CustomLocationsListOperationsResult` has been removed
- Struct `CustomLocationsUpdateOptions` has been removed
- Struct `CustomLocationsUpdateResponse` has been removed
- Struct `CustomLocationsUpdateResult` has been removed
- Field `Resource` of struct `TrackedResource` has been removed
- Field `TrackedResource` of struct `CustomLocation` has been removed
- Field `Resource` of struct `ProxyResource` has been removed
- Field `ProxyResource` of struct `EnabledResourceType` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed

### Features Added

- New function `*CustomLocationsClientCreateOrUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*CustomLocationsClientListOperationsPager.Err() error`
- New function `*CustomLocationsClientListBySubscriptionPager.Err() error`
- New function `*CustomLocationsClientListOperationsPager.PageResponse() CustomLocationsClientListOperationsResponse`
- New function `*CustomLocationsClientListByResourceGroupPager.PageResponse() CustomLocationsClientListByResourceGroupResponse`
- New function `*CustomLocationsClientListByResourceGroupPager.Err() error`
- New function `CustomLocationsClientCreateOrUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (CustomLocationsClientCreateOrUpdateResponse, error)`
- New function `*CustomLocationsClientListBySubscriptionPager.PageResponse() CustomLocationsClientListBySubscriptionResponse`
- New function `*CustomLocationsClientCreateOrUpdatePollerResponse.Resume(context.Context, *CustomLocationsClient, string) error`
- New function `CustomLocationsClientDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (CustomLocationsClientDeleteResponse, error)`
- New function `*CustomLocationsClientDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `*CustomLocationsClientDeletePoller.ResumeToken() (string, error)`
- New function `*CustomLocationsClientDeletePollerResponse.Resume(context.Context, *CustomLocationsClient, string) error`
- New function `*CustomLocationsClientListOperationsPager.NextPage(context.Context) bool`
- New function `*CustomLocationsClientListEnabledResourceTypesPager.Err() error`
- New function `*CustomLocationsClientCreateOrUpdatePoller.ResumeToken() (string, error)`
- New function `*CustomLocationsClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*CustomLocationsClientDeletePoller.Done() bool`
- New function `*CustomLocationsClientListBySubscriptionPager.NextPage(context.Context) bool`
- New function `*CustomLocationsClientDeletePoller.FinalResponse(context.Context) (CustomLocationsClientDeleteResponse, error)`
- New function `*CustomLocationsClientCreateOrUpdatePoller.FinalResponse(context.Context) (CustomLocationsClientCreateOrUpdateResponse, error)`
- New function `*CustomLocationsClientListEnabledResourceTypesPager.PageResponse() CustomLocationsClientListEnabledResourceTypesResponse`
- New function `*CustomLocationsClientCreateOrUpdatePoller.Done() bool`
- New function `*CustomLocationsClientListEnabledResourceTypesPager.NextPage(context.Context) bool`
- New struct `CustomLocationsClientBeginCreateOrUpdateOptions`
- New struct `CustomLocationsClientBeginDeleteOptions`
- New struct `CustomLocationsClientCreateOrUpdatePoller`
- New struct `CustomLocationsClientCreateOrUpdatePollerResponse`
- New struct `CustomLocationsClientCreateOrUpdateResponse`
- New struct `CustomLocationsClientCreateOrUpdateResult`
- New struct `CustomLocationsClientDeletePoller`
- New struct `CustomLocationsClientDeletePollerResponse`
- New struct `CustomLocationsClientDeleteResponse`
- New struct `CustomLocationsClientGetOptions`
- New struct `CustomLocationsClientGetResponse`
- New struct `CustomLocationsClientGetResult`
- New struct `CustomLocationsClientListByResourceGroupOptions`
- New struct `CustomLocationsClientListByResourceGroupPager`
- New struct `CustomLocationsClientListByResourceGroupResponse`
- New struct `CustomLocationsClientListByResourceGroupResult`
- New struct `CustomLocationsClientListBySubscriptionOptions`
- New struct `CustomLocationsClientListBySubscriptionPager`
- New struct `CustomLocationsClientListBySubscriptionResponse`
- New struct `CustomLocationsClientListBySubscriptionResult`
- New struct `CustomLocationsClientListEnabledResourceTypesOptions`
- New struct `CustomLocationsClientListEnabledResourceTypesPager`
- New struct `CustomLocationsClientListEnabledResourceTypesResponse`
- New struct `CustomLocationsClientListEnabledResourceTypesResult`
- New struct `CustomLocationsClientListOperationsOptions`
- New struct `CustomLocationsClientListOperationsPager`
- New struct `CustomLocationsClientListOperationsResponse`
- New struct `CustomLocationsClientListOperationsResult`
- New struct `CustomLocationsClientUpdateOptions`
- New struct `CustomLocationsClientUpdateResponse`
- New struct `CustomLocationsClientUpdateResult`
- New field `Error` in struct `ErrorResponse`
- New field `Type` in struct `CustomLocation`
- New field `Location` in struct `CustomLocation`
- New field `Tags` in struct `CustomLocation`
- New field `ID` in struct `CustomLocation`
- New field `Name` in struct `CustomLocation`
- New field `ID` in struct `ProxyResource`
- New field `Name` in struct `ProxyResource`
- New field `Type` in struct `ProxyResource`
- New field `Name` in struct `EnabledResourceType`
- New field `Type` in struct `EnabledResourceType`
- New field `ID` in struct `EnabledResourceType`
- New field `ID` in struct `TrackedResource`
- New field `Name` in struct `TrackedResource`
- New field `Type` in struct `TrackedResource`


## 0.1.0 (2021-12-07)

- Init release.
