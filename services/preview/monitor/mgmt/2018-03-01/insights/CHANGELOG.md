
Generated from https://github.com/Azure/azure-rest-api-specs/tree/b97299c968df5f99b724bd1231fd2161731d3b8f

Code generator C:\Users\dapzhang\Documents\workspace\autorest.go

## Breaking Changes

- Function `NewEventDataCollectionPage` signature has been changed from `(func(context.Context, EventDataCollection) (EventDataCollection, error))` to `(EventDataCollection,func(context.Context, EventDataCollection) (EventDataCollection, error))`
- Function `NewAutoscaleSettingResourceCollectionPage` signature has been changed from `(func(context.Context, AutoscaleSettingResourceCollection) (AutoscaleSettingResourceCollection, error))` to `(AutoscaleSettingResourceCollection,func(context.Context, AutoscaleSettingResourceCollection) (AutoscaleSettingResourceCollection, error))`

## New Content

- Field `SkipMetricValidation` is added to struct `MetricCriteria`
- Field `SkipMetricValidation` is added to struct `MultiMetricCriteria`
- Field `SkipMetricValidation` is added to struct `DynamicMetricCriteria`

