//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package azquery

import "time"

// BatchQueryRequest - An single request in a batch.
type BatchQueryRequest struct {
	// REQUIRED; The Analytics query. Learn more about the Analytics query syntax [https://azure.microsoft.com/documentation/articles/app-insights-analytics-reference/]
	Body *Body

	// REQUIRED; Unique ID corresponding to each request in the batch.
	CorrelationID *string

	// REQUIRED; Primary Workspace ID of the query. This is the Workspace ID from the Properties blade in the Azure portal.
	WorkspaceID *string

	// Headers of the request. Can use prefer header to set server timeout and to query statistics and visualization information.
	Headers map[string]*string

	// The method of a single request in a batch, defaults to POST
	Method *BatchQueryRequestMethod

	// The query path of a single request in a batch, defaults to /query
	Path *BatchQueryRequestPath
}

// BatchQueryResponse - Contains the batch query response and the headers, id, and status of the request
type BatchQueryResponse struct {
	// Contains the tables, columns & rows resulting from a query.
	Body          *BatchQueryResults
	CorrelationID *string

	// Dictionary of
	Headers map[string]*string
	Status  *int32
}

// BatchQueryResults - Contains the tables, columns & rows resulting from a query.
type BatchQueryResults struct {
	// The code and message for an error.
	Error *ErrorInfo

	// Statistics represented in JSON format.
	Statistics []byte

	// The results of the query in tabular format.
	Tables []*Table

	// Visualization data in JSON format.
	Visualization []byte
}

// BatchRequest - An array of requests.
type BatchRequest struct {
	// REQUIRED; An single request in a batch.
	Requests []*BatchQueryRequest
}

// BatchResponse - Response to a batch query.
type BatchResponse struct {
	// An array of responses corresponding to each individual request in a batch.
	Responses []*BatchQueryResponse
}

// Body - The Analytics query. Learn more about the Analytics query syntax [https://azure.microsoft.com/documentation/articles/app-insights-analytics-reference/]
type Body struct {
	// REQUIRED; The query to execute.
	Query *string

	// A list of workspaces to query in addition to the primary workspace.
	AdditionalWorkspaces []*string

	// Optional. The timespan over which to query data. This is an ISO8601 time period value. This timespan is applied in addition
	// to any that are specified in the query expression.
	Timespan *TimeInterval
}

// Column - A column in a table.
type Column struct {
	// The name of this column.
	Name *string

	// The data type of this column.
	Type *LogsColumnType
}

// LocalizableString - The localizable string class.
type LocalizableString struct {
	// REQUIRED; The invariant value.
	Value *string

	// The display name.
	LocalizedValue *string
}

// MetadataValue - Represents a metric metadata value.
type MetadataValue struct {
	// The name of the metadata.
	Name *LocalizableString

	// The value of the metadata.
	Value *string
}

// Metric - The result data of a query.
type Metric struct {
	// REQUIRED; the metric Id.
	ID *string

	// REQUIRED; the name and the display name of the metric, i.e. it is localizable string.
	Name *LocalizableString

	// REQUIRED; the time series returned when a data query is performed.
	TimeSeries []*TimeSeriesElement

	// REQUIRED; the resource type of the metric resource.
	Type *string

	// REQUIRED; The unit of the metric.
	Unit *MetricUnit

	// Detailed description of this metric.
	DisplayDescription *string

	// 'Success' or the error details on query failures for this metric.
	ErrorCode *string

	// Error message encountered querying this specific metric.
	ErrorMessage *string
}

// MetricAvailability - Metric availability specifies the time grain (aggregation interval or frequency) and the retention
// period for that time grain.
type MetricAvailability struct {
	// the retention period for the metric at the specified timegrain. Expressed as a duration 'PT1M', 'P1D', etc.
	Retention *string

	// the time grain specifies the aggregation interval for the metric. Expressed as a duration 'PT1M', 'P1D', etc.
	TimeGrain *string
}

// MetricDefinition - Metric definition class specifies the metadata for a metric.
type MetricDefinition struct {
	// Custom category name for this metric.
	Category *string

	// the name and the display name of the dimension, i.e. it is a localizable string.
	Dimensions []*LocalizableString

	// Detailed description of this metric.
	DisplayDescription *string

	// the resource identifier of the metric definition.
	ID *string

	// Flag to indicate whether the dimension is required.
	IsDimensionRequired *bool

	// the collection of what aggregation intervals are available to be queried.
	MetricAvailabilities []*MetricAvailability

	// The class of the metric.
	MetricClass *MetricClass

	// the name and the display name of the metric, i.e. it is a localizable string.
	Name *LocalizableString

	// the namespace the metric belongs to.
	Namespace *string

	// the primary aggregation type value defining how to use the values for display.
	PrimaryAggregationType *AggregationType

	// the resource identifier of the resource that emitted the metric.
	ResourceID *string

	// the collection of what aggregation types are supported.
	SupportedAggregationTypes []*AggregationType

	// The unit of the metric.
	Unit *MetricUnit
}

// MetricDefinitionCollection - Represents collection of metric definitions.
type MetricDefinitionCollection struct {
	// REQUIRED; the values for the metric definitions.
	Value []*MetricDefinition
}

// MetricNamespace - Metric namespace class specifies the metadata for a metric namespace.
type MetricNamespace struct {
	// Kind of namespace
	Classification *NamespaceClassification

	// The ID of the metric namespace.
	ID *string

	// The escaped name of the namespace.
	Name *string

	// Properties which include the fully qualified namespace name.
	Properties *MetricNamespaceName

	// The type of the namespace.
	Type *string
}

// MetricNamespaceCollection - Represents collection of metric namespaces.
type MetricNamespaceCollection struct {
	// REQUIRED; The values for the metric namespaces.
	Value []*MetricNamespace
}

// MetricNamespaceName - The fully qualified metric namespace name.
type MetricNamespaceName struct {
	// The metric namespace name.
	MetricNamespaceName *string
}

// MetricResultsResponse - The metrics result for a resource.
type MetricResultsResponse struct {
	// The collection of metric data responses per resource, per metric.
	Values []*MetricResultsResponseValuesItem
}

type MetricResultsResponseValuesItem struct {
	// REQUIRED; The end time, in datetime format, for which the data was retrieved.
	EndTime *string

	// REQUIRED; The start time, in datetime format, for which the data was retrieved.
	StartTime *string

	// REQUIRED; The value of the collection.
	Values []*Metric

	// The interval (window size) for which the metric data was returned in. Follows the IS8601/RFC3339 duration format (e.g.
	// 'P1D' for 1 day). This may be adjusted in the future and returned back from what
	// was originally requested. This is not present if a metadata request was made.
	Interval *string

	// The namespace of the metrics been queried
	Namespace *string

	// The resource that has been queried for metrics.
	ResourceID *string

	// The region of the resource been queried for metrics.
	ResourceRegion *string
}

// MetricValue - Represents a metric value.
type MetricValue struct {
	// REQUIRED; The timestamp for the metric value in ISO 8601 format.
	TimeStamp *time.Time

	// The average value in the time range.
	Average *float64

	// The number of samples in the time range. Can be used to determine the number of values that contributed to the average
	// value.
	Count *float64

	// The greatest value in the time range.
	Maximum *float64

	// The least value in the time range.
	Minimum *float64

	// The sum of all of the values in the time range.
	Total *float64
}

// ResourceIDList - The comma separated list of resource IDs to query metrics for.
type ResourceIDList struct {
	// The list of resource IDs to query metrics for.
	ResourceIDs []*string
}

// Response - The response to a metrics query.
type Response struct {
	// REQUIRED; The timespan for which the data was retrieved. Its value consists of two datetimes concatenated, separated by
	// '/'. This may be adjusted in the future and returned back from what was originally
	// requested.
	Timespan *TimeInterval

	// REQUIRED; the value of the collection.
	Value []*Metric

	// The integer value representing the relative cost of the query.
	Cost *int32

	// The interval (window size) for which the metric data was returned in ISO 8601 duration format with a special case for 'FULL'
	// value that returns single datapoint for entire time span requested (
	// Examples: PT15M, PT1H, P1D, FULL). This may be adjusted and different from what was originally requested if AutoAdjustTimegrain=true
	// is specified. This is not present if a metadata request was made.
	Interval *string

	// The namespace of the metrics being queried
	Namespace *string

	// The region of the resource being queried for metrics.
	ResourceRegion *string
}

// Results - Contains the tables, columns & rows resulting from a query.
type Results struct {
	// REQUIRED; The results of the query in tabular format.
	Tables []*Table

	// The code and message for an error.
	Error *ErrorInfo

	// Statistics represented in JSON format.
	Statistics []byte

	// Visualization data in JSON format.
	Visualization []byte
}

// Table - Contains the columns and rows for one table in a query response.
type Table struct {
	// REQUIRED; The list of columns in this table.
	Columns []*Column

	// REQUIRED; The name of the table.
	Name *string

	// REQUIRED; The resulting rows from this query.
	Rows []Row
}

// TimeSeriesElement - A time series result type. The discriminator value is always TimeSeries in this case.
type TimeSeriesElement struct {
	// An array of data points representing the metric values. This is only returned if a result type of data is specified.
	Data []*MetricValue

	// the metadata values returned if $filter was specified in the call.
	MetadataValues []*MetadataValue
}
