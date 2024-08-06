// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "fmt"

const (
	otelSpanNameCreateDatabase              = "create_database %s"
	otelSpanNameReadDatabase                = "read_database %s"
	otelSpanNameDeleteDatabase              = "delete_database %s"
	otelSpanNameQueryDatabases              = "query_databases %s"
	otelSpanNameReadThroughputDatabase      = "read_database_throughput %s"
	otelSpanNameReplaceThroughputDatabase   = "replace_database_throughput %s"
	otelSpanNameCreateContainer             = "create_container %s"
	otelSpanNameReadContainer               = "read_container %s"
	otelSpanNameDeleteContainer             = "delete_container %s"
	otelSpanNameReplaceContainer            = "replace_container %s"
	otelSpanNameQueryContainers             = "query_containers %s"
	otelSpanNameReadThroughputContainer     = "read_container_throughput %s"
	otelSpanNameReaplaceThroughputContainer = "replace_container_throughput %s"
	otelSpanNameExecuteBatch                = "execute_batch %s"
	otelSpanNameCreateItem                  = "create_item %s"
	otelSpanNameReadItem                    = "read_item %s"
	otelSpanNameDeleteItem                  = "delete_item %s"
	otelSpanNameReplaceItem                 = "replace_item %s"
	otelSpanNameUpsertItem                  = "upsert_item %s"
	otelSpanNamePatchItem                   = "patch_item %s"
	otelSpanNameQueryItems                  = "query_items %s"
)

func getSpanNameForDatabases(operationType operationType, resourceType resourceType, id string) (string, error) {
	switch resourceType {
	case resourceTypeDatabase:
		switch operationType {
		case operationTypeCreate:
			return fmt.Sprintf(otelSpanNameCreateDatabase, id), nil
		case operationTypeRead:
			return fmt.Sprintf(otelSpanNameReadDatabase, id), nil
		case operationTypeDelete:
			return fmt.Sprintf(otelSpanNameDeleteDatabase, id), nil
		case operationTypeQuery:
			return fmt.Sprintf(otelSpanNameQueryDatabases, id), nil
		}
	case resourceTypeOffer:
		switch operationType {
		case operationTypeRead:
			return fmt.Sprintf(otelSpanNameReadThroughputDatabase, id), nil
		case operationTypeReplace:
			return fmt.Sprintf(otelSpanNameReplaceThroughputDatabase, id), nil
		}
	}
	return "", fmt.Errorf("undefined telemetry span for operationType %v and resourceType %v", operationType, resourceType)
}

func getSpanNameForContainers(operationType operationType, resourceType resourceType, id string) (string, error) {
	switch resourceType {
	case resourceTypeCollection:
		switch operationType {
		case operationTypeCreate:
			return fmt.Sprintf(otelSpanNameCreateContainer, id), nil
		case operationTypeRead:
			return fmt.Sprintf(otelSpanNameReadContainer, id), nil
		case operationTypeDelete:
			return fmt.Sprintf(otelSpanNameDeleteContainer, id), nil
		case operationTypeReplace:
			return fmt.Sprintf(otelSpanNameReplaceContainer, id), nil
		case operationTypeQuery:
			return fmt.Sprintf(otelSpanNameQueryContainers, id), nil
		case operationTypeBatch:
			return fmt.Sprintf(otelSpanNameExecuteBatch, id), nil
		}
	case resourceTypeOffer:
		switch operationType {
		case operationTypeRead:
			return fmt.Sprintf(otelSpanNameReadThroughputContainer, id), nil
		case operationTypeReplace:
			return fmt.Sprintf(otelSpanNameReaplaceThroughputContainer, id), nil
		}
	}
	return "", fmt.Errorf("undefined telemetry span for operationType %v and resourceType %v", operationType, resourceType)
}

func getSpanNameForItems(operationType operationType, id string) (string, error) {
	switch operationType {
	case operationTypeCreate:
		return fmt.Sprintf(otelSpanNameCreateItem, id), nil
	case operationTypeRead:
		return fmt.Sprintf(otelSpanNameReadItem, id), nil
	case operationTypeDelete:
		return fmt.Sprintf(otelSpanNameDeleteItem, id), nil
	case operationTypeReplace:
		return fmt.Sprintf(otelSpanNameReplaceItem, id), nil
	case operationTypeUpsert:
		return fmt.Sprintf(otelSpanNameUpsertItem, id), nil
	case operationTypePatch:
		return fmt.Sprintf(otelSpanNamePatchItem, id), nil
	case operationTypeQuery:
		return fmt.Sprintf(otelSpanNameQueryItems, id), nil
	}
	return "", fmt.Errorf("undefined telemetry span for operationType %v and resourceType %v", operationType, resourceTypeDocument)
}
