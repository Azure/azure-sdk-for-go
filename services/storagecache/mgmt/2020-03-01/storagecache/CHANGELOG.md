
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewResourceSkusResultPage` signature has been changed from `(func(context.Context, ResourceSkusResult) (ResourceSkusResult, error))` to `(ResourceSkusResult,func(context.Context, ResourceSkusResult) (ResourceSkusResult, error))`
- Function `NewStorageTargetsResultPage` signature has been changed from `(func(context.Context, StorageTargetsResult) (StorageTargetsResult, error))` to `(StorageTargetsResult,func(context.Context, StorageTargetsResult) (StorageTargetsResult, error))`
- Function `NewUsageModelsResultPage` signature has been changed from `(func(context.Context, UsageModelsResult) (UsageModelsResult, error))` to `(UsageModelsResult,func(context.Context, UsageModelsResult) (UsageModelsResult, error))`
- Function `NewAPIOperationListResultPage` signature has been changed from `(func(context.Context, APIOperationListResult) (APIOperationListResult, error))` to `(APIOperationListResult,func(context.Context, APIOperationListResult) (APIOperationListResult, error))`
- Function `NewCachesListResultPage` signature has been changed from `(func(context.Context, CachesListResult) (CachesListResult, error))` to `(CachesListResult,func(context.Context, CachesListResult) (CachesListResult, error))`

## New Content

- Const `MetricAggregationTypeAverage` is added
- Const `MetricAggregationTypeMaximum` is added
- Const `MetricAggregationTypeNotSpecified` is added
- Const `MetricAggregationTypeMinimum` is added
- Const `Application` is added
- Const `MetricAggregationTypeCount` is added
- Const `MetricAggregationTypeNone` is added
- Const `MetricAggregationTypeTotal` is added
- Const `User` is added
- Const `ManagedIdentity` is added
- Const `Key` is added
- Function `*APIOperation.UnmarshalJSON([]byte) error` is added
- Function `PossibleCreatedByTypeValues() []CreatedByType` is added
- Function `PossibleMetricAggregationTypeValues() []MetricAggregationType` is added
- Function `APIOperation.MarshalJSON() ([]byte,error)` is added
- Struct `APIOperationProperties` is added
- Struct `APIOperationPropertiesServiceSpecification` is added
- Struct `MetricDimension` is added
- Struct `MetricSpecification` is added
- Struct `SystemData` is added
- Field `SystemData` is added to struct `Cache`
- Anonymous field `*APIOperationProperties` is added to struct `APIOperation`
- Field `Origin` is added to struct `APIOperation`
- Field `IsDataAction` is added to struct `APIOperation`
- Field `Location` is added to struct `StorageTargetResource`
- Field `SystemData` is added to struct `StorageTargetResource`
- Field `Description` is added to struct `APIOperationDisplay`
- Field `SystemData` is added to struct `StorageTarget`
- Field `Location` is added to struct `StorageTarget`

