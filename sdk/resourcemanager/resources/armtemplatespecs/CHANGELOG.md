# Release History

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*TemplateSpecVersionsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, TemplateSpecVersion, *TemplateSpecVersionsCreateOrUpdateOptions)` to `(context.Context, string, string, string, TemplateSpecVersion, *TemplateSpecVersionsClientCreateOrUpdateOptions)`
- Function `*TemplateSpecVersionsClient.CreateOrUpdate` return value(s) have been changed from `(TemplateSpecVersionsCreateOrUpdateResponse, error)` to `(TemplateSpecVersionsClientCreateOrUpdateResponse, error)`
- Function `*TemplateSpecVersionsClient.Update` parameter(s) have been changed from `(context.Context, string, string, string, *TemplateSpecVersionsUpdateOptions)` to `(context.Context, string, string, string, *TemplateSpecVersionsClientUpdateOptions)`
- Function `*TemplateSpecVersionsClient.Update` return value(s) have been changed from `(TemplateSpecVersionsUpdateResponse, error)` to `(TemplateSpecVersionsClientUpdateResponse, error)`
- Function `*TemplateSpecVersionsClient.List` parameter(s) have been changed from `(string, string, *TemplateSpecVersionsListOptions)` to `(string, string, *TemplateSpecVersionsClientListOptions)`
- Function `*TemplateSpecVersionsClient.List` return value(s) have been changed from `(*TemplateSpecVersionsListPager)` to `(*TemplateSpecVersionsClientListPager)`
- Function `*TemplateSpecVersionsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, string, *TemplateSpecVersionsDeleteOptions)` to `(context.Context, string, string, string, *TemplateSpecVersionsClientDeleteOptions)`
- Function `*TemplateSpecVersionsClient.Delete` return value(s) have been changed from `(TemplateSpecVersionsDeleteResponse, error)` to `(TemplateSpecVersionsClientDeleteResponse, error)`
- Function `*TemplateSpecVersionsClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *TemplateSpecVersionsGetOptions)` to `(context.Context, string, string, string, *TemplateSpecVersionsClientGetOptions)`
- Function `*TemplateSpecVersionsClient.Get` return value(s) have been changed from `(TemplateSpecVersionsGetResponse, error)` to `(TemplateSpecVersionsClientGetResponse, error)`
- Function `AzureResourceBase.MarshalJSON` has been removed
- Function `*TemplateSpecsClient.Update` has been removed
- Function `*TemplateSpecVersionsListPager.PageResponse` has been removed
- Function `*TemplateSpecsListBySubscriptionPager.PageResponse` has been removed
- Function `*TemplateSpecsClient.Delete` has been removed
- Function `*TemplateSpecsClient.Get` has been removed
- Function `*TemplateSpecsListBySubscriptionPager.NextPage` has been removed
- Function `*TemplateSpecsClient.ListBySubscription` has been removed
- Function `*TemplateSpecVersionsListPager.NextPage` has been removed
- Function `*TemplateSpecsListByResourceGroupPager.PageResponse` has been removed
- Function `*TemplateSpecsClient.ListByResourceGroup` has been removed
- Function `TemplateSpecsError.Error` has been removed
- Function `*TemplateSpecsListBySubscriptionPager.Err` has been removed
- Function `*TemplateSpecsClient.CreateOrUpdate` has been removed
- Function `*TemplateSpecsListByResourceGroupPager.NextPage` has been removed
- Function `TemplateSpecsListResult.MarshalJSON` has been removed
- Function `*TemplateSpecsListByResourceGroupPager.Err` has been removed
- Function `*TemplateSpecVersionsListPager.Err` has been removed
- Function `NewTemplateSpecsClient` has been removed
- Struct `TemplateSpecVersionsCreateOrUpdateOptions` has been removed
- Struct `TemplateSpecVersionsCreateOrUpdateResponse` has been removed
- Struct `TemplateSpecVersionsCreateOrUpdateResult` has been removed
- Struct `TemplateSpecVersionsDeleteOptions` has been removed
- Struct `TemplateSpecVersionsDeleteResponse` has been removed
- Struct `TemplateSpecVersionsGetOptions` has been removed
- Struct `TemplateSpecVersionsGetResponse` has been removed
- Struct `TemplateSpecVersionsGetResult` has been removed
- Struct `TemplateSpecVersionsListOptions` has been removed
- Struct `TemplateSpecVersionsListPager` has been removed
- Struct `TemplateSpecVersionsListResponse` has been removed
- Struct `TemplateSpecVersionsListResultEnvelope` has been removed
- Struct `TemplateSpecVersionsUpdateOptions` has been removed
- Struct `TemplateSpecVersionsUpdateResponse` has been removed
- Struct `TemplateSpecVersionsUpdateResult` has been removed
- Struct `TemplateSpecsClient` has been removed
- Struct `TemplateSpecsCreateOrUpdateOptions` has been removed
- Struct `TemplateSpecsCreateOrUpdateResponse` has been removed
- Struct `TemplateSpecsCreateOrUpdateResult` has been removed
- Struct `TemplateSpecsDeleteOptions` has been removed
- Struct `TemplateSpecsDeleteResponse` has been removed
- Struct `TemplateSpecsError` has been removed
- Struct `TemplateSpecsGetOptions` has been removed
- Struct `TemplateSpecsGetResponse` has been removed
- Struct `TemplateSpecsGetResult` has been removed
- Struct `TemplateSpecsListByResourceGroupOptions` has been removed
- Struct `TemplateSpecsListByResourceGroupPager` has been removed
- Struct `TemplateSpecsListByResourceGroupResponse` has been removed
- Struct `TemplateSpecsListByResourceGroupResult` has been removed
- Struct `TemplateSpecsListBySubscriptionOptions` has been removed
- Struct `TemplateSpecsListBySubscriptionPager` has been removed
- Struct `TemplateSpecsListBySubscriptionResponse` has been removed
- Struct `TemplateSpecsListBySubscriptionResult` has been removed
- Struct `TemplateSpecsListResult` has been removed
- Struct `TemplateSpecsUpdateOptions` has been removed
- Struct `TemplateSpecsUpdateResponse` has been removed
- Struct `TemplateSpecsUpdateResult` has been removed
- Field `AzureResourceBase` of struct `TemplateSpec` has been removed
- Field `AzureResourceBase` of struct `TemplateSpecVersionUpdateModel` has been removed
- Field `AzureResourceBase` of struct `TemplateSpecUpdateModel` has been removed
- Field `AzureResourceBase` of struct `TemplateSpecVersion` has been removed

### Features Added

- New function `*ClientListByResourceGroupPager.Err() error`
- New function `*Client.Get(context.Context, string, string, *ClientGetOptions) (ClientGetResponse, error)`
- New function `*Client.Delete(context.Context, string, string, *ClientDeleteOptions) (ClientDeleteResponse, error)`
- New function `*ClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*TemplateSpecVersionsClientListPager.NextPage(context.Context) bool`
- New function `NewClient(string, azcore.TokenCredential, *arm.ClientOptions) *Client`
- New function `*Client.ListByResourceGroup(string, *ClientListByResourceGroupOptions) *ClientListByResourceGroupPager`
- New function `*Client.Update(context.Context, string, string, *ClientUpdateOptions) (ClientUpdateResponse, error)`
- New function `*ClientListByResourceGroupPager.PageResponse() ClientListByResourceGroupResponse`
- New function `*ClientListBySubscriptionPager.PageResponse() ClientListBySubscriptionResponse`
- New function `*TemplateSpecVersionsClientListPager.Err() error`
- New function `*TemplateSpecVersionsClientListPager.PageResponse() TemplateSpecVersionsClientListResponse`
- New function `ListResult.MarshalJSON() ([]byte, error)`
- New function `*Client.CreateOrUpdate(context.Context, string, string, TemplateSpec, *ClientCreateOrUpdateOptions) (ClientCreateOrUpdateResponse, error)`
- New function `*Client.ListBySubscription(*ClientListBySubscriptionOptions) *ClientListBySubscriptionPager`
- New function `*ClientListBySubscriptionPager.NextPage(context.Context) bool`
- New function `*ClientListBySubscriptionPager.Err() error`
- New struct `Client`
- New struct `ClientCreateOrUpdateOptions`
- New struct `ClientCreateOrUpdateResponse`
- New struct `ClientCreateOrUpdateResult`
- New struct `ClientDeleteOptions`
- New struct `ClientDeleteResponse`
- New struct `ClientGetOptions`
- New struct `ClientGetResponse`
- New struct `ClientGetResult`
- New struct `ClientListByResourceGroupOptions`
- New struct `ClientListByResourceGroupPager`
- New struct `ClientListByResourceGroupResponse`
- New struct `ClientListByResourceGroupResult`
- New struct `ClientListBySubscriptionOptions`
- New struct `ClientListBySubscriptionPager`
- New struct `ClientListBySubscriptionResponse`
- New struct `ClientListBySubscriptionResult`
- New struct `ClientUpdateOptions`
- New struct `ClientUpdateResponse`
- New struct `ClientUpdateResult`
- New struct `Error`
- New struct `ListResult`
- New struct `TemplateSpecVersionsClientCreateOrUpdateOptions`
- New struct `TemplateSpecVersionsClientCreateOrUpdateResponse`
- New struct `TemplateSpecVersionsClientCreateOrUpdateResult`
- New struct `TemplateSpecVersionsClientDeleteOptions`
- New struct `TemplateSpecVersionsClientDeleteResponse`
- New struct `TemplateSpecVersionsClientGetOptions`
- New struct `TemplateSpecVersionsClientGetResponse`
- New struct `TemplateSpecVersionsClientGetResult`
- New struct `TemplateSpecVersionsClientListOptions`
- New struct `TemplateSpecVersionsClientListPager`
- New struct `TemplateSpecVersionsClientListResponse`
- New struct `TemplateSpecVersionsClientListResult`
- New struct `TemplateSpecVersionsClientUpdateOptions`
- New struct `TemplateSpecVersionsClientUpdateResponse`
- New struct `TemplateSpecVersionsClientUpdateResult`
- New field `SystemData` in struct `TemplateSpec`
- New field `Type` in struct `TemplateSpec`
- New field `ID` in struct `TemplateSpec`
- New field `Name` in struct `TemplateSpec`
- New field `Type` in struct `TemplateSpecUpdateModel`
- New field `ID` in struct `TemplateSpecUpdateModel`
- New field `Name` in struct `TemplateSpecUpdateModel`
- New field `SystemData` in struct `TemplateSpecUpdateModel`
- New field `Type` in struct `TemplateSpecVersion`
- New field `ID` in struct `TemplateSpecVersion`
- New field `Name` in struct `TemplateSpecVersion`
- New field `SystemData` in struct `TemplateSpecVersion`
- New field `Type` in struct `TemplateSpecVersionUpdateModel`
- New field `ID` in struct `TemplateSpecVersionUpdateModel`
- New field `Name` in struct `TemplateSpecVersionUpdateModel`
- New field `SystemData` in struct `TemplateSpecVersionUpdateModel`


## 0.1.1 (2021-12-13)

### Other Changes

- Fix the go minimum version to `1.16`

## 0.1.0 (2021-11-16)

- Initial preview release.
