# Release History

## 0.5.0 (2022-05-18)
### Breaking Changes

- Function `*APISClient.BeginManageInventoryMetadata` return value(s) have been changed from `(*armruntime.Poller[APISClientManageInventoryMetadataResponse], error)` to `(*runtime.Poller[APISClientManageInventoryMetadataResponse], error)`
- Function `InventoryAdditionalDetails.MarshalJSON` has been removed
- Function `OperationListResult.MarshalJSON` has been removed
- Function `AdditionalInventoryDetails.MarshalJSON` has been removed
- Function `ConfigurationDetails.MarshalJSON` has been removed
- Function `PartnerInventoryList.MarshalJSON` has been removed
- Function `ErrorDetail.MarshalJSON` has been removed
- Function `StageDetails.MarshalJSON` has been removed
- Function `CloudError.MarshalJSON` has been removed


## 0.4.0 (2022-04-15)
### Breaking Changes

- Function `*APISClient.ListOperationsPartner` has been removed
- Function `*APISClient.SearchInventories` has been removed

### Features Added

- New function `*APISClient.NewListOperationsPartnerPager(*APISClientListOperationsPartnerOptions) *runtime.Pager[APISClientListOperationsPartnerResponse]`
- New function `*APISClient.NewSearchInventoriesPager(SearchInventoriesRequest, *APISClientSearchInventoriesOptions) *runtime.Pager[APISClientSearchInventoriesResponse]`


## 0.3.0 (2022-04-11)
### Breaking Changes

- Function `NewAPISClient` return value(s) have been changed from `(*APISClient)` to `(*APISClient, error)`
- Function `*APISClient.BeginManageInventoryMetadata` return value(s) have been changed from `(APISClientManageInventoryMetadataPollerResponse, error)` to `(*armruntime.Poller[APISClientManageInventoryMetadataResponse], error)`
- Function `*APISClient.SearchInventories` return value(s) have been changed from `(*APISClientSearchInventoriesPager)` to `(*runtime.Pager[APISClientSearchInventoriesResponse])`
- Function `*APISClient.ListOperationsPartner` return value(s) have been changed from `(*APISClientListOperationsPartnerPager)` to `(*runtime.Pager[APISClientListOperationsPartnerResponse])`
- Type of `ErrorAdditionalInfo.Info` has been changed from `map[string]interface{}` to `interface{}`
- Type of `AdditionalErrorInfo.Info` has been changed from `map[string]interface{}` to `interface{}`
- Function `*APISClientManageInventoryMetadataPoller.Done` has been removed
- Function `ActionType.ToPtr` has been removed
- Function `ManageLinkOperation.ToPtr` has been removed
- Function `*APISClientSearchInventoriesPager.Err` has been removed
- Function `*APISClientManageInventoryMetadataPoller.Poll` has been removed
- Function `*APISClientListOperationsPartnerPager.NextPage` has been removed
- Function `StageStatus.ToPtr` has been removed
- Function `*APISClientManageInventoryMetadataPollerResponse.Resume` has been removed
- Function `Origin.ToPtr` has been removed
- Function `OrderItemType.ToPtr` has been removed
- Function `*APISClientSearchInventoriesPager.PageResponse` has been removed
- Function `APISClientManageInventoryMetadataPollerResponse.PollUntilDone` has been removed
- Function `*APISClientSearchInventoriesPager.NextPage` has been removed
- Function `*APISClientManageInventoryMetadataPoller.FinalResponse` has been removed
- Function `*APISClientListOperationsPartnerPager.Err` has been removed
- Function `StageName.ToPtr` has been removed
- Function `*APISClientManageInventoryMetadataPoller.ResumeToken` has been removed
- Function `*APISClientListOperationsPartnerPager.PageResponse` has been removed
- Struct `APISClientListOperationsPartnerPager` has been removed
- Struct `APISClientListOperationsPartnerResult` has been removed
- Struct `APISClientManageInventoryMetadataPoller` has been removed
- Struct `APISClientManageInventoryMetadataPollerResponse` has been removed
- Struct `APISClientSearchInventoriesPager` has been removed
- Struct `APISClientSearchInventoriesResult` has been removed
- Field `APISClientSearchInventoriesResult` of struct `APISClientSearchInventoriesResponse` has been removed
- Field `RawResponse` of struct `APISClientSearchInventoriesResponse` has been removed
- Field `RawResponse` of struct `APISClientManageInventoryMetadataResponse` has been removed
- Field `RawResponse` of struct `APISClientManageLinkResponse` has been removed
- Field `APISClientListOperationsPartnerResult` of struct `APISClientListOperationsPartnerResponse` has been removed
- Field `RawResponse` of struct `APISClientListOperationsPartnerResponse` has been removed

### Features Added

- New field `ResumeToken` in struct `APISClientBeginManageInventoryMetadataOptions`
- New anonymous field `OperationListResult` in struct `APISClientListOperationsPartnerResponse`
- New anonymous field `PartnerInventoryList` in struct `APISClientSearchInventoriesResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*EdgeOrderPartnerAPISClient.ListOperationsPartner` has been removed
- Function `EdgeOrderPartnerAPISManageInventoryMetadataPollerResponse.PollUntilDone` has been removed
- Function `*EdgeOrderPartnerAPISManageInventoryMetadataPollerResponse.Resume` has been removed
- Function `*EdgeOrderPartnerAPISClient.BeginManageInventoryMetadata` has been removed
- Function `*EdgeOrderPartnerAPISListOperationsPartnerPager.Err` has been removed
- Function `*EdgeOrderPartnerAPISSearchInventoriesPager.Err` has been removed
- Function `*EdgeOrderPartnerAPISSearchInventoriesPager.NextPage` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `*EdgeOrderPartnerAPISClient.SearchInventories` has been removed
- Function `*EdgeOrderPartnerAPISListOperationsPartnerPager.NextPage` has been removed
- Function `*EdgeOrderPartnerAPISSearchInventoriesPager.PageResponse` has been removed
- Function `*EdgeOrderPartnerAPISManageInventoryMetadataPoller.FinalResponse` has been removed
- Function `*EdgeOrderPartnerAPISManageInventoryMetadataPoller.Done` has been removed
- Function `*EdgeOrderPartnerAPISListOperationsPartnerPager.PageResponse` has been removed
- Function `*EdgeOrderPartnerAPISManageInventoryMetadataPoller.Poll` has been removed
- Function `*EdgeOrderPartnerAPISManageInventoryMetadataPoller.ResumeToken` has been removed
- Function `NewEdgeOrderPartnerAPISClient` has been removed
- Function `*EdgeOrderPartnerAPISClient.ManageLink` has been removed
- Struct `EdgeOrderPartnerAPISBeginManageInventoryMetadataOptions` has been removed
- Struct `EdgeOrderPartnerAPISClient` has been removed
- Struct `EdgeOrderPartnerAPISListOperationsPartnerOptions` has been removed
- Struct `EdgeOrderPartnerAPISListOperationsPartnerPager` has been removed
- Struct `EdgeOrderPartnerAPISListOperationsPartnerResponse` has been removed
- Struct `EdgeOrderPartnerAPISListOperationsPartnerResult` has been removed
- Struct `EdgeOrderPartnerAPISManageInventoryMetadataPoller` has been removed
- Struct `EdgeOrderPartnerAPISManageInventoryMetadataPollerResponse` has been removed
- Struct `EdgeOrderPartnerAPISManageInventoryMetadataResponse` has been removed
- Struct `EdgeOrderPartnerAPISManageLinkOptions` has been removed
- Struct `EdgeOrderPartnerAPISManageLinkResponse` has been removed
- Struct `EdgeOrderPartnerAPISSearchInventoriesOptions` has been removed
- Struct `EdgeOrderPartnerAPISSearchInventoriesPager` has been removed
- Struct `EdgeOrderPartnerAPISSearchInventoriesResponse` has been removed
- Struct `EdgeOrderPartnerAPISSearchInventoriesResult` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed

### Features Added

- New function `NewAPISClient(string, azcore.TokenCredential, *arm.ClientOptions) *APISClient`
- New function `*APISClientListOperationsPartnerPager.Err() error`
- New function `*APISClientManageInventoryMetadataPoller.ResumeToken() (string, error)`
- New function `*APISClientManageInventoryMetadataPoller.Done() bool`
- New function `*APISClientSearchInventoriesPager.NextPage(context.Context) bool`
- New function `*APISClientManageInventoryMetadataPoller.FinalResponse(context.Context) (APISClientManageInventoryMetadataResponse, error)`
- New function `*APISClientManageInventoryMetadataPoller.Poll(context.Context) (*http.Response, error)`
- New function `*APISClientSearchInventoriesPager.PageResponse() APISClientSearchInventoriesResponse`
- New function `APISClientManageInventoryMetadataPollerResponse.PollUntilDone(context.Context, time.Duration) (APISClientManageInventoryMetadataResponse, error)`
- New function `*APISClient.BeginManageInventoryMetadata(context.Context, string, string, string, ManageInventoryMetadataRequest, *APISClientBeginManageInventoryMetadataOptions) (APISClientManageInventoryMetadataPollerResponse, error)`
- New function `*APISClientManageInventoryMetadataPollerResponse.Resume(context.Context, *APISClient, string) error`
- New function `*APISClient.ManageLink(context.Context, string, string, string, ManageLinkRequest, *APISClientManageLinkOptions) (APISClientManageLinkResponse, error)`
- New function `*APISClient.SearchInventories(SearchInventoriesRequest, *APISClientSearchInventoriesOptions) *APISClientSearchInventoriesPager`
- New function `*APISClientSearchInventoriesPager.Err() error`
- New function `*APISClientListOperationsPartnerPager.NextPage(context.Context) bool`
- New function `*APISClient.ListOperationsPartner(*APISClientListOperationsPartnerOptions) *APISClientListOperationsPartnerPager`
- New function `*APISClientListOperationsPartnerPager.PageResponse() APISClientListOperationsPartnerResponse`
- New struct `APISClient`
- New struct `APISClientBeginManageInventoryMetadataOptions`
- New struct `APISClientListOperationsPartnerOptions`
- New struct `APISClientListOperationsPartnerPager`
- New struct `APISClientListOperationsPartnerResponse`
- New struct `APISClientListOperationsPartnerResult`
- New struct `APISClientManageInventoryMetadataPoller`
- New struct `APISClientManageInventoryMetadataPollerResponse`
- New struct `APISClientManageInventoryMetadataResponse`
- New struct `APISClientManageLinkOptions`
- New struct `APISClientManageLinkResponse`
- New struct `APISClientSearchInventoriesOptions`
- New struct `APISClientSearchInventoriesPager`
- New struct `APISClientSearchInventoriesResponse`
- New struct `APISClientSearchInventoriesResult`
- New field `Error` in struct `ErrorResponse`


## 0.1.0 (2021-12-07)

- Init release.
