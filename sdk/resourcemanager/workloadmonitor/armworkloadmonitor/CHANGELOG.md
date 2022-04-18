# Release History

## 0.4.0 (2022-04-18)
### Breaking Changes

- Function `*HealthMonitorsClient.ListStateChanges` has been removed
- Function `*OperationsClient.List` has been removed
- Function `*HealthMonitorsClient.List` has been removed

### Features Added

- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New function `*HealthMonitorsClient.NewListPager(string, string, string, string, string, *HealthMonitorsClientListOptions) *runtime.Pager[HealthMonitorsClientListResponse]`
- New function `*HealthMonitorsClient.NewListStateChangesPager(string, string, string, string, string, string, *HealthMonitorsClientListStateChangesOptions) *runtime.Pager[HealthMonitorsClientListStateChangesResponse]`


## 0.3.0 (2022-04-13)
### Breaking Changes

- Function `NewHealthMonitorsClient` return value(s) have been changed from `(*HealthMonitorsClient)` to `(*HealthMonitorsClient, error)`
- Function `*HealthMonitorsClient.List` return value(s) have been changed from `(*HealthMonitorsClientListPager)` to `(*runtime.Pager[HealthMonitorsClientListResponse])`
- Function `*HealthMonitorsClient.ListStateChanges` return value(s) have been changed from `(*HealthMonitorsClientListStateChangesPager)` to `(*runtime.Pager[HealthMonitorsClientListStateChangesResponse])`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Type of `HealthMonitorProperties.MonitorConfiguration` has been changed from `map[string]interface{}` to `interface{}`
- Type of `HealthMonitorProperties.Evidence` has been changed from `map[string]interface{}` to `interface{}`
- Type of `HealthMonitorStateChangeProperties.Evidence` has been changed from `map[string]interface{}` to `interface{}`
- Type of `HealthMonitorStateChangeProperties.MonitorConfiguration` has been changed from `map[string]interface{}` to `interface{}`
- Function `*OperationsClientListPager.Err` has been removed
- Function `HealthState.ToPtr` has been removed
- Function `*HealthMonitorsClientListStateChangesPager.PageResponse` has been removed
- Function `*HealthMonitorsClientListPager.Err` has been removed
- Function `*HealthMonitorsClientListPager.NextPage` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Function `*HealthMonitorsClientListPager.PageResponse` has been removed
- Function `*OperationsClientListPager.NextPage` has been removed
- Function `*HealthMonitorsClientListStateChangesPager.NextPage` has been removed
- Function `*HealthMonitorsClientListStateChangesPager.Err` has been removed
- Struct `HealthMonitorsClientGetResult` has been removed
- Struct `HealthMonitorsClientGetStateChangeResult` has been removed
- Struct `HealthMonitorsClientListPager` has been removed
- Struct `HealthMonitorsClientListResult` has been removed
- Struct `HealthMonitorsClientListStateChangesPager` has been removed
- Struct `HealthMonitorsClientListStateChangesResult` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Field `HealthMonitorsClientListStateChangesResult` of struct `HealthMonitorsClientListStateChangesResponse` has been removed
- Field `RawResponse` of struct `HealthMonitorsClientListStateChangesResponse` has been removed
- Field `HealthMonitorsClientListResult` of struct `HealthMonitorsClientListResponse` has been removed
- Field `RawResponse` of struct `HealthMonitorsClientListResponse` has been removed
- Field `HealthMonitorsClientGetStateChangeResult` of struct `HealthMonitorsClientGetStateChangeResponse` has been removed
- Field `RawResponse` of struct `HealthMonitorsClientGetStateChangeResponse` has been removed
- Field `HealthMonitorsClientGetResult` of struct `HealthMonitorsClientGetResponse` has been removed
- Field `RawResponse` of struct `HealthMonitorsClientGetResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed

### Features Added

- New anonymous field `OperationList` in struct `OperationsClientListResponse`
- New anonymous field `HealthMonitor` in struct `HealthMonitorsClientGetResponse`
- New anonymous field `HealthMonitorStateChange` in struct `HealthMonitorsClientGetStateChangeResponse`
- New anonymous field `HealthMonitorList` in struct `HealthMonitorsClientListResponse`
- New anonymous field `HealthMonitorStateChangeList` in struct `HealthMonitorsClientListStateChangesResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*HealthMonitorsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, string, string, string, *HealthMonitorsGetOptions)` to `(context.Context, string, string, string, string, string, string, *HealthMonitorsClientGetOptions)`
- Function `*HealthMonitorsClient.Get` return value(s) have been changed from `(HealthMonitorsGetResponse, error)` to `(HealthMonitorsClientGetResponse, error)`
- Function `*HealthMonitorsClient.GetStateChange` parameter(s) have been changed from `(context.Context, string, string, string, string, string, string, string, *HealthMonitorsGetStateChangeOptions)` to `(context.Context, string, string, string, string, string, string, string, *HealthMonitorsClientGetStateChangeOptions)`
- Function `*HealthMonitorsClient.GetStateChange` return value(s) have been changed from `(HealthMonitorsGetStateChangeResponse, error)` to `(HealthMonitorsClientGetStateChangeResponse, error)`
- Function `*HealthMonitorsClient.List` parameter(s) have been changed from `(string, string, string, string, string, *HealthMonitorsListOptions)` to `(string, string, string, string, string, *HealthMonitorsClientListOptions)`
- Function `*HealthMonitorsClient.List` return value(s) have been changed from `(*HealthMonitorsListPager)` to `(*HealthMonitorsClientListPager)`
- Function `*HealthMonitorsClient.ListStateChanges` parameter(s) have been changed from `(string, string, string, string, string, string, *HealthMonitorsListStateChangesOptions)` to `(string, string, string, string, string, string, *HealthMonitorsClientListStateChangesOptions)`
- Function `*HealthMonitorsClient.ListStateChanges` return value(s) have been changed from `(*HealthMonitorsListStateChangesPager)` to `(*HealthMonitorsClientListStateChangesPager)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*HealthMonitorsListPager.Err` has been removed
- Function `*HealthMonitorsListStateChangesPager.Err` has been removed
- Function `*HealthMonitorsListPager.NextPage` has been removed
- Function `*OperationsListPager.PageResponse` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `*HealthMonitorsListPager.PageResponse` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Function `*HealthMonitorsListStateChangesPager.NextPage` has been removed
- Function `*HealthMonitorsListStateChangesPager.PageResponse` has been removed
- Function `*OperationsListPager.Err` has been removed
- Struct `HealthMonitorsGetOptions` has been removed
- Struct `HealthMonitorsGetResponse` has been removed
- Struct `HealthMonitorsGetResult` has been removed
- Struct `HealthMonitorsGetStateChangeOptions` has been removed
- Struct `HealthMonitorsGetStateChangeResponse` has been removed
- Struct `HealthMonitorsGetStateChangeResult` has been removed
- Struct `HealthMonitorsListOptions` has been removed
- Struct `HealthMonitorsListPager` has been removed
- Struct `HealthMonitorsListResponse` has been removed
- Struct `HealthMonitorsListResult` has been removed
- Struct `HealthMonitorsListStateChangesOptions` has been removed
- Struct `HealthMonitorsListStateChangesPager` has been removed
- Struct `HealthMonitorsListStateChangesResponse` has been removed
- Struct `HealthMonitorsListStateChangesResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `Resource` of struct `HealthMonitorStateChange` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed
- Field `Resource` of struct `HealthMonitor` has been removed

### Features Added

- New function `*HealthMonitorsClientListStateChangesPager.NextPage(context.Context) bool`
- New function `*HealthMonitorsClientListStateChangesPager.Err() error`
- New function `*HealthMonitorsClientListStateChangesPager.PageResponse() HealthMonitorsClientListStateChangesResponse`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*HealthMonitorsClientListPager.NextPage(context.Context) bool`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*HealthMonitorsClientListPager.Err() error`
- New function `*OperationsClientListPager.Err() error`
- New function `*HealthMonitorsClientListPager.PageResponse() HealthMonitorsClientListResponse`
- New struct `HealthMonitorsClientGetOptions`
- New struct `HealthMonitorsClientGetResponse`
- New struct `HealthMonitorsClientGetResult`
- New struct `HealthMonitorsClientGetStateChangeOptions`
- New struct `HealthMonitorsClientGetStateChangeResponse`
- New struct `HealthMonitorsClientGetStateChangeResult`
- New struct `HealthMonitorsClientListOptions`
- New struct `HealthMonitorsClientListPager`
- New struct `HealthMonitorsClientListResponse`
- New struct `HealthMonitorsClientListResult`
- New struct `HealthMonitorsClientListStateChangesOptions`
- New struct `HealthMonitorsClientListStateChangesPager`
- New struct `HealthMonitorsClientListStateChangesResponse`
- New struct `HealthMonitorsClientListStateChangesResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New field `Error` in struct `ErrorResponse`
- New field `ID` in struct `HealthMonitorStateChange`
- New field `Name` in struct `HealthMonitorStateChange`
- New field `Type` in struct `HealthMonitorStateChange`
- New field `Name` in struct `HealthMonitor`
- New field `Type` in struct `HealthMonitor`
- New field `ID` in struct `HealthMonitor`


## 0.1.0 (2021-12-16)

- Init release.
