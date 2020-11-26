
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `FileSharesClient.Create` signature has been changed from `(context.Context,string,string,string,FileShare)` to `(context.Context,string,string,string,FileShare,PutSharesExpand)`
- Function `FileSharesClient.CreatePreparer` signature has been changed from `(context.Context,string,string,string,FileShare)` to `(context.Context,string,string,string,FileShare,PutSharesExpand)`
- Function `NewListQueueResourcePage` signature has been changed from `(func(context.Context, ListQueueResource) (ListQueueResource, error))` to `(ListQueueResource,func(context.Context, ListQueueResource) (ListQueueResource, error))`
- Function `FileSharesClient.Get` signature has been changed from `(context.Context,string,string,string,GetShareExpand)` to `(context.Context,string,string,string,GetShareExpand,string)`
- Function `NewDeletedAccountListResultPage` signature has been changed from `(func(context.Context, DeletedAccountListResult) (DeletedAccountListResult, error))` to `(DeletedAccountListResult,func(context.Context, DeletedAccountListResult) (DeletedAccountListResult, error))`
- Function `NewListTableResourcePage` signature has been changed from `(func(context.Context, ListTableResource) (ListTableResource, error))` to `(ListTableResource,func(context.Context, ListTableResource) (ListTableResource, error))`
- Function `FileSharesClient.Delete` signature has been changed from `(context.Context,string,string,string)` to `(context.Context,string,string,string,string)`
- Function `NewListContainerItemsPage` signature has been changed from `(func(context.Context, ListContainerItems) (ListContainerItems, error))` to `(ListContainerItems,func(context.Context, ListContainerItems) (ListContainerItems, error))`
- Function `NewFileShareItemsPage` signature has been changed from `(func(context.Context, FileShareItems) (FileShareItems, error))` to `(FileShareItems,func(context.Context, FileShareItems) (FileShareItems, error))`
- Function `FileSharesClient.DeletePreparer` signature has been changed from `(context.Context,string,string,string)` to `(context.Context,string,string,string,string)`
- Function `NewAccountListResultPage` signature has been changed from `(func(context.Context, AccountListResult) (AccountListResult, error))` to `(AccountListResult,func(context.Context, AccountListResult) (AccountListResult, error))`
- Function `FileSharesClient.GetPreparer` signature has been changed from `(context.Context,string,string,string,GetShareExpand)` to `(context.Context,string,string,string,GetShareExpand,string)`
- Function `NewEncryptionScopeListResultPage` signature has been changed from `(func(context.Context, EncryptionScopeListResult) (EncryptionScopeListResult, error))` to `(EncryptionScopeListResult,func(context.Context, EncryptionScopeListResult) (EncryptionScopeListResult, error))`

## New Content

- Const `ListSharesExpandSnapshots` is added
- Const `Snapshots` is added
- Function `PossiblePutSharesExpandValues() []PutSharesExpand` is added
- Field `SnapshotTime` is added to struct `FileShareProperties`

