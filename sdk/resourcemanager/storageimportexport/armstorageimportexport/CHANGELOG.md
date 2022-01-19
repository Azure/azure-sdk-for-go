# Release History

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*JobsClient.Update` parameter(s) have been changed from `(context.Context, string, string, UpdateJobParameters, *JobsUpdateOptions)` to `(context.Context, string, string, UpdateJobParameters, *JobsClientUpdateOptions)`
- Function `*JobsClient.Update` return value(s) have been changed from `(JobsUpdateResponse, error)` to `(JobsClientUpdateResponse, error)`
- Function `*JobsClient.ListByResourceGroup` parameter(s) have been changed from `(string, *JobsListByResourceGroupOptions)` to `(string, *JobsClientListByResourceGroupOptions)`
- Function `*JobsClient.ListByResourceGroup` return value(s) have been changed from `(*JobsListByResourceGroupPager)` to `(*JobsClientListByResourceGroupPager)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(context.Context, *OperationsListOptions)` to `(context.Context, *OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(OperationsListResponse, error)` to `(OperationsClientListResponse, error)`
- Function `*JobsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *JobsGetOptions)` to `(context.Context, string, string, *JobsClientGetOptions)`
- Function `*JobsClient.Get` return value(s) have been changed from `(JobsGetResponse, error)` to `(JobsClientGetResponse, error)`
- Function `*LocationsClient.List` parameter(s) have been changed from `(context.Context, *LocationsListOptions)` to `(context.Context, *LocationsClientListOptions)`
- Function `*LocationsClient.List` return value(s) have been changed from `(LocationsListResponse, error)` to `(LocationsClientListResponse, error)`
- Function `*JobsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *JobsDeleteOptions)` to `(context.Context, string, string, *JobsClientDeleteOptions)`
- Function `*JobsClient.Delete` return value(s) have been changed from `(JobsDeleteResponse, error)` to `(JobsClientDeleteResponse, error)`
- Function `*JobsClient.ListBySubscription` parameter(s) have been changed from `(*JobsListBySubscriptionOptions)` to `(*JobsClientListBySubscriptionOptions)`
- Function `*JobsClient.ListBySubscription` return value(s) have been changed from `(*JobsListBySubscriptionPager)` to `(*JobsClientListBySubscriptionPager)`
- Function `*JobsClient.Create` parameter(s) have been changed from `(context.Context, string, string, PutJobParameters, *JobsCreateOptions)` to `(context.Context, string, string, PutJobParameters, *JobsClientCreateOptions)`
- Function `*JobsClient.Create` return value(s) have been changed from `(JobsCreateResponse, error)` to `(JobsClientCreateResponse, error)`
- Function `*LocationsClient.Get` parameter(s) have been changed from `(context.Context, string, *LocationsGetOptions)` to `(context.Context, string, *LocationsClientGetOptions)`
- Function `*LocationsClient.Get` return value(s) have been changed from `(LocationsGetResponse, error)` to `(LocationsClientGetResponse, error)`
- Function `*BitLockerKeysClient.List` parameter(s) have been changed from `(context.Context, string, string, *BitLockerKeysListOptions)` to `(context.Context, string, string, *BitLockerKeysClientListOptions)`
- Function `*BitLockerKeysClient.List` return value(s) have been changed from `(BitLockerKeysListResponse, error)` to `(BitLockerKeysClientListResponse, error)`
- Function `*JobsListByResourceGroupPager.PageResponse` has been removed
- Function `*JobsListByResourceGroupPager.Err` has been removed
- Function `*JobsListByResourceGroupPager.NextPage` has been removed
- Function `*JobsListBySubscriptionPager.Err` has been removed
- Function `*JobsListBySubscriptionPager.NextPage` has been removed
- Function `*JobsListBySubscriptionPager.PageResponse` has been removed
- Function `ErrorResponse.Error` has been removed
- Struct `BitLockerKeysListOptions` has been removed
- Struct `BitLockerKeysListResponse` has been removed
- Struct `BitLockerKeysListResult` has been removed
- Struct `JobsCreateOptions` has been removed
- Struct `JobsCreateResponse` has been removed
- Struct `JobsCreateResult` has been removed
- Struct `JobsDeleteOptions` has been removed
- Struct `JobsDeleteResponse` has been removed
- Struct `JobsGetOptions` has been removed
- Struct `JobsGetResponse` has been removed
- Struct `JobsGetResult` has been removed
- Struct `JobsListByResourceGroupOptions` has been removed
- Struct `JobsListByResourceGroupPager` has been removed
- Struct `JobsListByResourceGroupResponse` has been removed
- Struct `JobsListByResourceGroupResult` has been removed
- Struct `JobsListBySubscriptionOptions` has been removed
- Struct `JobsListBySubscriptionPager` has been removed
- Struct `JobsListBySubscriptionResponse` has been removed
- Struct `JobsListBySubscriptionResult` has been removed
- Struct `JobsUpdateOptions` has been removed
- Struct `JobsUpdateResponse` has been removed
- Struct `JobsUpdateResult` has been removed
- Struct `LocationsGetOptions` has been removed
- Struct `LocationsGetResponse` has been removed
- Struct `LocationsGetResult` has been removed
- Struct `LocationsListOptions` has been removed
- Struct `LocationsListResponse` has been removed
- Struct `LocationsListResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed

### Features Added

- New function `*JobsClientListBySubscriptionPager.Err() error`
- New function `*JobsClientListBySubscriptionPager.NextPage(context.Context) bool`
- New function `*JobsClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*JobsClientListByResourceGroupPager.PageResponse() JobsClientListByResourceGroupResponse`
- New function `*JobsClientListByResourceGroupPager.Err() error`
- New function `*JobsClientListBySubscriptionPager.PageResponse() JobsClientListBySubscriptionResponse`
- New struct `BitLockerKeysClientListOptions`
- New struct `BitLockerKeysClientListResponse`
- New struct `BitLockerKeysClientListResult`
- New struct `JobsClientCreateOptions`
- New struct `JobsClientCreateResponse`
- New struct `JobsClientCreateResult`
- New struct `JobsClientDeleteOptions`
- New struct `JobsClientDeleteResponse`
- New struct `JobsClientGetOptions`
- New struct `JobsClientGetResponse`
- New struct `JobsClientGetResult`
- New struct `JobsClientListByResourceGroupOptions`
- New struct `JobsClientListByResourceGroupPager`
- New struct `JobsClientListByResourceGroupResponse`
- New struct `JobsClientListByResourceGroupResult`
- New struct `JobsClientListBySubscriptionOptions`
- New struct `JobsClientListBySubscriptionPager`
- New struct `JobsClientListBySubscriptionResponse`
- New struct `JobsClientListBySubscriptionResult`
- New struct `JobsClientUpdateOptions`
- New struct `JobsClientUpdateResponse`
- New struct `JobsClientUpdateResult`
- New struct `LocationsClientGetOptions`
- New struct `LocationsClientGetResponse`
- New struct `LocationsClientGetResult`
- New struct `LocationsClientListOptions`
- New struct `LocationsClientListResponse`
- New struct `LocationsClientListResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New field `Error` in struct `ErrorResponse`


## 0.1.0 (2021-12-22)

- Init release.
