# Release History

## 0.10.0 (2026-03-17)
### Breaking Changes

- Enum `AuthorizationScopeFilter` has been removed
- Enum `ColumnDataType` has been removed
- Enum `FacetSortOrder` has been removed
- Enum `ResultTruncated` has been removed
- Function `*Client.Resources` has been removed
- Function `*ClientFactory.NewOperationsClient` has been removed
- Function `*Facet.GetFacet` has been removed
- Function `*FacetError.GetFacet` has been removed
- Function `*FacetResult.GetFacet` has been removed
- Function `NewOperationsClient` has been removed
- Function `*OperationsClient.NewListPager` has been removed
- Struct `Column` has been removed
- Struct `Error` has been removed
- Struct `ErrorDetails` has been removed
- Struct `ErrorResponse` has been removed
- Struct `FacetError` has been removed
- Struct `FacetRequest` has been removed
- Struct `FacetRequestOptions` has been removed
- Struct `FacetResult` has been removed
- Struct `Operation` has been removed
- Struct `OperationDisplay` has been removed
- Struct `OperationListResult` has been removed
- Struct `QueryRequest` has been removed
- Struct `QueryRequestOptions` has been removed
- Struct `QueryResponse` has been removed
- Struct `Table` has been removed
- Field `Interface` of struct `ClientResourcesHistoryResponse` has been removed

### Features Added

- New enum type `ChangeCategory` with values `ChangeCategorySystem`, `ChangeCategoryUser`
- New enum type `ChangeType` with values `ChangeTypeCreate`, `ChangeTypeDelete`, `ChangeTypeUpdate`
- New enum type `PropertyChangeType` with values `PropertyChangeTypeInsert`, `PropertyChangeTypeRemove`, `PropertyChangeTypeUpdate`
- New function `*Client.ResourceChangeDetails(ctx context.Context, parameters ResourceChangeDetailsRequestParameters, options *ClientResourceChangeDetailsOptions) (ClientResourceChangeDetailsResponse, error)`
- New function `*Client.ResourceChanges(ctx context.Context, parameters ResourceChangesRequestParameters, options *ClientResourceChangesOptions) (ClientResourceChangesResponse, error)`
- New struct `ResourceChangeData`
- New struct `ResourceChangeDataAfterSnapshot`
- New struct `ResourceChangeDataBeforeSnapshot`
- New struct `ResourceChangeDetailsRequestParameters`
- New struct `ResourceChangeList`
- New struct `ResourceChangesRequestParameters`
- New struct `ResourceChangesRequestParametersInterval`
- New struct `ResourcePropertyChange`
- New field `Value` in struct `ClientResourcesHistoryResponse`


## 0.9.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.8.2 (2023-10-09)

### Other Changes

- Updated to latest `azcore` beta.

## 0.8.1 (2023-07-19)

### Bug Fixes

- Fixed a potential panic in faked paged and long-running operations.

## 0.8.0 (2023-06-13)

### Features Added

- Support for test fakes and OpenTelemetry trace spans.

## 0.7.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 0.7.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.6.0 (2022-05-18)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcegraph/armresourcegraph` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.6.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).