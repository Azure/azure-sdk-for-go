# Release History

## 0.6.0 (2022-05-17)
### Breaking Changes

- Function `*SolutionsClient.BeginCreateOrUpdate` return value(s) have been changed from `(*armruntime.Poller[SolutionsClientCreateOrUpdateResponse], error)` to `(*runtime.Poller[SolutionsClientCreateOrUpdateResponse], error)`
- Function `*SolutionsClient.BeginDelete` return value(s) have been changed from `(*armruntime.Poller[SolutionsClientDeleteResponse], error)` to `(*runtime.Poller[SolutionsClientDeleteResponse], error)`
- Function `*SolutionsClient.BeginUpdate` return value(s) have been changed from `(*armruntime.Poller[SolutionsClientUpdateResponse], error)` to `(*runtime.Poller[SolutionsClientUpdateResponse], error)`
- Function `SolutionPropertiesList.MarshalJSON` has been removed
- Function `ManagementAssociationPropertiesList.MarshalJSON` has been removed
- Function `OperationListResult.MarshalJSON` has been removed
- Function `ManagementConfigurationPropertiesList.MarshalJSON` has been removed


## 0.5.0 (2022-04-18)
### Breaking Changes

- Function `*OperationsClient.List` has been removed

### Features Added

- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`


## 0.4.0 (2022-04-12)
### Breaking Changes

- Function `*SolutionsClient.BeginCreateOrUpdate` return value(s) have been changed from `(SolutionsClientCreateOrUpdatePollerResponse, error)` to `(*armruntime.Poller[SolutionsClientCreateOrUpdateResponse], error)`
- Function `NewManagementAssociationsClient` parameter(s) have been changed from `(string, string, string, string, azcore.TokenCredential, *arm.ClientOptions)` to `(string, azcore.TokenCredential, *arm.ClientOptions)`
- Function `NewManagementAssociationsClient` return value(s) have been changed from `(*ManagementAssociationsClient)` to `(*ManagementAssociationsClient, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(context.Context, *OperationsClientListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(OperationsClientListResponse, error)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `*ManagementAssociationsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *ManagementAssociationsClientDeleteOptions)` to `(context.Context, string, string, string, string, string, *ManagementAssociationsClientDeleteOptions)`
- Function `*SolutionsClient.BeginUpdate` return value(s) have been changed from `(SolutionsClientUpdatePollerResponse, error)` to `(*armruntime.Poller[SolutionsClientUpdateResponse], error)`
- Function `*ManagementAssociationsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, ManagementAssociation, *ManagementAssociationsClientCreateOrUpdateOptions)` to `(context.Context, string, string, string, string, string, ManagementAssociation, *ManagementAssociationsClientCreateOrUpdateOptions)`
- Function `*ManagementAssociationsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *ManagementAssociationsClientGetOptions)` to `(context.Context, string, string, string, string, string, *ManagementAssociationsClientGetOptions)`
- Function `NewSolutionsClient` return value(s) have been changed from `(*SolutionsClient)` to `(*SolutionsClient, error)`
- Function `*SolutionsClient.BeginDelete` return value(s) have been changed from `(SolutionsClientDeletePollerResponse, error)` to `(*armruntime.Poller[SolutionsClientDeleteResponse], error)`
- Function `NewManagementConfigurationsClient` return value(s) have been changed from `(*ManagementConfigurationsClient)` to `(*ManagementConfigurationsClient, error)`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Type of `ManagementConfigurationProperties.Template` has been changed from `map[string]interface{}` to `interface{}`
- Function `*SolutionsClientUpdatePoller.ResumeToken` has been removed
- Function `*SolutionsClientUpdatePoller.FinalResponse` has been removed
- Function `*SolutionsClientCreateOrUpdatePoller.Poll` has been removed
- Function `*SolutionsClientCreateOrUpdatePoller.ResumeToken` has been removed
- Function `*SolutionsClientCreateOrUpdatePollerResponse.Resume` has been removed
- Function `*SolutionsClientDeletePollerResponse.Resume` has been removed
- Function `*SolutionsClientDeletePoller.ResumeToken` has been removed
- Function `*SolutionsClientDeletePoller.Done` has been removed
- Function `SolutionsClientDeletePollerResponse.PollUntilDone` has been removed
- Function `*SolutionsClientDeletePoller.FinalResponse` has been removed
- Function `SolutionsClientUpdatePollerResponse.PollUntilDone` has been removed
- Function `*SolutionsClientUpdatePoller.Done` has been removed
- Function `*SolutionsClientCreateOrUpdatePoller.FinalResponse` has been removed
- Function `SolutionsClientCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `*SolutionsClientCreateOrUpdatePoller.Done` has been removed
- Function `*SolutionsClientDeletePoller.Poll` has been removed
- Function `*SolutionsClientUpdatePoller.Poll` has been removed
- Function `*SolutionsClientUpdatePollerResponse.Resume` has been removed
- Struct `ManagementAssociationsClientCreateOrUpdateResult` has been removed
- Struct `ManagementAssociationsClientGetResult` has been removed
- Struct `ManagementAssociationsClientListBySubscriptionResult` has been removed
- Struct `ManagementConfigurationsClientCreateOrUpdateResult` has been removed
- Struct `ManagementConfigurationsClientGetResult` has been removed
- Struct `ManagementConfigurationsClientListBySubscriptionResult` has been removed
- Struct `OperationsClientListResult` has been removed
- Struct `SolutionsClientCreateOrUpdatePoller` has been removed
- Struct `SolutionsClientCreateOrUpdatePollerResponse` has been removed
- Struct `SolutionsClientCreateOrUpdateResult` has been removed
- Struct `SolutionsClientDeletePoller` has been removed
- Struct `SolutionsClientDeletePollerResponse` has been removed
- Struct `SolutionsClientGetResult` has been removed
- Struct `SolutionsClientListByResourceGroupResult` has been removed
- Struct `SolutionsClientListBySubscriptionResult` has been removed
- Struct `SolutionsClientUpdatePoller` has been removed
- Struct `SolutionsClientUpdatePollerResponse` has been removed
- Struct `SolutionsClientUpdateResult` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `ManagementConfigurationsClientCreateOrUpdateResult` of struct `ManagementConfigurationsClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `ManagementConfigurationsClientCreateOrUpdateResponse` has been removed
- Field `ManagementAssociationsClientListBySubscriptionResult` of struct `ManagementAssociationsClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `ManagementAssociationsClientListBySubscriptionResponse` has been removed
- Field `ManagementConfigurationsClientGetResult` of struct `ManagementConfigurationsClientGetResponse` has been removed
- Field `RawResponse` of struct `ManagementConfigurationsClientGetResponse` has been removed
- Field `ManagementAssociationsClientGetResult` of struct `ManagementAssociationsClientGetResponse` has been removed
- Field `RawResponse` of struct `ManagementAssociationsClientGetResponse` has been removed
- Field `RawResponse` of struct `ManagementConfigurationsClientDeleteResponse` has been removed
- Field `SolutionsClientCreateOrUpdateResult` of struct `SolutionsClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `SolutionsClientCreateOrUpdateResponse` has been removed
- Field `SolutionsClientListBySubscriptionResult` of struct `SolutionsClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `SolutionsClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `SolutionsClientDeleteResponse` has been removed
- Field `ManagementConfigurationsClientListBySubscriptionResult` of struct `ManagementConfigurationsClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `ManagementConfigurationsClientListBySubscriptionResponse` has been removed
- Field `ManagementAssociationsClientCreateOrUpdateResult` of struct `ManagementAssociationsClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `ManagementAssociationsClientCreateOrUpdateResponse` has been removed
- Field `SolutionsClientListByResourceGroupResult` of struct `SolutionsClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `SolutionsClientListByResourceGroupResponse` has been removed
- Field `SolutionsClientUpdateResult` of struct `SolutionsClientUpdateResponse` has been removed
- Field `RawResponse` of struct `SolutionsClientUpdateResponse` has been removed
- Field `RawResponse` of struct `ManagementAssociationsClientDeleteResponse` has been removed
- Field `SolutionsClientGetResult` of struct `SolutionsClientGetResponse` has been removed
- Field `RawResponse` of struct `SolutionsClientGetResponse` has been removed

### Features Added

- New field `ResumeToken` in struct `SolutionsClientBeginUpdateOptions`
- New anonymous field `ManagementConfiguration` in struct `ManagementConfigurationsClientCreateOrUpdateResponse`
- New anonymous field `ManagementAssociation` in struct `ManagementAssociationsClientCreateOrUpdateResponse`
- New field `ResumeToken` in struct `SolutionsClientBeginCreateOrUpdateOptions`
- New anonymous field `SolutionPropertiesList` in struct `SolutionsClientListByResourceGroupResponse`
- New field `ResumeToken` in struct `SolutionsClientBeginDeleteOptions`
- New anonymous field `SolutionPropertiesList` in struct `SolutionsClientListBySubscriptionResponse`
- New anonymous field `ManagementConfiguration` in struct `ManagementConfigurationsClientGetResponse`
- New anonymous field `Solution` in struct `SolutionsClientCreateOrUpdateResponse`
- New anonymous field `Solution` in struct `SolutionsClientGetResponse`
- New anonymous field `ManagementConfigurationPropertiesList` in struct `ManagementConfigurationsClientListBySubscriptionResponse`
- New anonymous field `Solution` in struct `SolutionsClientUpdateResponse`
- New anonymous field `ManagementAssociation` in struct `ManagementAssociationsClientGetResponse`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New anonymous field `ManagementAssociationPropertiesList` in struct `ManagementAssociationsClientListBySubscriptionResponse`


## 0.3.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.3.0 (2022-01-13)
### Breaking Changes

- Function `*ManagementConfigurationsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *ManagementConfigurationsGetOptions)` to `(context.Context, string, string, *ManagementConfigurationsClientGetOptions)`
- Function `*ManagementConfigurationsClient.Get` return value(s) have been changed from `(ManagementConfigurationsGetResponse, error)` to `(ManagementConfigurationsClientGetResponse, error)`
- Function `*ManagementAssociationsClient.ListBySubscription` parameter(s) have been changed from `(context.Context, *ManagementAssociationsListBySubscriptionOptions)` to `(context.Context, *ManagementAssociationsClientListBySubscriptionOptions)`
- Function `*ManagementAssociationsClient.ListBySubscription` return value(s) have been changed from `(ManagementAssociationsListBySubscriptionResponse, error)` to `(ManagementAssociationsClientListBySubscriptionResponse, error)`
- Function `*ManagementConfigurationsClient.ListBySubscription` parameter(s) have been changed from `(context.Context, *ManagementConfigurationsListBySubscriptionOptions)` to `(context.Context, *ManagementConfigurationsClientListBySubscriptionOptions)`
- Function `*ManagementConfigurationsClient.ListBySubscription` return value(s) have been changed from `(ManagementConfigurationsListBySubscriptionResponse, error)` to `(ManagementConfigurationsClientListBySubscriptionResponse, error)`
- Function `*SolutionsClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *SolutionsBeginDeleteOptions)` to `(context.Context, string, string, *SolutionsClientBeginDeleteOptions)`
- Function `*SolutionsClient.BeginDelete` return value(s) have been changed from `(SolutionsDeletePollerResponse, error)` to `(SolutionsClientDeletePollerResponse, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(context.Context, *OperationsListOptions)` to `(context.Context, *OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(OperationsListResponse, error)` to `(OperationsClientListResponse, error)`
- Function `*ManagementConfigurationsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *ManagementConfigurationsDeleteOptions)` to `(context.Context, string, string, *ManagementConfigurationsClientDeleteOptions)`
- Function `*ManagementConfigurationsClient.Delete` return value(s) have been changed from `(ManagementConfigurationsDeleteResponse, error)` to `(ManagementConfigurationsClientDeleteResponse, error)`
- Function `*ManagementAssociationsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *ManagementAssociationsDeleteOptions)` to `(context.Context, string, string, *ManagementAssociationsClientDeleteOptions)`
- Function `*ManagementAssociationsClient.Delete` return value(s) have been changed from `(ManagementAssociationsDeleteResponse, error)` to `(ManagementAssociationsClientDeleteResponse, error)`
- Function `*ManagementConfigurationsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, ManagementConfiguration, *ManagementConfigurationsCreateOrUpdateOptions)` to `(context.Context, string, string, ManagementConfiguration, *ManagementConfigurationsClientCreateOrUpdateOptions)`
- Function `*ManagementConfigurationsClient.CreateOrUpdate` return value(s) have been changed from `(ManagementConfigurationsCreateOrUpdateResponse, error)` to `(ManagementConfigurationsClientCreateOrUpdateResponse, error)`
- Function `*SolutionsClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, SolutionPatch, *SolutionsBeginUpdateOptions)` to `(context.Context, string, string, SolutionPatch, *SolutionsClientBeginUpdateOptions)`
- Function `*SolutionsClient.BeginUpdate` return value(s) have been changed from `(SolutionsUpdatePollerResponse, error)` to `(SolutionsClientUpdatePollerResponse, error)`
- Function `*SolutionsClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, Solution, *SolutionsBeginCreateOrUpdateOptions)` to `(context.Context, string, string, Solution, *SolutionsClientBeginCreateOrUpdateOptions)`
- Function `*SolutionsClient.BeginCreateOrUpdate` return value(s) have been changed from `(SolutionsCreateOrUpdatePollerResponse, error)` to `(SolutionsClientCreateOrUpdatePollerResponse, error)`
- Function `*SolutionsClient.ListByResourceGroup` parameter(s) have been changed from `(context.Context, string, *SolutionsListByResourceGroupOptions)` to `(context.Context, string, *SolutionsClientListByResourceGroupOptions)`
- Function `*SolutionsClient.ListByResourceGroup` return value(s) have been changed from `(SolutionsListByResourceGroupResponse, error)` to `(SolutionsClientListByResourceGroupResponse, error)`
- Function `*ManagementAssociationsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *ManagementAssociationsGetOptions)` to `(context.Context, string, string, *ManagementAssociationsClientGetOptions)`
- Function `*ManagementAssociationsClient.Get` return value(s) have been changed from `(ManagementAssociationsGetResponse, error)` to `(ManagementAssociationsClientGetResponse, error)`
- Function `*SolutionsClient.ListBySubscription` parameter(s) have been changed from `(context.Context, *SolutionsListBySubscriptionOptions)` to `(context.Context, *SolutionsClientListBySubscriptionOptions)`
- Function `*SolutionsClient.ListBySubscription` return value(s) have been changed from `(SolutionsListBySubscriptionResponse, error)` to `(SolutionsClientListBySubscriptionResponse, error)`
- Function `*ManagementAssociationsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, ManagementAssociation, *ManagementAssociationsCreateOrUpdateOptions)` to `(context.Context, string, string, ManagementAssociation, *ManagementAssociationsClientCreateOrUpdateOptions)`
- Function `*ManagementAssociationsClient.CreateOrUpdate` return value(s) have been changed from `(ManagementAssociationsCreateOrUpdateResponse, error)` to `(ManagementAssociationsClientCreateOrUpdateResponse, error)`
- Function `*SolutionsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *SolutionsGetOptions)` to `(context.Context, string, string, *SolutionsClientGetOptions)`
- Function `*SolutionsClient.Get` return value(s) have been changed from `(SolutionsGetResponse, error)` to `(SolutionsClientGetResponse, error)`
- Function `*SolutionsCreateOrUpdatePollerResponse.Resume` has been removed
- Function `*SolutionsCreateOrUpdatePoller.Poll` has been removed
- Function `*SolutionsUpdatePoller.ResumeToken` has been removed
- Function `*SolutionsDeletePoller.Done` has been removed
- Function `*SolutionsDeletePoller.FinalResponse` has been removed
- Function `*SolutionsCreateOrUpdatePoller.Done` has been removed
- Function `*SolutionsDeletePoller.Poll` has been removed
- Function `*SolutionsCreateOrUpdatePoller.FinalResponse` has been removed
- Function `*SolutionsDeletePollerResponse.Resume` has been removed
- Function `*SolutionsCreateOrUpdatePoller.ResumeToken` has been removed
- Function `SolutionsCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `*SolutionsUpdatePollerResponse.Resume` has been removed
- Function `*SolutionsUpdatePoller.Done` has been removed
- Function `SolutionsDeletePollerResponse.PollUntilDone` has been removed
- Function `*SolutionsDeletePoller.ResumeToken` has been removed
- Function `SolutionsUpdatePollerResponse.PollUntilDone` has been removed
- Function `CodeMessageError.Error` has been removed
- Function `*SolutionsUpdatePoller.Poll` has been removed
- Function `*SolutionsUpdatePoller.FinalResponse` has been removed
- Struct `ManagementAssociationsCreateOrUpdateOptions` has been removed
- Struct `ManagementAssociationsCreateOrUpdateResponse` has been removed
- Struct `ManagementAssociationsCreateOrUpdateResult` has been removed
- Struct `ManagementAssociationsDeleteOptions` has been removed
- Struct `ManagementAssociationsDeleteResponse` has been removed
- Struct `ManagementAssociationsGetOptions` has been removed
- Struct `ManagementAssociationsGetResponse` has been removed
- Struct `ManagementAssociationsGetResult` has been removed
- Struct `ManagementAssociationsListBySubscriptionOptions` has been removed
- Struct `ManagementAssociationsListBySubscriptionResponse` has been removed
- Struct `ManagementAssociationsListBySubscriptionResult` has been removed
- Struct `ManagementConfigurationsCreateOrUpdateOptions` has been removed
- Struct `ManagementConfigurationsCreateOrUpdateResponse` has been removed
- Struct `ManagementConfigurationsCreateOrUpdateResult` has been removed
- Struct `ManagementConfigurationsDeleteOptions` has been removed
- Struct `ManagementConfigurationsDeleteResponse` has been removed
- Struct `ManagementConfigurationsGetOptions` has been removed
- Struct `ManagementConfigurationsGetResponse` has been removed
- Struct `ManagementConfigurationsGetResult` has been removed
- Struct `ManagementConfigurationsListBySubscriptionOptions` has been removed
- Struct `ManagementConfigurationsListBySubscriptionResponse` has been removed
- Struct `ManagementConfigurationsListBySubscriptionResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Struct `SolutionsBeginCreateOrUpdateOptions` has been removed
- Struct `SolutionsBeginDeleteOptions` has been removed
- Struct `SolutionsBeginUpdateOptions` has been removed
- Struct `SolutionsCreateOrUpdatePoller` has been removed
- Struct `SolutionsCreateOrUpdatePollerResponse` has been removed
- Struct `SolutionsCreateOrUpdateResponse` has been removed
- Struct `SolutionsCreateOrUpdateResult` has been removed
- Struct `SolutionsDeletePoller` has been removed
- Struct `SolutionsDeletePollerResponse` has been removed
- Struct `SolutionsDeleteResponse` has been removed
- Struct `SolutionsGetOptions` has been removed
- Struct `SolutionsGetResponse` has been removed
- Struct `SolutionsGetResult` has been removed
- Struct `SolutionsListByResourceGroupOptions` has been removed
- Struct `SolutionsListByResourceGroupResponse` has been removed
- Struct `SolutionsListByResourceGroupResult` has been removed
- Struct `SolutionsListBySubscriptionOptions` has been removed
- Struct `SolutionsListBySubscriptionResponse` has been removed
- Struct `SolutionsListBySubscriptionResult` has been removed
- Struct `SolutionsUpdatePoller` has been removed
- Struct `SolutionsUpdatePollerResponse` has been removed
- Struct `SolutionsUpdateResponse` has been removed
- Struct `SolutionsUpdateResult` has been removed
- Field `InnerError` of struct `CodeMessageError` has been removed

### Features Added

- New function `SolutionsClientDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (SolutionsClientDeleteResponse, error)`
- New function `*SolutionsClientDeletePoller.FinalResponse(context.Context) (SolutionsClientDeleteResponse, error)`
- New function `*SolutionsClientUpdatePoller.FinalResponse(context.Context) (SolutionsClientUpdateResponse, error)`
- New function `*SolutionsClientUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*SolutionsClientCreateOrUpdatePollerResponse.Resume(context.Context, *SolutionsClient, string) error`
- New function `SolutionsClientCreateOrUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (SolutionsClientCreateOrUpdateResponse, error)`
- New function `*SolutionsClientDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `*SolutionsClientUpdatePoller.Done() bool`
- New function `*SolutionsClientCreateOrUpdatePoller.ResumeToken() (string, error)`
- New function `*SolutionsClientCreateOrUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*SolutionsClientCreateOrUpdatePoller.Done() bool`
- New function `*SolutionsClientDeletePoller.ResumeToken() (string, error)`
- New function `*SolutionsClientUpdatePollerResponse.Resume(context.Context, *SolutionsClient, string) error`
- New function `*SolutionsClientDeletePollerResponse.Resume(context.Context, *SolutionsClient, string) error`
- New function `SolutionsClientUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (SolutionsClientUpdateResponse, error)`
- New function `*SolutionsClientUpdatePoller.ResumeToken() (string, error)`
- New function `*SolutionsClientDeletePoller.Done() bool`
- New function `*SolutionsClientCreateOrUpdatePoller.FinalResponse(context.Context) (SolutionsClientCreateOrUpdateResponse, error)`
- New struct `ManagementAssociationsClientCreateOrUpdateOptions`
- New struct `ManagementAssociationsClientCreateOrUpdateResponse`
- New struct `ManagementAssociationsClientCreateOrUpdateResult`
- New struct `ManagementAssociationsClientDeleteOptions`
- New struct `ManagementAssociationsClientDeleteResponse`
- New struct `ManagementAssociationsClientGetOptions`
- New struct `ManagementAssociationsClientGetResponse`
- New struct `ManagementAssociationsClientGetResult`
- New struct `ManagementAssociationsClientListBySubscriptionOptions`
- New struct `ManagementAssociationsClientListBySubscriptionResponse`
- New struct `ManagementAssociationsClientListBySubscriptionResult`
- New struct `ManagementConfigurationsClientCreateOrUpdateOptions`
- New struct `ManagementConfigurationsClientCreateOrUpdateResponse`
- New struct `ManagementConfigurationsClientCreateOrUpdateResult`
- New struct `ManagementConfigurationsClientDeleteOptions`
- New struct `ManagementConfigurationsClientDeleteResponse`
- New struct `ManagementConfigurationsClientGetOptions`
- New struct `ManagementConfigurationsClientGetResponse`
- New struct `ManagementConfigurationsClientGetResult`
- New struct `ManagementConfigurationsClientListBySubscriptionOptions`
- New struct `ManagementConfigurationsClientListBySubscriptionResponse`
- New struct `ManagementConfigurationsClientListBySubscriptionResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New struct `SolutionsClientBeginCreateOrUpdateOptions`
- New struct `SolutionsClientBeginDeleteOptions`
- New struct `SolutionsClientBeginUpdateOptions`
- New struct `SolutionsClientCreateOrUpdatePoller`
- New struct `SolutionsClientCreateOrUpdatePollerResponse`
- New struct `SolutionsClientCreateOrUpdateResponse`
- New struct `SolutionsClientCreateOrUpdateResult`
- New struct `SolutionsClientDeletePoller`
- New struct `SolutionsClientDeletePollerResponse`
- New struct `SolutionsClientDeleteResponse`
- New struct `SolutionsClientGetOptions`
- New struct `SolutionsClientGetResponse`
- New struct `SolutionsClientGetResult`
- New struct `SolutionsClientListByResourceGroupOptions`
- New struct `SolutionsClientListByResourceGroupResponse`
- New struct `SolutionsClientListByResourceGroupResult`
- New struct `SolutionsClientListBySubscriptionOptions`
- New struct `SolutionsClientListBySubscriptionResponse`
- New struct `SolutionsClientListBySubscriptionResult`
- New struct `SolutionsClientUpdatePoller`
- New struct `SolutionsClientUpdatePollerResponse`
- New struct `SolutionsClientUpdateResponse`
- New struct `SolutionsClientUpdateResult`
- New field `Error` in struct `CodeMessageError`


## 0.2.0 (2021-10-29)

### Breaking Changes

- `arm.Connection` has been removed in `github.com/Azure/azure-sdk-for-go/sdk/azcore/v0.20.0`
- The parameters of `NewXXXClient` has been changed from `(con *arm.Connection, subscriptionID string)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`

## 0.1.0 (2021-10-20)

- Initial preview release.
