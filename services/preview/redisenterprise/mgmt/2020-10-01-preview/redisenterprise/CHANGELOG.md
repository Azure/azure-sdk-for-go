
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewClusterListPage` signature has been changed from `(func(context.Context, ClusterList) (ClusterList, error))` to `(ClusterList,func(context.Context, ClusterList) (ClusterList, error))`
- Function `NewOperationListResultPage` signature has been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult,func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `NewDatabaseListPage` signature has been changed from `(func(context.Context, DatabaseList) (DatabaseList, error))` to `(DatabaseList,func(context.Context, DatabaseList) (DatabaseList, error))`
- Function `DatabaseClient.RegenerateKeyResponder` has been removed
- Function `NewDatabaseClientWithBaseURI` has been removed
- Function `*DatabaseRegenerateKeyFuture.Result` has been removed
- Function `DatabaseClient.ImportResponder` has been removed
- Function `DatabaseClient.ListKeys` has been removed
- Function `NewDatabaseClient` has been removed
- Function `DatabaseClient.ImportPreparer` has been removed
- Function `DatabaseClient.RegenerateKeySender` has been removed
- Function `DatabaseClient.ListKeysPreparer` has been removed
- Function `*DatabaseImportFuture.Result` has been removed
- Function `DatabaseClient.ImportSender` has been removed
- Function `DatabaseClient.ListKeysSender` has been removed
- Function `DatabaseClient.Export` has been removed
- Function `DatabaseClient.ExportSender` has been removed
- Function `DatabaseClient.RegenerateKeyPreparer` has been removed
- Function `DatabaseClient.ListKeysResponder` has been removed
- Function `DatabaseClient.Import` has been removed
- Function `DatabaseClient.ExportPreparer` has been removed
- Function `DatabaseClient.ExportResponder` has been removed
- Function `DatabaseClient.RegenerateKey` has been removed
- Function `*DatabaseExportFuture.Result` has been removed
- Struct `DatabaseClient` has been removed
- Struct `DatabaseExportFuture` has been removed
- Struct `DatabaseImportFuture` has been removed
- Struct `DatabaseRegenerateKeyFuture` has been removed

## New Content

- Function `DatabasesClient.RegenerateKeyResponder(*http.Response) (AccessKeys,error)` is added
- Function `DatabasesClient.RegenerateKeyPreparer(context.Context,string,string,string,RegenerateKeyParameters) (*http.Request,error)` is added
- Function `DatabasesClient.ListKeys(context.Context,string,string,string) (AccessKeys,error)` is added
- Function `*DatabasesImportFuture.Result(DatabasesClient) (autorest.Response,error)` is added
- Function `DatabasesClient.ExportPreparer(context.Context,string,string,string,ExportClusterParameters) (*http.Request,error)` is added
- Function `DatabasesClient.RegenerateKey(context.Context,string,string,string,RegenerateKeyParameters) (DatabasesRegenerateKeyFuture,error)` is added
- Function `*DatabasesRegenerateKeyFuture.Result(DatabasesClient) (AccessKeys,error)` is added
- Function `DatabasesClient.RegenerateKeySender(*http.Request) (DatabasesRegenerateKeyFuture,error)` is added
- Function `DatabasesClient.ListKeysPreparer(context.Context,string,string,string) (*http.Request,error)` is added
- Function `DatabasesClient.ImportResponder(*http.Response) (autorest.Response,error)` is added
- Function `DatabasesClient.ListKeysResponder(*http.Response) (AccessKeys,error)` is added
- Function `DatabasesClient.ListKeysSender(*http.Request) (*http.Response,error)` is added
- Function `DatabasesClient.ExportResponder(*http.Response) (autorest.Response,error)` is added
- Function `*DatabasesExportFuture.Result(DatabasesClient) (autorest.Response,error)` is added
- Function `DatabasesClient.Export(context.Context,string,string,string,ExportClusterParameters) (DatabasesExportFuture,error)` is added
- Function `DatabasesClient.Import(context.Context,string,string,string,ImportClusterParameters) (DatabasesImportFuture,error)` is added
- Function `DatabasesClient.ExportSender(*http.Request) (DatabasesExportFuture,error)` is added
- Function `DatabasesClient.ImportPreparer(context.Context,string,string,string,ImportClusterParameters) (*http.Request,error)` is added
- Function `DatabasesClient.ImportSender(*http.Request) (DatabasesImportFuture,error)` is added
- Struct `DatabasesExportFuture` is added
- Struct `DatabasesImportFuture` is added
- Struct `DatabasesRegenerateKeyFuture` is added

