# Release History

## 0.4.0 (2022-04-18)
### Breaking Changes

- Function `*OperationsClient.List` has been removed
- Function `*PrivateEndpointConnectionsClient.ListByResource` has been removed
- Function `*PrivateLinkResourcesClient.ListByResource` has been removed

### Features Added

- New function `*PrivateLinkResourcesClient.NewListByResourcePager(*PrivateLinkResourcesClientListByResourceOptions) *runtime.Pager[PrivateLinkResourcesClientListByResourceResponse]`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New function `*PrivateEndpointConnectionsClient.NewListByResourcePager(string, string, *PrivateEndpointConnectionsClientListByResourceOptions) *runtime.Pager[PrivateEndpointConnectionsClientListByResourceResponse]`


## 0.3.0 (2022-04-12)
### Breaking Changes

- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `NewPrivateLinkServicesForPowerBIClient` return value(s) have been changed from `(*PrivateLinkServicesForPowerBIClient)` to `(*PrivateLinkServicesForPowerBIClient, error)`
- Function `NewPrivateLinkServiceResourceOperationResultsClient` return value(s) have been changed from `(*PrivateLinkServiceResourceOperationResultsClient)` to `(*PrivateLinkServiceResourceOperationResultsClient, error)`
- Function `NewPrivateLinkServicesClient` return value(s) have been changed from `(*PrivateLinkServicesClient)` to `(*PrivateLinkServicesClient, error)`
- Function `NewPrivateLinkResourcesClient` return value(s) have been changed from `(*PrivateLinkResourcesClient)` to `(*PrivateLinkResourcesClient, error)`
- Function `*PrivateLinkServiceResourceOperationResultsClient.BeginGet` return value(s) have been changed from `(PrivateLinkServiceResourceOperationResultsClientGetPollerResponse, error)` to `(*armruntime.Poller[PrivateLinkServiceResourceOperationResultsClientGetResponse], error)`
- Function `NewPowerBIResourcesClient` return value(s) have been changed from `(*PowerBIResourcesClient)` to `(*PowerBIResourcesClient, error)`
- Function `NewPrivateEndpointConnectionsClient` return value(s) have been changed from `(*PrivateEndpointConnectionsClient)` to `(*PrivateEndpointConnectionsClient, error)`
- Function `*PrivateEndpointConnectionsClient.ListByResource` return value(s) have been changed from `(*PrivateEndpointConnectionsClientListByResourcePager)` to `(*runtime.Pager[PrivateEndpointConnectionsClientListByResourceResponse])`
- Function `*PrivateLinkResourcesClient.ListByResource` return value(s) have been changed from `(*PrivateLinkResourcesClientListByResourcePager)` to `(*runtime.Pager[PrivateLinkResourcesClientListByResourceResponse])`
- Function `*PrivateEndpointConnectionsClient.BeginDelete` return value(s) have been changed from `(PrivateEndpointConnectionsClientDeletePollerResponse, error)` to `(*armruntime.Poller[PrivateEndpointConnectionsClientDeleteResponse], error)`
- Type of `ErrorAdditionalInfo.Info` has been changed from `map[string]interface{}` to `interface{}`
- Function `*PrivateLinkServiceResourceOperationResultsClientGetPoller.Done` has been removed
- Function `ActionType.ToPtr` has been removed
- Function `ResourceProvisioningState.ToPtr` has been removed
- Function `Origin.ToPtr` has been removed
- Function `*PrivateEndpointConnectionsClientDeletePoller.Done` has been removed
- Function `PrivateEndpointConnectionsClientDeletePollerResponse.PollUntilDone` has been removed
- Function `*OperationsClientListPager.NextPage` has been removed
- Function `*PrivateEndpointConnectionsClientListByResourcePager.PageResponse` has been removed
- Function `*PrivateEndpointConnectionsClientDeletePoller.Poll` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Function `*PrivateLinkServiceResourceOperationResultsClientGetPoller.Poll` has been removed
- Function `*PrivateLinkResourcesClientListByResourcePager.PageResponse` has been removed
- Function `*PrivateLinkServiceResourceOperationResultsClientGetPollerResponse.Resume` has been removed
- Function `*PrivateEndpointConnectionsClientDeletePoller.FinalResponse` has been removed
- Function `PersistedConnectionStatus.ToPtr` has been removed
- Function `ActionsRequired.ToPtr` has been removed
- Function `*PrivateLinkResourcesClientListByResourcePager.Err` has been removed
- Function `*PrivateEndpointConnectionsClientDeletePoller.ResumeToken` has been removed
- Function `*PrivateEndpointConnectionsClientListByResourcePager.NextPage` has been removed
- Function `*PrivateLinkServiceResourceOperationResultsClientGetPoller.ResumeToken` has been removed
- Function `*OperationsClientListPager.Err` has been removed
- Function `*PrivateEndpointConnectionsClientListByResourcePager.Err` has been removed
- Function `PrivateLinkServiceResourceOperationResultsClientGetPollerResponse.PollUntilDone` has been removed
- Function `*PrivateLinkResourcesClientListByResourcePager.NextPage` has been removed
- Function `*PrivateLinkServiceResourceOperationResultsClientGetPoller.FinalResponse` has been removed
- Function `CreatedByType.ToPtr` has been removed
- Function `*PrivateEndpointConnectionsClientDeletePollerResponse.Resume` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Struct `PowerBIResourcesClientCreateResult` has been removed
- Struct `PowerBIResourcesClientListByResourceNameResult` has been removed
- Struct `PowerBIResourcesClientUpdateResult` has been removed
- Struct `PrivateEndpointConnectionsClientCreateResult` has been removed
- Struct `PrivateEndpointConnectionsClientDeletePoller` has been removed
- Struct `PrivateEndpointConnectionsClientDeletePollerResponse` has been removed
- Struct `PrivateEndpointConnectionsClientGetResult` has been removed
- Struct `PrivateEndpointConnectionsClientListByResourcePager` has been removed
- Struct `PrivateEndpointConnectionsClientListByResourceResult` has been removed
- Struct `PrivateLinkResourcesClientGetResult` has been removed
- Struct `PrivateLinkResourcesClientListByResourcePager` has been removed
- Struct `PrivateLinkResourcesClientListByResourceResult` has been removed
- Struct `PrivateLinkServiceResourceOperationResultsClientGetPoller` has been removed
- Struct `PrivateLinkServiceResourceOperationResultsClientGetPollerResponse` has been removed
- Struct `PrivateLinkServiceResourceOperationResultsClientGetResult` has been removed
- Struct `PrivateLinkServicesClientListByResourceGroupResult` has been removed
- Struct `PrivateLinkServicesForPowerBIClientListBySubscriptionIDResult` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `PrivateLinkServicesForPowerBIClientListBySubscriptionIDResult` of struct `PrivateLinkServicesForPowerBIClientListBySubscriptionIDResponse` has been removed
- Field `RawResponse` of struct `PrivateLinkServicesForPowerBIClientListBySubscriptionIDResponse` has been removed
- Field `PrivateEndpointConnectionsClientListByResourceResult` of struct `PrivateEndpointConnectionsClientListByResourceResponse` has been removed
- Field `RawResponse` of struct `PrivateEndpointConnectionsClientListByResourceResponse` has been removed
- Field `PowerBIResourcesClientUpdateResult` of struct `PowerBIResourcesClientUpdateResponse` has been removed
- Field `RawResponse` of struct `PowerBIResourcesClientUpdateResponse` has been removed
- Field `RawResponse` of struct `PrivateEndpointConnectionsClientDeleteResponse` has been removed
- Field `PrivateLinkServicesClientListByResourceGroupResult` of struct `PrivateLinkServicesClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `PrivateLinkServicesClientListByResourceGroupResponse` has been removed
- Field `PowerBIResourcesClientCreateResult` of struct `PowerBIResourcesClientCreateResponse` has been removed
- Field `RawResponse` of struct `PowerBIResourcesClientCreateResponse` has been removed
- Field `PrivateLinkResourcesClientGetResult` of struct `PrivateLinkResourcesClientGetResponse` has been removed
- Field `RawResponse` of struct `PrivateLinkResourcesClientGetResponse` has been removed
- Field `PrivateLinkServiceResourceOperationResultsClientGetResult` of struct `PrivateLinkServiceResourceOperationResultsClientGetResponse` has been removed
- Field `RawResponse` of struct `PrivateLinkServiceResourceOperationResultsClientGetResponse` has been removed
- Field `PrivateEndpointConnectionsClientCreateResult` of struct `PrivateEndpointConnectionsClientCreateResponse` has been removed
- Field `RawResponse` of struct `PrivateEndpointConnectionsClientCreateResponse` has been removed
- Field `PrivateEndpointConnectionsClientGetResult` of struct `PrivateEndpointConnectionsClientGetResponse` has been removed
- Field `RawResponse` of struct `PrivateEndpointConnectionsClientGetResponse` has been removed
- Field `RawResponse` of struct `PowerBIResourcesClientDeleteResponse` has been removed
- Field `PrivateLinkResourcesClientListByResourceResult` of struct `PrivateLinkResourcesClientListByResourceResponse` has been removed
- Field `RawResponse` of struct `PrivateLinkResourcesClientListByResourceResponse` has been removed
- Field `PowerBIResourcesClientListByResourceNameResult` of struct `PowerBIResourcesClientListByResourceNameResponse` has been removed
- Field `RawResponse` of struct `PowerBIResourcesClientListByResourceNameResponse` has been removed

### Features Added

- New field `TenantResourceArray` in struct `PowerBIResourcesClientListByResourceNameResponse`
- New anonymous field `TenantResource` in struct `PowerBIResourcesClientCreateResponse`
- New field `ResumeToken` in struct `PrivateEndpointConnectionsClientBeginDeleteOptions`
- New field `TenantResourceArray` in struct `PrivateLinkServicesForPowerBIClientListBySubscriptionIDResponse`
- New anonymous field `PrivateEndpointConnection` in struct `PrivateEndpointConnectionsClientGetResponse`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New anonymous field `AsyncOperationDetail` in struct `PrivateLinkServiceResourceOperationResultsClientGetResponse`
- New anonymous field `PrivateEndpointConnection` in struct `PrivateEndpointConnectionsClientCreateResponse`
- New field `TenantResourceArray` in struct `PrivateLinkServicesClientListByResourceGroupResponse`
- New field `ResumeToken` in struct `PrivateLinkServiceResourceOperationResultsClientBeginGetOptions`
- New anonymous field `PrivateLinkResource` in struct `PrivateLinkResourcesClientGetResponse`
- New anonymous field `TenantResource` in struct `PowerBIResourcesClientUpdateResponse`
- New anonymous field `PrivateEndpointConnectionListResult` in struct `PrivateEndpointConnectionsClientListByResourceResponse`
- New anonymous field `PrivateLinkResourcesListResult` in struct `PrivateLinkResourcesClientListByResourceResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*PrivateLinkResourcesClient.ListByResource` parameter(s) have been changed from `(*PrivateLinkResourcesListByResourceOptions)` to `(*PrivateLinkResourcesClientListByResourceOptions)`
- Function `*PrivateLinkResourcesClient.ListByResource` return value(s) have been changed from `(*PrivateLinkResourcesListByResourcePager)` to `(*PrivateLinkResourcesClientListByResourcePager)`
- Function `*PowerBIResourcesClient.Delete` parameter(s) have been changed from `(context.Context, *PowerBIResourcesDeleteOptions)` to `(context.Context, *PowerBIResourcesClientDeleteOptions)`
- Function `*PowerBIResourcesClient.Delete` return value(s) have been changed from `(PowerBIResourcesDeleteResponse, error)` to `(PowerBIResourcesClientDeleteResponse, error)`
- Function `*PrivateLinkResourcesClient.Get` parameter(s) have been changed from `(context.Context, string, *PrivateLinkResourcesGetOptions)` to `(context.Context, string, *PrivateLinkResourcesClientGetOptions)`
- Function `*PrivateLinkResourcesClient.Get` return value(s) have been changed from `(PrivateLinkResourcesGetResponse, error)` to `(PrivateLinkResourcesClientGetResponse, error)`
- Function `*PrivateLinkServicesForPowerBIClient.ListBySubscriptionID` parameter(s) have been changed from `(context.Context, *PrivateLinkServicesForPowerBIListBySubscriptionIDOptions)` to `(context.Context, *PrivateLinkServicesForPowerBIClientListBySubscriptionIDOptions)`
- Function `*PrivateLinkServicesForPowerBIClient.ListBySubscriptionID` return value(s) have been changed from `(PrivateLinkServicesForPowerBIListBySubscriptionIDResponse, error)` to `(PrivateLinkServicesForPowerBIClientListBySubscriptionIDResponse, error)`
- Function `*PowerBIResourcesClient.Update` parameter(s) have been changed from `(context.Context, TenantResource, *PowerBIResourcesUpdateOptions)` to `(context.Context, TenantResource, *PowerBIResourcesClientUpdateOptions)`
- Function `*PowerBIResourcesClient.Update` return value(s) have been changed from `(PowerBIResourcesUpdateResponse, error)` to `(PowerBIResourcesClientUpdateResponse, error)`
- Function `*PrivateEndpointConnectionsClient.Get` parameter(s) have been changed from `(context.Context, *PrivateEndpointConnectionsGetOptions)` to `(context.Context, *PrivateEndpointConnectionsClientGetOptions)`
- Function `*PrivateEndpointConnectionsClient.Get` return value(s) have been changed from `(PrivateEndpointConnectionsGetResponse, error)` to `(PrivateEndpointConnectionsClientGetResponse, error)`
- Function `*PowerBIResourcesClient.Create` parameter(s) have been changed from `(context.Context, TenantResource, *PowerBIResourcesCreateOptions)` to `(context.Context, TenantResource, *PowerBIResourcesClientCreateOptions)`
- Function `*PowerBIResourcesClient.Create` return value(s) have been changed from `(PowerBIResourcesCreateResponse, error)` to `(PowerBIResourcesClientCreateResponse, error)`
- Function `*PrivateEndpointConnectionsClient.Create` parameter(s) have been changed from `(context.Context, PrivateEndpointConnection, *PrivateEndpointConnectionsCreateOptions)` to `(context.Context, PrivateEndpointConnection, *PrivateEndpointConnectionsClientCreateOptions)`
- Function `*PrivateEndpointConnectionsClient.Create` return value(s) have been changed from `(PrivateEndpointConnectionsCreateResponse, error)` to `(PrivateEndpointConnectionsClientCreateResponse, error)`
- Function `*PowerBIResourcesClient.ListByResourceName` parameter(s) have been changed from `(context.Context, *PowerBIResourcesListByResourceNameOptions)` to `(context.Context, *PowerBIResourcesClientListByResourceNameOptions)`
- Function `*PowerBIResourcesClient.ListByResourceName` return value(s) have been changed from `(PowerBIResourcesListByResourceNameResponse, error)` to `(PowerBIResourcesClientListByResourceNameResponse, error)`
- Function `*PrivateEndpointConnectionsClient.BeginDelete` parameter(s) have been changed from `(context.Context, *PrivateEndpointConnectionsBeginDeleteOptions)` to `(context.Context, *PrivateEndpointConnectionsClientBeginDeleteOptions)`
- Function `*PrivateEndpointConnectionsClient.BeginDelete` return value(s) have been changed from `(PrivateEndpointConnectionsDeletePollerResponse, error)` to `(PrivateEndpointConnectionsClientDeletePollerResponse, error)`
- Function `*PrivateLinkServicesClient.ListByResourceGroup` parameter(s) have been changed from `(context.Context, *PrivateLinkServicesListByResourceGroupOptions)` to `(context.Context, *PrivateLinkServicesClientListByResourceGroupOptions)`
- Function `*PrivateLinkServicesClient.ListByResourceGroup` return value(s) have been changed from `(PrivateLinkServicesListByResourceGroupResponse, error)` to `(PrivateLinkServicesClientListByResourceGroupResponse, error)`
- Function `*PrivateEndpointConnectionsClient.ListByResource` parameter(s) have been changed from `(string, string, *PrivateEndpointConnectionsListByResourceOptions)` to `(string, string, *PrivateEndpointConnectionsClientListByResourceOptions)`
- Function `*PrivateEndpointConnectionsClient.ListByResource` return value(s) have been changed from `(*PrivateEndpointConnectionsListByResourcePager)` to `(*PrivateEndpointConnectionsClientListByResourcePager)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*PrivateLinkServiceResourceOperationResultsClient.BeginGet` parameter(s) have been changed from `(context.Context, *PrivateLinkServiceResourceOperationResultsBeginGetOptions)` to `(context.Context, *PrivateLinkServiceResourceOperationResultsClientBeginGetOptions)`
- Function `*PrivateLinkServiceResourceOperationResultsClient.BeginGet` return value(s) have been changed from `(PrivateLinkServiceResourceOperationResultsGetPollerResponse, error)` to `(PrivateLinkServiceResourceOperationResultsClientGetPollerResponse, error)`
- Function `*PrivateLinkServiceResourceOperationResultsGetPoller.Done` has been removed
- Function `*PrivateLinkServiceResourceOperationResultsGetPoller.ResumeToken` has been removed
- Function `*OperationsListPager.PageResponse` has been removed
- Function `*PrivateLinkResourcesListByResourcePager.NextPage` has been removed
- Function `PrivateEndpointConnectionsDeletePollerResponse.PollUntilDone` has been removed
- Function `*PrivateEndpointConnectionsDeletePoller.Done` has been removed
- Function `*PrivateEndpointConnectionsListByResourcePager.NextPage` has been removed
- Function `*PrivateLinkServiceResourceOperationResultsGetPoller.FinalResponse` has been removed
- Function `PrivateLinkServiceResourceOperationResultsGetPollerResponse.PollUntilDone` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Function `*PrivateEndpointConnectionsDeletePoller.Poll` has been removed
- Function `*PrivateEndpointConnectionsDeletePoller.FinalResponse` has been removed
- Function `*PrivateLinkResourcesListByResourcePager.Err` has been removed
- Function `*PrivateLinkResourcesListByResourcePager.PageResponse` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `*PrivateLinkServiceResourceOperationResultsGetPoller.Poll` has been removed
- Function `*PrivateEndpointConnectionsDeletePollerResponse.Resume` has been removed
- Function `*PrivateEndpointConnectionsListByResourcePager.PageResponse` has been removed
- Function `*PrivateEndpointConnectionsListByResourcePager.Err` has been removed
- Function `*PrivateEndpointConnectionsDeletePoller.ResumeToken` has been removed
- Function `*PrivateLinkServiceResourceOperationResultsGetPollerResponse.Resume` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Struct `PowerBIResourcesCreateOptions` has been removed
- Struct `PowerBIResourcesCreateResponse` has been removed
- Struct `PowerBIResourcesCreateResult` has been removed
- Struct `PowerBIResourcesDeleteOptions` has been removed
- Struct `PowerBIResourcesDeleteResponse` has been removed
- Struct `PowerBIResourcesListByResourceNameOptions` has been removed
- Struct `PowerBIResourcesListByResourceNameResponse` has been removed
- Struct `PowerBIResourcesListByResourceNameResult` has been removed
- Struct `PowerBIResourcesUpdateOptions` has been removed
- Struct `PowerBIResourcesUpdateResponse` has been removed
- Struct `PowerBIResourcesUpdateResult` has been removed
- Struct `PrivateEndpointConnectionsBeginDeleteOptions` has been removed
- Struct `PrivateEndpointConnectionsCreateOptions` has been removed
- Struct `PrivateEndpointConnectionsCreateResponse` has been removed
- Struct `PrivateEndpointConnectionsCreateResult` has been removed
- Struct `PrivateEndpointConnectionsDeletePoller` has been removed
- Struct `PrivateEndpointConnectionsDeletePollerResponse` has been removed
- Struct `PrivateEndpointConnectionsDeleteResponse` has been removed
- Struct `PrivateEndpointConnectionsGetOptions` has been removed
- Struct `PrivateEndpointConnectionsGetResponse` has been removed
- Struct `PrivateEndpointConnectionsGetResult` has been removed
- Struct `PrivateEndpointConnectionsListByResourceOptions` has been removed
- Struct `PrivateEndpointConnectionsListByResourcePager` has been removed
- Struct `PrivateEndpointConnectionsListByResourceResponse` has been removed
- Struct `PrivateEndpointConnectionsListByResourceResult` has been removed
- Struct `PrivateLinkResourcesGetOptions` has been removed
- Struct `PrivateLinkResourcesGetResponse` has been removed
- Struct `PrivateLinkResourcesGetResult` has been removed
- Struct `PrivateLinkResourcesListByResourceOptions` has been removed
- Struct `PrivateLinkResourcesListByResourcePager` has been removed
- Struct `PrivateLinkResourcesListByResourceResponse` has been removed
- Struct `PrivateLinkResourcesListByResourceResult` has been removed
- Struct `PrivateLinkServiceResourceOperationResultsBeginGetOptions` has been removed
- Struct `PrivateLinkServiceResourceOperationResultsGetPoller` has been removed
- Struct `PrivateLinkServiceResourceOperationResultsGetPollerResponse` has been removed
- Struct `PrivateLinkServiceResourceOperationResultsGetResponse` has been removed
- Struct `PrivateLinkServiceResourceOperationResultsGetResult` has been removed
- Struct `PrivateLinkServicesForPowerBIListBySubscriptionIDOptions` has been removed
- Struct `PrivateLinkServicesForPowerBIListBySubscriptionIDResponse` has been removed
- Struct `PrivateLinkServicesForPowerBIListBySubscriptionIDResult` has been removed
- Struct `PrivateLinkServicesListByResourceGroupOptions` has been removed
- Struct `PrivateLinkServicesListByResourceGroupResponse` has been removed
- Struct `PrivateLinkServicesListByResourceGroupResult` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed

### Features Added

- New function `*PrivateLinkServiceResourceOperationResultsClientGetPoller.Poll(context.Context) (*http.Response, error)`
- New function `*PrivateLinkResourcesClientListByResourcePager.Err() error`
- New function `PrivateLinkServiceResourceOperationResultsClientGetPollerResponse.PollUntilDone(context.Context, time.Duration) (PrivateLinkServiceResourceOperationResultsClientGetResponse, error)`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*PrivateLinkServiceResourceOperationResultsClientGetPoller.Done() bool`
- New function `PrivateEndpointConnectionsClientDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (PrivateEndpointConnectionsClientDeleteResponse, error)`
- New function `*PrivateEndpointConnectionsClientDeletePoller.ResumeToken() (string, error)`
- New function `*PrivateEndpointConnectionsClientListByResourcePager.Err() error`
- New function `*PrivateEndpointConnectionsClientDeletePoller.FinalResponse(context.Context) (PrivateEndpointConnectionsClientDeleteResponse, error)`
- New function `*PrivateLinkResourcesClientListByResourcePager.NextPage(context.Context) bool`
- New function `*PrivateLinkServiceResourceOperationResultsClientGetPollerResponse.Resume(context.Context, *PrivateLinkServiceResourceOperationResultsClient, string) error`
- New function `*PrivateEndpointConnectionsClientDeletePoller.Done() bool`
- New function `*PrivateLinkResourcesClientListByResourcePager.PageResponse() PrivateLinkResourcesClientListByResourceResponse`
- New function `*PrivateEndpointConnectionsClientListByResourcePager.PageResponse() PrivateEndpointConnectionsClientListByResourceResponse`
- New function `*PrivateEndpointConnectionsClientListByResourcePager.NextPage(context.Context) bool`
- New function `*PrivateEndpointConnectionsClientDeletePollerResponse.Resume(context.Context, *PrivateEndpointConnectionsClient, string) error`
- New function `*OperationsClientListPager.Err() error`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*PrivateLinkServiceResourceOperationResultsClientGetPoller.FinalResponse(context.Context) (PrivateLinkServiceResourceOperationResultsClientGetResponse, error)`
- New function `*PrivateEndpointConnectionsClientDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `*PrivateLinkServiceResourceOperationResultsClientGetPoller.ResumeToken() (string, error)`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New struct `PowerBIResourcesClientCreateOptions`
- New struct `PowerBIResourcesClientCreateResponse`
- New struct `PowerBIResourcesClientCreateResult`
- New struct `PowerBIResourcesClientDeleteOptions`
- New struct `PowerBIResourcesClientDeleteResponse`
- New struct `PowerBIResourcesClientListByResourceNameOptions`
- New struct `PowerBIResourcesClientListByResourceNameResponse`
- New struct `PowerBIResourcesClientListByResourceNameResult`
- New struct `PowerBIResourcesClientUpdateOptions`
- New struct `PowerBIResourcesClientUpdateResponse`
- New struct `PowerBIResourcesClientUpdateResult`
- New struct `PrivateEndpointConnectionsClientBeginDeleteOptions`
- New struct `PrivateEndpointConnectionsClientCreateOptions`
- New struct `PrivateEndpointConnectionsClientCreateResponse`
- New struct `PrivateEndpointConnectionsClientCreateResult`
- New struct `PrivateEndpointConnectionsClientDeletePoller`
- New struct `PrivateEndpointConnectionsClientDeletePollerResponse`
- New struct `PrivateEndpointConnectionsClientDeleteResponse`
- New struct `PrivateEndpointConnectionsClientGetOptions`
- New struct `PrivateEndpointConnectionsClientGetResponse`
- New struct `PrivateEndpointConnectionsClientGetResult`
- New struct `PrivateEndpointConnectionsClientListByResourceOptions`
- New struct `PrivateEndpointConnectionsClientListByResourcePager`
- New struct `PrivateEndpointConnectionsClientListByResourceResponse`
- New struct `PrivateEndpointConnectionsClientListByResourceResult`
- New struct `PrivateLinkResourcesClientGetOptions`
- New struct `PrivateLinkResourcesClientGetResponse`
- New struct `PrivateLinkResourcesClientGetResult`
- New struct `PrivateLinkResourcesClientListByResourceOptions`
- New struct `PrivateLinkResourcesClientListByResourcePager`
- New struct `PrivateLinkResourcesClientListByResourceResponse`
- New struct `PrivateLinkResourcesClientListByResourceResult`
- New struct `PrivateLinkServiceResourceOperationResultsClientBeginGetOptions`
- New struct `PrivateLinkServiceResourceOperationResultsClientGetPoller`
- New struct `PrivateLinkServiceResourceOperationResultsClientGetPollerResponse`
- New struct `PrivateLinkServiceResourceOperationResultsClientGetResponse`
- New struct `PrivateLinkServiceResourceOperationResultsClientGetResult`
- New struct `PrivateLinkServicesClientListByResourceGroupOptions`
- New struct `PrivateLinkServicesClientListByResourceGroupResponse`
- New struct `PrivateLinkServicesClientListByResourceGroupResult`
- New struct `PrivateLinkServicesForPowerBIClientListBySubscriptionIDOptions`
- New struct `PrivateLinkServicesForPowerBIClientListBySubscriptionIDResponse`
- New struct `PrivateLinkServicesForPowerBIClientListBySubscriptionIDResult`
- New field `Error` in struct `ErrorResponse`


## 0.1.0 (2021-12-22)

- Init release.
