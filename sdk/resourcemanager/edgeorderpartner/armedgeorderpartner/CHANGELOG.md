# Release History

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
