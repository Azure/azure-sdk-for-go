# Release History

## 0.4.0 (2022-04-18)
### Breaking Changes

- Function `*WorkspacesClient.ListBySubscription` has been removed
- Function `*OfferingsClient.List` has been removed
- Function `*OperationsClient.List` has been removed
- Function `*WorkspacesClient.ListByResourceGroup` has been removed

### Features Added

- New function `*OfferingsClient.NewListPager(string, *OfferingsClientListOptions) *runtime.Pager[OfferingsClientListResponse]`
- New function `*WorkspacesClient.NewListBySubscriptionPager(*WorkspacesClientListBySubscriptionOptions) *runtime.Pager[WorkspacesClientListBySubscriptionResponse]`
- New function `*WorkspacesClient.NewListByResourceGroupPager(string, *WorkspacesClientListByResourceGroupOptions) *runtime.Pager[WorkspacesClientListByResourceGroupResponse]`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`


## 0.3.0 (2022-04-12)
### Breaking Changes

- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `*WorkspacesClient.BeginCreateOrUpdate` return value(s) have been changed from `(WorkspacesClientCreateOrUpdatePollerResponse, error)` to `(*armruntime.Poller[WorkspacesClientCreateOrUpdateResponse], error)`
- Function `*OfferingsClient.List` return value(s) have been changed from `(*OfferingsClientListPager)` to `(*runtime.Pager[OfferingsClientListResponse])`
- Function `*WorkspacesClient.BeginDelete` return value(s) have been changed from `(WorkspacesClientDeletePollerResponse, error)` to `(*armruntime.Poller[WorkspacesClientDeleteResponse], error)`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `NewWorkspaceClient` return value(s) have been changed from `(*WorkspaceClient)` to `(*WorkspaceClient, error)`
- Function `*WorkspacesClient.ListBySubscription` return value(s) have been changed from `(*WorkspacesClientListBySubscriptionPager)` to `(*runtime.Pager[WorkspacesClientListBySubscriptionResponse])`
- Function `NewOfferingsClient` return value(s) have been changed from `(*OfferingsClient)` to `(*OfferingsClient, error)`
- Function `NewWorkspacesClient` return value(s) have been changed from `(*WorkspacesClient)` to `(*WorkspacesClient, error)`
- Function `*WorkspacesClient.ListByResourceGroup` return value(s) have been changed from `(*WorkspacesClientListByResourceGroupPager)` to `(*runtime.Pager[WorkspacesClientListByResourceGroupResponse])`
- Type of `ErrorAdditionalInfo.Info` has been changed from `map[string]interface{}` to `interface{}`
- Function `*WorkspacesClientListBySubscriptionPager.Err` has been removed
- Function `*WorkspacesClientDeletePoller.Done` has been removed
- Function `ResourceIdentityType.ToPtr` has been removed
- Function `*WorkspacesClientListBySubscriptionPager.PageResponse` has been removed
- Function `*OperationsClientListPager.Err` has been removed
- Function `*WorkspacesClientCreateOrUpdatePollerResponse.Resume` has been removed
- Function `WorkspacesClientDeletePollerResponse.PollUntilDone` has been removed
- Function `Status.ToPtr` has been removed
- Function `*WorkspacesClientListBySubscriptionPager.NextPage` has been removed
- Function `*OfferingsClientListPager.PageResponse` has been removed
- Function `ProvisioningStatus.ToPtr` has been removed
- Function `*WorkspacesClientListByResourceGroupPager.NextPage` has been removed
- Function `*WorkspacesClientDeletePoller.Poll` has been removed
- Function `WorkspacesClientCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `UsableStatus.ToPtr` has been removed
- Function `*WorkspacesClientDeletePoller.ResumeToken` has been removed
- Function `*WorkspacesClientListByResourceGroupPager.Err` has been removed
- Function `*WorkspacesClientDeletePoller.FinalResponse` has been removed
- Function `*WorkspacesClientCreateOrUpdatePoller.ResumeToken` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Function `*WorkspacesClientCreateOrUpdatePoller.FinalResponse` has been removed
- Function `*OfferingsClientListPager.Err` has been removed
- Function `*WorkspacesClientDeletePollerResponse.Resume` has been removed
- Function `*WorkspacesClientListByResourceGroupPager.PageResponse` has been removed
- Function `*WorkspacesClientCreateOrUpdatePoller.Done` has been removed
- Function `*OperationsClientListPager.NextPage` has been removed
- Function `*WorkspacesClientCreateOrUpdatePoller.Poll` has been removed
- Function `CreatedByType.ToPtr` has been removed
- Function `*OfferingsClientListPager.NextPage` has been removed
- Struct `OfferingsClientListPager` has been removed
- Struct `OfferingsClientListResult` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Struct `WorkspaceClientCheckNameAvailabilityResult` has been removed
- Struct `WorkspacesClientCreateOrUpdatePoller` has been removed
- Struct `WorkspacesClientCreateOrUpdatePollerResponse` has been removed
- Struct `WorkspacesClientCreateOrUpdateResult` has been removed
- Struct `WorkspacesClientDeletePoller` has been removed
- Struct `WorkspacesClientDeletePollerResponse` has been removed
- Struct `WorkspacesClientGetResult` has been removed
- Struct `WorkspacesClientListByResourceGroupPager` has been removed
- Struct `WorkspacesClientListByResourceGroupResult` has been removed
- Struct `WorkspacesClientListBySubscriptionPager` has been removed
- Struct `WorkspacesClientListBySubscriptionResult` has been removed
- Struct `WorkspacesClientUpdateTagsResult` has been removed
- Field `WorkspacesClientUpdateTagsResult` of struct `WorkspacesClientUpdateTagsResponse` has been removed
- Field `RawResponse` of struct `WorkspacesClientUpdateTagsResponse` has been removed
- Field `WorkspacesClientListByResourceGroupResult` of struct `WorkspacesClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `WorkspacesClientListByResourceGroupResponse` has been removed
- Field `WorkspacesClientGetResult` of struct `WorkspacesClientGetResponse` has been removed
- Field `RawResponse` of struct `WorkspacesClientGetResponse` has been removed
- Field `OfferingsClientListResult` of struct `OfferingsClientListResponse` has been removed
- Field `RawResponse` of struct `OfferingsClientListResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `WorkspacesClientListBySubscriptionResult` of struct `WorkspacesClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `WorkspacesClientListBySubscriptionResponse` has been removed
- Field `WorkspacesClientCreateOrUpdateResult` of struct `WorkspacesClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `WorkspacesClientCreateOrUpdateResponse` has been removed
- Field `WorkspaceClientCheckNameAvailabilityResult` of struct `WorkspaceClientCheckNameAvailabilityResponse` has been removed
- Field `RawResponse` of struct `WorkspaceClientCheckNameAvailabilityResponse` has been removed
- Field `RawResponse` of struct `WorkspacesClientDeleteResponse` has been removed

### Features Added

- New anonymous field `Workspace` in struct `WorkspacesClientGetResponse`
- New field `ResumeToken` in struct `WorkspacesClientBeginDeleteOptions`
- New anonymous field `WorkspaceListResult` in struct `WorkspacesClientListByResourceGroupResponse`
- New anonymous field `OperationsList` in struct `OperationsClientListResponse`
- New anonymous field `CheckNameAvailabilityResult` in struct `WorkspaceClientCheckNameAvailabilityResponse`
- New anonymous field `WorkspaceListResult` in struct `WorkspacesClientListBySubscriptionResponse`
- New anonymous field `Workspace` in struct `WorkspacesClientUpdateTagsResponse`
- New field `AutoAdd` in struct `SKUDescription`
- New field `ResumeToken` in struct `WorkspacesClientBeginCreateOrUpdateOptions`
- New anonymous field `Workspace` in struct `WorkspacesClientCreateOrUpdateResponse`
- New anonymous field `OfferingsListResult` in struct `OfferingsClientListResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*WorkspacesClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *WorkspacesBeginDeleteOptions)` to `(context.Context, string, string, *WorkspacesClientBeginDeleteOptions)`
- Function `*WorkspacesClient.BeginDelete` return value(s) have been changed from `(WorkspacesDeletePollerResponse, error)` to `(WorkspacesClientDeletePollerResponse, error)`
- Function `*WorkspacesClient.ListByResourceGroup` parameter(s) have been changed from `(string, *WorkspacesListByResourceGroupOptions)` to `(string, *WorkspacesClientListByResourceGroupOptions)`
- Function `*WorkspacesClient.ListByResourceGroup` return value(s) have been changed from `(*WorkspacesListByResourceGroupPager)` to `(*WorkspacesClientListByResourceGroupPager)`
- Function `*WorkspacesClient.UpdateTags` parameter(s) have been changed from `(context.Context, string, string, TagsObject, *WorkspacesUpdateTagsOptions)` to `(context.Context, string, string, TagsObject, *WorkspacesClientUpdateTagsOptions)`
- Function `*WorkspacesClient.UpdateTags` return value(s) have been changed from `(WorkspacesUpdateTagsResponse, error)` to `(WorkspacesClientUpdateTagsResponse, error)`
- Function `*WorkspaceClient.CheckNameAvailability` parameter(s) have been changed from `(context.Context, string, CheckNameAvailabilityParameters, *WorkspaceCheckNameAvailabilityOptions)` to `(context.Context, string, CheckNameAvailabilityParameters, *WorkspaceClientCheckNameAvailabilityOptions)`
- Function `*WorkspaceClient.CheckNameAvailability` return value(s) have been changed from `(WorkspaceCheckNameAvailabilityResponse, error)` to `(WorkspaceClientCheckNameAvailabilityResponse, error)`
- Function `*WorkspacesClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, QuantumWorkspace, *WorkspacesBeginCreateOrUpdateOptions)` to `(context.Context, string, string, Workspace, *WorkspacesClientBeginCreateOrUpdateOptions)`
- Function `*WorkspacesClient.BeginCreateOrUpdate` return value(s) have been changed from `(WorkspacesCreateOrUpdatePollerResponse, error)` to `(WorkspacesClientCreateOrUpdatePollerResponse, error)`
- Function `*OfferingsClient.List` parameter(s) have been changed from `(string, *OfferingsListOptions)` to `(string, *OfferingsClientListOptions)`
- Function `*OfferingsClient.List` return value(s) have been changed from `(*OfferingsListPager)` to `(*OfferingsClientListPager)`
- Function `*WorkspacesClient.ListBySubscription` parameter(s) have been changed from `(*WorkspacesListBySubscriptionOptions)` to `(*WorkspacesClientListBySubscriptionOptions)`
- Function `*WorkspacesClient.ListBySubscription` return value(s) have been changed from `(*WorkspacesListBySubscriptionPager)` to `(*WorkspacesClientListBySubscriptionPager)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*WorkspacesClient.Get` parameter(s) have been changed from `(context.Context, string, string, *WorkspacesGetOptions)` to `(context.Context, string, string, *WorkspacesClientGetOptions)`
- Function `*WorkspacesClient.Get` return value(s) have been changed from `(WorkspacesGetResponse, error)` to `(WorkspacesClientGetResponse, error)`
- Type of `WorkspaceListResult.Value` has been changed from `[]*QuantumWorkspace` to `[]*Workspace`
- Function `*WorkspacesListBySubscriptionPager.PageResponse` has been removed
- Function `*OfferingsListPager.PageResponse` has been removed
- Function `*WorkspacesCreateOrUpdatePoller.ResumeToken` has been removed
- Function `*WorkspacesDeletePollerResponse.Resume` has been removed
- Function `*OfferingsListPager.Err` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `WorkspacesCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `*WorkspacesCreateOrUpdatePoller.FinalResponse` has been removed
- Function `*WorkspacesListByResourceGroupPager.PageResponse` has been removed
- Function `*WorkspacesCreateOrUpdatePoller.Done` has been removed
- Function `*OfferingsListPager.NextPage` has been removed
- Function `*OperationsListPager.PageResponse` has been removed
- Function `*WorkspacesDeletePoller.FinalResponse` has been removed
- Function `QuantumWorkspace.MarshalJSON` has been removed
- Function `*WorkspacesListBySubscriptionPager.Err` has been removed
- Function `*WorkspacesDeletePoller.Done` has been removed
- Function `*WorkspacesListBySubscriptionPager.NextPage` has been removed
- Function `*WorkspacesCreateOrUpdatePoller.Poll` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `*WorkspacesDeletePoller.ResumeToken` has been removed
- Function `*WorkspacesListByResourceGroupPager.NextPage` has been removed
- Function `*WorkspacesDeletePoller.Poll` has been removed
- Function `Resource.MarshalJSON` has been removed
- Function `*WorkspacesListByResourceGroupPager.Err` has been removed
- Function `WorkspacesDeletePollerResponse.PollUntilDone` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Function `*WorkspacesCreateOrUpdatePollerResponse.Resume` has been removed
- Struct `OfferingsListOptions` has been removed
- Struct `OfferingsListPager` has been removed
- Struct `OfferingsListResponse` has been removed
- Struct `OfferingsListResultEnvelope` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Struct `QuantumWorkspace` has been removed
- Struct `QuantumWorkspaceIdentity` has been removed
- Struct `WorkspaceCheckNameAvailabilityOptions` has been removed
- Struct `WorkspaceCheckNameAvailabilityResponse` has been removed
- Struct `WorkspaceCheckNameAvailabilityResult` has been removed
- Struct `WorkspacesBeginCreateOrUpdateOptions` has been removed
- Struct `WorkspacesBeginDeleteOptions` has been removed
- Struct `WorkspacesCreateOrUpdatePoller` has been removed
- Struct `WorkspacesCreateOrUpdatePollerResponse` has been removed
- Struct `WorkspacesCreateOrUpdateResponse` has been removed
- Struct `WorkspacesCreateOrUpdateResult` has been removed
- Struct `WorkspacesDeletePoller` has been removed
- Struct `WorkspacesDeletePollerResponse` has been removed
- Struct `WorkspacesDeleteResponse` has been removed
- Struct `WorkspacesGetOptions` has been removed
- Struct `WorkspacesGetResponse` has been removed
- Struct `WorkspacesGetResult` has been removed
- Struct `WorkspacesListByResourceGroupOptions` has been removed
- Struct `WorkspacesListByResourceGroupPager` has been removed
- Struct `WorkspacesListByResourceGroupResponse` has been removed
- Struct `WorkspacesListByResourceGroupResult` has been removed
- Struct `WorkspacesListBySubscriptionOptions` has been removed
- Struct `WorkspacesListBySubscriptionPager` has been removed
- Struct `WorkspacesListBySubscriptionResponse` has been removed
- Struct `WorkspacesListBySubscriptionResult` has been removed
- Struct `WorkspacesUpdateTagsOptions` has been removed
- Struct `WorkspacesUpdateTagsResponse` has been removed
- Struct `WorkspacesUpdateTagsResult` has been removed
- Field `Resource` of struct `TrackedResource` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed

### Features Added

- New function `*WorkspacesClientCreateOrUpdatePollerResponse.Resume(context.Context, *WorkspacesClient, string) error`
- New function `*OfferingsClientListPager.Err() error`
- New function `*WorkspacesClientListBySubscriptionPager.PageResponse() WorkspacesClientListBySubscriptionResponse`
- New function `WorkspacesClientCreateOrUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (WorkspacesClientCreateOrUpdateResponse, error)`
- New function `*WorkspacesClientCreateOrUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `WorkspacesClientDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (WorkspacesClientDeleteResponse, error)`
- New function `*WorkspacesClientListByResourceGroupPager.PageResponse() WorkspacesClientListByResourceGroupResponse`
- New function `*OperationsClientListPager.Err() error`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*WorkspacesClientDeletePoller.ResumeToken() (string, error)`
- New function `Workspace.MarshalJSON() ([]byte, error)`
- New function `*WorkspacesClientListBySubscriptionPager.NextPage(context.Context) bool`
- New function `*WorkspacesClientDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `*OfferingsClientListPager.PageResponse() OfferingsClientListResponse`
- New function `*WorkspacesClientListByResourceGroupPager.Err() error`
- New function `*WorkspacesClientDeletePoller.Done() bool`
- New function `*WorkspacesClientListBySubscriptionPager.Err() error`
- New function `*WorkspacesClientCreateOrUpdatePoller.Done() bool`
- New function `*WorkspacesClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*WorkspacesClientCreateOrUpdatePoller.ResumeToken() (string, error)`
- New function `*WorkspacesClientCreateOrUpdatePoller.FinalResponse(context.Context) (WorkspacesClientCreateOrUpdateResponse, error)`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*WorkspacesClientDeletePollerResponse.Resume(context.Context, *WorkspacesClient, string) error`
- New function `*WorkspacesClientDeletePoller.FinalResponse(context.Context) (WorkspacesClientDeleteResponse, error)`
- New function `*OfferingsClientListPager.NextPage(context.Context) bool`
- New struct `OfferingsClientListOptions`
- New struct `OfferingsClientListPager`
- New struct `OfferingsClientListResponse`
- New struct `OfferingsClientListResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New struct `Workspace`
- New struct `WorkspaceClientCheckNameAvailabilityOptions`
- New struct `WorkspaceClientCheckNameAvailabilityResponse`
- New struct `WorkspaceClientCheckNameAvailabilityResult`
- New struct `WorkspaceIdentity`
- New struct `WorkspacesClientBeginCreateOrUpdateOptions`
- New struct `WorkspacesClientBeginDeleteOptions`
- New struct `WorkspacesClientCreateOrUpdatePoller`
- New struct `WorkspacesClientCreateOrUpdatePollerResponse`
- New struct `WorkspacesClientCreateOrUpdateResponse`
- New struct `WorkspacesClientCreateOrUpdateResult`
- New struct `WorkspacesClientDeletePoller`
- New struct `WorkspacesClientDeletePollerResponse`
- New struct `WorkspacesClientDeleteResponse`
- New struct `WorkspacesClientGetOptions`
- New struct `WorkspacesClientGetResponse`
- New struct `WorkspacesClientGetResult`
- New struct `WorkspacesClientListByResourceGroupOptions`
- New struct `WorkspacesClientListByResourceGroupPager`
- New struct `WorkspacesClientListByResourceGroupResponse`
- New struct `WorkspacesClientListByResourceGroupResult`
- New struct `WorkspacesClientListBySubscriptionOptions`
- New struct `WorkspacesClientListBySubscriptionPager`
- New struct `WorkspacesClientListBySubscriptionResponse`
- New struct `WorkspacesClientListBySubscriptionResult`
- New struct `WorkspacesClientUpdateTagsOptions`
- New struct `WorkspacesClientUpdateTagsResponse`
- New struct `WorkspacesClientUpdateTagsResult`
- New field `Error` in struct `ErrorResponse`
- New field `ID` in struct `TrackedResource`
- New field `Name` in struct `TrackedResource`
- New field `Type` in struct `TrackedResource`


## 0.1.0 (2021-12-16)

- Init release.
