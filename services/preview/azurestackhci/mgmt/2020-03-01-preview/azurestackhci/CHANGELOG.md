Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `OperationsClient.List` return value(s) have been changed from `(OperationList, error)` to `(AvailableOperations, error)`
- Function `NewClusterListPage` parameter(s) have been changed from `(func(context.Context, ClusterList) (ClusterList, error))` to `(ClusterList, func(context.Context, ClusterList) (ClusterList, error))`
- Function `OperationsClient.ListResponder` return value(s) have been changed from `(OperationList, error)` to `(AvailableOperations, error)`
- Type of `ErrorResponse.Error` has been changed from `*ErrorResponseError` to `*ErrorDetail`
- Const `NeverConnected` has been removed
- Const `Expired` has been removed
- Function `Operation.MarshalJSON` has been removed
- Function `OperationList.MarshalJSON` has been removed
- Struct `ErrorResponseError` has been removed
- Struct `Operation` has been removed
- Struct `OperationList` has been removed

## New Content

- New const `NotYetRegistered`
- New const `Disconnected`
- New struct `AvailableOperations`
- New struct `ErrorDetail`
- New struct `OperationDetail`
- New field `LastSyncTimestamp` in struct `ClusterProperties`
- New field `LastBillingTimestamp` in struct `ClusterProperties`
- New field `RegistrationTimestamp` in struct `ClusterProperties`
