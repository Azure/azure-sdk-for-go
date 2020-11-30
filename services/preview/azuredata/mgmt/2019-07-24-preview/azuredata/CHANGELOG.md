Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `DataControllersClient.PatchDataControllerPreparer` parameter(s) have been changed from `(context.Context, string, string)` to `(context.Context, string, string, DataControllerUpdate)`
- Function `DataControllersClient.PatchDataController` parameter(s) have been changed from `(context.Context, string, string)` to `(context.Context, string, string, DataControllerUpdate)`
- Function `NewSQLServerRegistrationListResultPage` parameter(s) have been changed from `(func(context.Context, SQLServerRegistrationListResult) (SQLServerRegistrationListResult, error))` to `(SQLServerRegistrationListResult, func(context.Context, SQLServerRegistrationListResult) (SQLServerRegistrationListResult, error))`
- Function `NewOperationListResultPage` parameter(s) have been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `NewPageOfDataControllerResourcePage` parameter(s) have been changed from `(func(context.Context, PageOfDataControllerResource) (PageOfDataControllerResource, error))` to `(PageOfDataControllerResource, func(context.Context, PageOfDataControllerResource) (PageOfDataControllerResource, error))`
- Function `NewSQLServerInstanceListResultPage` parameter(s) have been changed from `(func(context.Context, SQLServerInstanceListResult) (SQLServerInstanceListResult, error))` to `(SQLServerInstanceListResult, func(context.Context, SQLServerInstanceListResult) (SQLServerInstanceListResult, error))`
- Function `NewPostgresInstanceListResultPage` parameter(s) have been changed from `(func(context.Context, PostgresInstanceListResult) (PostgresInstanceListResult, error))` to `(PostgresInstanceListResult, func(context.Context, PostgresInstanceListResult) (PostgresInstanceListResult, error))`
- Function `NewSQLServerListResultPage` parameter(s) have been changed from `(func(context.Context, SQLServerListResult) (SQLServerListResult, error))` to `(SQLServerListResult, func(context.Context, SQLServerListResult) (SQLServerListResult, error))`
- Function `NewSQLManagedInstanceListResultPage` parameter(s) have been changed from `(func(context.Context, SQLManagedInstanceListResult) (SQLManagedInstanceListResult, error))` to `(SQLManagedInstanceListResult, func(context.Context, SQLManagedInstanceListResult) (SQLManagedInstanceListResult, error))`
- Const `UsageUploadStatusFailed` has been removed
- Const `UsageUploadStatusCompleted` has been removed
- Const `Handshake` has been removed
- Const `UsageUpload` has been removed
- Const `UsageUploadStatusUnknown` has been removed
- Const `UsageUploadStatusPartialSuccess` has been removed
- Const `Unknown` has been removed
- Function `PossibleUsageUploadStatusValues` has been removed
- Function `AzureResource.MarshalJSON` has been removed
- Function `PossibleRequestTypeValues` has been removed
- Struct `AzureResource` has been removed
- Struct `HandshakeResponse` has been removed
- Struct `UsageRecord` has been removed
- Struct `UsageUploadRequest` has been removed
- Struct `UsageUploadResponse` has been removed
- Field `HandshakeResponse` of struct `DataControllerProperties` has been removed
- Field `RequestType` of struct `DataControllerProperties` has been removed
- Field `UploadRequest` of struct `DataControllerProperties` has been removed
- Field `UploadResponse` of struct `DataControllerProperties` has been removed
- Field `HandshakeRequest` of struct `DataControllerProperties` has been removed

## New Content

- New function `DataControllerUpdate.MarshalJSON() ([]byte, error)`
- New struct `DataControllerUpdate`
- New field `EndTime` in struct `SQLManagedInstanceProperties`
- New field `VCore` in struct `SQLManagedInstanceProperties`
- New field `InstanceEndpoint` in struct `SQLManagedInstanceProperties`
- New field `Admin` in struct `SQLManagedInstanceProperties`
- New field `StartTime` in struct `SQLManagedInstanceProperties`
