# Release History

## 0.5.0 (2022-05-17)
### Breaking Changes

- Function `ResourceModelWithAllowedPropertySet.MarshalJSON` has been removed
- Function `CloudErrorBody.MarshalJSON` has been removed
- Function `ODataError.MarshalJSON` has been removed
- Function `Operation.MarshalJSON` has been removed
- Function `SQLServerListResult.MarshalJSON` has been removed
- Function `SQLServerRegistrationListResult.MarshalJSON` has been removed
- Function `OperationListResult.MarshalJSON` has been removed


## 0.4.0 (2022-04-15)
### Breaking Changes

- Function `*SQLServerRegistrationsClient.ListByResourceGroup` has been removed
- Function `*OperationsClient.List` has been removed
- Function `*SQLServersClient.ListByResourceGroup` has been removed
- Function `*SQLServerRegistrationsClient.List` has been removed

### Features Added

- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`
- New function `*SQLServerRegistrationsClient.NewListByResourceGroupPager(string, *SQLServerRegistrationsClientListByResourceGroupOptions) *runtime.Pager[SQLServerRegistrationsClientListByResourceGroupResponse]`
- New function `*SQLServersClient.NewListByResourceGroupPager(string, string, *SQLServersClientListByResourceGroupOptions) *runtime.Pager[SQLServersClientListByResourceGroupResponse]`
- New function `*SQLServerRegistrationsClient.NewListPager(*SQLServerRegistrationsClientListOptions) *runtime.Pager[SQLServerRegistrationsClientListResponse]`


## 0.3.0 (2022-04-11)
### Breaking Changes

- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `*SQLServerRegistrationsClient.List` return value(s) have been changed from `(*SQLServerRegistrationsClientListPager)` to `(*runtime.Pager[SQLServerRegistrationsClientListResponse])`
- Function `NewSQLServersClient` return value(s) have been changed from `(*SQLServersClient)` to `(*SQLServersClient, error)`
- Function `*SQLServerRegistrationsClient.ListByResourceGroup` return value(s) have been changed from `(*SQLServerRegistrationsClientListByResourceGroupPager)` to `(*runtime.Pager[SQLServerRegistrationsClientListByResourceGroupResponse])`
- Function `*SQLServersClient.ListByResourceGroup` return value(s) have been changed from `(*SQLServersClientListByResourceGroupPager)` to `(*runtime.Pager[SQLServersClientListByResourceGroupResponse])`
- Function `NewSQLServerRegistrationsClient` return value(s) have been changed from `(*SQLServerRegistrationsClient)` to `(*SQLServerRegistrationsClient, error)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsClientListPager)` to `(*runtime.Pager[OperationsClientListResponse])`
- Type of `Operation.Properties` has been changed from `map[string]map[string]interface{}` to `map[string]interface{}`
- Function `*SQLServersClientListByResourceGroupPager.Err` has been removed
- Function `IdentityType.ToPtr` has been removed
- Function `*SQLServerRegistrationsClientListByResourceGroupPager.Err` has been removed
- Function `*SQLServersClientListByResourceGroupPager.NextPage` has been removed
- Function `*OperationsClientListPager.PageResponse` has been removed
- Function `SKUTier.ToPtr` has been removed
- Function `OperationOrigin.ToPtr` has been removed
- Function `*SQLServerRegistrationsClientListPager.PageResponse` has been removed
- Function `*OperationsClientListPager.NextPage` has been removed
- Function `*SQLServerRegistrationsClientListPager.NextPage` has been removed
- Function `*SQLServerRegistrationsClientListByResourceGroupPager.PageResponse` has been removed
- Function `*SQLServerRegistrationsClientListPager.Err` has been removed
- Function `*OperationsClientListPager.Err` has been removed
- Function `*SQLServersClientListByResourceGroupPager.PageResponse` has been removed
- Function `*SQLServerRegistrationsClientListByResourceGroupPager.NextPage` has been removed
- Struct `OperationsClientListPager` has been removed
- Struct `OperationsClientListResult` has been removed
- Struct `SQLServerRegistrationsClientCreateOrUpdateResult` has been removed
- Struct `SQLServerRegistrationsClientGetResult` has been removed
- Struct `SQLServerRegistrationsClientListByResourceGroupPager` has been removed
- Struct `SQLServerRegistrationsClientListByResourceGroupResult` has been removed
- Struct `SQLServerRegistrationsClientListPager` has been removed
- Struct `SQLServerRegistrationsClientListResult` has been removed
- Struct `SQLServerRegistrationsClientUpdateResult` has been removed
- Struct `SQLServersClientCreateOrUpdateResult` has been removed
- Struct `SQLServersClientGetResult` has been removed
- Struct `SQLServersClientListByResourceGroupPager` has been removed
- Struct `SQLServersClientListByResourceGroupResult` has been removed
- Field `SQLServerRegistrationsClientListResult` of struct `SQLServerRegistrationsClientListResponse` has been removed
- Field `RawResponse` of struct `SQLServerRegistrationsClientListResponse` has been removed
- Field `SQLServerRegistrationsClientCreateOrUpdateResult` of struct `SQLServerRegistrationsClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `SQLServerRegistrationsClientCreateOrUpdateResponse` has been removed
- Field `SQLServersClientGetResult` of struct `SQLServersClientGetResponse` has been removed
- Field `RawResponse` of struct `SQLServersClientGetResponse` has been removed
- Field `SQLServerRegistrationsClientUpdateResult` of struct `SQLServerRegistrationsClientUpdateResponse` has been removed
- Field `RawResponse` of struct `SQLServerRegistrationsClientUpdateResponse` has been removed
- Field `SQLServerRegistrationsClientGetResult` of struct `SQLServerRegistrationsClientGetResponse` has been removed
- Field `RawResponse` of struct `SQLServerRegistrationsClientGetResponse` has been removed
- Field `SQLServersClientListByResourceGroupResult` of struct `SQLServersClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `SQLServersClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `SQLServerRegistrationsClientDeleteResponse` has been removed
- Field `SQLServerRegistrationsClientListByResourceGroupResult` of struct `SQLServerRegistrationsClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `SQLServerRegistrationsClientListByResourceGroupResponse` has been removed
- Field `RawResponse` of struct `SQLServersClientDeleteResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `SQLServersClientCreateOrUpdateResult` of struct `SQLServersClientCreateOrUpdateResponse` has been removed
- Field `RawResponse` of struct `SQLServersClientCreateOrUpdateResponse` has been removed

### Features Added

- New anonymous field `SQLServerRegistrationListResult` in struct `SQLServerRegistrationsClientListResponse`
- New anonymous field `SQLServerListResult` in struct `SQLServersClientListByResourceGroupResponse`
- New anonymous field `SQLServerRegistration` in struct `SQLServerRegistrationsClientUpdateResponse`
- New anonymous field `SQLServer` in struct `SQLServersClientCreateOrUpdateResponse`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New anonymous field `SQLServer` in struct `SQLServersClientGetResponse`
- New anonymous field `SQLServerRegistrationListResult` in struct `SQLServerRegistrationsClientListByResourceGroupResponse`
- New anonymous field `SQLServerRegistration` in struct `SQLServerRegistrationsClientGetResponse`
- New anonymous field `SQLServerRegistration` in struct `SQLServerRegistrationsClientCreateOrUpdateResponse`


## 0.2.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*SQLServersClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, *SQLServersGetOptions)` to `(context.Context, string, string, string, *SQLServersClientGetOptions)`
- Function `*SQLServersClient.Get` return value(s) have been changed from `(SQLServersGetResponse, error)` to `(SQLServersClientGetResponse, error)`
- Function `*OperationsClient.List` parameter(s) have been changed from `(*OperationsListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(*OperationsListPager)` to `(*OperationsClientListPager)`
- Function `*SQLServerRegistrationsClient.Get` parameter(s) have been changed from `(context.Context, string, string, *SQLServerRegistrationsGetOptions)` to `(context.Context, string, string, *SQLServerRegistrationsClientGetOptions)`
- Function `*SQLServerRegistrationsClient.Get` return value(s) have been changed from `(SQLServerRegistrationsGetResponse, error)` to `(SQLServerRegistrationsClientGetResponse, error)`
- Function `*SQLServerRegistrationsClient.Update` parameter(s) have been changed from `(context.Context, string, string, SQLServerRegistrationUpdate, *SQLServerRegistrationsUpdateOptions)` to `(context.Context, string, string, SQLServerRegistrationUpdate, *SQLServerRegistrationsClientUpdateOptions)`
- Function `*SQLServerRegistrationsClient.Update` return value(s) have been changed from `(SQLServerRegistrationsUpdateResponse, error)` to `(SQLServerRegistrationsClientUpdateResponse, error)`
- Function `*SQLServerRegistrationsClient.List` parameter(s) have been changed from `(*SQLServerRegistrationsListOptions)` to `(*SQLServerRegistrationsClientListOptions)`
- Function `*SQLServerRegistrationsClient.List` return value(s) have been changed from `(*SQLServerRegistrationsListPager)` to `(*SQLServerRegistrationsClientListPager)`
- Function `*SQLServersClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, string, SQLServer, *SQLServersCreateOrUpdateOptions)` to `(context.Context, string, string, string, SQLServer, *SQLServersClientCreateOrUpdateOptions)`
- Function `*SQLServersClient.CreateOrUpdate` return value(s) have been changed from `(SQLServersCreateOrUpdateResponse, error)` to `(SQLServersClientCreateOrUpdateResponse, error)`
- Function `*SQLServersClient.Delete` parameter(s) have been changed from `(context.Context, string, string, string, *SQLServersDeleteOptions)` to `(context.Context, string, string, string, *SQLServersClientDeleteOptions)`
- Function `*SQLServersClient.Delete` return value(s) have been changed from `(SQLServersDeleteResponse, error)` to `(SQLServersClientDeleteResponse, error)`
- Function `*SQLServerRegistrationsClient.CreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, SQLServerRegistration, *SQLServerRegistrationsCreateOrUpdateOptions)` to `(context.Context, string, string, SQLServerRegistration, *SQLServerRegistrationsClientCreateOrUpdateOptions)`
- Function `*SQLServerRegistrationsClient.CreateOrUpdate` return value(s) have been changed from `(SQLServerRegistrationsCreateOrUpdateResponse, error)` to `(SQLServerRegistrationsClientCreateOrUpdateResponse, error)`
- Function `*SQLServersClient.ListByResourceGroup` parameter(s) have been changed from `(string, string, *SQLServersListByResourceGroupOptions)` to `(string, string, *SQLServersClientListByResourceGroupOptions)`
- Function `*SQLServersClient.ListByResourceGroup` return value(s) have been changed from `(*SQLServersListByResourceGroupPager)` to `(*SQLServersClientListByResourceGroupPager)`
- Function `*SQLServerRegistrationsClient.ListByResourceGroup` parameter(s) have been changed from `(string, *SQLServerRegistrationsListByResourceGroupOptions)` to `(string, *SQLServerRegistrationsClientListByResourceGroupOptions)`
- Function `*SQLServerRegistrationsClient.ListByResourceGroup` return value(s) have been changed from `(*SQLServerRegistrationsListByResourceGroupPager)` to `(*SQLServerRegistrationsClientListByResourceGroupPager)`
- Function `*SQLServerRegistrationsClient.Delete` parameter(s) have been changed from `(context.Context, string, string, *SQLServerRegistrationsDeleteOptions)` to `(context.Context, string, string, *SQLServerRegistrationsClientDeleteOptions)`
- Function `*SQLServerRegistrationsClient.Delete` return value(s) have been changed from `(SQLServerRegistrationsDeleteResponse, error)` to `(SQLServerRegistrationsClientDeleteResponse, error)`
- Function `*SQLServersListByResourceGroupPager.PageResponse` has been removed
- Function `Resource.MarshalJSON` has been removed
- Function `*SQLServerRegistrationsListPager.Err` has been removed
- Function `CloudError.Error` has been removed
- Function `*SQLServerRegistrationsListPager.NextPage` has been removed
- Function `*SQLServersListByResourceGroupPager.NextPage` has been removed
- Function `*SQLServerRegistrationsListByResourceGroupPager.Err` has been removed
- Function `*SQLServersListByResourceGroupPager.Err` has been removed
- Function `*SQLServerRegistrationsListByResourceGroupPager.PageResponse` has been removed
- Function `*SQLServerRegistrationsListPager.PageResponse` has been removed
- Function `*SQLServerRegistrationsListByResourceGroupPager.NextPage` has been removed
- Function `*OperationsListPager.Err` has been removed
- Function `SQLServer.MarshalJSON` has been removed
- Function `*OperationsListPager.PageResponse` has been removed
- Function `*OperationsListPager.NextPage` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListPager` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Struct `SQLServerRegistrationsCreateOrUpdateOptions` has been removed
- Struct `SQLServerRegistrationsCreateOrUpdateResponse` has been removed
- Struct `SQLServerRegistrationsCreateOrUpdateResult` has been removed
- Struct `SQLServerRegistrationsDeleteOptions` has been removed
- Struct `SQLServerRegistrationsDeleteResponse` has been removed
- Struct `SQLServerRegistrationsGetOptions` has been removed
- Struct `SQLServerRegistrationsGetResponse` has been removed
- Struct `SQLServerRegistrationsGetResult` has been removed
- Struct `SQLServerRegistrationsListByResourceGroupOptions` has been removed
- Struct `SQLServerRegistrationsListByResourceGroupPager` has been removed
- Struct `SQLServerRegistrationsListByResourceGroupResponse` has been removed
- Struct `SQLServerRegistrationsListByResourceGroupResult` has been removed
- Struct `SQLServerRegistrationsListOptions` has been removed
- Struct `SQLServerRegistrationsListPager` has been removed
- Struct `SQLServerRegistrationsListResponse` has been removed
- Struct `SQLServerRegistrationsListResult` has been removed
- Struct `SQLServerRegistrationsUpdateOptions` has been removed
- Struct `SQLServerRegistrationsUpdateResponse` has been removed
- Struct `SQLServerRegistrationsUpdateResult` has been removed
- Struct `SQLServersCreateOrUpdateOptions` has been removed
- Struct `SQLServersCreateOrUpdateResponse` has been removed
- Struct `SQLServersCreateOrUpdateResult` has been removed
- Struct `SQLServersDeleteOptions` has been removed
- Struct `SQLServersDeleteResponse` has been removed
- Struct `SQLServersGetOptions` has been removed
- Struct `SQLServersGetResponse` has been removed
- Struct `SQLServersGetResult` has been removed
- Struct `SQLServersListByResourceGroupOptions` has been removed
- Struct `SQLServersListByResourceGroupPager` has been removed
- Struct `SQLServersListByResourceGroupResponse` has been removed
- Struct `SQLServersListByResourceGroupResult` has been removed
- Field `TrackedResource` of struct `SQLServerRegistration` has been removed
- Field `Resource` of struct `ProxyResource` has been removed
- Field `Identity` of struct `ResourceModelWithAllowedPropertySetIdentity` has been removed
- Field `ProxyResource` of struct `SQLServer` has been removed
- Field `SKU` of struct `ResourceModelWithAllowedPropertySetSKU` has been removed
- Field `Resource` of struct `TrackedResource` has been removed
- Field `InnerError` of struct `CloudError` has been removed
- Field `Plan` of struct `ResourceModelWithAllowedPropertySetPlan` has been removed

### Features Added

- New function `*SQLServerRegistrationsClientListByResourceGroupPager.Err() error`
- New function `*SQLServerRegistrationsClientListPager.PageResponse() SQLServerRegistrationsClientListResponse`
- New function `*SQLServerRegistrationsClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*OperationsClientListPager.PageResponse() OperationsClientListResponse`
- New function `*SQLServersClientListByResourceGroupPager.NextPage(context.Context) bool`
- New function `*SQLServerRegistrationsClientListPager.NextPage(context.Context) bool`
- New function `*OperationsClientListPager.NextPage(context.Context) bool`
- New function `*OperationsClientListPager.Err() error`
- New function `*SQLServersClientListByResourceGroupPager.PageResponse() SQLServersClientListByResourceGroupResponse`
- New function `*SQLServersClientListByResourceGroupPager.Err() error`
- New function `*SQLServerRegistrationsClientListPager.Err() error`
- New function `*SQLServerRegistrationsClientListByResourceGroupPager.PageResponse() SQLServerRegistrationsClientListByResourceGroupResponse`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListPager`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New struct `SQLServerRegistrationsClientCreateOrUpdateOptions`
- New struct `SQLServerRegistrationsClientCreateOrUpdateResponse`
- New struct `SQLServerRegistrationsClientCreateOrUpdateResult`
- New struct `SQLServerRegistrationsClientDeleteOptions`
- New struct `SQLServerRegistrationsClientDeleteResponse`
- New struct `SQLServerRegistrationsClientGetOptions`
- New struct `SQLServerRegistrationsClientGetResponse`
- New struct `SQLServerRegistrationsClientGetResult`
- New struct `SQLServerRegistrationsClientListByResourceGroupOptions`
- New struct `SQLServerRegistrationsClientListByResourceGroupPager`
- New struct `SQLServerRegistrationsClientListByResourceGroupResponse`
- New struct `SQLServerRegistrationsClientListByResourceGroupResult`
- New struct `SQLServerRegistrationsClientListOptions`
- New struct `SQLServerRegistrationsClientListPager`
- New struct `SQLServerRegistrationsClientListResponse`
- New struct `SQLServerRegistrationsClientListResult`
- New struct `SQLServerRegistrationsClientUpdateOptions`
- New struct `SQLServerRegistrationsClientUpdateResponse`
- New struct `SQLServerRegistrationsClientUpdateResult`
- New struct `SQLServersClientCreateOrUpdateOptions`
- New struct `SQLServersClientCreateOrUpdateResponse`
- New struct `SQLServersClientCreateOrUpdateResult`
- New struct `SQLServersClientDeleteOptions`
- New struct `SQLServersClientDeleteResponse`
- New struct `SQLServersClientGetOptions`
- New struct `SQLServersClientGetResponse`
- New struct `SQLServersClientGetResult`
- New struct `SQLServersClientListByResourceGroupOptions`
- New struct `SQLServersClientListByResourceGroupPager`
- New struct `SQLServersClientListByResourceGroupResponse`
- New struct `SQLServersClientListByResourceGroupResult`
- New field `Type` in struct `ResourceModelWithAllowedPropertySetIdentity`
- New field `PrincipalID` in struct `ResourceModelWithAllowedPropertySetIdentity`
- New field `TenantID` in struct `ResourceModelWithAllowedPropertySetIdentity`
- New field `Type` in struct `SQLServerRegistration`
- New field `Location` in struct `SQLServerRegistration`
- New field `Tags` in struct `SQLServerRegistration`
- New field `ID` in struct `SQLServerRegistration`
- New field `Name` in struct `SQLServerRegistration`
- New field `SystemData` in struct `SQLServerRegistration`
- New field `ID` in struct `TrackedResource`
- New field `Name` in struct `TrackedResource`
- New field `Type` in struct `TrackedResource`
- New field `Error` in struct `CloudError`
- New field `Name` in struct `ResourceModelWithAllowedPropertySetPlan`
- New field `Product` in struct `ResourceModelWithAllowedPropertySetPlan`
- New field `Publisher` in struct `ResourceModelWithAllowedPropertySetPlan`
- New field `PromotionCode` in struct `ResourceModelWithAllowedPropertySetPlan`
- New field `Version` in struct `ResourceModelWithAllowedPropertySetPlan`
- New field `Family` in struct `ResourceModelWithAllowedPropertySetSKU`
- New field `Size` in struct `ResourceModelWithAllowedPropertySetSKU`
- New field `Tier` in struct `ResourceModelWithAllowedPropertySetSKU`
- New field `Name` in struct `ResourceModelWithAllowedPropertySetSKU`
- New field `Capacity` in struct `ResourceModelWithAllowedPropertySetSKU`
- New field `Type` in struct `SQLServer`
- New field `ID` in struct `SQLServer`
- New field `Name` in struct `SQLServer`
- New field `ID` in struct `ProxyResource`
- New field `Name` in struct `ProxyResource`
- New field `Type` in struct `ProxyResource`


## 0.1.0 (2021-12-01)

- Initial preview release.
