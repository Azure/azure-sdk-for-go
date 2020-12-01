Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `FileSharesClient.CreatePreparer` parameter(s) have been changed from `(context.Context, string, string, string, FileShare)` to `(context.Context, string, string, string, FileShare, PutSharesExpand)`
- Function `FileSharesClient.GetPreparer` parameter(s) have been changed from `(context.Context, string, string, string, GetShareExpand)` to `(context.Context, string, string, string, GetShareExpand, string)`
- Function `NewEncryptionScopeListResultPage` parameter(s) have been changed from `(func(context.Context, EncryptionScopeListResult) (EncryptionScopeListResult, error))` to `(EncryptionScopeListResult, func(context.Context, EncryptionScopeListResult) (EncryptionScopeListResult, error))`
- Function `NewDeletedAccountListResultPage` parameter(s) have been changed from `(func(context.Context, DeletedAccountListResult) (DeletedAccountListResult, error))` to `(DeletedAccountListResult, func(context.Context, DeletedAccountListResult) (DeletedAccountListResult, error))`
- Function `NewListContainerItemsPage` parameter(s) have been changed from `(func(context.Context, ListContainerItems) (ListContainerItems, error))` to `(ListContainerItems, func(context.Context, ListContainerItems) (ListContainerItems, error))`
- Function `FileSharesClient.Create` parameter(s) have been changed from `(context.Context, string, string, string, FileShare)` to `(context.Context, string, string, string, FileShare, PutSharesExpand)`
- Function `NewListTableResourcePage` parameter(s) have been changed from `(func(context.Context, ListTableResource) (ListTableResource, error))` to `(ListTableResource, func(context.Context, ListTableResource) (ListTableResource, error))`
- Function `FileSharesClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, GetShareExpand)` to `(context.Context, string, string, string, GetShareExpand, string)`
- Function `FileSharesClient.DeletePreparer` parameter(s) have been changed from `(context.Context, string, string, string)` to `(context.Context, string, string, string, string)`
- Function `NewAccountListResultPage` parameter(s) have been changed from `(func(context.Context, AccountListResult) (AccountListResult, error))` to `(AccountListResult, func(context.Context, AccountListResult) (AccountListResult, error))`
- Function `NewFileShareItemsPage` parameter(s) have been changed from `(func(context.Context, FileShareItems) (FileShareItems, error))` to `(FileShareItems, func(context.Context, FileShareItems) (FileShareItems, error))`
- Function `NewListQueueResourcePage` parameter(s) have been changed from `(func(context.Context, ListQueueResource) (ListQueueResource, error))` to `(ListQueueResource, func(context.Context, ListQueueResource) (ListQueueResource, error))`
- Function `FileSharesClient.Delete` parameter(s) have been changed from `(context.Context, string, string, string)` to `(context.Context, string, string, string, string)`

## New Content

- New const `ListSharesExpandSnapshots`
- New const `Snapshots`
- New function `PossiblePutSharesExpandValues() []PutSharesExpand`
- New field `SnapshotTime` in struct `FileShareProperties`
