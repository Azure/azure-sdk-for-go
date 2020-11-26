
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewListAmlUserFeatureResultPage` signature has been changed from `(func(context.Context, ListAmlUserFeatureResult) (ListAmlUserFeatureResult, error))` to `(ListAmlUserFeatureResult,func(context.Context, ListAmlUserFeatureResult) (ListAmlUserFeatureResult, error))`
- Function `NewWorkspaceListResultPage` signature has been changed from `(func(context.Context, WorkspaceListResult) (WorkspaceListResult, error))` to `(WorkspaceListResult,func(context.Context, WorkspaceListResult) (WorkspaceListResult, error))`
- Function `NewListWorkspaceQuotasPage` signature has been changed from `(func(context.Context, ListWorkspaceQuotas) (ListWorkspaceQuotas, error))` to `(ListWorkspaceQuotas,func(context.Context, ListWorkspaceQuotas) (ListWorkspaceQuotas, error))`
- Function `NewListUsagesResultPage` signature has been changed from `(func(context.Context, ListUsagesResult) (ListUsagesResult, error))` to `(ListUsagesResult,func(context.Context, ListUsagesResult) (ListUsagesResult, error))`
- Function `NewSkuListResultPage` signature has been changed from `(func(context.Context, SkuListResult) (SkuListResult, error))` to `(SkuListResult,func(context.Context, SkuListResult) (SkuListResult, error))`
- Function `NewPaginatedComputeResourcesListPage` signature has been changed from `(func(context.Context, PaginatedComputeResourcesList) (PaginatedComputeResourcesList, error))` to `(PaginatedComputeResourcesList,func(context.Context, PaginatedComputeResourcesList) (PaginatedComputeResourcesList, error))`

## New Content

- Const `ComputeInstanceAuthorizationTypePersonal` is added
- Function `PossibleComputeInstanceAuthorizationTypeValues() []ComputeInstanceAuthorizationType` is added
- Struct `AssignedUser` is added
- Struct `PersonalComputeInstanceSettings` is added
- Field `ComputeInstanceAuthorizationType` is added to struct `ComputeInstanceProperties`
- Field `PersonalComputeInstanceSettings` is added to struct `ComputeInstanceProperties`

