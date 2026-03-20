# Release History

## 0.10.0 (2026-03-20)
### Breaking Changes

- Function `NewClientFactory` parameter(s) have been changed from `(credential azcore.TokenCredential, options *arm.ClientOptions)` to `(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions)`
- Enum `ColumnDataType` has been removed
- Struct `Column` has been removed
- Struct `Error` has been removed
- Struct `ErrorResponse` has been removed
- Struct `Table` has been removed
- Field `Interface` of struct `ClientResourcesHistoryResponse` has been removed

### Features Added

- New enum type `ChangeCategory` with values `ChangeCategorySystem`, `ChangeCategoryUser`
- New enum type `ChangeType` with values `ChangeTypeCreate`, `ChangeTypeDelete`, `ChangeTypeUpdate`
- New enum type `CreatedByType` with values `CreatedByTypeApplication`, `CreatedByTypeKey`, `CreatedByTypeManagedIdentity`, `CreatedByTypeUser`
- New enum type `PropertyChangeType` with values `PropertyChangeTypeInsert`, `PropertyChangeTypeRemove`, `PropertyChangeTypeUpdate`
- New enum type `ResultKind` with values `ResultKindBasic`
- New function `*Client.NewGraphQueryClient() *GraphQueryClient`
- New function `*Client.NewOperationsClient() *OperationsClient`
- New function `*Client.ResourceChangeDetails(ctx context.Context, parameters ResourceChangeDetailsRequestParameters, options *ClientResourceChangeDetailsOptions) (ClientResourceChangeDetailsResponse, error)`
- New function `*Client.ResourceChanges(ctx context.Context, parameters ResourceChangesRequestParameters, options *ClientResourceChangesOptions) (ClientResourceChangesResponse, error)`
- New function `*ClientFactory.NewGraphQueryClient() *GraphQueryClient`
- New function `NewGraphQueryClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*GraphQueryClient, error)`
- New function `*GraphQueryClient.CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, properties GraphQueryResource, options *GraphQueryClientCreateOrUpdateOptions) (GraphQueryClientCreateOrUpdateResponse, error)`
- New function `*GraphQueryClient.Delete(ctx context.Context, resourceGroupName string, resourceName string, options *GraphQueryClientDeleteOptions) (GraphQueryClientDeleteResponse, error)`
- New function `*GraphQueryClient.Get(ctx context.Context, resourceGroupName string, resourceName string, options *GraphQueryClientGetOptions) (GraphQueryClientGetResponse, error)`
- New function `*GraphQueryClient.NewListBySubscriptionPager(options *GraphQueryClientListBySubscriptionOptions) *runtime.Pager[GraphQueryClientListBySubscriptionResponse]`
- New function `*GraphQueryClient.NewListPager(resourceGroupName string, options *GraphQueryClientListOptions) *runtime.Pager[GraphQueryClientListResponse]`
- New function `*GraphQueryClient.Update(ctx context.Context, resourceGroupName string, resourceName string, body GraphQueryUpdateParameters, options *GraphQueryClientUpdateOptions) (GraphQueryClientUpdateResponse, error)`
- New struct `GraphQueryListResult`
- New struct `GraphQueryProperties`
- New struct `GraphQueryPropertiesUpdateParameters`
- New struct `GraphQueryResource`
- New struct `GraphQueryUpdateParameters`
- New struct `ResourceChangeData`
- New struct `ResourceChangeDataAfterSnapshot`
- New struct `ResourceChangeDataBeforeSnapshot`
- New struct `ResourceChangeDetailsRequestParameters`
- New struct `ResourceChangeList`
- New struct `ResourceChangesRequestParameters`
- New struct `ResourceChangesRequestParametersInterval`
- New struct `ResourcePropertyChange`
- New struct `SystemData`
- New field `Value` in struct `ClientResourcesHistoryResponse`
- New field `NextLink` in struct `OperationListResult`


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