
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewSQLServerListResultPage` signature has been changed from `(func(context.Context, SQLServerListResult) (SQLServerListResult, error))` to `(SQLServerListResult,func(context.Context, SQLServerListResult) (SQLServerListResult, error))`
- Function `NewPageOfDataControllerResourcePage` signature has been changed from `(func(context.Context, PageOfDataControllerResource) (PageOfDataControllerResource, error))` to `(PageOfDataControllerResource,func(context.Context, PageOfDataControllerResource) (PageOfDataControllerResource, error))`
- Function `DataControllersClient.PatchDataController` signature has been changed from `(context.Context,string,string)` to `(context.Context,string,string,DataControllerUpdate)`
- Function `DataControllersClient.PatchDataControllerPreparer` signature has been changed from `(context.Context,string,string)` to `(context.Context,string,string,DataControllerUpdate)`
- Function `NewSQLServerRegistrationListResultPage` signature has been changed from `(func(context.Context, SQLServerRegistrationListResult) (SQLServerRegistrationListResult, error))` to `(SQLServerRegistrationListResult,func(context.Context, SQLServerRegistrationListResult) (SQLServerRegistrationListResult, error))`
- Function `NewSQLManagedInstanceListResultPage` signature has been changed from `(func(context.Context, SQLManagedInstanceListResult) (SQLManagedInstanceListResult, error))` to `(SQLManagedInstanceListResult,func(context.Context, SQLManagedInstanceListResult) (SQLManagedInstanceListResult, error))`
- Function `NewOperationListResultPage` signature has been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult,func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `NewPostgresInstanceListResultPage` signature has been changed from `(func(context.Context, PostgresInstanceListResult) (PostgresInstanceListResult, error))` to `(PostgresInstanceListResult,func(context.Context, PostgresInstanceListResult) (PostgresInstanceListResult, error))`
- Function `NewSQLServerInstanceListResultPage` signature has been changed from `(func(context.Context, SQLServerInstanceListResult) (SQLServerInstanceListResult, error))` to `(SQLServerInstanceListResult,func(context.Context, SQLServerInstanceListResult) (SQLServerInstanceListResult, error))`
- Const `UsageUploadStatusPartialSuccess` has been removed
- Const `UsageUploadStatusUnknown` has been removed
- Const `Handshake` has been removed
- Const `Unknown` has been removed
- Const `UsageUpload` has been removed
- Const `UsageUploadStatusCompleted` has been removed
- Const `UsageUploadStatusFailed` has been removed
- Function `AzureResource.MarshalJSON` has been removed
- Function `PossibleRequestTypeValues` has been removed
- Function `PossibleUsageUploadStatusValues` has been removed
- Struct `AzureResource` has been removed
- Struct `HandshakeResponse` has been removed
- Struct `UsageRecord` has been removed
- Struct `UsageUploadRequest` has been removed
- Struct `UsageUploadResponse` has been removed
- Field `HandshakeRequest` of struct `DataControllerProperties` has been removed
- Field `HandshakeResponse` of struct `DataControllerProperties` has been removed
- Field `RequestType` of struct `DataControllerProperties` has been removed
- Field `UploadRequest` of struct `DataControllerProperties` has been removed
- Field `UploadResponse` of struct `DataControllerProperties` has been removed

## New Content

- Function `DataControllerUpdate.MarshalJSON() ([]byte,error)` is added
- Struct `DataControllerUpdate` is added
- Field `VCore` is added to struct `SQLManagedInstanceProperties`
- Field `InstanceEndpoint` is added to struct `SQLManagedInstanceProperties`
- Field `Admin` is added to struct `SQLManagedInstanceProperties`
- Field `StartTime` is added to struct `SQLManagedInstanceProperties`
- Field `EndTime` is added to struct `SQLManagedInstanceProperties`

