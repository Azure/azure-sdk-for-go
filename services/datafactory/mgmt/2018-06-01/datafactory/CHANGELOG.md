
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewDatasetListResponsePage` signature has been changed from `(func(context.Context, DatasetListResponse) (DatasetListResponse, error))` to `(DatasetListResponse,func(context.Context, DatasetListResponse) (DatasetListResponse, error))`
- Function `NewOperationListResponsePage` signature has been changed from `(func(context.Context, OperationListResponse) (OperationListResponse, error))` to `(OperationListResponse,func(context.Context, OperationListResponse) (OperationListResponse, error))`
- Function `NewLinkedServiceListResponsePage` signature has been changed from `(func(context.Context, LinkedServiceListResponse) (LinkedServiceListResponse, error))` to `(LinkedServiceListResponse,func(context.Context, LinkedServiceListResponse) (LinkedServiceListResponse, error))`
- Function `NewManagedVirtualNetworkListResponsePage` signature has been changed from `(func(context.Context, ManagedVirtualNetworkListResponse) (ManagedVirtualNetworkListResponse, error))` to `(ManagedVirtualNetworkListResponse,func(context.Context, ManagedVirtualNetworkListResponse) (ManagedVirtualNetworkListResponse, error))`
- Function `NewPipelineListResponsePage` signature has been changed from `(func(context.Context, PipelineListResponse) (PipelineListResponse, error))` to `(PipelineListResponse,func(context.Context, PipelineListResponse) (PipelineListResponse, error))`
- Function `NewFactoryListResponsePage` signature has been changed from `(func(context.Context, FactoryListResponse) (FactoryListResponse, error))` to `(FactoryListResponse,func(context.Context, FactoryListResponse) (FactoryListResponse, error))`
- Function `NewQueryDataFlowDebugSessionsResponsePage` signature has been changed from `(func(context.Context, QueryDataFlowDebugSessionsResponse) (QueryDataFlowDebugSessionsResponse, error))` to `(QueryDataFlowDebugSessionsResponse,func(context.Context, QueryDataFlowDebugSessionsResponse) (QueryDataFlowDebugSessionsResponse, error))`
- Function `NewManagedPrivateEndpointListResponsePage` signature has been changed from `(func(context.Context, ManagedPrivateEndpointListResponse) (ManagedPrivateEndpointListResponse, error))` to `(ManagedPrivateEndpointListResponse,func(context.Context, ManagedPrivateEndpointListResponse) (ManagedPrivateEndpointListResponse, error))`
- Function `NewTriggerListResponsePage` signature has been changed from `(func(context.Context, TriggerListResponse) (TriggerListResponse, error))` to `(TriggerListResponse,func(context.Context, TriggerListResponse) (TriggerListResponse, error))`
- Function `NewDataFlowListResponsePage` signature has been changed from `(func(context.Context, DataFlowListResponse) (DataFlowListResponse, error))` to `(DataFlowListResponse,func(context.Context, DataFlowListResponse) (DataFlowListResponse, error))`
- Function `NewIntegrationRuntimeListResponsePage` signature has been changed from `(func(context.Context, IntegrationRuntimeListResponse) (IntegrationRuntimeListResponse, error))` to `(IntegrationRuntimeListResponse,func(context.Context, IntegrationRuntimeListResponse) (IntegrationRuntimeListResponse, error))`
- Type of `SapHanaSource.PartitionOption` has been changed from `SapHanaPartitionOption` to `interface{}`
- Type of `TeradataSource.PartitionOption` has been changed from `TeradataPartitionOption` to `interface{}`
- Type of `NetezzaSource.PartitionOption` has been changed from `NetezzaPartitionOption` to `interface{}`
- Type of `OracleSource.PartitionOption` has been changed from `OraclePartitionOption` to `interface{}`
- Type of `SapTableSource.PartitionOption` has been changed from `SapTablePartitionOption` to `interface{}`
- Type of `SQLServerSource.PartitionOption` has been changed from `SQLPartitionOption` to `interface{}`
- Type of `SQLSource.PartitionOption` has been changed from `SQLPartitionOption` to `interface{}`
- Type of `SQLMISource.PartitionOption` has been changed from `SQLPartitionOption` to `interface{}`
- Type of `SQLDWSource.PartitionOption` has been changed from `SQLPartitionOption` to `interface{}`
- Type of `AzureSQLSource.PartitionOption` has been changed from `SQLPartitionOption` to `interface{}`
- Type of `ExecuteDataFlowActivityTypePropertiesCompute.ComputeType` has been changed from `DataFlowComputeType` to `interface{}`
- Type of `ExecuteDataFlowActivityTypePropertiesCompute.CoreCount` has been changed from `*int32` to `interface{}`

## New Content

- Const `TumblingWindowFrequencyMonth` is added
- Struct `CopyActivityLogSettings` is added
- Struct `LogLocationSettings` is added
- Struct `LogSettings` is added
- Field `LogSettings` is added to struct `CopyActivityTypeProperties`
- Field `SessionToken` is added to struct `AmazonS3LinkedServiceTypeProperties`
- Field `AuthenticationType` is added to struct `AmazonS3LinkedServiceTypeProperties`
- Field `TraceLevel` is added to struct `ExecuteDataFlowActivityTypeProperties`
- Field `ContinueOnError` is added to struct `ExecuteDataFlowActivityTypeProperties`
- Field `RunConcurrently` is added to struct `ExecuteDataFlowActivityTypeProperties`
- Field `ConnectionProperties` is added to struct `ConcurLinkedServiceTypeProperties`

