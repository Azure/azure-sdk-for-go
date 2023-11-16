//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package azquery

// LogsClientQueryBatchResponse contains the response from method LogsClient.QueryBatch.
type LogsClientQueryBatchResponse struct {
	// Response to a batch query.
	BatchResponse
}

// LogsClientQueryResourceResponse contains the response from method LogsClient.QueryResource.
type LogsClientQueryResourceResponse struct {
	// Contains the tables, columns & rows resulting from a query.
	Results
}

// LogsClientQueryWorkspaceResponse contains the response from method LogsClient.QueryWorkspace.
type LogsClientQueryWorkspaceResponse struct {
	// Contains the tables, columns & rows resulting from a query.
	Results
}

// MetricsBatchClientQueryBatchResponse contains the response from method MetricsBatchClient.QueryBatch.
type MetricsBatchClientQueryBatchResponse struct {
	// The metrics result for a resource.
	MetricResults
}

// MetricsClientListDefinitionsResponse contains the response from method MetricsClient.NewListDefinitionsPager.
type MetricsClientListDefinitionsResponse struct {
	// Represents collection of metric definitions.
	MetricDefinitionCollection
}

// MetricsClientListNamespacesResponse contains the response from method MetricsClient.NewListNamespacesPager.
type MetricsClientListNamespacesResponse struct {
	// Represents collection of metric namespaces.
	MetricNamespaceCollection
}

// MetricsClientQueryResourceResponse contains the response from method MetricsClient.QueryResource.
type MetricsClientQueryResourceResponse struct {
	// The response to a metrics query.
	Response
}
