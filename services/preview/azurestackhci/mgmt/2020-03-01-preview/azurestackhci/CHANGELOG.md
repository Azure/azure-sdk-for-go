
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `OperationsClient.List` return values have been changed from `(OperationList,error)` to `(AvailableOperations,error)`
- Function `NewClusterListPage` signature has been changed from `(func(context.Context, ClusterList) (ClusterList, error))` to `(ClusterList,func(context.Context, ClusterList) (ClusterList, error))`
- Function `OperationsClient.ListResponder` return values have been changed from `(OperationList,error)` to `(AvailableOperations,error)`
- Type of `ErrorResponse.Error` has been changed from `*ErrorResponseError` to `*ErrorDetail`
- Const `Expired` has been removed
- Const `NeverConnected` has been removed
- Function `Operation.MarshalJSON` has been removed
- Function `OperationList.MarshalJSON` has been removed
- Struct `ErrorResponseError` has been removed
- Struct `Operation` has been removed
- Struct `OperationList` has been removed

## New Content

- Const `Disconnected` is added
- Const `NotYetRegistered` is added
- Struct `AvailableOperations` is added
- Struct `ErrorDetail` is added
- Struct `OperationDetail` is added
- Field `LastSyncTimestamp` is added to struct `ClusterProperties`
- Field `RegistrationTimestamp` is added to struct `ClusterProperties`
- Field `LastBillingTimestamp` is added to struct `ClusterProperties`

