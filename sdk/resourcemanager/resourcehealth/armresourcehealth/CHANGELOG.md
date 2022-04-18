# Release History

## 0.6.0 (2022-04-18)
### Breaking Changes

- Function `*AvailabilityStatusesClient.List` has been removed
- Function `*AvailabilityStatusesClient.ListBySubscriptionID` has been removed
- Function `*AvailabilityStatusesClient.ListByResourceGroup` has been removed

### Features Added

- New function `*AvailabilityStatusesClient.NewListPager(string, *AvailabilityStatusesClientListOptions) *runtime.Pager[AvailabilityStatusesClientListResponse]`
- New function `*AvailabilityStatusesClient.NewListBySubscriptionIDPager(*AvailabilityStatusesClientListBySubscriptionIDOptions) *runtime.Pager[AvailabilityStatusesClientListBySubscriptionIDResponse]`
- New function `*AvailabilityStatusesClient.NewListByResourceGroupPager(string, *AvailabilityStatusesClientListByResourceGroupOptions) *runtime.Pager[AvailabilityStatusesClientListByResourceGroupResponse]`


## 0.5.0 (2022-04-13)
### Breaking Changes

- Function `NewAvailabilityStatusesClient` return value(s) have been changed from `(*AvailabilityStatusesClient)` to `(*AvailabilityStatusesClient, error)`
- Function `*AvailabilityStatusesClient.ListByResourceGroup` return value(s) have been changed from `(*AvailabilityStatusesClientListByResourceGroupPager)` to `(*runtime.Pager[AvailabilityStatusesClientListByResourceGroupResponse])`
- Function `*AvailabilityStatusesClient.ListBySubscriptionID` return value(s) have been changed from `(*AvailabilityStatusesClientListBySubscriptionIDPager)` to `(*runtime.Pager[AvailabilityStatusesClientListBySubscriptionIDResponse])`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `*AvailabilityStatusesClient.List` return value(s) have been changed from `(*AvailabilityStatusesClientListPager)` to `(*runtime.Pager[AvailabilityStatusesClientListResponse])`
- Const `SeverityValuesWarning` has been removed
- Const `StageValuesActive` has been removed
- Const `StageValuesArchived` has been removed
- Const `StageValuesResolve` has been removed
- Const `SeverityValuesError` has been removed
- Const `SeverityValuesInformation` has been removed
- Function `*AvailabilityStatusesClientListByResourceGroupPager.PageResponse` has been removed
- Function `*ChildAvailabilityStatusesClientListPager.PageResponse` has been removed
- Function `NewChildResourcesClient` has been removed
- Function `*AvailabilityStatusesClientListBySubscriptionIDPager.NextPage` has been removed
- Function `*ChildResourcesClientListPager.NextPage` has been removed
- Function `NewChildAvailabilityStatusesClient` has been removed
- Function `*AvailabilityStatusesClientListBySubscriptionIDPager.PageResponse` has been removed
- Function `*ChildAvailabilityStatusesClient.GetByResource` has been removed
- Function `*AvailabilityStatusesClientListPager.NextPage` has been removed
- Function `SeverityValues.ToPtr` has been removed
- Function `PossibleStageValuesValues` has been removed
- Function `*ChildResourcesClient.List` has been removed
- Function `*EmergingIssuesClientListPager.NextPage` has been removed
- Function `*AvailabilityStatusPropertiesRecentlyResolvedState.UnmarshalJSON` has been removed
- Function `*ChildAvailabilityStatusesClient.List` has been removed
- Function `*EmergingIssuesClient.Get` has been removed
- Function `*AvailabilityStatusesClientListPager.PageResponse` has been removed
- Function `*EmergingIssuesClient.List` has been removed
- Function `EmergingIssueListResult.MarshalJSON` has been removed
- Function `EmergingIssueImpact.MarshalJSON` has been removed
- Function `*EmergingIssuesClientListPager.PageResponse` has been removed
- Function `PossibleSeverityValuesValues` has been removed
- Function `EmergingIssue.MarshalJSON` has been removed
- Function `*AvailabilityStatusesClientListBySubscriptionIDPager.Err` has been removed
- Function `AvailabilityStatusPropertiesRecentlyResolvedState.MarshalJSON` has been removed
- Function `*ChildAvailabilityStatusesClientListPager.Err` has been removed
- Function `*AvailabilityStatusesClientListPager.Err` has been removed
- Function `AvailabilityStateValues.ToPtr` has been removed
- Function `*EmergingIssuesClientListPager.Err` has been removed
- Function `*AvailabilityStatusesClientListByResourceGroupPager.Err` has been removed
- Function `StageValues.ToPtr` has been removed
- Function `*ChildAvailabilityStatusesClientListPager.NextPage` has been removed
- Function `*ChildResourcesClientListPager.Err` has been removed
- Function `*EmergingIssue.UnmarshalJSON` has been removed
- Function `NewEmergingIssuesClient` has been removed
- Function `*ChildResourcesClientListPager.PageResponse` has been removed
- Function `ReasonChronicityTypes.ToPtr` has been removed
- Function `*StatusActiveEvent.UnmarshalJSON` has been removed
- Function `StatusActiveEvent.MarshalJSON` has been removed
- Function `*AvailabilityStatusesClientListByResourceGroupPager.NextPage` has been removed
- Struct `AvailabilityStatusPropertiesRecentlyResolvedState` has been removed
- Struct `AvailabilityStatusesClientGetByResourceResult` has been removed
- Struct `AvailabilityStatusesClientListByResourceGroupPager` has been removed
- Struct `AvailabilityStatusesClientListByResourceGroupResult` has been removed
- Struct `AvailabilityStatusesClientListBySubscriptionIDPager` has been removed
- Struct `AvailabilityStatusesClientListBySubscriptionIDResult` has been removed
- Struct `AvailabilityStatusesClientListPager` has been removed
- Struct `AvailabilityStatusesClientListResult` has been removed
- Struct `ChildAvailabilityStatusesClient` has been removed
- Struct `ChildAvailabilityStatusesClientGetByResourceOptions` has been removed
- Struct `ChildAvailabilityStatusesClientGetByResourceResponse` has been removed
- Struct `ChildAvailabilityStatusesClientGetByResourceResult` has been removed
- Struct `ChildAvailabilityStatusesClientListOptions` has been removed
- Struct `ChildAvailabilityStatusesClientListPager` has been removed
- Struct `ChildAvailabilityStatusesClientListResponse` has been removed
- Struct `ChildAvailabilityStatusesClientListResult` has been removed
- Struct `ChildResourcesClient` has been removed
- Struct `ChildResourcesClientListOptions` has been removed
- Struct `ChildResourcesClientListPager` has been removed
- Struct `ChildResourcesClientListResponse` has been removed
- Struct `ChildResourcesClientListResult` has been removed
- Struct `EmergingIssue` has been removed
- Struct `EmergingIssueImpact` has been removed
- Struct `EmergingIssueListResult` has been removed
- Struct `EmergingIssuesClient` has been removed
- Struct `EmergingIssuesClientGetOptions` has been removed
- Struct `EmergingIssuesClientGetResponse` has been removed
- Struct `EmergingIssuesClientGetResult` has been removed
- Struct `EmergingIssuesClientListOptions` has been removed
- Struct `EmergingIssuesClientListPager` has been removed
- Struct `EmergingIssuesClientListResponse` has been removed
- Struct `EmergingIssuesClientListResult` has been removed
- Struct `EmergingIssuesGetResult` has been removed
- Struct `OperationsClientListResult` has been removed
- Struct `StatusActiveEvent` has been removed
- Field `AvailabilityStatusesClientListBySubscriptionIDResult` of struct `AvailabilityStatusesClientListBySubscriptionIDResponse` has been removed
- Field `RawResponse` of struct `AvailabilityStatusesClientListBySubscriptionIDResponse` has been removed
- Field `RecentlyResolvedState` of struct `AvailabilityStatusProperties` has been removed
- Field `OccuredTime` of struct `AvailabilityStatusProperties` has been removed
- Field `AvailabilityStatusesClientListByResourceGroupResult` of struct `AvailabilityStatusesClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `AvailabilityStatusesClientListByResourceGroupResponse` has been removed
- Field `AvailabilityStatusesClientListResult` of struct `AvailabilityStatusesClientListResponse` has been removed
- Field `RawResponse` of struct `AvailabilityStatusesClientListResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `AvailabilityStatusesClientGetByResourceResult` of struct `AvailabilityStatusesClientGetByResourceResponse` has been removed
- Field `RawResponse` of struct `AvailabilityStatusesClientGetByResourceResponse` has been removed

### Features Added

- New const `ReasonTypeValuesUserInitiated`
- New const `AvailabilityStateValuesDegraded`
- New const `ReasonTypeValuesPlanned`
- New const `ReasonTypeValuesUnplanned`
- New function `PossibleReasonTypeValuesValues() []ReasonTypeValues`
- New function `AvailabilityStatusPropertiesRecentlyResolved.MarshalJSON() ([]byte, error)`
- New function `ImpactedResourceStatusProperties.MarshalJSON() ([]byte, error)`
- New function `*AvailabilityStatusPropertiesRecentlyResolved.UnmarshalJSON([]byte) error`
- New function `*ImpactedResourceStatusProperties.UnmarshalJSON([]byte) error`
- New struct `AvailabilityStatusPropertiesRecentlyResolved`
- New struct `ErrorResponse`
- New struct `ErrorResponseError`
- New struct `ImpactedResourceStatus`
- New struct `ImpactedResourceStatusProperties`
- New anonymous field `AvailabilityStatusListResult` in struct `AvailabilityStatusesClientListBySubscriptionIDResponse`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New anonymous field `AvailabilityStatusListResult` in struct `AvailabilityStatusesClientListResponse`
- New anonymous field `AvailabilityStatus` in struct `AvailabilityStatusesClientGetByResourceResponse`
- New anonymous field `AvailabilityStatusListResult` in struct `AvailabilityStatusesClientListByResourceGroupResponse`
- New field `Title` in struct `AvailabilityStatusProperties`
- New field `RecentlyResolved` in struct `AvailabilityStatusProperties`
- New field `OccurredTime` in struct `AvailabilityStatusProperties`


## 0.4.0 (2022-02-22)
### Breaking Changes

- Function `*EmergingIssuesClient.Get` parameter(s) have been changed from `(context.Context, Enum0, *EmergingIssuesClientGetOptions)` to `(context.Context, *EmergingIssuesClientGetOptions)`
- Const `Enum0Default` has been removed
- Function `Enum0.ToPtr` has been removed
- Function `PossibleEnum0Values` has been removed
- Struct `ErrorResponse` has been removed


## 0.3.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.3.0 (2022-01-13)
### Breaking Changes

- Function `*AvailabilityStatusesClient.GetByResource` parameter(s) have been changed from `(context.Context, string, *AvailabilityStatusesGetByResourceOptions)` to `(context.Context, string, *AvailabilityStatusesClientGetByResourceOptions)`
- Function `*AvailabilityStatusesClient.GetByResource` return value(s) have been changed from `(AvailabilityStatusesGetByResourceResponse, error)` to `(AvailabilityStatusesClientGetByResourceResponse, error)`
- Function `*AvailabilityStatusesClient.ListBySubscriptionID` parameter(s) have been changed from `(*AvailabilityStatusesListBySubscriptionIDOptions)` to `(*AvailabilityStatusesClientListBySubscriptionIDOptions)`
- Function `*AvailabilityStatusesClient.ListBySubscriptionID` return value(s) have been changed from `(*AvailabilityStatusesListBySubscriptionIDPager)` to `(*AvailabilityStatusesClientListBySubscriptionIDPager)`
- Function `*ChildAvailabilityStatusesClient.List` parameter(s) have been changed from `(string, *ChildAvailabilityStatusesListOptions)` to `(string, *ChildAvailabilityStatusesClientListOptions)`
- Function `*ChildAvailabilityStatusesClient.List` return value(s) have been changed from `(*ChildAvailabilityStatusesListPager)` to `(*ChildAvailabilityStatusesClientListPager)`
- Function `*EmergingIssuesClient.List` parameter(s) have been changed from `(*EmergingIssuesListOptions)` to `(*EmergingIssuesClientListOptions)`
- Function `*EmergingIssuesClient.List` return value(s) have been changed from `(*EmergingIssuesListPager)` to `(*EmergingIssuesClientListPager)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(context.Context, *OperationsListOptions)` to `(context.Context, *OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(OperationsListResponse, error)` to `(OperationsClientListResponse, error)`
- Function `*AvailabilityStatusesClient.List` parameter(s) have been changed from `(string, *AvailabilityStatusesListOptions)` to `(string, *AvailabilityStatusesClientListOptions)`
- Function `*AvailabilityStatusesClient.List` return value(s) have been changed from `(*AvailabilityStatusesListPager)` to `(*AvailabilityStatusesClientListPager)`
- Function `*ChildAvailabilityStatusesClient.GetByResource` parameter(s) have been changed from `(context.Context, string, *ChildAvailabilityStatusesGetByResourceOptions)` to `(context.Context, string, *ChildAvailabilityStatusesClientGetByResourceOptions)`
- Function `*ChildAvailabilityStatusesClient.GetByResource` return value(s) have been changed from `(ChildAvailabilityStatusesGetByResourceResponse, error)` to `(ChildAvailabilityStatusesClientGetByResourceResponse, error)`
- Function `*EmergingIssuesClient.Get` parameter(s) have been changed from `(context.Context, Enum0, *EmergingIssuesGetOptions)` to `(context.Context, Enum0, *EmergingIssuesClientGetOptions)`
- Function `*EmergingIssuesClient.Get` return value(s) have been changed from `(EmergingIssuesGetResponse, error)` to `(EmergingIssuesClientGetResponse, error)`
- Function `*ChildResourcesClient.List` parameter(s) have been changed from `(string, *ChildResourcesListOptions)` to `(string, *ChildResourcesClientListOptions)`
- Function `*ChildResourcesClient.List` return value(s) have been changed from `(*ChildResourcesListPager)` to `(*ChildResourcesClientListPager)`
- Function `*AvailabilityStatusesClient.ListByResourceGroup` parameter(s) have been changed from `(string, *AvailabilityStatusesListByResourceGroupOptions)` to `(string, *AvailabilityStatusesClientListByResourceGroupOptions)`
- Function `*AvailabilityStatusesClient.ListByResourceGroup` return value(s) have been changed from `(*AvailabilityStatusesListByResourceGroupPager)` to `(*AvailabilityStatusesClientListByResourceGroupPager)`
- Function `*AvailabilityStatusesListPager.PageResponse` has been removed
- Function `*ChildAvailabilityStatusesListPager.Err` has been removed
- Function `*ChildResourcesListPager.NextPage` has been removed
- Function `*AvailabilityStatusesListByResourceGroupPager.NextPage` has been removed
- Function `*AvailabilityStatusesListBySubscriptionIDPager.NextPage` has been removed
- Function `*ChildAvailabilityStatusesListPager.PageResponse` has been removed
- Function `*AvailabilityStatusesListBySubscriptionIDPager.Err` has been removed
- Function `*ChildResourcesListPager.Err` has been removed
- Function `*EmergingIssuesListPager.PageResponse` has been removed
- Function `*EmergingIssuesListPager.Err` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `*ChildAvailabilityStatusesListPager.NextPage` has been removed
- Function `*AvailabilityStatusesListBySubscriptionIDPager.PageResponse` has been removed
- Function `*AvailabilityStatusesListByResourceGroupPager.PageResponse` has been removed
- Function `*EmergingIssuesListPager.NextPage` has been removed
- Function `*ChildResourcesListPager.PageResponse` has been removed
- Function `*AvailabilityStatusesListPager.NextPage` has been removed
- Function `*AvailabilityStatusesListPager.Err` has been removed
- Function `*AvailabilityStatusesListByResourceGroupPager.Err` has been removed
- Struct `AvailabilityStatusesGetByResourceOptions` has been removed
- Struct `AvailabilityStatusesGetByResourceResponse` has been removed
- Struct `AvailabilityStatusesGetByResourceResult` has been removed
- Struct `AvailabilityStatusesListByResourceGroupOptions` has been removed
- Struct `AvailabilityStatusesListByResourceGroupPager` has been removed
- Struct `AvailabilityStatusesListByResourceGroupResponse` has been removed
- Struct `AvailabilityStatusesListByResourceGroupResult` has been removed
- Struct `AvailabilityStatusesListBySubscriptionIDOptions` has been removed
- Struct `AvailabilityStatusesListBySubscriptionIDPager` has been removed
- Struct `AvailabilityStatusesListBySubscriptionIDResponse` has been removed
- Struct `AvailabilityStatusesListBySubscriptionIDResult` has been removed
- Struct `AvailabilityStatusesListOptions` has been removed
- Struct `AvailabilityStatusesListPager` has been removed
- Struct `AvailabilityStatusesListResponse` has been removed
- Struct `AvailabilityStatusesListResult` has been removed
- Struct `ChildAvailabilityStatusesGetByResourceOptions` has been removed
- Struct `ChildAvailabilityStatusesGetByResourceResponse` has been removed
- Struct `ChildAvailabilityStatusesGetByResourceResult` has been removed
- Struct `ChildAvailabilityStatusesListOptions` has been removed
- Struct `ChildAvailabilityStatusesListPager` has been removed
- Struct `ChildAvailabilityStatusesListResponse` has been removed
- Struct `ChildAvailabilityStatusesListResult` has been removed
- Struct `ChildResourcesListOptions` has been removed
- Struct `ChildResourcesListPager` has been removed
- Struct `ChildResourcesListResponse` has been removed
- Struct `ChildResourcesListResult` has been removed
- Struct `EmergingIssuesGetOptions` has been removed
- Struct `EmergingIssuesGetResponse` has been removed
- Struct `EmergingIssuesGetResultEnvelope` has been removed
- Struct `EmergingIssuesListOptions` has been removed
- Struct `EmergingIssuesListPager` has been removed
- Struct `EmergingIssuesListResponse` has been removed
- Struct `EmergingIssuesListResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `Resource` of struct `EmergingIssuesGetResult` has been removed

### Features Added

- New function `*EmergingIssuesClientListPager.NextPage(context.Context) bool`
- New function `*ChildResourcesClientListPager.Err() error`
- New function `*AvailabilityStatusesClientListPager.Err() error`
- New function `*EmergingIssuesClientListPager.PageResponse() EmergingIssuesClientListResponse`
- New function `*AvailabilityStatusesClientListByResourceGroupPager.Err() error`
- New function `*EmergingIssuesClientListPager.Err() error`
- New function `*AvailabilityStatusesClientListPager.PageResponse() AvailabilityStatusesClientListResponse`
- New function `*AvailabilityStatusesClientListByResourceGroupPager.PageResponse() AvailabilityStatusesClientListByResourceGroupResponse`
- New function `*AvailabilityStatusesClientListBySubscriptionIDPager.Err() error`
- New function `*AvailabilityStatusesClientListBySubscriptionIDPager.PageResponse() AvailabilityStatusesClientListBySubscriptionIDResponse`
- New function `*ChildAvailabilityStatusesClientListPager.NextPage(context.Context) bool`
- New function `*ChildResourcesClientListPager.PageResponse() ChildResourcesClientListResponse`
- New function `*AvailabilityStatusesClientListBySubscriptionIDPager.NextPage(context.Context) bool`
- New function `*ChildResourcesClientListPager.NextPage(context.Context) bool`
- New function `*AvailabilityStatusesClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*AvailabilityStatusesClientListPager.NextPage(context.Context) bool`
- New function `*ChildAvailabilityStatusesClientListPager.Err() error`
- New function `*ChildAvailabilityStatusesClientListPager.PageResponse() ChildAvailabilityStatusesClientListResponse`
- New struct `AvailabilityStatusesClientGetByResourceOptions`
- New struct `AvailabilityStatusesClientGetByResourceResponse`
- New struct `AvailabilityStatusesClientGetByResourceResult`
- New struct `AvailabilityStatusesClientListByResourceGroupOptions`
- New struct `AvailabilityStatusesClientListByResourceGroupPager`
- New struct `AvailabilityStatusesClientListByResourceGroupResponse`
- New struct `AvailabilityStatusesClientListByResourceGroupResult`
- New struct `AvailabilityStatusesClientListBySubscriptionIDOptions`
- New struct `AvailabilityStatusesClientListBySubscriptionIDPager`
- New struct `AvailabilityStatusesClientListBySubscriptionIDResponse`
- New struct `AvailabilityStatusesClientListBySubscriptionIDResult`
- New struct `AvailabilityStatusesClientListOptions`
- New struct `AvailabilityStatusesClientListPager`
- New struct `AvailabilityStatusesClientListResponse`
- New struct `AvailabilityStatusesClientListResult`
- New struct `ChildAvailabilityStatusesClientGetByResourceOptions`
- New struct `ChildAvailabilityStatusesClientGetByResourceResponse`
- New struct `ChildAvailabilityStatusesClientGetByResourceResult`
- New struct `ChildAvailabilityStatusesClientListOptions`
- New struct `ChildAvailabilityStatusesClientListPager`
- New struct `ChildAvailabilityStatusesClientListResponse`
- New struct `ChildAvailabilityStatusesClientListResult`
- New struct `ChildResourcesClientListOptions`
- New struct `ChildResourcesClientListPager`
- New struct `ChildResourcesClientListResponse`
- New struct `ChildResourcesClientListResult`
- New struct `EmergingIssuesClientGetOptions`
- New struct `EmergingIssuesClientGetResponse`
- New struct `EmergingIssuesClientGetResult`
- New struct `EmergingIssuesClientListOptions`
- New struct `EmergingIssuesClientListPager`
- New struct `EmergingIssuesClientListResponse`
- New struct `EmergingIssuesClientListResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New field `ID` in struct `EmergingIssuesGetResult`
- New field `Name` in struct `EmergingIssuesGetResult`
- New field `Type` in struct `EmergingIssuesGetResult`


## 0.2.0 (2021-10-29)

### Breaking Changes

- `arm.Connection` has been removed in `github.com/Azure/azure-sdk-for-go/sdk/azcore/v0.20.0`
- The parameters of `NewXXXClient` has been changed from `(con *arm.Connection, subscriptionID string)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`

## 0.1.0 (2021-10-26)

- Initial preview release.
