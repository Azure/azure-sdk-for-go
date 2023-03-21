## Go

``` yaml
title: MonitorQueryClient
description: Azure Monitor Query Go Client
generated-metadata: false

clear-output-folder: false
export-clients: true
go: true
input-file: 
    - https://github.com/Azure/azure-rest-api-specs/blob/dba6ed1f03bda88ac6884c0a883246446cc72495/specification/operationalinsights/data-plane/Microsoft.OperationalInsights/preview/2021-05-19_Preview/OperationalInsights.json
    - https://github.com/Azure/azure-rest-api-specs/blob/dba6ed1f03bda88ac6884c0a883246446cc72495/specification/monitor/resource-manager/Microsoft.Insights/stable/2018-01-01/metricDefinitions_API.json
    - https://github.com/Azure/azure-rest-api-specs/blob/dba6ed1f03bda88ac6884c0a883246446cc72495/specification/monitor/resource-manager/Microsoft.Insights/stable/2018-01-01/metrics_API.json
    - https://github.com/Azure/azure-rest-api-specs/blob/dba6ed1f03bda88ac6884c0a883246446cc72495/specification/monitor/resource-manager/Microsoft.Insights/preview/2017-12-01-preview/metricNamespaces_API.json
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery
openapi-type: "data-plane"
output-folder: ../azquery
override-client-name: LogsClient
security: "AADToken"
use: "@autorest/go@4.0.0-preview.44"
version: "^3.0.0"

directive:
  # delete metadata endpoint
  - from: swagger-document
    where: $["paths"]
    transform: >
        delete $["/workspaces/{workspaceId}/metadata"];

  # delete metadata operations
  - remove-operation: Metadata_Post
  - remove-operation: Metadata_Get
  - remove-operation: Query_Get

  # delete metadata models
  - remove-model: metadataResults
  - remove-model: metadataCategory
  - remove-model: metadataSolution
  - remove-model: metadataResourceType
  - remove-model: metadataTable
  - remove-model: metadataFunction
  - remove-model: metadataQuery
  - remove-model: metadataApplication
  - remove-model: metadataWorkspace
  - remove-model: metadataResource
  - remove-model: metadataPermissions

 # rename log queries
  - rename-operation:
      from: Query_Execute
      to: Logs_QueryWorkspace
  - rename-operation:
      from: Query_Batch
      to: Logs_QueryBatch

  # rename metric list to QueryResource
  - rename-operation:
      from: Metrics_List
      to: Metrics_QueryResource

 # rename ListMetricDefinitions and ListMetricNamespaces to generate in metrics_client.go
  - rename-operation:
      from: MetricDefinitions_List
      to: Metrics_ListDefinitions
  - rename-operation:
      from: MetricNamespaces_List
      to: Metrics_ListNamespaces

  # rename some metrics fields
  - from: swagger-document
    where: $.definitions.Metric.properties.timeseries
    transform: $["x-ms-client-name"] = "TimeSeries"
  - from: swagger-document
    where: $.definitions.TimeSeriesElement.properties.metadatavalues
    transform: $["x-ms-client-name"] = "MetadataValues"
  - from: swagger-document
    where: $.definitions.Response.properties.resourceregion
    transform: $["x-ms-client-name"] = "ResourceRegion"
  - from: swagger-document
    where: $.parameters.MetricNamespaceParameter
    transform: $["x-ms-client-name"] = "MetricNamespace"
  - from: swagger-document
    where: $.parameters.MetricNamesParameter
    transform: $["x-ms-client-name"] = "MetricNames"
  - from: swagger-document
    where: $.parameters.OrderByParameter
    transform: $["x-ms-client-name"] = "OrderBy"

  # rename Body.Workspaces to Body.AdditionalWorkspaces
  - from: swagger-document
    where: $.definitions.queryBody.properties.workspaces
    transform: $["x-ms-client-name"] = "AdditionalWorkspaces"
  
  # rename Render to Visualization
  - from: swagger-document
    where: $.definitions.queryResults.properties.render
    transform: $["x-ms-client-name"] = "Visualization"
  - from: swagger-document
    where: $.definitions.batchQueryResults.properties.render
    transform: $["x-ms-client-name"] = "Visualization"

  # rename BatchQueryRequest.ID to BatchQueryRequest.CorrelationID
  - from: swagger-document
    where: $.definitions.batchQueryRequest.properties.id
    transform: $["x-ms-client-name"] = "CorrelationID"
  - from: swagger-document
    where: $.definitions.batchQueryResponse.properties.id
    transform: $["x-ms-client-name"] = "CorrelationID"

  # rename BatchQueryRequest.Workspace to BatchQueryRequest.WorkspaceID
  - from: swagger-document
    where: $.definitions.batchQueryRequest.properties.workspace
    transform: $["x-ms-client-name"] = "WorkspaceID"
  
  # rename Prefer to Options
  - from: swagger-document
    where: $.parameters.PreferHeaderParameter
    transform: $["x-ms-client-name"] = "Options"
  - from: models.go
    where: $
    transform: return $.replace(/Options \*string/g, "Options *LogsQueryOptions");
  - from: logs_client.go
    where: $
    transform: return $.replace(/\*options\.Options/, "options.Options.preferHeader()");
  
  # add default values for batch request path and method attributes
  - from: swagger-document
    where: $.definitions.batchQueryRequest.properties.path
    transform: $["x-ms-client-default"] = "/query"
  - from: swagger-document
    where: $.definitions.batchQueryRequest.properties.method
    transform: $["x-ms-client-default"] = "POST"

  # add descriptions for models and constants that don't have them
  - from: swagger-document
    where: $.definitions.batchQueryRequest.properties.path
    transform: $["description"] = "The query path of a single request in a batch, defaults to /query"
  - from: swagger-document
    where: $.definitions.batchQueryRequest.properties.method
    transform: $["description"] = "The method of a single request in a batch, defaults to POST"
  - from: swagger-document
    where: $.definitions.batchQueryResponse
    transform: $["description"] = "Contains the batch query response and the headers, id, and status of the request"
  - from: constants.go
    where: $
    transform: return $.replace(/type ResultType string/, "//ResultType - Reduces the set of data collected. The syntax allowed depends on the operation. See the operation's description for details.\ntype ResultType string");

  # update doc comments
  - from: swagger-document
    where: $.paths["/workspaces/{workspaceId}/query"].post
    transform: $["description"] = "Executes an Analytics query for data."
  - from: swagger-document
    where: $.paths["/$batch"].post
    transform: $["description"] = "Executes a batch of Analytics queries for data."
  - from: swagger-document
    where: $.definitions.queryResults.properties.tables
    transform: $["description"] = "The results of the query in tabular format."
  - from: swagger-document
    where: $.definitions.batchQueryResults.properties.tables
    transform: $["description"] = "The results of the query in tabular format."
  - from: swagger-document
    where: $.definitions.queryBody.properties.workspaces
    transform: $["description"] = "A list of workspaces to query in addition to the primary workspace."
  - from: swagger-document
    where: $.definitions.batchQueryRequest.properties.headers
    transform: $["description"] = "Optional. Headers of the request. Can use prefer header to set server timeout, query statistics and visualization information. For more information, see https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery#readme-increase-wait-time-include-statistics-include-render-visualization"
  - from: swagger-document
    where: $.definitions.batchQueryRequest.properties.workspace
    transform: $["description"] = "Primary Workspace ID of the query"
  - from: swagger-document
    where: $.definitions.batchQueryRequest.properties.id
    transform: $["description"] = "Unique ID corresponding to each request in the batch"
  - from: swagger-document
    where: $.parameters.workspaceId
    transform: $["description"] = "Primary Workspace ID of the query. This is Workspace ID from the Properties blade in the Azure portal"

  # delete unused error models
  - from: models.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+type (?:ErrorResponse|ErrorResponseAutoGenerated).+\{(?:\s.+\s)+\}\s/g, "");
  - from: models_serde.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func \(\w \*?(?:ErrorResponse|ErrorResponseAutoGenerated)\).*\{\s(?:.+\s)+\}\s/g, "");
  - from: models.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+type (?:ErrorInfo|ErrorDetail).+\{(?:\s.+\s)+\}\s/g, "");
  - from: models_serde.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func \(\w \*?(?:ErrorInfo|ErrorDetail)\).*\{\s(?:.+\s)+\}\s/g, "");

  # delete generated constructor and client
  - from: logs_client.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func NewLogsClient.+\{\s(?:.+\s)+\}\s/, "");
  - from: logs_client.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+type LogsClient struct.+\{\s(?:.+\s)+\}\s/, "");
  - from: metrics_client.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+func NewMetricsClient.+\{\s(?:.+\s)+\}\s/, "");
  - from: metrics_client.go
    where: $
    transform: return $.replace(/(?:\/\/.*\s)+type MetricsClient.+\{\s(?:.+\s)+\}\s/, "");

  # point the clients to the correct host url
  - from: logs_client.go
    where: $
    transform: return $.replace(/host/g, "client.host");
  - from: metrics_client.go
    where: $
    transform: return $.replace(/host/g, "client.host");

  # delete generated host url
  - from: constants.go
    where: $
    transform: return $.replace(/const host = "(.*?)"/, "");

  # change Table.Rows from type [][]interface{} to type []Row
  - from: models.go
    where: $
    transform: return $.replace(/Rows \[\]\[\]interface{}/, "Rows []Row");

  # change render and statistics type to []byte
  - from: models.go
    where: $
    transform: return $.replace(/Statistics interface{}/g, "Statistics []byte");
  - from: models.go
    where: $
    transform: return $.replace(/Visualization interface{}/g, "Visualization []byte");
  - from: models_serde.go
    where: $
    transform: return 
      $.replace(/err(.*)r\.Statistics\)/, "r.Statistics = val") 
  - from: models_serde.go
    where: $
    transform: return $.replace(/err(.*)r\.Visualization\)/, "r.Visualization = val");
  - from: models_serde.go
    where: $
    transform: return 
      $.replace(/err(.*)b\.Statistics\)/, "b.Statistics = val") 
  - from: models_serde.go
    where: $
    transform: return $.replace(/err(.*)b\.Visualization\)/, "b.Visualization = val");

  # change type of timespan from *string to *TimeInterval
  - from: models.go
    where: $
    transform: return $.replace(/Timespan \*string/g, "Timespan *TimeInterval");
  - from: metrics_client.go
    where: $
    transform: return $.replace(/reqQP\.Set\(\"timespan\", \*options\.Timespan\)/g, "reqQP.Set(\"timespan\", string(*options.Timespan))");

  # change type of MetricsClientQueryResourceOptions.Aggregation from *string to []*AggregationType
  - from: models.go
    where: $
    transform: return $.replace(/Aggregation \*string/g, "Aggregation []*AggregationType");
  - from: metrics_client.go
    where: $
    transform: return $.replace(/\*options.Aggregation/g, "aggregationTypeToString(options.Aggregation)");
  - from: swagger-document
    where: $.parameters.AggregationsParameter
    transform: $["description"] = "The list of aggregation types to retrieve"