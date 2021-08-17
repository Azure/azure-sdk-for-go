// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
)

const (
	pathSegmentDatabase   string = "dbs"
	pathSegmentCollection string = "colls"
)

func getResourcePath(resourceType resourceType) (string, error) {
	switch resourceType {
	case resourceTypeDatabase:
		return pathSegmentDatabase, nil
	case resourceTypeCollection:
		return pathSegmentCollection, nil
	default:
		return "", fmt.Errorf("%v is not a valid resource type", resourceType)
	}
}
