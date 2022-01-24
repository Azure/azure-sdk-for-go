# Release History

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
