# Release History

## 0.3.0 (2022-04-11)
### Breaking Changes

- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `NewPipelinesClient` return value(s) have been changed from `(*PipelinesClient)` to `(*PipelinesClient, error)`
- Function `*PipelinesClient.ListByResourceGroup` return value(s) have been changed from `(*PipelinesClientListByResourceGroupPager)` to `(*runtime.Pager[PipelinesClientListByResourceGroupResponse])`
- Function `*PipelinesClient.ListBySubscription` return value(s) have been changed from `(*PipelinesClientListBySubscriptionPager)` to `(*runtime.Pager[PipelinesClientListBySubscriptionResponse])`
- Function `*PipelineTemplateDefinitionsClient.List` return value(s) have been changed from `(*PipelineTemplateDefinitionsClientListPager)` to `(*runtime.Pager[PipelineTemplateDefinitionsClientListResponse])`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `*PipelinesClient.BeginCreateOrUpdate` return value(s) have been changed from `(PipelinesClientCreateOrUpdatePollerResponse, error)` to `(*armruntime.Poller[PipelinesClientCreateOrUpdateResponse], error)`
- Function `NewPipelineTemplateDefinitionsClient` return value(s) have been changed from `(*PipelineTemplateDefinitionsClient)` to `(*PipelineTemplateDefinitionsClient, error)`
- Function `*OperationsClientListPager.PageResponse` has been removed
- Function `*PipelineTemplateDefinitionsClientListPager.NextPage` has been removed
- Function `*PipelinesClientListBySubscriptionPager.NextPage` has been removed
- Function `*PipelinesClientCreateOrUpdatePoller.Done` has been removed
- Function `*PipelinesClientCreateOrUpdatePollerResponse.Resume` has been removed
- Function `*PipelinesClientListBySubscriptionPager.Err` has been removed
- Function `*PipelinesClientListByResourceGroupPager.NextPage` has been removed
- Function `InputDataType.ToPtr` has been removed
- Function `*OperationsClientListPager.Err` has been removed
- Function `CodeRepositoryType.ToPtr` has been removed
- Function `*PipelinesClientCreateOrUpdatePoller.ResumeToken` has been removed
- Function `PipelinesClientCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `*PipelinesClientListByResourceGroupPager.PageResponse` has been removed
- Function `AuthorizationType.ToPtr` has been removed
- Function `*PipelineTemplateDefinitionsClientListPager.Err` has been removed
- Function `*PipelinesClientListByResourceGroupPager.Err` has been removed
- Function `*PipelinesClientCreateOrUpdatePoller.Poll` has been removed
- Function `*PipelinesClientCreateOrUpdatePoller.FinalResponse` has been removed
- Function `*PipelinesClientListBySubscriptionPager.PageResponse` has been removed
- Function `*PipelineTemplateDefinitionsClientListPager.PageResponse` has been removed
- Function `*OperationsClientListPager.NextPage` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Struct `PipelineTemplateDefinitionsClientListPager` has been removed
- Struct `PipelineTemplateDefinitionsClientListResult` has been removed
- Struct `PipelinesClientCreateOrUpdatePoller` has been removed
- Struct `PipelinesClientCreateOrUpdatePollerResponse` has been removed
- Struct `PipelinesClientCreateOrUpdateResult` has been removed
- Struct `PipelinesClientGetResult` has been removed
- Struct `PipelinesClientListByResourceGroupPager` has been removed
- Struct `PipelinesClientListByResourceGroupResult` has been removed
- Struct `PipelinesClientListBySubscriptionPager` has been removed
- Struct `PipelinesClientListBySubscriptionResult` has been removed
- Struct `PipelinesClientUpdateResult` has been removed
- Field `PipelinesClientGetResult` of struct `PipelinesClientGetResponse` has been removed
- Field `RawResponse` of struct `PipelinesClientGetResponse` has been removed
- Field `PipelinesClientListBySubscriptionResult` of struct `PipelinesClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `PipelinesClientListBySubscriptionResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `PipelinesClientListByResourceGroupResult` of struct `PipelinesClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `PipelinesClientListByResourceGroupResponse` has been removed
- Field `PipelinesClientCreateOrUpdateResult` of struct `PipelinesClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `PipelinesClientCreateOrUpdateResponse` has been removed
- Field `PipelinesClientUpdateResult` of struct `PipelinesClientUpdateResponse` has been removed
- Field `RawResponse` of struct `PipelinesClientUpdateResponse` has been removed
- Field `PipelineTemplateDefinitionsClientListResult` of struct `PipelineTemplateDefinitionsClientListResponse` has been removed
- Field `RawResponse` of struct `PipelineTemplateDefinitionsClientListResponse` has been removed
- Field `RawResponse` of struct `PipelinesClientDeleteResponse` has been removed

### Features Added

- New anonymous field `Pipeline` in struct `PipelinesClientUpdateResponse`
- New anonymous field `PipelineTemplateDefinitionListResult` in struct `PipelineTemplateDefinitionsClientListResponse`
- New anonymous field `Pipeline` in struct `PipelinesClientCreateOrUpdateResponse`
- New anonymous field `PipelineListResult` in struct `PipelinesClientListBySubscriptionResponse`
- New anonymous field `Pipeline` in struct `PipelinesClientGetResponse`
- New field `ResumeToken` in struct `PipelinesClientBeginCreateOrUpdateOptions`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New anonymous field `PipelineListResult` in struct `PipelinesClientListByResourceGroupResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*PipelineTemplateDefinitionsClient.List` parameter(s) have been changed from `(*PipelineTemplateDefinitionsListOptions)` to `(*PipelineTemplateDefinitionsClientListOptions)`
- Function `*PipelineTemplateDefinitionsClient.List` return value(s) have been changed from `(*PipelineTemplateDefinitionsListPager)` to `(*PipelineTemplateDefinitionsClientListPager)`
- Function `*PipelinesClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *PipelinesDeleteOptions)` to `(context.Context, string, string, *PipelinesClientDeleteOptions)`
- Function `*PipelinesClient.Delete` return value(s) have been changed from `(PipelinesDeleteResponse, error)` to `(PipelinesClientDeleteResponse, error)`
- Function `*PipelinesClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, Pipeline, *PipelinesBeginCreateOrUpdateOptions)` to `(context.Context, string, string, Pipeline, *PipelinesClientBeginCreateOrUpdateOptions)`
- Function `*PipelinesClient.BeginCreateOrUpdate` return value(s) have been changed from `(PipelinesCreateOrUpdatePollerResponse, error)` to `(PipelinesClientCreateOrUpdatePollerResponse, error)`
- Function `*PipelinesClient.Get` parameter(s) have been changed from `(context.Context, string, string, *PipelinesGetOptions)` to `(context.Context, string, string, *PipelinesClientGetOptions)`
- Function `*PipelinesClient.Get` return value(s) have been changed from `(PipelinesGetResponse, error)` to `(PipelinesClientGetResponse, error)`
- Function `*PipelinesClient.ListBySubscription` parameter(s) have been changed from `(*PipelinesListBySubscriptionOptions)` to `(*PipelinesClientListBySubscriptionOptions)`
- Function `*PipelinesClient.ListBySubscription` return value(s) have been changed from `(*PipelinesListBySubscriptionPager)` to `(*PipelinesClientListBySubscriptionPager)`
- Function `*PipelinesClient.ListByResourceGroup` parameter(s) have been changed from `(string, *PipelinesListByResourceGroupOptions)` to `(string, *PipelinesClientListByResourceGroupOptions)`
- Function `*PipelinesClient.ListByResourceGroup` return value(s) have been changed from `(*PipelinesListByResourceGroupPager)` to `(*PipelinesClientListByResourceGroupPager)`
- Function `*PipelinesClient.Update` parameter(s) have been changed from `(context.Context, string, string, PipelineUpdateParameters, *PipelinesUpdateOptions)` to `(context.Context, string, string, PipelineUpdateParameters, *PipelinesClientUpdateOptions)`
- Function `*PipelinesClient.Update` return value(s) have been changed from `(PipelinesUpdateResponse, error)` to `(PipelinesClientUpdateResponse, error)`
- Function `*PipelinesCreateOrUpdatePoller.Poll` has been removed
- Function `*PipelinesCreateOrUpdatePoller.ResumeToken` has been removed
- Function `*PipelinesCreateOrUpdatePoller.FinalResponse` has been removed
- Function `*PipelinesCreateOrUpdatePoller.Done` has been removed
- Function `*PipelineTemplateDefinitionsListPager.NextPage` has been removed
- Function `*PipelinesListBySubscriptionPager.PageResponse` has been removed
- Function `*PipelineTemplateDefinitionsListPager.Err` has been removed
- Function `*PipelineTemplateDefinitionsListPager.PageResponse` has been removed
- Function `*PipelinesListByResourceGroupPager.PageResponse` has been removed
- Function `*PipelinesCreateOrUpdatePollerResponse.Resume` has been removed
- Function `*PipelinesListBySubscriptionPager.Err` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `PipelinesCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `*OperationsListPager.PageResponse` has been removed
- Function `*PipelinesListByResourceGroupPager.NextPage` has been removed
- Function `*PipelinesListBySubscriptionPager.NextPage` has been removed
- Function `CloudError.Error` has been removed
- Function `*PipelinesListByResourceGroupPager.Err` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Struct `PipelineTemplateDefinitionsListOptions` has been removed
- Struct `PipelineTemplateDefinitionsListPager` has been removed
- Struct `PipelineTemplateDefinitionsListResponse` has been removed
- Struct `PipelineTemplateDefinitionsListResult` has been removed
- Struct `PipelinesBeginCreateOrUpdateOptions` has been removed
- Struct `PipelinesCreateOrUpdatePoller` has been removed
- Struct `PipelinesCreateOrUpdatePollerResponse` has been removed
- Struct `PipelinesCreateOrUpdateResponse` has been removed
- Struct `PipelinesCreateOrUpdateResult` has been removed
- Struct `PipelinesDeleteOptions` has been removed
- Struct `PipelinesDeleteResponse` has been removed
- Struct `PipelinesGetOptions` has been removed
- Struct `PipelinesGetResponse` has been removed
- Struct `PipelinesGetResult` has been removed
- Struct `PipelinesListByResourceGroupOptions` has been removed
- Struct `PipelinesListByResourceGroupPager` has been removed
- Struct `PipelinesListByResourceGroupResponse` has been removed
- Struct `PipelinesListByResourceGroupResult` has been removed
- Struct `PipelinesListBySubscriptionOptions` has been removed
- Struct `PipelinesListBySubscriptionPager` has been removed
- Struct `PipelinesListBySubscriptionResponse` has been removed
- Struct `PipelinesListBySubscriptionResult` has been removed
- Struct `PipelinesUpdateOptions` has been removed
- Struct `PipelinesUpdateResponse` has been removed
- Struct `PipelinesUpdateResult` has been removed
- Field `InnerError` of struct `CloudError` has been removed
- Field `Resource` of struct `Pipeline` has been removed

### Features Added

- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*PipelinesClientCreateOrUpdatePoller.FinalResponse(context.Context) (PipelinesClientCreateOrUpdateResponse, error)`
- New function `*PipelinesClientListByResourceGroupPager.Err() error`
- New function `*OperationsClientListPager.Err() error`
- New function `*PipelinesClientCreateOrUpdatePollerResponse.Resume(context.Context, *PipelinesClient, string) error`
- New function `*PipelinesClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*PipelineTemplateDefinitionsClientListPager.PageResponse() PipelineTemplateDefinitionsClientListResponse`
- New function `*PipelinesClientCreateOrUpdatePoller.ResumeToken() (string, error)`
- New function `*PipelineTemplateDefinitionsClientListPager.NextPage(context.Context) bool`
- New function `*PipelinesClientListBySubscriptionPager.Err() error`
- New function `*PipelineTemplateDefinitionsClientListPager.Err() error`
- New function `*PipelinesClientListByResourceGroupPager.PageResponse() PipelinesClientListByResourceGroupResponse`
- New function `PipelinesClientCreateOrUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (PipelinesClientCreateOrUpdateResponse, error)`
- New function `*PipelinesClientCreateOrUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*PipelinesClientCreateOrUpdatePoller.Done() bool`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*PipelinesClientListBySubscriptionPager.PageResponse() PipelinesClientListBySubscriptionResponse`
- New function `*PipelinesClientListBySubscriptionPager.NextPage(context.Context) bool`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New struct `PipelineTemplateDefinitionsClientListOptions`
- New struct `PipelineTemplateDefinitionsClientListPager`
- New struct `PipelineTemplateDefinitionsClientListResponse`
- New struct `PipelineTemplateDefinitionsClientListResult`
- New struct `PipelinesClientBeginCreateOrUpdateOptions`
- New struct `PipelinesClientCreateOrUpdatePoller`
- New struct `PipelinesClientCreateOrUpdatePollerResponse`
- New struct `PipelinesClientCreateOrUpdateResponse`
- New struct `PipelinesClientCreateOrUpdateResult`
- New struct `PipelinesClientDeleteOptions`
- New struct `PipelinesClientDeleteResponse`
- New struct `PipelinesClientGetOptions`
- New struct `PipelinesClientGetResponse`
- New struct `PipelinesClientGetResult`
- New struct `PipelinesClientListByResourceGroupOptions`
- New struct `PipelinesClientListByResourceGroupPager`
- New struct `PipelinesClientListByResourceGroupResponse`
- New struct `PipelinesClientListByResourceGroupResult`
- New struct `PipelinesClientListBySubscriptionOptions`
- New struct `PipelinesClientListBySubscriptionPager`
- New struct `PipelinesClientListBySubscriptionResponse`
- New struct `PipelinesClientListBySubscriptionResult`
- New struct `PipelinesClientUpdateOptions`
- New struct `PipelinesClientUpdateResponse`
- New struct `PipelinesClientUpdateResult`
- New field `Error` in struct `CloudError`
- New field `Location` in struct `Pipeline`
- New field `Tags` in struct `Pipeline`
- New field `ID` in struct `Pipeline`
- New field `Name` in struct `Pipeline`
- New field `Type` in struct `Pipeline`


## 0.1.0 (2021-12-07)

- Initial preview release.
