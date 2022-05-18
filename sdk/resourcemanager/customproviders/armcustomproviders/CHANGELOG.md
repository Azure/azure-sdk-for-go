# Release History

## 0.5.0 (2022-05-18)
### Breaking Changes

- Function `*CustomResourceProviderClient.BeginCreateOrUpdate` return value(s) have been changed from `(*armruntime.Poller[CustomResourceProviderClientCreateOrUpdateResponse], error)` to `(*runtime.Poller[CustomResourceProviderClientCreateOrUpdateResponse], error)`
- Function `*CustomResourceProviderClient.BeginDelete` return value(s) have been changed from `(*armruntime.Poller[CustomResourceProviderClientDeleteResponse], error)` to `(*runtime.Poller[CustomResourceProviderClientDeleteResponse], error)`
- Function `*AssociationsClient.BeginCreateOrUpdate` return value(s) have been changed from `(*armruntime.Poller[AssociationsClientCreateOrUpdateResponse], error)` to `(*runtime.Poller[AssociationsClientCreateOrUpdateResponse], error)`
- Function `*AssociationsClient.BeginDelete` return value(s) have been changed from `(*armruntime.Poller[AssociationsClientDeleteResponse], error)` to `(*runtime.Poller[AssociationsClientDeleteResponse], error)`
- Function `AssociationsList.MarshalJSON` has been removed
- Function `ResourceProviderOperationList.MarshalJSON` has been removed
- Function `ListByCustomRPManifest.MarshalJSON` has been removed
- Function `ErrorDefinition.MarshalJSON` has been removed


## 0.4.0 (2022-04-15)
### Breaking Changes

- Function `*CustomResourceProviderClient.ListBySubscription` has been removed
- Function `*OperationsClient.List` has been removed
- Function `*AssociationsClient.ListAll` has been removed
- Function `*CustomResourceProviderClient.ListByResourceGroup` has been removed

### Features Added

- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New function `*AssociationsClient.NewListAllPager(string, *AssociationsClientListAllOptions) *runtime.Pager[AssociationsClientListAllResponse]`
- New function `*CustomResourceProviderClient.NewListBySubscriptionPager(*CustomResourceProviderClientListBySubscriptionOptions) *runtime.Pager[CustomResourceProviderClientListBySubscriptionResponse]`
- New function `*CustomResourceProviderClient.NewListByResourceGroupPager(string, *CustomResourceProviderClientListByResourceGroupOptions) *runtime.Pager[CustomResourceProviderClientListByResourceGroupResponse]`


## 0.3.0 (2022-04-11)
### Breaking Changes

- Function `*CustomResourceProviderClient.ListByResourceGroup` return value(s) have been changed from `(*CustomResourceProviderClientListByResourceGroupPager)` to `(*runtime.Pager[CustomResourceProviderClientListByResourceGroupResponse])`
- Function `*AssociationsClient.BeginDelete` return value(s) have been changed from `(AssociationsClientDeletePollerResponse, error)` to `(*armruntime.Poller[AssociationsClientDeleteResponse], error)`
- Function `*AssociationsClient.ListAll` return value(s) have been changed from `(*AssociationsClientListAllPager)` to `(*runtime.Pager[AssociationsClientListAllResponse])`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `NewAssociationsClient` return value(s) have been changed from `(*AssociationsClient)` to `(*AssociationsClient, error)`
- Function `*CustomResourceProviderClient.BeginDelete` return value(s) have been changed from `(CustomResourceProviderClientDeletePollerResponse, error)` to `(*armruntime.Poller[CustomResourceProviderClientDeleteResponse], error)`
- Function `*CustomResourceProviderClient.ListBySubscription` return value(s) have been changed from `(*CustomResourceProviderClientListBySubscriptionPager)` to `(*runtime.Pager[CustomResourceProviderClientListBySubscriptionResponse])`
- Function `*AssociationsClient.BeginCreateOrUpdate` return value(s) have been changed from `(AssociationsClientCreateOrUpdatePollerResponse, error)` to `(*armruntime.Poller[AssociationsClientCreateOrUpdateResponse], error)`
- Function `*CustomResourceProviderClient.BeginCreateOrUpdate` return value(s) have been changed from `(CustomResourceProviderClientCreateOrUpdatePollerResponse, error)` to `(*armruntime.Poller[CustomResourceProviderClientCreateOrUpdateResponse], error)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `NewCustomResourceProviderClient` return value(s) have been changed from `(*CustomResourceProviderClient)` to `(*CustomResourceProviderClient, error)`
- Function `*OperationsClientListPager.NextPage` has been removed
- Function `*CustomResourceProviderClientCreateOrUpdatePoller.FinalResponse` has been removed
- Function `AssociationsClientDeletePollerResponse.PollUntilDone` has been removed
- Function `*AssociationsClientCreateOrUpdatePoller.FinalResponse` has been removed
- Function `ActionRouting.ToPtr` has been removed
- Function `*CustomResourceProviderClientDeletePoller.Poll` has been removed
- Function `*CustomResourceProviderClientListBySubscriptionPager.NextPage` has been removed
- Function `*CustomResourceProviderClientDeletePollerResponse.Resume` has been removed
- Function `*CustomResourceProviderClientListByResourceGroupPager.PageResponse` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Function `*CustomResourceProviderClientListBySubscriptionPager.Err` has been removed
- Function `*CustomResourceProviderClientListByResourceGroupPager.Err` has been removed
- Function `ResourceTypeRouting.ToPtr` has been removed
- Function `*CustomResourceProviderClientCreateOrUpdatePoller.Done` has been removed
- Function `*AssociationsClientCreateOrUpdatePoller.Poll` has been removed
- Function `ProvisioningState.ToPtr` has been removed
- Function `*AssociationsClientCreateOrUpdatePoller.Done` has been removed
- Function `*OperationsClientListPager.Err` has been removed
- Function `*AssociationsClientCreateOrUpdatePoller.ResumeToken` has been removed
- Function `CustomResourceProviderClientDeletePollerResponse.PollUntilDone` has been removed
- Function `*AssociationsClientListAllPager.NextPage` has been removed
- Function `*CustomResourceProviderClientCreateOrUpdatePollerResponse.Resume` has been removed
- Function `CustomResourceProviderClientCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `*AssociationsClientListAllPager.PageResponse` has been removed
- Function `*AssociationsClientDeletePoller.FinalResponse` has been removed
- Function `*AssociationsClientListAllPager.Err` has been removed
- Function `*AssociationsClientDeletePoller.Poll` has been removed
- Function `*AssociationsClientDeletePollerResponse.Resume` has been removed
- Function `*CustomResourceProviderClientCreateOrUpdatePoller.ResumeToken` has been removed
- Function `*CustomResourceProviderClientDeletePoller.Done` has been removed
- Function `*CustomResourceProviderClientCreateOrUpdatePoller.Poll` has been removed
- Function `*AssociationsClientCreateOrUpdatePollerResponse.Resume` has been removed
- Function `*CustomResourceProviderClientListByResourceGroupPager.NextPage` has been removed
- Function `*CustomResourceProviderClientDeletePoller.FinalResponse` has been removed
- Function `*AssociationsClientDeletePoller.ResumeToken` has been removed
- Function `*CustomResourceProviderClientListBySubscriptionPager.PageResponse` has been removed
- Function `*CustomResourceProviderClientDeletePoller.ResumeToken` has been removed
- Function `*AssociationsClientDeletePoller.Done` has been removed
- Function `AssociationsClientCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `ValidationType.ToPtr` has been removed
- Struct `AssociationsClientCreateOrUpdatePoller` has been removed
- Struct `AssociationsClientCreateOrUpdatePollerResponse` has been removed
- Struct `AssociationsClientCreateOrUpdateResult` has been removed
- Struct `AssociationsClientDeletePoller` has been removed
- Struct `AssociationsClientDeletePollerResponse` has been removed
- Struct `AssociationsClientGetResult` has been removed
- Struct `AssociationsClientListAllPager` has been removed
- Struct `AssociationsClientListAllResult` has been removed
- Struct `CustomResourceProviderClientCreateOrUpdatePoller` has been removed
- Struct `CustomResourceProviderClientCreateOrUpdatePollerResponse` has been removed
- Struct `CustomResourceProviderClientCreateOrUpdateResult` has been removed
- Struct `CustomResourceProviderClientDeletePoller` has been removed
- Struct `CustomResourceProviderClientDeletePollerResponse` has been removed
- Struct `CustomResourceProviderClientGetResult` has been removed
- Struct `CustomResourceProviderClientListByResourceGroupPager` has been removed
- Struct `CustomResourceProviderClientListByResourceGroupResult` has been removed
- Struct `CustomResourceProviderClientListBySubscriptionPager` has been removed
- Struct `CustomResourceProviderClientListBySubscriptionResult` has been removed
- Struct `CustomResourceProviderClientUpdateResult` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Field `AssociationsClientCreateOrUpdateResult` of struct `AssociationsClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `AssociationsClientCreateOrUpdateResponse` has been removed
- Field `CustomResourceProviderClientUpdateResult` of struct `CustomResourceProviderClientUpdateResponse` has been removed
- Field `RawResponse` of struct `CustomResourceProviderClientUpdateResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `CustomResourceProviderClientCreateOrUpdateResult` of struct `CustomResourceProviderClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `CustomResourceProviderClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `AssociationsClientDeleteResponse` has been removed
- Field `AssociationsClientListAllResult` of struct `AssociationsClientListAllResponse` has been removed
- Field `RawResponse` of struct `AssociationsClientListAllResponse` has been removed
- Field `CustomResourceProviderClientListBySubscriptionResult` of struct `CustomResourceProviderClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `CustomResourceProviderClientListBySubscriptionResponse` has been removed
- Field `CustomResourceProviderClientListByResourceGroupResult` of struct `CustomResourceProviderClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `CustomResourceProviderClientListByResourceGroupResponse` has been removed
- Field `CustomResourceProviderClientGetResult` of struct `CustomResourceProviderClientGetResponse` has been removed
- Field `RawResponse` of struct `CustomResourceProviderClientGetResponse` has been removed
- Field `AssociationsClientGetResult` of struct `AssociationsClientGetResponse` has been removed
- Field `RawResponse` of struct `AssociationsClientGetResponse` has been removed
- Field `RawResponse` of struct `CustomResourceProviderClientDeleteResponse` has been removed

### Features Added

- New field `ResumeToken` in struct `AssociationsClientBeginCreateOrUpdateOptions`
- New field `ResumeToken` in struct `CustomResourceProviderClientBeginDeleteOptions`
- New anonymous field `CustomRPManifest` in struct `CustomResourceProviderClientUpdateResponse`
- New anonymous field `Association` in struct `AssociationsClientCreateOrUpdateResponse`
- New anonymous field `ResourceProviderOperationList` in struct `OperationsClientListResponse`
- New anonymous field `CustomRPManifest` in struct `CustomResourceProviderClientGetResponse`
- New anonymous field `Association` in struct `AssociationsClientGetResponse`
- New field `ResumeToken` in struct `AssociationsClientBeginDeleteOptions`
- New anonymous field `AssociationsList` in struct `AssociationsClientListAllResponse`
- New anonymous field `ListByCustomRPManifest` in struct `CustomResourceProviderClientListByResourceGroupResponse`
- New anonymous field `ListByCustomRPManifest` in struct `CustomResourceProviderClientListBySubscriptionResponse`
- New field `ResumeToken` in struct `CustomResourceProviderClientBeginCreateOrUpdateOptions`
- New anonymous field `CustomRPManifest` in struct `CustomResourceProviderClientCreateOrUpdateResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*AssociationsClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *AssociationsBeginDeleteOptions)` to `(context.Context, string, string, *AssociationsClientBeginDeleteOptions)`
- Function `*AssociationsClient.BeginDelete` return value(s) have been changed from `(AssociationsDeletePollerResponse, error)` to `(AssociationsClientDeletePollerResponse, error)`
- Function `*CustomResourceProviderClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, CustomRPManifest, *CustomResourceProviderBeginCreateOrUpdateOptions)` to `(context.Context, string, string, CustomRPManifest, *CustomResourceProviderClientBeginCreateOrUpdateOptions)`
- Function `*CustomResourceProviderClient.BeginCreateOrUpdate` return value(s) have been changed from `(CustomResourceProviderCreateOrUpdatePollerResponse, error)` to `(CustomResourceProviderClientCreateOrUpdatePollerResponse, error)`
- Function `*CustomResourceProviderClient.ListBySubscription` parameter(s) have been changed from `(*CustomResourceProviderListBySubscriptionOptions)` to `(*CustomResourceProviderClientListBySubscriptionOptions)`
- Function `*CustomResourceProviderClient.ListBySubscription` return value(s) have been changed from `(*CustomResourceProviderListBySubscriptionPager)` to `(*CustomResourceProviderClientListBySubscriptionPager)`
- Function `*CustomResourceProviderClient.Update` parameter(s) have been changed from `(context.Context, string, string, ResourceProvidersUpdate, *CustomResourceProviderUpdateOptions)` to `(context.Context, string, string, ResourceProvidersUpdate, *CustomResourceProviderClientUpdateOptions)`
- Function `*CustomResourceProviderClient.Update` return value(s) have been changed from `(CustomResourceProviderUpdateResponse, error)` to `(CustomResourceProviderClientUpdateResponse, error)`
- Function `*CustomResourceProviderClient.Get` parameter(s) have been changed from `(context.Context, string, string, *CustomResourceProviderGetOptions)` to `(context.Context, string, string, *CustomResourceProviderClientGetOptions)`
- Function `*CustomResourceProviderClient.Get` return value(s) have been changed from `(CustomResourceProviderGetResponse, error)` to `(CustomResourceProviderClientGetResponse, error)`
- Function `*AssociationsClient.ListAll` parameter(s) have been changed from `(string, *AssociationsListAllOptions)` to `(string, *AssociationsClientListAllOptions)`
- Function `*AssociationsClient.ListAll` return value(s) have been changed from `(*AssociationsListAllPager)` to `(*AssociationsClientListAllPager)`
- Function `*AssociationsClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, Association, *AssociationsBeginCreateOrUpdateOptions)` to `(context.Context, string, string, Association, *AssociationsClientBeginCreateOrUpdateOptions)`
- Function `*AssociationsClient.BeginCreateOrUpdate` return value(s) have been changed from `(AssociationsCreateOrUpdatePollerResponse, error)` to `(AssociationsClientCreateOrUpdatePollerResponse, error)`
- Function `*AssociationsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *AssociationsGetOptions)` to `(context.Context, string, string, *AssociationsClientGetOptions)`
- Function `*AssociationsClient.Get` return value(s) have been changed from `(AssociationsGetResponse, error)` to `(AssociationsClientGetResponse, error)`
- Function `*CustomResourceProviderClient.ListByResourceGroup` parameter(s) have been changed from `(string, *CustomResourceProviderListByResourceGroupOptions)` to `(string, *CustomResourceProviderClientListByResourceGroupOptions)`
- Function `*CustomResourceProviderClient.ListByResourceGroup` return value(s) have been changed from `(*CustomResourceProviderListByResourceGroupPager)` to `(*CustomResourceProviderClientListByResourceGroupPager)`
- Function `*CustomResourceProviderClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *CustomResourceProviderBeginDeleteOptions)` to `(context.Context, string, string, *CustomResourceProviderClientBeginDeleteOptions)`
- Function `*CustomResourceProviderClient.BeginDelete` return value(s) have been changed from `(CustomResourceProviderDeletePollerResponse, error)` to `(CustomResourceProviderClientDeletePollerResponse, error)`
- Function `*CustomResourceProviderCreateOrUpdatePoller.FinalResponse` has been removed
- Function `*AssociationsDeletePoller.FinalResponse` has been removed
- Function `*AssociationsCreateOrUpdatePollerResponse.Resume` has been removed
- Function `CustomResourceProviderCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `*CustomResourceProviderCreateOrUpdatePoller.ResumeToken` has been removed
- Function `*AssociationsListAllPager.NextPage` has been removed
- Function `*CustomResourceProviderCreateOrUpdatePoller.Poll` has been removed
- Function `AssociationsCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `*CustomResourceProviderDeletePoller.ResumeToken` has been removed
- Function `*AssociationsDeletePoller.Done` has been removed
- Function `*CustomResourceProviderDeletePoller.Poll` has been removed
- Function `*AssociationsCreateOrUpdatePoller.ResumeToken` has been removed
- Function `*AssociationsCreateOrUpdatePoller.FinalResponse` has been removed
- Function `*CustomResourceProviderListByResourceGroupPager.NextPage` has been removed
- Function `*CustomResourceProviderListByResourceGroupPager.PageResponse` has been removed
- Function `*CustomResourceProviderDeletePoller.FinalResponse` has been removed
- Function `*AssociationsDeletePoller.ResumeToken` has been removed
- Function `*AssociationsDeletePollerResponse.Resume` has been removed
- Function `*CustomResourceProviderListBySubscriptionPager.Err` has been removed
- Function `*CustomResourceProviderListBySubscriptionPager.NextPage` has been removed
- Function `*AssociationsDeletePoller.Poll` has been removed
- Function `CustomResourceProviderDeletePollerResponse.PollUntilDone` has been removed
- Function `*CustomResourceProviderListBySubscriptionPager.PageResponse` has been removed
- Function `AssociationsDeletePollerResponse.PollUntilDone` has been removed
- Function `*AssociationsListAllPager.PageResponse` has been removed
- Function `*CustomResourceProviderDeletePoller.Done` has been removed
- Function `*CustomResourceProviderCreateOrUpdatePollerResponse.Resume` has been removed
- Function `*CustomResourceProviderListByResourceGroupPager.Err` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Function `*CustomResourceProviderDeletePollerResponse.Resume` has been removed
- Function `*OperationsListPager.PageResponse` has been removed
- Function `*AssociationsCreateOrUpdatePoller.Poll` has been removed
- Function `*CustomResourceProviderCreateOrUpdatePoller.Done` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `*AssociationsListAllPager.Err` has been removed
- Function `*AssociationsCreateOrUpdatePoller.Done` has been removed
- Struct `AssociationsBeginCreateOrUpdateOptions` has been removed
- Struct `AssociationsBeginDeleteOptions` has been removed
- Struct `AssociationsCreateOrUpdatePoller` has been removed
- Struct `AssociationsCreateOrUpdatePollerResponse` has been removed
- Struct `AssociationsCreateOrUpdateResponse` has been removed
- Struct `AssociationsCreateOrUpdateResult` has been removed
- Struct `AssociationsDeletePoller` has been removed
- Struct `AssociationsDeletePollerResponse` has been removed
- Struct `AssociationsDeleteResponse` has been removed
- Struct `AssociationsGetOptions` has been removed
- Struct `AssociationsGetResponse` has been removed
- Struct `AssociationsGetResult` has been removed
- Struct `AssociationsListAllOptions` has been removed
- Struct `AssociationsListAllPager` has been removed
- Struct `AssociationsListAllResponse` has been removed
- Struct `AssociationsListAllResult` has been removed
- Struct `CustomResourceProviderBeginCreateOrUpdateOptions` has been removed
- Struct `CustomResourceProviderBeginDeleteOptions` has been removed
- Struct `CustomResourceProviderCreateOrUpdatePoller` has been removed
- Struct `CustomResourceProviderCreateOrUpdatePollerResponse` has been removed
- Struct `CustomResourceProviderCreateOrUpdateResponse` has been removed
- Struct `CustomResourceProviderCreateOrUpdateResult` has been removed
- Struct `CustomResourceProviderDeletePoller` has been removed
- Struct `CustomResourceProviderDeletePollerResponse` has been removed
- Struct `CustomResourceProviderDeleteResponse` has been removed
- Struct `CustomResourceProviderGetOptions` has been removed
- Struct `CustomResourceProviderGetResponse` has been removed
- Struct `CustomResourceProviderGetResult` has been removed
- Struct `CustomResourceProviderListByResourceGroupOptions` has been removed
- Struct `CustomResourceProviderListByResourceGroupPager` has been removed
- Struct `CustomResourceProviderListByResourceGroupResponse` has been removed
- Struct `CustomResourceProviderListByResourceGroupResult` has been removed
- Struct `CustomResourceProviderListBySubscriptionOptions` has been removed
- Struct `CustomResourceProviderListBySubscriptionPager` has been removed
- Struct `CustomResourceProviderListBySubscriptionResponse` has been removed
- Struct `CustomResourceProviderListBySubscriptionResult` has been removed
- Struct `CustomResourceProviderUpdateOptions` has been removed
- Struct `CustomResourceProviderUpdateResponse` has been removed
- Struct `CustomResourceProviderUpdateResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `CustomRPRouteDefinition` of struct `CustomRPActionRouteDefinition` has been removed
- Field `CustomRPRouteDefinition` of struct `CustomRPResourceTypeRouteDefinition` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed
- Field `Resource` of struct `CustomRPManifest` has been removed

### Features Added

- New function `*CustomResourceProviderClientDeletePollerResponse.Resume(context.Context, *CustomResourceProviderClient, string) error`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*CustomResourceProviderClientCreateOrUpdatePoller.Done() bool`
- New function `*AssociationsClientDeletePoller.ResumeToken() (string, error)`
- New function `*CustomResourceProviderClientCreateOrUpdatePoller.FinalResponse(context.Context) (CustomResourceProviderClientCreateOrUpdateResponse, error)`
- New function `*CustomResourceProviderClientListBySubscriptionPager.PageResponse() CustomResourceProviderClientListBySubscriptionResponse`
- New function `*AssociationsClientCreateOrUpdatePollerResponse.Resume(context.Context, *AssociationsClient, string) error`
- New function `CustomResourceProviderClientCreateOrUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (CustomResourceProviderClientCreateOrUpdateResponse, error)`
- New function `*CustomResourceProviderClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*CustomResourceProviderClientCreateOrUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*CustomResourceProviderClientListByResourceGroupPager.PageResponse() CustomResourceProviderClientListByResourceGroupResponse`
- New function `*CustomResourceProviderClientCreateOrUpdatePollerResponse.Resume(context.Context, *CustomResourceProviderClient, string) error`
- New function `*OperationsClientListPager.Err() error`
- New function `*AssociationsClientListAllPager.NextPage(context.Context) bool`
- New function `*AssociationsClientCreateOrUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*AssociationsClientDeletePollerResponse.Resume(context.Context, *AssociationsClient, string) error`
- New function `*CustomResourceProviderClientListBySubscriptionPager.NextPage(context.Context) bool`
- New function `CustomResourceProviderClientDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (CustomResourceProviderClientDeleteResponse, error)`
- New function `*AssociationsClientListAllPager.Err() error`
- New function `*CustomResourceProviderClientDeletePoller.Done() bool`
- New function `*CustomResourceProviderClientListByResourceGroupPager.Err() error`
- New function `AssociationsClientCreateOrUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (AssociationsClientCreateOrUpdateResponse, error)`
- New function `*AssociationsClientCreateOrUpdatePoller.FinalResponse(context.Context) (AssociationsClientCreateOrUpdateResponse, error)`
- New function `AssociationsClientDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (AssociationsClientDeleteResponse, error)`
- New function `*AssociationsClientListAllPager.PageResponse() AssociationsClientListAllResponse`
- New function `*CustomResourceProviderClientListBySubscriptionPager.Err() error`
- New function `*CustomResourceProviderClientCreateOrUpdatePoller.ResumeToken() (string, error)`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*CustomResourceProviderClientDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `*CustomResourceProviderClientDeletePoller.FinalResponse(context.Context) (CustomResourceProviderClientDeleteResponse, error)`
- New function `*AssociationsClientDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `*AssociationsClientDeletePoller.Done() bool`
- New function `*AssociationsClientCreateOrUpdatePoller.Done() bool`
- New function `*AssociationsClientCreateOrUpdatePoller.ResumeToken() (string, error)`
- New function `*CustomResourceProviderClientDeletePoller.ResumeToken() (string, error)`
- New function `*AssociationsClientDeletePoller.FinalResponse(context.Context) (AssociationsClientDeleteResponse, error)`
- New struct `AssociationsClientBeginCreateOrUpdateOptions`
- New struct `AssociationsClientBeginDeleteOptions`
- New struct `AssociationsClientCreateOrUpdatePoller`
- New struct `AssociationsClientCreateOrUpdatePollerResponse`
- New struct `AssociationsClientCreateOrUpdateResponse`
- New struct `AssociationsClientCreateOrUpdateResult`
- New struct `AssociationsClientDeletePoller`
- New struct `AssociationsClientDeletePollerResponse`
- New struct `AssociationsClientDeleteResponse`
- New struct `AssociationsClientGetOptions`
- New struct `AssociationsClientGetResponse`
- New struct `AssociationsClientGetResult`
- New struct `AssociationsClientListAllOptions`
- New struct `AssociationsClientListAllPager`
- New struct `AssociationsClientListAllResponse`
- New struct `AssociationsClientListAllResult`
- New struct `CustomResourceProviderClientBeginCreateOrUpdateOptions`
- New struct `CustomResourceProviderClientBeginDeleteOptions`
- New struct `CustomResourceProviderClientCreateOrUpdatePoller`
- New struct `CustomResourceProviderClientCreateOrUpdatePollerResponse`
- New struct `CustomResourceProviderClientCreateOrUpdateResponse`
- New struct `CustomResourceProviderClientCreateOrUpdateResult`
- New struct `CustomResourceProviderClientDeletePoller`
- New struct `CustomResourceProviderClientDeletePollerResponse`
- New struct `CustomResourceProviderClientDeleteResponse`
- New struct `CustomResourceProviderClientGetOptions`
- New struct `CustomResourceProviderClientGetResponse`
- New struct `CustomResourceProviderClientGetResult`
- New struct `CustomResourceProviderClientListByResourceGroupOptions`
- New struct `CustomResourceProviderClientListByResourceGroupPager`
- New struct `CustomResourceProviderClientListByResourceGroupResponse`
- New struct `CustomResourceProviderClientListByResourceGroupResult`
- New struct `CustomResourceProviderClientListBySubscriptionOptions`
- New struct `CustomResourceProviderClientListBySubscriptionPager`
- New struct `CustomResourceProviderClientListBySubscriptionResponse`
- New struct `CustomResourceProviderClientListBySubscriptionResult`
- New struct `CustomResourceProviderClientUpdateOptions`
- New struct `CustomResourceProviderClientUpdateResponse`
- New struct `CustomResourceProviderClientUpdateResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New field `Tags` in struct `CustomRPManifest`
- New field `ID` in struct `CustomRPManifest`
- New field `Name` in struct `CustomRPManifest`
- New field `Type` in struct `CustomRPManifest`
- New field `Location` in struct `CustomRPManifest`
- New field `Error` in struct `ErrorResponse`
- New field `Endpoint` in struct `CustomRPResourceTypeRouteDefinition`
- New field `Name` in struct `CustomRPResourceTypeRouteDefinition`
- New field `Endpoint` in struct `CustomRPActionRouteDefinition`
- New field `Name` in struct `CustomRPActionRouteDefinition`


## 0.1.0 (2021-12-22)

- Init release.
