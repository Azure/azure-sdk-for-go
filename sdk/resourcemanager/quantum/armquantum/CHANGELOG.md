# Release History

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
