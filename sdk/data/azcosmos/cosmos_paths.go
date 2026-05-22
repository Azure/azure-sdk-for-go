// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	pathSegmentDatabase            string = "dbs"
	pathSegmentCollection          string = "colls"
	pathSegmentUser                string = "users"
	pathSegmentPermission          string = "permissions"
	pathSegmentStoredProcedure     string = "sprocs"
	pathSegmentTrigger             string = "triggers"
	pathSegmentUserDefinedFunction string = "udfs"
	pathSegmentConflict            string = "conflicts"
	pathSegmentDocument            string = "docs"
	pathSegmentClientEncryptionKey string = "clientencryptionkeys"
	pathSegmentOffer               string = "offers"
	pathSegmentDatabaseAccount     string = ""
	pathSegmentPartitionKeyRange   string = "pkranges"
)

// generatePathForNameBased generates the URL path for a request based on the current operation
func generatePathForNameBased(resourceType resourceType, ownerOrResourceId string, isFeed bool) (string, error) {
	if isFeed && ownerOrResourceId == "" &&
		resourceType != resourceTypeDatabase &&
		resourceType != resourceTypeOffer &&
		resourceType != resourceTypeDatabaseAccount {
		return "", errors.New("INVALID resource type")
	}

	if resourceType == resourceTypeDatabase {
		if isFeed {
			return pathSegmentDatabase, nil
		}
		return ownerOrResourceId, nil
	}

	if resourceType == resourceTypeCollection {
		if isFeed {
			return ownerOrResourceId + "/" + pathSegmentCollection, nil
		}
		return ownerOrResourceId, nil
	}

	if resourceType == resourceTypeOffer {
		if isFeed {
			return pathSegmentOffer, nil
		}
		return ownerOrResourceId, nil
	}

	if resourceType == resourceTypeStoredProcedure {
		if isFeed {
			return ownerOrResourceId + "/" + pathSegmentStoredProcedure, nil
		}
		return ownerOrResourceId, nil
	}

	if resourceType == resourceTypeUserDefinedFunction {
		if isFeed {
			return ownerOrResourceId + "/" + pathSegmentUserDefinedFunction, nil
		}
		return ownerOrResourceId, nil
	}

	if resourceType == resourceTypeTrigger {
		if isFeed {
			return ownerOrResourceId + "/" + pathSegmentTrigger, nil
		}
		return ownerOrResourceId, nil
	}

	if resourceType == resourceTypeConflict {
		if isFeed {
			return ownerOrResourceId + "/" + pathSegmentConflict, nil
		}
		return ownerOrResourceId, nil
	}

	if resourceType == resourceTypePartitionKeyRange {
		if isFeed {
			return ownerOrResourceId + "/" + pathSegmentPartitionKeyRange, nil
		}
		return ownerOrResourceId, nil
	}

	if resourceType == resourceTypeUser {
		if isFeed {
			return ownerOrResourceId + "/" + pathSegmentUser, nil
		}
		return ownerOrResourceId, nil
	}

	if resourceType == resourceTypePermission {
		if isFeed {
			return ownerOrResourceId + "/" + pathSegmentPermission, nil
		}
		return ownerOrResourceId, nil
	}

	if resourceType == resourceTypeDocument {
		if isFeed {
			return ownerOrResourceId + "/" + pathSegmentDocument, nil
		}
		return ownerOrResourceId, nil
	}

	if resourceType == resourceTypeDatabaseAccount {
		return pathSegmentDatabaseAccount + "/" + ownerOrResourceId, nil
	}

	if resourceType == resourceTypeClientEncryptionKey {
		return ownerOrResourceId, nil
	}

	return "", fmt.Errorf("INVALID resource type %v, isFeed %v, resourceId %v", resourceType, isFeed, ownerOrResourceId)
}

// getResourcePath is used in auth flows.
func getResourcePath(resourceType resourceType) (string, error) {
	switch resourceType {
	case resourceTypeDatabase:
		return pathSegmentDatabase, nil
	case resourceTypeCollection:
		return pathSegmentCollection, nil
	case resourceTypeDocument:
		return pathSegmentDocument, nil
	case resourceTypeDatabaseAccount:
		return pathSegmentDatabaseAccount, nil
	case resourceTypeOffer:
		return pathSegmentOffer, nil
	case resourceTypeUser:
		return pathSegmentUser, nil
	case resourceTypeStoredProcedure:
		return pathSegmentStoredProcedure, nil
	case resourceTypeUserDefinedFunction:
		return pathSegmentUserDefinedFunction, nil
	case resourceTypeTrigger:
		return pathSegmentTrigger, nil
	case resourceTypePermission:
		return pathSegmentPermission, nil
	case resourceTypeConflict:
		return pathSegmentConflict, nil
	case resourceTypePartitionKeyRange:
		return pathSegmentPartitionKeyRange, nil
	case resourceTypeClientEncryptionKey:
		return pathSegmentClientEncryptionKey, nil
	default:
		return "", fmt.Errorf("%v is not a valid resource type", resourceType)
	}
}

// createLink generates a url link for a resource base on the parent paths
func createLink(parentPath string, pathSegment string, id string) string {
	var completePath strings.Builder
	parentPathLength := len(parentPath)
	completePath.Grow(parentPathLength + 2 + len(pathSegment) + len(id))
	if parentPathLength > 0 {
		fmt.Fprint(&completePath, parentPath)
		fmt.Fprint(&completePath, "/")
	}
	fmt.Fprint(&completePath, pathSegment)
	fmt.Fprint(&completePath, "/")
	fmt.Fprint(&completePath, url.PathEscape(id))
	return completePath.String()
}
