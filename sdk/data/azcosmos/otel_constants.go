// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

const (
	otelSpanNameCreateDatabase              = "create_database"
	otelSpanNameReadDatabase                = "read_database"
	otelSpanNameDeleteDatabase              = "delete_database"
	otelSpanNameQueryDatabases              = "query_databases"
	otelSpanNameReadThroughputDatabase      = "read_database_throughput"
	otelSpanNameReplaceThroughputDatabase   = "replace_database_throughput"
	otelSpanNameCreateContainer             = "create_container"
	otelSpanNameReadContainer               = "read_container"
	otelSpanNameDeleteContainer             = "delete_container"
	otelSpanNameReplaceContainer            = "replace_container"
	otelSpanNameQueryContainers             = "query_containers"
	otelSpanNameReadThroughputContainer     = "read_container_throughput"
	otelSpanNameReaplaceThroughputContainer = "replace_container_throughput"
	otelSpanNameExecuteBatch                = "execute_batch"
	otelSpanNameCreateItem                  = "create_item"
	otelSpanNameReadItem                    = "read_item"
	otelSpanNameDeleteItem                  = "delete_item"
	otelSpanNameReplaceItem                 = "replace_item"
	otelSpanNameUpsertItem                  = "upsert_item"
	otelSpanNamePatchItem                   = "patch_item"
	otelSpanNameQueryItems                  = "query_items"
	otelSpanNamePartitionKeyRanges          = "read_partition_key_ranges"
)

type span struct {
	name    string
	options runtime.StartSpanOptions
}

func getSpanNameForClient(endpoint *url.URL, operationType operationType, resourceType resourceType, id string) (span, error) {
	var spanName string
	if resourceType == resourceTypeDatabase && operationType == operationTypeQuery {
		spanName = otelSpanNameQueryDatabases
	}
	if spanName == "" {
		return span{}, fmt.Errorf("undefined telemetry span for operationType %v and resourceType %v", operationType, resourceType)
	}

	return span{name: fmt.Sprintf("%s %s", spanName, id), options: getSpanPropertiesForClient(endpoint, spanName)}, nil
}

func getSpanNameForDatabases(endpoint *url.URL, operationType operationType, resourceType resourceType, id string) (span, error) {
	var spanName string
	switch resourceType {
	case resourceTypeDatabase:
		switch operationType {
		case operationTypeCreate:
			spanName = otelSpanNameCreateDatabase
		case operationTypeRead:
			spanName = otelSpanNameReadDatabase
		case operationTypeDelete:
			spanName = otelSpanNameDeleteDatabase
		}
	case resourceTypeCollection:
		if operationType == operationTypeQuery {
			spanName = otelSpanNameQueryContainers
		}
	case resourceTypeOffer:
		switch operationType {
		case operationTypeRead:
			spanName = otelSpanNameReadThroughputDatabase
		case operationTypeReplace:
			spanName = otelSpanNameReplaceThroughputDatabase
		}
	}

	if spanName == "" {
		return span{}, fmt.Errorf("undefined telemetry span for operationType %v and resourceType %v", operationType, resourceType)
	}

	return span{name: fmt.Sprintf("%s %s", spanName, id), options: getSpanPropertiesForDatabase(endpoint, spanName, id)}, nil
}

func getSpanNameForContainers(endpoint *url.URL, operationType operationType, resourceType resourceType, database string, id string) (span, error) {
	var spanName string
	switch resourceType {
	case resourceTypeCollection:
		switch operationType {
		case operationTypeCreate:
			spanName = otelSpanNameCreateContainer
		case operationTypeRead:
			spanName = otelSpanNameReadContainer
		case operationTypeDelete:
			spanName = otelSpanNameDeleteContainer
		case operationTypeReplace:
			spanName = otelSpanNameReplaceContainer
		case operationTypeBatch:
			spanName = otelSpanNameExecuteBatch
		}
	case resourceTypePartitionKeyRange:
		if operationType == operationTypeRead {
			spanName = otelSpanNamePartitionKeyRanges
		}
	case resourceTypeOffer:
		switch operationType {
		case operationTypeRead:
			spanName = otelSpanNameReadThroughputContainer
		case operationTypeReplace:
			spanName = otelSpanNameReaplaceThroughputContainer
		}
	}

	if spanName == "" {
		return span{}, fmt.Errorf("undefined telemetry span for operationType %v and resourceType %v", operationType, resourceType)
	}

	return span{name: fmt.Sprintf("%s %s", spanName, id), options: getSpanPropertiesForContainer(endpoint, spanName, database, id)}, nil
}

func getSpanNameForItems(endpoint *url.URL, operationType operationType, database string, id string) (span, error) {
	var spanName string
	switch operationType {
	case operationTypeCreate:
		spanName = otelSpanNameCreateItem
	case operationTypeRead:
		spanName = otelSpanNameReadItem
	case operationTypeDelete:
		spanName = otelSpanNameDeleteItem
	case operationTypeReplace:
		spanName = otelSpanNameReplaceItem
	case operationTypeUpsert:
		spanName = otelSpanNameUpsertItem
	case operationTypePatch:
		spanName = otelSpanNamePatchItem
	case operationTypeQuery:
		spanName = otelSpanNameQueryItems
	}

	if spanName == "" {
		return span{}, fmt.Errorf("undefined telemetry span for operationType %v and resourceType %v", operationType, resourceTypeDocument)
	}

	return span{name: fmt.Sprintf("%s %s", spanName, id), options: getSpanPropertiesForContainer(endpoint, spanName, database, id)}, nil
}

func getSpanPropertiesForClient(endpoint *url.URL, operationName string) runtime.StartSpanOptions {
	options := runtime.StartSpanOptions{
		Kind: tracing.SpanKindClient,
		Attributes: []tracing.Attribute{
			{Key: "db.system", Value: "cosmosdb"},
			{Key: "db.cosmosdb.connection_mode", Value: "gateway"},
			{Key: "db.operation.name", Value: operationName},
			{Key: "server.address", Value: endpoint.Hostname()},
		},
	}

	if endpoint.Port() != "443" {
		options.Attributes = append(options.Attributes, tracing.Attribute{Key: "server.port", Value: endpoint.Port()})
	}

	return options
}

func getSpanPropertiesForDatabase(endpoint *url.URL, operationName string, id string) runtime.StartSpanOptions {
	options := runtime.StartSpanOptions{
		Kind: tracing.SpanKindClient,
		Attributes: []tracing.Attribute{
			{Key: "db.system", Value: "cosmosdb"},
			{Key: "db.cosmosdb.connection_mode", Value: "gateway"},
			{Key: "db.namespace", Value: id},
			{Key: "db.operation.name", Value: operationName},
			{Key: "server.address", Value: endpoint.Hostname()},
		},
	}

	if endpoint.Port() != "443" {
		options.Attributes = append(options.Attributes, tracing.Attribute{Key: "server.port", Value: endpoint.Port()})
	}

	return options
}

func getSpanPropertiesForContainer(endpoint *url.URL, operationName string, database string, id string) runtime.StartSpanOptions {
	options := runtime.StartSpanOptions{
		Kind: tracing.SpanKindClient,
		Attributes: []tracing.Attribute{
			{Key: "db.system", Value: "cosmosdb"},
			{Key: "db.cosmosdb.connection_mode", Value: "gateway"},
			{Key: "db.namespace", Value: database},
			{Key: "db.collection.name", Value: id},
			{Key: "db.operation.name", Value: operationName},
			{Key: "server.address", Value: endpoint.Hostname()},
		},
	}

	if endpoint.Port() != "443" {
		options.Attributes = append(options.Attributes, tracing.Attribute{Key: "server.port", Value: endpoint.Port()})
	}

	return options
}
