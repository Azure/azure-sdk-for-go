# Unreleased

## Breaking Changes

### Removed Constants

1. PrimaryAggregationType.Average
1. PrimaryAggregationType.Maximum
1. PrimaryAggregationType.Minimum
1. PrimaryAggregationType.Total

### Removed Funcs

1. *ManagedDatabaseRestoreDetailsResult.UnmarshalJSON([]byte) error
1. ManagedDatabaseRestoreDetailsClient.Get(context.Context, string, string, string) (ManagedDatabaseRestoreDetailsResult, error)
1. ManagedDatabaseRestoreDetailsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. ManagedDatabaseRestoreDetailsClient.GetResponder(*http.Response) (ManagedDatabaseRestoreDetailsResult, error)
1. ManagedDatabaseRestoreDetailsClient.GetSender(*http.Request) (*http.Response, error)
1. ManagedDatabaseRestoreDetailsProperties.MarshalJSON() ([]byte, error)
1. ManagedDatabaseRestoreDetailsResult.MarshalJSON() ([]byte, error)
1. NewManagedDatabaseRestoreDetailsClient(string) ManagedDatabaseRestoreDetailsClient
1. NewManagedDatabaseRestoreDetailsClientWithBaseURI(string, string) ManagedDatabaseRestoreDetailsClient

### Struct Changes

#### Removed Structs

1. ManagedDatabaseRestoreDetailsClient
1. ManagedDatabaseRestoreDetailsProperties
1. ManagedDatabaseRestoreDetailsResult

#### Removed Struct Fields

1. DatabaseUpdate.*DatabaseProperties

### Signature Changes

#### Const Types

1. Count changed type from PrimaryAggregationType to QueryMetricUnitType
1. None changed type from PrimaryAggregationType to IdentityType

## Additive Changes

### New Constants

1. AggregationFunctionType.Avg
1. AggregationFunctionType.Max
1. AggregationFunctionType.Min
1. AggregationFunctionType.Stdev
1. AggregationFunctionType.Sum
1. IdentityType.SystemAssignedUserAssigned
1. IdentityType.UserAssigned
1. MetricType.CPU
1. MetricType.Dtu
1. MetricType.Duration
1. MetricType.Io
1. MetricType.LogIo
1. PrimaryAggregationType.PrimaryAggregationTypeAverage
1. PrimaryAggregationType.PrimaryAggregationTypeCount
1. PrimaryAggregationType.PrimaryAggregationTypeMaximum
1. PrimaryAggregationType.PrimaryAggregationTypeMinimum
1. PrimaryAggregationType.PrimaryAggregationTypeNone
1. PrimaryAggregationType.PrimaryAggregationTypeTotal
1. QueryMetricUnitType.KB
1. QueryMetricUnitType.Microseconds
1. QueryMetricUnitType.Percentage
1. QueryTimeGrainType.P1D
1. QueryTimeGrainType.PT1H

### New Funcs

1. *TopQueriesListResultIterator.Next() error
1. *TopQueriesListResultIterator.NextWithContext(context.Context) error
1. *TopQueriesListResultPage.Next() error
1. *TopQueriesListResultPage.NextWithContext(context.Context) error
1. DatabaseUpdateProperties.MarshalJSON() ([]byte, error)
1. ManagedInstancePecProperty.MarshalJSON() ([]byte, error)
1. ManagedInstancePrivateEndpointConnectionProperties.MarshalJSON() ([]byte, error)
1. ManagedInstancesClient.ListByManagedInstance(context.Context, string, string, *int32, string, string, string, QueryTimeGrainType, AggregationFunctionType, MetricType) (TopQueriesListResultPage, error)
1. ManagedInstancesClient.ListByManagedInstanceComplete(context.Context, string, string, *int32, string, string, string, QueryTimeGrainType, AggregationFunctionType, MetricType) (TopQueriesListResultIterator, error)
1. ManagedInstancesClient.ListByManagedInstancePreparer(context.Context, string, string, *int32, string, string, string, QueryTimeGrainType, AggregationFunctionType, MetricType) (*http.Request, error)
1. ManagedInstancesClient.ListByManagedInstanceResponder(*http.Response) (TopQueriesListResult, error)
1. ManagedInstancesClient.ListByManagedInstanceSender(*http.Request) (*http.Response, error)
1. NewTopQueriesListResultIterator(TopQueriesListResultPage) TopQueriesListResultIterator
1. NewTopQueriesListResultPage(TopQueriesListResult, func(context.Context, TopQueriesListResult) (TopQueriesListResult, error)) TopQueriesListResultPage
1. PossibleAggregationFunctionTypeValues() []AggregationFunctionType
1. PossibleMetricTypeValues() []MetricType
1. PossibleQueryMetricUnitTypeValues() []QueryMetricUnitType
1. PossibleQueryTimeGrainTypeValues() []QueryTimeGrainType
1. QueryMetricInterval.MarshalJSON() ([]byte, error)
1. QueryMetricProperties.MarshalJSON() ([]byte, error)
1. QueryStatisticsProperties.MarshalJSON() ([]byte, error)
1. TopQueries.MarshalJSON() ([]byte, error)
1. TopQueriesListResult.IsEmpty() bool
1. TopQueriesListResult.MarshalJSON() ([]byte, error)
1. TopQueriesListResultIterator.NotDone() bool
1. TopQueriesListResultIterator.Response() TopQueriesListResult
1. TopQueriesListResultIterator.Value() TopQueries
1. TopQueriesListResultPage.NotDone() bool
1. TopQueriesListResultPage.Response() TopQueriesListResult
1. TopQueriesListResultPage.Values() []TopQueries

### Struct Changes

#### New Structs

1. DatabaseUpdateProperties
1. ManagedInstancePecProperty
1. ManagedInstancePrivateEndpointConnectionProperties
1. ManagedInstancePrivateEndpointProperty
1. QueryMetricInterval
1. QueryMetricProperties
1. QueryStatisticsProperties
1. TopQueries
1. TopQueriesListResult
1. TopQueriesListResultIterator
1. TopQueriesListResultPage

#### New Struct Fields

1. DatabaseProperties.SourceResourceID
1. DatabaseUpdate.*DatabaseUpdateProperties
1. ManagedInstanceProperties.PrivateEndpointConnections
1. ManagedInstanceProperties.ZoneRedundant
1. ManagedInstanceUpdate.Identity
