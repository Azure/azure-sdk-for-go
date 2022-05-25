# Unreleased

## Breaking Changes

### Signature Changes

#### Funcs

1. AlertsClient.ChangeState
	- Params
		- From: context.Context, string, AlertState
		- To: context.Context, string, AlertState, string
1. AlertsClient.ChangeStatePreparer
	- Params
		- From: context.Context, string, AlertState
		- To: context.Context, string, AlertState, string
1. AlertsClient.GetAll
	- Params
		- From: context.Context, string, string, string, MonitorService, MonitorCondition, Severity, AlertState, string, string, *bool, *bool, *int32, AlertsSortByFields, string, string, TimeRange, string
		- To: context.Context, string, string, string, MonitorService, MonitorCondition, Severity, AlertState, string, string, *bool, *bool, *int64, AlertsSortByFields, SortOrder, string, TimeRange, string
1. AlertsClient.GetAllComplete
	- Params
		- From: context.Context, string, string, string, MonitorService, MonitorCondition, Severity, AlertState, string, string, *bool, *bool, *int32, AlertsSortByFields, string, string, TimeRange, string
		- To: context.Context, string, string, string, MonitorService, MonitorCondition, Severity, AlertState, string, string, *bool, *bool, *int64, AlertsSortByFields, SortOrder, string, TimeRange, string
1. AlertsClient.GetAllPreparer
	- Params
		- From: context.Context, string, string, string, MonitorService, MonitorCondition, Severity, AlertState, string, string, *bool, *bool, *int32, AlertsSortByFields, string, string, TimeRange, string
		- To: context.Context, string, string, string, MonitorService, MonitorCondition, Severity, AlertState, string, string, *bool, *bool, *int64, AlertsSortByFields, SortOrder, string, TimeRange, string
1. SmartGroupsClient.GetAll
	- Params
		- From: context.Context, string, string, string, MonitorService, MonitorCondition, Severity, AlertState, TimeRange, *int32, SmartGroupsSortByFields, string
		- To: context.Context, string, string, string, MonitorService, MonitorCondition, Severity, AlertState, TimeRange, *int64, SmartGroupsSortByFields, SortOrder
1. SmartGroupsClient.GetAllComplete
	- Params
		- From: context.Context, string, string, string, MonitorService, MonitorCondition, Severity, AlertState, TimeRange, *int32, SmartGroupsSortByFields, string
		- To: context.Context, string, string, string, MonitorService, MonitorCondition, Severity, AlertState, TimeRange, *int64, SmartGroupsSortByFields, SortOrder
1. SmartGroupsClient.GetAllPreparer
	- Params
		- From: context.Context, string, string, string, MonitorService, MonitorCondition, Severity, AlertState, TimeRange, *int32, SmartGroupsSortByFields, string
		- To: context.Context, string, string, string, MonitorService, MonitorCondition, Severity, AlertState, TimeRange, *int64, SmartGroupsSortByFields, SortOrder

#### Struct Fields

1. AlertsSummaryGroup.SmartGroupsCount changed type from *int32 to *int64
1. AlertsSummaryGroup.Total changed type from *int32 to *int64
1. AlertsSummaryGroupItem.Count changed type from *int32 to *int64
1. SmartGroupAggregatedProperty.Count changed type from *int32 to *int64
1. SmartGroupProperties.AlertsCount changed type from *int32 to *int64

## Additive Changes

### New Constants

1. SortOrder.Asc
1. SortOrder.Desc

### New Funcs

1. PossibleSortOrderValues() []SortOrder

### Struct Changes

#### New Structs

1. ActionStatus

#### New Struct Fields

1. Essentials.ActionStatus
1. Essentials.Description
1. Operation.Origin
