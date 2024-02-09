``` yaml
title: Metrics Query Client
clear-output-folder: false
go: true
input-file: https://github.com/Azure/azure-rest-api-specs/blob/0373f0edc4414fd402603fac51d0df93f1f70507/specification/monitor/data-plane/Microsoft.Insights/stable/2023-10-01/metricBatch.json
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/monitor/query/azmetrics
openapi-type: "data-plane"
output-folder: ../azmetrics
security: "AADToken"
use: "@autorest/go@4.0.0-preview.61"
inject-spans: true
version: "^3.0.0"

directive:
  # rename Batch to QueryResources
  - rename-operation:
        from: MetricsBatch_Batch
        to: Metrics_QueryResources

  # remove unused error models
  - from: swagger-document
    where: 
     - $.definitions..ErrorResponse
     - $.definitions..ErrorDetail
     - $.definitions..ErrorAdditionalInfo
    transform: $["x-ms-external"] = true

  # Rename MetricResultsResponse
  - rename-model:
      from: MetricResultsResponse
      to: MetricResults
  - from: 
        - models.go
        - models_serde.go
    where: $
    transform: return $.replace(/MetricResultsValuesItem/g, "MetricValues");
  - from: swagger-document
    where: $.definitions.MetricResults.properties.values.items
    transform: $["description"] = "Metric data values."

  # renaming or fixing the casing of struct fields and parameters
  - from: swagger-document
    where: $.definitions.Metric.properties.timeseries
    transform: $["x-ms-client-name"] = "TimeSeries"
  - from: swagger-document
    where: $.parameters.MetricNamespaceParameter
    transform: $["x-ms-client-name"] = "metricNamespace"
  - from: swagger-document
    where: $.parameters.MetricNamesParameter
    transform: $["x-ms-client-name"] = "metricNames"
  - from: swagger-document
    where: $.parameters.StartTimeParameter
    transform: $["x-ms-client-name"] = "StartTime"
  - from: swagger-document
    where: $.parameters.EndTimeParameter
    transform: $["x-ms-client-name"] = "EndTime"
  - from: swagger-document
    where: $.definitions.ResourceIdList.properties.resourceids
    transform: $["x-ms-client-name"] = "ResourceIDs"
  - from: swagger-document
    where: $.definitions.MetricResults.properties.values.items.properties.starttime
    transform: $["x-ms-client-name"] = "StartTime"
  - from: swagger-document
    where: $.definitions.MetricResults.properties.values.items.properties.endtime
    transform: $["x-ms-client-name"] = "EndTime"
  - from: swagger-document
    where: $.definitions.MetricResults.properties.values.items.properties.resourceid
    transform: $["x-ms-client-name"] = "ResourceID"
  - from: swagger-document
    where: $.definitions.MetricResults.properties.values.items.properties.resourceregion
    transform: $["x-ms-client-name"] = "ResourceRegion"
  - from: swagger-document
    where: $.definitions.MetricResults.properties.values.items.properties.value
    transform: $["x-ms-client-name"] = "Values"
  - from: swagger-document
    where: $.definitions.TimeSeriesElement.properties.metadatavalues
    transform: $["x-ms-client-name"] = "MetadataValues"
  - from: swagger-document
    where: $.parameters.OrderByParameter
    transform: $["x-ms-client-name"] = "OrderBy"
  - from: swagger-document
    where: $.parameters.RollUpByParameter
    transform: $["x-ms-client-name"] = "RollUpBy"
  - from: client.go
    where: $
    transform: return $.replace(/batchRequest/g, "resourceIDs");

  # delete client name prefix from method options and response types
  - from:
      - client.go
      - options.go
      - response_types.go
    where: $
    transform: return $.replace(/Client(\w+)((?:Options|Response))/g, "$1$2");
```