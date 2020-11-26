
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewListPage` signature has been changed from `(func(context.Context, List) (List, error))` to `(List,func(context.Context, List) (List, error))`
- Function `Client.AvailableScopesPreparer` signature has been changed from `(context.Context,string,string,[]string)` to `(context.Context,string,string,AvailableScopeRequest)`
- Function `Client.AvailableScopes` signature has been changed from `(context.Context,string,string,[]string)` to `(context.Context,string,string,AvailableScopeRequest)`
- Function `NewOrderListPage` signature has been changed from `(func(context.Context, OrderList) (OrderList, error))` to `(OrderList,func(context.Context, OrderList) (OrderList, error))`
- Function `NewOperationListPage` signature has been changed from `(func(context.Context, OperationList) (OperationList, error))` to `(OperationList,func(context.Context, OperationList) (OperationList, error))`

## New Content

- Const `DedicatedHost` is added
- Const `ManagedDisk` is added
- Const `RedisCache` is added
- Const `BlockBlob` is added
- Const `PostgreSQL` is added
- Const `SapHana` is added
- Const `MariaDb` is added
- Const `SQLAzureHybridBenefit` is added
- Const `Databricks` is added
- Const `AzureDataExplorer` is added
- Const `MySQL` is added
- Const `AppService` is added
- Function `Client.UnarchiveResponder(*http.Response) (autorest.Response,error)` is added
- Function `Client.Unarchive(context.Context,string,string) (autorest.Response,error)` is added
- Function `Client.UnarchiveSender(*http.Request) (*http.Response,error)` is added
- Function `Client.Archive(context.Context,string,string) (autorest.Response,error)` is added
- Function `Client.ArchiveSender(*http.Request) (*http.Response,error)` is added
- Function `Client.UnarchivePreparer(context.Context,string,string) (*http.Request,error)` is added
- Function `Client.ArchiveResponder(*http.Response) (autorest.Response,error)` is added
- Function `Client.ArchivePreparer(context.Context,string,string) (*http.Request,error)` is added
- Struct `AvailableScopeRequest` is added
- Struct `AvailableScopeRequestProperties` is added

