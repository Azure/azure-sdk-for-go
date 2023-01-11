# Azure Monitor Query client module for Go

The Azure Monitor Query client library is used to execute read-only queries against [Azure Monitor][azure_monitor_overview]'s two data platforms:

- [Logs](https://docs.microsoft.com/azure/azure-monitor/logs/data-platform-logs) - Collects and organizes log and performance data from monitored resources. Data from different sources such as platform logs from Azure services, log and performance data from virtual machines agents, and usage and performance data from apps can be consolidated into a single [Azure Log Analytics workspace](https://docs.microsoft.com/azure/azure-monitor/logs/data-platform-logs#log-analytics-and-workspaces). The various data types can be analyzed together using the [Kusto Query Language][kusto_query_language]. See the [Kusto to SQL cheat sheet][kusto_to_sql] for more information.
- [Metrics](https://docs.microsoft.com/azure/azure-monitor/essentials/data-platform-metrics) - Collects numeric data from monitored resources into a time series database. Metrics are numerical values that are collected at regular intervals and describe some aspect of a system at a particular time. Metrics are lightweight and capable of supporting near real-time scenarios, making them particularly useful for alerting and fast detection of issues.

**NOTE**: This library is currently a beta. There may be breaking changes until it reaches semantic version `v1.0.0`.

## Getting started

### Install packages

Install `azquery` and `azidentity` with `go get`:
```Bash
go get github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
```
[azidentity][azure_identity] is used for Azure Active Directory authentication as demonstrated below.

### Prerequisites

* An [Azure subscription][azure_sub]
* A supported Go version (the Azure SDK supports the two most recent Go releases)
* For log queries, an [Azure Log Analytics workspace][log_analytics_workspace_create] ID. 
* For metric queries, the Resource URI of any Azure resource (Storage Account, Key Vault, CosmosDB, etc).

### Authentication

This document demonstrates using [azidentity.NewDefaultAzureCredential][default_cred_ref] to authenticate. The client accepts any [azidentity][azure_identity] credential. See the [azidentity][azure_identity] documentation for more information about other credential types.

The clients default to the Azure Public Cloud. See the [cloud][cloud_documentation] documentation for more information about other cloud configurations. 

#### Create a logs client

Example logs client: [link][example_logs_client]

#### Create a metrics client

Example metrics client: [link][example_metrics_client]

### Execute the query

For examples of Logs and Metrics queries, see the [Examples](#examples) section of this readme or in the example_test.go file of our GitHub repo for [azquery](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/monitor/azquery).

## Key concepts

### Logs query rate limits and throttling

The Log Analytics service applies throttling when the request rate is too high. Limits, such as the maximum number of rows returned, are also applied on the Kusto queries. For more information, see [Query API](https://docs.microsoft.com/azure/azure-monitor/service-limits#la-query-api).

If you're executing a batch logs query, a throttled request will return a `ErrorInfo` object. That object's `code` value will be `ThrottledError`.

### Metrics data structure

Each set of metric values is a time series with the following characteristics:

- The time the value was collected
- The resource associated with the value
- A namespace that acts like a category for the metric
- A metric name
- The value itself
- Some metrics may have multiple dimensions as described in multi-dimensional metrics. Custom metrics can have up to 10 dimensions.

### Timespan

It's best practice to always query with a timespan (type `ISO8601TimeInterval`) to prevent excessive queries of the entire logs or metrics data set. Logs uses the ISO8601 Time Interval Standard. All time should be represented in UTC. If the timespan  is included in both the kusto query string and `Timespan` field, the timespan will be the intersection of the two values.

Use the `NewISO8601TimeInterval()` method for easy creation.

Example timespan: [link][example_query_workspace]

## Examples

- [Logs query](#logs-query)
    - [Logs query body structure](#logs-query-body-structure)
	- [Logs query result structure](#logs-query-result-structure)
- [Batch logs query](#batch-query)
	- [Batch query request structure](#batch-query-request-structure)
	- [Batch query result structure](#batch-query-result-structure)
- [Advanced logs query](#advanced-logs-query)
	- [Query multiple workspaces](#query-multiple-workspaces)
	- [Increase wait time, include statistics, include render (visualization)](#increase-wait-time-include-statistics-include-render-visualization)
- [Metrics query](#metrics-query)
  - [Metrics result structure](#metrics-result-structure)
  - [List Metric Definitions](#list-metric-definitions)
  - [List Metric Namespaces](#list-metric-namespaces)

### Logs query
The example below shows a basic logs query using the `QueryWorkspace` method. `QueryWorkspace` takes in a [context][context], a [Log Analytics Workspace][log_analytics_workspace] ID string, a [Body](#logs-query-body-structure) struct, and a [LogsClientQueryWorkspaceOptions](#increase-wait-time-include-statistics-include-render-visualization) struct and returns a [Results](#logs-query-result-structure) struct.

A workspace ID is required to query logs. To find the workspace ID:

1. If not already made, [create a Log Analytics workspace][create_workspace].
1. Navigate to your workspace's page in the Azure portal.
2. From the **Overview** blade, copy the value of the `Workspace ID` property.

Example QueryWorkspace: [link][example_query_workspace]

#### Logs query body structure
```
Body
|---Query *string                  // Kusto Query
|---Timespan *ISO8601TimeInterval  // ISO8601 Standard Time Interval
|---AdditionalWorkspaces []*string // Optional- additional workspaces to query
```

#### Logs query result structure
```
Results
|---Tables []*Table
	|---Columns []*Column
		|---Name *string
		|---Type *LogsColumnType
	|---Name *string
	|---Rows []Row               // Rows contain the actual results of the query
|---Error *ErrorInfo
	|---Code *string
|---Visualization []byte
|---Statistics []byte
```

### Batch query
`QueryBatch` is an advanced method allowing users to execute multiple logs queries in a single request. The method accepts a [BatchRequest](#batch-query-request-structure) and returns a [BatchResponse](#batch-query-result-structure). `QueryBatch` can return results in any order (usually in order of completion/success). Use the `CorrelationID` field to identify the correct response. 

Example QueryBatch: [link][example_batch]

#### Batch query request structure

```
BatchRequest
|---Body *Body
	|---Query *string                 // Kusto Query
	|---Timespan *ISO8601TimeInterval // ISO8601 Standard Time Interval
	|---Workspaces []*string          // Optional- additional workspaces to query
|---CorrelationID *string             // unique identifier for each query in batch
|---WorkspaceID *string
|---Headers map[string]*string        // Optional- advanced query options in prefer header
|---Method *BatchQueryRequestMethod  // Optional- defaults to POST
|---Path *BatchQueryRequestPath      // Optional- defaults to /query
```

#### Batch query result structure

```
BatchResponse
|---Responses []*BatchQueryResponse
	|---Body *BatchQueryResults
		|---Error *ErrorInfo
			|---Code *string
		|---Visualization []byte
		|---Statistics []byte
		|---Tables []*Table
			|---Columns []*Column
				|---Name *string
				|---Type *LogsColumnType
			|---Name *string
			|---Rows []Row
	|---Headers map[string]*string
	|---CorrelationID *string
	|---Status *int32
```

### Advanced logs query

#### Query multiple workspaces

To run the same query against multiple Log Analytics workspaces, add the additional workspace ID strings to the `AdditionalWorkspaces` slice in the `Body` struct. 

When multiple workspaces are included in the query, the logs in the result table are not grouped according to the workspace from which it was retrieved.

Example additional workspaces: [link][example_queryworkspace_2]

#### Increase wait time, include statistics, include render (visualization)

By default, the Azure Monitor Query service will run your query for up to three minutes. To increase the default timeout, set `wait` to desired number of seconds in LogsClientQueryWorkspaceOptions Prefer string. Max wait time the service will allow is ten minutes (600 seconds).

To get logs query execution statistics, such as CPU and memory consumption, set `include-statistics` to true in LogsClientQueryWorkspaceOptions Prefer string.

To get visualization data for logs queries, set `include-render` to `true` in the `LogsClientQueryWorkspaceOptions` `Prefer` string.

```go
azquery.LogsClientQueryWorkspaceOptions{Prefer: to.Ptr("wait=600,include-statistics=true,include-render=true")}
```

Example QueryWorkspace options: [link][example_queryworkspace_2]

To do the same with `QueryBatch`, set the values in the `BatchQueryRequest.Headers` map with a key of "prefer".

### Metrics query

You can query metrics on an Azure resource using the `MetricsClient.QueryResource` method. For each requested metric, a set of aggregated values is returned inside the `Timeseries` collection.

A resource ID is required to query metrics. To find the resource ID:

1. Navigate to your resource's page in the Azure portal.
2. From the **Overview** blade, select the **JSON View** link.
3. In the resulting JSON, copy the value of the `id` property.

Example QueryResource example: [link][example_metrics_queryresource]

#### Metrics result structure
```
Response
|---Timespan *string
|---Value []*Metric
	|---ID *string
	|---Name *LocalizableString
	|---Timeseries []*TimeSeriesElement
		|---Data []*MetricValue
			|---TimeStamp *time.Time
			|---Average *float64
			|---Count *float64
			|---Maximum *float64
			|---Minimum *float64
			|---Total *float64
		|---Metadatavalues []*MetadataValue
			|---Name *LocalizableString
			|---Value *string 
	|---Type *string
	|---Unit *MetricUnit
	|---DisplayDescription *string
	|---ErrorCode *string
	|---ErrorMessage *string
|---Cost *int32
|---Interval *string
|---Namespace *string
|---Resourceregion *string
```

#### List Metric Definitions

To list the metric definitions for the resource, use the `NewListDefinitionsPager` method.

Example NewListDefinitionsPager: [link][example_metrics_listdefinitions]

#### List Metric Namespaces

To list the metric namespaces for the resource, use the `NewListNamespacesPager` method.

Example NewListNamespacesPager: [link][example_metrics_listnamespaces]

## Troubleshooting

See our [troubleshooting guide][troubleshooting_guide] for details on how to diagnose various failure scenarios.

## Next steps

To learn more about Azure Monitor, see the [Azure Monitor service documentation][azure_monitor_overview].

## Contributing

This project welcomes contributions and suggestions. Most contributions require you to agree to a [Contributor License Agreement (CLA)][cla] declaring that you have the right to, and actually do, grant us the rights to use your contribution.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide a CLA and decorate
the PR appropriately (e.g., label, comment). Simply follow the instructions provided by the bot. You will only need to
do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct][coc]. For more information, see
the [Code of Conduct FAQ][coc_faq] or contact [opencode@microsoft.com][coc_contact] with any additional questions or
comments.

<!-- LINKS -->
[managed_identity]: https://docs.microsoft.com/azure/active-directory/managed-identities-azure-resources/overview
[azquery]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/monitor/azquery
[azure_identity]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity
[azure_sub]: https://azure.microsoft.com/free/
[azure_monitor_create_using_portal]: https://docs.microsoft.com/azure/azure-monitor/logs/quick-create-workspace
[azure_monitor_overview]: https://docs.microsoft.com/azure/azure-monitor/overview
[context]: https://pkg.go.dev/context
[cloud_documentation]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud
[create_workspace]: https://learn.microsoft.com/azure/azure-monitor/logs/quick-create-workspace
[default_cred_ref]: https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity#defaultazurecredential
[example_batch]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery#example-LogsClient.QueryBatch
[example_query_workspace]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery#example-LogsClient.QueryWorkspace
[example_queryworkspace_2]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery#example-LogsClient.QueryWorkspace-Second
[example_logs_client]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery#NewLogsClient
[example_metrics_client]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery#NewMetricsClient
[example_metrics_listdefinitions]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery#example-MetricsClient.NewListDefinitionsPager
[example_metrics_listnamespaces]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery#example-MetricsClient.NewListNamespacesPager
[example_metrics_queryresource]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery#example-MetricsClient.QueryResource
[kusto_query_language]: https://learn.microsoft.com/azure/data-explorer/kusto/query/
[kusto_to_sql]: https://learn.microsoft.com/azure/data-explorer/kusto/query/sqlcheatsheet
[log_analytics_workspace]: https://learn.microsoft.com/azure/azure-monitor/logs/log-analytics-workspace-overview
[log_analytics_workspace_create]: https://learn.microsoft.com/azure/azure-monitor/logs/quick-create-workspace?tabs=azure-portal
[time_go]: https://pkg.go.dev/time
[time_intervals]: https://en.wikipedia.org/wiki/ISO_8601#Time_intervals
[troubleshooting_guide]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/monitor/azquery/TROUBLESHOOTING.md
[cla]: https://cla.microsoft.com
[coc]: https://opensource.microsoft.com/codeofconduct/
[coc_faq]: https://opensource.microsoft.com/codeofconduct/faq/
[coc_contact]: mailto:opencode@microsoft.com