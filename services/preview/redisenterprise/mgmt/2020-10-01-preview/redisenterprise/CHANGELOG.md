Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewOperationListResultPage` parameter(s) have been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `NewClusterListPage` parameter(s) have been changed from `(func(context.Context, ClusterList) (ClusterList, error))` to `(ClusterList, func(context.Context, ClusterList) (ClusterList, error))`
- Function `NewDatabaseListPage` parameter(s) have been changed from `(func(context.Context, DatabaseList) (DatabaseList, error))` to `(DatabaseList, func(context.Context, DatabaseList) (DatabaseList, error))`
- Function `DatabaseClient.RegenerateKeyResponder` has been removed
- Function `DatabaseClient.Import` has been removed
- Function `DatabaseClient.ExportPreparer` has been removed
- Function `DatabaseClient.RegenerateKeyPreparer` has been removed
- Function `DatabaseClient.ImportPreparer` has been removed
- Function `DatabaseClient.ImportSender` has been removed
- Function `*DatabaseImportFuture.Result` has been removed
- Function `*DatabaseRegenerateKeyFuture.Result` has been removed
- Function `DatabaseClient.RegenerateKey` has been removed
- Function `DatabaseClient.ExportSender` has been removed
- Function `NewDatabaseClient` has been removed
- Function `DatabaseClient.ListKeysSender` has been removed
- Function `DatabaseClient.Export` has been removed
- Function `DatabaseClient.ExportResponder` has been removed
- Function `DatabaseClient.ListKeysResponder` has been removed
- Function `DatabaseClient.RegenerateKeySender` has been removed
- Function `*DatabaseExportFuture.Result` has been removed
- Function `NewDatabaseClientWithBaseURI` has been removed
- Function `DatabaseClient.ListKeysPreparer` has been removed
- Function `DatabaseClient.ImportResponder` has been removed
- Function `DatabaseClient.ListKeys` has been removed
- Struct `DatabaseClient` has been removed
- Struct `DatabaseExportFuture` has been removed
- Struct `DatabaseImportFuture` has been removed
- Struct `DatabaseRegenerateKeyFuture` has been removed

## New Content

- New function `DatabasesClient.ImportSender(*http.Request) (DatabasesImportFuture, error)`
- New function `DatabasesClient.ExportSender(*http.Request) (DatabasesExportFuture, error)`
- New function `DatabasesClient.ImportPreparer(context.Context, string, string, string, ImportClusterParameters) (*http.Request, error)`
- New function `*DatabasesImportFuture.Result(DatabasesClient) (autorest.Response, error)`
- New function `DatabasesClient.ExportPreparer(context.Context, string, string, string, ExportClusterParameters) (*http.Request, error)`
- New function `DatabasesClient.ListKeys(context.Context, string, string, string) (AccessKeys, error)`
- New function `DatabasesClient.ImportResponder(*http.Response) (autorest.Response, error)`
- New function `*DatabasesRegenerateKeyFuture.Result(DatabasesClient) (AccessKeys, error)`
- New function `DatabasesClient.ExportResponder(*http.Response) (autorest.Response, error)`
- New function `DatabasesClient.ListKeysPreparer(context.Context, string, string, string) (*http.Request, error)`
- New function `DatabasesClient.Export(context.Context, string, string, string, ExportClusterParameters) (DatabasesExportFuture, error)`
- New function `DatabasesClient.ListKeysResponder(*http.Response) (AccessKeys, error)`
- New function `DatabasesClient.Import(context.Context, string, string, string, ImportClusterParameters) (DatabasesImportFuture, error)`
- New function `DatabasesClient.RegenerateKeySender(*http.Request) (DatabasesRegenerateKeyFuture, error)`
- New function `DatabasesClient.RegenerateKeyPreparer(context.Context, string, string, string, RegenerateKeyParameters) (*http.Request, error)`
- New function `*DatabasesExportFuture.Result(DatabasesClient) (autorest.Response, error)`
- New function `DatabasesClient.ListKeysSender(*http.Request) (*http.Response, error)`
- New function `DatabasesClient.RegenerateKey(context.Context, string, string, string, RegenerateKeyParameters) (DatabasesRegenerateKeyFuture, error)`
- New function `DatabasesClient.RegenerateKeyResponder(*http.Response) (AccessKeys, error)`
- New struct `DatabasesExportFuture`
- New struct `DatabasesImportFuture`
- New struct `DatabasesRegenerateKeyFuture`
