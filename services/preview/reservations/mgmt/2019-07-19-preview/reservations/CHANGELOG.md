Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewOperationListPage` parameter(s) have been changed from `(func(context.Context, OperationList) (OperationList, error))` to `(OperationList, func(context.Context, OperationList) (OperationList, error))`
- Function `NewListPage` parameter(s) have been changed from `(func(context.Context, List) (List, error))` to `(List, func(context.Context, List) (List, error))`
- Function `NewQuotaLimitsPage` parameter(s) have been changed from `(func(context.Context, QuotaLimits) (QuotaLimits, error))` to `(QuotaLimits, func(context.Context, QuotaLimits) (QuotaLimits, error))`
- Function `Client.AvailableScopes` parameter(s) have been changed from `(context.Context, string, string, []string)` to `(context.Context, string, string, AvailableScopeRequest)`
- Function `Client.AvailableScopesPreparer` parameter(s) have been changed from `(context.Context, string, string, []string)` to `(context.Context, string, string, AvailableScopeRequest)`
- Function `NewQuotaRequestDetailsListPage` parameter(s) have been changed from `(func(context.Context, QuotaRequestDetailsList) (QuotaRequestDetailsList, error))` to `(QuotaRequestDetailsList, func(context.Context, QuotaRequestDetailsList) (QuotaRequestDetailsList, error))`
- Function `NewOrderListPage` parameter(s) have been changed from `(func(context.Context, OrderList) (OrderList, error))` to `(OrderList, func(context.Context, OrderList) (OrderList, error))`

## New Content

- New const `SapHana`
- New const `PostgreSQL`
- New const `ManagedDisk`
- New const `MariaDb`
- New const `SQLAzureHybridBenefit`
- New const `BlockBlob`
- New const `DedicatedHost`
- New const `AppService`
- New const `Databricks`
- New const `AzureDataExplorer`
- New const `RedisCache`
- New const `MySQL`
- New function `Client.ArchivePreparer(context.Context, string, string) (*http.Request, error)`
- New function `Client.Unarchive(context.Context, string, string) (autorest.Response, error)`
- New function `Client.UnarchiveSender(*http.Request) (*http.Response, error)`
- New function `Client.UnarchivePreparer(context.Context, string, string) (*http.Request, error)`
- New function `Client.UnarchiveResponder(*http.Response) (autorest.Response, error)`
- New function `Client.Archive(context.Context, string, string) (autorest.Response, error)`
- New function `Client.ArchiveResponder(*http.Response) (autorest.Response, error)`
- New function `Client.ArchiveSender(*http.Request) (*http.Response, error)`
- New struct `AvailableScopeRequest`
- New struct `AvailableScopeRequestProperties`
