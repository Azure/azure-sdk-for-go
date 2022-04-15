# Release History

## 0.4.0 (2022-04-15)
### Breaking Changes

- Function `*OperationsClient.List` has been removed
- Function `*LedgerClient.ListByResourceGroup` has been removed
- Function `*LedgerClient.ListBySubscription` has been removed

### Features Added

- New function `*LedgerClient.NewListBySubscriptionPager(*LedgerClientListBySubscriptionOptions) *runtime.Pager[LedgerClientListBySubscriptionResponse]`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New function `*LedgerClient.NewListByResourceGroupPager(string, *LedgerClientListByResourceGroupOptions) *runtime.Pager[LedgerClientListByResourceGroupResponse]`


## 0.3.0 (2022-04-11)
### Breaking Changes

- Function `*LedgerClient.ListByResourceGroup` return value(s) have been changed from `(*LedgerClientListByResourceGroupPager)` to `(*runtime.Pager[LedgerClientListByResourceGroupResponse])`
- Function `*LedgerClient.BeginCreate` return value(s) have been changed from `(LedgerClientCreatePollerResponse, error)` to `(*armruntime.Poller[LedgerClientCreateResponse], error)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `NewLedgerClient` return value(s) have been changed from `(*LedgerClient)` to `(*LedgerClient, error)`
- Function `*LedgerClient.ListBySubscription` return value(s) have been changed from `(*LedgerClientListBySubscriptionPager)` to `(*runtime.Pager[LedgerClientListBySubscriptionResponse])`
- Function `*LedgerClient.BeginDelete` return value(s) have been changed from `(LedgerClientDeletePollerResponse, error)` to `(*armruntime.Poller[LedgerClientDeleteResponse], error)`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `*LedgerClient.BeginUpdate` return value(s) have been changed from `(LedgerClientUpdatePollerResponse, error)` to `(*armruntime.Poller[LedgerClientUpdateResponse], error)`
- Function `NewClient` return value(s) have been changed from `(*Client)` to `(*Client, error)`
- Type of `ErrorAdditionalInfo.Info` has been changed from `map[string]interface{}` to `interface{}`
- Function `CreatedByType.ToPtr` has been removed
- Function `*LedgerClientCreatePoller.FinalResponse` has been removed
- Function `*LedgerClientCreatePoller.Poll` has been removed
- Function `*LedgerClientUpdatePoller.Poll` has been removed
- Function `*LedgerClientDeletePoller.Poll` has been removed
- Function `CheckNameAvailabilityReason.ToPtr` has been removed
- Function `*LedgerClientCreatePoller.ResumeToken` has been removed
- Function `*LedgerClientCreatePollerResponse.Resume` has been removed
- Function `*LedgerClientListByResourceGroupPager.Err` has been removed
- Function `*LedgerClientListByResourceGroupPager.PageResponse` has been removed
- Function `*LedgerClientUpdatePollerResponse.Resume` has been removed
- Function `*LedgerClientUpdatePoller.FinalResponse` has been removed
- Function `*LedgerClientCreatePoller.Done` has been removed
- Function `*LedgerClientDeletePoller.ResumeToken` has been removed
- Function `*OperationsClientListPager.NextPage` has been removed
- Function `LedgerType.ToPtr` has been removed
- Function `LedgerClientCreatePollerResponse.PollUntilDone` has been removed
- Function `*OperationsClientListPager.Err` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Function `ProvisioningState.ToPtr` has been removed
- Function `LedgerClientDeletePollerResponse.PollUntilDone` has been removed
- Function `LedgerRoleName.ToPtr` has been removed
- Function `*LedgerClientListByResourceGroupPager.NextPage` has been removed
- Function `*LedgerClientListBySubscriptionPager.NextPage` has been removed
- Function `*LedgerClientUpdatePoller.Done` has been removed
- Function `*LedgerClientListBySubscriptionPager.PageResponse` has been removed
- Function `*LedgerClientDeletePollerResponse.Resume` has been removed
- Function `*LedgerClientDeletePoller.Done` has been removed
- Function `*LedgerClientDeletePoller.FinalResponse` has been removed
- Function `*LedgerClientUpdatePoller.ResumeToken` has been removed
- Function `LedgerClientUpdatePollerResponse.PollUntilDone` has been removed
- Function `*LedgerClientListBySubscriptionPager.Err` has been removed
- Struct `ClientCheckNameAvailabilityResult` has been removed
- Struct `LedgerClientCreatePoller` has been removed
- Struct `LedgerClientCreatePollerResponse` has been removed
- Struct `LedgerClientCreateResult` has been removed
- Struct `LedgerClientDeletePoller` has been removed
- Struct `LedgerClientDeletePollerResponse` has been removed
- Struct `LedgerClientGetResult` has been removed
- Struct `LedgerClientListByResourceGroupPager` has been removed
- Struct `LedgerClientListByResourceGroupResult` has been removed
- Struct `LedgerClientListBySubscriptionPager` has been removed
- Struct `LedgerClientListBySubscriptionResult` has been removed
- Struct `LedgerClientUpdatePoller` has been removed
- Struct `LedgerClientUpdatePollerResponse` has been removed
- Struct `LedgerClientUpdateResult` has been removed
- Struct `Location` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Field `LedgerClientCreateResult` of struct `LedgerClientCreateResponse` has been removed
- Field `RawResponse` of struct `LedgerClientCreateResponse` has been removed
- Field `LedgerClientGetResult` of struct `LedgerClientGetResponse` has been removed
- Field `RawResponse` of struct `LedgerClientGetResponse` has been removed
- Field `LedgerClientUpdateResult` of struct `LedgerClientUpdateResponse` has been removed
- Field `RawResponse` of struct `LedgerClientUpdateResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `LedgerClientDeleteResponse` has been removed
- Field `LedgerClientListByResourceGroupResult` of struct `LedgerClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `LedgerClientListByResourceGroupResponse` has been removed
- Field `LedgerClientListBySubscriptionResult` of struct `LedgerClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `LedgerClientListBySubscriptionResponse` has been removed
- Field `ClientCheckNameAvailabilityResult` of struct `ClientCheckNameAvailabilityResponse` has been removed
- Field `RawResponse` of struct `ClientCheckNameAvailabilityResponse` has been removed

### Features Added

- New struct `ResourceLocation`
- New field `ResumeToken` in struct `LedgerClientBeginCreateOptions`
- New anonymous field `ConfidentialLedger` in struct `LedgerClientCreateResponse`
- New field `ResumeToken` in struct `LedgerClientBeginDeleteOptions`
- New field `ResumeToken` in struct `LedgerClientBeginUpdateOptions`
- New anonymous field `List` in struct `LedgerClientListBySubscriptionResponse`
- New anonymous field `CheckNameAvailabilityResponse` in struct `ClientCheckNameAvailabilityResponse`
- New anonymous field `ConfidentialLedger` in struct `LedgerClientGetResponse`
- New anonymous field `ResourceProviderOperationList` in struct `OperationsClientListResponse`
- New anonymous field `ConfidentialLedger` in struct `LedgerClientUpdateResponse`
- New anonymous field `List` in struct `LedgerClientListByResourceGroupResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*LedgerClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, ConfidentialLedger, *LedgerBeginUpdateOptions)` to `(context.Context, string, string, ConfidentialLedger, *LedgerClientBeginUpdateOptions)`
- Function `*LedgerClient.BeginUpdate` return value(s) have been changed from `(LedgerUpdatePollerResponse, error)` to `(LedgerClientUpdatePollerResponse, error)`
- Function `*LedgerClient.ListBySubscription` parameter(s) have been changed from `(*LedgerListBySubscriptionOptions)` to `(*LedgerClientListBySubscriptionOptions)`
- Function `*LedgerClient.ListBySubscription` return value(s) have been changed from `(*LedgerListBySubscriptionPager)` to `(*LedgerClientListBySubscriptionPager)`
- Function `*LedgerClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *LedgerBeginDeleteOptions)` to `(context.Context, string, string, *LedgerClientBeginDeleteOptions)`
- Function `*LedgerClient.BeginDelete` return value(s) have been changed from `(LedgerDeletePollerResponse, error)` to `(LedgerClientDeletePollerResponse, error)`
- Function `*LedgerClient.BeginCreate` parameter(s) have been changed from `(context.Context, string, string, ConfidentialLedger, *LedgerBeginCreateOptions)` to `(context.Context, string, string, ConfidentialLedger, *LedgerClientBeginCreateOptions)`
- Function `*LedgerClient.BeginCreate` return value(s) have been changed from `(LedgerCreatePollerResponse, error)` to `(LedgerClientCreatePollerResponse, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*LedgerClient.Get` parameter(s) have been changed from `(context.Context, string, string, *LedgerGetOptions)` to `(context.Context, string, string, *LedgerClientGetOptions)`
- Function `*LedgerClient.Get` return value(s) have been changed from `(LedgerGetResponse, error)` to `(LedgerClientGetResponse, error)`
- Function `*LedgerClient.ListByResourceGroup` parameter(s) have been changed from `(string, *LedgerListByResourceGroupOptions)` to `(string, *LedgerClientListByResourceGroupOptions)`
- Function `*LedgerClient.ListByResourceGroup` return value(s) have been changed from `(*LedgerListByResourceGroupPager)` to `(*LedgerClientListByResourceGroupPager)`
- Function `*LedgerDeletePollerResponse.Resume` has been removed
- Function `LedgerUpdatePollerResponse.PollUntilDone` has been removed
- Function `*LedgerUpdatePoller.Done` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `*LedgerUpdatePoller.Poll` has been removed
- Function `LedgerCreatePollerResponse.PollUntilDone` has been removed
- Function `*LedgerUpdatePollerResponse.Resume` has been removed
- Function `*ConfidentialLedgerClient.CheckNameAvailability` has been removed
- Function `*LedgerDeletePoller.FinalResponse` has been removed
- Function `*LedgerListByResourceGroupPager.Err` has been removed
- Function `*LedgerCreatePoller.ResumeToken` has been removed
- Function `*LedgerCreatePollerResponse.Resume` has been removed
- Function `*OperationsListPager.PageResponse` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Function `Location.MarshalJSON` has been removed
- Function `ConfidentialLedgerList.MarshalJSON` has been removed
- Function `*LedgerListBySubscriptionPager.Err` has been removed
- Function `*LedgerCreatePoller.Poll` has been removed
- Function `NewConfidentialLedgerClient` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `*LedgerUpdatePoller.FinalResponse` has been removed
- Function `*LedgerListBySubscriptionPager.PageResponse` has been removed
- Function `*LedgerCreatePoller.FinalResponse` has been removed
- Function `Resource.MarshalJSON` has been removed
- Function `*LedgerDeletePoller.Done` has been removed
- Function `*LedgerListByResourceGroupPager.PageResponse` has been removed
- Function `*LedgerListByResourceGroupPager.NextPage` has been removed
- Function `*LedgerCreatePoller.Done` has been removed
- Function `*LedgerUpdatePoller.ResumeToken` has been removed
- Function `*LedgerDeletePoller.ResumeToken` has been removed
- Function `*LedgerDeletePoller.Poll` has been removed
- Function `LedgerDeletePollerResponse.PollUntilDone` has been removed
- Function `*LedgerListBySubscriptionPager.NextPage` has been removed
- Struct `ConfidentialLedgerCheckNameAvailabilityOptions` has been removed
- Struct `ConfidentialLedgerCheckNameAvailabilityResponse` has been removed
- Struct `ConfidentialLedgerCheckNameAvailabilityResult` has been removed
- Struct `ConfidentialLedgerClient` has been removed
- Struct `ConfidentialLedgerList` has been removed
- Struct `LedgerBeginCreateOptions` has been removed
- Struct `LedgerBeginDeleteOptions` has been removed
- Struct `LedgerBeginUpdateOptions` has been removed
- Struct `LedgerCreatePoller` has been removed
- Struct `LedgerCreatePollerResponse` has been removed
- Struct `LedgerCreateResponse` has been removed
- Struct `LedgerCreateResult` has been removed
- Struct `LedgerDeletePoller` has been removed
- Struct `LedgerDeletePollerResponse` has been removed
- Struct `LedgerDeleteResponse` has been removed
- Struct `LedgerGetOptions` has been removed
- Struct `LedgerGetResponse` has been removed
- Struct `LedgerGetResult` has been removed
- Struct `LedgerListByResourceGroupOptions` has been removed
- Struct `LedgerListByResourceGroupPager` has been removed
- Struct `LedgerListByResourceGroupResponse` has been removed
- Struct `LedgerListByResourceGroupResult` has been removed
- Struct `LedgerListBySubscriptionOptions` has been removed
- Struct `LedgerListBySubscriptionPager` has been removed
- Struct `LedgerListBySubscriptionResponse` has been removed
- Struct `LedgerListBySubscriptionResult` has been removed
- Struct `LedgerUpdatePoller` has been removed
- Struct `LedgerUpdatePollerResponse` has been removed
- Struct `LedgerUpdateResponse` has been removed
- Struct `LedgerUpdateResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed
- Field `Location` of struct `ConfidentialLedger` has been removed
- Field `Resource` of struct `ConfidentialLedger` has been removed
- Field `Tags` of struct `ConfidentialLedger` has been removed

### Features Added

- New function `*LedgerClientListByResourceGroupPager.PageResponse() LedgerClientListByResourceGroupResponse`
- New function `*LedgerClientDeletePoller.FinalResponse(context.Context) (LedgerClientDeleteResponse, error)`
- New function `*LedgerClientDeletePoller.Done() bool`
- New function `*LedgerClientDeletePollerResponse.Resume(context.Context, *LedgerClient, string) error`
- New function `*OperationsClientListPager.Err() error`
- New function `*LedgerClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `LedgerClientUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (LedgerClientUpdateResponse, error)`
- New function `*LedgerClientDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `LedgerClientDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (LedgerClientDeleteResponse, error)`
- New function `*LedgerClientCreatePoller.ResumeToken() (string, error)`
- New function `*LedgerClientCreatePoller.FinalResponse(context.Context) (LedgerClientCreateResponse, error)`
- New function `*LedgerClientCreatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*LedgerClientDeletePoller.ResumeToken() (string, error)`
- New function `*LedgerClientListBySubscriptionPager.PageResponse() LedgerClientListBySubscriptionResponse`
- New function `*LedgerClientListBySubscriptionPager.NextPage(context.Context) bool`
- New function `*LedgerClientUpdatePoller.ResumeToken() (string, error)`
- New function `*LedgerClientListByResourceGroupPager.Err() error`
- New function `LedgerClientCreatePollerResponse.PollUntilDone(context.Context, time.Duration) (LedgerClientCreateResponse, error)`
- New function `*LedgerClientUpdatePollerResponse.Resume(context.Context, *LedgerClient, string) error`
- New function `*Client.CheckNameAvailability(context.Context, CheckNameAvailabilityRequest, *ClientCheckNameAvailabilityOptions) (ClientCheckNameAvailabilityResponse, error)`
- New function `*LedgerClientUpdatePoller.FinalResponse(context.Context) (LedgerClientUpdateResponse, error)`
- New function `*LedgerClientCreatePoller.Done() bool`
- New function `NewClient(string, azcore.TokenCredential, *arm.ClientOptions) *Client`
- New function `*LedgerClientListBySubscriptionPager.Err() error`
- New function `*LedgerClientUpdatePoller.Done() bool`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*LedgerClientUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `List.MarshalJSON() ([]byte, error)`
- New function `*LedgerClientCreatePollerResponse.Resume(context.Context, *LedgerClient, string) error`
- New struct `Client`
- New struct `ClientCheckNameAvailabilityOptions`
- New struct `ClientCheckNameAvailabilityResponse`
- New struct `ClientCheckNameAvailabilityResult`
- New struct `LedgerClientBeginCreateOptions`
- New struct `LedgerClientBeginDeleteOptions`
- New struct `LedgerClientBeginUpdateOptions`
- New struct `LedgerClientCreatePoller`
- New struct `LedgerClientCreatePollerResponse`
- New struct `LedgerClientCreateResponse`
- New struct `LedgerClientCreateResult`
- New struct `LedgerClientDeletePoller`
- New struct `LedgerClientDeletePollerResponse`
- New struct `LedgerClientDeleteResponse`
- New struct `LedgerClientGetOptions`
- New struct `LedgerClientGetResponse`
- New struct `LedgerClientGetResult`
- New struct `LedgerClientListByResourceGroupOptions`
- New struct `LedgerClientListByResourceGroupPager`
- New struct `LedgerClientListByResourceGroupResponse`
- New struct `LedgerClientListByResourceGroupResult`
- New struct `LedgerClientListBySubscriptionOptions`
- New struct `LedgerClientListBySubscriptionPager`
- New struct `LedgerClientListBySubscriptionResponse`
- New struct `LedgerClientListBySubscriptionResult`
- New struct `LedgerClientUpdatePoller`
- New struct `LedgerClientUpdatePollerResponse`
- New struct `LedgerClientUpdateResponse`
- New struct `LedgerClientUpdateResult`
- New struct `List`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New field `Error` in struct `ErrorResponse`
- New field `ID` in struct `ConfidentialLedger`
- New field `Name` in struct `ConfidentialLedger`
- New field `SystemData` in struct `ConfidentialLedger`
- New field `Type` in struct `ConfidentialLedger`
- New field `Location` in struct `ConfidentialLedger`
- New field `Tags` in struct `ConfidentialLedger`


## 0.1.0 (2021-12-01)

- Initial preview release.
