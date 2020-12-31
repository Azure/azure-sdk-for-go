Generated from https://github.com/Azure/azure-rest-api-specs/tree/b08824e05817297a4b2874d8db5e6fc8c29349c9

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

### Removed Funcs

1. *MoveCollectionsBulkRemoveFuture.Result(MoveCollectionsClient) (OperationStatus, error)
1. *MoveCollectionsCommitFuture.Result(MoveCollectionsClient) (OperationStatus, error)
1. *MoveCollectionsDeleteFuture.Result(MoveCollectionsClient) (OperationStatus, error)
1. *MoveCollectionsDiscardFuture.Result(MoveCollectionsClient) (OperationStatus, error)
1. *MoveCollectionsInitiateMoveFuture.Result(MoveCollectionsClient) (OperationStatus, error)
1. *MoveCollectionsPrepareFuture.Result(MoveCollectionsClient) (OperationStatus, error)
1. *MoveCollectionsResolveDependenciesFuture.Result(MoveCollectionsClient) (OperationStatus, error)
1. *MoveResourcesCreateFuture.Result(MoveResourcesClient) (MoveResource, error)
1. *MoveResourcesDeleteFuture.Result(MoveResourcesClient) (OperationStatus, error)

## Struct Changes

### Removed Structs

1. SummaryItem

### Removed Struct Fields

1. MoveCollectionsBulkRemoveFuture.azure.Future
1. MoveCollectionsCommitFuture.azure.Future
1. MoveCollectionsDeleteFuture.azure.Future
1. MoveCollectionsDiscardFuture.azure.Future
1. MoveCollectionsInitiateMoveFuture.azure.Future
1. MoveCollectionsPrepareFuture.azure.Future
1. MoveCollectionsResolveDependenciesFuture.azure.Future
1. MoveResourceCollection.Summary
1. MoveResourcesCreateFuture.azure.Future
1. MoveResourcesDeleteFuture.azure.Future

### New Funcs

1. MoveCollectionProperties.MarshalJSON() ([]byte, error)
1. MoveResourceCollection.MarshalJSON() ([]byte, error)

## Struct Changes

### New Structs

1. MoveCollectionPropertiesErrors
1. Summary
1. SummaryCollection

### New Struct Fields

1. MoveCollection.Etag
1. MoveCollectionProperties.Errors
1. MoveCollectionsBulkRemoveFuture.Result
1. MoveCollectionsBulkRemoveFuture.azure.FutureAPI
1. MoveCollectionsCommitFuture.Result
1. MoveCollectionsCommitFuture.azure.FutureAPI
1. MoveCollectionsDeleteFuture.Result
1. MoveCollectionsDeleteFuture.azure.FutureAPI
1. MoveCollectionsDiscardFuture.Result
1. MoveCollectionsDiscardFuture.azure.FutureAPI
1. MoveCollectionsInitiateMoveFuture.Result
1. MoveCollectionsInitiateMoveFuture.azure.FutureAPI
1. MoveCollectionsPrepareFuture.Result
1. MoveCollectionsPrepareFuture.azure.FutureAPI
1. MoveCollectionsResolveDependenciesFuture.Result
1. MoveCollectionsResolveDependenciesFuture.azure.FutureAPI
1. MoveResourceCollection.SummaryCollection
1. MoveResourceCollection.TotalCount
1. MoveResourceProperties.IsResolveRequired
1. MoveResourcesCreateFuture.Result
1. MoveResourcesCreateFuture.azure.FutureAPI
1. MoveResourcesDeleteFuture.Result
1. MoveResourcesDeleteFuture.azure.FutureAPI
