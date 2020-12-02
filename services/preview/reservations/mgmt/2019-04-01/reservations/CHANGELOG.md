Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewListPage` parameter(s) have been changed from `(func(context.Context, List) (List, error))` to `(List, func(context.Context, List) (List, error))`
- Function `NewOrderListPage` parameter(s) have been changed from `(func(context.Context, OrderList) (OrderList, error))` to `(OrderList, func(context.Context, OrderList) (OrderList, error))`
- Function `NewOperationListPage` parameter(s) have been changed from `(func(context.Context, OperationList) (OperationList, error))` to `(OperationList, func(context.Context, OperationList) (OperationList, error))`
- Function `Client.AvailableScopesPreparer` parameter(s) have been changed from `(context.Context, string, string, []string)` to `(context.Context, string, string, AvailableScopeRequest)`
- Function `Client.AvailableScopes` parameter(s) have been changed from `(context.Context, string, string, []string)` to `(context.Context, string, string, AvailableScopeRequest)`

## New Content

- New const `ManagedDisk`
- New const `Databricks`
- New const `SQLAzureHybridBenefit`
- New const `AzureDataExplorer`
- New const `BlockBlob`
- New const `SapHana`
- New const `PostgreSQL`
- New const `MySQL`
- New const `DedicatedHost`
- New const `RedisCache`
- New const `AppService`
- New const `MariaDb`
- New function `Client.ArchiveSender(*http.Request) (*http.Response, error)`
- New function `Client.UnarchivePreparer(context.Context, string, string) (*http.Request, error)`
- New function `Client.UnarchiveSender(*http.Request) (*http.Response, error)`
- New function `Client.ArchivePreparer(context.Context, string, string) (*http.Request, error)`
- New function `Client.ArchiveResponder(*http.Response) (autorest.Response, error)`
- New function `Client.Unarchive(context.Context, string, string) (autorest.Response, error)`
- New function `Client.Archive(context.Context, string, string) (autorest.Response, error)`
- New function `Client.UnarchiveResponder(*http.Response) (autorest.Response, error)`
- New struct `AvailableScopeRequest`
- New struct `AvailableScopeRequestProperties`
