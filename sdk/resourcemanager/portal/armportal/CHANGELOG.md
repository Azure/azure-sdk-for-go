# Release History

## 0.5.0 (2022-05-17)
### Breaking Changes

- Function `ViolationsList.MarshalJSON` has been removed
- Function `ErrorDefinition.MarshalJSON` has been removed
- Function `DashboardListResult.MarshalJSON` has been removed
- Function `ResourceProviderOperationList.MarshalJSON` has been removed
- Function `ConfigurationList.MarshalJSON` has been removed


## 0.4.0 (2022-04-18)
### Breaking Changes

- Function `*DashboardsClient.ListBySubscription` has been removed
- Function `*ListTenantConfigurationViolationsClient.List` has been removed
- Function `*TenantConfigurationsClient.List` has been removed
- Function `*OperationsClient.List` has been removed
- Function `*DashboardsClient.ListByResourceGroup` has been removed

### Features Added

- New function `*ListTenantConfigurationViolationsClient.NewListPager(*ListTenantConfigurationViolationsClientListOptions) *runtime.Pager[ListTenantConfigurationViolationsClientListResponse]`
- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New function `*DashboardsClient.NewListBySubscriptionPager(*DashboardsClientListBySubscriptionOptions) *runtime.Pager[DashboardsClientListBySubscriptionResponse]`
- New function `*TenantConfigurationsClient.NewListPager(*TenantConfigurationsClientListOptions) *runtime.Pager[TenantConfigurationsClientListResponse]`
- New function `*DashboardsClient.NewListByResourceGroupPager(string, *DashboardsClientListByResourceGroupOptions) *runtime.Pager[DashboardsClientListByResourceGroupResponse]`


## 0.3.0 (2022-04-12)
### Breaking Changes

- Function `*DashboardsClient.ListBySubscription` return value(s) have been changed from `(*DashboardsClientListBySubscriptionPager)` to `(*runtime.Pager[DashboardsClientListBySubscriptionResponse])`
- Function `NewListTenantConfigurationViolationsClient` return value(s) have been changed from `(*ListTenantConfigurationViolationsClient)` to `(*ListTenantConfigurationViolationsClient, error)`
- Function `*TenantConfigurationsClient.List` return value(s) have been changed from `(*TenantConfigurationsClientListPager)` to `(*runtime.Pager[TenantConfigurationsClientListResponse])`
- Function `NewTenantConfigurationsClient` return value(s) have been changed from `(*TenantConfigurationsClient)` to `(*TenantConfigurationsClient, error)`
- Function `*ListTenantConfigurationViolationsClient.List` return value(s) have been changed from `(*ListTenantConfigurationViolationsClientListPager)` to `(*runtime.Pager[ListTenantConfigurationViolationsClientListResponse])`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `*DashboardsClient.ListByResourceGroup` return value(s) have been changed from `(*DashboardsClientListByResourceGroupPager)` to `(*runtime.Pager[DashboardsClientListByResourceGroupResponse])`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `NewDashboardsClient` return value(s) have been changed from `(*DashboardsClient)` to `(*DashboardsClient, error)`
- Type of `DashboardProperties.Metadata` has been changed from `map[string]map[string]interface{}` to `map[string]interface{}`
- Type of `DashboardPartsPosition.Metadata` has been changed from `map[string]map[string]interface{}` to `map[string]interface{}`
- Type of `MarkdownPartMetadata.Inputs` has been changed from `[]map[string]interface{}` to `[]interface{}`
- Type of `DashboardLens.Metadata` has been changed from `map[string]map[string]interface{}` to `map[string]interface{}`
- Function `*OperationsClientListPager.Err` has been removed
- Function `*DashboardsClientListBySubscriptionPager.NextPage` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Function `*ListTenantConfigurationViolationsClientListPager.NextPage` has been removed
- Function `*TenantConfigurationsClientListPager.PageResponse` has been removed
- Function `ConfigurationName.ToPtr` has been removed
- Function `*ListTenantConfigurationViolationsClientListPager.Err` has been removed
- Function `*DashboardsClientListBySubscriptionPager.PageResponse` has been removed
- Function `*DashboardsClientListBySubscriptionPager.Err` has been removed
- Function `*DashboardsClientListByResourceGroupPager.NextPage` has been removed
- Function `*OperationsClientListPager.NextPage` has been removed
- Function `*ListTenantConfigurationViolationsClientListPager.PageResponse` has been removed
- Function `*DashboardsClientListByResourceGroupPager.PageResponse` has been removed
- Function `*TenantConfigurationsClientListPager.NextPage` has been removed
- Function `*TenantConfigurationsClientListPager.Err` has been removed
- Function `*DashboardsClientListByResourceGroupPager.Err` has been removed
- Struct `DashboardsClientCreateOrUpdateResult` has been removed
- Struct `DashboardsClientGetResult` has been removed
- Struct `DashboardsClientListByResourceGroupPager` has been removed
- Struct `DashboardsClientListByResourceGroupResult` has been removed
- Struct `DashboardsClientListBySubscriptionPager` has been removed
- Struct `DashboardsClientListBySubscriptionResult` has been removed
- Struct `DashboardsClientUpdateResult` has been removed
- Struct `ListTenantConfigurationViolationsClientListPager` has been removed
- Struct `ListTenantConfigurationViolationsClientListResult` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Struct `TenantConfigurationsClientCreateResult` has been removed
- Struct `TenantConfigurationsClientGetResult` has been removed
- Struct `TenantConfigurationsClientListPager` has been removed
- Struct `TenantConfigurationsClientListResult` has been removed
- Field `DashboardsClientListBySubscriptionResult` of struct `DashboardsClientListBySubscriptionResponse` has been removed
- Field `RawResponse` of struct `DashboardsClientListBySubscriptionResponse` has been removed
- Field `DashboardsClientCreateOrUpdateResult` of struct `DashboardsClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `DashboardsClientCreateOrUpdateResponse` has been removed
- Field `DashboardsClientListByResourceGroupResult` of struct `DashboardsClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `DashboardsClientListByResourceGroupResponse` has been removed
- Field `ListTenantConfigurationViolationsClientListResult` of struct `ListTenantConfigurationViolationsClientListResponse` has been removed
- Field `RawResponse` of struct `ListTenantConfigurationViolationsClientListResponse` has been removed
- Field `RawResponse` of struct `TenantConfigurationsClientDeleteResponse` has been removed
- Field `TenantConfigurationsClientGetResult` of struct `TenantConfigurationsClientGetResponse` has been removed
- Field `RawResponse` of struct `TenantConfigurationsClientGetResponse` has been removed
- Field `TenantConfigurationsClientCreateResult` of struct `TenantConfigurationsClientCreateResponse` has been removed
- Field `RawResponse` of struct `TenantConfigurationsClientCreateResponse` has been removed
- Field `RawResponse` of struct `DashboardsClientDeleteResponse` has been removed
- Field `DashboardsClientUpdateResult` of struct `DashboardsClientUpdateResponse` has been removed
- Field `RawResponse` of struct `DashboardsClientUpdateResponse` has been removed
- Field `DashboardsClientGetResult` of struct `DashboardsClientGetResponse` has been removed
- Field `RawResponse` of struct `DashboardsClientGetResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `TenantConfigurationsClientListResult` of struct `TenantConfigurationsClientListResponse` has been removed
- Field `RawResponse` of struct `TenantConfigurationsClientListResponse` has been removed

### Features Added

- New anonymous field `Dashboard` in struct `DashboardsClientGetResponse`
- New anonymous field `DashboardListResult` in struct `DashboardsClientListBySubscriptionResponse`
- New anonymous field `ConfigurationList` in struct `TenantConfigurationsClientListResponse`
- New anonymous field `Dashboard` in struct `DashboardsClientCreateOrUpdateResponse`
- New anonymous field `ResourceProviderOperationList` in struct `OperationsClientListResponse`
- New anonymous field `Configuration` in struct `TenantConfigurationsClientCreateResponse`
- New anonymous field `Configuration` in struct `TenantConfigurationsClientGetResponse`
- New anonymous field `DashboardListResult` in struct `DashboardsClientListByResourceGroupResponse`
- New anonymous field `Dashboard` in struct `DashboardsClientUpdateResponse`
- New anonymous field `ViolationsList` in struct `ListTenantConfigurationViolationsClientListResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*TenantConfigurationsClient.List` parameter(s) have been changed from `(*TenantConfigurationsListOptions)` to `(*TenantConfigurationsClientListOptions)`
- Function `*TenantConfigurationsClient.List` return value(s) have been changed from `(*TenantConfigurationsListPager)` to `(*TenantConfigurationsClientListPager)`
- Function `*DashboardsClient.Update` parameter(s) have been changed from `(context.Context, string, string, PatchableDashboard, *DashboardsUpdateOptions)` to `(context.Context, string, string, PatchableDashboard, *DashboardsClientUpdateOptions)`
- Function `*DashboardsClient.Update` return value(s) have been changed from `(DashboardsUpdateResponse, error)` to `(DashboardsClientUpdateResponse, error)`
- Function `*DashboardsClient.ListByResourceGroup` parameter(s) have been changed from `(string, *DashboardsListByResourceGroupOptions)` to `(string, *DashboardsClientListByResourceGroupOptions)`
- Function `*DashboardsClient.ListByResourceGroup` return value(s) have been changed from `(*DashboardsListByResourceGroupPager)` to `(*DashboardsClientListByResourceGroupPager)`
- Function `*DashboardsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *DashboardsGetOptions)` to `(context.Context, string, string, *DashboardsClientGetOptions)`
- Function `*DashboardsClient.Get` return value(s) have been changed from `(DashboardsGetResponse, error)` to `(DashboardsClientGetResponse, error)`
- Function `*TenantConfigurationsClient.Delete` parameter(s) have been changed from `(context.Context, ConfigurationName, *TenantConfigurationsDeleteOptions)` to `(context.Context, ConfigurationName, *TenantConfigurationsClientDeleteOptions)`
- Function `*TenantConfigurationsClient.Delete` return value(s) have been changed from `(TenantConfigurationsDeleteResponse, error)` to `(TenantConfigurationsClientDeleteResponse, error)`
- Function `*TenantConfigurationsClient.Get` parameter(s) have been changed from `(context.Context, ConfigurationName, *TenantConfigurationsGetOptions)` to `(context.Context, ConfigurationName, *TenantConfigurationsClientGetOptions)`
- Function `*TenantConfigurationsClient.Get` return value(s) have been changed from `(TenantConfigurationsGetResponse, error)` to `(TenantConfigurationsClientGetResponse, error)`
- Function `*DashboardsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, Dashboard, *DashboardsCreateOrUpdateOptions)` to `(context.Context, string, string, Dashboard, *DashboardsClientCreateOrUpdateOptions)`
- Function `*DashboardsClient.CreateOrUpdate` return value(s) have been changed from `(DashboardsCreateOrUpdateResponse, error)` to `(DashboardsClientCreateOrUpdateResponse, error)`
- Function `*DashboardsClient.ListBySubscription` parameter(s) have been changed from `(*DashboardsListBySubscriptionOptions)` to `(*DashboardsClientListBySubscriptionOptions)`
- Function `*DashboardsClient.ListBySubscription` return value(s) have been changed from `(*DashboardsListBySubscriptionPager)` to `(*DashboardsClientListBySubscriptionPager)`
- Function `*TenantConfigurationsClient.Create` parameter(s) have been changed from `(context.Context, ConfigurationName, Configuration, *TenantConfigurationsCreateOptions)` to `(context.Context, ConfigurationName, Configuration, *TenantConfigurationsClientCreateOptions)`
- Function `*TenantConfigurationsClient.Create` return value(s) have been changed from `(TenantConfigurationsCreateResponse, error)` to `(TenantConfigurationsClientCreateResponse, error)`
- Function `*ListTenantConfigurationViolationsClient.List` parameter(s) have been changed from `(*ListTenantConfigurationViolationsListOptions)` to `(*ListTenantConfigurationViolationsClientListOptions)`
- Function `*ListTenantConfigurationViolationsClient.List` return value(s) have been changed from `(*ListTenantConfigurationViolationsListPager)` to `(*ListTenantConfigurationViolationsClientListPager)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*DashboardsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *DashboardsDeleteOptions)` to `(context.Context, string, string, *DashboardsClientDeleteOptions)`
- Function `*DashboardsClient.Delete` return value(s) have been changed from `(DashboardsDeleteResponse, error)` to `(DashboardsClientDeleteResponse, error)`
- Function `*TenantConfigurationsListPager.PageResponse` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Function `*ListTenantConfigurationViolationsListPager.Err` has been removed
- Function `*ListTenantConfigurationViolationsListPager.PageResponse` has been removed
- Function `*DashboardsListByResourceGroupPager.Err` has been removed
- Function `*DashboardsListBySubscriptionPager.Err` has been removed
- Function `*DashboardsListBySubscriptionPager.NextPage` has been removed
- Function `*DashboardsListBySubscriptionPager.PageResponse` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `*TenantConfigurationsListPager.Err` has been removed
- Function `*ListTenantConfigurationViolationsListPager.NextPage` has been removed
- Function `*DashboardsListByResourceGroupPager.PageResponse` has been removed
- Function `*OperationsListPager.PageResponse` has been removed
- Function `*TenantConfigurationsListPager.NextPage` has been removed
- Function `*DashboardsListByResourceGroupPager.NextPage` has been removed
- Function `ErrorResponse.Error` has been removed
- Struct `DashboardsCreateOrUpdateOptions` has been removed
- Struct `DashboardsCreateOrUpdateResponse` has been removed
- Struct `DashboardsCreateOrUpdateResult` has been removed
- Struct `DashboardsDeleteOptions` has been removed
- Struct `DashboardsDeleteResponse` has been removed
- Struct `DashboardsGetOptions` has been removed
- Struct `DashboardsGetResponse` has been removed
- Struct `DashboardsGetResult` has been removed
- Struct `DashboardsListByResourceGroupOptions` has been removed
- Struct `DashboardsListByResourceGroupPager` has been removed
- Struct `DashboardsListByResourceGroupResponse` has been removed
- Struct `DashboardsListByResourceGroupResult` has been removed
- Struct `DashboardsListBySubscriptionOptions` has been removed
- Struct `DashboardsListBySubscriptionPager` has been removed
- Struct `DashboardsListBySubscriptionResponse` has been removed
- Struct `DashboardsListBySubscriptionResult` has been removed
- Struct `DashboardsUpdateOptions` has been removed
- Struct `DashboardsUpdateResponse` has been removed
- Struct `DashboardsUpdateResult` has been removed
- Struct `ListTenantConfigurationViolationsListOptions` has been removed
- Struct `ListTenantConfigurationViolationsListPager` has been removed
- Struct `ListTenantConfigurationViolationsListResponse` has been removed
- Struct `ListTenantConfigurationViolationsListResult` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Struct `TenantConfigurationsCreateOptions` has been removed
- Struct `TenantConfigurationsCreateResponse` has been removed
- Struct `TenantConfigurationsCreateResult` has been removed
- Struct `TenantConfigurationsDeleteOptions` has been removed
- Struct `TenantConfigurationsDeleteResponse` has been removed
- Struct `TenantConfigurationsGetOptions` has been removed
- Struct `TenantConfigurationsGetResponse` has been removed
- Struct `TenantConfigurationsGetResult` has been removed
- Struct `TenantConfigurationsListOptions` has been removed
- Struct `TenantConfigurationsListPager` has been removed
- Struct `TenantConfigurationsListResponse` has been removed
- Struct `TenantConfigurationsListResult` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed
- Field `Resource` of struct `ProxyResource` has been removed
- Field `ProxyResource` of struct `Configuration` has been removed
- Field `DashboardPartMetadata` of struct `MarkdownPartMetadata` has been removed

### Features Added

- New function `*DashboardsClientListBySubscriptionPager.NextPage(context.Context) bool`
- New function `*ListTenantConfigurationViolationsClientListPager.NextPage(context.Context) bool`
- New function `*DashboardsClientListBySubscriptionPager.Err() error`
- New function `DashboardPartMetadata.MarshalJSON() ([]byte, error)`
- New function `*ListTenantConfigurationViolationsClientListPager.Err() error`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*TenantConfigurationsClientListPager.Err() error`
- New function `*DashboardsClientListByResourceGroupPager.PageResponse() DashboardsClientListByResourceGroupResponse`
- New function `*TenantConfigurationsClientListPager.PageResponse() TenantConfigurationsClientListResponse`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*MarkdownPartMetadata.GetDashboardPartMetadata() *DashboardPartMetadata`
- New function `*OperationsClientListPager.Err() error`
- New function `*DashboardsClientListByResourceGroupPager.Err() error`
- New function `*DashboardsClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*DashboardsClientListBySubscriptionPager.PageResponse() DashboardsClientListBySubscriptionResponse`
- New function `*TenantConfigurationsClientListPager.NextPage(context.Context) bool`
- New function `*ListTenantConfigurationViolationsClientListPager.PageResponse() ListTenantConfigurationViolationsClientListResponse`
- New struct `DashboardsClientCreateOrUpdateOptions`
- New struct `DashboardsClientCreateOrUpdateResponse`
- New struct `DashboardsClientCreateOrUpdateResult`
- New struct `DashboardsClientDeleteOptions`
- New struct `DashboardsClientDeleteResponse`
- New struct `DashboardsClientGetOptions`
- New struct `DashboardsClientGetResponse`
- New struct `DashboardsClientGetResult`
- New struct `DashboardsClientListByResourceGroupOptions`
- New struct `DashboardsClientListByResourceGroupPager`
- New struct `DashboardsClientListByResourceGroupResponse`
- New struct `DashboardsClientListByResourceGroupResult`
- New struct `DashboardsClientListBySubscriptionOptions`
- New struct `DashboardsClientListBySubscriptionPager`
- New struct `DashboardsClientListBySubscriptionResponse`
- New struct `DashboardsClientListBySubscriptionResult`
- New struct `DashboardsClientUpdateOptions`
- New struct `DashboardsClientUpdateResponse`
- New struct `DashboardsClientUpdateResult`
- New struct `ListTenantConfigurationViolationsClientListOptions`
- New struct `ListTenantConfigurationViolationsClientListPager`
- New struct `ListTenantConfigurationViolationsClientListResponse`
- New struct `ListTenantConfigurationViolationsClientListResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New struct `TenantConfigurationsClientCreateOptions`
- New struct `TenantConfigurationsClientCreateResponse`
- New struct `TenantConfigurationsClientCreateResult`
- New struct `TenantConfigurationsClientDeleteOptions`
- New struct `TenantConfigurationsClientDeleteResponse`
- New struct `TenantConfigurationsClientGetOptions`
- New struct `TenantConfigurationsClientGetResponse`
- New struct `TenantConfigurationsClientGetResult`
- New struct `TenantConfigurationsClientListOptions`
- New struct `TenantConfigurationsClientListPager`
- New struct `TenantConfigurationsClientListResponse`
- New struct `TenantConfigurationsClientListResult`
- New field `Name` in struct `Configuration`
- New field `Type` in struct `Configuration`
- New field `ID` in struct `Configuration`
- New field `Type` in struct `MarkdownPartMetadata`
- New field `AdditionalProperties` in struct `MarkdownPartMetadata`
- New field `Name` in struct `ProxyResource`
- New field `Type` in struct `ProxyResource`
- New field `ID` in struct `ProxyResource`
- New field `Error` in struct `ErrorResponse`


## 0.1.0 (2021-11-16)

- Initial preview release.
