
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewClusterListPage` signature has been changed from `(func(context.Context, ClusterList) (ClusterList, error))` to `(ClusterList,func(context.Context, ClusterList) (ClusterList, error))`
- Type of `ErrorResponse.Error` has been changed from `*ErrorResponseError` to `*ErrorDetail`
- Struct `ErrorResponseError` has been removed

## New Content

- Struct `ErrorDetail` is added

