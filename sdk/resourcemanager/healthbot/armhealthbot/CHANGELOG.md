# Release History

## 1.0.0 (2022-05-18)
### Breaking Changes

- Function `*BotsClient.BeginCreate` return value(s) have been changed from `(*armruntime.Poller[BotsClientCreateResponse], error)` to `(*runtime.Poller[BotsClientCreateResponse], error)`
- Function `*BotsClient.BeginDelete` return value(s) have been changed from `(*armruntime.Poller[BotsClientDeleteResponse], error)` to `(*runtime.Poller[BotsClientDeleteResponse], error)`
- Function `BotResponseList.MarshalJSON` has been removed
- Function `ErrorError.MarshalJSON` has been removed
- Function `AvailableOperations.MarshalJSON` has been removed


## 0.4.0 (2022-04-15)
### Breaking Changes

- Function `*BotsClient.List` has been removed
- Function `*BotsClient.ListByResourceGroup` has been removed
- Function `*OperationsClient.List` has been removed

### Features Added

- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New function `*BotsClient.NewListPager(*BotsClientListOptions) *runtime.Pager[BotsClientListResponse]`
- New function `*BotsClient.NewListByResourceGroupPager(string, *BotsClientListByResourceGroupOptions) *runtime.Pager[BotsClientListByResourceGroupResponse]`


## 0.3.0 (2022-04-11)
### Breaking Changes

- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `*BotsClient.ListByResourceGroup` return value(s) have been changed from `(*BotsClientListByResourceGroupPager)` to `(*runtime.Pager[BotsClientListByResourceGroupResponse])`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `*BotsClient.BeginCreate` return value(s) have been changed from `(BotsClientCreatePollerResponse, error)` to `(*armruntime.Poller[BotsClientCreateResponse], error)`
- Function `*BotsClient.List` return value(s) have been changed from `(*BotsClientListPager)` to `(*runtime.Pager[BotsClientListResponse])`
- Function `*BotsClient.BeginDelete` return value(s) have been changed from `(BotsClientDeletePollerResponse, error)` to `(*armruntime.Poller[BotsClientDeleteResponse], error)`
- Function `NewBotsClient` return value(s) have been changed from `(*BotsClient)` to `(*BotsClient, error)`
- Type of `ErrorAdditionalInfo.Info` has been changed from `map[string]interface{}` to `interface{}`
- Type of `OperationDetail.Properties` has been changed from `map[string]interface{}` to `interface{}`
- Function `SKUName.ToPtr` has been removed
- Function `BotsClientCreatePollerResponse.PollUntilDone` has been removed
- Function `*BotsClientCreatePoller.ResumeToken` has been removed
- Function `*BotsClientDeletePoller.ResumeToken` has been removed
- Function `IdentityType.ToPtr` has been removed
- Function `*OperationsClientListPager.Err` has been removed
- Function `*BotsClientListByResourceGroupPager.PageResponse` has been removed
- Function `*BotsClientCreatePoller.Poll` has been removed
- Function `*OperationsClientListPager.NextPage` has been removed
- Function `*BotsClientDeletePoller.FinalResponse` has been removed
- Function `BotsClientDeletePollerResponse.PollUntilDone` has been removed
- Function `*BotsClientListPager.PageResponse` has been removed
- Function `*BotsClientCreatePoller.Done` has been removed
- Function `*BotsClientDeletePollerResponse.Resume` has been removed
- Function `*BotsClientDeletePoller.Poll` has been removed
- Function `*BotsClientListByResourceGroupPager.Err` has been removed
- Function `*BotsClientCreatePoller.FinalResponse` has been removed
- Function `ResourceIdentityType.ToPtr` has been removed
- Function `*BotsClientDeletePoller.Done` has been removed
- Function `*BotsClientListByResourceGroupPager.NextPage` has been removed
- Function `*BotsClientCreatePollerResponse.Resume` has been removed
- Function `*BotsClientListPager.Err` has been removed
- Function `*BotsClientListPager.NextPage` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Struct `BotsClientCreatePoller` has been removed
- Struct `BotsClientCreatePollerResponse` has been removed
- Struct `BotsClientCreateResult` has been removed
- Struct `BotsClientDeletePoller` has been removed
- Struct `BotsClientDeletePollerResponse` has been removed
- Struct `BotsClientGetResult` has been removed
- Struct `BotsClientListByResourceGroupPager` has been removed
- Struct `BotsClientListByResourceGroupResult` has been removed
- Struct `BotsClientListPager` has been removed
- Struct `BotsClientListResult` has been removed
- Struct `BotsClientUpdateResult` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Field `BotsClientListResult` of struct `BotsClientListResponse` has been removed
- Field `RawResponse` of struct `BotsClientListResponse` has been removed
- Field `BotsClientGetResult` of struct `BotsClientGetResponse` has been removed
- Field `RawResponse` of struct `BotsClientGetResponse` has been removed
- Field `BotsClientListByResourceGroupResult` of struct `BotsClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `BotsClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `BotsClientDeleteResponse` has been removed
- Field `BotsClientUpdateResult` of struct `BotsClientUpdateResponse` has been removed
- Field `RawResponse` of struct `BotsClientUpdateResponse` has been removed
- Field `BotsClientCreateResult` of struct `BotsClientCreateResponse` has been removed
- Field `RawResponse` of struct `BotsClientCreateResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed

### Features Added

- New anonymous field `BotResponseList` in struct `BotsClientListResponse`
- New field `ResumeToken` in struct `BotsClientBeginCreateOptions`
- New anonymous field `BotResponseList` in struct `BotsClientListByResourceGroupResponse`
- New anonymous field `AvailableOperations` in struct `OperationsClientListResponse`
- New anonymous field `HealthBot` in struct `BotsClientGetResponse`
- New anonymous field `HealthBot` in struct `BotsClientUpdateResponse`
- New field `ResumeToken` in struct `BotsClientBeginDeleteOptions`
- New anonymous field `HealthBot` in struct `BotsClientCreateResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*BotsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *BotsGetOptions)` to `(context.Context, string, string, *BotsClientGetOptions)`
- Function `*BotsClient.Get` return value(s) have been changed from `(BotsGetResponse, error)` to `(BotsClientGetResponse, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*BotsClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, HealthBot, *BotsBeginCreateOptions)` to `(context.Context, string, string, HealthBot, *BotsClientBeginCreateOptions)`
- Function `*BotsClient.BeginCreate` return value(s) have been changed from `(BotsCreatePollerResponse, error)` to `(BotsClientCreatePollerResponse, error)`
- Function `*BotsClient.Update` parameter(s) have been changed from `(context.Context, string, string, HealthBotUpdateParameters, *BotsUpdateOptions)` to `(context.Context, string, string, UpdateParameters, *BotsClientUpdateOptions)`
- Function `*BotsClient.Update` return value(s) have been changed from `(BotsUpdateResponse, error)` to `(BotsClientUpdateResponse, error)`
- Function `*BotsClient.List` parameter(s) have been changed from `(*BotsListOptions)` to `(*BotsClientListOptions)`
- Function `*BotsClient.List` return value(s) have been changed from `(*BotsListPager)` to `(*BotsClientListPager)`
- Function `*BotsClient.ListByResourceGroup` parameter(s) have been changed from `(string, *BotsListByResourceGroupOptions)` to `(string, *BotsClientListByResourceGroupOptions)`
- Function `*BotsClient.ListByResourceGroup` return value(s) have been changed from `(*BotsListByResourceGroupPager)` to `(*BotsClientListByResourceGroupPager)`
- Function `*BotsClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *BotsBeginDeleteOptions)` to `(context.Context, string, string, *BotsClientBeginDeleteOptions)`
- Function `*BotsClient.BeginDelete` return value(s) have been changed from `(BotsDeletePollerResponse, error)` to `(BotsClientDeletePollerResponse, error)`
- Type of `HealthBot.Properties` has been changed from `*HealthBotProperties` to `*Properties`
- Function `Error.Error` has been removed
- Function `Resource.MarshalJSON` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Function `*OperationsListPager.PageResponse` has been removed
- Function `*BotsCreatePoller.FinalResponse` has been removed
- Function `BotsCreatePollerResponse.PollUntilDone` has been removed
- Function `*BotsCreatePoller.ResumeToken` has been removed
- Function `*BotsCreatePoller.Poll` has been removed
- Function `*BotsCreatePollerResponse.Resume` has been removed
- Function `*BotsDeletePoller.ResumeToken` has been removed
- Function `*BotsListByResourceGroupPager.Err` has been removed
- Function `*BotsListPager.Err` has been removed
- Function `*BotsDeletePoller.FinalResponse` has been removed
- Function `*BotsListPager.NextPage` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `*BotsListByResourceGroupPager.NextPage` has been removed
- Function `*BotsCreatePoller.Done` has been removed
- Function `*BotsDeletePollerResponse.Resume` has been removed
- Function `*BotsDeletePoller.Poll` has been removed
- Function `*BotsDeletePoller.Done` has been removed
- Function `BotsDeletePollerResponse.PollUntilDone` has been removed
- Function `*BotsListByResourceGroupPager.PageResponse` has been removed
- Function `HealthBotUpdateParameters.MarshalJSON` has been removed
- Function `*BotsListPager.PageResponse` has been removed
- Struct `BotsBeginCreateOptions` has been removed
- Struct `BotsBeginDeleteOptions` has been removed
- Struct `BotsCreatePoller` has been removed
- Struct `BotsCreatePollerResponse` has been removed
- Struct `BotsCreateResponse` has been removed
- Struct `BotsCreateResult` has been removed
- Struct `BotsDeletePoller` has been removed
- Struct `BotsDeletePollerResponse` has been removed
- Struct `BotsDeleteResponse` has been removed
- Struct `BotsGetOptions` has been removed
- Struct `BotsGetResponse` has been removed
- Struct `BotsGetResult` has been removed
- Struct `BotsListByResourceGroupOptions` has been removed
- Struct `BotsListByResourceGroupPager` has been removed
- Struct `BotsListByResourceGroupResponse` has been removed
- Struct `BotsListByResourceGroupResult` has been removed
- Struct `BotsListOptions` has been removed
- Struct `BotsListPager` has been removed
- Struct `BotsListResponse` has been removed
- Struct `BotsListResult` has been removed
- Struct `BotsUpdateOptions` has been removed
- Struct `BotsUpdateResponse` has been removed
- Struct `BotsUpdateResult` has been removed
- Struct `HealthBotProperties` has been removed
- Struct `HealthBotUpdateParameters` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `Resource` of struct `TrackedResource` has been removed
- Field `TrackedResource` of struct `HealthBot` has been removed
- Field `InnerError` of struct `Error` has been removed

### Features Added

- New function `*BotsClientListByResourceGroupPager.PageResponse() BotsClientListByResourceGroupResponse`
- New function `*BotsClientListPager.Err() error`
- New function `*BotsClientDeletePoller.ResumeToken() (string, error)`
- New function `*BotsClientCreatePoller.FinalResponse(context.Context) (BotsClientCreateResponse, error)`
- New function `*BotsClientCreatePoller.Done() bool`
- New function `*BotsClientListPager.PageResponse() BotsClientListResponse`
- New function `*BotsClientCreatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*BotsClientCreatePoller.ResumeToken() (string, error)`
- New function `BotsClientCreatePollerResponse.PollUntilDone(context.Context, time.Duration) (BotsClientCreateResponse, error)`
- New function `*BotsClientDeletePoller.Done() bool`
- New function `BotsClientDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (BotsClientDeleteResponse, error)`
- New function `*BotsClientDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `*BotsClientListPager.NextPage(context.Context) bool`
- New function `UpdateParameters.MarshalJSON() ([]byte, error)`
- New function `*BotsClientDeletePoller.FinalResponse(context.Context) (BotsClientDeleteResponse, error)`
- New function `*BotsClientListByResourceGroupPager.Err() error`
- New function `*BotsClientDeletePollerResponse.Resume(context.Context, *BotsClient, string) error`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*OperationsClientListPager.Err() error`
- New function `*BotsClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*BotsClientCreatePollerResponse.Resume(context.Context, *BotsClient, string) error`
- New struct `BotsClientBeginCreateOptions`
- New struct `BotsClientBeginDeleteOptions`
- New struct `BotsClientCreatePoller`
- New struct `BotsClientCreatePollerResponse`
- New struct `BotsClientCreateResponse`
- New struct `BotsClientCreateResult`
- New struct `BotsClientDeletePoller`
- New struct `BotsClientDeletePollerResponse`
- New struct `BotsClientDeleteResponse`
- New struct `BotsClientGetOptions`
- New struct `BotsClientGetResponse`
- New struct `BotsClientGetResult`
- New struct `BotsClientListByResourceGroupOptions`
- New struct `BotsClientListByResourceGroupPager`
- New struct `BotsClientListByResourceGroupResponse`
- New struct `BotsClientListByResourceGroupResult`
- New struct `BotsClientListOptions`
- New struct `BotsClientListPager`
- New struct `BotsClientListResponse`
- New struct `BotsClientListResult`
- New struct `BotsClientUpdateOptions`
- New struct `BotsClientUpdateResponse`
- New struct `BotsClientUpdateResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New struct `Properties`
- New struct `UpdateParameters`
- New field `ID` in struct `TrackedResource`
- New field `Name` in struct `TrackedResource`
- New field `SystemData` in struct `TrackedResource`
- New field `Type` in struct `TrackedResource`
- New field `Error` in struct `Error`
- New field `Name` in struct `HealthBot`
- New field `Type` in struct `HealthBot`
- New field `Tags` in struct `HealthBot`
- New field `SystemData` in struct `HealthBot`
- New field `Location` in struct `HealthBot`
- New field `ID` in struct `HealthBot`


## 0.1.0 (2021-12-07)

- Init release.
