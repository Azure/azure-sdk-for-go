# Release History

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
