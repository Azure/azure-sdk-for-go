# Release History

## 1.0.0 (2022-05-18)
### Breaking Changes

- Function `*OpenShiftClustersClient.BeginDelete` return value(s) have been changed from `(*armruntime.Poller[OpenShiftClustersClientDeleteResponse], error)` to `(*runtime.Poller[OpenShiftClustersClientDeleteResponse], error)`
- Function `*OpenShiftClustersClient.BeginUpdate` return value(s) have been changed from `(*armruntime.Poller[OpenShiftClustersClientUpdateResponse], error)` to `(*runtime.Poller[OpenShiftClustersClientUpdateResponse], error)`
- Function `*OpenShiftClustersClient.BeginCreateOrUpdate` return value(s) have been changed from `(*armruntime.Poller[OpenShiftClustersClientCreateOrUpdateResponse], error)` to `(*runtime.Poller[OpenShiftClustersClientCreateOrUpdateResponse], error)`
- Type of `WorkerProfile.VMSize` has been changed from `*VMSize` to `*string`
- Type of `MasterProfile.VMSize` has been changed from `*VMSize` to `*string`
- Const `VMSizeStandardD4SV3` has been removed
- Const `VMSizeStandardD8SV3` has been removed
- Const `VMSizeStandardD2SV3` has been removed
- Function `CloudErrorBody.MarshalJSON` has been removed
- Function `PossibleVMSizeValues` has been removed
- Function `OperationList.MarshalJSON` has been removed
- Function `OpenShiftClusterList.MarshalJSON` has been removed

### Features Added

- New const `CreatedByTypeApplication`
- New const `FipsValidatedModulesEnabled`
- New const `CreatedByTypeKey`
- New const `CreatedByTypeUser`
- New const `EncryptionAtHostDisabled`
- New const `EncryptionAtHostEnabled`
- New const `CreatedByTypeManagedIdentity`
- New const `FipsValidatedModulesDisabled`
- New function `*OpenShiftClustersClient.ListAdminCredentials(context.Context, string, string, *OpenShiftClustersClientListAdminCredentialsOptions) (OpenShiftClustersClientListAdminCredentialsResponse, error)`
- New function `SystemData.MarshalJSON() ([]byte, error)`
- New function `PossibleFipsValidatedModulesValues() []FipsValidatedModules`
- New function `PossibleCreatedByTypeValues() []CreatedByType`
- New function `*SystemData.UnmarshalJSON([]byte) error`
- New function `PossibleEncryptionAtHostValues() []EncryptionAtHost`
- New struct `OpenShiftClusterAdminKubeconfig`
- New struct `OpenShiftClustersClientListAdminCredentialsOptions`
- New struct `OpenShiftClustersClientListAdminCredentialsResponse`
- New struct `SystemData`
- New field `SystemData` in struct `OpenShiftClusterUpdate`
- New field `SystemData` in struct `OpenShiftCluster`
- New field `DiskEncryptionSetID` in struct `MasterProfile`
- New field `EncryptionAtHost` in struct `MasterProfile`
- New field `FipsValidatedModules` in struct `ClusterProfile`
- New field `EncryptionAtHost` in struct `WorkerProfile`
- New field `DiskEncryptionSetID` in struct `WorkerProfile`


## 0.4.0 (2022-04-18)
### Breaking Changes

- Function `*OpenShiftClustersClient.List` has been removed
- Function `*OpenShiftClustersClient.ListByResourceGroup` has been removed
- Function `*OperationsClient.List` has been removed

### Features Added

- New function `*OpenShiftClustersClient.NewListByResourceGroupPager(string, *OpenShiftClustersClientListByResourceGroupOptions) *runtime.Pager[OpenShiftClustersClientListByResourceGroupResponse]`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New function `*OpenShiftClustersClient.NewListPager(*OpenShiftClustersClientListOptions) *runtime.Pager[OpenShiftClustersClientListResponse]`


## 0.3.0 (2022-04-12)
### Breaking Changes

- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `*OpenShiftClustersClient.BeginCreateOrUpdate` return value(s) have been changed from `(OpenShiftClustersClientCreateOrUpdatePollerResponse, error)` to `(*armruntime.Poller[OpenShiftClustersClientCreateOrUpdateResponse], error)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `NewOpenShiftClustersClient` return value(s) have been changed from `(*OpenShiftClustersClient)` to `(*OpenShiftClustersClient, error)`
- Function `*OpenShiftClustersClient.ListByResourceGroup` return value(s) have been changed from `(*OpenShiftClustersClientListByResourceGroupPager)` to `(*runtime.Pager[OpenShiftClustersClientListByResourceGroupResponse])`
- Function `*OpenShiftClustersClient.BeginUpdate` return value(s) have been changed from `(OpenShiftClustersClientUpdatePollerResponse, error)` to `(*armruntime.Poller[OpenShiftClustersClientUpdateResponse], error)`
- Function `*OpenShiftClustersClient.BeginDelete` return value(s) have been changed from `(OpenShiftClustersClientDeletePollerResponse, error)` to `(*armruntime.Poller[OpenShiftClustersClientDeleteResponse], error)`
- Function `*OpenShiftClustersClient.List` return value(s) have been changed from `(*OpenShiftClustersClientListPager)` to `(*runtime.Pager[OpenShiftClustersClientListResponse])`
- Function `*OpenShiftClustersClientListByResourceGroupPager.PageResponse` has been removed
- Function `OpenShiftClustersClientDeletePollerResponse.PollUntilDone` has been removed
- Function `*OpenShiftClustersClientListByResourceGroupPager.Err` has been removed
- Function `*OpenShiftClustersClientListPager.PageResponse` has been removed
- Function `*OpenShiftClustersClientDeletePollerResponse.Resume` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Function `*OpenShiftClustersClientCreateOrUpdatePoller.FinalResponse` has been removed
- Function `*OpenShiftClustersClientListPager.NextPage` has been removed
- Function `*OpenShiftClustersClientUpdatePoller.ResumeToken` has been removed
- Function `*OpenShiftClustersClientDeletePoller.FinalResponse` has been removed
- Function `*OperationsClientListPager.Err` has been removed
- Function `VMSize.ToPtr` has been removed
- Function `*OpenShiftClustersClientDeletePoller.Done` has been removed
- Function `*OperationsClientListPager.NextPage` has been removed
- Function `Visibility.ToPtr` has been removed
- Function `*OpenShiftClustersClientDeletePoller.Poll` has been removed
- Function `*OpenShiftClustersClientDeletePoller.ResumeToken` has been removed
- Function `OpenShiftClustersClientCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `OpenShiftClustersClientUpdatePollerResponse.PollUntilDone` has been removed
- Function `*OpenShiftClustersClientCreateOrUpdatePollerResponse.Resume` has been removed
- Function `*OpenShiftClustersClientListByResourceGroupPager.NextPage` has been removed
- Function `*OpenShiftClustersClientCreateOrUpdatePoller.ResumeToken` has been removed
- Function `ProvisioningState.ToPtr` has been removed
- Function `*OpenShiftClustersClientCreateOrUpdatePoller.Done` has been removed
- Function `*OpenShiftClustersClientUpdatePoller.FinalResponse` has been removed
- Function `*OpenShiftClustersClientCreateOrUpdatePoller.Poll` has been removed
- Function `*OpenShiftClustersClientListPager.Err` has been removed
- Function `*OpenShiftClustersClientUpdatePoller.Done` has been removed
- Function `*OpenShiftClustersClientUpdatePollerResponse.Resume` has been removed
- Function `*OpenShiftClustersClientUpdatePoller.Poll` has been removed
- Struct `OpenShiftClustersClientCreateOrUpdatePoller` has been removed
- Struct `OpenShiftClustersClientCreateOrUpdatePollerResponse` has been removed
- Struct `OpenShiftClustersClientCreateOrUpdateResult` has been removed
- Struct `OpenShiftClustersClientDeletePoller` has been removed
- Struct `OpenShiftClustersClientDeletePollerResponse` has been removed
- Struct `OpenShiftClustersClientGetResult` has been removed
- Struct `OpenShiftClustersClientListByResourceGroupPager` has been removed
- Struct `OpenShiftClustersClientListByResourceGroupResult` has been removed
- Struct `OpenShiftClustersClientListCredentialsResult` has been removed
- Struct `OpenShiftClustersClientListPager` has been removed
- Struct `OpenShiftClustersClientListResult` has been removed
- Struct `OpenShiftClustersClientUpdatePoller` has been removed
- Struct `OpenShiftClustersClientUpdatePollerResponse` has been removed
- Struct `OpenShiftClustersClientUpdateResult` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Field `RawResponse` of struct `OpenShiftClustersClientDeleteResponse` has been removed
- Field `OpenShiftClustersClientListCredentialsResult` of struct `OpenShiftClustersClientListCredentialsResponse` has been removed
- Field `RawResponse` of struct `OpenShiftClustersClientListCredentialsResponse` has been removed
- Field `OpenShiftClustersClientCreateOrUpdateResult` of struct `OpenShiftClustersClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `OpenShiftClustersClientCreateOrUpdateResponse` has been removed
- Field `OpenShiftClustersClientGetResult` of struct `OpenShiftClustersClientGetResponse` has been removed
- Field `RawResponse` of struct `OpenShiftClustersClientGetResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `OpenShiftClustersClientUpdateResult` of struct `OpenShiftClustersClientUpdateResponse` has been removed
- Field `RawResponse` of struct `OpenShiftClustersClientUpdateResponse` has been removed
- Field `OpenShiftClustersClientListByResourceGroupResult` of struct `OpenShiftClustersClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `OpenShiftClustersClientListByResourceGroupResponse` has been removed
- Field `OpenShiftClustersClientListResult` of struct `OpenShiftClustersClientListResponse` has been removed
- Field `RawResponse` of struct `OpenShiftClustersClientListResponse` has been removed

### Features Added

- New field `ResumeToken` in struct `OpenShiftClustersClientBeginCreateOrUpdateOptions`
- New anonymous field `OpenShiftCluster` in struct `OpenShiftClustersClientCreateOrUpdateResponse`
- New anonymous field `OpenShiftClusterCredentials` in struct `OpenShiftClustersClientListCredentialsResponse`
- New anonymous field `OpenShiftClusterList` in struct `OpenShiftClustersClientListResponse`
- New anonymous field `OpenShiftClusterList` in struct `OpenShiftClustersClientListByResourceGroupResponse`
- New field `ResumeToken` in struct `OpenShiftClustersClientBeginDeleteOptions`
- New anonymous field `OpenShiftCluster` in struct `OpenShiftClustersClientGetResponse`
- New anonymous field `OperationList` in struct `OperationsClientListResponse`
- New field `ResumeToken` in struct `OpenShiftClustersClientBeginUpdateOptions`
- New anonymous field `OpenShiftCluster` in struct `OpenShiftClustersClientUpdateResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*OpenShiftClustersClient.Get` parameter(s) have been changed from `(context.Context, string, string, *OpenShiftClustersGetOptions)` to `(context.Context, string, string, *OpenShiftClustersClientGetOptions)`
- Function `*OpenShiftClustersClient.Get` return value(s) have been changed from `(OpenShiftClustersGetResponse, error)` to `(OpenShiftClustersClientGetResponse, error)`
- Function `*OpenShiftClustersClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, OpenShiftCluster, *OpenShiftClustersBeginCreateOrUpdateOptions)` to `(context.Context, string, string, OpenShiftCluster, *OpenShiftClustersClientBeginCreateOrUpdateOptions)`
- Function `*OpenShiftClustersClient.BeginCreateOrUpdate` return value(s) have been changed from `(OpenShiftClustersCreateOrUpdatePollerResponse, error)` to `(OpenShiftClustersClientCreateOrUpdatePollerResponse, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*OpenShiftClustersClient.ListByResourceGroup` parameter(s) have been changed from `(string, *OpenShiftClustersListByResourceGroupOptions)` to `(string, *OpenShiftClustersClientListByResourceGroupOptions)`
- Function `*OpenShiftClustersClient.ListByResourceGroup` return value(s) have been changed from `(*OpenShiftClustersListByResourceGroupPager)` to `(*OpenShiftClustersClientListByResourceGroupPager)`
- Function `*OpenShiftClustersClient.ListCredentials` parameter(s) have been changed from `(context.Context, string, string, *OpenShiftClustersListCredentialsOptions)` to `(context.Context, string, string, *OpenShiftClustersClientListCredentialsOptions)`
- Function `*OpenShiftClustersClient.ListCredentials` return value(s) have been changed from `(OpenShiftClustersListCredentialsResponse, error)` to `(OpenShiftClustersClientListCredentialsResponse, error)`
- Function `*OpenShiftClustersClient.List` parameter(s) have been changed from `(*OpenShiftClustersListOptions)` to `(*OpenShiftClustersClientListOptions)`
- Function `*OpenShiftClustersClient.List` return value(s) have been changed from `(*OpenShiftClustersListPager)` to `(*OpenShiftClustersClientListPager)`
- Function `*OpenShiftClustersClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, OpenShiftClusterUpdate, *OpenShiftClustersBeginUpdateOptions)` to `(context.Context, string, string, OpenShiftClusterUpdate, *OpenShiftClustersClientBeginUpdateOptions)`
- Function `*OpenShiftClustersClient.BeginUpdate` return value(s) have been changed from `(OpenShiftClustersUpdatePollerResponse, error)` to `(OpenShiftClustersClientUpdatePollerResponse, error)`
- Function `*OpenShiftClustersClient.BeginDelete` parameter(s) have been changed from `(context.Context, string, string, *OpenShiftClustersBeginDeleteOptions)` to `(context.Context, string, string, *OpenShiftClustersClientBeginDeleteOptions)`
- Function `*OpenShiftClustersClient.BeginDelete` return value(s) have been changed from `(OpenShiftClustersDeletePollerResponse, error)` to `(OpenShiftClustersClientDeletePollerResponse, error)`
- Function `*OpenShiftClustersListPager.PageResponse` has been removed
- Function `OpenShiftClustersUpdatePollerResponse.PollUntilDone` has been removed
- Function `*OpenShiftClustersDeletePoller.FinalResponse` has been removed
- Function `*OpenShiftClustersUpdatePoller.FinalResponse` has been removed
- Function `*OpenShiftClustersCreateOrUpdatePoller.Poll` has been removed
- Function `*OpenShiftClustersDeletePoller.ResumeToken` has been removed
- Function `*OpenShiftClustersListPager.NextPage` has been removed
- Function `*OpenShiftClustersCreateOrUpdatePoller.Done` has been removed
- Function `OpenShiftClustersDeletePollerResponse.PollUntilDone` has been removed
- Function `*OpenShiftClustersListPager.Err` has been removed
- Function `*OpenShiftClustersCreateOrUpdatePollerResponse.Resume` has been removed
- Function `*OpenShiftClustersUpdatePoller.ResumeToken` has been removed
- Function `*OpenShiftClustersListByResourceGroupPager.NextPage` has been removed
- Function `*OpenShiftClustersListByResourceGroupPager.PageResponse` has been removed
- Function `Resource.MarshalJSON` has been removed
- Function `*OpenShiftClustersDeletePollerResponse.Resume` has been removed
- Function `*OpenShiftClustersUpdatePollerResponse.Resume` has been removed
- Function `CloudError.Error` has been removed
- Function `*OpenShiftClustersCreateOrUpdatePoller.ResumeToken` has been removed
- Function `*OperationsListPager.PageResponse` has been removed
- Function `OpenShiftClustersCreateOrUpdatePollerResponse.PollUntilDone` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Function `*OpenShiftClustersDeletePoller.Done` has been removed
- Function `*OpenShiftClustersCreateOrUpdatePoller.FinalResponse` has been removed
- Function `*OpenShiftClustersDeletePoller.Poll` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `*OpenShiftClustersUpdatePoller.Done` has been removed
- Function `*OpenShiftClustersListByResourceGroupPager.Err` has been removed
- Function `*OpenShiftClustersUpdatePoller.Poll` has been removed
- Struct `OpenShiftClustersBeginCreateOrUpdateOptions` has been removed
- Struct `OpenShiftClustersBeginDeleteOptions` has been removed
- Struct `OpenShiftClustersBeginUpdateOptions` has been removed
- Struct `OpenShiftClustersCreateOrUpdatePoller` has been removed
- Struct `OpenShiftClustersCreateOrUpdatePollerResponse` has been removed
- Struct `OpenShiftClustersCreateOrUpdateResponse` has been removed
- Struct `OpenShiftClustersCreateOrUpdateResult` has been removed
- Struct `OpenShiftClustersDeletePoller` has been removed
- Struct `OpenShiftClustersDeletePollerResponse` has been removed
- Struct `OpenShiftClustersDeleteResponse` has been removed
- Struct `OpenShiftClustersGetOptions` has been removed
- Struct `OpenShiftClustersGetResponse` has been removed
- Struct `OpenShiftClustersGetResult` has been removed
- Struct `OpenShiftClustersListByResourceGroupOptions` has been removed
- Struct `OpenShiftClustersListByResourceGroupPager` has been removed
- Struct `OpenShiftClustersListByResourceGroupResponse` has been removed
- Struct `OpenShiftClustersListByResourceGroupResult` has been removed
- Struct `OpenShiftClustersListCredentialsOptions` has been removed
- Struct `OpenShiftClustersListCredentialsResponse` has been removed
- Struct `OpenShiftClustersListCredentialsResult` has been removed
- Struct `OpenShiftClustersListOptions` has been removed
- Struct `OpenShiftClustersListPager` has been removed
- Struct `OpenShiftClustersListResponse` has been removed
- Struct `OpenShiftClustersListResult` has been removed
- Struct `OpenShiftClustersUpdatePoller` has been removed
- Struct `OpenShiftClustersUpdatePollerResponse` has been removed
- Struct `OpenShiftClustersUpdateResponse` has been removed
- Struct `OpenShiftClustersUpdateResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Field `TrackedResource` of struct `OpenShiftCluster` has been removed
- Field `Resource` of struct `TrackedResource` has been removed
- Field `InnerError` of struct `CloudError` has been removed

### Features Added

- New function `*OpenShiftClustersClientListPager.PageResponse() OpenShiftClustersClientListResponse`
- New function `*OpenShiftClustersClientListByResourceGroupPager.PageResponse() OpenShiftClustersClientListByResourceGroupResponse`
- New function `*OpenShiftClustersClientCreateOrUpdatePoller.Done() bool`
- New function `*OpenShiftClustersClientDeletePoller.Poll(context.Context) (*http.Response, error)`
- New function `*OpenShiftClustersClientCreateOrUpdatePollerResponse.Resume(context.Context, *OpenShiftClustersClient, string) error`
- New function `*OpenShiftClustersClientUpdatePoller.Done() bool`
- New function `OpenShiftClustersClientCreateOrUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (OpenShiftClustersClientCreateOrUpdateResponse, error)`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*OpenShiftClustersClientDeletePollerResponse.Resume(context.Context, *OpenShiftClustersClient, string) error`
- New function `*OpenShiftClustersClientDeletePoller.Done() bool`
- New function `*OperationsClientListPager.Err() error`
- New function `*OpenShiftClustersClientDeletePoller.ResumeToken() (string, error)`
- New function `*OpenShiftClustersClientUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*OpenShiftClustersClientUpdatePoller.ResumeToken() (string, error)`
- New function `*OpenShiftClustersClientUpdatePoller.FinalResponse(context.Context) (OpenShiftClustersClientUpdateResponse, error)`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*OpenShiftClustersClientListPager.Err() error`
- New function `OpenShiftClustersClientUpdatePollerResponse.PollUntilDone(context.Context, time.Duration) (OpenShiftClustersClientUpdateResponse, error)`
- New function `*OpenShiftClustersClientDeletePoller.FinalResponse(context.Context) (OpenShiftClustersClientDeleteResponse, error)`
- New function `*OpenShiftClustersClientCreateOrUpdatePoller.FinalResponse(context.Context) (OpenShiftClustersClientCreateOrUpdateResponse, error)`
- New function `*OpenShiftClustersClientListPager.NextPage(context.Context) bool`
- New function `OpenShiftClustersClientDeletePollerResponse.PollUntilDone(context.Context, time.Duration) (OpenShiftClustersClientDeleteResponse, error)`
- New function `*OpenShiftClustersClientCreateOrUpdatePoller.Poll(context.Context) (*http.Response, error)`
- New function `*OpenShiftClustersClientCreateOrUpdatePoller.ResumeToken() (string, error)`
- New function `*OpenShiftClustersClientUpdatePollerResponse.Resume(context.Context, *OpenShiftClustersClient, string) error`
- New function `*OpenShiftClustersClientListByResourceGroupPager.Err() error`
- New function `*OpenShiftClustersClientListByResourceGroupPager.NextPage(context.Context) bool`
- New struct `OpenShiftClustersClientBeginCreateOrUpdateOptions`
- New struct `OpenShiftClustersClientBeginDeleteOptions`
- New struct `OpenShiftClustersClientBeginUpdateOptions`
- New struct `OpenShiftClustersClientCreateOrUpdatePoller`
- New struct `OpenShiftClustersClientCreateOrUpdatePollerResponse`
- New struct `OpenShiftClustersClientCreateOrUpdateResponse`
- New struct `OpenShiftClustersClientCreateOrUpdateResult`
- New struct `OpenShiftClustersClientDeletePoller`
- New struct `OpenShiftClustersClientDeletePollerResponse`
- New struct `OpenShiftClustersClientDeleteResponse`
- New struct `OpenShiftClustersClientGetOptions`
- New struct `OpenShiftClustersClientGetResponse`
- New struct `OpenShiftClustersClientGetResult`
- New struct `OpenShiftClustersClientListByResourceGroupOptions`
- New struct `OpenShiftClustersClientListByResourceGroupPager`
- New struct `OpenShiftClustersClientListByResourceGroupResponse`
- New struct `OpenShiftClustersClientListByResourceGroupResult`
- New struct `OpenShiftClustersClientListCredentialsOptions`
- New struct `OpenShiftClustersClientListCredentialsResponse`
- New struct `OpenShiftClustersClientListCredentialsResult`
- New struct `OpenShiftClustersClientListOptions`
- New struct `OpenShiftClustersClientListPager`
- New struct `OpenShiftClustersClientListResponse`
- New struct `OpenShiftClustersClientListResult`
- New struct `OpenShiftClustersClientUpdatePoller`
- New struct `OpenShiftClustersClientUpdatePollerResponse`
- New struct `OpenShiftClustersClientUpdateResponse`
- New struct `OpenShiftClustersClientUpdateResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New field `ID` in struct `TrackedResource`
- New field `Name` in struct `TrackedResource`
- New field `Type` in struct `TrackedResource`
- New field `Location` in struct `OpenShiftCluster`
- New field `Tags` in struct `OpenShiftCluster`
- New field `ID` in struct `OpenShiftCluster`
- New field `Name` in struct `OpenShiftCluster`
- New field `Type` in struct `OpenShiftCluster`
- New field `Error` in struct `CloudError`


## 0.1.0 (2021-12-09)

- Init release.
