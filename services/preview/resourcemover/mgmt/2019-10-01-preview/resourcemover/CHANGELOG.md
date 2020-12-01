Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `MoveCollectionsClient.InitiateMovePreparer` parameter(s) have been changed from `(context.Context, string, string, *ResourceMoveRequest)` to `(context.Context, string, string, *ResourceMoveRequestType)`
- Function `NewMoveResourceCollectionPage` parameter(s) have been changed from `(func(context.Context, MoveResourceCollection) (MoveResourceCollection, error))` to `(MoveResourceCollection, func(context.Context, MoveResourceCollection) (MoveResourceCollection, error))`
- Function `NewMoveCollectionResultListPage` parameter(s) have been changed from `(func(context.Context, MoveCollectionResultList) (MoveCollectionResultList, error))` to `(MoveCollectionResultList, func(context.Context, MoveCollectionResultList) (MoveCollectionResultList, error))`
- Function `MoveCollectionsClient.InitiateMove` parameter(s) have been changed from `(context.Context, string, string, *ResourceMoveRequest)` to `(context.Context, string, string, *ResourceMoveRequestType)`
- Type of `MoveResourceProperties.Errors` has been changed from `*MoveResourceError` to `*MoveResourcePropertiesErrors`
- Type of `MoveResourceProperties.SourceResourceSettings` has been changed from `*MoveResourcePropertiesSourceResourceSettings` to `BasicResourceSettings`
- Const `ResourceTypeMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsNetworkInterfaceResourceSettings` has been removed
- Function `VirtualNetworkResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `SQLServerResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsResourceSettings` has been removed
- Function `VirtualMachineResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsLoadBalancerResourceSettings` has been removed
- Function `ResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsPublicIPAddressResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsSQLElasticPoolResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsSQLServerResourceSettings` has been removed
- Function `LoadBalancerResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsVirtualMachineResourceSettings` has been removed
- Function `NetworkSecurityGroupResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsResourceGroupResourceSettings` has been removed
- Function `SQLDatabaseResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsVirtualNetworkResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsAvailabilitySetResourceSettings` has been removed
- Function `ResourceGroupResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsBasicResourceSettings` has been removed
- Function `NetworkInterfaceResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsNetworkSecurityGroupResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsSQLDatabaseResourceSettings` has been removed
- Function `AvailabilitySetResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `PublicIPAddressResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.MarshalJSON` has been removed
- Function `SQLElasticPoolResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Struct `MoveResourcePropertiesSourceResourceSettings` has been removed
- Struct `ResourceMoveRequest` has been removed

## New Content

- New function `JobStatus.MarshalJSON() ([]byte, error)`
- New function `MoveResourcePropertiesMoveStatus.MarshalJSON() ([]byte, error)`
- New function `MoveCollectionsClient.BulkRemoveResponder(*http.Response) (OperationStatus, error)`
- New function `*MoveCollectionsBulkRemoveFuture.Result(MoveCollectionsClient) (OperationStatus, error)`
- New function `MoveCollectionsClient.BulkRemovePreparer(context.Context, string, string, *BulkRemoveRequest) (*http.Request, error)`
- New function `MoveResourceStatus.MarshalJSON() ([]byte, error)`
- New function `MoveCollectionsClient.BulkRemove(context.Context, string, string, *BulkRemoveRequest) (MoveCollectionsBulkRemoveFuture, error)`
- New function `MoveCollectionsClient.BulkRemoveSender(*http.Request) (MoveCollectionsBulkRemoveFuture, error)`
- New struct `BulkRemoveRequest`
- New struct `MoveCollectionsBulkRemoveFuture`
- New struct `MoveResourcePropertiesErrors`
- New struct `ResourceMoveRequestType`
- New struct `SummaryItem`
- New field `IsDataAction` in struct `OperationsDiscovery`
- New field `Summary` in struct `MoveResourceCollection`
