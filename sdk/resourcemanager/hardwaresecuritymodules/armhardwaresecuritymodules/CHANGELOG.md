# Release History

## 0.4.0 (2022-04-15)
### Breaking Changes

- Function `*DedicatedHsmClient.ListBySubscription` has been removed
- Function `*DedicatedHsmClient.ListOutboundNetworkDependenciesEndpoints` has been removed
- Function `*OperationsClient.List` has been removed
- Function `*DedicatedHsmClient.ListByResourceGroup` has been removed

### Features Added

- New function `*DedicatedHsmClient.NewListOutboundNetworkDependenciesEndpointsPager(string, string, *DedicatedHsmClientListOutboundNetworkDependenciesEndpointsOptions) *runtime.Pager[DedicatedHsmClientListOutboundNetworkDependenciesEndpointsResponse]`
- New function `*DedicatedHsmClient.NewListBySubscriptionPager(*DedicatedHsmClientListBySubscriptionOptions) *runtime.Pager[DedicatedHsmClientListBySubscriptionResponse]`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New function `*DedicatedHsmClient.NewListByResourceGroupPager(string, *DedicatedHsmClientListByResourceGroupOptions) *runtime.Pager[DedicatedHsmClientListByResourceGroupResponse]`


## 0.3.0 (2022-04-11)
### Breaking Changes

- Function `*DedicatedHsmClient.ListByResourceGroup` return value(s) have been changed from `(*DedicatedHsmClientListByResourceGroupPager)` to `(*runtime.Pager[DedicatedHsmClientListByResourceGroupResponse])`
- Function `NewDedicatedHsmClient` return value(s) have been changed from `(*DedicatedHsmClient)` to `(*DedicatedHsmClient, error)`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(context.Context, *OperationsClientListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(OperationsClientListResponse, error)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `*DedicatedHsmClient.BeginDelete` return value(s) have been changed from `(DedicatedHsmClientDeletePollerResponse, error)` to `(*armruntime.Poller[DedicatedHsmClientDeleteResponse], error)`
- Function `*DedicatedHsmClient.ListOutboundNetworkDependenciesEndpoints` return value(s) have been changed from `(*DedicatedHsmClientListOutboundNetworkDependenciesEndpointsPager)` to `(*runtime.Pager[DedicatedHsmClientListOutboundNetworkDependenciesEndpointsResponse])`
- Function `*DedicatedHsmClient.BeginCreateOrUpdate` return value(s) have been changed from `(DedicatedHsmClientCreateOrUpdatePollerResponse, error)` to `(*armruntime.Poller[DedicatedHsmClientCreateOrUpdateResponse], error)`
- Function `*DedicatedHsmClient.ListBySubscription` return value(s) have been changed from `(*DedicatedHsmClientListBySubscriptionPager)` to `(*runtime.Pager[DedicatedHsmClientListBySubscriptionResponse])`
- Function `*DedicatedHsmClient.BeginUpdate` return value(s) have been changed from `(DedicatedHsmClientUpdatePollerResponse, error)` to `(*armruntime.Poller[DedicatedHsmClientUpdateResponse], error)`
- Function `*DedicatedHsmClientDeletePollerResponse.Resume` has been removed
- Function `IdentityType.ToPtr` has been removed
- Function `*DedicatedHsmClientListOutboundNetworkDependenciesEndpointsPager.PageResponse` has been removed
- Function `*DedicatedHsmClientCreateOrUpdatePoller.ResumeToken` has been removed
- Function `*DedicatedHsmClientListByResourceGroupPager.Err` has been removed
- Function `*DedicatedHsmClientUpdatePoller.ResumeToken` has been removed
- Function `*DedicatedHsmClientUpdatePoller.Done` has been removed
- Function `*DedicatedHsmClientUpdatePoller.Poll` has been removed
- Function `*DedicatedHsmClientListByResourceGroupPager.NextPage` has been removed
- Function `*DedicatedHsmClientCreateOrUpdatePoller.Poll` has been removed
- Function `*DedicatedHsmClientListByResourceGroupPager.PageResponse` has been removed
- Function `*DedicatedHsmClientListOutboundNetworkDependenciesEndpointsPager.Err` has been removed
- Function `*DedicatedHsmClientListBySubscriptionPager.NextPage` has been removed
- Function `*DedicatedHsmClientDeletePoller.FinalResponse` has been removed
- Function `DedicatedHsmClientCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `*DedicatedHsmClientDeletePoller.Done` has been removed
- Function `*DedicatedHsmClientCreateOrUpdatePollerResponse.Resume` has been removed
- Function `DedicatedHsmClientDeletePollerResponse.PollUntilDone` has been removed
- Function `JSONWebKeyType.ToPtr` has been removed
- Function `*DedicatedHsmClientUpdatePoller.FinalResponse` has been removed
- Function `*DedicatedHsmClientCreateOrUpdatePoller.Done` has been removed
- Function `*DedicatedHsmClientDeletePoller.Poll` has been removed
- Function `DedicatedHsmClientUpdatePollerResponse.PollUntilDone` has been removed
- Function `*DedicatedHsmClientUpdatePollerResponse.Resume` has been removed
- Function `*DedicatedHsmClientListOutboundNetworkDependenciesEndpointsPager.NextPage` has been removed
- Function `*DedicatedHsmClientListBySubscriptionPager.Err` has been removed
- Function `*DedicatedHsmClientDeletePoller.ResumeToken` has been removed
- Function `*DedicatedHsmClientListBySubscriptionPager.PageResponse` has been removed
- Function `*DedicatedHsmClientCreateOrUpdatePoller.FinalResponse` has been removed
- Function `SKUName.ToPtr` has been removed
- Struct `DedicatedHsmClientCreateOrUpdatePoller` has been removed
- Struct `DedicatedHsmClientCreateOrUpdatePollerResponse` has been removed
- Struct `DedicatedHsmClientCreateOrUpdateResult` has been removed
- Struct `DedicatedHsmClientDeletePoller` has been removed
- Struct `DedicatedHsmClientDeletePollerResponse` has been removed
- Struct `DedicatedHsmClientGetResult` has been removed
- Struct `DedicatedHsmClientListByResourceGroupPager` has been removed
- Struct `DedicatedHsmClientListByResourceGroupResult` has been removed
- Struct `DedicatedHsmClientListBySubscriptionPager` has been removed
- Struct `DedicatedHsmClientListBySubscriptionResult` has been removed
- Struct `DedicatedHsmClientListOutboundNetworkDependenciesEndpointsPager` has been removed
- Struct `DedicatedHsmClientListOutboundNetworkDependenciesEndpointsResult` has been removed
- Struct `DedicatedHsmClientUpdatePoller` has been removed
- Struct `DedicatedHsmClientUpdatePollerResponse` has been removed
- Struct `DedicatedHsmClientUpdateResult` has been removed
- Struct `OperationsClientListResult` has been removed
- Field `DedicatedHsmClientCreateOrUpdateResult` of struct `DedicatedHsmClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `DedicatedHsmClientCreateOrUpdateResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `DedicatedHsmClientListBySubscriptionResult` of struct `DedicatedHsmClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `DedicatedHsmClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `DedicatedHsmClientDeleteResponse` has been removed
- Field `DedicatedHsmClientListOutboundNetworkDependenciesEndpointsResult` of struct `DedicatedHsmClientListOutboundNetworkDependenciesEndpointsResponse` has been removed
- Field `RawResponse` of struct `DedicatedHsmClientListOutboundNetworkDependenciesEndpointsResponse` has been removed
- Field `DedicatedHsmClientGetResult` of struct `DedicatedHsmClientGetResponse` has been removed
- Field `RawResponse` of struct `DedicatedHsmClientGetResponse` has been removed
- Field `DedicatedHsmClientUpdateResult` of struct `DedicatedHsmClientUpdateResponse` has been removed
- Field `RawResponse` of struct `DedicatedHsmClientUpdateResponse` has been removed
- Field `DedicatedHsmClientListByResourceGroupResult` of struct `DedicatedHsmClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `DedicatedHsmClientListByResourceGroupResponse` has been removed

### Features Added

- New field `ResumeToken` in struct `DedicatedHsmClientBeginDeleteOptions`
- New anonymous field `DedicatedHsmListResult` in struct `DedicatedHsmClientListBySubscriptionResponse`
- New anonymous field `DedicatedHsmOperationListResult` in struct `OperationsClientListResponse`
- New anonymous field `DedicatedHsm` in struct `DedicatedHsmClientGetResponse`
- New anonymous field `OutboundEnvironmentEndpointCollection` in struct `DedicatedHsmClientListOutboundNetworkDependenciesEndpointsResponse`
- New field `ResumeToken` in struct `DedicatedHsmClientBeginCreateOrUpdateOptions`
- New anonymous field `DedicatedHsm` in struct `DedicatedHsmClientCreateOrUpdateResponse`
- New anonymous field `DedicatedHsmListResult` in struct `DedicatedHsmClientListByResourceGroupResponse`
- New anonymous field `DedicatedHsm` in struct `DedicatedHsmClientUpdateResponse`
- New field `ResumeToken` in struct `DedicatedHsmClientBeginUpdateOptions`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*DedicatedHsmClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, DedicatedHsm, *DedicatedHsmBeginCreateOrUpdateOptions)` to `(context.Context, string, string, DedicatedHsm, *DedicatedHsmClientBeginCreateOrUpdateOptions)`
- Function `*DedicatedHsmClient.BeginCreateOrUpdate` return value(s) have been changed from `(DedicatedHsmCreateOrUpdatePollerResponse, error)` to `(DedicatedHsmClientCreateOrUpdatePollerResponse, error)`
- Function `*DedicatedHsmClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, DedicatedHsmPatchParameters, *DedicatedHsmBeginUpdateOptions)` to `(context.Context, string, string, DedicatedHsmPatchParameters, *DedicatedHsmClientBeginUpdateOptions)`
- Function `*DedicatedHsmClient.BeginUpdate` return value(s) have been changed from `(DedicatedHsmUpdatePollerResponse, error)` to `(DedicatedHsmClientUpdatePollerResponse, error)`
- Function `*DedicatedHsmClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *DedicatedHsmBeginDeleteOptions)` to `(context.Context, string, string, *DedicatedHsmClientBeginDeleteOptions)`
- Function `*DedicatedHsmClient.BeginDelete` return value(s) have been changed from `(DedicatedHsmDeletePollerResponse, error)` to `(DedicatedHsmClientDeletePollerResponse, error)`
- Function `*DedicatedHsmClient.Get` parameter(s) have been changed from `(context.Context, string, string, *DedicatedHsmGetOptions)` to `(context.Context, string, string, *DedicatedHsmClientGetOptions)`
- Function `*DedicatedHsmClient.Get` return value(s) have been changed from `(DedicatedHsmGetResponse, error)` to `(DedicatedHsmClientGetResponse, error)`
- Function `*DedicatedHsmClient.ListByResourceGroup` parameter(s) have been changed from `(string, *DedicatedHsmListByResourceGroupOptions)` to `(string, *DedicatedHsmClientListByResourceGroupOptions)`
- Function `*DedicatedHsmClient.ListByResourceGroup` return value(s) have been changed from `(*DedicatedHsmListByResourceGroupPager)` to `(*DedicatedHsmClientListByResourceGroupPager)`
- Function `*DedicatedHsmClient.ListBySubscription` parameter(s) have been changed from `(*DedicatedHsmListBySubscriptionOptions)` to `(*DedicatedHsmClientListBySubscriptionOptions)`
- Function `*DedicatedHsmClient.ListBySubscription` return value(s) have been changed from `(*DedicatedHsmListBySubscriptionPager)` to `(*DedicatedHsmClientListBySubscriptionPager)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(context.Context, *OperationsListOptions)` to `(context.Context, *OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(OperationsListResponse, error)` to `(OperationsClientListResponse, error)`
- Function `*DedicatedHsmUpdatePoller.Done` has been removed
- Function `*DedicatedHsmUpdatePoller.FinalResponse` has been removed
- Function `*DedicatedHsmCreateOrUpdatePoller.FinalResponse` has been removed
- Function `*DedicatedHsmListBySubscriptionPager.NextPage` has been removed
- Function `*DedicatedHsmDeletePoller.ResumeToken` has been removed
- Function `*DedicatedHsmListBySubscriptionPager.Err` has been removed
- Function `DedicatedHsmDeletePollerResponse.PollUntilDone` has been removed
- Function `*DedicatedHsmCreateOrUpdatePoller.ResumeToken` has been removed
- Function `DedicatedHsmUpdatePollerResponse.PollUntilDone` has been removed
- Function `*DedicatedHsmUpdatePoller.Poll` has been removed
- Function `DedicatedHsmCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `*DedicatedHsmDeletePoller.Poll` has been removed
- Function `*DedicatedHsmCreateOrUpdatePoller.Done` has been removed
- Function `*DedicatedHsmDeletePoller.Done` has been removed
- Function `*DedicatedHsmUpdatePoller.ResumeToken` has been removed
- Function `*DedicatedHsmListBySubscriptionPager.PageResponse` has been removed
- Function `*DedicatedHsmCreateOrUpdatePoller.Poll` has been removed
- Function `*DedicatedHsmDeletePoller.FinalResponse` has been removed
- Function `*DedicatedHsmUpdatePollerResponse.Resume` has been removed
- Function `*DedicatedHsmCreateOrUpdatePollerResponse.Resume` has been removed
- Function `*DedicatedHsmListByResourceGroupPager.PageResponse` has been removed
- Function `*DedicatedHsmListByResourceGroupPager.Err` has been removed
- Function `DedicatedHsmError.Error` has been removed
- Function `*DedicatedHsmListByResourceGroupPager.NextPage` has been removed
- Function `*DedicatedHsmDeletePollerResponse.Resume` has been removed
- Struct `DedicatedHsmBeginCreateOrUpdateOptions` has been removed
- Struct `DedicatedHsmBeginDeleteOptions` has been removed
- Struct `DedicatedHsmBeginUpdateOptions` has been removed
- Struct `DedicatedHsmCreateOrUpdatePoller` has been removed
- Struct `DedicatedHsmCreateOrUpdatePollerResponse` has been removed
- Struct `DedicatedHsmCreateOrUpdateResponse` has been removed
- Struct `DedicatedHsmCreateOrUpdateResult` has been removed
- Struct `DedicatedHsmDeletePoller` has been removed
- Struct `DedicatedHsmDeletePollerResponse` has been removed
- Struct `DedicatedHsmDeleteResponse` has been removed
- Struct `DedicatedHsmGetOptions` has been removed
- Struct `DedicatedHsmGetResponse` has been removed
- Struct `DedicatedHsmGetResult` has been removed
- Struct `DedicatedHsmListByResourceGroupOptions` has been removed
- Struct `DedicatedHsmListByResourceGroupPager` has been removed
- Struct `DedicatedHsmListByResourceGroupResponse` has been removed
- Struct `DedicatedHsmListByResourceGroupResult` has been removed
- Struct `DedicatedHsmListBySubscriptionOptions` has been removed
- Struct `DedicatedHsmListBySubscriptionPager` has been removed
- Struct `DedicatedHsmListBySubscriptionResponse` has been removed
- Struct `DedicatedHsmListBySubscriptionResult` has been removed
- Struct `DedicatedHsmUpdatePoller` has been removed
- Struct `DedicatedHsmUpdatePollerResponse` has been removed
- Struct `DedicatedHsmUpdateResponse` has been removed
- Struct `DedicatedHsmUpdateResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `InnerError` of struct `DedicatedHsmError` has been removed
- Field `Resource` of struct `DedicatedHsm` has been removed

### Features Added

- New const `IdentityTypeManagedIdentity`
- New const `IdentityTypeKey`
- New const `SKUNamePayShield10KLMK1CPS60`
- New const `IdentityTypeUser`
- New const `SKUNamePayShield10KLMK1CPS250`
- New const `SKUNamePayShield10KLMK2CPS250`
- New const `SKUNamePayShield10KLMK2CPS2500`
- New const `IdentityTypeApplication`
- New const `SKUNamePayShield10KLMK1CPS2500`
- New const `SKUNamePayShield10KLMK2CPS60`
- New function `DedicatedHsmClientCreateOrUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (DedicatedHsmClientCreateOrUpdateResponse, error)`
- New function `*DedicatedHsmClientUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*DedicatedHsmClientListBySubscriptionPager.NextPage(context.Context) bool`
- New function `*DedicatedHsmClientListOutboundNetworkDependenciesEndpointsPager.Err() error`
- New function `DedicatedHsmClientDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (DedicatedHsmClientDeleteResponse, error)`
- New function `*timeRFC3339.Parse(string) error`
- New function `*SystemData.UnmarshalJSON([]byte) error`
- New function `*DedicatedHsmClientDeletePollerResponse.Resume(context.Context, *DedicatedHsmClient, string) error`
- New function `timeRFC3339.MarshalText() ([]byte, error)`
- New function `*DedicatedHsmClientUpdatePoller.FinalResponse(context.Context) (DedicatedHsmClientUpdateResponse, error)`
- New function `OutboundEnvironmentEndpointCollection.MarshalJSON() ([]byte, error)`
- New function `*DedicatedHsmClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*DedicatedHsmClientUpdatePoller.Done() bool`
- New function `*DedicatedHsmClientUpdatePoller.ResumeToken() (string, error)`
- New function `*DedicatedHsmClientUpdatePollerResponse.Resume(context.Context, *DedicatedHsmClient, string) error`
- New function `*DedicatedHsmClientDeletePoller.Done() bool`
- New function `OutboundEnvironmentEndpoint.MarshalJSON() ([]byte, error)`
- New function `timeRFC3339.MarshalJSON() ([]byte, error)`
- New function `*DedicatedHsmClientCreateOrUpdatePoller.FinalResponse(context.Context) (DedicatedHsmClientCreateOrUpdateResponse, error)`
- New function `IdentityType.ToPtr() *IdentityType`
- New function `*DedicatedHsmClientListByResourceGroupPager.PageResponse() DedicatedHsmClientListByResourceGroupResponse`
- New function `EndpointDependency.MarshalJSON() ([]byte, error)`
- New function `*DedicatedHsmClientCreateOrUpdatePollerResponse.Resume(context.Context, *DedicatedHsmClient, string) error`
- New function `*DedicatedHsmClientListOutboundNetworkDependenciesEndpointsPager.PageResponse() DedicatedHsmClientListOutboundNetworkDependenciesEndpointsResponse`
- New function `*DedicatedHsmClientListOutboundNetworkDependenciesEndpointsPager.NextPage(context.Context) bool`
- New function `*DedicatedHsmClientListBySubscriptionPager.PageResponse() DedicatedHsmClientListBySubscriptionResponse`
- New function `*timeRFC3339.UnmarshalText([]byte) error`
- New function `*DedicatedHsmClientDeletePoller.ResumeToken() (string, error)`
- New function `*timeRFC3339.UnmarshalJSON([]byte) error`
- New function `PossibleIdentityTypeValues() []IdentityType`
- New function `*DedicatedHsmClientListByResourceGroupPager.Err() error`
- New function `*DedicatedHsmClientListBySubscriptionPager.Err() error`
- New function `*DedicatedHsmClientCreateOrUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `SystemData.MarshalJSON() ([]byte, error)`
- New function `*DedicatedHsmClientCreateOrUpdatePoller.Done() bool`
- New function `*DedicatedHsmClientDeletePoller.FinalResponse(context.Context) (DedicatedHsmClientDeleteResponse, error)`
- New function `DedicatedHsmClientUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (DedicatedHsmClientUpdateResponse, error)`
- New function `*DedicatedHsmClientCreateOrUpdatePoller.ResumeToken() (string, error)`
- New function `*DedicatedHsmClient.ListOutboundNetworkDependenciesEndpoints(string, string, *DedicatedHsmClientListOutboundNetworkDependenciesEndpointsOptions) *DedicatedHsmClientListOutboundNetworkDependenciesEndpointsPager`
- New function `*DedicatedHsmClientDeletePoller.Poll(context.Context) (*http.Response, error)`
- New struct `DedicatedHsmClientBeginCreateOrUpdateOptions`
- New struct `DedicatedHsmClientBeginDeleteOptions`
- New struct `DedicatedHsmClientBeginUpdateOptions`
- New struct `DedicatedHsmClientCreateOrUpdatePoller`
- New struct `DedicatedHsmClientCreateOrUpdatePollerResponse`
- New struct `DedicatedHsmClientCreateOrUpdateResponse`
- New struct `DedicatedHsmClientCreateOrUpdateResult`
- New struct `DedicatedHsmClientDeletePoller`
- New struct `DedicatedHsmClientDeletePollerResponse`
- New struct `DedicatedHsmClientDeleteResponse`
- New struct `DedicatedHsmClientGetOptions`
- New struct `DedicatedHsmClientGetResponse`
- New struct `DedicatedHsmClientGetResult`
- New struct `DedicatedHsmClientListByResourceGroupOptions`
- New struct `DedicatedHsmClientListByResourceGroupPager`
- New struct `DedicatedHsmClientListByResourceGroupResponse`
- New struct `DedicatedHsmClientListByResourceGroupResult`
- New struct `DedicatedHsmClientListBySubscriptionOptions`
- New struct `DedicatedHsmClientListBySubscriptionPager`
- New struct `DedicatedHsmClientListBySubscriptionResponse`
- New struct `DedicatedHsmClientListBySubscriptionResult`
- New struct `DedicatedHsmClientListOutboundNetworkDependenciesEndpointsOptions`
- New struct `DedicatedHsmClientListOutboundNetworkDependenciesEndpointsPager`
- New struct `DedicatedHsmClientListOutboundNetworkDependenciesEndpointsResponse`
- New struct `DedicatedHsmClientListOutboundNetworkDependenciesEndpointsResult`
- New struct `DedicatedHsmClientUpdatePoller`
- New struct `DedicatedHsmClientUpdatePollerResponse`
- New struct `DedicatedHsmClientUpdateResponse`
- New struct `DedicatedHsmClientUpdateResult`
- New struct `EndpointDependency`
- New struct `EndpointDetail`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New struct `OutboundEnvironmentEndpoint`
- New struct `OutboundEnvironmentEndpointCollection`
- New struct `SystemData`
- New field `SystemData` in struct `DedicatedHsm`
- New field `Type` in struct `DedicatedHsm`
- New field `Location` in struct `DedicatedHsm`
- New field `Tags` in struct `DedicatedHsm`
- New field `Zones` in struct `DedicatedHsm`
- New field `Name` in struct `DedicatedHsm`
- New field `SKU` in struct `DedicatedHsm`
- New field `ID` in struct `DedicatedHsm`
- New field `Error` in struct `DedicatedHsmError`
- New field `ManagementNetworkProfile` in struct `DedicatedHsmProperties`


## 0.1.0 (2021-12-07)

- Init release.
