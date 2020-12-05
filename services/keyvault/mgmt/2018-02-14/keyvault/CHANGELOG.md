Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewOperationListResultPage` parameter(s) have been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `NewDeletedVaultListResultPage` parameter(s) have been changed from `(func(context.Context, DeletedVaultListResult) (DeletedVaultListResult, error))` to `(DeletedVaultListResult, func(context.Context, DeletedVaultListResult) (DeletedVaultListResult, error))`
- Function `NewResourceListResultPage` parameter(s) have been changed from `(func(context.Context, ResourceListResult) (ResourceListResult, error))` to `(ResourceListResult, func(context.Context, ResourceListResult) (ResourceListResult, error))`
- Function `NewVaultListResultPage` parameter(s) have been changed from `(func(context.Context, VaultListResult) (VaultListResult, error))` to `(VaultListResult, func(context.Context, VaultListResult) (VaultListResult, error))`

## New Content

- New function `DeletedVault.MarshalJSON() ([]byte, error)`
- New function `VaultProperties.MarshalJSON() ([]byte, error)`
- New function `PrivateLinkResourceProperties.MarshalJSON() ([]byte, error)`
- New function `VaultAccessPolicyParameters.MarshalJSON() ([]byte, error)`
