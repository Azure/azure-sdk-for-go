# Release History

## 0.6.0 (2022-05-18)
### Breaking Changes

- Function `ErrorDetails.MarshalJSON` has been removed
- Function `Table.MarshalJSON` has been removed
- Function `Error.MarshalJSON` has been removed
- Function `FacetError.MarshalJSON` has been removed
- Function `FacetResult.MarshalJSON` has been removed
- Function `QueryResponse.MarshalJSON` has been removed
- Function `OperationListResult.MarshalJSON` has been removed


## 0.5.0 (2022-04-18)
### Breaking Changes

- Function `*OperationsClient.List` has been removed

### Features Added

- New function `*OperationsClient.NewListPager(*OperationsClientListOptions) *runtime.Pager[OperationsClientListResponse]`


## 0.4.0 (2022-04-13)
### Breaking Changes

- Function `*OperationsClient.List` parameter(s) have been changed from `(context.Context, *OperationsClientListOptions)` to `(*OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(OperationsClientListResponse, error)` to `(*runtime.Pager[OperationsClientListResponse])`
- Function `NewClient` return value(s) have been changed from `(*Client)` to `(*Client, error)`
- Function `NewOperationsClient` return value(s) have been changed from `(*OperationsClient)` to `(*OperationsClient, error)`
- Function `AuthorizationScopeFilter.ToPtr` has been removed
- Function `ColumnDataType.ToPtr` has been removed
- Function `FacetSortOrder.ToPtr` has been removed
- Function `ResultFormat.ToPtr` has been removed
- Function `ResultTruncated.ToPtr` has been removed
- Struct `ClientResourcesHistoryResult` has been removed
- Struct `ClientResourcesResult` has been removed
- Struct `OperationsClientListResult` has been removed
- Field `ClientResourcesResult` of struct `ClientResourcesResponse` has been removed
- Field `RawResponse` of struct `ClientResourcesResponse` has been removed
- Field `OperationsClientListResult` of struct `OperationsClientListResponse` has been removed
- Field `RawResponse` of struct `OperationsClientListResponse` has been removed
- Field `ClientResourcesHistoryResult` of struct `ClientResourcesHistoryResponse` has been removed
- Field `RawResponse` of struct `ClientResourcesHistoryResponse` has been removed

### Features Added

- New function `Table.MarshalJSON() ([]byte, error)`
- New function `Error.MarshalJSON() ([]byte, error)`
- New struct `Column`
- New struct `Error`
- New struct `ErrorResponse`
- New struct `Table`
- New anonymous field `QueryResponse` in struct `ClientResourcesResponse`
- New anonymous field `OperationListResult` in struct `OperationsClientListResponse`
- New field `Interface` in struct `ClientResourcesHistoryResponse`


## 0.3.1 (2022-02-22)

### Other Changes

- Remove the go_mod_tidy_hack.go file.

## 0.3.0 (2022-01-20)
### Breaking Changes

- Type of `QueryResponse.Data` has been changed from `map[string]interface{}` to `interface{}`
- Type of `FacetResult.Data` has been changed from `map[string]interface{}` to `interface{}`
- Type of `ErrorDetails.AdditionalProperties` has been changed from `map[string]map[string]interface{}` to `map[string]interface{}`
- Function `Error.MarshalJSON` has been removed
- Function `Table.MarshalJSON` has been removed
- Struct `Column` has been removed
- Struct `Error` has been removed
- Struct `ErrorResponse` has been removed
- Struct `Table` has been removed
- Field `Object` of struct `ClientResourcesHistoryResult` has been removed

### Features Added

- New field `Interface` in struct `ClientResourcesHistoryResult`


## 0.2.0 (2022-01-13)
### Breaking Changes

- Function `*OperationsClient.List` parameter(s) have been changed from `(context.Context, *OperationsListOptions)` to `(context.Context, *OperationsClientListOptions)`
- Function `*OperationsClient.List` return value(s) have been changed from `(OperationsListResponse, error)` to `(OperationsClientListResponse, error)`
- Function `*ResourceGraphClient.Resources` has been removed
- Function `NewResourceGraphClient` has been removed
- Function `ErrorResponse.Error` has been removed
- Function `*Facet.UnmarshalJSON` has been removed
- Function `*ResourceGraphClient.ResourcesHistory` has been removed
- Struct `OperationsListOptions` has been removed
- Struct `OperationsListResponse` has been removed
- Struct `OperationsListResult` has been removed
- Struct `ResourceGraphClient` has been removed
- Struct `ResourceGraphClientResourcesHistoryOptions` has been removed
- Struct `ResourceGraphClientResourcesHistoryResponse` has been removed
- Struct `ResourceGraphClientResourcesHistoryResult` has been removed
- Struct `ResourceGraphClientResourcesOptions` has been removed
- Struct `ResourceGraphClientResourcesResponse` has been removed
- Struct `ResourceGraphClientResourcesResult` has been removed
- Field `Facet` of struct `FacetError` has been removed
- Field `Facet` of struct `FacetResult` has been removed
- Field `InnerError` of struct `ErrorResponse` has been removed

### Features Added

- New function `NewClient(azcore.TokenCredential, *arm.ClientOptions) *Client`
- New function `*Client.ResourcesHistory(context.Context, ResourcesHistoryRequest, *ClientResourcesHistoryOptions) (ClientResourcesHistoryResponse, error)`
- New function `*Client.Resources(context.Context, QueryRequest, *ClientResourcesOptions) (ClientResourcesResponse, error)`
- New function `*FacetResult.GetFacet() *Facet`
- New function `*FacetError.GetFacet() *Facet`
- New struct `Client`
- New struct `ClientResourcesHistoryOptions`
- New struct `ClientResourcesHistoryResponse`
- New struct `ClientResourcesHistoryResult`
- New struct `ClientResourcesOptions`
- New struct `ClientResourcesResponse`
- New struct `ClientResourcesResult`
- New struct `OperationsClientListOptions`
- New struct `OperationsClientListResponse`
- New struct `OperationsClientListResult`
- New field `Error` in struct `ErrorResponse`
- New field `Expression` in struct `FacetResult`
- New field `ResultType` in struct `FacetResult`
- New field `Expression` in struct `FacetError`
- New field `ResultType` in struct `FacetError`


## 0.1.0 (2021-12-09)

- Init release.
