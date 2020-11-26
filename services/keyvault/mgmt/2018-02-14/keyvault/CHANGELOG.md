
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewResourceListResultPage` signature has been changed from `(func(context.Context, ResourceListResult) (ResourceListResult, error))` to `(ResourceListResult,func(context.Context, ResourceListResult) (ResourceListResult, error))`
- Function `NewVaultListResultPage` signature has been changed from `(func(context.Context, VaultListResult) (VaultListResult, error))` to `(VaultListResult,func(context.Context, VaultListResult) (VaultListResult, error))`
- Function `NewOperationListResultPage` signature has been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult,func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `NewDeletedVaultListResultPage` signature has been changed from `(func(context.Context, DeletedVaultListResult) (DeletedVaultListResult, error))` to `(DeletedVaultListResult,func(context.Context, DeletedVaultListResult) (DeletedVaultListResult, error))`

## New Content

- Function `VaultAccessPolicyParameters.MarshalJSON() ([]byte,error)` is added
- Function `DeletedVault.MarshalJSON() ([]byte,error)` is added
- Function `PrivateLinkResourceProperties.MarshalJSON() ([]byte,error)` is added
- Function `VaultProperties.MarshalJSON() ([]byte,error)` is added

