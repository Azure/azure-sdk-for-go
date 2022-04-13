# Release History

## 0.3.0 (2022-04-13)
### Breaking Changes

- Function `*MultipleActivationKeysClient.ListByResourceGroup` return value(s) have been changed from `(*MultipleActivationKeysClientListByResourceGroupPager)` to `(*runtime.Pager[MultipleActivationKeysClientListByResourceGroupResponse])`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `NewMultipleActivationKeysClient` return value(s) have been changed from `(*MultipleActivationKeysClient)` to `(*MultipleActivationKeysClient, error)`
- Function `*MultipleActivationKeysClient.BeginCreate` return value(s) have been changed from `(MultipleActivationKeysClientCreatePollerResponse, error)` to `(*armruntime.Poller[MultipleActivationKeysClientCreateResponse], error)`
- Function `*MultipleActivationKeysClient.List` return value(s) have been changed from `(*MultipleActivationKeysClientListPager)` to `(*runtime.Pager[MultipleActivationKeysClientListResponse])`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `MultipleActivationKeysClientCreatePollerResponse.PollUntilDone` has been removed
- Function `*OperationsClientListPager.Err` has been removed
- Function `*MultipleActivationKeysClientListByResourceGroupPager.Err` has been removed
- Function `*MultipleActivationKeysClientListPager.PageResponse` has been removed
- Function `*MultipleActivationKeysClientCreatePollerResponse.Resume` has been removed
- Function `*MultipleActivationKeysClientListPager.NextPage` has been removed
- Function `*MultipleActivationKeysClientCreatePoller.ResumeToken` has been removed
- Function `OsType.ToPtr` has been removed
- Function `*MultipleActivationKeysClientCreatePoller.FinalResponse` has been removed
- Function `*MultipleActivationKeysClientListByResourceGroupPager.NextPage` has been removed
- Function `SupportType.ToPtr` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Function `*MultipleActivationKeysClientListPager.Err` has been removed
- Function `*OperationsClientListPager.NextPage` has been removed
- Function `*MultipleActivationKeysClientCreatePoller.Poll` has been removed
- Function `ProvisioningState.ToPtr` has been removed
- Function `*MultipleActivationKeysClientCreatePoller.Done` has been removed
- Function `*MultipleActivationKeysClientListByResourceGroupPager.PageResponse` has been removed
- Struct `MultipleActivationKeysClientCreatePoller` has been removed
- Struct `MultipleActivationKeysClientCreatePollerResponse` has been removed
- Struct `MultipleActivationKeysClientCreateResult` has been removed
- Struct `MultipleActivationKeysClientGetResult` has been removed
- Struct `MultipleActivationKeysClientListByResourceGroupPager` has been removed
- Struct `MultipleActivationKeysClientListByResourceGroupResult` has been removed
- Struct `MultipleActivationKeysClientListPager` has been removed
- Struct `MultipleActivationKeysClientListResult` has been removed
- Struct `MultipleActivationKeysClientUpdateResult` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Field `MultipleActivationKeysClientListByResourceGroupResult` of struct `MultipleActivationKeysClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `MultipleActivationKeysClientListByResourceGroupResponse` has been removed
- Field `MultipleActivationKeysClientCreateResult` of struct `MultipleActivationKeysClientCreateResponse` has been removed
- Field `RawResponse` of struct `MultipleActivationKeysClientCreateResponse` has been removed
- Field `MultipleActivationKeysClientUpdateResult` of struct `MultipleActivationKeysClientUpdateResponse` has been removed
- Field `RawResponse` of struct `MultipleActivationKeysClientUpdateResponse` has been removed
- Field `MultipleActivationKeysClientGetResult` of struct `MultipleActivationKeysClientGetResponse` has been removed
- Field `RawResponse` of struct `MultipleActivationKeysClientGetResponse` has been removed
- Field `MultipleActivationKeysClientListResult` of struct `MultipleActivationKeysClientListResponse` has been removed
- Field `RawResponse` of struct `MultipleActivationKeysClientListResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `MultipleActivationKeysClientDeleteResponse` has been removed

### Features Added

- New anonymous field `MultipleActivationKeyList` in struct `MultipleActivationKeysClientListResponse`
- New anonymous field `MultipleActivationKeyList` in struct `MultipleActivationKeysClientListByResourceGroupResponse`
- New anonymous field `MultipleActivationKey` in struct `MultipleActivationKeysClientGetResponse`
- New anonymous field `MultipleActivationKey` in struct `MultipleActivationKeysClientUpdateResponse`
- New field `ResumeToken` in struct `MultipleActivationKeysClientBeginCreateOptions`
- New anonymous field `MultipleActivationKey` in struct `MultipleActivationKeysClientCreateResponse`
- New anonymous field `OperationList` in struct `OperationsClientListResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*MultipleActivationKeysClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, MultipleActivationKey, *MultipleActivationKeysBeginCreateOptions)` to `(context.Context, string, string, MultipleActivationKey, *MultipleActivationKeysClientBeginCreateOptions)`
- Function `*MultipleActivationKeysClient.BeginCreate` return value(s) have been changed from `(MultipleActivationKeysCreatePollerResponse, error)` to `(MultipleActivationKeysClientCreatePollerResponse, error)`
- Function `*MultipleActivationKeysClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *MultipleActivationKeysDeleteOptions)` to `(context.Context, string, string, *MultipleActivationKeysClientDeleteOptions)`
- Function `*MultipleActivationKeysClient.Delete` return value(s) have been changed from `(MultipleActivationKeysDeleteResponse, error)` to `(MultipleActivationKeysClientDeleteResponse, error)`
- Function `*MultipleActivationKeysClient.ListByResourceGroup` parameter(s) have been changed from `(string, *MultipleActivationKeysListByResourceGroupOptions)` to `(string, *MultipleActivationKeysClientListByResourceGroupOptions)`
- Function `*MultipleActivationKeysClient.ListByResourceGroup` return value(s) have been changed from `(*MultipleActivationKeysListByResourceGroupPager)` to `(*MultipleActivationKeysClientListByResourceGroupPager)`
- Function `*MultipleActivationKeysClient.Get` parameter(s) have been changed from `(context.Context, string, string, *MultipleActivationKeysGetOptions)` to `(context.Context, string, string, *MultipleActivationKeysClientGetOptions)`
- Function `*MultipleActivationKeysClient.Get` return value(s) have been changed from `(MultipleActivationKeysGetResponse, error)` to `(MultipleActivationKeysClientGetResponse, error)`
- Function `*MultipleActivationKeysClient.List` parameter(s) have been changed from `(*MultipleActivationKeysListOptions)` to `(*MultipleActivationKeysClientListOptions)`
- Function `*MultipleActivationKeysClient.List` return value(s) have been changed from `(*MultipleActivationKeysListPager)` to `(*MultipleActivationKeysClientListPager)`
- Function `*MultipleActivationKeysClient.Update` parameter(s) have been changed from `(context.Context, string, string, MultipleActivationKeyUpdate, *MultipleActivationKeysUpdateOptions)` to `(context.Context, string, string, MultipleActivationKeyUpdate, *MultipleActivationKeysClientUpdateOptions)`
- Function `*MultipleActivationKeysClient.Update` return value(s) have been changed from `(MultipleActivationKeysUpdateResponse, error)` to `(MultipleActivationKeysClientUpdateResponse, error)`
- Function `Resource.MarshalJSON` has been removed
- Function `*MultipleActivationKeysCreatePoller.Poll` has been removed
- Function `MultipleActivationKeysCreatePollerResponse.PollUntilDone` has been removed
- Function `*MultipleActivationKeysListPager.Err` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `*MultipleActivationKeysListPager.NextPage` has been removed
- Function `*OperationsListPager.PageResponse` has been removed
- Function `*MultipleActivationKeysCreatePoller.Done` has been removed
- Function `*MultipleActivationKeysListByResourceGroupPager.Err` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Function `*MultipleActivationKeysListPager.PageResponse` has been removed
- Function `*MultipleActivationKeysCreatePoller.ResumeToken` has been removed
- Function `*MultipleActivationKeysCreatePoller.FinalResponse` has been removed
- Function `*MultipleActivationKeysListByResourceGroupPager.PageResponse` has been removed
- Function `*MultipleActivationKeysCreatePollerResponse.Resume` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `*MultipleActivationKeysListByResourceGroupPager.NextPage` has been removed
- Struct `MultipleActivationKeysBeginCreateOptions` has been removed
- Struct `MultipleActivationKeysCreatePoller` has been removed
- Struct `MultipleActivationKeysCreatePollerResponse` has been removed
- Struct `MultipleActivationKeysCreateResponse` has been removed
- Struct `MultipleActivationKeysCreateResult` has been removed
- Struct `MultipleActivationKeysDeleteOptions` has been removed
- Struct `MultipleActivationKeysDeleteResponse` has been removed
- Struct `MultipleActivationKeysGetOptions` has been removed
- Struct `MultipleActivationKeysGetResponse` has been removed
- Struct `MultipleActivationKeysGetResult` has been removed
- Struct `MultipleActivationKeysListByResourceGroupOptions` has been removed
- Struct `MultipleActivationKeysListByResourceGroupPager` has been removed
- Struct `MultipleActivationKeysListByResourceGroupResponse` has been removed
- Struct `MultipleActivationKeysListByResourceGroupResult` has been removed
- Struct `MultipleActivationKeysListOptions` has been removed
- Struct `MultipleActivationKeysListPager` has been removed
- Struct `MultipleActivationKeysListResponse` has been removed
- Struct `MultipleActivationKeysListResult` has been removed
- Struct `MultipleActivationKeysUpdateOptions` has been removed
- Struct `MultipleActivationKeysUpdateResponse` has been removed
- Struct `MultipleActivationKeysUpdateResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `TrackedResource` of struct `MultipleActivationKey` has been removed
- Field `Resource` of struct `TrackedResource` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed

### Features Added

- New function `*MultipleActivationKeysClientListByResourceGroupPager.Err() error`
- New function `*MultipleActivationKeysClientListPager.PageResponse() MultipleActivationKeysClientListResponse`
- New function `*MultipleActivationKeysClientListPager.NextPage(context.Context) bool`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*MultipleActivationKeysClientCreatePoller.ResumeToken() (string, error)`
- New function `*OperationsClientListPager.Err() error`
- New function `*MultipleActivationKeysClientCreatePoller.Done() bool`
- New function `*MultipleActivationKeysClientListByResourceGroupPager.PageResponse() MultipleActivationKeysClientListByResourceGroupResponse`
- New function `*MultipleActivationKeysClientListPager.Err() error`
- New function `*MultipleActivationKeysClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*MultipleActivationKeysClientCreatePollerResponse.Resume(context.Context, *MultipleActivationKeysClient, string) error`
- New function `MultipleActivationKeysClientCreatePollerResponse.PollUntilDone(context.Context, time.Duration) (MultipleActivationKeysClientCreateResponse, error)`
- New function `*MultipleActivationKeysClientCreatePoller.FinalResponse(context.Context) (MultipleActivationKeysClientCreateResponse, error)`
- New function `*MultipleActivationKeysClientCreatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New struct `MultipleActivationKeysClientBeginCreateOptions`
- New struct `MultipleActivationKeysClientCreatePoller`
- New struct `MultipleActivationKeysClientCreatePollerResponse`
- New struct `MultipleActivationKeysClientCreateResponse`
- New struct `MultipleActivationKeysClientCreateResult`
- New struct `MultipleActivationKeysClientDeleteOptions`
- New struct `MultipleActivationKeysClientDeleteResponse`
- New struct `MultipleActivationKeysClientGetOptions`
- New struct `MultipleActivationKeysClientGetResponse`
- New struct `MultipleActivationKeysClientGetResult`
- New struct `MultipleActivationKeysClientListByResourceGroupOptions`
- New struct `MultipleActivationKeysClientListByResourceGroupPager`
- New struct `MultipleActivationKeysClientListByResourceGroupResponse`
- New struct `MultipleActivationKeysClientListByResourceGroupResult`
- New struct `MultipleActivationKeysClientListOptions`
- New struct `MultipleActivationKeysClientListPager`
- New struct `MultipleActivationKeysClientListResponse`
- New struct `MultipleActivationKeysClientListResult`
- New struct `MultipleActivationKeysClientUpdateOptions`
- New struct `MultipleActivationKeysClientUpdateResponse`
- New struct `MultipleActivationKeysClientUpdateResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New field `Name` in struct `TrackedResource`
- New field `Type` in struct `TrackedResource`
- New field `ID` in struct `TrackedResource`
- New field `Tags` in struct `MultipleActivationKey`
- New field `ID` in struct `MultipleActivationKey`
- New field `Name` in struct `MultipleActivationKey`
- New field `Type` in struct `MultipleActivationKey`
- New field `Location` in struct `MultipleActivationKey`
- New field `Error` in struct `ErrorResponse`


## 0.1.0 (2021-12-22)

- Init release.
