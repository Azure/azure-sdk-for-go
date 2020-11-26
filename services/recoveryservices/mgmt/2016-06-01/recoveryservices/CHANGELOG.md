
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewPrivateLinkResourcesPage` signature has been changed from `(func(context.Context, PrivateLinkResources) (PrivateLinkResources, error))` to `(PrivateLinkResources,func(context.Context, PrivateLinkResources) (PrivateLinkResources, error))`
- Function `NewClientDiscoveryResponsePage` signature has been changed from `(func(context.Context, ClientDiscoveryResponse) (ClientDiscoveryResponse, error))` to `(ClientDiscoveryResponse,func(context.Context, ClientDiscoveryResponse) (ClientDiscoveryResponse, error))`
- Function `NewVaultListPage` signature has been changed from `(func(context.Context, VaultList) (VaultList, error))` to `(VaultList,func(context.Context, VaultList) (VaultList, error))`

## New Content

- Field `Identity` is added to struct `PatchVault`

