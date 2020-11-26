
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewMoveCollectionResultListPage` signature has been changed from `(func(context.Context, MoveCollectionResultList) (MoveCollectionResultList, error))` to `(MoveCollectionResultList,func(context.Context, MoveCollectionResultList) (MoveCollectionResultList, error))`
- Function `NewMoveResourceCollectionPage` signature has been changed from `(func(context.Context, MoveResourceCollection) (MoveResourceCollection, error))` to `(MoveResourceCollection,func(context.Context, MoveResourceCollection) (MoveResourceCollection, error))`
- Function `MoveCollectionsClient.InitiateMove` signature has been changed from `(context.Context,string,string,*ResourceMoveRequest)` to `(context.Context,string,string,*ResourceMoveRequestType)`
- Function `MoveCollectionsClient.InitiateMovePreparer` signature has been changed from `(context.Context,string,string,*ResourceMoveRequest)` to `(context.Context,string,string,*ResourceMoveRequestType)`
- Type of `MoveResourceProperties.SourceResourceSettings` has been changed from `*MoveResourcePropertiesSourceResourceSettings` to `BasicResourceSettings`
- Type of `MoveResourceProperties.Errors` has been changed from `*MoveResourceError` to `*MoveResourcePropertiesErrors`
- Const `ResourceTypeMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsBasicResourceSettings` has been removed
- Function `ResourceGroupResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `SQLServerResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsSQLDatabaseResourceSettings` has been removed
- Function `VirtualNetworkResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsPublicIPAddressResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsAvailabilitySetResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsNetworkSecurityGroupResourceSettings` has been removed
- Function `LoadBalancerResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsResourceGroupResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsVirtualNetworkResourceSettings` has been removed
- Function `SQLDatabaseResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.MarshalJSON` has been removed
- Function `NetworkInterfaceResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsSQLElasticPoolResourceSettings` has been removed
- Function `NetworkSecurityGroupResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsNetworkInterfaceResourceSettings` has been removed
- Function `ResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `SQLElasticPoolResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `PublicIPAddressResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `AvailabilitySetResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsSQLServerResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsVirtualMachineResourceSettings` has been removed
- Function `MoveResourcePropertiesSourceResourceSettings.AsLoadBalancerResourceSettings` has been removed
- Function `VirtualMachineResourceSettings.AsMoveResourcePropertiesSourceResourceSettings` has been removed
- Struct `MoveResourcePropertiesSourceResourceSettings` has been removed
- Struct `ResourceMoveRequest` has been removed

## New Content

- Function `MoveResourcePropertiesMoveStatus.MarshalJSON() ([]byte,error)` is added
- Function `MoveCollectionsClient.BulkRemove(context.Context,string,string,*BulkRemoveRequest) (MoveCollectionsBulkRemoveFuture,error)` is added
- Function `MoveCollectionsClient.BulkRemovePreparer(context.Context,string,string,*BulkRemoveRequest) (*http.Request,error)` is added
- Function `MoveCollectionsClient.BulkRemoveResponder(*http.Response) (OperationStatus,error)` is added
- Function `MoveCollectionsClient.BulkRemoveSender(*http.Request) (MoveCollectionsBulkRemoveFuture,error)` is added
- Function `*MoveCollectionsBulkRemoveFuture.Result(MoveCollectionsClient) (OperationStatus,error)` is added
- Function `JobStatus.MarshalJSON() ([]byte,error)` is added
- Function `MoveResourceStatus.MarshalJSON() ([]byte,error)` is added
- Struct `BulkRemoveRequest` is added
- Struct `MoveCollectionsBulkRemoveFuture` is added
- Struct `MoveResourcePropertiesErrors` is added
- Struct `ResourceMoveRequestType` is added
- Struct `SummaryItem` is added
- Field `IsDataAction` is added to struct `OperationsDiscovery`
- Field `Summary` is added to struct `MoveResourceCollection`

