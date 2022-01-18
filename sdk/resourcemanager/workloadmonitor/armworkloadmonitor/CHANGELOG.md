# Release History

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
