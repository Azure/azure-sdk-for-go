# Release History

## 0.3.0 (2022-05-17)
### Breaking Changes

- Function `*GrafanaClient.BeginCreate` return value(s) have been changed from `(*armruntime.Poller[GrafanaClientCreateResponse], error)` to `(*runtime.Poller[GrafanaClientCreateResponse], error)`
- Function `*GrafanaClient.BeginDelete` return value(s) have been changed from `(*armruntime.Poller[GrafanaClientDeleteResponse], error)` to `(*runtime.Poller[GrafanaClientDeleteResponse], error)`
- Function `ErrorDetail.MarshalJSON` has been removed
- Function `OperationListResult.MarshalJSON` has been removed
- Function `ManagedGrafanaListResponse.MarshalJSON` has been removed


## 0.2.0 (2022-04-15)
### Breaking Changes

- Function `*GrafanaClient.ListByResourceGroup` has been removed
- Function `*GrafanaClient.List` has been removed
- Function `*OperationsClient.List` has been removed

### Features Added

- New function `*GrafanaClient.NewListByResourceGroupPager(string, *GrafanaClientListByResourceGroupOptions) *runtime.Pager[GrafanaClientListByResourceGroupResponse]`
- New function `*GrafanaClient.NewListPager(*GrafanaClientListOptions) *runtime.Pager[GrafanaClientListResponse]`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`


## 0.1.0 (2022-04-11)

- Init release.